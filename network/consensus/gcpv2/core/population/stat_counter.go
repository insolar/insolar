// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package population

import (
	"sync/atomic"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
)

var _ EventDispatcher = &AtomicEventStats{}

type AtomicEventStats struct {
	purgatoryCounts  uint32
	dynamicsCounts   uint32
	trustLevelCounts uint64
}

func (p *AtomicEventStats) OnTrustUpdated(populationVersion uint32, n *NodeAppearance,
	trustBefore member.TrustLevel, trustAfter member.TrustLevel, hasFullProfile bool) {

	delta := uint64(0)

	switch {
	case trustBefore == trustAfter:
		return
	case trustAfter.IsNegative():
		if trustBefore.IsNegative() {
			return
		}
		delta |= 1
	default:
		if trustBefore == member.UnknownTrust && trustAfter >= member.TrustBySelf {
			delta |= 1 << 16
		}
		if trustBefore < member.TrustBySome && trustAfter >= member.TrustBySome {
			delta |= 1 << 32
		}
		if trustBefore < member.TrustByNeighbors && trustAfter >= member.TrustByNeighbors {
			delta |= 1 << 48
		}
		if delta == 0 {
			return
		}
	}
	atomic.AddUint64(&p.trustLevelCounts, delta)
}

func (p *AtomicEventStats) GetTrustCounts() (fraudCount, bySelfCount, bySomeCount, byNeighborsCount uint16) {
	dc := atomic.LoadUint64(&p.trustLevelCounts)
	return uint16(dc), uint16(dc >> 16), uint16(dc >> 32), uint16(dc >> 48)
}

func (p *AtomicEventStats) OnDynamicNodeUpdate(populationVersion uint32, n *NodeAppearance, flags UpdateFlags) {

	if flags&(FlagFixedInit) != 0 {
		return // not a dynamic node
	}
	delta := uint32(0)
	if flags&(FlagCreated) != 0 {
		delta |= 1
	}
	if flags&FlagUpdatedProfile != 0 {
		delta |= 1 << 16
	}
	atomic.AddUint32(&p.dynamicsCounts, delta)
}

func (p *AtomicEventStats) GetDynamicCounts() (briefCount, fullCount uint16) {
	dc := atomic.LoadUint32(&p.dynamicsCounts)
	return uint16(dc), uint16(dc >> 16)
}

func (p *AtomicEventStats) OnPurgatoryNodeUpdate(populationVersion uint32, n MemberPacketSender, flags UpdateFlags) {
	delta := uint32(0)
	if flags&FlagCreated != 0 {
		delta |= 1
	}
	if flags&FlagAscent != 0 {
		delta |= 1 << 16
	}
	if delta != 0 {
		atomic.AddUint32(&p.purgatoryCounts, delta)
	}
}

func (p *AtomicEventStats) GetPurgatoryCounts() (addedCount, ascentCount uint16) {
	dc := atomic.LoadUint32(&p.purgatoryCounts)
	return uint16(dc), uint16(dc >> 16)
}

func (p *AtomicEventStats) OnCustomEvent(populationVersion uint32, n *NodeAppearance, event interface{}) {
}

func (p *AtomicEventStats) OnDynamicPopulationCompleted(populationVersion uint32, indexedCount int) {
}

func (p *AtomicEventStats) OnNodeStateAssigned(populationVersion uint32, n *NodeAppearance) {
}
