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
	"os"
	"os/signal"

	"github.com/insolar/insolar/configuration"
	log "github.com/sirupsen/logrus"
	jww "github.com/spf13/jwalterweatherman"
)

func main() {
	/*
		var addr = flag.String("addr", "0.0.0.0:0", "IP Address and port to use")
		var bootstrapAddress = flag.String("bootstrap", "", "IP Address and port to bootstrap against")
		var help = flag.Bool("help", false, "Display Help")
		var stun = flag.Bool("stun", true, "Use STUN")
		var metrics = flag.Bool("metrics", false, "Prometheus metrics")
		var transportType = flag.String("transport", "UTP", "Network transport protocol UTP or KCP")

		flag.Parse()

		if *help {
			displayCLIHelp()
			os.Exit(0)
		}
	*/
	jww.SetStdoutThreshold(jww.LevelTrace)
	jww.SetLogOutput(log.StandardLogger().Out)

	cfgHolder := configuration.NewHolder()
	cfgHolder.Load()
	cfgHolder.Configuration.Host.Transport.BehindNAT = false

	network, err := StartNetwork(cfgHolder.Configuration)
	if err != nil {
		log.Fatalln("Failed to start network: ", err.Error())
	}

	defer func() {
		network.closeNetwork()
	}()
	handleSignals()

	handleStats(cfgHolder.Configuration.Stats, network)

	//Println("Running interactive mode:")
	//repl(dhtNetwork, ctx)
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			//closeNetwork()
		}
	}()
}

/*
func displayCLIHelp() {
	fmt.Println(`example

Usage:
	example --addr [addr]

Options:
	--help Show this screen.
	--addr=<ip> Local IP and Port [default: 0.0.0.0]
	--bootstrap=<ip> Bootstrap IP and Port
	--stun=<bool> Use STUN protocol for public addr discovery [default: true]
    --relay=<ip> send relay request
	--metrics=<bool> [default: false]
	--transport=<protocol> Network transport protocol UTP or KCP [default: UTP]`)
}
*/
