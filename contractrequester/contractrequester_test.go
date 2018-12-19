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
	"runtime"
	"testing"

	"github.com/insolar/insolar/core/message"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func mockMessageBus(t *testing.T, result core.Reply) *testutils.MessageBusMock {
	mbMock := testutils.NewMessageBusMock(t)
	mbMock.SendFunc = func(c context.Context, m core.Message, _ core.Pulse, o *core.MessageSendOptions) (r core.Reply, r1 error) {
		return result, nil
	}
	return mbMock
}

func mockMessageBusError(t *testing.T) *testutils.MessageBusMock {
	mbMock := testutils.NewMessageBusMock(t)
	mbMock.SendFunc = func(c context.Context, m core.Message, _ core.Pulse, o *core.MessageSendOptions) (r core.Reply, r1 error) {
		return nil, errors.New("test error message")
	}
	return mbMock
}

func TestNew(t *testing.T) {
	ps := testutils.NewPulseStorageMock(t)
	messageBus := mockMessageBus(t, nil)

	contractRequester, err := New()

	cm := &component.Manager{}
	cm.Inject(ps, messageBus, contractRequester)

	require.NoError(t, err)
	require.Equal(t, messageBus, contractRequester.MessageBus)
	require.Equal(t, ps, contractRequester.PulseStorage)
}

func TestContractRequester_SendRequest(t *testing.T) {
	ctx := inslogger.TestContext(t)
	ref := testutils.RandomRef()

	pm := testutils.NewPulseStorageMock(t)
	pm.CurrentMock.Return(core.GenesisPulse, nil)

	mbm := mockMessageBus(t, &reply.RegisterRequest{})
	cReq, err := New()
	assert.NoError(t, err)
	cReq.MessageBus = mbm
	cReq.PulseStorage = pm

	mbm.MustRegisterMock.Return()
	cReq.Start(ctx)

	go func() {
		for len(cReq.ResultMap) == 0 {
			runtime.Gosched()
		}
		cReq.ResultMutex.Lock()
		for k, v := range cReq.ResultMap {
			v <- &message.ReturnResults{
				Request: k,
				Reply:   &reply.CallMethod{},
			}
		}
		cReq.ResultMutex.Unlock()
	}()
	result, err := cReq.SendRequest(ctx, &ref, "TestMethod", []interface{}{})

	require.NoError(t, err)
	require.Equal(t, &reply.CallMethod{}, result)
}

func TestContractRequester_SendRequest_RouteError(t *testing.T) {
	ctx := inslogger.TestContext(t)
	ref := testutils.RandomRef()

	pm := testutils.NewPulseStorageMock(t)
	pm.CurrentMock.Return(core.GenesisPulse, nil)

	mbm := mockMessageBus(t, &reply.CallMethod{})
	cReq, err := New()
	assert.NoError(t, err)
	cReq.MessageBus = mbm
	cReq.PulseStorage = pm

	mbm.MustRegisterMock.Return()
	cReq.Start(ctx)

	go func() {
		for len(cReq.ResultMap) == 0 {
			runtime.Gosched()

		}
		cReq.ResultMutex.Lock()
		for k, v := range cReq.ResultMap {
			v <- &message.ReturnResults{
				Request: k,
				Reply:   nil,
			}
		}
		cReq.ResultMutex.Unlock()
	}()

	result, err := cReq.SendRequest(ctx, &ref, "TestMethod", []interface{}{})
	require.Nil(t, result)
}
