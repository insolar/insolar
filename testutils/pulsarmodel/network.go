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
