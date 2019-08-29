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

func TestActivateObject_NilMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   nil,
	}

	handler := handle.NewActivateObject(nil, msg, false)

	err := handler.Present(ctx, flow.NewFlowMock(t))
	require.Error(t, err)
}

func TestActivateObject_BadMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   []byte{1, 2, 3, 4, 5},
	}

	handler := handle.NewActivateObject(nil, msg, false)

	err := handler.Present(ctx, flow.NewFlowMock(t))
	require.Error(t, err)
}

func TestActivateObject_IncorrectTypeMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	f := flow.NewFlowMock(t)

	// Incorrect type (SetIncomingRequest instead of Activate).
	result := payload.SetIncomingRequest{
		Polymorph: uint32(payload.TypeSetIncomingRequest),
		Request:   record.Virtual{},
	}
	buf, err := result.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   buf,
	}

	handler := handle.NewActivateObject(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestActivateObject_BadWrappedVirtualRecord(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	f := flow.NewFlowMock(t)

	activate := payload.Activate{
		Polymorph: uint32(payload.TypeActivate),
		// Just a byte slice, not a correct virtual record.
		Record: []byte{1, 2, 3, 4, 5},
	}
	buf, err := activate.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
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

	// Incorrect record type in virtual.
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
		Polymorph: uint32(payload.TypeActivate),
		Record:    virtualBuf,
	}
	activateBuf, err := activate.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   activateBuf,
	}

	handler := handle.NewActivateObject(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestActivateObject_EmptyActivateRequestField(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	f := flow.NewFlowMock(t)

	// Activate.Request object is nil.
	virtual := record.Virtual{
		Union: &record.Virtual_Activate{
			Activate: &record.Activate{
				Request: *insolar.NewEmptyReference(),
			},
		},
	}
	virtualBuf, err := virtual.Marshal()
	require.NoError(t, err)

	activate := payload.Activate{
		Polymorph: uint32(payload.TypeActivate),
		Record:    virtualBuf,
	}
	activateBuf, err := activate.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   activateBuf,
	}

	handler := handle.NewActivateObject(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	assert.EqualError(t, err, "request is empty")
}

func TestActivateObject_IncorrectActivateResultPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	f := flow.NewFlowMock(t)

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
		Polymorph: uint32(payload.TypeActivate),
		Record:    virtualActivateBuf,
		// Just a byte slice, not a correct virtual record.
		Result: []byte{1, 2, 3, 4, 5},
	}
	activateBuf, err := activate.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   activateBuf,
	}

	handler := handle.NewActivateObject(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestActivateObject_WrongTypeActivateResultInVirtual(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	f := flow.NewFlowMock(t)

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
		Polymorph: uint32(payload.TypeActivate),
		Record:    virtualActivateBuf,
		// Incorrect value.
		Result: virtualResultBuf,
	}
	activateBuf, err := activate.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   activateBuf,
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
			case *proc.FetchJet:
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
			case *proc.FetchJet:
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
			case *proc.FetchJet:
				return proc.ErrNotExecutor
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
			case *proc.FetchJet:
				return nil
			case *proc.WaitHot:
				return errors.New("error from waithot")
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewActivateObject(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "error from waithot")
	})

	// Happy path, everything is fine.
	t.Run("waithot procedure returns nil err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.FetchJet:
				return nil
			case *proc.WaitHot:
				return nil
			case *proc.SetResult:
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

	t.Run("SetResult procedure returns err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.FetchJet:
				return nil
			case *proc.WaitHot:
				return nil
			case *proc.SetResult:
				return errors.New("error from SetResult")
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewActivateObject(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "error from SetResult")
	})

	// Happy path, everything is fine.
	t.Run("SetResult procedure returns nil err", func(t *testing.T) {
		t.Parallel()
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.FetchJet:
				return nil
			case *proc.WaitHot:
				return nil
			case *proc.SetResult:
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
		Polymorph: uint32(payload.TypeActivate),
		Record:    virtualActivateBuf,
		Result:    virtualResultBuf,
	}
	activateBuf, err := activate.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   activateBuf,
	}
	return msg
}
