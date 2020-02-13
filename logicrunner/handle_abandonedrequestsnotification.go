// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package logicrunner

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

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

	objRef := *insolar.NewReference(abandoned.ObjectID)

	ctx, logger := inslogger.WithField(ctx, "object", objRef.String())

	logger.Debug("got abandoned requests notification")

	ctx, span := instracer.StartSpan(ctx, "HandleAbandonedRequestsNotification.Present")
	span.SetTag("msg.Type", payload.TypeAbandonedRequestsNotification.String())
	defer span.Finish()

	done, err := h.dep.WriteAccessor.Begin(ctx, flow.Pulse(ctx))
	if err != nil {
		logger.Warn("late notification about abandoned, ignoring: ", err.Error())
		return nil
	}
	defer done()

	broker := h.dep.StateStorage.UpsertExecutionState(objRef)
	broker.AbandonedRequestsOnLedger(ctx)

	return nil
}
