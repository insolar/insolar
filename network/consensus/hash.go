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

package consensus

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash"
	"sort"

	"github.com/insolar/insolar/core"
	"golang.org/x/crypto/sha3"
)

func hashWriteChecked(hash hash.Hash, data []byte) {
	n, err := hash.Write(data)
	if n != len(data) {
		panic(fmt.Sprintf("Error writing hash. Bytes expected: %d; bytes actual: %d", len(data), n))
	}
	if err != nil {
		panic(err.Error())
	}
}

func calculateNodeHash(node *core.ActiveNode) []byte {
	hash := sha3.New224()
	hashWriteChecked(hash, node.NodeID[:])
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(node.Role))
	hashWriteChecked(hash, b[:])
	binary.LittleEndian.PutUint32(b, uint32(node.PulseNum))
	hashWriteChecked(hash, b[:4])
	b[0] = byte(node.State)
	hashWriteChecked(hash, b[:1])
	// hashWriteChecked(hash, node.PublicKey)
	return hash.Sum(nil)
}

// CalculateHash calculates hash of active node list
func CalculateHash(list []*core.ActiveNode) (result []byte, err error) {
	sort.Slice(list[:], func(i, j int) bool {
		return bytes.Compare(list[i].NodeID[:], list[j].NodeID[:]) < 0
	})

	// catch possible panic from hashWriteChecked in this function and in all calculateNodeHash funcs
	defer func() {
		if r := recover(); r != nil {
			result, err = nil, fmt.Errorf("error calculating hash: %s", r)
		}
	}()

	hash := sha3.New224()
	for _, node := range list {
		nodeHash := calculateNodeHash(node)
		hashWriteChecked(hash, nodeHash)
	}
	return hash.Sum(nil), nil
}

// CalculateNodeUnsyncHash calculates hash for a NodeUnsyncHash
func CalculateNodeUnsyncHash(nodeID core.RecordRef, list []*core.ActiveNode) (*NodeUnsyncHash, error) {
	hash, err := CalculateHash(list)
	if err != nil {
		return nil, err
	}
	return &NodeUnsyncHash{nodeID, hash}, nil
}
