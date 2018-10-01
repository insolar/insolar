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
	"io/ioutil"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/insolar/insolar/log"
	"github.com/spf13/pflag"

	"github.com/insolar/insolar/logicrunner/goplugin/ginsider"
	"github.com/insolar/insolar/logicrunner/goplugin/proxyctx"
)

func main() {
	listen := pflag.StringP("listen", "l", ":7777", "address and port to listen")
	protocol := pflag.String("proto", "tcp", "listen protocol")
	path := pflag.StringP("directory", "d", "", "directory where to store code of go plugins")
	rpcAddress := pflag.String("rpc", "localhost:7778", "address and port of RPC API")
	rpcProtocol := pflag.String("rpc-proto", "tcp", "protocol of RPC API")
	pflag.Parse()

	err := log.SetLevel("Debug")
	if err != nil {
		log.Errorln(err.Error())
	}

	if *path == "" {
		tmpDir, err := ioutil.TempDir("", "contractcache-")
		if err != nil {
			log.Fatal("Couldn't create temp cache dir: ", err)
			os.Exit(1)
		}
		defer os.RemoveAll(tmpDir)
		*path = tmpDir
	}

	insider := ginsider.NewGoInsider(*path, *rpcProtocol, *rpcAddress)
	proxyctx.Current = insider

	err = rpc.Register(&ginsider.RPC{GI: insider})
	if err != nil {
		log.Fatal("Couldn't register RPC interface: ", err)
		os.Exit(1)
	}

	rpc.HandleHTTP()
	listener, err := net.Listen(*protocol, *listen)
	if err != nil {
		log.Fatal("couldn't setup listener on '"+*listen+"':", err)
		os.Exit(1)
	}

	log.Debug("ginsider launched, listens " + *listen)
	err = http.Serve(listener, nil)
	if err != nil {
		log.Fatal("couldn't start server: ", err)
		os.Exit(1)
	}
	log.Debug("bye\n")
}
