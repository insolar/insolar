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

type RunnerOptions struct {
	Listen      string
	StoragePath string
}

type GoPlugin struct {
	ListenAddr    string
	sock          net.Listener
	Runner        *exec.Cmd
	RunnerOptions RunnerOptions
	CodeDir       string
}

type GoPluginRPC struct {
	gp *GoPlugin
}

// returns code for
func (gpr *GoPluginRPC) GetObject(ref logicrunner.Reference, reply *logicrunner.Object) error {
	f, err := os.Open(gpr.gp.CodeDir + string(ref) + ".so")
	if err != nil {
		return err
	}
	reply.Data, err = ioutil.ReadAll(f)
	return err
}

func NewGoPlugin(addr string, runner_options RunnerOptions) (*GoPlugin, error) {
	gp := GoPlugin{
		ListenAddr:    addr,
		RunnerOptions: runner_options,
	}

	var runner_arguments []string
	if runner_options.Listen != "" {
		runner_arguments = append(runner_arguments, "-s", runner_options.Listen)
	} else {
		return nil, errors.New("Listen is not optional in runner_options")
	}
	if runner_options.StoragePath != "" {
		runner_arguments = append(runner_arguments, "-d", runner_options.StoragePath)
	}
	runner := exec.Command("ginsider/ginsider", runner_arguments...)
	err := runner.Start()
	if err != nil {
		return nil, err
	}

	gp.Runner = runner

	time.Sleep(2000 * time.Millisecond)
	go gp.Start()
	return &gp, nil
}

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

func (gp *GoPlugin) Stop() {
	gp.Runner.Process.Kill()
	gp.sock.Close()
}

type CallReq struct {
	Object logicrunner.Object
	Method string
	Args   logicrunner.Arguments
}

type CallResp struct {
	Data []byte
	Ret  logicrunner.Arguments
	Err  error
}

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
