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
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/object"
)

type GetCode struct {
	message payload.Meta
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

func NewGetCode(msg payload.Meta, codeID insolar.ID, pass bool) *GetCode {
	return &GetCode{
		message: msg,
		codeID:  codeID,
		pass:    pass,
	}
}

func (p *GetCode) Proceed(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	sendCode := func(rec record.Material) error {
		buf, err := rec.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal record")
		}
		msg, err := payload.NewMessage(&payload.Code{
			Record: buf,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create message")
		}

		go p.Dep.Sender.Reply(ctx, p.message, msg)

		return nil
	}

	sendPassCode := func() error {
		originMeta, err := p.message.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal origin meta message")
		}
		msg, err := payload.NewMessage(&payload.Pass{
			Origin: originMeta,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create reply")
		}

		onHeavy, err := p.Dep.Coordinator.IsBeyondLimit(ctx, p.codeID.Pulse())
		if err != nil {
			return errors.Wrap(err, "failed to calculate pulse")
		}
		var node insolar.Reference
		if onHeavy {
			h, err := p.Dep.Coordinator.Heavy(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to calculate heavy")
			}
			node = *h
		} else {
			jetID, err := p.Dep.JetFetcher.Fetch(ctx, p.codeID, p.codeID.Pulse())
			if err != nil {
				return errors.Wrap(err, "failed to fetch jet")
			}
			logger.Debug("calculated jet for pass: %s", jetID.DebugString())
			l, err := p.Dep.Coordinator.LightExecutorForJet(ctx, *jetID, p.codeID.Pulse())
			if err != nil {
				return errors.Wrap(err, "failed to calculate role")
			}
			node = *l
		}

		go func() {
			_, done := p.Dep.Sender.SendTarget(ctx, msg, node)
			done()
			logger.Debug("passed GetCode")
		}()
		return nil
	}

	rec, err := p.Dep.RecordAccessor.ForID(ctx, p.codeID)
	switch err {
	case nil:
		logger.Debug("sending code")
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
