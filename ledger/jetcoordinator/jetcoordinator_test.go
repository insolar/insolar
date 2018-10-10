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
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/ledgertestutil"
	"github.com/insolar/insolar/logicrunner"
	"github.com/stretchr/testify/assert"
)

func TestJetCoordinator_QueryRole(t *testing.T) {
	lr, err := logicrunner.NewLogicRunner(&configuration.LogicRunner{
		BuiltIn: &configuration.BuiltIn{},
	})
	assert.NoError(t, err)
	ledger, cleaner := ledgertestutil.TmpLedger(t, lr, "")
	defer cleaner()

	am := ledger.GetArtifactManager()
	pm := ledger.GetPulseManager()
	jc := ledger.GetJetCoordinator()

	pulse, err := pm.Current()
	assert.NoError(t, err)

	selected, err := jc.QueryRole(core.RoleVirtualExecutor, *am.GenesisRef(), pulse.PulseNumber)
	assert.NoError(t, err)
	assert.Equal(t, []core.RecordRef{core.NewRefFromBase58("ve2")}, selected)

	selected, err = jc.QueryRole(core.RoleVirtualValidator, *am.GenesisRef(), pulse.PulseNumber)
	assert.NoError(t, err)
	assert.Equal(t, []core.RecordRef{
		core.NewRefFromBase58("vv3"),
		core.NewRefFromBase58("vv1"),
		core.NewRefFromBase58("vv4"),
	}, selected)
}

func TestJetCoordinator_IsAuthorized(t *testing.T) {
	lr, err := logicrunner.NewLogicRunner(&configuration.LogicRunner{
		BuiltIn: &configuration.BuiltIn{},
	})
	assert.NoError(t, err)
	ledger, cleaner := ledgertestutil.TmpLedger(t, lr, "")
	defer cleaner()

	am := ledger.GetArtifactManager()
	pm := ledger.GetPulseManager()
	jc := ledger.GetJetCoordinator()

	pulse, err := pm.Current()
	assert.NoError(t, err)

	authorized, err := jc.IsAuthorized(
		core.RoleVirtualExecutor, *am.GenesisRef(), pulse.PulseNumber, core.NewRefFromBase58("ve1"),
	)
	assert.NoError(t, err)
	assert.Equal(t, false, authorized)

	authorized, err = jc.IsAuthorized(
		core.RoleVirtualExecutor, *am.GenesisRef(), pulse.PulseNumber, core.NewRefFromBase58("ve2"),
	)
	assert.NoError(t, err)
	assert.Equal(t, true, authorized)
}
