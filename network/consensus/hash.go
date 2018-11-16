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

func calculateNodeHash(scheme core.PlatformCryptographyScheme, node core.Node) []byte {
	h := scheme.IntegrityHasher()
	hashWriteChecked(h, node.ID().Bytes())
	b := make([]byte, 8)
	nodeRoles := make([]core.NodeRole, len(node.Roles()))
	copy(nodeRoles, node.Roles())
	sort.Slice(nodeRoles[:], func(i, j int) bool {
		return nodeRoles[i] < nodeRoles[j]
	})
	for _, nodeRole := range nodeRoles {
		binary.LittleEndian.PutUint32(b, uint32(nodeRole))
		hashWriteChecked(h, b[:4])
	}
	hashWriteChecked(h, b[:])
	binary.LittleEndian.PutUint32(b, uint32(node.Pulse()))
	hashWriteChecked(h, b[:4])
	// TODO: pass correctly public key to active node
	// publicKey, err := ecdsa.ExportPublicKey(node.PublicKey)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// hashWriteChecked(h, []byte(publicKey))
	hashWriteChecked(h, []byte(node.PhysicalAddress()))
	hashWriteChecked(h, []byte(node.Version()))
	return h.Sum(nil)
}

// CalculateHash calculates hash of active node list
func CalculateHash(scheme core.PlatformCryptographyScheme, list []core.Node) (result []byte, err error) {
	sort.Slice(list[:], func(i, j int) bool {
		return bytes.Compare(list[i].ID().Bytes(), list[j].ID().Bytes()) < 0
	})

	// catch possible panic from hashWriteChecked in this function and in all calculateNodeHash funcs
	defer func() {
		if r := recover(); r != nil {
			result, err = nil, fmt.Errorf("error calculating h: %s", r)
		}
	}()

	h := scheme.IntegrityHasher()
	for _, node := range list {
		nodeHash := calculateNodeHash(scheme, node)
		hashWriteChecked(h, nodeHash)
	}
	return h.Sum(nil), nil
}

// CalculateNodeUnsyncHash calculates hash for a NodeUnsyncHash
func CalculateNodeUnsyncHash(scheme core.PlatformCryptographyScheme, nodeID core.RecordRef, list []core.Node) (*NodeUnsyncHash, error) {
	h, err := CalculateHash(scheme, list)
	if err != nil {
		return nil, err
	}
	return &NodeUnsyncHash{NodeID: nodeID, Hash: h}, nil
}
