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
	"errors"
	"testing"

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

func TestSetRequest_NilMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   nil,
	}

	handler := handle.NewSetIncomingRequest(nil, msg, false)

	err := handler.Present(ctx, flow.NewFlowMock(t))
	require.Error(t, err)
}

func TestSetRequest_BadMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	msg := payload.Meta{
		Payload: []byte{1, 2, 3, 4, 5},
	}

	handler := handle.NewSetIncomingRequest(nil, msg, false)

	err := handler.Present(ctx, flow.NewFlowMock(t))
	require.Error(t, err)
}

func TestSetRequest_IncorrectTypeMsgPayload(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	f := flow.NewFlowMock(t)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		// Incorrect type (SetOutgoingRequest instead of SetIncomingRequest).
		Payload: payload.MustMarshal(&payload.SetOutgoingRequest{
			Polymorph: uint32(payload.TypeSetOutgoingRequest),
			Request:   record.Virtual{},
		}),
	}

	handler := handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, false)

	err := handler.Present(ctx, f)
	require.Error(t, err)
}

func TestSetRequest_IncorrectRecordInVirtual(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	f := flow.NewFlowMock(t)

	f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
		switch p.(type) {
		case *proc.CalculateID:
			return nil
		default:
			panic("unknown procedure")
		}
	})

	// Incorrect record in virtual (Genesis instead of IncomingRequest).
	virtual := record.Virtual{
		Union: &record.Virtual_Genesis{
			Genesis: &record.Genesis{
				Hash: []byte{1, 2, 3, 4, 5},
			},
		},
	}

	request := payload.SetIncomingRequest{
		Polymorph: uint32(payload.TypeSetIncomingRequest),
		Request:   virtual,
	}
	requestBuf, err := request.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   requestBuf,
	}

	handler := handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	require.Error(t, err)
}

func TestSetRequest_EmptyRequestObject(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	f := flow.NewFlowMock(t)

	// IncomingRequest object is nil.
	virtual := record.Virtual{
		Union: &record.Virtual_IncomingRequest{
			IncomingRequest: &record.IncomingRequest{
				Object: nil,
			},
		},
	}

	request := payload.SetIncomingRequest{
		Polymorph: uint32(payload.TypeSetIncomingRequest),
		Request:   virtual,
	}
	requestBuf, err := request.Marshal()
	require.NoError(t, err)

	msg := payload.Meta{
		Polymorph: uint32(payload.TypeMeta),
		Payload:   requestBuf,
	}

	handler := handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, false)

	err = handler.Present(ctx, f)
	assert.EqualError(t, err, "object is nil")
}

func TestSetIncomingRequest_ErrorCalculateID(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	var (
		msg     payload.Meta
		err     error
		handler *handle.SetIncomingRequest
	)

	t.Run("calculateID procedure returns err", func(t *testing.T) {
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return errors.New("something strange from calculateID")
			default:
				panic("unknown procedure")
			}
		})

		// Creation incoming request.
		msg = metaRequestMsg(t, true)
		handler = handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, false)
		err = handler.Present(ctx, f)
		assert.EqualError(t, err, "something strange from calculateID")

		// Non-creation incoming request.
		msg = metaRequestMsg(t, false)
		handler = handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, false)
		err = handler.Present(ctx, f)
		assert.EqualError(t, err, "something strange from calculateID")
	})

	// Happy path, everything is fine.
	t.Run("calculateID procedure returns nil err", func(t *testing.T) {
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

		// Creation incoming request.
		msg = metaRequestMsg(t, true)
		handler = handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, false)
		err = handler.Present(ctx, f)
		require.NoError(t, err)

		// Non-creation incoming request.
		msg = metaRequestMsg(t, false)
		handler = handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, false)
		err = handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func TestSetRequest_FlowWithPassedFlag(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	var (
		msg     payload.Meta
		err     error
		handler *handle.SetIncomingRequest
	)

	t.Run("FetchJet procedure returns unknown err", func(t *testing.T) {
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

		// Creation incoming request.
		msg = metaRequestMsg(t, true)
		handler = handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, false)
		err = handler.Present(ctx, f)
		assert.EqualError(t, err, "something strange from FetchJet")

		// Non-creation incoming request.
		msg = metaRequestMsg(t, false)
		handler = handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, false)
		err = handler.Present(ctx, f)
		assert.EqualError(t, err, "something strange from FetchJet")
	})

	t.Run("passed flag is false and FetchJet returns ErrNotExecutor", func(t *testing.T) {
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

		// Creation incoming request.
		msg = metaRequestMsg(t, true)
		handler = handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, false)
		err = handler.Present(ctx, f)
		require.NoError(t, err)

		// Non-creation incoming request.
		msg = metaRequestMsg(t, false)
		handler = handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, false)
		err = handler.Present(ctx, f)
		require.NoError(t, err)
	})

	t.Run("passed flag is true and FetchJet returns ErrNotExecutor", func(t *testing.T) {
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return nil
			case *proc.FetchJet:
				return proc.ErrNotExecutor
			case *proc.EnsureIndex:
				return nil
			default:
				panic("unknown procedure")
			}
		})

		// Creation incoming request.
		msg = metaRequestMsg(t, true)
		handler = handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, true)
		err = handler.Present(ctx, f)
		require.Error(t, err)
		assert.Equal(t, proc.ErrNotExecutor, err)

		// Non-creation incoming request.
		msg = metaRequestMsg(t, false)
		handler = handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, true)
		err = handler.Present(ctx, f)
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

	var (
		msg     payload.Meta
		err     error
		handler *handle.SetIncomingRequest
	)

	t.Run("WaitHot procedure returns err", func(t *testing.T) {
		f := flow.NewFlowMock(t)
		f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
			switch p.(type) {
			case *proc.CalculateID:
				return nil
			case *proc.FetchJet:
				return nil
			case *proc.WaitHot:
				return errors.New("error from WaitHot")
			default:
				panic("unknown procedure")
			}
		})

		// Creation incoming request.
		msg = metaRequestMsg(t, true)
		handler = handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, false)
		err = handler.Present(ctx, f)
		assert.EqualError(t, err, "error from WaitHot")

		// Non-creation incoming request.
		msg = metaRequestMsg(t, false)
		handler = handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, false)
		err = handler.Present(ctx, f)
		assert.EqualError(t, err, "error from WaitHot")
	})

	// Happy path, everything is fine.
	t.Run("WaitHot procedure returns nil err", func(t *testing.T) {
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

		// Creation incoming request.
		msg = metaRequestMsg(t, true)
		handler = handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, false)
		err = handler.Present(ctx, f)
		require.NoError(t, err)

		// Non-creation incoming request.
		msg = metaRequestMsg(t, false)
		handler = handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, false)
		err = handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func TestSetRequest_ErrorFromEnsureIndex(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	var (
		msg     payload.Meta
		err     error
		handler *handle.SetIncomingRequest
	)

	// EnsureIndex procedure is called only for non-creation IncomingRequest.

	t.Run("EnsureIndex procedure returns err", func(t *testing.T) {
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
				return errors.New("can't get index: error from EnsureIndex")
			default:
				panic("unknown procedure")
			}
		})

		// Non-creation incoming request.
		msg = metaRequestMsg(t, false)
		handler = handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, false)
		err = handler.Present(ctx, f)
		assert.EqualError(t, err, "can't get index: error from EnsureIndex")
	})

	// Happy path, everything is fine.
	t.Run("EnsureIndex procedure returns nil err", func(t *testing.T) {
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
		// Non-creation incoming request.
		msg = metaRequestMsg(t, false)
		handler = handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, false)
		err = handler.Present(ctx, f)
		require.NoError(t, err)
	})
}

