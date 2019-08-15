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

package executor_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/testutils/network"
	"github.com/pkg/errors"
)

func TestFinalizePulse_HappyPath(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	testPulse := insolar.PulseNumber(insolar.FirstPulseNumber)
	targetPulse := testPulse + 1

	pc := network.NewPulseCalculatorMock(mc)
	pc.ForwardsMock.Return(insolar.Pulse{PulseNumber: targetPulse}, nil)

	bkp := executor.NewBackupMakerMock(mc)
	bkp.MakeBackupMock.Return(nil)

	jk := executor.NewJetKeeperMock(mc)
	var hasConfirmCount uint32
	hasConfirm := func(ctx context.Context, pulse insolar.PulseNumber) bool {
		var p bool
		switch atomic.LoadUint32(&hasConfirmCount) {
		case 0:
			p = true
		case 1:
			p = false
		}
		atomic.AddUint32(&hasConfirmCount, 1)
		return p
	}

	jk.HasAllJetConfirmsMock.Set(hasConfirm)

	var topSyncCount uint32
	topSync := func() insolar.PulseNumber {
		var p insolar.PulseNumber
		switch atomic.LoadUint32(&topSyncCount) {
		case 0:
			p = testPulse
		case 1:
			p = targetPulse
		}
		atomic.AddUint32(&topSyncCount, 1)
		return p
	}

	jk.TopSyncPulseMock.Set(topSync)
	jk.AddBackupConfirmationMock.Return(nil)

	indexes := object.NewIndexModifierMock(mc)
	indexes.UpdateLastKnownPulseMock.Return(nil)

	executor.FinalizePulse(ctx, pc, bkp, jk, indexes, targetPulse)
}

func TestFinalizePulse_JetIsNotConfirmed(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	testPulse := insolar.PulseNumber(insolar.FirstPulseNumber)

	jk := executor.NewJetKeeperMock(mc)
	jk.HasAllJetConfirmsMock.Return(false)

	executor.FinalizePulse(ctx, nil, nil, jk, nil, testPulse)
}

func TestFinalizePulse_CantGteNextPulse(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	testPulse := insolar.PulseNumber(insolar.FirstPulseNumber)

	jk := executor.NewJetKeeperMock(mc)
	jk.HasAllJetConfirmsMock.Return(true)
	jk.TopSyncPulseMock.Return(testPulse)

	pc := network.NewPulseCalculatorMock(mc)
	pc.ForwardsMock.Return(insolar.Pulse{}, errors.New("Test"))

	executor.FinalizePulse(ctx, pc, nil, jk, nil, testPulse)
}

func TestFinalizePulse_BackupError(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	testPulse := insolar.PulseNumber(insolar.FirstPulseNumber)
	targetPulse := testPulse + 1

	jk := executor.NewJetKeeperMock(mc)
	jk.HasAllJetConfirmsMock.Return(true)
	jk.TopSyncPulseMock.Return(targetPulse)

	pc := network.NewPulseCalculatorMock(mc)
	pc.ForwardsMock.Return(insolar.Pulse{PulseNumber: targetPulse}, nil)

	bkp := executor.NewBackupMakerMock(mc)
	bkp.MakeBackupMock.Return(executor.ErrAlreadyDone)

	executor.FinalizePulse(ctx, pc, bkp, jk, nil, targetPulse)

}
