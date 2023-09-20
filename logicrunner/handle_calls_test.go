package logicrunner

import (
	"context"
	"testing"
	"time"

	wmMessage "github.com/ThreeDotsLabs/watermill/message"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/testutils"

	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/writecontroller"
)

func TestHandleCall_Present(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		if useLeakTest {
			defer testutils.LeakTester(t)
		} else {
			t.Parallel()
		}

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
			default:
				t.Fatalf("Unknown procedure: %T", proc)
			}
			return nil
		})

		handler := HandleCall{
			dep: &Dependencies{
				Publisher: nil,
				StateStorage: NewStateStorageMock(mc).
					UpsertExecutionStateMock.Expect(objRef).Return(NewExecutionBrokerIMock(mc).HasMoreRequestsMock.Return()),
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
		if useLeakTest {
			defer testutils.LeakTester(t)
		} else {
			t.Parallel()
		}

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
						c := make(chan *wmMessage.Message, 1)
						c <- &wmMessage.Message{Payload: payload.MustMarshal(&payload.ID{})}
						return c, func() {}
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
		if useLeakTest {
			defer testutils.LeakTester(t)
		} else {
			t.Parallel()
		}

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)

		fm := flow.NewFlowMock(mc)
		fm.ProcedureMock.Set(func(ctx context.Context, proc flow.Procedure, cancelable bool) (err error) {
			switch proc.(type) {
			case *CheckOurRole:
				return ErrCantExecute
			case *RegisterIncomingRequest:
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
		if useLeakTest {
			defer testutils.LeakTester(t)
		} else {
			t.Parallel()
		}

		ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
		mc := minimock.NewController(t)

		fm := flow.NewFlowMock(mc)
		fm.ProcedureMock.Set(func(ctx context.Context, proc flow.Procedure, cancelable bool) (err error) {
			switch proc.(type) {
			case *CheckOurRole:
				return nil
			case *RegisterIncomingRequest:
				return flow.ErrCancelled
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

	t.Run("write accessor failed to fetch lock AND registry is empty after on pulse", func(t *testing.T) {
		if useLeakTest {
			defer testutils.LeakTester(t)
		} else {
			t.Parallel()
		}

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
						c := make(chan *wmMessage.Message, 1)
						c <- &wmMessage.Message{Payload: payload.MustMarshal(&payload.ID{})}
						return c, func() {}
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

	t.Run("object not found during request registration", func(t *testing.T) {
		if useLeakTest {
			defer testutils.LeakTester(t)
		} else {
			t.Parallel()
		}

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
		if useLeakTest {
			defer testutils.LeakTester(t)
		} else {
			t.Parallel()
		}

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
