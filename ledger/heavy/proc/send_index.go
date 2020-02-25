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
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type SendIndex struct {
	meta payload.Meta

	dep struct {
		indexes object.MemoryIndexAccessor
		sender  bus.Sender
	}
}

func (p *SendIndex) Dep(
	indexes object.MemoryIndexAccessor,
	sender bus.Sender,
) {
	p.dep.indexes = indexes
	p.dep.sender = sender
}

func NewSendIndex(meta payload.Meta) *SendIndex {
	return &SendIndex{
		meta: meta,
	}
}

func (p *SendIndex) Proceed(ctx context.Context) error {
	ensureIndex := payload.GetIndex{}
	err := ensureIndex.Unmarshal(p.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal ensureIndex message")
	}

	idx, err := p.dep.indexes.ForID(ctx, p.meta.Pulse, ensureIndex.ObjectID)
	if err == object.ErrIndexNotFound {
		return &payload.CodedError{
			Code: payload.CodeNotFound,
			Text: fmt.Sprintf("index not found for %v", ensureIndex.ObjectID.DebugString()),
		}
	}
	if err != nil {
		return errors.Wrapf(
			err,
			"failed to fetch object index for %v", ensureIndex.ObjectID.String(),
		)
	}

	buf := object.EncodeLifeline(idx.Lifeline)
	msg, err := payload.NewMessage(&payload.Index{
		Index: buf,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}

	p.dep.sender.Reply(ctx, p.meta, msg)
	return nil
}
