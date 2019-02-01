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
	"crypto/rand"
	"errors"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func GenerateRandomString(size int) ([]byte, error) {
	data := make([]byte, size)
	bytesRead, err := rand.Read(data)
	if bytesRead != size {
		return nil, err
	}
	return data, nil
}

type LogicRunnerTestSuite struct {
	suite.Suite

	mc  *minimock.Controller
	ctx context.Context
	am  *testutils.ArtifactManagerMock
	mb  *testutils.MessageBusMock
	jc  *testutils.JetCoordinatorMock
	lr  *LogicRunner
	es  ExecutionState
	ps  *testutils.PulseStorageMock
	mle *testutils.MachineLogicExecutorMock
}

func (suite *LogicRunnerTestSuite) BeforeTest(suiteName, testName string) {
	// testing context
	suite.ctx = inslogger.TestContext(suite.T())

	// initialize minimock and mocks
	suite.mc = minimock.NewController(suite.T())
	suite.am = testutils.NewArtifactManagerMock(suite.mc)
	suite.mb = testutils.NewMessageBusMock(suite.mc)
	suite.jc = testutils.NewJetCoordinatorMock(suite.mc)
	suite.ps = testutils.NewPulseStorageMock(suite.mc)

	suite.SetupLogicRunner()
}

func (suite *LogicRunnerTestSuite) SetupLogicRunner() {
	suite.lr, _ = NewLogicRunner(&configuration.LogicRunner{})
	suite.lr.ArtifactManager = suite.am
	suite.lr.MessageBus = suite.mb
	suite.lr.JetCoordinator = suite.jc
	suite.lr.PulseStorage = suite.ps
}

func (suite *LogicRunnerTestSuite) AfterTest(suiteName, testName string) {
	suite.mc.Finish()
}

func (suite *LogicRunnerTestSuite) TestOnPulse() {
	suite.jc.IsAuthorizedMock.Return(false, nil)
	suite.jc.MeMock.Return(core.RecordRef{})
	suite.mb.SendMock.Return(&reply.ID{}, nil)

	pulse := core.Pulse{}

	objectRef := testutils.RandomRef()

	// test empty lr
	err := suite.lr.OnPulse(suite.ctx, pulse)
	suite.NoError(err)

	// test empty ExecutionState
	suite.lr.state[objectRef] = &ObjectState{
		ExecutionState: &ExecutionState{
			Behaviour: &ValidationSaver{},
		},
	}
	err = suite.lr.OnPulse(suite.ctx, pulse)
	suite.NoError(err)
	suite.Nil(suite.lr.state[objectRef])

	// test empty ExecutionState but not empty ValidationSaver/Consensus
	suite.lr.state[objectRef] = &ObjectState{
		ExecutionState: &ExecutionState{
			Behaviour: &ValidationSaver{},
		},
		Validation: &ExecutionState{},
		Consensus:  &Consensus{},
	}
	err = suite.lr.OnPulse(suite.ctx, pulse)
	suite.NoError(err)
	suite.NotNil(suite.lr.state[objectRef])
	suite.Nil(suite.lr.state[objectRef].ExecutionState)

	// test empty ExecutionState with query in action
	suite.lr.state[objectRef] = &ObjectState{
		ExecutionState: &ExecutionState{
			Behaviour: &ValidationSaver{},
			Current:   &CurrentExecution{},
		},
	}
	err = suite.lr.OnPulse(suite.ctx, pulse)
	suite.NoError(err)
	suite.Equal(message.InPending, suite.lr.state[objectRef].ExecutionState.pending)

	//
	suite.lr.state[objectRef] = &ObjectState{
		ExecutionState: &ExecutionState{
			Behaviour: &ValidationSaver{},
			Current:   &CurrentExecution{},
			Queue:     append(make([]ExecutionQueueElement, 0), ExecutionQueueElement{}),
		},
	}
	err = suite.lr.OnPulse(suite.ctx, pulse)
	suite.NoError(err)
	suite.Equal(message.InPending, suite.lr.state[objectRef].ExecutionState.pending)

	// Executor in new pulse is same node
	suite.jc.IsAuthorizedMock.Return(true, nil)

	suite.lr.state[objectRef] = &ObjectState{
		ExecutionState: &ExecutionState{
			Behaviour: &ValidationSaver{},
			Current:   &CurrentExecution{},
			Queue:     append(make([]ExecutionQueueElement, 0), ExecutionQueueElement{}),
		},
	}
	err = suite.lr.OnPulse(suite.ctx, pulse)
	suite.NoError(err)
	suite.Equal(message.PendingUnknown, suite.lr.state[objectRef].ExecutionState.pending)

	suite.mc.Wait(time.Second)
}

