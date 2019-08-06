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
	"strings"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/transcript"
)

type ExecutionArchiveSuite struct{ suite.Suite }

func TestExecutionArchive(t *testing.T) { suite.Run(t, new(ExecutionArchiveSuite)) }

func (s *ExecutionArchiveSuite) genTranscriptForObject() *transcript.Transcript {
	ctx := inslogger.TestContext(s.T())
	return transcript.NewTranscript(ctx, gen.Reference(), record.IncomingRequest{
		ReturnMode:   record.ReturnResult,
		APIRequestID: s.genAPIRequestID(),
	})
}

func (s *ExecutionArchiveSuite) genAPIRequestID() string {
	APIRequestID := utils.RandTraceID()
	if strings.Contains(APIRequestID, "createRandomTraceIDFailed") {
		panic("Failed to generate uuid: " + APIRequestID)
	}
	return APIRequestID
}

func (s *ExecutionArchiveSuite) TestArchive() {
	mc := minimock.NewController(s.T())
	ctx := inslogger.TestContext(s.T())

	objectRef := gen.Reference()
	jc := jet.NewCoordinatorMock(mc)

	archiveI := NewExecutionArchive(objectRef, jc)
	archive := archiveI.(*executionArchive)
	firstTranscript := s.genTranscriptForObject()

	// successful archiving
	archiveI.Archive(ctx, firstTranscript)
	s.Len(archive.archive, 1)

	// duplicate
	archiveI.Archive(ctx, firstTranscript)
	s.Len(archive.archive, 1)

	// successful archiving
	archiveI.Archive(ctx, s.genTranscriptForObject())
	s.Len(archive.archive, 2)

	mc.Finish()
}

func (s *ExecutionArchiveSuite) TestDone() {
	mc := minimock.NewController(s.T())
	ctx := inslogger.TestContext(s.T())

	objectRef := gen.Reference()
	jc := jet.NewCoordinatorMock(mc)

	archiveI := NewExecutionArchive(objectRef, jc)
	archive := archiveI.(*executionArchive)
	T1, T2, T3 := s.genTranscriptForObject(), s.genTranscriptForObject(), s.genTranscriptForObject()

	archiveI.Archive(ctx, T1)
	archiveI.Archive(ctx, T2)
	s.Len(archive.archive, 2)

	s.False(archiveI.Done(T3))
	s.True(archiveI.Done(T2))
	s.False(archiveI.Done(T2))
	s.True(archiveI.Done(T1))
	s.False(archiveI.Done(T1))

	mc.Finish()
}

func (s *ExecutionArchiveSuite) TestIsEmpty() {
	mc := minimock.NewController(s.T())
	ctx := inslogger.TestContext(s.T())

	objectRef := gen.Reference()
	jc := jet.NewCoordinatorMock(mc)

	archiveI := NewExecutionArchive(objectRef, jc)
	archive := archiveI.(*executionArchive)

	s.True(archiveI.IsEmpty())

	T := s.genTranscriptForObject()
	archiveI.Archive(ctx, T)
	s.Len(archive.archive, 1)
	s.False(archiveI.IsEmpty())

	s.True(archiveI.Done(T))
	s.Len(archive.archive, 0)
	s.True(archiveI.IsEmpty())

	mc.Finish()
}

