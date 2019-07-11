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
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	message2 "github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/infrastructure/gochannel"
	"github.com/gojuno/minimock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/pulsar"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/network"
)

type LogicRunnerCommonTestSuite struct {
	suite.Suite

	mc  *minimock.Controller
	ctx context.Context
	am  *artifacts.ClientMock
	dc  *artifacts.DescriptorsCacheMock
	mb  *testutils.MessageBusMock
	jc  *jet.CoordinatorMock
	mm  *mmanager
	lr  *LogicRunner
	re  *RequestsExecutorMock
	es  ExecutionState
	ps  *pulse.AccessorMock
	mle *testutils.MachineLogicExecutorMock
	nn  *network.NodeNetworkMock
}

func (suite *LogicRunnerCommonTestSuite) BeforeTest(suiteName, testName string) {
	// testing context
	suite.ctx = inslogger.TestContext(suite.T())

	// initialize minimock and mocks
	suite.mc = minimock.NewController(suite.T())
	suite.am = artifacts.NewClientMock(suite.mc)
	suite.dc = artifacts.NewDescriptorsCacheMock(suite.mc)
	suite.mm = &mmanager{}
	suite.re = NewRequestsExecutorMock(suite.mc)
	suite.mb = testutils.NewMessageBusMock(suite.mc)
	suite.jc = jet.NewCoordinatorMock(suite.mc)
	suite.ps = pulse.NewAccessorMock(suite.mc)
	suite.nn = network.NewNodeNetworkMock(suite.mc)

	suite.SetupLogicRunner()
}

func (suite *LogicRunnerCommonTestSuite) SetupLogicRunner() {
	suite.lr, _ = NewLogicRunner(&configuration.LogicRunner{})
	suite.lr.ArtifactManager = suite.am
	suite.lr.DescriptorsCache = suite.dc
	suite.lr.MessageBus = suite.mb
	suite.lr.MachinesManager = suite.mm
	suite.lr.JetCoordinator = suite.jc
	suite.lr.PulseAccessor = suite.ps
	suite.lr.NodeNetwork = suite.nn
	suite.lr.RequestsExecutor = suite.re

	_ = suite.lr.Init(suite.ctx)
}

func (suite *LogicRunnerCommonTestSuite) AfterTest(suiteName, testName string) {
	suite.mc.Wait(2 * time.Minute)
	suite.mc.Finish()

	// LogicRunner created a number of goroutines (in watermill, for example)
	// that weren't shut down in case no Stop was called
	// Do what we must, stop server
	_ = suite.lr.Stop(suite.ctx)
}

type LogicRunnerTestSuite struct {
	LogicRunnerCommonTestSuite
}

func (suite *LogicRunnerTestSuite) BeforeTest(suiteName, testName string) {
	suite.LogicRunnerCommonTestSuite.BeforeTest(suiteName, testName)
}

func (suite *LogicRunnerTestSuite) SetupLogicRunner() {
	suite.LogicRunnerCommonTestSuite.SetupLogicRunner()
}

func (suite *LogicRunnerTestSuite) AfterTest(suiteName, testName string) {
	suite.LogicRunnerCommonTestSuite.AfterTest(suiteName, testName)
}

func (suite *LogicRunnerTestSuite) TestPendingFinished() {
	pulseNum := insolar.Pulse{}
	objectRef := testutils.RandomRef()
	meRef := testutils.RandomRef()

	suite.jc.MeMock.Return(meRef)
	suite.ps.LatestFunc = func(p context.Context) (r insolar.Pulse, r1 error) {
		return pulseNum, nil
	}

	es, broker := suite.lr.StateStorage.UpsertExecutionState(objectRef)
	broker.currentList.Set(objectRef, &Transcript{})
	es.pending = message.NotPending

	// make sure that if there is no pending finishPendingIfNeeded returns false,
	// doesn't send PendingFinished message and doesn't change ExecutionState.pending
	broker.finishPendingIfNeeded(suite.ctx)
	suite.Require().Zero(suite.mb.SendCounter)
	suite.Require().Equal(message.NotPending, es.pending)

	es.pending = message.InPending
	suite.mb.SendMock.ExpectOnce(suite.ctx, &message.PendingFinished{Reference: objectRef}, nil).Return(&reply.ID{}, nil)
	suite.jc.IsAuthorizedMock.Return(false, nil)
	broker.finishPendingIfNeeded(suite.ctx)
	suite.Require().Equal(message.NotPending, es.pending)

	suite.mc.Wait(time.Minute) // message bus' send is called in a goroutine

	es.pending = message.InPending
	suite.jc.IsAuthorizedMock.Return(true, nil)
	broker.finishPendingIfNeeded(suite.ctx)
	suite.Require().Equal(message.NotPending, es.pending)
}

func (suite *LogicRunnerTestSuite) TestHandleAdditionalCallFromPreviousExecutor() {
	table := []struct {
		name                           string
		clarifyPendingStateResult      error
		startQueueProcessorResult      error
		expectedClarifyPendingStateCtr int32
		expectedStartQueueProcessorCtr int32
	}{
		{
			name:                           "Happy path",
			expectedClarifyPendingStateCtr: 1,
			expectedStartQueueProcessorCtr: 1,
		},
		{
			name:                           "ClarifyPendingState failed",
			clarifyPendingStateResult:      fmt.Errorf("ClarifyPendingState failed"),
			expectedClarifyPendingStateCtr: 1,
		},
		{
			name:                           "StartQueueProcessorIfNeeded failed",
			startQueueProcessorResult:      fmt.Errorf("StartQueueProcessorIfNeeded failed"),
			expectedClarifyPendingStateCtr: 1,
			expectedStartQueueProcessorCtr: 1,
		},
		{
			name:                           "Both procedures fail",
			clarifyPendingStateResult:      fmt.Errorf("ClarifyPendingState failed"),
			startQueueProcessorResult:      fmt.Errorf("StartQueueProcessorIfNeeded failed"),
			expectedClarifyPendingStateCtr: 1,
			expectedStartQueueProcessorCtr: 0,
		},
	}

	for _, test := range table {
		test := test
		suite.T().Run(test.name, func(t *testing.T) {
			h := HandleAdditionalCallFromPreviousExecutor{
				dep: &Dependencies{
					lr: suite.lr,
				},
			}
			f := flow.NewFlowMock(suite.T())
			reqRef := gen.Reference()
			msg := message.AdditionalCallFromPreviousExecutor{
				ObjectReference: gen.Reference(),
				RequestRef:      reqRef,
				Request:         record.IncomingRequest{},
			}

			var clarifyPendingStateCtr int32
			f.ProcedureFunc = func(ctx context.Context, proc flow.Procedure, cancelable bool) error {
				atomic.AddInt32(&clarifyPendingStateCtr, 1)
				_, ok := proc.(*ClarifyPendingState)
				require.True(suite.T(), ok)
				return test.clarifyPendingStateResult
			}

			var startQueueProcessorCtr int32
			f.HandleFunc = func(ctx context.Context, handle flow.Handle) error {
				atomic.AddInt32(&startQueueProcessorCtr, 1)
				return test.startQueueProcessorResult
			}

			h.handleActual(suite.ctx, &msg, f)

			assert.Equal(suite.T(), test.expectedClarifyPendingStateCtr, atomic.LoadInt32(&clarifyPendingStateCtr))
		})
	}
}

