// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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

type GCRunner interface {
	// RunValueGC run badger values garbage collection
	RunValueGC(ctx context.Context)
}

type PostgresGCRunInfo struct {
}

func (p *PostgresGCRunInfo) RunValueGC(ctx context.Context) {
}

func (p *PostgresGCRunInfo) RunGCIfNeeded(ctx context.Context) <-chan struct{} {
	c := make(chan struct{}, 1)
	c <- struct{}{}
	return c
}

type GCRunInfo interface {
	RunGCIfNeeded(ctx context.Context) <-chan struct{}
}

type BadgerGCRunInfo struct {
	runner GCRunner
	// runFrequency is period of running gc (in number of pulses)
	runFrequency uint

	callCounter uint
	tryLock     chan struct{}
}

func NewBadgerGCRunInfo(runner GCRunner, runFrequency uint) *BadgerGCRunInfo {
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

func shouldStartFinalization(ctx context.Context, logger insolar.Logger, jetKeeper JetKeeper, pulses pulse.Calculator, pulseToFinalize insolar.PulseNumber) bool {
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
func FinalizePulse(ctx context.Context, pulses pulse.Calculator, backuper BackupMaker, jetKeeper JetKeeper, indexes object.IndexModifier, newPulse insolar.PulseNumber, gcRunner GCRunInfo) {
	finPulse := &newPulse
	for {
		finPulse = finalizePulseStep(ctx, pulses, backuper, jetKeeper, indexes, *finPulse, gcRunner)
		if finPulse == nil {
			break
		}
	}
}

var finalizationLock sync.Mutex

func finalizePulseStep(ctx context.Context, pulses pulse.Calculator, backuper BackupMaker, jetKeeper JetKeeper, indexes object.IndexModifier, pulseToFinalize insolar.PulseNumber, gcRunner GCRunInfo) *insolar.PulseNumber {
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"pulse_to_finalize": pulseToFinalize,
	})

	logger.Info("finalizePulseStep: begin")

	if !shouldStartFinalization(ctx, logger, jetKeeper, pulses, pulseToFinalize) {
		logger.Info("finalizePulseStep: skip finalization")
		return nil
	}

	// record all jets count
	stats.Record(ctx, statJets.M(int64(len(jetKeeper.Storage().All(ctx, pulseToFinalize)))))

	if backuper != nil {
		// Badger backend is used and backups are enabled
		startedAt := time.Now()
		logger.Infof("finalizePulseStep: calling backuper.MakeBackup()...")
		bkpError := backuper.MakeBackup(ctx, pulseToFinalize)
		if bkpError != nil && bkpError != ErrAlreadyDone && bkpError != ErrBackupDisabled {
			logger.Fatal("finalizePulseStep: MakeBackup() failed: " + bkpError.Error())
		}
		logger.Infof("finalizePulseStep: MakeBackup() done!")
		stats.Record(ctx, statBackupTime.M(time.Since(startedAt).Nanoseconds()))

		if bkpError == ErrAlreadyDone {
			logger.Info("finalizePulseStep: pulse already backuped: ", pulseToFinalize, bkpError)
			return nil
		}
	}

	logger.Info("finalizePulseStep: before getting lock")
	finalizationLock.Lock()
	defer finalizationLock.Unlock()
	logger.Info("finalizePulseStep: lock acquired, calling AddBackupConfirmation()...")

	err := jetKeeper.AddBackupConfirmation(ctx, pulseToFinalize)
	if err != nil {
		logger.Fatal("finalizePulseStep: can't add backup confirmation: " + err.Error())
	}

	logger.Info("finalizePulseStep: AddBackupConfirmation() done, calling jetKeeper.TopSyncPulse()...")
	newTopSyncPulse := jetKeeper.TopSyncPulse()
	if pulseToFinalize != newTopSyncPulse {
		logger.Fatal("finalizePulseStep: pulse has not been changed after adding backup confirmation. newTopSyncPulse: ", newTopSyncPulse, ", pulseToFinalize: ", pulseToFinalize)
	}

	logger.Info("finalizePulseStep: jetKeeper.TopSyncPulse() done, calling indexes.UpdateLastKnownPulse()...")

	if err := indexes.UpdateLastKnownPulse(ctx, newTopSyncPulse); err != nil {
		logger.Fatal("finalizePulseStep: can't update indexes for last sync pulse: ", err)
	}

	logger.Infof("finalizePulseStep: pulse completely finalized ( drops + hots + backup )")
	stats.Record(ctx, statFinalizedPulse.M(int64(pulseToFinalize)))

	// We run value GC here ( and only here ) implicitly since we want to
	// exclude running GC during process of backup-replication
	// Skip return value - we don't want to wait completion
	_ = gcRunner.RunGCIfNeeded(ctx)

	nextTop, err := pulses.Forwards(ctx, newTopSyncPulse, 1)
	if err != nil && err != pulse.ErrNotFound {
		logger.Fatal("finalizePulseStep: pulses.Forwards topSynс: " + newTopSyncPulse.String())
	}
	if err == pulse.ErrNotFound {
		logger.Info("finalizePulseStep: done! Stop propagating of backups")
		return nil
	}

	logger.Info("finalizePulseStep: done! Propagating finalization to next pulse: ", nextTop.PulseNumber)
	pulseCopy := nextTop.PulseNumber
	return &pulseCopy
}
