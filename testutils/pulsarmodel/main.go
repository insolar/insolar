package main

import (
	"fmt"
	"math/rand"
)

func main() {
	rand.Seed(31)

	fmt.Println("---------- Start to init -----------")
	for networkNodeIndex := 0; networkNodeIndex < len(network); networkNodeIndex++ {
		network[networkNodeIndex] = &NetworkNode{Neighbours: map[int]struct{}{}, PreviousSignals: map[int]struct{}{}}
		network[networkNodeIndex].Channel = make(chan *Signal, 100000)
		go network[networkNodeIndex].Listen(networkNodeIndex)
	}
	fmt.Println("---------- finish to init-----------")

	fmt.Println("---------- Start to configure -----------")
	for networkNodeIndex := 0; networkNodeIndex < len(network); networkNodeIndex++ {
		for len(network[networkNodeIndex].Neighbours) < 300 {
			nextEdge := rand.Intn(len(network))

			_, ok := network[networkNodeIndex].Neighbours[nextEdge]
			if !ok {
				network[networkNodeIndex].Neighbours[nextEdge] = struct{}{}
				network[nextEdge].Neighbours[nextEdge] = struct{}{}
			}
		}
	}
	fmt.Println("---------- finish to configure -----------")

	pulse()
	calculateLoad()
}
