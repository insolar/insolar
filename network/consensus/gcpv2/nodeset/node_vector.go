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

package nodeset

import (
	"math"

	"github.com/insolar/insolar/network/consensus/common/cryptography_containers"
	"github.com/insolar/insolar/network/consensus/gcpv2/gcp_types"

	"github.com/insolar/insolar/network/consensus/gcpv2/stats"
)

type NodeVectorHelper struct {
	digestFactory     cryptography_containers.DigestFactory
	signatureVerifier cryptography_containers.SignatureVerifier
	entryScanner      VectorEntryScanner
	bitset            gcp_types.NodeBitset
	parentBitset      gcp_types.NodeBitset
}

func NewLocalNodeVector(digestFactory cryptography_containers.DigestFactory,
	entryScanner VectorEntryScanner) NodeVectorHelper {

	p := NodeVectorHelper{digestFactory, nil,
		entryScanner, make(gcp_types.NodeBitset, entryScanner.GetIndexedCount()), nil,
	}

	entryScanner.ScanIndexed(func(idx int, nodeData VectorEntryData) {
		p.bitset[idx] = mapVectorEntryDataToNodesetEntry(nodeData)
	})
	return p
}

func mapVectorEntryDataToNodesetEntry(nodeData VectorEntryData) gcp_types.NodeBitsetEntry {
	switch {
	case nodeData.IsEmpty():
		return gcp_types.NbsTimeout
	case nodeData.TrustLevel.IsNegative():
		return gcp_types.NbsFraud
	case nodeData.TrustLevel == gcp_types.UnknownTrust:
		return gcp_types.NbsBaselineTrust
	case nodeData.TrustLevel < gcp_types.TrustByNeighbors:
		return gcp_types.NbsLimitedTrust
	default:
		return gcp_types.NbsHighTrust
	}
}

func (p *NodeVectorHelper) CreateDerivedVector(signatureVerifier cryptography_containers.SignatureVerifier) NodeVectorHelper {
	return NodeVectorHelper{p.digestFactory, signatureVerifier,
		p.entryScanner, nil, p.bitset}
}

func (p *NodeVectorHelper) PrepareDerivedVector(statRow *stats.Row) {
	if p.bitset != nil && p.parentBitset == nil {
		panic("illegal state")
	}

	p.bitset = make(gcp_types.NodeBitset, len(p.parentBitset))

	for idx := range p.parentBitset {
		switch statRow.Get(idx) {
		case NodeBitMissingHere:
			p.bitset[idx] = gcp_types.NbsTimeout // we don't have it
		case NodeBitDoubtedMissingHere:
			p.bitset[idx] = gcp_types.NbsTimeout // we don't have it
		case NodeBitSame:
			// ok, use as-is
			p.bitset[idx] = p.parentBitset[idx]
		case NodeBitLessTrustedThere:
			// ok - exclude for trusted
			p.bitset[idx] = gcp_types.NbsBaselineTrust
		case NodeBitLessTrustedHere:
			// ok - use for both
			p.bitset[idx] = gcp_types.NbsHighTrust
		case NodeBitMissingThere:
			p.bitset[idx] = gcp_types.NbsTimeout // we have it, but the other's doesn't
		default:
			panic("unexpected")
		}
	}
}

