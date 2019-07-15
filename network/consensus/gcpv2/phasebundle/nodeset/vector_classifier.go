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
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/statevector"
	"strings"
)

type LocalHashedNodeVector struct {
	statevector.Vector
	TrustedGlobulaStateVector proofs.GlobulaStateHash
	DoubtedGlobulaStateVector proofs.GlobulaStateHash
}

func ClassifyByNodeGsh(selfData LocalHashedNodeVector, otherData statevector.Vector,
	derivedVector *NodeVectorHelper) (NodeVerificationResult, ConsensusStatRow) {

	sr := CompareToStatRow(selfData.Bitset, otherData.Bitset)

	if sr.HasValues(ComparedMissingHere) {
		// we can't validate anything without data
		// ...  check for updates or/and send requests
		return NvrMissingNodes, ConsensusStatRow{}
	}

	trustedPart, doubtedPart := PrepareSubVectorsComparison(sr,
		otherData.Trusted.AnnouncementHash != nil,
		otherData.Doubted.AnnouncementHash != nil)

	if trustedPart == verifyRecalc || doubtedPart == verifyRecalc {
		//It does remap the original bitset with the given stats
		derivedVector.PrepareDerivedVector(sr)
	}

	verifyRes := doVerifyVectorHashes(trustedPart, doubtedPart, selfData, otherData, derivedVector)

	if verifyRes == NvrNotVerified || verifyRes == NvrSenderFault {
		return verifyRes, ConsensusStatRow{}
	}

	nodeStats := SummarizeStats(otherData.Bitset, verifyRes&^NvrHashlessFlags, sr)
	return verifyRes, nodeStats
}

type SubVectorCompared uint8

const (
	SvcIgnore SubVectorCompared = iota
	verify
	verifyAsIs
	verifyRecalc
)

func (v SubVectorCompared) UpdateVerify(n SubVectorCompared) SubVectorCompared {
	if v == SvcIgnore {
		return SvcIgnore
	}
	if v != verify {
		panic("illegal state")
	}
	return n
}

func (v SubVectorCompared) IsNeeded() bool {
	return v != SvcIgnore
}

func (v SubVectorCompared) IsRecalc() bool {
	return v == verifyRecalc
}

func initVerify(needed bool) SubVectorCompared {
	if needed {
		return verify
	}
	return SvcIgnore
}

func PrepareSubVectorsComparison(sr ComparedBitsetRow, hasOtherTrusted, hasOtherDoubted bool) (SubVectorCompared, SubVectorCompared) {
	trustedPart := initVerify(hasOtherTrusted)
	doubtedPart := initVerify(hasOtherDoubted)

	if !trustedPart.IsNeeded() {
		//Trusted is always present as there is always at least one node - the sender
		panic("illegal state")
	}

	if sr.HasValues(ComparedMissingHere) {
		// we can't validate anything without data
		// ...  check for updates or/and send requests
		return SvcIgnore, SvcIgnore
	}

	switch {
	case sr.HasAllValuesOf(ComparedSame, ComparedLessTrustedHere):
		// check DoubtedGsh as is, if not then TrustedGSH with some locally-known NSH included
		fallthrough
	case sr.HasAllValuesOf(ComparedSame, ComparedLessTrustedThere):
		// check DoubtedGsh as is, if not then TrustedGSH with some locally-known NSH excluded
		if sr.HasAllValues(ComparedSame) {
			trustedPart = trustedPart.UpdateVerify(verifyAsIs)
		} else {
			trustedPart = trustedPart.UpdateVerify(verifyRecalc)
		}
		doubtedPart = doubtedPart.UpdateVerify(verifyAsIs)
	case sr.HasValues(ComparedDoubtedMissingHere):
		// check TrustedGSH only
		// validation of DoubtedGsh needs requests
		// ...  check for updates and send requests

		// if HasValues(ComparedLessTrustedThere) then ... exclude some locally-known NSH
		// if HasValues(ComparedMissingThere) then ... exclude some locally-known NSH
		// if HasValues(ComparedLessTrustedHere) then ... include some locally-known NSH
		trustedPart = trustedPart.UpdateVerify(verifyRecalc)
		doubtedPart = doubtedPart.UpdateVerify(SvcIgnore)
	default:
		// if HasValues(ComparedLessTrustedThere) then ... exclude some locally-known NSH from TrustedGSH
		// if HasValues(ComparedLessTrustedHere) then ... include some locally-known NSH to TrustedGSH

		trustedPart = trustedPart.UpdateVerify(verifyRecalc)
		if sr.HasValues(ComparedMissingThere) {
			// check DoubtedGsh with exclusions, then TrustedGSH with exclusions/inclusions
			doubtedPart = doubtedPart.UpdateVerify(verifyRecalc)
		} else {
			// check DoubtedGsh as-is, then TrustedGSH with exclusions/inclusions
			doubtedPart = doubtedPart.UpdateVerify(verifyAsIs)
		}
	}

	if trustedPart == verify || doubtedPart == verify {
		panic("illegal state")
	}

	return trustedPart, doubtedPart
}

