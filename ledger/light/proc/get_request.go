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
	wmbus "github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type GetRequestWM struct {
	message   payload.Meta
	requestID insolar.ID

	Dep struct {
		RecordAccessor object.RecordAccessor
		Sender         wmbus.Sender
	}
}

func NewGetRequestWM(msg payload.Meta, requestID insolar.ID) *GetRequestWM {
	return &GetRequestWM{
		requestID: requestID,
		message:   msg,
	}
}

func (p *GetRequestWM) Proceed(ctx context.Context) error {
	rec, err := p.Dep.RecordAccessor.ForID(ctx, p.requestID)
	if err != nil {
		return errors.Wrap(err, "failed to fetch request")
	}
	msg, err := payload.NewMessage(&payload.Request{
		RequestID: p.requestID,
		Request:   record.Wrap(rec),
	})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}

	go p.Dep.Sender.Reply(ctx, p.message, msg)
	inslogger.FromContext(ctx).Info("sending request")

	return nil
}

type GetRequest struct {
	message payload.Meta
	request insolar.ID

	Dep struct {
		RecordAccessor object.RecordAccessor
		Sender         bus.Sender
	}
}

func NewGetRequest(request insolar.ID, message payload.Meta) *GetRequest {
	return &GetRequest{
		request: request,
		message: message,
	}
}

func (p *GetRequest) Proceed(ctx context.Context) error {
	rec, err := p.Dep.RecordAccessor.ForID(ctx, p.request)
	if err != nil {
		return errors.Wrap(err, "failed to fetch request")
	}

	virtRec := rec.Virtual
	concrete := record.Unwrap(virtRec)
	_, ok := concrete.(*record.IncomingRequest)
	if !ok {
		return errors.New("failed to decode request")
	}

	data, err := virtRec.Marshal()
	if err != nil {
		return errors.Wrap(err, "can't serialize record")
	}

	rep := &reply.Request{
		ID:     p.request,
		Record: data,
	}

	msg := bus.ReplyAsMessage(ctx, rep)
	go p.Dep.Sender.Reply(ctx, p.message, msg)
	return nil
}
