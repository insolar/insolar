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
