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

package host

import (
	"fmt"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

// Host is the over-the-wire representation of a host.
type Host struct {
	// NodeID is unique identifier of the node
	NodeID core.RecordRef
	// ShortID is shortened unique identifier of the node inside the globe
	ShortID core.ShortNodeID
	// Address is IP and port.
	Address *Address
}

// NewHost creates a new Host with specified physical address.
func NewHost(address string) (*Host, error) {
	addr, err := NewAddress(address)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create Host")
	}
	return &Host{Address: addr}, nil
}

// NewHostN creates a new Host with specified physical address and NodeID.
func NewHostN(address string, nodeID core.RecordRef) (*Host, error) {
	h, err := NewHost(address)
	if err != nil {
		return nil, err
	}
	h.NodeID = nodeID
	return h, nil
}

// NewHostNS creates a new Host with specified physical address, NodeID and ShortID.
func NewHostNS(address string, nodeID core.RecordRef, shortID core.ShortNodeID) (*Host, error) {
	h, err := NewHostN(address, nodeID)
	if err != nil {
		return nil, err
	}
	h.ShortID = shortID
	return h, nil
}

// String representation of Host.
func (host Host) String() string {
	return fmt.Sprintf("%s (%s)", host.NodeID.String(), host.Address.String())
}

// Equal checks if host equals to other host (e.g. hosts' IDs and network addresses match).
func (host Host) Equal(other Host) bool {
	return host.NodeID.Equal(other.NodeID) && (other.Address != nil) && host.Address.Equal(*other.Address)
}
