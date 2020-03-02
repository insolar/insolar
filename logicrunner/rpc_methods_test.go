// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package logicrunner

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/insolar/go-actors/actor/system"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/executionregistry"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestRPCMethods_New(t *testing.T) {
	m := NewRPCMethods(
		artifacts.NewClientMock(t),
		artifacts.NewDescriptorsCacheMock(t),
		testutils.NewContractRequesterMock(t),
		NewStateStorageMock(t),
		NewOutgoingRequestSenderMock(t),
	)
	require.NotNil(t, m)
}

func TestRPCMethods_DeactivateObject(t *testing.T) {
	defer testutils.LeakTester(t)

	mc := minimock.NewController(t)
	defer mc.Finish()

	reqRef := gen.Reference()
	objRef := gen.Reference()

	tr := &common.Transcript{RequestRef: reqRef}

	executionRegistry := executionregistry.NewExecutionRegistryMock(mc).GetActiveTranscriptMock.Set(
		func(ref insolar.Reference) (r *common.Transcript) {
			if ref.Equal(reqRef) {
				return tr
			} else {
				return nil
			}
		},
	)

	ss := NewStateStorageMock(t).GetExecutionRegistryMock.Set(
		func(ref insolar.Reference) (r executionregistry.ExecutionRegistry) {
			if ref.Equal(objRef) {
				return executionRegistry
			} else {
				return nil
			}
		},
	)

	m := &RPCMethods{
		ss:        ss,
		execution: NewProxyImplementationMock(t),
	}

	m.execution.(*ProxyImplementationMock).
		DeactivateObjectMock.Return(nil).
		GetCodeMock.Return(nil).
		RouteCallMock.Return(nil).
		SaveAsChildMock.Return(nil)

	table := []struct {
		name string
		f    func(rpctypes.UpBaseReq) error
	}{
		{
			name: "deactivate",
			f: func(baseReq rpctypes.UpBaseReq) error {
				return m.DeactivateObject(
					rpctypes.UpDeactivateObjectReq{UpBaseReq: baseReq},
					&rpctypes.UpDeactivateObjectResp{},
				)
			},
		},
		{
			name: "code",
			f: func(baseReq rpctypes.UpBaseReq) error {
				return m.GetCode(
					rpctypes.UpGetCodeReq{UpBaseReq: baseReq},
					&rpctypes.UpGetCodeResp{},
				)
			},
		},
		{
			name: "call",
			f: func(baseReq rpctypes.UpBaseReq) error {
				return m.RouteCall(
					rpctypes.UpRouteReq{UpBaseReq: baseReq},
					&rpctypes.UpRouteResp{},
				)
			},
		},
		{
			name: "as child",
			f: func(baseReq rpctypes.UpBaseReq) error {
				return m.SaveAsChild(
					rpctypes.UpSaveAsChildReq{UpBaseReq: baseReq},
					&rpctypes.UpSaveAsChildResp{},
				)
			},
		},
	}

	for _, test := range table {
		t.Run(test.name, func(t *testing.T) {
			err := test.f(rpctypes.UpBaseReq{Callee: objRef, Request: reqRef})
			require.NoError(t, err)

			err = test.f(rpctypes.UpBaseReq{Callee: objRef, Request: gen.Reference()})
			require.Error(t, err)

			err = test.f(rpctypes.UpBaseReq{Callee: gen.Reference(), Request: reqRef})
			require.Error(t, err)
		})
	}
}

