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
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

type ProcessExecutionQueue struct {
	dep *Dependencies

	Message *watermillMsg.Message
}

func (p *ProcessExecutionQueue) Present(ctx context.Context, f flow.Flow) error {
	lr := p.dep.lr
	es := lr.getExecStateFromRef(ctx, p.Message.Payload)
	if es == nil {
		return nil
	}

	for {
		es.Lock()
		if len(es.Queue) == 0 && es.LedgerQueueElement == nil {
			inslogger.FromContext(ctx).Debug("Quiting queue processing, empty")
			es.QueueProcessorActive = false
			es.Current = nil
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

		sender := qe.parcel.GetSender()
		current := CurrentExecution{
			Request:       qe.request,
			RequesterNode: &sender,
			Context:       qe.ctx,
		}
		es.Current = &current

		if msg, ok := qe.parcel.Message().(*message.CallMethod); ok {
			current.ReturnMode = msg.ReturnMode
		}
		if msg, ok := qe.parcel.Message().(message.IBaseLogicMessage); ok {
			current.Sequence = msg.GetBaseLogicMessage().Sequence
		}

		es.Unlock()

		lr.executeOrValidate(current.Context, es, qe.parcel)

		if qe.fromLedger {
			go lr.getLedgerPendingRequest(ctx, es)
		}

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

	inslogger.FromContext(ctx).Debug("Starting a new queue processor")
	s.es.QueueProcessorActive = true

	pub := s.dep.Publisher
	rawRef := s.ref.Bytes()
	err := pub.Publish(InnerMsgTopic, makeWMMessage(ctx, rawRef, "ProcessExecutionQueue"))
	if err != nil {
		return errors.Wrap(err, "can't send ProcessExecutionQueue msg")
	}
	err = pub.Publish(InnerMsgTopic, makeWMMessage(ctx, rawRef, "getLedgerPendingRequest"))
	if err != nil {
		return errors.Wrap(err, "can't send getLedgerPendingRequest msg")
	}

	return nil
}
