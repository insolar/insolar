/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package nodenetwork

import (
	"github.com/insolar/insolar/core"
)

type Accessor struct {
	snapshot  *Snapshot
	refIndex  map[core.RecordRef]core.Node
	sidIndex  map[core.ShortNodeID]core.Node
	roleIndex map[core.StaticRole]*recordRefSet
	// should be removed in future
	active []core.Node
}

func (a *Accessor) GetActiveNodeByShortID(shortID core.ShortNodeID) core.Node {
	return a.sidIndex[shortID]
}

func (a *Accessor) GetActiveNodes() []core.Node {
	return a.active
}

func (a *Accessor) GetActiveNode(ref core.RecordRef) core.Node {
	return a.refIndex[ref]
}

func (a *Accessor) GetWorkingNode(ref core.RecordRef) core.Node {
	node := a.GetActiveNode(ref)
	if node == nil || node.GetState() != core.NodeReady {
		return nil
	}
	return node
}

func (a *Accessor) GetWorkingNodes() []core.Node {
	return a.snapshot.nodeList[ListWorking]
}

func (a *Accessor) GetWorkingNodesByRole(role core.DynamicRole) []core.RecordRef {
	staticRole := jetRoleToNodeRole(role)
	return a.roleIndex[staticRole].Collect()
}

func GetSnapshotActiveNodes(snapshot *Snapshot) []core.Node {
	joining := snapshot.nodeList[ListJoiner]
	working := snapshot.nodeList[ListWorking]
	leaving := snapshot.nodeList[ListLeaving]

	result := make([]core.Node, len(joining)+len(working)+len(leaving))
	copy(result[:len(joining)], joining[:])
	copy(result[len(joining):len(joining)+len(working)], working[:])
	copy(result[len(joining)+len(working):], leaving[:])
	return result
}

func (a *Accessor) addToRoleIndex(node core.Node) {
	if node.GetState() != core.NodeReady {
		return
	}

	list, ok := a.roleIndex[node.Role()]
	if !ok {
		list = newRecordRefSet()
	}

	list.Add(node.ID())
	a.roleIndex[node.Role()] = list
}

func NewAccessor(snapshot *Snapshot) *Accessor {
	result := &Accessor{
		snapshot:  snapshot,
		refIndex:  make(map[core.RecordRef]core.Node),
		sidIndex:  make(map[core.ShortNodeID]core.Node),
		roleIndex: make(map[core.StaticRole]*recordRefSet),
	}
	result.active = GetSnapshotActiveNodes(snapshot)
	for _, node := range result.active {
		result.refIndex[node.ID()] = node
		result.sidIndex[node.ShortID()] = node
		if node.GetState() == core.NodeReady {
			result.addToRoleIndex(node)
		}
	}
	return result
}
