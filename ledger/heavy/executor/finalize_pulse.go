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
	"sync"
	"time"

	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/object"
)

type BadgerGCRunner interface {
	// RunValueGC run badger values garbage collection
	RunValueGC(ctx context.Context)
}

type BadgerGCRunInfo struct {
	runner BadgerGCRunner
	// runFrequency is period of running gc (in number of pulses)
	runFrequency uint

	callCounter uint
	tryLock     chan struct{}
}

func NewBadgerGCRunInfo(runner BadgerGCRunner, runFrequency uint) *BadgerGCRunInfo {
	tryLock := make(chan struct{}, 1)
	tryLock <- struct{}{}
	return &BadgerGCRunInfo{
		runner:       runner,
		runFrequency: runFrequency,
		tryLock:      tryLock,
	}
}

func (b *BadgerGCRunInfo) RunGCIfNeeded(ctx context.Context) (doneWaiter <-chan struct{}) {
	done := make(chan struct{}, 1)
	go func() {
		defer func() {
			done <- struct{}{}
		}()
		select {
		case v := <-b.tryLock:
			b.callCounter++
			if (b.runFrequency > 0) && (b.callCounter >= b.runFrequency) && (b.callCounter%b.runFrequency == 0) {
				startedAt := time.Now().Second()
				b.runner.RunValueGC(ctx)
				stats.Record(ctx, statBadgerValueGCTime.M(int64(time.Now().Second()-startedAt)))
			} else {
				inslogger.FromContext(ctx).Info("values GC is not called")
			}
			b.tryLock <- v
		default:
			inslogger.FromContext(ctx).Info("values GC in progress. Skip It")
		}
	}()

	return done
}

func shouldStartFinalization(ctx context.Context, jetKeeper JetKeeper, pulses pulse.Calculator, pulseToFinalize insolar.PulseNumber) bool {
	logger := inslogger.FromContext(ctx)
	if !jetKeeper.HasAllJetConfirms(ctx, pulseToFinalize) {
		logger.Debug("not all jets confirmed. Do nothing. Pulse: ", pulseToFinalize)
		return false
	}

	nextTop, err := pulses.Forwards(ctx, jetKeeper.TopSyncPulse(), 1)
	if err != nil {
		logger.Warn("Can't get next pulse for topSynk: ", jetKeeper.TopSyncPulse())
		return false
	}

	if !nextTop.PulseNumber.Equal(pulseToFinalize) {
		logger.Infof("Try to finalize not sequential pulse. newTop: %d, target: %d", nextTop.PulseNumber, pulseToFinalize)
		return false
	}

	return true
}

// FinalizePulse starts backup process if needed
func FinalizePulse(ctx context.Context, pulses pulse.Calculator, backuper BackupMaker, jetKeeper JetKeeper, indexes object.IndexModifier, newPulse insolar.PulseNumber, gcRunner *BadgerGCRunInfo) {
	finPulse := &newPulse
	for {
		finPulse = finalizePulseStep(ctx, pulses, backuper, jetKeeper, indexes, *finPulse, gcRunner)
		if finPulse == nil {
			break
		}
	}
}

var finalizationLock sync.Mutex

func finalizePulseStep(ctx context.Context, pulses pulse.Calculator, backuper BackupMaker, jetKeeper JetKeeper, indexes object.IndexModifier, pulseToFinalize insolar.PulseNumber, gcRunner *BadgerGCRunInfo) *insolar.PulseNumber {
	logger := inslogger.FromContext(ctx)
	if !shouldStartFinalization(ctx, jetKeeper, pulses, pulseToFinalize) {
		logger.Info("Skip finalization. pulse: ", pulseToFinalize)
		return nil
	}

	// record all jets count
	stats.Record(ctx, statJets.M(int64(len(jetKeeper.Storage().All(ctx, pulseToFinalize)))))

	logger.Debug("FinalizePulse starts. pulse: ", pulseToFinalize)
	startedAt := time.Now().Second()
	bkpError := backuper.MakeBackup(ctx, pulseToFinalize)
	if bkpError != nil && bkpError != ErrAlreadyDone && bkpError != ErrBackupDisabled {
		logger.Fatal("Can't do backup: " + bkpError.Error())
	}
	stats.Record(ctx, statBackupTime.M(int64(time.Now().Second()-startedAt)))

	if bkpError == ErrAlreadyDone {
		logger.Info("Pulse already backuped: ", pulseToFinalize, bkpError)
		return nil
	}

	logger.Debug("FinalizePulse: before getting lock. pulse: ", pulseToFinalize)
	finalizationLock.Lock()
	defer finalizationLock.Unlock()
	logger.Debug("FinalizePulse: after getting lock. pulse: ", pulseToFinalize)

	err := jetKeeper.AddBackupConfirmation(ctx, pulseToFinalize)
	if err != nil {
		logger.Fatal("Can't add backup confirmation: " + err.Error())
	}

	newTopSyncPulse := jetKeeper.TopSyncPulse()

	if pulseToFinalize != newTopSyncPulse {
		logger.Fatal("Pulse has not been changed after adding backup confirmation. newTopSyncPulse: ", newTopSyncPulse, ", pulseToFinalize: ", pulseToFinalize)
	}
	if err := indexes.UpdateLastKnownPulse(ctx, newTopSyncPulse); err != nil {
		logger.Fatal("Can't update indexes for last sync pulse: ", err)
	}

	inslogger.FromContext(ctx).Infof("Pulse %d completely finalized ( drops + hots + backup )", pulseToFinalize)
	stats.Record(ctx, statFinalizedPulse.M(int64(pulseToFinalize)))

	// We run value GC here ( and only here ) implicitly since we want to
	// exclude running GC during process of backup-replication
	// Skip return value - we don't want to wait completion
	_ = gcRunner.RunGCIfNeeded(ctx)

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
