//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package replication

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
	"github.com/stretchr/testify/require"
)

func TestCleaner_cleanPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)

	inputPulse := insolar.Pulse{PulseNumber: insolar.PulseNumber(111)}

	ctrl := minimock.NewController(t)

	jm := jet.NewCleanerMock(ctrl)
	jm.DeleteForPNMock.Expect(ctx, inputPulse.PulseNumber)

	nm := node.NewModifierMock(ctrl)
	nm.DeleteForPNMock.Expect(inputPulse.PulseNumber)

	dc := drop.NewCleanerMock(ctrl)
	dc.DeleteForPNMock.Expect(ctx, inputPulse.PulseNumber)

	bc := blob.NewCleanerMock(ctrl)
	bc.DeleteForPNMock.Expect(ctx, inputPulse.PulseNumber)

	rc := object.NewRecordCleanerMock(ctrl)
	rc.DeleteForPNMock.Expect(ctx, inputPulse.PulseNumber)

	ic := object.NewIndexCleanerMock(ctrl)
	ic.DeleteForPNMock.Expect(ctx, inputPulse.PulseNumber)

	ps := pulse.NewShifterMock(ctrl)
	ps.ShiftMock.Expect(ctx, inputPulse.PulseNumber).Return(nil)

	cleaner := NewCleaner(jm, nm, dc, bc, rc, ic, ps, nil, 0)

	cleaner.cleanPulse(ctx, inputPulse.PulseNumber)

	ctrl.Finish()
}

func DeleteForPNMock(t *testing.T, expected insolar.PulseNumber) func(p context.Context, p1 insolar.PulseNumber) {
	return func(ctx context.Context, actual insolar.PulseNumber) {
		require.Equal(t, expected, actual)
	}
}

func TestCleaner_clean(t *testing.T) {
	ctx := inslogger.TestContext(t)

	inputPulse := insolar.Pulse{PulseNumber: insolar.PulseNumber(111)}
	calculatedPulse := insolar.Pulse{PulseNumber: insolar.PulseNumber(98765)}
	limit := 123

	ctrl := minimock.NewController(t)

	jm := jet.NewCleanerMock(ctrl)
	jm.DeleteForPNFunc = DeleteForPNMock(t, calculatedPulse.PulseNumber)

	nm := node.NewModifierMock(ctrl)
	nm.DeleteForPNMock.Expect(calculatedPulse.PulseNumber)

	dc := drop.NewCleanerMock(ctrl)
	dc.DeleteForPNFunc = DeleteForPNMock(t, calculatedPulse.PulseNumber)

	bc := blob.NewCleanerMock(ctrl)
	bc.DeleteForPNFunc = func(p context.Context, pn insolar.PulseNumber) {
		require.Equal(t, calculatedPulse.PulseNumber, pn)
	}

	rc := object.NewRecordCleanerMock(ctrl)
	rc.DeleteForPNFunc = DeleteForPNMock(t, calculatedPulse.PulseNumber)

	ic := object.NewIndexCleanerMock(ctrl)
	ic.DeleteForPNFunc = DeleteForPNMock(t, calculatedPulse.PulseNumber)

	ps := pulse.NewShifterMock(ctrl)
	ps.ShiftFunc = func(p context.Context, pn insolar.PulseNumber) error {
		require.Equal(t, calculatedPulse.PulseNumber, pn)
		return nil
	}

	pc := pulse.NewCalculatorMock(ctrl)
	pc.BackwardsFunc = func(p context.Context, pn insolar.PulseNumber, l int) (r insolar.Pulse, r1 error) {
		require.Equal(t, inputPulse.PulseNumber, pn)
		require.Equal(t, limit+1, l)
		return calculatedPulse, nil
	}

	cleaner := NewCleaner(jm, nm, dc, bc, rc, ic, ps, pc, limit)
	defer close(cleaner.pulseForClean)

	go cleaner.clean(ctx)
	cleaner.pulseForClean <- inputPulse.PulseNumber

	ctrl.Wait(time.Minute)
}

func TestLightCleaner_NotifyAboutPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)

	inputPulse := insolar.Pulse{PulseNumber: insolar.PulseNumber(111)}
	calculatedPulse := insolar.Pulse{PulseNumber: insolar.PulseNumber(98765)}
	limit := 123

	ctrl := minimock.NewController(t)

	jm := jet.NewCleanerMock(ctrl)
	jm.DeleteForPNFunc = DeleteForPNMock(t, calculatedPulse.PulseNumber)

	nm := node.NewModifierMock(ctrl)
	nm.DeleteForPNMock.Expect(calculatedPulse.PulseNumber)

	dc := drop.NewCleanerMock(ctrl)
	dc.DeleteForPNFunc = DeleteForPNMock(t, calculatedPulse.PulseNumber)

	bc := blob.NewCleanerMock(ctrl)
	bc.DeleteForPNFunc = DeleteForPNMock(t, calculatedPulse.PulseNumber)

	rc := object.NewRecordCleanerMock(ctrl)
	rc.DeleteForPNFunc = DeleteForPNMock(t, calculatedPulse.PulseNumber)

	ic := object.NewIndexCleanerMock(ctrl)
	ic.DeleteForPNFunc = DeleteForPNMock(t, calculatedPulse.PulseNumber)

	ps := pulse.NewShifterMock(ctrl)
	ps.ShiftFunc = func(p context.Context, pn insolar.PulseNumber) error {
		require.Equal(t, calculatedPulse.PulseNumber, pn)
		return nil
	}

	pc := pulse.NewCalculatorMock(ctrl)
	pc.BackwardsFunc = func(p context.Context, pn insolar.PulseNumber, l int) (r insolar.Pulse, r1 error) {
		require.Equal(t, inputPulse.PulseNumber, pn)
		require.Equal(t, limit+1, l)
		return calculatedPulse, nil
	}

	cleaner := NewCleaner(jm, nm, dc, bc, rc, ic, ps, pc, limit)
	defer close(cleaner.pulseForClean)

	go cleaner.NotifyAboutPulse(ctx, inputPulse.PulseNumber)

	ctrl.Wait(time.Minute)
}
