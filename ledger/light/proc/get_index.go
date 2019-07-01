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
	"github.com/insolar/insolar/insolar/bus"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type EnsureIndex struct {
	object  insolar.Reference
	jet     insolar.JetID
	message payload.Meta
	pn      insolar.PulseNumber

	Dep struct {
		IndexLocker   object.IndexLocker
		IndexAccessor object.IndexAccessor
		IndexModifier object.IndexModifier
		Coordinator   jet.Coordinator
		Bus           insolar.MessageBus
		Sender        bus.Sender
	}
}

func NewEnsureIndex(obj insolar.Reference, jetID insolar.JetID, msg payload.Meta, pn insolar.PulseNumber) *EnsureIndex {
	return &EnsureIndex{
		object:  obj,
		jet:     jetID,
		message: msg,
		pn:      pn,
	}
}

func (p *EnsureIndex) Proceed(ctx context.Context) error {
	err := p.process(ctx)
	if err != nil {
		msg := bus.ErrorAsMessage(ctx, err)
		p.Dep.Sender.Reply(ctx, p.message, msg)
	}
	return err
}

func (p *EnsureIndex) process(ctx context.Context) error {
	objectID := *p.object.Record()
	logger := inslogger.FromContext(ctx)

	p.Dep.IndexLocker.Lock(&objectID)
	defer p.Dep.IndexLocker.Unlock(&objectID)

	idx, err := p.Dep.IndexAccessor.ForID(ctx, p.pn, objectID)
	if err == nil {
		if flow.Pulse(ctx) == p.pn {
			idx.LifelineLastUsed = p.pn
			err = p.Dep.IndexModifier.SetIndex(ctx, p.pn, idx)
			if err != nil {
				return errors.Wrap(err, "failed to update lifeline usage")
			}
		}
		return nil
	}
	if err != object.ErrIndexNotFound {
		return errors.Wrap(err, "failed to fetch index")
	}

	logger.Debug("failed to fetch index (fetching from heavy)")
	heavy, err := p.Dep.Coordinator.Heavy(ctx, flow.Pulse(ctx))
	if err != nil {
		return errors.Wrap(err, "failed to calculate heavy")
	}
	genericReply, err := p.Dep.Bus.Send(ctx, &message.GetObjectIndex{
		Object: p.object,
	}, &insolar.MessageSendOptions{
		Receiver: heavy,
	})
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"jet": p.jet.DebugString(),
			"pn":  flow.Pulse(ctx),
		}).Error(errors.Wrapf(err, "failed to fetch index from heavy - %v", p.object.Record().DebugString()))
		return errors.Wrap(err, "failed to fetch index from heavy")
	}
	rep, ok := genericReply.(*reply.ObjectIndex)
	if !ok {
		return fmt.Errorf("failed to fetch index from heavy: unexpected reply type %T", genericReply)
	}

	lfl, err := object.DecodeLifeline(rep.Index)
	if err != nil {
		return errors.Wrap(err, "failed to decode index")
	}

	lfl.JetID = p.jet
	err = p.Dep.IndexModifier.SetIndex(ctx, flow.Pulse(ctx), object.FilamentIndex{
		LifelineLastUsed: p.pn,
		Lifeline:         lfl,
		PendingRecords:   []insolar.ID{},
		ObjID:            *p.object.Record(),
	})
	if err != nil {
		return errors.Wrap(err, "failed to save lifeline")
	}

	return nil
}

type EnsureIndexWM struct {
	object  insolar.ID
	jet     insolar.JetID
	message payload.Meta

	Result struct {
		Lifeline object.Lifeline
	}

	Dep struct {
		IndexLocker   object.IndexLocker
		IndexModifier object.IndexModifier
		IndexAccessor object.IndexAccessor

		Coordinator jet.Coordinator
		Bus         insolar.MessageBus
		Sender      bus.Sender
	}
}

func NewEnsureIndexWM(obj insolar.ID, jetID insolar.JetID, msg payload.Meta) *EnsureIndexWM {
	return &EnsureIndexWM{
		object:  obj,
		jet:     jetID,
		message: msg,
	}
}

func (p *EnsureIndexWM) Proceed(ctx context.Context) error {
	err := p.process(ctx)
	if err != nil {
		msg, err := payload.NewMessage(&payload.Error{Text: err.Error()})
		if err != nil {
			return err
		}
		go p.Dep.Sender.Reply(ctx, p.message, msg)
	}
	return err
}

func (p *EnsureIndexWM) process(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	p.Dep.IndexLocker.Lock(&p.object)
	defer p.Dep.IndexLocker.Unlock(&p.object)

	idx, err := p.Dep.IndexAccessor.ForID(ctx, flow.Pulse(ctx), p.object)
	if err == nil {
		p.Result.Lifeline = idx.Lifeline

		idx.LifelineLastUsed = flow.Pulse(ctx)
		err = p.Dep.IndexModifier.SetIndex(ctx, flow.Pulse(ctx), idx)
		if err != nil {
			return errors.Wrap(err, "failed to update lifeline usage")
		}
		return nil
	}
	if err != object.ErrIndexNotFound {
		return errors.Wrap(err, "failed to fetch index")
	}

	logger.Debug("failed to fetch index (fetching from heavy)")
	heavy, err := p.Dep.Coordinator.Heavy(ctx, flow.Pulse(ctx))
	if err != nil {
		return errors.Wrap(err, "failed to calculate heavy")
	}
	genericReply, err := p.Dep.Bus.Send(ctx, &message.GetObjectIndex{
		Object: *insolar.NewReference(p.object),
	}, &insolar.MessageSendOptions{
		Receiver: heavy,
	})
	if err != nil {
		logger.WithFields(map[string]interface{}{
			"jet": p.jet.DebugString(),
			"pn":  flow.Pulse(ctx),
		}).Error(errors.Wrapf(err, "failed to fetch index from heavy - %v", p.object.DebugString()))
		return errors.Wrap(err, "failed to fetch index from heavy")
	}
	rep, ok := genericReply.(*reply.ObjectIndex)
	if !ok {
		return fmt.Errorf("failed to fetch index from heavy: unexpected reply type %T", genericReply)
	}

	p.Result.Lifeline, err = object.DecodeLifeline(rep.Index)
	if err != nil {
		return errors.Wrap(err, "failed to decode index")
	}

	p.Result.Lifeline.JetID = p.jet
	err = p.Dep.IndexModifier.SetIndex(ctx, flow.Pulse(ctx), object.FilamentIndex{
		LifelineLastUsed: flow.Pulse(ctx),
		Lifeline:         p.Result.Lifeline,
		PendingRecords:   []insolar.ID{},
		ObjID:            p.object,
	})
	if err != nil {
		return errors.Wrap(err, "failed to save lifeline")
	}

	return nil
}
