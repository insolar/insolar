//
// Copyright 2019 Insolar Technologies GmbH
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
//

package logicrunner

import (
	"context"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

type initializeAbandonedRequestsNotificationExecutionState struct {
	dep *Dependencies
	msg payload.AbandonedRequestsNotification
}

// Proceed initializes or sets LedgerHasMoreRequests to right value
func (p *initializeAbandonedRequestsNotificationExecutionState) Proceed(ctx context.Context) error {
	ref := *insolar.NewReference(p.msg.ObjectID)

	broker := p.dep.StateStorage.UpsertExecutionState(ref)
	broker.AbandonedRequestsOnLedger(ctx)

	return nil
}

type HandleAbandonedRequestsNotification struct {
	dep  *Dependencies
	meta payload.Meta
}

func (h *HandleAbandonedRequestsNotification) Present(ctx context.Context, f flow.Flow) error {
	abandoned := payload.AbandonedRequestsNotification{}
	err := abandoned.Unmarshal(h.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal AbandonedRequestsNotification message")
	}

	ctx, _ = inslogger.WithField(ctx, "targetid", abandoned.ObjectID.String())
	logger := inslogger.FromContext(ctx)

	logger.Debug("HandleAbandonedRequestsNotification.Present starts ...")

	ctx, span := instracer.StartSpan(ctx, "HandleAbandonedRequestsNotification.Present")
	span.AddAttributes(trace.StringAttribute("msg.Type", payload.TypeAbandonedRequestsNotification.String()))
	defer span.End()

	done, err := h.dep.WriteAccessor.Begin(ctx, flow.Pulse(ctx))
	defer done()

	if err != nil {
		return nil
	}

	procInitializeExecutionState := initializeAbandonedRequestsNotificationExecutionState{
		dep: h.dep,
		msg: abandoned,
	}
	err = f.Procedure(ctx, &procInitializeExecutionState, false)

	if err != nil {
		return err
	}
	go h.dep.Sender.Reply(ctx, h.meta, bus.ReplyAsMessage(ctx, &reply.OK{}))
	return nil
}
