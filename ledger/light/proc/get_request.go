// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package proc

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/object"
)

type GetRequest struct {
	message             payload.Meta
	objectID, requestID insolar.ID
	pass                bool

	dep struct {
		records     object.RecordAccessor
		sender      bus.Sender
		coordinator jet.Coordinator
		fetcher     executor.JetFetcher
	}
}

func NewGetRequest(msg payload.Meta, objectID, requestID insolar.ID, pass bool) *GetRequest {
	return &GetRequest{
		requestID: requestID,
		objectID:  objectID,
		message:   msg,
		pass:      pass,
	}
}

func (p *GetRequest) Dep(
	records object.RecordAccessor,
	sender bus.Sender,
	coordinator jet.Coordinator,
	fetcher executor.JetFetcher,
) {
	p.dep.records = records
	p.dep.sender = sender
	p.dep.coordinator = coordinator
	p.dep.fetcher = fetcher
}

func (p *GetRequest) Proceed(ctx context.Context) error {
	sendRequest := func(rec record.Material) error {
		concrete := record.Unwrap(&rec.Virtual)
		_, isIncoming := concrete.(*record.IncomingRequest)
		_, isOutgoing := concrete.(*record.OutgoingRequest)
		if !isIncoming && !isOutgoing {
			return fmt.Errorf("unexpected request type")
		}

		msg, err := payload.NewMessage(&payload.Request{
			RequestID: p.requestID,
			Request:   rec.Virtual,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}

		p.dep.sender.Reply(ctx, p.message, msg)
		return nil
	}

	sendPassRequest := func() error {
		buf, err := p.message.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal origin meta message")
		}
		msg, err := payload.NewMessage(&payload.Pass{
			Origin: buf,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}

		onHeavy, err := p.dep.coordinator.IsBeyondLimit(ctx, p.requestID.Pulse())
		if err != nil {
			return errors.Wrap(err, "failed to calculate pulse")
		}
		var node insolar.Reference
		if onHeavy {
			h, err := p.dep.coordinator.Heavy(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to calculate heavy")
			}
			node = *h
		} else {
			jetID, err := p.dep.fetcher.Fetch(ctx, p.objectID, p.requestID.Pulse())
			if err != nil {
				return errors.Wrap(err, "failed to fetch jet")
			}
			l, err := p.dep.coordinator.LightExecutorForJet(ctx, *jetID, p.requestID.Pulse())
			if err != nil {
				return errors.Wrap(err, "failed to calculate role")
			}
			node = *l
		}

		_, done := p.dep.sender.SendTarget(ctx, msg, node)
		done()
		return nil
	}

	rec, err := p.dep.records.ForID(ctx, p.requestID)
	switch err {
	case nil:
		return sendRequest(rec)

	case object.ErrNotFound:
		if p.pass {
			return sendPassRequest()
		}

		return &payload.CodedError{
			Text: "request not found",
			Code: payload.CodeNotFound,
		}

	default:
		return errors.Wrap(err, "failed to fetch record")
	}
}
