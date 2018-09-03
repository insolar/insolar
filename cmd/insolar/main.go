/*
 *    Copyright 2018 INS Ecosystem
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
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/insolar/insolar/configuration"
	log "github.com/sirupsen/logrus"
	jww "github.com/spf13/jwalterweatherman"
)

func main() {
	jww.SetStdoutThreshold(jww.LevelDebug)
	cfgHolder := configuration.NewHolder()
	err := cfgHolder.Load()
	if err != nil {
		log.Warnln(err.Error())
	}

	cfgHolder.Configuration.Host.Transport.BehindNAT = false
	cfgHolder.Configuration.Ledger.DataDirectory = "./data"

	initLogger(cfgHolder.Configuration.Log)

	fmt.Print("Starts with configuration:\n", configuration.ToString(cfgHolder.Configuration))
	network, err := StartNetwork(cfgHolder.Configuration)
	if err != nil {
		log.Fatalln("Failed to start network: ", err.Error())
	}

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		log.Debugln("caught sig: ", sig)
		network.closeNetwork()
		os.Exit(0)
	}()

	/*
		ledger, err := ledger.NewLedger(cfgHolder.Configuration.Ledger)
		if err != nil {
			log.Fatalln("Failed to create ledger: ", err.Error())
		}

		ledger.GetManager()
		//logicrunner.NewLogicRunner(ledger.GetManager())
		//messagerouter.MessageRouter{}
	*/

	go handleStats(cfgHolder.Configuration.Stats, network)

	fmt.Println("Running interactive mode:")
	repl(network.HostNetwork, network.ctx)
}

func initLogger(cfg configuration.Log) {

	cfg.Level = "debug"
	level, err := log.ParseLevel(strings.ToLower(cfg.Level))
	if err != nil {
		log.Warnln(err.Error())
	}
	jww.SetLogOutput(log.StandardLogger().Out)
	log.SetLevel(level)
}
