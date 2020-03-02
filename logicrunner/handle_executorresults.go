// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package logicrunner

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/logicrunner/writecontroller"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

type HandleExecutorResults struct {
	dep *Dependencies

	meta payload.Meta
}

func checkPayloadExecutorResults(ctx context.Context, results payload.ExecutorResults) error {
	if !results.Caller.IsEmpty() && !results.Caller.IsObjectReference() {
		return errors.Errorf("results.Caller should be ObjectReference; ref=%s", results.Caller.String())
	}
	if !results.RecordRef.IsObjectReference() {
		return errors.Errorf("results.RecordRef should be ObjectReference; ref=%s", results.RecordRef.String())
	}

	return nil
}

func (h *HandleExecutorResults) Present(ctx context.Context, f flow.Flow) error {
	message := payload.ExecutorResults{}
	err := message.Unmarshal(h.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal message")
	}

	ctx, logger := inslogger.WithField(ctx, "object", message.RecordRef.String())
	logger.Debug("handling ExecutorResults")

	ctx, span := instracer.StartSpan(ctx, "HandleExecutorResults.Present")
	defer span.Finish()

	if err := checkPayloadExecutorResults(ctx, message); err != nil {
		return err
	}

	return h.handleMessage(ctx, message)
}

func (h *HandleExecutorResults) handleMessage(ctx context.Context, msg payload.ExecutorResults) error {
	done, err := h.dep.WriteAccessor.Begin(ctx, flow.Pulse(ctx))
	if err != nil {
		if err == writecontroller.ErrWriteClosed {
			return flow.ErrCancelled
		}
		return errors.Wrap(err, "failed to acquire write access")
	}
	defer done()

	broker := h.dep.StateStorage.UpsertExecutionState(msg.RecordRef)
	broker.PrevExecutorPendingResult(ctx, msg.Pending)

	if msg.LedgerHasMoreRequests {
		broker.HasMoreRequests(ctx)
	}

	return nil
}
