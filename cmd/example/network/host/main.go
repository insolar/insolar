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
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/hostnetwork/host"

	"github.com/chzyer/readline"
)

var dhtNetwork *hostnetwork.DHT

func main() {
	var addr = flag.String("addr", "0.0.0.0:0", "IP Address and port to use")
	var bootstrapAddress = flag.String("bootstrap", "", "IP Address and port to bootstrap against")
	var help = flag.Bool("help", false, "Display Help")
	var stun = flag.Bool("stun", true, "Use STUN")

	flag.Parse()

	if *help {
		displayCLIHelp()
		os.Exit(0)
	}

	cfg := configuration.NewConfiguration()
	cfg.Host.Transport.Address = *addr
	cfg.Host.Transport.BehindNAT = *stun

	dhtNetwork, err := hostnetwork.NewHostNetwork(cfg.Host)
	if err != nil {
		log.Fatalln("Failed to create network:", err.Error())
	}

	defer closeNetwork()

	ctx := createContext(dhtNetwork)

	go listen(dhtNetwork)

	if len(*bootstrapAddress) > 0 {
		bootstrap(dhtNetwork)
	}

	handleSignals()

	err = dhtNetwork.ObtainIP(ctx)
	if err != nil {
		log.Println(err)
	}
	err = dhtNetwork.AnalyzeNetwork(ctx)
	if err != nil {
		log.Println(err)
	}
	repl(dhtNetwork, ctx)
}

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			closeNetwork()
		}
	}()
}

func createContext(dhtNetwork *hostnetwork.DHT) hostnetwork.Context {
	ctx, err := hostnetwork.NewContextBuilder(dhtNetwork).SetDefaultHost().Build()
	if err != nil {
		log.Fatalln("Failed to create context:", err.Error())
	}
	return ctx
}

func bootstrap(dhtNetwork *hostnetwork.DHT) {
	err := dhtNetwork.Bootstrap()
	if err != nil {
		log.Fatalln("Failed to bootstrap network", err.Error())
	}
}

func listen(dhtNetwork *hostnetwork.DHT) {
	func() {
		err := dhtNetwork.Listen()
		if err != nil {
			log.Fatalln("Listen failed:", err.Error())
		}
	}()
}

func closeNetwork() {
	func() {
		dhtNetwork.Disconnect()
	}()
}

func repl(dhtNetwork *hostnetwork.DHT, ctx hostnetwork.Context) {
	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer func() {
		errRlClose := rl.Close()
		if errRlClose != nil {
			panic(errRlClose)
		}
	}()
	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
		}
		input := strings.Split(line, " ")

		switch input[0] {
		case "help":
			displayInteractiveHelp()
		case "findhost":
			doFindHost(input, dhtNetwork, ctx)
		case "info":
			doInfo(dhtNetwork, ctx)
		case "relay":
			doSendRelay(input[2], input[1], dhtNetwork, ctx)
		default:
			doRPC(input, dhtNetwork, ctx)
		}
	}
}

func doFindHost(input []string, dhtNetwork *hostnetwork.DHT, ctx hostnetwork.Context) {
	if len(input) != 2 {
		displayInteractiveHelp()
		return
	}
	fmt.Println("Searching for targetHost", input[1])
	targetHost, exists, err := dhtNetwork.FindHost(ctx, input[1])
	if err != nil {
		fmt.Println(err.Error())
	}
	if exists {
		fmt.Println("..Found targetHost:", targetHost)
	} else {
		fmt.Println("..Nothing found for this id!")
	}
}

func doInfo(dhtNetwork *hostnetwork.DHT, ctx hostnetwork.Context) {
	hosts := dhtNetwork.NumHosts(ctx)
	originID := dhtNetwork.GetOriginID(ctx)
	fmt.Println("ID: " + originID.HashString())
	fmt.Println("Known hosts: " + strconv.Itoa(hosts))
}

func doSendRelay(command, relayAddr string, dhtNetwork *hostnetwork.DHT, ctx hostnetwork.Context) {
	err := dhtNetwork.RelayRequest(ctx, command, relayAddr)
	if err != nil {
		log.Println(err)
	}
}

func doRPC(input []string, dhtNetwork *hostnetwork.DHT, ctx hostnetwork.Context) {
	if len(input) < 2 || len(input[0]) == 0 || len(input[1]) == 0 {
		if len(input) > 0 && len(input[0]) > 0 {
			displayInteractiveHelp()
		}
		return
	}

	method, target := input[0], input[1]
	args := make([][]byte, 0, 4)
	for _, arg := range input[2:] {
		args = append(args, []byte(arg))
	}

	fmt.Printf("Running remote method %s on %s with args %v \n", method, target, args)

	result, err := dhtNetwork.RemoteProcedureCall(ctx, target, method, args)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(string(result))
	}
}

func displayCLIHelp() {
	fmt.Println(`example

Usage:
	example --addr [addr]

Options:
	--help Show this screen.
	--addr=<ip> Local IP and Port [default: 0.0.0.0]
	--bootstrap=<ip> Bootstrap IP and Port
	--stun=<bool> Use STUN protocol for public addr discovery [default: true]
    --relay=<ip> send relay request`)
}

func displayInteractiveHelp() {
	fmt.Println(`
help - This packet
findhost <key> - Find host's real network address
info - Display information about this host

<method> <target> <args...> - Remote procedure call`)
}

func send(sender *host.Host, args [][]byte) ([]byte, error) {
	bs := append([]byte{}, []byte(time.Now().Format(time.Kitchen))...)
	bs = append(bs, ' ')
	bs = append(bs, sender.ID.HashString()...)

	for _, item := range args {
		bs = append(bs, ' ')
		bs = append(bs, item...)
	}

	fmt.Println(string(bs))

	return bs, nil
}
