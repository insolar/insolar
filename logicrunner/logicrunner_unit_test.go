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
	"bytes"
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	message2 "github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/infrastructure/gochannel"
	"github.com/gojuno/minimock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/pulsar"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/network"
)

type LogicRunnerCommonTestSuite struct {
	suite.Suite

	mc     *minimock.Controller
	ctx    context.Context
	am     *artifacts.ClientMock
	dc     *artifacts.DescriptorsCacheMock
	mb     *testutils.MessageBusMock
	jc     *jet.CoordinatorMock
	mm     *mmanager
	lr     *LogicRunner
	re     *RequestsExecutorMock
	ps     *pulse.AccessorMock
	mle    *testutils.MachineLogicExecutorMock
	nn     *network.NodeNetworkMock
	sender *bus.SenderMock
	pub    message2.Publisher
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
	suite.sender = bus.NewSenderMock(suite.mc)
	suite.pub = &publisherMock{}

	suite.SetupLogicRunner()
}

func (suite *LogicRunnerCommonTestSuite) SetupLogicRunner() {
	suite.sender = bus.NewSenderMock(suite.mc)
	suite.pub = &publisherMock{}
	suite.lr, _ = NewLogicRunner(&configuration.LogicRunner{}, suite.pub, suite.sender)
	suite.lr.ArtifactManager = suite.am
	suite.lr.DescriptorsCache = suite.dc
	suite.lr.MessageBus = suite.mb
	suite.lr.MachinesManager = suite.mm
	suite.lr.JetCoordinator = suite.jc
	suite.lr.PulseAccessor = suite.ps
	suite.lr.NodeNetwork = suite.nn
	suite.lr.Sender = suite.sender
	suite.lr.Publisher = suite.pub
	suite.lr.RequestsExecutor = suite.re

	_ = suite.lr.Init(suite.ctx)

	suite.lr.FlowDispatcher.PulseAccessor = suite.ps
}

