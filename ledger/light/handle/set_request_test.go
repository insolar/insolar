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

package handle

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func TestSetRequest_BadMsgPayload(t *testing.T) {
	ctx := inslogger.TestContext(t)
	msg := payload.Meta{
		Payload: []byte{1, 2, 3, 4, 5},
	}

	handler := NewSetRequest(nil, msg, false)

	err := handler.Present(ctx, flow.NewFlowMock(t))
	require.Error(t, err)
}

func TestSetRequest_BadWrappedVirtualRecord(t *testing.T) {
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)
	pcs := testutils.NewPlatformCryptographyScheme()

	f := flow.NewFlowMock(t)
	f.ProcedureMock.Return(nil)

	request := payload.SetRequest{
		Request: []byte{1, 2, 3, 4, 5},
	}
	buf, err := request.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		// this buf is not wrapped as virtual record
		Payload: buf,
	}
	dep := &proc.Dependencies{
		CalculateID: func(p *proc.CalculateID) {
			p.Dep(pcs)
		},
	}

	handler := NewSetRequest(dep, msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestSetRequest_IncorrectRecordInVirtual(t *testing.T) {
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)
	pcs := testutils.NewPlatformCryptographyScheme()

	f := flow.NewFlowMock(t)
	f.ProcedureMock.Return(nil)

	// incorrect record in virtual
	virtual := record.Virtual{
		Union: &record.Virtual_Genesis{
			Genesis: &record.Genesis{
				Hash: []byte{1, 2, 3, 4, 5},
			},
		},
	}
	virtualBufbuf, err := virtual.Marshal()
	require.NoError(t, err)

	request := payload.SetRequest{
		Request: virtualBufbuf,
	}
	requestBuf, err := request.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: requestBuf,
	}
	dep := &proc.Dependencies{
		CalculateID: func(p *proc.CalculateID) {
			p.Dep(pcs)
		},
	}

	handler := NewSetRequest(dep, msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}
