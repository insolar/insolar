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
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/currentexecution"
	"github.com/insolar/insolar/logicrunner/executionarchive"
	"github.com/insolar/insolar/testutils"
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

	ea := executionarchive.NewExecutionArchiveMock(s.T()).ArchiveMock.Return().DoneMock.Return(true)

	objectRef := gen.Reference()
	b := NewExecutionBroker(objectRef, nil, rem, nil, nil, nil, nil, ea, nil)
	b.pending = insolar.NotPending

	tr := common.NewTranscript(s.Context, gen.Reference(), record.IncomingRequest{})

	b.Put(s.Context, false, tr)
	s.Equal(b.mutable.Length(), 1)

	tr = common.NewTranscript(s.Context, gen.Reference(), record.IncomingRequest{})

	b.Put(s.Context, true, tr)
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

	ea := executionarchive.NewExecutionArchiveMock(s.T()).ArchiveMock.Return().DoneMock.Return(true)

	objectRef := gen.Reference()
	b := NewExecutionBroker(objectRef, nil, rem, nil, nil, nil, nil, ea, nil)
	b.pending = insolar.NotPending

	reqRef1 := gen.Reference()
	tr := common.NewTranscript(s.Context, reqRef1, record.IncomingRequest{})
	b.Prepend(s.Context, false, tr)
	s.Equal(b.mutable.Length(), 1)

	reqRef2 := gen.Reference()
	tr = common.NewTranscript(s.Context, reqRef2, record.IncomingRequest{})
	b.Prepend(s.Context, true, tr)
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
	rem.ExecuteAndSaveMock.Set(func(_ context.Context, t *common.Transcript) (r insolar.Reply, r1 error) {
		if !t.Request.Immutable {
			waitMutableChannel <- struct{}{}
		} else {
			waitImmutableChannel <- struct{}{}
		}
		return nil, nil
	})
	rem.SendReplyMock.Return()

	ea := executionarchive.NewExecutionArchiveMock(s.T()).ArchiveMock.Return().DoneMock.Return(true)

	objectRef := gen.Reference()
	b := NewExecutionBroker(objectRef, nil, rem, nil, nil, nil, nil, ea, nil)
	b.pending = insolar.NotPending

	tr := common.NewTranscript(s.Context, gen.Reference(), record.IncomingRequest{Immutable: true})

	b.Prepend(s.Context, false, tr)
	s.Equal(b.immutable.Length(), 1)

	reqRef2 := gen.Reference()
	tr = common.NewTranscript(s.Context, reqRef2, record.IncomingRequest{Immutable: true})

	b.Prepend(s.Context, true, tr)
	s.Require().True(waitOnChannel(waitImmutableChannel), "failed to wait while processing is finished")
	s.Require().True(waitOnChannel(waitImmutableChannel), "failed to wait while processing is finished")
	s.Require().True(wait(processorStatus, b, false))
	s.Require().Empty(waitMutableChannel)

	s.Equal(b.mutable.Length(), 0)
	s.Equal(b.immutable.Length(), 0)
	s.Equal(b.finished.Length(), 2)

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

	ea := executionarchive.NewExecutionArchiveMock(s.T()).ArchiveMock.Return().DoneMock.Return(true)

	objectRef := gen.Reference()
	b := NewExecutionBroker(objectRef, nil, rem, nil, nil, nil, nil, ea, nil)
	b.pending = insolar.InPending

	tr := common.NewTranscript(s.Context, gen.Reference(), record.IncomingRequest{Immutable: true})

	b.Prepend(s.Context, false, tr)
	s.Require().True(wait(immutableCount, b, 1), "failed to wait until immutable was put")
	s.Require().True(wait(processorStatus, b, false))

	tr = common.NewTranscript(s.Context, gen.Reference(), record.IncomingRequest{Immutable: true})

	b.Prepend(s.Context, true, tr)

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
	b := NewExecutionBroker(objectRef, nil, nil, nil, nil, nil, nil, nil, nil)
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
	b := NewExecutionBroker(objectRef, nil, nil, nil, nil, nil, nil, nil, nil)

	b.pending = insolar.InPending

	reqRef1 := gen.Reference()
	b.Put(s.Context, false, common.NewTranscript(s.Context, reqRef1, record.IncomingRequest{})) // no duplication
	s.Equal(b.mutable.Length(), 1)

	b.Put(s.Context, false, common.NewTranscript(s.Context, reqRef1, record.IncomingRequest{})) // duplication
	s.Equal(b.mutable.Length(), 1)
}

