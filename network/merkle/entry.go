// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/network/LICENSE.md.

package merkle

import (
	"sort"

	"github.com/insolar/insolar/insolar"
	"github.com/pkg/errors"
)

type PulseEntry struct {
	Pulse *insolar.Pulse
}

func (pe *PulseEntry) hash(helper *merkleHelper) []byte {
	return helper.pulseHash(pe.Pulse)
}

type GlobuleEntry struct {
	*PulseEntry
	ProofSet      map[insolar.NetworkNode]*PulseProof
	PulseHash     []byte
	PrevCloudHash []byte
	GlobuleID     insolar.GlobuleID
}

func (ge *GlobuleEntry) hash(helper *merkleHelper) ([]byte, error) {
	neByRole := nodeEntryByRole(ge.PulseEntry, ge.ProofSet)
	bucketHashes := make([][]byte, 0, len(insolar.AllStaticRoles))

	for _, role := range insolar.AllStaticRoles {
		roleEntries, ok := neByRole[role]
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
	result := make([][]byte, 0, len(ce.ProofSet))

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
	Node       insolar.NetworkNode
}

func (ne *nodeEntry) hash(helper *merkleHelper) []byte {
	pulseHash := ne.PulseEntry.hash(helper)
	nodeInfoHash := helper.nodeInfoHash(pulseHash, ne.PulseProof.StateHash)
	return helper.nodeHash(ne.PulseProof.Signature.Bytes(), nodeInfoHash)
}

func nodeEntryByRole(pulseEntry *PulseEntry, nodeProofs map[insolar.NetworkNode]*PulseProof) map[insolar.StaticRole][]*nodeEntry {
	roleMap := make(map[insolar.StaticRole][]*nodeEntry)
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
	roleEntriesHashes := make([][]byte, 0, len(roleEntries))
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
