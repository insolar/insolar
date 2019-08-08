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
	"github.com/gojuno/minimock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
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
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/writecontroller"
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
	cr     *testutils.ContractRequesterMock
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
	suite.cr = testutils.NewContractRequesterMock(suite.mc)
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
	suite.lr.ContractRequester = suite.cr
	suite.lr.PulseAccessor = suite.ps

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

func mockSender(suite *LogicRunnerTestSuite) chan *message2.Message {
	replyChan := make(chan *message2.Message, 1)
	suite.sender.ReplyMock.Set(func(p context.Context, p1 payload.Meta, p2 *message2.Message) {
		replyChan <- p2
	})
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
		Caller:     gen.Reference(),
		Reason:     gen.Reference(),
		ReturnMode: record.ReturnSaga,
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

	suite.cr.CallMock.Set(func(ctx context.Context, msg insolar.Message) (insolar.Reply, error) {
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
	})

	registerResultChan := make(chan struct{})
	var usedRequestRef insolar.Reference
	var usedResult []byte

	suite.am.RegisterResultMock.Set(func(ctx context.Context, reqRef insolar.Reference, reqResults artifacts.RequestResult) (r error) {
		usedRequestRef = reqRef
		usedResult = reqResults.Result()
		registerResultChan <- struct{}{}
		return nil
	})

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

	notMeRef := testutils.RandomRef()

	pulseNum := insolar.PulseNumber(insolar.FirstPulseNumber)

	suite.jc.IsMeAuthorizedNowMock.Return(true, nil)

	syncT := &utils.SyncT{T: suite.T()}
	meRef := testutils.RandomRef()
	nodeMock := network.NewNetworkNodeMock(syncT)
	nodeMock.IDMock.Return(meRef)

	od := artifacts.NewObjectDescriptorMock(syncT)
	od.PrototypeMock.Return(&protoRef, nil)
	od.MemoryMock.Return([]byte{1, 2, 3})
	od.ParentMock.Return(&parentRef)
	od.HeadRefMock.Return(&objectRef)

	pd := artifacts.NewObjectDescriptorMock(syncT)
	pd.CodeMock.Return(&codeRef, nil)
	pd.HeadRefMock.Return(&protoRef)

	cd := artifacts.NewCodeDescriptorMock(syncT)
	cd.MachineTypeMock.Return(insolar.MachineTypeBuiltin)
	cd.RefMock.Return(&codeRef)

	suite.am.HasPendingsMock.Return(false, nil)

	suite.am.RegisterIncomingRequestMock.Set(func(ctx context.Context, r *record.IncomingRequest) (*payload.RequestInfo, error) {
		reqId := testutils.RandomID()
		return &payload.RequestInfo{RequestID: reqId, ObjectID: *objectRef.Record()}, nil
	})

	suite.re.ExecuteAndSaveMock.Return(nil, nil)
	suite.re.SendReplyMock.Return()

	num := 100
	wg := sync.WaitGroup{}
	wg.Add(num)

	suite.sender.ReplyMock.Set(func(p context.Context, p1 payload.Meta, p2 *message2.Message) {
		wg.Done()
	})

	suite.ps.LatestMock.Set(func(p context.Context) (r insolar.Pulse, r1 error) {
		return insolar.Pulse{PulseNumber: pulseNum}, nil
	})
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
			require.NoError(syncT, err)

			wmMsg := message2.NewMessage(watermill.NewUUID(), buf)
			wmMsg.Metadata.Set(bus.MetaPulse, pulseNum.String())
			sp, err := instracer.Serialize(context.Background())
			require.NoError(syncT, err)
			wmMsg.Metadata.Set(bus.MetaSpanData, string(sp))
			wmMsg.Metadata.Set(bus.MetaType, fmt.Sprintf("%s", msg.Type()))
			wmMsg.Metadata.Set(bus.MetaTraceID, "req-"+strconv.Itoa(i))

			_, err = suite.lr.FlowDispatcher.Process(wmMsg)
			require.NoError(syncT, err)
		}(i)
	}

	suite.Require().True(WaitGroup_TimeoutWait(&wg, 2*time.Minute),
		"Failed to wait for all requests to be processed")
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

				lr.MessageBus = testutils.NewMessageBusMock(mc).
					SendMock.Return(&reply.OK{}, nil)

				lr.StateStorage = NewStateStorageMock(mc).
					LockMock.Return().
					UnlockMock.Return().
					IsEmptyMock.Return(false).
					OnPulseMock.Return([]insolar.Message{&message.ExecutorResults{}})

				lr.WriteController = writecontroller.NewWriteController()
				_ = lr.WriteController.Open(ctx, insolar.FirstPulseNumber)

				return lr
			},
		},
		{
			name: "broker that goes way",
			mocks: func(ctx context.Context, mc minimock.Tester) *LogicRunner {
				lr, err := NewLogicRunner(&configuration.LogicRunner{}, nil, nil)
				require.NoError(t, err)

				lr.initHandlers()

				lr.StateStorage = NewStateStorageMock(mc).
					LockMock.Return().
					UnlockMock.Return().
					IsEmptyMock.Return(true).
					OnPulseMock.Return([]insolar.Message{})

				lr.WriteController = writecontroller.NewWriteController()
				_ = lr.WriteController.Open(ctx, insolar.FirstPulseNumber)

				return lr
			},
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			ctx := inslogger.TestContext(t)
			mc := minimock.NewController(t)

			lr := test.mocks(ctx, mc)
			err := lr.OnPulse(ctx, insolar.Pulse{PulseNumber: insolar.FirstPulseNumber}, insolar.Pulse{PulseNumber: insolar.FirstPulseNumber + 1})
			require.NoError(t, err)

			mc.Wait(3 * time.Minute)
			mc.Finish()
		})
	}
}

