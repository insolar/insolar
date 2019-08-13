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
	"math"
	"testing"
	"time"

	wmMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/executionregistry"
)

type publisherMock struct{}

func (p *publisherMock) Publish(topic string, messages ...*wmMessage.Message) error {
	return nil
}

func (p *publisherMock) Close() error {
	return nil
}

// wait is Exponential retries waiting function
// example usage: require.True(wait(func))
func wait(check func(...interface{}) bool, args ...interface{}) bool {
	for i := 0; i < 16; i++ {
		time.Sleep(time.Millisecond * time.Duration(math.Pow(2, float64(i))))
		if check(args...) {
			return true
		}
	}
	return false
}

type ExecutionBrokerSuite struct {
	suite.Suite
	Context    context.Context
	Controller *minimock.Controller
}

func TestExecutionBroker(t *testing.T) { suite.Run(t, new(ExecutionBrokerSuite)) }

func (s *ExecutionBrokerSuite) BeforeTest(suiteName, testName string) {
	s.Context = inslogger.TestContext(s.T())
	s.Controller = minimock.NewController(s.T())
}

func waitOnChannel(channel chan struct{}) bool {
	select {
	case <-channel:
		return true
	case <-time.After(1 * time.Minute):
		return false
	}
}

func immutableCount(args ...interface{}) bool {
	broker := args[0].(*ExecutionBroker)
	count := args[1].(int)

	return broker.immutable.Length() >= count
}

func finishedCount(args ...interface{}) bool {
	broker := args[0].(*ExecutionBroker)
	count := args[1].(int)

	return broker.finished.Length() >= count
}

func processorStatus(args ...interface{}) bool {
	broker := args[0].(*ExecutionBroker)
	status := args[1].(bool)

	return broker.isActiveProcessor() == status
}

func (s *ExecutionBrokerSuite) TestPut() {
	waitChannel := make(chan struct{})
	rem := NewRequestsExecutorMock(s.T())
	rem.ExecuteAndSaveMock.Set(func(_ context.Context, t *common.Transcript) (r insolar.Reply, r1 error) {
		if !t.Request.Immutable {
			waitChannel <- struct{}{}
		}
		return nil, nil
	})
	rem.SendReplyMock.Return()

	er := executionregistry.NewExecutionRegistryMock(s.T()).
		RegisterMock.Return().
		DoneMock.Return(true)

	objectRef := gen.Reference()
	b := NewExecutionBroker(objectRef, nil, rem, nil, nil, er, nil)
	b.pending = insolar.NotPending

	reqRef1 := gen.Reference()
	tr1 := common.NewTranscript(s.Context, reqRef1, record.IncomingRequest{})
	er.GetActiveTranscriptMock.When(reqRef1).Then(tr1)
	b.Put(s.Context, false, tr1)

	s.Equal(b.mutable.Length(), 1)

	reqRef2 := gen.Reference()
	tr2 := common.NewTranscript(s.Context, reqRef2, record.IncomingRequest{})
	er.GetActiveTranscriptMock.When(reqRef2).Then(tr2)
	b.Put(s.Context, true, tr2)

	s.True(waitOnChannel(waitChannel), "failed to wait until put triggers start of queue processor")
	s.True(waitOnChannel(waitChannel), "failed to wait until queue processor'll finish processing")
	s.Require().True(wait(processorStatus, b, false))
	s.Empty(waitChannel)

	s.Equal(b.mutable.Length(), 0)
	s.Equal(b.immutable.Length(), 0)
	s.Equal(b.finished.Length(), 2)

	rotationResults := b.rotate(10)
	s.Len(rotationResults.Requests, 0)
	s.Equal(rotationResults.LedgerHasMoreRequests, false)
	s.Len(rotationResults.Finished, 2)
}

