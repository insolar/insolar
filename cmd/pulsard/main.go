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
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/pulsar"
	"github.com/insolar/insolar/pulsar/storage"
	jww "github.com/spf13/jwalterweatherman"
)

func main() {
	jww.SetStdoutThreshold(jww.LevelDebug)
	cfgHolder := configuration.NewHolder()
	err := cfgHolder.Load()
	if err != nil {
		log.Warnln("Failed to load configuration from file: ", err.Error())
	}

	err = cfgHolder.LoadEnv()
	if err != nil {
		log.Warnln("Failed to load configuration from env:", err.Error())
	}

	initLogger(cfgHolder.Configuration.Log)
	fmt.Print("Starts with configuration:\n", configuration.ToString(cfgHolder.Configuration))

	storage, err := pulsarstorage.NewStorageBadger(cfgHolder.Configuration.Pulsar, nil)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	server, err := pulsar.NewPulsar(cfgHolder.Configuration.Pulsar, storage, net.Listen)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	go server.StartServer()

	server.RefreshConnections()
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			server.RefreshConnections()
		}
	}()

	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}

	fmt.Println("Press any button to exit")
	_, _ = rl.Readline()

	// Need to think about the shutdown mechanism
	ticker.Stop()
	defer server.StopServer()
}

func initLogger(cfg configuration.Log) {
	err := log.SetLevel(strings.ToLower(cfg.Level))
	if err != nil {
		log.Errorln(err.Error())
	}
}
