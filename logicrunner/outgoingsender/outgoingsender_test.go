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

package outgoingsender

import (
	"context"
	"math/rand"
	"testing"

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
		Reason:          gen.Reference(),
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
	require.Equal(t, outgoing.APIRequestID, incoming.APIRequestID)
	require.Equal(t, outgoing.Reason, incoming.Reason)
}

func TestOutgoingSenderSendRegularOutgoing(t *testing.T) {
	t.Parallel()

	cr := testutils.NewContractRequesterMock(t)
	am := artifacts.NewClientMock(t)

	sender := newOutgoingSenderActorState(cr, am)
	resultChan := make(chan sendOutgoingResult, 1)
	outgoing := randomOutgoingRequest()
	msg := sendOutgoingRequestMessage{
		ctx:              context.Background(),
		requestReference: gen.Reference(),
		outgoingRequest:  outgoing,
		resultChan:       resultChan,
	}

	cr.CallMock.Return(&reply.CallMethod{}, nil)
	am.RegisterResultMock.Return(nil)

	_, err := sender.Receive(msg)
	require.NoError(t, err)

	res := <-resultChan
	require.NoError(t, res.err)
	checkIncomingAndOutgoingMatch(t, res.incoming, outgoing)
	require.Equal(t, outgoing.ReturnMode, res.incoming.ReturnMode)
}

// Special case: outgoing request is marked with ReturnMode = ReturnSaga.
// A corresponding incoming request is not an exact copy of the outgoing request,
// it has ReturnMode = ReturnNoWait.
func TestOutgoingSenderSendSagaOutgoing(t *testing.T) {
	t.Parallel()

	cr := testutils.NewContractRequesterMock(t)
	am := artifacts.NewClientMock(t)

	sender := newOutgoingSenderActorState(cr, am)
	resultChan := make(chan sendOutgoingResult, 1)
	outgoing := randomOutgoingRequest()
	outgoing.ReturnMode = record.ReturnSaga

	msg := sendOutgoingRequestMessage{
		ctx:              context.Background(),
		requestReference: gen.Reference(),
		outgoingRequest:  outgoing,
		resultChan:       resultChan,
	}

	cr.CallMock.Return(&reply.CallMethod{}, nil)
	am.RegisterResultMock.Return(nil)

	_, err := sender.Receive(msg)
	require.NoError(t, err)

	res := <-resultChan
	require.NoError(t, res.err)
	checkIncomingAndOutgoingMatch(t, res.incoming, outgoing)
	require.Equal(t, record.ReturnNoWait, res.incoming.ReturnMode, "ReturnMode is no ReturnNoWait")
}

func TestOutgoingSenderSendAbandonedOutgoing(t *testing.T) {
	t.Parallel()

	cr := testutils.NewContractRequesterMock(t)
	am := artifacts.NewClientMock(t)

	sender := newOutgoingSenderActorState(cr, am)
	outgoing := randomOutgoingRequest()
	msg := sendAbandonedOutgoingRequestMessage{
		ctx:              context.Background(),
		requestReference: gen.Reference(),
		outgoingRequest:  outgoing,
	}

	cr.CallMock.Return(&reply.CallMethod{}, nil)
	am.RegisterResultMock.Return(nil)

	_, err := sender.Receive(msg)
	require.NoError(t, err)
}
