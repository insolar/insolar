/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package logicrunner

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
)

func TestOnPulse(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	mb := testutils.NewMessageBusMock(t)
	mb.SendMock.Return(&reply.ID{}, nil)
	jc := testutils.NewJetCoordinatorMock(mc)

	lr, _ := NewLogicRunner(&configuration.LogicRunner{})
	lr.MessageBus = mb
	lr.JetCoordinator = jc

	jc.IsAuthorizedMock.Return(false, nil)
	jc.MeMock.Return(core.RecordRef{})

	// test empty lr
	pulse := core.Pulse{}

	err := lr.OnPulse(ctx, pulse)
	require.NoError(t, err)

	objectRef := testutils.RandomRef()

	// test empty ExecutionState
	lr.state[objectRef] = &ObjectState{ExecutionState: &ExecutionState{Behaviour: &ValidationSaver{}}}
	err = lr.OnPulse(ctx, pulse)
	require.NoError(t, err)
	assert.Nil(t, lr.state[objectRef])

	// test empty ExecutionState but not empty validation/consensus
	lr.state[objectRef] = &ObjectState{
		ExecutionState: &ExecutionState{
			Behaviour: &ValidationSaver{},
		},
		Validation: &ExecutionState{},
		Consensus:  &Consensus{},
	}
	err = lr.OnPulse(ctx, pulse)
	require.NoError(t, err)
	require.NotNil(t, lr.state[objectRef])
	assert.Nil(t, lr.state[objectRef].ExecutionState)

	// test empty es with query in current
	lr.state[objectRef] = &ObjectState{
		ExecutionState: &ExecutionState{
			Behaviour: &ValidationSaver{},
			Current:   &CurrentExecution{},
		},
	}
	err = lr.OnPulse(ctx, pulse)
	require.NoError(t, err)
	assert.Equal(t, message.InPending, lr.state[objectRef].ExecutionState.pending)
	qe := ExecutionQueueElement{}

	queue := append(make([]ExecutionQueueElement, 0), qe)

	lr.state[objectRef] = &ObjectState{
		ExecutionState: &ExecutionState{
			Behaviour: &ValidationSaver{},
			Current:   &CurrentExecution{},
			Queue:     queue,
		},
	}

	err = lr.OnPulse(ctx, pulse)
	require.NoError(t, err)
	require.Equal(t, message.InPending, lr.state[objectRef].ExecutionState.pending)

	// Executor in new pulse is same node
	jc.IsAuthorizedMock.Return(true, nil)
	lr.state[objectRef].ExecutionState.pending = message.PendingUnknown

	lr.state[objectRef] = &ObjectState{
		ExecutionState: &ExecutionState{
			Behaviour: &ValidationSaver{},
			Current:   &CurrentExecution{},
			Queue:     queue,
		},
	}

	err = lr.OnPulse(ctx, pulse)
	require.NoError(t, err)
	require.Equal(t, message.PendingUnknown, lr.state[objectRef].ExecutionState.pending)

	jc.IsAuthorizedMock.Return(true, nil)
	lr.state[objectRef].ExecutionState.pending = message.InPending

	err = lr.OnPulse(ctx, pulse)
	require.NoError(t, err)
	require.Equal(t, message.NotPending, lr.state[objectRef].ExecutionState.pending)

	lr.state[objectRef].ExecutionState.Current = nil
	lr.state[objectRef].ExecutionState.pending = message.InPending
	lr.state[objectRef].ExecutionState.PendingConfirmed = false

	err = lr.OnPulse(ctx, pulse)
	require.NoError(t, err)
	assert.Equal(t, message.NotPending, lr.state[objectRef].ExecutionState.pending)
	assert.Nil(t, lr.state[objectRef].ExecutionState.objectbody)
}

