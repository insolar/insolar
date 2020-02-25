// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package handle_test

import (
	"context"
	"testing"

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

func TestGetRequestInfo_NilMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	meta := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   nil,
	}

	handler := handle.NewGetRequestInfo(nil, meta)

	err := handler.Present(ctx, flow.NewFlowMock(t))
	require.Error(t, err)
}

func TestGetRequestInfo_BadMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	meta := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   []byte{1, 2, 3, 4, 5},
	}

	handler := handle.NewGetRequestInfo(nil, meta)

	err := handler.Present(ctx, flow.NewFlowMock(t))
	require.Error(t, err)
}

func TestGetRequestInfo_IncorrectTypeMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	f := flow.NewFlowMock(t)

	meta := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		// Incorrect type (SetIncomingRequest instead of GetRequestInfo).
		Payload: payload.MustMarshal(&payload.SetIncomingRequest{
			Polymorph: uint32(payload.TypeSetIncomingRequest),
			Request:   record.Virtual{},
		}),
		ID: []byte{1, 1, 1},
	}

	handler := handle.NewGetRequestInfo(proc.NewDependenciesMock(), meta)

	err := handler.Present(ctx, f)
	require.Error(t, err)
}

func TestGetRequestInfo_ErrorFromFetchJet(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload: payload.MustMarshal(&payload.GetRequestInfo{
			Polymorph: uint32(payload.TypeGetRequestInfo),
			ObjectID:  insolar.ID{},
			RequestID: insolar.ID{},
		}),
	}

	t.Run("FetchJet procedure returns unknown err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.FetchJet:
				return errors.New("something strange from FetchJet")
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewGetRequestInfo(proc.NewDependenciesMock(), msg)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "something strange from FetchJet")
	})

	t.Run("FetchJet returns ErrNotExecutor", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.FetchJet:
				return proc.ErrNotExecutor
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewGetRequestInfo(proc.NewDependenciesMock(), msg)
		err := handler.Present(ctx, f)
		require.Error(t, err)
		assert.Equal(t, proc.ErrNotExecutor, err)
	})

	// Happy path, everything is fine.
	t.Run("FetchJet procedure returns nil err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.FetchJet:
				return nil
			case *proc.WaitHot:
				return nil
			case *proc.EnsureIndex:
				return nil
			case *proc.SendRequestInfo:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewGetRequestInfo(proc.NewDependenciesMock(), msg)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func TestGetRequestInfo_ErrorFromWaitHot(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload: payload.MustMarshal(&payload.GetRequestInfo{
			Polymorph: uint32(payload.TypeGetRequestInfo),
			ObjectID:  insolar.ID{},
			RequestID: insolar.ID{},
		}),
	}

	t.Run("WaitHot procedure returns err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
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

		handler := handle.NewGetRequestInfo(proc.NewDependenciesMock(), msg)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "error from WaitHot")
	})

	// Happy path, everything is fine.
	t.Run("WaitHot procedure returns nil err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.FetchJet:
				return nil
			case *proc.WaitHot:
				return nil
			case *proc.EnsureIndex:
				return nil
			case *proc.SendRequestInfo:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewGetRequestInfo(proc.NewDependenciesMock(), msg)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func TestGetRequestInfo_ErrorFromEnsureIndex(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload: payload.MustMarshal(&payload.GetRequestInfo{
			Polymorph: uint32(payload.TypeGetRequestInfo),
			ObjectID:  insolar.ID{},
			RequestID: insolar.ID{},
		}),
	}

	t.Run("EnsureIndex procedure returns err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.FetchJet:
				return nil
			case *proc.WaitHot:
				return nil
			case *proc.EnsureIndex:
				return errors.New("error from EnsureIndex")

			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewGetRequestInfo(proc.NewDependenciesMock(), msg)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "error from EnsureIndex")
	})

	// Happy path, everything is fine.
	t.Run("EnsureIndex procedure returns nil err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.FetchJet:
				return nil
			case *proc.WaitHot:
				return nil
			case *proc.EnsureIndex:
				return nil
			case *proc.SendRequestInfo:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewGetRequestInfo(proc.NewDependenciesMock(), msg)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func TestGetRequestInfo_ErrorFromSendRequestInfo(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload: payload.MustMarshal(&payload.GetRequestInfo{
			Polymorph: uint32(payload.TypeGetRequestInfo),
			ObjectID:  insolar.ID{},
			RequestID: insolar.ID{},
		}),
	}

	t.Run("SendRequestInfo procedure returns err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.FetchJet:
				return nil
			case *proc.WaitHot:
				return nil
			case *proc.EnsureIndex:
				return nil
			case *proc.SendRequestInfo:
				return errors.New("error from SendRequestInfo")
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewGetRequestInfo(proc.NewDependenciesMock(), msg)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "error from SendRequestInfo")
	})

	// Happy path, everything is fine.
	t.Run("SendRequestInfo procedure returns nil err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.FetchJet:
				return nil
			case *proc.WaitHot:
				return nil
			case *proc.EnsureIndex:
				return nil
			case *proc.SendRequestInfo:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewGetRequestInfo(proc.NewDependenciesMock(), msg)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}
