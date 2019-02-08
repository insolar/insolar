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
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
)

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
	suite.mc.Wait(60 * time.Second)
	suite.mc.Finish()
}

func (suite *LogicRunnerTestSuite) TestOnPulse() {
	suite.T().Skip()
	// TODO in test case where we are executor again need check for queue start, or make it active before test
	suite.mb.SendMock.Return(&reply.ID{}, nil)
	suite.am.GetPendingRequestMock.Return(nil, nil)

	suite.jc.IsAuthorizedMock.Return(false, nil)
	suite.jc.MeMock.Return(core.RecordRef{})

	// test empty lr
	pulse := core.Pulse{}
	suite.ps.CurrentMock.Return(&pulse, nil)

	err := suite.lr.OnPulse(suite.ctx, pulse)
	suite.Require().NoError(err)

	objectRef := testutils.RandomRef()

	// test empty ExecutionState
	suite.lr.state[objectRef] = &ObjectState{ExecutionState: &ExecutionState{Behaviour: &ValidationSaver{}}}
	err = suite.lr.OnPulse(suite.ctx, pulse)
	suite.Require().NoError(err)
	suite.Nil(suite.lr.state[objectRef])

	// test empty ExecutionState but not empty validation/consensus
	suite.lr.state[objectRef] = &ObjectState{
		ExecutionState: &ExecutionState{
			Behaviour: &ValidationSaver{},
		},
		Validation: &ExecutionState{},
		Consensus:  &Consensus{},
	}
	err = suite.lr.OnPulse(suite.ctx, pulse)
	suite.Require().NoError(err)
	suite.Require().NotNil(suite.lr.state[objectRef])
	suite.Nil(suite.lr.state[objectRef].ExecutionState)

	// test empty es with query in current
	suite.lr.state[objectRef] = &ObjectState{
		ExecutionState: &ExecutionState{
			Behaviour: &ValidationSaver{},
			Current:   &CurrentExecution{},
		},
	}
	err = suite.lr.OnPulse(suite.ctx, pulse)
	suite.Require().NoError(err)
	suite.Equal(message.InPending, suite.lr.state[objectRef].ExecutionState.pending)
	qe := ExecutionQueueElement{}

	queue := append(make([]ExecutionQueueElement, 0), qe)

	suite.lr.state[objectRef] = &ObjectState{
		ExecutionState: &ExecutionState{
			Behaviour: &ValidationSaver{},
			Current:   &CurrentExecution{},
			Queue:     queue,
		},
	}

	err = suite.lr.OnPulse(suite.ctx, pulse)
	suite.Require().NoError(err)
	suite.Equal(message.InPending, suite.lr.state[objectRef].ExecutionState.pending)

	// Executor in new pulse is same node
	suite.jc.IsAuthorizedMock.Return(true, nil)
	suite.lr.state[objectRef].ExecutionState.pending = message.PendingUnknown

	suite.lr.state[objectRef] = &ObjectState{
		ExecutionState: &ExecutionState{
			Behaviour: &ValidationSaver{},
			Current:   &CurrentExecution{},
			Queue:     queue,
		},
	}

	err = suite.lr.OnPulse(suite.ctx, pulse)
	suite.Require().NoError(err)
	suite.Require().Equal(message.PendingUnknown, suite.lr.state[objectRef].ExecutionState.pending)

	suite.lr.state[objectRef].ExecutionState.pending = message.InPending

	err = suite.lr.OnPulse(suite.ctx, pulse)
	suite.Require().NoError(err)
	suite.Require().Equal(message.NotPending, suite.lr.state[objectRef].ExecutionState.pending)

	suite.jc.IsAuthorizedMock.Return(true, nil)
	suite.lr.state[objectRef].ExecutionState.pending = message.InPending

	err = suite.lr.OnPulse(suite.ctx, pulse)
	suite.Require().NoError(err)
	suite.Require().Equal(message.NotPending, suite.lr.state[objectRef].ExecutionState.pending)

	suite.lr.state[objectRef].ExecutionState.Current = nil
	suite.lr.state[objectRef].ExecutionState.pending = message.InPending
	suite.lr.state[objectRef].ExecutionState.PendingConfirmed = false

	err = suite.lr.OnPulse(suite.ctx, pulse)
	suite.Require().NoError(err)
	suite.Equal(message.NotPending, suite.lr.state[objectRef].ExecutionState.pending)

	suite.jc.IsAuthorizedMock.Return(false, nil)
	suite.lr.state[objectRef].ExecutionState.Current = nil
	suite.lr.state[objectRef].ExecutionState.pending = message.InPending
	suite.lr.state[objectRef].ExecutionState.PendingConfirmed = false

	err = suite.lr.OnPulse(suite.ctx, pulse)
	suite.Require().NoError(err)

	_, ok := suite.lr.state[objectRef]
	suite.Equal(false, ok)
}

