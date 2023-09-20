package network

import (
	"crypto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/network/node"
)

func newTestNode() insolar.NetworkNode {
	return node.NewNode(gen.Reference(), insolar.StaticRoleUnknown, nil, "127.0.0.1:5432", "")
}

func newTestNodeWithShortID(id insolar.ShortNodeID) insolar.NetworkNode {
	n := newTestNode()
	n.(node.MutableNode).SetShortID(id)
	return n
}

func TestCorrectShortIDCollision(t *testing.T) {

	nodes := []insolar.NetworkNode{
		newTestNodeWithShortID(0),
		newTestNodeWithShortID(1),
		newTestNodeWithShortID(30),
		newTestNodeWithShortID(32),
		newTestNodeWithShortID(33),
		newTestNodeWithShortID(34),
		newTestNodeWithShortID(64),
		newTestNodeWithShortID(1<<32 - 2),
		newTestNodeWithShortID(1<<32 - 1),
	}

	require.False(t, CheckShortIDCollision(nodes, insolar.ShortNodeID(2)))
	require.False(t, CheckShortIDCollision(nodes, insolar.ShortNodeID(31)))
	require.False(t, CheckShortIDCollision(nodes, insolar.ShortNodeID(35)))
	require.False(t, CheckShortIDCollision(nodes, insolar.ShortNodeID(65)))

	require.True(t, CheckShortIDCollision(nodes, insolar.ShortNodeID(30)))
	require.Equal(t, insolar.ShortNodeID(31), regenerateShortID(nodes, insolar.ShortNodeID(30)))

	require.True(t, CheckShortIDCollision(nodes, insolar.ShortNodeID(32)))
	require.Equal(t, insolar.ShortNodeID(35), regenerateShortID(nodes, insolar.ShortNodeID(32)))

	require.True(t, CheckShortIDCollision(nodes, insolar.ShortNodeID(64)))
	require.Equal(t, insolar.ShortNodeID(65), regenerateShortID(nodes, insolar.ShortNodeID(64)))

	require.True(t, CheckShortIDCollision(nodes, insolar.ShortNodeID(1<<32-2)))
	require.Equal(t, insolar.ShortNodeID(2), regenerateShortID(nodes, insolar.ShortNodeID(1<<32-2)))
}

type testNode struct {
	ref insolar.Reference
}

func (t *testNode) GetNodeRef() *insolar.Reference {
	return &t.ref
}

func (t *testNode) GetPublicKey() crypto.PublicKey {
	return nil
}

func (t *testNode) GetHost() string {
	return ""
}

func (t *testNode) GetBriefDigest() []byte {
	return nil
}

func (t *testNode) GetBriefSign() []byte {
	return nil
}

func (t *testNode) GetRole() insolar.StaticRole {
	return insolar.StaticRoleVirtual
}

func TestExcludeOrigin(t *testing.T) {
	origin := gen.Reference()
	originNode := &testNode{origin}
	first := &testNode{gen.Reference()}
	second := &testNode{gen.Reference()}

	discoveryNodes := []insolar.DiscoveryNode{first, originNode, second}
	result := ExcludeOrigin(discoveryNodes, origin)
	assert.Equal(t, []insolar.DiscoveryNode{first, second}, result)

	discoveryNodes = []insolar.DiscoveryNode{first, second}
	result = ExcludeOrigin(discoveryNodes, origin)
	assert.Equal(t, discoveryNodes, result)

	discoveryNodes = []insolar.DiscoveryNode{first, originNode}
	result = ExcludeOrigin(discoveryNodes, origin)
	assert.Equal(t, []insolar.DiscoveryNode{first}, result)

	discoveryNodes = []insolar.DiscoveryNode{originNode, first}
	result = ExcludeOrigin(discoveryNodes, origin)
	assert.Equal(t, []insolar.DiscoveryNode{first}, result)

	discoveryNodes = []insolar.DiscoveryNode{originNode}
	result = ExcludeOrigin(discoveryNodes, origin)
	assert.Empty(t, result)

	discoveryNodes = []insolar.DiscoveryNode{originNode, first, second}
	result = ExcludeOrigin(discoveryNodes, origin)
	assert.Equal(t, []insolar.DiscoveryNode{first, second}, result)

	discoveryNodes = []insolar.DiscoveryNode{first, second, originNode}
	result = ExcludeOrigin(discoveryNodes, origin)
	assert.Equal(t, []insolar.DiscoveryNode{first, second}, result)

}
