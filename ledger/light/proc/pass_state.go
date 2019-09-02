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

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type PassState struct {
	message payload.Meta

	dep struct {
		sender  bus.Sender
		records object.RecordAccessor
	}
}

func NewPassState(meta payload.Meta) *PassState {
	return &PassState{
		message: meta,
	}
}

func (p *PassState) Dep(
	records object.RecordAccessor,
	sender bus.Sender,
) {
	p.dep.records = records
	p.dep.sender = sender
}

func (p *PassState) Proceed(ctx context.Context) error {

	sendObject := func(rec record.Material, origin payload.Meta) error {
		virtual := rec.Virtual
		concrete := record.Unwrap(&virtual)
		state, ok := concrete.(record.State)
		if !ok {
			return p.replyError(ctx, fmt.Sprintf("invalid object record %#v", virtual), payload.CodeInvalidRequest)
		}

		if state.ID() == record.StateDeactivation {
			return p.replyError(ctx, "object is deactivated", payload.CodeDeactivated)
		}

		buf, err := rec.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal state record")
		}
		msg, err := payload.NewMessage(&payload.State{
			Record: buf,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create message")
		}

		p.dep.sender.Reply(ctx, origin, msg)

		return nil
	}

	pass := payload.PassState{}
	err := pass.Unmarshal(p.message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to decode PassState payload")
	}

	origin := payload.Meta{}
	err = origin.Unmarshal(pass.Origin)
	if err != nil {
		return errors.Wrap(err, "failed to decode origin message")
	}

	rec, err := p.dep.records.ForID(ctx, pass.StateID)
	switch err {
	case nil:
		return sendObject(rec, origin)
	case object.ErrNotFound:
		return p.replyError(ctx, "state not found", payload.CodeNotFound)
	default:
		return errors.Wrap(err, "failed to fetch object state")
	}
}

func (p *PassState) replyError(
	ctx context.Context,
	text string,
	code uint32,
) error {
	msg, err := payload.NewMessage(&payload.Error{
		Text: text,
		Code: code,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}

	p.dep.sender.Reply(ctx, p.message, msg)
	return nil
}