func (suite *LogicRunnerTestSuite) TestPendingFinished() {
	pulse := core.Pulse{}
	objectRef := testutils.RandomRef()
	meRef := testutils.RandomRef()

	suite.jc.MeMock.Return(meRef)
	suite.ps.CurrentMock.Return(&pulse, nil)

	es := &ExecutionState{
		Behaviour: &ValidationSaver{},
		Current:   &CurrentExecution{},
		pending:   message.NotPending,
	}

	// make sure that if there is no pending finishPendingIfNeeded returns false,
	// doesn't send PendingFinished message and doesn't change ExecutionState.pending
	suite.lr.finishPendingIfNeeded(suite.ctx, es, objectRef)
	suite.Require().Zero(suite.mb.SendCounter)
	suite.Require().Equal(message.NotPending, es.pending)

	es.pending = message.InPending
	es.objectbody = &ObjectBody{}
	suite.mb.SendMock.ExpectOnce(suite.ctx, &message.PendingFinished{Reference: objectRef}, nil).Return(&reply.ID{}, nil)
	suite.jc.IsAuthorizedMock.Return(false, nil)
	suite.lr.finishPendingIfNeeded(suite.ctx, es, objectRef)
	suite.Require().Equal(message.NotPending, es.pending)
	suite.Require().Nil(es.objectbody)

	suite.mc.Wait(time.Second) // message bus' send is called in a goroutine

	es.pending = message.InPending
	es.objectbody = &ObjectBody{}
	suite.jc.IsAuthorizedMock.Return(true, nil)
	suite.lr.finishPendingIfNeeded(suite.ctx, es, objectRef)
	suite.Require().Equal(message.NotPending, es.pending)
	suite.Require().NotNil(es.objectbody)
}

func (suite *LogicRunnerTestSuite) TestStartQueueProcessorIfNeeded_DontStartQueueProcessorWhenPending() {
	objectRef := testutils.RandomRef()

	suite.am.HasPendingRequestsMock.Return(true, nil)
	es := &ExecutionState{ArtifactManager: suite.am, Queue: make([]ExecutionQueueElement, 0)}
	es.Queue = append(es.Queue, ExecutionQueueElement{})
	err := suite.lr.StartQueueProcessorIfNeeded(
		suite.ctx,
		es,
		&message.CallMethod{
			ObjectRef: objectRef,
			Method:    "some",
		},
	)
	suite.Require().NoError(err)
	suite.Require().Equal(message.InPending, es.pending)
}

func (suite *LogicRunnerTestSuite) TestCheckPendingRequests() {
	objectRef := testutils.RandomRef()

	es := &ExecutionState{ArtifactManager: suite.am}
	pending, err := es.CheckPendingRequests(
		suite.ctx, &message.CallConstructor{},
	)
	suite.Require().NoError(err)
	suite.Require().Equal(message.NotPending, pending)

	suite.am.HasPendingRequestsMock.Return(false, nil)
	es = &ExecutionState{ArtifactManager: suite.am}
	pending, err = es.CheckPendingRequests(
		suite.ctx, &message.CallMethod{
			ObjectRef: objectRef,
		},
	)
	suite.Require().NoError(err)
	suite.Require().Equal(message.NotPending, pending)

	suite.am.HasPendingRequestsMock.Return(true, nil)
	es = &ExecutionState{ArtifactManager: suite.am}
	pending, err = es.CheckPendingRequests(
		suite.ctx, &message.CallMethod{
			ObjectRef: objectRef,
		},
	)
	suite.Require().NoError(err)
	suite.Require().Equal(message.InPending, pending)

	suite.am.HasPendingRequestsMock.Return(false, errors.New("some"))
	es = &ExecutionState{ArtifactManager: suite.am}
	pending, err = es.CheckPendingRequests(
		suite.ctx, &message.CallMethod{
			ObjectRef: objectRef,
		},
	)
	suite.Require().Error(err)
	suite.Require().Equal(message.NotPending, pending)
}

