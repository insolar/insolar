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
	"sort"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

type PulseEntry struct {
	Pulse *core.Pulse
}

func (pe *PulseEntry) hash(helper *merkleHelper) []byte {
	return helper.pulseHash(pe.Pulse)
}

type GlobuleEntry struct {
	*PulseEntry
	ProofSet      map[core.Node]*PulseProof
	PulseHash     []byte
	PrevCloudHash []byte
	GlobuleID     core.GlobuleID
}

func (ge *GlobuleEntry) hash(helper *merkleHelper) ([]byte, error) {
	nodeEntryByRole := nodeEntryByRole(ge.PulseEntry, ge.ProofSet)
	var bucketHashes [][]byte

	for _, role := range core.AllStaticRoles {
		roleEntries, ok := nodeEntryByRole[role]
		if !ok {
			continue
		}

		sortEntries(roleEntries)
		bucketEntryRoot, err := roleEntryRoot(roleEntries, helper)

		if err != nil {
			return nil, errors.Wrap(err, "[ hash ] Failed to create tree for bucket role entry")
		}

		bucketInfoHash := helper.bucketInfoHash(role, uint32(len(roleEntries)))
		bucketHash := helper.bucketHash(bucketInfoHash, bucketEntryRoot)
		bucketHashes = append(bucketHashes, bucketHash)
	}

	tree, err := treeFromHashList(bucketHashes, helper.scheme.IntegrityHasher())

	if err != nil {
		return nil, errors.Wrap(err, "[ hash ] Failed to create tree for bucket hashes")
	}

	return tree.Root(), nil
}

type CloudEntry struct {
	ProofSet      []*GlobuleProof
	PrevCloudHash []byte
}

func (ce *CloudEntry) hash(helper *merkleHelper) ([]byte, error) {
	var result [][]byte

	for _, proof := range ce.ProofSet {
		globuleInfoHash := helper.globuleInfoHash(ce.PrevCloudHash, uint32(proof.GlobuleID), proof.NodeCount)
		globuleHash := helper.globuleHash(globuleInfoHash, proof.NodeRoot)
		result = append(result, globuleHash)
	}

	tree, err := treeFromHashList(result, helper.scheme.IntegrityHasher())
	if err != nil {
		return nil, errors.Wrap(err, "[ hash ] Failed to create tree")
	}

	return tree.Root(), nil
}

type nodeEntry struct {
	*PulseEntry
	PulseProof *PulseProof
	Node       core.Node
}

func (ne *nodeEntry) hash(helper *merkleHelper) []byte {
	pulseHash := ne.PulseEntry.hash(helper)
	nodeInfoHash := helper.nodeInfoHash(pulseHash, ne.PulseProof.StateHash)
	return helper.nodeHash(ne.PulseProof.Signature.Bytes(), nodeInfoHash)
}

func nodeEntryByRole(pulseEntry *PulseEntry, nodeProofs map[core.Node]*PulseProof) map[core.StaticRole][]*nodeEntry {
	roleMap := make(map[core.StaticRole][]*nodeEntry)
	for node, pulseProof := range nodeProofs {
		role := node.Role()
		roleMap[role] = append(roleMap[role], &nodeEntry{
			PulseEntry: pulseEntry,
			Node:       node,
			PulseProof: pulseProof,
		})
	}
	return roleMap
}

func sortEntries(roleEntries []*nodeEntry) {
	sort.SliceStable(roleEntries, func(i, j int) bool {
		return roleEntries[i].Node.ID().Compare(roleEntries[j].Node.ID()) < 0
	})
}

func roleEntryRoot(roleEntries []*nodeEntry, helper *merkleHelper) ([]byte, error) {
	var roleEntriesHashes [][]byte
	for index, entry := range roleEntries {
		bucketEntryHash := helper.bucketEntryHash(uint32(index), entry.hash(helper))
		roleEntriesHashes = append(roleEntriesHashes, bucketEntryHash)
	}

	tree, err := treeFromHashList(roleEntriesHashes, helper.scheme.IntegrityHasher())
	if err != nil {
		return nil, errors.Wrap(err, "[ hash ] Failed to create tree")
	}

	return tree.Root(), nil
}
