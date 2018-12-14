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

package bootstrap

import (
	"crypto"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestNode() core.Node {
	return nodenetwork.NewNode(testutils.RandomRef(), core.StaticRoleUnknown, nil, "", "")
}

func newTestNodeWithShortID(id core.ShortNodeID) core.Node {
	node := newTestNode()
	node.(nodenetwork.MutableNode).SetShortID(id)
	return node
}

func TestCorrectShortIDCollision(t *testing.T) {
	keeper := nodenetwork.NewNodeKeeper(newTestNode())
	keeper.AddActiveNodes([]core.Node{
		newTestNodeWithShortID(0),
		newTestNodeWithShortID(1),
		newTestNodeWithShortID(30),
		newTestNodeWithShortID(32),
		newTestNodeWithShortID(33),
		newTestNodeWithShortID(34),
		newTestNodeWithShortID(64),
		newTestNodeWithShortID(1<<32 - 2),
		newTestNodeWithShortID(1<<32 - 1),
	})

	require.False(t, checkShortIDCollision(keeper, core.ShortNodeID(2)))
	require.False(t, checkShortIDCollision(keeper, core.ShortNodeID(31)))
	require.False(t, checkShortIDCollision(keeper, core.ShortNodeID(35)))
	require.False(t, checkShortIDCollision(keeper, core.ShortNodeID(65)))

	require.True(t, checkShortIDCollision(keeper, core.ShortNodeID(30)))
	require.Equal(t, core.ShortNodeID(31), regenerateShortID(keeper, core.ShortNodeID(30)))

	require.True(t, checkShortIDCollision(keeper, core.ShortNodeID(32)))
	require.Equal(t, core.ShortNodeID(35), regenerateShortID(keeper, core.ShortNodeID(32)))

	require.True(t, checkShortIDCollision(keeper, core.ShortNodeID(64)))
	require.Equal(t, core.ShortNodeID(65), regenerateShortID(keeper, core.ShortNodeID(64)))

	require.True(t, checkShortIDCollision(keeper, core.ShortNodeID(1<<32-2)))
	require.Equal(t, core.ShortNodeID(2), regenerateShortID(keeper, core.ShortNodeID(1<<32-2)))
}

type testNode struct {
	ref core.RecordRef
}

func (t *testNode) GetNodeRef() *core.RecordRef {
	return &t.ref
}

func (t *testNode) GetPublicKey() crypto.PublicKey {
	return nil
}

func (t *testNode) GetHost() string {
	return ""
}

func TestRemoveOrigin(t *testing.T) {
	origin := testutils.RandomRef()
	originNode := &testNode{origin}
	first := &testNode{testutils.RandomRef()}
	second := &testNode{testutils.RandomRef()}

	discoveryNodes := []core.DiscoveryNode{first, originNode, second}
	result, err := RemoveOrigin(discoveryNodes, origin)
	require.NoError(t, err)
	assert.Equal(t, []core.DiscoveryNode{first, second}, result)

	discoveryNodes = []core.DiscoveryNode{first, second}
	_, err = RemoveOrigin(discoveryNodes, origin)
	assert.Error(t, err)

	discoveryNodes = []core.DiscoveryNode{first, originNode}
	result, err = RemoveOrigin(discoveryNodes, origin)
	require.NoError(t, err)
	assert.Equal(t, []core.DiscoveryNode{first}, result)

	discoveryNodes = []core.DiscoveryNode{originNode, first}
	result, err = RemoveOrigin(discoveryNodes, origin)
	require.NoError(t, err)
	assert.Equal(t, []core.DiscoveryNode{first}, result)

	discoveryNodes = []core.DiscoveryNode{originNode}
	result, err = RemoveOrigin(discoveryNodes, origin)
	require.NoError(t, err)
	assert.Empty(t, result)
}