func (suite *LogicRunnerTestSuite) TestPrepareState() {
	object := testutils.RandomRef()
	msg := &message.ExecutorResults{
		Caller:    testutils.RandomRef(),
		RecordRef: object,
	}

	// not pending
	// it's a first call, it's also initialize lr.state[object].ExecutionState
	// also check for empty Queue
	msg.Pending = message.NotPending
	_ = suite.lr.prepareObjectState(suite.ctx, msg)
	suite.Require().Equal(message.NotPending, suite.lr.state[object].ExecutionState.pending)
	suite.Require().Equal(0, len(suite.lr.state[object].ExecutionState.Queue))

	// pending without queue
	suite.lr.state[object].ExecutionState.pending = message.PendingUnknown
	msg.Pending = message.InPending
	_ = suite.lr.prepareObjectState(suite.ctx, msg)
	suite.Require().Equal(message.InPending, suite.lr.state[object].ExecutionState.pending)

	// do not change pending status if it isn't unknown
	suite.lr.state[object].ExecutionState.pending = message.NotPending
	msg.Pending = message.InPending
	_ = suite.lr.prepareObjectState(suite.ctx, msg)
	suite.Require().Equal(message.NotPending, suite.lr.state[object].ExecutionState.pending)

	// do not change pending status if it isn't unknown
	suite.lr.state[object].ExecutionState.pending = message.InPending
	msg.Pending = message.InPending
	_ = suite.lr.prepareObjectState(suite.ctx, msg)
	suite.Require().Equal(message.InPending, suite.lr.state[object].ExecutionState.pending)

	parcel := testutils.NewParcelMock(suite.mc)
	parcel.ContextMock.Expect(context.Background()).Return(context.Background())
	// brand new queue from message
	msg.Queue = []message.ExecutionQueueElement{
		message.ExecutionQueueElement{Parcel: parcel},
	}
	_ = suite.lr.prepareObjectState(suite.ctx, msg)
	suite.Require().Equal(1, len(suite.lr.state[object].ExecutionState.Queue))

	testMsg := message.CallMethod{ReturnMode: message.ReturnNoWait}
	parcel = testutils.NewParcelMock(suite.mc)
	parcel.ContextMock.Expect(context.Background()).Return(context.Background())
	parcel.MessageMock.Return(&testMsg) // mock message that returns NoWait

	queueElementRequest := testutils.RandomRef()
	msg.Queue = []message.ExecutionQueueElement{message.ExecutionQueueElement{Request: &queueElementRequest, Parcel: parcel}}
	_ = suite.lr.prepareObjectState(suite.ctx, msg)
	suite.Require().Equal(2, len(suite.lr.state[object].ExecutionState.Queue))
	suite.Require().Equal(&queueElementRequest, suite.lr.state[object].ExecutionState.Queue[0].request)
	suite.Require().Equal(&testMsg, suite.lr.state[object].ExecutionState.Queue[0].parcel.Message())
}

func (suite *LogicRunnerTestSuite) TestHandlePendingFinishedMessage() {
	objectRef := testutils.RandomRef()

	parcel := testutils.NewParcelMock(suite.mc).MessageMock.Return(
		&message.PendingFinished{Reference: objectRef},
	)

	parcel.DefaultTargetMock.Return(&core.RecordRef{})

	re, err := suite.lr.HandlePendingFinishedMessage(suite.ctx, parcel)
	suite.Require().NoError(err)
	suite.Require().Equal(&reply.OK{}, re)

	st := suite.lr.MustObjectState(objectRef)

	es := st.ExecutionState
	suite.Require().NotNil(es)
	suite.Require().Equal(message.NotPending, es.pending)

	es.Current = &CurrentExecution{}
	re, err = suite.lr.HandlePendingFinishedMessage(suite.ctx, parcel)
	suite.Require().Error(err)

	es.Current = nil

	re, err = suite.lr.HandlePendingFinishedMessage(suite.ctx, parcel)
	suite.Require().NoError(err)
	suite.Require().Equal(&reply.OK{}, re)
}

