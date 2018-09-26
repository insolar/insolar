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
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/eventbus/event"
	"github.com/insolar/insolar/eventbus/reaction"
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
	EventBus        core.EventBus
	ArtifactManager core.ArtifactManager
	runner          *exec.Cmd
	client          *rpc.Client
}

// RPC is a RPC interface for runner to use for various tasks, e.g. code fetching
type RPC struct {
	gp *GoPlugin
}

// GetCode is an RPC retrieving a code by its reference
func (gpr *RPC) GetCode(req rpctypes.UpGetCodeReq, reply *rpctypes.UpGetCodeResp) error {
	am := gpr.gp.ArtifactManager
	codeDescriptor, err := am.GetCode(req.Code, []core.MachineType{req.MType})
	if err != nil {
		return err
	}
	reply.Code = codeDescriptor.Code()
	return nil
}

// RouteCall routes call from a contract to a contract through event bus.
func (gpr *RPC) RouteCall(req rpctypes.UpRouteReq, reply *rpctypes.UpRouteResp) error {
	if gpr.gp.EventBus == nil {
		return errors.New("event bus was not set during initialization")
	}

	var mode event.MethodReturnMode
	if req.Wait {
		mode = event.ReturnResult
	} else {
		mode = event.ReturnNoWait
	}

	e := &event.CallMethod{
		ReturnMode: mode,
		ObjectRef:  req.Object,
		Method:     req.Method,
		Arguments:  req.Arguments,
	}

	res, err := gpr.gp.EventBus.Dispatch(e)
	if err != nil {
		return errors.Wrap(err, "couldn't dispatch event")
	}

	reply.Result = res.(*reaction.CommonReaction).Result

	return nil
}

// RouteConstructorCall routes call from a contract to a constructor of another contract
func (gpr *RPC) RouteConstructorCall(req rpctypes.UpRouteConstructorReq, reply *rpctypes.UpRouteConstructorResp) error {
	if gpr.gp.EventBus == nil {
		return errors.New("event bus was not set during initialization")
	}

	e := &event.CallConstructor{
		ClassRef:  req.Reference,
		Name:      req.Constructor,
		Arguments: req.Arguments,
	}

	res, err := gpr.gp.EventBus.Dispatch(e)
	if err != nil {
		return errors.Wrap(err, "couldn't dispatch event")
	}

	reply.Data = res.(*reaction.CommonReaction).Data
	return nil
}

// SaveAsChild is an RPC saving data as memory of a contract as child a parent
func (gpr *RPC) SaveAsChild(req rpctypes.UpSaveAsChildReq, reply *rpctypes.UpSaveAsChildResp) error {
	ref, err := gpr.gp.ArtifactManager.ActivateObject(
		core.RecordRef{}, core.RecordRef{}, req.Class, req.Parent, req.Data,
	)
	if err != nil {
		return errors.Wrap(err, "couldn't save new object")
	}
	reply.Reference = *ref
	return nil
}

// GetObjChildren is an RPC returns set of object children
func (gpr *RPC) GetObjChildren(req rpctypes.UpGetObjChildrenReq, reply *rpctypes.UpGetObjChildrenResp) error {
	// TODO: INS-408
	am := gpr.gp.ArtifactManager
	obj, err := am.GetObject(req.Obj, nil)
	if err != nil {
		return errors.Wrap(err, "am.GetObjChildren failed")
	}
	i := obj.Children()
	for i.HasNext() {
		r, err := i.Next()
		if err != nil {
			return err
		}
		o, err := am.GetObject(r, nil)
		if err != nil {
			return errors.Wrap(err, "Have ref, have no object")
		}
		cd, err := o.ClassDescriptor(nil)
		if err != nil {
			return errors.Wrap(err, "Have ref, have no object")
		}
		ref := cd.HeadRef()
		if ref.Equal(req.Class) {
			reply.Children = append(reply.Children, r)
		}
	}
	return nil
}

// SaveAsDelegate is an RPC saving data as memory of a contract as child a parent
func (gpr *RPC) SaveAsDelegate(req rpctypes.UpSaveAsDelegateReq, reply *rpctypes.UpSaveAsDelegateResp) error {
	ref, err := gpr.gp.ArtifactManager.ActivateObjectDelegate(
		core.RecordRef{}, core.RecordRef{}, req.Class, req.Into, req.Data,
	)
	if err != nil {
		return errors.Wrap(err, "couldn't save delegate")
	}
	reply.Reference = *ref
	return nil
}

// GetDelegate is an RPC saving data as memory of a contract as child a parent
func (gpr *RPC) GetDelegate(req rpctypes.UpGetDelegateReq, reply *rpctypes.UpGetDelegateResp) error {
	am := gpr.gp.ArtifactManager
	ref, err := am.GetDelegate(req.Object, req.OfType)
	if err != nil {
		return err
	}

	reply.Object = *ref
	return nil
}

// NewGoPlugin returns a new started GoPlugin
func NewGoPlugin(conf *configuration.LogicRunner, eb core.EventBus, am core.ArtifactManager) (*GoPlugin, error) {
	gp := GoPlugin{
		Cfg:             conf,
		EventBus:        eb,
		ArtifactManager: am,
	}

	err := gp.StartRunner()
	if err != nil {
		return nil, err
	}
	return &gp, nil
}

// StartRunner starts ginsider process

func (gp *GoPlugin) StartRunner() error {
	var runnerArguments []string
	if gp.Cfg.GoPlugin.RunnerListen != "" {
		runnerArguments = append(runnerArguments, "-l", gp.Cfg.GoPlugin.RunnerListen)
	} else {
		return errors.New("RunnerListen is not set in the configuration of GoPlugin")
	}
	if gp.Cfg.GoPlugin.RunnerCodePath != "" {
		runnerArguments = append(runnerArguments, "-d", gp.Cfg.GoPlugin.RunnerCodePath)
	}
	runnerArguments = append(runnerArguments, "--rpc", gp.Cfg.RpcListen)

	if gp.Cfg.GoPlugin.RunnerPath == "" {
		return errors.New("RunnerPath is not set in the configuration of GoPlugin")
	}

	runner := exec.Command(gp.Cfg.GoPlugin.RunnerPath, runnerArguments...)
	runner.Stdout = os.Stdout
	runner.Stderr = os.Stderr
	err := runner.Start()
	if err != nil {
		return err
	}
	time.Sleep(200 * time.Millisecond)
	gp.runner = runner

	return nil
}

// Stop stops runner(s) and RPC service
func (gp *GoPlugin) Stop() error {
	err := gp.runner.Process.Signal(syscall.SIGINT)
	if err != nil {
		return err
	}
	time.Sleep(200 * time.Millisecond)

	return nil
}

// Downstream returns a connection to `ginsider`
func (gp *GoPlugin) Downstream() (*rpc.Client, error) {
	if gp.client != nil {
		return gp.client, nil
	}

	client, err := rpc.DialHTTP("tcp", gp.Cfg.GoPlugin.RunnerListen)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't dial '%s'", gp.Cfg.GoPlugin.RunnerListen)
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
