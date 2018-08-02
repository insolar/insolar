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
func (gpr *GoPluginRPC) GetObject(ref logicrunner.Reference, reply *logicrunner.Object) error {
	f, err := os.Open(gpr.gp.CodeDir + string(ref) + ".so")
	if err != nil {
		return err
	}
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
	Data []byte
	Ret  logicrunner.Arguments
	Err  error
}

func (gp *GoPlugin) Exec(object logicrunner.Object, method string, args logicrunner.Arguments) ([]byte, logicrunner.Arguments, error) {
	client, err := rpc.DialHTTP("tcp", gp.DockerAddr)
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
