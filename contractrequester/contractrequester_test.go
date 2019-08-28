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

package contractrequester

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/bus/meta"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulse"
	"github.com/insolar/insolar/testutils"
)

// ensure that ContractRequester implements insolar.ContractRequester
var _ insolar.ContractRequester = &ContractRequester{}

func TestNew(t *testing.T) {
	sender := bus.NewSenderMock(t)
	pulseAccessor := insolarPulse.NewAccessorMock(t)
	jetCoordinator := jet.NewCoordinatorMock(t)
	pcs := platformpolicy.NewPlatformCryptographyScheme()

	_, err := New(sender,
		pulseAccessor,
		jetCoordinator,
		pcs)
	require.NoError(t, err)
}

func mockPulseAccessor(t minimock.Tester) insolarPulse.Accessor {
	pulseAccessor := insolarPulse.NewAccessorMock(t)
	currentPulse := pulse.MinTimePulse
	pulseAccessor.LatestMock.Set(func(p context.Context) (r insolar.Pulse, r1 error) {
		return insolar.Pulse{
			PulseNumber:     insolar.PulseNumber(currentPulse),
			NextPulseNumber: insolar.PulseNumber(currentPulse + 1),
		}, nil
	})

	return pulseAccessor
}

func mockJetCoordinator(t minimock.Tester) jet.Coordinator {
	coordinator := jet.NewCoordinatorMock(t)
	coordinator.MeMock.Set(func() (r insolar.Reference) {
		return gen.Reference()
	})
	return coordinator
}

func TestContractRequester_SendRequest(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	ref := gen.Reference()

	cReq, err := New(bus.NewSenderMock(t),
		mockPulseAccessor(mc),
		mockJetCoordinator(mc),
		testutils.NewPlatformCryptographyScheme())
	require.NoError(t, err)

	table := []struct {
		name          string
		resultMessage payload.ReturnResults
		earlyResult   bool
	}{
		{
			name: "success",
			resultMessage: payload.ReturnResults{
				Reply: reply.ToBytes(&reply.CallMethod{}),
			},
		},
		{
			name: "earlyResult",
			resultMessage: payload.ReturnResults{
				Reply: reply.ToBytes(&reply.CallMethod{}),
			},
			earlyResult: true,
		},
		{
			name: "early result, before registration",
			resultMessage: payload.ReturnResults{
				Reply: reply.ToBytes(&reply.CallMethod{}),
			},
		},
	}

	for _, test := range table {
		test := test
		t.Run(test.name, func(t *testing.T) {

			cReq.Sender = bus.NewSenderMock(t).SendRoleMock.
				Set(func(ctx context.Context, msg *message.Message, role insolar.DynamicRole, obj insolar.Reference) (<-chan *message.Message, func()) {
					data, err := payload.Unmarshal(msg.Payload)
					require.NoError(t, err)

					request := data.(*payload.CallMethod).Request

					hash, err := cReq.calcRequestHash(*request)
					require.NoError(t, err)
					requestRef := insolar.NewReference(*insolar.NewID(gen.PulseNumber(), hash[:]))

					resultSender := func() {
						res := test.resultMessage
						res.RequestRef = *requestRef
						cReq.result(ctx, &res)
					}

					resChan := make(chan *message.Message)

					res, err := serializeReply(bus.ReplyAsMessage(ctx, &reply.RegisterRequest{Request: *requestRef}))
					require.NoError(t, err)

					go func() {
						resChan <- res
					}()

					if test.earlyResult {
						resultSender()
					} else {
						go func() {
							time.Sleep(time.Millisecond)
							resultSender()
						}()
					}

					return resChan, func() {
						close(resChan)
					}
				})

			result, err := cReq.SendRequest(ctx, &ref, "TestMethod", []interface{}{})
			require.NoError(t, err)
			require.Equal(t, &reply.CallMethod{}, result)
		})
	}
}

func TestContractRequester_Call_Timeout(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	cReq, err := New(
		bus.NewSenderMock(t),
		mockPulseAccessor(mc),
		jet.NewCoordinatorMock(t),
		testutils.NewPlatformCryptographyScheme(),
	)
	require.NoError(t, err)

	cReq.callTimeout = 1 * time.Nanosecond

	ref := gen.Reference()
	prototypeRef := gen.Reference()
	method := testutils.RandomString()

	request := &record.IncomingRequest{
		Caller:    gen.Reference(),
		Object:    &ref,
		Prototype: &prototypeRef,
		Method:    method,
		Arguments: insolar.Arguments{},
	}

	cReq.Sender = bus.NewSenderMock(t).SendRoleMock.Set(
		func(ctx context.Context, msg *message.Message, role insolar.DynamicRole, obj insolar.Reference) (<-chan *message.Message, func()) {
			resChan := make(chan *message.Message)

			data, err := payload.Unmarshal(msg.Payload)
			require.NoError(t, err)

			request := data.(*payload.CallMethod).Request

			hash, err := cReq.calcRequestHash(*request)
			require.NoError(t, err)
			requestRef := insolar.NewReference(*insolar.NewID(pulse.MinTimePulse, hash[:]))

			res, err := serializeReply(bus.ReplyAsMessage(ctx, &reply.RegisterRequest{
				Request: *requestRef,
			}))
			require.NoError(t, err)

			go func() {
				resChan <- res
			}()
			return resChan, func() {
				close(resChan)
			}
		})

	msg := &payload.CallMethod{
		Request: request,
	}

	_, _, err = cReq.Call(ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "canceled")
	require.Contains(t, err.Error(), "timeout")
}