func TestPendingFinished(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	jc := testutils.NewJetCoordinatorMock(mc)
	mb := testutils.NewMessageBusMock(mc)
	ps := testutils.NewPulseStorageMock(mc)

	lr, _ := NewLogicRunner(&configuration.LogicRunner{})

	lr.JetCoordinator = jc
	lr.MessageBus = mb
	lr.PulseStorage = ps

	pulse := core.Pulse{}
	objectRef := testutils.RandomRef()
	meRef := testutils.RandomRef()

	jc.MeMock.Return(meRef)

	ps.CurrentMock.Return(&pulse, nil)

	es := &ExecutionState{
		Behaviour: &ValidationSaver{},
		Current:   &CurrentExecution{},
		pending:   message.NotPending,
	}

	// make sure that if there is no pending finishPendingIfNeeded returns false,
	// doesn't send PendingFinished message and doesn't change ExecutionState.pending
	lr.finishPendingIfNeeded(ctx, es, objectRef)
	require.Zero(t, mb.SendCounter)
	require.Equal(t, message.NotPending, es.pending)

	es.pending = message.InPending
	es.objectbody = &ObjectBody{}
	mb.SendMock.ExpectOnce(ctx, &message.PendingFinished{Reference: objectRef}, nil).Return(&reply.ID{}, nil)
	jc.IsAuthorizedMock.Return(false, nil)
	lr.finishPendingIfNeeded(ctx, es, objectRef)
	require.Equal(t, message.NotPending, es.pending)
	require.Nil(t, es.objectbody)

	mc.Wait(time.Second) // message bus' send is called in a goroutine

	es.pending = message.InPending
	es.objectbody = &ObjectBody{}
	jc.IsAuthorizedMock.Return(true, nil)
	lr.finishPendingIfNeeded(ctx, es, objectRef)
	require.Equal(t, message.NotPending, es.pending)
	require.NotNil(t, es.objectbody)
}

func TestStartQueueProcessorIfNeeded_DontStartQueueProcessorWhenPending(
	t *testing.T,
) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	am := testutils.NewArtifactManagerMock(t)
	lr, _ := NewLogicRunner(&configuration.LogicRunner{})
	lr.ArtifactManager = am

	objectRef := testutils.RandomRef()

	am.HasPendingRequestsMock.Return(true, nil)
	es := &ExecutionState{ArtifactManager: am, Queue: make([]ExecutionQueueElement, 0)}
	es.Queue = append(es.Queue, ExecutionQueueElement{})
	err := lr.StartQueueProcessorIfNeeded(
		ctx,
		es,
		&message.CallMethod{
			ObjectRef: objectRef,
			Method:    "some",
		},
	)
	require.NoError(t, err)
	require.Equal(t, message.InPending, es.pending)
}

func TestCheckPendingRequests(
	t *testing.T,
) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	objectRef := testutils.RandomRef()

	am := testutils.NewArtifactManagerMock(t)

	es := &ExecutionState{ArtifactManager: am}
	pending, err := es.CheckPendingRequests(
		ctx, &message.CallConstructor{},
	)
	require.NoError(t, err)
	require.Equal(t, message.NotPending, pending)

	am.HasPendingRequestsMock.Return(false, nil)
	es = &ExecutionState{ArtifactManager: am}
	pending, err = es.CheckPendingRequests(
		ctx, &message.CallMethod{
			ObjectRef: objectRef,
		},
	)
	require.NoError(t, err)
	require.Equal(t, message.NotPending, pending)

	am.HasPendingRequestsMock.Return(true, nil)
	es = &ExecutionState{ArtifactManager: am}
	pending, err = es.CheckPendingRequests(
		ctx, &message.CallMethod{
			ObjectRef: objectRef,
		},
	)
	require.NoError(t, err)
	require.Equal(t, message.InPending, pending)

	am.HasPendingRequestsMock.Return(false, errors.New("some"))
	es = &ExecutionState{ArtifactManager: am}
	pending, err = es.CheckPendingRequests(
		ctx, &message.CallMethod{
			ObjectRef: objectRef,
		},
	)
	require.Error(t, err)
	require.Equal(t, message.NotPending, pending)
}

