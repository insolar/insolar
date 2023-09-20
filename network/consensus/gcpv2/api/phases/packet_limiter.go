package phases

import (
	"math"
	"strings"
	"sync/atomic"
)

type PacketLimiter struct {
	received uint16
	sent     uint16
	counters LimitCounters
}

func NewPacketLimiter(maxExtPhase2 uint8) PacketLimiter {
	limits, _ := CreateLimitCounters(maxExtPhase2)
	return PacketLimiter{counters: limits}
}

func NewLocalPacketLimiter() PacketLimiter {
	return PacketLimiter{received: math.MaxUint16}
}

func NewPacketWithOptions(isJoiner bool, maxExtPhase2 uint8) PacketLimiter {
	limits, joinerInits := CreateLimitCounters(maxExtPhase2)
	if !isJoiner {
		joinerInits = 0
	}
	return PacketLimiter{counters: limits, received: joinerInits} // , sent: joinerInits}
}

func (p PacketLimiter) ForJoiner() PacketLimiter {
	_, joinerInits := CreateLimitCounters(0)
	p.received |= joinerInits
	return p
}

func (p PacketLimiter) HasAnyPacketSent() bool {
	return p.sent != 0
}

func (p PacketLimiter) IsLocal() bool {
	return ^p.received == 0
}

func (p PacketLimiter) HasPacketSent(pt PacketType) bool {
	return p.sent&(1<<pt) != 0
}

func (p PacketLimiter) HasAnyPacketReceived() bool {
	return p.received != 0
}

func (p PacketLimiter) HasPacketReceived(pt PacketType) bool {
	return p.received&(1<<pt) != 0
}

func (p PacketLimiter) SetPacketSent(pt PacketType) (bool, PacketLimiter) {
	res := p.sent&(1<<pt) == 0
	p.sent |= 1 << pt
	return res, p
}

func (p PacketLimiter) GetRemainingPacketCount(replaceUnlimitedWith uint8) uint8 {
	count := uint8(0)
	for pt := PacketType(0); pt < PacketType(PacketTypeCount); pt++ {
		switch pt.GetLimitPerSender() {
		case 0:
			continue
		case 1:
			if !p.HasPacketReceived(pt) {
				count++
			}
		default:
			limit := p.counters[pt.GetLimitCounterIndex()]
			if limit == UnlimitedPackets {
				if replaceUnlimitedWith == UnlimitedPackets {
					return UnlimitedPackets
				}
				limit = replaceUnlimitedWith
			}
			count += limit
		}
	}
	return count
}

func (p PacketLimiter) CanReceivePacket(pt PacketType) bool {
	switch pt.GetLimitPerSender() {
	case 1:
		return p.received&(1<<pt) == 0
	case 0:
		return false
	default:
		idx := pt.GetLimitCounterIndex()
		return p.counters[idx] > 0
	}
}

func (p PacketLimiter) SetPacketReceived(pt PacketType) (bool, PacketLimiter) {
	res := p.received&(1<<pt) == 0
	p.received |= 1 << pt
	switch pt.GetLimitPerSender() {
	case 1:
		return res, p
	case 0:
		return false, p
	default:
		idx := pt.GetLimitCounterIndex()
		switch p.counters[idx] {
		case 0:
			return false, p
		case UnlimitedPackets:
			return true, p
		default:
			p.counters[idx]--
			return true, p
		}
	}
}

func (p *PacketLimiter) HasReceivedOrSent() bool {
	return p.received != 0 || p.sent != 0
}

func (p PacketLimiter) asUint64() uint64 {
	v := uint64(p.received)
	v |= uint64(p.sent) << 16
	v |= uint64(p.counters.asUint32()) << 32
	return v
}

func packetLimiterOfUint64(v uint64) PacketLimiter {
	var p PacketLimiter
	p.received = uint16(v)
	p.sent = uint16(v >> 16)
	p.counters = limitCountersOfUint32(uint32(v >> 32))
	return p
}

func (p PacketLimiter) String() string {
	var mode string
	if p.IsLocal() {
		mode = "local"
		if !p.HasAnyPacketSent() {
			return mode + ":idle"
		}
	} else {
		mode = "rmt"
		if !p.HasReceivedOrSent() {
			return mode + ":idle"
		}
	}

	b := strings.Builder{}
	b.WriteString(mode)

	fmtNodeStatePhases(&b, 'S', p.sent, nil)
	if !p.IsLocal() {
		fmtNodeStatePhases(&b, 'R', p.received, &p.counters)
	}

	return b.String()
}

func (p PacketLimiter) GetRemainingPacketCountDefault() uint8 {
	return p.GetRemainingPacketCount(5)
}

func (p PacketLimiter) MergeSent(limiter PacketLimiter) PacketLimiter {
	p.sent |= limiter.sent
	return p
}

func fmtNodeStatePhases(b *strings.Builder, prefix byte, ns uint16, limits *LimitCounters) {

	if ns == 0 {
		return
	}
	b.WriteByte(':')
	b.WriteByte(prefix)
	fmtNodeStatePhasesSubset(ns&(1<<PacketOffPhase-1), b, 0, nil)

	ns >>= PacketOffPhase
	if ns == 0 {
		return
	}

	b.WriteByte('-')
	fmtNodeStatePhasesSubset(ns, b, PacketOffPhase, limits)
}

func fmtNodeStatePhasesSubset(ns0 uint16, b *strings.Builder, pt PacketType, limits *LimitCounters) {

	for ; ns0 != 0; ns0 >>= 1 {
		if ns0&1 != 0 {
			b.WriteRune(pt.RuneName())
			if limits != nil {
				limitIndex := pt.GetLimitCounterIndex()
				if limitIndex >= 0 && (*limits)[limitIndex] > 0 {
					b.WriteByte('[')
					limits.WriteLimitTo(b, limitIndex)
					b.WriteByte(']')
				}
			}
		} else {
			b.WriteByte('_')
		}
		pt++
	}
}

func NewAtomicPacketLimiter(initial PacketLimiter) *AtomicPacketLimiter {
	return &AtomicPacketLimiter{initial.asUint64()}
}

type AtomicPacketLimiter struct {
	packetLimiter uint64
}

func (p *AtomicPacketLimiter) GetPacketLimiter() PacketLimiter {
	return packetLimiterOfUint64(atomic.LoadUint64(&p.packetLimiter))
}

func (p *AtomicPacketLimiter) UpdatePacketLimiter(prev PacketLimiter, new PacketLimiter) bool {
	return atomic.CompareAndSwapUint64(&p.packetLimiter, prev.asUint64(), new.asUint64())
}

func (p *AtomicPacketLimiter) SetPacketReceived(pt PacketType) bool {

	for {
		prev := atomic.LoadUint64(&p.packetLimiter)
		res, upd := packetLimiterOfUint64(prev).SetPacketReceived(pt)
		if !res {
			return false
		}

		if atomic.CompareAndSwapUint64(&p.packetLimiter, prev, upd.asUint64()) {
			return true
		}
	}
}
