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
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

type HandleCall struct {
	dep *Dependencies

	Message bus.Message
}

func (h *HandleCall) sendToNextExecutor(
	ctx context.Context,
	es *ExecutionState,
	broker *ExecutionBroker,
	requestRef *Ref,
) {
	// If the flow has canceled during ClarifyPendingState there are two possibilities.
	// 1. It's possible that we started to execute ClarifyPendingState, the pulse has
	// changed and the execution queue was sent to the next executor in OnPulse method.
	// This means that the next executor already has the last queue element, it's OK.
	// 2. It's also possible that the pulse has changed after RegisterIncomingRequest but
	// before adding an item to the execution queue. In this case the queue was sent to the
	// next executor without the last item. We could just return ErrCanceled to make the
	// caller to resend the request. However this will cause a slow request deduplication
	// process on the LME side (it will be caused anyway by another modifying request if
	// the object is used a lot, but many requests are read-only and don't cause deduplication).
	// As an optimization we decided to send a special message type to the next executor.
	// It's possible that while we send the message the pulse will change once again and
	// the receiver will be not an executor of the object anymore. However in this case
	// MessageBus will automatically resend the message to the right VE.

	logger := inslogger.FromContext(ctx)

	es.Lock()
	// We want to remove element we have just added to es.Queue to eliminate doubling
	transcript := broker.GetByReference(ctx, requestRef)
	es.Unlock()

	// it might be already collected in OnPulse, that is why it already might not be in es.Queue
	if transcript != nil {
		logger.Debug("Sending additional request to next executor")
		additionalCallMsg := message.AdditionalCallFromPreviousExecutor{
			ObjectReference: es.Ref,
			RequestRef:      *requestRef,
			Request:         *transcript.Request,
		}
		if es.pending == message.PendingUnknown {
			additionalCallMsg.Pending = message.NotPending
		} else {
			additionalCallMsg.Pending = es.pending
		}

		if _, err := h.dep.lr.MessageBus.Send(ctx, &additionalCallMsg, nil); err != nil {
			logger.Error("[ HandleCall.handleActual.sendToNextExecutor ] mb.Send failed to send AdditionalCallFromPreviousExecutor, ", err)
		}
	}
}

func (h *HandleCall) handleActual(
	ctx context.Context,
	msg *message.CallMethod,
	f flow.Flow,
) (insolar.Reply, error) {

	lr := h.dep.lr

	procCheckRole := CheckOurRole{
		msg:         msg,
		role:        insolar.DynamicRoleVirtualExecutor,
		lr:          lr,
		pulseNumber: flow.Pulse(ctx),
	}

	if err := f.Procedure(ctx, &procCheckRole, true); err != nil {
		// rewrite "can't execute this object" to "flow cancelled" for force retry message
		// just temporary fix till mb moved to watermill
		if err == flow.ErrCancelled || err == ErrCantExecute {
			return nil, flow.ErrCancelled
		}
		return nil, errors.Wrap(err, "[ HandleCall.handleActual ] can't play role")
	}

	request := msg.IncomingRequest

	if lr.CheckExecutionLoop(ctx, request) {
		return nil, errors.New("loop detected")
	}

	procRegisterRequest := NewRegisterIncomingRequest(request, h.dep)

	if err := f.Procedure(ctx, procRegisterRequest, true); err != nil {
		if err == flow.ErrCancelled {
			// Requests need to be deduplicated. For now in case of ErrCancelled we may have 2 registered requests
			return nil, err // message bus will retry on the calling side in ContractRequester
		}
		return nil, errors.Wrap(err, "[ HandleCall.handleActual ] can't create request")
	}
	requestRef := procRegisterRequest.getResult()

	objRef := request.Object
	if request.CallType != record.CTMethod {
		objRef = requestRef
	}
	if objRef == nil {
		return nil, errors.New("can't get object reference")
	}

	es, broker := lr.StateStorage.UpsertExecutionState(*objRef)

	es.Lock()
	broker.Put(ctx, false, NewTranscript(ctx, requestRef, request))
	es.Unlock()

	procClarifyPendingState := ClarifyPendingState{
		es:              es,
		request:         &request,
		ArtifactManager: lr.ArtifactManager,
	}

	if err := f.Procedure(ctx, &procClarifyPendingState, true); err != nil {
		if err == flow.ErrCancelled {
			h.sendToNextExecutor(ctx, es, broker, requestRef)
		} else {
			inslogger.FromContext(ctx).Error(" HandleCall.handleActual ] ClarifyPendingState returns error: ", err)
		}
		// and return the reply as usual
	} else {
		// it's 'fast' operation, so we don't need to check that pulse ends
		broker.StartProcessorIfNeeded(ctx)
	}

	return &reply.RegisterRequest{
		Request: *requestRef,
	}, nil
}

