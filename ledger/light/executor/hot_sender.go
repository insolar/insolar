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
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/reply"
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
	bus                 insolar.MessageBus
	dropAccessor        drop.Accessor
	indexBucketAccessor object.IndexBucketAccessor
	pulseCalculator     pulse.Calculator
	jetAccessor         jet.Accessor

	// lightChainLimit is the LM-node cache limit configuration (how long index could be unused)
	lightChainLimit int
}

// NewHotSender returns a new instance of a default HotSender implementation.
func NewHotSender(
	bus insolar.MessageBus,
	dropAccessor drop.Accessor,
	indexBucketAccessor object.IndexBucketAccessor,
	pulseCalculator pulse.Calculator,
	jetAccessor jet.Accessor,
	lightChainLimit int,
) *HotSenderDefault {
	return &HotSenderDefault{
		bus:                 bus,
		dropAccessor:        dropAccessor,
		indexBucketAccessor: indexBucketAccessor,
		pulseCalculator:     pulseCalculator,
		jetAccessor:         jetAccessor,

		lightChainLimit: lightChainLimit,
	}
}

func (m *HotSenderDefault) filterAndGroupIndexes(
	ctx context.Context, currentPulse, newPulse insolar.PulseNumber,
) (map[insolar.JetID][]object.FilamentIndex, error) {
	limitPN, err := m.pulseCalculator.Backwards(ctx, currentPulse, m.lightChainLimit)
	if err == pulse.ErrNotFound {
		limitPN = *insolar.GenesisPulse
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to fetch starting pulse for getting filaments")
	}

	// filter out inactive indexes
	indexes := m.indexBucketAccessor.ForPulse(ctx, currentPulse)
	// filtering in-place (optimization to avoid double allocation)
	filtered := indexes[:0]
	for _, idx := range indexes {
		if idx.LifelineLastUsed < limitPN.PulseNumber {
			continue
		}
		filtered = append(filtered, idx)
	}

	byJet := map[insolar.JetID][]object.FilamentIndex{}
	for _, idx := range filtered {
		jetID, _ := m.jetAccessor.ForID(ctx, newPulse, idx.ObjID)
		byJet[jetID] = append(byJet[jetID], idx)
	}
	return byJet, nil
}

// SendHot send hot records from oldPulse to all jets in newPulse.
func (m *HotSenderDefault) SendHot(
	ctx context.Context, currentPulse, newPulse insolar.PulseNumber, jets []insolar.JetID,
) error {
	ctx, span := instracer.StartSpan(ctx, "hot_sender.start")
	defer span.End()
	logger := inslogger.FromContext(ctx)

	byJet, err := m.filterAndGroupIndexes(ctx, currentPulse, newPulse)
	if err != nil {
		return errors.Wrapf(err, "failed to get filament indexes for %v pulse", newPulse)
	}

	for _, id := range jets {
		jetID := id
		logger := logger.WithSkipFrameCount(1).WithField("jetID", jetID.DebugString())

		block, err := m.findDrop(ctx, currentPulse, jetID)
		if err != nil {
			return errors.Wrapf(err, "get drop for pulse %v and jet %v failed", currentPulse, jetID.DebugString())
		}
		logger.Infof("save drop for pulse %v", currentPulse)

		// send data for every jet asynchronously
		go func() {
			logger.Infof("SPLIT> fire sendForJet in goroutine")
			err := m.sendForJet(ctx, jetID, newPulse, byJet[jetID], block)
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
	indexes []object.FilamentIndex,
	block drop.Drop,
) error {
	ctx, span := instracer.StartSpan(ctx, "hot_sender.send_hot")
	defer span.End()

	hots := m.hotDataForJet(ctx, indexes)
	stats.Record(ctx, statHotObjectsTotal.M(int64(len(hots))))

	msg := &message.HotData{
		Jet:         *insolar.NewReference(insolar.ID(jetID)),
		Drop:        block,
		HotIndexes:  hots,
		PulseNumber: pn,
	}

	genericRep, err := m.bus.Send(ctx, msg, nil)
	if err != nil {
		return errors.Wrapf(err, "failed to send hot data, method Send failed")
	}
	if _, ok := genericRep.(*reply.OK); !ok {
		return errors.Wrapf(err, "failed to send hot data, not OK reply")
	}

	stats.Record(ctx, statHotObjectsSend.M(int64(len(hots))))
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
		// try to get parent
		block, err = m.dropAccessor.ForPulse(ctx, jetID, pn)
		if err == drop.ErrNotFound {
			err = errors.Wrap(err, "drop for parent jet not found too")
		}
	}
	return block, err
}

// hotDataForJet prepares list of HotIndex struct from provided FilamentIndex struct.
// if jetID provided, filters output records what only match this jet.
func (m *HotSenderDefault) hotDataForJet(
	ctx context.Context,
	indexes []object.FilamentIndex,
) []message.HotIndex {
	hotIndexes := make([]message.HotIndex, 0, len(indexes))
	for _, meta := range indexes {
		encoded, err := meta.Lifeline.Marshal()
		if err != nil {
			inslogger.FromContext(ctx).Errorf("failed to marshal lifeline: %v", err)
			continue
		}

		hotIndexes = append(hotIndexes, message.HotIndex{
			LifelineLastUsed: meta.LifelineLastUsed,
			ObjID:            meta.ObjID,
			Index:            encoded,
		})
	}
	return hotIndexes
}
