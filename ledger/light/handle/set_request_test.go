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
	f := flow.NewFlowMock(t)
	f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
		switch p.(type) {
		case *proc.CalculateID:
			return nil
		default:
			panic("unknown procedure")
		}
	})

	request := payload.SetRequest{
		Request: []byte{1, 2, 3, 4, 5},
	}
	buf, err := request.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		// This buf is not wrapped as virtual record.
		Payload: buf,
	}

	handler := NewSetRequest(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestSetRequest_IncorrectRecordInVirtual(t *testing.T) {
	t.Parallel()

	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)
	f := flow.NewFlowMock(t)
	f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
		switch p.(type) {
		case *proc.CalculateID:
			return nil
		default:
			panic("unknown procedure")
		}
	})

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

	handler := NewSetRequest(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestSetRequest_EmptyRequestObject(t *testing.T) {
	t.Parallel()

	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)
	f := flow.NewFlowMock(t)
	f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
		switch p.(type) {
		case *proc.CalculateID:
			return nil
		default:
			panic("unknown procedure")
		}
	})

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

	handler := NewSetRequest(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	assert.EqualError(t, err, "object is nil")
}

func TestSetRequest_FlowWithPassedFlag(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := metaRequestMsg(t)

	t.Run("checkjet procedure returns unknown err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return nil
			case *proc.CheckJet:
				return errors.New("something strange from checkjet")
			default:
				panic("unknown procedure")
			}
		})

		handler := NewSetRequest(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "something strange from checkjet")
	})

	t.Run("passed flag is false and checkjet returns ErrNotExecutor", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return nil
			case *proc.CheckJet:
				return proc.ErrNotExecutor
			default:
				panic("unknown procedure")
			}
		})

		handler := NewSetRequest(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})

	t.Run("passed flag is true and checkjet returns ErrNotExecutor", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return nil
			case *proc.CheckJet:
				return proc.ErrNotExecutor
			case *proc.GetIndexWM:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := NewSetRequest(proc.NewDependenciesMock(), msg, true)
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

	msg := metaRequestMsg(t)

	t.Run("waithot procedure returns err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return nil
			case *proc.CheckJet:
				return nil
			case *proc.WaitHotWM:
				return errors.New("error from waithot")
			default:
				panic("unknown procedure")
			}
		})

		handler := NewSetRequest(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "error from waithot")
	})

	t.Run("waithot procedure returns nil err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return nil
			case *proc.CheckJet:
				return nil
			case *proc.WaitHotWM:
				return nil
			case *proc.SetRequest:
				return nil
			case *proc.GetIndexWM:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := NewSetRequest(proc.NewDependenciesMock(), msg, false)
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

	msg := metaRequestMsg(t)

	t.Run("setrequest procedure returns err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return nil
			case *proc.CheckJet:
				return nil
			case *proc.WaitHotWM:
				return nil
			case *proc.SetRequest:
				return errors.New("error from setrequest")
			case *proc.GetIndexWM:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := NewSetRequest(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "error from setrequest")
	})

	t.Run("setrequest procedure returns nil err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return nil
			case *proc.CheckJet:
				return nil
			case *proc.WaitHotWM:
				return nil
			case *proc.SetRequest:
				return nil
			case *proc.GetIndexWM:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := NewSetRequest(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func metaRequestMsg(t *testing.T) payload.Meta {
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
