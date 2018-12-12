/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

// Package goplugin - golang plugin in docker runner
package goplugin

import (
	"context"
	"net/rpc"
	"os/exec"
	"sync"
	"time"

	"github.com/insolar/insolar/metrics"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
	"github.com/pkg/errors"
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
	MessageBus      core.MessageBus
	ArtifactManager core.ArtifactManager
	runner          *exec.Cmd

	clientMutex sync.Mutex
	client      *rpc.Client
}

// NewGoPlugin returns a new started GoPlugin
func NewGoPlugin(conf *configuration.LogicRunner, eb core.MessageBus, am core.ArtifactManager) (*GoPlugin, error) {
	gp := GoPlugin{
		Cfg:             conf,
		MessageBus:      eb,
		ArtifactManager: am,
	}

	return &gp, nil
}

// Stop stops runner(s) and RPC service
func (gp *GoPlugin) Stop() error {
	return nil
}

const timeout = time.Minute * 10

// Downstream returns a connection to `ginsider`
func (gp *GoPlugin) Downstream(ctx context.Context) (*rpc.Client, error) {
	gp.clientMutex.Lock()
	defer gp.clientMutex.Unlock()

	if gp.client != nil {
		return gp.client, nil
	}

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

	gp.client.Close()
	gp.client = nil
}

func (gp *GoPlugin) callClientWithReconnect(ctx context.Context, method string, req interface{}, res interface{}) error {
	inslogger.FromContext(ctx).Debug("GoPlugin.callClientWithReconnect starts")
	var err error
	var client *rpc.Client

	for {
		inslogger.FromContext(ctx).Info("Connect to insgorund")
		client, err = gp.Downstream(ctx)
		if err == nil {
			call := <-client.Go(method, req, res, nil).Done
			err = call.Error

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
	ctx context.Context, callContext *core.LogicCallContext,
	code core.RecordRef, data []byte,
	method string, args core.Arguments,
) (
	[]byte, core.Arguments, error,
) {
	inslogger.FromContext(ctx).Debug("GoPlugin.CallMethod starts")
	start := time.Now()

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
		callTime := time.Since(start)
		metrics.GopluginContractExecutionTime.Observe(callTime.Seconds())
		inslogger.FromContext(ctx).Debugf("CallMethod done work, time spend in here - %s", callTime)
		if callResult.Error != nil {
			return nil, nil, errors.Wrap(callResult.Error, "problem with API call")
		}
		return callResult.Response.Data, callResult.Response.Ret, nil
	case <-time.After(timeout):
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
	ctx context.Context, callContext *core.LogicCallContext,
	code core.RecordRef, name string, args core.Arguments,
) (
	[]byte, error,
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
			return nil, errors.Wrap(callResult.Error, "problem with API call")
		}
		return callResult.Response.Ret, nil
	case <-time.After(timeout):
		return nil, errors.New("logicrunner execution timeout")
	}
}
