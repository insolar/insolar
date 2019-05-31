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

	wmessage "github.com/ThreeDotsLabs/watermill/message"

	"github.com/insolar/insolar/insolar"
	wbus "github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type GetIndex struct {
	object  insolar.Reference
	jet     insolar.JetID
	replyTo chan<- bus.Reply
	pn      insolar.PulseNumber

	Result struct {
		Index object.Lifeline
	}

	Dep struct {
		Index       object.LifelineIndex
		IndexState  object.LifelineStateModifier
		Locker      object.IDLocker
		Coordinator jet.Coordinator
		Bus         insolar.MessageBus
	}
}

func NewGetIndex(obj insolar.Reference, jetID insolar.JetID, rep chan<- bus.Reply, pn insolar.PulseNumber) *GetIndex {
	return &GetIndex{
		object:  obj,
		jet:     jetID,
		replyTo: rep,
		pn:      pn,
	}
}

func (p *GetIndex) Proceed(ctx context.Context) error {
	err := p.process(ctx)
	if err != nil {
		p.replyTo <- bus.Reply{Err: err}
	}
	return err
}

func (p *GetIndex) process(ctx context.Context) error {
	objectID := *p.object.Record()
	logger := inslogger.FromContext(ctx)

	p.Dep.Locker.Lock(&objectID)
	defer p.Dep.Locker.Unlock(&objectID)

	idx, err := p.Dep.Index.ForID(ctx, p.pn, objectID)
	if err == nil {
		p.Result.Index = idx
		if flow.Pulse(ctx) == p.pn {
			err = p.Dep.IndexState.SetLifelineUsage(ctx, p.pn, objectID)
			if err != nil {
				return errors.Wrap(err, "failed to update lifeline usage")
			}
		}
		return nil
	}
	if err != object.ErrLifelineNotFound {
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

	p.Result.Index, err = object.DecodeIndex(rep.Index)
	if err != nil {
		return errors.Wrap(err, "failed to decode index")
	}

	p.Result.Index.JetID = p.jet
	err = p.Dep.Index.Set(ctx, flow.Pulse(ctx), objectID, p.Result.Index)
	if err != nil {
		return errors.Wrap(err, "failed to save lifeline")
	}
	err = p.Dep.IndexState.SetLifelineUsage(ctx, flow.Pulse(ctx), objectID)
	if err != nil {
		return errors.Wrap(err, "failed to update lifeline usage")
	}

	return nil
}

type GetIndexWM struct {
	object  insolar.ID
	jet     insolar.JetID
	message *wmessage.Message

	Result struct {
		Index object.Lifeline
	}

	Dep struct {
		Index       object.LifelineIndex
		IndexState  object.LifelineStateModifier
		Locker      object.IDLocker
		Coordinator jet.Coordinator
		Bus         insolar.MessageBus
		Sender      wbus.Sender
	}
}

func NewGetIndexWM(obj insolar.ID, jetID insolar.JetID, msg *wmessage.Message) *GetIndexWM {
	return &GetIndexWM{
		object:  obj,
		jet:     jetID,
		message: msg,
	}
}

func (p *GetIndexWM) Proceed(ctx context.Context) error {
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

func (p *GetIndexWM) process(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	p.Dep.Locker.Lock(&p.object)
	defer p.Dep.Locker.Unlock(&p.object)

	idx, err := p.Dep.Index.ForID(ctx, flow.Pulse(ctx), p.object)
	if err == nil {
		p.Result.Index = idx
		err = p.Dep.IndexState.SetLifelineUsage(ctx, flow.Pulse(ctx), p.object)
		if err != nil {
			return errors.Wrap(err, "failed to update lifeline usage")
		}
		return nil
	}
	if err != object.ErrLifelineNotFound {
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

	p.Result.Index, err = object.DecodeIndex(rep.Index)
	if err != nil {
		return errors.Wrap(err, "failed to decode index")
	}

	p.Result.Index.JetID = p.jet
	err = p.Dep.Index.Set(ctx, flow.Pulse(ctx), p.object, p.Result.Index)
	if err != nil {
		return errors.Wrap(err, "failed to save lifeline")
	}
	err = p.Dep.IndexState.SetLifelineUsage(ctx, flow.Pulse(ctx), p.object)
	if err != nil {
		return errors.Wrap(err, "failed to update lifeline usage")
	}

	return nil
}