func (h *HandleCall) Present(ctx context.Context, f flow.Flow) error {
	parcel := h.Message.Parcel
	ctx = loggerWithTargetID(ctx, parcel)
	inslogger.FromContext(ctx).Debug("HandleCall.Present starts ...")

	msg, ok := parcel.Message().(*message.CallMethod)
	if !ok {
		return errors.New("is not CallMethod message")
	}

	ctx, span := instracer.StartSpan(ctx, "HandleCall.Present")
	span.AddAttributes(
		trace.StringAttribute("msg.Type", msg.Type().String()),
	)
	defer span.End()

	r := bus.Reply{}
	r.Reply, r.Err = h.handleActual(ctx, msg, f)

	h.Message.ReplyTo <- r
	return nil
}

type HandleAdditionalCallFromPreviousExecutor struct {
	dep *Dependencies

	Message bus.Message
}

// This is basically a simplified version of HandleCall.handleActual().
// Please note that currently we lack any fraud detection here.
// Ideally we should check that the previous executor was really an executor
// during previous pulse, that the request was really registered, etc.
// Also we don't handle case when pulse changes during execution of this handle.
// In this scenario user is in a bad luck. The request will be lost and user will have
// to re-send it after some timeout.
func (h *HandleAdditionalCallFromPreviousExecutor) handleActual(
	ctx context.Context,
	msg *message.AdditionalCallFromPreviousExecutor,
	f flow.Flow,
) {
	lr := h.dep.lr
	es, broker := lr.StateStorage.UpsertExecutionState(msg.ObjectReference)

	es.Lock()
	if msg.Pending == message.NotPending {
		es.pending = message.NotPending
	}

	broker.Put(ctx, false, NewTranscript(ctx, &msg.RequestRef, msg.Request))
	es.Unlock()

	procClarifyPendingState := ClarifyPendingState{
		es:              es,
		request:         &msg.Request,
		ArtifactManager: lr.ArtifactManager,
	}

	if err := f.Procedure(ctx, &procClarifyPendingState, true); err != nil {
		inslogger.FromContext(ctx).Warn("[ HandleAdditionalCallFromPreviousExecutor.handleActual ] ClarifyPendingState returns error: ", err)
		// We intentionally report OK to the previous executor here. There is no point
		// in resending the message or anything.
		return
	}

	// it's 'fast' operation, so we don't need to check that pulse ends
	broker.StartProcessorIfNeeded(ctx)
}

func (h *HandleAdditionalCallFromPreviousExecutor) Present(ctx context.Context, f flow.Flow) error {
	parcel := h.Message.Parcel
	ctx = loggerWithTargetID(ctx, parcel)
	inslogger.FromContext(ctx).Debug("HandleAdditionalCallFromPreviousExecutor.Present starts ...")

	msg, ok := parcel.Message().(*message.AdditionalCallFromPreviousExecutor)
	if !ok {
		return errors.New("is not AdditionalCallFromPreviousExecutor message")
	}

	ctx, span := instracer.StartSpan(ctx, "HandleAdditionalCallFromPreviousExecutor.Present")
	span.AddAttributes(
		trace.StringAttribute("msg.Type", msg.Type().String()),
	)
	defer span.End()

	h.handleActual(ctx, msg, f)

	// we never return any other replies
	h.Message.ReplyTo <- bus.Reply{Reply: &reply.OK{}}
	return nil
}
