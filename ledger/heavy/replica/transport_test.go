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

package replica

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/hostnetwork/future"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/testutils/network"
)

func TestTransport_Send(t *testing.T) {
	var (
		ctx               = inslogger.TestContext(t)
		partyAddress      = "127.0.0.1:8080"
		method            = "foo.Bar"
		msg               = []byte{4, 5, 6}
		expected          = []byte{1, 2, 3}
		serviceNetwork, _ = servicenetwork.NewServiceNetwork(configuration.Configuration{}, nil, false)
		trans             = NewTransport(serviceNetwork)
		net               = network.NewHostNetworkMock(t)
	)
	net.RegisterRequestHandlerMock.Return()
	net.SendRequestToHostMock.Return(makeFuture(expected), nil)
	serviceNetwork.HostNetwork = net
	trans.(*internalTransport).Init(ctx)

	actual, err := trans.Send(ctx, partyAddress, method, msg)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestTransport_Register(t *testing.T) {
	var (
		ctx               = inslogger.TestContext(t)
		partyAddress      = "127.0.0.1:8080"
		method            = "foo.Bar"
		expectedMsg       = []byte{1, 2, 3}
		expectedReply     = []byte{4, 5, 6}
		serviceNetwork, _ = servicenetwork.NewServiceNetwork(configuration.Configuration{}, nil, false)
		trans             = NewTransport(serviceNetwork)
		net               = network.NewHostNetworkMock(t)
	)
	net.RegisterRequestHandlerMock.Return()
	net.SendRequestToHostMock.Return(makeFuture(expectedReply), nil)
	serviceNetwork.HostNetwork = net
	trans.(*internalTransport).Init(ctx)

	trans.Register(method, func(data []byte) ([]byte, error) {
		require.Equal(t, expectedMsg, data)
		return nil, nil
	})
	reply, err := trans.Send(ctx, partyAddress, method, expectedMsg)
	require.NoError(t, err)
	require.Equal(t, expectedReply, reply)
}

func TestTransport_Me(t *testing.T) {
	var (
		ctx               = inslogger.TestContext(t)
		expected          = "127.0.0.1:8080"
		serviceNetwork, _ = servicenetwork.NewServiceNetwork(configuration.Configuration{}, nil, false)
		trans             = NewTransport(serviceNetwork)
		net               = network.NewHostNetworkMock(t)
	)
	net.RegisterRequestHandlerMock.Return()
	net.PublicAddressMock.Return(expected)
	serviceNetwork.HostNetwork = net
	trans.(*internalTransport).Init(ctx)

	actual := trans.Me()
	require.Equal(t, expected, actual)
}

func makeFuture(expected []byte) future.Future {
	request := packet.NewPacket(nil, nil, types.Replication, 0)
	response := packet.NewPacket(nil, nil, types.Replication, 0)
	response.SetResponse(&packet.RPCResponse{Result: expected})
	future := future.NewFuture(0, nil, request, func(f future.Future) {})
	future.SetResponse(response)
	return future
}