func TestSetRequest_ErrorFromSetRequest(t *testing.T) {
	t.Parallel()
	ctx := flow.TestContextWithPulse(
		inslogger.TestContext(t),
		insolar.GenesisPulse.PulseNumber+10,
	)

	var (
		msg     payload.Meta
		err     error
		handler *handle.SetIncomingRequest
	)

	t.Run("SetRequest procedure returns err", func(t *testing.T) {

		t.Run("creation request", func(t *testing.T) {
			f := flow.NewFlowMock(t)
			f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
				switch p.(type) {
				case *proc.CalculateID:
					return nil
				case *proc.FetchJet:
					return nil
				case *proc.WaitHot:
					return nil
				case *proc.SetRequest:
					return errors.New("error from SetRequest")
				default:
					panic("unknown procedure")
				}
			})

			// Creation incoming request.
			msg = metaRequestMsg(t, true)
			handler = handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, false)
			err = handler.Present(ctx, f)
			assert.EqualError(t, err, "error from SetRequest")
		})

		t.Run("non-creation request", func(t *testing.T) {
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
					return errors.New("error from SetRequest")
				default:
					panic("unknown procedure")
				}
			})

			// Creation incoming request.
			msg = metaRequestMsg(t, false)
			handler = handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, false)
			err = handler.Present(ctx, f)
			assert.EqualError(t, err, "error from SetRequest")
		})
	})

	// Happy path, everything is fine.
	t.Run("SetRequest procedure returns nil err", func(t *testing.T) {

		t.Run("creation request", func(t *testing.T) {
			f := flow.NewFlowMock(t)
			f.ProcedureMock.Set(func(ctx context.Context, p flow.Procedure, passed bool) (r error) {
				switch p.(type) {
				case *proc.CalculateID:
					return nil
				case *proc.FetchJet:
					return nil
				case *proc.WaitHot:
					return nil
				case *proc.SetRequest:
					return nil
				default:
					panic("unknown procedure")
				}
			})

			// Creation incoming request.
			msg = metaRequestMsg(t, true)
			handler = handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, false)
			err = handler.Present(ctx, f)
			require.NoError(t, err)
		})

		t.Run("non-creation request", func(t *testing.T) {
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

			// Non-creation incoming request.
			msg = metaRequestMsg(t, false)
			handler = handle.NewSetIncomingRequest(proc.NewDependenciesMock(), msg, false)
			err = handler.Present(ctx, f)
			require.NoError(t, err)
		})
	})
}

func metaRequestMsg(t *testing.T, isCreation bool) payload.Meta {
	ref := gen.Reference()

	var callType record.CallType
	if isCreation {
		callType = record.CTSaveAsChild
	} else {
		callType = record.CTMethod
	}

	virtual := record.Virtual{
		Union: &record.Virtual_IncomingRequest{
			IncomingRequest: &record.IncomingRequest{
				Object:   &ref,
				CallType: callType,
			},
		},
	}

	request := payload.SetIncomingRequest{
		Polymorph: uint32(payload.TypeSetIncomingRequest),
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
