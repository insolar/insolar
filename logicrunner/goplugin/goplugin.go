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
)

type GoPlugin struct {
	DockerAddr string
	DockerCmd  *exec.Cmd
	ListenAddr string
	sock       net.Listener
	CodeDir    string
}

type GoPluginRPC struct {
	gp *GoPlugin
}

// returns code for
func (gpr *GoPluginRPC) GetObject(args logicrunner.Object, reply *logicrunner.Object) error {
	f, err := os.Open(gpr.gp.CodeDir + args.Reference + ".so")
	if err != nil {
		return err
	}
	reply.MachineType = args.MachineType
	reply.Data, err = ioutil.ReadAll(f)
	return err
}

func NewGoPlugin(addr string, myaddr string) (*GoPlugin, error) {
	gp := GoPlugin{
		DockerAddr: addr,
		DockerCmd:  exec.Command("ginsider/ginsider"),
		ListenAddr: myaddr,
	}
	gp.DockerCmd.Start()
	time.Sleep(200 * time.Millisecond)
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
	gp.DockerCmd.Process.Kill()
	gp.sock.Close()
}

type CallReq struct {
	Object logicrunner.Object
	Method string
	Args   logicrunner.Arguments
}

type CallResp struct {
	Object logicrunner.Object
	Ret    logicrunner.Arguments
	Err    error
}

func (gp *GoPlugin) Exec(object logicrunner.Object, method string, args logicrunner.Arguments) (retdata logicrunner.Arguments, ret logicrunner.Arguments, err error) {
	client, err := rpc.DialHTTP("tcp", gp.DockerAddr)
	if err != nil {
		return nil, nil, err
	}
	Ret := CallResp{}
	err = client.Call("GoInsider.Call", CallReq{Object: object, Method: method, Args: args}, &Ret)
	if err != nil {
		return nil, nil, err
	}
	retdata = Ret.Object.Data
	ret = Ret.Ret
	err = Ret.Err
	return retdata, ret, err
}
