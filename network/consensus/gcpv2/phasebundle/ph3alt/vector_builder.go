package ph3alt

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
	return statevector.CalcStateWithRank{tHash, tRank.AsMembershipRank(tCount)}
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
