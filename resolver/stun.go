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

package resolver

import (
	"net"

	"github.com/ccding/go-stun/stun"
)

type stunResolver struct {
	stunAddress string
}

// NewStunResolver returns new STUN network address resolver
func NewStunResolver(stunAddress string) PublicAddressResolver {
	return newStunResolver(stunAddress)
}

func newStunResolver(stunAddress string) *stunResolver {
	return &stunResolver{
		stunAddress: stunAddress,
	}
}

// Resolve returns node's public network address as it seen from Internet
func (sr *stunResolver) Resolve(conn net.PacketConn) (string, error) {
	client := stun.NewClientWithConnection(conn)

	if sr.stunAddress != "" {
		client.SetServerAddr(sr.stunAddress)
	}

	_, host, err := client.Discover()
	if err != nil {
		return "", err
	}

	_, err = client.Keepalive()
	if err != nil {
		return "", err
	}

	return host.TransportAddr(), nil
}