func TestProxyImplementation_GetCode(t *testing.T) {
	defer testutils.LeakTester(t)

	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	table := []struct {
		name       string
		transcript *common.Transcript
		req        rpctypes.UpGetCodeReq
		dc         artifacts.DescriptorsCache
		error      bool
		result     rpctypes.UpGetCodeResp
	}{
		{
			name: "success",
			transcript: &common.Transcript{
				Request: genIncomingRequest(),
			},
			req: rpctypes.UpGetCodeReq{Code: gen.Reference()},
			dc: artifacts.NewDescriptorsCacheMock(mc).
				GetCodeMock.
				Return(
					artifacts.NewCodeDescriptorMock(mc).
						CodeMock.Return([]byte{3, 2, 1}, nil),
					nil,
				),
			result: rpctypes.UpGetCodeResp{Code: []byte{3, 2, 1}},
		},
		{
			name: "no code descriptor",
			transcript: &common.Transcript{
				Request: genIncomingRequest(),
			},
			req: rpctypes.UpGetCodeReq{Code: gen.Reference()},
			dc: artifacts.NewDescriptorsCacheMock(mc).
				GetCodeMock.Return(nil, errors.New("some")),
			error: true,
		},
		{
			name: "no code",
			transcript: &common.Transcript{
				Request: genIncomingRequest(),
			},
			req: rpctypes.UpGetCodeReq{Code: gen.Reference()},
			dc: artifacts.NewDescriptorsCacheMock(mc).
				GetCodeMock.Return(
				artifacts.NewCodeDescriptorMock(mc).
					CodeMock.Return(nil, errors.New("some")),
				nil,
			),
			error: true,
		},
	}

	for _, test := range table {
		test := test
		t.Run(test.name, func(t *testing.T) {
			implementations := []ProxyImplementation{
				&executionProxyImplementation{dc: test.dc},
				&validationProxyImplementation{dc: test.dc},
			}
			for _, impl := range implementations {
				result := rpctypes.UpGetCodeResp{}
				err := impl.GetCode(ctx, test.transcript, test.req, &result)
				if !test.error {
					require.NoError(t, err)
					require.NotNil(t, result)
					require.Equal(t, test.result, result)
				} else {
					require.Error(t, err)
					require.Equal(t, test.result, rpctypes.UpGetCodeResp{})
				}
			}
		})
	}
}

