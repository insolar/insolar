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
	"os/exec"
	"path/filepath"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
	"github.com/insolar/insolar/testutils"
	"github.com/spf13/pflag"
)

func main() {
	contractPath := pflag.StringP("contract-path", "c", "", "path to healthcheck contract")
	cacheDir := pflag.StringP("cache-dir", "d", "", "path to insgorund cache directory")
	rpcAddress := pflag.StringP("rpc", "a", "", "address and port of RPC API")
	rpcProtocol := pflag.StringP("rpc-proto", "p", "", "protocol of RPC API")
	pflag.Parse()

	err := log.SetLevel("Debug")
	if err != nil {
		log.Errorln(err.Error())
	}

	if *cacheDir == "" {
		log.Error("need to provide path to insgorund cache directory")
		os.Exit(2)
	}

	ref := core.RecordRef{}.FromSlice(append(make([]byte, 63), 1))
	destination := filepath.Join(*cacheDir, ref.String())

	log.Error("destination: " + destination)
	log.Error("contractPath: " + *contractPath)

	_, err = exec.Command("./bin/insgocc", "compile", "-o", destination, *contractPath).CombinedOutput()
	if err != nil {
		log.Errorln(err.Error())
		os.Exit(2)
	}

	client, err := rpc.Dial(*rpcProtocol, *rpcAddress)
	//_, err = rpc.Dial(*rpcProtocol, *rpcAddress)
	if err != nil {
		log.Errorln(err.Error())
		os.Exit(2)
	}

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

	os.Exit(0)
}