func (s *ExecutionBrokerSuite) TestPrepend() {

	waitChannel := make(chan struct{})
	rem := NewRequestsExecutorMock(s.T())
	rem.ExecuteAndSaveMock.Set(func(_ context.Context, t *common.Transcript) (r insolar.Reply, r1 error) {
		if !t.Request.Immutable {
			waitChannel <- struct{}{}
		}
		return nil, nil
	})
	rem.SendReplyMock.Return()

	er := executionregistry.NewExecutionRegistryMock(s.T()).
		RegisterMock.Return().
		DoneMock.Return(true)

	objectRef := gen.Reference()
	b := NewExecutionBroker(objectRef, nil, rem, nil, nil, er, nil)
	b.pending = insolar.NotPending

	reqRef1 := gen.Reference()
	tr1 := common.NewTranscript(s.Context, reqRef1, record.IncomingRequest{})
	er.GetActiveTranscriptMock.When(reqRef1).Then(tr1)
	b.Prepend(s.Context, false, tr1)

	s.Equal(b.mutable.Length(), 1)

	reqRef2 := gen.Reference()
	tr2 := common.NewTranscript(s.Context, reqRef2, record.IncomingRequest{})
	er.GetActiveTranscriptMock.When(reqRef2).Then(tr2)
	b.Prepend(s.Context, true, tr2)

	s.Require().True(waitOnChannel(waitChannel), "failed to wait until put triggers start of queue processor")
	s.Require().True(waitOnChannel(waitChannel), "failed to wait until queue processor'll finish processing")
	s.Require().True(wait(processorStatus, b, false))
	s.Require().Empty(waitChannel)

	s.Equal(b.mutable.Length(), 0)
	s.Equal(b.immutable.Length(), 0)
	s.Equal(b.finished.Length(), 2)

	rotationResults := b.rotate(10)
	s.Len(rotationResults.Requests, 0)
	s.Equal(rotationResults.LedgerHasMoreRequests, false)
	s.Len(rotationResults.Finished, 2)
}

func (s *ExecutionBrokerSuite) TestImmutable_NotPending() {
	waitMutableChannel := make(chan struct{})
	waitImmutableChannel := make(chan struct{})

	rem := NewRequestsExecutorMock(s.T())
	rem.ExecuteAndSaveMock.Return(nil, nil)

	rem.SendReplyMock.Set(func(ctx context.Context, current *common.Transcript, re insolar.Reply, err error) {
		if !current.Request.Immutable {
			waitMutableChannel <- struct{}{}
		} else {
			waitImmutableChannel <- struct{}{}
		}
	})

	er := executionregistry.NewExecutionRegistryMock(s.T()).
		RegisterMock.Return().
		DoneMock.Return(true)

	objectRef := gen.Reference()
	b := NewExecutionBroker(objectRef, nil, rem, nil, nil, er, nil)
	b.pending = insolar.NotPending

	reqRef1 := gen.Reference()
	tr1 := common.NewTranscript(s.Context, reqRef1, record.IncomingRequest{Immutable: true})
	er.GetActiveTranscriptMock.When(reqRef1).Then(tr1)
	b.Prepend(s.Context, false, tr1)
	s.Equal(b.immutable.Length(), 1)

	reqRef2 := gen.Reference()
	tr2 := common.NewTranscript(s.Context, reqRef2, record.IncomingRequest{Immutable: true})
	er.GetActiveTranscriptMock.When(reqRef2).Then(tr2)
	b.Prepend(s.Context, true, tr2)

	s.Require().True(waitOnChannel(waitImmutableChannel), "failed to wait while processing is finished")
	s.Require().True(waitOnChannel(waitImmutableChannel), "failed to wait while processing is finished")
	s.Require().True(wait(processorStatus, b, false))
	s.Require().Empty(waitMutableChannel)

	s.Equal(0, b.mutable.Length())
	s.Equal(0, b.immutable.Length())
	s.Equal(2, b.finished.Length())

	rotationResults := b.rotate(10)
	s.Len(rotationResults.Requests, 0)
	s.Equal(rotationResults.LedgerHasMoreRequests, false)
	s.Len(rotationResults.Finished, 2)
}

