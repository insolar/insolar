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

package serialization

import "github.com/insolar/insolar/network/consensus/common"

/*
	ByteSize<=16
		4+1+1+2+4 + (optional)4 = 12 + (optional)4
	NB! For embedded payloads or digests it MUST be treated as Bits128, padded with zeros
*/
type UnifiedProtocolPacketHeader struct {
	/*
		Functions of TargetID, SourceId and RelayId depends on ProtocolType
	*/
	ReceiverID            uint32 //NB! MUST not be included into packet Signature, MUST NOT =0
	ProtocolAndPacketType uint8  `insolar-transport:"[0:3]=header:Packet;[4:7]=header:Protocol"` //[00-03]PacketType [04-07]ProtocolType
	PacketFlags           uint8  `insolar-transport:"[0:0]=flags:relay;[1:]=flags:PacketFlags"`
	// bit[0] RelayFlag =1 when RelayTargetId is !=0 (otherwise that field is excluded)
	//
	HeaderAndPayloadLength uint16 //[00-13] ByteLength of Payload, [14-15] reserved = 0
	SourceId               uint32 //may differ from actual sender when relay is in use, MUST NOT =0
	RelayTargetId          uint32 `insolar-transport:"optional=relay[0]"` //indicates final destination, MUST NOT =0
}
type EmbeddedUnifiedProtocolPacketHeader common.Bits128

type ProtocolPacketExample struct {
	Header UnifiedProtocolPacketHeader
	//Protocol Custom Payload - length is in the header
	PacketSignature common.Bits512 //
}
