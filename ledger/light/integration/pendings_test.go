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

package integration_test

import (
	"crypto/rand"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/require"
)

func Test_Pending_RequestRegistration_Incoming(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	cfg := DefaultLightConfig()
	s, err := NewServer(ctx, cfg, nil)
	require.NoError(t, err)

	// First pulse goes in storage then interrupts.
	s.Pulse(ctx)
	// Second pulse goes in storage and starts processing, including pulse change in flow dispatcher.
	s.Pulse(ctx)

	t.Run("pending was added", func(t *testing.T) {
		p, _ := setIncomingRequest(ctx, t, s, gen.ID(), gen.ID(), record.CTSaveAsChild)
		requirePayloadNotError(t, p)
		reqInfo := p.(*payload.RequestInfo)
		p = fetchPendings(ctx, t, s, reqInfo.RequestID)
		requirePayloadNotError(t, p)

		ids := p.(*payload.IDs)
		require.Equal(t, 1, len(ids.IDs))
		require.Equal(t, reqInfo.RequestID, ids.IDs[0])
	})

	t.Run("pending was added and closed", func(t *testing.T) {
		p, _ := setIncomingRequest(ctx, t, s, gen.ID(), gen.ID(), record.CTSaveAsChild)
		requirePayloadNotError(t, p)
		reqInfo := p.(*payload.RequestInfo)

		p, _ = activateObject(ctx, t, s, reqInfo.RequestID)
		requirePayloadNotError(t, p)

		p = fetchPendings(ctx, t, s, reqInfo.RequestID)

		err := p.(*payload.Error)
		require.Equal(t, insolar.ErrNoPendingRequest.Error(), err.Text)
	})

	t.Run("reason on the another object", func(t *testing.T) {
		firstObjP, _ := setIncomingRequest(ctx, t, s, gen.ID(), gen.ID(), record.CTSaveAsChild)
		requirePayloadNotError(t, firstObjP)
		reqInfo := firstObjP.(*payload.RequestInfo)
		firstObjP, _ = activateObject(ctx, t, s, reqInfo.RequestID)
		requirePayloadNotError(t, firstObjP)

		secondObjP, _ := setIncomingRequest(ctx, t, s, gen.ID(), reqInfo.RequestID, record.CTSaveAsChild)
		requirePayloadNotError(t, secondObjP)
		secondReqInfo := secondObjP.(*payload.RequestInfo)
		secondPendings := fetchPendings(ctx, t, s, secondReqInfo.RequestID)
		requirePayloadNotError(t, secondPendings)

		ids := secondPendings.(*payload.IDs)
		require.Equal(t, 1, len(ids.IDs))
		require.Equal(t, secondReqInfo.RequestID, ids.IDs[0])
	})

	t.Run("try to register request twice. no result", func(t *testing.T) {
		args := make([]byte, 100)
		_, err := rand.Read(args)
		initReq := record.IncomingRequest{
			Object:    insolar.NewReference(gen.ID()),
			Arguments: args,
			CallType:  record.CTSaveAsChild,
			Reason:    *insolar.NewReference(*insolar.NewID(s.pulse.PulseNumber, []byte{1, 2, 3})),
			APINode:   gen.Reference(),
		}
		initReqMsg := &payload.SetIncomingRequest{
			Request: record.Wrap(initReq),
		}

		// Set first request
		p := setRequest(ctx, t, s, initReqMsg)
		requirePayloadNotError(t, p)
		reqInfo := p.(*payload.RequestInfo)
		require.Nil(t, reqInfo.Request)
		require.Nil(t, reqInfo.Result)

		// Try to set it again
		secondP := setRequest(ctx, t, s, initReqMsg)
		requirePayloadNotError(t, secondP)
		reqInfo = secondP.(*payload.RequestInfo)
		require.NotNil(t, reqInfo.Request)
		require.Nil(t, reqInfo.Result)

		// Check for the result
		compositeRec := record.CompositeFilamentRecord{}
		err = compositeRec.Unmarshal(reqInfo.Request)
		require.NoError(t, err)
		returnedReq := record.Unwrap(compositeRec.Record.Virtual)
		require.Equal(t, &initReq, returnedReq)
	})

	t.Run("try to register request twice. when there is result", func(t *testing.T) {
		args := make([]byte, 100)
		_, err := rand.Read(args)
		initReq := record.IncomingRequest{
			Object:    insolar.NewReference(gen.ID()),
			Arguments: args,
			CallType:  record.CTSaveAsChild,
			Reason:    *insolar.NewReference(*insolar.NewID(s.pulse.PulseNumber, []byte{1, 2, 3})),
			APINode:   gen.Reference(),
		}
		initReqMsg := &payload.SetIncomingRequest{
			Request: record.Wrap(initReq),
		}

		// Set first request
		p := setRequest(ctx, t, s, initReqMsg)
		requirePayloadNotError(t, p)
		reqInfo := p.(*payload.RequestInfo)

		// Set result for the request
		p, _ = activateObject(ctx, t, s, reqInfo.RequestID)
		requirePayloadNotError(t, p)

		// Try to set it again
		secondP := setRequest(ctx, t, s, initReqMsg)
		requirePayloadNotError(t, secondP)
		secondReqInfo := secondP.(*payload.RequestInfo)
		require.NotNil(t, secondReqInfo.Request)
		require.NotNil(t, secondReqInfo.Result)

		// Check for the request
		compositeReq := record.CompositeFilamentRecord{}
		err = compositeReq.Unmarshal(secondReqInfo.Request)
		require.NoError(t, err)
		returnedReq := record.Unwrap(compositeReq.Record.Virtual)
		require.Equal(t, &initReq, returnedReq)

		// Check for the result
		compositeRes := record.CompositeFilamentRecord{}
		err = compositeRes.Unmarshal(secondReqInfo.Result)
		require.NoError(t, err)
		returnedRes := record.Unwrap(compositeRes.Record.Virtual).(*record.Result)
		require.Equal(t, *insolar.NewReference(reqInfo.RequestID), returnedRes.Request)
		require.Equal(t, reqInfo.RequestID, returnedRes.Object)
	})
}
