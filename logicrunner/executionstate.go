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
	"sync"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/message"
)

type ExecutionState struct {
	sync.Mutex

	Ref Ref

	objectbody *ObjectBody
	deactivate bool
	nonce      uint64

	Current               *CurrentExecution
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

func (es *ExecutionState) WrapError(err error, message string) error {
	if err == nil {
		err = errors.New(message)
	} else {
		err = errors.Wrap(err, message)
	}
	res := Error{Err: err}
	if es.objectbody != nil {
		res.Contract = es.objectbody.objDescriptor.HeadRef()
	}
	if es.Current != nil {
		res.Request = es.Current.Request
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
