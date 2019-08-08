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

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/object"
)

type SendObject struct {
	message  payload.Meta
	objectID insolar.ID

	dep struct {
		coordinator jet.Coordinator
		jets        jet.Storage
		jetFetcher  executor.JetFetcher
		records     object.RecordAccessor
		indices     object.IndexAccessor
		sender      bus.Sender
	}
}

func NewSendObject(
	msg payload.Meta,
	id insolar.ID,
) *SendObject {
	return &SendObject{
		message:  msg,
		objectID: id,
	}
}

func (p *SendObject) Dep(
	coordinator jet.Coordinator,
	jets jet.Storage,
	jetFetcher executor.JetFetcher,
	records object.RecordAccessor,
	indices object.IndexAccessor,
	sender bus.Sender,
) {
	p.dep.coordinator = coordinator
	p.dep.jets = jets
	p.dep.jetFetcher = jetFetcher
	p.dep.records = records
	p.dep.indices = indices
	p.dep.sender = sender
}

func (p *SendObject) Proceed(ctx context.Context) error {
	sendState := func(rec record.Material) error {
		virtual := rec.Virtual
		concrete := record.Unwrap(&virtual)
		state, ok := concrete.(record.State)
		if !ok {
			return fmt.Errorf("invalid object record %#v", virtual)
		}

		if state.ID() == record.StateDeactivation {
			msg, err := payload.NewMessage(&payload.Error{Text: "object is deactivated", Code: payload.CodeDeactivated})
			if err != nil {
				return errors.Wrap(err, "failed to create reply")
			}
			p.dep.sender.Reply(ctx, p.message, msg)
			return nil
		}

		buf, err := rec.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal state record")
		}
		msg, err := payload.NewMessage(&payload.State{
			Record: buf,
			Memory: state.GetMemory(),
		})
		if err != nil {
			return errors.Wrap(err, "failed to create message")
		}
		p.dep.sender.Reply(ctx, p.message, msg)

		return nil
	}

	sendPassState := func(stateID insolar.ID) error {
		ctx, span := instracer.StartSpan(ctx, "SendObject.sendPassState")
		defer span.End()

		buf, err := p.message.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal origin meta message")
		}
		msg, err := payload.NewMessage(&payload.PassState{
			Origin:  buf,
			StateID: stateID,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}

		onHeavy, err := p.dep.coordinator.IsBeyondLimit(ctx, stateID.Pulse())
		if err != nil {
			return errors.Wrap(err, "failed to calculate pulse")
		}
		var node insolar.Reference
		if onHeavy {
			inslogger.FromContext(ctx).Warnf("State not found on light. Go to heavy. StateID:%v, CurrentPN:%v", stateID.DebugString(), flow.Pulse(ctx))
			h, err := p.dep.coordinator.Heavy(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to calculate heavy")
			}
			node = *h
			span.Annotate(nil, fmt.Sprintf("Send StateID:%v to heavy", stateID.DebugString()))
		} else {
			inslogger.FromContext(ctx).Warnf("State not found on light. Go to light. StateID:%v, CurrentPN:%v", stateID.DebugString(), flow.Pulse(ctx))
			jetID, err := p.dep.jetFetcher.Fetch(ctx, p.objectID, stateID.Pulse())
			if err != nil {
				return errors.Wrap(err, "failed to fetch jet")
			}
			l, err := p.dep.coordinator.LightExecutorForJet(ctx, *jetID, stateID.Pulse())
			if err != nil {
				return errors.Wrap(err, "failed to calculate role")
			}
			node = *l
			span.Annotate(nil, fmt.Sprintf("Send StateID:%v to light", stateID.DebugString()))
		}

		go func() {
			_, done := p.dep.sender.SendTarget(ctx, msg, node)
			done()
		}()
		return nil
	}

	idx, err := p.dep.indices.ForID(ctx, flow.Pulse(ctx), p.objectID)
	if err != nil {
		return errors.Wrap(err, "can't get index from storage")
	}

	lifeline := idx.Lifeline

	if lifeline.StateID == record.StateDeactivation {
		return errors.New("object is deactivated")
	}

	if lifeline.LatestState == nil {
		return ErrNotActivated
	}

	{
		buf, err := lifeline.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal index")
		}
		msg, err := payload.NewMessage(&payload.Index{
			Index: buf,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}

		p.dep.sender.Reply(ctx, p.message, msg)
	}

	rec, err := p.dep.records.ForID(ctx, *lifeline.LatestState)
	switch err {
	case nil:
		return sendState(rec)
	case object.ErrNotFound:
		return sendPassState(*lifeline.LatestState)
	default:
		return errors.Wrap(err, "failed to fetch record")
	}
}
