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
	"strings"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/statevector"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/core"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/nodeset"
)

type VectorInspectorFactory interface {
	CreateInspector(scanner nodeset.VectorEntryScanner, digestFactory transport.ConsensusDigestFactory,
		nodeID insolar.ShortNodeID) VectorInspector
}

type VectorInspector interface {
	CreateVector(cryptkit.DigestSigner) statevector.Vector
	InspectVector(sender *core.NodeAppearance, otherData statevector.Vector) InspectedVector
	GetBitset() member.StateBitset
}

type InspectedVector interface {
	IsInspected() bool
	HasMissingMembers() bool
	HasSenderFault() bool
	Reinspect(inspector VectorInspector) InspectedVector

	Inspect()
	GetInspectionResults() (*nodeset.ConsensusStatRow, nodeset.NodeVerificationResult)

	GetBitset() member.StateBitset
	GetNode() *core.NodeAppearance
}

func NewVectorInspectionFactory(maxPopulationForInlineHashing int) VectorInspectorFactory {
	return &VectorInspectionFactory{maxPopulationForInlineHashing}
}

type VectorInspectionFactory struct {
	maxPopulationForInlineHashing int
}

func (p *VectorInspectionFactory) CreateInspector(scanner nodeset.VectorEntryScanner, digestFactory transport.ConsensusDigestFactory,
	nodeID insolar.ShortNodeID) VectorInspector {

	r := &vectorInspectorImpl{VectorBuilder: VectorBuilder{entryScanner: scanner, digestFactory: digestFactory,
		bitset: make(member.StateBitset, scanner.GetIndexedCount())},
		maxPopulationForInlineHashing: p.maxPopulationForInlineHashing,
		nodeID:                        nodeID}

	scanner.ScanIndexed(func(idx int, nodeData nodeset.VectorEntryData) {
		r.bitset[idx] = mapVectorEntryDataToNodesetEntry(nodeData)
	})

	return r
}

type vectorInspectorImpl struct {
	nodeID                        insolar.ShortNodeID
	maxPopulationForInlineHashing int
	VectorBuilder
	Trusted statevector.CalcSubVector
	Doubted statevector.CalcSubVector
}

func (p *vectorInspectorImpl) ensureHashes() {
	if p.Trusted.AnnouncementHash != nil {
		return
	}
	p.Trusted.AnnouncementHash, p.Doubted.AnnouncementHash = p.buildGlobulaAnnouncementHashes()

	p.Trusted.CalcStateWithRank, p.Doubted.CalcStateWithRank =
		p.BuildGlobulaStateHashesAndRanks(true, p.Doubted.AnnouncementHash != nil, p.nodeID,
			statevector.CalcStateWithRank{}, statevector.CalcStateWithRank{})
}

func (p *vectorInspectorImpl) CreateVector(signer cryptkit.DigestSigner) statevector.Vector {
	p.ensureHashes()
	return statevector.NewVector(p.bitset, p.Trusted.Sign(signer), p.Doubted.Sign(signer))
}

func (p *vectorInspectorImpl) InspectVector(sender *core.NodeAppearance, otherVector statevector.Vector) InspectedVector {

	if p.bitset.Len() != otherVector.Bitset.Len() {
		panic("illegal state - StateBitset length mismatch")
	}
	p.ensureHashes()

	r := newInspectedVectorAndPreInspect(p, sender, otherVector)
	return &r
}

func (p *vectorInspectorImpl) GetBitset() member.StateBitset {
	return p.bitset
}

func newInspectedVectorAndPreInspect(p *vectorInspectorImpl, sender *core.NodeAppearance,
	otherVector statevector.Vector) inspectedVector {

	r := inspectedVector{parent: p, node: sender, otherData: otherVector}
	r.comparedStats = nodeset.CompareToStatRow(p.bitset, otherVector.Bitset)

	r.verifyResult = nodeset.NvrNotVerified

	if r.comparedStats.HasValues(nodeset.ComparedMissingHere) {
		// we can't validate anything without data
		// ...  check for updates or/and send requests
		r.verifyResult |= nodeset.NvrMissingNodes
		return r
	}

	if otherVector.Trusted.AnnouncementHash == nil && otherVector.Doubted.AnnouncementHash == nil {
		r.verifyResult |= nodeset.NvrSenderFault
		return r
	}

	r.trustedPart, r.doubtedPart = nodeset.PrepareSubVectorsComparison(r.comparedStats,
		otherVector.Trusted.AnnouncementHash != nil,
		otherVector.Doubted.AnnouncementHash != nil)

	if p.maxPopulationForInlineHashing == 0 || p.maxPopulationForInlineHashing >= p.entryScanner.GetIndexedCount() {
		r.Inspect()
	}
	return r
}

type inspectedVector struct {
	parent                   *vectorInspectorImpl
	node                     *core.NodeAppearance
	otherData                statevector.Vector
	trustedPart, doubtedPart nodeset.SubVectorCompared
	verifyResult             nodeset.NodeVerificationResult
	comparedStats            nodeset.ComparedBitsetRow
	nodeStats                *nodeset.ConsensusStatRow
}

func (p *inspectedVector) GetNode() *core.NodeAppearance {
	return p.node
}

func (p *inspectedVector) GetInspectionResults() (*nodeset.ConsensusStatRow, nodeset.NodeVerificationResult) {
	return p.nodeStats, p.verifyResult
}

