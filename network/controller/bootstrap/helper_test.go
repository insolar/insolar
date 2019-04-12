//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package bootstrap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/platformpolicy/keys"
	"github.com/insolar/insolar/testutils"
)

func newTestNode() insolar.NetworkNode {
	return node.NewNode(testutils.RandomRef(), insolar.StaticRoleUnknown, nil, "127.0.0.1:5432", "")
}

func newTestNodeWithShortID(id insolar.ShortNodeID) insolar.NetworkNode {
	n := newTestNode()
	n.(node.MutableNode).SetShortID(id)
	return n
}

func TestCorrectShortIDCollision(t *testing.T) {
	keeper := nodenetwork.NewNodeKeeper(newTestNode())
	keeper.SetInitialSnapshot([]insolar.NetworkNode{
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

	require.False(t, CheckShortIDCollision(keeper, insolar.ShortNodeID(2)))
	require.False(t, CheckShortIDCollision(keeper, insolar.ShortNodeID(31)))
	require.False(t, CheckShortIDCollision(keeper, insolar.ShortNodeID(35)))
	require.False(t, CheckShortIDCollision(keeper, insolar.ShortNodeID(65)))

	require.True(t, CheckShortIDCollision(keeper, insolar.ShortNodeID(30)))
	require.Equal(t, insolar.ShortNodeID(31), regenerateShortID(keeper, insolar.ShortNodeID(30)))

	require.True(t, CheckShortIDCollision(keeper, insolar.ShortNodeID(32)))
	require.Equal(t, insolar.ShortNodeID(35), regenerateShortID(keeper, insolar.ShortNodeID(32)))

	require.True(t, CheckShortIDCollision(keeper, insolar.ShortNodeID(64)))
	require.Equal(t, insolar.ShortNodeID(65), regenerateShortID(keeper, insolar.ShortNodeID(64)))

	require.True(t, CheckShortIDCollision(keeper, insolar.ShortNodeID(1<<32-2)))
	require.Equal(t, insolar.ShortNodeID(2), regenerateShortID(keeper, insolar.ShortNodeID(1<<32-2)))
}

type testNode struct {
	ref insolar.Reference
}

func (t *testNode) GetNodeRef() *insolar.Reference {
	return &t.ref
}

func (t *testNode) GetPublicKey() keys.PublicKey {
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

	discoveryNodes := []insolar.DiscoveryNode{first, originNode, second}
	result := RemoveOrigin(discoveryNodes, origin)
	assert.Equal(t, []insolar.DiscoveryNode{first, second}, result)

	discoveryNodes = []insolar.DiscoveryNode{first, second}
	result = RemoveOrigin(discoveryNodes, origin)
	assert.Equal(t, discoveryNodes, result)

	discoveryNodes = []insolar.DiscoveryNode{first, originNode}
	result = RemoveOrigin(discoveryNodes, origin)
	assert.Equal(t, []insolar.DiscoveryNode{first}, result)

	discoveryNodes = []insolar.DiscoveryNode{originNode, first}
	result = RemoveOrigin(discoveryNodes, origin)
	assert.Equal(t, []insolar.DiscoveryNode{first}, result)

	discoveryNodes = []insolar.DiscoveryNode{originNode}
	result = RemoveOrigin(discoveryNodes, origin)
	assert.Empty(t, result)
}