func (suite *LogicRunnerTestSuite) TestCheckPendingRequests() {
	objectRef := testutils.RandomRef()

	table := []struct {
		name     string
		inState  message.PendingState
		outState message.PendingState
		request  bool
		callType record.CallType
		amReply  *struct {
			has bool
			err error
		}
		isError bool
	}{
		{
			name:     "already in pending",
			inState:  message.InPending,
			outState: message.InPending,
		},
		{
			name:     "already not in pending",
			inState:  message.NotPending,
			outState: message.NotPending,
		},
		{
			name:     "constructor call",
			inState:  message.PendingUnknown,
			request:  true,
			callType: record.CTSaveAsChild,
			outState: message.NotPending,
		},
		{
			name:    "method call, not pending",
			inState: message.PendingUnknown,
			request: true,
			amReply: &struct {
				has bool
				err error
			}{false, nil},
			outState: message.NotPending,
		},
		{
			name:    "method call, in pending",
			inState: message.PendingUnknown,
			request: true,
			amReply: &struct {
				has bool
				err error
			}{true, nil},
			outState: message.InPending,
		},
		{
			name:    "method call, in pending",
			inState: message.PendingUnknown,
			request: true,
			amReply: &struct {
				has bool
				err error
			}{true, errors.New("some")},
			outState: message.PendingUnknown,
			isError:  true,
		},
	}

	for _, test := range table {
		suite.T().Run(test.name, func(t *testing.T) {
			var request *record.IncomingRequest
			if test.request {
				request = &record.IncomingRequest{CallType: test.callType}
			}
			es := NewExecutionState(objectRef)
			es.pending = test.inState
			if test.amReply != nil {
				suite.am.HasPendingRequestsMock.Return(test.amReply.has, test.amReply.err)
			}
			proc := ClarifyPendingState{
				es:              es,
				request:         request,
				ArtifactManager: suite.lr.ArtifactManager,
			}
			err := proc.Proceed(suite.ctx)
			if test.isError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, test.outState, es.pending)
		})
	}

	suite.T().Run("method call, AM error", func(t *testing.T) {
		request := record.IncomingRequest{CallType: record.CTMethod}

		es := NewExecutionState(objectRef)
		es.pending = message.PendingUnknown

		suite.am.HasPendingRequestsMock.Return(false, errors.New("some"))

		proc := ClarifyPendingState{
			es:              es,
			request:         &request,
			ArtifactManager: suite.lr.ArtifactManager,
		}
		err := proc.Proceed(suite.ctx)
		require.Error(t, err)
		require.Equal(t, message.PendingUnknown, es.pending)
	})
}

func prepareParcel(t minimock.Tester, msg insolar.Message, needType bool, needSender bool) insolar.Parcel {
	parcel := testutils.NewParcelMock(t)
	parcel.MessageMock.Return(msg)
	if needType {
		parcel.TypeMock.Return(msg.Type())
	}
	if needSender {
		parcel.GetSenderMock.Return(gen.Reference())
	}
	return parcel
}

func prepareWatermill(suite *LogicRunnerTestSuite) (flow.Flow, message2.PubSub) {
	flowMock := flow.NewFlowMock(suite.mc)
	flowMock.ProcedureMock.Set(func(p context.Context, p1 flow.Procedure, p2 bool) (r error) {
		return p1.Proceed(p)
	})

	wmLogger := log.NewWatermillLogAdapter(inslogger.FromContext(suite.ctx))
	pubSub := gochannel.NewGoChannel(gochannel.Config{}, wmLogger)

	return flowMock, pubSub
}

func (suite *LogicRunnerTestSuite) TestPrepareState() {
	type msgt struct {
		pending  message.PendingState
		queueLen int
	}
	type exp struct {
		pending        message.PendingState
		queueLen       int
		hasPendingCall bool
	}
	type obj struct {
		pending  message.PendingState
		queueLen int
	}
	table := []struct {
		name           string
		existingObject bool
		object         obj
		message        msgt
		expected       exp
	}{
		{
			name:     "first call, NotPending in message",
			message:  msgt{pending: message.NotPending},
			expected: exp{pending: message.NotPending},
		},
		{
			name:     "message says InPending, no object",
			message:  msgt{pending: message.InPending},
			expected: exp{pending: message.InPending},
		},
		{
			name:           "message says InPending, with object",
			existingObject: true,
			message:        msgt{pending: message.InPending},
			expected:       exp{pending: message.InPending},
		},
		{
			name:           "do not change pending status if existing says NotPending",
			existingObject: true,
			object:         obj{pending: message.NotPending},
			message:        msgt{pending: message.InPending},
			expected:       exp{pending: message.NotPending},
		},
		{
			name:           "message changes to NotPending, prev executor forces",
			existingObject: true,
			object:         obj{pending: message.InPending},
			message:        msgt{pending: message.NotPending},
			expected:       exp{pending: message.NotPending},
		},
		{
			name: "message has queue, no existing object",
			message: msgt{
				pending:  message.InPending,
				queueLen: 1,
			},
			expected: exp{
				pending:  message.InPending,
				queueLen: 1,
			},
		},
		{
			name:           "message has queue and object has queue",
			existingObject: true,
			object: obj{
				pending:  message.InPending,
				queueLen: 1,
			},
			message: msgt{
				pending:  message.InPending,
				queueLen: 1,
			},
			expected: exp{
				pending:  message.InPending,
				queueLen: 2,
			},
		},
		{
			name: "message has queue, but unknown pending state",
			message: msgt{
				pending:  message.PendingUnknown,
				queueLen: 1,
			},
			expected: exp{
				pending:        message.InPending,
				queueLen:       1,
				hasPendingCall: true,
			},
		},
	}

	for _, test := range table {
		test := test
		suite.T().Run(test.name, func(t *testing.T) {
			pulseObj := insolar.Pulse{}
			pulseObj.PulseNumber = insolar.FirstPulseNumber

			object := testutils.RandomRef()
			defer delete(*suite.lr.StateStorage.StateMap(), object)

			msg := &message.ExecutorResults{
				Caller:    testutils.RandomRef(),
				RecordRef: object,
				Pending:   test.message.pending,
				Queue:     []message.ExecutionQueueElement{},
			}

			for test.message.queueLen > 0 {
				test.message.queueLen--

				reqRef := gen.Reference()
				msg.Queue = append(
					msg.Queue,
					message.ExecutionQueueElement{RequestRef: reqRef, Request: record.IncomingRequest{}},
				)
			}

			if test.existingObject {
				es, broker := suite.lr.StateStorage.UpsertExecutionState(object)
				es.pending = test.object.pending

				for test.object.queueLen > 0 {
					test.object.queueLen--

					reqRef := gen.Reference()
					broker.mutable.Push(&Transcript{RequestRef: reqRef})
				}
			}

			if test.expected.hasPendingCall {
				suite.am.HasPendingRequestsMock.Return(true, nil)
			}

			flowMock, pubSub := prepareWatermill(suite)
			fakeParcel := prepareParcel(suite.mc, msg, false, false)

			h := HandleExecutorResults{
				dep:     &Dependencies{Publisher: pubSub, lr: suite.lr},
				Message: bus.Message{Parcel: fakeParcel, ReplyTo: make(chan bus.Reply)},
			}
			err := h.realHandleExecutorState(suite.ctx, flowMock)
			suite.mc.Wait(time.Minute)
			suite.Require().NoError(err)

			es, broker := suite.lr.StateStorage.UpsertExecutionState(object)
			suite.Require().Equal(test.expected.pending, es.pending)
			suite.Require().Equal(test.expected.queueLen, broker.mutable.Length())
		})
	}
}

