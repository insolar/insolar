///
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
///

package gcp_types

import (
	"math"
)

type PacketType uint8

const (
	/* Phased Packets - these are SENT by a node in the given sequence */

	PacketPhase0 PacketType = iota
	PacketPhase1            /* Namely Phases0 and Phases1 are actually variations of Phase 1 */

	PacketPhase2 /* Phase2 has no phased variations, so =3 is reserved */
	_

	PacketPhase3 /* Namely Phases3 and Phases4 are actually variations of Phase 3 */
	PacketPhase4

	/* Off-phase Packets - these packets can be sent at any moment */

	PacketPulse /* Triggers Phase0-1 */
	PacketFraud /* Delivers fraud proof, by request only */

	PacketReqPhase1 /* Request to resend own NSH - will be replied with PacketPhase1 without PulseData.
	The reply MUST include all claims presented in the original Phase1 packet.
	This request MUST be replied not more than 1 time per requesting node per consensus round,
	otherwise is ignored.
	*/
	PacketReqIntro /* Request to resend other's (NSH + intro) - will be replied with PacketPhase2.
	Only joiners can send this request, and only to anyone in a relevant neighbourhood.
	Limited by 1 times per requesting node per consensus round per requested intro,
	otherwise is ignored.
	*/
	PacketReqFraud /* Requests fraud proof */

	MaxPacketType
)

func (p PacketType) IsPhasedPacket() bool {
	return p <= PacketPhase4
}

func (p PacketType) IsMemberPacket() bool {
	return p != PacketPulse
}

func (p PacketType) GetPayloadEquivalent() PacketType {
	switch p {
	case PacketReqPhase1:
		return PacketPhase1
	case PacketReqIntro:
		return PacketPhase2
	default:
		return p
	}
}

func (p PacketType) ToPhaseNumber() (PhaseNumber, bool) {
	switch p {
	case PacketPhase0:
		return Phase0, true
	case PacketPhase1:
		return Phase1, true
	case PacketPhase2:
		return Phase2, true
	case PacketPhase3:
		return Phase3, true
	case PacketPhase4:
		return Phase4, true
	default:
		return math.MaxUint8, false
	}
}
