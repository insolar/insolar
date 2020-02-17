// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package nodeset

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/statevector"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
)

type VectorBuilder struct {
	digestFactory transport.ConsensusDigestFactory
	entryScanner  VectorEntryScanner
	bitset        member.StateBitset
}

func NewVectorBuilder(df transport.ConsensusDigestFactory, s VectorEntryScanner, bitset member.StateBitset) VectorBuilder {
	return VectorBuilder{df, s, bitset}
}

func (p *VectorBuilder) FillBitset() {
	p.entryScanner.ScanIndexed(func(idx int, nodeData VectorEntryData) {
		p.bitset[idx] = mapVectorEntryDataToNodesetEntry(nodeData)
	})
}

func mapVectorEntryDataToNodesetEntry(nodeData VectorEntryData) member.BitsetEntry {
	switch {
	case nodeData.IsEmpty():
		return member.BeTimeout
	case nodeData.TrustLevel.IsNegative():
		return member.BeFraud
	case nodeData.TrustLevel == member.UnknownTrust:
		return member.BeBaselineTrust
	case nodeData.TrustLevel < member.TrustByNeighbors:
		return member.BeLimitedTrust
	default:
		return member.BeHighTrust
	}
}

func (p *VectorBuilder) GetBitset() member.StateBitset {
	return p.bitset
}

func (p *VectorBuilder) GetEntryScanner() VectorEntryScanner {
	return p.entryScanner
}

func (p *VectorBuilder) CreateDerived(comparedStats ComparedBitsetRow) VectorBuilder {
	return VectorBuilder{p.digestFactory, p.entryScanner,
		p.buildDerivedBitset(comparedStats)}
}

func (p *VectorBuilder) buildDerivedBitset(comparedStats ComparedBitsetRow) member.StateBitset {

	bitset := make(member.StateBitset, len(p.bitset))

	for idx := range p.bitset {
		switch comparedStats.Get(idx) {
		case ComparedMissingHere:
			bitset[idx] = member.BeTimeout // we don't have it
		case ComparedDoubtedMissingHere:
			bitset[idx] = member.BeTimeout // we don't have it
		case ComparedSame:
			// ok, use as-is
			bitset[idx] = p.bitset[idx]
		case ComparedLessTrustedThere:
			// ok - exclude for trusted
			bitset[idx] = member.BeBaselineTrust
		case ComparedLessTrustedHere:
			// ok - use for both
			bitset[idx] = member.BeHighTrust
		case ComparedMissingThere:
			bitset[idx] = member.BeTimeout // we have it, but the other's doesn't
		default:
			panic("unexpected")
		}
	}

	return bitset
}

func (p *VectorBuilder) buildGlobulaAnnouncementHash(trusted bool) proofs.GlobulaAnnouncementHash {

	calc := NewAnnouncementSequenceCalc(p.digestFactory)
	p.entryScanner.ScanIndexed(func(idx int, nodeData VectorEntryData) {
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

func (p *VectorBuilder) BuildAllGlobulaAnnouncementHashes() (proofs.GlobulaAnnouncementHash, proofs.GlobulaAnnouncementHash) {

	calcTrusted := NewAnnouncementSequenceCalc(p.digestFactory)
	var calcDoubted AnnouncementSequenceCalc

	p.entryScanner.ScanIndexed(
		func(idx int, nodeData VectorEntryData) {
			b := p.bitset[idx]
			if b.IsTimeout() {
				return
			}
			if b.IsTrusted() {
				calcTrusted.AddNext(nodeData, false)
				if calcDoubted.IsEmpty() {
					return
				}
			} else if calcDoubted.IsEmpty() {
				calcDoubted.ForkSequenceOf(calcTrusted)
			}
			calcDoubted.AddNext(nodeData, false)
		})

	return calcTrusted.FinishSequence(), calcDoubted.FinishSequence()
}

// TODO reuse BuildGlobulaStateHashWithFilter
func (p *VectorBuilder) buildGlobulaStateHash(trusted bool, nodeID insolar.ShortNodeID) statevector.CalcStateWithRank {

	calc := NewStateAndRankSequenceCalc(p.digestFactory, nodeID,
		1+p.entryScanner.GetSortedCount()>>1)

	p.entryScanner.ScanSortedWithFilter(0,
		func(nodeData VectorEntryData, postponed bool, filter uint32) {
			calc.AddNext(nodeData, postponed)
		},
		func(idx int, nodeData VectorEntryData, parentFilter uint32) (bool, uint32) {
			if nodeData.IsEmpty() || nodeData.RequestedPower == 0 || nodeData.RequestedMode.IsPowerless() {
				return true, 0
			}

			if idx < 0 { // this is joiner check - it should only indicate "postpone"
				return false, 0
			}

			if trusted && !p.bitset[idx].IsTrusted() {
				return true, 0
			}
			return false, 0
		})

	tHash, tRank, tCount := calc.FinishSequence()
	return statevector.CalcStateWithRank{StateHash: tHash, ExpectedRank: tRank.AsMembershipRank(tCount)}
}

func (p *VectorBuilder) BuildGlobulaAnnouncementHashes(buildTrusted, buildDoubted bool,
	defaultTrusted, defaultDoubted proofs.GlobulaAnnouncementHash) (trustedHash, doubtedHash proofs.GlobulaAnnouncementHash) {

	if buildTrusted && buildDoubted {
		t, d := p.BuildAllGlobulaAnnouncementHashes()
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

func (p *VectorBuilder) BuildGlobulaStateHashWithFilter(nodeID insolar.ShortNodeID, apply EntryApplyFunc,
	filter EntryFilterFunc) statevector.CalcStateWithRank {

	calc := NewStateAndRankSequenceCalc(p.digestFactory, nodeID,
		1+p.entryScanner.GetSortedCount()>>1)

	p.entryScanner.ScanSortedWithFilter(0,
		func(nodeData VectorEntryData, postponed bool, filter uint32) {
			apply(nodeData, postponed, filter)
			calc.AddNext(nodeData, postponed)
		},
		filter)

	tHash, tRank, tCount := calc.FinishSequence()
	return statevector.CalcStateWithRank{StateHash: tHash, ExpectedRank: tRank.AsMembershipRank(tCount)}
}
