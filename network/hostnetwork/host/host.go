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
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/reference"
)

// Host is the over-the-wire representation of a host.
type Host struct {
	// NodeID is unique identifier of the node
	NodeID insolar.Reference
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
func NewHostN(address string, nodeID insolar.Reference) (*Host, error) {
	h, err := NewHost(address)
	if err != nil {
		return nil, err
	}
	h.NodeID = nodeID
	return h, nil
}

// NewHostNS creates a new Host with specified physical address, NodeID and ShortID.
func NewHostNS(address string, nodeID insolar.Reference, shortID insolar.ShortNodeID) (*Host, error) {
	h, err := NewHostN(address, nodeID)
	if err != nil {
		return nil, err
	}
	h.ShortID = shortID
	return h, nil
}

// String representation of Host.
func (host Host) String() string {
	return fmt.Sprintf("id: %d ref: %s addr: %s", host.ShortID, host.NodeID.String(), host.Address.String())
}

// Equal checks if host equals to other host (e.g. hosts' IDs and network addresses match).
func (host Host) Equal(other Host) bool {
	return host.NodeID.Equal(other.NodeID) && (other.Address != nil) && host.Address.Equal(*other.Address)
}

// Host serialization:
// 1. NodeID (64 bytes)
// 2. ShortID (4 bytes)
// 3. Address.UDP.IP byteslice size (8 bit)
// next fields are present if flag #3 is not 0
// 4. Address.UDP.IP (length = value from #3)
// 5. Address.UDP.Port (2 bytes)
// 6. Address.UDP.Zone (till the end of the input data)

func (host *Host) Marshal() ([]byte, error) {
	length := host.Size()
	buffer := bytes.NewBuffer(make([]byte, 0, length))
	if err := binary.Write(buffer, binary.BigEndian, host.NodeID); err != nil {
		return nil, errors.Wrap(err, "failed to marshal protobuf host NodeID")
	}
	if err := binary.Write(buffer, binary.BigEndian, host.ShortID); err != nil {
		return nil, errors.Wrap(err, "failed to marshal protobuf host ShortID")
	}
	var header byte
	if host.Address != nil {
		header = byte(len(host.Address.IP))
	}
	if err := binary.Write(buffer, binary.BigEndian, header); err != nil {
		return nil, errors.Wrap(err, "failed to marshal protobuf host header")
	}
	if header == 0 {
		// Address is not present, marshalling is finished
		return buffer.Bytes(), nil
	}
	if err := binary.Write(buffer, binary.BigEndian, host.Address.IP); err != nil {
		return nil, errors.Wrap(err, "failed to marshal protobuf host IP")
	}
	if err := binary.Write(buffer, binary.BigEndian, uint16(host.Address.Port)); err != nil {
		return nil, errors.Wrap(err, "failed to marshal protobuf host port")
	}
	if err := binary.Write(buffer, binary.BigEndian, []byte(host.Address.Zone)); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal protobuf host zone")
	}
	return buffer.Bytes(), nil
}

func (host *Host) MarshalTo(data []byte) (int, error) {
	tmp, err := host.Marshal()
	if err != nil {
		return 0, err
	}
	copy(data, tmp)
	return len(tmp), nil
}

func (host *Host) Unmarshal(data []byte) error {
	reader := bytes.NewReader(data)

	var nodeIDBinary [reference.GlobalBinarySize]byte
	if err := binary.Read(reader, binary.BigEndian, &nodeIDBinary); err != nil {
		return errors.Wrap(err, "failed to unmarshal protobuf host NodeID")
	}
	host.NodeID = *insolar.NewReferenceFromBytes(nodeIDBinary[:])

	if err := binary.Read(reader, binary.BigEndian, &host.ShortID); err != nil {
		return errors.Wrap(err, "failed to unmarshal protobuf host ShortID")
	}
	var header byte
	if err := binary.Read(reader, binary.BigEndian, &header); err != nil {
		return errors.Wrap(err, "failed to unmarshal protobuf host header")
	}
	if header == 0 {
		// Address is not present, unmarshalling is finished
		return nil
	}
	ip := make([]byte, header)
	if err := binary.Read(reader, binary.BigEndian, ip); err != nil {
		return errors.Wrap(err, "failed to unmarshal protobuf host IP")
	}
	var port uint16
	if err := binary.Read(reader, binary.BigEndian, &port); err != nil {
		return errors.Wrap(err, "failed to unmarshal protobuf host port")
	}
	zone, err := ioutil.ReadAll(reader)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal protobuf host zone")
	}
	host.Address = &Address{UDPAddr: net.UDPAddr{IP: ip, Port: int(port), Zone: string(zone)}}
	return nil
}

func (host *Host) Size() int {
	if host.Address == nil {
		return host.basicSize()
	}
	return host.basicSize() + len(host.Address.IP) + 2 /* UDP.Port size */ + len(host.Address.Zone)
}

func (host *Host) basicSize() int {
	return insolar.RecordRefSize + insolar.ShortNodeIDSize + 1
}
