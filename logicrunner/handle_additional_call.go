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

type AdditionalCallFromPreviousExecutor struct {
	stateStorage StateStorage

	message *message.AdditionalCallFromPreviousExecutor
}

func (p *AdditionalCallFromPreviousExecutor) Proceed(ctx context.Context) error {
	broker := p.stateStorage.UpsertExecutionState(p.message.ObjectReference)

	if p.message.Pending == insolar.NotPending {
		broker.SetNotPending(ctx)
	}

	tr := NewTranscript(freshContextFromContext(ctx), p.message.RequestRef, p.message.Request)
	broker.AddAdditionalRequestFromPrevExecutor(ctx, tr)
	return nil
}

type HandleAdditionalCallFromPreviousExecutor struct {
	dep *Dependencies

	Message payload.Meta
	Parcel  insolar.Parcel
}

// Please note that currently we lack any fraud detection here.
// Ideally we should check that the previous executor was really an executor during previous pulse,
// that the request was really registered, etc. Also we don't handle case when pulse changes during
// execution of this handle. In this scenario user is in a bad luck. The request will be lost and
// user will have to re-send it after some timeout.
func (h *HandleAdditionalCallFromPreviousExecutor) Present(ctx context.Context, f flow.Flow) error {
	ctx = loggerWithTargetID(ctx, h.Parcel)

	logger := inslogger.FromContext(ctx).WithField("handler", "HandleAdditionalCallFromPreviousExecutor")
	logger.Debug("Handler.Present starts")

	msg := h.Parcel.Message().(*message.AdditionalCallFromPreviousExecutor)
	ctx = contextWithServiceData(ctx, msg.ServiceData)

	ctx, span := instracer.StartSpan(ctx, "HandleAdditionalCallFromPreviousExecutor.Present")
	span.AddAttributes(
		trace.StringAttribute("msg.Type", msg.Type().String()),
	)
	defer span.End()

	done, err := h.dep.WriteAccessor.Begin(ctx, flow.Pulse(ctx))
	if err != nil { // pulse changed, send that message to next executor
		_, err := h.dep.lr.MessageBus.Send(ctx, msg, nil)
		if err != nil {
			logger.Error("MessageBus.Send failed: ", err)
		}
		return nil
	}
	defer done()

	proc := &AdditionalCallFromPreviousExecutor{
		stateStorage: h.dep.StateStorage,
		message:      msg,
	}
	if err := f.Procedure(ctx, proc, false); err != nil {
		return err
	}

	// we never return any other replies
	h.dep.Sender.Reply(ctx, h.Message, bus.ReplyAsMessage(ctx, &reply.OK{}))
	return nil
}
