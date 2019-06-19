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
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
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

	// light limit configuration
	lightChainLimit int
}

// NewHotSender returns a new instance of a default HotSender implementation.
func NewHotSender(
	bus insolar.MessageBus,
	dropModifier drop.Modifier,
	indexBucketAccessor object.IndexBucketAccessor,
	pulseCalculator pulse.Calculator,
	lightChainLimit int,
) *HotSenderDefault {
	return &HotSenderDefault{
		bus:                 bus,
		dropModifier:        dropModifier,
		indexBucketAccessor: indexBucketAccessor,
		pulseCalculator:     pulseCalculator,
		lightChainLimit:     lightChainLimit,
	}
}

func (m *HotSenderDefault) SendHot(
	ctx context.Context, jets []JetInfo, oldPulse, newPulse insolar.PulseNumber,
) error {
	var g errgroup.Group
	// ctx, span := instracer.StartSpan(ctx, "pulse.process_end")
	// defer span.End()

	// logger := inslogger.FromContext(ctx)
	for _, i := range jets {
		info := i

		g.Go(func() error {
			block, err := m.createDrop(ctx, info, oldPulse)
			if err != nil {
				return errors.Wrapf(err, "create drop on pulse %v failed", oldPulse)
			}

			sender := func(msg message.HotData, jetID insolar.JetID) {
				// ctx, span := instracer.StartSpan(ctx, "pulse.send_hot")
				// defer span.End()
				msg.Jet = *insolar.NewReference(insolar.ID(jetID))
				genericRep, err := m.bus.Send(ctx, &msg, nil)
				if err != nil {
					// logger.WithField("err", err).Error("failed to send hot data")
					return
				}
				if _, ok := genericRep.(*reply.OK); !ok {
					// logger.WithField(
					// 	"err",
					// 	fmt.Sprintf("unexpected reply: %T", genericRep),
					// ).Error("failed to send hot data")
					return
				}
			}

			if info.SplitPerformed {
				msg, err := m.executorHotData(
					ctx, info.ID, oldPulse, newPulse, block,
				)
				if err != nil {
					return errors.Wrapf(err, "getExecutorData failed for jet ID %v", info.ID)
				}
				// No Split happened.
				go sender(*msg, info.ID)
			} else {
				msg, err := m.executorHotData(ctx, info.ID, oldPulse, newPulse, block)
				if err != nil {
					return errors.Wrapf(err, "getExecutorData failed for jet ID %v", info.ID)
				}

				// Split happened.
				left, right := jet.Siblings(info.ID)
				go sender(*msg, left)
				go sender(*msg, right)
			}

			return nil
		})
	}
	err := g.Wait()
	if err != nil {
		return errors.Wrap(err, "got error on jets sync")
	}

	return nil
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

func (m *HotSenderDefault) executorHotData(
	ctx context.Context,
	jetID insolar.JetID,
	oldP insolar.PulseNumber,
	newP insolar.PulseNumber,
	drop *drop.Drop,
) (*message.HotData, error) {
	// ctx, span := instracer.StartSpan(ctx, "hot_sender.executorHotData")
	// defer span.End()

	bucks := m.indexBucketAccessor.ForPNAndJet(ctx, oldP, jetID)
	limitPN, err := m.pulseCalculator.Backwards(ctx, oldP, m.lightChainLimit)
	if err == pulse.ErrNotFound {
		limitPN = *insolar.GenesisPulse
	} else if err != nil {
		// inslogger.FromContext(ctx).Errorf("failed to fetch limit %v", err)
		return nil, err
	}

	hotIndexes := []message.HotIndex{}
	for _, meta := range bucks {
		encoded, err := meta.Lifeline.Marshal()
		if err != nil {
			// inslogger.FromContext(ctx).WithField("id", meta.ObjID.DebugString()).Error("failed to marshal lifeline")
			continue
		}
		if meta.LifelineLastUsed < limitPN.PulseNumber {
			continue
		}

		hotIndexes = append(hotIndexes, message.HotIndex{
			LifelineLastUsed: meta.LifelineLastUsed,
			ObjID:            meta.ObjID,
			Index:            encoded,
		})
	}

	// stats.Record(
	// 	ctx,
	// 	statHotObjectsSent.M(int64(len(hotIndexes))),
	// )

	msg := &message.HotData{
		Drop:        *drop,
		PulseNumber: newP,
		HotIndexes:  hotIndexes,
	}
	return msg, nil
}
