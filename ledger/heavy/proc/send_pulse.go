// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package proc

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
)

type SendPulse struct {
	meta payload.Meta

	dep struct {
		pulses pulse.Accessor
		sender bus.Sender
	}
}

func (p *SendPulse) Dep(
	pulses pulse.Accessor,
	sender bus.Sender,
) {
	p.dep.pulses = pulses
	p.dep.sender = sender
}

func NewSendPulse(meta payload.Meta) *SendPulse {
	return &SendPulse{
		meta: meta,
	}
}

func (p *SendPulse) Proceed(ctx context.Context) error {
	getPulse := payload.GetPulse{}
	err := getPulse.Unmarshal(p.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal GetPulse message")
	}

	foundPulse, err := p.dep.pulses.ForPulseNumber(ctx, getPulse.PulseNumber)
	if err != nil {
		return errors.Wrap(err, "failed to fetch pulse data from storage")
	}

	msg, err := payload.NewMessage(&payload.Pulse{
		Pulse: *pulse.ToProto(&foundPulse),
	})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}

	p.dep.sender.Reply(ctx, p.meta, msg)
	return nil
}
