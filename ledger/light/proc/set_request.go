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
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/light/recentstorage"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type SetRequest struct {
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

func NewSetRequest(
	msg payload.Meta,
	rec record.Virtual,
	recID insolar.ID,
	jetID insolar.JetID,
) *SetRequest {
	return &SetRequest{
		message:   msg,
		request:   rec,
		requestID: recID,
		jetID:     jetID,
	}
}

func (p *SetRequest) Dep(
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

func (p *SetRequest) Proceed(ctx context.Context) error {
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
	if err == object.ErrOverride {
		inslogger.FromContext(ctx).Errorf("can't save record into storage: %s", err)
		// Since there is no deduplication yet it's quite possible that there will be
		// two writes by the same key. For this reason currently instead of reporting
		// an error we return OK (nil error). When deduplication will be implemented
		// we should change `nil` to `ErrOverride` here.
		return nil
	} else if err != nil {
		return errors.Wrap(err, "failed to store record")
	}

	err = p.handlePendings(ctx, p.requestID, p.request)
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

func (p *SetRequest) handlePendings(ctx context.Context, id insolar.ID, virtReq record.Virtual) error {
	concrete := record.Unwrap(&virtReq)
	req := concrete.(*record.Request)

	// Skip object creation and genesis
	if req.CallType == record.CTMethod {
		if p.dep.recentStorage.Count() > recentstorage.PendingRequestsLimit {
			return insolar.ErrTooManyPendingRequests
		}
		recentStorage := p.dep.recentStorage.GetPendingStorage(ctx, insolar.ID(p.jetID))
		recentStorage.AddPendingRequest(ctx, *req.Object.Record(), id)

		// err := p.dep.pendings.SetRequest(ctx, flow.Pulse(ctx), *req.Object.Record(), id)
		// if err != nil {
		// 	return errors.Wrap(err, "can't save result into filament-index")
		// }
	}

	return nil
}
