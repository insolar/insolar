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

	"go.opencensus.io/trace"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

type HandleAdditionalCallFromPreviousExecutor struct {
	dep *Dependencies

	Message payload.Meta
	Parcel  insolar.Parcel
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
	_ flow.Flow,
) {
	broker := h.dep.StateStorage.UpsertExecutionState(msg.ObjectReference)

	if msg.Pending == insolar.NotPending {
		broker.SetNotPending(ctx)
	}

	tr := NewTranscript(freshContextFromContext(ctx), msg.RequestRef, msg.Request)
	broker.AddAdditionalRequestFromPrevExecutor(ctx, tr)
}

func (h *HandleAdditionalCallFromPreviousExecutor) Present(ctx context.Context, f flow.Flow) error {
	ctx = loggerWithTargetID(ctx, h.Parcel)
	inslogger.FromContext(ctx).Debug("HandleAdditionalCallFromPreviousExecutor.Present starts ...")

	msg := h.Parcel.Message().(*message.AdditionalCallFromPreviousExecutor)

	ctx, span := instracer.StartSpan(ctx, "HandleAdditionalCallFromPreviousExecutor.Present")
	span.AddAttributes(
		trace.StringAttribute("msg.Type", msg.Type().String()),
	)
	defer span.End()

	h.handleActual(ctx, msg, f)

	// we never return any other replies
	repMsg := bus.ReplyAsMessage(ctx, &reply.OK{})
	h.dep.Sender.Reply(ctx, h.Message, repMsg)

	return nil
}
