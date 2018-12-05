/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package contractrequester

import (
	"context"
	"errors"
	"testing"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func mockMessageBus(t *testing.T, result core.Reply) *testutils.MessageBusMock {
	mbMock := testutils.NewMessageBusMock(t)
	mbMock.SendFunc = func(c context.Context, m core.Message, o *core.MessageSendOptions) (r core.Reply, r1 error) {
		return result, nil
	}
	return mbMock
}

func mockMessageBusError(t *testing.T) *testutils.MessageBusMock {
	mbMock := testutils.NewMessageBusMock(t)
	mbMock.SendFunc = func(c context.Context, m core.Message, o *core.MessageSendOptions) (r core.Reply, r1 error) {
		return nil, errors.New("test error message")
	}
	return mbMock
}

func TestNew(t *testing.T) {
	messageBus := mockMessageBus(t, nil)

	contractRequester, err := New()

	cm := &component.Manager{}
	cm.Inject(messageBus, contractRequester)

	require.NoError(t, err)
	require.Equal(t, messageBus, contractRequester.MessageBus)
}

func TestContractRequester_routeCall(t *testing.T) {
	ctx := inslogger.TestContext(t)
	testResult := &reply.CallMethod{}

	cReq := &ContractRequester{
		MessageBus: mockMessageBus(t, testResult),
	}

	routResult, err := cReq.routeCall(ctx, testutils.RandomRef(), "TestMethod", core.Arguments{})

	require.NoError(t, err)
	require.Equal(t, testResult, routResult)
}

func TestContractRequester_routeCall_SendError(t *testing.T) {
	ctx := inslogger.TestContext(t)

	cReq := &ContractRequester{
		MessageBus: mockMessageBusError(t),
	}

	routResult, err := cReq.routeCall(ctx, testutils.RandomRef(), "TestMethod", core.Arguments{})

	require.Contains(t, err.Error(), "couldn't send message")
	require.Contains(t, err.Error(), "test error message")
	require.Nil(t, routResult)
}

func TestContractRequester_routeCall_MessageBusNil(t *testing.T) {
	ctx := inslogger.TestContext(t)

	cReq := &ContractRequester{}

	routResult, err := cReq.routeCall(ctx, testutils.RandomRef(), "TestMethod", core.Arguments{})

	require.Contains(t, err.Error(), "message bus was not set during initialization")
	require.Nil(t, routResult)
}

func TestContractRequester_SendRequest(t *testing.T) {
	ctx := inslogger.TestContext(t)
	ref := testutils.RandomRef()
	testResult := &reply.CallMethod{}

	cReq := &ContractRequester{
		MessageBus: mockMessageBus(t, testResult),
	}

	result, err := cReq.SendRequest(ctx, &ref, "TestMethod", []interface{}{})

	require.NoError(t, err)
	require.Equal(t, testResult, result)
}

func TestContractRequester_SendRequest_RouteError(t *testing.T) {
	ctx := inslogger.TestContext(t)
	ref := testutils.RandomRef()

	cReq := &ContractRequester{
		MessageBus: mockMessageBusError(t),
	}

	result, err := cReq.SendRequest(ctx, &ref, "TestMethod", []interface{}{})

	require.Contains(t, err.Error(), "Can't route call")
	require.Contains(t, err.Error(), "test error message")
	require.Nil(t, result)
}
