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
	"net/rpc"
	"os/exec"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
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
	client          *rpc.Client
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

// Downstream returns a connection to `ginsider`
func (gp *GoPlugin) Downstream() (*rpc.Client, error) {
	if gp.client != nil {
		return gp.client, nil
	}

	client, err := rpc.DialHTTP(gp.Cfg.GoPlugin.RunnerProtocol, gp.Cfg.GoPlugin.RunnerListen)
	if err != nil {
		return nil, errors.Wrapf(
			err, "couldn't dial '%s' over %s",
			gp.Cfg.GoPlugin.RunnerListen, gp.Cfg.GoPlugin.RunnerProtocol,
		)
	}

	gp.client = client
	return gp.client, nil
}

const timeout = time.Second * 60

// CallMethod runs a method on an object in controlled environment
func (gp *GoPlugin) CallMethod(ctx *core.LogicCallContext, code core.RecordRef, data []byte, method string, args core.Arguments) ([]byte, core.Arguments, error) {
	client, err := gp.Downstream()
	if err != nil {
		return nil, nil, errors.Wrap(err, "problem with rpc connection")
	}

	res := rpctypes.DownCallMethodResp{}
	req := rpctypes.DownCallMethodReq{
		Context:   ctx,
		Code:      code,
		Data:      data,
		Method:    method,
		Arguments: args,
	}

	select {
	case call := <-client.Go("RPC.CallMethod", req, &res, nil).Done:
		if call.Error != nil {
			return nil, nil, errors.Wrap(call.Error, "problem with API call")
		}
	case <-time.After(timeout):
		return nil, nil, errors.New("timeout")
	}
	return res.Data, res.Ret, nil
}

// CallConstructor runs a constructor of a contract in controlled environment
func (gp *GoPlugin) CallConstructor(ctx *core.LogicCallContext, code core.RecordRef, name string, args core.Arguments) ([]byte, error) {
	client, err := gp.Downstream()
	if err != nil {
		return nil, errors.Wrap(err, "problem with rpc connection")
	}

	res := rpctypes.DownCallConstructorResp{}
	req := rpctypes.DownCallConstructorReq{Code: code, Name: name, Arguments: args}

	select {
	case call := <-client.Go("RPC.CallConstructor", req, &res, nil).Done:
		if call.Error != nil {
			return nil, errors.Wrap(call.Error, "problem with API call")
		}
	case <-time.After(timeout):
		return nil, errors.New("timeout")
	}
	return res.Ret, nil
}
