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
	lifeline record.Lifeline

	Dep struct {
		Coordinator    jet.Coordinator
		Jets           jet.Storage
		JetFetcher     executor.JetFetcher
		RecordAccessor object.RecordAccessor
		Bus            insolar.MessageBus
		Sender         bus.Sender
	}
}

func NewSendObject(
	msg payload.Meta,
	id insolar.ID,
	lifeline record.Lifeline,
) *SendObject {
	return &SendObject{
		message:  msg,
		lifeline: lifeline,
		objectID: id,
	}
}

func (p *SendObject) Proceed(ctx context.Context) error {
	sendState := func(rec record.Material) error {
		virtual := rec.Virtual
		concrete := record.Unwrap(virtual)
		state, ok := concrete.(record.State)
		if !ok {
			return fmt.Errorf("invalid object record %#v", virtual)
		}

		if state.ID() == record.StateDeactivation {
			msg, err := payload.NewMessage(&payload.Error{Text: "object is deactivated", Code: payload.CodeDeactivated})
			if err != nil {
				return errors.Wrap(err, "failed to create reply")
			}
			go p.Dep.Sender.Reply(ctx, p.message, msg)
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
		go p.Dep.Sender.Reply(ctx, p.message, msg)

		return nil
	}

	sendPassState := func(stateID insolar.ID) error {
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

		onHeavy, err := p.Dep.Coordinator.IsBeyondLimit(ctx, stateID.Pulse())
		if err != nil {
			return errors.Wrap(err, "failed to calculate pulse")
		}
		var node insolar.Reference
		if onHeavy {
			h, err := p.Dep.Coordinator.Heavy(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to calculate heavy")
			}
			node = *h
		} else {
			jetID, err := p.Dep.JetFetcher.Fetch(ctx, p.objectID, stateID.Pulse())
			if err != nil {
				return errors.Wrap(err, "failed to fetch jet")
			}
			l, err := p.Dep.Coordinator.LightExecutorForJet(ctx, *jetID, stateID.Pulse())
			if err != nil {
				return errors.Wrap(err, "failed to calculate role")
			}
			node = *l
		}

		go func() {
			_, done := p.Dep.Sender.SendTarget(ctx, msg, node)
			done()
		}()
		return nil
	}

	if p.lifeline.LatestState == nil {
		return ErrNotActivated
	}

	{
		buf, err := p.lifeline.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal index")
		}
		msg, err := payload.NewMessage(&payload.Index{
			Index: buf,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}

		p.Dep.Sender.Reply(ctx, p.message, msg)
	}

	rec, err := p.Dep.RecordAccessor.ForID(ctx, *p.lifeline.LatestState)
	switch err {
	case nil:
		return sendState(rec)
	case object.ErrNotFound:
		return sendPassState(*p.lifeline.LatestState)
	default:
		return errors.Wrap(err, "failed to fetch record")
	}
}