func (suite *LogicRunnerTestSuite) TestCheckExecutionLoop() {
	es := &ExecutionState{
		Current: nil,
	}

	loop := suite.lr.CheckExecutionLoop(suite.ctx, es, nil)
	suite.Require().False(loop)

	ctxA, _ := inslogger.WithTraceField(suite.ctx, "a")
	ctxB, _ := inslogger.WithTraceField(suite.ctx, "b")

	parcel := testutils.NewParcelMock(suite.mc).MessageMock.Return(
		&message.CallMethod{ReturnMode: message.ReturnResult},
	)
	es.Current = &CurrentExecution{
		ReturnMode: message.ReturnResult,
		Context:    ctxA,
	}

	loop = suite.lr.CheckExecutionLoop(ctxA, es, parcel)
	suite.Require().True(loop)

	loop = suite.lr.CheckExecutionLoop(ctxB, es, parcel)
	suite.Require().False(loop)

	parcel = testutils.NewParcelMock(suite.mc).MessageMock.Return(
		&message.CallMethod{ReturnMode: message.ReturnNoWait},
	)
	es.Current = &CurrentExecution{
		ReturnMode: message.ReturnResult,
		Context:    ctxA,
	}
	loop = suite.lr.CheckExecutionLoop(ctxA, es, parcel)
	suite.Require().False(loop)

	parcel = testutils.NewParcelMock(suite.mc)
	es.Current = &CurrentExecution{
		ReturnMode: message.ReturnNoWait,
		Context:    ctxA,
	}
	loop = suite.lr.CheckExecutionLoop(ctxA, es, parcel)
	suite.Require().False(loop)
}

func (suite *LogicRunnerTestSuite) TestHandleStillExecutingMessage() {
	objectRef := testutils.RandomRef()

	parcel := testutils.NewParcelMock(suite.mc).MessageMock.Return(
		&message.StillExecuting{Reference: objectRef},
	)

	parcel.DefaultTargetMock.Return(&core.RecordRef{})

	re, err := suite.lr.HandleStillExecutingMessage(suite.ctx, parcel)
	suite.Require().NoError(err)
	suite.Require().Equal(&reply.OK{}, re)

	st := suite.lr.MustObjectState(objectRef)
	suite.Require().NotNil(st.ExecutionState)
	suite.Require().Equal(message.InPending, st.ExecutionState.pending)
	suite.Require().Equal(true, st.ExecutionState.PendingConfirmed)

	st.ExecutionState.pending = message.NotPending
	st.ExecutionState.PendingConfirmed = false

	re, err = suite.lr.HandleStillExecutingMessage(suite.ctx, parcel)
	suite.Require().NoError(err)
	suite.Require().Equal(&reply.OK{}, re)

	st = suite.lr.MustObjectState(objectRef)
	suite.Require().NotNil(st.ExecutionState)
	suite.Require().Equal(message.NotPending, st.ExecutionState.pending)
	suite.Require().Equal(false, st.ExecutionState.PendingConfirmed)
}

func (suite *LogicRunnerTestSuite) Test_OnPulse_StillExecuting() {
	suite.jc.IsAuthorizedMock.Return(false, nil)
	suite.jc.MeMock.Return(core.RecordRef{})

	// test empty lr
	pulse := core.Pulse{}

	objectRef := testutils.RandomRef()

	suite.lr.state[objectRef] = &ObjectState{
		ExecutionState: &ExecutionState{
			Behaviour: &ValidationSaver{},
			Current:   &CurrentExecution{},
		},
	}
	suite.mb.SendMock.Return(&reply.OK{}, nil)
	err := suite.lr.OnPulse(suite.ctx, pulse)
	suite.Require().NoError(err)
	suite.NotNil(suite.lr.state[objectRef].ExecutionState)
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
		{maxQueueLength + 1, expected{maxQueueLength, true}},
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
		mbMock                        *testutils.MessageBusMock
		ExpectedLedgerHasMoreRequests bool
	}
	testCases := []testCase{
		{make([]ExecutionQueueElement, maxQueueLength+1), testutils.NewMessageBusMock(suite.mc), true},
		{make([]ExecutionQueueElement, maxQueueLength), testutils.NewMessageBusMock(suite.mc), false},
	}

	suite.jc.IsAuthorizedMock.Return(false, nil)
	suite.jc.MeMock.Return(core.RecordRef{})

	pulse := core.Pulse{}

	for _, test := range testCases {
		suite.SetupLogicRunner()
		queue := test.queue

		messagesQueue := convertQueueToMessageQueue(queue[:maxQueueLength])

		ref := testutils.RandomRef()

		suite.lr.JetCoordinator = suite.jc

		suite.lr.state[ref] = &ObjectState{
			ExecutionState: &ExecutionState{
				Behaviour: &ValidationSaver{},
				Queue:     queue,
			},
		}

		mb := test.mbMock
		suite.lr.MessageBus = mb

		expectedMessage := &message.ExecutorResults{
			RecordRef:             ref,
			Requests:              make([]message.CaseBindRequest, 0),
			Queue:                 messagesQueue,
			LedgerHasMoreRequests: test.ExpectedLedgerHasMoreRequests,
		}

		// defer new SendFunc before calling OnPulse
		mb.SendFunc = func(p context.Context, p1 core.Message, p2 *core.MessageSendOptions) (r core.Reply, r1 error) {
			suite.Equal(expectedMessage, p1)
			return nil, nil
		}

		err := suite.lr.OnPulse(suite.ctx, pulse)
		suite.NoError(err)
	}

	// waiting for all goroutines with Send() processing
	suite.mc.Wait(10 * time.Second)
	for _, test := range testCases {
		suite.Equal(1, int(test.mbMock.SendCounter))
	}
}

