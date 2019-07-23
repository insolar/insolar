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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/pkg/errors"
)

type initializeAbandonedRequestsNotificationExecutionState struct {
	LR  *LogicRunner
	msg payload.AbandonedRequestsNotification
}

// Proceed initializes or sets LedgerHasMoreRequests to right value
func (p *initializeAbandonedRequestsNotificationExecutionState) Proceed(ctx context.Context) error {
	ref := *insolar.NewReference(p.msg.ObjectID)

	broker := p.LR.StateStorage.UpsertExecutionState(ref)

	broker.executionState.Lock()
	if broker.executionState.pending == insolar.PendingUnknown {
		broker.executionState.pending = insolar.InPending
		broker.executionState.PendingConfirmed = false
	}
	broker.executionState.Unlock()

	broker.MoreRequestsOnLedger(ctx)
	broker.FetchMoreRequestsFromLedger(ctx)

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

	return nil

	// FIXME: uncomment and fix this (if needed) after enabling abandoned requests
	// ctx, _ = inslogger.WithField(ctx, "targetid", abandoned.ObjectID.String())
	// logger := inslogger.FromContext(ctx)
	//
	// logger.Debug("HandleAbandonedRequestsNotification.Present starts ...")
	//
	// ctx, span := instracer.StartSpan(ctx, "HandleAbandonedRequestsNotification.Present")
	// span.AddAttributes(trace.StringAttribute("msg.Type", payload.TypeAbandonedRequestsNotification.String()))
	// defer span.End()
	//
	// procInitializeExecutionState := initializeAbandonedRequestsNotificationExecutionState{
	// 	LR:  h.dep.lr,
	// 	msg: abandoned,
	// }
	// if err := f.Procedure(ctx, &procInitializeExecutionState, false); err != nil {
	// 	err := errors.Wrap(err, "[ HandleExecutorResults ] Failed to initialize execution state")
	// 	rep, newErr := payload.NewMessage(&payload.Error{Text: err.Error()})
	// 	if newErr != nil {
	// 		return newErr
	// 	}
	// 	go h.dep.Sender.Reply(ctx, h.meta, rep)
	// 	return err
	// }
	// replyOk := bus.ReplyAsMessage(ctx, &reply.OK{})
	// go h.dep.Sender.Reply(ctx, h.meta, replyOk)
	// return nil
}
