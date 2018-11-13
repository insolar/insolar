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

package auth

import (
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
)

func newTestNode() core.Node {
	return nodenetwork.NewNode(testutils.RandomRef(), nil, nil,
		core.PulseNumber(0), "", "")
}

func newTestNodeWithShortID(id core.ShortNodeID) core.Node {
	node := newTestNode()
	node.(nodenetwork.MutableNode).SetShortID(id)
	return node
}

func TestCorrectShortIDCollision(t *testing.T) {
	keeper := nodenetwork.NewNodeKeeper(newTestNode())
	keeper.AddActiveNodes([]core.Node{
		newTestNodeWithShortID(30),
		newTestNodeWithShortID(32),
		newTestNodeWithShortID(33),
		newTestNodeWithShortID(34),
		newTestNodeWithShortID(64),
	})

	assert.False(t, CheckShortIDCollision(keeper, core.ShortNodeID(0)))
	assert.False(t, CheckShortIDCollision(keeper, core.ShortNodeID(31)))
	assert.False(t, CheckShortIDCollision(keeper, core.ShortNodeID(35)))
	assert.False(t, CheckShortIDCollision(keeper, core.ShortNodeID(65)))

	assert.True(t, CheckShortIDCollision(keeper, core.ShortNodeID(30)))
	node := newTestNodeWithShortID(30)
	CorrectShortIDCollision(keeper, node)
	assert.Equal(t, core.ShortNodeID(31), node.ShortID())

	assert.True(t, CheckShortIDCollision(keeper, core.ShortNodeID(32)))
	node = newTestNodeWithShortID(32)
	CorrectShortIDCollision(keeper, node)
	assert.Equal(t, core.ShortNodeID(35), node.ShortID())

	assert.True(t, CheckShortIDCollision(keeper, core.ShortNodeID(64)))
	node = newTestNodeWithShortID(64)
	CorrectShortIDCollision(keeper, node)
	assert.Equal(t, core.ShortNodeID(65), node.ShortID())
}
