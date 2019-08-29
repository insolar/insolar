package logicrunner

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/gojuno/minimock"

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

	cr.CallMock.Return(&reply.CallMethod{}, insolar.NewEmptyReference(), nil)
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

	cr.CallMock.Return(&reply.CallMethod{}, insolar.NewEmptyReference(), nil)
	am.RegisterResultMock.Return(nil)

	_, err := sender.Receive(msg)
	require.NoError(t, err)

	res := <-resultChan
	require.NoError(t, res.err)
	checkIncomingAndOutgoingMatch(t, res.incoming, outgoing)
	require.Equal(t, record.ReturnNoWait, res.incoming.ReturnMode)
}

func TestOutgoingSenderSendAbandonedOutgoing(t *testing.T) {
	t.Parallel()

	mc := minimock.NewController(t)
	defer mc.Wait(2 * time.Minute)

	cr := testutils.NewContractRequesterMock(mc)
	am := artifacts.NewClientMock(mc)

	pulseObject := insolar.Pulse{PulseNumber: gen.PulseNumber()}
	pa := pulse.NewAccessorMock(mc).LatestMock.Return(pulseObject, nil)

	sender := newOutgoingSenderActorState(cr, am, pa)
	outgoing := randomOutgoingRequest()
	msg := sendAbandonedOutgoingRequestMessage{
		ctx:              context.Background(),
		requestReference: gen.Reference(),
		outgoingRequest:  outgoing,
	}

	cr.CallMock.Return(&reply.CallMethod{}, insolar.NewEmptyReference(), nil)
	am.RegisterResultMock.Return(nil)

	_, err := sender.Receive(msg)
	require.NoError(t, err)
}
