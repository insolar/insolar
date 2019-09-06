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

package executor_test

import (
	"context"
	"testing"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
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
		assert.Equal(t, uint32(payload.CodeInvalidRequest), coded.Code)
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
		assert.Equal(t, uint32(payload.CodeInvalidRequest), coded.Code)
	})

	t.Run("incoming API request is ok", func(t *testing.T) {
		setup()
		defer mc.Finish()

		req := record.IncomingRequest{
			Caller:  gen.ReferenceWithPulse(pulse.MinTimePulse + 1),
			APINode: gen.Reference(),
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
				}),
			})
			return ch, func() {}
		})

		err := checker.CheckRequest(ctx, requestID, &req)
		assert.Nil(t, err)
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

		filament.RequestInfoMock.Set(func(_ context.Context, objectID insolar.ID, reqID insolar.ID, pulse insolar.PulseNumber) (foundRequest *record.CompositeFilamentRecord, foundResult *record.CompositeFilamentRecord, err error) {
			require.Equal(t, reasonObjectRef.GetLocal(), &objectID)
			require.Equal(t, reasonRef.GetLocal(), &reqID)
			require.Equal(t, requestID.Pulse(), pulse)

			request := record.CompositeFilamentRecord{
				Record: record.Material{Virtual: record.Wrap(&record.IncomingRequest{})},
			}
			return &request, nil, nil
		})

		err := checker.CheckRequest(ctx, requestID, &req)
		assert.Nil(t, err)
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

		filament.RequestInfoMock.Set(func(_ context.Context, objectID insolar.ID, reqID insolar.ID, pulse insolar.PulseNumber) (foundRequest *record.CompositeFilamentRecord, foundResult *record.CompositeFilamentRecord, err error) {
			require.Equal(t, reasonObjectRef.GetLocal(), &objectID)
			require.Equal(t, reasonRef.GetLocal(), &reqID)
			require.Equal(t, requestID.Pulse(), pulse)

			request := record.CompositeFilamentRecord{
				Record: record.Material{Virtual: record.Wrap(&record.IncomingRequest{})},
			}
			result := record.CompositeFilamentRecord{
				Record: record.Material{Virtual: record.Wrap(&record.Result{})},
			}
			return &request, &result, nil
		})

		err := checker.CheckRequest(ctx, requestID, &req)
		assert.EqualError(t, err, "reason request is closed for a regular (not detached) call")
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
			ReturnMode: record.ReturnNoWait,
		}

		filament.RequestInfoMock.Set(func(_ context.Context, objectID insolar.ID, reqID insolar.ID, pulse insolar.PulseNumber) (foundRequest *record.CompositeFilamentRecord, foundResult *record.CompositeFilamentRecord, err error) {
			require.Equal(t, reasonObjectRef.GetLocal(), &objectID)
			require.Equal(t, reasonRef.GetLocal(), &reqID)
			require.Equal(t, requestID.Pulse(), pulse)

			request := record.CompositeFilamentRecord{
				Record: record.Material{Virtual: record.Wrap(&record.IncomingRequest{})},
			}
			return &request, nil, nil
		})

		err := checker.CheckRequest(ctx, requestID, &req)
		assert.EqualError(t, err, "reason request is not closed for a detached call")
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
		coded, ok := err.(*payload.CodedError)
		require.True(t, ok, "should be coded error")
		assert.Equal(t, uint32(payload.CodeReasonNotFound), coded.Code)
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
			}
			return []record.CompositeFilamentRecord{req}, nil
		})

		err := checker.CheckRequest(ctx, requestID, &req)
		assert.Nil(t, err)
	})
}
