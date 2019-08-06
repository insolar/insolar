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

package executionarchive

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/transcript"
)

type Archiver interface {
	Archive(transcript *transcript.Transcript)
}

//go:generate minimock -i github.com/insolar/insolar/logicrunner/executionarchive.ExecutionArchive -o ./ -s _mock.go
type ExecutionArchive interface {
	Archive(ctx context.Context, transcript *transcript.Transcript)
	Done(transcript *transcript.Transcript) bool

	IsEmpty() bool
	OnPulse(ctx context.Context) []insolar.Message
	FindRequestLoop(ctx context.Context, reqRef insolar.Reference, apiRequestID string) bool
	GetActiveTranscript(req insolar.Reference) *transcript.Transcript
}

type executionArchive struct {
	// maps requestReference -> request to notify
	// (StillExecuting message) that we've started to work on in previous pulses
	objectRef   insolar.Reference
	archiveLock sync.Mutex
	archive     map[insolar.Reference]*transcript.Transcript

	jetCoordinator jet.Coordinator
}

func NewExecutionArchive(
	objectRef insolar.Reference,
	jetCoordinator jet.Coordinator,
) ExecutionArchive {

	return &executionArchive{
		objectRef:   objectRef,
		archiveLock: sync.Mutex{},
		archive:     make(map[insolar.Reference]*transcript.Transcript),

		jetCoordinator: jetCoordinator,
	}
}

func (ea *executionArchive) IsEmpty() bool {
	ea.archiveLock.Lock()
	defer ea.archiveLock.Unlock()

	return len(ea.archive) == 0
}

func (ea *executionArchive) Archive(ctx context.Context, transcript *transcript.Transcript) {
	requestRef := transcript.RequestRef

	ea.archiveLock.Lock()
	defer ea.archiveLock.Unlock()

	if _, ok := ea.archive[requestRef]; !ok {
		ea.archive[requestRef] = transcript
	} else {
		inslogger.FromContext(ctx).Error("Trying to archive task that is already archived")
	}
}

func (ea *executionArchive) Done(transcript *transcript.Transcript) bool {
	requestRef := transcript.RequestRef

	ea.archiveLock.Lock()
	defer ea.archiveLock.Unlock()

	if _, ok := ea.archive[requestRef]; !ok {
		return false
	}

	delete(ea.archive, requestRef)
	return true
}

// constructs all StillExecuting messages
func (ea *executionArchive) OnPulse(_ context.Context) []insolar.Message {
	ea.archiveLock.Lock()
	defer ea.archiveLock.Unlock()

	// TODO: this should return delegation token to continue execution of the pending
	messages := make([]insolar.Message, 0)
	if len(ea.archive) != 0 {
		requestRefs := make([]insolar.Reference, 0, len(ea.archive))
		for requestRef := range ea.archive {
			requestRefs = append(requestRefs, requestRef)
		}

		messages = append(messages, &message.StillExecuting{
			Reference:   ea.objectRef,
			Executor:    ea.jetCoordinator.Me(),
			RequestRefs: requestRefs,
		})
	}
	return messages
}

func (ea *executionArchive) FindRequestLoop(ctx context.Context, reqRef insolar.Reference, apiRequestID string) bool {
	ea.archiveLock.Lock()
	defer ea.archiveLock.Unlock()

	if _, ok := ea.archive[reqRef]; ok {
		// we're executing this request right now
		// de-duplication should deal with this case
		return false
	}

	for _, transcript := range ea.archive {
		req := transcript.Request
		if req.APIRequestID == apiRequestID && req.ReturnMode != record.ReturnNoWait {
			return true
		}
	}

	return false
}

func (ea *executionArchive) GetActiveTranscript(request insolar.Reference) *transcript.Transcript {
	return ea.archive[request]
}
