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

package census

import (
	"math"

	"github.com/insolar/insolar/network/consensus/common"
	common2 "github.com/insolar/insolar/network/consensus/gcpv2/common"
)

var _ common2.LocalNodeProfile = &NodeProfileSlot{}

const joinerIndex = 0x8000

type NodeProfileSlot struct {
	common2.NodeIntroProfile
	verifier common.SignatureVerifier
	index    uint16
	mode     common2.MemberOpMode
	power    common2.MemberPower
}

func NewNodeProfile(index int, p common2.NodeIntroProfile, verifier common.SignatureVerifier, pw common2.MemberPower) NodeProfileSlot {

	if index < 0 || index > common2.MaxNodeIndex {
		panic("illegal value")
	}
	return NodeProfileSlot{index: uint16(index), NodeIntroProfile: p, verifier: verifier, power: pw}
}

func NewJoinerProfile(p common2.NodeIntroProfile, verifier common.SignatureVerifier, pw common2.MemberPower) NodeProfileSlot {

	return NodeProfileSlot{index: joinerIndex, NodeIntroProfile: p, verifier: verifier, power: pw}
}

func (c *NodeProfileSlot) GetDeclaredPower() common2.MemberPower {
	return c.power
}

func (c *NodeProfileSlot) GetOpMode() common2.MemberOpMode {
	return c.mode
}

func (c *NodeProfileSlot) LocalNodeProfile() {
}

func (c *NodeProfileSlot) GetIndex() int {
	return int(c.index & common2.NodeIndexMask)
}

func (c *NodeProfileSlot) IsJoiner() bool {
	return c.index == joinerIndex
}

func (c *NodeProfileSlot) GetSignatureVerifier() common.SignatureVerifier {
	return c.verifier
}

var _ common2.UpdatableNodeProfile = &updatableSlot{}

type updatableSlot struct {
	NodeProfileSlot
	leaveReason uint32
}

func (c *updatableSlot) SetRank(index int, m common2.MemberOpMode, power common2.MemberPower) {
	c.SetIndex(index)
	c.power = power
	c.mode = m
}

func (c *updatableSlot) SetPower(power common2.MemberPower) {
	c.power = power
}

func (c *updatableSlot) SetOpMode(m common2.MemberOpMode) {
	c.mode = m
}

func (c *updatableSlot) SetOpModeAndLeaveReason(leaveReason uint32) {
	c.mode = common2.MemberModeEvictedGracefully
	c.leaveReason = leaveReason
}

func (c *updatableSlot) GetLeaveReason() uint32 {
	if c.mode != common2.MemberModeEvictedGracefully {
		return 0
	}
	return c.leaveReason
}

func (c *updatableSlot) SetIndex(index int) {
	if index < 0 || index > math.MaxUint16 {
		panic("wrong index")
	}
	c.index = uint16(index)
}

func (c *updatableSlot) SetSignatureVerifier(verifier common.SignatureVerifier) {
	c.verifier = verifier
}