func TestExecutionBroker_FinishPendingIfNeed(t *testing.T) {
	notEmptyCurrentList := currentexecution.NewList()
	transcript := common.NewTranscript(inslogger.TestContext(t), gen.Reference(), record.IncomingRequest{})
	err := notEmptyCurrentList.SetOnce(transcript)
	require.NoError(t, err)

	tests := []struct {
		name             string
		mocks            func(t minimock.Tester) *ExecutionBroker
		pending          insolar.PendingState
		pendingConfirmed bool
	}{
		{
			name: "success, complete",
			mocks: func(t minimock.Tester) *ExecutionBroker {
				obj := gen.Reference()
				broker := &ExecutionBroker{
					Ref:         obj,
					currentList: currentexecution.NewList(),
					pending:     insolar.InPending,

					jetCoordinator: jet.NewCoordinatorMock(t).
						IsMeAuthorizedNowMock.Return(false, nil),

					messageBus: testutils.NewMessageBusMock(t).SendMock.Return(&reply.OK{}, nil),
				}

				return broker
			},
			pending: insolar.NotPending,
		},
		{
			name: "success, me is next executor",
			mocks: func(t minimock.Tester) *ExecutionBroker {
				obj := gen.Reference()
				broker := &ExecutionBroker{
					Ref:         obj,
					currentList: currentexecution.NewList(),
					pending:     insolar.InPending,

					jetCoordinator: jet.NewCoordinatorMock(t).
						IsMeAuthorizedNowMock.Return(true, nil),
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
					Ref:         obj,
					currentList: currentexecution.NewList(),
					pending:     insolar.NotPending,
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
					Ref:         obj,
					currentList: notEmptyCurrentList,
					pending:     insolar.InPending,
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

	ea := executionarchive.NewExecutionArchiveMock(mc).ArchiveMock.Return().DoneMock.Return(true)

	// prepare default object and execution state
	objectRef := gen.Reference()
	re := NewRequestsExecutorMock(mc)
	broker := NewExecutionBroker(objectRef, nil, re, nil, nil, nil, nil, ea, nil)
	broker.pending = insolar.NotPending

	immutableRequestRef1 := gen.Reference()
	immutableRequest1 := record.IncomingRequest{
		ReturnMode:   record.ReturnResult,
		Object:       &objectRef,
		APIRequestID: utils.RandTraceID(),
		Immutable:    true,
	}
	immutableTranscript1 := common.NewTranscript(ctx, immutableRequestRef1, immutableRequest1)

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

		meNext           bool
		mocks            func(ctx context.Context, t minimock.Tester) *ExecutionBroker
		numberOfMessages int
		pending          insolar.PendingState
		pendingConfirmed bool
		ledgerHasMore    bool
		end              bool
	}{
		{
			name: "next is not me, not active, queue",
			mocks: func(ctx context.Context, t minimock.Tester) *ExecutionBroker {
				objectRef := gen.Reference()
				broker := NewExecutionBroker(objectRef, nil, nil, nil, nil, nil, nil, nil, nil)
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
			name: "next is not me, active, no queue",
			mocks: func(ctx context.Context, t minimock.Tester) *ExecutionBroker {
				objectRef := gen.Reference()
				broker := NewExecutionBroker(objectRef, nil, nil, nil, nil, nil, nil, nil, nil)
				broker.currentList.SetOnce(randTranscript(ctx))
				return broker
			},
			numberOfMessages: 1,
			pending:          insolar.InPending,
		},
		{
			name: "next is not me, not confirmed pending",
			mocks: func(ctx context.Context, t minimock.Tester) *ExecutionBroker {
				objectRef := gen.Reference()
				broker := NewExecutionBroker(objectRef, nil, nil, nil, nil, nil, nil, nil, nil)
				broker.pending = insolar.InPending
				return broker
			},
			numberOfMessages: 1,
			pending:          insolar.NotPending,
			ledgerHasMore:    true,
			end:              true,
		},
		{
			name: "next is not me, not active, no pending, finished a request",
			mocks: func(ctx context.Context, t minimock.Tester) *ExecutionBroker {
				objectRef := gen.Reference()
				broker := NewExecutionBroker(objectRef, nil, nil, nil, nil, nil, nil, nil, nil)
				broker.finished.Push(randTranscript(ctx), randTranscript(ctx))
				return broker
			},
			numberOfMessages: 1,
			end:              true,
		},
		{
			name: "next is not me, did nothing",
			mocks: func(ctx context.Context, t minimock.Tester) *ExecutionBroker {
				objectRef := gen.Reference()
				broker := NewExecutionBroker(objectRef, nil, nil, nil, nil, nil, nil, nil, nil)
				return broker
			},
			numberOfMessages: 0,
			end:              true,
		},
		{
			name:   "next is me, active",
			meNext: true,
			mocks: func(ctx context.Context, t minimock.Tester) *ExecutionBroker {
				objectRef := gen.Reference()
				broker := NewExecutionBroker(objectRef, nil, nil, nil, nil, nil, nil, nil, nil)
				broker.requestsFetcher = NewRequestsFetcherMock(t)
				broker.currentList.SetOnce(randTranscript(ctx))
				return broker
			},
			numberOfMessages: 0,
			pending:          insolar.InPending,
			pendingConfirmed: true,
		},
		{
			// We're the next executor, previous executor confirmed that this task
			// is executing and still in pending. We expect that previous executor
			// come to the current executor every pulse to confirm pending
			name:   "next is me, not active, in confirmed pending",
			meNext: true,
			mocks: func(ctx context.Context, t minimock.Tester) *ExecutionBroker {
				objectRef := gen.Reference()
				broker := NewExecutionBroker(objectRef, nil, nil, nil, nil, nil, nil, nil, nil)
				broker.requestsFetcher = NewRequestsFetcherMock(t)
				broker.pending = insolar.InPending
				broker.PendingConfirmed = true
				return broker
			},
			numberOfMessages: 0,
			pending:          insolar.InPending,
			end:              true,
		},
		{
			// We're the next executor and no one confirmed that this object is executing
			// restarting execution and fetching tasks off ledger
			name:   "next is me, not active, not confirmed pending",
			meNext: true,
			mocks: func(ctx context.Context, t minimock.Tester) *ExecutionBroker {
				objectRef := gen.Reference()
				broker := NewExecutionBroker(objectRef, nil, nil, nil, nil, nil, nil, nil, nil)
				broker.pending = insolar.InPending
				broker.requestsFetcher = NewRequestsFetcherMock(t).
					FetchPendingsMock.Return()
				return broker
			},
			numberOfMessages: 0,
			pending:          insolar.NotPending,
			ledgerHasMore:    true,
			end:              true,
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			ctx := inslogger.TestContext(t)
			mc := minimock.NewController(t)

			broker := test.mocks(ctx, mc)
			messages := broker.OnPulse(ctx, test.meNext)

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

	table := []struct {
		name   string
		mocks  func(ctx context.Context, t minimock.Tester) (*ExecutionBroker, *[]insolar.Message)
		checks func(ctx context.Context, t *testing.T, msgs []insolar.Message)
	}{
		{
			name: "pulse change in HasPendings",
			mocks: func(ctx context.Context, t minimock.Tester) (*ExecutionBroker, *[]insolar.Message) {
				am := artifacts.NewClientMock(t)

				broker := NewExecutionBroker(objectRef,
					nil, nil, nil,
					nil, nil, am, nil, nil)

				var msgs []insolar.Message
				am.HasPendingsMock.Set(func(ctx context.Context, ref insolar.Reference) (bool, error) {
					msgs = broker.OnPulse(ctx, false)
					return false, nil
				})
				return broker, &msgs
			},
			checks: func(ctx context.Context, t *testing.T, msgs []insolar.Message) {
				require.Len(t, msgs, 1)
				results, ok := msgs[0].(*message.ExecutorResults)
				require.True(t, ok)

				require.False(t, results.LedgerHasMoreRequests)
				require.Equal(t, insolar.PendingUnknown, results.Pending)
				require.Len(t, results.Queue, 1)
			},
		},
		{
			name: "pulse change in Execute",
			mocks: func(ctx context.Context, t minimock.Tester) (*ExecutionBroker, *[]insolar.Message) {
				ea := executionarchive.NewExecutionArchiveMock(t).
					ArchiveMock.Return().
					DoneMock.Return(true)
				am := artifacts.NewClientMock(t).
					HasPendingsMock.Return(false, nil)
				re := NewRequestsExecutorMock(t).
					SendReplyMock.Return()
				jc := jet.NewCoordinatorMock(t).
					IsMeAuthorizedNowMock.Return(true, nil)
				// pa := pulse.NewAccessorMock(t).LatestMock.Return(insolar.Pulse{}, nil)

				broker := NewExecutionBroker(objectRef, nil, re, nil, jc, nil, am, ea, nil)

				var msgs []insolar.Message
				re.ExecuteAndSaveMock.Set(func(ctx context.Context, tr *common.Transcript) (insolar.Reply, error) {
					msgs = broker.OnPulse(ctx, false)
					return &reply.OK{}, nil
				})
				return broker, &msgs
			},
			checks: func(ctx context.Context, t *testing.T, msgs []insolar.Message) {
				require.Len(t, msgs, 1)

				results, ok := msgs[0].(*message.ExecutorResults)
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

			reqRef := gen.Reference()
			broker.AddFreshRequest(ctx, common.NewTranscript(ctx, reqRef, record.IncomingRequest{}))

			mc.Wait(1 * time.Minute)
			mc.Finish()

			test.checks(ctx, t, *msgs)
		})
	}
}
