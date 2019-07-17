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

	"github.com/gojuno/minimock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/testutils"
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
		transcript *Transcript
		am         artifacts.Client
		le         LogicExecutor
		reply      insolar.Reply
		error      bool
	}{
		{
			name: "success, constructor",
			transcript: &Transcript{
				RequestRef: requestRef,
				Request: &record.IncomingRequest{
					CallType:  record.CTSaveAsChild,
					Base:      &baseRef,
					Prototype: &protoRef,
				},
			},
			le:    NewLogicExecutorMock(mc).ExecuteMock.Return(&requestResult{sideEffectType: artifacts.RequestSideEffectActivate}, nil),
			am:    artifacts.NewClientMock(mc).RegisterResultMock.Return(nil),
			reply: &reply.CallConstructor{Object: &requestRef},
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
		transcript *Transcript
		am         artifacts.Client
		le         LogicExecutor
		error      bool
		result     *requestResult
	}{
		{
			name: "success, constructor",
			transcript: &Transcript{
				Request: &record.IncomingRequest{
					CallType: record.CTSaveAsChild,
				},
			},
			le:     NewLogicExecutorMock(mc).ExecuteMock.Return(&requestResult{sideEffectType: artifacts.RequestSideEffectActivate}, nil),
			result: &requestResult{sideEffectType: artifacts.RequestSideEffectActivate},
		},
		{
			name: "success, method",
			transcript: &Transcript{
				Request: &record.IncomingRequest{
					Object: &objRef,
				},
			},
			am:     artifacts.NewClientMock(mc).GetObjectMock.Return(nil, nil),
			le:     NewLogicExecutorMock(mc).ExecuteMock.Return(&requestResult{sideEffectType: artifacts.RequestSideEffectActivate}, nil),
			result: &requestResult{sideEffectType: artifacts.RequestSideEffectActivate},
		},
		{
			name: "method, no object",
			transcript: &Transcript{
				Request: &record.IncomingRequest{
					Object: &objRef,
				},
			},
			am:    artifacts.NewClientMock(mc).GetObjectMock.Return(nil, errors.New("some")),
			error: true,
		},
		{
			name: "method, execution error",
			transcript: &Transcript{
				Request: &record.IncomingRequest{
					Object: &objRef,
				},
			},
			am:    artifacts.NewClientMock(mc).GetObjectMock.Return(nil, nil),
			le:    NewLogicExecutorMock(mc).ExecuteMock.Return(nil, errors.New("some")),
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
		result     *requestResult
		transcript *Transcript
		am         artifacts.Client
		error      bool
		reply      insolar.Reply
	}{
		{
			name: "activation",
			transcript: &Transcript{
				RequestRef: requestRef,
				Request: &record.IncomingRequest{
					Base:      &baseRef,
					Prototype: &protoRef,
				},
			},
			result: &requestResult{sideEffectType: artifacts.RequestSideEffectActivate},
			am:     artifacts.NewClientMock(mc).RegisterResultMock.Return(nil),
			reply:  &reply.CallConstructor{Object: &requestRef},
		},
		{
			name: "activation error",
			transcript: &Transcript{
				RequestRef: requestRef,
				Request: &record.IncomingRequest{
					Base:      &baseRef,
					Prototype: &protoRef,
				},
			},
			result: &requestResult{sideEffectType: artifacts.RequestSideEffectActivate},
			am:     artifacts.NewClientMock(mc).RegisterResultMock.Return(errors.New("some error")),
			error:  true,
		},
		{
			name: "deactivation",
			transcript: &Transcript{
				RequestRef: requestRef,
				Request:    &record.IncomingRequest{},
			},
			result: &requestResult{sideEffectType: artifacts.RequestSideEffectDeactivate, result: []byte{1, 2, 3}},
			am:     artifacts.NewClientMock(mc).RegisterResultMock.Return(nil),
			reply:  &reply.CallMethod{Result: []byte{1, 2, 3}},
		},
		{
			name: "deactivation error",
			transcript: &Transcript{
				RequestRef: requestRef,
				Request:    &record.IncomingRequest{},
			},
			result: &requestResult{sideEffectType: artifacts.RequestSideEffectDeactivate, result: []byte{1, 2, 3}},
			am:     artifacts.NewClientMock(mc).RegisterResultMock.Return(errors.New("some")),
			error:  true,
		},
		{
			name: "update",
			transcript: &Transcript{
				RequestRef: requestRef,
				Request:    &record.IncomingRequest{},
			},
			result: &requestResult{sideEffectType: artifacts.RequestSideEffectAmend, memory: []byte{3, 2, 1}, result: []byte{1, 2, 3}},
			am:     artifacts.NewClientMock(mc).RegisterResultMock.Return(nil),
			reply:  &reply.CallMethod{Result: []byte{1, 2, 3}},
		},
		{
			name: "update error",
			transcript: &Transcript{
				RequestRef: requestRef,
				Request:    &record.IncomingRequest{},
			},
			result: &requestResult{sideEffectType: artifacts.RequestSideEffectAmend, memory: []byte{3, 2, 1}, result: []byte{1, 2, 3}},
			am:     artifacts.NewClientMock(mc).RegisterResultMock.Return(errors.New("some")),
			error:  true,
		},
		{
			name: "result without update",
			transcript: &Transcript{
				RequestRef: requestRef,
				Request:    &record.IncomingRequest{Object: &objRef},
			},
			result: &requestResult{sideEffectType: artifacts.RequestSideEffectNone, result: []byte{1, 2, 3}},
			am:     artifacts.NewClientMock(mc).RegisterResultMock.Return(nil),
			reply:  &reply.CallMethod{Result: []byte{1, 2, 3}},
		},
		{
			name: "result without update, error",
			transcript: &Transcript{
				RequestRef: requestRef,
				Request:    &record.IncomingRequest{Object: &objRef},
			},
			result: &requestResult{sideEffectType: artifacts.RequestSideEffectNone, result: []byte{1, 2, 3}},
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
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	defer mc.Finish()
	defer mc.Wait(time.Minute)

	requestRef := gen.Reference()
	nodeRef := gen.Reference()

	reqRef := testutils.RandomRef()

	table := []struct {
		name       string
		reply      insolar.Reply
		err        error
		transcript *Transcript
		mb         insolar.MessageBus
	}{
		{
			name: "success",
			transcript: &Transcript{
				RequesterNode: &nodeRef,
				RequestRef:    reqRef,
				Request:       &record.IncomingRequest{},
			},
			reply: &reply.CallConstructor{Object: &requestRef},
			mb: testutils.NewMessageBusMock(mc).SendMock.Set(
				func(
					ctx context.Context, msg insolar.Message, opt *insolar.MessageSendOptions,
				) (insolar.Reply, error) {
					return nil, nil
				},
			),
		},
		{
			name: "error",
			transcript: &Transcript{
				RequesterNode: &nodeRef,
				RequestRef:    reqRef,
				Request:       &record.IncomingRequest{},
			},
			reply: &reply.CallConstructor{Object: &requestRef},
			mb: testutils.NewMessageBusMock(mc).SendMock.Set(
				func(
					ctx context.Context, msg insolar.Message, opt *insolar.MessageSendOptions,
				) (insolar.Reply, error) {
					return nil, errors.New("some error")
				},
			),
		},
	}

	for _, test := range table {
		test := test
		t.Run(test.name, func(t *testing.T) {
			re := &requestsExecutor{MessageBus: test.mb}
			re.SendReply(ctx, test.transcript, test.reply, test.err)
		})
	}
}
