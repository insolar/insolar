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

package handles

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/executionarchive"
	"github.com/insolar/insolar/logicrunner/procs"
	"github.com/insolar/insolar/logicrunner/requestsexecutor"
	"github.com/insolar/insolar/logicrunner/statestorage"
	"github.com/insolar/insolar/logicrunner/writecontroller"
	"github.com/insolar/insolar/testutils"
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
						StateStorage: statestorage.NewStateStorageMock(t).
							GetExecutionArchiveMock.Expect(obj).
							Return(
								executionarchive.NewExecutionArchiveMock(t).
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
						StateStorage: statestorage.NewStateStorageMock(t).
							GetExecutionArchiveMock.Expect(obj).
							Return(
								executionarchive.NewExecutionArchiveMock(t).
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
						StateStorage: statestorage.NewStateStorageMock(t).
							GetExecutionArchiveMock.Expect(obj).
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
		defer mc.Wait(time.Second)

		objRef := gen.Reference()

		fm := flow.NewFlowMock(mc)
		fm.ProcedureMock.Set(func(ctx context.Context, proc flow.Procedure, cancelable bool) (err error) {
			switch p := proc.(type) {
			case *procs.CheckOurRole:
				return nil
			case *procs.RegisterIncomingRequest:
				p.SetResult(&payload.RequestInfo{
					RequestID: gen.ID(),
					ObjectID:  *objRef.Record(),
				})
				return nil
			case *procs.AddFreshRequest:
				return nil
			default:
				t.Fatalf("Unknown procedure: %T", proc)
			}
			return nil
		})

		handler := HandleCall{
			dep: &Dependencies{
				Publisher: nil,
				StateStorage: statestorage.NewStateStorageMock(mc).
					GetExecutionArchiveMock.Expect(objRef).Return(
					executionarchive.NewExecutionArchiveMock(mc).FindRequestLoopMock.Return(false),
				).
					UpsertExecutionStateMock.Expect(objRef).Return(nil),
				ResultsMatcher:  nil,
				ArtifactManager: artifacts.NewClientMock(mc),
				Sender:          nil,
				JetStorage:      nil,
				WriteAccessor:   writecontroller.NewAccessorMock(mc).BeginMock.Return(func() {}, nil),
			},
			Message: payload.Meta{},
			Parcel:  nil,
		}

		msg := message.CallMethod{
			IncomingRequest: record.IncomingRequest{
				CallType: record.CTMethod,
				Object:   &objRef,
			},
		}

		reply, err := handler.handleActual(ctx, &msg, fm)
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
			case *procs.CheckOurRole:
				return nil
			case *procs.RegisterIncomingRequest:
				p.SetResult(&payload.RequestInfo{RequestID: gen.ID(), ObjectID: *objRef.Record()})
				return nil
			case *procs.AddFreshRequest:
				return nil
			default:
				t.Fatalf("Unknown procedure: %T", proc)
			}
			return nil
		})

		handler := HandleCall{
			dep: &Dependencies{
				Publisher: nil,
				StateStorage: statestorage.NewStateStorageMock(mc).
					GetExecutionArchiveMock.Expect(objRef).Return(
					executionarchive.NewExecutionArchiveMock(mc).FindRequestLoopMock.Return(false),
				),
				ResultsMatcher:  nil,
				ArtifactManager: artifacts.NewClientMock(mc),
				MessageBus: testutils.NewMessageBusMock(mc).SendMock.Set(
					func(_ context.Context, m1 insolar.Message, _ *insolar.MessageSendOptions) (insolar.Reply, error) {
						assert.IsType(t, &message.AdditionalCallFromPreviousExecutor{}, m1)
						return nil, nil
					}),
				Sender:     nil,
				JetStorage: nil,
				WriteAccessor: writecontroller.NewAccessorMock(mc).
					BeginMock.Return(func() {}, writecontroller.ErrWriteClosed),
			},
			Message: payload.Meta{},
			Parcel:  nil,
		}

		msg := message.CallMethod{
			IncomingRequest: record.IncomingRequest{
				CallType: record.CTMethod,
				Object:   &objRef,
			},
		}

		reply, err := handler.handleActual(ctx, &msg, fm)
		assert.NotNil(t, reply)
		assert.NoError(t, err)
	})

	t.Run("failed to authorize", func(t *testing.T) {
		t.Parallel()

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)
		defer mc.Wait(time.Second)

		fm := flow.NewFlowMock(mc)

		fm.ProcedureMock.Set(func(ctx context.Context, proc flow.Procedure, cancelable bool) (err error) {
			switch proc.(type) {
			case *procs.CheckOurRole:
				return procs.ErrCantExecute
			case *procs.RegisterIncomingRequest:
				t.Fatalf("Shouldn't be called: %T", proc)
			case *procs.AddFreshRequest:
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
				StateStorage:    statestorage.NewStateStorageMock(mc),
				ResultsMatcher:  nil,
				ArtifactManager: artifacts.NewClientMock(mc),
				MessageBus:      testutils.NewMessageBusMock(mc),
				Sender:          nil,
				JetStorage:      nil,
				WriteAccessor:   writecontroller.NewAccessorMock(mc),
			},
			Message: payload.Meta{},
			Parcel:  nil,
		}

		msg := message.CallMethod{
			IncomingRequest: record.IncomingRequest{
				CallType: record.CTMethod,
				Object:   &objRef,
			},
		}

		reply, err := handler.handleActual(ctx, &msg, fm)
		assert.Nil(t, reply)
		assert.EqualError(t, err, flow.ErrCancelled.Error())
	})

	t.Run("failed to register incoming request", func(t *testing.T) {
		t.Parallel()

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)
		defer mc.Wait(time.Second)

		fm := flow.NewFlowMock(mc)

		fm.ProcedureMock.Set(func(ctx context.Context, proc flow.Procedure, cancelable bool) (err error) {
			switch proc.(type) {
			case *procs.CheckOurRole:
				return nil
			case *procs.RegisterIncomingRequest:
				return flow.ErrCancelled
			case *procs.AddFreshRequest:
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
				MessageBus:      testutils.NewMessageBusMock(mc),
				Sender:          nil,
				JetStorage:      nil,
				WriteAccessor:   writecontroller.NewAccessorMock(mc),
			},
			Message: payload.Meta{},
			Parcel:  nil,
		}

		msg := message.CallMethod{
			IncomingRequest: record.IncomingRequest{
				CallType: record.CTMethod,
				Object:   &objRef,
			},
		}

		reply, err := handler.handleActual(ctx, &msg, fm)
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
			case *procs.CheckOurRole:
				return nil
			case *procs.RegisterIncomingRequest:
				return flow.ErrCancelled
			case *procs.AddFreshRequest:
				t.Fatalf("Shouldn't be called: %T", proc)
			default:
				t.Fatalf("Unknown procedure: %T", proc)
			}
			return nil
		})

		handler := HandleCall{
			dep: &Dependencies{
				Publisher:       nil,
				StateStorage:    statestorage.NewStateStorageMock(mc),
				ResultsMatcher:  nil,
				ArtifactManager: artifacts.NewClientMock(mc),
				MessageBus:      testutils.NewMessageBusMock(mc),
				Sender:          nil,
				JetStorage:      nil,
				WriteAccessor:   writecontroller.NewAccessorMock(mc),
			},
			Message: payload.Meta{},
			Parcel:  nil,
		}

		msg := message.CallMethod{
			IncomingRequest: record.IncomingRequest{
				CallType: record.CTMethod,
				Object:   nil,
			},
		}

		reply, err := handler.handleActual(ctx, &msg, fm)
		assert.Nil(t, reply)
		assert.Error(t, err)
	})

	t.Run("write accessor failed to fetch lock AND archive is empty after on pulse", func(t *testing.T) {
		t.Parallel()

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)
		defer mc.Wait(time.Second)

		objRef := gen.Reference()

		fm := flow.NewFlowMock(mc)
		fm.ProcedureMock.Set(func(ctx context.Context, proc flow.Procedure, cancelable bool) (err error) {
			switch p := proc.(type) {
			case *procs.CheckOurRole:
				return nil
			case *procs.RegisterIncomingRequest:
				p.SetResult(&payload.RequestInfo{RequestID: gen.ID(), ObjectID: *objRef.Record()})
				return nil
			case *procs.AddFreshRequest:
				return nil
			default:
				t.Fatalf("Unknown procedure: %T", proc)
			}
			return nil
		})

		handler := HandleCall{
			dep: &Dependencies{
				Publisher: nil,
				StateStorage: statestorage.NewStateStorageMock(mc).
					GetExecutionArchiveMock.Expect(objRef).Return(nil),
				ResultsMatcher:  nil,
				ArtifactManager: artifacts.NewClientMock(mc),
				MessageBus: testutils.NewMessageBusMock(mc).SendMock.Set(
					func(_ context.Context, m1 insolar.Message, _ *insolar.MessageSendOptions) (insolar.Reply, error) {
						assert.IsType(t, &message.AdditionalCallFromPreviousExecutor{}, m1)
						return nil, nil
					}),
				Sender:        nil,
				JetStorage:    nil,
				WriteAccessor: writecontroller.NewAccessorMock(mc).BeginMock.Return(func() {}, writecontroller.ErrWriteClosed),
			},
			Message: payload.Meta{},
			Parcel:  nil,
		}

		msg := message.CallMethod{
			IncomingRequest: record.IncomingRequest{
				CallType: record.CTMethod,
				Object:   &objRef,
			},
		}

		reply, err := handler.handleActual(ctx, &msg, fm)
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
			case *procs.CheckOurRole:
				return nil
			case *procs.RegisterIncomingRequest:
				p.SetResult(&payload.RequestInfo{
					RequestID: *reqRef.Record(),
					ObjectID:  *objRef.Record(),
					Request:   []byte{1, 2, 3},
					Result:    matRecordSerialized,
				})
				return nil
			case *procs.AddFreshRequest:
				return nil
			default:
				t.Fatalf("Unknown procedure: %T", proc)
			}
			return nil
		})

		handler := HandleCall{
			dep: &Dependencies{
				ArtifactManager:  artifacts.NewClientMock(mc),
				RequestsExecutor: requestsexecutor.NewRequestsExecutorMock(mc).SendReplyMock.Return(),
			},
			Message: payload.Meta{},
			Parcel:  nil,
		}

		msg := message.CallMethod{
			IncomingRequest: record.IncomingRequest{
				CallType: record.CTMethod,
				Object:   &objRef,
			},
		}

		gotReply, err := handler.handleActual(ctx, &msg, fm)
		require.NoError(t, err)
		require.Equal(t, &reply.RegisterRequest{Request: reqRef}, gotReply)
	})
}
