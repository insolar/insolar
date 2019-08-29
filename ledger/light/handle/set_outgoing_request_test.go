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

func TestSetOutgoingRequest_NilMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   nil,
	}

	handler := handle.NewSetOutgoingRequest(nil, msg, false)

	err := handler.Present(ctx, flow.NewFlowMock(t))
	require.Error(t, err)
}

func TestSetOutgoingRequest_BadMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	msg := payload.Meta{
		Payload: []byte{1, 2, 3, 4, 5},
	}

	handler := handle.NewSetOutgoingRequest(nil, msg, false)

	err := handler.Present(ctx, flow.NewFlowMock(t))
	require.Error(t, err)
}

func TestSetOutgoingRequest_IncorrectTypeMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	f := flow.NewFlowMock(t)

	// Incorrect type.
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

	handler := handle.NewSetOutgoingRequest(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestSetOutgoingRequest_IncorrectRecordInVirtual(t *testing.T) {
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

	request := payload.SetOutgoingRequest{
		Polymorph: uint32(payload.TypeSetOutgoingRequest),
		Request:   virtual,
	}
	requestBuf, err := request.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   requestBuf,
	}

	handler := handle.NewSetOutgoingRequest(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestSetOutgoingRequest_EmptyRequestAffinityRef(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	f := flow.NewFlowMock(t)

	// AffinityRef object (request.Caller) is empty.
	virtual := record.Virtual{
		Union: &record.Virtual_OutgoingRequest{
			OutgoingRequest: &record.OutgoingRequest{
				Caller: insolar.Reference{},
			},
		},
	}

	request := payload.SetOutgoingRequest{
		Polymorph: uint32(payload.TypeSetOutgoingRequest),
		Request:   virtual,
	}
	requestBuf, err := request.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   requestBuf,
	}

	handler := handle.NewSetOutgoingRequest(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	assert.EqualError(t, err, "object is nil")
}

func TestSetOutgoingRequest_ErrorCalculateID(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := metaOutgoingRequestMsg(t)

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

		handler := handle.NewSetOutgoingRequest(proc.NewDependenciesMock(), msg, false)
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
			case *proc.WaitHot:
				return nil
			case *proc.EnsureIndex:
				return nil
			case *proc.SetRequest:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewSetOutgoingRequest(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func TestSetOutgoingRequest_FlowWithPassedFlag(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := metaOutgoingRequestMsg(t)

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

		handler := handle.NewSetOutgoingRequest(proc.NewDependenciesMock(), msg, false)
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

		handler := handle.NewSetOutgoingRequest(proc.NewDependenciesMock(), msg, false)
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

		handler := handle.NewSetOutgoingRequest(proc.NewDependenciesMock(), msg, true)
		err := handler.Present(ctx, f)
		require.Error(t, err)
		assert.Equal(t, proc.ErrNotExecutor, err)
	})
}

func TestSetOutgoingRequest_ErrorFromWaitHot(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := metaOutgoingRequestMsg(t)

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

		handler := handle.NewSetOutgoingRequest(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "error from waithot")
	})

	// Happy path, everything is fine.
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
			case *proc.SetRequest:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewSetOutgoingRequest(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func TestSetOutgoingRequest_ErrorFromEnsureIndex(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := metaOutgoingRequestMsg(t)

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
				return errors.New("can't get index: error from ensureindex")
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewSetOutgoingRequest(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "can't get index: error from ensureindex")
	})

	// Happy path, everything is fine.
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
			case *proc.SetRequest:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewSetOutgoingRequest(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func TestSetOutgoingRequest_ErrorFromSetRequest(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	msg := metaOutgoingRequestMsg(t)

	t.Run("setrequest procedure returns err", func(t *testing.T) {
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
			case *proc.SetRequest:
				return errors.New("error from setrequest")
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewSetOutgoingRequest(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		assert.EqualError(t, err, "error from setrequest")
	})

	// Happy path, everything is fine.
	t.Run("setrequest procedure returns nil err", func(t *testing.T) {
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
			case *proc.SetRequest:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		handler := handle.NewSetOutgoingRequest(proc.NewDependenciesMock(), msg, false)
		err := handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func metaOutgoingRequestMsg(t *testing.T) payload.Meta {
	ref := gen.Reference()

	virtual := record.Virtual{
		Union: &record.Virtual_OutgoingRequest{
			OutgoingRequest: &record.OutgoingRequest{
				Caller: ref,
			},
		},
	}

	request := payload.SetOutgoingRequest{
		Polymorph: uint32(payload.TypeSetOutgoingRequest),
		Request:   virtual,
	}
	requestBuf, err := request.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   requestBuf,
	}

	return msg
}
