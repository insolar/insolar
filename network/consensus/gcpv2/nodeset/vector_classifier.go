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
	"github.com/insolar/insolar/network/consensus/gcpv2/gcp_types"
	"github.com/insolar/insolar/network/consensus/gcpv2/stats"
)

type LocalHashedNodeVector struct {
	gcp_types.HashedNodeVector
	TrustedGlobulaStateVector gcp_types.GlobulaStateHash
	DoubtedGlobulaStateVector gcp_types.GlobulaStateHash
}

func ClassifyByNodeGsh(selfData LocalHashedNodeVector, otherData gcp_types.HashedNodeVector, nodeStats *stats.Row, derivedVector *NodeVectorHelper) NodeVerificationResult {

	if selfData.Bitset.Len() != nodeStats.ColumnCount() {
		panic("bitset length mismatch")
	}

	sr := CompareToStatRow(selfData.Bitset, otherData.Bitset)
	verifyRes := verifyVectorHashes(selfData, otherData, sr, derivedVector)

	if verifyRes == norNotVerified || verifyRes == NvrSenderFault {
		return verifyRes
	}

	summarize(otherData.Bitset, verifyRes&^NvrSenderFault, sr, nodeStats)
	return verifyRes
}

type subVectorVerifyMode uint8

const (
	ignore subVectorVerifyMode = iota
	verify
	verifyAsIs
	verifyRecalc
)

func (v subVectorVerifyMode) UpdateVerify(n subVectorVerifyMode) subVectorVerifyMode {
	if v == ignore {
		return ignore
	}
	if v != verify {
		panic("illegal state")
	}
	return n
}

func (v subVectorVerifyMode) IsNeeded() bool {
	return v != ignore
}

func initVerify(needed bool) subVectorVerifyMode {
	if needed {
		return verify
	}
	return ignore
}

