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

package network

import (
	"errors"
	"net"

	"github.com/insolar/network/connection"
	"github.com/insolar/network/node"
	"github.com/insolar/network/relay"
	"github.com/insolar/network/resolver"
	"github.com/insolar/network/rpc"
	"github.com/insolar/network/store"
	"github.com/insolar/network/transport"
)

// Configuration is a helper to initialize network easily
type Configuration struct {
	addressResolver resolver.PublicAddressResolver

	connectionFactory connection.Factory
	transportFactory  transport.Factory
	storeFactory      store.Factory
	rpcFactory        rpc.Factory

	network *DHT
	conn    net.PacketConn

	proxy relay.Proxy
}

// NewNetworkConfiguration creates new Configuration
func NewNetworkConfiguration(
	addressResolver resolver.PublicAddressResolver,
	connectionFactory connection.Factory,
	transportFactory transport.Factory,
	storeFactory store.Factory,
	rpcFactory rpc.Factory,
	proxy relay.Proxy,
) *Configuration {
	return &Configuration{
		addressResolver:   addressResolver,
		connectionFactory: connectionFactory,
		transportFactory:  transportFactory,
		storeFactory:      storeFactory,
		rpcFactory:        rpcFactory,
		proxy:             proxy,
	}
}

// CreateNetwork creates and returns DHT network with parameters stored in Configuration
func (cfg *Configuration) CreateNetwork(address string, options *Options) (*DHT, error) {
	var err error

	if cfg.network != nil {
		return nil, errors.New("already created")
	}

	cfg.conn, err = cfg.connectionFactory.Create(address)
	if err != nil {
		return nil, err
	}

	publicAddress, err := cfg.addressResolver.Resolve(cfg.conn)
	if err != nil {
		return nil, err
	}

	originAddress, err := node.NewAddress(publicAddress)
	if err != nil {
		return nil, err
	}

	origin, err := node.NewOrigin(nil, originAddress)
	if err != nil {
		return nil, err
	}

	tp, err := cfg.transportFactory.Create(cfg.conn, cfg.proxy)
	if err != nil {
		return nil, err
	}

	cfg.network, err = NewDHT(
		cfg.storeFactory.Create(),
		origin,
		tp,
		cfg.rpcFactory.Create(),
		options,
		cfg.proxy)
	if err != nil {
		return nil, err
	}

	return cfg.network, nil
}

// CloseNetwork stops networking
func (cfg *Configuration) CloseNetwork() error {
	cfg.network.Disconnect()
	return cfg.conn.Close()
}
