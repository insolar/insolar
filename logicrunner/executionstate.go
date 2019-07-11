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

	"github.com/insolar/insolar/insolar/message"
)

type ExecutionState struct {
	sync.Mutex

	LedgerHasMoreRequests bool
	getLedgerPendingMutex sync.Mutex

	// TODO not using in validation, need separate ObjectState.ExecutionState and ObjectState.Validation from ExecutionState struct
	pending              message.PendingState
	PendingConfirmed     bool
	HasPendingCheckMutex sync.Mutex
}

func NewExecutionState() *ExecutionState {
	es := &ExecutionState{
		pending: message.PendingUnknown,
	}

	return es
}

// PendingNotConfirmed checks that we were in pending and waiting
// for previous executor to notify us that he still executes it
//
// Used in OnPulse to tell next executor, that it's time to continue
// work on this object
func (es *ExecutionState) InPendingNotConfirmed() bool {
	return es.pending == message.InPending && !es.PendingConfirmed
}
