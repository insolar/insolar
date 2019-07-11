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

package api

import (
	"fmt"
	"math/bits"
)

type PhaseNumber uint8

const (
	Phase0 PhaseNumber = iota
	Phase1
	Phase2
	Phase3
	Phase4
	upperPhaseNumber
)

type NodeState uint16

const (
	NodeStateInactive NodeState = 0

	NodeStateLocalActive NodeState = 0x8000 /* Not applicable to remote node */
	// NodeStateFraud			NodeState = 0x4000 /* Not applicable to local node */

	NodeStateReceivedPhases NodeState = 1<<upperPhaseNumber - 1
	NodeStateSentPhases               = NodeStateReceivedPhases << shiftNodeStateSent
)
const shiftNodeStateSent = uint(upperPhaseNumber)

func (ns NodeState) setNodeStateMask(m NodeState) (rt NodeState, modified bool) {
	if ns&m != 0 {
		return ns, false
	}
	return ns | m, true
}

func (ns NodeState) SetReceivedPhase(pn PhaseNumber) (rt NodeState, modified bool) {
	return ns.setNodeStateMask(1 << pn)
}

func (ns *NodeState) UpdReceivedPhase(pn PhaseNumber) (modified bool) {
	*ns, modified = ns.SetReceivedPhase(pn)
	return
}

func (ns *NodeState) UpdReceivedPacket(packetType PacketType) bool {
	pn, ok := packetType.ToPhaseNumber()
	if ok {
		return ns.UpdReceivedPhase(pn)
	}
	// panic("packet type cant be mapped to a phase")
	return false
}

func (ns NodeState) HasReceivedPhase(pn PhaseNumber) bool {
	return ns&(1<<pn) != 0
}

func (ns NodeState) MaxReceivedPhase() (pn PhaseNumber, ok bool) {
	i := bits.Len8(uint8(ns & NodeStateReceivedPhases))
	if i == 0 {
		return 0, false
	}
	return PhaseNumber(i - 1), true
}

func (ns NodeState) SetSentPhase(pn PhaseNumber) (rt NodeState, modified bool) {
	return ns.setNodeStateMask(1 << (uint(pn) + shiftNodeStateSent))
}

func (ns *NodeState) UpdSentPhase(pn PhaseNumber) (modified bool) {
	*ns, modified = ns.SetSentPhase(pn)
	return
}

func (ns NodeState) UpdSentPacket(packetType PacketType) bool {
	pn, ok := packetType.ToPhaseNumber()
	if ok {
		return ns.UpdSentPhase(pn)
	}
	// panic("packet type cant be mapped to a phase")
	return false
}

func (ns NodeState) HasSentPhase(pn PhaseNumber) bool {
	return ns&(1<<(uint(pn)+shiftNodeStateSent)) != 0
}

func (ns NodeState) MaxSentPhase() (pn PhaseNumber, ok bool) {
	i := bits.Len8(uint8((ns & NodeStateSentPhases) >> shiftNodeStateSent))
	if i == 0 {
		return 0, false
	}
	return PhaseNumber(i - 1), true
}

func (ns NodeState) HasReceived() bool {
	return ns&NodeStateReceivedPhases != 0
}

func (ns NodeState) HasSent() bool {
	return ns&NodeStateSentPhases != 0
}

func (ns NodeState) HasReceivedOrSent() bool {
	return ns&(NodeStateSentPhases|NodeStateReceivedPhases) != 0
}

func (ns NodeState) SetLocalActive() (rt NodeState, duplicate bool) {
	return ns.setNodeStateMask(NodeStateLocalActive)
}

func (ns NodeState) IsLocalActive() bool {
	return ns&NodeStateLocalActive != 0
}

func (ns NodeState) IsOperational() bool {
	// if ns & NodeStateFraud != 0 { return false }
	return ns&(NodeStateLocalActive|NodeStateReceivedPhases) != 0
}

func (ns NodeState) String() string {
	switch ns {
	case NodeStateInactive:
		return "inactive"
	case NodeStateLocalActive:
		return "active"
	default:
		mode := "rmt"
		if ns.IsLocalActive() {
			mode = "act"
		}
		// fraud := ",fraud"
		// if !ns.HasFraud() {
		// 	fraud = ""
		// }
		return fmt.Sprintf("%s%s%s", mode, // fraud,
			fmtNodeStatePhases('R', false, ns&NodeStateReceivedPhases),
			fmtNodeStatePhases('S', false, (ns&NodeStateSentPhases)>>shiftNodeStateSent),
		)
	}
}

func fmtNodeStatePhases(p byte, suffix bool, ns NodeState) string {
	if ns == 0 {
		return ""
	}
	buf := [upperPhaseNumber + 2]byte{}
	var o = 0
	buf[0] = ','
	o++
	if suffix {
		buf[len(buf)-1] = p
	} else {
		buf[o] = p
		o++
	}
	var i byte = '0'
	for ; o < len(buf); o++ {
		if ns&NodeState(1) == 0 {
			buf[o] = '_'
		} else {
			buf[o] = i
		}
		i++
		ns >>= 1
	}
	return string(buf[:])
}
