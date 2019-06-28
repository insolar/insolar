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
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/light/recentstorage"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type SetResult struct {
	message   payload.Meta
	request   record.Virtual
	requestID insolar.ID
	jetID     insolar.JetID

	dep struct {
		writer        hot.WriteAccessor
		records       object.RecordModifier
		recentStorage recentstorage.Provider
		pendings      object.PendingModifier
		sender        bus.Sender
	}
}

func NewSetResult(
	msg payload.Meta,
	rec record.Virtual,
	recID insolar.ID,
	jetID insolar.JetID,
) *SetResult {
	return &SetResult{
		message:   msg,
		request:   rec,
		requestID: recID,
		jetID:     jetID,
	}
}

func (p *SetResult) Dep(
	w hot.WriteAccessor,
	r object.RecordModifier,
	rs recentstorage.Provider,
	pnds object.PendingModifier,
	s bus.Sender,
) {
	p.dep.writer = w
	p.dep.records = r
	p.dep.recentStorage = rs
	p.dep.pendings = pnds
	p.dep.sender = s
}

func (p *SetResult) Proceed(ctx context.Context) error {
	done, err := p.dep.writer.Begin(ctx, flow.Pulse(ctx))
	if err != nil {
		if err == hot.ErrWriteClosed {
			return flow.ErrCancelled
		}
		return err
	}
	defer done()

	material := record.Material{
		Virtual: &p.request,
		JetID:   p.jetID,
	}
	err = p.dep.records.Set(ctx, p.requestID, material)
	if err != nil {
		return errors.Wrap(err, "failed to store record")
	}

	err = p.handlePendings(ctx, p.request)
	if err != nil {
		return err
	}

	msg, err := payload.NewMessage(&payload.ID{ID: p.requestID})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}

	go p.dep.sender.Reply(ctx, p.message, msg)

	return nil
}

func (p *SetResult) handlePendings(ctx context.Context, virtRec record.Virtual) error {
	// TODO: check it after INS-1939
	// concrete := record.Unwrap(&virtRec)
	//
	// rec := concrete.(*record.Result)
	// recentStorage := p.dep.recentStorage.GetPendingStorage(ctx, insolar.ID(p.jetID))
	// recentStorage.RemovePendingRequest(ctx, rec.Object, *rec.Request.Record())
	// err := p.Dep.PendingModifier.SetResult(ctx, flow.Pulse(ctx), r.Object, calculatedID, *r)
	// if err != nil {
	// 	return &bus.Reply{Err: errors.Wrap(err, "can't save result into filament-index")}
	// }

	return nil
}
