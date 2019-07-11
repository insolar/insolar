package nodeset

import (
	common2 "github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/stats"
	"math"
)

type NodeVectorHelper struct {
	digestFactory     common2.DigestFactory
	signatureVerifier common2.SignatureVerifier
	entryScanner      VectorEntryScanner
	bitset            NodeBitset
}

func NewLocalNodeVector(digestFactory common2.DigestFactory, signatureVerifier common2.SignatureVerifier,
	entryScanner VectorEntryScanner) NodeVectorHelper {

	p := NodeVectorHelper{digestFactory, signatureVerifier,
		entryScanner, make(NodeBitset, entryScanner.GetIndexedCount())}

	entryScanner.ScanIndexed(func(idx int, nodeData VectorEntryData) {
		p.bitset[idx] = mapVectorEntryDataToNodesetEntry(nodeData)
	})
	return p
}

func mapVectorEntryDataToNodesetEntry(nodeData VectorEntryData) NodeBitsetEntry {
	switch {
	case nodeData.IsEmpty():
		return NbsTimeout
	case nodeData.TrustLevel.IsNegative():
		return NbsFraud
	case nodeData.TrustLevel == common.UnknownTrust:
		return NbsBaselineTrust
	case nodeData.TrustLevel < common.TrustByNeighbors:
		return NbsLimitedTrust
	default:
		return NbsHighTrust
	}
}

func (p *NodeVectorHelper) CreateDerivedVector(statRow *stats.Row) NodeVectorHelper {
	dv := NodeVectorHelper{p.digestFactory, p.signatureVerifier,
		p.entryScanner, make(NodeBitset, len(p.bitset))}

	for idx := range p.bitset {
		switch statRow.Get(idx) {
		case NodeBitMissingHere:
			dv.bitset[idx] = NbsTimeout // we don't have it
		case NodeBitDoubtedMissingHere:
			dv.bitset[idx] = NbsTimeout // we don't have it
		case NodeBitSame:
			// ok, use as-is
			dv.bitset[idx] = p.bitset[idx]
		case NodeBitLessTrustedThere:
			// ok - exclude for trusted
			dv.bitset[idx] = NbsBaselineTrust
		case NodeBitLessTrustedHere:
			// ok - use for both
			dv.bitset[idx] = NbsHighTrust
		case NodeBitMissingThere:
			dv.bitset[idx] = NbsTimeout // we have it, but the other's doesn't
		default:
			panic("unexpected")
		}
	}

	return dv
}

func (p *NodeVectorHelper) buildGlobulaAnnouncementHash(trusted bool) common.GlobulaAnnouncementHash {
	hasEntries := false
	agg := p.digestFactory.GetGshDigester()

	p.entryScanner.ScanIndexed(func(idx int, nodeData VectorEntryData) {
		filter := p.bitset[idx]
		if trusted && !filter.IsTrusted() {
			return
		}
		agg.AddNext(nodeData.StateEvidence.GetNodeStateHash())
		hasEntries = true
	})

	if hasEntries {
		return agg.FinishSequence().AsDigestHolder()
	}
	return nil
}

func (p *NodeVectorHelper) buildGlobulaAnnouncementHashes() (common.GlobulaAnnouncementHash, common.GlobulaAnnouncementHash) {
	/*
		NB! SequenceDigester requires at least one hash to be added. So to avoid errors, local node MUST always
		have trust level set high enough to get bitset[i].IsTrusted() == true
	*/

	aggTrusted := p.digestFactory.GetGshDigester()
	var aggDoubted common2.SequenceDigester

	p.entryScanner.ScanIndexed(
		func(idx int, nodeData VectorEntryData) {
			b := p.bitset[idx]
			if b.IsTrusted() {
				aggTrusted.AddNext(nodeData.StateEvidence.GetNodeStateHash())
				if aggDoubted == nil {
					return
				}
			} else if aggDoubted == nil {
				aggDoubted = aggTrusted.ForkSequence()
			}
			aggDoubted.AddNext(nodeData.StateEvidence.GetNodeStateHash())
		})

	trustedResult := aggTrusted.FinishSequence().AsDigestHolder()
	if aggDoubted != nil {
		return trustedResult, aggDoubted.FinishSequence().AsDigestHolder()
	}
	return trustedResult, trustedResult
}

func (p *NodeVectorHelper) buildGlobulaStateHash(trusted bool) common.GlobulaAnnouncementHash {
	hasEntries := false
	agg := p.digestFactory.GetGshDigester()

	const skip = math.MaxUint32

	p.entryScanner.ScanSortedWithFilter(
		func(nodeData VectorEntryData, filter uint32) {
			if filter == skip {
				return
			}
			b := p.bitset[filter]
			if trusted && !b.IsTrusted() {
				return
			}
			digest := common2.NewDigest(nodeData.AnnounceSignature, "").AsDigestHolder()

			agg.AddNext(digest)
			hasEntries = true
		},
		func(idx int, nodeData VectorEntryData) (bool, uint32) {
			postpone := p.bitset[idx].IsTimeout() || nodeData.RequestedPower == 0
			return postpone, uint32(idx)
		})

	if hasEntries {
		return agg.FinishSequence().AsDigestHolder()
	}
	return nil
}

