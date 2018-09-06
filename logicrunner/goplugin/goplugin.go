/*
 *    Copyright 2018 INS Ecosystem
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
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
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
	Cfg             *configuration.GoPlugin
	MessageRouter   core.MessageRouter
	ArtifactManager core.ArtifactManager
	sock            net.Listener
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
	am.SetArchPref([]core.MachineType{core.MachineTypeGoPlugin})
	codeDescriptor, err := am.GetCode(req.Reference)
	if err != nil {
		return err
	}

	code, err := codeDescriptor.Code()
	if err != nil {
		return err
	}

	reply.Code = code
	return nil
}

// RouteCall routes call from a contract to a contract through message router
func (gpr *RPC) RouteCall(req rpctypes.UpRouteReq, reply *rpctypes.UpRouteResp) error {
	if gpr.gp.MessageRouter == nil {
		return errors.New("message router was not set during initialization")
	}

	msg := core.Message{
		Reference: req.Reference,
		Method:    req.Method,
		Arguments: req.Arguments,
	}

	res, err := gpr.gp.MessageRouter.Route(msg)
	if err != nil {
		return errors.Wrap(err, "couldn't route message")
	}
	if reply.Err != nil {
		return errors.Wrap(reply.Err, "couldn't route message (error in respone)")
	}

	reply.Result = res.Result

	return nil
}

// RouteConstructorCall routes call from a contract to a constructor of another contract
func (gpr *RPC) RouteConstructorCall(req rpctypes.UpRouteConstructorReq, reply *rpctypes.UpRouteConstructorResp) error {
	if gpr.gp.MessageRouter == nil {
		return errors.New("message router was not set during initialization")
	}

	msg := core.Message{
		Constructor: true,
		Reference:   req.Reference,
		Method:      req.Constructor,
		Arguments:   req.Arguments,
	}

	res, err := gpr.gp.MessageRouter.Route(msg)
	if err != nil {
		return errors.Wrap(err, "couldn't route message")
	}
	if reply.Err != nil {
		return errors.Wrap(reply.Err, "couldn't route message (error in respone)")
	}

	// TODO: store data on ledger via artifact manager
	_ = res.Data

	reply.Reference = core.String2Ref("some-ref")

	return nil
}

// NewGoPlugin returns a new started GoPlugin
func NewGoPlugin(conf *configuration.GoPlugin, mr core.MessageRouter, am core.ArtifactManager) (*GoPlugin, error) {
	gp := GoPlugin{
		Cfg:             conf,
		MessageRouter:   mr,
		ArtifactManager: am,
	}
	if gp.Cfg.MainListen == "" {
		gp.Cfg.MainListen = "127.0.0.1:7777"
	}

	err := gp.StartRunner()
	if err != nil {
		return nil, err
	}

	go gp.Start()
	return &gp, nil
}

var rpcService *RPC

// Start starts RPC interface to help runner, note that NewGoPlugin does
// this for you
func (gp *GoPlugin) Start() {
	if rpcService == nil {
		rpcService = &RPC{}
		_ = rpc.Register(rpcService)
		rpc.HandleHTTP()
	}
	rpcService.gp = gp

	l, e := net.Listen("tcp", gp.Cfg.MainListen)
	if e != nil {
		log.Fatal("couldn't setup listener on '"+gp.Cfg.MainListen+"': ", e)
	}
	gp.sock = l
	log.Printf("starting goplugin RPC service on %q", gp.Cfg.MainListen)
	_ = http.Serve(l, nil)
	log.Printf("STOP")
}

// StartRunner starts ginsider process
func (gp *GoPlugin) StartRunner() error {
	var runnerArguments []string
	if gp.Cfg.RunnerListen != "" {
		runnerArguments = append(runnerArguments, "-l", gp.Cfg.RunnerListen)
	} else {
		return errors.New("listen is not optional in gp.RunnerOptions")
	}
	if gp.Cfg.RunnerCodePath != "" {
		runnerArguments = append(runnerArguments, "-d", gp.Cfg.RunnerCodePath)
	}
	runnerArguments = append(runnerArguments, "--rpc", gp.Cfg.MainListen)

	runner := exec.Command("ginsider-cli/ginsider-cli", runnerArguments...)
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
	err := gp.runner.Process.Kill()
	if err != nil {
		return err
	}
	time.Sleep(200 * time.Millisecond)

	if gp.sock != nil {
		err = gp.sock.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// Downstream returns a connection to `ginsider`
func (gp *GoPlugin) Downstream() (*rpc.Client, error) {
	if gp.client != nil {
		return gp.client, nil
	}

	client, err := rpc.DialHTTP("tcp", gp.Cfg.RunnerListen)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't dial '%s'", gp.Cfg.RunnerListen)
	}

	gp.client = client
	return gp.client, nil
}

const timeout = time.Second * 5

// CallMethod runs a method on an object in controlled environment
func (gp *GoPlugin) CallMethod(codeRef core.RecordRef, data []byte, method string, args core.Arguments) ([]byte, core.Arguments, error) {
	client, err := gp.Downstream()
	if err != nil {
		return nil, nil, errors.Wrap(err, "problem with rpc connection")
	}

	res := rpctypes.DownCallMethodResp{}
	req := rpctypes.DownCallMethodReq{Reference: codeRef, Data: data, Method: method, Arguments: args}

	select {
	case call := <-client.Go("RPC.CallMethod", req, &res, nil).Done:
		if call.Error != nil {
			return nil, nil, errors.Wrap(call.Error, "problem with API call")
		}
	case <-time.After(timeout):
		return nil, nil, errors.New("timeout")
	}
	return res.Data, res.Ret, res.Err
}

// CallConstructor runs a constructor of a contract in controlled environment
func (gp *GoPlugin) CallConstructor(codeRef core.RecordRef, name string, args core.Arguments) ([]byte, error) {
	client, err := gp.Downstream()
	if err != nil {
		return nil, errors.Wrap(err, "problem with rpc connection")
	}

	res := rpctypes.DownCallConstructorResp{}
	req := rpctypes.DownCallConstructorReq{Reference: codeRef, Name: name, Arguments: args}

	select {
	case call := <-client.Go("RPC.CallConstructor", req, &res, nil).Done:
		if call.Error != nil {
			return nil, errors.Wrap(call.Error, "problem with API call")
		}
	case <-time.After(timeout):
		return nil, errors.New("timeout")
	}
	return res.Ret, res.Err
}
