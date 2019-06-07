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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
)

type ExecutionState struct {
	sync.Mutex

	Ref Ref // Object reference

	ObjectDescriptor    artifacts.ObjectDescriptor
	PrototypeDescriptor artifacts.ObjectDescriptor
	CodeDescriptor      artifacts.CodeDescriptor

	Finished              []*CurrentExecution
	CurrentList           *CurrentExecutionList
	Queue                 []ExecutionQueueElement
	QueueProcessorActive  bool
	LedgerHasMoreRequests bool
	LedgerQueueElement    *ExecutionQueueElement
	getLedgerPendingMutex sync.Mutex

	// TODO not using in validation, need separate ObjectState.ExecutionState and ObjectState.Validation from ExecutionState struct
	pending              message.PendingState
	PendingConfirmed     bool
	HasPendingCheckMutex sync.Mutex
}

func NewExecutionState(ref insolar.Reference) *ExecutionState {
	return &ExecutionState{
		Ref:         ref,
		CurrentList: NewCurrentExecutionList(),
		Queue:       make([]ExecutionQueueElement, 0),
	}
}

func (es *ExecutionState) WrapError(current *CurrentExecution, err error, message string) error {
	if err == nil {
		err = errors.New(message)
	} else {
		err = errors.Wrap(err, message)
	}
	res := Error{Err: err}
	if es.ObjectDescriptor != nil {
		res.Contract = es.ObjectDescriptor.HeadRef()
	}
	if current != nil {
		res.Request = current.RequestRef
	}
	return res
}

// releaseQueue must be calling only with es.Lock
func (es *ExecutionState) releaseQueue() ([]ExecutionQueueElement, bool) {
	ledgerHasMoreRequest := false
	q := es.Queue

	if len(q) > maxQueueLength {
		q = q[:maxQueueLength]
		ledgerHasMoreRequest = true
	}

	es.Queue = make([]ExecutionQueueElement, 0)

	return q, ledgerHasMoreRequest
}

func (es *ExecutionState) haveSomeToProcess() bool {
	return len(es.Queue) > 0 || es.LedgerHasMoreRequests || es.LedgerQueueElement != nil
}

func (es *ExecutionState) OnPulse(ctx context.Context, meNext bool) []insolar.Message {
	messages := make([]insolar.Message, 0)

	ref := es.Ref

	// if we are executor again we just continue working
	// without sending data on next executor (because we are next executor)
	if !meNext {
		sendExecResults := false

		if !es.CurrentList.Empty() {
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
				"looks like pending executor died, continuing execution",
			)
			es.pending = message.NotPending
			sendExecResults = true
			es.LedgerHasMoreRequests = true
		}

		queue, ledgerHasMoreRequest := es.releaseQueue()
		if len(queue) > 0 || sendExecResults {
			// TODO: we also should send when executed something for validation
			// TODO: now validation is disabled
			messagesQueue := convertQueueToMessageQueue(ctx, queue)

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
					LedgerHasMoreRequests: es.LedgerHasMoreRequests || ledgerHasMoreRequest,
				},
			)
		}
	} else {
		if !es.CurrentList.Empty() {
			// no pending should be as we are executing
			if es.pending == message.InPending {
				inslogger.FromContext(ctx).Warn(
					"we are executing ATM, but ES marked as pending, shouldn't be",
				)
				es.pending = message.NotPending
			}
		} else if es.pending == message.InPending && !es.PendingConfirmed {
			inslogger.FromContext(ctx).Warn(
				"looks like pending executor died, continuing execution",
			)
			es.pending = message.NotPending
			es.LedgerHasMoreRequests = true
		}
		es.PendingConfirmed = false
	}
	es.Finished = nil

	return messages
}
