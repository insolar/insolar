package goplugin

import (
	"net/rpc"

	"github.com/insolar/insolar/logicrunner"
)

type GoPlugin struct {
	DockerAddr string
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
}

func New() GoPlugin {
	return GoPlugin{} // TODO
}
