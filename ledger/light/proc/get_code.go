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

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/object"
)

type GetCode struct {
	message *message.Message
	codeID  insolar.ID
	pass    bool

	Dep struct {
		RecordAccessor object.RecordAccessor
		Coordinator    jet.Coordinator
		BlobAccessor   blob.Accessor
		Sender         bus.Sender
		JetFetcher     jet.Fetcher
	}
}

func NewGetCode(msg *message.Message, codeID insolar.ID, pass bool) *GetCode {
	return &GetCode{
		message: msg,
		codeID:  codeID,
		pass:    pass,
	}
}

func (p *GetCode) Proceed(ctx context.Context) error {
	sendCode := func(rec record.Material) error {
		virtual := record.Unwrap(rec.Virtual)
		code, ok := virtual.(*record.Code)
		if !ok {
			return fmt.Errorf("invalid code record %#v", virtual)
		}
		b, err := p.Dep.BlobAccessor.ForID(ctx, code.Code)
		if err != nil {
			return errors.Wrap(err, "failed to fetch code blob")
		}
		buf, err := rec.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal record")
		}
		msg, err := payload.NewMessage(&payload.Code{
			Record: buf,
			Code:   b.Value,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create message")
		}
		go p.Dep.Sender.Reply(ctx, p.message, msg)

		return nil
	}

	sendPassCode := func() error {
		msg, err := payload.NewMessage(&payload.Pass{
			Origin:        p.message.Payload,
			CorrelationID: []byte(middleware.MessageCorrelationID(p.message)),
		})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}

		onHeavy, err := p.Dep.Coordinator.IsBeyondLimit(ctx, flow.Pulse(ctx), p.codeID.Pulse())
		if err != nil {
			return errors.Wrap(err, "failed to calculate pulse")
		}
		var node insolar.Reference
		if onHeavy {
			h, err := p.Dep.Coordinator.Heavy(ctx, flow.Pulse(ctx))
			if err != nil {
				return errors.Wrap(err, "failed to calculate heavy")
			}
			node = *h
		} else {
			jetID, err := p.Dep.JetFetcher.Fetch(ctx, p.codeID, p.codeID.Pulse())
			if err != nil {
				return errors.Wrap(err, "failed to fetch jet")
			}
			l, err := p.Dep.Coordinator.LightExecutorForJet(ctx, *jetID, p.codeID.Pulse())
			if err != nil {
				return errors.Wrap(err, "failed to calculate role")
			}
			node = *l
		}

		go func() {
			_, done := p.Dep.Sender.SendTarget(ctx, msg, node)
			done()
		}()
		return nil
	}

	logger := inslogger.FromContext(ctx)
	rec, err := p.Dep.RecordAccessor.ForID(ctx, p.codeID)
	switch err {
	case nil:
		logger.Info("sending code")
		return sendCode(rec)
	case object.ErrNotFound:
		if p.pass {
			logger.Info("code not found (sending pass)")
			return sendPassCode()
		}
		return errors.Wrap(err, "failed to fetch record")
	default:
		return errors.Wrap(err, "failed to fetch record")
	}
}
