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
	"github.com/gojuno/minimock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/logicexecutor"
	"github.com/insolar/insolar/logicrunner/requestresult"
)

func TestRequestsExecutor_ExecuteAndSave(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	requestRef := gen.Reference()
	baseRef := gen.Reference()
	protoRef := gen.Reference()

	table := []struct {
		name       string
		transcript *common.Transcript
		am         artifacts.Client
		le         logicexecutor.LogicExecutor
		reply      insolar.Reply
		error      bool
	}{
		{
			name: "success, constructor",
			transcript: &common.Transcript{
				RequestRef: requestRef,
				Request: &record.IncomingRequest{
					CallType:  record.CTSaveAsChild,
					Base:      &baseRef,
					Prototype: &protoRef,
				},
			},
			le: logicexecutor.NewLogicExecutorMock(mc).
				ExecuteMock.
				Return(
					&requestresult.RequestResult{
						SideEffectType:     artifacts.RequestSideEffectActivate,
						RawObjectReference: requestRef,
					},
					nil,
				),
			am:    artifacts.NewClientMock(mc).RegisterResultMock.Return(nil),
			reply: &reply.CallMethod{Object: &requestRef},
		},
	}

	for _, test := range table {
		test := test
		t.Run(test.name, func(t *testing.T) {
			re := &requestsExecutor{ArtifactManager: test.am, LogicExecutor: test.le}
			res, err := re.ExecuteAndSave(ctx, test.transcript)
			if !test.error {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.Equal(t, test.reply, res)
			} else {
				require.Error(t, err)
				require.Nil(t, res)
			}
		})
	}
}

func TestRequestsExecutor_Execute(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	objRef := gen.Reference()

	table := []struct {
		name       string
		transcript *common.Transcript
		am         artifacts.Client
		le         logicexecutor.LogicExecutor
		error      bool
		result     *requestresult.RequestResult
	}{
		{
			name: "success, constructor",
			transcript: &common.Transcript{
				Request: &record.IncomingRequest{
					CallType: record.CTSaveAsChild,
				},
			},
			le:     logicexecutor.NewLogicExecutorMock(mc).ExecuteMock.Return(&requestresult.RequestResult{SideEffectType: artifacts.RequestSideEffectActivate}, nil),
			result: &requestresult.RequestResult{SideEffectType: artifacts.RequestSideEffectActivate},
		},
		{
			name: "success, method",
			transcript: &common.Transcript{
				Request: &record.IncomingRequest{
					Object: &objRef,
				},
			},
			am:     artifacts.NewClientMock(mc).GetObjectMock.Return(nil, nil),
			le:     logicexecutor.NewLogicExecutorMock(mc).ExecuteMock.Return(&requestresult.RequestResult{SideEffectType: artifacts.RequestSideEffectActivate}, nil),
			result: &requestresult.RequestResult{SideEffectType: artifacts.RequestSideEffectActivate},
		},
		{
			name: "method, no object",
			transcript: &common.Transcript{
				Request: &record.IncomingRequest{
					Object: &objRef,
				},
			},
			am:    artifacts.NewClientMock(mc).GetObjectMock.Return(nil, errors.New("some")),
			error: true,
		},
		{
			name: "method, execution error",
			transcript: &common.Transcript{
				Request: &record.IncomingRequest{
					Object: &objRef,
				},
			},
			am:    artifacts.NewClientMock(mc).GetObjectMock.Return(nil, nil),
			le:    logicexecutor.NewLogicExecutorMock(mc).ExecuteMock.Return(nil, errors.New("some")),
			error: true,
		},
	}

	for _, test := range table {
		test := test
		t.Run(test.name, func(t *testing.T) {
			re := &requestsExecutor{ArtifactManager: test.am, LogicExecutor: test.le}
			result, err := re.Execute(ctx, test.transcript)
			if !test.error {
				require.NoError(t, err)
				require.NotNil(t, result)
				require.Equal(t, test.result, result)
			} else {
				require.Error(t, err)
				require.Nil(t, result)
			}
		})
	}
}