type OnPulseCallOrderEnum int

const (
	OrderInitial OnPulseCallOrderEnum = iota
	OrderWriteControllerClose
	OrderResultsMatcherClear
	OrderStateStorageOnPulse
	OrderWriteControllerOpen
	OrderMAX
)

func TestLogicRunner_OnPulse_Order(t *testing.T) {
	ctx := inslogger.TestContext(t)
	lr, err := NewLogicRunner(&configuration.LogicRunner{}, nil, nil)
	require.NoError(t, err)

	mc := minimock.NewController(t)
	defer mc.Wait(time.Second)

	orderChan := make(chan OnPulseCallOrderEnum, 6)

	lr.WriteController = writecontroller.NewWriteControllerMock(mc).
		CloseAndWaitMock.Set(
		func(_ context.Context, _ insolar.PulseNumber) error {
			orderChan <- OrderWriteControllerClose
			return nil
		}).
		OpenMock.Set(
		func(_ context.Context, _ insolar.PulseNumber) error {
			orderChan <- OrderWriteControllerOpen
			return nil
		})
	lr.ResultsMatcher = NewResultMatcherMock(mc).
		ClearMock.Set(
		func() {
			orderChan <- OrderResultsMatcherClear
		})
	lr.StateStorage = NewStateStorageMock(mc).
		OnPulseMock.Set(
		func(_ context.Context, _ insolar.Pulse) []insolar.Message {
			orderChan <- OrderStateStorageOnPulse
			return []insolar.Message{}
		}).
		LockMock.Return().
		UnlockMock.Return().
		IsEmptyMock.Return(true)

	oldPulse := insolar.Pulse{PulseNumber: insolar.FirstPulseNumber}
	newPulse := insolar.Pulse{PulseNumber: insolar.FirstPulseNumber + 1}
	require.NoError(t, lr.OnPulse(ctx, oldPulse, newPulse))
	require.Len(t, orderChan, int(OrderMAX-1))

	previousOrderElement := OrderInitial
	for {
		var orderElement OnPulseCallOrderEnum
		select {
		case orderElement = <-orderChan:
			if orderElement <= previousOrderElement {
				t.Fatalf("Wrong execution order of OnPulse")
			}
			previousOrderElement = orderElement
		default:
			return
		}
	}
}
