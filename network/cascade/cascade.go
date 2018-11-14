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

package cascade

import (
	"bytes"
	"fmt"
	"math"
	"sort"

	"github.com/insolar/insolar/core"
)

func min(a, b int) int {
	if a >= b {
		return b
	}
	return a
}

// Cascade is struct to hold callback that sends cascade messages to next layers of cascade
type Cascade struct {
	SendMessage func(data core.Cascade, method string, args [][]byte) error
}

// SendToNextLayer sends data to callback.
func (casc *Cascade) SendToNextLayer(data core.Cascade, method string, args [][]byte) error {
	return casc.SendMessage(data, method, args)
}

// a - scale factor
// r - common ratio
// n - length of progression
func geometricProgressionSum(a int, r int, n int) int {
	S := int(math.Pow(float64(r), float64(n)))
	return a * (1 - S) / (1 - r)
}

func calcHash(scheme core.PlatformCryptographyScheme, nodeID core.RecordRef, entropy core.Entropy) []byte {
	data := make([]byte, core.RecordRefSize)
	copy(data, nodeID[:])
	for i, d := range data {
		data[i] = entropy[i%core.EntropySize] ^ d
	}

	h := scheme.IntegrityHasher()
	_, err := h.Write(data)
	if err != nil {
		panic(err)
	}
	return h.Sum(nil)
}

func getNextCascadeLayerIndexes(nodeIds []core.RecordRef, currentNode core.RecordRef, replicationFactor uint) (startIndex, endIndex int) {
	depth := 0
	j := 0
	layerWidth := replicationFactor
	found := false
	// iterate to find current node in the nodes slice, incrementing j and depth according to replicationFactor
	for _, nodeID := range nodeIds {
		if nodeID == currentNode {
			found = true
			break
		}
		j++
		if j == int(layerWidth) {
			layerWidth *= replicationFactor
			depth++
			j = 0
		}
	}

	if !found {
		return len(nodeIds), len(nodeIds)
	}

	// calculate count of the all nodes that have depth less or equal to the current node
	n := int(replicationFactor)
	var layerWeight int
	if n == 1 {
		layerWeight = depth + 1
	} else {
		layerWeight = geometricProgressionSum(n, n, depth+1)
	}
	// calculate children subtree of the current node
	startIndex = layerWeight + j*n
	endIndex = startIndex + n
	return
}

// CalculateNextNodes get nodes of the next cascade layer from the input nodes slice
func CalculateNextNodes(scheme core.PlatformCryptographyScheme, data core.Cascade, currentNode *core.RecordRef) (nextNodeIds []core.RecordRef, err error) {
	nodeIds := make([]core.RecordRef, len(data.NodeIds))
	copy(nodeIds, data.NodeIds)

	// catching possible panic from calcHash
	defer func() {
		if r := recover(); r != nil {
			nextNodeIds, err = nil, fmt.Errorf("panic: %s", r)
		}
	}()

	sort.SliceStable(nodeIds, func(i, j int) bool {
		return bytes.Compare(
			calcHash(scheme, nodeIds[i], data.Entropy),
			calcHash(scheme, nodeIds[j], data.Entropy)) < 0
	})

	if currentNode == nil {
		length := min(int(data.ReplicationFactor), len(nodeIds))
		return nodeIds[:length], nil
	}

	// get indexes of the next layer nodes from the sorted nodes slice
	startIndex, endIndex := getNextCascadeLayerIndexes(nodeIds, *currentNode, data.ReplicationFactor)

	if startIndex >= len(nodeIds) {
		return nil, nil
	}
	return nodeIds[startIndex:min(endIndex, len(nodeIds))], nil
}