func (suite *LogicRunnerTestSuite) TestPendingFinished() {
	pulse := core.Pulse{}
	objectRef := testutils.RandomRef()
	meRef := testutils.RandomRef()

	suite.jc.MeMock.Return(meRef)
	suite.ps.CurrentMock.Return(&pulse, nil)

	// make sure that if there is no pending finishPendingIfNeeded returns false,
	// doesn't send PendingFinished message and doesn't change ExecutionState.pending
	es := &ExecutionState{
		Behaviour: &ValidationSaver{},
		Current:   &CurrentExecution{},
		pending:   message.NotPending,
	}

	suite.lr.finishPendingIfNeeded(suite.ctx, es, objectRef)
	suite.Zero(suite.mb.SendCounter)
	suite.Equal(message.NotPending, es.pending)

	//
	es = &ExecutionState{
		Behaviour:  &ValidationSaver{},
		Current:    &CurrentExecution{},
		pending:    message.InPending,
		objectbody: &ObjectBody{},
	}

	suite.mb.SendMock.ExpectOnce(suite.ctx, &message.PendingFinished{Reference: objectRef}, nil).Return(&reply.ID{}, nil)
	suite.jc.IsAuthorizedMock.Return(false, nil)

	suite.lr.finishPendingIfNeeded(suite.ctx, es, objectRef)
	suite.Equal(message.NotPending, es.pending)
	suite.Nil(es.objectbody)

	// message bus' send is called in a goroutine
	suite.mc.Wait(time.Second)

	//
	es = &ExecutionState{
		Behaviour:  &ValidationSaver{},
		Current:    &CurrentExecution{},
		pending:    message.InPending,
		objectbody: &ObjectBody{},
	}

	suite.jc.IsAuthorizedMock.Return(true, nil)
	suite.lr.finishPendingIfNeeded(suite.ctx, es, objectRef)

	suite.Equal(message.NotPending, es.pending)
	suite.NotNil(es.objectbody)
}

func (suite *LogicRunnerTestSuite) TestStartQueueProcessorIfNeeded_DontStartQueueProcessorWhenPending() {
	objectRef := testutils.RandomRef()
	suite.am.HasPendingRequestsMock.Return(true, nil)
	es := &ExecutionState{
		ArtifactManager: suite.am,
		Queue:           append(make([]ExecutionQueueElement, 0), ExecutionQueueElement{}),
	}
	msg := &message.CallMethod{ObjectRef: objectRef, Method: "some"}

	err := suite.lr.StartQueueProcessorIfNeeded(suite.ctx, es, msg)
	suite.NoError(err)
	suite.Equal(message.InPending, es.pending)
}

func (suite *LogicRunnerTestSuite) TestCheckPendingRequests() {
	objectRef := testutils.RandomRef()
	msg := &message.CallMethod{ObjectRef: objectRef}
	es := &ExecutionState{ArtifactManager: suite.am}

	suite.am.HasPendingRequestsMock.Return(false, nil)
	pending, err := es.CheckPendingRequests(suite.ctx, msg)
	suite.NoError(err)
	suite.Equal(message.NotPending, pending)

	suite.am.HasPendingRequestsMock.Return(true, nil)
	pending, err = es.CheckPendingRequests(suite.ctx, msg)
	suite.NoError(err)
	suite.Equal(message.InPending, pending)

	suite.am.HasPendingRequestsMock.Return(false, errors.New("some"))
	pending, err = es.CheckPendingRequests(suite.ctx, msg)
	suite.Error(err)
	suite.Equal(message.NotPending, pending)
}