func (p *NodeVectorHelper) buildGlobulaAnnouncementHash(trusted bool) gcp_types.GlobulaAnnouncementHash {
	hasEntries := false
	agg := p.digestFactory.GetGshDigester()

	p.entryScanner.ScanIndexed(func(idx int, nodeData VectorEntryData) {
		b := p.bitset[idx]
		if b.IsTimeout() {
			return
		}
		if trusted && !b.IsTrusted() {
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

func (p *NodeVectorHelper) buildGlobulaAnnouncementHashes() (gcp_types.GlobulaAnnouncementHash, gcp_types.GlobulaAnnouncementHash) {
	/*
		NB! SequenceDigester requires at least one hash to be added. So to avoid errors, local node MUST always
		have trust level set high enough to get bitset[i].IsTrusted() == true
	*/

	aggTrusted := p.digestFactory.GetGshDigester()
	var aggDoubted cryptography_containers.SequenceDigester

	p.entryScanner.ScanIndexed(
		func(idx int, nodeData VectorEntryData) {
			b := p.bitset[idx]
			if b.IsTimeout() {
				return
			}
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

func (p *NodeVectorHelper) buildGlobulaStateHash(trusted bool) gcp_types.GlobulaAnnouncementHash {
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
			digest := cryptography_containers.NewDigest(nodeData.AnnounceSignature, "").AsDigestHolder()

			agg.AddNext(digest)
			hasEntries = true
		},
		func(idx int, nodeData VectorEntryData) (bool, uint32) {
			b := p.bitset[idx]
			if b.IsTimeout() {
				return false, skip
			}
			postpone := b.IsTimeout() || nodeData.RequestedPower == 0
			return postpone, uint32(idx)
		})

	if hasEntries {
		return agg.FinishSequence().AsDigestHolder()
	}
	return nil
}

func (p *NodeVectorHelper) buildGlobulaStateHashes() (gcp_types.GlobulaAnnouncementHash, gcp_types.GlobulaAnnouncementHash) {
	/*
		NB! SequenceDigester requires at least one hash to be added. So to avoid errors, local node MUST always
		have trust level set high enough to get bitset[i].IsTrusted() == true
	*/

	aggTrusted := p.digestFactory.GetGshDigester()
	var aggDoubted cryptography_containers.SequenceDigester

	const skip = math.MaxUint32

	p.entryScanner.ScanSortedWithFilter(
		func(nodeData VectorEntryData, filter uint32) {
			if filter == skip {
				return
			}
			b := p.bitset[filter]
			digest := cryptography_containers.NewDigest(nodeData.AnnounceSignature, "").AsDigestHolder()

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
			b := p.bitset[idx]
			if b.IsTimeout() {
				return false, skip
			}
			postpone := b.IsTimeout() || nodeData.RequestedPower == 0
			return postpone, uint32(idx)
		})

	trustedResult := aggTrusted.FinishSequence().AsDigestHolder()
	if aggDoubted != nil {
		return trustedResult, aggDoubted.FinishSequence().AsDigestHolder()
	}
	return trustedResult, trustedResult
}

func (p *NodeVectorHelper) BuildGlobulaAnnouncementHashes(buildTrusted, buildDoubted bool,
	defaultTrusted, defaultDoubted gcp_types.GlobulaAnnouncementHash) (trustedHash, doubtedHash gcp_types.GlobulaAnnouncementHash) {

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
	defaultTrusted, defaultDoubted gcp_types.GlobulaStateHash) (trustedHash, doubtedHash gcp_types.GlobulaStateHash) {

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

func (p *NodeVectorHelper) VerifyGlobulaStateSignature(localHash gcp_types.GlobulaStateHash, remoteSignature cryptography_containers.SignatureHolder) bool {
	if p.signatureVerifier == nil {
		panic("illegal state - helper must be initialized as a derived one")
	}
	return localHash != nil && p.signatureVerifier.IsValidDigestSignature(localHash, remoteSignature)
}

func (p *NodeVectorHelper) GetNodeBitset() gcp_types.NodeBitset {
	return p.bitset
}

// type vectorSubFilter uint8
//
// func (v vectorSubFilter) CanApply(filter NodeBitsetEntry) bool {
//
// }
//
// type vectorBuilder struct {
//	empty common2.SequenceDigester
//	inclusive bool
//	filters [2]vectorSubFilter
//	vectors [2]common2.SequenceDigester
//	extractor func(nodeData VectorEntryData) common2.DigestHolder
//
// }
//
// func (p *vectorBuilder) applyExclusive(nodeData VectorEntryData, filter NodeBitsetEntry) {
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
// }
//
// func (p *NodeVectorHelper) buildHashes() (common.GlobulaAnnouncementHash, common.GlobulaAnnouncementHash) {
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
// }
