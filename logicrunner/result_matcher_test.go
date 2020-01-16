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

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock/v3"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
)

func replyMessage(msg *message.Message) *message.Message {
	replyMsg := payload.MustNewMessage(&payload.Error{Text: "test error", Code: payload.CodeUnknown})
	meta := payload.Meta{
		Payload: msg.Payload,
	}
	buf, _ := meta.Marshal()
	replyMsg.Payload = buf
	return replyMsg
}

func sendTargetHelper(
	ctx context.Context, msg *message.Message, target insolar.Reference,
) (<-chan *message.Message, func()) {
	res := make(chan *message.Message)
	go func() { res <- replyMessage(msg) }()
	return res, func() { close(res) }
}

func TestResultsMatcher_AddStillExecution(t *testing.T) {
	defer testutils.LeakTester(t)

	tests := []struct {
		name  string
		mocks func(ctx context.Context, t minimock.Tester) (ResultMatcher, payload.StillExecuting)
	}{
		{
			name: "empty matcher",
			mocks: func(ctx context.Context, t minimock.Tester) (ResultMatcher, payload.StillExecuting) {
				rm := newResultsMatcher(nil, nil)
				msg := payload.StillExecuting{}
				return rm, msg
			},
		},
		{
			name: "match",
			mocks: func(ctx context.Context, t minimock.Tester) (ResultMatcher, payload.StillExecuting) {
				reqRef := gen.Reference()
				nodeRef := gen.Reference()

				pa := pulse.NewAccessorMock(t)
				pa.LatestMock.Set(func(p context.Context) (insolar.Pulse, error) {
					return insolar.Pulse{
						PulseNumber: 1000,
					}, nil
				})

				rm := newResultsMatcher(
					bus.NewSenderMock(t).SendTargetMock.Set(sendTargetHelper), pa,
				)

				rm.AddUnwantedResponse(ctx, payload.ReturnResults{
					Reason: reqRef,
				})

				msg := payload.StillExecuting{
					Executor:    nodeRef,
					RequestRefs: []insolar.Reference{reqRef},
				}
				return rm, msg
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
			mc := minimock.NewController(t)

			rm, msg := test.mocks(ctx, mc)
			rm.AddStillExecution(ctx, msg)

			mc.Wait(1 * time.Second)
			mc.Finish()
		})
	}
}

func TestResultsMatcher_AddUnwantedResponse(t *testing.T) {
	defer testutils.LeakTester(t)

	tests := []struct {
		name  string
		mocks func(ctx context.Context, t minimock.Tester) (ResultMatcher, payload.ReturnResults)
		error bool
	}{
		{
			name: "empty matcher",
			mocks: func(ctx context.Context, t minimock.Tester) (ResultMatcher, payload.ReturnResults) {
				rm := newResultsMatcher(nil, nil)
				msg := payload.ReturnResults{}
				return rm, msg
			},
		},
		{
			name: "match",
			mocks: func(ctx context.Context, t minimock.Tester) (ResultMatcher, payload.ReturnResults) {
				reqRef := gen.Reference()
				nodeRef := gen.Reference()

				pa := pulse.NewAccessorMock(t)
				pa.LatestMock.Set(func(p context.Context) (insolar.Pulse, error) {
					return insolar.Pulse{
						PulseNumber: 1000,
					}, nil
				})

				rm := newResultsMatcher(
					bus.NewSenderMock(t).SendTargetMock.Set(sendTargetHelper), pa,
				)

				rm.AddStillExecution(ctx, payload.StillExecuting{
					Executor:    nodeRef,
					RequestRefs: []insolar.Reference{reqRef},
				})
				msg := payload.ReturnResults{
					Reason: reqRef,
				}
				return rm, msg
			},
		},
		{
			name: "loop detected",
			mocks: func(ctx context.Context, t minimock.Tester) (ResultMatcher, payload.ReturnResults) {
				reqRef := gen.Reference()
				nodeRef := gen.Reference()

				pa := pulse.NewAccessorMock(t)

				rm := newResultsMatcher(
					bus.NewSenderMock(t), pa,
				)

				rm.AddStillExecution(ctx, payload.StillExecuting{
					Executor:    nodeRef,
					RequestRefs: []insolar.Reference{reqRef},
				})
				msg := payload.ReturnResults{
					Reason:      reqRef,
					ResendCount: 2,
				}
				return rm, msg
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
			mc := minimock.NewController(t)

			rm, msg := test.mocks(ctx, mc)
			rm.AddUnwantedResponse(ctx, msg)

			mc.Wait(1 * time.Second)
			mc.Finish()
		})
	}
}

func TestResultsMatcher_Clear(t *testing.T) {
	defer testutils.LeakTester(t)

	tests := []struct {
		name  string
		mocks func(ctx context.Context, t minimock.Tester) ResultMatcher
		error bool
	}{
		{
			name: "success",
			mocks: func(ctx context.Context, t minimock.Tester) ResultMatcher {
				reqRef1 := gen.Reference()
				reqRef2 := gen.Reference()

				rm := newResultsMatcher(nil, nil)

				rm.AddStillExecution(ctx, payload.StillExecuting{
					RequestRefs: []insolar.Reference{reqRef1},
				})

				rm.AddUnwantedResponse(ctx, payload.ReturnResults{
					Reason: reqRef2,
				})

				return rm
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := flow.TestContextWithPulse(inslogger.TestContext(t), gen.PulseNumber())
			mc := minimock.NewController(t)

			rm := test.mocks(ctx, mc)
			rm.Clear(ctx)

			mc.Wait(1 * time.Second)
			mc.Finish()
		})
	}
}
