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
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/light/recentstorage"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type UpdateObject struct {
	JetID       insolar.JetID
	Message     *message.UpdateObject
	ReplyTo     chan<- bus.Reply
	PulseNumber insolar.PulseNumber

	Dep struct {
		RecordModifier        object.RecordModifier
		Bus                   insolar.MessageBus
		Coordinator           jet.Coordinator
		BlobModifier          blob.Modifier
		RecentStorageProvider recentstorage.Provider
		PCS                   insolar.PlatformCryptographyScheme
		IDLocker              object.IDLocker
		LifelineIndex         object.LifelineIndex
		LifelineStateModifier object.LifelineStateModifier
		WriteAccessor         hot.WriteAccessor
	}
}

func NewUpdateObject(jetID insolar.JetID, message *message.UpdateObject, pulseNumber insolar.PulseNumber, replyTo chan<- bus.Reply) *UpdateObject {
	return &UpdateObject{
		JetID:       jetID,
		Message:     message,
		ReplyTo:     replyTo,
		PulseNumber: pulseNumber,
	}
}

func (p *UpdateObject) Proceed(ctx context.Context) error {
	p.ReplyTo <- p.handle(ctx)
	return nil
}

func (p *UpdateObject) handle(ctx context.Context) bus.Reply {
	done, err := p.Dep.WriteAccessor.Begin(ctx, p.PulseNumber)
	if err == hot.ErrWriteClosed {
		return bus.Reply{Err: flow.ErrCancelled}
	}
	defer done()

	logger := inslogger.FromContext(ctx)
	if p.Message.Object.Record() == nil {
		return bus.Reply{
			Err: errors.New("updateObject message object is nil"),
		}
	}

	virtRec := record.Virtual{}
	err = virtRec.Unmarshal(p.Message.Record)
	if err != nil {
		return bus.Reply{Err: errors.Wrap(err, "can't deserialize record")}
	}
	concreteRec := record.Unwrap(&virtRec)
	state, ok := concreteRec.(record.State)
	if !ok {
		return bus.Reply{Err: errors.New("wrong object state record")}
	}

	calculatedID := object.CalculateIDForBlob(p.Dep.PCS, p.PulseNumber, p.Message.Memory)
	// FIXME: temporary fix. If we calculate blob id on the client, pulse can change before message sending and this
	//  id will not match the one calculated on the server.
	err = p.Dep.BlobModifier.Set(ctx, *calculatedID, blob.Blob{JetID: p.JetID, Value: p.Message.Memory})
	if err != nil && err != blob.ErrOverride {
		return bus.Reply{Err: errors.Wrap(err, "failed to set blob")}
	}

	switch s := state.(type) {
	case *record.Activate:
		s.Memory = *calculatedID
	case *record.Amend:
		s.Memory = *calculatedID
	}

	p.Dep.IDLocker.Lock(p.Message.Object.Record())
	defer p.Dep.IDLocker.Unlock(p.Message.Object.Record())

	idx, err := p.Dep.LifelineIndex.ForID(ctx, p.PulseNumber, *p.Message.Object.Record())
	// No index on our node.
	if err == object.ErrLifelineNotFound {
		if state.ID() == record.StateActivation {
			// We are activating the object. There is no index for it anywhere.
			idx = object.Lifeline{StateID: record.StateUndefined}
			logger.Debugf("new lifeline created")
		} else {
			logger.Debug("failed to fetch index (fetching from heavy)")
			// We are updating object. LifelineIndex should be on the heavy executor.
			heavy, err := p.Dep.Coordinator.Heavy(ctx, p.PulseNumber)
			if err != nil {
				return bus.Reply{Err: err}
			}
			idx, err = p.saveIndexFromHeavy(ctx, p.JetID, p.Message.Object, heavy)
			if err != nil {
				logger.WithFields(map[string]interface{}{
					"jet": p.JetID.DebugString(),
					"pn":  flow.Pulse(ctx),
				}).Error(errors.Wrapf(err, "failed to fetch index from heavy - %v", p.Message.Object.Record().DebugString()))
				return bus.Reply{Err: errors.Wrapf(err, "failed to fetch index from heavy")}
			}
		}
	} else if err != nil {
		return bus.Reply{Err: err}
	}

	if err = validateState(idx.StateID, state.ID()); err != nil {
		return bus.Reply{Reply: &reply.Error{ErrType: reply.ErrDeactivated}}
	}

	hash := record.HashVirtual(p.Dep.PCS.ReferenceHasher(), virtRec)
	recID := insolar.NewID(p.PulseNumber, hash)

	// LifelineIndex exists and latest record id does not match (preserving chain consistency).
	// For the case when vm can't save or send result to another vm and it tries to update the same record again
	if idx.LatestState != nil && !state.PrevStateID().Equal(*idx.LatestState) && idx.LatestState != recID {
		return bus.Reply{Err: errors.New("invalid state record")}
	}

	hash = record.HashVirtual(p.Dep.PCS.ReferenceHasher(), virtRec)
	id := insolar.NewID(p.PulseNumber, hash)
	rec := record.Material{
		Virtual: &virtRec,
		JetID:   p.JetID,
	}

	err = p.Dep.RecordModifier.Set(ctx, *id, rec)

	if err == object.ErrOverride {
		logger.WithField("type", fmt.Sprintf("%T", virtRec)).Warn("set record override (#1)")
		id = recID
	} else if err != nil {
		return bus.Reply{Err: errors.Wrap(err, "can't save record into storage")}
	}
	idx.LatestState = id
	idx.StateID = state.ID()
	if state.ID() == record.StateActivation {
		idx.Parent = state.(*record.Activate).Parent
	}

	idx.LatestUpdate = p.PulseNumber
	idx.JetID = p.JetID
	err = p.Dep.LifelineIndex.Set(ctx, p.PulseNumber, *p.Message.Object.Record(), idx)
	if err != nil {
		return bus.Reply{Err: err}
	}
	err = p.Dep.LifelineStateModifier.SetLifelineUsage(ctx, p.PulseNumber, *p.Message.Object.Record())
	if err != nil {
		return bus.Reply{Err: errors.Wrap(err, "failed to update lifeline usage state")}
	}

	logger.WithField("state", idx.LatestState.DebugString()).Debug("saved object")

	rep := reply.Object{
		Head:         p.Message.Object,
		State:        *idx.LatestState,
		Prototype:    state.GetImage(),
		IsPrototype:  state.GetIsPrototype(),
		ChildPointer: idx.ChildPointer,
		Parent:       idx.Parent,
	}
	return bus.Reply{Reply: &rep}
}

