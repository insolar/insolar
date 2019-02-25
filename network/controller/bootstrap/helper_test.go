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
 *
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
	return nodenetwork.NewNode(testutils.RandomRef(), core.StaticRoleUnknown, nil, "127.0.0.1:5432", "")
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

	require.False(t, CheckShortIDCollision(keeper, core.ShortNodeID(2)))
	require.False(t, CheckShortIDCollision(keeper, core.ShortNodeID(31)))
	require.False(t, CheckShortIDCollision(keeper, core.ShortNodeID(35)))
	require.False(t, CheckShortIDCollision(keeper, core.ShortNodeID(65)))

	require.True(t, CheckShortIDCollision(keeper, core.ShortNodeID(30)))
	require.Equal(t, core.ShortNodeID(31), regenerateShortID(keeper, core.ShortNodeID(30)))

	require.True(t, CheckShortIDCollision(keeper, core.ShortNodeID(32)))
	require.Equal(t, core.ShortNodeID(35), regenerateShortID(keeper, core.ShortNodeID(32)))

	require.True(t, CheckShortIDCollision(keeper, core.ShortNodeID(64)))
	require.Equal(t, core.ShortNodeID(65), regenerateShortID(keeper, core.ShortNodeID(64)))

	require.True(t, CheckShortIDCollision(keeper, core.ShortNodeID(1<<32-2)))
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
	result, err := removeOrigin(discoveryNodes, origin)
	require.NoError(t, err)
	assert.Equal(t, []core.DiscoveryNode{first, second}, result)

	discoveryNodes = []core.DiscoveryNode{first, second}
	_, err = removeOrigin(discoveryNodes, origin)
	assert.Error(t, err)

	discoveryNodes = []core.DiscoveryNode{first, originNode}
	result, err = removeOrigin(discoveryNodes, origin)
	require.NoError(t, err)
	assert.Equal(t, []core.DiscoveryNode{first}, result)

	discoveryNodes = []core.DiscoveryNode{originNode, first}
	result, err = removeOrigin(discoveryNodes, origin)
	require.NoError(t, err)
	assert.Equal(t, []core.DiscoveryNode{first}, result)

	discoveryNodes = []core.DiscoveryNode{originNode}
	result, err = removeOrigin(discoveryNodes, origin)
	require.NoError(t, err)
	assert.Empty(t, result)
}
