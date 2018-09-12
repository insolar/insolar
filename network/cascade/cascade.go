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
	"encoding/binary"
	"golang.org/x/crypto/sha3"
	"math"
	"sort"
)

type SendData struct {
	NodeIds           []string
	Entropy           uint64
	ReplicationFactor uint
}

func min(a, b int) int {
	if a >= b {
		return b
	}
	return a
}

// a - scale factor
// r - common ratio
// n - length of progression
func geometricProgressionSum(a int, r int, n int) int {
	S := int(math.Pow(float64(r), float64(n)))
	return a * (1 - S) / (1 - r)
}

func calcHash(nodeID string, entropy uint64) []byte {
	data := []byte(nodeID)

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, entropy)
	for i, d := range data {
		data[i] = b[i%8] ^ d
	}

	hash := sha3.New224()
	hash.Write(data)
	return hash.Sum(nil)
}

func getNextCascadeLayerIndexes(nodeIds []string, currentNode string, replicationFactor uint) (startIndex, endIndex int) {
	depth := 0
	j := 0
	layerWidth := replicationFactor
	found := false
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

	n := int(replicationFactor)
	var layerWeight int
	if n == 1 {
		layerWeight = depth + 1
	} else {
		layerWeight = geometricProgressionSum(n, n, depth+1)
	}
	startIndex = layerWeight + j*n
	endIndex = startIndex + n
	return
}

func CalculateNextNodes(data SendData, currentNode string) (nextNodeIds []string) {
	nodeIds := make([]string, len(data.NodeIds))
	copy(nodeIds, data.NodeIds)

	sort.SliceStable(nodeIds, func(i, j int) bool {
		return bytes.Compare(
			calcHash(nodeIds[i], data.Entropy),
			calcHash(nodeIds[j], data.Entropy)) < 0
	})

	if currentNode == "" {
		l := min(len(nodeIds), int(data.ReplicationFactor))
		return nodeIds[:l]
	}

	startIndex, endIndex := getNextCascadeLayerIndexes(nodeIds, currentNode, data.ReplicationFactor)

	if startIndex >= len(nodeIds) {
		return nil
	}
	return nodeIds[startIndex:min(endIndex, len(nodeIds))]
}