func doVerifyVectorHashes(trustedPart, doubtedPart SubVectorCompared,
	selfData LocalHashedNodeVector, otherData statevector.Vector, derivedVector *NodeVectorHelper) NodeVerificationResult {

	if doubtedPart.IsNeeded() && selfData.Doubted.AnnouncementHash == nil {
		//special case when all our nodes are in trusted, so other's doubted vector will be matched with the trusted one of ours
		selfData.Doubted.AnnouncementHash = selfData.Trusted.AnnouncementHash
		selfData.DoubtedGlobulaStateVector = selfData.TrustedGlobulaStateVector
		//selfData.DoubtedGlobulaStateVectorSignature = selfData.TrustedGlobulaStateVectorSignature
	}

	gahTrusted, gahDoubted := selfData.Trusted.AnnouncementHash, selfData.Doubted.AnnouncementHash

	if trustedPart.IsRecalc() || doubtedPart.IsRecalc() {
		//It does remap the original bitset with the given stats
		//derivedVector.PrepareDerivedVector(sr)
		gahTrusted, gahDoubted = derivedVector.BuildGlobulaAnnouncementHashes(
			trustedPart.IsRecalc(), doubtedPart.IsRecalc(), gahTrusted, gahDoubted)
	}

	validTrusted := trustedPart.IsNeeded() && gahTrusted.Equals(otherData.Trusted.AnnouncementHash)
	validDoubted := doubtedPart.IsNeeded() && gahDoubted.Equals(otherData.Doubted.AnnouncementHash)

	verifyRes := NvrNotVerified
	if validDoubted && !validTrusted {
		// As Trusted is a subset of Doubted, then Doubted can't be valid if Trusted is not.
		// This is an evident fraud/error by the sender.
		// Use status for doubted, but SvcIgnore results for Trusted check
		// TODO report fraud
		verifyRes |= NvrSenderFault
		trustedPart = SvcIgnore
	}

	if validTrusted || validDoubted {
		recalcTrusted := trustedPart.IsRecalc() && validTrusted
		recalcDoubted := doubtedPart.IsRecalc() && validDoubted

		gshTrusted, gshDoubted := selfData.TrustedGlobulaStateVector, selfData.DoubtedGlobulaStateVector
		if recalcTrusted || recalcDoubted {
			gshTrusted, gshDoubted = derivedVector.BuildGlobulaStateHashes(
				recalcTrusted, recalcDoubted, gshTrusted, gshDoubted)
		}

		validTrusted = validTrusted && derivedVector.VerifyGlobulaStateSignature(gshTrusted, otherData.Trusted.StateSignature)
		validDoubted = validDoubted && derivedVector.VerifyGlobulaStateSignature(gshDoubted, otherData.Doubted.StateSignature)
	}

	if trustedPart.IsNeeded() {
		verifyRes.SetTrusted(validTrusted, trustedPart.IsRecalc())
	}
	if doubtedPart.IsNeeded() {
		verifyRes.SetDoubted(validDoubted, doubtedPart.IsRecalc())
	}

	return verifyRes
}

