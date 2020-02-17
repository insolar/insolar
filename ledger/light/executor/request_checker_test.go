// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package executor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/pulse"
	"github.com/insolar/insolar/testutils"
)

func TestRequestCheckerDefault_CheckRequest(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	pcs := testutils.NewPlatformCryptographyScheme()

	var (
		filament *executor.FilamentCalculatorMock
		jets     *jet.CoordinatorMock
		fetcher  *executor.JetFetcherMock
		sender   *bus.SenderMock
		checker  *executor.RequestCheckerDefault
	)

	setup := func() {
		filament = executor.NewFilamentCalculatorMock(mc)
		jets = jet.NewCoordinatorMock(mc)
		fetcher = executor.NewJetFetcherMock(mc)
		sender = bus.NewSenderMock(mc)
		checker = executor.NewRequestChecker(filament, jets, fetcher, pcs, sender)
	}

	t.Run("invalid request returns error", func(t *testing.T) {
		setup()
		defer mc.Finish()

		req := record.NewRequestMock(mc)
		req.ValidateMock.Return(errors.New("test error"))

		err := checker.CheckRequest(ctx, gen.ID(), req)
		coded, ok := err.(*payload.CodedError)
		require.True(t, ok, "should be coded error")
		assert.Equal(t, payload.CodeRequestInvalid, coded.Code)
	})

	t.Run("incoming, reason is empty returns error", func(t *testing.T) {
		setup()
		defer mc.Finish()

		req := record.IncomingRequest{
			Caller:  gen.ReferenceWithPulse(pulse.MinTimePulse + 1),
			APINode: gen.Reference(),
			Reason:  insolar.Reference{},
		}

		err := checker.CheckRequest(ctx, gen.IDWithPulse(pulse.MinTimePulse+2), &req)
		assert.Error(t, err)
	})

	t.Run("reason is older than request returns error", func(t *testing.T) {
		setup()
		defer mc.Finish()

		req := record.NewRequestMock(mc)
		req.ValidateMock.Return(nil)
		req.ReasonRefMock.Return(*insolar.NewRecordReference(gen.IDWithPulse(pulse.MinTimePulse + 2)))

		err := checker.CheckRequest(ctx, gen.IDWithPulse(pulse.MinTimePulse+1), req)
		coded, ok := err.(*payload.CodedError)
		require.True(t, ok, "should be coded error")
		assert.Equal(t, payload.CodeRequestInvalid, coded.Code)
	})

	t.Run("incoming API request is ok", func(t *testing.T) {
		setup()
		defer mc.Finish()

		filament.OpenedRequestsMock.Return(nil, nil)

		objectRef := gen.Reference()
		req := record.IncomingRequest{
			Caller:  gen.ReferenceWithPulse(pulse.MinTimePulse + 1),
			APINode: gen.Reference(),
			Reason:  gen.ReferenceWithPulse(pulse.MinTimePulse + 1),
			Object:  &objectRef,
		}

		err := checker.CheckRequest(ctx, gen.IDWithPulse(pulse.MinTimePulse+2), &req)
		assert.Nil(t, err)
	})

	t.Run("incoming network reason check is ok (creation request)", func(t *testing.T) {
		setup()
		defer mc.Finish()

		requestID := gen.IDWithPulse(pulse.MinTimePulse + 2)
		reasonObjectRef := gen.ReferenceWithPulse(pulse.MinTimePulse + 1)
		jetID := gen.JetID()
		nodeRef := gen.Reference()
		req := record.IncomingRequest{
			Caller:   reasonObjectRef,
			CallType: record.CTSaveAsChild,
			Reason:   gen.ReferenceWithPulse(pulse.MinTimePulse + 1),
		}

		fetcher.FetchMock.Inspect(func(_ context.Context, target insolar.ID, pulse insolar.PulseNumber) {
			require.Equal(t, reasonObjectRef.GetLocal(), &target)
			require.Equal(t, requestID.Pulse(), pulse)
		}).Return((*insolar.ID)(&jetID), nil)

		jets.LightExecutorForJetMock.Inspect(func(_ context.Context, j insolar.ID, pulse insolar.PulseNumber) {
			require.Equal(t, insolar.ID(jetID), j)
			require.Equal(t, requestID.Pulse(), pulse)
		}).Return(&nodeRef, nil)

		sender.SendTargetMock.Set(func(_ context.Context, msg *message.Message, target insolar.Reference) (<-chan *message.Message, func()) {
			ch := make(chan *message.Message, 1)
			ch <- payload.MustNewMessage(&payload.Meta{
				Payload: payload.MustMarshal(&payload.RequestInfo{
					Request: func() []byte {
						material := record.Material{Virtual: record.Wrap(&record.IncomingRequest{})}
						buf, err := material.Marshal()
						if err != nil {
							panic(err)
						}
						return buf
					}(),
					OldestMutable: true,
				}),
			})
			return ch, func() {}
		})

		err := checker.CheckRequest(ctx, requestID, &req)
		assert.Nil(t, err)
	})

	t.Run("incoming network reason check failed, request not found error", func(t *testing.T) {
		setup()
		defer mc.Finish()

		requestID := gen.IDWithPulse(pulse.MinTimePulse + 2)
		reasonObjectRef := gen.ReferenceWithPulse(pulse.MinTimePulse + 1)
		jetID := gen.JetID()
		nodeRef := gen.Reference()
		req := record.IncomingRequest{
			Caller:   reasonObjectRef,
			CallType: record.CTSaveAsChild,
			Reason:   gen.ReferenceWithPulse(pulse.MinTimePulse + 1),
		}

		fetcher.FetchMock.Inspect(func(_ context.Context, target insolar.ID, pulse insolar.PulseNumber) {
			require.Equal(t, reasonObjectRef.GetLocal(), &target)
			require.Equal(t, requestID.Pulse(), pulse)
		}).Return((*insolar.ID)(&jetID), nil)

		jets.LightExecutorForJetMock.Inspect(func(_ context.Context, j insolar.ID, pulse insolar.PulseNumber) {
			require.Equal(t, insolar.ID(jetID), j)
			require.Equal(t, requestID.Pulse(), pulse)
		}).Return(&nodeRef, nil)

		sender.SendTargetMock.Set(func(_ context.Context, msg *message.Message, target insolar.Reference) (<-chan *message.Message, func()) {
			ch := make(chan *message.Message, 1)
			ch <- payload.MustNewMessage(&payload.Meta{
				Payload: payload.MustMarshal(&payload.Error{
					Text: "not found",
					Code: payload.CodeRequestNotFound,
				}),
			})
			return ch, func() {}
		})

		err := checker.CheckRequest(ctx, requestID, &req)
		require.Error(t, err)
		insError, ok := errors.Cause(err).(*payload.CodedError)
		require.True(t, ok)
		require.Equal(t, payload.CodeRequestNotFound, insError.GetCode())
	})

	t.Run("incoming local reason check is ok", func(t *testing.T) {
		setup()
		defer mc.Finish()

		requestID := gen.IDWithPulse(pulse.MinTimePulse + 2)
		reasonObjectRef := gen.ReferenceWithPulse(pulse.MinTimePulse + 1)
		reasonRef := gen.ReferenceWithPulse(pulse.MinTimePulse + 1)
		req := record.IncomingRequest{
			Caller: reasonObjectRef,
			Object: &reasonObjectRef,
			Reason: reasonRef,
		}
		filament.OpenedRequestsMock.Inspect(func(_ context.Context, pn insolar.PulseNumber, objectID insolar.ID, pendingOnly bool) {
			require.Equal(t, requestID.Pulse(), pn)
			require.Equal(t, *reasonObjectRef.GetLocal(), objectID)
			require.False(t, pendingOnly)
		}).Return(nil, nil)

		filament.RequestInfoMock.Set(func(ctx context.Context, objectID insolar.ID, reqID insolar.ID, pulse insolar.PulseNumber) (requestInfo executor.FilamentsRequestInfo, err error) {

			require.Equal(t, reasonObjectRef.GetLocal(), &objectID)
			require.Equal(t, reasonRef.GetLocal(), &reqID)
			require.Equal(t, requestID.Pulse(), pulse)

			request := executor.FilamentsRequestInfo{
				Request: &record.CompositeFilamentRecord{
					Record: record.Material{Virtual: record.Wrap(&record.IncomingRequest{})},
				},
				OldestMutable: true,
			}
			return request, nil
		})

		err := checker.CheckRequest(ctx, requestID, &req)
		assert.Nil(t, err)
	})

	t.Run("incoming local reason check failed, returns request not found error", func(t *testing.T) {
		setup()
		defer mc.Finish()

		requestID := gen.IDWithPulse(pulse.MinTimePulse + 2)
		reasonObjectRef := gen.ReferenceWithPulse(pulse.MinTimePulse + 1)
		reasonRef := gen.ReferenceWithPulse(pulse.MinTimePulse + 1)
		req := record.IncomingRequest{
			Caller: reasonObjectRef,
			Object: &reasonObjectRef,
			Reason: reasonRef,
		}

		filament.OpenedRequestsMock.Inspect(func(_ context.Context, pn insolar.PulseNumber, objectID insolar.ID, pendingOnly bool) {
			require.Equal(t, requestID.Pulse(), pn)
			require.Equal(t, *reasonObjectRef.GetLocal(), objectID)
			require.False(t, pendingOnly)
		}).Return(nil, nil)

		filament.RequestInfoMock.Set(func(_ context.Context, objectID insolar.ID, reqID insolar.ID, pulse insolar.PulseNumber) (requestInfo executor.FilamentsRequestInfo, err error) {
			require.Equal(t, reasonObjectRef.GetLocal(), &objectID)
			require.Equal(t, reasonRef.GetLocal(), &reqID)
			require.Equal(t, requestID.Pulse(), pulse)

			return executor.FilamentsRequestInfo{}, &payload.CodedError{
				Text: fmt.Sprintf("requestInfo not found request %s", requestID.DebugString()),
				Code: payload.CodeRequestNotFound,
			}
		})

		err := checker.CheckRequest(ctx, requestID, &req)
		require.Error(t, err)
		insError, ok := errors.Cause(err).(*payload.CodedError)
		require.True(t, ok)
		require.Equal(t, payload.CodeRequestNotFound, insError.GetCode())
	})

	t.Run("incoming reason is closed for regular request", func(t *testing.T) {
		setup()
		defer mc.Finish()

		requestID := gen.IDWithPulse(pulse.MinTimePulse + 2)
		reasonObjectRef := gen.ReferenceWithPulse(pulse.MinTimePulse + 1)
		reasonRef := gen.ReferenceWithPulse(pulse.MinTimePulse + 1)
		req := record.IncomingRequest{
			Caller: reasonObjectRef,
			Object: &reasonObjectRef,
			Reason: reasonRef,
		}

		filament.OpenedRequestsMock.Inspect(func(_ context.Context, pn insolar.PulseNumber, objectID insolar.ID, pendingOnly bool) {
			require.Equal(t, requestID.Pulse(), pn)
			require.Equal(t, *reasonObjectRef.GetLocal(), objectID)
			require.False(t, pendingOnly)
		}).Return(nil, nil)

		filament.RequestInfoMock.Set(func(_ context.Context, objectID insolar.ID, reqID insolar.ID, pulse insolar.PulseNumber) (requestInfo executor.FilamentsRequestInfo, err error) {
			require.Equal(t, reasonObjectRef.GetLocal(), &objectID)
			require.Equal(t, reasonRef.GetLocal(), &reqID)
			require.Equal(t, requestID.Pulse(), pulse)

			request := record.CompositeFilamentRecord{
				Record: record.Material{Virtual: record.Wrap(&record.IncomingRequest{})},
			}
			result := record.CompositeFilamentRecord{
				Record: record.Material{Virtual: record.Wrap(&record.Result{})},
			}
			requestInfo = executor.FilamentsRequestInfo{
				Request: &request,
				Result:  &result,
			}
			return requestInfo, nil
		})

		err := checker.CheckRequest(ctx, requestID, &req)
		require.Error(t, err)
		insError, ok := errors.Cause(err).(*payload.CodedError)
		require.True(t, ok)
		require.Equal(t, payload.CodeReasonIsWrong, insError.GetCode())
	})

	t.Run("incoming reason is not closed for detached request", func(t *testing.T) {
		setup()
		defer mc.Finish()

		requestID := gen.IDWithPulse(pulse.MinTimePulse + 2)
		reasonObjectRef := gen.ReferenceWithPulse(pulse.MinTimePulse + 1)
		reasonRef := gen.ReferenceWithPulse(pulse.MinTimePulse + 1)
		req := record.IncomingRequest{
			Caller:     reasonObjectRef,
			Object:     &reasonObjectRef,
			Reason:     reasonRef,
			ReturnMode: record.ReturnSaga,
		}

		filament.OpenedRequestsMock.Inspect(func(_ context.Context, pn insolar.PulseNumber, objectID insolar.ID, pendingOnly bool) {
			require.Equal(t, requestID.Pulse(), pn)
			require.Equal(t, *reasonObjectRef.GetLocal(), objectID)
			require.False(t, pendingOnly)
		}).Return(nil, nil)

		filament.RequestInfoMock.Set(func(_ context.Context, objectID insolar.ID, reqID insolar.ID, pulse insolar.PulseNumber) (requestInfo executor.FilamentsRequestInfo, err error) {
			require.Equal(t, reasonObjectRef.GetLocal(), &objectID)
			require.Equal(t, reasonRef.GetLocal(), &reqID)
			require.Equal(t, requestID.Pulse(), pulse)

			request := executor.FilamentsRequestInfo{
				Request: &record.CompositeFilamentRecord{
					Record: record.Material{Virtual: record.Wrap(&record.IncomingRequest{})},
				},
			}
			return request, nil
		})

		err := checker.CheckRequest(ctx, requestID, &req)
		require.Error(t, err)
		insError, ok := errors.Cause(err).(*payload.CodedError)
		require.True(t, ok)
		require.Equal(t, payload.CodeReasonIsWrong, insError.GetCode())
	})

	t.Run("outgoing reason is not found returns error", func(t *testing.T) {
		setup()
		defer mc.Finish()

		requestID := gen.IDWithPulse(pulse.MinTimePulse + 2)
		reasonRef := gen.ReferenceWithPulse(pulse.MinTimePulse + 1)
		req := record.OutgoingRequest{
			Reason: reasonRef,
		}

		filament.OpenedRequestsMock.Return(nil, nil)

		err := checker.CheckRequest(ctx, requestID, &req)
		require.Error(t, err)
		insError, ok := errors.Cause(err).(*payload.CodedError)
		require.True(t, ok)
		require.Equal(t, payload.CodeReasonIsWrong, insError.GetCode())
	})

	t.Run("outgoing, reason is immutable does not have to be the latest", func(t *testing.T) {
		setup()
		defer mc.Finish()

		requestID := gen.IDWithPulse(pulse.MinTimePulse + 2)
		reasonRef := gen.ReferenceWithPulse(pulse.MinTimePulse + 1)
		req := record.OutgoingRequest{
			Reason: reasonRef,
		}

		filament.OpenedRequestsMock.Set(func(_ context.Context, pulse insolar.PulseNumber, objectID insolar.ID, pendingOnly bool) (ca1 []record.CompositeFilamentRecord, err error) {
			require.Equal(t, requestID.Pulse(), pulse)

			req := record.CompositeFilamentRecord{
				RecordID: *reasonRef.GetLocal(),
				Record: record.Material{
					Virtual: record.Wrap(&record.IncomingRequest{
						Immutable: true,
					}),
				},
			}
			return []record.CompositeFilamentRecord{
				// garbage
				{
					RecordID: gen.ID(),
					Record: record.Material{
						Virtual: record.Wrap(&record.OutgoingRequest{
							ReturnMode: record.ReturnSaga,
						}),
					},
				},
				// oldest mutable
				{
					RecordID: gen.ID(),
					Record: record.Material{
						Virtual: record.Wrap(&record.IncomingRequest{
							Immutable: false,
						}),
					},
				},
				// garbage
				{
					RecordID: gen.ID(),
					Record: record.Material{
						Virtual: record.Wrap(&record.IncomingRequest{
							Immutable: true,
						}),
					},
				},
				// reason
				req,
			}, nil
		})
		err := checker.CheckRequest(ctx, requestID, &req)
		require.NoError(t, err)
	})

	t.Run("outgoing, reason is the oldest mutable", func(t *testing.T) {
		setup()
		defer mc.Finish()

		requestID := gen.IDWithPulse(pulse.MinTimePulse + 2)
		reasonRef := gen.ReferenceWithPulse(pulse.MinTimePulse + 1)
		req := record.OutgoingRequest{
			Reason: reasonRef,
		}

		filament.OpenedRequestsMock.Set(func(_ context.Context, pulse insolar.PulseNumber, objectID insolar.ID, pendingOnly bool) (ca1 []record.CompositeFilamentRecord, err error) {
			require.Equal(t, requestID.Pulse(), pulse)

			req := record.CompositeFilamentRecord{
				RecordID: *reasonRef.GetLocal(),
				Record: record.Material{
					Virtual: record.Wrap(&record.IncomingRequest{
						Immutable: false,
					}),
				},
			}
			return []record.CompositeFilamentRecord{
				// garbage
				{
					RecordID: gen.ID(),
					Record: record.Material{
						Virtual: record.Wrap(&record.OutgoingRequest{
							ReturnMode: record.ReturnSaga,
						}),
					},
				},
				{
					RecordID: gen.ID(),
					Record: record.Material{
						Virtual: record.Wrap(&record.IncomingRequest{
							Immutable: true,
						}),
					},
				},
				// reason
				req,
			}, nil
		})
		err := checker.CheckRequest(ctx, requestID, &req)
		assert.Nil(t, err)
	})

	t.Run("outgoing, reason is the oldest mutable", func(t *testing.T) {
		setup()
		defer mc.Finish()

		requestID := gen.IDWithPulse(pulse.MinTimePulse + 2)
		reasonRef := gen.ReferenceWithPulse(pulse.MinTimePulse + 1)
		req := record.OutgoingRequest{
			Reason: reasonRef,
		}

		filament.OpenedRequestsMock.Set(func(_ context.Context, pulse insolar.PulseNumber, objectID insolar.ID, pendingOnly bool) (ca1 []record.CompositeFilamentRecord, err error) {
			require.Equal(t, requestID.Pulse(), pulse)

			req := record.CompositeFilamentRecord{
				RecordID: *reasonRef.GetLocal(),
				Record: record.Material{
					Virtual: record.Wrap(&record.IncomingRequest{
						Immutable: false,
					}),
				},
			}
			return []record.CompositeFilamentRecord{
				// garbage
				{
					RecordID: gen.ID(),
					Record: record.Material{
						Virtual: record.Wrap(&record.OutgoingRequest{
							ReturnMode: record.ReturnSaga,
						}),
					},
				},
				{
					RecordID: gen.ID(),
					Record: record.Material{
						Virtual: record.Wrap(&record.IncomingRequest{
							Immutable: true,
						}),
					},
				},
				// reason
				req,
			}, nil
		})
		err := checker.CheckRequest(ctx, requestID, &req)
		assert.Nil(t, err)
	})

	t.Run("outgoing is ok", func(t *testing.T) {
		setup()
		defer mc.Finish()

		requestID := gen.IDWithPulse(pulse.MinTimePulse + 2)
		reasonRef := gen.ReferenceWithPulse(pulse.MinTimePulse + 1)
		req := record.OutgoingRequest{
			Reason: reasonRef,
		}

		filament.OpenedRequestsMock.Set(func(_ context.Context, pulse insolar.PulseNumber, objectID insolar.ID, pendingOnly bool) (ca1 []record.CompositeFilamentRecord, err error) {
			require.Equal(t, requestID.Pulse(), pulse)

			req := record.CompositeFilamentRecord{
				RecordID: *reasonRef.GetLocal(),
				Record: record.Material{
					Virtual: record.Wrap(&record.IncomingRequest{}),
				},
			}
			return []record.CompositeFilamentRecord{req}, nil
		})

		err := checker.CheckRequest(ctx, requestID, &req)
		assert.Nil(t, err)
	})
}
