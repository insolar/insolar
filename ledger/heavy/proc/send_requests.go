/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

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
		indexes object.IndexAccessor
	}
}

func NewSendRequests(meta payload.Meta) *SendRequests {
	return &SendRequests{
		meta: meta,
	}
}

func (p *SendRequests) Dep(sender bus.Sender, records object.RecordAccessor, indexes object.IndexAccessor) {
	p.dep.sender = sender
	p.dep.records = records
	p.dep.indexes = indexes
}

func (p *SendRequests) Proceed(ctx context.Context) error {
	ctx, span := instracer.StartSpan(ctx, fmt.Sprintf("GetPendingFilament"))
	defer span.End()

	msg := payload.GetFilament{}
	err := msg.Unmarshal(p.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to decode PassState payload")
	}

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
			return err
		}
		composite.MetaID = *iter
		composite.Meta = filamentRecord

		// Fetching primary record.
		virtual := record.Unwrap(filamentRecord.Virtual)
		filament, ok := virtual.(*record.PendingFilament)
		if !ok {
			return errors.New("failed to convert filament record")
		}
		rec, err := p.dep.records.ForID(ctx, filament.RecordID)
		if err != nil {
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
		return errors.Wrap(err, "failed to create a PendingFilament message")
	}
	go p.dep.sender.Reply(ctx, p.meta, rep)
	return nil
}
