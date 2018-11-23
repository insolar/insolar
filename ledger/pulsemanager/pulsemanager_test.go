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

package pulsemanager_test

import (
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/ledgertestutils"
	"github.com/insolar/insolar/logicrunner"
	"github.com/stretchr/testify/assert"
)

func TestPulseManager_Current(t *testing.T) {
	ctx := inslogger.TestContext(t)
	lr, err := logicrunner.NewLogicRunner(&configuration.LogicRunner{
		BuiltIn: &configuration.BuiltIn{},
	})
	assert.NoError(t, err)
	c := core.Components{LogicRunner: lr}
	// FIXME: TmpLedger is deprecated. Use mocks instead.
	ledger, cleaner := ledgertestutils.TmpLedger(t, "", c)
	defer cleaner()

	pm := ledger.GetPulseManager()

	pulse, err := pm.Current(ctx)
	assert.NoError(t, err)
	assert.Equal(t, core.Pulse{PulseNumber: core.FirstPulseNumber}, *pulse)
}