func (suite *LogicRunnerTestSuite) TestPrepareState() {
	objectRef := testutils.RandomRef()
	msg := &message.ExecutorResults{Caller: testutils.RandomRef(), RecordRef: objectRef}

	// not pending
	// it's a first call, it's also initialize lr.state[objectRef].ExecutionState
	// also check for empty Queue
	msg.Pending = message.NotPending
	_ = suite.lr.prepareObjectState(suite.ctx, msg)
	suite.Equal(message.NotPending, suite.lr.state[objectRef].ExecutionState.pending)
	suite.Equal(0, len(suite.lr.state[objectRef].ExecutionState.Queue))

	// pending without queue
	suite.lr.state[objectRef].ExecutionState.pending = message.PendingUnknown
	msg.Pending = message.InPending
	_ = suite.lr.prepareObjectState(suite.ctx, msg)
	suite.Equal(message.InPending, suite.lr.state[objectRef].ExecutionState.pending)

	// do not change pending status if it isn't unknown
	suite.lr.state[objectRef].ExecutionState.pending = message.NotPending
	msg.Pending = message.InPending
	_ = suite.lr.prepareObjectState(suite.ctx, msg)
	suite.Equal(message.NotPending, suite.lr.state[objectRef].ExecutionState.pending)

	// do not change pending status if it isn't unknown
	suite.lr.state[objectRef].ExecutionState.pending = message.InPending
	msg.Pending = message.InPending
	_ = suite.lr.prepareObjectState(suite.ctx, msg)
	suite.Equal(message.InPending, suite.lr.state[objectRef].ExecutionState.pending)

	parcel := testutils.NewParcelMock(suite.mc)
	parcel.ContextMock.Expect(context.Background()).Return(context.Background())
	// brand new queue from message
	msg.Queue = []message.ExecutionQueueElement{
		message.ExecutionQueueElement{Parcel: parcel},
	}
	_ = suite.lr.prepareObjectState(suite.ctx, msg)
	suite.Equal(1, len(suite.lr.state[objectRef].ExecutionState.Queue))

	testMsg := message.CallMethod{ReturnMode: message.ReturnNoWait}
	parcel = testutils.NewParcelMock(suite.mc)
	parcel.ContextMock.Expect(context.Background()).Return(context.Background())
	parcel.MessageMock.Return(&testMsg) // mock message that returns NoWait

	queueElementRequest := testutils.RandomRef()
	msg.Queue = []message.ExecutionQueueElement{
		message.ExecutionQueueElement{Request: &queueElementRequest, Parcel: parcel},
	}
	_ = suite.lr.prepareObjectState(suite.ctx, msg)
	suite.Equal(2, len(suite.lr.state[objectRef].ExecutionState.Queue))
	suite.Equal(&queueElementRequest, suite.lr.state[objectRef].ExecutionState.Queue[0].request)
	suite.Equal(&testMsg, suite.lr.state[objectRef].ExecutionState.Queue[0].parcel.Message())
}

func (suite *LogicRunnerTestSuite) TestHandlePendingFinishedMessage() {
	objectRef := testutils.RandomRef()
	msg := &message.PendingFinished{Reference: objectRef}
	parcel := testutils.NewParcelMock(suite.mc)
	parcel.MessageMock.Return(msg)
	parcel.DefaultTargetMock.Return(&core.RecordRef{})

	re, err := suite.lr.HandlePendingFinishedMessage(suite.ctx, parcel)
	suite.NoError(err)
	suite.Equal(&reply.OK{}, re)

	st := suite.lr.MustObjectState(objectRef)
	suite.Require().NotNil(st)

	es := st.ExecutionState
	suite.Require().NotNil(es)
	suite.Equal(message.NotPending, es.pending)

	es.Current = nil
	re, err = suite.lr.HandlePendingFinishedMessage(suite.ctx, parcel)
	suite.NoError(err)
	suite.Equal(&reply.OK{}, re)
}