func TestReceiveResult_UnwantedResult(t *testing.T) {
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*10)
	defer cancelFunc()

	cReq, err := New(
		bus.NewSenderMock(t),
		insolarPulse.NewAccessorMock(t),
		jet.NewCoordinatorMock(t),
		testutils.NewPlatformCryptographyScheme(),
	)
	require.NoError(t, err)

	mc := minimock.NewController(t)
	defer mc.Finish()

	reqRef := gen.Reference()
	var reqHash [insolar.RecordHashSize]byte
	copy(reqHash[:], reqRef.GetLocal().Hash())

	msgPayload := &payload.ReturnResults{
		RequestRef: reqRef,
	}

	msg := payload.MustNewMessage(msgPayload)

	sp, err := instracer.Serialize(ctx)
	require.NoError(t, err)
	msg.Metadata.Set(meta.SpanData, string(sp))

	msg.Metadata.Set(meta.TraceID, "OK_unwanted_test")
	cReq.Sender = bus.NewSenderMock(t).ReplyMock.Set(
		func(_ context.Context, origin payload.Meta, replyMsg *message.Message) {
			replyData, err := reply.Deserialize(bytes.NewBuffer(replyMsg.Payload))
			require.NoError(t, err)
			require.Equal(t, &reply.OK{}, replyData)
		})

	// unexpected result
	res, err := serializeReply(msg)
	require.NoError(t, err)
	err = cReq.ReceiveResult(res)
	require.NoError(t, err)
}

func TestReceiveResult_WantedResult(t *testing.T) {
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*10)
	defer cancelFunc()

	cReq, err := New(
		bus.NewSenderMock(t),
		insolarPulse.NewAccessorMock(t),
		jet.NewCoordinatorMock(t),
		testutils.NewPlatformCryptographyScheme())
	require.NoError(t, err)

	mc := minimock.NewController(t)
	defer mc.Finish()

	reqRef := gen.Reference()
	var reqHash [insolar.RecordHashSize]byte
	copy(reqHash[:], reqRef.GetLocal().Hash())

	msgPayload := &payload.ReturnResults{
		RequestRef: reqRef,
	}

	msg := payload.MustNewMessage(msgPayload)

	sp, err := instracer.Serialize(ctx)
	require.NoError(t, err)
	msg.Metadata.Set(meta.SpanData, string(sp))

	msg.Metadata.Set(meta.TraceID, "OK_wanted_test")
	cReq.Sender = bus.NewSenderMock(t).ReplyMock.Set(
		func(_ context.Context, origin payload.Meta, replyMsg *message.Message) {
			replyData, err := reply.Deserialize(bytes.NewBuffer(replyMsg.Payload))
			require.NoError(t, err)
			require.Equal(t, &reply.OK{}, replyData)
		})

	// expected result
	resChan := make(chan *payload.ReturnResults)
	chanResult := make(chan *payload.ReturnResults)
	cReq.ResultMap[reqHash] = resChan

	go func() {
		chanResult <- <-cReq.ResultMap[reqHash]
	}()

	res, err := serializeReply(msg)
	require.NoError(t, err)
	err = cReq.ReceiveResult(res)

	require.NoError(t, err)
	require.Equal(t, 0, len(cReq.ResultMap))
	require.Equal(t, msgPayload, <-chanResult)
}

func TestReceiveResult_UnwantedResultWithError(t *testing.T) {
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*10)
	defer cancelFunc()

	cReq, err := New(bus.NewSenderMock(t),
		insolarPulse.NewAccessorMock(t),
		jet.NewCoordinatorMock(t),
		testutils.NewPlatformCryptographyScheme())
	require.NoError(t, err)

	mc := minimock.NewController(t)
	defer mc.Finish()

	reqRef := gen.Reference()
	var reqHash [insolar.RecordHashSize]byte
	copy(reqHash[:], reqRef.GetLocal().Hash())

	msgPayload := &payload.ReturnResults{
		RequestRef: reqRef,
	}

	msg := payload.MustNewMessage(msgPayload)

	sp, err := instracer.Serialize(ctx)
	require.NoError(t, err)
	msg.Metadata.Set(meta.SpanData, string(sp))

	msg.Metadata.Set(meta.TraceID, "handle_flow_cancelled")
	cReq.LR = testutils.NewLogicRunnerMock(t).AddUnwantedResponseMock.Set(
		func(ctx context.Context, msg insolar.Payload) error {
			return flow.ErrCancelled
		})
	cReq.Sender = bus.NewSenderMock(t).ReplyMock.Set(
		func(_ context.Context, origin payload.Meta, replyMsg *message.Message) {
			payloadError := &payload.Error{}
			err := payloadError.Unmarshal(replyMsg.Payload)
			require.NoError(t, err)
			require.Equal(t, &payload.Error{
				Polymorph: uint32(payload.TypeError),
				Code:      uint32(payload.CodeFlowCanceled),
				Text:      errors.Wrap(flow.ErrCancelled, "[ ReceiveResult ]").Error(),
			}, payloadError)
		})

	res, err := serializeReply(msg)
	require.NoError(t, err)
	err = cReq.ReceiveResult(res)
	require.NoError(t, err)
}

func serializeReply(msg *message.Message) (*message.Message, error) {
	msg = msg.Copy()

	meta := payload.Meta{
		Payload: msg.Payload,
		ID:      []byte(msg.UUID),
	}

	buf, err := meta.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "serializePayload. failed to wrap message")
	}
	msg.Payload = buf

	return msg, nil
}
