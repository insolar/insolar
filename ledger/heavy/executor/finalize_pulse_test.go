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
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/executor"
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
	jk.HasAllJetConfirmsMock.Return(true)
	jk.TopSyncPulseMock.Return(testPulse)
	jk.AddBackupConfirmationMock.Return(nil)

	executor.FinalizePulse(ctx, pc, bkp, jk, targetPulse)
}

func TestFinalizePulse_JetIsNotConfirmed(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	testPulse := insolar.PulseNumber(insolar.FirstPulseNumber)

	jk := executor.NewJetKeeperMock(mc)
	jk.HasAllJetConfirmsMock.Return(false)

	executor.FinalizePulse(ctx, nil, nil, jk, testPulse)
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

	executor.FinalizePulse(ctx, pc, nil, jk, testPulse)
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

	executor.FinalizePulse(ctx, pc, bkp, jk, targetPulse)

}
