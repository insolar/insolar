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

package pulsar

import (
	"crypto/ecdsa"
	"net"
	"net/rpc"
	"sync"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/log"
)

// Interface for factory of rpc wrappers
// Needed for creation wrappers objects
type RpcClientWrapperFactory interface {
	CreateWrapper() RpcClientWrapper
}

// Standard factory implementation
// Returns RpcClientWrapperImpl
type RpcClientWrapperFactoryImpl struct {
}

// Standard factory implementation
// Returns RpcClientWrapperImpl
func (RpcClientWrapperFactoryImpl) CreateWrapper() RpcClientWrapper {
	return &RpcClientWrapperImpl{}
}

// Interface for wrappers around rpc clients
type RpcClientWrapper interface {
	// Take current neighbour's lock
	Lock()
	// Release current neighbour's lock
	Unlock()

	// Check if client initialised
	IsInitialised() bool
	// Set wrapper's undercover rpc client
	SetRpcClient(client *rpc.Client)
	// Create connection and reinit client
	CreateConnection(connectionType configuration.ConnectionType, connectionAddress string) error
	// Close wrapped client
	Close() error

	// Make rpc call in async-style
	Go(serviceMethod string, args interface{}, reply interface{}, done chan *rpc.Call) *rpc.Call
}

// Base RpcClientWrapper impl for rpc.Client
type RpcClientWrapperImpl struct {
	*sync.Mutex
	*rpc.Client
}

func (impl *RpcClientWrapperImpl) IsInitialised() bool {
	return impl.Client != nil
}

// Close wrapped client
func (impl *RpcClientWrapperImpl) Close() error {
	return impl.Client.Close()
}

// Take current neighbour's lock
func (impl *RpcClientWrapperImpl) Lock() {
	impl.Lock()
}

// Release current neighbour's lock
func (impl *RpcClientWrapperImpl) Unlock() {
	impl.Unlock()
}

// Set wrapper's undercover rpc client
func (impl *RpcClientWrapperImpl) SetRpcClient(client *rpc.Client) {
	impl.Client = client
}

// Create base rpc connection
func (impl *RpcClientWrapperImpl) CreateConnection(connectionType configuration.ConnectionType, connectionAddress string) error {
	conn, err := net.Dial(connectionType.String(), connectionAddress)
	if err != nil {
		return err
	}
	impl.Client = rpc.NewClient(conn)
	return nil
}

// Make a call in async style, with done channel as async-marker
func (impl *RpcClientWrapperImpl) Go(serviceMethod string, args interface{}, reply interface{}, done chan *rpc.Call) *rpc.Call {
	return impl.Client.Go(serviceMethod, args, reply, done)
}

// Helper for functionality of connection to another pulsar
type Neighbour struct {
	ConnectionType    configuration.ConnectionType
	ConnectionAddress string
	OutgoingClient    RpcClientWrapper
	PublicKey         *ecdsa.PublicKey
}

// Check connection error, write it to the log and try to refresh connection
func (neighbour *Neighbour) CheckAndRefreshConnection(rpcErr error) error {
	log.Infof("Restarting RPC Connection to %v due to error %v", neighbour.ConnectionAddress, rpcErr)

	neighbour.OutgoingClient.Lock()

	err := neighbour.OutgoingClient.CreateConnection(neighbour.ConnectionType, neighbour.ConnectionAddress)
	if err != nil {
		log.Errorf("Refreshing connection to %v failed due to error %v", neighbour.ConnectionAddress, err)
	}

	neighbour.OutgoingClient.Unlock()

	return err
}
