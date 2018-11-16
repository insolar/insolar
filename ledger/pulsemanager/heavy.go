package pulsemanager

import (
	"context"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
)

// HeavySync syncs records from light to heavy node, returns last synced pulse and error.
//
// It syncs records from start to end of provided pulse numbers.
func (m *PulseManager) HeavySync(
	ctx context.Context,
	start core.PulseNumber,
	end core.PulseNumber,
) (core.PulseNumber, error) {
	replicator := storage.NewReplicaIter(
		ctx, m.db, start, end, m.options.syncmessagelimit)
	for i := 0; ; i++ {
		recs, err := replicator.NextRecords()
		if err == storage.ErrReplicatorDone {
			break
		}
		if err != nil {
			panic(err)
		}
		msg := &message.HeavyRecords{Records: recs}
		reply, senderr := m.Bus.Send(ctx, msg)
		if senderr != nil {
			return core.PulseNumber(0), senderr
		}
		// TODO: check reply?
		_ = reply
	}
	inslogger.FromContext(ctx).Debugf(
		"synchronize on [%v:%v) finised (maximum record pulse is %v)",
		start, end, replicator.LastPulse())
	return replicator.LastPulse(), nil
}

// NextSyncPulses returns pulse numbers diapasone for syncing to heavy node.
// If nothing to sync it returns 0, 0, nil.
func (m *PulseManager) NextSyncPulses(ctx context.Context) (start, end core.PulseNumber, err error) {
	var (
		replicated core.PulseNumber
		last       core.PulseNumber
	)
	if replicated, err = m.db.GetReplicatedPulse(ctx); err != nil {
		return
	}
	if last, err = m.db.GetLastPulseAsLightMaterial(ctx); err != nil {
		return
	}
	// if replicated pulse is not less than "last light material pulse", return zero
	if !(replicated < last) {
		return
	}

	// start should be after replicated pulse or at least from FirstPulseNumber (for zero case)
	start = replicated + 1
	if replicated == 0 {
		start = core.FirstPulseNumber
	}
	// end should be after "last light material pulse" + 1 or at least next pulse after start
	end = last + 1
	if last == 0 {
		end = start + 1
	}
	return
}
