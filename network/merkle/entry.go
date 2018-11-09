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

package merkle

import (
	"bytes"
	"sort"

	"github.com/insolar/insolar/core"
)

type PulseEntry struct {
	Pulse *core.Pulse
}

func (pe *PulseEntry) hash() []byte {
	return pulseHash(pe.Pulse)
}

type GlobuleEntry struct {
	PulseEntry
	ProofSet      map[core.Node]*PulseProof
	PulseHash     []byte
	PrevCloudHash []byte
	GlobuleIndex  uint32
}

func (ge *GlobuleEntry) hash() []byte {
	nodeEntryByRole := nodeEntryByRole(ge.ProofSet)
	var bucketHashes [][]byte

	for role, roleEntries := range nodeEntryByRole {
		sortEntries(roleEntries)

		bucketEntryRoot := roleEntryRoot(roleEntries)
		bucketInfoHash := bucketInfoHash(role, uint32(len(roleEntries)))
		bucketHash := bucketHash(bucketInfoHash, bucketEntryRoot)

		bucketHashes = append(bucketHashes, bucketHash)
	}

	return fromList(bucketHashes).MerkleRoot()
}

type CloudEntry struct {
	ProofSet      []*GlobuleProof
	PrevCloudHash []byte
	// TODO: implement later
	// ProofSet map[core.Globule]*GlobuleProof
}

func (ce *CloudEntry) hash() []byte {
	var result [][]byte

	for _, proof := range ce.ProofSet {
		globuleInfoHash := globuleInfoHash(ce.PrevCloudHash, proof.GlobuleIndex, proof.NodeCount)
		globuleHash := globuleHash(globuleInfoHash, proof.NodeRoot)
		result = append(result, globuleHash)
	}

	mt := fromList(result)
	return mt.MerkleRoot()
}

type nodeEntry struct {
	PulseEntry
	PulseProof *PulseProof
	Node       core.Node
}

func (ne *nodeEntry) hash() []byte {
	pulseHash := ne.PulseEntry.hash()
	nodeInfoHash := nodeInfoHash(pulseHash, ne.PulseProof.StateHash)
	return nodeHash(ne.PulseProof.Signature, nodeInfoHash)
}

func nodeEntryByRole(nodeProofs map[core.Node]*PulseProof) map[core.NodeRole][]*nodeEntry {
	roleMap := make(map[core.NodeRole][]*nodeEntry)
	for node, pulseProof := range nodeProofs {
		role := node.Role()
		roleMap[role] = append(roleMap[role], &nodeEntry{
			Node:       node,
			PulseProof: pulseProof,
		})
	}
	return roleMap
}

func sortEntries(roleEntries []*nodeEntry) {
	sort.SliceStable(roleEntries, func(i, j int) bool {
		return bytes.Compare(
			roleEntries[i].Node.ID().Bytes(),
			roleEntries[j].Node.ID().Bytes()) < 0
	})
}

func roleEntryRoot(roleEntries []*nodeEntry) []byte {
	var roleEntriesHashes [][]byte
	for index, entry := range roleEntries {
		bucketEntryHash := bucketEntryHash(uint32(index), entry.hash())
		roleEntriesHashes = append(roleEntriesHashes, bucketEntryHash)
	}
	return fromList(roleEntriesHashes).MerkleRoot()
}
