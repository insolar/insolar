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
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
)

func TestValidationProxyImplementation_GetCode(t *testing.T) {
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
			impl := &validationProxyImplementation{dc: test.dc}
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
				RequestRef:   &reqRef1,
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
