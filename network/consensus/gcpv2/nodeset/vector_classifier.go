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
	"fmt"

	"github.com/insolar/insolar/network/consensus/common"
	"github.com/insolar/insolar/network/consensus/gcpv2/stats"
)

type HashedNodeVector struct {
	Bitset        NodeBitset
	TrustedVector common.DigestHolder
	DoubtedVector common.DigestHolder
}

func ClassifyByNodeGsh(selfData HashedNodeVector, otherData HashedNodeVector, nodeStats *stats.Row, hasher FilteredSequenceHasher) NodeVerificationResult {

	if selfData.DoubtedVector == nil {
		selfData.DoubtedVector = selfData.TrustedVector
	}
	if otherData.DoubtedVector == nil {
		otherData.DoubtedVector = otherData.TrustedVector
	}

	if selfData.Bitset.Len() != nodeStats.ColumnCount() {
		panic("bitset length mismatch")
	}

	// var sr *stats.Row
	// var verifyRes NodeVerificationResult
	// for {
	// 	repeat := false
	// 	verifyRes, repeat = catcher(selfData, otherData, hasher, &sr)
	// 	if repeat {
	// 		continue
	// 	} else {
	// 		break
	// 	}
	// }
	sr := selfData.Bitset.CompareToStatRow(otherData.Bitset)
	verifyRes := verifyVectorHashes(selfData, otherData, sr, hasher)

	if verifyRes == norNotVerified || verifyRes == NvrSenderFault {
		return verifyRes
	}

	summarize(otherData.Bitset, verifyRes&^NvrSenderFault, sr, nodeStats)
	fmt.Printf("%v\n%v\n%v\n%v\n", selfData.Bitset, otherData.Bitset, sr, nodeStats)
	return verifyRes
}

// func catcher(selfData HashedNodeVector, otherData HashedNodeVector, hasher *FilteredSequenceHasher, sr **stats.Row) (verifyRes NodeVerificationResult, repeat bool) {
// 	defer func() {
// 		repeat = recover() != nil
// 	}()
// 	*sr = selfData.Bitset.CompareToStatRow(otherData.Bitset)
// 	return verifyVectorHashes(selfData, otherData, *sr, hasher), false
// }

func verifyVectorHashes(selfData HashedNodeVector, otherData HashedNodeVector, sr *stats.Row, hasher FilteredSequenceHasher) NodeVerificationResult {
	// TODO All GSH comparisons should be based on SIGNATURES! not on pure hashes

	verifyRes := norNotVerified
	valid := false
	gsh := selfData.TrustedVector
	altered := false

	switch {
	case sr.HasValues(NodeBitMissingHere):
		// we can't validate anything without requests
		// ...  check for updates and send requests
		return norNotVerified
	case sr.HasAllValuesOf(NodeBitSame, NodeBitLessTrustedHere):
		// check DoubtedGsh as is, if not then TrustedGSH with some locally-known NSH included
		fallthrough
	case sr.HasAllValuesOf(NodeBitSame, NodeBitLessTrustedThere):
		// check DoubtedGsh as is, if not then TrustedGSH with some locally-known NSH excluded
		valid = verifyRes.SetDoubted(selfData.DoubtedVector.Equals(otherData.DoubtedVector), false)
		if !sr.HasAllValues(NodeBitSame) {
			gsh = hasher.BuildHashByFilter(otherData.Bitset, sr, true)
			altered = true
		}
	case sr.HasValues(NodeBitDoubtedMissingHere):
		// check TrustedGSH only
		// validation of DoubtedGsh needs requests
		// ...  check for updates and send requests

		// if HasValues(NodeBitLessTrustedThere) then ... exclude some locally-known NSH
		// if HasValues(NodeBitMissingThere) then ... exclude some locally-known NSH
		// if HasValues(NodeBitLessTrustedHere) then ... include some locally-known NSH
		gsh = hasher.BuildHashByFilter(otherData.Bitset, sr, true)
		verifyRes.SetTrusted(gsh.Equals(otherData.TrustedVector), true)
		return verifyRes
	default:
		if sr.HasValues(NodeBitMissingThere) {
			// check DoubtedGsh with exclusions, then TrustedGSH with exclusions/inclusions
			gsh = hasher.BuildHashByFilter(otherData.Bitset, sr, false)
			valid = verifyRes.SetDoubted(gsh.Equals(otherData.DoubtedVector), true)
		} else {
			// check DoubtedGsh, then TrustedGSH with exclusions/inclusions
			gsh = selfData.DoubtedVector
			valid = verifyRes.SetDoubted(gsh.Equals(otherData.DoubtedVector), false)
		}

		// if HasValues(NodeBitLessTrustedThere) then ... exclude some locally-known NSH from TrustedGSH
		// if HasValues(NodeBitLessTrustedHere) then ... include some locally-known NSH to TrustedGSH
		gsh = hasher.BuildHashByFilter(otherData.Bitset, sr, true)
		altered = true
	}

	switch {
	case valid == gsh.Equals(otherData.TrustedVector):
		verifyRes.SetTrusted(valid, altered)
	case valid:
		verifyRes.SetTrusted(false, altered)
	default:
		// Dont set valid/fraud as there is an evident fraud/error by the sender
		// TODO fraud - there must be match when a wider set is in match
		verifyRes |= NvrSenderFault
	}
	return verifyRes
}

func summarize(otherDataBitset NodeBitset, verifyRes NodeVerificationResult, sr *stats.Row, nodeStats *stats.Row) {

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
