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
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/executionregistry"
	"github.com/insolar/insolar/logicrunner/writecontroller"
)

func TestHandleCall_CheckExecutionLoop(t *testing.T) {
	obj := gen.Reference()

	tests := []struct {
		name  string
		mocks func(t minimock.Tester) (*HandleCall, *record.IncomingRequest)
		loop  bool
	}{
		{
			name: "loop detected",
			loop: true,
			mocks: func(t minimock.Tester) (*HandleCall, *record.IncomingRequest) {
				h := &HandleCall{
					dep: &Dependencies{
						StateStorage: NewStateStorageMock(t).
							GetExecutionRegistryMock.Expect(obj).
							Return(
								executionregistry.NewExecutionRegistryMock(t).
									FindRequestLoopMock.Return(true),
							),
					},
				}
				req := &record.IncomingRequest{
					Object: &obj,
				}
				return h, req
			},
		},
		{
			name: "no loop, broker check",
			mocks: func(t minimock.Tester) (*HandleCall, *record.IncomingRequest) {
				h := &HandleCall{
					dep: &Dependencies{
						StateStorage: NewStateStorageMock(t).
							GetExecutionRegistryMock.Expect(obj).
							Return(
								executionregistry.NewExecutionRegistryMock(t).
									FindRequestLoopMock.Return(false),
							),
					},
				}
				req := &record.IncomingRequest{
					Object: &obj,
				}
				return h, req
			},
		},
		{
			name: "no loop, not executing",
			mocks: func(t minimock.Tester) (*HandleCall, *record.IncomingRequest) {
				h := &HandleCall{
					dep: &Dependencies{
						StateStorage: NewStateStorageMock(t).
							GetExecutionRegistryMock.Expect(obj).
							Return(nil),
					},
				}
				req := &record.IncomingRequest{
					Object: &obj,
				}
				return h, req
			},
		},
		{
			name: "no loop, nil object",
			mocks: func(t minimock.Tester) (*HandleCall, *record.IncomingRequest) {
				h := &HandleCall{
					dep: &Dependencies{},
				}
				req := &record.IncomingRequest{}
				return h, req
			},
		},
		{
			name: "no loop, constructor",
			mocks: func(t minimock.Tester) (*HandleCall, *record.IncomingRequest) {
				h := &HandleCall{
					dep: &Dependencies{},
				}
				req := &record.IncomingRequest{
					CallType: record.CTSaveAsChild,
				}
				return h, req
			},
		},
		{
			name: "no loop, no wait call",
			mocks: func(t minimock.Tester) (*HandleCall, *record.IncomingRequest) {
				h := &HandleCall{
					dep: &Dependencies{},
				}
				req := &record.IncomingRequest{
					ReturnMode: record.ReturnNoWait,
				}
				return h, req
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := inslogger.TestContext(t)
			mc := minimock.NewController(t)

			h, req := test.mocks(mc)
			loop := h.checkExecutionLoop(ctx, gen.Reference(), *req)
			require.Equal(t, test.loop, loop)

			mc.Wait(1 * time.Minute)
			mc.Finish()
		})
	}
}

