// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package executionregistry

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/logicrunner/common"
)

var (
	ErrAlreadyRegistered = errors.New("trying to register task that is executing right now")
)

type Registry interface {
	Register(transcript *common.Transcript)
}

//go:generate minimock -i github.com/insolar/insolar/logicrunner/executionregistry.ExecutionRegistry -o ./ -s _mock.go -g
type ExecutionRegistry interface {
	Register(ctx context.Context, transcript *common.Transcript) error
	Done(transcript *common.Transcript) bool

	Length() int
	IsEmpty() bool
	OnPulse(ctx context.Context) []payload.Payload
	GetActiveTranscript(req insolar.Reference) *common.Transcript
}

type executionRegistry struct {
	// maps requestReference -> request to notify
	// (StillExecuting message) that we've started to work on in previous pulses
	objectRef    insolar.Reference
	registryLock sync.Mutex
	registry     map[insolar.Reference]*common.Transcript

	jetCoordinator jet.Coordinator
}

func New(
	objectRef insolar.Reference,
	jetCoordinator jet.Coordinator,
) ExecutionRegistry {

	return &executionRegistry{
		objectRef:    objectRef,
		registryLock: sync.Mutex{},
		registry:     make(map[insolar.Reference]*common.Transcript),

		jetCoordinator: jetCoordinator,
	}
}

func (r *executionRegistry) Length() int {
	r.registryLock.Lock()
	defer r.registryLock.Unlock()

	return len(r.registry)
}

func (r *executionRegistry) IsEmpty() bool {
	return r.Length() == 0
}

func (r *executionRegistry) Register(ctx context.Context, transcript *common.Transcript) error {
	requestRef := transcript.RequestRef

	r.registryLock.Lock()
	defer r.registryLock.Unlock()

	if _, ok := r.registry[requestRef]; ok {
		return ErrAlreadyRegistered
	}

	r.registry[requestRef] = transcript
	return nil
}

func (r *executionRegistry) Done(transcript *common.Transcript) bool {
	requestRef := transcript.RequestRef

	r.registryLock.Lock()
	defer r.registryLock.Unlock()

	if _, ok := r.registry[requestRef]; !ok {
		return false
	}

	delete(r.registry, requestRef)
	return true
}

// constructs all StillExecuting messages
func (r *executionRegistry) OnPulse(_ context.Context) []payload.Payload {
	r.registryLock.Lock()
	defer r.registryLock.Unlock()

	// TODO: this should return delegation token to continue execution of the pending
	messages := make([]payload.Payload, 0)
	if len(r.registry) != 0 {
		requestRefs := make([]insolar.Reference, 0, len(r.registry))
		for requestRef := range r.registry {
			requestRefs = append(requestRefs, requestRef)
		}

		messages = append(messages, &payload.StillExecuting{
			ObjectRef:   r.objectRef,
			Executor:    r.jetCoordinator.Me(),
			RequestRefs: requestRefs,
		})
	}
	return messages
}

func (r *executionRegistry) GetActiveTranscript(request insolar.Reference) *common.Transcript {
	r.registryLock.Lock()
	defer r.registryLock.Unlock()

	return r.registry[request]
}
