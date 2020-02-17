// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package logicrunner

import (
	"context"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/metrics"
	"github.com/insolar/insolar/logicrunner/writecontroller"
)

func checkPayloadPendingFinished(finished payload.PendingFinished) error {
	if finished.ObjectRef.IsEmpty() {
		return errors.New("Got PendingFinished message, but field ObjectRef is empty")
	}
	if !finished.ObjectRef.IsObjectReference() {
		return errors.Errorf("PendingFinished.ObjectRef should be RecordReference; ref=%s", finished.ObjectRef.String())
	}

	return nil
}

type HandlePendingFinished struct {
	dep *Dependencies

	Message payload.Meta
}

func (h *HandlePendingFinished) Present(ctx context.Context, _ flow.Flow) error {
	message := payload.PendingFinished{}
	err := message.Unmarshal(h.Message.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal message")
	}

	ctx, logger := inslogger.WithFields(ctx, map[string]interface{}{
		"object": message.ObjectRef.String(),
		"sender": h.Message.Sender.String(),
	})
	logger.Debug("handle PendingFinished message")

	if err := checkPayloadPendingFinished(message); err != nil {
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

	broker := h.dep.StateStorage.UpsertExecutionState(message.ObjectRef)

	err = broker.PrevExecutorSentPendingFinished(ctx)
	if err != nil {
		if err == ErrAlreadyExecuting {
			stats.Record(ctx, metrics.PendingFinishedAlreadyExecuting.M(1))
		}
		err = errors.Wrap(err, "handle PendingFinished failed")
		inslogger.FromContext(ctx).Error(err.Error())
		return err
	}

	replyOk := bus.ReplyAsMessage(ctx, &reply.OK{})
	h.dep.Sender.Reply(ctx, h.Message, replyOk)
	return nil
}
