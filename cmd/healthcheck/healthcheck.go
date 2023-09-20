package main

import (
	"net/rpc"
	"os"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
)

func main() {
	rpcAddress := pflag.StringP("rpc", "a", "", "address and port of RPC API")
	rpcProtocol := pflag.StringP("rpc-proto", "p", "tcp", "protocol of RPC API, tcp by default")
	refString := pflag.StringP("ref", "r", "", "ref of healthcheck contract")
	pflag.Parse()

	if *rpcAddress == "" || *rpcProtocol == "" || *refString == "" {
		log.Error(errors.New("need to provide all params"))
		os.Exit(2)
	}

	client, err := rpc.Dial(*rpcProtocol, *rpcAddress)
	if err != nil {
		log.Error(err.Error())
		os.Exit(2)
	}

	ref, err := insolar.NewReferenceFromString(*refString)
	if err != nil {
		log.Errorf("Failed to parse healthcheck contract ref: %s", err.Error())
		os.Exit(2)
	}

	empty, _ := insolar.Serialize([]interface{}{})

	caller := gen.Reference()
	res := rpctypes.DownCallMethodResp{}
	req := rpctypes.DownCallMethodReq{
		Context:   &insolar.LogicCallContext{Caller: &caller},
		Code:      *ref,
		Data:      empty,
		Method:    "Check",
		Arguments: empty,
	}

	err = client.Call("RPC.CallMethod", req, &res)
	if err != nil {
		log.Error(err.Error())
		os.Exit(2)
	}
}
