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

package transport

import (
	"net"

	"github.com/insolar/insolar/network/hostnetwork/relay"
)

// Factory allows to create new Transport.
type Factory interface {
	Create(conn net.PacketConn, proxy relay.Proxy) (Transport, error)
}

type utpTransportFactory struct{}

// NewUTPTransportFactory creates new Factory of utpTransport.
func NewUTPTransportFactory() Factory {
	return &utpTransportFactory{}
}

// Create creates new UTP Transport.
func (utpTransportFactory *utpTransportFactory) Create(conn net.PacketConn, proxy relay.Proxy) (Transport, error) {
	return NewUTPTransport(conn, proxy)
}

type kcpTransportFactory struct{}

// NewKCPTransportFactory creates new Factory of kcpTransport.
func NewKCPTransportFactory( /* security params, keys, etc */ ) Factory {
	return &kcpTransportFactory{}
}

// Create creates new KCP Transport.
func (kcpTransportFactory *kcpTransportFactory) Create(conn net.PacketConn, proxy relay.Proxy) (Transport, error) {
	return NewKCPTransport(conn, proxy)
}
