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

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
)

type Archiver interface {
	Archive(transcript *Transcript)
}

//go:generate minimock -i github.com/insolar/insolar/logicrunner.ExecutionArchive -o ./ -s _mock.go
type ExecutionArchive interface {
	Archive(transcript *Transcript)
	Done(transcript *Transcript) bool

	IsEmpty() bool
	OnPulse(ctx context.Context) []insolar.Message
	FindRequestLoop(ctx context.Context, apiRequestID string) bool
	GetActiveTranscript(req insolar.Reference) *Transcript
}

type executionArchive struct {
	// maps requestReference -> request to notify
	// (StillExecuting message) that we've started to work on in previous pulses
	objectRef   insolar.Reference
	archiveLock sync.Mutex
	archive     map[insolar.Reference]*Transcript

	jetCoordinator jet.Coordinator
}

func NewExecutionArchive(
	objectRef insolar.Reference,
	jetCoordinator jet.Coordinator,
) ExecutionArchive {

	return &executionArchive{
		objectRef:   objectRef,
		archiveLock: sync.Mutex{},
		archive:     make(map[insolar.Reference]*Transcript),

		jetCoordinator: jetCoordinator,
	}
}

func (sa *executionArchive) IsEmpty() bool {
	sa.archiveLock.Lock()
	defer sa.archiveLock.Unlock()

	return len(sa.archive) == 0
}

func (sa *executionArchive) Archive(transcript *Transcript) {
	requestRef := transcript.RequestRef

	sa.archiveLock.Lock()
	defer sa.archiveLock.Unlock()

	if _, ok := sa.archive[requestRef]; !ok {
		sa.archive[requestRef] = transcript
	}
}

func (sa *executionArchive) Done(transcript *Transcript) bool {
	requestRef := transcript.RequestRef

	sa.archiveLock.Lock()
	defer sa.archiveLock.Unlock()

	if _, ok := sa.archive[requestRef]; !ok {
		return false
	}

	delete(sa.archive, requestRef)
	return true
}

// constructs all StillExecuting messages
func (sa *executionArchive) OnPulse(_ context.Context) []insolar.Message {
	if sa == nil {
		return nil
	}

	sa.archiveLock.Lock()
	defer sa.archiveLock.Unlock()

	// TODO: this should return delegation token to continue execution of the pending
	messages := make([]insolar.Message, 0)
	if len(sa.archive) != 0 {
		requestRefs := make([]insolar.Reference, 0, len(sa.archive))
		for requestRef := range sa.archive {
			requestRefs = append(requestRefs, requestRef)
		}

		messages = append(messages, &message.StillExecuting{
			Reference:   sa.objectRef,
			Executor:    sa.jetCoordinator.Me(),
			RequestRefs: requestRefs,
		})
	}
	return messages
}

func (sa *executionArchive) FindRequestLoop(ctx context.Context, apiRequestID string) bool {
	sa.archiveLock.Lock()
	defer sa.archiveLock.Unlock()

	for _, transcript := range sa.archive {
		req := transcript.Request
		if req.APIRequestID == apiRequestID && req.ReturnMode != record.ReturnNoWait {
			return true
		}
	}

	return false
}

func (sa *executionArchive) GetActiveTranscript(request insolar.Reference) *Transcript {
	return sa.archive[request]
}
