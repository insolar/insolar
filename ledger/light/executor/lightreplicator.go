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
	"fmt"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"go.opencensus.io/trace"
)

// LightReplicator is a base interface for a sync component
type LightReplicator interface {
	// NotifyAboutPulse is method for notifying a sync component about new pulse
	NotifyAboutPulse(ctx context.Context, pn insolar.PulseNumber)

	// Stop stops the component
	Stop()
}

// LightReplicatorDefault is a base impl of LightReplicator
type LightReplicatorDefault struct {
	once sync.Once
	done chan struct{}

	jetCalculator   JetCalculator
	cleaner         Cleaner
	sender          bus.Sender
	pulseCalculator pulse.Calculator

	dropAccessor drop.Accessor
	recsAccessor object.RecordCollectionAccessor
	idxAccessor  object.IndexAccessor
	jetAccessor  jet.Accessor

	syncWaitingPulses chan insolar.PulseNumber
}

// NewReplicatorDefault creates new instance of LightReplicator
func NewReplicatorDefault(
	jetCalculator JetCalculator,
	cleaner Cleaner,
	sender bus.Sender,
	calculator pulse.Calculator,
	dropAccessor drop.Accessor,
	recsAccessor object.RecordCollectionAccessor,
	idxAccessor object.IndexAccessor,
	jetAccessor jet.Accessor,
) *LightReplicatorDefault {
	return &LightReplicatorDefault{
		jetCalculator:   jetCalculator,
		cleaner:         cleaner,
		sender:          sender,
		pulseCalculator: calculator,

		dropAccessor: dropAccessor,
		recsAccessor: recsAccessor,
		idxAccessor:  idxAccessor,
		jetAccessor:  jetAccessor,

		syncWaitingPulses: make(chan insolar.PulseNumber),
		done:              make(chan struct{}),
	}
}

// NotifyAboutPulse is method for notifying a sync component about new pulse
// When it's called, a provided pulse is added to a channel.
// There is a special gorutine that is reading that channel. When a new pulse is being received,
// the routine starts to gather data (with using of LightDataGatherer). After gathering all the data,
// it attempts to send it to the heavy. After sending a heavy payload to a heavy, data is deleted
// with help of Cleaner
func (lr *LightReplicatorDefault) NotifyAboutPulse(ctx context.Context, pn insolar.PulseNumber) {
	lr.once.Do(func() {
		go lr.sync(context.Background())
	})

	logger := inslogger.FromContext(ctx)
	logger.Debugf("[Replicator][NotifyAboutPulse] received pulse - %v", pn)

	prevPN, err := lr.pulseCalculator.Backwards(ctx, pn, 1)
	if err != nil {
		logger.Error("[Replicator][NotifyAboutPulse]", err)
		return
	}

	logger.Debugf("[Replicator][NotifyAboutPulse] start replication, pulse - %v", prevPN.PulseNumber)
	lr.syncWaitingPulses <- prevPN.PulseNumber
}

func (lr *LightReplicatorDefault) Stop() {
	close(lr.done)
}

func (lr *LightReplicatorDefault) sync(ctx context.Context) {
	work := func(pn insolar.PulseNumber) {
		ctx, logger := inslogger.WithTraceField(ctx, utils.RandTraceID())
		logger.Debugf("[Replicator][sync] pn received - %v", pn)

		ctx, span := instracer.StartSpan(ctx, "LightReplicatorDefault.sync")
		span.AddAttributes(
			trace.Int64Attribute("pulse", int64(pn)),
		)
		defer span.End()

		allIndexes := lr.filterAndGroupIndexes(ctx, pn)
		jets, err := lr.jetCalculator.MineForPulse(ctx, pn)
		if err != nil {
			panic(errors.Wrap(err, "failed to calculate jets to sync"))
		}
		logger.Debugf("[Replicator][sync] founds %d jets", len(jets), ". Jets: ", insolar.JetIDCollection(jets).DebugString())

		for _, jetID := range jets {
			msg, err := lr.heavyPayload(ctx, pn, jetID, allIndexes[jetID])
			if err != nil {
				instracer.AddError(span, err)
				panic(
					fmt.Sprintf(
						"[Replicator][sync] Problems with gather data for a pulse - %v and jet - %v. err - %v",
						pn,
						jetID.DebugString(),
						err,
					),
				)
			}
			err = lr.sendToHeavy(ctx, msg)
			if err != nil {
				instracer.AddError(span, err)
				logger.Fatalf("[Replicator][sync] Problem with sending payload to a heavy node", err)
			} else {
				logger.Debugf("[Replicator][sync] Data has been sent to a heavy. pn - %v, jetID - %v", msg.Pulse, msg.JetID.DebugString())
			}
		}
		lr.cleaner.NotifyAboutPulse(ctx, pn)

		stats.Record(ctx, statLastReplicatedPulse.M(int64(pn)))
	}

	for {
		select {
		case pn, ok := <-lr.syncWaitingPulses:
			if !ok {
				return
			}
			work(pn)
		case <-lr.done:
			inslogger.FromContext(ctx).Info("light replicator stopped")
			return
		}
	}
}

func (lr *LightReplicatorDefault) sendToHeavy(ctx context.Context, pl payload.Replication) error {
	msg, err := payload.NewMessage(&pl)
	if err != nil {
		return err
	}

	inslogger.FromContext(ctx).Debug("send drop to heavy. pulse: ", pl.Pulse, ". jet: ", pl.JetID.DebugString())

	_, done := lr.sender.SendRole(ctx, msg, insolar.DynamicRoleHeavyExecutor, *insolar.NewReference(insolar.ID(pl.JetID)))
	done()

	return nil
}

func (lr *LightReplicatorDefault) filterAndGroupIndexes(
	ctx context.Context, pn insolar.PulseNumber,
) map[insolar.JetID][]record.Index {
	byJet := map[insolar.JetID][]record.Index{}
	indexes, err := lr.idxAccessor.ForPulse(ctx, pn)
	if err == nil {
		for _, idx := range indexes {
			jetID, _ := lr.jetAccessor.ForID(ctx, pn, idx.ObjID)
			byJet[jetID] = append(byJet[jetID], idx)
		}
	} else if err != object.ErrIndexNotFound {
		inslogger.FromContext(ctx).Errorf("Can't get indexes: %s", err)
	}
	return byJet
}

// ForPulseAndJet returns HeavyPayload message for a provided pulse and a jetID
func (lr *LightReplicatorDefault) heavyPayload(
	ctx context.Context,
	pn insolar.PulseNumber,
	jetID insolar.JetID,
	indexes []record.Index,
) (payload.Replication, error) {
	dr, err := lr.dropAccessor.ForPulse(ctx, jetID, pn)
	if err != nil {
		return payload.Replication{}, errors.Wrap(err, "failed to fetch drop")
	}

	records := lr.recsAccessor.ForPulse(ctx, jetID, pn)

	return payload.Replication{
		JetID:   jetID,
		Pulse:   pn,
		Indexes: indexes,
		Drop:    dr,
		Records: records,
	}, nil
}
