// Copyright 2020 Insolar Network Ltd.
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

package proc

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type PassState struct {
	message payload.Meta

	Dep struct {
		Sender  bus.Sender
		Records object.RecordAccessor
		Pulses  pulse.Accessor
	}
}

func NewPassState(msg payload.Meta) *PassState {
	return &PassState{
		message: msg,
	}
}

func (p *PassState) Proceed(ctx context.Context) error {
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

	rec, err := p.Dep.Records.ForID(ctx, pass.StateID)
	if err == object.ErrNotFound {
		var latestPulse insolar.PulseNumber
		latest, err := p.Dep.Pulses.Latest(ctx)
		if err == nil {
			latestPulse = latest.PulseNumber
		}
		inslogger.FromContext(ctx).Errorf(
			"state not found. StateID: %s, messagePN: %v, latestPN: %v",
			pass.StateID.DebugString(),
			origin.Pulse,
			latestPulse,
		)
		msg, err := payload.NewMessage(&payload.Error{Text: "state not found"})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}
		p.Dep.Sender.Reply(ctx, origin, msg)
		return nil
	}
	if err != nil {
		return err
	}

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

		p.Dep.Sender.Reply(ctx, origin, msg)
		return nil
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

	p.Dep.Sender.Reply(ctx, origin, msg)

	return nil
}
