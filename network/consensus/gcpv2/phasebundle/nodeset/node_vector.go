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
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
)

type NodeVectorHelper struct {
	digestFactory     transport.ConsensusDigestFactory
	signatureVerifier cryptkit.SignatureVerifier
	entryScanner      VectorEntryScanner
	bitset            member.StateBitset
	parentBitset      member.StateBitset
}

func NewLocalNodeVector(digestFactory transport.ConsensusDigestFactory,
	entryScanner VectorEntryScanner) NodeVectorHelper {

	p := NodeVectorHelper{digestFactory, nil,
		entryScanner, make(member.StateBitset, entryScanner.GetIndexedCount()), nil,
	}

	entryScanner.ScanIndexed(func(idx int, nodeData VectorEntryData) {
		p.bitset[idx] = mapVectorEntryDataToNodesetEntry(nodeData)
	})
	return p
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

func (p *NodeVectorHelper) CreateDerivedVector(signatureVerifier cryptkit.SignatureVerifier) NodeVectorHelper {
	return NodeVectorHelper{p.digestFactory, signatureVerifier,
		p.entryScanner, nil, p.bitset}
}

func (p *NodeVectorHelper) PrepareDerivedVector(statRow ComparedBitsetRow) {
	if p.bitset != nil && p.parentBitset == nil {
		panic("illegal state")
	}

	p.bitset = make(member.StateBitset, len(p.parentBitset))

	for idx := range p.parentBitset {
		switch statRow.Get(idx) {
		case ComparedMissingHere:
			p.bitset[idx] = member.BeTimeout // we don't have it
		case ComparedDoubtedMissingHere:
			p.bitset[idx] = member.BeTimeout // we don't have it
		case ComparedSame:
			// ok, use as-is
			p.bitset[idx] = p.parentBitset[idx]
		case ComparedLessTrustedThere:
			// ok - exclude for trusted
			p.bitset[idx] = member.BeBaselineTrust
		case ComparedLessTrustedHere:
			// ok - use for both
			p.bitset[idx] = member.BeHighTrust
		case ComparedMissingThere:
			p.bitset[idx] = member.BeTimeout // we have it, but the other's doesn't
		default:
			panic("unexpected")
		}
	}
}

func (p *NodeVectorHelper) buildGlobulaAnnouncementHash(trusted bool) proofs.GlobulaAnnouncementHash {
	hasEntries := false
	agg := p.digestFactory.GetAnnouncementDigester()

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

func (p *NodeVectorHelper) buildGlobulaAnnouncementHashes() (proofs.GlobulaAnnouncementHash, proofs.GlobulaAnnouncementHash) {
	/*
		NB! SequenceDigester requires at least one hash to be added. So to avoid errors, local node MUST always
		have trust level set high enough to get bitset[i].IsTrusted() == true
	*/

	aggTrusted := p.digestFactory.GetAnnouncementDigester()
	var aggDoubted cryptkit.SequenceDigester

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

const (
	skipEntry = iota
	//missingEntry
	trustedEntry
	doubtedEntry
)

func (p *NodeVectorHelper) stateFilter(idx int, nodeData VectorEntryData) (bool, uint32) {
	b := p.bitset[idx]
	if b.IsTimeout() {
		return false, skipEntry
	}
	postpone := nodeData.RequestedMode.IsPowerless() || nodeData.RequestedPower == 0
	if b.IsTrusted() {
		return postpone, trustedEntry
	}
	return postpone, doubtedEntry
}

func (p *NodeVectorHelper) buildGlobulaStateHash(trusted bool) proofs.GlobulaAnnouncementHash {
	hasEntries := false
	agg := p.digestFactory.GetSequenceDigester()

	p.entryScanner.ScanSortedWithFilter(
		func(nodeData VectorEntryData, filter uint32) {
			if filter == skipEntry {
				return
			}
			if filter == doubtedEntry && !trusted {
				return
			}
			digest := cryptkit.NewDigest(nodeData.AnnounceSignature, "").AsDigestHolder()

			agg.AddNext(digest)
			hasEntries = true
		}, p.stateFilter)

	if hasEntries {
		return agg.FinishSequence().AsDigestHolder()
	}
	return nil
}

func (p *NodeVectorHelper) buildGlobulaStateHashes() (proofs.GlobulaAnnouncementHash, proofs.GlobulaAnnouncementHash) {
	/*
		NB! SequenceDigester requires at least one hash to be added. So to avoid errors, local node MUST always
		have trust level set high enough to get bitset[i].IsTrusted() == true
	*/

	aggTrusted := p.digestFactory.GetSequenceDigester()
	var aggDoubted cryptkit.SequenceDigester

	p.entryScanner.ScanSortedWithFilter(
		func(nodeData VectorEntryData, filter uint32) {
			if filter == skipEntry {
				return
			}

			digest := cryptkit.NewDigest(nodeData.AnnounceSignature, "").AsDigestHolder()
			switch filter {
			case trustedEntry:
				aggTrusted.AddNext(digest)
				if aggDoubted == nil {
					return
				}
			case doubtedEntry:
				if aggDoubted == nil {
					aggDoubted = aggTrusted.ForkSequence()
				}
			}
			aggDoubted.AddNext(digest)
		}, p.stateFilter)

	trustedResult := aggTrusted.FinishSequence().AsDigestHolder()
	if aggDoubted != nil {
		return trustedResult, aggDoubted.FinishSequence().AsDigestHolder()
	}
	return trustedResult, trustedResult
}

func (p *NodeVectorHelper) BuildGlobulaAnnouncementHashes(buildTrusted, buildDoubted bool,
	defaultTrusted, defaultDoubted proofs.GlobulaAnnouncementHash) (trustedHash, doubtedHash proofs.GlobulaAnnouncementHash) {

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
	defaultTrusted, defaultDoubted proofs.GlobulaStateHash) (trustedHash, doubtedHash proofs.GlobulaStateHash) {

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

func (p *NodeVectorHelper) VerifyGlobulaStateSignature(localHash proofs.GlobulaStateHash, remoteSignature cryptkit.SignatureHolder) bool {
	if p.signatureVerifier == nil {
		panic("illegal state - helper must be initialized as a derived one")
	}
	return localHash != nil && p.signatureVerifier.IsValidDigestSignature(localHash, remoteSignature)
}

func (p *NodeVectorHelper) GetNodeBitset() member.StateBitset {
	return p.bitset
}
