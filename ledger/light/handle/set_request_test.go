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
	"context"
	"errors"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetRequest_BadMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	msg := payload.Meta{
		Payload: []byte{1, 2, 3, 4, 5},
	}

	handler := NewSetRequest(nil, msg, false)

	err := handler.Present(ctx, flow.NewFlowMock(t))
	require.Error(t, err)
}

func TestSetRequest_BadWrappedVirtualRecord(t *testing.T) {
	t.Parallel()

	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)
	f := mockProcedures(t, nil, nil, nil, nil)

	request := payload.SetRequest{
		Request: []byte{1, 2, 3, 4, 5},
	}
	buf, err := request.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		// This buf is not wrapped as virtual record.
		Payload: buf,
	}

	handler := NewSetRequest(emptyDeps(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestSetRequest_IncorrectRecordInVirtual(t *testing.T) {
	t.Parallel()

	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)
	f := mockProcedures(t, nil, nil, nil, nil)

	// Incorrect record in virtual.
	virtual := record.Virtual{
		Union: &record.Virtual_Genesis{
			Genesis: &record.Genesis{
				Hash: []byte{1, 2, 3, 4, 5},
			},
		},
	}
	virtualBuf, err := virtual.Marshal()
	require.NoError(t, err)

	request := payload.SetRequest{
		Request: virtualBuf,
	}
	requestBuf, err := request.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: requestBuf,
	}

	handler := NewSetRequest(emptyDeps(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestSetRequest_EmptyRequestObject(t *testing.T) {
	t.Parallel()

	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)
	f := mockProcedures(t, nil, nil, nil, nil)

	// Request object is nil
	virtual := record.Virtual{
		Union: &record.Virtual_Request{
			Request: &record.Request{
				Object: nil,
			},
		},
	}
	virtualBuf, err := virtual.Marshal()
	require.NoError(t, err)

	request := payload.SetRequest{
		Request: virtualBuf,
	}
	requestBuf, err := request.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: requestBuf,
	}

	handler := NewSetRequest(emptyDeps(), msg, false)

	err = handler.Present(ctx, f)
	assert.EqualError(t, err, "object is nil")
}

func TestSetRequest_FlowWithPassedFlag(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := correctMetaMsg(t)

	t.Run("checkjet procedure returns unknown err", func(t *testing.T) {
		t.Parallel()
		f := mockProcedures(t, nil, errors.New("something strange from checkjet"), nil, nil)

		handler := NewSetRequest(emptyDeps(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "something strange from checkjet")
	})

	t.Run("passed flag is false and checkjet returns ErrNotExecutor", func(t *testing.T) {
		t.Parallel()
		f := mockProcedures(t, nil, proc.ErrNotExecutor, nil, nil)

		handler := NewSetRequest(emptyDeps(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})

	t.Run("passed flag is true and checkjet returns ErrNotExecutor", func(t *testing.T) {
		t.Parallel()
		f := mockProcedures(t, nil, proc.ErrNotExecutor, nil, nil)

		handler := NewSetRequest(emptyDeps(), msg, true)
		err := handler.Present(ctx, f)
		require.Error(t, err)
		assert.Equal(t, proc.ErrNotExecutor, err)
	})
}

func TestSetRequest_ErrorFromWaitHot(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := correctMetaMsg(t)

	t.Run("waithot procedure returns err", func(t *testing.T) {
		t.Parallel()
		f := mockProcedures(t, nil, nil, errors.New("error from waithot"), nil)

		handler := NewSetRequest(emptyDeps(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "error from waithot")
	})

	t.Run("waithot procedure returns nil err", func(t *testing.T) {
		t.Parallel()
		f := mockProcedures(t, nil, nil, nil, nil)

		handler := NewSetRequest(emptyDeps(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func TestSetRequest_ErrorFromSetRequest(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := correctMetaMsg(t)

	t.Run("setrequest procedure returns err", func(t *testing.T) {
		t.Parallel()
		f := mockProcedures(t, nil, nil, nil, errors.New("error from setrequest"))

		handler := NewSetRequest(emptyDeps(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "error from setrequest")
	})

	t.Run("setrequest procedure returns nil err", func(t *testing.T) {
		t.Parallel()
		f := mockProcedures(t, nil, nil, nil, nil)

		handler := NewSetRequest(emptyDeps(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func correctMetaMsg(t *testing.T) payload.Meta {
	ref := gen.Reference()

	virtual := record.Virtual{
		Union: &record.Virtual_Request{
			Request: &record.Request{
				Object: &ref,
			},
		},
	}
	virtualBuf, err := virtual.Marshal()
	require.NoError(t, err)

	request := payload.SetRequest{
		Request: virtualBuf,
	}
	requestBuf, err := request.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: requestBuf,
	}

	return msg
}

func emptyDeps() *proc.Dependencies {
	return &proc.Dependencies{
		CalculateID: func(p *proc.CalculateID) {},
		CheckJet:    func(p *proc.CheckJet) {},
		WaitHotWM:   func(p *proc.WaitHotWM) {},
		SetRequest:  func(p *proc.SetRequest) {},
	}
}

func mockProcedures(
	t *testing.T,
	calcErr error,
	jetErr error,
	hotErr error,
	reqErr error,
) *flow.FlowMock {
	f := flow.NewFlowMock(t)
	f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
		switch p.(type) {
		case *proc.CalculateID:
			return calcErr
		case *proc.CheckJet:
			return jetErr
		case *proc.WaitHotWM:
			return hotErr
		case *proc.SetRequest:
			return reqErr
		default:
			panic("unknown procedure")
		}
	})
	return f
}
