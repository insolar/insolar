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

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/ledgertestutil"
	"github.com/stretchr/testify/assert"
)

func TestJetCoordinator_QueryRole(t *testing.T) {
	ledger, cleaner := ledgertestutil.TmpLedger(t, "")
	defer cleaner()

	var err error

	am := ledger.GetArtifactManager()
	pm := ledger.GetPulseManager()
	jc := ledger.GetJetCoordinator()

	pulse, err := pm.Current()
	assert.NoError(t, err)

	selected, err := jc.QueryRole(core.RoleVirtualExecutor, *am.RootRef(), pulse.PulseNumber)
	assert.Equal(t, []core.RecordRef{core.String2Ref("ve2")}, selected)

	selected, err = jc.QueryRole(core.RoleVirtualValidator, *am.RootRef(), pulse.PulseNumber)
	assert.Equal(t, []core.RecordRef{
		core.String2Ref("vv3"),
		core.String2Ref("vv1"),
		core.String2Ref("vv4"),
	}, selected)
}

func TestJetCoordinator_IsAuthorized(t *testing.T) {
	ledger, cleaner := ledgertestutil.TmpLedger(t, "")
	defer cleaner()

	var err error

	am := ledger.GetArtifactManager()
	pm := ledger.GetPulseManager()
	jc := ledger.GetJetCoordinator()

	pulse, err := pm.Current()
	assert.NoError(t, err)

	authorized, err := jc.IsAuthorized(
		core.RoleVirtualExecutor, *am.RootRef(), pulse.PulseNumber, core.String2Ref("ve1"),
	)
	assert.Equal(t, false, authorized)

	authorized, err = jc.IsAuthorized(
		core.RoleVirtualExecutor, *am.RootRef(), pulse.PulseNumber, core.String2Ref("ve2"),
	)
	assert.Equal(t, true, authorized)
}
