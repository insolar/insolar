// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package proc

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/object"
)

type SetCode struct {
	message  payload.Meta
	record   record.Virtual
	code     []byte
	recordID insolar.ID
	jetID    insolar.JetID

	dep struct {
		writer  executor.WriteAccessor
		records object.AtomicRecordModifier
		pcs     insolar.PlatformCryptographyScheme
		sender  bus.Sender
	}
}

func NewSetCode(msg payload.Meta, rec record.Virtual, recID insolar.ID, jetID insolar.JetID) *SetCode {
	return &SetCode{
		message:  msg,
		record:   rec,
		recordID: recID,
		jetID:    jetID,
	}
}

func (p *SetCode) Dep(
	w executor.WriteAccessor,
	r object.AtomicRecordModifier,
	pcs insolar.PlatformCryptographyScheme,
	s bus.Sender,
) {
	p.dep.writer = w
	p.dep.records = r
	p.dep.pcs = pcs
	p.dep.sender = s
}

func (p *SetCode) Proceed(ctx context.Context) error {
	done, err := p.dep.writer.Begin(ctx, flow.Pulse(ctx))
	if err != nil {
		if err == executor.ErrWriteClosed {
			return flow.ErrCancelled
		}
		return err
	}
	defer done()

	material := record.Material{
		Virtual: p.record,
		JetID:   p.jetID,
		ID:      p.recordID,
	}

	err = p.dep.records.SetAtomic(ctx, material)
	if err != nil {
		return errors.Wrap(err, "failed to store record")
	}

	msg, err := payload.NewMessage(&payload.ID{ID: p.recordID})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}

	p.dep.sender.Reply(ctx, p.message, msg)

	return nil
}
