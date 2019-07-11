package logicrunner

import (
	"context"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
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

	requestRef := gen.Reference()

	rpcm := NewExecutionProxyImplementation(dc, cr, am)
	ctx := context.Background()
	transcript := NewTranscript(ctx, &requestRef, record.IncomingRequest{})
	reason := gen.Reference()
	transcript.RequestRef = &reason
	req := rpctypes.UpRouteReq{}
	resp := &rpctypes.UpRouteResp{}

	var outreq *record.OutgoingRequest
	outgoingReqID := gen.ID()
	outgoingReqRef := insolar.NewReference(outgoingReqID)
	// Make sure an outgoing request is registered
	am.RegisterOutgoingRequestFunc = func(ctx context.Context, r *record.OutgoingRequest) (*insolar.ID, error) {
		require.Nil(t, outreq)
		outreq = r
		id := outgoingReqID
		return &id, nil
	}

	cr.CallMethodMock.Return(&reply.OK{}, nil)
	// Make sure the result of the outgoing request is registered as well
	am.RegisterResultFunc = func(ctx context.Context, objref insolar.Reference, reqref insolar.Reference, result []byte) (r *insolar.ID, r1 error) {
		require.Equal(t, outgoingReqRef, &reqref)
		id := gen.ID()
		return &id, nil
	}

	err := rpcm.RouteCall(ctx, transcript, req, resp)
	require.NoError(t, err)
	require.NotNil(t, outreq)
	require.Equal(t, reason, outreq.Reason)
}

func TestSaveAsChildRegistersOutgoingRequestWithValidReason(t *testing.T) {
	t.Parallel()

	am := artifacts.NewClientMock(t)
	dc := artifacts.NewDescriptorsCacheMock(t)
	cr := testutils.NewContractRequesterMock(t)

	requestRef := gen.Reference()

	rpcm := NewExecutionProxyImplementation(dc, cr, am)
	ctx := context.Background()
	transcript := NewTranscript(ctx, &requestRef, record.IncomingRequest{})
	reason := gen.Reference()
	transcript.RequestRef = &reason
	req := rpctypes.UpSaveAsChildReq{}
	resp := &rpctypes.UpSaveAsChildResp{}

	var outreq *record.OutgoingRequest
	outgoingReqID := gen.ID()
	outgoingReqRef := insolar.NewReference(outgoingReqID)
	// Make sure an outgoing request is registered
	am.RegisterOutgoingRequestFunc = func(ctx context.Context, r *record.OutgoingRequest) (*insolar.ID, error) {
		require.Nil(t, outreq)
		outreq = r
		id := outgoingReqID
		return &id, nil
	}

	newObjRef := gen.Reference()
	cr.CallConstructorMock.Return(&newObjRef, nil)

	// Make sure the result of the outgoing request is registered as well
	am.RegisterResultFunc = func(ctx context.Context, objref insolar.Reference, reqref insolar.Reference, result []byte) (r *insolar.ID, r1 error) {
		require.Equal(t, outgoingReqRef, &reqref)
		require.Equal(t, newObjRef.Bytes(), result)
		id := gen.ID()
		return &id, nil
	}

	err := rpcm.SaveAsChild(ctx, transcript, req, resp)
	require.NoError(t, err)
	require.NotNil(t, outreq)
	require.Equal(t, reason, outreq.Reason)
}

func TestSaveAsDelegateRegistersOutgoingRequestWithValidReason(t *testing.T) {
	t.Parallel()

	am := artifacts.NewClientMock(t)
	dc := artifacts.NewDescriptorsCacheMock(t)
	cr := testutils.NewContractRequesterMock(t)

	requestRef := gen.Reference()

	rpcm := NewExecutionProxyImplementation(dc, cr, am)
	ctx := context.Background()
	transcript := NewTranscript(ctx, &requestRef, record.IncomingRequest{})
	reason := gen.Reference()
	transcript.RequestRef = &reason
	req := rpctypes.UpSaveAsDelegateReq{}
	resp := &rpctypes.UpSaveAsDelegateResp{}

	var outreq *record.OutgoingRequest
	outgoingReqID := gen.ID()
	outgoingReqRef := insolar.NewReference(outgoingReqID)
	// Make sure an outgoing request is registered
	am.RegisterOutgoingRequestFunc = func(ctx context.Context, r *record.OutgoingRequest) (*insolar.ID, error) {
		require.Nil(t, outreq)
		outreq = r
		id := outgoingReqID
		return &id, nil
	}

	newObjRef := gen.Reference()
	cr.CallConstructorMock.Return(&newObjRef, nil)

	// Make sure the result of the outgoing request is registered as well
	am.RegisterResultFunc = func(ctx context.Context, objref insolar.Reference, reqref insolar.Reference, result []byte) (r *insolar.ID, r1 error) {
		require.Equal(t, outgoingReqRef, &reqref)
		require.Equal(t, newObjRef.Bytes(), result)
		id := gen.ID()
		return &id, nil
	}

	err := rpcm.SaveAsDelegate(ctx, transcript, req, resp)
	require.NoError(t, err)
	require.NotNil(t, outreq)
	require.Equal(t, reason, outreq.Reason)
}