func TestValidationProxyImplementation_DeactivateObject(t *testing.T) {
	defer testutils.LeakTester(t)

	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	table := []struct {
		name       string
		transcript *common.Transcript
		req        rpctypes.UpDeactivateObjectReq
		error      bool
	}{
		{
			name:       "success",
			transcript: &common.Transcript{},
			req:        rpctypes.UpDeactivateObjectReq{},
		},
	}

	for _, test := range table {
		test := test
		t.Run(test.name, func(t *testing.T) {
			impl := &validationProxyImplementation{}
			result := rpctypes.UpDeactivateObjectResp{}
			err := impl.DeactivateObject(ctx, test.transcript, test.req, &result)
			if !test.error {
				require.NoError(t, err)
				require.True(t, test.transcript.Deactivate)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestValidationProxyImplementation_RouteCall(t *testing.T) {
	defer testutils.LeakTester(t)

	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	objRef1 := gen.Reference()
	protoRef1 := gen.Reference()
	reqRef1 := gen.Reference()

	table := []struct {
		name       string
		transcript *common.Transcript
		req        rpctypes.UpRouteReq
		error      bool
		result     rpctypes.UpRouteResp
	}{
		{
			name: "success",
			transcript: &common.Transcript{
				LogicContext: &insolar.LogicCallContext{},
				Request:      &record.IncomingRequest{},
				RequestRef:   reqRef1,
				OutgoingRequests: []common.OutgoingRequest{
					{
						Request: record.IncomingRequest{
							Nonce: 1, Reason: reqRef1, Object: &objRef1, Prototype: &protoRef1,
						},
						Response: []byte{1, 2, 3},
					},
				},
			},
			req:    rpctypes.UpRouteReq{Object: objRef1, Prototype: protoRef1},
			result: rpctypes.UpRouteResp{Result: []byte{1, 2, 3}},
		},
	}

	for _, test := range table {
		test := test
		t.Run(test.name, func(t *testing.T) {
			impl := &validationProxyImplementation{}
			result := rpctypes.UpRouteResp{}
			err := impl.RouteCall(ctx, test.transcript, test.req, &result)
			if !test.error {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, test.result, result)
			} else {
				require.Error(t, err)
				require.Equal(t, test.result, rpctypes.UpGetCodeResp{})
			}
		})
	}
}
func TestRouteCallRegistersOutgoingRequestWithValidReason(t *testing.T) {
	if useLeakTest {
		defer testutils.LeakTester(t)
	} else {
		t.Parallel()
	}

	mc := minimock.NewController(t)
	defer mc.Finish()

	ctx := inslogger.TestContext(t)

	pulseObject := insolar.Pulse{PulseNumber: gen.PulseNumber()}
	pa := pulse.NewAccessorMock(mc).LatestMock.Return(pulseObject, nil)

	am := artifacts.NewClientMock(mc)
	dc := artifacts.NewDescriptorsCacheMock(mc)
	cr := testutils.NewContractRequesterMock(mc)
	as := system.New()
	os := NewOutgoingRequestSender(as, cr, am, pa)

	objectRef := gen.Reference()
	requestRef := gen.RecordReference()

	rpcm := NewExecutionProxyImplementation(dc, cr, am, os)
	transcript := common.NewTranscript(ctx, requestRef, record.IncomingRequest{
		Object: &objectRef,
	})
	req := rpctypes.UpRouteReq{}
	resp := &rpctypes.UpRouteResp{}

	var outreq *record.OutgoingRequest
	outgoingReqRef := gen.RecordReference()
	// Make sure an outgoing request is registered
	am.RegisterOutgoingRequestMock.Set(func(ctx context.Context, r *record.OutgoingRequest) (*payload.RequestInfo, error) {
		require.Nil(t, outreq)
		require.Equal(t, record.ReturnResult, r.ReturnMode)
		outreq = r
		id := *outgoingReqRef.GetLocal()
		return &payload.RequestInfo{RequestID: id}, nil
	})

	ref := gen.Reference()
	cr.SendRequestMock.Return(&reply.CallMethod{}, &ref, nil)
	// Make sure the result of the outgoing request is registered as well
	am.RegisterResultMock.Set(func(ctx context.Context, reqref insolar.Reference, result artifacts.RequestResult) (r error) {
		if outgoingReqRef != reqref {
			return errors.Errorf("outgoingReqRef != reqref, ref1=%s, ref2=%s", outgoingReqRef.String(), reqref.String())
		}
		return nil
	})

	err := rpcm.RouteCall(ctx, transcript, req, resp)
	require.NoError(t, err)
	require.NotNil(t, outreq)
	require.Equal(t, requestRef, outreq.Reason)

	os.Stop(ctx)
	as.AwaitTermination()
}

func TestRouteCallRegistersOutgoingRequestAlreadyHasResult(t *testing.T) {
	if useLeakTest {
		defer testutils.LeakTester(t)
	} else {
		t.Parallel()
	}

	mc := minimock.NewController(t)
	defer mc.Finish()

	ctx := inslogger.TestContext(t)

	pa := pulse.NewAccessorMock(mc)

	am := artifacts.NewClientMock(mc)
	dc := artifacts.NewDescriptorsCacheMock(mc)
	cr := testutils.NewContractRequesterMock(mc)
	as := system.New()
	os := NewOutgoingRequestSender(as, cr, am, pa)

	objectRef := gen.Reference()
	requestRef := gen.RecordReference()

	rpcm := NewExecutionProxyImplementation(dc, cr, am, os)
	transcript := common.NewTranscript(ctx, requestRef, record.IncomingRequest{
		Object: &objectRef,
	})
	req := rpctypes.UpRouteReq{}
	resp := &rpctypes.UpRouteResp{}

	var outreq *record.OutgoingRequest
	outgoingReqRef := gen.RecordReference()
	// Make sure an outgoing request is registered
	am.RegisterOutgoingRequestMock.Set(func(ctx context.Context, r *record.OutgoingRequest) (*payload.RequestInfo, error) {
		require.Nil(t, outreq)
		require.Equal(t, record.ReturnResult, r.ReturnMode)
		outreq = r
		id := *outgoingReqRef.GetLocal()
		result := append(make([]byte, 1), 1)
		resRecord := &record.Result{Payload: result}
		virtResRecord := record.Wrap(resRecord)
		matRecord := record.Material{Virtual: virtResRecord}
		matRecordSerialized, err := matRecord.Marshal()
		require.NoError(t, err)
		return &payload.RequestInfo{RequestID: id, Result: matRecordSerialized}, nil
	})

	err := rpcm.RouteCall(ctx, transcript, req, resp)
	require.NoError(t, err)
	require.NotNil(t, outreq)
	require.Equal(t, requestRef, outreq.Reason)

	os.Stop(ctx)
	as.AwaitTermination()
}

func TestRouteCallRegistersSaga(t *testing.T) {
	if useLeakTest {
		defer testutils.LeakTester(t)
	} else {
		t.Parallel()
	}

	am := artifacts.NewClientMock(t)
	dc := artifacts.NewDescriptorsCacheMock(t)
	cr := testutils.NewContractRequesterMock(t)
	os := NewOutgoingRequestSenderMock(t)

	requestRef := gen.Reference()

	rpcm := NewExecutionProxyImplementation(dc, cr, am, os)
	ctx := context.Background()
	transcript := common.NewTranscript(ctx, requestRef, record.IncomingRequest{})
	req := rpctypes.UpRouteReq{Saga: true}
	resp := &rpctypes.UpRouteResp{}

	var outreq *record.OutgoingRequest
	outgoingReqID := gen.ID()
	// Make sure an outgoing request is registered
	am.RegisterOutgoingRequestMock.Set(func(ctx context.Context, r *record.OutgoingRequest) (*payload.RequestInfo, error) {
		require.Nil(t, outreq)
		require.Equal(t, record.ReturnSaga, r.ReturnMode)
		outreq = r
		id := outgoingReqID
		return &payload.RequestInfo{RequestID: id}, nil
	})

	// cr.CallMethod and am.RegisterResults are NOT called

	err := rpcm.RouteCall(ctx, transcript, req, resp)
	require.NoError(t, err)
	require.NotNil(t, outreq)
	require.Equal(t, requestRef, outreq.Reason)
}

func TestRouteCallFailedAfterReturningResultForSaga(t *testing.T) {
	if useLeakTest {
		defer testutils.LeakTester(t)
	} else {
		t.Parallel()
	}

	am := artifacts.NewClientMock(t)
	dc := artifacts.NewDescriptorsCacheMock(t)
	cr := testutils.NewContractRequesterMock(t)
	os := NewOutgoingRequestSenderMock(t)

	requestRef := gen.Reference()

	rpcm := NewExecutionProxyImplementation(dc, cr, am, os)
	ctx := context.Background()
	transcript := common.NewTranscript(ctx, requestRef, record.IncomingRequest{})
	req := rpctypes.UpRouteReq{Saga: true}
	resp := &rpctypes.UpRouteResp{}

	var outreq *record.OutgoingRequest
	outgoingReqID := gen.ID()
	// Make sure an outgoing request is registered
	am.RegisterOutgoingRequestMock.Set(func(ctx context.Context, r *record.OutgoingRequest) (*payload.RequestInfo, error) {
		require.Nil(t, outreq)
		require.Equal(t, record.ReturnSaga, r.ReturnMode)
		outreq = r
		id := outgoingReqID
		result := append(make([]byte, 1), 1)
		return &payload.RequestInfo{RequestID: id, Result: result}, nil
	})

	err := rpcm.RouteCall(ctx, transcript, req, resp)
	require.Error(t, err)
	require.Equal(t, requestRef, outreq.Reason)
}

func TestSaveAsChildRegistersOutgoingRequestWithValidReason(t *testing.T) {
	if useLeakTest {
		defer testutils.LeakTester(t)
	} else {
		t.Parallel()
	}

	am := artifacts.NewClientMock(t)
	dc := artifacts.NewDescriptorsCacheMock(t)
	cr := testutils.NewContractRequesterMock(t)
	os := NewOutgoingRequestSenderMock(t)

	requestRef := gen.Reference()

	rpcm := NewExecutionProxyImplementation(dc, cr, am, os)
	ctx := context.Background()
	transcript := common.NewTranscript(ctx, requestRef, record.IncomingRequest{})
	req := rpctypes.UpSaveAsChildReq{}
	resp := &rpctypes.UpSaveAsChildResp{}

	// Make sure the outgoing request was registered
	var registeredReq *record.OutgoingRequest
	am.RegisterOutgoingRequestMock.Set(func(ctx context.Context, r *record.OutgoingRequest) (*payload.RequestInfo, error) {
		require.Nil(t, registeredReq)
		registeredReq = r
		id := gen.ID()
		return &payload.RequestInfo{RequestID: id}, nil
	})

	// Make sure the result of the outgoing request was sent
	var sentReq *record.OutgoingRequest
	os.SendOutgoingRequestMock.Set(func(ctx context.Context, reqRef insolar.Reference, req *record.OutgoingRequest) (
		insolar.Arguments, *record.IncomingRequest, error) {
		require.Nil(t, sentReq)
		sentReq = req
		return []byte{3, 2, 1}, &record.IncomingRequest{}, nil
	})

	err := rpcm.SaveAsChild(ctx, transcript, req, resp)
	require.NoError(t, err)
	require.NotNil(t, registeredReq)
	require.Equal(t, requestRef, registeredReq.Reason)
	require.NotNil(t, sentReq)
	require.Equal(t, registeredReq, sentReq)
}
