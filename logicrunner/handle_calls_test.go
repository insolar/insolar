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

package logicrunner

import (
	"context"
	"testing"
	"time"

	wmMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/writecontroller"
)

func TestHandleCall_Present(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)

		objRef := gen.Reference()

		fm := flow.NewFlowMock(mc)
		fm.ProcedureMock.Set(func(ctx context.Context, proc flow.Procedure, cancelable bool) (err error) {
			switch p := proc.(type) {
			case *CheckOurRole:
				return nil
			case *RegisterIncomingRequest:
				p.result <- &payload.RequestInfo{RequestID: gen.ID(), ObjectID: *objRef.GetLocal()}
				return nil
			case *AddFreshRequest:
				return nil
			default:
				t.Fatalf("Unknown procedure: %T", proc)
			}
			return nil
		})

		handler := HandleCall{
			dep: &Dependencies{
				Publisher: nil,
				StateStorage: NewStateStorageMock(mc).
					UpsertExecutionStateMock.Expect(objRef).Return(nil),
				ResultsMatcher:  nil,
				ArtifactManager: artifacts.NewClientMock(mc),
				Sender:          nil,
				JetStorage:      nil,
				WriteAccessor:   writecontroller.NewAccessorMock(mc).BeginMock.Return(func() {}, nil),
			},
			Message: payload.Meta{},
		}

		msg := payload.CallMethod{
			Request: &record.IncomingRequest{
				CallType: record.CTMethod,
				Object:   &objRef,
			},
		}

		reply, err := handler.handleActual(ctx, msg, fm)
		assert.NotNil(t, reply)
		assert.NoError(t, err)

		mc.Wait(time.Minute)
		mc.Finish()
	})

	t.Run("write accessor failed to fetch lock", func(t *testing.T) {
		t.Parallel()

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)

		objRef := gen.Reference()

		fm := flow.NewFlowMock(mc)
		fm.ProcedureMock.Set(func(ctx context.Context, proc flow.Procedure, cancelable bool) (err error) {
			switch p := proc.(type) {
			case *CheckOurRole:
				return nil
			case *RegisterIncomingRequest:
				p.result <- &payload.RequestInfo{RequestID: gen.ID(), ObjectID: *objRef.GetLocal()}
				return nil
			case *AddFreshRequest:
				return nil
			default:
				t.Fatalf("Unknown procedure: %T", proc)
			}
			return nil
		})

		handler := HandleCall{
			dep: &Dependencies{
				Publisher:      nil,
				StateStorage:   NewStateStorageMock(mc),
				ResultsMatcher: nil,
				Sender: bus.NewSenderMock(mc).SendRoleMock.Set(
					func(ctx context.Context, msg *wmMessage.Message, role insolar.DynamicRole, obj insolar.Reference) (<-chan *wmMessage.Message, func()) {
						payloadType, err := payload.UnmarshalType(msg.Payload)
						require.NoError(t, err, "unmarshalType")
						require.Equal(t, payload.TypeAdditionalCallFromPreviousExecutor, payloadType)
						return nil, func() {}
					}),
				ArtifactManager: artifacts.NewClientMock(mc),
				JetStorage:      nil,
				WriteAccessor: writecontroller.NewAccessorMock(mc).
					BeginMock.Return(func() {}, writecontroller.ErrWriteClosed),
				PulseAccessor: pulse.NewAccessorMock(mc).
					LatestMock.Return(insolar.Pulse{PulseNumber: 100}, nil),
			},
			Message: payload.Meta{},
		}

		msg := payload.CallMethod{
			Request: &record.IncomingRequest{
				CallType: record.CTMethod,
				Object:   &objRef,
			},
		}

		reply, err := handler.handleActual(ctx, msg, fm)
		assert.NotNil(t, reply)
		assert.NoError(t, err)

		mc.Wait(time.Minute)
		mc.Finish()
	})

	t.Run("failed to authorize", func(t *testing.T) {
		t.Parallel()

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)

		fm := flow.NewFlowMock(mc)
		fm.ProcedureMock.Set(func(ctx context.Context, proc flow.Procedure, cancelable bool) (err error) {
			switch proc.(type) {
			case *CheckOurRole:
				return ErrCantExecute
			case *RegisterIncomingRequest:
				t.Fatalf("Shouldn't be called: %T", proc)
			case *AddFreshRequest:
				t.Fatalf("Shouldn't be called: %T", proc)
			default:
				t.Fatalf("Unknown procedure: %T", proc)
			}
			return nil
		})

		objRef := gen.Reference()
		handler := HandleCall{
			dep: &Dependencies{
				Publisher:       nil,
				StateStorage:    NewStateStorageMock(mc),
				ResultsMatcher:  nil,
				ArtifactManager: artifacts.NewClientMock(mc),
				Sender:          nil,
				JetStorage:      nil,
				WriteAccessor:   writecontroller.NewAccessorMock(mc),
			},
			Message: payload.Meta{},
		}

		msg := payload.CallMethod{
			Request: &record.IncomingRequest{
				CallType: record.CTMethod,
				Object:   &objRef,
			},
		}

		reply, err := handler.handleActual(ctx, msg, fm)
		assert.Nil(t, reply)
		assert.EqualError(t, err, flow.ErrCancelled.Error())

		mc.Wait(time.Minute)
		mc.Finish()
	})

	t.Run("failed to register incoming request", func(t *testing.T) {
		t.Parallel()

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)

		fm := flow.NewFlowMock(mc)
		fm.ProcedureMock.Set(func(ctx context.Context, proc flow.Procedure, cancelable bool) (err error) {
			switch proc.(type) {
			case *CheckOurRole:
				return nil
			case *RegisterIncomingRequest:
				return flow.ErrCancelled
			case *AddFreshRequest:
				t.Fatalf("Shouldn't be called: %T", proc)
			default:
				t.Fatalf("Unknown procedure: %T", proc)
			}
			return nil
		})

		objRef := gen.Reference()
		handler := HandleCall{
			dep: &Dependencies{
				Publisher:       nil,
				ResultsMatcher:  nil,
				ArtifactManager: artifacts.NewClientMock(mc),
				Sender:          nil,
				JetStorage:      nil,
				WriteAccessor:   writecontroller.NewAccessorMock(mc),
			},
			Message: payload.Meta{},
		}

		msg := payload.CallMethod{
			Request: &record.IncomingRequest{
				CallType: record.CTMethod,
				Object:   &objRef,
			},
		}

		reply, err := handler.handleActual(ctx, msg, fm)
		assert.Nil(t, reply)
		assert.EqualError(t, err, flow.ErrCancelled.Error())

		mc.Wait(time.Minute)
		mc.Finish()
	})

	t.Run("objectRef for CTMethod is empty", func(t *testing.T) {
		t.Parallel()

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)

		fm := flow.NewFlowMock(mc)
		fm.ProcedureMock.Set(func(ctx context.Context, proc flow.Procedure, cancelable bool) (err error) {
			switch proc.(type) {
			case *CheckOurRole:
				return nil
			case *RegisterIncomingRequest:
				return flow.ErrCancelled
			case *AddFreshRequest:
				t.Fatalf("Shouldn't be called: %T", proc)
			default:
				t.Fatalf("Unknown procedure: %T", proc)
			}
			return nil
		})

		handler := HandleCall{
			dep: &Dependencies{
				Publisher:       nil,
				StateStorage:    NewStateStorageMock(mc),
				ResultsMatcher:  nil,
				ArtifactManager: artifacts.NewClientMock(mc),
				Sender:          nil,
				JetStorage:      nil,
				WriteAccessor:   writecontroller.NewAccessorMock(mc),
			},
			Message: payload.Meta{},
		}

		msg := payload.CallMethod{
			Request: &record.IncomingRequest{
				CallType: record.CTMethod,
				Object:   nil,
			},
			PulseNumber: gen.PulseNumber(),
		}

		reply, err := handler.handleActual(ctx, msg, fm)
		assert.Nil(t, reply)
		assert.Error(t, err)

		mc.Wait(time.Minute)
		mc.Finish()
	})

	t.Run("write accessor failed to fetch lock AND registry is empty after on pulse", func(t *testing.T) {
		t.Parallel()

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)

		objRef := gen.Reference()

		fm := flow.NewFlowMock(mc)
		fm.ProcedureMock.Set(func(ctx context.Context, proc flow.Procedure, cancelable bool) (err error) {
			switch p := proc.(type) {
			case *CheckOurRole:
				return nil
			case *RegisterIncomingRequest:
				p.result <- &payload.RequestInfo{RequestID: gen.ID(), ObjectID: *objRef.GetLocal()}
				return nil
			case *AddFreshRequest:
				return nil
			default:
				t.Fatalf("Unknown procedure: %T", proc)
			}
			return nil
		})

		handler := HandleCall{
			dep: &Dependencies{
				Publisher:       nil,
				ResultsMatcher:  nil,
				ArtifactManager: artifacts.NewClientMock(mc),
				Sender: bus.NewSenderMock(mc).SendRoleMock.Set(
					func(ctx context.Context, msg *wmMessage.Message, role insolar.DynamicRole, obj insolar.Reference) (<-chan *wmMessage.Message, func()) {
						payloadType, err := payload.UnmarshalType(msg.Payload)
						require.NoError(t, err, "unmarshalType")
						require.Equal(t, payload.TypeAdditionalCallFromPreviousExecutor, payloadType)
						return nil, func() {}
					}),
				JetStorage:    nil,
				WriteAccessor: writecontroller.NewAccessorMock(mc).BeginMock.Return(func() {}, writecontroller.ErrWriteClosed),
				PulseAccessor: pulse.NewAccessorMock(mc).
					LatestMock.Return(insolar.Pulse{PulseNumber: 100}, nil),
			},
			Message: payload.Meta{},
		}

		msg := payload.CallMethod{
			Request: &record.IncomingRequest{
				CallType: record.CTMethod,
				Object:   &objRef,
			},
			PulseNumber: gen.PulseNumber(),
		}

		reply, err := handler.handleActual(ctx, msg, fm)
		assert.NotNil(t, reply)
		assert.NoError(t, err)

		mc.Wait(time.Minute)
		mc.Finish()
	})

	t.Run("already completed request", func(t *testing.T) {
		t.Parallel()

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)

		objRef := gen.Reference()
		reqRef := gen.RecordReference()

		resRecord := &record.Result{Payload: []byte{3, 2, 1}}
		virtResRecord := record.Wrap(resRecord)
		matRecord := record.Material{Virtual: virtResRecord}
		matRecordSerialized, err := matRecord.Marshal()
		require.NoError(t, err)

		fm := flow.NewFlowMock(mc)
		fm.ProcedureMock.Set(func(ctx context.Context, proc flow.Procedure, cancelable bool) (err error) {
			switch p := proc.(type) {
			case *CheckOurRole:
				return nil
			case *RegisterIncomingRequest:
				p.result <- &payload.RequestInfo{
					RequestID: *reqRef.GetLocal(),
					ObjectID:  *objRef.GetLocal(),
					Request:   []byte{1, 2, 3},
					Result:    matRecordSerialized,
				}
				return nil
			case *AddFreshRequest:
				return nil
			default:
				t.Fatalf("Unknown procedure: %T", proc)
			}
			return nil
		})

		handler := HandleCall{
			dep: &Dependencies{
				ArtifactManager:  artifacts.NewClientMock(mc),
				RequestsExecutor: NewRequestsExecutorMock(mc).SendReplyMock.Return(),
			},
			Message: payload.Meta{},
		}

		msg := payload.CallMethod{
			Request: &record.IncomingRequest{
				CallType: record.CTMethod,
				Object:   &objRef,
			},
		}

		gotReply, err := handler.handleActual(ctx, msg, fm)
		require.NoError(t, err)
		require.Equal(t, &reply.RegisterRequest{Request: reqRef}, gotReply)

		mc.Wait(time.Minute)
		mc.Finish()
	})

	t.Run("object not found during request registration", func(t *testing.T) {
		t.Parallel()

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)

		objRef := gen.Reference()

		fm := flow.NewFlowMock(mc)
		fm.ProcedureMock.Set(func(ctx context.Context, proc flow.Procedure, cancelable bool) (err error) {
			switch proc.(type) {
			case *CheckOurRole:
				return nil
			case *RegisterIncomingRequest:
				return errors.Wrap(
					&payload.CodedError{Code: payload.CodeNotFound, Text: "index not found"},
					"RegisterIncomingRequest")
			case *AddFreshRequest:
				return nil
			default:
				t.Fatalf("Unknown procedure: %T", proc)
			}
			return nil
		})

		handler := HandleCall{
			dep: &Dependencies{
				ArtifactManager: artifacts.NewClientMock(mc),
			},
			Message: payload.Meta{},
		}

		msg := payload.CallMethod{
			Request: &record.IncomingRequest{
				CallType: record.CTMethod,
				Object:   &objRef,
			},
		}

		expectedResult, err := foundation.MarshalMethodErrorResult(errors.New("RegisterIncomingRequest: index not found"))
		require.NoError(t, err)

		expectedReply := &reply.CallMethod{Result: expectedResult}
		gotReply, err := handler.handleActual(ctx, msg, fm)
		require.NoError(t, err)
		require.Equal(t, expectedReply, gotReply)

		mc.Wait(time.Minute)
		mc.Finish()
	})

	t.Run("loop detected", func(t *testing.T) {
		t.Parallel()

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)

		objRef := gen.Reference()

		fm := flow.NewFlowMock(mc)
		fm.ProcedureMock.Set(func(ctx context.Context, proc flow.Procedure, cancelable bool) (err error) {
			switch p := proc.(type) {
			case *CheckOurRole:
				return nil
			case *RegisterIncomingRequest:
				return &payload.CodedError{Code: payload.CodeLoopDetected, Text: "loop detected"}
			case *RecordErrorResult:
				p.result = []byte{3, 2, 1}
				return nil
			default:
				t.Fatalf("Unknown procedure: %T", proc)
			}
			return nil
		})

		handler := HandleCall{
			dep: &Dependencies{
				Publisher:       nil,
				StateStorage:    NewStateStorageMock(mc),
				ResultsMatcher:  nil,
				ArtifactManager: artifacts.NewClientMock(mc),
			},
			Message: payload.Meta{},
		}

		msg := payload.CallMethod{
			Request: &record.IncomingRequest{
				CallType: record.CTMethod,
				Object:   &objRef,
			},
		}

		resultWithErr, err := foundation.MarshalMethodErrorResult(&payload.CodedError{Code: payload.CodeLoopDetected, Text: "loop detected"})
		require.NoError(t, err)

		expectedReply := &reply.CallMethod{Result: resultWithErr}
		gotReply, err := handler.handleActual(ctx, msg, fm)
		require.NoError(t, err)
		require.Equal(t, expectedReply, gotReply)

		mc.Wait(time.Minute)
		mc.Finish()
	})
}
