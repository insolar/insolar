// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package handle_test

import (
	"context"
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/handle"
	"github.com/insolar/insolar/ledger/light/proc"
)

func TestGetPendings_NilMsgPayload(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)

	ctx := inslogger.TestContext(t)
	meta := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   nil,
	}

	handler := handle.NewGetPendings(nil, meta, false)

	err := handler.Present(ctx, flow.NewFlowMock(mc))
	require.Error(t, err)
}

func TestGetPendings_BadMsgPayload(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)

	ctx := inslogger.TestContext(t)
	meta := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   []byte{1, 2, 3, 4, 5},
	}

	handler := handle.NewGetPendings(nil, meta, false)

	err := handler.Present(ctx, flow.NewFlowMock(mc))
	require.Error(t, err)
}

func TestGetPendings_IncorrectTypeMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	f := flow.NewFlowMock(mc)

	meta := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		// Incorrect type (SetIncomingRequest instead of GetPendings).
		Payload: payload.MustMarshal(&payload.SetIncomingRequest{
			Polymorph: uint32(payload.TypeSetIncomingRequest),
			Request:   record.Virtual{},
		}),
		ID: []byte{1, 1, 1},
	}

	handler := handle.NewGetPendings(proc.NewDependenciesMock(), meta, false)

	err := handler.Present(ctx, f)
	require.Error(t, err)
}

func TestGetPendings_FlowWithPassedFlag(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload: payload.MustMarshal(&payload.GetPendings{
			Polymorph: uint32(payload.TypeGetPendings),
			ObjectID:  insolar.ID{},
		}),
	}

	t.Run("FetchJet procedure returns unknown err", func(t *testing.T) {
		t.Parallel()
		mc := minimock.NewController(t)
		f := flow.NewFlowMock(mc)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.FetchJet:
				return errors.New("something strange from FetchJet")
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewGetPendings(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "something strange from FetchJet")
	})

	t.Run("passed flag is false and FetchJet returns ErrNotExecutor", func(t *testing.T) {
		t.Parallel()
		mc := minimock.NewController(t)
		f := flow.NewFlowMock(mc)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.FetchJet:
				return proc.ErrNotExecutor
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewGetPendings(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})

	t.Run("passed flag is true and FetchJet returns ErrNotExecutor", func(t *testing.T) {
		t.Parallel()
		mc := minimock.NewController(t)
		f := flow.NewFlowMock(mc)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.FetchJet:
				return proc.ErrNotExecutor
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewGetPendings(proc.NewDependenciesMock(), msg, true)
		err := handler.Present(ctx, f)
		require.Error(t, err)
		assert.Equal(t, proc.ErrNotExecutor, err)
	})
}

func TestGetPendings_ErrorFromWaitHot(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload: payload.MustMarshal(&payload.GetPendings{
			Polymorph: uint32(payload.TypeGetPendings),
			ObjectID:  insolar.ID{},
		}),
	}

	t.Run("WaitHot procedure returns err", func(t *testing.T) {
		t.Parallel()
		mc := minimock.NewController(t)
		f := flow.NewFlowMock(mc)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.FetchJet:
				return nil
			case *proc.WaitHot:
				return errors.New("error from WaitHot")
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewGetPendings(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "error from WaitHot")
	})

	// Happy path, everything is fine.
	t.Run("WaitHot procedure returns nil err", func(t *testing.T) {
		t.Parallel()
		mc := minimock.NewController(t)
		f := flow.NewFlowMock(mc)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.FetchJet:
				return nil
			case *proc.WaitHot:
				return nil
			case *proc.GetPendings:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewGetPendings(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func TestGetPendings_ErrorFromGetPendings(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload: payload.MustMarshal(&payload.GetPendings{
			Polymorph: uint32(payload.TypeGetPendings),
			ObjectID:  insolar.ID{},
		}),
	}

	t.Run("GetPendings procedure returns err", func(t *testing.T) {
		t.Parallel()
		mc := minimock.NewController(t)
		f := flow.NewFlowMock(mc)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.FetchJet:
				return nil
			case *proc.WaitHot:
				return nil
			case *proc.GetPendings:
				return errors.New("error from GetPendings")
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewGetPendings(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "error from GetPendings")
	})

	// Happy path, everything is fine.
	t.Run("GetPendings procedure returns nil err", func(t *testing.T) {
		t.Parallel()
		mc := minimock.NewController(t)
		f := flow.NewFlowMock(mc)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.FetchJet:
				return nil
			case *proc.WaitHot:
				return nil
			case *proc.GetPendings:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewGetPendings(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}
