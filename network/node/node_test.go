// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/network/LICENSE.md.

package node

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/stretchr/testify/assert"
)

func TestNode_Version(t *testing.T) {
	n := NewNode(gen.Reference(), insolar.StaticRoleVirtual, nil, "127.0.0.1", "123")
	assert.Equal(t, "123", n.Version())
	n.(MutableNode).SetVersion("234")
	assert.Equal(t, "234", n.Version())
}

func TestNode_GetState(t *testing.T) {
	n := NewNode(gen.Reference(), insolar.StaticRoleVirtual, nil, "127.0.0.1", "123")
	assert.Equal(t, insolar.NodeReady, n.GetState())
	n.(MutableNode).SetState(insolar.NodeUndefined)
	assert.Equal(t, insolar.NodeUndefined, n.GetState())
	n.(MutableNode).ChangeState()
	assert.Equal(t, insolar.NodeJoining, n.GetState())
	n.(MutableNode).ChangeState()
	assert.Equal(t, insolar.NodeReady, n.GetState())
	n.(MutableNode).ChangeState()
	assert.Equal(t, insolar.NodeReady, n.GetState())
}

func TestNode_GetGlobuleID(t *testing.T) {
	n := NewNode(gen.Reference(), insolar.StaticRoleVirtual, nil, "127.0.0.1", "123")
	assert.EqualValues(t, 0, n.GetGlobuleID())
}

func TestNode_LeavingETA(t *testing.T) {
	n := NewNode(gen.Reference(), insolar.StaticRoleVirtual, nil, "127.0.0.1", "123")
	assert.Equal(t, insolar.NodeReady, n.GetState())
	n.(MutableNode).SetLeavingETA(25)
	assert.Equal(t, insolar.NodeLeaving, n.GetState())
	assert.EqualValues(t, 25, n.LeavingETA())
}

func TestNode_ShortID(t *testing.T) {
	n := NewNode(gen.Reference(), insolar.StaticRoleVirtual, nil, "127.0.0.1", "123")
	assert.EqualValues(t, GenerateUintShortID(n.ID()), n.ShortID())
	n.(MutableNode).SetShortID(11)
	assert.EqualValues(t, 11, n.ShortID())
}
