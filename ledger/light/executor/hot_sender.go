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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
)

// HotSender provides sending hot records send for provided pulse.
type HotSender interface {
	SendHot(ctx context.Context, old, new insolar.PulseNumber, jets []insolar.JetID) error
}

// HotSenderDefault implements HotSender.
type HotSenderDefault struct {
	dropAccessor    drop.Accessor
	indexAccessor   object.IndexAccessor
	pulseCalculator pulse.Calculator
	jetAccessor     jet.Accessor
	sender          bus.Sender

	// lightChainLimit is the LM-node cache limit configuration (how long index could be unused)
	lightChainLimit int
}

// NewHotSender returns a new instance of a default HotSender implementation.
func NewHotSender(
	dropAccessor drop.Accessor,
	indexAccessor object.IndexAccessor,
	pulseCalculator pulse.Calculator,
	jetAccessor jet.Accessor,
	lightChainLimit int,
	sender bus.Sender,
) *HotSenderDefault {
	return &HotSenderDefault{
		dropAccessor:    dropAccessor,
		indexAccessor:   indexAccessor,
		pulseCalculator: pulseCalculator,
		jetAccessor:     jetAccessor,
		sender:          sender,

		lightChainLimit: lightChainLimit,
	}
}

func (m *HotSenderDefault) filterAndGroupIndexes(
	ctx context.Context, currentPulse, newPulse insolar.PulseNumber,
) (map[insolar.JetID][]record.Index, error) {
	limitPN, err := m.pulseCalculator.Backwards(ctx, currentPulse, m.lightChainLimit)
	if err == pulse.ErrNotFound {
		limitPN = *insolar.GenesisPulse
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to fetch starting pulse for getting filaments")
	}

	byJet := map[insolar.JetID][]record.Index{}

	// filter out inactive indexes
	indexes, err := m.indexAccessor.ForPulse(ctx, currentPulse)
	if err == nil {

		// filtering in-place (optimization to avoid double allocation)
		filtered := indexes[:0]
		for _, idx := range indexes {
			if idx.LifelineLastUsed < limitPN.PulseNumber && idx.Lifeline.EarliestOpenRequest == nil {
				continue
			}
			filtered = append(filtered, record.Index{
				Lifeline:         idx.Lifeline,
				ObjID:            idx.ObjID,
				LifelineLastUsed: idx.LifelineLastUsed,
			})
		}

		for _, idx := range filtered {
			jetID, _ := m.jetAccessor.ForID(ctx, newPulse, idx.ObjID)
			byJet[jetID] = append(byJet[jetID], idx)
		}
	} else if err != object.ErrIndexNotFound {
		inslogger.FromContext(ctx).Errorf("Can't get indexes for pulse: %s", err)
	}

	return byJet, nil
}

// SendHot send hot records from oldPulse to all jets in newPulse.
func (m *HotSenderDefault) SendHot(
	ctx context.Context, currentPulse, newPulse insolar.PulseNumber, jets []insolar.JetID,
) error {
	ctx, span := instracer.StartSpan(ctx, "HotSenderDefault.SendHot")
	defer span.End()
	logger := inslogger.FromContext(ctx)

	idxByJet, err := m.filterAndGroupIndexes(ctx, currentPulse, newPulse)
	if err != nil {
		err = errors.Wrapf(err, "failed to get filament indexes for %v pulse", newPulse)
		instracer.AddError(span, err)
		return err
	}

	for _, id := range jets {
		jetID := id
		logger := logger.WithField("jetID", jetID.DebugString())

		block, err := m.findDrop(ctx, currentPulse, jetID)
		if err != nil {
			err = errors.Wrapf(err, "get drop for pulse %v and jet %v failed", currentPulse, jetID.DebugString())
			instracer.AddError(span, err)
			return err
		}
		logger.Infof("save drop for pulse %v", currentPulse)

		// send data for every jet asynchronously
		go func() {
			err := m.sendForJet(ctx, jetID, newPulse, idxByJet[jetID], block)
			if err != nil {
				logger.WithField("error", err.Error()).Error("hot sender: sendForJet failed")
			} else {
				logger.Info("hot sender: sendForJet OK")
			}
		}()
	}
	return nil
}

func (m *HotSenderDefault) sendForJet(
	ctx context.Context,
	jetID insolar.JetID,
	pn insolar.PulseNumber,
	indexes []record.Index,
	block drop.Drop,
) error {
	ctx, span := instracer.StartSpan(ctx, "hot_sender.send_hot")
	defer span.End()

	stats.Record(ctx, statHotObjectsTotal.M(int64(len(indexes))))

	buf := drop.MustEncode(&block)
	msg, err := payload.NewMessage(&payload.HotObjects{
		JetID:   jetID,
		Drop:    buf,
		Pulse:   pn,
		Indexes: indexes,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create message")
	}
	_, done := m.sender.SendRole(ctx, msg, insolar.DynamicRoleLightExecutor, *insolar.NewReference(insolar.ID(jetID)))
	done()

	stats.Record(ctx, statHotObjectsSend.M(int64(len(indexes))))
	return nil
}

// findDrop try to get drop for provided jet and if not found tries
// to find Parent's jet (if jet have been split and we have no previous drop for it by this reason)
func (m *HotSenderDefault) findDrop(
	ctx context.Context, pn insolar.PulseNumber, jetID insolar.JetID,
) (drop.Drop, error) {
	block, err := m.dropAccessor.ForPulse(ctx, jetID, pn)
	if err == drop.ErrNotFound {
		jetID = jet.Parent(jetID)
		// try to get parent's drop
		block, err = m.dropAccessor.ForPulse(ctx, jetID, pn)
		if err == drop.ErrNotFound {
			err = errors.Wrap(err, "drop for parent jet not found too")
		}
	}
	return block, err
}
