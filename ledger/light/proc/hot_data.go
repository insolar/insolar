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
	"github.com/insolar/insolar/ledger/light/recentstorage"
	"github.com/insolar/insolar/ledger/object"
)

type HotData struct {
	replyTo chan<- bus.Reply
	msg     *message.HotData

	Dep struct {
		DropModifier          drop.Modifier
		RecentStorageProvider recentstorage.Provider
		MessageBus            insolar.MessageBus
		IndexStateModifier    object.ExtendedIndexModifier
		JetStorage            jet.Storage
		JetFetcher            jet.Fetcher
		JetReleaser           hot.JetReleaser
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

	pendingStorage := p.Dep.RecentStorageProvider.GetPendingStorage(ctx, insolar.ID(jetID))
	logger.Debugf("received %d pending requests", len(p.msg.PendingRequests))

	var notificationList []insolar.ID
	for objID, objContext := range p.msg.PendingRequests {
		if !objContext.Active {
			notificationList = append(notificationList, objID)
		}

		objContext.Active = false
		pendingStorage.SetContextToObject(ctx, objID, objContext)
	}

	go func() {
		for _, objID := range notificationList {
			go func(objID insolar.ID) {
				rep, err := p.Dep.MessageBus.Send(ctx, &message.AbandonedRequestsNotification{
					Object: objID,
				}, nil)

				if err != nil {
					logger.Error("failed to notify about pending requests")
					return
				}
				if _, ok := rep.(*reply.OK); !ok {
					logger.Error("received unexpected reply on pending notification")
				}
			}(objID)
		}
	}()

	for id, meta := range p.msg.HotIndexes {
		decodedIndex, err := object.DecodeIndex(meta.Index)
		if err != nil {
			logger.Error(err)
			continue
		}

		err = p.Dep.IndexStateModifier.SetWithMeta(ctx, id, meta.LastUsed, decodedIndex)
		if err != nil {
			logger.Error(err)
			continue
		}
	}

	p.Dep.JetStorage.Update(
		ctx, p.msg.PulseNumber, true, jetID,
	)

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