func TestHandleCall_Present(t *testing.T) {

	t.Run("happy path", func(t *testing.T) {
		t.Parallel()

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)
		defer mc.Wait(time.Minute)

		objRef := gen.Reference()

		fm := flow.NewFlowMock(mc)
		fm.ProcedureMock.Set(func(ctx context.Context, proc flow.Procedure, cancelable bool) (err error) {
			switch p := proc.(type) {
			case *CheckOurRole:
				return nil
			case *RegisterIncomingRequest:
				p.result <- &payload.RequestInfo{RequestID: gen.ID(), ObjectID: *objRef.Record()}
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
					GetExecutionRegistryMock.Expect(objRef).Return(
					executionregistry.NewExecutionRegistryMock(mc).FindRequestLoopMock.Return(false),
				).
					UpsertExecutionStateMock.Expect(objRef).Return(nil),
				ResultsMatcher: nil,
				lr: &LogicRunner{
					ArtifactManager: artifacts.NewClientMock(mc),
				},
				Sender:        nil,
				JetStorage:    nil,
				WriteAccessor: writecontroller.NewAccessorMock(mc).BeginMock.Return(func() {}, nil),
			},
			Message: payload.Meta{},
			Parcel:  nil,
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
	})

	t.Run("write accessor failed to fetch lock", func(t *testing.T) {
		t.Parallel()

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)
		defer mc.Wait(time.Second)

		objRef := gen.Reference()

		fm := flow.NewFlowMock(mc)
		fm.ProcedureMock.Set(func(ctx context.Context, proc flow.Procedure, cancelable bool) (err error) {
			switch p := proc.(type) {
			case *CheckOurRole:
				return nil
			case *RegisterIncomingRequest:
				p.result <- &payload.RequestInfo{RequestID: gen.ID(), ObjectID: *objRef.Record()}
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
					GetExecutionRegistryMock.Expect(objRef).Return(
					executionregistry.NewExecutionRegistryMock(mc).FindRequestLoopMock.Return(false),
				),
				ResultsMatcher: nil,
				Sender: bus.NewSenderMock(mc).SendRoleMock.Set(
					func(ctx context.Context, msg *wmMessage.Message, role insolar.DynamicRole, obj insolar.Reference) (<-chan *wmMessage.Message, func()) {
						payloadType, err := payload.UnmarshalType(msg.Payload)
						require.NoError(t, err, "unmarshalType")
						require.Equal(t, payload.TypeAdditionalCallFromPreviousExecutor, payloadType)
						return nil, func() {}
					}),
				lr: &LogicRunner{
					ArtifactManager: artifacts.NewClientMock(mc),
				},
				JetStorage: nil,
				WriteAccessor: writecontroller.NewAccessorMock(mc).
					BeginMock.Return(func() {}, writecontroller.ErrWriteClosed),
			},
			Message: payload.Meta{},
			Parcel:  nil,
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
	})

	t.Run("failed to authorize", func(t *testing.T) {
		t.Parallel()

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)
		defer mc.Wait(time.Minute)

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
				Publisher:      nil,
				StateStorage:   NewStateStorageMock(mc),
				ResultsMatcher: nil,
				lr: &LogicRunner{
					ArtifactManager: artifacts.NewClientMock(mc),
				},
				Sender:        nil,
				JetStorage:    nil,
				WriteAccessor: writecontroller.NewAccessorMock(mc),
			},
			Message: payload.Meta{},
			Parcel:  nil,
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
	})

	t.Run("failed to register incoming request", func(t *testing.T) {
		t.Parallel()

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)
		defer mc.Wait(time.Minute)

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
				Publisher:      nil,
				ResultsMatcher: nil,
				lr: &LogicRunner{
					ArtifactManager: artifacts.NewClientMock(mc),
				},
				Sender:        nil,
				JetStorage:    nil,
				WriteAccessor: writecontroller.NewAccessorMock(mc),
			},
			Message: payload.Meta{},
			Parcel:  nil,
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
	})

	t.Run("objectRef for CTMethod is empty", func(t *testing.T) {
		t.Parallel()

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)
		defer mc.Wait(time.Second)

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
				Publisher:      nil,
				StateStorage:   NewStateStorageMock(mc),
				ResultsMatcher: nil,
				lr: &LogicRunner{
					ArtifactManager: artifacts.NewClientMock(mc),
				},
				Sender:        nil,
				JetStorage:    nil,
				WriteAccessor: writecontroller.NewAccessorMock(mc),
			},
			Message: payload.Meta{},
			Parcel:  nil,
		}

		msg := payload.CallMethod{
			Request: &record.IncomingRequest{
				CallType: record.CTMethod,
				Object:   nil,
			},
		}

		reply, err := handler.handleActual(ctx, msg, fm)
		assert.Nil(t, reply)
		assert.Error(t, err)
	})

	t.Run("write accessor failed to fetch lock AND registry is empty after on pulse", func(t *testing.T) {
		t.Parallel()

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)
		defer mc.Wait(time.Minute)

		objRef := gen.Reference()

		fm := flow.NewFlowMock(mc)
		fm.ProcedureMock.Set(func(ctx context.Context, proc flow.Procedure, cancelable bool) (err error) {
			switch p := proc.(type) {
			case *CheckOurRole:
				return nil
			case *RegisterIncomingRequest:
				p.result <- &payload.RequestInfo{RequestID: gen.ID(), ObjectID: *objRef.Record()}
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
					GetExecutionRegistryMock.Expect(objRef).Return(nil),
				ResultsMatcher: nil,
				lr: &LogicRunner{
					ArtifactManager: artifacts.NewClientMock(mc),
				},
				Sender: bus.NewSenderMock(mc).SendRoleMock.Set(
					func(ctx context.Context, msg *wmMessage.Message, role insolar.DynamicRole, obj insolar.Reference) (<-chan *wmMessage.Message, func()) {
						payloadType, err := payload.UnmarshalType(msg.Payload)
						require.NoError(t, err, "unmarshalType")
						require.Equal(t, payload.TypeAdditionalCallFromPreviousExecutor, payloadType)
						return nil, func() {}
					}),
				JetStorage:    nil,
				WriteAccessor: writecontroller.NewAccessorMock(mc).BeginMock.Return(func() {}, writecontroller.ErrWriteClosed),
			},
			Message: payload.Meta{},
			Parcel:  nil,
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
	})

	t.Run("already completed request", func(t *testing.T) {
		t.Parallel()

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)
		defer mc.Wait(time.Minute)

		objRef := gen.Reference()
		reqRef := gen.Reference()

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
					RequestID: *reqRef.Record(),
					ObjectID:  *objRef.Record(),
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
				lr: &LogicRunner{
					ArtifactManager: artifacts.NewClientMock(mc),
				},
				RequestsExecutor: NewRequestsExecutorMock(mc).SendReplyMock.Return(),
			},
			Message: payload.Meta{},
			Parcel:  nil,
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
	})

	t.Run("object not found during request registration", func(t *testing.T) {
		t.Parallel()

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)
		defer mc.Wait(time.Minute)

		objRef := gen.Reference()

		fm := flow.NewFlowMock(mc)
		fm.ProcedureMock.Set(func(ctx context.Context, proc flow.Procedure, cancelable bool) (err error) {
			switch proc.(type) {
			case *CheckOurRole:
				return nil
			case *RegisterIncomingRequest:
				return &payload.CodedError{Code: payload.CodeNotFound, Text: "index not found"}
			case *AddFreshRequest:
				return nil
			default:
				t.Fatalf("Unknown procedure: %T", proc)
			}
			return nil
		})

		handler := HandleCall{
			dep: &Dependencies{
				lr: &LogicRunner{
					ArtifactManager: artifacts.NewClientMock(mc),
				},
			},
			Message: payload.Meta{},
			Parcel:  nil,
		}

		msg := payload.CallMethod{
			Request: &record.IncomingRequest{
				CallType: record.CTMethod,
				Object:   &objRef,
			},
		}

		expectedResult, err := foundation.MarshalMethodErrorResult(errors.New("index not found"))
		require.NoError(t, err)

		expectedReply := &reply.CallMethod{Result: expectedResult}
		gotReply, err := handler.handleActual(ctx, msg, fm)
		require.NoError(t, err)
		require.Equal(t, expectedReply, gotReply)

	})

	t.Run("loop detected", func(t *testing.T) {
		t.Parallel()

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)
		defer mc.Wait(time.Minute)

		objRef := gen.Reference()

		fm := flow.NewFlowMock(mc)
		fm.ProcedureMock.Set(func(ctx context.Context, proc flow.Procedure, cancelable bool) (err error) {
			switch p := proc.(type) {
			case *CheckOurRole:
				return nil
			case *RegisterIncomingRequest:
				p.result <- &payload.RequestInfo{RequestID: gen.ID(), ObjectID: *objRef.Record()}
				return nil
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
				Publisher: nil,
				StateStorage: NewStateStorageMock(mc).
					GetExecutionRegistryMock.Expect(objRef).Return(
					executionregistry.NewExecutionRegistryMock(mc).FindRequestLoopMock.Return(true),
				),
				ResultsMatcher: nil,
				lr: &LogicRunner{
					ArtifactManager: artifacts.NewClientMock(mc),
				},
			},
			Message: payload.Meta{},
			Parcel:  nil,
		}

		msg := payload.CallMethod{
			Request: &record.IncomingRequest{
				CallType: record.CTMethod,
				Object:   &objRef,
			},
		}

		expectedReply := &reply.CallMethod{Object: &objRef, Result: []byte{3, 2, 1}}
		gotReply, err := handler.handleActual(ctx, msg, fm)
		require.NoError(t, err)
		require.Equal(t, expectedReply, gotReply)

	})
}
