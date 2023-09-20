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
