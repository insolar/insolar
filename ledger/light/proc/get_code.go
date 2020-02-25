// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package proc

import (
	"context"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/object"
)

type GetCode struct {
	message payload.Meta
	codeID  insolar.ID
	pass    bool

	dep struct {
		records     object.RecordAccessor
		coordinator jet.Coordinator
		jetFetcher  executor.JetFetcher
		sender      bus.Sender
	}
}

func NewGetCode(msg payload.Meta, codeID insolar.ID, pass bool) *GetCode {
	return &GetCode{
		message: msg,
		codeID:  codeID,
		pass:    pass,
	}
}

func (p *GetCode) Dep(
	r object.RecordAccessor,
	c jet.Coordinator,
	f executor.JetFetcher,
	s bus.Sender,
) {
	p.dep.records = r
	p.dep.coordinator = c
	p.dep.jetFetcher = f
	p.dep.sender = s
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

		p.dep.sender.Reply(ctx, p.message, msg)

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

		onHeavy, err := p.dep.coordinator.IsBeyondLimit(ctx, p.codeID.Pulse())
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
			jetID, err := p.dep.jetFetcher.Fetch(ctx, p.codeID, p.codeID.Pulse())
			if err != nil {
				return errors.Wrap(err, "failed to fetch jet")
			}
			logger.Debug("calculated jet for pass: %s", jetID.DebugString())
			l, err := p.dep.coordinator.LightExecutorForJet(ctx, *jetID, p.codeID.Pulse())
			if err != nil {
				return errors.Wrap(err, "failed to calculate role")
			}
			node = *l
		}

		go func() {
			_, done := p.dep.sender.SendTarget(ctx, msg, node)
			done()
			logger.Debug("passed GetCode")
		}()
		return nil
	}

	rec, err := p.dep.records.ForID(ctx, p.codeID)
	switch err {
	case nil:
		logger.Debug("sending code")
		return sendCode(rec)
	case object.ErrNotFound:
		if p.pass {
			logger.Info("code not found (sending pass)")
			return sendPassCode()
		}
		return &payload.CodedError{
			Text: "code not found",
			Code: payload.CodeNotFound,
		}

	default:
		return errors.Wrap(err, "failed to fetch record")
	}
}
