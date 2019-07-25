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
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
	"github.com/insolar/insolar/testutils"
)

func TestRPCMethods_New(t *testing.T) {
	m := NewRPCMethods(
		artifacts.NewClientMock(t),
		artifacts.NewDescriptorsCacheMock(t),
		testutils.NewContractRequesterMock(t),
		NewStateStorageMock(t),
	)
	require.NotNil(t, m)
}

func TestRPCMethods_DeactivateObject(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	reqRef := gen.Reference()
	objRef := gen.Reference()

	tr := &Transcript{RequestRef: reqRef}

	execList := NewCurrentExecutionList()
	err := execList.SetOnce(tr)
	require.NoError(t, err)

	ss := NewStateStorageMock(t).
		GetExecutionStateMock.Set(func(ref insolar.Reference) (r ExecutionBrokerI) {
		if ref.Equal(objRef) {
			return &ExecutionBroker{currentList: execList}
		} else {
			return nil
		}
	})

	m := &RPCMethods{
		ss:        ss,
		execution: NewProxyImplementationMock(t),
	}

	m.execution.(*ProxyImplementationMock).
		DeactivateObjectMock.Return(nil).
		GetCodeMock.Return(nil).
		RouteCallMock.Return(nil).
		SaveAsChildMock.Return(nil).
		SaveAsDelegateMock.Return(nil).
		GetObjChildrenIteratorMock.Return(nil).
		GetDelegateMock.Return(nil)

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
		{
			name: "as delegate",
			f: func(baseReq rpctypes.UpBaseReq) error {
				return m.SaveAsDelegate(
					rpctypes.UpSaveAsDelegateReq{UpBaseReq: baseReq},
					&rpctypes.UpSaveAsDelegateResp{},
				)
			},
		},
		{
			name: "child iter",
			f: func(baseReq rpctypes.UpBaseReq) error {
				return m.GetObjChildrenIterator(
					rpctypes.UpGetObjChildrenIteratorReq{UpBaseReq: baseReq},
					&rpctypes.UpGetObjChildrenIteratorResp{},
				)
			},
		},
		{
			name: "get delegate",
			f: func(baseReq rpctypes.UpBaseReq) error {
				return m.GetDelegate(
					rpctypes.UpGetDelegateReq{UpBaseReq: baseReq},
					&rpctypes.UpGetDelegateResp{},
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
		transcript *Transcript
		req        rpctypes.UpGetCodeReq
		dc         artifacts.DescriptorsCache
		error      bool
		result     rpctypes.UpGetCodeResp
	}{
		{
			name:       "success",
			transcript: &Transcript{},
			req:        rpctypes.UpGetCodeReq{Code: insolar.Reference{1, 2, 3}},
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
			transcript: &Transcript{},
			req:        rpctypes.UpGetCodeReq{Code: insolar.Reference{1, 2, 3}},
			dc: artifacts.NewDescriptorsCacheMock(mc).
				GetCodeMock.Return(nil, errors.New("some")),
			error: true,
		},
		{
			name:       "no code",
			transcript: &Transcript{},
			req:        rpctypes.UpGetCodeReq{Code: insolar.Reference{1, 2, 3}},
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
		transcript *Transcript
		req        rpctypes.UpDeactivateObjectReq
		error      bool
	}{
		{
			name:       "success",
			transcript: &Transcript{},
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
		transcript *Transcript
		req        rpctypes.UpRouteReq
		error      bool
		result     rpctypes.UpRouteResp
	}{
		{
			name: "success",
			transcript: &Transcript{
				LogicContext: &insolar.LogicCallContext{},
				Request:      &record.IncomingRequest{},
				RequestRef:   reqRef1,
				OutgoingRequests: []OutgoingRequest{
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
