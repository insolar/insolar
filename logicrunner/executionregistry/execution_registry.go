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

package executionregistry

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/common"
)

type Registry interface {
	Register(transcript *common.Transcript)
}

//go:generate minimock -i github.com/insolar/insolar/logicrunner/executionregistry.ExecutionRegistry -o ./ -s _mock.go -g
type ExecutionRegistry interface {
	Register(ctx context.Context, transcript *common.Transcript)
	Done(transcript *common.Transcript) bool

	Length() int
	IsEmpty() bool
	OnPulse(ctx context.Context) []insolar.Message
	FindRequestLoop(ctx context.Context, reqRef insolar.Reference, apiRequestID string) bool
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

func (r *executionRegistry) Register(ctx context.Context, transcript *common.Transcript) {
	requestRef := transcript.RequestRef

	r.registryLock.Lock()
	defer r.registryLock.Unlock()

	if _, ok := r.registry[requestRef]; !ok {
		r.registry[requestRef] = transcript
	} else {
		inslogger.FromContext(ctx).Error("Trying to register task that is already registered")
	}
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
func (r *executionRegistry) OnPulse(_ context.Context) []insolar.Message {
	r.registryLock.Lock()
	defer r.registryLock.Unlock()

	// TODO: this should return delegation token to continue execution of the pending
	messages := make([]insolar.Message, 0)
	if len(r.registry) != 0 {
		requestRefs := make([]insolar.Reference, 0, len(r.registry))
		for requestRef := range r.registry {
			requestRefs = append(requestRefs, requestRef)
		}

		messages = append(messages, &message.StillExecuting{
			Reference:   r.objectRef,
			Executor:    r.jetCoordinator.Me(),
			RequestRefs: requestRefs,
		})
	}
	return messages
}

func (r *executionRegistry) FindRequestLoop(ctx context.Context, reqRef insolar.Reference, apiRequestID string) bool {
	r.registryLock.Lock()
	defer r.registryLock.Unlock()

	if _, ok := r.registry[reqRef]; ok {
		// we're executing this request right now
		// de-duplication should deal with this case
		return false
	}

	for _, transcript := range r.registry {
		req := transcript.Request
		if req.APIRequestID == apiRequestID && req.ReturnMode != record.ReturnNoWait {
			return true
		}
	}

	return false
}

func (r *executionRegistry) GetActiveTranscript(request insolar.Reference) *common.Transcript {
	r.registryLock.Lock()
	defer r.registryLock.Unlock()

	return r.registry[request]
}
