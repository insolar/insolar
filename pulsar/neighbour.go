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
	"crypto"
	"net"
	"net/rpc"
	"sync"

	"github.com/insolar/insolar/configuration"
)

// RPCClientWrapperFactory describes interface for the wrappers factory
type RPCClientWrapperFactory interface {
	CreateWrapper() RPCClientWrapper
}

// RPCClientWrapperFactoryImpl is a base impl of the RPCClientWrapperFactory
type RPCClientWrapperFactoryImpl struct {
}

// CreateWrapper return new RPCClientWrapper
func (RPCClientWrapperFactoryImpl) CreateWrapper() RPCClientWrapper {
	return &RPCClientWrapperImpl{Mutex: &sync.Mutex{}}
}

// RPCClientWrapper describes interface of the wrapper around rpc-client
type RPCClientWrapper interface {
	// Lock takes current neighbour's lock
	Lock()
	// Unlock releases current neighbour's lock
	Unlock()

	// IsInitialised compares underhood rpc-client with nil
	IsInitialised() bool
	// CreateConnection creates connection to an another pulsar
	CreateConnection(connectionType configuration.ConnectionType, connectionAddress string) error
	// Close closes connection
	Close() error

	// Go makes rpc-call to an another pulsar
	Go(serviceMethod string, args interface{}, reply interface{}, done chan *rpc.Call) *rpc.Call

	// ResetClient clears rpc-client
	ResetClient()
}

// RPCClientWrapperImpl is a standard impl of RPCClientWrapper
type RPCClientWrapperImpl struct {
	*sync.Mutex
	*rpc.Client
}

// IsInitialised compares underhood rpc-client with nil
func (impl *RPCClientWrapperImpl) IsInitialised() bool {
	return impl.Client != nil
}

// Close closes connection
func (impl *RPCClientWrapperImpl) Close() error {
	return impl.Client.Close()
}

// Lock takes current neighbour's lock
func (impl *RPCClientWrapperImpl) Lock() {
	impl.Mutex.Lock()
}

// Unlock releases current neighbour's lock
func (impl *RPCClientWrapperImpl) Unlock() {
	impl.Mutex.Unlock()
}

// CreateConnection creates connection to an another pulsar
func (impl *RPCClientWrapperImpl) CreateConnection(connectionType configuration.ConnectionType, connectionAddress string) error {
	conn, err := net.Dial(connectionType.String(), connectionAddress)
	if err != nil {
		return err
	}
	impl.Client = rpc.NewClient(conn)
	return nil
}

// Go makes rpc-call to an another pulsar
func (impl *RPCClientWrapperImpl) Go(serviceMethod string, args interface{}, reply interface{}, done chan *rpc.Call) *rpc.Call {
	return impl.Client.Go(serviceMethod, args, reply, done)
}

// ResetClient clears rpc-client
func (impl *RPCClientWrapperImpl) ResetClient() {
	impl.Lock()
	impl.Client = nil
	impl.Unlock()
}

// Neighbour is a helper struct, which contains info about pulsar-neighbour
type Neighbour struct {
	ConnectionType    configuration.ConnectionType
	ConnectionAddress string
	OutgoingClient    RPCClientWrapper
	PublicKey         crypto.PublicKey
}
