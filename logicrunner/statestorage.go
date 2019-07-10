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

	"github.com/insolar/insolar/insolar"
)

//go:generate minimock -i github.com/insolar/insolar/logicrunner.StateStorage -o ./ -s _mock.go
type StateStorage interface {
	sync.Locker

	UpsertExecutionState(lr *LogicRunner, ref insolar.Reference) (*ExecutionState, *ExecutionBroker)
	GetExecutionState(ref insolar.Reference) (*ExecutionState, *ExecutionBroker)

	UpsertValidationState(ref insolar.Reference) *ExecutionState
	GetValidationState(ref insolar.Reference) *ExecutionState

	DeleteObjectState(ref insolar.Reference)

	StateMap() *map[insolar.Reference]*ObjectState
}

type stateStorage struct {
	sync.RWMutex
	state map[insolar.Reference]*ObjectState // if object exists, we are validating or executing it right now
}

func (ss *stateStorage) UpsertValidationState(ref insolar.Reference) *ExecutionState {
	os := ss.upsertObjectState(ref)

	os.Lock()
	defer os.Unlock()

	os.Validation = NewExecutionState(ref)
	return os.Validation
}

func (ss *stateStorage) GetValidationState(ref insolar.Reference) *ExecutionState {
	os := ss.getObjectState(ref)
	if os == nil {
		return nil
	}

	os.Lock()
	defer os.Unlock()

	return os.Validation
}

func NewStateStorage() StateStorage {
	ss := &stateStorage{
		state: make(map[insolar.Reference]*ObjectState),
	}
	return ss
}

func (ss *stateStorage) UpsertExecutionState(
	lr *LogicRunner,
	ref insolar.Reference,
) (*ExecutionState, *ExecutionBroker) {
	os := ss.upsertObjectState(ref)

	os.Lock()
	defer os.Unlock()

	if os.ExecutionState == nil {
		os.ExecutionState = NewExecutionState(ref)
		os.ExecutionBroker = NewExecutionBroker(lr, os.ExecutionState)
	}
	return os.ExecutionState, os.ExecutionBroker
}

func (ss *stateStorage) GetExecutionState(ref insolar.Reference) (*ExecutionState, *ExecutionBroker) {
	os := ss.getObjectState(ref)
	if os == nil {
		return nil, nil
	}

	os.Lock()
	defer os.Unlock()

	return os.ExecutionState, os.ExecutionBroker
}

func (ss *stateStorage) getObjectState(ref insolar.Reference) *ObjectState {
	ss.RLock()
	res, ok := ss.state[ref]
	ss.RUnlock()
	if !ok {
		return nil
	}
	return res
}

func (ss *stateStorage) upsertObjectState(ref insolar.Reference) *ObjectState {
	ss.RLock()
	if res, ok := ss.state[ref]; ok {
		ss.RUnlock()
		return res
	}
	ss.RUnlock()

	ss.Lock()
	defer ss.Unlock()
	if _, ok := ss.state[ref]; !ok {
		ss.state[ref] = &ObjectState{}
	}
	return ss.state[ref]
}

func (ss *stateStorage) DeleteObjectState(ref insolar.Reference) {
	delete(ss.state, ref)
}

func (ss *stateStorage) StateMap() *map[insolar.Reference]*ObjectState {
	return &ss.state
}