func (suite *LogicRunnerTestSuite) TestCheckExecutionLoop() {
	es := &ExecutionState{Current: nil}

	loop := suite.lr.CheckExecutionLoop(suite.ctx, es, nil)
	suite.False(loop)

	ctxA, _ := inslogger.WithTraceField(suite.ctx, "a")
	ctxB, _ := inslogger.WithTraceField(suite.ctx, "b")

	msg := &message.CallMethod{ReturnMode: message.ReturnResult}
	parcel := testutils.NewParcelMock(suite.mc)
	parcel.MessageMock.Return(msg)

	es.Current = &CurrentExecution{
		ReturnMode: message.ReturnResult,
		Context:    ctxA,
	}

	suite.True(suite.lr.CheckExecutionLoop(ctxA, es, parcel))
	suite.False(suite.lr.CheckExecutionLoop(ctxB, es, parcel))

	msg = &message.CallMethod{ReturnMode: message.ReturnNoWait}
	parcel = testutils.NewParcelMock(suite.mc)
	parcel.MessageMock.Return(msg)

	es.Current = &CurrentExecution{
		ReturnMode: message.ReturnResult,
		Context:    ctxA,
	}
	suite.False(suite.lr.CheckExecutionLoop(ctxA, es, parcel))

	msg = &message.CallMethod{ReturnMode: message.ReturnResult}
	parcel = testutils.NewParcelMock(suite.mc)

	es.Current = &CurrentExecution{
		ReturnMode: message.ReturnNoWait,
		Context:    ctxA,
	}
	suite.False(suite.lr.CheckExecutionLoop(ctxA, es, parcel))
}

func (suite *LogicRunnerTestSuite) SetupCreateObjectAndSet() core.RecordRef {
	objectRef := testutils.RandomRef()

	suite.lr.state[objectRef] = &ObjectState{
		ExecutionState: &ExecutionState{
			Behaviour: &ValidationSaver{},
			Current:   &CurrentExecution{},
		},
	}

	return objectRef
}

func (suite *LogicRunnerTestSuite) Test_OnPulse_StillExecuting() {
	suite.jc.IsAuthorizedMock.Return(false, nil)
	suite.jc.MeMock.Return(core.RecordRef{})

	pulse := core.Pulse{}
	objectRef := suite.SetupCreateObjectAndSet()

	err := suite.lr.OnPulse(suite.ctx, pulse)
	suite.Require().NoError(err)
	suite.Assert().NotNil(suite.lr.state[objectRef].ExecutionState)
}

