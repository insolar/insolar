package executionregistry

import (
	"strings"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/common"
)

type ExecutionRegistrySuite struct{ suite.Suite }

func TestExecutionRegistry(t *testing.T) { suite.Run(t, new(ExecutionRegistrySuite)) }

func (s *ExecutionRegistrySuite) genTranscriptForObject() *common.Transcript {
	ctx := inslogger.TestContext(s.T())
	return common.NewTranscript(ctx, gen.Reference(), record.IncomingRequest{
		ReturnMode:   record.ReturnResult,
		APIRequestID: s.genAPIRequestID(),
	})
}

func (s *ExecutionRegistrySuite) genAPIRequestID() string {
	APIRequestID := utils.RandTraceID()
	if strings.Contains(APIRequestID, "createRandomTraceIDFailed") {
		panic("Failed to generate uuid: " + APIRequestID)
	}
	return APIRequestID
}

func (s *ExecutionRegistrySuite) TestRegister() {
	mc := minimock.NewController(s.T())
	ctx := inslogger.TestContext(s.T())

	objectRef := gen.Reference()
	jc := jet.NewCoordinatorMock(mc)

	registryI := New(objectRef, jc)
	registry := registryI.(*executionRegistry)
	firstTranscript := s.genTranscriptForObject()

	// successful archiving
	err := registryI.Register(ctx, firstTranscript)
	s.NoError(err)
	s.Len(registry.registry, 1)

	// duplicate
	err = registryI.Register(ctx, firstTranscript)
	s.Error(err)
	s.Len(registry.registry, 1)

	// successful archiving
	err = registryI.Register(ctx, s.genTranscriptForObject())
	s.NoError(err)
	s.Len(registry.registry, 2)

	mc.Finish()
}

func (s *ExecutionRegistrySuite) TestDone() {
	mc := minimock.NewController(s.T())
	ctx := inslogger.TestContext(s.T())

	objectRef := gen.Reference()
	jc := jet.NewCoordinatorMock(mc)

	registryI := New(objectRef, jc)
	registry := registryI.(*executionRegistry)
	T1, T2, T3 := s.genTranscriptForObject(), s.genTranscriptForObject(), s.genTranscriptForObject()

	err := registryI.Register(ctx, T1)
	s.NoError(err)

	err = registryI.Register(ctx, T2)
	s.NoError(err)
	s.Len(registry.registry, 2)

	s.False(registryI.Done(T3))
	s.True(registryI.Done(T2))
	s.False(registryI.Done(T2))
	s.True(registryI.Done(T1))
	s.False(registryI.Done(T1))

	mc.Finish()
}

func (s *ExecutionRegistrySuite) TestIsEmpty() {
	mc := minimock.NewController(s.T())
	ctx := inslogger.TestContext(s.T())

	objectRef := gen.Reference()
	jc := jet.NewCoordinatorMock(mc)

	registryI := New(objectRef, jc)
	registry := registryI.(*executionRegistry)

	s.True(registryI.IsEmpty())

	T := s.genTranscriptForObject()
	err := registryI.Register(ctx, T)
	s.NoError(err)
	s.Len(registry.registry, 1)
	s.False(registryI.IsEmpty())

	s.True(registryI.Done(T))
	s.Len(registry.registry, 0)
	s.True(registryI.IsEmpty())

	mc.Finish()
}

func (s *ExecutionRegistrySuite) TestOnPulse() {
	ctx := inslogger.TestContext(s.T())
	mc := minimock.NewController(s.T())

	meRef := gen.Reference()
	objectRef := gen.Reference()
	jc := jet.NewCoordinatorMock(mc).
		MeMock.Return(meRef)

	registryI := New(objectRef, jc)
	{
		msgs := registryI.OnPulse(ctx)
		s.Len(msgs, 0)
	}

	T1 := s.genTranscriptForObject()
	{
		err := registryI.Register(ctx, T1)
		s.NoError(err)
		msgs := registryI.OnPulse(ctx)
		s.Len(msgs, 1)
		msg, ok := msgs[0].(*payload.StillExecuting)
		s.Truef(ok, "expected message to be payload.StillExecuting, got %T", msgs[0])
		s.Len(msg.RequestRefs, 1)
		s.Contains(msg.RequestRefs, T1.RequestRef)
		s.Equal(meRef, msg.Executor)
	}

	T2 := s.genTranscriptForObject()
	{
		err := registryI.Register(ctx, T2)
		s.NoError(err)
		msgs := registryI.OnPulse(ctx)
		s.Len(msgs, 1)
		msg, ok := msgs[0].(*payload.StillExecuting)
		s.Truef(ok, "expected message to be message.StillExecuting, got %T", msgs[0])
		s.Len(msg.RequestRefs, 2)
		s.Contains(msg.RequestRefs, T1.RequestRef)
		s.Contains(msg.RequestRefs, T2.RequestRef)
		s.Equal(meRef, msg.Executor)
	}

	registryI.Done(T2)
	registryI.Done(T1)
	{
		msgs := registryI.OnPulse(ctx)
		s.Len(msgs, 0)
	}

	mc.Finish()
}

func (s *ExecutionRegistrySuite) TestGetActiveTranscript() {
	ctx := inslogger.TestContext(s.T())
	mc := minimock.NewController(s.T())

	jc := jet.NewCoordinatorMock(mc)
	objRef := gen.Reference()

	T := s.genTranscriptForObject()
	registryI := New(objRef, jc)
	err := registryI.Register(ctx, T)
	s.NoError(err)

	{ // have (put before)
		s.NotNil(registryI.GetActiveTranscript(T.RequestRef))
	}
	{ // don't have
		s.Nil(registryI.GetActiveTranscript(gen.Reference()))
	}

	registryI.Done(T)
	{ // don't have (done task)
		s.Nil(registryI.GetActiveTranscript(T.RequestRef))
	}
}
