package phases

import (
	"fmt"
	"math"
	"strings"
)

type PacketType uint8

const (
	/* Phased Packets - these are SENT by a node only ONCE and in the given sequence */

	PacketPhase0 PacketType = iota
	PacketPhase1
	PacketPhase2
	PacketPhase3
	PacketPhase4
	_ // 5
	_ // 6
	_ // 7

	PacketPulsarPulse // PacketPhase0 | PacketOffPhase /* Triggers Phase0-1 */

	PacketReqPhase1 // PacketPhase1 | PacketOffPhase
	/*  Request to resend own NSH - will be replied with PacketPhase1 without PulseData.
	The reply MUST include all data (except for PulseData) as it was presented in the original Phase1 packet.
	This request MUST be replied not more than 1-2 times per requesting node per consensus round,
	otherwise is ignored.
	*/

	PacketExtPhase2 // PacketPhase2 | PacketOffPhase
	/*	And additional Phase 2 packet to improve coverage for fraud detection,
		but it doesn't increase trust-level, as can be exploited by sending multiple times.
	*/

	PacketFastPhase3 // PacketPhase3 | PacketOffPhase
	/* Out-of-order Phase3 packet that can be sent before or during Phase 2. Can only be sent once.	*/

	// PacketReqIntro /* Request to resend other's (NSH + intro) - will be replied with PacketPhase2.
	// Only joiners can send this request, and only to anyone in a relevant neighbourhood.
	// Limited by 1 times per requesting node per consensus round per requested intro,
	// otherwise is ignored.
	// PacketReqFraud /* Requests fraud proof */
	// PacketFraud /* Delivers fraud proof, by request only */

	maxPacketType
)

const PacketOffPhase = 8
const PacketTypeCount = int(maxPacketType)
const UnlimitedPackets = math.MaxUint8

// TODO TEST must correlate: (p.GetLimitPerSender()<=1)==(p.GetLimitCounterIndex()<0)
func (p PacketType) GetLimitPerSender() uint8 {
	switch p {
	case PacketPhase0, PacketPhase1, PacketPhase2, PacketPhase3, PacketPhase4:
		return 1
	case PacketPulsarPulse:
		return 1
	case PacketReqPhase1:
		return 2
	case PacketExtPhase2:
		return UnlimitedPackets
	case PacketFastPhase3:
		return 1
	default:
		return 0 // packet is not allowed
	}
}

// TODO TEST must correlate: GetLimitCounterIndex() must be unique for every p, and less than PacketCountedLimits
func (p PacketType) GetLimitCounterIndex() int {
	switch p {
	case PacketReqPhase1:
		return 0
	case PacketExtPhase2:
		return 1
	default:
		return -1
	}
}

func (p PacketType) IsAllowedForJoiner() bool {
	switch p {
	case PacketPulsarPulse, PacketPhase0:
		return false
	default:
		return true
	}
}

const PacketCountedLimits = 2

type LimitCounters [PacketCountedLimits]uint8

func (v LimitCounters) WriteLimitTo(b *strings.Builder, index int) {
	vv := v[index]
	switch vv {
	case UnlimitedPackets:
		b.WriteRune('∞')
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9:
		b.WriteByte('0' + vv)
	default:
		b.WriteString(fmt.Sprint(vv))
	}
}

func limitCountersOfUint32(r uint32) LimitCounters {

	var v LimitCounters
	for i := range v {
		if i >= 4 {
			panic("unsupported")
		}
		v[i] = uint8(r >> uint32(i<<3))
	}
	return v
}

func (v LimitCounters) asUint32() uint32 {
	r := uint32(0)
	for i, vv := range v {
		if i >= 4 {
			panic("unsupported")
		}
		r |= uint32(vv) << uint32(i<<3)
	}
	return r
}

func (v LimitCounters) String() string {
	res := strings.Builder{}
	res.WriteByte('[')
	for i := range v {
		if i != 0 {
			res.WriteByte(' ')
		}
		v.WriteLimitTo(&res, i)
	}
	res.WriteByte(']')
	return res.String()
}

var limitCounters LimitCounters
var joinerInits uint16

func fillLimitCounters() (LimitCounters, uint16) {
	var limits LimitCounters
	var inits = uint16(0)
	idx := 0
	for i := PacketType(0); i < maxPacketType; i++ {
		limit := i.GetLimitPerSender()
		if limit <= 1 {
			continue
		}
		if idx != i.GetLimitCounterIndex() {
			panic("illegal state")
		}
		limits[idx] = limit
		idx++

		if !i.IsAllowedForJoiner() {
			inits |= 1 << i
		}
	}
	return limits, inits
}

func CreateLimitCounters(maxExtPhase2 uint8) (LimitCounters, uint16) {
	r := limitCounters
	r[PacketExtPhase2.GetLimitCounterIndex()] = maxExtPhase2
	return r, joinerInits
}

func (p PacketType) IsPhasedPacket() bool {
	return p < PacketOffPhase
}

func (p PacketType) IsMemberPacket() bool {
	return p < maxPacketType && p != PacketPulsarPulse
}

func (p PacketType) IsEphemeralPacket() bool {
	return p != PacketPulsarPulse
}

func (p PacketType) GetPayloadEquivalent() PacketType {
	switch p {
	case PacketReqPhase1:
		return PacketPhase1
	case PacketExtPhase2:
		return PacketPhase2
	case PacketFastPhase3:
		return PacketPhase3
	default:
		return p
	}
}

func (p PacketType) ToPhaseNumber() (Number, bool) {
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

func (p PacketType) String() string {
	switch p {
	case PacketPhase0:
		return "ph0"
	case PacketPhase1:
		return "ph1"
	case PacketPhase2:
		return "ph2"
	case PacketPhase3:
		return "ph3"
	case PacketPhase4:
		return "ph4"
	case PacketPulsarPulse:
		return "pulse"
	case PacketReqPhase1:
		return "ph1rq"
	case PacketExtPhase2:
		return "ph2ex"
	case PacketFastPhase3:
		return "ph3ft"
	default:
		return fmt.Sprintf("packetType%d", p)
	}
}

func (p PacketType) RuneName() rune {
	return rune([]byte("01234   prxf")[p])
}

func init() {
	limitCounters, joinerInits = fillLimitCounters()
}
