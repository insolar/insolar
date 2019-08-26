//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package goplugin - golang plugin in docker runner
package goplugin

import (
	"context"
	"net/rpc"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
)

// Options of the GoPlugin
type Options struct {
	// Listen  is address `GoPlugin` listens on and provides RPC interface for runner(s)
	Listen string
}

// RunnerOptions - set of options to control internal isolated code runner(s)
type RunnerOptions struct {
	// Listen is address the runner listens on and provides RPC interface for the `GoPlugin`
	Listen string
	// CodeStoragePath is path to directory where the runner caches code
	CodeStoragePath string
}

// GoPlugin is a logic runner of code written in golang and compiled as go plugins
type GoPlugin struct {
	Cfg             *configuration.LogicRunner
	ArtifactManager artifacts.Client

	clientMutex sync.Mutex
	client      *rpc.Client
}

// NewGoPlugin returns a new started GoPlugin
func NewGoPlugin(conf *configuration.LogicRunner, am artifacts.Client) (*GoPlugin, error) {
	gp := GoPlugin{
		Cfg:             conf,
		ArtifactManager: am,
	}

	return &gp, nil
}

const timeout = time.Minute * 10

// Downstream returns a connection to `ginsider`
func (gp *GoPlugin) Downstream(ctx context.Context) (*rpc.Client, error) {
	_, span := instracer.StartSpan(ctx, "GoPlugin.Downstream")
	defer span.End()

	gp.clientMutex.Lock()
	defer gp.clientMutex.Unlock()

	if gp.client != nil {
		return gp.client, nil
	}

	inslogger.FromContext(ctx).Debug("dialing insgorund")
	client, err := rpc.Dial(gp.Cfg.GoPlugin.RunnerProtocol, gp.Cfg.GoPlugin.RunnerListen)
	if err != nil {
		return nil, errors.Wrapf(
			err, "couldn't dial '%s' over %s",
			gp.Cfg.GoPlugin.RunnerListen, gp.Cfg.GoPlugin.RunnerProtocol,
		)
	}

	gp.client = client
	return gp.client, nil
}

func (gp *GoPlugin) CloseDownstream() {
	gp.clientMutex.Lock()
	defer gp.clientMutex.Unlock()

	// this method can be called multiple times from callClientWithReconnect
	if gp.client != nil {
		gp.client.Close()
		gp.client = nil
	}
}

func (gp *GoPlugin) callClientWithReconnect(ctx context.Context, method string, req interface{}, res interface{}) error {
	var err error
	var client *rpc.Client

	for {
		client, err = gp.Downstream(ctx)
		if err == nil {
			ctx, span := instracer.StartSpan(ctx, "GoPlugin callClientWithReconnect")
			defer span.End()

			inslogger.FromContext(ctx).Debug("Sending request to insgorund")

			call := <-client.Go(method, req, res, nil).Done
			err = call.Error

			inslogger.FromContext(ctx).Debug("insgorund replied")

			if err != rpc.ErrShutdown {
				break
			} else {
				inslogger.FromContext(ctx).Debug("Connection to insgorund is closed, need to reconnect")
				gp.CloseDownstream()
				inslogger.FromContext(ctx).Debugf("Reconnecting...")
			}
		} else {
			inslogger.FromContext(ctx).Debugf("Can't connect to to insgorund, err: %s", err.Error())
			inslogger.FromContext(ctx).Debugf("Reconnecting...")
		}
	}

	return err
}

type CallMethodResult struct {
	Response rpctypes.DownCallMethodResp
	Error    error
}

func (gp *GoPlugin) CallMethodRPC(ctx context.Context, req rpctypes.DownCallMethodReq, res rpctypes.DownCallMethodResp, resultChan chan CallMethodResult) {
	inslogger.FromContext(ctx).Debug("GoPlugin.CallMethodRPC starts ...")
	method := "RPC.CallMethod"
	callClientError := gp.callClientWithReconnect(ctx, method, req, &res)
	resultChan <- CallMethodResult{Response: res, Error: callClientError}
}

// CallMethod runs a method on an object in controlled environment
func (gp *GoPlugin) CallMethod(
	ctx context.Context, callContext *insolar.LogicCallContext,
	code insolar.Reference, data []byte,
	method string, args insolar.Arguments,
) (
	[]byte, insolar.Arguments, error,
) {
	ctx = insmetrics.InsertTag(ctx, tagMethodName, method)

	ctx, span := instracer.StartSpan(ctx, "GoPlugin.CallMethod "+method)
	defer span.End()

	inslogger.FromContext(ctx).Debug("GoPlugin.CallMethod starts")
	start := time.Now()
	defer func() {
		stats.Record(ctx, statGopluginContractMethodTime.M(
			float64(time.Since(start).Nanoseconds())/1e6,
		))
	}()

	res := rpctypes.DownCallMethodResp{}
	req := rpctypes.DownCallMethodReq{
		Context:   callContext,
		Code:      code,
		Data:      data,
		Method:    method,
		Arguments: args,
	}

	resultChan := make(chan CallMethodResult)
	go gp.CallMethodRPC(ctx, req, res, resultChan)

	select {
	case callResult := <-resultChan:
		if callResult.Error != nil {
			return nil, nil, errors.Wrap(callResult.Error, "problem with API call")
		}
		return callResult.Response.Data, callResult.Response.Ret, nil
	case <-time.After(timeout):
		inslogger.FromContext(ctx).Debug("CallMethodRPC waiting results timeout")
		return nil, nil, errors.New("logicrunner execution timeout")
	}
}

type CallConstructorResult struct {
	Response rpctypes.DownCallConstructorResp
	Error    error
}

func (gp *GoPlugin) CallConstructorRPC(ctx context.Context, req rpctypes.DownCallConstructorReq, res rpctypes.DownCallConstructorResp, resultChan chan CallConstructorResult) {
	method := "RPC.CallConstructor"
	callClientError := gp.callClientWithReconnect(ctx, method, req, &res)
	resultChan <- CallConstructorResult{Response: res, Error: callClientError}
}

// CallConstructor runs a constructor of a contract in controlled environment
func (gp *GoPlugin) CallConstructor(
	ctx context.Context, callContext *insolar.LogicCallContext,
	code insolar.Reference, name string, args insolar.Arguments,
) (
	[]byte, insolar.Arguments, error,
) {

	res := rpctypes.DownCallConstructorResp{}
	req := rpctypes.DownCallConstructorReq{
		Context:   callContext,
		Code:      code,
		Name:      name,
		Arguments: args,
	}

	resultChan := make(chan CallConstructorResult)
	go gp.CallConstructorRPC(ctx, req, res, resultChan)

	select {
	case callResult := <-resultChan:
		if callResult.Error != nil {
			return nil, nil, errors.Wrap(callResult.Error, "problem with API call")
		}
		return callResult.Response.Data, callResult.Response.Ret, nil
	case <-time.After(timeout):
		inslogger.FromContext(ctx).Debug("CallConstructor waiting results timeout")
		return nil, nil, errors.New("logicrunner execution timeout")
	}
}