func SummarizeStats(otherDataBitset member.StateBitset, verifyRes NodeVerificationResult, sr ComparedBitsetRow) ConsensusStatRow {

	nodeStats := NewConsensusStatRow(sr.ColumnCount())

	for i := 0; i < sr.ColumnCount(); i++ {
		nodeResult := ConsensusStatMissingHere
		fraudEnforcementCheck := NvrNotVerified

		switch sr.Get(i) {
		case ComparedMissingHere:
			// missing here and present there in "doubted"
			// so we can't build GSH without it anyway
			if verifyRes != NvrNotVerified {
				panic("unexpected")
			}
			nodeResult = ConsensusStatMissingHere
		case ComparedMissingThere:
			nodeResult = ConsensusStatMissingThere
		case ComparedDoubtedMissingHere:
			// missed by us, so we can't build doubted GSH without it anyway
			if verifyRes.AnyOf(NvrDoubtedFraud | NvrDoubtedValid) {
				panic("unexpected")
			}
			nodeResult = ConsensusStatMissingHere
		case ComparedSame:
			b := otherDataBitset[i]
			switch {
			case b.IsTimeout():
				// it was missing on both sides
				nodeResult = ConsensusStatMissingThere
			case b.IsFraud():
				// we don't need checks to agree on fraud mutually detected
				nodeResult = ConsensusStatFraud
			case b.IsTrusted() && verifyRes.AnyOf(NvrTrustedFraud):
				nodeResult = ConsensusStatFraudSuspect
				fraudEnforcementCheck = NvrTrustedFraud
			case b.IsTrusted() && verifyRes.AnyOf(NvrTrustedValid):
				nodeResult = ConsensusStatTrusted
			case verifyRes.AnyOf(NvrDoubtedValid):
				nodeResult = ConsensusStatDoubted
			case verifyRes.AnyOf(NvrDoubtedFraud):
				fraudEnforcementCheck = NvrDoubtedFraud
				nodeResult = ConsensusStatFraudSuspect
			}
		case ComparedLessTrustedThere:
			switch {
			case verifyRes.AnyOf(NvrDoubtedValid):
				nodeResult = ConsensusStatDoubted
			case verifyRes.AnyOf(NvrDoubtedFraud):
				fraudEnforcementCheck = NvrDoubtedFraud
				nodeResult = ConsensusStatFraudSuspect
			}
		case ComparedLessTrustedHere:
			switch {
			case verifyRes.AllOf(NvrTrustedValid | NvrTrustedAlteredNodeSet):
				nodeResult = ConsensusStatTrusted
			case verifyRes.AllOf(NvrTrustedFraud | NvrTrustedAlteredNodeSet):
				fraudEnforcementCheck = NvrTrustedFraud
				nodeResult = ConsensusStatFraudSuspect
			case verifyRes.AnyOf(NvrDoubtedValid):
				nodeResult = ConsensusStatDoubted
			case verifyRes.AnyOf(NvrDoubtedFraud):
				fraudEnforcementCheck = NvrDoubtedFraud
				nodeResult = ConsensusStatFraudSuspect
			}
		default:
			panic("unexpected")
		}
		if nodeResult == ConsensusStatFraudSuspect {
			switch fraudEnforcementCheck {
			case NvrTrustedFraud:
				// TODO check if there is the only one trusted, then set ConsensusStatFraud
			case NvrDoubtedFraud:
				// TODO check if it is the only one in doubted different from trusted, then set ConsensusStatFraud
			}
		}
		nodeStats.Set(i, nodeResult)
	}
	return nodeStats
}

type NodeVerificationResult uint16

const NvrNotVerified NodeVerificationResult = 0
const NvrHashlessFlags = NvrSenderFault | NvrMissingNodes

const (
	NvrSenderFault NodeVerificationResult = 1 << iota
	NvrMissingNodes
	NvrTrustedValid
	NvrTrustedFraud
	NvrTrustedAlteredNodeSet
	NvrDoubtedValid
	NvrDoubtedFraud
	NvrDoubtedAlteredNodeSet
)

func (v NodeVerificationResult) AnyOf(f NodeVerificationResult) bool {
	return v&f != 0
}

func (v NodeVerificationResult) AllOf(f NodeVerificationResult) bool {
	return v&f == f
}

func (v *NodeVerificationResult) setOnce(f NodeVerificationResult, valid bool, altered bool) {
	if *v&(f|(f<<1)) != 0 {
		panic("repeated set")
	}
	if altered {
		*v |= f << 2
	}
	if !valid {
		f <<= 1
	}
	*v |= f
}

func (v *NodeVerificationResult) SetTrusted(valid bool, altered bool) bool {
	v.setOnce(NvrTrustedValid, valid, altered)
	return valid
}

func (v *NodeVerificationResult) SetDoubted(valid bool, altered bool) bool {
	v.setOnce(NvrDoubtedValid, valid, altered)
	return valid
}

func (v NodeVerificationResult) StringPart(b *strings.Builder) {

	b.WriteString("verified")
	if v.AnyOf(NvrSenderFault) {
		b.WriteString(" sender-fault")
	}
	if v.AnyOf(NvrMissingNodes) {
		b.WriteString(" missing")
	}
	if v.AnyOf(NvrTrustedValid | NvrTrustedFraud) {
		b.WriteByte(' ')
		if v.AnyOf(NvrTrustedAlteredNodeSet) {
			b.WriteRune('≈')
		}
		if v.AnyOf(NvrTrustedFraud) {
			b.WriteByte('!')
		}
		b.WriteByte('T')
	}
	if v.AnyOf(NvrDoubtedValid | NvrDoubtedFraud) {
		b.WriteByte(' ')
		if v.AnyOf(NvrDoubtedAlteredNodeSet) {
			b.WriteRune('≈')
		}
		if v.AnyOf(NvrDoubtedFraud) {
			b.WriteByte('!')
		}
		b.WriteByte('D')
	}
}

func (v NodeVerificationResult) String() string {
	switch v {
	case NvrNotVerified:
		return "[unverified]"
	}

	b := strings.Builder{}
	b.WriteByte('[')
	v.StringPart(&b)
	b.WriteByte(']')
	return b.String()
}
