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
	"time"

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

type ProcessExecutionQueue struct {
	dep *Dependencies

	Message *watermillMsg.Message
}

func (p *ProcessExecutionQueue) Present(ctx context.Context, f flow.Flow) error {
	ctx, span := instracer.StartSpan(ctx, "ProcessExecutionQueue")
	defer span.End()

	inslogger.FromContext(ctx).Debug("ProcessExecutionQueue")

	lr := p.dep.lr
	es := lr.GetExecutionState(Ref{}.FromSlice(p.Message.Payload))
	if es == nil {
		return nil
	}

	for {
		es.Lock()
		if len(es.Queue) == 0 && es.LedgerQueueElement == nil {
			inslogger.FromContext(ctx).Debug("Quiting queue processing, empty. Ref: ", es.Ref.String())
			es.QueueProcessorActive = false

			if mutable := es.CurrentList.GetMutable(); mutable != nil {
				es.CurrentList.Delete(*mutable.LogicContext.Request)
			}
			es.Unlock()
			return nil
		}

		var qe ExecutionQueueElement
		if es.LedgerQueueElement != nil {
			qe = *es.LedgerQueueElement
			es.LedgerQueueElement = nil
		} else {
			qe, es.Queue = es.Queue[0], es.Queue[1:]
		}

		msg, ok := qe.parcel.Message().(*message.CallMethod)
		if !ok {
			panic("Not a call method message, should never happen")
		}

		sender := qe.parcel.GetSender()
		current := CurrentExecution{
			RequestRef:    qe.request,
			RequesterNode: &sender,
			Context:       qe.ctx,
			Request:       &msg.Request,
			LogicContext: &insolar.LogicCallContext{
				Mode:            "execution",
				Caller:          msg.GetCaller(),
				Callee:          &es.Ref,
				Request:         qe.request,
				Time:            time.Now(), // TODO: probably we should take it earlier
				Pulse:           *lr.pulse(ctx),
				TraceID:         inslogger.TraceID(ctx),
				CallerPrototype: &msg.CallerPrototype,
				Immutable:       msg.Immutable,
			},
		}

		es.CurrentList.Set(*qe.request, &current)
		es.Unlock()

		lr.executeOrValidate(current.Context, es, &current)

		if qe.fromLedger {
			pub := p.dep.Publisher
			err := pub.Publish(InnerMsgTopic, makeWMMessage(ctx, p.Message.Payload, getLedgerPendingRequestMsg))
			if err != nil {
				inslogger.FromContext(ctx).Warnf("can't send processExecutionQueueMsg: ", err)
			}
		}
		es.Finished = append(es.Finished, &current)

		lr.finishPendingIfNeeded(ctx, es)
	}
}

// ---------------- StartQueueProcessorIfNeeded

type StartQueueProcessorIfNeeded struct {
	es  *ExecutionState
	ref *insolar.Reference
	dep *Dependencies
}

func (s *StartQueueProcessorIfNeeded) Present(ctx context.Context, f flow.Flow) error {
	ctx, span := instracer.StartSpan(ctx, "StartQueueProcessorIfNeeded")
	defer span.End()

	s.es.Lock()
	defer s.es.Unlock()

	if !s.es.haveSomeToProcess() {
		inslogger.FromContext(ctx).Debug("queue is empty. processor is not needed")
		return nil
	}

	if s.es.QueueProcessorActive {
		inslogger.FromContext(ctx).Debug("queue processor is already active. processor is not needed")
		return nil
	}

	if s.es.pending == message.PendingUnknown {
		return errors.New("shouldn't start queue processor with unknown pending state")
	} else if s.es.pending == message.InPending {
		inslogger.FromContext(ctx).Debug("object in pending. not starting queue processor")
		return nil
	}

	pub := s.dep.Publisher
	rawRef := s.ref.Bytes()
	err := pub.Publish(InnerMsgTopic, makeWMMessage(ctx, rawRef, processExecutionQueueMsg))
	if err != nil {
		return errors.Wrap(err, "can't send processExecutionQueueMsg")
	}
	err = pub.Publish(InnerMsgTopic, makeWMMessage(ctx, rawRef, getLedgerPendingRequestMsg))
	if err != nil {
		return errors.Wrap(err, "can't send getLedgerPendingRequestMsg")
	}

	inslogger.FromContext(ctx).Debug("Starting a new queue processor")
	s.es.QueueProcessorActive = true

	return nil
}
