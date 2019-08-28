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
	"errors"
	"math"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	wmMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/executionregistry"
	"github.com/insolar/insolar/logicrunner/requestsqueue"
	"github.com/insolar/insolar/pulse"
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

func waitOnChannel(channel chan struct{}) bool {
	select {
	case <-channel:
		return true
	case <-time.After(1 * time.Minute):
		return false
	}
}

func finishedCount(args ...interface{}) bool {
	broker := args[0].(*ExecutionBroker)
	count := args[1].(int)

	broker.stateLock.Lock()
	defer broker.stateLock.Unlock()
	return len(broker.finished) >= count
}

func TestExecutionBroker_AddFreshRequest(t *testing.T) {
	objectRef := gen.Reference()

	ctx := inslogger.TestContext(t)
	reqRef := gen.Reference()
	transcript := common.NewTranscript(ctx, reqRef, record.IncomingRequest{})

	table := []struct {
		name  string
		mocks func(ctx context.Context, t minimock.Tester) *ExecutionBroker
	}{
		{
			name: "happy path",
			mocks: func(ctx context.Context, t minimock.Tester) *ExecutionBroker {
				er := executionregistry.NewExecutionRegistryMock(t).
					RegisterMock.Return(nil).
					DoneMock.Return(true).
					GetActiveTranscriptMock.When(reqRef).Then(nil)
				am := artifacts.NewClientMock(t).
					HasPendingsMock.Return(false, nil)
				re := NewRequestsExecutorMock(t).
					SendReplyMock.Return()
				broker := NewExecutionBroker(objectRef, nil, re, nil, am, er, nil, nil)

				re.ExecuteAndSaveMock.Set(func(ctx context.Context, tr *common.Transcript) (insolar.Reply, error) {
					return &reply.OK{}, nil
				})
				return broker
			},
		},
	}

	for _, test := range table {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx := inslogger.TestContext(t)
			mc := minimock.NewController(t)

			broker := test.mocks(ctx, mc)
			broker.AddFreshRequest(ctx, transcript)

			mc.Wait(1 * time.Minute)
			mc.Finish()
		})
	}
}

func TestExecutionBroker_Deduplication(t *testing.T) {
	objectRef := gen.Reference()
	reqRef := gen.Reference()

	tests := []struct {
		name   string
		mocks  func(ctx context.Context, t minimock.Tester) *ExecutionBroker
		checks func(t *testing.T, b *ExecutionBroker)
	}{
		{
			name: "request known to broker, gets deduplicated",
			mocks: func(ctx context.Context, t minimock.Tester) *ExecutionBroker {
				er := executionregistry.NewExecutionRegistryMock(t).
					GetActiveTranscriptMock.Expect(reqRef).Return(nil)

				b := NewExecutionBroker(
					objectRef, nil, nil, nil, nil, er, nil, nil,
				)

				queueMock := requestsqueue.NewRequestsQueueMock(t).AppendMock.Return()
				b.mutable.queue = queueMock

				tr := common.NewTranscript(ctx, reqRef, record.IncomingRequest{})
				b.add(ctx, requestsqueue.FromLedger, tr)

				return b
			},
			checks: func(t *testing.T, b *ExecutionBroker) {
				appended := b.mutable.queue.(*requestsqueue.RequestsQueueMock).AppendAfterCounter()
				require.Equal(t, 1, int(appended))
			},
		},
		{
			name: "request NOT known to broker, but registry, gets deduplicated",
			mocks: func(ctx context.Context, t minimock.Tester) *ExecutionBroker {
				tr := common.NewTranscript(ctx, reqRef, record.IncomingRequest{})
				er := executionregistry.NewExecutionRegistryMock(t).
					GetActiveTranscriptMock.Expect(reqRef).Return(tr)

				b := NewExecutionBroker(
					objectRef, nil, nil, nil, nil, er, nil, nil,
				)

				queueMock := requestsqueue.NewRequestsQueueMock(t)
				b.mutable.queue = queueMock

				return b
			},
			checks: func(t *testing.T, b *ExecutionBroker) {
				appended := b.mutable.queue.(*requestsqueue.RequestsQueueMock).AppendAfterCounter()
				require.Equal(t, 0, int(appended))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := inslogger.TestContext(t)
			mc := minimock.NewController(t)

			broker := test.mocks(ctx, mc)

			tr := common.NewTranscript(ctx, reqRef, record.IncomingRequest{})
			broker.add(ctx, requestsqueue.FromLedger, tr)

			mc.Wait(1 * time.Minute)
			mc.Finish()

			test.checks(t, broker)
		})
	}

}

