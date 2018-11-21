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
	"fmt"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/insolar/insolar/core"
)

// Repl is "read-eval-print loop" interactive console
type Repl struct {
	NodeNetwork core.NodeNetwork
	Manager     core.PulseManager
}

// Start method starts interactive console
func (r *Repl) Start(ctx context.Context) {
	displayInteractiveHelp()

	doInfo(r.NodeNetwork)

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
			return
		case "help":
			displayInteractiveHelp()
		case "findhost":
			doFindHost(input, r.NodeNetwork)
		case "info":
			doInfo(r.NodeNetwork)
		case "activenodes":
			doActiveNodes(r.NodeNetwork)
		case "pulse":
			doPulse(ctx, r.Manager)
		default:
			displayInteractiveHelp()
		}
	}
}

func doPulse(ctx context.Context, pm core.PulseManager) {
	pulse, err := pm.Current(ctx)
	if err != nil {
		fmt.Println("Failed to get pulse")
	} else {
		fmt.Printf("Current pulse number: %d \n", pulse.PulseNumber)
	}
}

func doActiveNodes(network core.NodeNetwork) {
	nodes := network.GetActiveNodes()
	fmt.Println("Active nodes:")
	for _, n := range nodes {
		fmt.Println(n.ID().String())
	}
}

func doFindHost(input []string, network core.NodeNetwork) {
	if len(input) != 2 {
		displayInteractiveHelp()
		return
	}
	fmt.Println("Searching for NodeID:", input[1])
	nodeID := core.NewRefFromBase58(input[1])
	node := network.GetActiveNode(nodeID)
	if node != nil {
		fmt.Println("..Found targetHost:", node.PhysicalAddress())
	} else {
		fmt.Println("..Nothing found for this id!")
	}
}

func doInfo(network core.NodeNetwork) {
	hosts := len(network.GetActiveNodes())
	fmt.Println("======= Host info ======")
	fmt.Println("ID: " + network.GetOrigin().ID().String())
	fmt.Println("Known hosts: " + strconv.Itoa(hosts))
	fmt.Println("Address: " + network.GetOrigin().PhysicalAddress())
}

func displayInteractiveHelp() {
	fmt.Println(`
help - This message
findhost <key> - Find node's real network address
info - Display information about this node
activenodes - Shows active node list for current pulse
pulse - Shows current pulse number
exit - Exit programm`)
}
