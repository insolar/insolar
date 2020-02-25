// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package handle_test

import (
	"context"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/handle"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateObject_BadMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	msg := payload.Meta{
		Payload: []byte{1, 2, 3, 4, 5},
	}

	handler := handle.NewUpdateObject(nil, msg, false)

	err := handler.Present(ctx, flow.NewFlowMock(t))
	require.Error(t, err)
}

func TestUpdateObject_BadWrappedVirtualRecord(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	f := flow.NewFlowMock(t)

	update := payload.Update{
		Record: []byte{1, 2, 3, 4, 5},
	}
	buf, err := update.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		// This buf is not wrapped as virtual record.
		Payload: buf,
	}

	handler := handle.NewUpdateObject(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestUpdateObject_IncorrectUpdateRecordInVirtual(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	f := flow.NewFlowMock(t)

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

	update := payload.Update{
		Record: virtualBuf,
	}
	updateBuf, err := update.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: updateBuf,
	}

	handler := handle.NewUpdateObject(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestUpdateObject_IncorrectUpdateResultPayload(t *testing.T) {
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

	virtualUpdate := record.Virtual{
		Union: &record.Virtual_Amend{
			Amend: &record.Amend{},
		},
	}
	virtualUpdateBuf, err := virtualUpdate.Marshal()
	require.NoError(t, err)

	update := payload.Update{
		Record: virtualUpdateBuf,
		Result: []byte{1, 2, 3, 4, 5},
	}
	updateBuf, err := update.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: updateBuf,
	}

	handler := handle.NewUpdateObject(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestUpdateObject_EmptyUpdateResultObject(t *testing.T) {
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

	virtualUpdate := record.Virtual{
		Union: &record.Virtual_Amend{
			Amend: &record.Amend{},
		},
	}
	virtualUpdateBuf, err := virtualUpdate.Marshal()
	require.NoError(t, err)

	// Update.Result object is empty
	virtualResult := record.Virtual{
		Union: &record.Virtual_Result{
			Result: &record.Result{
				Object: insolar.ID{},
			},
		},
	}
	virtualResultBuf, err := virtualResult.Marshal()
	require.NoError(t, err)

	update := payload.Update{
		Record: virtualUpdateBuf,
		Result: virtualResultBuf,
	}
	updateBuf, err := update.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: updateBuf,
	}

	handler := handle.NewUpdateObject(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	assert.EqualError(t, err, "object is nil")
}

func TestUpdateObject_WrongTypeUpdateResultInVirtual(t *testing.T) {
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

	virtualUpdate := record.Virtual{
		Union: &record.Virtual_Amend{
			Amend: &record.Amend{},
		},
	}
	virtualUpdateBuf, err := virtualUpdate.Marshal()
	require.NoError(t, err)

	// Incorrect record in virtual.
	virtualResult := record.Virtual{
		Union: &record.Virtual_Genesis{
			Genesis: &record.Genesis{
				Hash: []byte{1, 2, 3, 4, 5},
			},
		},
	}
	virtualResultBuf, err := virtualResult.Marshal()
	require.NoError(t, err)

	update := payload.Update{
		Record: virtualUpdateBuf,
		Result: virtualResultBuf,
	}
	updateBuf, err := update.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: updateBuf,
	}

	handler := handle.NewUpdateObject(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestUpdateObject_FlowWithPassedFlag(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := metaUpdateMsg(t)

	t.Run("checkjet procedure returns unknown err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return nil
			case *proc.FetchJet:
				return errors.New("something strange from checkjet")
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewUpdateObject(proc.NewDependenciesMock(), msg, false)
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
			case *proc.FetchJet:
				return proc.ErrNotExecutor
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewUpdateObject(proc.NewDependenciesMock(), msg, false)
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
			case *proc.FetchJet:
				return proc.ErrNotExecutor
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewUpdateObject(proc.NewDependenciesMock(), msg, true)
		err := handler.Present(ctx, f)
		require.Error(t, err)
		assert.Equal(t, proc.ErrNotExecutor, err)
	})
}

func TestUpdateObject_ErrorFromWaitHot(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := metaUpdateMsg(t)

	t.Run("waithot procedure returns err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return nil
			case *proc.FetchJet:
				return nil
			case *proc.WaitHot:
				return errors.New("error from waithot")
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewUpdateObject(proc.NewDependenciesMock(), msg, false)
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
			case *proc.FetchJet:
				return nil
			case *proc.WaitHot:
				return nil
			case *proc.EnsureIndex:
				return nil
			case *proc.SetResult:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewUpdateObject(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func TestUpdateObject_ErrorFromEnsureIndex(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := metaUpdateMsg(t)

	t.Run("ensureindex procedure returns err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return nil
			case *proc.FetchJet:
				return nil
			case *proc.WaitHot:
				return nil
			case *proc.EnsureIndex:
				return errors.New("error from ensureindex")

			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewUpdateObject(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "error from ensureindex")
	})

	t.Run("ensureindex procedure returns nil err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return nil
			case *proc.FetchJet:
				return nil
			case *proc.WaitHot:
				return nil
			case *proc.EnsureIndex:
				return nil
			case *proc.SetResult:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewUpdateObject(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func TestUpdateObject_ErrorFromUpdateObject(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := metaUpdateMsg(t)

	t.Run("SetResult procedure returns err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return nil
			case *proc.FetchJet:
				return nil
			case *proc.WaitHot:
				return nil
			case *proc.EnsureIndex:
				return nil
			case *proc.SetResult:
				return errors.New("error from SetResult")
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewUpdateObject(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "error from SetResult")
	})

	t.Run("SetResult procedure returns nil err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return nil
			case *proc.FetchJet:
				return nil
			case *proc.WaitHot:
				return nil
			case *proc.EnsureIndex:
				return nil
			case *proc.SetResult:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewUpdateObject(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func metaUpdateMsg(t *testing.T) payload.Meta {
	virtualUpdate := record.Virtual{
		Union: &record.Virtual_Amend{
			Amend: &record.Amend{},
		},
	}
	virtualUpdateBuf, err := virtualUpdate.Marshal()
	require.NoError(t, err)

	// Update.Result object is ok.
	virtualResult := record.Virtual{
		Union: &record.Virtual_Result{
			Result: &record.Result{
				Object: gen.ID(),
			},
		},
	}
	virtualResultBuf, err := virtualResult.Marshal()
	require.NoError(t, err)

	update := payload.Update{
		Record: virtualUpdateBuf,
		Result: virtualResultBuf,
	}
	updateBuf, err := update.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: updateBuf,
	}
	return msg
}