func (suite *LogicRunnerTestSuite) TestHandlePendingFinishedMessage() {
	objectRef := testutils.RandomRef()
	p := insolar.Pulse{PulseNumber: 100}

	parcel := testutils.NewParcelMock(suite.mc).MessageMock.Return(
		&message.PendingFinished{Reference: objectRef},
	)

	parcel.DefaultTargetMock.Return(&insolar.Reference{})
	parcel.PulseFunc = func() insolar.PulseNumber { return p.PulseNumber }

	re, err := suite.lr.FlowDispatcher.WrapBusHandle(suite.ctx, parcel)
	suite.Require().NoError(err)
	suite.Require().Equal(&reply.OK{}, re)

	es, broker := suite.lr.StateStorage.GetExecutionState(objectRef)
	suite.Require().NotNil(es)
	suite.Require().Equal(message.NotPending, es.pending)

	broker.currentList.Set(objectRef, &Transcript{})
	re, err = suite.lr.FlowDispatcher.WrapBusHandle(suite.ctx, parcel)
	suite.Require().Error(err)

	broker.currentList.Cleanup()

	re, err = suite.lr.FlowDispatcher.WrapBusHandle(suite.ctx, parcel)
	suite.Require().NoError(err)
	suite.Require().Equal(&reply.OK{}, re)
}

// TODO: move test to executionBroker package
func (suite *LogicRunnerTestSuite) TestCheckExecutionLoop() {
	objectRef := gen.Reference()
	_, broker := suite.lr.StateStorage.UpsertExecutionState(objectRef)

	reqIdA := utils.RandTraceID()
	reqIdB := utils.RandTraceID()
	request := record.IncomingRequest{
		ReturnMode:   record.ReturnResult,
		Object:       &objectRef,
		APIRequestID: reqIdA,
	}

	executingReqRef := gen.Reference()
	broker.currentList.Set(executingReqRef, &Transcript{
		Request: &record.IncomingRequest{ReturnMode: record.ReturnResult, APIRequestID: reqIdA},
	})
	loop := suite.lr.CheckExecutionLoop(suite.ctx, request)
	suite.Require().True(loop)

	broker.currentList.Set(executingReqRef, &Transcript{
		Request: &record.IncomingRequest{ReturnMode: record.ReturnResult, APIRequestID: reqIdB},
	})
	loop = suite.lr.CheckExecutionLoop(suite.ctx, request)

	suite.Require().False(loop)

	// intermediate env cleanup
	broker.currentList.Cleanup()

	request = record.IncomingRequest{
		ReturnMode: record.ReturnNoWait,
		Object:     &objectRef,
	}
	broker.currentList.Set(executingReqRef, &Transcript{
		Request: &record.IncomingRequest{ReturnMode: record.ReturnResult},
	})
	loop = suite.lr.CheckExecutionLoop(suite.ctx, request)

	suite.Require().False(loop)
	broker.currentList.Cleanup()

	broker.currentList.Set(executingReqRef, &Transcript{
		Request: &record.IncomingRequest{ReturnMode: record.ReturnNoWait},
	})
	loop = suite.lr.CheckExecutionLoop(suite.ctx, request)

	suite.Require().False(loop)
	broker.currentList.Cleanup()

	broker.currentList.Set(executingReqRef, &Transcript{
		Request: &record.IncomingRequest{ReturnMode: record.ReturnNoWait},
	})
	loop = suite.lr.CheckExecutionLoop(suite.ctx, request)

	suite.Require().False(loop)
}

func (suite *LogicRunnerTestSuite) TestHandleStillExecutingMessage() {
	objectRef := testutils.RandomRef()

	parcel := testutils.NewParcelMock(suite.mc).MessageMock.Return(
		&message.StillExecuting{Reference: objectRef},
	)

	p := insolar.Pulse{PulseNumber: 100}
	parcel.PulseFunc = func() insolar.PulseNumber { return p.PulseNumber }

	// check that creation of new execution state is handled (on StillExecuting Message)
	re, err := suite.lr.FlowDispatcher.WrapBusHandle(suite.ctx, parcel)
	suite.Require().NoError(err)
	suite.Require().Equal(&reply.OK{}, re)

	es, _ := suite.lr.StateStorage.GetExecutionState(objectRef)
	suite.Require().NotNil(es)
	suite.Require().Equal(message.InPending, es.pending)
	suite.Require().Equal(true, es.PendingConfirmed)

	es.pending = message.NotPending
	es.PendingConfirmed = false

	re, err = suite.lr.FlowDispatcher.WrapBusHandle(suite.ctx, parcel)
	suite.Require().NoError(err)
	suite.Require().Equal(&reply.OK{}, re)

	es, _ = suite.lr.StateStorage.GetExecutionState(objectRef)
	suite.Require().NotNil(es)
	suite.Require().Equal(message.NotPending, es.pending)
	suite.Require().Equal(false, es.PendingConfirmed)

	// If we already have task in InPending, but it wasn't confirmed
	suite.lr.StateStorage.DeleteObjectState(objectRef)

	es, _ = suite.lr.StateStorage.UpsertExecutionState(objectRef)
	es.pending = message.InPending
	es.PendingConfirmed = false

	re, err = suite.lr.FlowDispatcher.WrapBusHandle(suite.ctx, parcel)
	suite.Require().NoError(err)

	es, _ = suite.lr.StateStorage.GetExecutionState(objectRef)
	suite.Equal(message.InPending, es.pending)
	suite.Equal(true, es.PendingConfirmed)
}

// TODO: move this test to EB tests (TestRotate)
func (suite *LogicRunnerTestSuite) TestReleaseQueue() {
	tests := map[string]struct {
		QueueLength     int
		ExpectedLength  int
		ExpectedHasMore bool
	}{
		"zero":  {0, 0, false},
		"one":   {1, 1, false},
		"max":   {maxQueueLength, maxQueueLength, false},
		"max+1": {maxQueueLength + 1, maxQueueLength, true},
	}
	for name, tc := range tests {
		suite.T().Run(name, func(t *testing.T) {
			a := assert.New(t)

			objectRef := gen.Reference()
			_, broker := suite.lr.StateStorage.UpsertExecutionState(objectRef)
			InitBroker(t, suite.ctx, tc.QueueLength, broker, false)

			rotationResults := broker.Rotate(maxQueueLength)
			a.Equal(tc.ExpectedLength, len(rotationResults.Requests))
			a.Equal(tc.ExpectedHasMore, rotationResults.LedgerHasMoreRequests)
		})
	}
}

func (suite *LogicRunnerTestSuite) TestHandleAbandonedRequestsNotificationMessage() {
	suite.T().Skip("we disabled handling of this notification for now")

	objectId := testutils.RandomID()
	objectRef := *insolar.NewReference(objectId)
	msg := &message.AbandonedRequestsNotification{Object: objectId}
	parcel := &message.Parcel{Msg: msg}

	_, err := suite.lr.HandleAbandonedRequestsNotificationMessage(suite.ctx, parcel)
	suite.Require().NoError(err)
	es, _ := suite.lr.StateStorage.GetExecutionState(objectRef)
	suite.Equal(true, es.LedgerHasMoreRequests)
	_ = suite.lr.Stop(suite.ctx)

	// LedgerHasMoreRequests false
	suite.lr, _ = NewLogicRunner(&configuration.LogicRunner{})

	es, _ = suite.lr.StateStorage.UpsertExecutionState(objectRef)
	es.LedgerHasMoreRequests = false

	_, err = suite.lr.HandleAbandonedRequestsNotificationMessage(suite.ctx, parcel)
	suite.Require().NoError(err)
	es, _ = suite.lr.StateStorage.GetExecutionState(objectRef)
	suite.Equal(true, es.LedgerHasMoreRequests)
	_ = suite.lr.Stop(suite.ctx)

	// LedgerHasMoreRequests already true
	suite.lr, _ = NewLogicRunner(&configuration.LogicRunner{})

	es, _ = suite.lr.StateStorage.UpsertExecutionState(objectRef)
	es.LedgerHasMoreRequests = true

	_, err = suite.lr.HandleAbandonedRequestsNotificationMessage(suite.ctx, parcel)
	suite.Require().NoError(err)
	es, _ = suite.lr.StateStorage.GetExecutionState(objectRef)
	suite.Equal(true, es.LedgerHasMoreRequests)
	_ = suite.lr.Stop(suite.ctx)
}

