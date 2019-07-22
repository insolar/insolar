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
	"strings"

	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/statevector"
)

type LocalHashedNodeVector struct {
	statevector.Vector
	TrustedGlobulaStateVector proofs.GlobulaStateHash
	DoubtedGlobulaStateVector proofs.GlobulaStateHash
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
		// Trusted is always present as there is always at least one node - the sender
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

func SummarizeStatsOld(hereDataBitset, otherDataBitset member.StateBitset, verifyRes NodeVerificationResult, sr ComparedBitsetRow) ConsensusStatRow {

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
			nodeResult = ConsensusStatMissingHere // TODO wrong
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
				if hereDataBitset[i].IsFraud() {
					nodeResult = ConsensusStatFraud
				} else {
					nodeResult = ConsensusStatFraudSuspect
				}
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

func SummarizeStats(otherDataBitset member.StateBitset, verifyRes NodeVerificationResult, sr ComparedBitsetRow) ConsensusStatRow {

	nodeStats := NewConsensusStatRow(sr.ColumnCount())

	for i := 0; i < sr.ColumnCount(); i++ {
		nodeResult, fraudEnforcementCheck := summaryByEntry(otherDataBitset[i], verifyRes, sr.Get(i))

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

func summaryByEntry(otherEntry member.BitsetEntry, verifyRes NodeVerificationResult,
	compared ComparedState) (ConsensusStat, NodeVerificationResult) {

	if otherEntry.IsFraud() {
		return ConsensusStatFraud, NvrNotVerified
	}

	switch compared {
	case ComparedMissingHere:
		// missing here and present there in "doubted"
		// so we can't build GSH without it anyway
		if verifyRes != NvrNotVerified {
			panic("unexpected")
		}
	case ComparedDoubtedMissingHere:
		// missed by us, so we can't build doubted GSH without it anyway
		switch {
		case verifyRes.AnyOf(NvrDoubtedFraud | NvrDoubtedValid):
			panic("unexpected")
		case otherEntry.IsTrusted():
			switch {
			case verifyRes.AnyOf(NvrTrustedValid):
				return ConsensusStatTrusted, NvrNotVerified
			case verifyRes.AllOf(NvrTrustedFraud):
				return ConsensusStatFraudSuspect, NvrTrustedFraud
			}
		}
	case ComparedMissingThere:
		return ConsensusStatMissingThere, NvrNotVerified
	case ComparedSame:
		switch {
		case otherEntry.IsTimeout():
			// it was missing on both sides
			return ConsensusStatMissingThere, NvrNotVerified
		case otherEntry.IsFraud():
			return ConsensusStatFraud, NvrNotVerified
		case otherEntry.IsTrusted() && verifyRes.AnyOf(NvrTrustedFraud):
			return ConsensusStatFraudSuspect, NvrTrustedFraud
		case otherEntry.IsTrusted() && verifyRes.AnyOf(NvrTrustedValid):
			return ConsensusStatTrusted, NvrNotVerified
		case verifyRes.AnyOf(NvrDoubtedValid):
			return ConsensusStatDoubted, NvrNotVerified
		case verifyRes.AnyOf(NvrDoubtedFraud):
			return ConsensusStatFraudSuspect, NvrDoubtedFraud
		}
	case ComparedLessTrustedThere:
		switch {
		case otherEntry.IsFraud():
			return ConsensusStatFraudSuspect, NvrNotVerified
		case verifyRes.AnyOf(NvrDoubtedValid):
			return ConsensusStatDoubted, NvrNotVerified
		case verifyRes.AnyOf(NvrDoubtedFraud):
			return ConsensusStatFraudSuspect, NvrDoubtedFraud
		}
	case ComparedLessTrustedHere:
		switch {
		case verifyRes.AllOf(NvrTrustedValid | NvrTrustedAlteredNodeSet):
			return ConsensusStatTrusted, NvrNotVerified
		case verifyRes.AllOf(NvrTrustedFraud | NvrTrustedAlteredNodeSet):
			return ConsensusStatFraudSuspect, NvrTrustedFraud
		case verifyRes.AnyOf(NvrDoubtedValid):
			return ConsensusStatDoubted, NvrNotVerified
		case verifyRes.AnyOf(NvrDoubtedFraud):
			return ConsensusStatFraudSuspect, NvrDoubtedFraud
		}
	default:
		panic("unexpected")
	}
	return ConsensusStatMissingHere, NvrNotVerified
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
	if v == NvrNotVerified {
		return "[unverified]"
	}

	b := strings.Builder{}
	b.WriteByte('[')
	v.StringPart(&b)
	b.WriteByte(']')
	return b.String()
}
