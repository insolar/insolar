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

package member

import (
	"fmt"
)

/*
	Power      common2.Power // serialized to [00-07]
	Index      uint16              // serialized to [08-17]
	TotalCount uint16              // serialized to [18-27]
	OpMode 	   //serialized to [28-31]
*/type Rank uint32

const JoinerRank Rank = 0

func (v Rank) GetPower() Power {
	return Power(v)
}

func (v Rank) GetIndex() Index {
	if v == JoinerRank {
		return JoinerIndex
	}
	return AsIndex(int(v>>8) & NodeIndexMask)
}

func (v Rank) GetTotalCount() uint16 {
	return AsIndex(int(v>>18) & NodeIndexMask).AsUint16()
}

func (v Rank) GetMode() OpMode {
	return OpMode(v >> 28)
}

func (v Rank) IsJoiner() bool {
	return v == JoinerRank
}

func (v Rank) String() string {
	if v.IsJoiner() {
		return "{joiner}"
	}
	return fmt.Sprintf("{%v %d/%d pw:%v}", v.GetMode(), v.GetIndex(), v.GetTotalCount(), v.GetPower())
}

func NewMembershipRank(mode OpMode, pw Power, idx Index, count Index) Rank {
	if idx >= count {
		panic("illegal value")
	}

	r := uint32(pw)
	r |= idx.AsUint32() << 8
	r |= count.AsUint32() << 18
	r |= mode.AsUnit32() << 28
	return Rank(r)
}

type RankCursor struct {
	Role           PrimaryRole
	RoleIndex      Index
	RolePowerIndex uint32
	TotalIndex     Index
}

type InterimRank struct {
	RankCursor
	SpecialRoles SpecialRole
	Power        Power
	OpMode       OpMode
}

type FullRank struct {
	InterimRank
	RoleCount uint16
	RolePower uint32
}

func (v FullRank) AsMembershipRank(totalCount Index) Rank {
	return NewMembershipRank(v.OpMode, v.Power, v.TotalIndex, totalCount)
}
