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

// Need to fix problem with start pulsar
func main() {
	jww.SetStdoutThreshold(jww.LevelDebug)
	cfgHolder := configuration.NewHolder()
	err := cfgHolder.Load()
	if err != nil {
		log.Warnln("failed to load configuration from file: ", err.Error())
	}

	cfgHolder.Configuration.Log.Level = "Debug"
	err = cfgHolder.LoadEnv()
	if err != nil {
		log.Warnln("failed to load configuration from env:", err.Error())
	}
	initLogger(cfgHolder.Configuration.Log)
	server, storage := initPulsar(cfgHolder.Configuration)

	go server.StartServer()
	pulseTicker, refreshTicker := runPulsar(server, cfgHolder.Configuration.Pulsar)

	fmt.Println("Press any button to exit")
	time.Sleep(2 * time.Hour)
	//rl, err := readline.New("> ")
	//_, err = rl.Readline()
	//if err != nil {
	//	log.Warn(err)
	//}

	defer func() {
		pulseTicker.Stop()
		refreshTicker.Stop()
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

func initPulsar(cfg configuration.Configuration) (*pulsar.Pulsar, pulsarstorage.PulsarStorage) {
	fmt.Print("Starts with configuration:\n", configuration.ToString(cfg))
	fmt.Println("Version: ", version.GetFullVersion())

	storage, err := pulsarstorage.NewStorageBadger(cfg.Pulsar, nil)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	switcher := &pulsar.StateSwitcherImpl{}
	server, err := pulsar.NewPulsar(
		cfg.Pulsar,
		storage,
		&pulsar.RPCClientWrapperFactoryImpl{},
		&pulsar.StandardEntropyGenerator{},
		switcher,
		net.Listen,
		cfg.PrivateKey)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	switcher.SetPulsar(server)

	return server, storage
}

func runPulsar(server *pulsar.Pulsar, cfg configuration.Pulsar) (pulseTicker *time.Ticker, refreshTicker *time.Ticker) {
	server.CheckConnectionsToPulsars()

	nextPulseNumber := core.CalculatePulseNumber(time.Now())
	err := server.StartConsensusProcess(nextPulseNumber)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	pulseTicker = time.NewTicker(time.Duration(cfg.PulseTime) * time.Millisecond)
	go func() {
		for range pulseTicker.C {
			err = server.StartConsensusProcess(core.PulseNumber(server.LastPulse.PulseNumber + 10))
			if err != nil {
				log.Fatal(err)
				panic(err)
			}
		}
	}()

	refreshTicker = time.NewTicker(1 * time.Second)
	go func() {
		for range refreshTicker.C {
			server.CheckConnectionsToPulsars()
		}
	}()

	return
}
