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

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/logicrunner/artifacts"
)

// Context of one contract execution
type ObjectState struct {
	sync.Mutex

	ExecutionBroker  ExecutionBrokerI
	ExecutionArchive ExecutionArchive
}

//go:generate minimock -i github.com/insolar/insolar/logicrunner.StateStorage -o ./ -s _mock.go -g
type StateStorage interface {
	sync.Locker

	UpsertExecutionState(ref insolar.Reference) ExecutionBrokerI
	GetExecutionState(ref insolar.Reference) ExecutionBrokerI
	GetExecutionArchive(ref insolar.Reference) ExecutionArchive

	IsEmpty() bool
	OnPulse(ctx context.Context, pulse insolar.Pulse) []insolar.Message
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

func (ss *stateStorage) UpsertExecutionState(ref insolar.Reference) ExecutionBrokerI {
	os := ss.upsertObjectState(ref)

	os.Lock()
	defer os.Unlock()

	if os.ExecutionArchive == nil {
		os.ExecutionArchive = NewExecutionArchive(ref, ss.jetCoordinator)
	}
	if os.ExecutionBroker == nil {
		os.ExecutionBroker = NewExecutionBroker(
			ref,
			ss.publisher,
			ss.requestsExecutor,
			ss.messageBus,
			ss.jetCoordinator,
			ss.pulseAccessor,
			ss.artifactsManager,
			os.ExecutionArchive,
		)
	}
	return os.ExecutionBroker
}

func (ss *stateStorage) GetExecutionState(ref insolar.Reference) ExecutionBrokerI {
	os := ss.getObjectState(ref)
	if os == nil {
		return nil
	}

	os.Lock()
	defer os.Unlock()

	return os.ExecutionBroker
}

func (ss *stateStorage) GetExecutionArchive(ref insolar.Reference) ExecutionArchive {
	os := ss.getObjectState(ref)
	if os == nil {
		return nil
	}

	os.Lock()
	defer os.Unlock()

	return os.ExecutionArchive
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

func (ss *stateStorage) IsEmpty() bool {
	return len(ss.state) == 0
}

func (ss *stateStorage) OnPulse(ctx context.Context, pulse insolar.Pulse) []insolar.Message {
	ss.Lock()
	defer ss.Unlock()

	onPulseMessages := make([]insolar.Message, 0)

	oldState := ss.state
	ss.state = make(map[insolar.Reference]*ObjectState)
	for objectRef, objectState := range oldState {
		objectState.Lock()

		meNext, _ := ss.jetCoordinator.IsAuthorized(
			ctx, insolar.DynamicRoleVirtualExecutor, *objectRef.Record(), pulse.PulseNumber, ss.jetCoordinator.Me(),
		)

		if broker := objectState.ExecutionBroker; broker != nil {
			onPulseMessages = append(onPulseMessages, broker.OnPulse(ctx, meNext)...)
		}

		if archive := objectState.ExecutionArchive; archive != nil {
			onPulseMessages = append(onPulseMessages, archive.OnPulse(ctx)...)
		}

		if meNext && objectState.ExecutionArchive != nil {
			ss.state[objectRef] = &ObjectState{
				ExecutionArchive: objectState.ExecutionArchive,
				ExecutionBroker:  objectState.ExecutionBroker,
			}
		} else if objectState.ExecutionArchive != nil && !objectState.ExecutionArchive.IsEmpty() {
			// import previous ExecutionArchive if it's not empty
			ss.state[objectRef] = &ObjectState{
				ExecutionArchive: objectState.ExecutionArchive,
			}
		}

		objectState.Unlock()
	}

	return onPulseMessages
}
