package goplugin

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/exec"
	"time"

	"github.com/insolar/insolar/logicrunner"
)

// RunnerOptions - set of options to control internal isolated code runner(s)
type RunnerOptions struct {
	Listen      string
	StoragePath string
}

// GoPlugin is a logic runner of code written in golang and compiled as go plugins
type GoPlugin struct {
	ListenAddr    string
	sock          net.Listener
	Runner        *exec.Cmd
	RunnerOptions RunnerOptions
	CodeDir       string
}

// GoPluginRPC is a RPC interface for runner to use for variouse tasks, e.g. code fetching
type GoPluginRPC struct {
	gp *GoPlugin
}

// GetObject is an RPC retriving an object by its reference, so far short circueted to return
// code of the plugin
func (gpr *GoPluginRPC) GetObject(ref logicrunner.Reference, reply *logicrunner.Object) error {
	f, err := os.Open(gpr.gp.CodeDir + string(ref) + ".so")
	if err != nil {
		return err
	}
	reply.Data, err = ioutil.ReadAll(f)
	return err
}

// NewGoPlugin returns a new started GoPlugin
func NewGoPlugin(addr string, runnerOptions RunnerOptions) (*GoPlugin, error) {
	gp := GoPlugin{
		ListenAddr:    addr,
		RunnerOptions: runnerOptions,
	}

	var runnerArguments []string
	if runnerOptions.Listen != "" {
		runnerArguments = append(runnerArguments, "-s", runnerOptions.Listen)
	} else {
		return nil, errors.New("Listen is not optional in runnerOptions")
	}
	if runnerOptions.StoragePath != "" {
		runnerArguments = append(runnerArguments, "-d", runnerOptions.StoragePath)
	}
	runner := exec.Command("ginsider/ginsider", runnerArguments...)
	err := runner.Start()
	if err != nil {
		return nil, err
	}

	gp.Runner = runner

	time.Sleep(2000 * time.Millisecond)
	go gp.Start()
	return &gp, nil
}

// Start starts runner and RPC interface to help runner, note that NewGoPlugin does
// this for you
func (gp *GoPlugin) Start() {
	r := GoPluginRPC{gp: gp}
	rpc.Register(r)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", gp.ListenAddr)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	gp.sock = l
	http.Serve(l, nil)
}

// Stop stops runner(s) and RPC service
func (gp *GoPlugin) Stop() {
	gp.Runner.Process.Kill()
	gp.sock.Close()
}

// CallReq is a set of arguments for Call RPC in the runner
type CallReq struct {
	Object logicrunner.Object
	Method string
	Args   logicrunner.Arguments
}

// CallResp is response from Call RPC in the runner
type CallResp struct {
	Data []byte
	Ret  logicrunner.Arguments
	Err  error
}

// Exec runs a method on an object in controlled environment
func (gp *GoPlugin) Exec(object logicrunner.Object, method string, args logicrunner.Arguments) ([]byte, logicrunner.Arguments, error) {
	client, err := rpc.DialHTTP("tcp", gp.RunnerOptions.Listen)
	if err != nil {
		return nil, nil, err
	}
	res := CallResp{}
	err = client.Call("GoInsider.Call", CallReq{Object: object, Method: method, Args: args}, &res)
	if err != nil {
		return nil, nil, err
	}
	return res.Data, res.Ret, res.Err
}
