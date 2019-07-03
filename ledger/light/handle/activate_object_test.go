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

func TestActivateObject_BadMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	msg := payload.Meta{
		Payload: []byte{1, 2, 3, 4, 5},
	}

	handler := handle.NewActivateObject(nil, msg, false)

	err := handler.Present(ctx, flow.NewFlowMock(t))
	require.Error(t, err)
}

func TestActivateObject_BadWrappedVirtualRecord(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	f := flow.NewFlowMock(t)

	activate := payload.Activate{
		Record: []byte{1, 2, 3, 4, 5},
	}
	buf, err := activate.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		// This buf is not wrapped as virtual record.
		Payload: buf,
	}

	handler := handle.NewActivateObject(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestActivateObject_IncorrectActivateRecordInVirtual(t *testing.T) {
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

	activate := payload.Activate{
		Record: virtualBuf,
	}
	activateBuf, err := activate.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: activateBuf,
	}

	handler := handle.NewActivateObject(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestActivateObject_EmptyActivateRequestField(t *testing.T) {
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

	// Activate.Request object is nil.
	virtual := record.Virtual{
		Union: &record.Virtual_Activate{
			Activate: &record.Activate{
				Request: insolar.Reference{},
			},
		},
	}
	virtualBuf, err := virtual.Marshal()
	require.NoError(t, err)

	activate := payload.Activate{
		Record: virtualBuf,
	}
	activateBuf, err := activate.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: activateBuf,
	}

	handler := handle.NewActivateObject(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	assert.EqualError(t, err, "request is nil")
}

func TestActivateObject_IncorrectActivateResultPayload(t *testing.T) {
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

	// Activate.Request is ok.
	virtualActivate := record.Virtual{
		Union: &record.Virtual_Activate{
			Activate: &record.Activate{
				Request: gen.Reference(),
			},
		},
	}
	virtualActivateBuf, err := virtualActivate.Marshal()
	require.NoError(t, err)

	activate := payload.Activate{
		Record: virtualActivateBuf,
		Result: []byte{1, 2, 3, 4, 5},
	}
	activateBuf, err := activate.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: activateBuf,
	}

	handler := handle.NewActivateObject(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestActivateObject_WrongTypeActivateResultInVirtual(t *testing.T) {
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

	// Activate.Request is ok.
	virtualActivate := record.Virtual{
		Union: &record.Virtual_Activate{
			Activate: &record.Activate{
				Request: gen.Reference(),
			},
		},
	}
	virtualActivateBuf, err := virtualActivate.Marshal()
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

	activate := payload.Activate{
		Record: virtualActivateBuf,
		Result: virtualResultBuf,
	}
	activateBuf, err := activate.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: activateBuf,
	}

	handler := handle.NewActivateObject(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestActivateObject_FlowWithPassedFlag(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := metaActivateMsg(t)

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

		handler := handle.NewActivateObject(proc.NewDependenciesMock(), msg, false)
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

		handler := handle.NewActivateObject(proc.NewDependenciesMock(), msg, false)
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
			case *proc.EnsureIndexWM:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewActivateObject(proc.NewDependenciesMock(), msg, true)
		err := handler.Present(ctx, f)
		require.Error(t, err)
		assert.Equal(t, proc.ErrNotExecutor, err)
	})
}

func TestActivateObject_ErrorFromWaitHot(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := metaActivateMsg(t)

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

		handler := handle.NewActivateObject(proc.NewDependenciesMock(), msg, false)
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
			case *proc.ActivateObject:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewActivateObject(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func TestActivateObject_ErrorFromActivateObject(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := metaActivateMsg(t)

	t.Run("activateobject procedure returns err", func(t *testing.T) {
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
			case *proc.ActivateObject:
				return errors.New("error from activateobject")
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewActivateObject(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "error from activateobject")
	})

	t.Run("activateobject procedure returns nil err", func(t *testing.T) {
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
			case *proc.ActivateObject:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewActivateObject(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func metaActivateMsg(t *testing.T) payload.Meta {
	// Activate.Request is ok.
	virtualActivate := record.Virtual{
		Union: &record.Virtual_Activate{
			Activate: &record.Activate{
				Request: gen.Reference(),
			},
		},
	}
	virtualActivateBuf, err := virtualActivate.Marshal()
	require.NoError(t, err)

	// Activate.Result is ok.
	virtualResult := record.Virtual{
		Union: &record.Virtual_Result{
			Result: &record.Result{
				Object: gen.ID(),
			},
		},
	}
	virtualResultBuf, err := virtualResult.Marshal()
	require.NoError(t, err)

	activate := payload.Activate{
		Record: virtualActivateBuf,
		Result: virtualResultBuf,
	}
	activateBuf, err := activate.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: activateBuf,
	}
	return msg
}
