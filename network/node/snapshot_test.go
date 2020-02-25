// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package node

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulse"

	"github.com/stretchr/testify/assert"
)

func NewMutator(snapshot *Snapshot) *Mutator {
	return &Mutator{Accessor: NewAccessor(snapshot)}
}

type Mutator struct {
	*Accessor
}

func (m *Mutator) AddWorkingNode(n insolar.NetworkNode) {
	if _, ok := m.refIndex[n.ID()]; ok {
		return
	}
	mutableNode := n.(MutableNode)
	m.addToIndex(mutableNode)
	listType := nodeStateToListType(mutableNode)
	m.snapshot.nodeList[listType] = append(m.snapshot.nodeList[listType], n)
	m.active = append(m.active, n)
}

func TestSnapshotEncodeDecode(t *testing.T) {

	ks := platformpolicy.NewKeyProcessor()
	p1, err := ks.GeneratePrivateKey()
	p2, err := ks.GeneratePrivateKey()
	assert.NoError(t, err)

	n1 := newMutableNode(gen.Reference(), insolar.StaticRoleVirtual, ks.ExtractPublicKey(p1), insolar.NodeReady, "127.0.0.1:22", "ver2")
	n2 := newMutableNode(gen.Reference(), insolar.StaticRoleHeavyMaterial, ks.ExtractPublicKey(p2), insolar.NodeLeaving, "127.0.0.1:33", "ver5")

	s := Snapshot{}
	s.pulse = 22
	s.state = insolar.CompleteNetworkState
	s.nodeList[ListLeaving] = []insolar.NetworkNode{n1, n2}
	s.nodeList[ListJoiner] = []insolar.NetworkNode{n2}

	buff, err := s.Encode()
	assert.NoError(t, err)
	assert.NotEmptyf(t, buff, "should not be empty")

	s2 := Snapshot{}
	err = s2.Decode(buff)
	assert.NoError(t, err)
	assert.True(t, s.Equal(&s2))
}

func TestSnapshot_Decode(t *testing.T) {
	buff := []byte("hohoho i'm so broken")
	s := Snapshot{}
	assert.Error(t, s.Decode(buff))
}

func TestSnapshot_Copy(t *testing.T) {
	snapshot := NewSnapshot(pulse.MinTimePulse, nil)
	mutator := NewMutator(snapshot)
	ref1 := gen.Reference()
	node1 := newMutableNode(ref1, insolar.StaticRoleVirtual, nil, insolar.NodeReady, "127.0.0.1:0", "")
	mutator.AddWorkingNode(node1)

	snapshot2 := snapshot.Copy()
	accessor := NewAccessor(snapshot2)

	ref2 := gen.Reference()
	node2 := newMutableNode(ref2, insolar.StaticRoleLightMaterial, nil, insolar.NodeReady, "127.0.0.1:0", "")
	mutator.AddWorkingNode(node2)

	// mutator and accessor observe different copies of snapshot and don't affect each other
	assert.Equal(t, 2, len(mutator.GetActiveNodes()))
	assert.Equal(t, 1, len(accessor.GetActiveNodes()))
	assert.False(t, snapshot.Equal(snapshot2))
}

func TestSnapshot_GetPulse(t *testing.T) {
	snapshot := NewSnapshot(10, nil)
	assert.EqualValues(t, 10, snapshot.GetPulse())
	snapshot = NewSnapshot(152, nil)
	assert.EqualValues(t, 152, snapshot.GetPulse())
}

func TestSnapshot_Equal(t *testing.T) {
	snapshot := NewSnapshot(10, nil)
	snapshot2 := NewSnapshot(11, nil)
	assert.False(t, snapshot.Equal(snapshot2))

	snapshot2.pulse = insolar.PulseNumber(10)

	genNodeCopy := func(reference insolar.Reference) insolar.NetworkNode {
		return newMutableNode(reference, insolar.StaticRoleLightMaterial,
			nil, insolar.NodeReady, "127.0.0.1:0", "")
	}

	refs := gen.UniqueReferences(2)
	node1 := genNodeCopy(refs[0])
	node2 := genNodeCopy(refs[1])
	node3 := genNodeCopy(refs[1])

	mutator := NewMutator(snapshot)
	mutator.AddWorkingNode(node1)
	mutator.AddWorkingNode(node2)

	mutator2 := NewMutator(snapshot2)
	mutator2.AddWorkingNode(node1)
	mutator2.AddWorkingNode(node3)

	assert.True(t, snapshot.Equal(snapshot2))
}
