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

package jetcoordinator_test

import (
	"bytes"
	"sort"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/ledgertestutils"
	"github.com/insolar/insolar/logicrunner"
	"github.com/stretchr/testify/assert"
)

func TestJetCoordinator_QueryRole(t *testing.T) {
	lr, err := logicrunner.NewLogicRunner(&configuration.LogicRunner{
		BuiltIn: &configuration.BuiltIn{},
	})
	assert.NoError(t, err)
	ledger, cleaner, keeper := ledgertestutils.TmpLedger(t, lr, "")
	defer cleaner()

	am := ledger.GetArtifactManager()
	pm := ledger.GetPulseManager()
	jc := ledger.GetJetCoordinator()

	pulse, err := pm.Current()
	assert.NoError(t, err)

	ref := func(r string) core.RecordRef { return core.NewRefFromBase58(r) }

	keeper.AddActiveNodes([]*core.ActiveNode{
		{NodeID: ref("v1"), Roles: []core.NodeRole{core.RoleVirtual}},
		{NodeID: ref("v2"), Roles: []core.NodeRole{core.RoleVirtual}},
		{NodeID: ref("l1"), Roles: []core.NodeRole{core.RoleLightMaterial}},
		{NodeID: ref("l2"), Roles: []core.NodeRole{core.RoleLightMaterial}},
		{NodeID: ref("l3"), Roles: []core.NodeRole{core.RoleLightMaterial}},
	})

	sorted := func(list []core.RecordRef) []core.RecordRef {
		sort.Slice(list, func(i, j int) bool {
			return bytes.Compare(list[i][:], list[j][:]) < 0
		})
		return list
	}

	selected, err := jc.QueryRole(core.RoleVirtualExecutor, *am.GenesisRef(), pulse.PulseNumber)
	assert.NoError(t, err)
	assert.Equal(t, []core.RecordRef{ref("v2")}, selected)

	selected, err = jc.QueryRole(core.RoleVirtualValidator, *am.GenesisRef(), pulse.PulseNumber)
	assert.Error(t, err)

	selected, err = jc.QueryRole(core.RoleLightValidator, *am.GenesisRef(), pulse.PulseNumber)
	assert.NoError(t, err)
	assert.Equal(t, sorted([]core.RecordRef{ref("l1"), ref("l2"), ref("l3")}), sorted(selected))

	selected, err = jc.QueryRole(core.RoleHeavyExecutor, *am.GenesisRef(), pulse.PulseNumber)
	assert.Error(t, err)
}
