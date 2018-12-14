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

package jetcoordinator

import (
	"testing"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJetCoordinator_QueryRole(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	nodeNet := network.NewNodeNetworkMock(mc)
	jc := JetCoordinator{
		db:                         db,
		NodeNet:                    nodeNet,
		PlatformCryptographyScheme: testutils.NewPlatformCryptographyScheme(),
	}
	err := db.AddPulse(ctx, core.Pulse{PulseNumber: 0, Entropy: core.Entropy{1, 2, 3}})
	require.NoError(t, err)
	var nodes []core.RecordRef
	for i := 0; i < 100; i++ {
		nodes = append(nodes, testutils.RandomRef())
	}

	t.Run("without object returns correct nodes", func(t *testing.T) {
		jc.roleCounts = map[core.DynamicRole]int{core.DynamicRoleVirtualExecutor: 3}
		nodeNet.GetActiveNodesByRoleMock.Expect(core.DynamicRoleVirtualExecutor).Return(nodes)

		selected, err := jc.QueryRole(ctx, core.DynamicRoleVirtualExecutor, nil, 0)
		require.NoError(t, err)
		assert.Equal(t, 3, len(selected))
		// Indexes are hard-coded from previously calculated values.
		assert.Equal(t, []core.RecordRef{nodes[25], nodes[78], nodes[36]}, selected)
	})

	t.Run("virtual returns correct nodes", func(t *testing.T) {
		jc.roleCounts = map[core.DynamicRole]int{core.DynamicRoleVirtualExecutor: 1}
		obj := core.RecordRef{}
		obj.SetRecord(*core.NewRecordID(0, []byte{3, 14, 15, 92}))
		nodeNet.GetActiveNodesByRoleMock.Expect(core.DynamicRoleVirtualExecutor).Return(nodes)

		selected, err := jc.QueryRole(ctx, core.DynamicRoleVirtualExecutor, obj.Record(), 0)
		require.NoError(t, err)
		assert.Equal(t, 1, len(selected))
		// Indexes are hard-coded from previously calculated values.
		assert.Equal(t, []core.RecordRef{nodes[22]}, selected)
	})

	t.Run("material returns correct nodes", func(t *testing.T) {
		objID := core.NewRecordID(0, []byte{1, 42, 123})
		jc.roleCounts = map[core.DynamicRole]int{core.DynamicRoleLightExecutor: 1}
		err := db.UpdateJetTree(ctx, 0, *jet.NewID(1, []byte{1, 42, 123}))
		require.NoError(t, err)
		nodeNet.GetActiveNodesByRoleMock.Expect(core.DynamicRoleLightExecutor).Return(nodes)

		selected, err := jc.QueryRole(ctx, core.DynamicRoleLightExecutor, objID, 0)
		require.NoError(t, err)
		assert.Equal(t, 1, len(selected))
		// Indexes are hard-coded from previously calculated values.
		assert.Equal(t, []core.RecordRef{nodes[25]}, selected)
	})
}
