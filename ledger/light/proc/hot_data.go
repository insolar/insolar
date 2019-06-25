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

package proc

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/object"
)

type HotData struct {
	replyTo chan<- bus.Reply
	msg     *message.HotData

	Dep struct {
		DropModifier        drop.Modifier
		MessageBus          insolar.MessageBus
		IndexModifier        object.IndexModifier
		JetStorage           jet.Storage
		JetFetcher           jet.Fetcher
		JetReleaser          hot.JetReleaser
		Coordinator          jet.Coordinator

	}
}

func NewHotData(msg *message.HotData, replyTo chan<- bus.Reply) *HotData {
	return &HotData{
		msg:     msg,
		replyTo: replyTo,
	}
}

func (p *HotData) Proceed(ctx context.Context) error {
	err := p.process(ctx)
	if err != nil {
		p.replyTo <- bus.Reply{Err: err}
	}
	return err
}

func (p *HotData) process(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	jetID := insolar.JetID(*p.msg.Jet.Record())

	logger.WithFields(map[string]interface{}{
		"jet": jetID.DebugString(),
	}).Info("received hot data")

	err := p.Dep.DropModifier.Set(ctx, p.msg.Drop)
	if err == drop.ErrOverride {
		err = nil
	}
	if err != nil {
		return errors.Wrapf(err, "[jet]: drop error (pulse: %v)", p.msg.Drop.Pulse)
	}

	p.Dep.JetStorage.Update(
		ctx, p.msg.PulseNumber, true, jetID,
	)

	logger.Debugf("[handleHotRecords] received %v hot indexes", len(p.msg.HotIndexes))
	for _, meta := range p.msg.HotIndexes {
		decodedIndex, err := object.DecodeLifeline(meta.Index)
		if err != nil {
			logger.Error(err)
			continue
		}

		objJetID, _ := p.Dep.JetStorage.ForID(ctx, p.msg.PulseNumber, meta.ObjID)
		if objJetID != jetID {
			logger.Warn("received wrong id")
			continue
		}

		decodedIndex.JetID = jetID
		err = p.Dep.IndexModifier.SetIndex(
			ctx,
			p.msg.PulseNumber,
			object.FilamentIndex{
				ObjID:            meta.ObjID,
				Lifeline:         decodedIndex,
				LifelineLastUsed: meta.LifelineLastUsed,
				PendingRecords:   []insolar.ID{},
			},
		)
		if err != nil {
			logger.Error(errors.Wrapf(err, "[handleHotRecords] failed to save index - %v", meta.ObjID.DebugString()))
			continue
		}
		logger.Debugf("[handleHotRecords] lifeline with id - %v saved", meta.ObjID.DebugString())

		// go func(objID insolar.ID, pn insolar.PulseNumber) {
		// 	err = p.Dep.FilamentCacheManager.SendAbandonedNotification(ctx, pn, objID)
		// 	if err != nil {
		// 		logger.Errorf("failed to notify about abandoned notification %v", err)
		// 	}
		// }(meta.ObjID, flow.Pulse(ctx))
	}

	p.Dep.JetFetcher.Release(ctx, jetID, p.msg.PulseNumber)

	p.replyTo <- bus.Reply{Reply: &reply.OK{}}

	p.releaseHotDataWaiters(ctx)
	return nil
}

func (p *HotData) releaseHotDataWaiters(ctx context.Context) {
	jetID := p.msg.Jet.Record()
	err := p.Dep.JetReleaser.Unlock(ctx, *jetID)
	if err != nil {
		inslogger.FromContext(ctx).Error(err)
	}
}
