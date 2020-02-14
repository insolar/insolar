// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package logicrunner

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/writecontroller"
)

type AdditionalCallFromPreviousExecutor struct {
	stateStorage StateStorage

	message *payload.AdditionalCallFromPreviousExecutor
}

func (p *AdditionalCallFromPreviousExecutor) Proceed(ctx context.Context) error {
	broker := p.stateStorage.UpsertExecutionState(p.message.ObjectReference)

	if p.message.Pending == insolar.NotPending {
		broker.SetNotPending(ctx)
	}
	broker.HasMoreRequests(ctx)

	return nil
}

func checkPayloadAdditionalCallFromPreviousExecutor(ctx context.Context, msg payload.AdditionalCallFromPreviousExecutor) error {
	if !msg.ObjectReference.IsObjectReference() {
		return errors.Errorf("StillExecuting.ObjectReference should be ObjectReference; ref=%s", msg.ObjectReference.String())
	}
	if !msg.RequestRef.IsRecordScope() {
		return errors.Errorf("StillExecuting.RequestRef should be RecordReference; ref=%s", msg.RequestRef.String())
	}
	if err := checkIncomingRequest(ctx, msg.Request); err != nil {
		return errors.Wrap(err, "failed to check IncomingRequest of AdditionalCallFromPreviousExecutor")
	}
	return nil
}

type HandleAdditionalCallFromPreviousExecutor struct {
	dep *Dependencies

	Message payload.Meta
}

// Please note that currently we lack any fraud detection here.
// Ideally we should check that the previous executor was really an executor during previous pulse,
// that the request was really registered, etc. Also we don't handle case when pulse changes during
// execution of this handle. In this scenario user is in a bad luck. The request will be lost and
// user will have to re-send it after some timeout.
func (h *HandleAdditionalCallFromPreviousExecutor) Present(ctx context.Context, f flow.Flow) error {
	ctx, logger := inslogger.WithFields(ctx, map[string]interface{}{
		"handler": "HandleAdditionalCallFromPreviousExecutor",
	})

	logger.Debug("Handler.Present starts")

	message := payload.AdditionalCallFromPreviousExecutor{}
	err := message.Unmarshal(h.Message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal message")
	}

	ctx = contextWithServiceData(ctx, message.ServiceData)

	ctx, _ = inslogger.WithFields(ctx, map[string]interface{}{
		"object":  message.ObjectReference.String(),
		"request": message.Request.String(),
	})

	ctx = contextWithServiceData(ctx, message.ServiceData)

	if err := checkPayloadAdditionalCallFromPreviousExecutor(ctx, message); err != nil {
		return err
	}

	done, err := h.dep.WriteAccessor.Begin(ctx, flow.Pulse(ctx))
	if err != nil {
		if err == writecontroller.ErrWriteClosed {
			return flow.ErrCancelled
		}
		return errors.Wrap(err, "failed to acquire write access")
	}
	defer done()

	proc := &AdditionalCallFromPreviousExecutor{
		stateStorage: h.dep.StateStorage,
		message:      &message,
	}
	if err := f.Procedure(ctx, proc, false); err != nil {
		return err
	}

	// we never return any other replies
	h.dep.Sender.Reply(ctx, h.Message, bus.ReplyAsMessage(ctx, &reply.OK{}))
	return nil
}