func (suite *LogicRunnerTestSuite) TestNoExcessiveAmends() {
	suite.am.UpdateObjectMock.Return(nil, nil)

	randRef := testutils.RandomRef()

	es := &ExecutionState{ArtifactManager: suite.am, Queue: make([]ExecutionQueueElement, 0)}
	es.Queue = append(es.Queue, ExecutionQueueElement{})
	es.objectbody = &ObjectBody{}
	es.objectbody.CodeMachineType = core.MachineTypeBuiltin
	es.Current = &CurrentExecution{}
	es.Current.LogicContext = &core.LogicCallContext{}
	es.Current.Request = &randRef
	es.objectbody.CodeRef = &randRef

	data := []byte(testutils.RandomString())
	es.objectbody.Object = data

	mle := testutils.NewMachineLogicExecutorMock(suite.mc)
	suite.lr.Executors[core.MachineTypeBuiltin] = mle
	mle.CallMethodMock.Return(data, nil, nil)

	msg := &message.CallMethod{
		ObjectRef: randRef,
		Method:    "some",
	}

	// In this case Update isn't send to ledger (objects data/newData are the same)
	suite.am.RegisterResultMock.Return(nil, nil)

	_, err := suite.lr.executeMethodCall(suite.ctx, es, msg)
	suite.Require().NoError(err)
	suite.Require().Equal(uint64(0), suite.am.UpdateObjectCounter)

	// In this case Update is send to ledger (objects data/newData are different)
	newData := make([]byte, 5, 5)
	mle.CallMethodMock.Return(newData, nil, nil)

	_, err = suite.lr.executeMethodCall(suite.ctx, es, msg)
	suite.Require().NoError(err)
	suite.Require().Equal(uint64(1), suite.am.UpdateObjectCounter)
}

func (suite *LogicRunnerTestSuite) TestHandleAbandonedRequestsNotificationMessage() {
	objectId := testutils.RandomID()
	msg := &message.AbandonedRequestsNotification{Object: objectId}
	parcel := &message.Parcel{Msg: msg}

	_, err := suite.lr.HandleAbandonedRequestsNotificationMessage(suite.ctx, parcel)
	suite.Require().NoError(err)
	suite.Equal(true, suite.lr.state[*msg.DefaultTarget()].ExecutionState.LedgerHasMoreRequests)

	// LedgerHasMoreRequests false
	suite.lr, _ = NewLogicRunner(&configuration.LogicRunner{})
	suite.lr.state[*msg.DefaultTarget()] = &ObjectState{ExecutionState: &ExecutionState{LedgerHasMoreRequests: false}}

	_, err = suite.lr.HandleAbandonedRequestsNotificationMessage(suite.ctx, parcel)
	suite.Require().NoError(err)
	suite.Equal(true, suite.lr.state[*msg.DefaultTarget()].ExecutionState.LedgerHasMoreRequests)

	// LedgerHasMoreRequests already true
	suite.lr, _ = NewLogicRunner(&configuration.LogicRunner{})
	suite.lr.state[*msg.DefaultTarget()] = &ObjectState{ExecutionState: &ExecutionState{LedgerHasMoreRequests: true}}

	_, err = suite.lr.HandleAbandonedRequestsNotificationMessage(suite.ctx, parcel)
	suite.Require().NoError(err)
	suite.Equal(true, suite.lr.state[*msg.DefaultTarget()].ExecutionState.LedgerHasMoreRequests)
}

