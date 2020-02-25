// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package node

import (
	"sort"

	"github.com/insolar/insolar/insolar"
)

type Accessor struct {
	snapshot  *Snapshot
	refIndex  map[insolar.Reference]insolar.NetworkNode
	sidIndex  map[insolar.ShortNodeID]insolar.NetworkNode
	addrIndex map[string]insolar.NetworkNode
	// should be removed in future
	active []insolar.NetworkNode
}

func (a *Accessor) GetActiveNodeByShortID(shortID insolar.ShortNodeID) insolar.NetworkNode {
	return a.sidIndex[shortID]
}

func (a *Accessor) GetActiveNodeByAddr(address string) insolar.NetworkNode {
	return a.addrIndex[address]
}

func (a *Accessor) GetActiveNodes() []insolar.NetworkNode {
	result := make([]insolar.NetworkNode, len(a.active))
	copy(result, a.active)
	return result
}

func (a *Accessor) GetActiveNode(ref insolar.Reference) insolar.NetworkNode {
	return a.refIndex[ref]
}

func (a *Accessor) GetWorkingNode(ref insolar.Reference) insolar.NetworkNode {
	node := a.GetActiveNode(ref)
	if node == nil || node.GetPower() == 0 {
		return nil
	}
	return node
}

func (a *Accessor) GetWorkingNodes() []insolar.NetworkNode {
	workingList := a.snapshot.nodeList[ListWorking]
	result := make([]insolar.NetworkNode, len(workingList))
	copy(result, workingList)
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID().Compare(result[j].ID()) < 0
	})
	return result
}

func GetSnapshotActiveNodes(snapshot *Snapshot) []insolar.NetworkNode {
	joining := snapshot.nodeList[ListJoiner]
	idle := snapshot.nodeList[ListIdle]
	working := snapshot.nodeList[ListWorking]
	leaving := snapshot.nodeList[ListLeaving]

	joinersCount := len(joining)
	idlersCount := len(idle)
	workingCount := len(working)
	leavingCount := len(leaving)

	result := make([]insolar.NetworkNode, joinersCount+idlersCount+workingCount+leavingCount)

	copy(result[:joinersCount], joining)
	copy(result[joinersCount:joinersCount+idlersCount], idle)
	copy(result[joinersCount+idlersCount:joinersCount+idlersCount+workingCount], working)
	copy(result[joinersCount+idlersCount+workingCount:], leaving)

	return result
}

func (a *Accessor) addToIndex(node insolar.NetworkNode) {
	a.refIndex[node.ID()] = node
	a.sidIndex[node.ShortID()] = node
	a.addrIndex[node.Address()] = node

	if node.GetPower() == 0 {
		return
	}
}

func NewAccessor(snapshot *Snapshot) *Accessor {
	result := &Accessor{
		snapshot:  snapshot,
		refIndex:  make(map[insolar.Reference]insolar.NetworkNode),
		sidIndex:  make(map[insolar.ShortNodeID]insolar.NetworkNode),
		addrIndex: make(map[string]insolar.NetworkNode),
	}
	result.active = GetSnapshotActiveNodes(snapshot)
	for _, node := range result.active {
		result.addToIndex(node)
	}
	return result
}