func TestPrepareState(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	lr, _ := NewLogicRunner(&configuration.LogicRunner{})

	object := testutils.RandomRef()
	msg := &message.ExecutorResults{
		Caller:    testutils.RandomRef(),
		RecordRef: object,
	}

	// not pending
	// it's a first call, it's also initialize lr.state[object].ExecutionState
	// also check for empty Queue
	msg.Pending = message.NotPending
	_ = lr.prepareObjectState(ctx, msg)
	require.Equal(t, message.NotPending, lr.state[object].ExecutionState.pending)
	require.Equal(t, 0, len(lr.state[object].ExecutionState.Queue))

	// pending without queue
	lr.state[object].ExecutionState.pending = message.PendingUnknown
	msg.Pending = message.InPending
	_ = lr.prepareObjectState(ctx, msg)
	require.Equal(t, message.InPending, lr.state[object].ExecutionState.pending)

	// do not change pending status if it isn't unknown
	lr.state[object].ExecutionState.pending = message.NotPending
	msg.Pending = message.InPending
	_ = lr.prepareObjectState(ctx, msg)
	require.Equal(t, message.NotPending, lr.state[object].ExecutionState.pending)

	// do not change pending status if it isn't unknown
	lr.state[object].ExecutionState.pending = message.InPending
	msg.Pending = message.InPending
	_ = lr.prepareObjectState(ctx, msg)
	require.Equal(t, message.InPending, lr.state[object].ExecutionState.pending)

	parcel := testutils.NewParcelMock(t)
	parcel.ContextMock.Expect(context.Background()).Return(context.Background())
	// brand new queue from message
	msg.Queue = []message.ExecutionQueueElement{
		message.ExecutionQueueElement{Parcel: parcel},
	}
	_ = lr.prepareObjectState(ctx, msg)
	require.Equal(t, 1, len(lr.state[object].ExecutionState.Queue))

	testMsg := message.CallMethod{ReturnMode: message.ReturnNoWait}
	parcel = testutils.NewParcelMock(t)
	parcel.ContextMock.Expect(context.Background()).Return(context.Background())
	parcel.MessageMock.Return(&testMsg) // mock message that returns NoWait

	queueElementRequest := testutils.RandomRef()
	msg.Queue = []message.ExecutionQueueElement{message.ExecutionQueueElement{Request: &queueElementRequest, Parcel: parcel}}
	_ = lr.prepareObjectState(ctx, msg)
	require.Equal(t, 2, len(lr.state[object].ExecutionState.Queue))
	require.Equal(t, &queueElementRequest, lr.state[object].ExecutionState.Queue[0].request)
	require.Equal(t, &testMsg, lr.state[object].ExecutionState.Queue[0].parcel.Message())

}

func TestHandlePendingFinishedMessage(
	t *testing.T,
) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	objectRef := testutils.RandomRef()

	lr, _ := NewLogicRunner(&configuration.LogicRunner{})

	parcel := testutils.NewParcelMock(t).MessageMock.Return(
		&message.PendingFinished{Reference: objectRef},
	)

	parcel.DefaultTargetMock.Return(&core.RecordRef{})

	re, err := lr.HandlePendingFinishedMessage(ctx, parcel)
	require.NoError(t, err)
	require.Equal(t, &reply.OK{}, re)

	st := lr.MustObjectState(objectRef)

	es := st.ExecutionState
	require.NotNil(t, es)
	require.Equal(t, message.NotPending, es.pending)

	es.Current = &CurrentExecution{}
	re, err = lr.HandlePendingFinishedMessage(ctx, parcel)
	require.Error(t, err)

	es.Current = nil

	re, err = lr.HandlePendingFinishedMessage(ctx, parcel)
	require.NoError(t, err)
	require.Equal(t, &reply.OK{}, re)

}