func (suite *LogicRunnerTestSuite) TestPrepareObjectStateChangePendingStatus() {
	ref := testutils.RandomRef()

	msg := &message.ExecutorResults{RecordRef: ref}

	// we are in pending and come to ourselves again
	suite.lr.state[ref] = &ObjectState{ExecutionState: &ExecutionState{
		pending: message.InPending, Current: &CurrentExecution{}},
	}
	err := suite.lr.prepareObjectState(suite.ctx, msg)
	suite.Require().NoError(err)
	suite.Equal(message.NotPending, suite.lr.state[ref].ExecutionState.pending)
	suite.Equal(false, suite.lr.state[ref].ExecutionState.PendingConfirmed)

	// previous executor decline pending, trust him
	msg = &message.ExecutorResults{RecordRef: ref, Pending: message.NotPending}
	suite.lr.state[ref] = &ObjectState{ExecutionState: &ExecutionState{
		pending: message.InPending, Current: nil},
	}
	err = suite.lr.prepareObjectState(suite.ctx, msg)
	suite.Require().NoError(err)
	suite.Equal(message.NotPending, suite.lr.state[ref].ExecutionState.pending)
}

func (suite *LogicRunnerTestSuite) TestPrepareObjectStateChangeLedgerHasMoreRequests() {
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
		suite.lr.state[ref] = &ObjectState{ExecutionState: &ExecutionState{QueueProcessorActive: true, LedgerHasMoreRequests: test.objectStateStatus}}
		err := suite.lr.prepareObjectState(suite.ctx, msg)
		suite.Require().NoError(err)
		suite.Equal(test.expectedObjectStateStatue, suite.lr.state[ref].ExecutionState.LedgerHasMoreRequests)
	}
}

func TestLogicRunner(t *testing.T) {
	suite.Run(t, new(LogicRunnerTestSuite))
}

func TestGetLedgerPendingRequest(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	ctx := inslogger.TestContext(t)
	lr, _ := NewLogicRunner(&configuration.LogicRunner{})
	//ref := testutils.RandomRef()
	id := testutils.RandomID()

	const currentPulseNumber = core.PulseNumber(3)
	const oldRequestPulseNumber = core.PulseNumber(1)

	// no need to change
	es := &ExecutionState{LedgerQueueElement: &ExecutionQueueElement{pulse: currentPulseNumber}}
	lr.getLedgerPendingRequest(ctx, es, id)
	assert.Equal(t, es.LedgerQueueElement.pulse, currentPulseNumber)

	// no need to change
	es = &ExecutionState{LedgerHasMoreRequests: false}
	lr.getLedgerPendingRequest(ctx, es, id)
	assert.Nil(t, es.LedgerQueueElement)

	// there is no more requests
	es = &ExecutionState{LedgerHasMoreRequests: true}
	am := testutils.NewArtifactManagerMock(t)
	am.GetPendingRequestMock.Return(nil, core.ErrNoPendingRequest)
	lr.ArtifactManager = am
	lr.getLedgerPendingRequest(ctx, es, id)
	assert.Equal(t, false, es.LedgerHasMoreRequests)

	// there is requests
	es = &ExecutionState{LedgerHasMoreRequests: true}
	parcel := &message.Parcel{
		PulseNumber: oldRequestPulseNumber,
		Msg:         &message.CallMethod{},
	}
	am.GetPendingRequestMock.Return(parcel, nil)
	lr.ArtifactManager = am

	// preset some data, that doesn't need before
	ps := testutils.NewPulseStorageMock(mc)
	ps.CurrentMock.Return(&core.Pulse{PulseNumber: currentPulseNumber}, nil)
	lr.PulseStorage = ps

	// we doesn't authorized (pulse change in time we process function)
	jc := testutils.NewJetCoordinatorMock(mc)
	jc.IsAuthorizedMock.Return(false, nil)
	jc.MeMock.Return(core.RecordRef{})
	lr.JetCoordinator = jc
	lr.getLedgerPendingRequest(ctx, es, id)
	assert.Nil(t, es.LedgerQueueElement)

	// happy path
	jc = testutils.NewJetCoordinatorMock(mc)
	jc.IsAuthorizedMock.Return(true, nil)
	jc.MeMock.Return(core.RecordRef{})
	lr.JetCoordinator = jc
	lr.getLedgerPendingRequest(ctx, es, id)
	assert.Equal(t, true, es.LedgerHasMoreRequests)
	assert.Equal(t, parcel, es.LedgerQueueElement.parcel)
	assert.Equal(t, currentPulseNumber, es.LedgerQueueElement.pulse)
}
