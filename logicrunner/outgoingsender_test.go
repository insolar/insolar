package logicrunner

import (
	"context"
	"testing"

	"github.com/insolar/insolar/insolar/reply"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/testutils"
)

func TestOutgoingSenderSendRegularOutgoing(t *testing.T) {
	t.Parallel()

	cr := testutils.NewContractRequesterMock(t)
	am := artifacts.NewClientMock(t)

	sender := newOutgoingSenderActorState(cr, am)
	resultChan := make(chan sendOutgoingResult, 1)
	req := &record.OutgoingRequest{
		Method: "TestOutgoingSenderSendRegularOutgoing",
	}
	msg := sendOutgoingRequestMessage{
		ctx:              context.Background(),
		requestReference: gen.Reference(),
		outgoingRequest:  req,
		resultChan:       resultChan,
	}

	cr.CallMethodMock.Return(&reply.CallMethod{}, nil)
	am.RegisterResultMock.Return(nil)

	_, err := sender.Receive(msg)
	require.NoError(t, err)

	res := <-resultChan
	require.NoError(t, res.err)
	require.Equal(t, res.incoming.Method, "TestOutgoingSenderSendRegularOutgoing")
}

func TestOutgoingSenderSendAbandonedOutgoing(t *testing.T) {
	t.Parallel()

	cr := testutils.NewContractRequesterMock(t)
	am := artifacts.NewClientMock(t)

	sender := newOutgoingSenderActorState(cr, am)
	req := &record.OutgoingRequest{
		Method: "TestOutgoingSenderSendAbandonedOutgoing",
	}
	msg := sendAbandonedOutgoingRequestMessage{
		ctx:              context.Background(),
		requestReference: gen.Reference(),
		outgoingRequest:  req,
	}

	cr.CallMethodMock.Return(&reply.CallMethod{}, nil)
	am.RegisterResultMock.Return(nil)

	_, err := sender.Receive(msg)
	require.NoError(t, err)
}
