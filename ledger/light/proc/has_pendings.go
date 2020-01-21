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