func TestExecutionBroker_PendingFinishedIfNeed(t *testing.T) {
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

					pulseAccessor: insolarPulse.NewAccessorMock(t).LatestMock.Set(func(p context.Context) (r insolar.Pulse, r1 error) {
						return insolar.Pulse{
							PulseNumber:     insolar.PulseNumber(pulse.MinTimePulse),
							NextPulseNumber: insolar.PulseNumber(pulse.MinTimePulse + 1),
						}, nil
					}),

					sender: bus.NewSenderMock(mc).SendRoleMock.Set(
						func(_ context.Context, pendingMsg *wmMessage.Message, role insolar.DynamicRole, obj insolar.Reference) (<-chan *wmMessage.Message, func()) {
							r := make(chan *wmMessage.Message)
							go func() {
								r <- wmMessage.NewMessage(watermill.NewUUID(), nil)
							}()
							return r, func() {
								require.Equal(t, obj, objRef)
								require.Equal(t, insolar.DynamicRoleVirtualExecutor, role, "role")
								require.Equal(t, msg.Payload, pendingMsg.Payload)
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
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(1 * time.Minute)

	er := executionregistry.NewExecutionRegistryMock(mc).
		RegisterMock.Return(nil).
		DoneMock.Return(true)

	// prepare default object and execution state
	objectRef := gen.Reference()
	re := NewRequestsExecutorMock(mc)
	broker := NewExecutionBroker(objectRef, nil, re, nil, nil, er, nil, nil)
	broker.pending = insolar.NotPending

	immutableRequestRef1 := gen.Reference()
	immutableRequest1 := record.IncomingRequest{
		ReturnMode:   record.ReturnResult,
		Object:       &objectRef,
		APIRequestID: utils.RandTraceID(),
		Immutable:    true,
	}
	immutableTranscript1 := common.NewTranscript(ctx, immutableRequestRef1, immutableRequest1)
	er.GetActiveTranscriptMock.When(immutableRequestRef1).Then(nil).
		DoneMock.Return(true)

	re.ExecuteAndSaveMock.Return(&reply.CallMethod{Result: []byte{1, 2, 3}}, nil)
	re.SendReplyMock.Return()

	broker.AddFreshRequest(ctx, immutableTranscript1)
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
				broker := NewExecutionBroker(objectRef, nil, nil, nil, nil, er, nil, nil)
				// fetcher is stopped
				broker.requestsFetcher = NewRequestsFetcherMock(t).AbortMock.Return()
				broker.mutable.queue.Append(ctx, requestsqueue.FromLedger, randTranscript(ctx), randTranscript(ctx))
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
				broker := NewExecutionBroker(objectRef, nil, nil, nil, nil, er, nil, nil)

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
				broker := NewExecutionBroker(objectRef, nil, nil, nil, nil, er, nil, nil)
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
				broker := NewExecutionBroker(objectRef, nil, nil, nil, nil, er, nil, nil)
				broker.mutable.queue.Append(ctx, requestsqueue.FromLedger, randTranscript(ctx), randTranscript(ctx))
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
				broker := NewExecutionBroker(objectRef, nil, nil, nil, nil, er, nil, nil)
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
					GetActiveTranscriptMock.Return(nil).
					IsEmptyMock.Return(true)
				broker := NewExecutionBroker(objectRef, nil, nil, nil, am, er, nil, nil)

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
					RegisterMock.Return(nil).
					DoneMock.Set(func(_ *common.Transcript) bool { doneCalled = true; return true }).
					GetActiveTranscriptMock.When(reqRef).Then(nil)
				am := artifacts.NewClientMock(t).
					HasPendingsMock.Return(false, nil)
				re := NewRequestsExecutorMock(t).
					SendReplyMock.Return()
				sender := bus.NewSenderMock(t).SendRoleMock.Return(nil, func() { return })
				pulseMock := insolarPulse.NewAccessorMock(t).LatestMock.Set(func(p context.Context) (r insolar.Pulse, r1 error) {
					return insolar.Pulse{
						PulseNumber:     insolar.PulseNumber(pulse.MinTimePulse),
						NextPulseNumber: insolar.PulseNumber(pulse.MinTimePulse + 1),
					}, nil
				})
				broker := NewExecutionBroker(objectRef, nil, re, sender, am, er, nil, pulseMock)

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

func TestExecutionBroker_IsKnownRequest(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	reqRef1 := gen.Reference()
	er := executionregistry.NewExecutionRegistryMock(mc).GetActiveTranscriptMock.Expect(reqRef1).Return(nil)

	objectRef := gen.Reference()
	b := NewExecutionBroker(
		objectRef, nil, nil, nil, nil, er, nil, nil,
	)

	queueMock := requestsqueue.NewRequestsQueueMock(mc).AppendMock.Return()
	b.mutable.queue = queueMock

	tr := common.NewTranscript(ctx, reqRef1, record.IncomingRequest{})
	b.add(ctx, requestsqueue.FromLedger, tr)

	require.True(t, b.IsKnownRequest(ctx, reqRef1))

	require.False(t, b.IsKnownRequest(ctx, gen.Reference()))
}

func TestExecutionBroker_MoreRequestsOnLedger(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	objectRef := gen.Reference()
	b := NewExecutionBroker(
		objectRef, nil, nil, nil, nil, nil, nil, nil,
	)
	b.MoreRequestsOnLedger(ctx)
	require.True(t, b.ledgerHasMoreRequests)
	require.Empty(t, b.requestsFetcher)
}

func TestExecutionBroker_NoMoreRequestsOnLedger(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	objectRef := gen.Reference()
	b := NewExecutionBroker(
		objectRef, nil, nil, nil, nil, nil, nil, nil,
	)

	b.ledgerHasMoreRequests = true
	b.requestsFetcher = NewRequestsFetcherMock(mc).AbortMock.Return()
	b.NoMoreRequestsOnLedger(ctx)

	require.False(t, b.ledgerHasMoreRequests)
}

func TestExecutionBroker_AbandonedRequestsOnLedger(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	objectRef := gen.Reference()
	b := NewExecutionBroker(
		objectRef, nil, nil, nil, nil, nil, nil, nil,
	)

	b.requestsFetcher = NewRequestsFetcherMock(mc).FetchPendingsMock.Return()
	b.AbandonedRequestsOnLedger(ctx)
}

func TestExecutionBroker_AbandonedRequestsOnLedger_Integration(t *testing.T) {
	mc := minimock.NewController(t)

	objectRef := gen.Reference()

	tests := []struct {
		name  string
		mocks func(t minimock.Tester) *ExecutionBroker
	}{
		{
			name: "no requests on ledger",
			mocks: func(t minimock.Tester) *ExecutionBroker {
				am := artifacts.NewClientMock(mc).GetPendingsMock.
					Return([]insolar.Reference{}, insolar.ErrNoPendingRequest)

				b := NewExecutionBroker(
					objectRef, nil, nil, nil, am, nil, nil, nil,
				)
				return b
			},
		},
		{
			name: "request on ledger, abort during fetch",
			mocks: func(t minimock.Tester) *ExecutionBroker {
				reqRef := gen.Reference()
				am := artifacts.NewClientMock(mc).
					GetPendingsMock.
					Return([]insolar.Reference{reqRef}, nil)
				b := NewExecutionBroker(
					objectRef, nil, nil, nil, am, nil, nil, nil,
				)
				am.GetAbandonedRequestMock.Set(func(ctx context.Context, o insolar.Reference, r insolar.Reference) (record.Request, error) {
					b.stopRequestsFetcher(ctx)
					return nil, nil
				})

				return b
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := inslogger.TestContext(t)
			mc := minimock.NewController(t)

			broker := test.mocks(mc)
			broker.AbandonedRequestsOnLedger(ctx)

			mc.Wait(1 * time.Minute)
			mc.Finish()
		})
	}
}

func TestExecutionBroker_PrevExecutorPendingResult(t *testing.T) {
	objectRef := gen.Reference()

	tests := []struct {
		name   string
		state  insolar.PendingState
		mocks  func(t minimock.Tester) *ExecutionBroker
		checks func(t *testing.T, b *ExecutionBroker)
	}{
		{
			name:  "local unknown, using remote",
			state: insolar.NotPending,
			mocks: func(t minimock.Tester) *ExecutionBroker {
				b := NewExecutionBroker(
					objectRef, nil, nil, nil, nil, nil, nil, nil,
				)
				return b
			},
			checks: func(t *testing.T, b *ExecutionBroker) {
				require.Equal(t, insolar.NotPending, b.pending)
			},
		},
		{
			name:  "in pending, no executions, prev said continue",
			state: insolar.NotPending,
			mocks: func(t minimock.Tester) *ExecutionBroker {
				er := executionregistry.NewExecutionRegistryMock(t).IsEmptyMock.Return(true)
				b := NewExecutionBroker(
					objectRef, nil, nil, nil, nil, er, nil, nil,
				)
				b.pending = insolar.InPending
				return b
			},
			checks: func(t *testing.T, b *ExecutionBroker) {
				require.Equal(t, insolar.NotPending, b.pending)
				require.False(t, b.PendingConfirmed)
			},
		},
		{
			name:  "local execution, in pending, ignoring",
			state: insolar.NotPending,
			mocks: func(t minimock.Tester) *ExecutionBroker {
				er := executionregistry.NewExecutionRegistryMock(t).IsEmptyMock.Return(false)
				b := NewExecutionBroker(
					objectRef, nil, nil, nil, nil, er, nil, nil,
				)
				b.pending = insolar.InPending
				return b
			},
			checks: func(t *testing.T, b *ExecutionBroker) {
				require.Equal(t, insolar.InPending, b.pending)
				require.False(t, b.PendingConfirmed)
			},
		},
		{
			name:  "local not pending, ignoring remote",
			state: insolar.InPending,
			mocks: func(t minimock.Tester) *ExecutionBroker {
				b := NewExecutionBroker(
					objectRef, nil, nil, nil, nil, nil, nil, nil,
				)
				b.pending = insolar.NotPending
				return b
			},
			checks: func(t *testing.T, b *ExecutionBroker) {
				require.Equal(t, insolar.NotPending, b.pending)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := inslogger.TestContext(t)
			mc := minimock.NewController(t)

			broker := test.mocks(mc)
			broker.PrevExecutorPendingResult(ctx, test.state)

			mc.Wait(1 * time.Minute)
			mc.Finish()

			test.checks(t, broker)
		})
	}
}

func TestExecutionBroker_getTask(t *testing.T) {
	tests := []struct {
		name    string
		mocks   func(ctx context.Context, t minimock.Tester) *ExecutionBroker
		hasTask bool
	}{
		{
			name:    "happy path, got task",
			hasTask: true,
			mocks: func(ctx context.Context, t minimock.Tester) *ExecutionBroker {
				er := executionregistry.NewExecutionRegistryMock(t).
					GetActiveTranscriptMock.Return(nil).
					RegisterMock.Return(nil)

				objectRef := gen.Reference()
				b := NewExecutionBroker(
					objectRef, nil, nil, nil, nil, er, nil, nil,
				)

				reqRef := gen.Reference()
				tr := common.NewTranscript(ctx, reqRef, record.IncomingRequest{})
				b.add(ctx, requestsqueue.FromLedger, tr)

				return b
			},
		},
		{
			name:    "no task, empty queue",
			hasTask: false,
			mocks: func(ctx context.Context, t minimock.Tester) *ExecutionBroker {
				objectRef := gen.Reference()
				b := NewExecutionBroker(
					objectRef, nil, nil, nil, nil, nil, nil, nil,
				)
				return b
			},
		},
		{
			name:    "no task, already in the registry",
			hasTask: false,
			mocks: func(ctx context.Context, t minimock.Tester) *ExecutionBroker {
				er := executionregistry.NewExecutionRegistryMock(t).
					GetActiveTranscriptMock.Return(nil).
					RegisterMock.Return(errors.New("some"))

				objectRef := gen.Reference()
				b := NewExecutionBroker(
					objectRef, nil, nil, nil, nil, er, nil, nil,
				)

				reqRef := gen.Reference()
				tr := common.NewTranscript(ctx, reqRef, record.IncomingRequest{})
				b.add(ctx, requestsqueue.FromLedger, tr)

				return b
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := inslogger.TestContext(t)
			mc := minimock.NewController(t)

			broker := test.mocks(ctx, mc)
			task := broker.getTask(ctx, broker.mutable.queue)

			mc.Wait(1 * time.Minute)
			mc.Finish()

			require.Equal(t, test.hasTask, task != nil)
		})
	}
}
