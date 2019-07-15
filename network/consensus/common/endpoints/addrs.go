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

package endpoints

import (
	"fmt"

	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensusv1/packets"
)

type Name string

func (addr *Name) IsLocalHost() bool {
	return addr != nil && len(*addr) == 0
}

func (addr *Name) Equals(o Name) bool {
	return addr != nil && *addr == o
}

func (addr *Name) EqualsToString(o string) bool {
	return addr.Equals(Name(o))
}

func (addr Name) String() string {
	return string(addr)
}

//go:generate minimock -i github.com/insolar/insolar/network/consensus/common/endpoints.Outbound -o . -s _mock.go

type Outbound interface {
	GetEndpointType() NodeEndpointType
	GetRelayID() insolar.ShortNodeID
	GetNameAddress() Name
	GetIPAddress() packets.NodeAddress
	AsByteString() string
	CanAccept(connection Inbound) bool
}

func EqualEndpoints(p, o Outbound) bool {
	if p == nil || o == nil {
		return false
	}
	if p == o {
		return true
	}

	if p.GetEndpointType() != o.GetEndpointType() {
		return false
	}
	switch p.GetEndpointType() {
	case NameEndpoint:
		return p.GetNameAddress() == o.GetNameAddress()
	case IPEndpoint:
		return p.GetIPAddress() == o.GetIPAddress()
	case RelayEndpoint:
		return p.GetRelayID() == o.GetRelayID()
	}
	panic("missing")
}

type NodeEndpointType uint8

const (
	IPEndpoint NodeEndpointType = iota
	NameEndpoint
	RelayEndpoint
)

//go:generate minimock -i github.com/insolar/insolar/network/consensus/common/endpoints.Inbound -o . -s _mock.go

type Inbound interface {
	GetNameAddress() Name
	//	GetIPAddress() packets.NodeAddress // TODO
	GetTransportKey() cryptkit.SignatureKeyHolder
	GetTransportCert() cryptkit.CertificateHolder
	AsByteString() string
}

var _ Inbound = &InboundConnection{}

func NewHostIdentityFromHolder(h Inbound) InboundConnection {
	return InboundConnection{
		Addr: h.GetNameAddress(),
		Key:  h.GetTransportKey(),
		Cert: h.GetTransportCert(),
	}
}

type InboundConnection struct {
	Addr Name
	Key  cryptkit.SignatureKeyHolder
	Cert cryptkit.CertificateHolder
}

func ShortNodeIDAsByteString(nodeID insolar.ShortNodeID) string {
	return fmt.Sprintf("node:%s",
		string([]byte{byte(nodeID), byte(nodeID >> 8), byte(nodeID >> 16), byte(nodeID >> 24)}))
}

func (v *InboundConnection) AsByteString() string {
	return fmt.Sprintf("name:%s", v.Addr)
}

func (v *InboundConnection) GetNameAddress() Name {
	return v.Addr
}

func (v *InboundConnection) GetTransportKey() cryptkit.SignatureKeyHolder {
	return v.Key
}

func (v *InboundConnection) GetTransportCert() cryptkit.CertificateHolder {
	return v.Cert
}
