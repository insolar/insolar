// Copyright 2020 Insolar Network Ltd.
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

package handle_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/handle"
	"github.com/insolar/insolar/ledger/light/proc"
)

func TestSetCode_NilMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   nil,
	}

	handler := handle.NewSetCode(nil, msg, false)

	err := handler.Present(ctx, flow.NewFlowMock(t))
	require.Error(t, err)
}

func TestSetCode_BadMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	msg := payload.Meta{
		Payload: []byte{1, 2, 3, 4, 5},
	}

	handler := handle.NewSetCode(nil, msg, false)

	err := handler.Present(ctx, flow.NewFlowMock(t))
	require.Error(t, err)
}

func TestSetCode_IncorrectTypeMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	f := flow.NewFlowMock(t)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		// Incorrect type (SetIncomingRequest instead of SetCode).
		Payload: payload.MustMarshal(&payload.SetIncomingRequest{
			Polymorph: uint32(payload.TypeSetIncomingRequest),
			Request:   record.Virtual{},
		}),
	}

	handler := handle.NewSetCode(proc.NewDependenciesMock(), msg, false)

	err := handler.Present(ctx, f)
	require.Error(t, err)
}

func TestSetCode_WrongRecordField(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	t.Run("Record is nil", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)

		msg := payload.Meta{
			Polymorph: uint32(payload.TypeMeta),
			Payload: payload.MustMarshal(&payload.SetCode{
				Polymorph: uint32(payload.TypeSetCode),
				Record:    nil,
			}),
		}

		handler := handle.NewSetCode(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "empty record")
	})

	t.Run("Record is empty", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)

		msg := payload.Meta{
			Polymorph: uint32(payload.TypeMeta),
			Payload: payload.MustMarshal(&payload.SetCode{
				Polymorph: uint32(payload.TypeSetCode),
				Record:    []byte{},
			}),
		}

		handler := handle.NewSetCode(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "empty record")
	})
}

func TestSetCode_NonVirtualInRecord(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	f := flow.NewFlowMock(t)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload: payload.MustMarshal(&payload.SetCode{
			Polymorph: uint32(payload.TypeSetCode),
			Record:    []byte{1, 2, 3, 4, 5},
		}),
	}

	handler := handle.NewSetCode(proc.NewDependenciesMock(), msg, false)
	err := handler.Present(ctx, f)
	require.Error(t, err)
}

func TestSetCode_ErrorCalculateID(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := metaSetCodeMsg(t)

	t.Run("calculateID procedure returns err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return errors.New("something strange from calculateID")
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewSetCode(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "something strange from calculateID")
	})

	// Happy path, everything is fine.
	t.Run("calculateID procedure returns nil err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return nil
			case *proc.FetchJet:
				return nil
			case *proc.SetCode:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewSetCode(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func TestSetCode_FlowWithPassedFlag(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := metaSetCodeMsg(t)

	t.Run("FetchJet procedure returns unknown err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return nil
			case *proc.FetchJet:
				return errors.New("something strange from FetchJet")
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewSetCode(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "something strange from FetchJet")
	})

	t.Run("passed flag is false and FetchJet returns ErrNotExecutor", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return nil
			case *proc.FetchJet:
				return proc.ErrNotExecutor
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewSetCode(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})

	t.Run("passed flag is true and FetchJet returns ErrNotExecutor", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return nil
			case *proc.FetchJet:
				return proc.ErrNotExecutor
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewSetCode(proc.NewDependenciesMock(), msg, true)
		err := handler.Present(ctx, f)
		require.Error(t, err)
		assert.Equal(t, proc.ErrNotExecutor, err)
	})
}

func TestSetCode_ErrorFromSetCode(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := metaSetCodeMsg(t)

	t.Run("SetCode procedure returns err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return nil
			case *proc.FetchJet:
				return nil
			case *proc.SetCode:
				return errors.New("error from SetCode")
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewSetCode(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "error from SetCode")
	})

	// Happy path, everything is fine.
	t.Run("SetCode procedure returns nil err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return nil
			case *proc.FetchJet:
				return nil
			case *proc.SetCode:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewSetCode(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func metaSetCodeMsg(t *testing.T) payload.Meta {
	ref := gen.Reference()

	virtual := record.Virtual{
		Union: &record.Virtual_Code{
			Code: &record.Code{
				Request: ref,
			},
		},
	}

	virtBuf, err := virtual.Marshal()
	require.NoError(t, err)

	setCode := payload.SetCode{
		Polymorph: uint32(payload.TypeSetCode),
		Record:    virtBuf,
	}
	buf, err := setCode.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   buf,
	}

	return msg
}