func verifyVectorHashes(selfData LocalHashedNodeVector, otherData gcp_types.HashedNodeVector, sr *stats.Row, derivedVector *NodeVectorHelper) NodeVerificationResult {
	// TODO All GSH comparisons should be based on SIGNATURES! not on pure hashes

	verifyRes := norNotVerified

	trustedPart := initVerify(otherData.TrustedAnnouncementVector != nil)
	doubtedPart := initVerify(otherData.DoubtedAnnouncementVector != nil)

	if !trustedPart.IsNeeded() {
		// Trusted is always present as there is always at least one node - the sender
		panic("illegal state")
	}
	if doubtedPart.IsNeeded() && selfData.DoubtedAnnouncementVector == nil {
		// special case when all our nodes are in trusted, so other's doubted vector will be matched with the trusted one of ours
		selfData.DoubtedAnnouncementVector = selfData.TrustedAnnouncementVector
		selfData.DoubtedGlobulaStateVector = selfData.TrustedGlobulaStateVector
		// selfData.DoubtedGlobulaStateVectorSignature = selfData.TrustedGlobulaStateVectorSignature
	}

	if sr.HasValues(NodeBitMissingHere) {
		// we can't validate anything without data
		// ...  check for updates or/and send requests
		return norNotVerified
	}

	switch {
	case sr.HasAllValuesOf(NodeBitSame, NodeBitLessTrustedHere):
		// check DoubtedGsh as is, if not then TrustedGSH with some locally-known NSH included
		fallthrough
	case sr.HasAllValuesOf(NodeBitSame, NodeBitLessTrustedThere):
		// check DoubtedGsh as is, if not then TrustedGSH with some locally-known NSH excluded
		if sr.HasAllValues(NodeBitSame) {
			trustedPart = trustedPart.UpdateVerify(verifyAsIs)
		} else {
			trustedPart = trustedPart.UpdateVerify(verifyRecalc)
		}
		doubtedPart = doubtedPart.UpdateVerify(verifyAsIs)
	case sr.HasValues(NodeBitDoubtedMissingHere):
		// check TrustedGSH only
		// validation of DoubtedGsh needs requests
		// ...  check for updates and send requests

		// if HasValues(NodeBitLessTrustedThere) then ... exclude some locally-known NSH
		// if HasValues(NodeBitMissingThere) then ... exclude some locally-known NSH
		// if HasValues(NodeBitLessTrustedHere) then ... include some locally-known NSH
		trustedPart = trustedPart.UpdateVerify(verifyRecalc)
		doubtedPart = doubtedPart.UpdateVerify(ignore)
	default:
		// if HasValues(NodeBitLessTrustedThere) then ... exclude some locally-known NSH from TrustedGSH
		// if HasValues(NodeBitLessTrustedHere) then ... include some locally-known NSH to TrustedGSH

		trustedPart = trustedPart.UpdateVerify(verifyRecalc)
		if sr.HasValues(NodeBitMissingThere) {
			// check DoubtedGsh with exclusions, then TrustedGSH with exclusions/inclusions
			doubtedPart = doubtedPart.UpdateVerify(verifyRecalc)
		} else {
			// check DoubtedGsh as-is, then TrustedGSH with exclusions/inclusions
			doubtedPart = doubtedPart.UpdateVerify(verifyAsIs)
		}
	}

	gahTrusted, gahDoubted := selfData.TrustedAnnouncementVector, selfData.DoubtedAnnouncementVector

	switch {
	case trustedPart == verify || doubtedPart == verify:
		panic("illegal state")
	case trustedPart == verifyRecalc || doubtedPart == verifyRecalc:

		// It does remap the original bitset with the given stats
		derivedVector.PrepareDerivedVector(sr)

		gahTrusted, gahDoubted = derivedVector.BuildGlobulaAnnouncementHashes(
			trustedPart == verifyRecalc, doubtedPart == verifyRecalc, gahTrusted, gahDoubted)
	}

	validTrusted := trustedPart.IsNeeded() && gahTrusted.Equals(otherData.TrustedAnnouncementVector)
	validDoubted := doubtedPart.IsNeeded() && gahDoubted.Equals(otherData.DoubtedAnnouncementVector)

	if validDoubted && !validTrusted {
		// As Trusted is a subset of Doubted, then Doubted can't be valid if Trusted is not.
		// This is an evident fraud/error by the sender.
		// Use status for doubted, but ignore results for Trusted check
		// TODO report fraud
		verifyRes |= NvrSenderFault
		trustedPart = ignore
	}

	if validTrusted || validDoubted {
		recalcTrusted := trustedPart == verifyRecalc && validTrusted
		recalcDoubted := doubtedPart == verifyRecalc && validDoubted

		gshTrusted, gshDoubted := selfData.TrustedGlobulaStateVector, selfData.DoubtedGlobulaStateVector
		if recalcTrusted || recalcDoubted {
			gshTrusted, gshDoubted = derivedVector.BuildGlobulaStateHashes(
				recalcTrusted, recalcDoubted, gshTrusted, gshDoubted)
		}

		validTrusted = validTrusted && derivedVector.VerifyGlobulaStateSignature(gshTrusted, otherData.TrustedGlobulaStateVectorSignature)
		validDoubted = validDoubted && derivedVector.VerifyGlobulaStateSignature(gshDoubted, otherData.DoubtedGlobulaStateVectorSignature)
	}

	if doubtedPart.IsNeeded() {
		verifyRes.SetDoubted(validDoubted, doubtedPart == verifyRecalc)
	}
	if trustedPart.IsNeeded() {
		verifyRes.SetTrusted(validTrusted, trustedPart == verifyRecalc)
	}

	return verifyRes
}

func summarize(otherDataBitset gcp_types.NodeBitset, verifyRes NodeVerificationResult, sr *stats.Row, nodeStats *stats.Row) {

	for i := 0; i < sr.ColumnCount(); i++ {
		nodeResult := ConsensusStatMissingHere
		fraudEnforcementCheck := norNotVerified

		switch sr.Get(i) {
		case NodeBitMissingHere:
			// missing here and present there in "doubted"
			// so we can't build GSH without it anyway
			if verifyRes != norNotVerified {
				panic("unexpected")
			}
			nodeResult = ConsensusStatMissingHere
		case NodeBitMissingThere:
			nodeResult = ConsensusStatMissingThere
		case NodeBitDoubtedMissingHere:
			// missed by us, so we can't build doubted GSH without it anyway
			if verifyRes.AnyOf(NvrDoubtedFraud | NvrDoubtedValid) {
				panic("unexpected")
			}
			nodeResult = ConsensusStatMissingHere
		case NodeBitSame:
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
		case NodeBitLessTrustedThere:
			switch {
			case verifyRes.AnyOf(NvrDoubtedValid):
				nodeResult = ConsensusStatDoubted
			case verifyRes.AnyOf(NvrDoubtedFraud):
				fraudEnforcementCheck = NvrDoubtedFraud
				nodeResult = ConsensusStatFraudSuspect
			}
		case NodeBitLessTrustedHere:
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
}

type NodeVerificationResult int

const norNotVerified NodeVerificationResult = 0

const (
	NvrSenderFault NodeVerificationResult = 1 << iota
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
