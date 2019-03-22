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
