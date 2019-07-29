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

func TestDeactivateObject_BadMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	msg := payload.Meta{
		Payload: []byte{1, 2, 3, 4, 5},
	}

	handler := handle.NewDeactivateObject(nil, msg, false)

	err := handler.Present(ctx, flow.NewFlowMock(t))
	require.Error(t, err)
}

func TestDeactivateObject_BadWrappedVirtualRecord(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	f := flow.NewFlowMock(t)

	deactivate := payload.Deactivate{
		Record: []byte{1, 2, 3, 4, 5},
	}
	buf, err := deactivate.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		// This buf is not wrapped as virtual record.
		Payload: buf,
	}

	handler := handle.NewDeactivateObject(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestDeactivateObject_IncorrectDeactivateRecordInVirtual(t *testing.T) {
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

	deactivate := payload.Deactivate{
		Record: virtualBuf,
	}
	deactivateBuf, err := deactivate.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: deactivateBuf,
	}

	handler := handle.NewDeactivateObject(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestDeactivateObject_IncorrectDeactivateResultPayload(t *testing.T) {
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

	virtualDeactivate := record.Virtual{
		Union: &record.Virtual_Deactivate{
			Deactivate: &record.Deactivate{},
		},
	}
	virtualDeactivateBuf, err := virtualDeactivate.Marshal()
	require.NoError(t, err)

	deactivate := payload.Deactivate{
		Record: virtualDeactivateBuf,
		Result: []byte{1, 2, 3, 4, 5},
	}
	deactivateBuf, err := deactivate.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: deactivateBuf,
	}

	handler := handle.NewDeactivateObject(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestDeactivateObject_EmptyDeactivateResultObject(t *testing.T) {
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

	virtualDeactivate := record.Virtual{
		Union: &record.Virtual_Deactivate{
			Deactivate: &record.Deactivate{},
		},
	}
	virtualDeactivateBuf, err := virtualDeactivate.Marshal()
	require.NoError(t, err)

	// Deactivate.Result object is empty
	virtualResult := record.Virtual{
		Union: &record.Virtual_Result{
			Result: &record.Result{
				Object: insolar.ID{},
			},
		},
	}
	virtualResultBuf, err := virtualResult.Marshal()
	require.NoError(t, err)

	deactivate := payload.Deactivate{
		Record: virtualDeactivateBuf,
		Result: virtualResultBuf,
	}
	deactivateBuf, err := deactivate.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: deactivateBuf,
	}

	handler := handle.NewDeactivateObject(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	assert.EqualError(t, err, "object is nil")
}

func TestDeactivateObject_WrongTypeDeactivateResultInVirtual(t *testing.T) {
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

	virtualDeactivate := record.Virtual{
		Union: &record.Virtual_Deactivate{
			Deactivate: &record.Deactivate{},
		},
	}
	virtualDeactivateBuf, err := virtualDeactivate.Marshal()
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

	deactivate := payload.Deactivate{
		Record: virtualDeactivateBuf,
		Result: virtualResultBuf,
	}
	deactivateBuf, err := deactivate.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: deactivateBuf,
	}

	handler := handle.NewDeactivateObject(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestDeactivateObject_FlowWithPassedFlag(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := metaDeactivateMsg(t)

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

		handler := handle.NewDeactivateObject(proc.NewDependenciesMock(), msg, false)
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

		handler := handle.NewDeactivateObject(proc.NewDependenciesMock(), msg, false)
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

		handler := handle.NewDeactivateObject(proc.NewDependenciesMock(), msg, true)
		err := handler.Present(ctx, f)
		require.Error(t, err)
		assert.Equal(t, proc.ErrNotExecutor, err)
	})
}

func TestDeactivateObject_ErrorFromWaitHot(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := metaDeactivateMsg(t)

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

		handler := handle.NewDeactivateObject(proc.NewDependenciesMock(), msg, false)
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

		handler := handle.NewDeactivateObject(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func TestDeactivateObject_ErrorFromEnsureIndex(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := metaDeactivateMsg(t)

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

		handler := handle.NewDeactivateObject(proc.NewDependenciesMock(), msg, false)
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

		handler := handle.NewDeactivateObject(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func TestDeactivateObject_ErrorFromDeactivateObject(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := metaDeactivateMsg(t)

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

		handler := handle.NewDeactivateObject(proc.NewDependenciesMock(), msg, false)
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

		handler := handle.NewDeactivateObject(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func metaDeactivateMsg(t *testing.T) payload.Meta {
	virtualDeactivate := record.Virtual{
		Union: &record.Virtual_Deactivate{
			Deactivate: &record.Deactivate{},
		},
	}
	virtualDeactivateBuf, err := virtualDeactivate.Marshal()
	require.NoError(t, err)

	// Deactivate.Result object is ok.
	virtualResult := record.Virtual{
		Union: &record.Virtual_Result{
			Result: &record.Result{
				Object: gen.ID(),
			},
		},
	}
	virtualResultBuf, err := virtualResult.Marshal()
	require.NoError(t, err)

	deactivate := payload.Deactivate{
		Record: virtualDeactivateBuf,
		Result: virtualResultBuf,
	}
	deactivateBuf, err := deactivate.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Payload: deactivateBuf,
	}
	return msg
}