func TestLogicRunner_CheckExecutionLoop(
	t *testing.T,
) {
	t.Parallel()
	ctx := inslogger.TestContext(t)

	mc := minimock.NewController(t)
	defer mc.Finish()

	lr, _ := NewLogicRunner(&configuration.LogicRunner{})

	es := &ExecutionState{
		Current: nil,
	}

	loop := lr.CheckExecutionLoop(ctx, es, nil)
	require.False(t, loop)

	ctxA, _ := inslogger.WithTraceField(ctx, "a")
	ctxB, _ := inslogger.WithTraceField(ctx, "b")

	parcel := testutils.NewParcelMock(t).MessageMock.Return(
		&message.CallMethod{ReturnMode: message.ReturnResult},
	)
	es.Current = &CurrentExecution{
		ReturnMode: message.ReturnResult,
		Context:    ctxA,
	}

	loop = lr.CheckExecutionLoop(ctxA, es, parcel)
	require.True(t, loop)

	loop = lr.CheckExecutionLoop(ctxB, es, parcel)
	require.False(t, loop)

	parcel = testutils.NewParcelMock(t).MessageMock.Return(
		&message.CallMethod{ReturnMode: message.ReturnNoWait},
	)
	es.Current = &CurrentExecution{
		ReturnMode: message.ReturnResult,
		Context:    ctxA,
	}
	loop = lr.CheckExecutionLoop(ctxA, es, parcel)
	require.False(t, loop)

	parcel = testutils.NewParcelMock(t).MessageMock.Return(
		&message.CallMethod{ReturnMode: message.ReturnResult},
	)
	es.Current = &CurrentExecution{
		ReturnMode: message.ReturnNoWait,
		Context:    ctxA,
	}
	loop = lr.CheckExecutionLoop(ctxA, es, parcel)
	require.False(t, loop)
}

func TestLogicRunner_HandleStillExecutingMessage(
	t *testing.T,
) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	objectRef := testutils.RandomRef()

	lr, _ := NewLogicRunner(&configuration.LogicRunner{})

	parcel := testutils.NewParcelMock(t).MessageMock.Return(
		&message.StillExecuting{Reference: objectRef},
	)

	parcel.DefaultTargetMock.Return(&core.RecordRef{})

	re, err := lr.HandleStillExecutingMessage(ctx, parcel)
	require.NoError(t, err)
	require.Equal(t, &reply.OK{}, re)

	st := lr.MustObjectState(objectRef)
	require.NotNil(t, st.ExecutionState)
	require.Equal(t, message.InPending, st.ExecutionState.pending)
	require.Equal(t, true, st.ExecutionState.PendingConfirmed)

	st.ExecutionState.pending = message.NotPending
	st.ExecutionState.PendingConfirmed = false

	re, err = lr.HandleStillExecutingMessage(ctx, parcel)
	require.NoError(t, err)
	require.Equal(t, &reply.OK{}, re)

	st = lr.MustObjectState(objectRef)
	require.NotNil(t, st.ExecutionState)
	require.Equal(t, message.NotPending, st.ExecutionState.pending)
	require.Equal(t, false, st.ExecutionState.PendingConfirmed)
}

func TestLogicRunner_OnPulse_StillExecuting(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	mb := testutils.NewMessageBusMock(t)
	jc := testutils.NewJetCoordinatorMock(mc)

	lr, _ := NewLogicRunner(&configuration.LogicRunner{})
	lr.MessageBus = mb
	lr.JetCoordinator = jc

	jc.IsAuthorizedMock.Return(false, nil)
	jc.MeMock.Return(core.RecordRef{})

	// test empty lr
	pulse := core.Pulse{}

	objectRef := testutils.RandomRef()

	lr.state[objectRef] = &ObjectState{
		ExecutionState: &ExecutionState{
			Behaviour: &ValidationSaver{},
			Current:   &CurrentExecution{},
		},
	}
	mb.SendMock.Return(&reply.OK{}, nil)
	err := lr.OnPulse(ctx, pulse)
	require.NoError(t, err)
	assert.NotNil(t, lr.state[objectRef].ExecutionState)
}

func TestReleaseQueue(t *testing.T) {
	t.Parallel()
	type expected struct {
		Length  int
		HasMore bool
	}
	type testCase struct {
		QueueLength int
		Expected    expected
	}
	var testCases = []testCase{
		{0, expected{0, false}},
		{1, expected{1, false}},
		{maxQueueLength, expected{maxQueueLength, false}},

		// TODO fix expected count to maxQueueLength after start taking data from ledger
		{maxQueueLength + 1, expected{maxQueueLength + 1, true}},
	}

	for _, tc := range testCases {
		es := ExecutionState{Queue: make([]ExecutionQueueElement, tc.QueueLength)}
		mq, hasMore := es.releaseQueue()
		assert.Equal(t, tc.Expected.Length, len(mq))
		assert.Equal(t, tc.Expected.HasMore, hasMore)
	}
}

