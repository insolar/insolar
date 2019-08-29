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
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/handle"
	"github.com/insolar/insolar/ledger/light/proc"
)

func TestGetPendings_NilMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	meta := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   nil,
	}

	handler := handle.NewGetPendings(nil, meta, false)

	err := handler.Present(ctx, flow.NewFlowMock(t))
	require.Error(t, err)
}

func TestGetPendings_BadMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	meta := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   []byte{1, 2, 3, 4, 5},
	}

	handler := handle.NewGetPendings(nil, meta, false)

	err := handler.Present(ctx, flow.NewFlowMock(t))
	require.Error(t, err)
}

func TestGetPendings_IncorrectTypeMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	f := flow.NewFlowMock(t)

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
		f := flow.NewFlowMock(t)
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
		f := flow.NewFlowMock(t)
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
		f := flow.NewFlowMock(t)
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

		handler := handle.NewGetPendings(proc.NewDependenciesMock(), msg, false)
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
		f := flow.NewFlowMock(t)
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
		f := flow.NewFlowMock(t)
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
