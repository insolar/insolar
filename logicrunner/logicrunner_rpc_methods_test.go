package logicrunner

import (
	"context"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
	"github.com/insolar/insolar/testutils"

	"github.com/stretchr/testify/require"
)

func TestRouteCallRegistersOutgoingRequestWithValidReason(t *testing.T) {
	t.Parallel()

	am := artifacts.NewClientMock(t)
	dc := artifacts.NewDescriptorsCacheMock(t)
	cr := testutils.NewContractRequesterMock(t)

	parcel := testutils.NewParcelMock(t)
	parcel.MessageMock.Expect().Return(&message.CallMethod{})
	parcel.GetSenderMock.Expect().Return(gen.Reference())

	requestRef := gen.Reference()
	pulse := insolar.Pulse{}
	callee := gen.Reference()

	rpcm := NewExecutionProxyImplementation(dc, cr, am)
	ctx := context.Background()
	transcript := NewTranscript(ctx, parcel, &requestRef, &pulse, callee)
	reason := gen.Reference()
	transcript.RequestRef = &reason
	req := rpctypes.UpRouteReq{}
	resp := &rpctypes.UpRouteResp{}

	var outreq *record.OutgoingRequest

	am.RegisterOutgoingRequestFunc = func(ctx context.Context, r record.OutgoingRequest) (*insolar.ID, error) {
		require.Nil(t, outreq)
		outreq = &r
		id := gen.ID()
		return &id, nil
	}

	cr.CallMethodMock.Return(&reply.OK{}, nil)

	err := rpcm.RouteCall(ctx, transcript, req, resp)
	require.NoError(t, err)
	require.NotNil(t, outreq)
	require.Equal(t, reason, outreq.Reason)
}
