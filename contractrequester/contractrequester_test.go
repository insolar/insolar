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
	"runtime"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
)

func mockMessageBus(t *testing.T, result insolar.Reply) *testutils.MessageBusMock {
	mbMock := testutils.NewMessageBusMock(t)
	mbMock.SendFunc = func(c context.Context, m insolar.Message, o *insolar.MessageSendOptions) (r insolar.Reply, r1 error) {
		return result, nil
	}
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

func mockPulseAccessor(t *testing.T) pulse.Accessor {
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

func mockJetCoordinator(t *testing.T) jet.Coordinator {
	coordinator := jet.NewCoordinatorMock(t)
	coordinator.MeFunc = func() (r insolar.Reference) {
		return testutils.RandomRef()
	}
	return coordinator
}

func TestContractRequester_SendRequest(t *testing.T) {
	ctx := inslogger.TestContext(t)
	ref := testutils.RandomRef()

	mbm := mockMessageBus(t, &reply.RegisterRequest{})
	cReq, err := New()
	assert.NoError(t, err)
	cReq.MessageBus = mbm

	cReq.PulseAccessor = mockPulseAccessor(t)

	cReq.JetCoordinator = mockJetCoordinator(t)
	cReq.PCS = platformpolicy.NewPlatformCryptographyScheme()

	mbm.MustRegisterMock.Return()
	cReq.Start(ctx)

	go func() {
		resLen := 0
		for resLen == 0 {
			cReq.ResultMutex.Lock()
			resLen = len(cReq.ResultMap)
			cReq.ResultMutex.Unlock()
			runtime.Gosched()
		}

		cReq.ResultMutex.Lock()
		for _, v := range cReq.ResultMap {
			v <- &message.ReturnResults{
				Reply: &reply.CallMethod{},
			}
		}
		cReq.ResultMutex.Unlock()
	}()
	result, err := cReq.SendRequest(ctx, &ref, "TestMethod", []interface{}{})

	require.NoError(t, err)
	require.Equal(t, &reply.CallMethod{}, result)
}

func TestContractRequester_SendRequest_RouteError(t *testing.T) {
	ctx := inslogger.TestContext(t)
	ref := testutils.RandomRef()

	mbm := mockMessageBus(t, &reply.CallMethod{})
	cReq, err := New()
	assert.NoError(t, err)
	cReq.MessageBus = mbm
	cReq.PulseAccessor = mockPulseAccessor(t)
	cReq.JetCoordinator = mockJetCoordinator(t)
	cReq.PCS = platformpolicy.NewPlatformCryptographyScheme()

	mbm.MustRegisterMock.Return()
	err = cReq.Start(ctx)
	require.NoError(t, err)

	ifResultMapEmpty := func() bool {
		cReq.ResultMutex.Lock()
		defer cReq.ResultMutex.Unlock()
		return len(cReq.ResultMap) == 0
	}

	go func() {
		for ifResultMapEmpty() {
			runtime.Gosched()
		}
		cReq.ResultMutex.Lock()
		for _, v := range cReq.ResultMap {
			v <- &message.ReturnResults{
				Reply: nil,
			}
		}
		cReq.ResultMutex.Unlock()
	}()

	result, err := cReq.SendRequest(ctx, &ref, "TestMethod", []interface{}{})
	require.Nil(t, result)
}

func TestCallMethodCanceled(t *testing.T) {
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second)
	defer cancelFunc()

	cr, err := New()
	require.NoError(t, err)

	mc := minimock.NewController(t)
	defer mc.Finish()

	mb := testutils.NewMessageBusMock(mc)
	cr.MessageBus = mb
	cr.PulseAccessor = mockPulseAccessor(t)
	cr.JetCoordinator = mockJetCoordinator(t)
	cr.PCS = platformpolicy.NewPlatformCryptographyScheme()

	ref := testutils.RandomRef()
	prototypeRef := testutils.RandomRef()
	method := testutils.RandomString()

	var reqHash [insolar.RecordHashSize]byte

	mb.SendFunc = func(p context.Context, p1 insolar.Message, p2 *insolar.MessageSendOptions) (r insolar.Reply, r1 error) {
		for k := range cr.ResultMap {
			reqHash = k
			break
		}
		return &reply.RegisterRequest{}, nil
	}

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
	assert.Contains(t, err.Error(), "canceled")

	_, ok := cr.ResultMap[reqHash]
	assert.Equal(t, false, ok)
}

func TestCallMethodWaitResults(t *testing.T) {
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*10)
	defer cancelFunc()

	cr, err := New()
	require.NoError(t, err)

	mc := minimock.NewController(t)
	defer mc.Finish()

	mb := testutils.NewMessageBusMock(mc)
	cr.MessageBus = mb
	cr.PulseAccessor = mockPulseAccessor(t)
	cr.JetCoordinator = mockJetCoordinator(t)
	cr.PCS = platformpolicy.NewPlatformCryptographyScheme()

	ref := testutils.RandomRef()
	prototypeRef := testutils.RandomRef()
	method := testutils.RandomString()

	var reqHash [insolar.RecordHashSize]byte

	mb.SendFunc = func(p context.Context, p1 insolar.Message, p2 *insolar.MessageSendOptions) (r insolar.Reply, r1 error) {
		for k := range cr.ResultMap {
			reqHash = k
			break
		}

		go func() {
			_, ok := p1.(*message.CallMethod)
			require.Equal(t, ok, true)

			cr.ResultMutex.Lock()
			defer cr.ResultMutex.Unlock()
			resChan, ok := cr.ResultMap[reqHash]
			resChan <- &message.ReturnResults{
				Reply: &reply.CallMethod{},
			}
		}()
		return &reply.RegisterRequest{}, nil
	}

	msg := &message.CallMethod{
		IncomingRequest: record.IncomingRequest{
			Object:    &ref,
			Prototype: &prototypeRef,
			Method:    method,
			Arguments: insolar.Arguments{},
		},
	}

	_, err = cr.CallMethod(ctx, msg)
	require.NoError(t, err)
}

func TestReceiveResult(t *testing.T) {
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*10)
	defer cancelFunc()

	cr, err := New()
	cr.PCS = platformpolicy.NewPlatformCryptographyScheme()
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
