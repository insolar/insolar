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
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/object"
)

type EnsureIndex struct {
	object  insolar.ID
	jet     insolar.JetID
	message payload.Meta
	pulse   insolar.PulseNumber

	dep struct {
		indexes       object.MemoryIndexStorage
		coordinator   jet.Coordinator
		sender        bus.Sender
		writeAccessor executor.WriteAccessor
	}
}

func NewEnsureIndex(obj insolar.ID, jetID insolar.JetID, msg payload.Meta, pulse insolar.PulseNumber) *EnsureIndex {
	return &EnsureIndex{
		object:  obj,
		jet:     jetID,
		message: msg,
		pulse:   pulse,
	}
}

func (p *EnsureIndex) Dep(
	idxs object.MemoryIndexStorage,
	c jet.Coordinator,
	s bus.Sender,
	wc executor.WriteAccessor,
) {
	p.dep.indexes = idxs
	p.dep.coordinator = c
	p.dep.sender = s
	p.dep.writeAccessor = wc
}

func (p *EnsureIndex) Proceed(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	ctx, span := instracer.StartSpan(ctx, "EnsureIndex")
	defer span.Finish()

	span.SetTag("object_id", p.object.DebugString())

	_, err := p.dep.indexes.ForID(ctx, p.pulse, p.object)
	if err == nil {
		return nil
	}
	if err != object.ErrIndexNotFound {
		return errors.Wrap(err, "EnsureIndex: failed to fetch index")
	}

	logger.Debug("EnsureIndex: failed to fetch index (fetching from heavy)")
	heavy, err := p.dep.coordinator.Heavy(ctx)
	if err != nil {
		return errors.Wrap(err, "EnsureIndex: failed to calculate heavy")
	}

	ensureIndex, err := payload.NewMessage(&payload.GetIndex{
		ObjectID: p.object,
	})
	if err != nil {
		return errors.Wrap(err, "EnsureIndex: failed to create EnsureIndex message")
	}

	reps, done := p.dep.sender.SendTarget(ctx, ensureIndex, *heavy)
	defer done()

	res, ok := <-reps
	if !ok {
		return errors.New("EnsureIndex: no reply")
	}

	pl, err := payload.UnmarshalFromMeta(res.Payload)
	if err != nil {
		return errors.Wrap(err, "EnsureIndex: failed to unmarshal reply")
	}

	switch rep := pl.(type) {
	case *payload.Index:
		idx, err := object.DecodeLifeline(rep.Index)
		if err != nil {
			return errors.Wrap(err, "EnsureIndex: failed to decode index")
		}

		done, err := p.dep.writeAccessor.Begin(ctx, p.pulse)
		if err != nil {
			if err == executor.ErrWriteClosed {
				return flow.ErrCancelled
			}
			return errors.Wrap(err, "failed to write to db")
		}
		defer done()

		p.dep.indexes.SetIfNone(ctx, p.pulse, record.Index{
			LifelineLastUsed: p.pulse,
			Lifeline:         idx,
			PendingRecords:   []insolar.ID{},
			ObjID:            p.object,
		})
		return nil
	case *payload.Error:
		return &payload.CodedError{
			Text: fmt.Sprint("failed to fetch index from heavy: ", rep.Text),
			Code: rep.Code,
		}
	default:
		return fmt.Errorf("EnsureIndex: unexpected reply %T", pl)
	}
}
