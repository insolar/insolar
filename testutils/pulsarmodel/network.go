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
	"math/rand"
	"time"
)

var network [10000]*NetworkNode
var commonBus = make(chan *SignalStatistic, 10000000)

type NetworkNode struct {
	Neighbours      map[int]struct{}
	Channel         chan *Signal
	PreviousSignals map[int]struct{}
}

func (node *NetworkNode) Listen(current int) {
	for signal := range node.Channel {
		delay := time.Duration(rand.Intn(100)) * time.Millisecond
		time.Sleep(delay)
		commonBus <- &SignalStatistic{from: signal.from, to: current, id: signal.id, time: time.Now()}

		if _, ok := node.PreviousSignals[signal.id]; !ok {
			node.PreviousSignals[signal.id] = struct{}{}
			for neighbour := range node.Neighbours {
				if neighbour != signal.from {
					network[neighbour].Channel <- &Signal{from: current, id: signal.id}

				}
			}
		}
	}
}

func pulse() {
	id := 0
	pulseFunc := func() {
		id++
		nextPulsar := rand.Intn(len(network))
		network[nextPulsar].Channel <- &Signal{from: -1, id: id}
	}
	pulseFunc()

	ticker := time.NewTicker(time.Second * 10)
	go func() {
		for range ticker.C {
			pulseFunc()
		}
	}()
}