func TestOnPulseLedgerHasMoreRequests(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	type testCase struct {
		queue                         []ExecutionQueueElement
		mbMock                        *testutils.MessageBusMock
		ExpectedLedgerHasMoreRequests bool
	}
	testCases := []testCase{
		{make([]ExecutionQueueElement, maxQueueLength+1), testutils.NewMessageBusMock(mc), true},
		{make([]ExecutionQueueElement, maxQueueLength), testutils.NewMessageBusMock(mc), false},
	}

	jc := testutils.NewJetCoordinatorMock(mc)
	jc.IsAuthorizedMock.Return(false, nil)
	jc.MeMock.Return(core.RecordRef{})

	pulse := core.Pulse{}

	for _, test := range testCases {
		queue := test.queue

		// waiting for ledger implement fetch method
		// waiting for us implement fetching
		messagesQueue := convertQueueToMessageQueue(queue)
		//messagesQueue := convertQueueToMessageQueue(queue[:maxQueueLength])

		ref := testutils.RandomRef()

		lr, _ := NewLogicRunner(&configuration.LogicRunner{})
		lr.JetCoordinator = jc

		lr.state[ref] = &ObjectState{
			ExecutionState: &ExecutionState{
				Behaviour: &ValidationSaver{},
				Queue:     queue,
			},
		}

		mb := test.mbMock
		lr.MessageBus = mb

		expectedMessage := &message.ExecutorResults{
			RecordRef:             ref,
			Requests:              make([]message.CaseBindRequest, 0),
			Queue:                 messagesQueue,
			LedgerHasMoreRequests: test.ExpectedLedgerHasMoreRequests,
		}

		// defer new SendFunc before calling OnPulse
		mb.SendFunc = func(p context.Context, p1 core.Message, p2 *core.MessageSendOptions) (r core.Reply, r1 error) {
			assert.Equal(t, expectedMessage, p1)
			return nil, nil
		}

		err := lr.OnPulse(ctx, pulse)
		require.NoError(t, err)
	}

	// waiting for all goroutines with Send() processing
	mc.Wait(10 * time.Second)
	for _, test := range testCases {
		assert.Equal(t, 1, int(test.mbMock.SendCounter))
	}
}

func TestLogicRunner_NoExcessiveAmends(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	am := testutils.NewArtifactManagerMock(t)
	lr, _ := NewLogicRunner(&configuration.LogicRunner{})
	lr.ArtifactManager = am
	am.UpdateObjectMock.Return(nil, nil)

	randRef := testutils.RandomRef()

	es := &ExecutionState{ArtifactManager: am, Queue: make([]ExecutionQueueElement, 0)}
	es.Queue = append(es.Queue, ExecutionQueueElement{})
	es.objectbody = &ObjectBody{}
	es.objectbody.CodeMachineType = core.MachineTypeBuiltin
	es.Current = &CurrentExecution{}
	es.Current.LogicContext = &core.LogicCallContext{}
	es.Current.Request = &randRef
	es.objectbody.CodeRef = &randRef

	data := []byte(testutils.RandomString())
	es.objectbody.Object = data

	mle := testutils.NewMachineLogicExecutorMock(t)
	lr.Executors[core.MachineTypeBuiltin] = mle
	mle.CallMethodMock.Return(data, nil, nil)

	msg := &message.CallMethod{
		ObjectRef: randRef,
		Method:    "some",
	}

	// In this case Update isn't send to ledger (objects data/newData are the same)
	am.RegisterResultMock.Return(nil, nil)

	_, err := lr.executeMethodCall(ctx, es, msg)
	require.NoError(t, err)
	require.Equal(t, uint64(0), am.UpdateObjectCounter)

	// In this case Update is send to ledger (objects data/newData are different)
	newData := make([]byte, 5, 5)
	mle.CallMethodMock.Return(newData, nil, nil)

	_, err = lr.executeMethodCall(ctx, es, msg)
	require.NoError(t, err)
	require.Equal(t, uint64(1), am.UpdateObjectCounter)
}

