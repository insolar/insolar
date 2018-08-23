/*
 *    Copyright 2018 INS Ecosystem
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

package hostnetwork

import (
	"errors"
	"net"
	"testing"

	"github.com/insolar/insolar/network/hostnetwork/relay"
	"github.com/insolar/insolar/network/hostnetwork/rpc"
	"github.com/insolar/insolar/network/hostnetwork/store"
	"github.com/insolar/insolar/network/hostnetwork/transport"

	"github.com/stretchr/testify/assert"
)

type mockResolverOk struct{}

func (r *mockResolverOk) Resolve(conn net.PacketConn) (string, error) {
	return "127.0.0.1:31337", nil
}

type mockResolverFail struct{}

func (r *mockResolverFail) Resolve(conn net.PacketConn) (string, error) {
	return "", errors.New("mock resolver error")
}

type mockResolverInvalid struct{}

func (r *mockResolverInvalid) Resolve(conn net.PacketConn) (string, error) {
	return "invalid address", nil
}

type mockConnFactoryOk struct{}

func (cf *mockConnFactoryOk) Create(address string) (net.PacketConn, error) {
	return nil, nil
}

type mockConnFactoryFail struct{}

func (cf *mockConnFactoryFail) Create(address string) (net.PacketConn, error) {
	return nil, errors.New("mock conn factory error")
}

type mockTransportFactoryOk struct{}

func (tf *mockTransportFactoryOk) Create(conn net.PacketConn, proxy relay.Proxy) (transport.Transport, error) {
	return newMockTransport(), nil
}

type mockTransportFactoryFail struct{}

func (tf *mockTransportFactoryFail) Create(conn net.PacketConn, proxy relay.Proxy) (transport.Transport, error) {
	return nil, errors.New("mock transport factory error")
}

func TestNewNetworkConfiguration(t *testing.T) {
	cfg := NewNetworkConfiguration(
		&mockResolverOk{},
		&mockConnFactoryOk{},
		&mockTransportFactoryOk{},
		store.NewMemoryStoreFactory(),
		rpc.NewRPCFactory(map[string]rpc.RemoteProcedure{}),
		relay.NewProxy(),
	)

	expectedCfg := &Configuration{
		addressResolver:   &mockResolverOk{},
		connectionFactory: &mockConnFactoryOk{},
		transportFactory:  &mockTransportFactoryOk{},
		storeFactory:      store.NewMemoryStoreFactory(),
		rpcFactory:        rpc.NewRPCFactory(map[string]rpc.RemoteProcedure{}),
		proxy:             relay.NewProxy(),
	}

	assert.Equal(t, expectedCfg, cfg)
}

func TestConfiguration_CreateNetwork(t *testing.T) {
	cfg := NewNetworkConfiguration(
		&mockResolverOk{},
		&mockConnFactoryOk{},
		&mockTransportFactoryOk{},
		store.NewMemoryStoreFactory(),
		rpc.NewRPCFactory(map[string]rpc.RemoteProcedure{}),
		relay.NewProxy(),
	)

	network, err := cfg.CreateNetwork("127.0.0.1:31337", &Options{})

	assert.NotNil(t, network)
	assert.NoError(t, err)
	assert.Equal(t, cfg.network, network)
}

func TestConfiguration_CreateNetwork_AlreadyCreated(t *testing.T) {
	cfg := NewNetworkConfiguration(
		&mockResolverOk{},
		&mockConnFactoryOk{},
		&mockTransportFactoryOk{},
		store.NewMemoryStoreFactory(),
		rpc.NewRPCFactory(map[string]rpc.RemoteProcedure{}),
		relay.NewProxy(),
	)

	dht, err := cfg.CreateNetwork("127.0.0.1:31337", &Options{})

	assert.NotNil(t, dht)
	assert.NoError(t, err)

	_, err = cfg.CreateNetwork("127.0.0.1:31337", &Options{})

	assert.EqualError(t, err, "already created")
}

func TestConfiguration_CreateNetwork_ConnFactoryFail(t *testing.T) {
	cfg := NewNetworkConfiguration(
		&mockResolverOk{},
		&mockConnFactoryFail{},
		&mockTransportFactoryOk{},
		store.NewMemoryStoreFactory(),
		rpc.NewRPCFactory(map[string]rpc.RemoteProcedure{}),
		relay.NewProxy(),
	)

	_, err := cfg.CreateNetwork("127.0.0.1:31337", &Options{})

	assert.EqualError(t, err, "mock conn factory error")
}

func TestConfiguration_CreateNetwork_ResolverFail(t *testing.T) {
	cfg := NewNetworkConfiguration(
		&mockResolverFail{},
		&mockConnFactoryOk{},
		&mockTransportFactoryOk{},
		store.NewMemoryStoreFactory(),
		rpc.NewRPCFactory(map[string]rpc.RemoteProcedure{}),
		relay.NewProxy(),
	)

	_, err := cfg.CreateNetwork("127.0.0.1:31337", &Options{})

	assert.EqualError(t, err, "mock resolver error")
}

func TestConfiguration_CreateNetwork_InvalidAddress(t *testing.T) {
	cfg := NewNetworkConfiguration(
		&mockResolverInvalid{},
		&mockConnFactoryOk{},
		&mockTransportFactoryOk{},
		store.NewMemoryStoreFactory(),
		rpc.NewRPCFactory(map[string]rpc.RemoteProcedure{}),
		relay.NewProxy(),
	)

	_, err := cfg.CreateNetwork("127.0.0.1:31337", &Options{})

	assert.EqualError(t, err, "address invalid address: missing port in address")
}

func TestConfiguration_CreateNetwork_TransportFactoryFail(t *testing.T) {
	cfg := NewNetworkConfiguration(
		&mockResolverOk{},
		&mockConnFactoryOk{},
		&mockTransportFactoryFail{},
		store.NewMemoryStoreFactory(),
		rpc.NewRPCFactory(map[string]rpc.RemoteProcedure{}),
		relay.NewProxy(),
	)

	_, err := cfg.CreateNetwork("127.0.0.1:31337", &Options{})

	assert.EqualError(t, err, "mock transport factory error")
}
