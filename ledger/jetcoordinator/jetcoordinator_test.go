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
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJetCoordinator_QueryRole(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	jc := JetCoordinator{
		db:                         db,
		PlatformCryptographyScheme: testutils.NewPlatformCryptographyScheme(),
	}
	err := db.AddPulse(ctx, core.Pulse{PulseNumber: 0, Entropy: core.Entropy{1, 2, 3}})
	require.NoError(t, err)
	var nodes []core.Node
	var nodeRefs []core.RecordRef
	for i := 0; i < 100; i++ {
		ref := *core.NewRecordRef(core.DomainID, *core.NewRecordID(0, []byte{byte(i)}))
		nodes = append(nodes, storage.Node{FID: ref, FRole: core.StaticRoleLightMaterial})
		nodeRefs = append(nodeRefs, ref)
	}
	err = db.SetActiveNodes(0, nodes)
	require.NoError(t, err)

	objID := core.NewRecordID(0, []byte{1, 42, 123})
	jc.roleCounts = map[core.DynamicRole]int{core.DynamicRoleLightValidator: 3}
	err = db.UpdateJetTree(ctx, 0, true, *jet.NewID(50, []byte{1, 42, 123}))
	require.NoError(t, err)

	selected, err := jc.QueryRole(ctx, core.DynamicRoleLightValidator, *objID, 0)
	require.NoError(t, err)
	assert.Equal(t, 3, len(selected))

	// Indexes are hard-coded from previously calculated values.
	assert.Equal(t, []core.RecordRef{nodeRefs[16], nodeRefs[21], nodeRefs[78]}, selected)
}
