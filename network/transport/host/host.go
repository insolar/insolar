//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package host

import (
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/pkg/errors"
)

// Host is the over-the-wire representation of a host.
type Host struct {
	// NodeID is unique identifier of the node
	NodeID insolar.RecordRef
	// ShortID is shortened unique identifier of the node inside the globe
	ShortID insolar.ShortNodeID
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
func NewHostN(address string, nodeID insolar.RecordRef) (*Host, error) {
	h, err := NewHost(address)
	if err != nil {
		return nil, err
	}
	h.NodeID = nodeID
	return h, nil
}

// NewHostNS creates a new Host with specified physical address, NodeID and ShortID.
func NewHostNS(address string, nodeID insolar.RecordRef, shortID insolar.ShortNodeID) (*Host, error) {
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
