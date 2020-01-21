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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/light/executor"
)

type SendFilament struct {
	message          payload.Meta
	objID, startFrom insolar.ID
	readUntil        insolar.PulseNumber

	dep struct {
		sender    bus.Sender
		filaments executor.FilamentCalculator
	}
}

func NewSendFilament(msg payload.Meta, objID insolar.ID, startFrom insolar.ID, readUntil insolar.PulseNumber) *SendFilament {
	return &SendFilament{
		message:   msg,
		objID:     objID,
		startFrom: startFrom,
		readUntil: readUntil,
	}
}

func (p *SendFilament) Dep(sender bus.Sender, filaments executor.FilamentCalculator) {
	p.dep.sender = sender
	p.dep.filaments = filaments
}

func (p *SendFilament) Proceed(ctx context.Context) error {
	ctx, span := instracer.StartSpan(ctx, "SendFilament")
	defer span.Finish()

	span.SetTag("objID", p.objID.DebugString()).
		SetTag("startFrom", p.startFrom.DebugString()).
		SetTag("readUntil", p.readUntil.String())

	records, err := p.dep.filaments.Requests(ctx, p.objID, p.startFrom, p.readUntil)
	if err != nil {
		return errors.Wrap(err, "failed to fetch filament")
	}
	if len(records) == 0 {
		return &payload.CodedError{
			Text: "requests not found",
			Code: payload.CodeNotFound,
		}
	}

	msg, err := payload.NewMessage(&payload.FilamentSegment{
		ObjectID: p.objID,
		Records:  records,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create message")
	}
	p.dep.sender.Reply(ctx, p.message, msg)
	return nil
}
