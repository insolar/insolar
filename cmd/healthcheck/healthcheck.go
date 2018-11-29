/*
 *    Copyright 2018 Insolar
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

package main

import (
	"net/rpc"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
	"github.com/insolar/insolar/testutils"
)

func main() {
	rpcAddress := pflag.StringP("rpc", "a", "", "address and port of RPC API")
	rpcProtocol := pflag.StringP("rpc-proto", "p", "tcp", "protocol of RPC API, tcp by default")
	refString := pflag.StringP("ref", "r", "", "ref of healthcheck contract")
	pflag.Parse()

	if *rpcAddress == "" || *rpcProtocol == "" || *refString == "" {
		log.Errorln(errors.New("need to provide all params"))
		os.Exit(2)
	}

	client, err := rpc.Dial(*rpcProtocol, *rpcAddress)
	if err != nil {
		log.Errorln(err.Error())
		os.Exit(2)
	}

	ref := core.NewRefFromBase58(*refString)

	empty, _ := core.Serialize([]interface{}{})

	caller := testutils.RandomRef()
	res := rpctypes.DownCallMethodResp{}
	req := rpctypes.DownCallMethodReq{
		Context:   &core.LogicCallContext{Caller: &caller},
		Code:      ref,
		Data:      empty,
		Method:    "Check",
		Arguments: empty,
	}

	err = client.Call("RPC.CallMethod", req, &res)
	if err != nil {
		log.Errorln(err.Error())
		os.Exit(2)
	}
}
