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

package common

import (
	"math/bits"

	"github.com/insolar/insolar/network/consensus/common"
)

type HostProfile interface {
	GetDefaultEndpoint() common.HostAddress
	GetNodePublicKeyStore() common.PublicKeyStore
	IsAcceptableHost(from common.HostIdentityHolder) bool
	// GetHostType()
}

type NodeIntroduction interface {
	GetClaimEvidence() common.SignedEvidenceHolder
	GetShortNodeId() common.ShortNodeID
}

type NodeIntroProfile interface {
	HostProfile
	GetShortNodeId() common.ShortNodeID
	GetIntroduction() NodeIntroduction
}

type NodeProfile interface {
	NodeIntroProfile
	GetIndex() int
	GetDeclaredPower() MemberPower
	GetPower() MemberPower
	GetSignatureVerifier() common.SignatureVerifier
	GetState() MembershipState
	IsJoiner() bool
	HasWorkingPower() bool
}

type LocalNodeProfile interface {
	NodeProfile
	LocalNodeProfile()
}

type UpdatableNodeProfile interface {
	NodeProfile
	SetRank(index int, state MembershipState, declaredPower MemberPower)
	SetSignatureVerifier(verifier common.SignatureVerifier)
	// Update certificate / mandate
}

type MemberPower uint8

func MemberPowerOf(linearValue uint16) MemberPower { // TODO tests are needed
	linearValue += 32
	pwr := uint8(bits.Len16(linearValue))
	if pwr > 6 {
		pwr -= 6
		linearValue >>= pwr
	} else {
		pwr = 0
	}
	return MemberPower((pwr << 5) | uint8(linearValue-32))
}

func (v MemberPower) ToLinearValue() uint16 {
	return uint16(v&0x1F+32)<<(v>>5) - 32
}

type MembershipState int8

const (
	SuspectedOnce MembershipState = iota - 1
	Undefined
	Joining
	Working
	Leaving
)

func (v MembershipState) IsSuspect() bool {
	return v <= SuspectedOnce
}

func (v MembershipState) TimesAsSuspect() int {
	if !v.IsSuspect() {
		return 0
	}
	return 1 + int(SuspectedOnce-v)
}

func (v *MembershipState) IncrementSuspect() (becameSuspect bool) {
	if v.IsSuspect() {
		*v--
		return false
	} else {
		*v = SuspectedOnce
		return true
	}
}

func (v MembershipState) IsUndefined() bool {
	return v == Undefined
}

func (v MembershipState) IsWorking() bool {
	return v == Working
}

func (v MembershipState) IsJoining() bool {
	return v == Joining
}

func (v MembershipState) IsLeaving() bool {
	return v == Leaving
}

func LessForNodeProfile(c NodeProfile, o NodeProfile) bool {
	cP := c.GetPower()
	oP := o.GetPower()
	if cP < oP {
		return true
	} else if cP > oP {
		return false
	}

	cI := c.GetShortNodeId()
	oI := o.GetShortNodeId()
	if cI < oI {
		return true
	} else {
		return false
	}
}

type MembershipProfile struct {
	Index uint16
	Power MemberPower
	Nsh   NodeStateHash
}

func NewMembershipProfile(index uint16, power MemberPower, nsh NodeStateHash) MembershipProfile {
	return MembershipProfile{Index: index, Power: power, Nsh: nsh}
}

func (p MembershipProfile) IsEmpty() bool {
	return p.Nsh == nil
}

func (p MembershipProfile) Equals(o MembershipProfile) bool {
	if p.Index != o.Index || p.Power != o.Power || p.Nsh == nil || o.Nsh == nil {
		return false
	}
	return p.Nsh == o.Nsh || p.Nsh.Equals(o.Nsh)
}
