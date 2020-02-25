// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package proc

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type SendRequest struct {
	meta payload.Meta

	dep struct {
		records object.RecordAccessor
		sender  bus.Sender
	}
}

func NewSendRequest(meta payload.Meta) *SendRequest {
	return &SendRequest{
		meta: meta,
	}
}

func (p *SendRequest) Dep(records object.RecordAccessor, sender bus.Sender) {
	p.dep.records = records
	p.dep.sender = sender
}

func (p *SendRequest) Proceed(ctx context.Context) error {
	msg := payload.GetRequest{}
	err := msg.Unmarshal(p.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to decode GetRequest payload")
	}

	rec, err := p.dep.records.ForID(ctx, msg.RequestID)
	if err == object.ErrNotFound {
		msg, err := payload.NewMessage(&payload.Error{
			Text: object.ErrNotFound.Error(),
			Code: payload.CodeNotFound,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}

		p.dep.sender.Reply(ctx, p.meta, msg)
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "failed to find a request")
	}

	concrete := record.Unwrap(&rec.Virtual)
	_, isIncoming := concrete.(*record.IncomingRequest)
	_, isOutgoing := concrete.(*record.OutgoingRequest)
	if !isIncoming && !isOutgoing {
		return fmt.Errorf("unexpected request type")
	}

	rep, err := payload.NewMessage(&payload.Request{
		RequestID: msg.RequestID,
		Request:   rec.Virtual,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create a Request message")
	}
	p.dep.sender.Reply(ctx, p.meta, rep)
	return nil
}