func (s *ExecutionBrokerSuite) TestImmutable_InPending() {
	waitMutableChannel := make(chan struct{})
	waitImmutableChannel := make(chan struct{})

	rem := NewRequestsExecutorMock(s.T())
	rem.ExecuteAndSaveMock.Set(func(_ context.Context, t *common.Transcript) (r insolar.Reply, r1 error) {
		if !t.Request.Immutable {
			waitMutableChannel <- struct{}{}
		} else {
			waitImmutableChannel <- struct{}{}
		}
		return nil, nil
	})
	rem.SendReplyMock.Return()

	er := executionregistry.NewExecutionRegistryMock(s.T()).
		RegisterMock.Return().
		DoneMock.Return(true)

	objectRef := gen.Reference()
	b := NewExecutionBroker(objectRef, nil, rem, nil, nil, er, nil)
	b.pending = insolar.InPending

	tr1 := common.NewTranscript(s.Context, gen.Reference(), record.IncomingRequest{Immutable: true})
	er.GetActiveTranscriptMock.When(tr1.RequestRef).Then(tr1)
	b.Prepend(s.Context, false, tr1)

	s.Require().True(wait(immutableCount, b, 1), "failed to wait until immutable was put")
	s.Require().True(wait(processorStatus, b, false))

	tr2 := common.NewTranscript(s.Context, gen.Reference(), record.IncomingRequest{Immutable: true})
	er.GetActiveTranscriptMock.When(tr2.RequestRef).Then(tr2)
	b.Prepend(s.Context, true, tr2)

	s.Require().True(wait(immutableCount, b, 2), "failed to wait until immutable was put")
	s.Require().True(wait(processorStatus, b, false))
	s.Require().Empty(waitMutableChannel)

	b.StartProcessorsIfNeeded(s.Context)
	s.Require().True(wait(processorStatus, b, false))
	s.Empty(waitMutableChannel)
	s.Empty(waitImmutableChannel)

	rotationResults := b.rotate(10)
	s.Len(rotationResults.Requests, 2)
	s.Equal(rotationResults.LedgerHasMoreRequests, false)
	s.Len(rotationResults.Finished, 0)
}

func (s *ExecutionBrokerSuite) TestRotate() {
	objectRef := gen.Reference()
	b := NewExecutionBroker(objectRef, nil, nil, nil, nil, nil, nil)
	b.pending = insolar.NotPending

	for i := 0; i < 4; i++ {
		b.stateLock.Lock()
		b.immutable.Push(&common.Transcript{})
		b.mutable.Push(&common.Transcript{})
		b.stateLock.Unlock()
	}

	rotationResults := b.rotate(10)
	s.Equal(b.immutable.Length(), 0)
	s.Equal(b.mutable.Length(), 0)
	s.Equal(b.finished.Length(), 0)
	s.Len(rotationResults.Requests, 8)
	s.Len(rotationResults.Finished, 0)
	s.False(rotationResults.LedgerHasMoreRequests)

	for i := 0; i < 4; i++ {
		b.stateLock.Lock()
		b.immutable.Push(&common.Transcript{})
		b.stateLock.Unlock()
	}

	rotationResults = b.rotate(10)
	s.Equal(b.immutable.Length(), 0)
	s.Equal(b.mutable.Length(), 0)
	s.Equal(b.finished.Length(), 0)
	s.Len(rotationResults.Requests, 4)
	s.Len(rotationResults.Finished, 0)
	s.False(rotationResults.LedgerHasMoreRequests)

	for i := 0; i < 5; i++ {
		b.mutable.Push(&common.Transcript{})
		b.immutable.Push(&common.Transcript{})
	}

	rotationResults = b.rotate(10)
	s.Equal(b.immutable.Length(), 0)
	s.Equal(b.mutable.Length(), 0)
	s.Equal(b.finished.Length(), 0)
	s.Len(rotationResults.Requests, 10)
	s.Len(rotationResults.Finished, 0)
	s.False(rotationResults.LedgerHasMoreRequests)

	for i := 0; i < 6; i++ {
		b.mutable.Push(&common.Transcript{})
		b.immutable.Push(&common.Transcript{})
	}

	rotationResults = b.rotate(10)
	s.Equal(b.immutable.Length(), 0)
	s.Equal(b.mutable.Length(), 0)
	s.Equal(b.finished.Length(), 0)
	s.Len(rotationResults.Requests, 10)
	s.Len(rotationResults.Finished, 0)
	s.True(rotationResults.LedgerHasMoreRequests)
}

