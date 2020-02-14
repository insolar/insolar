// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package logicrunner

import (
	"context"
	"sync"

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/executionregistry"
	"github.com/insolar/insolar/logicrunner/shutdown"
)

//go:generate minimock -i github.com/insolar/insolar/logicrunner.StateStorage -o ./ -s _mock.go -g
type StateStorage interface {
	UpsertExecutionState(ref insolar.Reference) ExecutionBrokerI
	GetExecutionState(ref insolar.Reference) ExecutionBrokerI
	GetExecutionRegistry(ref insolar.Reference) executionregistry.ExecutionRegistry

	IsEmpty() bool
	OnPulse(ctx context.Context, pulse insolar.Pulse) map[insolar.Reference][]payload.Payload
}

type stateStorage struct {
	sync.RWMutex

	publisher        watermillMsg.Publisher
	requestsExecutor RequestsExecutor
	sender           bus.Sender
	jetCoordinator   jet.Coordinator
	pulseAccessor    pulse.Accessor
	artifactsManager artifacts.Client
	outgoingSender   OutgoingRequestSender
	shutdownFlag     shutdown.Flag

	brokers    map[insolar.Reference]ExecutionBrokerI
	registries map[insolar.Reference]executionregistry.ExecutionRegistry
}

func NewStateStorage(
	publisher watermillMsg.Publisher,
	requestsExecutor RequestsExecutor,
	sender bus.Sender,
	jetCoordinator jet.Coordinator,
	pulseAccessor pulse.Accessor,
	artifactsManager artifacts.Client,
	outgoingSender OutgoingRequestSender,
	shutdownFlag shutdown.Flag,
) StateStorage {
	ss := &stateStorage{
		brokers:    make(map[insolar.Reference]ExecutionBrokerI),
		registries: make(map[insolar.Reference]executionregistry.ExecutionRegistry),

		publisher:        publisher,
		requestsExecutor: requestsExecutor,
		sender:           sender,
		jetCoordinator:   jetCoordinator,
		pulseAccessor:    pulseAccessor,
		artifactsManager: artifactsManager,
		outgoingSender:   outgoingSender,
		shutdownFlag:     shutdownFlag,
	}
	return ss
}

func (ss *stateStorage) upsertExecutionRegistry(ref insolar.Reference) executionregistry.ExecutionRegistry {
	if res, ok := ss.registries[ref]; ok {
		return res
	}

	ss.registries[ref] = executionregistry.New(ref, ss.jetCoordinator)
	return ss.registries[ref]
}

func (ss *stateStorage) UpsertExecutionState(ref insolar.Reference) ExecutionBrokerI {
	if ss.shutdownFlag.IsStopped() {
		log.Warn("UpsertExecutionState after shutdown was triggered for ", ref.String())
	}

	ss.RLock()
	if res, ok := ss.brokers[ref]; ok {
		ss.RUnlock()
		return res
	}
	ss.RUnlock()

	ss.Lock()
	defer ss.Unlock()
	if _, ok := ss.brokers[ref]; !ok {
		registry := ss.upsertExecutionRegistry(ref)

		ss.brokers[ref] = NewExecutionBroker(ref, ss.publisher, ss.requestsExecutor, ss.sender, ss.artifactsManager, registry, ss.outgoingSender, ss.pulseAccessor)
	}
	return ss.brokers[ref]
}

func (ss *stateStorage) GetExecutionState(ref insolar.Reference) ExecutionBrokerI {
	if ss.shutdownFlag.IsStopped() {
		log.Warn("GetExecutionState after shutdown was triggered for ", ref.String())
	}

	ss.RLock()
	defer ss.RUnlock()
	return ss.brokers[ref]
}

func (ss *stateStorage) GetExecutionRegistry(ref insolar.Reference) executionregistry.ExecutionRegistry {
	if ss.shutdownFlag.IsStopped() {
		log.Warn("GetExecutionRegistry after shutdown was triggered for ", ref.String())
	}

	ss.RLock()
	defer ss.RUnlock()
	return ss.registries[ref]
}

func (ss *stateStorage) IsEmpty() bool {
	ss.RLock()
	defer ss.RUnlock()

	if len(ss.brokers) == 0 && len(ss.registries) == 0 {
		return true
	}

	for _, el := range ss.registries {
		if !el.IsEmpty() {
			return false
		}
	}
	return true
}

func (ss *stateStorage) OnPulse(ctx context.Context, pulse insolar.Pulse) map[insolar.Reference][]payload.Payload {
	onPulseMessages := make(map[insolar.Reference][]payload.Payload)

	ss.Lock()
	defer ss.Unlock()

	oldBrokers := ss.brokers
	ss.brokers = make(map[insolar.Reference]ExecutionBrokerI)
	for objectRef, broker := range oldBrokers {
		if _, ok := ss.registries[objectRef]; !ok {
			inslogger.FromContext(ctx).Error("exeuction broker exists, but registry doesn't")
		}

		messages := broker.OnPulse(ctx)
		if len(messages) > 0 {
			onPulseMessages[objectRef] = append(onPulseMessages[objectRef], messages...)
		}
	}

	for objectRef, registry := range ss.registries {
		messages := registry.OnPulse(ctx)
		if len(messages) > 0 {
			onPulseMessages[objectRef] = append(onPulseMessages[objectRef], messages...)
		}

		if registry.IsEmpty() {
			delete(ss.registries, objectRef)
		}
	}

	return onPulseMessages
}
