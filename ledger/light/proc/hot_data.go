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

	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
)

type HotData struct {
	replyTo chan<- bus.Reply
	msg     *message.HotData
	/*
		jet     insolar.JetID

		pulse   insolar.PulseNumber
		idx     object.Lifeline
	*/

	Dep struct {
		/*
			IDLocker                   object.IDLocker
			IndexStorage               object.IndexStorage
			JetCoordinator             jet.Coordinator
			RecordModifier             object.RecordModifier
			IndexStateModifier         object.ExtendedIndexModifier
			PlatformCryptographyScheme insolar.PlatformCryptographyScheme
		*/
	}
}

func NewHotRecords(msg *message.HotData, replyTo chan<- bus.Reply) *HotData {
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

/***
func (h *MessageHandler) handleHotRecords(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	logger := inslogger.FromContext(ctx)

	msg := parcel.Message().(*message.HotData)
	jetID := insolar.JetID(*msg.Jet.Record())

	logger.WithFields(map[string]interface{}{
		"jet": jetID.DebugString(),
	}).Info("received hot data")

	err := h.DropModifier.Set(ctx, msg.Drop)
	if err == drop.ErrOverride {
		err = nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, "[jet]: drop error (pulse: %v)", msg.Drop.Pulse)
	}

	pendingStorage := h.RecentStorageProvider.GetPendingStorage(ctx, insolar.ID(jetID))
	logger.Debugf("received %d pending requests", len(msg.PendingRequests))

	var notificationList []insolar.ID
	for objID, objContext := range msg.PendingRequests {
		if !objContext.Active {
			notificationList = append(notificationList, objID)
		}

		objContext.Active = false
		pendingStorage.SetContextToObject(ctx, objID, objContext)
	}

	go func() {
		for _, objID := range notificationList {
			go func(objID insolar.ID) {
				rep, err := h.Bus.Send(ctx, &message.AbandonedRequestsNotification{
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

	for id, meta := range msg.HotIndexes {
		decodedIndex, err := object.DecodeIndex(meta.Index)
		if err != nil {
			logger.Error(err)
			continue
		}

		err = h.IndexStateModifier.SetWithMeta(ctx, id, meta.LastUsed, decodedIndex)
		if err != nil {
			logger.Error(err)
			continue
		}
	}

	h.JetStorage.Update(
		ctx, msg.PulseNumber, true, insolar.JetID(jetID),
	)

	h.jetTreeUpdater.Release(ctx, jetID, msg.PulseNumber)

	return &reply.OK{}, nil
}

***/

func (p *HotData) process(ctx context.Context) error {
	return nil
	// TODO
	/*
		r, err := object.DecodeVirtual(p.msg.Record)
		if err != nil {
			return errors.Wrap(err, "can't deserialize record")
		}
		childRec, ok := r.(*object.ChildRecord)
		if !ok {
			return errors.New("wrong child record")
		}

		p.Dep.IDLocker.Lock(p.msg.Parent.Record())
		defer p.Dep.IDLocker.Unlock(p.msg.Parent.Record())
		p.Dep.IndexStateModifier.SetUsageForPulse(ctx, *p.msg.Parent.Record(), p.pulse)
		recID := object.NewRecordIDFromRecord(p.Dep.PlatformCryptographyScheme, p.pulse, childRec)

		// Children exist and pointer does not match (preserving chain consistency).
		// For the case when vm can't save or send result to another vm and it tries to update the same record again
		if p.idx.ChildPointer != nil && !childRec.PrevChild.Equal(*p.idx.ChildPointer) && p.idx.ChildPointer != recID {
			return errors.New("invalid child record")
		}

		child := object.NewRecordIDFromRecord(p.Dep.PlatformCryptographyScheme, p.pulse, childRec)
		rec := record.MaterialRecord{
			Record: childRec,
			JetID:  p.jet,
		}

		err = p.Dep.RecordModifier.Set(ctx, *child, rec)

		if err == object.ErrOverride {
			inslogger.FromContext(ctx).WithField("type", fmt.Sprintf("%T", r)).Warn("set record override (#2)")
			child = recID
		} else if err != nil {
			return errors.Wrap(err, "can't save record into storage")
		}

		p.idx.ChildPointer = child
		if p.msg.AsType != nil {
			p.idx.Delegates[*p.msg.AsType] = p.msg.Child
		}
		p.idx.LatestUpdate = p.pulse
		p.idx.JetID = p.jet
		err = p.Dep.IndexStorage.Set(ctx, *p.msg.Parent.Record(), p.idx)
		if err != nil {
			return err
		}

		p.replyTo <- bus.Reply{Reply: &reply.ID{ID: *child}}
		return nil*/
}
