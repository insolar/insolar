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

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

type initializeAbandonedRequestsNotificationExecutionState struct {
	LR  *LogicRunner
	msg *message.AbandonedRequestsNotification
}

// Proceed initializes or sets LedgerHasMoreRequests to right value
func (p *initializeAbandonedRequestsNotificationExecutionState) Proceed(ctx context.Context) error {
	ref := *p.msg.DefaultTarget()

	state := p.LR.UpsertObjectState(ref)

	state.Lock()
	if state.ExecutionState == nil {
		state.ExecutionState = NewExecutionState(ref)
		state.ExecutionState.pending = message.InPending
		state.ExecutionState.PendingConfirmed = false
		state.ExecutionState.LedgerHasMoreRequests = true
	} else {
		executionState := state.ExecutionState
		executionState.Lock()
		executionState.LedgerHasMoreRequests = true
		executionState.Unlock()

	}
	state.Unlock()

	return nil
}

type HandleAbandonedRequestsNotification struct {
	dep *Dependencies

	Message bus.Message
}

func (h *HandleAbandonedRequestsNotification) Present(ctx context.Context, f flow.Flow) error {
	parcel := h.Message.Parcel
	ctx = loggerWithTargetID(ctx, parcel)
	logger := inslogger.FromContext(ctx)

	logger.Debug("HandleAbandonedRequestsNotification.Present starts ...")

	msg, ok := parcel.Message().(*message.AbandonedRequestsNotification)
	if !ok {
		return errors.New("HandleAbandonedRequestsNotification( ! message.AbandonedRequestsNotification )")
	}

	ctx, span := instracer.StartSpan(ctx, "HandleAbandonedRequestsNotification.Present")
	span.AddAttributes(trace.StringAttribute("msg.Type", msg.Type().String()))
	defer span.End()

	procInitializeExecutionState := initializeAbandonedRequestsNotificationExecutionState{
		LR:  h.dep.lr,
		msg: msg,
	}
	if err := f.Procedure(ctx, &procInitializeExecutionState, false); err != nil {
		err := errors.Wrap(err, "[ HandleExecutorResults ] Failed to initialize execution state")
		h.Message.ReplyTo <- bus.Reply{Reply: &reply.Error{}, Err: err}
		return err
	}

	h.Message.ReplyTo <- bus.Reply{Reply: &reply.OK{}, Err: nil}
	return nil
}
