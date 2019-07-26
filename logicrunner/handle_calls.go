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

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

type HandleCall struct {
	dep *Dependencies

	Message payload.Meta
	Parcel  insolar.Parcel
}

func (h *HandleCall) sendToNextExecutor(
	ctx context.Context,
	objectRef insolar.Reference,
	requestRef insolar.Reference,
	request record.IncomingRequest,
	ps insolar.PendingState,
) {
	// If the flow has canceled during ClarifyPendingState there are two possibilities.
	// 1. It's possible that we were addding request to broker, the pulse has
	// changed and the execution queue was sent to the next executor.
	// This means that the next executor already queue element, it's OK.
	// 2. It's also possible that the pulse has changed after registration, but
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

	logger.Debug("Sending additional request to next executor")
	additionalCallMsg := message.AdditionalCallFromPreviousExecutor{
		ObjectReference: objectRef,
		RequestRef:      requestRef,
		Request:         request,
		ServiceData:     serviceDataFromContext(ctx),
	}

	if ps == insolar.PendingUnknown {
		additionalCallMsg.Pending = insolar.NotPending
	} else {
		additionalCallMsg.Pending = ps
	}

	_, err := h.dep.lr.MessageBus.Send(ctx, &additionalCallMsg, nil)
	if err != nil {
		logger.Error("[ HandleCall.handleActual.sendToNextExecutor ] mb.Send failed to send AdditionalCallFromPreviousExecutor, ", err)
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

	if h.checkExecutionLoop(ctx, request) {
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

	ctx, logger := inslogger.WithField(ctx, "request", requestRef.String())
	logger.Debug("Registered request")

	objRef := request.Object
	if request.CallType != record.CTMethod {
		objRef = requestRef
	}
	if objRef == nil {
		return nil, errors.New("can't get object reference")
	}

	broker := lr.StateStorage.UpsertExecutionState(*objRef)

	proc := AddFreshRequest{broker: broker, requestRef: *requestRef, request: request}
	err := f.Procedure(ctx, &proc, true)
	if err != nil {
		if err == flow.ErrCancelled {
			h.sendToNextExecutor(ctx, *objRef, *requestRef, request, broker.PendingState())
		}
		return nil, errors.Wrap(err, "couldn't pass request to broker")
	}

	return &reply.RegisterRequest{
		Request: *requestRef,
	}, nil
}

func (h *HandleCall) Present(ctx context.Context, f flow.Flow) error {
	ctx = loggerWithTargetID(ctx, h.Parcel)
	inslogger.FromContext(ctx).Debug("HandleCall.Present starts ...")

	msg, ok := h.Parcel.Message().(*message.CallMethod)
	if !ok {
		return errors.New("is not CallMethod message")
	}

	ctx, span := instracer.StartSpan(ctx, "HandleCall.Present")
	span.AddAttributes(
		trace.StringAttribute("msg.Type", msg.Type().String()),
	)
	defer span.End()

	rep, err := h.handleActual(ctx, msg, f)

	var repMsg *watermillMsg.Message
	if err != nil {
		var newErr error
		repMsg, newErr = payload.NewMessage(&payload.Error{Text: err.Error()})
		if newErr != nil {
			return newErr
		}
	} else {
		repMsg = bus.ReplyAsMessage(ctx, rep)
	}
	go h.dep.Sender.Reply(ctx, h.Message, repMsg)

	return nil
}

func (h *HandleCall) checkExecutionLoop(
	ctx context.Context, request record.IncomingRequest,
) bool {

	if request.ReturnMode == record.ReturnNoWait {
		return false
	}
	if request.CallType != record.CTMethod {
		return false
	}
	if request.Object == nil {
		// should be catched by other code
		return false
	}

	broker := h.dep.StateStorage.GetExecutionState(*request.Object)
	if broker == nil {
		return false
	}

	if !broker.CheckExecutionLoop(ctx, request.APIRequestID) {
		return false
	}

	inslogger.FromContext(ctx).Error("loop detected")
	return true
}


