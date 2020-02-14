// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package executor_test

import (
	"context"
	"sync"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/pulse"
)

type TestBadgerGCRunner struct {
	lock  sync.RWMutex
	count uint
}

func (t *TestBadgerGCRunner) RunValueGC(ctx context.Context) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.count++
}

func (t *TestBadgerGCRunner) getCount() uint {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.count
}

func TestBadgerGCRunInfo(t *testing.T) {

	ctx := inslogger.TestContext(t)

	t.Run("call every time if frequency equal 1", func(t *testing.T) {
		t.Parallel()
		runner := &TestBadgerGCRunner{}
		info := executor.NewBadgerGCRunInfo(runner, 1)
		for i := 1; i < 5; i++ {
			done := info.RunGCIfNeeded(ctx)
			<-done
			require.Equal(t, uint(i), runner.getCount())
		}
	})

	t.Run("no call if frequency equal 0", func(t *testing.T) {
		t.Parallel()
		runner := &TestBadgerGCRunner{}
		info := executor.NewBadgerGCRunInfo(runner, 0)
		for i := 1; i < 5; i++ {
			done := info.RunGCIfNeeded(ctx)
			<-done
			require.Equal(t, uint(0), runner.getCount())
		}
	})

	t.Run("even calls if frequency equal 2", func(t *testing.T) {
		t.Parallel()
		runner := &TestBadgerGCRunner{}
		info := executor.NewBadgerGCRunInfo(runner, 2)
		for i := 1; i < 5; i++ {
			done := info.RunGCIfNeeded(ctx)
			<-done
			require.Equal(t, uint(i/2), runner.getCount())
		}
	})
}

func TestFinalizePulse_HappyPath(t *testing.T) {
	ctx := inslogger.TestContext(t)

	testPulse := insolar.PulseNumber(pulse.MinTimePulse)
	targetPulse := testPulse + 1

	pc := insolarPulse.NewCalculatorMock(t)
	pc.ForwardsMock.Return(insolar.Pulse{PulseNumber: targetPulse}, nil)

	bkp := executor.NewBackupMakerMock(t)
	bkp.MakeBackupMock.Return(nil)

	jk := executor.NewJetKeeperMock(t)
	var hasConfirmCount uint32
	hasConfirm := func(ctx context.Context, pulse insolar.PulseNumber) bool {
		var p bool
		switch hasConfirmCount {
		case 0:
			p = true
		case 1:
			p = false
		}
		hasConfirmCount++
		return p
	}

	jk.HasAllJetConfirmsMock.Set(hasConfirm)

	js := jet.NewStorageMock(t)
	js.AllMock.Return(nil)
	jk.StorageMock.Return(js)

	var topSyncCount uint32
	topSync := func() insolar.PulseNumber {
		var p insolar.PulseNumber
		switch topSyncCount {
		case 0:
			p = testPulse
		case 1:
			p = targetPulse
		}
		topSyncCount++
		return p
	}

	jk.TopSyncPulseMock.Set(topSync)
	jk.AddBackupConfirmationMock.Return(nil)

	indexes := object.NewIndexModifierMock(t)
	indexes.UpdateLastKnownPulseMock.Return(nil)

	executor.FinalizePulse(ctx, pc, bkp, jk, indexes, targetPulse, testBadgerGCInfo())
}

func testBadgerGCInfo() *executor.BadgerGCRunInfo {
	return executor.NewBadgerGCRunInfo(&TestBadgerGCRunner{}, 1)
}

func TestFinalizePulse_JetIsNotConfirmed(t *testing.T) {
	ctx := inslogger.TestContext(t)

	testPulse := insolar.PulseNumber(pulse.MinTimePulse)

	jk := executor.NewJetKeeperMock(t)
	jk.HasAllJetConfirmsMock.Return(false)

	executor.FinalizePulse(ctx, nil, nil, jk, nil, testPulse, testBadgerGCInfo())
}

func TestFinalizePulse_CantGteNextPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)

	testPulse := insolar.PulseNumber(pulse.MinTimePulse)

	jk := executor.NewJetKeeperMock(t)
	jk.HasAllJetConfirmsMock.Return(true)
	jk.TopSyncPulseMock.Return(testPulse)

	pc := insolarPulse.NewCalculatorMock(t)
	pc.ForwardsMock.Return(insolar.Pulse{}, errors.New("Test"))

	executor.FinalizePulse(ctx, pc, nil, jk, nil, testPulse, testBadgerGCInfo())
}

func TestFinalizePulse_BackupError(t *testing.T) {
	ctx := inslogger.TestContext(t)

	testPulse := insolar.PulseNumber(pulse.MinTimePulse)
	targetPulse := testPulse + 1

	jk := executor.NewJetKeeperMock(t)
	jk.HasAllJetConfirmsMock.Return(true)
	jk.TopSyncPulseMock.Return(targetPulse)

	js := jet.NewStorageMock(t)
	js.AllMock.Return(nil)
	jk.StorageMock.Return(js)

	pc := insolarPulse.NewCalculatorMock(t)
	pc.ForwardsMock.Return(insolar.Pulse{PulseNumber: targetPulse}, nil)

	bkp := executor.NewBackupMakerMock(t)
	bkp.MakeBackupMock.Return(executor.ErrAlreadyDone)

	executor.FinalizePulse(ctx, pc, bkp, jk, nil, targetPulse, testBadgerGCInfo())
}

func TestFinalizePulse_NotNextPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)

	testPulse := insolar.PulseNumber(pulse.MinTimePulse)

	jk := executor.NewJetKeeperMock(t)
	jk.HasAllJetConfirmsMock.Return(true)
	jk.TopSyncPulseMock.Return(testPulse)

	pc := insolarPulse.NewCalculatorMock(t)
	pc.ForwardsMock.Return(insolar.Pulse{PulseNumber: testPulse}, nil)

	executor.FinalizePulse(ctx, pc, nil, jk, nil, testPulse+10, testBadgerGCInfo())
}
