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

package population

import (
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"sync/atomic"
)

var _ EventDispatcher = &AtomicEventStats{}

type AtomicEventStats struct {
	purgatoryCounts  uint32
	dynamicsCounts   uint32
	trustLevelCounts uint64
}

func (p *AtomicEventStats) OnTrustUpdated(populationVersion uint32, n *NodeAppearance, trustBefore member.TrustLevel, trustAfter member.TrustLevel) {
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
