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

type OpMode uint8 // 4-bit value
const (
	ModeBits                          = 4
	ModeFlagRestrictedBehavior OpMode = 4
	ModeFlagValidationWarning  OpMode = 2
	ModeFlagSuspendedOps       OpMode = 1

	ModeNormal                    OpMode = 0
	ModeSuspected                        = /* 0x01 */ ModeFlagSuspendedOps
	ModePossibleFraud                    = /* 0x02 */ ModeFlagValidationWarning
	ModePossibleFraudAndSuspected        = /* 0x03 */ ModeFlagSuspendedOps | ModeFlagValidationWarning
	ModeRestrictedAnnouncement           = /* 0x04 */ ModeFlagRestrictedBehavior
	ModeEvictedGracefully                = /* 0x05 */ ModeFlagRestrictedBehavior | ModeFlagSuspendedOps
	ModeEvictedAsFraud                   = /* 0x06 */ ModeFlagRestrictedBehavior | ModeFlagValidationWarning
	ModeEvictedAsSuspected               = /* 0x07 */ ModeFlagRestrictedBehavior | ModeFlagValidationWarning | ModeFlagSuspendedOps
)

func (v OpMode) IsEvicted() bool {
	return v >= ModeEvictedGracefully
}

func (v OpMode) IsJustJoined() bool {
	return v == ModeRestrictedAnnouncement
}

func (v OpMode) IsEvictedForcefully() bool {
	return v > ModeEvictedGracefully
}

func (v OpMode) IsEvictedGracefully() bool {
	return v == ModeEvictedGracefully
}

func (v OpMode) IsRestricted() bool {
	return v&ModeFlagRestrictedBehavior != 0
}

func (v OpMode) CanIntroduceJoiner(isJoiner bool) bool {
	return !v.IsRestricted() && !v.IsSuspended() && !isJoiner
}

func (v OpMode) IsMistrustful() bool {
	return v&ModeFlagValidationWarning != 0
}

func (v OpMode) IsSuspended() bool {
	return v&ModeFlagSuspendedOps != 0
}

func (v OpMode) IsPowerless() bool {
	return v.IsSuspended() || v.IsEvicted()
}

func (v OpMode) AsUnit32() uint32 {
	if v >= 1<<ModeBits {
		panic("illegal value")
	}
	return uint32(v)
}

func (v OpMode) CanVote() bool {
	return v == ModeNormal
}

func (v OpMode) CanHaveState() bool {
	// TODO: verify
	return /*!v.IsSuspended() && */ !v.IsEvicted()
}

func (v OpMode) String() string {
	switch v {
	case ModeNormal:
		return "mode:norm"
	case ModeSuspected:
		return "mode:susp"
	case ModePossibleFraud:
		return "mode:warn"
	case ModePossibleFraudAndSuspected:
		return "mode:warn+susp"
	case ModeRestrictedAnnouncement:
		return "mode:joiner"
	case ModeEvictedGracefully:
		return "evict:norm"
	case ModeEvictedAsFraud:
		return "evict:fraud"
	case ModeEvictedAsSuspected:
		return "evict:susp"
	default:
		return fmt.Sprintf("?%d?", v)
	}
}

func (v OpMode) IsClean() bool {
	return v == ModeNormal || v.IsJustJoined()
}
