// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package logicrunner

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/metrics"
	"github.com/insolar/insolar/logicrunner/writecontroller"
	"github.com/insolar/insolar/platformpolicy"

	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/common"
)

func checkPayloadCallMethod(ctx context.Context, callMethod payload.CallMethod) error {
	if err := checkIncomingRequest(ctx, callMethod.Request); err != nil {
		return errors.Wrap(err, "failed to verify callMethod.Request")
	}

	return nil
}

type HandleCall struct {
	dep *Dependencies

	Message payload.Meta
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
	msg, err := payload.NewMessage(&payload.AdditionalCallFromPreviousExecutor{
		ObjectReference: objectRef,
		RequestRef:      requestRef,
		Request:         &request,
		ServiceData:     common.ServiceDataFromContext(ctx),
	})
	if err != nil {
		logger.Error("[ HandleCall.handleActual.sendToNextExecutor ] failed to serialize payload message", err)
	}

	sender := bus.NewWaitOKWithRetrySender(h.dep.Sender, h.dep.PulseAccessor, 1)
	sender.SendRole(ctx, msg, insolar.DynamicRoleVirtualExecutor, objectRef)
}

func (h *HandleCall) handleActual(
	ctx context.Context,
	msg payload.CallMethod,
	f flow.Flow,
) (insolar.Reply, error) {
	var pcs = platformpolicy.NewPlatformCryptographyScheme() // TODO: create message factory
	target := record.CalculateRequestAffinityRef(msg.Request, msg.PulseNumber, pcs)

	request := msg.Request
	ctx, logger := inslogger.WithField(ctx, "method", request.Method)

	procCheckRole := CheckOurRole{
		target:         *target,
		role:           insolar.DynamicRoleVirtualExecutor,
		jetCoordinator: h.dep.JetCoordinator,
		pulseNumber:    flow.Pulse(ctx),
	}

	if err := f.Procedure(ctx, &procCheckRole, true); err != nil {
		// rewrite "can't execute this object" to "flow cancelled" for force retry message
		// just temporary fix till mb moved to watermill
		if err == flow.ErrCancelled || err == ErrCantExecute {
			return nil, flow.ErrCancelled
		}
		return nil, errors.Wrap(err, "[ HandleCall.handleActual ] can't play role")
	}

	logger.Debug("registering incoming request")

	procRegisterRequest := NewRegisterIncomingRequest(*request, h.dep)
	err := f.Procedure(ctx, procRegisterRequest, true)
	if err != nil {
		logger.WithField("error", err.Error()).Debug("failed to register incoming request")
		if err == flow.ErrCancelled {
			inslogger.FromContext(ctx).Info("pulse change during registration, asking caller for retry")
			// Requests need to be deduplicated. For now in case of ErrCancelled we may have 2 registered requests
			return nil, err // message bus will retry on the calling side in ContractRequester
		}
		if isLogicalError := ProcessLogicalError(ctx, err); isLogicalError {
			inslogger.FromContext(ctx).Warn("request to not existing object")

			resultWithErr, err := foundation.MarshalMethodErrorResult(err)
			if err != nil {
				return nil, errors.Wrap(err, "can't create error result")
			}
			stats.Record(ctx, metrics.CallMethodLogicalError.M(1))
			return &reply.CallMethod{Result: resultWithErr}, nil
		}
		return nil, errors.Wrap(err, "[ HandleCall.handleActual ] can't create request")
	}
	logger.Debug("registered request")

	reqInfo := procRegisterRequest.getResult()
	requestRef := *getRequestReference(reqInfo)

	if request.CallType != record.CTMethod {
		request.Object = insolar.NewReference(reqInfo.RequestID)
	}

	objectRef := request.Object

	if objectRef == nil || !objectRef.IsSelfScope() {
		logger.Debug("incoming request bad object reference")
		return nil, errors.New("can't get object reference")
	}

	ctx, logger = inslogger.WithFields(
		ctx,
		map[string]interface{}{
			"object":  objectRef.String(),
			"request": requestRef.String(),
		},
	)

	logger.Debug("registered incoming request")

	if !objectRef.GetLocal().Equal(reqInfo.ObjectID) {
		logger.Debug("incoming request invalid object reference")
		return nil, errors.New("object id we calculated doesn't match ledger")
	}

	registeredRequestReply := &reply.RegisterRequest{Request: requestRef}

	if len(reqInfo.Request) != 0 {
		logger.Debug("duplicated request")
	}

	if len(reqInfo.Result) != 0 {
		logger.Debug("incoming request already has result on ledger, returning it")
		go func() {
			err := h.sendRequestResult(ctx, *objectRef, requestRef, *request, *reqInfo)
			if err != nil {
				logger.Error("couldn't send request result: ", err.Error())
			}
		}()
		return registeredRequestReply, nil
	}

	done, err := h.dep.WriteAccessor.Begin(ctx, flow.Pulse(ctx))
	if err != nil {
		logger.WithField("error", err).Debug("failed to acquire write accessor")
		if err == writecontroller.ErrWriteClosed {
			stats.Record(ctx, metrics.CallMethodAdditionalCall.M(1))
			go h.sendToNextExecutor(ctx, *objectRef, requestRef, *request)
			return registeredRequestReply, nil
		}
		return nil, errors.Wrap(err, "failed to acquire write access")
	}
	defer done()

	broker := h.dep.StateStorage.UpsertExecutionState(*objectRef)
	broker.HasMoreRequests(ctx)

	return registeredRequestReply, nil
}

func (h *HandleCall) Present(ctx context.Context, f flow.Flow) error {
	inslogger.FromContext(ctx).Debug("HandleCall.Present starts ...")

	message := payload.CallMethod{}
	err := message.Unmarshal(h.Message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal message")
	}

	if err := checkPayloadCallMethod(ctx, message); err != nil {
		return err
	}

	rep, err := h.handleActual(ctx, message, f)
	if err != nil {
		return err
	}

	h.dep.Sender.Reply(ctx, h.Message, bus.ReplyAsMessage(ctx, rep))

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
	logger.Debug("sending earlier computed result")

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
	h.dep.RequestsExecutor.SendReply(ctx, reqRef, request, repl, nil)

	return nil
}