func (p *NodeVectorHelper) buildGlobulaStateHashes() (common.GlobulaAnnouncementHash, common.GlobulaAnnouncementHash) {
	/*
		NB! SequenceDigester requires at least one hash to be added. So to avoid errors, local node MUST always
		have trust level set high enough to get bitset[i].IsTrusted() == true
	*/

	aggTrusted := p.digestFactory.GetGshDigester()
	var aggDoubted common2.SequenceDigester

	const skip = math.MaxUint32

	p.entryScanner.ScanSortedWithFilter(
		func(nodeData VectorEntryData, filter uint32) {
			if filter == skip {
				return
			}
			b := p.bitset[filter]
			digest := common2.NewDigest(nodeData.AnnounceSignature, "").AsDigestHolder()

			if b.IsTrusted() {
				aggTrusted.AddNext(digest)
				if aggDoubted == nil {
					return
				}
			} else if aggDoubted == nil {
				aggDoubted = aggTrusted.ForkSequence()
			}
			aggDoubted.AddNext(digest)
		},
		func(idx int, nodeData VectorEntryData) (bool, uint32) {
			postpone := p.bitset[idx].IsTimeout() || nodeData.RequestedPower == 0
			return postpone, uint32(idx)
		})

	trustedResult := aggTrusted.FinishSequence().AsDigestHolder()
	if aggDoubted != nil {
		return trustedResult, aggDoubted.FinishSequence().AsDigestHolder()
	}
	return trustedResult, trustedResult
}

func (p *NodeVectorHelper) BuildGlobulaAnnouncementHashes(buildTrusted, buildDoubted bool,
	defaultTrusted, defaultDoubted common.GlobulaAnnouncementHash) (trustedHash, doubtedHash common.GlobulaAnnouncementHash) {

	if buildTrusted && buildDoubted {
		return p.buildGlobulaAnnouncementHashes()
	}
	if buildTrusted {
		return p.buildGlobulaAnnouncementHash(true), defaultDoubted
	}
	if buildDoubted {
		return defaultTrusted, p.buildGlobulaAnnouncementHash(false)
	}
	return defaultTrusted, defaultDoubted
}

func (p *NodeVectorHelper) BuildGlobulaStateHashes(buildTrusted, buildDoubted bool,
	defaultTrusted, defaultDoubted common.GlobulaStateHash) (trustedHash, doubtedHash common.GlobulaStateHash) {

	if buildTrusted && buildDoubted {
		return p.buildGlobulaStateHashes()
	}
	if buildTrusted {
		return p.buildGlobulaStateHash(true), defaultDoubted
	}
	if buildDoubted {
		return defaultTrusted, p.buildGlobulaStateHash(false)
	}
	return defaultTrusted, defaultDoubted
}

func (p *NodeVectorHelper) VerifyGlobulaStateSignature(localHash common.GlobulaStateHash, remoteSignature common2.SignatureHolder) bool {
	return localHash != nil && p.signatureVerifier.IsValidDigestSignature(localHash, remoteSignature)
}

func (p *NodeVectorHelper) GetNodeBitset() NodeBitset {
	return p.bitset
}

//type vectorSubFilter uint8
//
//func (v vectorSubFilter) CanApply(filter NodeBitsetEntry) bool {
//
//}
//
//type vectorBuilder struct {
//	empty common2.SequenceDigester
//	inclusive bool
//	filters [2]vectorSubFilter
//	vectors [2]common2.SequenceDigester
//	extractor func(nodeData VectorEntryData) common2.DigestHolder
//
//}
//
//func (p *vectorBuilder) applyExclusive(nodeData VectorEntryData, filter NodeBitsetEntry) {
//	value := p.extractor(nodeData)
//	prev := false
//
//	for i, f := range p.filters {
//		if f.CanApply(filter)
//	}
//	if filter.IsTrusted() {
//		if p.required[0] {
//
//		}
//	}
//	if p.required[1] {
//
//	}
//}
//
//func (p *NodeVectorHelper) buildHashes() (common.GlobulaAnnouncementHash, common.GlobulaAnnouncementHash) {
//	/*
//		NB! SequenceDigester requires at least one hash to be added. So to avoid errors, local node MUST always
//		have trust level set high enough to get bitset[i].IsTrusted() == true
//	*/
//
//	aggTrusted := p.digestFactory.GetGshDigester()
//	var aggDoubted common2.SequenceDigester
//
//	const skip = math.MaxUint32
//
//	p.entryScanner.ScanSortedWithFilter(
//		func(nodeData VectorEntryData, filter uint32) {
//			if filter == skip {	return }
//			b := p.bitset[filter]
//			if b.IsTrusted() {
//				aggTrusted.AddNext(nodeData.StateEvidence.GetNodeStateHash())
//				if aggDoubted == nil {
//					return
//				}
//			} else if aggDoubted == nil {
//				aggDoubted = aggTrusted.ForkSequence()
//			}
//			aggDoubted.AddNext(nodeData.StateEvidence.GetNodeStateHash())
//		},
//		func(idx int, nodeData VectorEntryData) (bool, uint32) {
//			postpone := p.bitset[idx].IsTimeout() || nodeData.RequestedPower == 0
//			return postpone, uint32(idx)
//		})
//
//	trustedResult := aggTrusted.FinishSequence().AsDigestHolder()
//	if aggDoubted != nil {
//		return trustedResult, aggDoubted.FinishSequence().AsDigestHolder()
//	}
//	return trustedResult, trustedResult
//}