func TestRequestsExecutor_Save(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	requestRef := gen.Reference()
	baseRef := gen.Reference()
	protoRef := gen.Reference()
	objRef := gen.Reference()

	table := []struct {
		name       string
		result     *requestresult.RequestResult
		transcript *common.Transcript
		am         artifacts.Client
		error      bool
		reply      insolar.Reply
	}{
		{
			name: "activation",
			transcript: &common.Transcript{
				RequestRef: requestRef,
				Request: &record.IncomingRequest{
					Base:      &baseRef,
					Prototype: &protoRef,
				},
			},
			result: &requestresult.RequestResult{
				SideEffectType:     artifacts.RequestSideEffectActivate,
				RawObjectReference: requestRef,
			},
			am:    artifacts.NewClientMock(mc).RegisterResultMock.Return(nil),
			reply: &reply.CallMethod{Object: &requestRef},
		},
		{
			name: "activation error",
			transcript: &common.Transcript{
				RequestRef: requestRef,
				Request: &record.IncomingRequest{
					Base:      &baseRef,
					Prototype: &protoRef,
				},
			},
			result: &requestresult.RequestResult{SideEffectType: artifacts.RequestSideEffectActivate},
			am:     artifacts.NewClientMock(mc).RegisterResultMock.Return(errors.New("some error")),
			error:  true,
		},
		{
			name: "deactivation",
			transcript: &common.Transcript{
				RequestRef: requestRef,
				Request:    &record.IncomingRequest{},
			},
			result: &requestresult.RequestResult{
				SideEffectType:     artifacts.RequestSideEffectDeactivate,
				RawResult:          []byte{1, 2, 3},
				RawObjectReference: requestRef,
			},
			am: artifacts.NewClientMock(mc).RegisterResultMock.Return(nil),
			reply: &reply.CallMethod{
				Result: []byte{1, 2, 3}, Object: &requestRef,
			},
		},
		{
			name: "deactivation error",
			transcript: &common.Transcript{
				RequestRef: requestRef,
				Request:    &record.IncomingRequest{},
			},
			result: &requestresult.RequestResult{SideEffectType: artifacts.RequestSideEffectDeactivate, RawResult: []byte{1, 2, 3}},
			am:     artifacts.NewClientMock(mc).RegisterResultMock.Return(errors.New("some")),
			error:  true,
		},
		{
			name: "update",
			transcript: &common.Transcript{
				RequestRef: requestRef,
				Request:    &record.IncomingRequest{},
			},
			result: &requestresult.RequestResult{
				SideEffectType:     artifacts.RequestSideEffectAmend,
				Memory:             []byte{3, 2, 1},
				RawResult:          []byte{1, 2, 3},
				RawObjectReference: requestRef,
			},
			am: artifacts.NewClientMock(mc).RegisterResultMock.Return(nil),
			reply: &reply.CallMethod{
				Result: []byte{1, 2, 3},
				Object: &requestRef,
			},
		},
		{
			name: "update error",
			transcript: &common.Transcript{
				RequestRef: requestRef,
				Request:    &record.IncomingRequest{},
			},
			result: &requestresult.RequestResult{SideEffectType: artifacts.RequestSideEffectAmend, Memory: []byte{3, 2, 1}, RawResult: []byte{1, 2, 3}},
			am:     artifacts.NewClientMock(mc).RegisterResultMock.Return(errors.New("some")),
			error:  true,
		},
		{
			name: "result without update",
			transcript: &common.Transcript{
				RequestRef: requestRef,
				Request:    &record.IncomingRequest{Object: &objRef},
			},
			result: &requestresult.RequestResult{
				SideEffectType:     artifacts.RequestSideEffectNone,
				RawResult:          []byte{1, 2, 3},
				RawObjectReference: requestRef,
			},
			am: artifacts.NewClientMock(mc).RegisterResultMock.Return(nil),
			reply: &reply.CallMethod{
				Result: []byte{1, 2, 3},
				Object: &requestRef,
			},
		},
		{
			name: "result without update, error",
			transcript: &common.Transcript{
				RequestRef: requestRef,
				Request:    &record.IncomingRequest{Object: &objRef},
			},
			result: &requestresult.RequestResult{SideEffectType: artifacts.RequestSideEffectNone, RawResult: []byte{1, 2, 3}},
			am:     artifacts.NewClientMock(mc).RegisterResultMock.Return(errors.New("some")),
			error:  true,
		},
	}

	for _, test := range table {
		test := test
		t.Run(test.name, func(t *testing.T) {
			re := &requestsExecutor{ArtifactManager: test.am}
			replyVal, err := re.Save(ctx, test.transcript, test.result)
			if !test.error {
				require.NoError(t, err)
				require.NotNil(t, replyVal)
				require.Equal(t, test.reply, replyVal)
			} else {
				require.Error(t, err)
				require.Nil(t, replyVal)
			}
		})
	}
}

