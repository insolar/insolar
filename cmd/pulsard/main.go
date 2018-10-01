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

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/pulsar"
	"github.com/insolar/insolar/pulsar/storage"
	"github.com/insolar/insolar/version"
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
	fmt.Println("Version: ", version.GetFullVersion())

	storage, err := pulsarstorage.NewStorageBadger(cfgHolder.Configuration.Pulsar, nil)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	server, err := pulsar.NewPulsar(cfgHolder.Configuration.Pulsar, storage, &pulsar.RPCClientWrapperFactoryImpl{}, &pulsar.StandardEntropyGenerator{}, net.Listen)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	go server.StartServer()

	server.RefreshConnections()

	var nextPulseNumber core.PulseNumber
	if server.LastPulse.PulseNumber == core.FirstPulseNumber {
		nextPulseNumber = core.CalculatePulseNumber(time.Now())
	} else {
		waitTime := core.CalculateMsToNextPulse(server.LastPulse.PulseNumber, time.Now())
		if waitTime != 0 {
			nextPulseNumber = server.LastPulse.PulseNumber + core.PulseNumber(cfgHolder.Configuration.Pulsar.PulseTime)
		}
		time.Sleep(waitTime)

	}
	err = server.StartConsensusProcess(nextPulseNumber)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	//pulseTicker := time.NewTicker(time.Duration(cfgHolder.Configuration.Pulsar.PulseTime) * time.Second)
	//go func() {
	//	for range pulseTicker.C {
	//		err = server.StartConsensusProcess(core.PulseNumber(server.LastPulse.PulseNumber + 10))
	//		if err != nil {
	//			log.Fatal(err)
	//			panic(err)
	//		}
	//	}
	//}()
	//
	//refreshTicker := time.NewTicker(1 * time.Second)
	//go func() {
	//	for range refreshTicker.C {
	//		server.RefreshConnections()
	//	}
	//}()

	time.Sleep(10 * time.Minute)

	//fmt.Println("Press any button to exit")
	//_, err = rl.Readline()
	//if err != nil {
	//	log.Warn(err)
	//}

	//refreshTicker.Stop()
	defer func() {
		err := storage.Close()
		if err != nil {
			log.Error(err)
		}
		server.StopServer()
	}()

}

func initLogger(cfg configuration.Log) {
	err := log.SetLevel(strings.ToLower(cfg.Level))
	if err != nil {
		log.Errorln(err.Error())
	}
}
