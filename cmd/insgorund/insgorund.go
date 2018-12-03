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
	"context"
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"strings"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/metrics"

	"github.com/spf13/pflag"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/goplugin/ginsider"
)

func main() {
	listen := pflag.StringP("listen", "l", ":7777", "address and port to listen")
	protocol := pflag.String("proto", "tcp", "listen protocol")
	path := pflag.StringP("directory", "d", "", "directory where to store code of go plugins")
	rpcAddress := pflag.String("rpc", "localhost:7778", "address and port of RPC API")
	rpcProtocol := pflag.String("rpc-proto", "tcp", "protocol of RPC API")
	code := pflag.String("code", "", "add pre-compiled code to cache (<ref>:</path/to/plugin.so>)")

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
		log.Debug("ginsider cache dir is " + tmpDir)
	}

	insider := ginsider.NewGoInsider(*path, *rpcProtocol, *rpcAddress)

	if *code != "" {
		codeSlice := strings.Split(*code, ":")
		if len(codeSlice) != 2 {
			log.Fatal("code param format is <ref>:</path/to/plugin.so>")
			os.Exit(1)
		}
		ref := core.NewRefFromBase58(codeSlice[0])
		pluginPath := codeSlice[1]

		err := insider.AddPlugin(ref, pluginPath)
		if err != nil {
			log.Fatalf("Couldn't add plugin by ref %s with .so from %s, err: %s ", ref, pluginPath, err.Error())
			os.Exit(1)
		}
	}

	err = rpc.Register(&ginsider.RPC{GI: insider})
	if err != nil {
		log.Fatal("Couldn't register RPC interface: ", err)
		os.Exit(1)
	}

	listener, err := net.Listen(*protocol, *listen)
	if err != nil {
		log.Fatal("couldn't setup listener on '"+*listen+"':", err)
		os.Exit(1)
	}

	runMetrics()

	log.Debug("ginsider launched, listens " + *listen)
	rpc.Accept(listener)
	log.Debug("bye\n")
}

func runMetrics() {
	log.Debug("ginsider start metrics")

	// TODO copy-pasted from configuration/metrics.NewMetrics()
	metricsConfiguration := configuration.Metrics{
		ListenAddress: "0.0.0.0:9090",
		Namespace:     "insolar",
		ZpagesEnabled: true,
	}
	ctx := context.TODO()

	// TODO make it right
	m, err := metrics.NewMetrics(ctx, metricsConfiguration)
	if err != nil {
		log.Fatal("couldn't setup metrics ", err)
		os.Exit(1)
	}
	err = m.Start(ctx)
	if err != nil {
		log.Fatal("couldn't setup metrics ", err)
		os.Exit(1)
	}

	log.Debug("ginsider metrics start successfully")
}
