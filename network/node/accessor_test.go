package node

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/pulse"

	"github.com/stretchr/testify/assert"
)

func TestAccessor(t *testing.T) {
	t.Skip("FIXME")

	node := newMutableNode(gen.Reference(), insolar.StaticRoleVirtual, nil, insolar.NodeReady, "127.0.0.1:0", "")

	node2 := newMutableNode(gen.Reference(), insolar.StaticRoleVirtual, nil, insolar.NodeJoining, "127.0.0.1:0", "")
	node2.SetShortID(11)

	node3 := newMutableNode(gen.Reference(), insolar.StaticRoleVirtual, nil, insolar.NodeLeaving, "127.0.0.1:0", "")
	node3.SetShortID(10)

	node4 := newMutableNode(gen.Reference(), insolar.StaticRoleVirtual, nil, insolar.NodeUndefined, "127.0.0.1:0", "")

	snapshot := NewSnapshot(pulse.MinTimePulse, []insolar.NetworkNode{node, node2, node3, node4})
	accessor := NewAccessor(snapshot)
	assert.Equal(t, 4, len(accessor.GetActiveNodes()))
	assert.Equal(t, 1, len(accessor.GetWorkingNodes()))
	assert.NotNil(t, accessor.GetWorkingNode(node.ID()))
	assert.Nil(t, accessor.GetWorkingNode(node2.ID()))
	assert.NotNil(t, accessor.GetActiveNode(node2.ID()))
	assert.NotNil(t, accessor.GetActiveNode(node3.ID()))
	assert.NotNil(t, accessor.GetActiveNode(node4.ID()))

	assert.NotNil(t, accessor.GetActiveNodeByShortID(10))
	assert.NotNil(t, accessor.GetActiveNodeByShortID(11))
	assert.Nil(t, accessor.GetActiveNodeByShortID(12))
}
