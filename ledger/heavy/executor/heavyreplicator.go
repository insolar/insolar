// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package executor

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

const (
	abandonedNotifyThreshold = 2
)

// HeavyReplicator is a base interface for a heavy sync component.
type HeavyReplicator interface {
	// NotifyAboutMessage is method for notifying a sync component about new data.
	NotifyAboutMessage(context.Context, *payload.Replication)

	// Stop stops the component.
	Stop()
}

// HeavyReplicatorDefault is a base impl for HeavyReplicator
type HeavyReplicatorDefault struct {
	txManager       object.TxManager
	records         object.RecordModifier
	indexes         object.IndexModifier
	pcs             insolar.PlatformCryptographyScheme
	pulseCalculator pulse.Calculator
	drops           drop.Modifier
	keeper          JetKeeper
	backuper        BackupMaker
	jets            jet.Modifier
	gcRunner        GCRunInfo
}

// NewHeavyReplicatorDefault creates new instance of HeavyReplicatorDefault.
func NewHeavyReplicatorDefault(
	txManager object.TxManager,
	records object.RecordModifier,
	indexes object.IndexModifier,
	pcs insolar.PlatformCryptographyScheme,
	pulseCalculator pulse.Calculator,
	drops drop.Modifier,
	keeper JetKeeper,
	backuper BackupMaker,
	jets jet.Modifier,
	gcRunner GCRunInfo,
) *HeavyReplicatorDefault {
	return &HeavyReplicatorDefault{
		txManager:       txManager,
		records:         records,
		indexes:         indexes,
		pcs:             pcs,
		pulseCalculator: pulseCalculator,
		drops:           drops,
		keeper:          keeper,
		backuper:        backuper,
		jets:            jets,
		gcRunner:        gcRunner,
	}
}

// Stop stops the component.
func (h *HeavyReplicatorDefault) Stop() {
	// do nothing
}

// NotifyAboutMessage is method for notifying a sync component about new data.
func (h *HeavyReplicatorDefault) NotifyAboutMessage(ctx context.Context, msg *payload.Replication) {
	// TODO: do all of this in a single transaction with retries!
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"jet_id": msg.JetID.DebugString(),
		"pulse":  msg.Pulse,
	})
	logger.Info("heavy replicator got a new message")

	logger.Debug("heavy replicator storing records")
	if err := storeRecords(ctx, h.records, h.pcs, msg.Pulse, msg.Records); err != nil {
		logger.Panic(errors.Wrap(err, "heavy replicator failed to store records"))
	}

	abandonedNotifyPulse, err := h.pulseCalculator.Backwards(ctx, msg.Pulse, abandonedNotifyThreshold)
	if err != nil {
		if err == pulse.ErrNotFound {
			abandonedNotifyPulse = *insolar.GenesisPulse
		} else {
			logger.Panic(errors.Wrap(err, "failed to calculate pending notify pulse"))
		}
	}

	logger.Debug("heavy replicator storing indexes")
	if err := storeIndexes(ctx, h.indexes, msg.Indexes, msg.Pulse, abandonedNotifyPulse.PulseNumber); err != nil {
		logger.Panic(errors.Wrap(err, "heavy replicator failed to store indexes"))
	}

	logger.Debug("heavy replicator storing drop")
	err = storeDrop(ctx, h.drops, msg.Drop)
	if err != nil {
		logger.Panic(errors.Wrap(err, "heavy replicator failed to store drop"))
	}
	logger = logger.WithField("drop_pulse", msg.Drop.Pulse)

	logger.Debug("heavy replicator storing drop confirmation")
	if err := h.keeper.AddDropConfirmation(ctx, msg.Drop.Pulse, msg.Drop.JetID, msg.Drop.Split); err != nil {
		logger.Panic(errors.Wrapf(err, "heavy replicator failed to add drop confirmation jet=%v", msg.Drop.JetID.DebugString()))
	}

	logger.Debug("heavy replicator update jets")
	err = h.jets.Update(ctx, msg.Drop.Pulse, true, msg.Drop.JetID)
	if err != nil {
		logger.Panic(errors.Wrapf(err, "heavy replicator failed to update jet %s", msg.Drop.JetID.DebugString()))
	}

	logger.Debug("heavy replicator finalize pulse")
	FinalizePulse(ctx, h.pulseCalculator, h.backuper, h.keeper, h.indexes, msg.Drop.Pulse, h.gcRunner)

	logger.Info("heavy replicator stops replication")
}

func storeIndexes(
	ctx context.Context,
	mod object.IndexModifier,
	indexes []record.Index,
	pn insolar.PulseNumber,
	abandonedNotifyPulse insolar.PulseNumber,
) error {
	for _, idx := range indexes {
		if idx.Lifeline.EarliestOpenRequest != nil && *idx.Lifeline.EarliestOpenRequest < abandonedNotifyPulse {
			stats.Record(ctx, statAbandonedRequests.M(1))
		}
		err := mod.SetIndex(ctx, pn, idx)
		if err != nil {
			return err
		}
	}
	return nil
}

func storeDrop(
	ctx context.Context,
	drops drop.Modifier,
	drop drop.Drop,
) error {
	err := drops.Set(ctx, drop)
	if err != nil {
		return err
	}

	return nil
}

func storeRecords(
	ctx context.Context,
	recordStorage object.RecordModifier,
	pcs insolar.PlatformCryptographyScheme,
	pn insolar.PulseNumber,
	records []record.Material,
) error {
	for _, rec := range records {
		hash := record.HashVirtual(pcs.ReferenceHasher(), rec.Virtual)
		id := *insolar.NewID(pn, hash)
		if rec.ID != id {
			return fmt.Errorf(
				"record id does not match (calculated: %s, received: %s)",
				id.DebugString(),
				rec.ID.DebugString(),
			)
		}
	}
	if err := recordStorage.BatchSet(ctx, records); err != nil {
		return errors.Wrap(err, "set method failed")
	}
	return nil
}
