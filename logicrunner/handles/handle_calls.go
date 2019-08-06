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

package handles

import (
	"context"
	"fmt"

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
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/procs"
	"github.com/insolar/insolar/logicrunner/transcript"
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
) {
	// If the flow has canceled during ClarifyPendingState there are two possibilities.
	// 1. It's possible that we were addding request to broker, the pulse has
	//    changed and the execution queue was sent to the next executor.
	//    This means that the next executor already queue element, it's OK.
	// 2. It's also possible that the pulse has changed after registration, but
	//    before adding an item to the execution queue. In this case the queue was sent to the
	//    next executor without the last item. We could just return ErrCanceled to make the
	//    caller to resend the request. However this will cause a slow request deduplication
	//    process on the LME side (it will be caused anyway by another modifying request if
	//    the object is used a lot, but many requests are read-only and don't cause deduplication).
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
		ServiceData:     common.ServiceDataFromContext(ctx),
	}

	_, err := h.dep.MessageBus.Send(ctx, &additionalCallMsg, nil)
	if err != nil {
		logger.Error("[ HandleCall.handleActual.sendToNextExecutor ] mb.Send failed to send AdditionalCallFromPreviousExecutor, ", err)
	}
}

func (h *HandleCall) checkExecutionLoop(
	ctx context.Context, reqRef insolar.Reference, request record.IncomingRequest,
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

	archive := h.dep.StateStorage.GetExecutionArchive(*request.Object)
	if archive == nil {
		return false
	}

	if !archive.FindRequestLoop(ctx, reqRef, request.APIRequestID) {
		return false
	}

	inslogger.FromContext(ctx).Error("loop detected")
	return true
}

func (h *HandleCall) handleActual(
	ctx context.Context,
	msg *message.CallMethod,
	f flow.Flow,
) (insolar.Reply, error) {

	procCheckRole := procs.CheckOurRole{
		Msg:            msg,
		Role:           insolar.DynamicRoleVirtualExecutor,
		JetCoordinator: h.dep.JetCoordinator,

		PulseNumber: flow.Pulse(ctx),
	}

	if err := f.Procedure(ctx, &procCheckRole, true); err != nil {
		// rewrite "can't execute this object" to "flow cancelled" for force retry message
		// just temporary fix till mb moved to watermill
		if err == flow.ErrCancelled || err == procs.ErrCantExecute {
			return nil, flow.ErrCancelled
		}
		return nil, errors.Wrap(err, "[ HandleCall.handleActual ] can't play role")
	}

	request := msg.IncomingRequest

	procRegisterRequest := procs.NewRegisterIncomingRequest(request, h.dep.ArtifactManager)
	if err := f.Procedure(ctx, procRegisterRequest, true); err != nil {
		if err == flow.ErrCancelled {
			inslogger.FromContext(ctx).Info("pulse change during registration, asking caller for retry")
			// Requests need to be deduplicated. For now in case of ErrCancelled we may have 2 registered requests
			return nil, err // message bus will retry on the calling side in ContractRequester
		}
		return nil, errors.Wrap(err, "[ HandleCall.handleActual ] can't create request")
	}

	reqInfo := procRegisterRequest.Result()
	requestRef := insolar.NewReference(reqInfo.RequestID)

	ctx, logger := inslogger.WithField(ctx, "request", requestRef.String())
	logger.Debug("registered request")

	objRef := request.Object
	if request.CallType != record.CTMethod {
		objRef = requestRef
	}
	if objRef == nil {
		return nil, errors.New("can't get object reference")
	}
	if !objRef.Record().Equal(reqInfo.ObjectID) {
		return nil, errors.New("object id we calculated doesn't match ledger")
	}

	registeredRequestReply := &reply.RegisterRequest{Request: *requestRef}

	if len(reqInfo.Request) != 0 {
		logger.Debug("duplicated request")
	}

	if len(reqInfo.Result) != 0 {
		logger.Debug("request has result already")
		go func() {
			err := h.sendRequestResult(ctx, *objRef, *requestRef, request, *reqInfo)
			if err != nil {
				logger.Error("couldn't send request result: ", err.Error())
			}
		}()
		return registeredRequestReply, nil
	}

	if h.checkExecutionLoop(ctx, *requestRef, request) {
		return nil, errors.New("loop detected")
	}

	done, err := h.dep.WriteAccessor.Begin(ctx, flow.Pulse(ctx))
	defer done()

	if err != nil {
		go h.sendToNextExecutor(ctx, *objRef, *requestRef, request)
		return registeredRequestReply, nil
	}

	broker := h.dep.StateStorage.UpsertExecutionState(*objRef)

	proc := procs.AddFreshRequest{Broker: broker, RequestRef: *requestRef, Request: request}
	if err := f.Procedure(ctx, &proc, true); err != nil {
		return nil, errors.Wrap(err, "couldn't pass request to broker")
	}

	return registeredRequestReply, nil
}

func (h *HandleCall) Present(ctx context.Context, f flow.Flow) error {
	ctx = common.LoggerWithTargetID(ctx, h.Parcel)
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

	if err != nil {
		return sendErrorMessage(ctx, h.dep.Sender, h.Message, err)
	}
	go h.dep.Sender.Reply(ctx, h.Message, bus.ReplyAsMessage(ctx, rep))
	return nil
}

func (h *HandleCall) sendRequestResult(
	ctx context.Context,
	objRef insolar.Reference,
	reqRef insolar.Reference,
	request record.IncomingRequest,
	reqInfo payload.RequestInfo,
) error {
	logger := inslogger.FromContext(ctx)
	logger.Debug("sending earlier")

	rec := record.Material{}
	err := rec.Unmarshal(reqInfo.Result)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal record")
	}
	virtual := record.Unwrap(&rec.Virtual)
	resultRecord, ok := virtual.(*record.Result)
	if !ok {
		return fmt.Errorf("unexpected record %T", virtual)
	}

	repl := &reply.CallMethod{Result: resultRecord.Payload, Object: &objRef}
	tr := transcript.NewTranscript(ctx, reqRef, request)
	h.dep.RequestsExecutor.SendReply(ctx, tr, repl, nil)

	return nil
}
