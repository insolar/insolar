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

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/logicrunner/artifacts"
)

// Context of one contract execution
type ObjectState struct {
	sync.Mutex

	ExecutionBroker *ExecutionBroker
	Validation      *ExecutionState
}

//go:generate minimock -i github.com/insolar/insolar/logicrunner.StateStorage -o ./ -s _mock.go
type StateStorage interface {
	sync.Locker

	UpsertExecutionState(ref insolar.Reference) *ExecutionBroker
	GetExecutionState(ref insolar.Reference) *ExecutionBroker

	UpsertValidationState(ref insolar.Reference) *ExecutionState
	GetValidationState(ref insolar.Reference) *ExecutionState

	DeleteObjectState(ref insolar.Reference)

	StateMap() *map[insolar.Reference]*ObjectState
}

type stateStorage struct {
	sync.RWMutex

	publisher        watermillMsg.Publisher
	requestsExecutor RequestsExecutor
	messageBus       insolar.MessageBus
	jetCoordinator   jet.Coordinator
	pulseAccessor    pulse.Accessor
	artifactsManager artifacts.Client

	state map[insolar.Reference]*ObjectState // if object exists, we are validating or executing it right now
}

func (ss *stateStorage) UpsertValidationState(ref insolar.Reference) *ExecutionState {
	os := ss.upsertObjectState(ref)

	os.Lock()
	defer os.Unlock()

	os.Validation = &ExecutionState{}
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

func NewStateStorage(
	publisher watermillMsg.Publisher,
	requestsExecutor RequestsExecutor,
	messageBus insolar.MessageBus,
	jetCoordinator jet.Coordinator,
	pulseAccessor pulse.Accessor,
	artifactsManager artifacts.Client,

) StateStorage {
	ss := &stateStorage{
		state: make(map[insolar.Reference]*ObjectState),

		publisher:        publisher,
		requestsExecutor: requestsExecutor,
		messageBus:       messageBus,
		jetCoordinator:   jetCoordinator,
		pulseAccessor:    pulseAccessor,
		artifactsManager: artifactsManager,
	}
	return ss
}

func (ss *stateStorage) UpsertExecutionState(ref insolar.Reference) *ExecutionBroker {
	os := ss.upsertObjectState(ref)

	os.Lock()
	defer os.Unlock()

	if os.ExecutionBroker == nil {
		os.ExecutionBroker = NewExecutionBroker(
			ref,
			ss.publisher,
			ss.requestsExecutor,
			ss.messageBus,
			ss.jetCoordinator,
			ss.pulseAccessor,
			ss.artifactsManager,
		)
	}
	return os.ExecutionBroker
}

func (ss *stateStorage) GetExecutionState(ref insolar.Reference) *ExecutionBroker {
	os := ss.getObjectState(ref)
	if os == nil {
		return nil
	}

	os.Lock()
	defer os.Unlock()

	return os.ExecutionBroker
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
