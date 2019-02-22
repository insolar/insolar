/*
 *    Copyright 2019 Insolar
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

package nodes

import (
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func TestNodeStorage_Set(t *testing.T) {
	t.Parallel()
	firstNode := Node{FID: testutils.RandomRef()}
	secondNode := Node{FID: testutils.RandomRef()}
	nodeStorage := Storage{
		nodes: map[core.PulseNumber][]Node{},
	}
	err := nodeStorage.Set(1, []core.Node{firstNode, secondNode})

	require.NoError(t, err)
	require.Equal(t, 1, len(nodeStorage.nodes))
	require.Equal(t, firstNode, nodeStorage.nodes[1][0])
	require.Equal(t, secondNode, nodeStorage.nodes[1][1])
}

func TestNodeStorage_Set_OverrideError(t *testing.T) {
	t.Parallel()
	firstNode := Node{FID: testutils.RandomRef()}
	secondNode := Node{FID: testutils.RandomRef()}
	nodeStorage := Storage{
		nodes: map[core.PulseNumber][]Node{},
	}

	err := nodeStorage.Set(1, []core.Node{firstNode, secondNode})
	require.NoError(t, err)
	err = nodeStorage.Set(1, []core.Node{firstNode, secondNode})
	require.Error(t, err)

	require.Equal(t, 1, len(nodeStorage.nodes))
	require.Equal(t, firstNode, nodeStorage.nodes[1][0])
	require.Equal(t, secondNode, nodeStorage.nodes[1][1])
}

func TestNodeStorage_GetActiveNodes(t *testing.T) {
	t.Parallel()
	firstNode := Node{FID: testutils.RandomRef()}
	secondNode := Node{FID: testutils.RandomRef()}
	nodeStorage := Storage{
		nodes: map[core.PulseNumber][]Node{
			1: {firstNode, secondNode},
		},
	}

	result, err := nodeStorage.All(1)

	require.NoError(t, err)
	require.Equal(t, 2, len(result))
	require.Equal(t, firstNode, result[0])
	require.Equal(t, secondNode, result[1])
}

func TestNodeStorage_GetActiveNodes_FailsWhenNoNodes(t *testing.T) {
	t.Parallel()

	nodeStorage := Storage{
		nodes: map[core.PulseNumber][]Node{},
	}

	result, err := nodeStorage.All(1)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestNodeStorage_GetActiveNodesByRole(t *testing.T) {
	t.Parallel()
	nodeWithouRole := Node{}
	light := Node{FID: testutils.RandomRef(), FRole: core.StaticRoleLightMaterial}
	heavy := Node{FID: testutils.RandomRef(), FRole: core.StaticRoleHeavyMaterial}
	nodeStorage := Storage{
		nodes: map[core.PulseNumber][]Node{
			1: {nodeWithouRole, light, heavy},
		},
	}

	lightResult, err := nodeStorage.InRole(1, core.StaticRoleLightMaterial)
	require.NoError(t, err)
	heavyResult, err := nodeStorage.InRole(1, core.StaticRoleHeavyMaterial)
	require.NoError(t, err)

	require.Equal(t, 1, len(lightResult))
	require.Equal(t, light, lightResult[0])
	require.Equal(t, 1, len(heavyResult))
	require.Equal(t, heavy, heavyResult[0])
}

func TestNodeStorage_GetActiveNodesByRole_FailsWhenNoNode(t *testing.T) {
	t.Parallel()
	nodeStorage := Storage{
		nodes: map[core.PulseNumber][]Node{},
	}

	result, err := nodeStorage.InRole(1, core.StaticRoleLightMaterial)

	require.Error(t, err)
	require.Nil(t, result)
}

func TestNodeStorage_RemoveActiveNodesUntil(t *testing.T) {
	t.Parallel()
	nodeStorage := Storage{
		nodes: map[core.PulseNumber][]Node{
			1:   {},
			2:   {},
			222: {},
			555: {},
			5:   {},
		},
	}

	nodeStorage.Delete(1)
	nodeStorage.Delete(2)
	nodeStorage.Delete(5)

	require.Equal(t, 2, len(nodeStorage.nodes))
	_, ok := nodeStorage.nodes[222]
	require.Equal(t, true, ok)
	_, ok = nodeStorage.nodes[555]
	require.Equal(t, true, ok)
}
