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
	"github.com/insolar/insolar/insolar"
)

type SortingRank struct {
	nodeID    insolar.ShortNodeID
	powerRole uint16
}

func NewSortingRank(nodeID insolar.ShortNodeID, role PrimaryRole, pw Power, mode OpMode) SortingRank {
	return SortingRank{nodeID, SortingPowerRole(role, pw, mode)}
}

func (v SortingRank) GetNodeID() insolar.ShortNodeID {
	return v.nodeID
}

func (v SortingRank) IsWorking() bool {
	return v.powerRole != 0
}

func (v SortingRank) GetWorkingRole() PrimaryRole {
	return PrimaryRole(v.powerRole >> 8)
}

func (v SortingRank) GetPower() Power {
	return Power(v.powerRole)
}

// NB! Sorting is REVERSED
func (v SortingRank) Less(o SortingRank) bool {
	if o.powerRole < v.powerRole {
		return true
	}
	if o.powerRole > v.powerRole {
		return false
	}
	return o.nodeID < v.nodeID
}

// NB! Sorting is REVERSED
func LessByID(vNodeID, oNodeID insolar.ShortNodeID) bool {
	return oNodeID < vNodeID
}

func SortingPowerRole(role PrimaryRole, pw Power, mode OpMode) uint16 {
	if role == 0 {
		panic("illegal value")
	}
	if pw == 0 || mode.IsPowerless() {
		return 0
	}
	return uint16(role)<<8 | uint16(pw)
}
