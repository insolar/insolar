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

package ph3ctl

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/statevector"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/nodeset"
)

type VectorBuilder struct {
	digestFactory transport.ConsensusDigestFactory
	entryScanner  nodeset.VectorEntryScanner
	bitset        member.StateBitset
}

func (p *VectorBuilder) CreateDerived(comparedStats nodeset.ComparedBitsetRow) VectorBuilder {
	return VectorBuilder{p.digestFactory, p.entryScanner,
		p.buildDerivedBitset(comparedStats)}
}

func (p *VectorBuilder) buildDerivedBitset(comparedStats nodeset.ComparedBitsetRow) member.StateBitset {

	bitset := make(member.StateBitset, len(p.bitset))

	for idx := range p.bitset {
		switch comparedStats.Get(idx) {
		case nodeset.ComparedMissingHere:
			bitset[idx] = member.BeTimeout // we don't have it
		case nodeset.ComparedDoubtedMissingHere:
			bitset[idx] = member.BeTimeout // we don't have it
		case nodeset.ComparedSame:
			// ok, use as-is
			bitset[idx] = p.bitset[idx]
		case nodeset.ComparedLessTrustedThere:
			// ok - exclude for trusted
			bitset[idx] = member.BeBaselineTrust
		case nodeset.ComparedLessTrustedHere:
			// ok - use for both
			bitset[idx] = member.BeHighTrust
		case nodeset.ComparedMissingThere:
			bitset[idx] = member.BeTimeout // we have it, but the other's doesn't
		default:
			panic("unexpected")
		}
	}

	return bitset
}

func (p *VectorBuilder) buildGlobulaAnnouncementHash(trusted bool) proofs.GlobulaAnnouncementHash {

	calc := nodeset.NewAnnouncementSequenceCalc(p.digestFactory)
	p.entryScanner.ScanIndexed(func(idx int, nodeData nodeset.VectorEntryData) {
		b := p.bitset[idx]
		if b.IsTimeout() {
			return
		}
		if trusted && !b.IsTrusted() {
			return
		}
		calc.AddNext(nodeData, false)
	})
	return calc.FinishSequence()
}

func (p *VectorBuilder) buildGlobulaAnnouncementHashes() (proofs.GlobulaAnnouncementHash, proofs.GlobulaAnnouncementHash) {

	calcTrusted := nodeset.NewAnnouncementSequenceCalc(p.digestFactory)
	var calcDoubted nodeset.AnnouncementSequenceCalc

	p.entryScanner.ScanIndexed(
		func(idx int, nodeData nodeset.VectorEntryData) {
			b := p.bitset[idx]
			if b.IsTimeout() {
				return
			}
			if b.IsTrusted() {
				calcTrusted.AddNext(nodeData, false)
				if calcDoubted.IsEmpty() {
					return
				}
			} else {
				if calcDoubted.IsEmpty() {
					calcDoubted.ForkSequenceOf(calcTrusted)
				}
			}
			calcDoubted.AddNext(nodeData, false)
		})

	return calcTrusted.FinishSequence(), calcDoubted.FinishSequence()
}

func (p *VectorBuilder) buildGlobulaStateHash(trusted bool, nodeID insolar.ShortNodeID) statevector.CalcStateWithRank {

	calc := nodeset.NewStateAndRankSequenceCalc(p.digestFactory, nodeID,
		1+p.entryScanner.GetSortedCount()>>1)

	const (
		skipEntry = iota
		postponeEntry
		normalEntry
	)

	p.entryScanner.ScanSortedWithFilter(
		func(nodeData nodeset.VectorEntryData, filter uint32) {
			if filter == skipEntry {
				return
			}
			// TODO use default state hash on missing data
			calc.AddNext(nodeData, filter == postponeEntry)
		},
		func(idx int, nodeData nodeset.VectorEntryData) (bool, uint32) {
			b := p.bitset[idx]
			if b.IsTimeout() || nodeData.RequestedPower == 0 || nodeData.RequestedMode.IsPowerless() || trusted && !b.IsTrusted() {
				return true, postponeEntry
			}
			return false, normalEntry
		})

	tHash, tRank, tCount := calc.FinishSequence()
	return statevector.CalcStateWithRank{StateHash: tHash, ExpectedRank: tRank.AsMembershipRank(tCount)}
}

func (p *VectorBuilder) BuildGlobulaAnnouncementHashes(buildTrusted, buildDoubted bool,
	defaultTrusted, defaultDoubted proofs.GlobulaAnnouncementHash) (trustedHash, doubtedHash proofs.GlobulaAnnouncementHash) {

	if buildTrusted && buildDoubted {
		t, d := p.buildGlobulaAnnouncementHashes()
		if d == nil {
			return t, t
		}
		return t, d
	}
	if buildTrusted {
		return p.buildGlobulaAnnouncementHash(true), defaultDoubted
	}
	if buildDoubted {
		return defaultTrusted, p.buildGlobulaAnnouncementHash(false)
	}
	return defaultTrusted, defaultDoubted
}

func (p *VectorBuilder) BuildGlobulaStateHashesAndRanks(buildTrusted, buildDoubted bool, nodeID insolar.ShortNodeID,
	defaultTrusted, defaultDoubted statevector.CalcStateWithRank) (trustedHash, doubtedHash statevector.CalcStateWithRank) {

	if buildTrusted {
		defaultTrusted = p.buildGlobulaStateHash(true, nodeID)
	}
	if buildDoubted {
		defaultDoubted = p.buildGlobulaStateHash(false, nodeID)
	}
	return defaultTrusted, defaultDoubted
}
