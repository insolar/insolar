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

package executor

import (
	"context"

	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/object"
)

func shouldStartFinalization(ctx context.Context, jetKeeper JetKeeper, pulses pulse.Calculator, newPulse insolar.PulseNumber) bool {
	logger := inslogger.FromContext(ctx)
	if !jetKeeper.HasAllJetConfirms(ctx, newPulse) {
		logger.Debug("not all jets confirmed. Do nothing. Pulse: ", newPulse)
		return false
	}

	nextTop, err := pulses.Forwards(ctx, jetKeeper.TopSyncPulse(), 1)
	if err != nil {
		logger.Warn("Can't get next pulse for topSynk: ", jetKeeper.TopSyncPulse())
		return false
	}

	if !nextTop.PulseNumber.Equal(newPulse) {
		logger.Infof("Try to finalize not sequential pulse. newTop: %d, target: %d", nextTop.PulseNumber, newPulse)
		return false
	}

	return true
}

// FinalizePulse starts backup process if needed
func FinalizePulse(ctx context.Context, pulses pulse.Calculator, backuper BackupMaker, jetKeeper JetKeeper, indexes object.IndexModifier, newPulse insolar.PulseNumber) {
	finPulse := &newPulse
	for {
		finPulse = finalizePulseStep(ctx, pulses, backuper, jetKeeper, indexes, *finPulse)
		if finPulse == nil {
			break
		}
	}
}

func finalizePulseStep(ctx context.Context, pulses pulse.Calculator, backuper BackupMaker, jetKeeper JetKeeper, indexes object.IndexModifier, newPulse insolar.PulseNumber) *insolar.PulseNumber {
	logger := inslogger.FromContext(ctx)
	if !shouldStartFinalization(ctx, jetKeeper, pulses, newPulse) {
		logger.Info("Skip finalization")
		return nil
	}

	// record all jets count
	stats.Record(ctx, statJets.M(int64(len(jetKeeper.Storage().All(ctx, newPulse)))))

	logger.Debug("FinalizePulse starts")
	bkpError := backuper.MakeBackup(ctx, newPulse)
	if bkpError != nil && bkpError != ErrAlreadyDone && bkpError != ErrBackupDisabled {
		logger.Fatal("Can't do backup: " + bkpError.Error())
	}

	if bkpError == ErrAlreadyDone {
		logger.Info("Pulse already backuped: ", newPulse, bkpError)
		return nil
	}

	err := jetKeeper.AddBackupConfirmation(ctx, newPulse)
	if err != nil {
		logger.Fatal("Can't add backup confirmation: " + err.Error())
	}

	newTopSyncPulse := jetKeeper.TopSyncPulse()

	if newPulse != newTopSyncPulse {
		logger.Fatal("Pulse has not been changed after adding backup confirmation. newTopSyncPulse: ", newTopSyncPulse, ", newPulse: ", newPulse)
	}
	if err := indexes.UpdateLastKnownPulse(ctx, newTopSyncPulse); err != nil {
		logger.Fatal("Can't update indexes for last sync pulse: ", err)
	}

	inslogger.FromContext(ctx).Infof("Pulse %d completely finalized ( drops + hots + backup )", newPulse)
	stats.Record(ctx, statFinalizedPulse.M(int64(newPulse)))

	nextTop, err := pulses.Forwards(ctx, newTopSyncPulse, 1)
	if err != nil && err != pulse.ErrNotFound {
		logger.Fatal("pulses.Forwards topSynс: " + newTopSyncPulse.String())
	}
	if err == pulse.ErrNotFound {
		logger.Info("Stop propagating of backups")
		return nil
	}
	logger.Info("Propagating finalization to next pulse: ", nextTop.PulseNumber)

	pulseCopy := nextTop.PulseNumber
	return &pulseCopy
}
