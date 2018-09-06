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
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/messagerouter"
	"github.com/insolar/insolar/network/servicenetwork"
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

	components := make(core.Components)

	cfgHolder.Configuration.Host.Transport.BehindNAT = false

	initLogger(cfgHolder.Configuration.Log)

	fmt.Print("Starts with configuration:\n", configuration.ToString(cfgHolder.Configuration))
	network, err := StartNetwork(cfgHolder.Configuration)
	if err != nil {
		log.Fatalln("Failed to start network: ", err.Error())
	}
	components["core.Network"] = network

	ledger := StartLedger(cfgHolder.Configuration.Ledger)
	components["core.Ledger"] = ledger

	logicrunner1, _ := logicrunner.NewLogicRunner(cfgHolder.Configuration.LogicRunner)
	// todo: check err
	components["core.LogicRunner"] = logicrunner1

	mr, _ := messagerouter.New(cfgHolder.Configuration)
	components["core.MessageRouter"] = mr

	// start all components
	for n, c := range components {
		err := c.Start(components)
		if err != nil {
			log.Errorln("failed to start component ", n, " : ", err.Error())
		}
	}

	// TODO: call Start() on all components.

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		log.Debugln("caught sig: ", sig)
		err := network.Stop()
		if err != nil {
			log.Println(err.Error())
		}
		// TODO: call Stop() on all components.
		lComp := ledger.(core.Component)
		if err := lComp.Stop(); err != nil {
			log.Warnln("ledger teardown failed:", err.Error())
		}
		os.Exit(0)
	}()

	go handleStats(cfgHolder.Configuration.Stats)

	fmt.Println("Running interactive mode:")
	dht, ctx := network.(*servicenetwork.ServiceNetwork).GetHostNetwork()
	repl(network, dht, ctx)
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