func (suite *LogicRunnerTestSuite) TestLogicRunner_HandleStillExecutingMessage() {
	objectRef := testutils.RandomRef()

	msg := &message.StillExecuting{Reference: objectRef}
	parcel := testutils.NewParcelMock(suite.mc)
	parcel.MessageMock.Return(msg)
	parcel.DefaultTargetMock.Return(&core.RecordRef{})

	//
	re, err := suite.lr.HandleStillExecutingMessage(suite.ctx, parcel)
	suite.NoError(err)
	suite.Equal(&reply.OK{}, re)

	st := suite.lr.MustObjectState(objectRef)
	suite.Require().NotNil(st)
	suite.Equal(message.InPending, st.ExecutionState.pending)
	suite.True(st.ExecutionState.PendingConfirmed)

	//
	st.ExecutionState.pending = message.NotPending
	st.ExecutionState.PendingConfirmed = false

	re, err = suite.lr.HandleStillExecutingMessage(suite.ctx, parcel)
	suite.NoError(err)
	suite.Equal(&reply.OK{}, re)

	st = suite.lr.MustObjectState(objectRef)
	suite.Require().NotNil(st)
	suite.Equal(message.NotPending, st.ExecutionState.pending)
	suite.False(st.ExecutionState.PendingConfirmed)
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

func (suite *LogicRunnerTestSuite) TestOnPulseLedgerHasMoreRequests() {
	type testCase struct {
		queue                         []ExecutionQueueElement
		ExpectedLedgerHasMoreRequests bool
	}
	testCases := []testCase{
		{make([]ExecutionQueueElement, maxQueueLength+1), true},
		{make([]ExecutionQueueElement, maxQueueLength), false},
	}

	suite.jc.IsAuthorizedMock.Return(false, nil)
	suite.jc.MeMock.Return(core.RecordRef{})

	pulse := core.Pulse{}

	for _, test := range testCases {
		queue := test.queue

		// waiting for ledger implement fetch method
		// waiting for us implement fetching
		messagesQueue := convertQueueToMessageQueue(queue)
		//messagesQueue := convertQueueToMessageQueue(queue[:maxQueueLength])

		ref := testutils.RandomRef()

		lr, _ := NewLogicRunner(&configuration.LogicRunner{})
		lr.JetCoordinator = suite.jc

		lr.state[ref] = &ObjectState{
			ExecutionState: &ExecutionState{
				Behaviour: &ValidationSaver{},
				Queue:     queue,
			},
		}

		mb := testutils.NewMessageBusMock(suite.mc)
		lr.MessageBus = mb

		expectedMessage := &message.ExecutorResults{
			RecordRef:             ref,
			Requests:              make([]message.CaseBindRequest, 0),
			Queue:                 messagesQueue,
			LedgerHasMoreRequests: test.ExpectedLedgerHasMoreRequests,
		}

		counter := 0
		// defer new SendFunc before calling OnPulse
		mb.SendMock.Set(func(p context.Context, p1 core.Message, p2 *core.MessageSendOptions) (r core.Reply, r1 error) {
			suite.Equal(expectedMessage, p1)
			suite.Equal(0, counter)
			counter = counter + 1
			return nil, nil
		})

		err := lr.OnPulse(suite.ctx, pulse)
		suite.Require().NoError(err)
	}

	// waiting for all goroutines with Send() processing
	suite.mc.Wait(10 * time.Second)
}

func (suite *LogicRunnerTestSuite) TestNoExcessiveAmends() {
	data, err := GenerateRandomString(128)
	suite.Require().NoError(err)

	requestRef := testutils.RandomRef()
	codeRef := testutils.RandomRef()

	es := ExecutionState{
		ArtifactManager: suite.am,
		Queue:           make([]ExecutionQueueElement, 0),
		Current: &CurrentExecution{
			LogicContext: &core.LogicCallContext{},
			Request:      &requestRef,
		},
		objectbody: &ObjectBody{
			CodeRef:         &codeRef,
			Object:          data,
			CodeMachineType: core.MachineTypeBuiltin,
		},
	}
	mle := testutils.NewMachineLogicExecutorMock(suite.mc)
	suite.lr.Executors[core.MachineTypeBuiltin] = mle

	msg := &message.CallMethod{ObjectRef: testutils.RandomRef(), Method: "some"}

	suite.am.RegisterResultMock.Return(nil, nil)
	suite.am.UpdateObjectMock.Return(nil, nil)

	// In this case Update isn't send to ledger (objects data/newData are the same)
	mle.CallMethodMock.Return(data, nil, nil)
	_, err = suite.lr.executeMethodCall(suite.ctx, &es, msg)
	suite.NoError(err)
	suite.Equal(0, int(suite.am.UpdateObjectCounter))

	// In this case Update is send to ledger (objects data/newData are different)
	changedData, err := GenerateRandomString(128)
	suite.Require().NoError(err)
	suite.Require().NotEqual(data, changedData)

	mle.CallMethodMock.Return(changedData, nil, nil)

	_, err = suite.lr.executeMethodCall(suite.ctx, &es, msg)
	suite.NoError(err)
	suite.Equal(1, int(suite.am.UpdateObjectCounter))
}

func (suite *LogicRunnerTestSuite) TestHandleAbandonedRequestsNotificationMessage() {
	objectId := testutils.RandomID()
	msg := &message.AbandonedRequestsNotification{Object: objectId}
	parcel := &message.Parcel{Msg: msg}

	// empty lr
	_, err := suite.lr.HandleAbandonedRequestsNotificationMessage(suite.ctx, parcel)
	suite.Require().NoError(err)
	suite.True(suite.lr.state[*msg.DefaultTarget()].ExecutionState.LedgerHasMoreRequests)

	// LedgerHasMoreRequests false
	suite.lr.state[*msg.DefaultTarget()] = &ObjectState{
		ExecutionState: &ExecutionState{
			LedgerHasMoreRequests: false,
		},
	}
	_, err = suite.lr.HandleAbandonedRequestsNotificationMessage(suite.ctx, parcel)
	suite.Require().NoError(err)
	suite.True(suite.lr.state[*msg.DefaultTarget()].ExecutionState.LedgerHasMoreRequests)

	// LedgerHasMoreRequests already true
	suite.lr.state[*msg.DefaultTarget()] = &ObjectState{
		ExecutionState: &ExecutionState{
			LedgerHasMoreRequests: true,
		},
	}
	_, err = suite.lr.HandleAbandonedRequestsNotificationMessage(suite.ctx, parcel)
	suite.Require().NoError(err)
	suite.True(suite.lr.state[*msg.DefaultTarget()].ExecutionState.LedgerHasMoreRequests)
}

func TestLogicRunner(t *testing.T) {
	suite.Run(t, new(LogicRunnerTestSuite))
}
