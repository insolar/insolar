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
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type SetRequest struct {
	message   payload.Meta
	request   record.Request
	requestID insolar.ID
	jetID     insolar.JetID

	dep struct {
		writer   hot.WriteAccessor
		records  object.RecordModifier
		filament executor.FilamentModifier
		sender   bus.Sender
		locker   object.IndexLocker
	}
}

func NewSetRequest(
	msg payload.Meta,
	rec record.Request,
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
	f executor.FilamentModifier,
	s bus.Sender,
	l object.IndexLocker,
) {
	p.dep.writer = w
	p.dep.records = r
	p.dep.filament = f
	p.dep.sender = s
	p.dep.locker = l
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

	var req *record.CompositeFilamentRecord
	var res *record.CompositeFilamentRecord
	var objID insolar.ID

	if p.request.GetCallType() == record.CTSaveAsChild || p.request.GetCallType() == record.CTSaveAsDelegate {
		p.dep.locker.Lock(&p.requestID)
		defer p.dep.locker.Unlock(&p.requestID)

		req, res, err = p.dep.filament.SetActivationRequest(ctx, p.requestID, p.jetID, p.request)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		objID = p.requestID
	} else {
		p.dep.locker.Lock(p.request.AffinityRef().Record())
		defer p.dep.locker.Unlock(p.request.AffinityRef().Record())

		req, res, err = p.dep.filament.SetRequest(ctx, p.requestID, p.jetID, p.request)
		if err != nil {
			return errors.Wrap(err, "failed to store record")
		}
		objID = *p.request.AffinityRef().Record()
	}

	var reqBuf []byte
	var resBuf []byte

	if req != nil {
		reqBuf, err = req.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal stored record")
		}
	}
	if res != nil {
		resBuf, err = res.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal stored record")
		}
	}

	msg, err := payload.NewMessage(&payload.RequestInfo{
		ObjectID:  objID,
		RequestID: p.requestID,
		Request:   reqBuf,
		Result:    resBuf,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}

	go p.dep.sender.Reply(ctx, p.message, msg)
	return nil
}