func (s *ExecutionArchiveSuite) TestOnPulse() {
	ctx := inslogger.TestContext(s.T())
	mc := minimock.NewController(s.T())

	meRef := gen.Reference()
	objectRef := gen.Reference()
	jc := jet.NewCoordinatorMock(mc).
		MeMock.Return(meRef)

	archiveI := NewExecutionArchive(objectRef, jc)
	{
		msgs := archiveI.OnPulse(ctx)
		s.Len(msgs, 0)
	}

	T1 := s.genTranscriptForObject()
	{
		archiveI.Archive(ctx, T1)
		msgs := archiveI.OnPulse(ctx)
		s.Len(msgs, 1)
		msg, ok := msgs[0].(*message.StillExecuting)
		s.Truef(ok, "expected message to be message.StillExecuting, got %T", msgs[0])
		s.Len(msg.RequestRefs, 1)
		s.Contains(msg.RequestRefs, T1.RequestRef)
		s.Equal(meRef, msg.Executor)
	}

	T2 := s.genTranscriptForObject()
	{
		archiveI.Archive(ctx, T2)
		msgs := archiveI.OnPulse(ctx)
		s.Len(msgs, 1)
		msg, ok := msgs[0].(*message.StillExecuting)
		s.Truef(ok, "expected message to be message.StillExecuting, got %T", msgs[0])
		s.Len(msg.RequestRefs, 2)
		s.Contains(msg.RequestRefs, T1.RequestRef)
		s.Contains(msg.RequestRefs, T2.RequestRef)
		s.Equal(meRef, msg.Executor)
	}

	archiveI.Done(T2)
	archiveI.Done(T1)
	{
		msgs := archiveI.OnPulse(ctx)
		s.Len(msgs, 0)
	}

	mc.Finish()
}

func (s *ExecutionArchiveSuite) TestFindRequestLoop() {
	ctx := inslogger.TestContext(s.T())
	mc := minimock.NewController(s.T())

	jc := jet.NewCoordinatorMock(mc)
	objRef := gen.Reference()
	reqRef := gen.Reference()

	archiveI := NewExecutionArchive(objRef, jc)
	{ // no requests with current apirequestid
		id := s.genAPIRequestID()

		s.False(archiveI.FindRequestLoop(ctx, reqRef, id))

		// cleanup after
		archiveI.(*executionArchive).archive = make(map[insolar.Reference]*transcript.Transcript)
	}

	T := s.genTranscriptForObject()
	{ // go request with current apirequestid (loop found)
		archiveI.Archive(ctx, T)

		s.True(archiveI.FindRequestLoop(ctx, reqRef, T.Request.APIRequestID))

		// cleanup after
		archiveI.(*executionArchive).archive = make(map[insolar.Reference]*transcript.Transcript)
	}

	{ // go request with current apirequestid, but record returnnowait (loop not found)
		id := s.genAPIRequestID()

		T.Request.ReturnMode = record.ReturnNoWait
		archiveI.Archive(ctx, T)

		s.False(archiveI.FindRequestLoop(ctx, reqRef, id))

		// cleanup after
		archiveI.(*executionArchive).archive = make(map[insolar.Reference]*transcript.Transcript)
	}

	T1 := s.genTranscriptForObject()
	T2 := s.genTranscriptForObject()
	T2.Request.ReturnMode = record.ReturnNoWait
	{ // combined test
		id := s.genAPIRequestID()

		archiveI.Archive(ctx, T1)
		archiveI.Archive(ctx, T2)

		s.False(archiveI.FindRequestLoop(ctx, reqRef, T2.Request.APIRequestID))
		s.True(archiveI.FindRequestLoop(ctx, reqRef, T1.Request.APIRequestID))
		s.False(archiveI.FindRequestLoop(ctx, reqRef, id))

		// cleanup after
		archiveI.(*executionArchive).archive = make(map[insolar.Reference]*transcript.Transcript)
	}

	mc.Finish()
}

func (s *ExecutionArchiveSuite) TestGetActiveTranscript() {
	ctx := inslogger.TestContext(s.T())
	mc := minimock.NewController(s.T())

	jc := jet.NewCoordinatorMock(mc)
	objRef := gen.Reference()

	T := s.genTranscriptForObject()
	archiveI := NewExecutionArchive(objRef, jc)
	archiveI.Archive(ctx, T)
	{ // have (put before)
		s.NotNil(archiveI.GetActiveTranscript(T.RequestRef))
	}
	{ // don't have
		s.Nil(archiveI.GetActiveTranscript(gen.Reference()))
	}

	archiveI.Done(T)
	{ // don't have (done task)
		s.Nil(archiveI.GetActiveTranscript(T.RequestRef))
	}
}
