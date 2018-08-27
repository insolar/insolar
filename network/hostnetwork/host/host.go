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

package host

import (
	"fmt"

	"github.com/insolar/insolar/network/hostnetwork/id"
)

// Host is the over-the-wire representation of a host.
type Host struct {
	// ID is a 20 byte unique identifier.
	ID id.ID

	// Address is IP and port.
	Address *Address
}

// NewHost creates a new Host for bootstrapping.
func NewHost(address *Address) *Host {
	return &Host{
		Address: address,
	}
}

// String representation of Host.
func (host Host) String() string {
	return fmt.Sprintf("%s (%s)", host.ID.HashString(), host.Address.String())
}

// Equal checks if host equals to other host (e.g. hosts' IDs and network addresses match).
func (host Host) Equal(other Host) bool {
	return host.ID.HashEqual(other.ID.GetHash()) && (other.Address != nil) && host.Address.Equal(*other.Address)
}
