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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/pkg/errors"
)

type SendRequests struct {
	message          payload.Meta
	objID, startFrom insolar.ID
	readUntil        insolar.PulseNumber

	dep struct {
		sender    bus.Sender
		filaments executor.FilamentCalculator
	}
}

func NewSendRequests(msg payload.Meta, objID insolar.ID, startFrom insolar.ID, readUntil insolar.PulseNumber) *SendRequests {
	return &SendRequests{
		message:   msg,
		objID:     objID,
		startFrom: startFrom,
		readUntil: readUntil,
	}
}

func (p *SendRequests) Dep(sender bus.Sender, filaments executor.FilamentCalculator) {
	p.dep.sender = sender
	p.dep.filaments = filaments
}

func (p *SendRequests) Proceed(ctx context.Context) error {
	records, err := p.dep.filaments.Requests(ctx, p.objID, p.startFrom, p.readUntil, flow.Pulse(ctx))
	if err != nil {
		return errors.Wrap(err, "failed to fetch filament")
	}

	msg, err := payload.NewMessage(&payload.FilamentSegment{
		ObjectID: p.objID,
		Records:  records,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create message")
	}
	go p.dep.sender.Reply(ctx, p.message, msg)
	return nil
}
