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

package logicrunner

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/go-actors/actor/system"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/executionregistry"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
	"github.com/insolar/insolar/testutils"
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
			name:       "success",
			transcript: &common.Transcript{},
			req:        rpctypes.UpGetCodeReq{Code: gen.Reference()},
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
			name:       "no code descriptor",
			transcript: &common.Transcript{},
			req:        rpctypes.UpGetCodeReq{Code: gen.Reference()},
			dc: artifacts.NewDescriptorsCacheMock(mc).
				GetCodeMock.Return(nil, errors.New("some")),
			error: true,
		},
		{
			name:       "no code",
			transcript: &common.Transcript{},
			req:        rpctypes.UpGetCodeReq{Code: gen.Reference()},
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
			req:    rpctypes.UpRouteReq{Wait: true, Object: objRef1, Prototype: protoRef1},
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
	t.Parallel()
	mc := minimock.NewController(t)
	defer mc.Finish()

	am := artifacts.NewClientMock(mc)
	dc := artifacts.NewDescriptorsCacheMock(mc)
	cr := testutils.NewContractRequesterMock(mc)
	as := system.New()
	os := NewOutgoingRequestSender(as, cr, am)

	requestRef := gen.Reference()

	rpcm := NewExecutionProxyImplementation(dc, cr, am, os)
	ctx := context.Background()
	transcript := common.NewTranscript(ctx, requestRef, record.IncomingRequest{})
	req := rpctypes.UpRouteReq{Wait: true}
	resp := &rpctypes.UpRouteResp{}

	var outreq *record.OutgoingRequest
	outgoingReqID := gen.ID()
	outgoingReqRef := insolar.NewReference(outgoingReqID)
	// Make sure an outgoing request is registered
	am.RegisterOutgoingRequestMock.Set(func(ctx context.Context, r *record.OutgoingRequest) (*payload.RequestInfo, error) {
		require.Nil(t, outreq)
		require.Equal(t, record.ReturnResult, r.ReturnMode)
		outreq = r
		id := outgoingReqID
		return &payload.RequestInfo{RequestID: id}, nil
	})

	ref := gen.Reference()
	cr.CallMock.Return(&reply.CallMethod{}, &ref, nil)
	// Make sure the result of the outgoing request is registered as well
	am.RegisterResultMock.Set(func(ctx context.Context, reqref insolar.Reference, result artifacts.RequestResult) (r error) {
		require.Equal(t, outgoingReqRef, &reqref)
		return nil
	})

	err := rpcm.RouteCall(ctx, transcript, req, resp)
	require.NoError(t, err)
	require.NotNil(t, outreq)
	require.Equal(t, requestRef, outreq.Reason)
	as.CloseAll()
}

func TestRouteCallRegistersSaga(t *testing.T) {
	t.Parallel()

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

func TestSaveAsChildRegistersOutgoingRequestWithValidReason(t *testing.T) {
	t.Parallel()

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
		*insolar.Reference, insolar.Arguments, *record.IncomingRequest, error) {
		require.Nil(t, sentReq)
		sentReq = req
		var newObjectRef = gen.Reference()
		return &newObjectRef, []byte{3, 2, 1}, &record.IncomingRequest{}, nil
	})

	err := rpcm.SaveAsChild(ctx, transcript, req, resp)
	require.NoError(t, err)
	require.NotNil(t, registeredReq)
	require.Equal(t, requestRef, registeredReq.Reason)
	require.NotNil(t, sentReq)
	require.Equal(t, registeredReq, sentReq)
}