func (s *ExecutionBrokerSuite) TestDeduplication() {
	objectRef := gen.Reference()
	b := NewExecutionBroker(objectRef, nil, nil, nil, nil, nil, nil)

	b.pending = insolar.InPending

	reqRef1 := gen.Reference()
	b.Put(s.Context, false, common.NewTranscript(s.Context, reqRef1, record.IncomingRequest{})) // no duplication
	s.Equal(b.mutable.Length(), 1)

	b.Put(s.Context, false, common.NewTranscript(s.Context, reqRef1, record.IncomingRequest{})) // duplication
	s.Equal(b.mutable.Length(), 1)
}

func TestExecutionBroker_FinishedPendingIfNeed(t *testing.T) {
	mc := minimock.NewController(t)

	tests := []struct {
		name             string
		mocks            func(t minimock.Tester) *ExecutionBroker
		pending          insolar.PendingState
		pendingConfirmed bool
	}{
		{
			name: "success, complete",
			mocks: func(t minimock.Tester) *ExecutionBroker {
				objRef := gen.Reference()

				msg, err := payload.NewMessage(&payload.PendingFinished{
					ObjectRef: objRef,
				})
				require.NoError(t, err, "NewMessage")
				broker := &ExecutionBroker{
					Ref:     objRef,
					pending: insolar.InPending,

					sender: bus.NewSenderMock(mc).SendRoleMock.Set(
						func(_ context.Context, pendingMsg *wmMessage.Message, role insolar.DynamicRole, obj insolar.Reference) (r <-chan *wmMessage.Message, r1 func()) {

							return nil, func() {
								require.Equal(t, obj, objRef)
								require.Equal(t, insolar.DynamicRoleVirtualExecutor, role, "role")
								require.Equal(t, msg, pendingMsg)
							}
						}),
					executionRegistry: executionregistry.NewExecutionRegistryMock(t).
						IsEmptyMock.Return(true),
				}

				return broker
			},
			pending: insolar.NotPending,
		},
		{
			name: "success, not in pending",
			mocks: func(t minimock.Tester) *ExecutionBroker {
				obj := gen.Reference()
				broker := &ExecutionBroker{
					Ref:     obj,
					pending: insolar.NotPending,

					executionRegistry: executionregistry.NewExecutionRegistryMock(t),
				}
				return broker
			},
			pending: insolar.NotPending,
		},
		{
			name: "we have more unfinished requests",
			mocks: func(t minimock.Tester) *ExecutionBroker {
				obj := gen.Reference()
				broker := &ExecutionBroker{
					Ref:     obj,
					pending: insolar.InPending,

					executionRegistry: executionregistry.NewExecutionRegistryMock(t).
						IsEmptyMock.Return(false).
						LengthMock.Return(1),
				}

				return broker
			},
			pending: insolar.InPending,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := inslogger.TestContext(t)
			mc := minimock.NewController(t)

			broker := test.mocks(mc)
			broker.finishPendingIfNeeded(ctx)

			mc.Wait(1 * time.Minute)
			mc.Finish()

			require.Equal(t, test.pending, broker.pending)
			require.Equal(t, test.pendingConfirmed, broker.PendingConfirmed)
		})
	}
}

