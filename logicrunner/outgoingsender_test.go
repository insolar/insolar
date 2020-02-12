// Copyright 2020 Insolar Network Ltd.
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

package logicrunner

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/insolar/go-actors/actor/errors"

	"github.com/gojuno/minimock/v3"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/reply"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/testutils"
)

func randomOutgoingRequest() *record.OutgoingRequest {
	object := gen.Reference()
	prototype := gen.Reference()
	arguments := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	rand.Shuffle(len(arguments), func(i, j int) {
		arguments[i], arguments[j] = arguments[j], arguments[i]
	})
	outgoing := &record.OutgoingRequest{
		CallType:        record.CTMethod,
		Caller:          gen.Reference(),
		CallerPrototype: gen.Reference(),
		Nonce:           rand.Uint64(),
		ReturnMode:      record.ReturnResult,
		Immutable:       rand.Int()&1 == 1,
		Object:          &object,
		Prototype:       &prototype,
		Method:          "RandomMethodName",
		Arguments:       arguments,
		APIRequestID:    "dummy-api-request-id",
		Reason:          gen.RecordReference(),
	}
	return outgoing
}

func checkIncomingAndOutgoingMatch(t *testing.T, incoming *record.IncomingRequest, outgoing *record.OutgoingRequest) {
	require.Equal(t, outgoing.CallType, incoming.CallType)
	require.Equal(t, outgoing.CallerPrototype, incoming.CallerPrototype)
	require.Equal(t, outgoing.Nonce, incoming.Nonce)
	require.Equal(t, outgoing.Immutable, incoming.Immutable)
	require.Equal(t, outgoing.Object, incoming.Object)
	require.Equal(t, outgoing.Prototype, incoming.Prototype)
	require.Equal(t, outgoing.Method, incoming.Method)
	require.Equal(t, outgoing.Arguments, incoming.Arguments)
	if outgoing.ReturnMode == record.ReturnSaga {
		require.Equal(t, fmt.Sprintf("%s-saga-%d", outgoing.APIRequestID, incoming.Nonce), incoming.APIRequestID)
	} else {
		require.Equal(t, outgoing.APIRequestID, incoming.APIRequestID)
	}
	require.Equal(t, outgoing.Reason, incoming.Reason)
}

func TestOutgoingSenderSendRegularOutgoing(t *testing.T) {
	if useLeakTest {
		defer testutils.LeakTester(t)
	} else {
		t.Parallel()
	}

	mc := minimock.NewController(t)
	defer mc.Wait(2 * time.Minute)

	cr := testutils.NewContractRequesterMock(mc)
	am := artifacts.NewClientMock(mc)

	pulseObject := insolar.Pulse{PulseNumber: gen.PulseNumber()}
	pa := pulse.NewAccessorMock(mc).LatestMock.Return(pulseObject, nil)

	sender := newOutgoingSenderActorState(cr, am, pa)
	resultChan := make(chan sendOutgoingResult, 1)
	outgoing := randomOutgoingRequest()
	msg := sendOutgoingRequestMessage{
		ctx:              context.Background(),
		requestReference: gen.Reference(),
		outgoingRequest:  outgoing,
		resultChan:       resultChan,
	}

	cr.SendRequestMock.Return(&reply.CallMethod{}, insolar.NewEmptyReference(), nil)
	am.RegisterResultMock.Return(nil)

	_, err := sender.Receive(msg)
	require.NoError(t, err)

	res := <-resultChan
	require.NoError(t, res.err)
	checkIncomingAndOutgoingMatch(t, res.incoming, outgoing)
	require.Equal(t, outgoing.ReturnMode, res.incoming.ReturnMode)
}

func TestOutgoingSenderSendSagaOutgoing(t *testing.T) {
	if useLeakTest {
		defer testutils.LeakTester(t)
	} else {
		t.Parallel()
	}

	mc := minimock.NewController(t)
	defer mc.Wait(2 * time.Minute)

	cr := testutils.NewContractRequesterMock(mc)
	am := artifacts.NewClientMock(mc)

	pulseObject := insolar.Pulse{PulseNumber: gen.PulseNumber()}
	pa := pulse.NewAccessorMock(mc).LatestMock.Return(pulseObject, nil)

	sender := newOutgoingSenderActorState(cr, am, pa)
	resultChan := make(chan sendOutgoingResult, 1)
	outgoing := randomOutgoingRequest()
	outgoing.ReturnMode = record.ReturnSaga

	msg := sendOutgoingRequestMessage{
		ctx:              context.Background(),
		requestReference: gen.Reference(),
		outgoingRequest:  outgoing,
		resultChan:       resultChan,
	}

	cr.SendRequestMock.Return(&reply.CallMethod{}, insolar.NewEmptyReference(), nil)
	am.RegisterResultMock.Return(nil)

	_, err := sender.Receive(msg)
	require.NoError(t, err)

	res := <-resultChan
	require.NoError(t, res.err)
	checkIncomingAndOutgoingMatch(t, res.incoming, outgoing)
	require.Equal(t, record.ReturnSaga, res.incoming.ReturnMode)
}

func TestOutgoingSenderSendAbandonedOutgoing(t *testing.T) {
	if useLeakTest {
		defer testutils.LeakTester(t)
	} else {
		t.Parallel()
	}

	mc := minimock.NewController(t)
	defer mc.Wait(2 * time.Minute)

	cr := testutils.NewContractRequesterMock(mc)
	am := artifacts.NewClientMock(mc)

	pulseObject := insolar.Pulse{PulseNumber: gen.PulseNumber()}
	pa := pulse.NewAccessorMock(mc).LatestMock.Return(pulseObject, nil)

	sender := newAbandonedSenderActorState(cr, am, pa)
	outgoing := randomOutgoingRequest()
	msg := sendAbandonedOutgoingRequestMessage{
		ctx:              context.Background(),
		requestReference: gen.RecordReference(),
		outgoingRequest:  outgoing,
	}

	cr.SendRequestMock.Return(&reply.CallMethod{}, insolar.NewEmptyReference(), nil)
	am.RegisterResultMock.Return(nil)

	_, err := sender.Receive(msg)
	require.NoError(t, err)
}

func TestOutgoingSenderStop(t *testing.T) {
	if useLeakTest {
		defer testutils.LeakTester(t)
	} else {
		t.Parallel()
	}

	mc := minimock.NewController(t)
	defer mc.Wait(2 * time.Minute)

	cr := testutils.NewContractRequesterMock(mc)
	am := artifacts.NewClientMock(mc)
	pa := pulse.NewAccessorMock(mc)

	sender := newOutgoingSenderActorState(cr, am, pa)
	resultChan := make(chan struct{}, 1)
	msg := stopRequestSenderMessage{
		resultChan: resultChan,
	}
	_, err := sender.Receive(msg)
	<-resultChan
	require.Equal(t, errors.Terminate, err)
}

func TestAbandonedSenderStop(t *testing.T) {
	if useLeakTest {
		defer testutils.LeakTester(t)
	} else {
		t.Parallel()
	}

	mc := minimock.NewController(t)
	defer mc.Wait(2 * time.Minute)

	cr := testutils.NewContractRequesterMock(mc)
	am := artifacts.NewClientMock(mc)
	pa := pulse.NewAccessorMock(mc)

	sender := newAbandonedSenderActorState(cr, am, pa)
	resultChan := make(chan struct{}, 1)
	msg := stopRequestSenderMessage{
		resultChan: resultChan,
	}
	_, err := sender.Receive(msg)
	<-resultChan
	require.Equal(t, errors.Terminate, err)
}
