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

package executor

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/light/executor.DetachedNotifier -o ./ -s _mock.go -g

type DetachedNotifier interface {
	Notify(
		ctx context.Context,
		openedRequests []record.CompositeFilamentRecord,
		objectID insolar.ID,
		closedRequestID insolar.ID,
	)
}

type DetachedNotifierDefault struct {
	sender bus.Sender
}

func NewDetachedNotifierDefault(
	sender bus.Sender,
) *DetachedNotifierDefault {
	return &DetachedNotifierDefault{
		sender: sender,
	}
}

// Notify sends notifications about detached requests that are ready for execution.
func (p *DetachedNotifierDefault) Notify(
	ctx context.Context,
	openedRequests []record.CompositeFilamentRecord,
	objectID insolar.ID,
	closedRequestID insolar.ID,
) {
	for _, req := range openedRequests {
		outgoing, ok := record.Unwrap(&req.Record.Virtual).(*record.OutgoingRequest)
		if !ok {
			continue
		}
		if !outgoing.IsDetached() {
			continue
		}
		if reasonRef := outgoing.ReasonRef(); *reasonRef.GetLocal() != closedRequestID {
			continue
		}

		buf, err := req.Record.Virtual.Marshal()
		if err != nil {
			inslogger.FromContext(ctx).Error(
				errors.Wrapf(err, "failed to notify about detached %s", req.RecordID.DebugString()),
			)
			continue
		}
		msg, err := payload.NewMessage(&payload.SagaCallAcceptNotification{
			ObjectID:          objectID,
			DetachedRequestID: req.RecordID,
			Request:           buf,
		})
		if err != nil {
			inslogger.FromContext(ctx).Error(
				errors.Wrapf(err, "failed to notify about detached %s", req.RecordID.DebugString()),
			)
			continue
		}
		_, done := p.sender.SendRole(ctx, msg, insolar.DynamicRoleVirtualExecutor, *insolar.NewReference(objectID))
		done()
	}
}
