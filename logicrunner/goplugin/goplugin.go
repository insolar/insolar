package goplugin

import (
	"net/rpc"

	"os/exec"

	"time"

	"github.com/insolar/insolar/logicrunner"
)

type GoPlugin struct {
	DockerAddr string
	DockerCmd  *exec.Cmd
}

func NewGoPlugin(addr string) (*GoPlugin, error) {
	gp := GoPlugin{
		DockerAddr: addr,
		DockerCmd:  exec.Command("ginsider/ginsider"),
	}
	gp.DockerCmd.Start()
	time.Sleep(100 * time.Millisecond)
	go gp.Start()
	return &gp, nil
}

func (gp *GoPlugin) Start() {
}

func (gp *GoPlugin) Stop() {
	gp.DockerCmd.Process.Kill()
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
