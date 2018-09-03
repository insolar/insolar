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
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/insolar/insolar/network/hostnetwork"
)

func repl(dhtNetwork *hostnetwork.DHT, ctx hostnetwork.Context) {
	displayInteractiveHelp()
	doInfo(dhtNetwork, ctx)

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
		case "exit":
			fallthrough
		case "quit":
			os.Exit(0)
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
	originID := dhtNetwork.GetOriginHost(ctx).ID
	fmt.Println("======= Host info ======")
	fmt.Println("ID key: " + originID.KeyString())
	fmt.Println("ID hash: " + originID.HashString())
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

func displayInteractiveHelp() {
	fmt.Println(`
help - This message
findnode <key> - Find node's real network address
info - Display information about this node
exit - Exit programm

<method> <target> <args...> - Remote procedure call`)
}
