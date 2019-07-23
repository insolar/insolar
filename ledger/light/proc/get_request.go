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
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type GetRequest struct {
	message             payload.Meta
	objectID, requestID insolar.ID
	passed              bool

	dep struct {
		records     object.RecordAccessor
		sender      bus.Sender
		coordinator jet.Coordinator
		fetcher     executor.JetFetcher
	}
}

func NewGetRequest(msg payload.Meta, objectID, requestID insolar.ID, passed bool) *GetRequest {
	return &GetRequest{
		requestID: requestID,
		objectID:  objectID,
		message:   msg,
		passed:    passed,
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
		concrete := record.Unwrap(rec.Virtual)
		_, isIncoming := concrete.(*record.IncomingRequest)
		_, isOutgoing := concrete.(*record.OutgoingRequest)
		if !isIncoming && !isOutgoing {
			return fmt.Errorf("unexpected request type")
		}

		msg, err := payload.NewMessage(&payload.Request{
			RequestID: p.requestID,
			Request:   *rec.Virtual,
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

	fmt.Printf("looking for %s", p.requestID.DebugString())
	fmt.Println()
	rec, err := p.dep.records.ForID(ctx, p.requestID)
	switch err {
	case nil:
		return sendRequest(rec)

	case object.ErrNotFound:
		if !p.passed {
			return sendPassRequest()
		}

		msg, err := payload.NewMessage(&payload.Error{
			Text: "request not found",
			Code: payload.CodeNotFound,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}

		p.dep.sender.Reply(ctx, p.message, msg)
		return nil

	default:
		return errors.Wrap(err, "failed to fetch record")
	}
}
