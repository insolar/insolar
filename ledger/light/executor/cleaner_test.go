package executor

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
)

func TestCleaner_cleanPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)

	inputPulse := insolar.Pulse{PulseNumber: insolar.PulseNumber(111)}
	latestPulse := insolar.PulseNumber(123)

	ctrl := minimock.NewController(t)

	jm := jet.NewCleanerMock(ctrl)
	jm.DeleteForPNMock.Expect(ctx, inputPulse.PulseNumber)

	nm := node.NewModifierMock(ctrl)
	nm.DeleteForPNMock.Expect(inputPulse.PulseNumber)

	dc := drop.NewCleanerMock(ctrl)
	dc.DeleteForPNMock.Expect(ctx, inputPulse.PulseNumber)

	rc := object.NewRecordCleanerMock(ctrl)
	rc.DeleteForPNMock.Expect(ctx, inputPulse.PulseNumber)

	ic := object.NewIndexCleanerMock(ctrl)
	ic.DeleteForPNMock.Expect(ctx, inputPulse.PulseNumber)

	ps := pulse.NewShifterMock(ctrl)
	ps.ShiftMock.Expect(ctx, inputPulse.PulseNumber).Return(nil)

	prevPulseFromInput := insolar.Pulse{PulseNumber: inputPulse.PulseNumber - 1}
	pc := pulse.NewCalculatorMock(ctrl)
	pc.BackwardsMock.Inspect(func(ctx context.Context, pn insolar.PulseNumber, steps int) {
		require.Equal(t, latestPulse, pn)
		require.Equal(t, 1, steps)
	}).Return(prevPulseFromInput, nil)

	objID := gen.ID()
	ia := object.NewIndexAccessorMock(ctrl)
	ia.ForPulseMock.Expect(ctx, prevPulseFromInput.PulseNumber).Return([]record.Index{
		{ObjID: objID, LifelineLastUsed: latestPulse},
	}, nil)

	fc := NewFilamentCleanerMock(ctrl)
	fc.ClearIfLongerMock.Expect(100)
	fc.ClearAllExceptMock.Inspect(func(ids []insolar.ID) {
		require.Equal(t, 1, len(ids))
		require.Equal(t, objID, ids[0])
	}).Return()

	cleaner := NewCleaner(jm, nm, dc, rc, ic, ps, pc, ia, fc, 0, 0, 100)

	cleaner.cleanPulse(ctx, inputPulse.PulseNumber, latestPulse)

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
	defer ctrl.Finish()

	jm := jet.NewCleanerMock(ctrl)
	jm.DeleteForPNMock.Set(DeleteForPNMock(t, calculatedPulse.PulseNumber))

	nm := node.NewModifierMock(ctrl)
	nm.DeleteForPNMock.Expect(calculatedPulse.PulseNumber)

	dc := drop.NewCleanerMock(ctrl)
	dc.DeleteForPNMock.Set(DeleteForPNMock(t, calculatedPulse.PulseNumber))

	rc := object.NewRecordCleanerMock(ctrl)
	rc.DeleteForPNMock.Set(DeleteForPNMock(t, calculatedPulse.PulseNumber))

	ic := object.NewIndexCleanerMock(ctrl)
	ic.DeleteForPNMock.Set(DeleteForPNMock(t, calculatedPulse.PulseNumber))

	ps := pulse.NewShifterMock(ctrl)
	ps.ShiftMock.Inspect(func(ctx context.Context, pn insolar.PulseNumber) {
		require.Equal(t, calculatedPulse.PulseNumber, pn)
	}).Return(nil)

	pc := pulse.NewCalculatorMock(ctrl)
	pc.BackwardsMock.Return(calculatedPulse, nil)

	objID := gen.ID()
	ia := object.NewIndexAccessorMock(ctrl)
	ia.ForPulseMock.Inspect(func(ctx context.Context, pn insolar.PulseNumber) {
		require.Equal(t, calculatedPulse.PulseNumber, pn)
	}).Return([]record.Index{
		{ObjID: objID, LifelineLastUsed: insolar.PulseNumber(110)},
	}, nil)

	fc := NewFilamentCleanerMock(ctrl)
	fc.ClearAllExceptMock.Expect([]insolar.ID{objID})
	fc.ClearIfLongerMock.Inspect(func(limit int) {
		require.Equal(t, limit, 100)
	}).Return()

	cleaner := NewCleaner(jm, nm, dc, rc, ic, ps, pc, ia, fc, limit, 1, 100)
	defer close(cleaner.pulseForClean)

	go cleaner.clean(ctx)
	cleaner.pulseForClean <- inputPulse.PulseNumber

	ctrl.Wait(time.Minute * 10)
}

func TestLightCleaner_NotifyAboutPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)

	inputPulse := insolar.Pulse{PulseNumber: insolar.PulseNumber(111)}
	calculatedPulse := insolar.Pulse{PulseNumber: insolar.PulseNumber(98765)}
	limit := 123

	ctrl := minimock.NewController(t)

	jm := jet.NewCleanerMock(ctrl)
	jm.DeleteForPNMock.Set(DeleteForPNMock(t, calculatedPulse.PulseNumber))

	nm := node.NewModifierMock(ctrl)
	nm.DeleteForPNMock.Expect(calculatedPulse.PulseNumber)

	dc := drop.NewCleanerMock(ctrl)
	dc.DeleteForPNMock.Set(DeleteForPNMock(t, calculatedPulse.PulseNumber))

	rc := object.NewRecordCleanerMock(ctrl)
	rc.DeleteForPNMock.Set(DeleteForPNMock(t, calculatedPulse.PulseNumber))

	ic := object.NewIndexCleanerMock(ctrl)
	ic.DeleteForPNMock.Set(DeleteForPNMock(t, calculatedPulse.PulseNumber))

	ps := pulse.NewShifterMock(ctrl)
	ps.ShiftMock.Inspect(func(ctx context.Context, pn insolar.PulseNumber) {
		require.Equal(t, calculatedPulse.PulseNumber, pn)
	}).Return(nil)

	pc := pulse.NewCalculatorMock(ctrl)
	pc.BackwardsMock.Inspect(func(ctx context.Context, pn insolar.PulseNumber, steps int) {
		switch pn {
		case inputPulse.PulseNumber:
		default:
			require.Fail(t, "wrong input")
		}
	}).Return(calculatedPulse, nil)

	objID := gen.ID()
	ia := object.NewIndexAccessorMock(ctrl)
	ia.ForPulseMock.Inspect(func(ctx context.Context, pn insolar.PulseNumber) {
		require.Equal(t, calculatedPulse.PulseNumber, pn)
	}).Return([]record.Index{
		{ObjID: objID, LifelineLastUsed: insolar.PulseNumber(110)},
	}, nil)

	fc := NewFilamentCleanerMock(ctrl)
	fc.ClearIfLongerMock.Expect(100)
	fc.ClearAllExceptMock.Inspect(func(ids []insolar.ID) {
		require.Equal(t, 1, len(ids))
		require.Equal(t, objID, ids[0])
	}).Return()

	cleaner := NewCleaner(jm, nm, dc, rc, ic, ps, pc, ia, fc, limit, 1, 100)
	defer close(cleaner.pulseForClean)

	go cleaner.NotifyAboutPulse(ctx, inputPulse.PulseNumber)

	ctrl.Wait(time.Minute)
}