func (suite *LogicRunnerTestSuite) TestPrepareObjectStateChangePendingStatus() {
	ref1 := testutils.RandomRef()

	flowMock, pubSub := prepareWatermill(suite)
	var fakeParcel insolar.Parcel
	var h HandleExecutorResults
	var err error

	msg := &message.ExecutorResults{RecordRef: ref1}
	fakeParcel = prepareParcel(suite.mc, msg, false, false)
	h = HandleExecutorResults{
		dep:     &Dependencies{Publisher: pubSub, lr: suite.lr},
		Message: bus.Message{Parcel: fakeParcel, ReplyTo: make(chan bus.Reply)},
	}

	es, broker := suite.lr.StateStorage.UpsertExecutionState(ref1)
	// st := suite.lr.StateStorage.UpsertObjectState(ref1)
	// es, broker := st.UpsertExecution(suite.lr, &ref1)
	broker.currentList.Set(ref1, &Transcript{})
	es.pending = message.InPending

	// we are in pending and come to ourselves again
	err = h.realHandleExecutorState(suite.ctx, flowMock)
	suite.Require().NoError(err)
	suite.Equal(message.NotPending, es.pending)
	suite.Equal(false, es.PendingConfirmed)

	ref2 := testutils.RandomRef()
	// previous executor decline pending, trust him
	msg = &message.ExecutorResults{RecordRef: ref2, Pending: message.NotPending}
	fakeParcel = prepareParcel(suite.mc, msg, false, false)
	h = HandleExecutorResults{
		dep:     &Dependencies{Publisher: pubSub, lr: suite.lr},
		Message: bus.Message{Parcel: fakeParcel, ReplyTo: make(chan bus.Reply)},
	}

	es, _ = suite.lr.StateStorage.UpsertExecutionState(ref2)
	es.pending = message.InPending

	err = h.realHandleExecutorState(suite.ctx, flowMock)
	suite.Require().NoError(err)
	suite.Equal(message.NotPending, es.pending)
}

func (suite *LogicRunnerTestSuite) TestPrepareObjectStateChangeLedgerHasMoreRequests() {
	ref := testutils.RandomRef()

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
		msg := &message.ExecutorResults{
			RecordRef:             ref,
			LedgerHasMoreRequests: test.messageStatus,
			Pending:               message.NotPending,
		}

		flowMock, pubSub := prepareWatermill(suite)
		fakeParcel := prepareParcel(suite.mc, msg, false, false)

		h := HandleExecutorResults{
			dep:     &Dependencies{Publisher: pubSub, lr: suite.lr},
			Message: bus.Message{Parcel: fakeParcel, ReplyTo: make(chan bus.Reply)},
		}

		es, broker := suite.lr.StateStorage.UpsertExecutionState(ref)
		broker.tryTakeProcessor(suite.ctx)
		es.LedgerHasMoreRequests = test.objectStateStatus

		err := h.realHandleExecutorState(suite.ctx, flowMock)
		suite.Require().NoError(err)
		es, _ = suite.lr.StateStorage.GetExecutionState(ref)
		suite.Equal(test.expectedObjectStateStatue, es.LedgerHasMoreRequests)
	}
}

func (suite *LogicRunnerTestSuite) TestNewLogicRunner() {
	lr, err := NewLogicRunner(nil)
	suite.Require().Error(err)
	suite.Require().Nil(lr)

	lr, err = NewLogicRunner(&configuration.LogicRunner{})
	suite.Require().NoError(err)
	suite.Require().NotNil(lr)
	_ = lr.Stop(context.Background())
}

func (suite *LogicRunnerTestSuite) TestStartStop() {
	lr, err := NewLogicRunner(&configuration.LogicRunner{
		BuiltIn: &configuration.BuiltIn{},
	})
	suite.Require().NoError(err)
	suite.Require().NotNil(lr)

	suite.mb.MustRegisterMock.Return()
	lr.MessageBus = suite.mb

	lr.MachinesManager = suite.mm

	suite.am.InjectCodeDescriptorMock.Return()
	suite.am.InjectObjectDescriptorMock.Return()
	suite.am.InjectFinishMock.Return()
	lr.ArtifactManager = suite.am

	err = lr.Start(suite.ctx)
	suite.Require().NoError(err)

	err = lr.Stop(suite.ctx)
	suite.Require().NoError(err)
}

