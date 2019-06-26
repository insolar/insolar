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
	"golang.org/x/sync/errgroup"
)

// HotSender provides sending hot records send for provided pulse.
type HotSender interface {
	SendHot(ctx context.Context, jets []JetInfo, old, new insolar.PulseNumber) error
}

// HotSenderDefault implements HotSender.
type HotSenderDefault struct {
	bus                 insolar.MessageBus
	dropModifier        drop.Modifier
	indexBucketAccessor object.IndexBucketAccessor
	pulseCalculator     pulse.Calculator
	jetAccessor         jet.Accessor

	// light limit configuration
	lightChainLimit int
}

// NewHotSender returns a new instance of a default HotSender implementation.
func NewHotSender(
	bus insolar.MessageBus,
	dropModifier drop.Modifier,
	indexBucketAccessor object.IndexBucketAccessor,
	pulseCalculator pulse.Calculator,
	jetAccessor jet.Accessor,
	lightChainLimit int,
) *HotSenderDefault {
	return &HotSenderDefault{
		bus:                 bus,
		dropModifier:        dropModifier,
		indexBucketAccessor: indexBucketAccessor,
		pulseCalculator:     pulseCalculator,
		jetAccessor:         jetAccessor,

		lightChainLimit: lightChainLimit,
	}
}

func (m *HotSenderDefault) filterAndGroupIndexes(
	ctx context.Context, oldPulse, newPulse insolar.PulseNumber,
) (map[insolar.JetID][]object.FilamentIndex, error) {
	limitPN, err := m.pulseCalculator.Backwards(ctx, oldPulse, m.lightChainLimit)
	if err == pulse.ErrNotFound {
		limitPN = *insolar.GenesisPulse
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to fetch starting pulse for getting filaments")
	}

	// filter out inactive indexes
	indexes := m.indexBucketAccessor.ForPulse(ctx, oldPulse)
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
	ctx context.Context, jets []JetInfo, oldPulse, newPulse insolar.PulseNumber,
) error {
	ctx, span := instracer.StartSpan(ctx, "hot_sender.start")
	defer span.End()
	logger := inslogger.FromContext(ctx)

	// get indexes grouped by actual jet ID
	byJet, err := m.filterAndGroupIndexes(ctx, oldPulse, newPulse)
	if err != nil {
		return errors.Wrapf(err, "failed to get filament indexes for %v pulse", newPulse)
	}

	// process every jet asynchronously
	var eg errgroup.Group
	for _, info := range jets {
		info := info
		eg.Go(func() error {
			err := m.sendForJet(ctx, info.ID, byJet[info.ID], oldPulse, newPulse, info.SplitIntent)
			if err != nil {
				logger.WithSkipFrameCount(1).WithFields(map[string]interface{}{
					"jetID": info.ID.DebugString(),
					"error": err.Error(),
				}).Error("hot sender: sendForJet failed")
			}
			return err
		})
	}
	return errors.Wrap(eg.Wait(), "got error on jets sync")
}

func (m *HotSenderDefault) sendForJet(
	ctx context.Context,
	jetID insolar.JetID,
	indexes []object.FilamentIndex,
	oldPulse, newPulse insolar.PulseNumber,
	splitIntent bool,
) error {
	ctx, span := instracer.StartSpan(ctx, "hot_sender.send_hot")
	defer span.End()

	block, err := m.createDrop(ctx, oldPulse, jetID, splitIntent)
	if err != nil {
		return errors.Wrapf(err, "create drop for pulse %v failed", oldPulse)
	}

	hots := m.hotDataForJet(ctx, indexes)
	stats.Record(ctx, statHotObjectsTotal.M(int64(len(hots))))

	msg := &message.HotData{
		Jet:         *insolar.NewReference(insolar.ID(jetID)),
		Drop:        *block,
		HotIndexes:  hots,
		PulseNumber: newPulse,
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

func (m *HotSenderDefault) createDrop(
	ctx context.Context,
	pn insolar.PulseNumber,
	jetID insolar.JetID,
	splitIntent bool,
) (
	block *drop.Drop,
	err error,
) {
	block = &drop.Drop{
		Pulse: pn,
		JetID: jetID,
		Split: splitIntent,
	}

	err = m.dropModifier.Set(ctx, *block)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to set drop %#v", block)
	}

	return block, nil
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
