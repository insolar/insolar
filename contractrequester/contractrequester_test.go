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
	"context"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/gojuno/minimock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	busMeta "github.com/insolar/insolar/insolar/bus/meta"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
)

func TestNew(t *testing.T) {
	sender := bus.NewSenderMock(t)
	pulseAccessor := pulse.NewAccessorMock(t)
	jetCoordinator := jet.NewCoordinatorMock(t)
	pcs := platformpolicy.NewPlatformCryptographyScheme()

	ctx := inslogger.TestContext(t)
	contractRequester, err := New(ctx, gochannel.NewGoChannel(gochannel.Config{}, nil), &bus.Bus{})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, contractRequester.Stop())
	}()

	cm := &component.Manager{}
	cm.Inject(sender, contractRequester, pulseAccessor, jetCoordinator, pcs)

	require.NoError(t, err)
	require.Equal(t, sender, contractRequester.Sender)
}

func mockPulseAccessor(t minimock.Tester) pulse.Accessor {
	pulseAccessor := pulse.NewAccessorMock(t)
	currentPulse := insolar.FirstPulseNumber
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

func TestContractRequester_Start(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	ctx := inslogger.TestContext(t)
	cReq, err := New(ctx, gochannel.NewGoChannel(gochannel.Config{}, nil), &bus.Bus{})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, cReq.Stop())
	}()

	cReq.Sender = bus.NewSenderMock(t)
}

func TestContractRequester_SendRequest(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()

	ref := gen.Reference()

	cReq, err := New(ctx, gochannel.NewGoChannel(gochannel.Config{}, nil), &bus.Bus{})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, cReq.Stop())
	}()

	cReq.PulseAccessor = mockPulseAccessor(mc)
	cReq.JetCoordinator = mockJetCoordinator(mc)
	cReq.PlatformCryptographyScheme = testutils.NewPlatformCryptographyScheme()

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
					requestRef := insolar.NewReference(*insolar.NewID(insolar.FirstPulseNumber, hash[:]))

					resultSender := func() {
						res := test.resultMessage
						res.RequestRef = *requestRef
						cReq.result(ctx, &res)
					}

					resChan := make(chan *message.Message)

					res, err := serializeReply(bus.ReplyAsMessage(ctx, &reply.RegisterRequest{Request: *requestRef}))
					require.NoError(t, err)

					if test.earlyResult {
						resultSender()
					} else {
						go func() {
							resChan <- res
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

	cr, err := New(ctx, gochannel.NewGoChannel(gochannel.Config{}, nil), &bus.Bus{})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, cr.Stop())
	}()

	cr.callTimeout = 1 * time.Nanosecond

	cr.PlatformCryptographyScheme = testutils.NewPlatformCryptographyScheme()

	cr.PulseAccessor = mockPulseAccessor(mc)
	cr.JetCoordinator = jet.NewCoordinatorMock(t)

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

	cr.Sender = bus.NewSenderMock(t).SendRoleMock.Set(
		func(ctx context.Context, msg *message.Message, role insolar.DynamicRole, obj insolar.Reference) (<-chan *message.Message, func()) {
			resChan := make(chan *message.Message)

			data, err := payload.Unmarshal(msg.Payload)
			require.NoError(t, err)

			request := data.(*payload.CallMethod).Request

			hash, err := cr.calcRequestHash(*request)
			require.NoError(t, err)
			requestRef := insolar.NewReference(*insolar.NewID(insolar.FirstPulseNumber, hash[:]))

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

	_, _, err = cr.Call(ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "canceled")
	require.Contains(t, err.Error(), "timeout")
}

func TestReceiveResult(t *testing.T) {
	ctx := context.Background()
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*10)
	defer cancelFunc()

	cr, err := New(ctx, gochannel.NewGoChannel(gochannel.Config{}, nil), &bus.Bus{})
	require.NoError(t, err)
	defer func() {
		require.NoError(t, cr.Stop())
	}()

	mc := minimock.NewController(t)
	defer mc.Finish()

	reqRef := gen.Reference()
	var reqHash [insolar.RecordHashSize]byte
	copy(reqHash[:], reqRef.Record().Hash())

	msg := &payload.ReturnResults{
		RequestRef: reqRef,
	}

	// unexpected result
	res, err := serializeReply(payload.MustNewMessage(msg))
	require.NoError(t, err)
	err = cr.ReceiveResult(res)
	require.NoError(t, err)

	// expected result
	resChan := make(chan *payload.ReturnResults)
	chanResult := make(chan *payload.ReturnResults)
	cr.ResultMap[reqHash] = resChan

	go func() {
		chanResult <- <-cr.ResultMap[reqHash]
	}()

	res, err = serializeReply(payload.MustNewMessage(msg))
	require.NoError(t, err)
	err = cr.ReceiveResult(res)

	require.NoError(t, err)
	require.Equal(t, 0, len(cr.ResultMap))
	require.Equal(t, msg, <-chanResult)
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

	msg.Metadata.Set(busMeta.Type, busMeta.TypeReply)

	return msg, nil
}
