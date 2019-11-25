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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

type GetPulse struct {
	message payload.Meta
	pulse   insolar.PulseNumber

	dep struct {
		coordinator jet.Coordinator
		sender      bus.Sender
	}
}

func NewGetPulse(msg payload.Meta, pulse insolar.PulseNumber) *GetPulse {
	return &GetPulse{
		message: msg,
		pulse:   pulse,
	}
}

func (p *GetPulse) Dep(
	c jet.Coordinator,
	s bus.Sender,
) {
	p.dep.coordinator = c
	p.dep.sender = s
}

func (p *GetPulse) Proceed(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	ctx, span := instracer.StartSpan(ctx, "GetPulse")
	defer span.Finish()

	logger.Debugf("GetPulse: go to heavy for pulse data, pulse: %v", p.pulse.String())
	heavy, err := p.dep.coordinator.Heavy(ctx)
	if err != nil {
		return errors.Wrap(err, "GetPulse: failed to calculate heavy")
	}

	getPulse, err := payload.NewMessage(&payload.GetPulse{
		PulseNumber: p.pulse,
	})
	if err != nil {
		return errors.Wrap(err, "GetPulse: failed to create GetPulse message")
	}

	reps, done := p.dep.sender.SendTarget(ctx, getPulse, *heavy)
	defer done()

	res, ok := <-reps
	if !ok {
		return errors.New("GetPulse: no reply")
	}

	pl, err := payload.UnmarshalFromMeta(res.Payload)
	if err != nil {
		return errors.Wrap(err, "GetPulse: failed to unmarshal reply")
	}

	switch rep := pl.(type) {
	case *payload.Pulse:
		msg, err := payload.NewMessage(rep)
		if err != nil {
			return errors.Wrap(err, "failed to create message")
		}
		p.dep.sender.Reply(ctx, p.message, msg)
		return nil
	case *payload.Error:
		return &payload.CodedError{
			Text: fmt.Sprintf("failed to fetch pulse data from heavy: %v, pulse=%v", rep.Text, p.pulse.String()),
			Code: rep.Code,
		}
	default:
		return fmt.Errorf("GetPulse: unexpected reply %T", pl)
	}
}
