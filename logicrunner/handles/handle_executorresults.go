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

	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/transcript"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

type initializeExecutionState struct {
	dep *Dependencies
	msg *message.ExecutorResults
}

func (p *initializeExecutionState) Proceed(ctx context.Context) error {
	ref := p.msg.GetReference()

	broker := p.dep.StateStorage.UpsertExecutionState(ref)
	broker.PrevExecutorPendingResult(ctx, p.msg.Pending)

	if p.msg.LedgerHasMoreRequests {
		broker.MoreRequestsOnLedger(ctx)
	}

	if len(p.msg.Queue) > 0 {
		transcripts := make([]*transcript.Transcript, len(p.msg.Queue))
		for i, qe := range p.msg.Queue {
			requestCtx := common.ContextFromServiceData(qe.ServiceData)
			transcripts[i] = transcript.NewTranscript(requestCtx, qe.RequestRef, qe.Request)
		}
		broker.AddRequestsFromPrevExecutor(ctx, transcripts...)
	}

	return nil
}

type HandleExecutorResults struct {
	dep *Dependencies

	Message payload.Meta
	Parcel  insolar.Parcel
}

func (h *HandleExecutorResults) realHandleExecutorState(ctx context.Context, f flow.Flow) error {
	msg := h.Parcel.Message().(*message.ExecutorResults)

	procInitializeExecutionState := initializeExecutionState{
		dep: h.dep,
		msg: msg,
	}
	err := f.Procedure(ctx, &procInitializeExecutionState, true)
	if err != nil {
		if err == flow.ErrCancelled {
			return nil
		}
		err := errors.Wrap(err, "[ HandleExecutorResults ] Failed to initialize execution state")
		return err
	}

	return nil
}

func (h *HandleExecutorResults) Present(ctx context.Context, f flow.Flow) error {
	ctx = common.LoggerWithTargetID(ctx, h.Parcel)
	logger := inslogger.FromContext(ctx)

	logger.Debug("HandleExecutorResults.Present starts ...")

	msg, ok := h.Parcel.Message().(*message.ExecutorResults)
	if !ok {
		return errors.New("HandleExecutorResults( ! message.ExecutorResults )")
	}

	ctx, span := instracer.StartSpan(ctx, "HandleExecutorResults.Present")
	span.AddAttributes(trace.StringAttribute("msg.Type", msg.Type().String()))
	defer span.End()

	err := h.realHandleExecutorState(ctx, f)

	if err != nil {
		return sendErrorMessage(ctx, h.dep.Sender, h.Message, err)
	}
	go h.dep.Sender.Reply(ctx, h.Message, bus.ReplyAsMessage(ctx, &reply.OK{}))
	return nil
}