func (p *inspectedVector) GetBitset() member.StateBitset {
	return p.otherData.Bitset
}

func (p *inspectedVector) String() string {
	switch p.verifyResult {
	case nodeset.NvrNotVerified, nodeset.NvrSenderFault:
		return p.verifyResult.String()
	}

	b := strings.Builder{}
	b.WriteByte('[')
	p.verifyResult.StringPart(&b)
	b.WriteString("]∑")
	b.WriteString(p.comparedStats.StringSummary())

	return b.String()
}

func (p *inspectedVector) verifySignature(localHash proofs.GlobulaStateHash, remoteSignature cryptkit.SignatureHolder) bool {
	return localHash != nil && p.node.GetSignatureVerifier().IsValidDigestSignature(localHash, remoteSignature)
}

func (p *inspectedVector) IsInspected() bool {
	return p.verifyResult != nodeset.NvrNotVerified
}

func (p *inspectedVector) HasMissingMembers() bool {
	return p.verifyResult&nodeset.NvrMissingNodes != 0
}

func (p *inspectedVector) HasSenderFault() bool {
	return p.verifyResult&nodeset.NvrSenderFault != 0
}

func (p *inspectedVector) Reinspect(inspector VectorInspector) InspectedVector {
	if newParent, ok := inspector.(*vectorInspectorImpl); ok && p.parent == newParent {
		return p
	}
	return inspector.InspectVector(p.node, p.otherData)
}

func (p *inspectedVector) Inspect() {

	if p.IsInspected() {
		return
	}
	if p.nodeStats != nil {
		panic("illegal state")
	}

	p.verifyResult = p.doVerifyVectorHashes()
	if p.verifyResult == nodeset.NvrNotVerified {
		panic("illegal state")
	}

	vr := p.verifyResult &^ nodeset.NvrHashlessFlags
	if p.verifyResult == nodeset.NvrNotVerified {
		panic("illegal state")
	}

	ns := nodeset.SummarizeStats(p.otherData.Bitset, vr, p.comparedStats)
	p.nodeStats = &ns
}

func (p *inspectedVector) doVerifyVectorHashes() nodeset.NodeVerificationResult {

	selfData := statevector.CalcVector{Trusted: p.parent.Trusted, Doubted: p.parent.Doubted}

	if p.doubtedPart.IsNeeded() && selfData.Doubted.AnnouncementHash == nil {
		// special case when all our nodes are in trusted, so other's doubted vector will be matched with the trusted one of ours
		selfData.Doubted = selfData.Trusted
	}

	gahTrusted, gahDoubted := selfData.Trusted.AnnouncementHash, selfData.Doubted.AnnouncementHash

	var vectorBuilder VectorBuilder //
	if p.trustedPart.IsRecalc() || p.doubtedPart.IsRecalc() {
		vectorBuilder = p.parent.CreateDerived(p.comparedStats)

		// It does remap the original bitset with the given stats
		gahTrusted, gahDoubted = vectorBuilder.BuildGlobulaAnnouncementHashes(
			p.trustedPart.IsRecalc(), p.doubtedPart.IsRecalc(), gahTrusted, gahDoubted)
	}

	validTrusted := p.trustedPart.IsNeeded() && gahTrusted.Equals(p.otherData.Trusted.AnnouncementHash)
	validDoubted := p.doubtedPart.IsNeeded() && gahDoubted.Equals(p.otherData.Doubted.AnnouncementHash)

	verifyRes := nodeset.NvrNotVerified
	if validDoubted && !validTrusted {
		// As Trusted is a subset of Doubted, then Doubted can't be valid if Trusted is not.
		// This is an evident fraud/error by the sender.
		// Use status for doubted, but ignore results for Trusted check
		// TODO report fraud
		verifyRes |= nodeset.NvrSenderFault
		p.trustedPart = nodeset.SvcIgnore
	}

	if validTrusted || validDoubted {
		recalcTrusted := p.trustedPart.IsRecalc() && validTrusted
		recalcDoubted := p.doubtedPart.IsRecalc() && validDoubted

		gshTrusted, gshDoubted := selfData.Trusted.CalcStateWithRank, selfData.Doubted.CalcStateWithRank
		if recalcTrusted || recalcDoubted {
			gshTrusted, gshDoubted = vectorBuilder.BuildGlobulaStateHashesAndRanks(recalcTrusted, recalcDoubted,
				p.node.GetNodeID(), gshTrusted, gshDoubted)

			if recalcTrusted {
				validTrusted = gshTrusted.ExpectedRank == p.otherData.Trusted.ExpectedRank
			}
			if recalcDoubted {
				validDoubted = gshDoubted.ExpectedRank == p.otherData.Doubted.ExpectedRank
			}
		}

		validTrusted = validTrusted && p.verifySignature(gshTrusted.StateHash, p.otherData.Trusted.StateSignature)
		validDoubted = validDoubted && p.verifySignature(gshDoubted.StateHash, p.otherData.Doubted.StateSignature)
	}

	if p.trustedPart.IsNeeded() {
		verifyRes.SetTrusted(validTrusted, p.trustedPart.IsRecalc())
	}
	if p.doubtedPart.IsNeeded() {
		verifyRes.SetDoubted(validDoubted, p.doubtedPart.IsRecalc())
	}

	return verifyRes
}

func mapVectorEntryDataToNodesetEntry(nodeData nodeset.VectorEntryData) member.BitsetEntry {
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
