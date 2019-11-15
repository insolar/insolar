//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"context"
	"io/ioutil"
	"net"
	"net/rpc"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/metrics"

	"github.com/spf13/pflag"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/goplugin/ginsider"
)

func main() {
	listen := pflag.StringP("listen", "l", ":7777", "address and port to listen")
	protocol := pflag.String("proto", "tcp", "listen protocol")
	path := pflag.StringP("directory", "d", "", "directory where to store code of go plugins")
	rpcAddress := pflag.String("rpc", "localhost:7778", "address and port of RPC API")
	rpcProtocol := pflag.String("rpc-proto", "tcp", "protocol of RPC API")
	metricsAddress := pflag.String("metrics", "", "address and port of prometheus metrics")
	code := pflag.String("code", "", "add pre-compiled code to cache (<ref>:</path/to/plugin.so>)")
	logLevel := pflag.String("log-level", "debug", "log level")

	pflag.Parse()

	err := log.SetLevel(*logLevel)
	if err != nil {
		log.Fatalf("Couldn't set log level to %q: %s", *logLevel, err)
	}
	log.InitTicker()

	if *path == "" {
		tmpDir, err := ioutil.TempDir("", "funcTestContractcache-")
		if err != nil {
			log.Fatalf("Couldn't create temp cache dir: %s", err.Error())
		}
		defer func() {
			err := os.RemoveAll(tmpDir)
			if err != nil {
				log.Fatalf("Failed to clean up tmp dir: %s", err.Error())
			}
		}()
		*path = tmpDir
		log.Debug("ginsider cache dir is " + tmpDir)
	}

	insider := ginsider.NewGoInsider(*path, *rpcProtocol, *rpcAddress)

	if *code != "" {
		codeSlice := strings.Split(*code, ":")
		if len(codeSlice) != 2 {
			log.Fatal("code param format is <ref>:</path/to/plugin.so>")
		}
		ref, err := insolar.NewReferenceFromString(codeSlice[0])
		if err != nil {
			log.Fatalf("Couldn't parse ref: %s", err.Error())
		}
		pluginPath := codeSlice[1]

		err = insider.AddPlugin(*ref, pluginPath)
		if err != nil {
			log.Fatalf("Couldn't add plugin by ref %s with .so from %s, err: %s ", ref, pluginPath, err.Error())
		}
	}

	err = rpc.Register(&ginsider.RPC{GI: insider})
	if err != nil {
		log.Fatal("Couldn't register RPC interface: ", err)
	}

	listener, err := net.Listen(*protocol, *listen)
	if err != nil {
		log.Fatal("couldn't setup listener on '"+*listen+"':", err)
	}

	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	var waitChannel = make(chan bool)

	go func() {
		sig := <-gracefulStop
		log.Info("ginsider get signal: ", sig)
		close(waitChannel)
	}()

	if *metricsAddress != "" {
		ctx := context.Background() // TODO add tradeId and logger

		metricsConfiguration := configuration.Metrics{
			ListenAddress: *metricsAddress,
			Namespace:     "insgorund",
			ZpagesEnabled: true,
		}

		m := metrics.NewMetrics(metricsConfiguration, metrics.GetInsgorundRegistry(), "virtual")
		err = m.Start(ctx)
		if err != nil {
			log.Fatal("couldn't setup metrics ", err)
		}

		defer m.Stop(ctx) // nolint: errcheck
	}

	log.Debug("ginsider launched, listens " + *listen)
	go rpc.Accept(listener)

	<-waitChannel
	log.Debug("bye\n")
}
