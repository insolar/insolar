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

package contractrequester

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
)

func mockMessageBus(t minimock.Tester, result insolar.Reply) *testutils.MessageBusMock {
	mbMock := testutils.NewMessageBusMock(t)
	mbMock.SendMock.Return(result, nil)
	return mbMock
}

func TestNew(t *testing.T) {
	messageBus := mockMessageBus(t, nil)
	pulseAccessor := pulse.NewAccessorMock(t)
	jetCoordinator := jet.NewCoordinatorMock(t)
	pcs := platformpolicy.NewPlatformCryptographyScheme()

	contractRequester, err := New()

	cm := &component.Manager{}
	cm.Inject(messageBus, contractRequester, pulseAccessor, jetCoordinator, pcs)

	require.NoError(t, err)
	require.Equal(t, messageBus, contractRequester.MessageBus)
}

func mockPulseAccessor(t minimock.Tester) pulse.Accessor {
	pulseAccessor := pulse.NewAccessorMock(t)
	currentPulse := insolar.FirstPulseNumber
	pulseAccessor.LatestFunc = func(p context.Context) (r insolar.Pulse, r1 error) {
		return insolar.Pulse{
			PulseNumber:     insolar.PulseNumber(currentPulse),
			NextPulseNumber: insolar.PulseNumber(currentPulse + 1),
		}, nil
	}

	return pulseAccessor
}

func mockJetCoordinator(t minimock.Tester) jet.Coordinator {
	coordinator := jet.NewCoordinatorMock(t)
	coordinator.MeFunc = func() (r insolar.Reference) {
		return testutils.RandomRef()
	}
	return coordinator
}

func TestContractRequester_Start(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	cReq, err := New()
	require.NoError(t, err)

	cReq.MessageBus = testutils.NewMessageBusMock(mc).
		MustRegisterMock.Return()

	err = cReq.Start(ctx)
	require.NoError(t, err)
}

func TestContractRequester_SendRequest(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	ref := gen.Reference()

	cReq, err := New()
	require.NoError(t, err)

	cReq.PulseAccessor = mockPulseAccessor(mc)
	cReq.JetCoordinator = mockJetCoordinator(mc)
	cReq.PlatformCryptographyScheme = testutils.NewPlatformCryptographyScheme()

	table := []struct {
		name          string
		resultMessage message.ReturnResults
		earlyResult   bool
	}{
		{
			name: "success",
			resultMessage: message.ReturnResults{
				Reply: &reply.CallMethod{},
			},
		},
		{
			name: "early result, before registration",
			resultMessage: message.ReturnResults{
				Reply: &reply.CallMethod{},
			},
		},
	}

	for _, test := range table {
		test := test
		t.Run(test.name, func(t *testing.T) {

			cReq.MessageBus = testutils.NewMessageBusMock(mc).SendMock.
				Set(func(ctx context.Context, m insolar.Message, opt *insolar.MessageSendOptions) (insolar.Reply, error) {
					request := m.(*message.CallMethod).IncomingRequest

					hash, err := cReq.calcRequestHash(request)
					require.NoError(t, err)
					requestRef := insolar.NewReference(*insolar.NewID(insolar.FirstPulseNumber, hash[:]))

					resultSender := func () {
						res := test.resultMessage
						res.RequestRef = *requestRef
						cReq.result(ctx, &res)
					}
					if test.earlyResult {
						resultSender()
					} else {
						go func() {
							time.Sleep(time.Millisecond)
							resultSender()
						}()
					}

					return &reply.RegisterRequest{Request: *requestRef}, nil
				})

			result, err := cReq.SendRequest(ctx, &ref, "TestMethod", []interface{}{})
			require.NoError(t, err)
			require.Equal(t, &reply.CallMethod{}, result)
		})
	}
}

func TestContractRequester_CallMethod_Timeout(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	cr, err := New()
	require.NoError(t, err)

	cr.callTimeout = 1*time.Nanosecond

	cr.PlatformCryptographyScheme = testutils.NewPlatformCryptographyScheme()

	cr.PulseAccessor = mockPulseAccessor(mc)
	cr.JetCoordinator = mockJetCoordinator(mc)

	ref := testutils.RandomRef()
	prototypeRef := testutils.RandomRef()
	method := testutils.RandomString()

	mb := testutils.NewMessageBusMock(mc)
	mb.SendFunc = func(ctx context.Context, m insolar.Message, opt *insolar.MessageSendOptions) (insolar.Reply, error) {
		request := m.(*message.CallMethod).IncomingRequest

		hash, err := cr.calcRequestHash(request)
		require.NoError(t, err)
		requestRef := insolar.NewReference(*insolar.NewID(insolar.FirstPulseNumber, hash[:]))

		return &reply.RegisterRequest{Request: *requestRef}, nil
	}
	cr.MessageBus = mb

	msg := &message.CallMethod{
		IncomingRequest: record.IncomingRequest{
			Object:    &ref,
			Prototype: &prototypeRef,
			Method:    method,
			Arguments: insolar.Arguments{},
		},
	}

	_, err = cr.CallMethod(ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "canceled")
	require.Contains(t, err.Error(), "timeout")
}

func TestReceiveResult(t *testing.T) {
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*10)
	defer cancelFunc()

	cr, err := New()
	require.NoError(t, err)

	mc := minimock.NewController(t)
	defer mc.Finish()

	reqRef := testutils.RandomRef()
	var reqHash [insolar.RecordHashSize]byte
	copy(reqHash[:], reqRef.Record().Hash())

	msg := &message.ReturnResults{RequestRef: reqRef}
	parcel := testutils.NewParcelMock(mc).MessageMock.Return(
		msg,
	)

	// unexpected result
	rep, err := cr.ReceiveResult(ctx, parcel)
	require.NoError(t, err)
	require.Equal(t, &reply.OK{}, rep)

	// expected result
	resChan := make(chan *message.ReturnResults)
	chanResult := make(chan *message.ReturnResults)
	cr.ResultMap[reqHash] = resChan

	go func() {
		chanResult <- <-cr.ResultMap[reqHash]
	}()

	rep, err = cr.ReceiveResult(ctx, parcel)

	require.NoError(t, err)
	require.Equal(t, &reply.OK{}, rep)
	require.Equal(t, 0, len(cr.ResultMap))
	require.Equal(t, msg, <-chanResult)
}
