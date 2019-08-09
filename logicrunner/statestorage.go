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
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
)

//go:generate minimock -i github.com/insolar/insolar/logicrunner.StateStorage -o ./ -s _mock.go -g
type StateStorage interface {
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
	outgoingSender   OutgoingRequestSender

	brokers  map[insolar.Reference]ExecutionBrokerI
	archives map[insolar.Reference]ExecutionArchive
}

func NewStateStorage(
	publisher watermillMsg.Publisher,
	requestsExecutor RequestsExecutor,
	messageBus insolar.MessageBus,
	jetCoordinator jet.Coordinator,
	pulseAccessor pulse.Accessor,
	artifactsManager artifacts.Client,
	outgoingSender OutgoingRequestSender,
) StateStorage {
	ss := &stateStorage{
		brokers:  make(map[insolar.Reference]ExecutionBrokerI),
		archives: make(map[insolar.Reference]ExecutionArchive),

		publisher:        publisher,
		requestsExecutor: requestsExecutor,
		messageBus:       messageBus,
		jetCoordinator:   jetCoordinator,
		pulseAccessor:    pulseAccessor,
		artifactsManager: artifactsManager,
		outgoingSender:   outgoingSender,
	}
	return ss
}

func (ss *stateStorage) upsertExecutionArchive(ref insolar.Reference) ExecutionArchive {
	if res, ok := ss.archives[ref]; ok {
		return res
	}

	ss.archives[ref] = NewExecutionArchive(ref, ss.jetCoordinator)
	return ss.archives[ref]
}

func (ss *stateStorage) UpsertExecutionState(ref insolar.Reference) ExecutionBrokerI {
	ss.RLock()
	if res, ok := ss.brokers[ref]; ok {
		ss.RUnlock()
		return res
	}
	ss.RUnlock()

	ss.Lock()
	defer ss.Unlock()
	if _, ok := ss.brokers[ref]; !ok {
		archive := ss.upsertExecutionArchive(ref)

		ss.brokers[ref] = NewExecutionBroker(
			ref,
			ss.publisher,
			ss.requestsExecutor,
			ss.messageBus,
			ss.jetCoordinator,
			ss.pulseAccessor,
			ss.artifactsManager,
			archive,
			ss.outgoingSender,
		)
	}
	return ss.brokers[ref]
}

func (ss *stateStorage) GetExecutionState(ref insolar.Reference) ExecutionBrokerI {
	ss.RLock()
	defer ss.RUnlock()
	return ss.brokers[ref]
}

func (ss *stateStorage) GetExecutionArchive(ref insolar.Reference) ExecutionArchive {
	ss.RLock()
	defer ss.RUnlock()
	return ss.archives[ref]
}

func (ss *stateStorage) IsEmpty() bool {
	ss.RLock()
	defer ss.RUnlock()

	if len(ss.brokers) == 0 && len(ss.archives) == 0 {
		return true
	}

	for _, el := range ss.archives {
		if !el.IsEmpty() {
			return false
		}
	}
	return true
}

func (ss *stateStorage) OnPulse(ctx context.Context, pulse insolar.Pulse) []insolar.Message {
	onPulseMessages := make([]insolar.Message, 0)

	ss.Lock()
	defer ss.Unlock()

	oldBrokers := ss.brokers
	ss.brokers = make(map[insolar.Reference]ExecutionBrokerI)
	for objectRef, broker := range oldBrokers {
		if _, ok := ss.archives[objectRef]; !ok {
			inslogger.FromContext(ctx).Error("exeuction broker exists, but archive doesn't")
		}

		meNext, _ := ss.jetCoordinator.IsMeAuthorizedNow(ctx, insolar.DynamicRoleVirtualExecutor, *objectRef.Record())

		onPulseMessages = append(onPulseMessages, broker.OnPulse(ctx, meNext)...)

		if meNext {
			ss.brokers[objectRef] = broker
		}
	}

	for objectRef, archive := range ss.archives {
		onPulseMessages = append(onPulseMessages, archive.OnPulse(ctx)...)

		meNext, _ := ss.jetCoordinator.IsMeAuthorizedNow(ctx, insolar.DynamicRoleVirtualExecutor, *objectRef.Record())
		if !meNext && archive.IsEmpty() {
			delete(ss.archives, objectRef)
		}
	}

	return onPulseMessages
}
