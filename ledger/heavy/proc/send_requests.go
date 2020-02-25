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
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type SendRequests struct {
	meta payload.Meta

	dep struct {
		sender  bus.Sender
		records object.RecordAccessor
		indexes object.MemoryIndexAccessor
	}
}

func NewSendRequests(meta payload.Meta) *SendRequests {
	return &SendRequests{
		meta: meta,
	}
}

func (p *SendRequests) Dep(sender bus.Sender, records object.RecordAccessor, indexes object.MemoryIndexAccessor) {
	p.dep.sender = sender
	p.dep.records = records
	p.dep.indexes = indexes
}

func (p *SendRequests) Proceed(ctx context.Context) error {
	ctx, span := instracer.StartSpan(ctx, fmt.Sprintf("SendRequests"))
	defer span.Finish()

	msg := payload.GetFilament{}
	err := msg.Unmarshal(p.meta.Payload)
	if err != nil {
		instracer.AddError(span, err)
		return errors.Wrap(err, "failed to decode GetFilament payload")
	}

	span.SetTag(
		"objID", msg.ObjectID.DebugString()).
		SetTag("startFrom", msg.StartFrom.DebugString()).
		SetTag("readUntil", msg.ReadUntil.String())

	_, err = p.dep.indexes.ForID(ctx, msg.StartFrom.Pulse(), msg.ObjectID)
	if err != nil {
		return errors.Wrap(err, "failed to find object")
	}

	var records []record.CompositeFilamentRecord
	iter := &msg.StartFrom
	for iter != nil && iter.Pulse() >= msg.ReadUntil {
		var composite record.CompositeFilamentRecord
		// Fetching filament record.
		filamentRecord, err := p.dep.records.ForID(ctx, *iter)
		if err != nil {
			instracer.AddError(span, err)
			return err
		}
		composite.MetaID = *iter
		composite.Meta = filamentRecord

		// Fetching primary record.
		virtual := record.Unwrap(&filamentRecord.Virtual)
		filament, ok := virtual.(*record.PendingFilament)
		if !ok {
			return errors.New("failed to convert filament record")
		}
		rec, err := p.dep.records.ForID(ctx, filament.RecordID)
		if err != nil {
			instracer.AddError(span, err)
			return err
		}
		composite.RecordID = filament.RecordID
		composite.Record = rec

		records = append(records, composite)

		// Iterating back.
		iter = filament.PreviousRecord
	}

	if len(records) == 0 {
		return errors.New("wrong filament request. empty segment")
	}

	rep, err := payload.NewMessage(&payload.FilamentSegment{
		ObjectID: msg.ObjectID,
		Records:  records,
	})
	if err != nil {
		instracer.AddError(span, err)
		return errors.Wrap(err, "failed to create a FilamentSegment message")
	}
	p.dep.sender.Reply(ctx, p.meta, rep)
	return nil
}
