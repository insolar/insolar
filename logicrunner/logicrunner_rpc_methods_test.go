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
	transcript := NewTranscript(ctx, requestRef, record.IncomingRequest{})
	req := rpctypes.UpRouteReq{Wait: true}
	resp := &rpctypes.UpRouteResp{}

	var outreq *record.OutgoingRequest
	outgoingReqID := gen.ID()
	outgoingReqRef := insolar.NewReference(outgoingReqID)
	// Make sure an outgoing request is registered
	am.RegisterOutgoingRequestFunc = func(ctx context.Context, r *record.OutgoingRequest) (*insolar.ID, error) {
		require.Nil(t, outreq)
		require.Equal(t, record.ReturnResult, r.ReturnMode)
		outreq = r
		id := outgoingReqID
		return &id, nil
	}

	cr.CallMethodMock.Return(&reply.CallMethod{}, nil)
	// Make sure the result of the outgoing request is registered as well
	am.RegisterResultMock.Set(func(ctx context.Context, reqref insolar.Reference, result artifacts.RequestResult) (r error) {
		require.Equal(t, outgoingReqRef, &reqref)
		return nil
	})

	err := rpcm.RouteCall(ctx, transcript, req, resp)
	require.NoError(t, err)
	require.NotNil(t, outreq)
	require.Equal(t, requestRef, outreq.Reason)
}

func TestRouteCallRegistersSaga(t *testing.T) {
	t.Parallel()

	am := artifacts.NewClientMock(t)
	dc := artifacts.NewDescriptorsCacheMock(t)
	cr := testutils.NewContractRequesterMock(t)

	requestRef := gen.Reference()

	rpcm := NewExecutionProxyImplementation(dc, cr, am)
	ctx := context.Background()
	transcript := NewTranscript(ctx, requestRef, record.IncomingRequest{})
	req := rpctypes.UpRouteReq{Saga: true}
	resp := &rpctypes.UpRouteResp{}

	var outreq *record.OutgoingRequest
	outgoingReqID := gen.ID()
	// Make sure an outgoing request is registered
	am.RegisterOutgoingRequestFunc = func(ctx context.Context, r *record.OutgoingRequest) (*insolar.ID, error) {
		require.Nil(t, outreq)
		require.Equal(t, record.ReturnSaga, r.ReturnMode)
		outreq = r
		id := outgoingReqID
		return &id, nil
	}

	// cr.CallMethod and am.RegisterResults are NOT called

	err := rpcm.RouteCall(ctx, transcript, req, resp)
	require.NoError(t, err)
	require.NotNil(t, outreq)
	require.Equal(t, requestRef, outreq.Reason)
}

func TestSaveAsChildRegistersOutgoingRequestWithValidReason(t *testing.T) {
	t.Parallel()

	am := artifacts.NewClientMock(t)
	dc := artifacts.NewDescriptorsCacheMock(t)
	cr := testutils.NewContractRequesterMock(t)

	requestRef := gen.Reference()

	rpcm := NewExecutionProxyImplementation(dc, cr, am)
	ctx := context.Background()
	transcript := NewTranscript(ctx, requestRef, record.IncomingRequest{})
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
	am.RegisterResultMock.Set(func(ctx context.Context, reqref insolar.Reference, result artifacts.RequestResult) (r error) {
		require.Equal(t, outgoingReqRef, &reqref)
		require.Equal(t, newObjRef.Bytes(), result.Result())
		return nil
	})

	err := rpcm.SaveAsChild(ctx, transcript, req, resp)
	require.NoError(t, err)
	require.NotNil(t, outreq)
	require.Equal(t, requestRef, outreq.Reason)
}

func TestSaveAsDelegateRegistersOutgoingRequestWithValidReason(t *testing.T) {
	t.Parallel()

	am := artifacts.NewClientMock(t)
	dc := artifacts.NewDescriptorsCacheMock(t)
	cr := testutils.NewContractRequesterMock(t)

	requestRef := gen.Reference()

	rpcm := NewExecutionProxyImplementation(dc, cr, am)
	ctx := context.Background()
	transcript := NewTranscript(ctx, requestRef, record.IncomingRequest{})
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
	am.RegisterResultMock.Set(func(ctx context.Context, reqref insolar.Reference, result artifacts.RequestResult) (r error) {
		require.Equal(t, outgoingReqRef, &reqref)
		require.Equal(t, newObjRef.Bytes(), result.Result())

		return nil
	})

	err := rpcm.SaveAsDelegate(ctx, transcript, req, resp)
	require.NoError(t, err)
	require.NotNil(t, outreq)
	require.Equal(t, requestRef, outreq.Reason)
}
