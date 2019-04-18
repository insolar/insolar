/*
 *    Copyright 2019 Insolar Technologies
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

package replication

import (
	"testing"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/blob"
	"github.com/insolar/insolar/ledger/storage/drop"
	"github.com/insolar/insolar/ledger/storage/node"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/insolar/insolar/ledger/storage/pulse"
)

func TestCleaner_Clean(t *testing.T) {
	ctx := inslogger.TestContext(t)

	inputPulse := insolar.Pulse{PulseNumber: insolar.FirstPulseNumber}
	calculatedPulse := insolar.Pulse{PulseNumber: insolar.PulseNumber(98765)}
	lcl := 123

	ctrl := minimock.NewController(t)

	jm := jet.NewStorageMock(ctrl)
	jm.DeleteForPNMock.Expect(ctx, calculatedPulse.PulseNumber)

	nm := node.NewModifierMock(ctrl)
	nm.DeleteForPNMock.Expect(calculatedPulse.PulseNumber)

	dc := drop.NewCleanerMock(ctrl)
	dc.DeleteForPNMock.Expect(ctx, calculatedPulse.PulseNumber)

	bc := blob.NewCleanerMock(ctrl)
	bc.DeleteForPNMock.Expect(ctx, calculatedPulse.PulseNumber)

	rc := object.NewRecordCleanerMock(ctrl)
	rc.DeleteForPNMock.Expect(ctx, calculatedPulse.PulseNumber)

	ic := object.NewIndexCleanerMock(ctrl)
	ic.RemoveForPulseMock.Expect(ctx, calculatedPulse.PulseNumber)

	ps := pulse.NewShifterMock(ctrl)
	ps.ShiftMock.Expect(ctx, calculatedPulse.PulseNumber).Return(nil)

	pc := pulse.NewCalculatorMock(ctrl)
	pc.BackwardsMock.Expect(ctx, inputPulse.PulseNumber, lcl).Return(calculatedPulse, nil)

	cleaner := NewCleaner(jm, nm, dc, bc, rc, ic, ps, pc, lcl)
	cleaner.NotifyAboutPulse(ctx, inputPulse.PulseNumber)

	ctrl.Finish()
}
