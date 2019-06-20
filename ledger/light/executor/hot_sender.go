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

func (m *HotSenderDefault) getFilamentIndexes(
	ctx context.Context, pn insolar.PulseNumber,
) ([]object.FilamentIndex, error) {
	bucks := m.indexBucketAccessor.ForPulse(ctx, pn)

	limitPN, err := m.pulseCalculator.Backwards(ctx, pn, m.lightChainLimit)
	if err == pulse.ErrNotFound {
		limitPN = *insolar.GenesisPulse
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to fetch starting pulse for getting filaments")
	}

	out := bucks[:0]
	for _, elem := range bucks {
		if elem.LifelineLastUsed < limitPN.PulseNumber {
			continue
		}
		out = append(out, elem)
	}
	return out, nil
}

func groupFilamentsByJet(filaments []object.FilamentIndex) map[insolar.JetID][]object.FilamentIndex {
	m := map[insolar.JetID][]object.FilamentIndex{}
	for _, f := range filaments {
		m[f.Lifeline.JetID] = append(m[f.Lifeline.JetID], f)
	}
	return m
}

func logErrorStr(ctx context.Context, jetID insolar.JetID, s string) {
	logger := inslogger.FromContext(ctx)
	logger.WithSkipFrameCount(1).WithFields(map[string]interface{}{
		"err":   s,
		"jetID": jetID.DebugString(),
	}).Error("failed to send hot data")
}

func (m *HotSenderDefault) sendForJet(
	ctx context.Context,
	info JetInfo,
	filaments []object.FilamentIndex,
	oldPulse, newPulse insolar.PulseNumber,
) error {
	block, err := m.createDrop(ctx, info, oldPulse)
	if err != nil {
		return errors.Wrapf(err, "create drop on pulse %v failed", oldPulse)
	}

	sender := func(hotIndexes []message.HotIndex, jetID insolar.JetID) {
		ctx, span := instracer.StartSpan(ctx, "hot_sender.send_hot")
		defer span.End()
		stats.Record(ctx, statHotObjectsTotal.M(int64(len(hotIndexes))))

		msg := &message.HotData{
			Jet:         *insolar.NewReference(insolar.ID(jetID)),
			Drop:        *block,
			HotIndexes:  hotIndexes,
			PulseNumber: newPulse,
		}

		genericRep, err := m.bus.Send(ctx, msg, nil)
		if err != nil {
			logErrorStr(ctx, jetID, err.Error())
			return
		}
		if _, ok := genericRep.(*reply.OK); !ok {
			logErrorStr(ctx, jetID, "failed to send hot data")
			return
		}

		logErrorStr(ctx, jetID, "test logErrorStr")
		stats.Record(ctx, statHotObjectsSend.M(int64(len(hotIndexes))))
	}

	if !info.SplitPerformed {
		hots := m.hotDataForJet(ctx, filaments, newPulse, nil)
		go sender(hots, info.ID)
		return nil
	}

	left, right := jet.Siblings(info.ID)
	hotsLeft := m.hotDataForJet(ctx, filaments, newPulse, &left)
	hotsRight := m.hotDataForJet(ctx, filaments, newPulse, &right)
	go sender(hotsLeft, left)
	go sender(hotsRight, right)
	return nil
}

func (m *HotSenderDefault) SendHot(
	ctx context.Context, jets []JetInfo, oldPulse, newPulse insolar.PulseNumber,
) error {
	ctx, span := instracer.StartSpan(ctx, "hot_sender.start")
	defer span.End()

	// get all indexes for old (current) pulse
	filaments, err := m.getFilamentIndexes(ctx, oldPulse)
	if err != nil {
		return errors.Wrapf(err, "failed to get filament indexes for %v pulse", newPulse)
	}
	filamentsByJet := groupFilamentsByJet(filaments)

	// process every jet asynchronously
	var eg errgroup.Group
	for _, info := range jets {
		info := info
		eg.Go(func() error {
			return m.sendForJet(ctx, info, filamentsByJet[info.ID], oldPulse, newPulse)
		})
	}
	return errors.Wrap(eg.Wait(), "got error on jets sync")
}

func (m *HotSenderDefault) createDrop(ctx context.Context, info JetInfo, p insolar.PulseNumber) (
	block *drop.Drop,
	err error,
) {
	block = &drop.Drop{
		Pulse: p,
		JetID: info.ID,
		Split: info.SplitIntent,
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
	filaments []object.FilamentIndex,
	pn insolar.PulseNumber,
	jetID *insolar.JetID,
) []message.HotIndex {
	filtered := 0
	hotIndexes := make([]message.HotIndex, 0, len(filaments))
	for _, meta := range filaments {
		encoded, err := meta.Lifeline.Marshal()
		if err != nil {
			inslogger.FromContext(ctx).Errorf("failed to marshal lifeline: %v", err)
			continue
		}

		// skip checks if not needed
		if jetID != nil {
			objJetID, _ := m.jetAccessor.ForID(ctx, pn, meta.ObjID)
			if objJetID != *jetID {
				filtered++
				continue
			}
		}

		hotIndexes = append(hotIndexes, message.HotIndex{
			LifelineLastUsed: meta.LifelineLastUsed,
			ObjID:            meta.ObjID,
			Index:            encoded,
		})
	}
	// inslogger.FromContext(ctx).Debugf("hot indexes filtered %v/%v", filtered, len(filaments))
	return hotIndexes
}