func TestExecutionBroker_ExecuteImmutable(t *testing.T) {
	// TODO .Put should become private, we should test interface
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(1 * time.Minute)

	er := executionregistry.NewExecutionRegistryMock(mc).
		RegisterMock.Return().
		DoneMock.Return(true)

	// prepare default object and execution state
	objectRef := gen.Reference()
	re := NewRequestsExecutorMock(mc)
	broker := NewExecutionBroker(objectRef, nil, re, nil, nil, er, nil)
	broker.pending = insolar.NotPending

	immutableRequestRef1 := gen.Reference()
	immutableRequest1 := record.IncomingRequest{
		ReturnMode:   record.ReturnResult,
		Object:       &objectRef,
		APIRequestID: utils.RandTraceID(),
		Immutable:    true,
	}
	immutableTranscript1 := common.NewTranscript(ctx, immutableRequestRef1, immutableRequest1)
	er.GetActiveTranscriptMock.When(immutableRequestRef1).Then(immutableTranscript1)

	re.ExecuteAndSaveMock.Return(&reply.CallMethod{Result: []byte{1, 2, 3}}, nil)
	re.SendReplyMock.Return()

	broker.Put(ctx, true, immutableTranscript1)
}

func TestExecutionBroker_OnPulse(t *testing.T) {
	randTranscript := func(ctx context.Context) *common.Transcript {
		reqRef := gen.Reference()
		return common.NewTranscript(ctx, reqRef, record.IncomingRequest{})
	}

	table := []struct {
		name string

		mocks            func(ctx context.Context, t minimock.Tester) *ExecutionBroker
		numberOfMessages int
		pending          insolar.PendingState
		pendingConfirmed bool
		ledgerHasMore    bool
		end              bool
	}{
		{
			name: "not active, queue",
			mocks: func(ctx context.Context, t minimock.Tester) *ExecutionBroker {
				objectRef := gen.Reference()
				er := executionregistry.NewExecutionRegistryMock(t).
					IsEmptyMock.Return(true)
				broker := NewExecutionBroker(objectRef, nil, nil, nil, nil, er, nil)
				// fetcher is stopped
				broker.requestsFetcher = NewRequestsFetcherMock(t).AbortMock.Return()
				broker.mutable.Push(randTranscript(ctx), randTranscript(ctx))
				return broker
			},
			numberOfMessages: 1,
			end:              true,
		},
		{
			// We aren't next executor but we're currently executing.
			// Expecting sending message to new executor and moving state to InPending
			name: "active, no queue",
			mocks: func(ctx context.Context, t minimock.Tester) *ExecutionBroker {
				objectRef := gen.Reference()
				er := executionregistry.NewExecutionRegistryMock(t).
					IsEmptyMock.Return(false)
				broker := NewExecutionBroker(objectRef, nil, nil, nil, nil, er, nil)

				return broker
			},
			numberOfMessages: 1,
			pending:          insolar.InPending,
		},
		{
			name: "not confirmed pending",
			mocks: func(ctx context.Context, t minimock.Tester) *ExecutionBroker {
				objectRef := gen.Reference()
				er := executionregistry.NewExecutionRegistryMock(t).
					IsEmptyMock.Return(true)
				broker := NewExecutionBroker(objectRef, nil, nil, nil, nil, er, nil)
				broker.pending = insolar.InPending
				return broker
			},
			numberOfMessages: 1,
			pending:          insolar.NotPending,
			ledgerHasMore:    true,
			end:              true,
		},
		{
			name: "not active, no pending, finished a request",
			mocks: func(ctx context.Context, t minimock.Tester) *ExecutionBroker {
				objectRef := gen.Reference()
				er := executionregistry.NewExecutionRegistryMock(t).
					IsEmptyMock.Return(true)
				broker := NewExecutionBroker(objectRef, nil, nil, nil, nil, er, nil)
				broker.finished.Push(randTranscript(ctx), randTranscript(ctx))
				return broker
			},
			numberOfMessages: 1,
			end:              true,
		},
		{
			name: "did nothing",
			mocks: func(ctx context.Context, t minimock.Tester) *ExecutionBroker {
				objectRef := gen.Reference()
				er := executionregistry.NewExecutionRegistryMock(t).
					IsEmptyMock.Return(true)
				broker := NewExecutionBroker(objectRef, nil, nil, nil, nil, er, nil)
				return broker
			},
			numberOfMessages: 0,
			end:              true,
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			ctx := inslogger.TestContext(t)
			mc := minimock.NewController(t)

			broker := test.mocks(ctx, mc)
			messages := broker.OnPulse(ctx)

			mc.Wait(1 * time.Minute)
			mc.Finish()

			require.Equal(t, test.pending, broker.pending)
			require.Equal(t, test.pendingConfirmed, broker.PendingConfirmed)
			require.Equal(t, test.end, !broker.isActive())
			require.Equal(t, test.ledgerHasMore, broker.ledgerHasMoreRequests)
			require.Len(t, messages, test.numberOfMessages)
		})
	}
}

