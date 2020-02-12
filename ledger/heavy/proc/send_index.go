// Copyright 2020 Insolar Network Ltd.
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
