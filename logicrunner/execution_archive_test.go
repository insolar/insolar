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
	"testing"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

type ExecutionArchiveSuite struct{ suite.Suite }

func TestExecutionArchive(t *testing.T) { suite.Run(t, new(ExecutionArchiveSuite)) }

func (s *ExecutionArchiveSuite) genTranscriptForObject() *Transcript {
	ctx := inslogger.TestContext(s.T())
	return NewTranscript(ctx, gen.Reference(), record.IncomingRequest{})
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
