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
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

type initializeExecutionState struct {
	LR  *LogicRunner
	msg *message.ExecutorResults

	Result struct {
		es             *ExecutionState
		clarifyPending bool
	}
}

func (p *initializeExecutionState) Proceed(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	ref := p.msg.GetReference()

	state := p.LR.StateStorage.UpsertObjectState(ref)

	state.Lock()
	if state.ExecutionState == nil {
		state.ExecutionState = NewExecutionState(ref)
		state.ExecutionState.RegisterLogicRunner(p.LR)
	}
	es := state.ExecutionState
	p.Result.es = es
	state.Unlock()

	p.Result.clarifyPending = false

	es.Lock()
	if es.pending == message.InPending {
		if !es.Broker.currentList.Empty() {
			logger.Debug("execution returned to node that is still executing pending")

			es.pending = message.NotPending
			es.PendingConfirmed = false
		} else if p.msg.Pending == message.NotPending {
			logger.Debug("executor we came to thinks that execution pending, but previous said to continue")

			es.pending = message.NotPending
		}
	} else if es.pending == message.PendingUnknown {
		es.pending = p.msg.Pending
		logger.Debug("pending state was unknown, setting from previous executor to ", es.pending)

		if es.pending == message.PendingUnknown {
			p.Result.clarifyPending = true
		}
	}

	// set false to true is good, set true to false may be wrong, better make unnecessary call
	if !es.LedgerHasMoreRequests {
		es.LedgerHasMoreRequests = p.msg.LedgerHasMoreRequests
	}

	// prepare Queue
	if p.msg.Queue != nil {
		for _, qe := range p.msg.Queue {
			ctxToSent := context.TODO() //FIXME: !!!!
			transcript := NewTranscript(ctxToSent, &qe.RequestRef, qe.Request)

			es.Broker.Prepend(ctx, false, transcript)
		}
	}
	es.Unlock()

	return nil
}

type HandleExecutorResults struct {
	dep *Dependencies

	Message payload.Meta
	Parcel  insolar.Parcel
}

func (h *HandleExecutorResults) realHandleExecutorState(ctx context.Context, f flow.Flow) error {
	msg := h.Parcel.Message().(*message.ExecutorResults)

	// now we have 2 different types of data in message.HandleExecutorResultsMessage
	// one part of it is about consensus
	// another one is about prepare state on new executor after pulse
	// TODO make it in different goroutines

	// prepare state after previous executor
	procInitializeExecutionState := initializeExecutionState{
		LR:  h.dep.lr,
		msg: msg,
	}
	if err := f.Procedure(ctx, &procInitializeExecutionState, true); err != nil {
		if err == flow.ErrCancelled {
			return nil
		}
		err := errors.Wrap(err, "[ HandleExecutorResults ] Failed to initialize execution state")
		return err
	}

	if procInitializeExecutionState.Result.clarifyPending {
		procClarifyPending := ClarifyPendingState{
			es:              procInitializeExecutionState.Result.es,
			ArtifactManager: h.dep.lr.ArtifactManager,
		}

		if err := f.Procedure(ctx, &procClarifyPending, true); err != nil {
			if err == flow.ErrCancelled {
				return nil
			}

			err := errors.Wrap(err, "[ HandleExecutorResults ] Failed to clarify pending")
			return err
		}
	}

	es := procInitializeExecutionState.Result.es
	es.Broker.StartProcessorIfNeeded(ctx)
	return nil
}

func (h *HandleExecutorResults) Present(ctx context.Context, f flow.Flow) error {
	ctx = loggerWithTargetID(ctx, h.Parcel)
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

	var rep *watermillMsg.Message
	if err != nil {
		var newErr error
		rep, newErr = payload.NewMessage(&payload.Error{Text: err.Error()})
		if newErr != nil {
			return newErr
		}
	} else {
		rep = bus.ReplyAsMessage(ctx, &reply.OK{})
	}
	h.dep.Sender.Reply(ctx, h.Message, rep)

	return err
}