func TestHandleAbandonedRequestsNotificationMessage(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	objectId := testutils.RandomID()
	msg := &message.AbandonedRequestsNotification{Object: objectId}
	parcel := &message.Parcel{Msg: msg}

	// empty lr
	lr, _ := NewLogicRunner(&configuration.LogicRunner{})

	_, err := lr.HandleAbandonedRequestsNotificationMessage(ctx, parcel)
	require.NoError(t, err)
	assert.Equal(t, true, lr.state[*msg.DefaultTarget()].ExecutionState.LedgerHasMoreRequests)

	// LedgerHasMoreRequests false
	lr, _ = NewLogicRunner(&configuration.LogicRunner{})
	lr.state[*msg.DefaultTarget()] = &ObjectState{ExecutionState: &ExecutionState{LedgerHasMoreRequests: false}}

	_, err = lr.HandleAbandonedRequestsNotificationMessage(ctx, parcel)
	require.NoError(t, err)
	assert.Equal(t, true, lr.state[*msg.DefaultTarget()].ExecutionState.LedgerHasMoreRequests)

	// LedgerHasMoreRequests already true
	lr, _ = NewLogicRunner(&configuration.LogicRunner{})
	lr.state[*msg.DefaultTarget()] = &ObjectState{ExecutionState: &ExecutionState{LedgerHasMoreRequests: true}}

	_, err = lr.HandleAbandonedRequestsNotificationMessage(ctx, parcel)
	require.NoError(t, err)
	assert.Equal(t, true, lr.state[*msg.DefaultTarget()].ExecutionState.LedgerHasMoreRequests)

}

func TestPrepareObjectStateChangePendingStatus(t *testing.T) {
	ctx := inslogger.TestContext(t)
	lr, _ := NewLogicRunner(&configuration.LogicRunner{})
	ref := testutils.RandomRef()

	msg := &message.ExecutorResults{RecordRef: ref}

	// we are in pending and come to ourselves again
	lr.state[ref] = &ObjectState{ExecutionState: &ExecutionState{
		pending: message.InPending, Current: &CurrentExecution{}},
	}
	err := lr.prepareObjectState(ctx, msg)
	require.NoError(t, err)
	assert.Equal(t, message.NotPending, lr.state[ref].ExecutionState.pending)
	assert.Equal(t, false, lr.state[ref].ExecutionState.PendingConfirmed)

	// previous executor decline pending, trust him
	msg = &message.ExecutorResults{RecordRef: ref, Pending: message.NotPending}
	lr.state[ref] = &ObjectState{ExecutionState: &ExecutionState{
		pending: message.InPending, Current: nil},
	}
	err = lr.prepareObjectState(ctx, msg)
	require.NoError(t, err)
	assert.Equal(t, message.NotPending, lr.state[ref].ExecutionState.pending)
}

func TestPrepareObjectStateChangeLedgerHasMoreRequests(t *testing.T) {
	ctx := inslogger.TestContext(t)
	lr, _ := NewLogicRunner(&configuration.LogicRunner{})
	ref := testutils.RandomRef()

	msg := &message.ExecutorResults{RecordRef: ref}

	type testCase struct {
		messageStatus             bool
		objectStateStatus         bool
		expectedObjectStateStatue bool
	}

	testCases := []testCase{
		{true, true, true},
		{true, false, true},
		{false, true, true},
		{false, false, false},
	}

	for _, test := range testCases {
		msg = &message.ExecutorResults{RecordRef: ref, LedgerHasMoreRequests: test.messageStatus}
		lr.state[ref] = &ObjectState{ExecutionState: &ExecutionState{LedgerHasMoreRequests: test.objectStateStatus}}
		err := lr.prepareObjectState(ctx, msg)
		require.NoError(t, err)
		assert.Equal(t, test.expectedObjectStateStatue, lr.state[ref].ExecutionState.LedgerHasMoreRequests)
	}

}
