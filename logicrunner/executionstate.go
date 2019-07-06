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
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
)

type ExecutionState struct {
	sync.Mutex

	Ref Ref // Object reference

	ObjectDescriptor artifacts.ObjectDescriptor

	Broker                *ExecutionBroker
	LedgerHasMoreRequests bool
	getLedgerPendingMutex sync.Mutex

	// TODO not using in validation, need separate ObjectState.ExecutionState and ObjectState.Validation from ExecutionState struct
	pending              message.PendingState
	PendingConfirmed     bool
	HasPendingCheckMutex sync.Mutex
}

func NewExecutionState(ref insolar.Reference) *ExecutionState {
	es := &ExecutionState{
		Ref:     ref,
		pending: message.PendingUnknown,
	}
	es.Broker = NewExecutionBroker(es)

	return es
}

func (es *ExecutionState) RegisterLogicRunner(lr *LogicRunner) {
	es.Broker.logicRunner = lr
}

func (es *ExecutionState) OnPulse(ctx context.Context, meNext bool) []insolar.Message {
	if es == nil {
		return nil
	}

	messages := make([]insolar.Message, 0)
	ref := es.Ref

	// if we are executor again we just continue working
	// without sending data on next executor (because we are next executor)
	if !meNext {
		sendExecResults := false

		if !es.Broker.currentList.Empty() {
			es.pending = message.InPending
			sendExecResults = true

			// TODO: this should return delegation token to continue execution of the pending
			messages = append(
				messages,
				&message.StillExecuting{
					Reference: ref,
				},
			)
		} else if es.pending == message.InPending && !es.PendingConfirmed {
			inslogger.FromContext(ctx).Warn(
				"looks like pending executor died, continuing execution on next executor",
			)
			es.pending = message.NotPending
			sendExecResults = true
			es.LedgerHasMoreRequests = true
		} else if es.Broker.finished.Len() > 0 {
			sendExecResults = true
		}

		// rotation results also contain finished requests
		rotationResults := es.Broker.Rotate(maxQueueLength)
		if len(rotationResults.Requests) > 0 || sendExecResults {
			// TODO: we also should send when executed something for validation
			// TODO: now validation is disabled
			messagesQueue := convertQueueToMessageQueue(ctx, rotationResults.Requests)

			messages = append(
				messages,
				//&message.ValidateCaseBind{
				//	Reference: ref,
				//	Requests:  requests,
				//	Pulse:     pulse,
				//},
				&message.ExecutorResults{
					RecordRef:             ref,
					Pending:               es.pending,
					Queue:                 messagesQueue,
					LedgerHasMoreRequests: es.LedgerHasMoreRequests || rotationResults.LedgerHasMoreRequests,
				},
			)
		}
	} else {
		if !es.Broker.currentList.Empty() {
			// no pending should be as we are executing
			if es.pending == message.InPending {
				inslogger.FromContext(ctx).Warn(
					"we are executing ATM, but ES marked as pending, shouldn't be",
				)
				es.pending = message.NotPending
			}
		} else if es.pending == message.InPending && !es.PendingConfirmed {
			inslogger.FromContext(ctx).Warn(
				"looks like pending executor died, re-starting execution",
			)
			es.pending = message.NotPending
			es.LedgerHasMoreRequests = true
		}
		es.PendingConfirmed = false
	}

	return messages
}