func (suite *LogicRunnerCommonTestSuite) AfterTest(suiteName, testName string) {
	suite.mc.Wait(2 * time.Second)
	//suite.mc.Finish()

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

func mockSender(suite *LogicRunnerTestSuite) chan *message2.Message {
	replyChan := make(chan *message2.Message, 1)
	suite.sender.ReplyFunc = func(p context.Context, p1 payload.Meta, p2 *message2.Message) {
		replyChan <- p2
	}
	return replyChan
}

func getReply(suite *LogicRunnerTestSuite, replyChan chan *message2.Message) (insolar.Reply, error) {
	res := <-replyChan
	re, err := reply.Deserialize(bytes.NewBuffer(res.Payload))
	if err != nil {
		payloadType, err := payload.UnmarshalType(res.Payload)
		suite.Require().NoError(err)
		suite.Require().EqualValues(payload.TypeError, payloadType)

		pl, err := payload.Unmarshal(res.Payload)
		suite.Require().NoError(err)
		p, ok := pl.(*payload.Error)
		suite.Require().True(ok)
		return nil, errors.New(p.Text)
	}
	return re, nil
}

func (suite *LogicRunnerTestSuite) TestSagaCallAcceptNotificationHandler() {
	outgoing := record.OutgoingRequest{
		Caller: gen.Reference(),
		Reason: gen.Reference(),
	}
	virtual := record.Wrap(&outgoing)
	outgoingBytes, err := virtual.Marshal()
	suite.Require().NoError(err)

	outgoingReqId := gen.ID()
	outgoingRequestRef := insolar.NewReference(outgoingReqId)

	pl := &payload.SagaCallAcceptNotification{
		DetachedRequestID: outgoingReqId,
		Request:           outgoingBytes,
	}
	msg, err := payload.NewMessage(pl)
	suite.Require().NoError(err)

	pulseNum := pulsar.NewPulse(0, insolar.FirstPulseNumber, &entropygenerator.StandardEntropyGenerator{})

	suite.ps.LatestMock.Return(*pulseNum, nil)

	msg.Metadata.Set(bus.MetaPulse, pulseNum.PulseNumber.String())
	sp, err := instracer.Serialize(context.Background())
	suite.Require().NoError(err)
	msg.Metadata.Set(bus.MetaSpanData, string(sp))

	meta := payload.Meta{
		Payload: msg.Payload,
	}
	buf, err := meta.Marshal()
	msg.Payload = buf

	dummyRequestRef := gen.Reference()
	callMethodChan := make(chan struct{})
	var usedCaller insolar.Reference
	var usedReason insolar.Reference
	var usedReturnMode record.ReturnMode

	cr := testutils.NewContractRequesterMock(suite.T())
	cr.CallMethodFunc = func(ctx context.Context, msg insolar.Message) (insolar.Reply, error) {
		suite.Require().Equal(insolar.TypeCallMethod, msg.Type())
		cm := msg.(*message.CallMethod)
		usedCaller = cm.Caller
		usedReason = cm.Reason
		usedReturnMode = cm.ReturnMode

		result := &reply.RegisterRequest{
			Request: dummyRequestRef,
		}
		callMethodChan <- struct{}{}
		return result, nil
	}
	suite.lr.ContractRequester = cr

	registerResultChan := make(chan struct{})
	var usedRequestRef insolar.Reference
	var usedResult []byte

	am := artifacts.NewClientMock(suite.T())
	am.RegisterResultMock.Set(func(ctx context.Context, reqRef insolar.Reference, reqResults artifacts.RequestResult) (r error) {
		usedRequestRef = reqRef
		usedResult = reqResults.Result()
		registerResultChan <- struct{}{}
		return nil
	})
	suite.lr.ArtifactManager = am

	_, err = suite.lr.FlowDispatcher.Process(msg)
	suite.Require().NoError(err)

	<-callMethodChan
	suite.Require().Equal(outgoing.Caller, usedCaller)
	suite.Require().Equal(outgoing.Reason, usedReason)
	suite.Require().Equal(record.ReturnNoWait, usedReturnMode)

	<-registerResultChan
	suite.Require().Equal(outgoingRequestRef, &usedRequestRef)
	suite.Require().Equal(dummyRequestRef.Bytes(), usedResult)

	// In this test LME doesn't need any reply from VE. But if an reply was
	// required you could check it like this:
	// ```
	// replyChan = mockSender(suite)
	// rep, err := getReply(suite, replyChan)
	// suite.Require().NoError(err)
	// suite.Require().Equal(&reply.OK{}, rep)
	// ```
}

func (suite *LogicRunnerTestSuite) TestNewLogicRunner() {
	lr, err := NewLogicRunner(nil, suite.pub, suite.sender)
	suite.Require().Error(err)
	suite.Require().Nil(lr)

	lr, err = NewLogicRunner(&configuration.LogicRunner{}, suite.pub, suite.sender)
	suite.Require().NoError(err)
	suite.Require().NotNil(lr)
	_ = lr.Stop(context.Background())
}

func (suite *LogicRunnerTestSuite) TestStartStop() {
	lr, err := NewLogicRunner(&configuration.LogicRunner{
		BuiltIn: &configuration.BuiltIn{},
	}, suite.pub, suite.sender)
	suite.Require().NoError(err)
	suite.Require().NotNil(lr)

	lr.MessageBus = suite.mb

	lr.MachinesManager = suite.mm

	suite.am.InjectCodeDescriptorMock.Return()
	suite.am.InjectObjectDescriptorMock.Return()
	suite.am.InjectFinishMock.Return()
	lr.ArtifactManager = suite.am

	err = lr.Init(suite.ctx)
	suite.Require().NoError(err)

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

	pulseNum := insolar.PulseNumber(insolar.FirstPulseNumber)

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

	suite.am.HasPendingsMock.Return(false, nil)

	suite.am.RegisterIncomingRequestMock.Set(func(ctx context.Context, r *record.IncomingRequest) (*insolar.ID, error) {
		reqId := testutils.RandomID()
		return &reqId, nil
	})

	suite.re.ExecuteAndSaveMock.Return(nil, nil)
	suite.re.SendReplyMock.Return()

	num := 100
	wg := sync.WaitGroup{}
	wg.Add(num)

	suite.sender.ReplyFunc = func(p context.Context, p1 payload.Meta, p2 *message2.Message) {
		wg.Done()
	}

	suite.ps.LatestFunc = func(p context.Context) (r insolar.Pulse, r1 error) {
		return insolar.Pulse{PulseNumber: pulseNum}, nil
	}
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

			parcel := &message.Parcel{
				Sender:      notMeRef,
				Msg:         msg,
				PulseNumber: pulseNum,
			}

			wrapper := payload.Meta{
				Payload: message.ParcelToBytes(parcel),
				Sender:  notMeRef,
				Pulse:   pulseNum,
			}
			buf, err := wrapper.Marshal()
			suite.Require().NoError(err)

			wmMsg := message2.NewMessage(watermill.NewUUID(), buf)
			wmMsg.Metadata.Set(bus.MetaPulse, pulseNum.String())
			sp, err := instracer.Serialize(context.Background())
			suite.Require().NoError(err)
			wmMsg.Metadata.Set(bus.MetaSpanData, string(sp))
			wmMsg.Metadata.Set(bus.MetaType, fmt.Sprintf("%s", msg.Type()))
			wmMsg.Metadata.Set(bus.MetaTraceID, "req-"+strconv.Itoa(i))

			_, err = suite.lr.FlowDispatcher.Process(wmMsg)
			suite.Require().NoError(err)
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
	var pn insolar.PulseNumber = insolar.FirstPulseNumber
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
		whenHasPendings
		whenCallMethod
	)

	table := []struct {
		name                      string
		when                      whenType
		messagesExpected          []insolar.MessageType
		errorExpected             bool
		flowCanceledExpected      bool
		pendingInExecutorResults  insolar.PendingState
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
		// two cases below two un-deterministic, created a task to write proper
		// test cases
		//{
		//	name: "pulse change in HasPendings",
		//	when: whenHasPendings,
		//	messagesExpected: []insolar.MessageType{
		//		insolar.TypeExecutorResults,
		//	},
		//	pendingInExecutorResults:  insolar.PendingUnknown,
		//	queueLenInExecutorResults: 1,
		//},
		//{
		//	name: "pulse change in CallMethod",
		//	when: whenCallMethod,
		//	messagesExpected: []insolar.MessageType{
		//		insolar.TypeExecutorResults, insolar.TypePendingFinished, insolar.TypeStillExecuting,
		//	},
		//	pendingInExecutorResults:  insolar.InPending,
		//	queueLenInExecutorResults: 0,
		//},
	}

	for _, test := range table {
		test := test
		suite.T().Run(test.name, func(t *testing.T) {
			lck.Lock()
			pn = insolar.FirstPulseNumber
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
				if pnArg == insolar.FirstPulseNumber+1 {
					return false, nil
				}

				if test.when == whenIsAuthorized {
					// Please note that changePulse calls LogicRunner.ChangePulse which calls IsAuthorized.
					// In other words this procedure is not called sequentially!
					changePulse()
				}

				lck.Lock()
				defer lck.Unlock()

				return pn == insolar.FirstPulseNumber, nil
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
				suite.am.HasPendingsFunc = func(ctx context.Context, r insolar.Reference) (bool, error) {
					if test.when == whenHasPendings {
						changePulse()

						// We have to implicitly return ErrCancelled to make f.Procedure return ErrCancelled as well
						// which will cause the correct code path to execute in logicrunner.HandleCall.
						// Otherwise the test has a race condition - f.Procedure can be cancelled or return normally.
						return false, flow.ErrCancelled
					}

					return false, nil
				}
			}

			if test.when > whenHasPendings {
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
					// AdditionalCallFromPreviousExecutor is not deterministic
					if msg.Type() == insolar.TypeAdditionalCallFromPreviousExecutor {
						return &reply.OK{}, nil
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
						insolar.TypeStillExecuting:
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

			parcel := &message.Parcel{
				Sender:      notMeRef,
				Msg:         msg,
				PulseNumber: insolar.PulseNumber(insolar.FirstPulseNumber),
			}

			ctx := inslogger.ContextWithTrace(suite.ctx, "req")

			pulseNum := pulsar.NewPulse(1, parcel.Pulse()-1, &entropygenerator.StandardEntropyGenerator{})
			err := suite.lr.OnPulse(ctx, *pulseNum)
			require.NoError(t, err)

			wrapper := payload.Meta{
				Payload: message.ParcelToBytes(parcel),
				Sender:  notMeRef,
				Pulse:   insolar.PulseNumber(insolar.FirstPulseNumber),
			}
			buf, err := wrapper.Marshal()
			suite.Require().NoError(err)

			wmMsg := message2.NewMessage(watermill.NewUUID(), buf)
			wmMsg.Metadata.Set(bus.MetaType, fmt.Sprintf("%s", msg.Type()))
			wmMsg.Metadata.Set(bus.MetaTraceID, inslogger.TraceID(ctx))
			wmMsg.Metadata.Set(bus.MetaPulse, pulseNum.PulseNumber.String())
			sp, err := instracer.Serialize(context.Background())
			suite.Require().NoError(err)
			wmMsg.Metadata.Set(bus.MetaSpanData, string(sp))

			replyChan := mockSender(suite)
			_, err = suite.lr.FlowDispatcher.Process(wmMsg)

			if test.flowCanceledExpected {
				_, err := getReply(suite, replyChan)
				require.EqualError(t, err, flow.ErrCancelled.Error())
			} else if test.errorExpected {
				_, err := getReply(suite, replyChan)
				require.Error(t, err)
			} else {
				_, err := getReply(suite, replyChan)
				require.NoError(t, err)
			}

			suite.Require().True(WaitGroup_TimeoutWait(&wg, 2*time.Second),
				"Failed to wait for all requests to be processed")
		})
	}
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

func TestLogicRunner_OnPulse(t *testing.T) {
	table := []struct {
		name  string
		mocks func(ctx context.Context, t minimock.Tester) *LogicRunner
	}{
		{
			name: "broker that stays and sends messages",
			mocks: func(ctx context.Context, mc minimock.Tester) *LogicRunner {
				lr, err := NewLogicRunner(&configuration.LogicRunner{}, nil, nil)
				require.NoError(t, err)

				lr.initHandlers()

				lr.JetCoordinator = jet.NewCoordinatorMock(mc).
					MeMock.Return(gen.Reference()).
					IsAuthorizedMock.Return(true, nil)

				lr.MessageBus = testutils.NewMessageBusMock(mc).
					SendMock.Return(&reply.OK{}, nil)

				broker := NewExecutionBrokerIMock(mc).
					OnPulseMock.Return(
					false,
					[]insolar.Message{&message.ExecutorResults{}},
				)
				stateMap := map[insolar.Reference]*ObjectState{
					gen.Reference(): {
						ExecutionBroker: broker,
					},
				}
				lr.StateStorage = NewStateStorageMock(mc).
					LockMock.Return().
					UnlockMock.Return().
					StateMapMock.Return(&stateMap)

				return lr
			},
		},
		{
			name: "broker that goes way",
			mocks: func(ctx context.Context, mc minimock.Tester) *LogicRunner {
				lr, err := NewLogicRunner(&configuration.LogicRunner{}, nil, nil)
				require.NoError(t, err)

				lr.initHandlers()

				lr.JetCoordinator = jet.NewCoordinatorMock(mc).
					MeMock.Return(gen.Reference()).
					IsAuthorizedMock.Return(true, nil)

				stateMap := map[insolar.Reference]*ObjectState{
					gen.Reference(): {
						ExecutionBroker: NewExecutionBrokerIMock(mc).
							OnPulseMock.Return(true, []insolar.Message{}),
					},
				}
				lr.StateStorage = NewStateStorageMock(mc).
					LockMock.Return().
					UnlockMock.Return().
					StateMapMock.Return(&stateMap).
					DeleteObjectStateMock.Return()

				return lr
			},
		},
		{
			name: "one empty object state record",
			mocks: func(ctx context.Context, mc minimock.Tester) *LogicRunner {
				lr, err := NewLogicRunner(&configuration.LogicRunner{}, nil, nil)
				require.NoError(t, err)

				lr.initHandlers()

				lr.JetCoordinator = jet.NewCoordinatorMock(mc).
					MeMock.Return(gen.Reference()).
					IsAuthorizedMock.Return(true, nil)

				stateMap := map[insolar.Reference]*ObjectState{
					gen.Reference(): {},
				}
				lr.StateStorage = NewStateStorageMock(mc).
					LockMock.Return().
					UnlockMock.Return().
					StateMapMock.Return(&stateMap).
					DeleteObjectStateMock.Return()

				return lr
			},
		},
		{
			name: "empty state map",
			mocks: func(ctx context.Context, mc minimock.Tester) *LogicRunner {
				lr, err := NewLogicRunner(&configuration.LogicRunner{}, nil, nil)
				require.NoError(t, err)

				lr.initHandlers()

				stateMap := map[insolar.Reference]*ObjectState{}
				lr.StateStorage = NewStateStorageMock(mc).
					LockMock.Return().
					UnlockMock.Return().
					StateMapMock.Return(&stateMap)

				return lr
			},
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			ctx := inslogger.TestContext(t)
			mc := minimock.NewController(t)

			lr := test.mocks(ctx, mc)
			err := lr.OnPulse(ctx, insolar.Pulse{})
			require.NoError(t, err)

			mc.Wait(3 * time.Second)
			mc.Finish()
		})
	}
}