func WaitGroup_TimeoutWait(wg *sync.WaitGroup, timeout time.Duration) bool {
	waitChannel := make(chan struct{}, 0)

	go func() {
		wg.Wait()
		waitChannel <- struct{}{}
	}()

	select {
	case <-waitChannel:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (suite *LogicRunnerTestSuite) TestConcurrency() {
	objectRef := testutils.RandomRef()
	parentRef := testutils.RandomRef()
	protoRef := testutils.RandomRef()
	codeRef := testutils.RandomRef()

	meRef := testutils.RandomRef()
	notMeRef := testutils.RandomRef()
	suite.jc.MeMock.Return(meRef)

	pulseNum := insolar.Pulse{PulseNumber: 100}

	suite.jc.IsAuthorizedFunc = func(
		ctx context.Context, role insolar.DynamicRole, id insolar.ID, pn insolar.PulseNumber, obj insolar.Reference,
	) (bool, error) {
		return true, nil
	}

	nodeMock := network.NewNetworkNodeMock(suite.T())
	nodeMock.IDMock.Return(meRef)

	od := artifacts.NewObjectDescriptorMock(suite.T())
	od.PrototypeMock.Return(&protoRef, nil)
	od.MemoryMock.Return([]byte{1, 2, 3})
	od.ParentMock.Return(&parentRef)
	od.HeadRefMock.Return(&objectRef)

	pd := artifacts.NewObjectDescriptorMock(suite.T())
	pd.CodeMock.Return(&codeRef, nil)
	pd.HeadRefMock.Return(&protoRef)

	cd := artifacts.NewCodeDescriptorMock(suite.T())
	cd.MachineTypeMock.Return(insolar.MachineTypeBuiltin)
	cd.RefMock.Return(&codeRef)

	suite.am.HasPendingRequestsMock.Return(false, nil)

	suite.am.RegisterIncomingRequestMock.Set(func(ctx context.Context, r *record.IncomingRequest) (*insolar.ID, error) {
		reqId := testutils.RandomID()
		return &reqId, nil
	})

	suite.re.ExecuteAndSaveMock.Return(nil, nil)
	suite.re.SendReplyMock.Return()

	num := 100
	wg := sync.WaitGroup{}
	wg.Add(num)

	for i := 0; i < num; i++ {
		go func(i int) {
			msg := &message.CallMethod{
				IncomingRequest: record.IncomingRequest{
					Prototype:    &protoRef,
					Object:       &objectRef,
					Method:       "some",
					APIRequestID: utils.RandTraceID(),
				},
			}

			parcel := testutils.NewParcelMock(suite.T())
			parcel.DefaultTargetMock.Return(&objectRef)
			parcel.MessageMock.Return(msg)
			parcel.TypeMock.Return(msg.Type())
			parcel.PulseMock.Return(pulseNum.PulseNumber)
			parcel.GetSenderMock.Return(notMeRef)

			ctx := inslogger.ContextWithTrace(suite.ctx, "req-"+strconv.Itoa(i))

			_, err := suite.lr.FlowDispatcher.WrapBusHandle(ctx, parcel)
			suite.Require().NoError(err)

			wg.Done()
		}(i)
	}

	suite.Require().True(WaitGroup_TimeoutWait(&wg, 2*time.Minute),
		"Failed to wait for all requests to be processed")
}

func (suite *LogicRunnerTestSuite) TestCallMethodWithOnPulse() {
	objectRef := testutils.RandomRef()
	protoRef := testutils.RandomRef()

	meRef := testutils.RandomRef()
	notMeRef := testutils.RandomRef()
	suite.jc.MeMock.Return(meRef)

	// If you think you are smart enough to make this test 'more effective'
	// by using atomic variables or goroutines or anything else, you are wrong.
	// Last time we spent two full workdays trying to find a race condition
	// in our code before we realized this test has a logic error related
	// to it concurrent nature. Keep the code as simple as possible. Don't be smart.
	var pn insolar.PulseNumber = 100
	var lck sync.Mutex

	suite.ps.LatestFunc = func(ctx context.Context) (insolar.Pulse, error) {
		lck.Lock()
		defer lck.Unlock()
		return insolar.Pulse{PulseNumber: pn}, nil
	}

	type whenType int
	const (
		whenIsAuthorized whenType = iota
		whenRegisterRequest
		whenHasPendingRequest
		whenCallMethod
	)

	table := []struct {
		name                      string
		when                      whenType
		messagesExpected          []insolar.MessageType
		errorExpected             bool
		flowCanceledExpected      bool
		pendingInExecutorResults  message.PendingState
		queueLenInExecutorResults int
	}{
		{
			name:                 "pulse change in IsAuthorized",
			when:                 whenIsAuthorized,
			flowCanceledExpected: true,
		},
		{
			name:                 "pulse change in RegisterIncomingRequest",
			when:                 whenRegisterRequest,
			flowCanceledExpected: true,
		},
		{
			name: "pulse change in HasPendingRequests",
			when: whenHasPendingRequest,
			messagesExpected: []insolar.MessageType{
				insolar.TypeAdditionalCallFromPreviousExecutor, insolar.TypeExecutorResults,
			},
			pendingInExecutorResults:  message.PendingUnknown,
			queueLenInExecutorResults: 1,
		},
		{
			name: "pulse change in CallMethod",
			when: whenCallMethod,
			messagesExpected: []insolar.MessageType{
				insolar.TypeExecutorResults, insolar.TypePendingFinished, insolar.TypeStillExecuting,
			},
			pendingInExecutorResults:  message.InPending,
			queueLenInExecutorResults: 0,
		},
	}

	for _, test := range table {
		test := test
		suite.T().Run(test.name, func(t *testing.T) {
			lck.Lock()
			pn = 100
			lck.Unlock()

			changePulse := func() {
				lck.Lock()
				defer lck.Unlock()
				pn += 1

				pulseNum := insolar.Pulse{PulseNumber: pn}
				ctx := inslogger.ContextWithTrace(suite.ctx, "pulse-"+strconv.Itoa(int(pn)))
				err := suite.lr.OnPulse(ctx, pulseNum)
				require.NoError(t, err)
				return
			}

			suite.jc.IsAuthorizedFunc = func(
				ctx context.Context, role insolar.DynamicRole, id insolar.ID, pnArg insolar.PulseNumber, obj insolar.Reference,
			) (bool, error) {
				if pnArg == 101 {
					return false, nil
				}

				if test.when == whenIsAuthorized {
					// Please note that changePulse calls LogicRunner.ChangePulse which calls IsAuthorized.
					// In other words this procedure is not called sequentially!
					changePulse()
				}

				lck.Lock()
				defer lck.Unlock()

				return pn == 100, nil
			}

			if test.when > whenIsAuthorized {
				suite.am.RegisterIncomingRequestFunc = func(ctx context.Context, req *record.IncomingRequest) (*insolar.ID, error) {
					if test.when == whenRegisterRequest {
						changePulse()
						// Due to specific implementation of HandleCall.handleActual
						// for this particular test we have to explicitly return
						// ErrCancelled. Otherwise it's possible that RegisterIncomingRequest
						// Procedure will return normally before Flow cancels it.
						return nil, flow.ErrCancelled
					}

					reqId := testutils.RandomID()
					return &reqId, nil
				}
			}

			if test.when > whenRegisterRequest {
				suite.am.HasPendingRequestsFunc = func(ctx context.Context, r insolar.Reference) (bool, error) {
					if test.when == whenHasPendingRequest {
						changePulse()

						// We have to implicitly return ErrCancelled to make f.Procedure return ErrCancelled as well
						// which will cause the correct code path to execute in logicrunner.HandleCall.
						// Otherwise the test has a race condition - f.Procedure can be cancelled or return normally.
						return false, flow.ErrCancelled
					}

					return false, nil
				}
			}

			if test.when > whenHasPendingRequest {
				suite.re.ExecuteAndSaveFunc = func(
					ctx context.Context, transcript *Transcript,
				) (insolar.Reply, error) {
					if test.when == whenCallMethod {
						changePulse()
					}

					return &reply.CallMethod{Result: []byte{3, 2, 1}}, nil
				}

				suite.re.SendReplyMock.Return()
			}

			wg := sync.WaitGroup{}
			wg.Add(len(test.messagesExpected))

			if len(test.messagesExpected) > 0 {
				suite.mb.SendFunc = func(
					ctx context.Context, msg insolar.Message, opts *insolar.MessageSendOptions,
				) (insolar.Reply, error) {

					if test.when == whenHasPendingRequest {
						// in case of whenHasPendingRequest we wait for at least one message from messagesExpected appear
						wg.Done()
					}
					wg.Done()

					if msg.Type() == insolar.TypeExecutorResults {
						require.Equal(t, test.pendingInExecutorResults, msg.(*message.ExecutorResults).Pending)
						require.Equal(t, test.queueLenInExecutorResults, len(msg.(*message.ExecutorResults).Queue))
					}

					switch msg.Type() {
					case insolar.TypeReturnResults,
						insolar.TypeExecutorResults,
						insolar.TypePendingFinished,
						insolar.TypeStillExecuting,
						insolar.TypeAdditionalCallFromPreviousExecutor:
						return &reply.OK{}, nil
					default:
						panic("no idea how to handle " + msg.Type().String())
					}
				}
			}

			msg := &message.CallMethod{
				IncomingRequest: record.IncomingRequest{
					Prototype: &protoRef,
					Object:    &objectRef,
					Method:    "some",
				},
			}

			parcel := testutils.NewParcelMock(suite.T())
			parcel.DefaultTargetMock.Return(&objectRef)
			parcel.MessageMock.Return(msg)
			parcel.TypeMock.Return(msg.Type())
			parcel.PulseMock.Return(insolar.PulseNumber(100))
			parcel.GetSenderMock.Return(notMeRef)

			ctx := inslogger.ContextWithTrace(suite.ctx, "req")
			pulseNum := pulsar.NewPulse(1, parcel.Pulse(), &entropygenerator.StandardEntropyGenerator{})
			err := suite.lr.OnPulse(ctx, *pulseNum)
			require.NoError(t, err)

			_, err = suite.lr.FlowDispatcher.WrapBusHandle(ctx, parcel)

			if test.flowCanceledExpected {
				require.EqualError(t, err, flow.ErrCancelled.Error())
				require.Equal(t, flow.ErrCancelled, err)
			} else if test.errorExpected {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			suite.Require().True(WaitGroup_TimeoutWait(&wg, 2*time.Minute),
				"Failed to wait for all requests to be processed")
		})
	}
}

func (s *LogicRunnerTestSuite) TestImmutableOrder() {
	// prepare default object and execution state
	objectRef := gen.Reference()
	es, broker := s.lr.StateStorage.UpsertExecutionState(objectRef)
	es.pending = message.NotPending

	// prepare request objects
	mutableRequestRef := gen.Reference()
	immutableRequestRef1 := gen.Reference()
	immutableRequestRef2 := gen.Reference()

	// prepare all three requests
	mutableRequest := record.IncomingRequest{
		ReturnMode:   record.ReturnResult,
		Object:       &objectRef,
		APIRequestID: utils.RandTraceID(),
		Immutable:    false,
	}
	mutableTranscript := NewTranscript(s.ctx, mutableRequestRef, mutableRequest)

	immutableRequest1 := record.IncomingRequest{
		ReturnMode:   record.ReturnResult,
		Object:       &objectRef,
		APIRequestID: utils.RandTraceID(),
		Immutable:    true,
	}
	immutableTranscript1 := NewTranscript(s.ctx, immutableRequestRef1, immutableRequest1)

	immutableRequest2 := record.IncomingRequest{
		ReturnMode:   record.ReturnResult,
		Object:       &objectRef,
		APIRequestID: utils.RandTraceID(),
		Immutable:    true,
	}
	immutableTranscript2 := NewTranscript(s.ctx, immutableRequestRef2, immutableRequest2)

	// Set custom executor, that'll:
	// 1) mutable will start execution and wait until something will ping it on channel 1
	// 2) immutable 1 will start execution and will wait on channel 2 until something will ping it
	// 3) immutable 2 will start execution and will ping on channel 2 and exit
	// 4) immutable 1 will ping on channel 1 and exit
	// 5) mutable request will continue execution and exit

	var mutableChan = make(chan interface{}, 1)
	var immutableChan chan interface{} = nil
	var immutableLock = sync.Mutex{}

	s.re.SendReplyMock.Return()
	s.re.ExecuteAndSaveMock.Set(func(ctx context.Context, transcript *Transcript) (insolar.Reply, error) {

		if transcript.RequestRef.Equal(mutableRequestRef) {
			log.Debug("mutableChan 1")
			select {
			case _ = <-mutableChan:

				log.Info("mutable got notifications")
				return &reply.CallMethod{Result: []byte{1, 2, 3}}, nil
			case <-time.After(2 * time.Minute):
				panic("timeout on waiting for immutable request 1 pinged us")
			}
		} else if transcript.RequestRef.Equal(immutableRequestRef1) || transcript.RequestRef.Equal(immutableRequestRef2) {
			newChan := false
			immutableLock.Lock()
			if immutableChan == nil {
				immutableChan = make(chan interface{}, 1)
				newChan = true
			}
			immutableLock.Unlock()
			if newChan {
				log.Debug("immutableChan 1")
				select {
				case _ = <-immutableChan:
					mutableChan <- struct{}{}
					log.Info("notify mutable chan and exit")
					return &reply.CallMethod{Result: []byte{1, 2, 3}}, nil
				case <-time.After(2 * time.Minute):
					panic("timeout on waiting for immutable request 2 pinged us")
				}
			} else {
				log.Info("notify immutable chan and exit")
				immutableChan <- struct{}{}
			}
		} else {
			panic("unreachable")
		}
		return &reply.CallMethod{Result: []byte{1, 2, 3}}, nil
	})

	broker.Put(s.ctx, true, mutableTranscript)
	broker.Put(s.ctx, true, immutableTranscript1, immutableTranscript2)

	s.True(wait(finishedCount, broker, 3))
}

func (s *LogicRunnerTestSuite) TestImmutableIsReal() {
	// prepare default object and execution state
	objectRef := gen.Reference()
	es, broker := s.lr.StateStorage.UpsertExecutionState(objectRef)
	es.pending = message.NotPending

	immutableRequestRef1 := gen.Reference()

	immutableRequest1 := record.IncomingRequest{
		ReturnMode:   record.ReturnResult,
		Object:       &objectRef,
		APIRequestID: utils.RandTraceID(),
		Immutable:    true,
	}
	immutableTranscript1 := NewTranscript(s.ctx, immutableRequestRef1, immutableRequest1)

	s.re.ExecuteAndSaveMock.Return(&reply.CallMethod{Result: []byte{1, 2, 3}}, nil)
	s.re.SendReplyMock.Return()

	broker.Put(s.ctx, true, immutableTranscript1)

	s.True(wait(finishedCount, broker, 1))
}

func TestLogicRunner(t *testing.T) {
	// Hello my friend! I bet you would like to place t.Parallel() here.
	// Of course this may sound as a good idea. This will run multiple
	// test in parallel which will make them execute faster. Right?
	// Wrong! You see, by historical reasons LogicRunnerTestSuite
	// is in fact 4 independent tests which share their state (suite.* fields).
	// Guess what happens when they run in parallel? Right, it seem to work
	// at first but after some time someone will spent a lot of exciting
	// days trying to figure out why these test sometimes fail (e.g. on CI).
	// In other words dont you dare to use t.Parallel() here unless you are
	// willing to completely rewrite the whole LogicRunnerTestSuite, OK?
	suite.Run(t, new(LogicRunnerTestSuite))
}

type LogicRunnerOnPulseTestSuite struct {
	LogicRunnerCommonTestSuite

	pulse     insolar.Pulse
	objectRef insolar.Reference
}

func (s *LogicRunnerOnPulseTestSuite) BeforeTest(suiteName, testName string) {
	s.LogicRunnerCommonTestSuite.BeforeTest(suiteName, testName)

	s.pulse = insolar.Pulse{}
	s.objectRef = testutils.RandomRef()
}

func (s *LogicRunnerOnPulseTestSuite) AfterTest(suiteName, testName string) {
	s.LogicRunnerCommonTestSuite.AfterTest(suiteName, testName)
}

// Empty state, expecting no error
func (s *LogicRunnerOnPulseTestSuite) TestEmptyLR() {
	err := s.lr.OnPulse(s.ctx, s.pulse)
	s.Require().NoError(err)
}

// We aren't next executor and we're not executing it
// Expecting empty state of object
func (s *LogicRunnerOnPulseTestSuite) TestEmptyES() {
	s.jc.MeMock.Return(insolar.Reference{})
	s.jc.IsAuthorizedMock.Return(false, nil)

	s.lr.StateStorage.UpsertExecutionState(s.objectRef)

	err := s.lr.OnPulse(s.ctx, s.pulse)
	s.Require().NoError(err)

	_, broker := s.lr.StateStorage.GetExecutionState(s.objectRef)
	s.Nil(broker)
}

// We aren't next executor and we're not executing it
// Expecting empty execution state
func (s *LogicRunnerOnPulseTestSuite) TestEmptyESWithValidation() {
	s.jc.MeMock.Return(insolar.Reference{})
	s.jc.IsAuthorizedMock.Return(false, nil)

	s.lr.StateStorage.UpsertExecutionState(s.objectRef)
	s.lr.StateStorage.UpsertValidationState(s.objectRef)

	err := s.lr.OnPulse(s.ctx, s.pulse)
	s.Require().NoError(err)

	_, broker := s.lr.StateStorage.GetExecutionState(s.objectRef)
	s.Nil(broker)
}

// We aren't next executor but we're currently executing
// Expecting we send message to new executor and moving state to InPending
func (s *LogicRunnerOnPulseTestSuite) TestESWithValidationCurrent() {
	s.jc.MeMock.Return(insolar.Reference{})
	s.jc.IsAuthorizedMock.Return(false, nil)
	s.mb.SendMock.Return(&reply.ID{}, nil)

	es, broker := s.lr.StateStorage.UpsertExecutionState(s.objectRef)
	es.pending = message.NotPending
	// we should set empty current execution here, since we added new
	// logic with not empty number of elements in CurrentList
	broker.currentList.Set(s.objectRef, &Transcript{})

	err := s.lr.OnPulse(s.ctx, s.pulse)
	s.Require().NoError(err)

	es, _ = s.lr.StateStorage.GetExecutionState(s.objectRef)
	s.Equal(message.InPending, es.pending)
	broker.currentList.Cleanup()
}

// We aren't next executor but we're currently executing and queue isn't empty.
// Expecting we send message to new executor and moving state to InPending
func (s *LogicRunnerOnPulseTestSuite) TestWithNotEmptyQueue() {
	s.jc.MeMock.Return(insolar.Reference{})
	s.jc.IsAuthorizedMock.Return(false, nil)
	s.mb.SendMock.Return(&reply.ID{}, nil)

	es, broker := s.lr.StateStorage.UpsertExecutionState(s.objectRef)
	broker.currentList.Set(s.objectRef, &Transcript{})

	reqRef := gen.Reference()
	broker.mutable.Push(&Transcript{
		Context:    s.ctx,
		RequestRef: reqRef,
		Request:    &record.IncomingRequest{},
	})
	es.pending = message.NotPending

	err := s.lr.OnPulse(s.ctx, s.pulse)
	s.Require().NoError(err)
	es, _ = s.lr.StateStorage.GetExecutionState(s.objectRef)
	s.Equal(message.InPending, es.pending)
}

// We aren't next executor but we're currently executing.
// Expecting sending message to new executor and moving state to InPending
func (s *LogicRunnerOnPulseTestSuite) TestWithEmptyQueue() {
	s.jc.MeMock.Return(insolar.Reference{})
	s.jc.IsAuthorizedMock.Return(false, nil)
	s.mb.SendMock.Return(&reply.ID{}, nil)

	es, broker := s.lr.StateStorage.UpsertExecutionState(s.objectRef)
	broker.currentList.Set(s.objectRef, &Transcript{})
	es.pending = message.NotPending

	err := s.lr.OnPulse(s.ctx, s.pulse)
	s.Require().NoError(err)
	es, _ = s.lr.StateStorage.GetExecutionState(s.objectRef)
	s.Equal(message.InPending, es.pending)
}

// Executor is on the same node and we're currently executing
// Expecting task to be moved to NotPending
func (s *LogicRunnerOnPulseTestSuite) TestExecutorSameNode() {
	s.jc.MeMock.Return(insolar.Reference{})
	s.jc.IsAuthorizedMock.Return(true, nil)

	es, broker := s.lr.StateStorage.UpsertExecutionState(s.objectRef)
	es.pending = message.NotPending
	broker.currentList.Set(s.objectRef, &Transcript{})
	InitBroker(s.T(), s.ctx, 0, broker, false)

	err := s.lr.OnPulse(s.ctx, s.pulse)
	s.Require().NoError(err)
	es, _ = s.lr.StateStorage.GetExecutionState(s.objectRef)
	s.Require().Equal(message.NotPending, es.pending)
	broker.currentList.Cleanup()
}

// We're the next executor, task was currently executing and in InPending.
// Expecting task to moved from InPending -> NotPending
func (s *LogicRunnerOnPulseTestSuite) TestStateTransfer1() {
	s.jc.MeMock.Return(insolar.Reference{})
	s.jc.IsAuthorizedMock.Return(true, nil)

	es, broker := s.lr.StateStorage.UpsertExecutionState(s.objectRef)
	broker.currentList.Set(s.objectRef, &Transcript{})
	es.pending = message.InPending

	err := s.lr.OnPulse(s.ctx, s.pulse)
	s.Require().NoError(err)
	es, _ = s.lr.StateStorage.GetExecutionState(s.objectRef)
	s.Require().Equal(message.NotPending, es.pending)
}

// We're the next executor and no one confirmed that this task is executing
// move task from InPending -> NotPending
func (s *LogicRunnerOnPulseTestSuite) TestStateTransfer2() {
	s.jc.MeMock.Return(insolar.Reference{})
	s.jc.IsAuthorizedMock.Return(true, nil)

	es, _ := s.lr.StateStorage.UpsertExecutionState(s.objectRef)
	es.pending = message.InPending
	es.PendingConfirmed = false

	err := s.lr.OnPulse(s.ctx, s.pulse)
	s.Require().NoError(err)
	es, _ = s.lr.StateStorage.GetExecutionState(s.objectRef)
	s.Require().Equal(message.NotPending, es.pending)
}

// We're the next executor and previous confirmed that this task is executing
// still in pending
// but we expect that previous executor come to us for token
func (s *LogicRunnerOnPulseTestSuite) TestStateTransfer3() {
	s.jc.MeMock.Return(insolar.Reference{})
	s.jc.IsAuthorizedMock.Return(true, nil)

	es, _ := s.lr.StateStorage.UpsertExecutionState(s.objectRef)
	es.pending = message.InPending
	es.PendingConfirmed = true

	err := s.lr.OnPulse(s.ctx, s.pulse)
	s.Require().NoError(err)

	es, _ = s.lr.StateStorage.GetExecutionState(s.objectRef)
	// we still in pending
	s.Equal(message.InPending, es.pending)
	// but we expect that previous executor come to us for token
	s.Equal(false, es.PendingConfirmed)
}

// We're not the next executor, so we must send this task to the next executor
func (s *LogicRunnerOnPulseTestSuite) TestSendTaskToNextExecutor() {
	s.jc.MeMock.Return(insolar.Reference{})
	s.jc.IsAuthorizedMock.Return(false, nil)
	s.mb.SendMock.Return(&reply.ID{}, nil)

	es, _ := s.lr.StateStorage.UpsertExecutionState(s.objectRef)
	es.pending = message.InPending
	es.PendingConfirmed = false

	err := s.lr.OnPulse(s.ctx, s.pulse)
	s.Require().NoError(err)

	_, broker := s.lr.StateStorage.GetExecutionState(s.objectRef)
	s.Nil(broker)
}

func (s *LogicRunnerOnPulseTestSuite) TestLedgerHasMoreRequests() {
	s.jc.IsAuthorizedMock.Return(false, nil)
	s.jc.MeMock.Return(insolar.Reference{})

	var testCases = map[string]struct {
		Count           int
		hasMoreRequests bool
	}{
		"Has": {
			maxQueueLength + 1,
			true,
		},
		"Don't": {
			maxQueueLength,
			false,
		},
	}

	for name, test := range testCases {
		s.T().Run(name, func(t *testing.T) {
			a := assert.New(t)

			_, broker := s.lr.StateStorage.UpsertExecutionState(s.objectRef)
			InitBroker(t, s.ctx, test.Count, broker, false)

			messagesQueue := convertQueueToMessageQueue(s.ctx, broker.mutable.Peek(maxQueueLength))

			expectedMessage := &message.ExecutorResults{
				RecordRef:             s.objectRef,
				Queue:                 messagesQueue,
				LedgerHasMoreRequests: test.hasMoreRequests,
			}

			wg := sync.WaitGroup{}
			wg.Add(1)
			s.mb.SendMock.Set(func(p context.Context, p1 insolar.Message, p2 *insolar.MessageSendOptions) (r insolar.Reply, r1 error) {
				a.Equal(expectedMessage, p1)
				wg.Done()
				return nil, nil
			})

			err := s.lr.OnPulse(s.ctx, s.pulse)
			a.NoError(err)

			wg.Wait()
		})
	}
}

func TestLogicRunnerOnPulse(t *testing.T) {
	suite.Run(t, new(LogicRunnerOnPulseTestSuite))
}

func TestLRUnsafeGetLedgerPendingRequest(t *testing.T) {
	suite.Run(t, new(LRUnsafeGetLedgerPendingRequestTestSuite))
}

type LRUnsafeGetLedgerPendingRequestTestSuite struct {
	LogicRunnerCommonTestSuite

	pulse                 insolar.Pulse
	ref                   insolar.Reference
	currentPulseNumber    insolar.PulseNumber
	oldRequestPulseNumber insolar.PulseNumber
}

func (s *LRUnsafeGetLedgerPendingRequestTestSuite) BeforeTest(suiteName, testName string) {
	s.LogicRunnerCommonTestSuite.BeforeTest(suiteName, testName)

	s.pulse = insolar.Pulse{}
	s.ref = testutils.RandomRef()
	s.currentPulseNumber = 3
	s.oldRequestPulseNumber = 1
}

func (s *LRUnsafeGetLedgerPendingRequestTestSuite) AfterTest(suiteName, testName string) {
	s.LogicRunnerCommonTestSuite.AfterTest(suiteName, testName)
}

func (s *LRUnsafeGetLedgerPendingRequestTestSuite) TestAlreadyHaveLedgerQueueElement() {
	objectRef := gen.Reference()
	es, broker := s.lr.StateStorage.UpsertExecutionState(objectRef)

	reqRef := gen.Reference()
	broker.Put(s.ctx, false, &Transcript{
		FromLedger:   true,
		LogicContext: &insolar.LogicCallContext{},
		RequestRef:   reqRef,
		Request:      &record.IncomingRequest{},
	})

	proc := UnsafeGetLedgerPendingRequest{es: es, broker: broker, dep: &Dependencies{lr: s.lr}}
	err := proc.Proceed(s.ctx)
	s.Require().NoError(err)

	// we check that there is no unexpected calls to A.M., as we already have element
	// from ledger another call to the ledger will return the same request, so we make
	// sure it doesn't happen
}

func (s *LRUnsafeGetLedgerPendingRequestTestSuite) TestNoMoreRequestsInExecutionState() {
	objectRef := gen.Reference()
	es, broker := s.lr.StateStorage.UpsertExecutionState(objectRef)
	es.LedgerHasMoreRequests = false

	proc := UnsafeGetLedgerPendingRequest{es: es, broker: broker, dep: &Dependencies{lr: s.lr}}
	err := proc.Proceed(s.ctx)
	s.Require().NoError(err)
	s.Require().Nil(broker.HasLedgerRequest(s.ctx))
}

func (s *LRUnsafeGetLedgerPendingRequestTestSuite) TestNoMoreRequestsInLedger() {
	objectRef := gen.Reference()
	es, broker := s.lr.StateStorage.UpsertExecutionState(objectRef)
	es.LedgerHasMoreRequests = true

	am := artifacts.NewClientMock(s.mc)
	am.GetPendingRequestMock.Return(nil, nil, insolar.ErrNoPendingRequest)
	s.lr.ArtifactManager = am
	proc := UnsafeGetLedgerPendingRequest{es: es, broker: broker, dep: &Dependencies{lr: s.lr}}
	err := proc.Proceed(s.ctx)
	s.Require().NoError(err)
}

func (s *LRUnsafeGetLedgerPendingRequestTestSuite) TestDoesNotAuthorized() {
	objectRef := gen.Reference()
	es, broker := s.lr.StateStorage.UpsertExecutionState(objectRef)
	es.LedgerHasMoreRequests = true

	s.am.GetPendingRequestMock.Return(nil, &record.IncomingRequest{}, nil)

	// we doesn't authorized (pulse change in time we process function)
	s.ps.LatestMock.Return(insolar.Pulse{PulseNumber: s.currentPulseNumber}, nil)
	s.jc.IsAuthorizedMock.Return(false, nil)
	s.jc.MeMock.Return(insolar.Reference{})

	proc := UnsafeGetLedgerPendingRequest{es: es, broker: broker, dep: &Dependencies{lr: s.lr}}
	err := proc.Proceed(s.ctx)
	s.Require().NoError(err)
	s.Require().Nil(broker.HasLedgerRequest(s.ctx))
}

func (s LRUnsafeGetLedgerPendingRequestTestSuite) TestUnsafeGetLedgerPendingRequest() {
	objectRef := gen.Reference()
	es, broker := s.lr.StateStorage.UpsertExecutionState(objectRef)
	es.LedgerHasMoreRequests = true

	testRequestRef := gen.Reference()
	testRequest := record.IncomingRequest{Object: &s.ref}

	s.am.GetPendingRequestMock.Return(&testRequestRef, &testRequest, nil)

	s.ps.LatestMock.Return(insolar.Pulse{PulseNumber: s.currentPulseNumber}, nil)
	s.jc.IsAuthorizedMock.Return(true, nil)
	s.jc.MeMock.Return(insolar.Reference{})

	proc := UnsafeGetLedgerPendingRequest{es: es, broker: broker, dep: &Dependencies{lr: s.lr}}
	err := proc.Proceed(s.ctx)
	s.Require().NoError(err)

	s.Require().Equal(true, es.LedgerHasMoreRequests)
	ledgerRequest := broker.HasLedgerRequest(s.ctx)
	s.Require().NotNil(ledgerRequest)
}

func (suite *LogicRunnerTestSuite) TestInitializeExecutionState() {
	suite.T().Run("InitializeExecutionState copy queue properly", func(t *testing.T) {
		pulseObj := insolar.Pulse{}
		pulseObj.PulseNumber = insolar.FirstPulseNumber

		object := testutils.RandomRef()
		defer delete(*suite.lr.StateStorage.StateMap(), object)

		firstRef := testutils.RandomRef()
		firstElement := message.ExecutionQueueElement{RequestRef: firstRef, Request: record.IncomingRequest{Immutable: false}}

		secondRef := testutils.RandomRef()
		secondElement := message.ExecutionQueueElement{RequestRef: secondRef, Request: record.IncomingRequest{Immutable: false}}

		msg := &message.ExecutorResults{
			Caller:    testutils.RandomRef(),
			RecordRef: object,
			Queue:     []message.ExecutionQueueElement{firstElement, secondElement},
		}

		proc := initializeExecutionState{
			LR:  suite.lr,
			msg: msg,
		}
		err := proc.Proceed(suite.ctx)
		require.NoError(t, err)

		require.Equal(t, &secondRef, proc.Result.broker.mutable.first.value.RequestRef)
		require.Equal(t, &firstRef, proc.Result.broker.mutable.last.value.RequestRef)
	})
}
