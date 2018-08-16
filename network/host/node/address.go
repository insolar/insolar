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

package node

import (
	"net"
)

// Address is node's real network address.
type Address struct {
	net.UDPAddr
}

// NewAddress is constructor.
func NewAddress(address string) (*Address, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	return &Address{UDPAddr: *udpAddr}, nil
}

// Equal checks if address is equal to another.
func (address Address) Equal(other Address) bool {
	return address.IP.Equal(other.IP) && address.Port == other.Port
}
