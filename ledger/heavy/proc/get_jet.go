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

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/pkg/errors"
)

type GetJet struct {
	meta payload.Meta

	dep struct {
		jets   jet.Accessor
		sender bus.Sender
	}
}

func (p *GetJet) Dep(
	jets jet.Accessor,
	sender bus.Sender,
) {
	p.dep.jets = jets
	p.dep.sender = sender
}

func NewGetJet(meta payload.Meta) *GetJet {
	return &GetJet{
		meta: meta,
	}
}

func (p *GetJet) Proceed(ctx context.Context) error {
	getJet := payload.GetJet{}
	err := getJet.Unmarshal(p.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal GetJet message")
	}

	jetID, actual := p.dep.jets.ForID(ctx, getJet.PulseNumber, getJet.ObjectID)

	msg, err := payload.NewMessage(&payload.Jet{
		JetID:  jetID,
		Actual: actual,
	})
	if err != nil {
		return errors.Wrap(err, "GetJet: failed to create reply")
	}

	p.dep.sender.Reply(ctx, p.meta, msg)
	return nil
}
