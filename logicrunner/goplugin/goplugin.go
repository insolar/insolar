package goplugin

import (
	"net/rpc"

	"os/exec"

	"time"

	"log"
	"net"

	"net/http"

	"github.com/insolar/insolar/logicrunner"
)

type GoPlugin struct {
	DockerAddr string
	DockerCmd  *exec.Cmd
	ListenAddr string
	sock       net.Listener
	CodeDir    string
}

type GoPluginRPC struct{
	gp *GoPlugin
}

type GetObjectReq struct {
	Reference string
}

type GetObjectResp struct {
	Object logicrunner.Object
}

func (gpr *GoPluginRPC) GetObject(args GetObjectReq, reply *GetObjectResp) error {
	addr := args.Reference
	fname := gpr.gp.CodeDir + addr // sorry generic
	reply.Object.Reference = fname // fix this
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
	r := GoPluginRPC{
		gp gp
	}
	rpc.Register(r)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", gp.ListenAddr)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	gp.sock = l
	gp.
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
	Ret logicrunner.Arguments
	Err error
}

func (gp *GoPlugin) Exec(object logicrunner.Object, method string, args logicrunner.Arguments) (ret logicrunner.Arguments, err error) {
	client, err := rpc.DialHTTP("tcp", gp.DockerAddr)
	if err != nil {
		return nil, err
	}
	Ret := CallResp{}
	err = client.Call("GoInsider.Call", CallReq{Object: object, Method: method, Args: args}, &Ret)
	if err != nil {
		return nil, err
	}
	ret = Ret.Ret
	err = Ret.Err
	return ret, err
}
