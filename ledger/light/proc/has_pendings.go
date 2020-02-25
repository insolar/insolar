// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package proc

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type HasPendings struct {
	message  payload.Meta
	objectID insolar.ID

	dep struct {
		index  object.MemoryIndexAccessor
		sender bus.Sender
	}
}

func NewHasPendings(msg payload.Meta, objectID insolar.ID) *HasPendings {
	return &HasPendings{
		message:  msg,
		objectID: objectID,
	}
}

func (hp *HasPendings) Dep(
	index object.MemoryIndexAccessor,
	sender bus.Sender,
) {
	hp.dep.index = index
	hp.dep.sender = sender
}

func (hp *HasPendings) Proceed(ctx context.Context) error {
	idx, err := hp.dep.index.ForID(ctx, flow.Pulse(ctx), hp.objectID)
	if err != nil {
		return err
	}

	msg, err := payload.NewMessage(&payload.PendingsInfo{
		HasPendings: idx.Lifeline.EarliestOpenRequest != nil && *idx.Lifeline.EarliestOpenRequest < flow.Pulse(ctx),
	})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}

	hp.dep.sender.Reply(ctx, hp.message, msg)
	return nil
}
