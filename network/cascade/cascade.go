/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification,
 * are permitted (subject to the limitations in the disclaimer below) provided that
 * the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *  * Neither the name of Insolar Technologies nor the names of its contributors
 *    may be used to endorse or promote products derived from this software
 *    without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
 * BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
 * CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING,
 * BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
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