func TestExecutionBroker_AddFreshRequestWithOnPulse(t *testing.T) {
	objectRef := gen.Reference()

	ctx := inslogger.TestContext(t)
	reqRef := gen.Reference()
	transcript := common.NewTranscript(ctx, reqRef, record.IncomingRequest{})

	table := []struct {
		name   string
		mocks  func(ctx context.Context, t minimock.Tester) (*ExecutionBroker, *[]payload.Payload)
		checks func(ctx context.Context, t *testing.T, msgs []payload.Payload)
	}{
		{
			name: "pulse change in HasPendings",
			mocks: func(ctx context.Context, t minimock.Tester) (*ExecutionBroker, *[]payload.Payload) {
				am := artifacts.NewClientMock(t)

				er := executionregistry.NewExecutionRegistryMock(t).
					IsEmptyMock.Return(true)
				broker := NewExecutionBroker(objectRef, nil, nil, nil, am, er, nil)

				var msgs []payload.Payload
				am.HasPendingsMock.Set(func(ctx context.Context, ref insolar.Reference) (bool, error) {
					msgs = broker.OnPulse(ctx)
					return false, nil
				})
				return broker, &msgs
			},
			checks: func(ctx context.Context, t *testing.T, msgs []payload.Payload) {
				require.Len(t, msgs, 1)
				results, ok := msgs[0].(*payload.ExecutorResults)
				require.True(t, ok)

				require.False(t, results.LedgerHasMoreRequests)
				require.Equal(t, insolar.PendingUnknown, results.Pending)
				require.Len(t, results.Queue, 1)
			},
		},
		{
			name: "pulse change in Execute",
			mocks: func(ctx context.Context, t minimock.Tester) (*ExecutionBroker, *[]payload.Payload) {
				doneCalled := false
				er := executionregistry.NewExecutionRegistryMock(t).
					IsEmptyMock.Set(func() bool { return doneCalled }).
					RegisterMock.Return().
					DoneMock.Set(func(_ *common.Transcript) bool { doneCalled = true; return true }).
					GetActiveTranscriptMock.When(reqRef).Then(transcript)
				am := artifacts.NewClientMock(t).
					HasPendingsMock.Return(false, nil)
				re := NewRequestsExecutorMock(t).
					SendReplyMock.Return()
				sender := bus.NewSenderMock(t).SendRoleMock.Return(nil, func() { return })

				broker := NewExecutionBroker(objectRef, nil, re, sender, am, er, nil)

				var msgs []payload.Payload
				re.ExecuteAndSaveMock.Set(func(ctx context.Context, tr *common.Transcript) (insolar.Reply, error) {
					msgs = broker.OnPulse(ctx)
					return &reply.OK{}, nil
				})
				return broker, &msgs
			},
			checks: func(ctx context.Context, t *testing.T, msgs []payload.Payload) {
				require.Len(t, msgs, 1)

				results, ok := msgs[0].(*payload.ExecutorResults)
				require.True(t, ok)
				require.False(t, results.LedgerHasMoreRequests)
				require.Equal(t, insolar.InPending, results.Pending)
				require.Len(t, results.Queue, 0)
			},
		},
	}

	for _, test := range table {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx := inslogger.TestContext(t)
			mc := minimock.NewController(t)

			broker, msgs := test.mocks(ctx, mc)
			broker.AddFreshRequest(ctx, transcript)

			mc.Wait(1 * time.Minute)
			mc.Finish()

			test.checks(ctx, t, *msgs)
		})
	}
}