func TestRequestsExecutor_SendReply(t *testing.T) {

	reqRef := gen.Reference()

	replyMessage := func(msg *message.Message) *message.Message {
		replyMsg := payload.MustNewMessage(&payload.Error{Text: "test error", Code: payload.CodeUnknown})
		meta := payload.Meta{
			Payload: msg.Payload,
		}
		buf, _ := meta.Marshal()
		replyMsg.Payload = buf
		return replyMsg
	}

	sendRoleHelper := func(ctx context.Context, msg *message.Message, role insolar.DynamicRole, target insolar.Reference) (<-chan *message.Message, func()) {
		res := make(chan *message.Message)
		go func() { res <- replyMessage(msg) }()
		return res, func() { close(res) }
	}
	sendTargetHelper :=	func(ctx context.Context, msg *message.Message, target insolar.Reference) (<-chan *message.Message, func()) {
		res := make(chan *message.Message)
		go func() { res <- replyMessage(msg) }()
		return res, func() { close(res) }
	}

	table := []struct {
		name    string
		mocks   func(ctx context.Context, mc minimock.Tester) RequestsExecutor
		reply   insolar.Reply
		request record.IncomingRequest
		err     error
	}{
		{
			name: "success, reply to caller",
			mocks: func(ctx context.Context, mc minimock.Tester) RequestsExecutor {
				pa := pulse.NewAccessorMock(t)
				pa.LatestMock.Set(func(p context.Context) (insolar.Pulse, error) {
					return insolar.Pulse{
						PulseNumber: 1000,
					}, nil
				})
				sender := bus.NewSenderMock(t).SendRoleMock.Set(sendRoleHelper)

				return &requestsExecutor{Sender: sender, PulseAccessor: pa}
			},
			request: record.IncomingRequest{
				Caller: gen.Reference(),
			},
			reply: &reply.CallMethod{Object: &reqRef},
		},
		{
			name: "success, reply to API",
			mocks: func(ctx context.Context, mc minimock.Tester) RequestsExecutor {
				pa := pulse.NewAccessorMock(t)
				pa.LatestMock.Set(func(p context.Context) (insolar.Pulse, error) {
					return insolar.Pulse{
						PulseNumber: 1000,
					}, nil
				})
				sender := bus.NewSenderMock(t).SendTargetMock.Set(sendTargetHelper)

				return &requestsExecutor{Sender: sender, PulseAccessor: pa}
			},
			request: record.IncomingRequest{
				APINode: gen.Reference(),
			},
			reply: &reply.CallMethod{Object: &reqRef},
		},
		{
			name: "success, reply with error",
			mocks: func(ctx context.Context, mc minimock.Tester) RequestsExecutor {
				pa := pulse.NewAccessorMock(t)
				pa.LatestMock.Set(func(p context.Context) (insolar.Pulse, error) {
					return insolar.Pulse{
						PulseNumber: 1000,
					}, nil
				})
				sender := bus.NewSenderMock(t).SendRoleMock.Set(sendRoleHelper)

				return &requestsExecutor{Sender: sender, PulseAccessor: pa}
			},
			request: record.IncomingRequest{
				Caller: gen.Reference(),
			},
			err:   errors.New("some"),
		},
		{
			name: "return mode NoWait, no reply required",
			mocks: func(ctx context.Context, mc minimock.Tester) RequestsExecutor {
				return &requestsExecutor{}
			},
			request: record.IncomingRequest{
				ReturnMode: record.ReturnNoWait,
			},
		},
		{
			name: "empty reply and no error",
			mocks: func(ctx context.Context, mc minimock.Tester) RequestsExecutor {
				return &requestsExecutor{}
			},
			request: record.IncomingRequest{
			},
		},
	}

	for _, test := range table {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx := inslogger.TestContext(t)
			mc := minimock.NewController(t)

			re := test.mocks(ctx, mc)
			re.SendReply(ctx, reqRef, test.request, test.reply, test.err)

			mc.Wait(time.Minute)
			mc.Finish()
		})
	}
}