func (p *UpdateObject) saveIndexFromHeavy(
	ctx context.Context, jetID insolar.JetID, obj insolar.Reference, heavy *insolar.Reference,
) (object.Lifeline, error) {
	genericReply, err := p.Dep.Bus.Send(ctx, &message.GetObjectIndex{
		Object: obj,
	}, &insolar.MessageSendOptions{
		Receiver: heavy,
	})
	if err != nil {
		return object.Lifeline{}, errors.Wrap(err, "failed to send")
	}
	rep, ok := genericReply.(*reply.ObjectIndex)
	if !ok {
		return object.Lifeline{}, fmt.Errorf("failed to fetch object index: unexpected reply type %T (reply=%+v)", genericReply, genericReply)
	}
	idx, err := object.DecodeIndex(rep.Index)
	if err != nil {
		return object.Lifeline{}, errors.Wrap(err, "failed to decode")
	}

	idx.JetID = jetID
	err = p.Dep.LifelineIndex.Set(ctx, p.PulseNumber, *obj.Record(), idx)
	if err != nil {
		return object.Lifeline{}, errors.Wrap(err, "failed to save")
	}
	return idx, nil
}

func validateState(old record.StateID, new record.StateID) error {
	if old == record.StateDeactivation {
		return ErrObjectDeactivated
	}
	if old == record.StateUndefined && new != record.StateActivation {
		return errors.New("object is not activated")
	}
	if old != record.StateUndefined && new == record.StateActivation {
		return errors.New("object is already activated")
	}
	return nil
}
