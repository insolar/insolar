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
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"

	"time"

	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
	"github.com/insolar/insolar/messagerouter"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/pkg/errors"
)

// Options of the GoPlugin
type Options struct {
	// Listen  is address `GoPlugin` listens on and provides RPC interface for runner(s)
	Listen string
	// CodePath is path to directory with plugin's code, this should go away at some point
	CodePath string
}

// RunnerOptions - set of options to control internal isolated code runner(s)
type RunnerOptions struct {
	// Listen is address the runner listens on and provides RPC interface for the `GoPlugin`
	Listen string
	// CodeStoragePath is path to directory where the runner caches code
	CodeStoragePath string
}

// MessageRouter interface
type MessageRouter interface {
	Route(msg servicenetwork.Message) (resp messagerouter.Response, err error)
}

// GoPlugin is a logic runner of code written in golang and compiled as go plugins
type GoPlugin struct {
	Options       Options
	RunnerOptions RunnerOptions
	MessageRouter MessageRouter
	sock          net.Listener
	runner        *exec.Cmd
}

// RPC is a RPC interface for runner to use for various tasks, e.g. code fetching
type RPC struct {
	gp *GoPlugin
}

// GetObject is an RPC retrieving an object by its reference, so far short circuted to return
// code of the plugin
func (gpr *RPC) GetObject(ref logicrunner.Reference, reply *logicrunner.Object) error {
	f, err := os.Open(gpr.gp.Options.CodePath + string(ref) + ".so")
	if err != nil {
		return err
	}
	reply.Data, err = ioutil.ReadAll(f)
	return err
}

// RouteCall routes call from a contract to a contract through message router
func (gpr *RPC) RouteCall(req rpctypes.UpRouteReq, reply *rpctypes.UpRouteResp) error {
	if gpr.gp.MessageRouter == nil {
		return errors.New("message router was not set during initialization")
	}

	msg := servicenetwork.Message{
		Reference: string(req.Reference),
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

// NewGoPlugin returns a new started GoPlugin
func NewGoPlugin(options Options, runnerOptions RunnerOptions, mr MessageRouter) (*GoPlugin, error) {
	gp := GoPlugin{
		Options:       options,
		RunnerOptions: runnerOptions,
		MessageRouter: mr,
	}

	if gp.Options.Listen == "" {
		gp.Options.Listen = "127.0.0.1:7777"
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

	l, e := net.Listen("tcp", gp.Options.Listen)
	if e != nil {
		log.Fatal("couldn't setup listener on '"+gp.Options.Listen+"': ", e)
	}
	gp.sock = l
	log.Printf("starting goplugin RPC service on %q", gp.Options.Listen)
	_ = http.Serve(l, nil)
	log.Printf("STOP")
}

// StartRunner starts ginsider process
func (gp *GoPlugin) StartRunner() error {
	var runnerArguments []string
	if gp.RunnerOptions.Listen != "" {
		runnerArguments = append(runnerArguments, "-l", gp.RunnerOptions.Listen)
	} else {
		return errors.New("listen is not optional in gp.RunnerOptions")
	}
	if gp.RunnerOptions.CodeStoragePath != "" {
		runnerArguments = append(runnerArguments, "-d", gp.RunnerOptions.CodeStoragePath)
	}
	runnerArguments = append(runnerArguments, "--rpc", gp.Options.Listen)

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
func (gp *GoPlugin) Stop() {
	err := gp.runner.Process.Kill()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(200 * time.Millisecond)

	if gp.sock != nil {
		err = gp.sock.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}

const timeout = time.Second * 5

// Exec runs a method on an object in controlled environment
func (gp *GoPlugin) Exec(codeRef logicrunner.Reference, data []byte, method string, args logicrunner.Arguments) ([]byte, logicrunner.Arguments, error) {
	client, err := rpc.DialHTTP("tcp", gp.RunnerOptions.Listen)
	if err != nil {
		return nil, nil, errors.Wrap(err, "problem with rpc connection")
	}
	res := rpctypes.DownCallResp{}

	req := rpctypes.DownCallReq{Reference: codeRef, Data: data, Method: method, Arguments: args}

	select {
	case call := <-client.Go("RPC.Call", req, &res, nil).Done:
		if call.Error != nil {
			return nil, nil, errors.Wrap(call.Error, "problem with API call")
		}
	case <-time.After(timeout):
		return nil, nil, errors.New("timeout")
	}
	return res.Data, res.Ret, res.Err
}
