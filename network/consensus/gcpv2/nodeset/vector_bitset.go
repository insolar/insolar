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

	"github.com/insolar/insolar/network/consensus/gcpv2/stats"
)

type NodeBitsetEntry uint8

const (
	NbsHighTrust NodeBitsetEntry = iota
	NbsLimitedTrust
	NbsBaselineTrust
	NbsTimeout
	NbsFraud
	maxNodeBitsetEntry
)

/* MUST be based on NodeBitsetEntry to reuse serialization */
type ConsensusBitsetEntry NodeBitsetEntry

const (
	CbsIncluded  = ConsensusBitsetEntry(NbsHighTrust)
	CbsSuspected = ConsensusBitsetEntry(NbsBaselineTrust)
	CbsExcluded  = ConsensusBitsetEntry(NbsTimeout)
	CbsFraud     = ConsensusBitsetEntry(NbsFraud)
)

func (v NodeBitsetEntry) IsTrusted() bool { return v < NbsBaselineTrust }
func (v NodeBitsetEntry) IsTimeout() bool { return v == NbsTimeout }
func (v NodeBitsetEntry) IsFraud() bool   { return v == NbsFraud }

func (s NodeBitsetEntry) String() string {
	return FmtNodeBitsetEntry(uint8(s))
}

func FmtNodeBitsetEntry(s uint8) string {
	switch NodeBitsetEntry(s) {
	case NbsHighTrust:
		return "H"
	case NbsLimitedTrust:
		return "L"
	case NbsBaselineTrust:
		return "B"
	case NbsTimeout:
		return "Ø"
	case NbsFraud:
		return "F"
	default:
		return fmt.Sprintf("?%d", s)
	}
}

func (s ConsensusBitsetEntry) String() string {
	return FmtNodeBitsetEntry(uint8(s))
}

type NodeBitset []NodeBitsetEntry

func (b *NodeBitset) Len() int {
	if b == nil {
		return 0
	}
	return len(*b)
}

const (
	NodeBitSame uint8 = iota
	NodeBitLessTrustedHere
	NodeBitLessTrustedThere
	NodeBitMissingThere
	NodeBitDoubtedMissingHere
	NodeBitMissingHere
	maxNodeBitClassification
)

func (b *NodeBitset) CompareToStatRow(otherDataBitset NodeBitset) *stats.Row {

	if otherDataBitset.Len() != b.Len() {
		// TODO handle different bitset size
		panic("different bitset size")
	}

	bitStats := stats.NewStatRow(maxNodeBitClassification-1, b.Len())

	for i, fHere := range *b {
		fThere := otherDataBitset[i]
		var bitStat uint8

		switch {
		case fHere == fThere:
			// all the same, proceed as it is
			bitStat = NodeBitSame
		case fThere.IsTimeout():
			// we can skip this NSH and recalculate
			bitStat = NodeBitMissingThere

		case fHere.IsTimeout():
			if fThere.IsTrusted() {
				bitStat = NodeBitMissingHere
			} else {
				bitStat = NodeBitDoubtedMissingHere
			}
		case fHere.IsTrusted() == fThere.IsTrusted():
			// fraud is considered as "doubted"
			bitStat = NodeBitSame
		case fThere.IsTrusted():
			// we don't trust, other one does
			bitStat = NodeBitLessTrustedHere
		default: // fHere.IsTrusted()
			// we trust, other one doesn't
			bitStat = NodeBitLessTrustedThere
		}
		bitStats.Set(i, bitStat)
	}

	return &bitStats
}

func (b *NodeBitset) LocalToConsensusStatRow() *stats.Row {

	nodeStats := stats.NewStatRow(maxConsensusStat-1, b.Len())

	for i, v := range *b {
		switch {
		case v.IsTimeout():
			nodeStats.Set(i, ConsensusStatMissingThere)
		case v.IsFraud():
			nodeStats.Set(i, ConsensusStatFraud)
		case v.IsTrusted():
			nodeStats.Set(i, ConsensusStatTrusted)
		default:
			nodeStats.Set(i, ConsensusStatDoubted)
		}
	}

	return &nodeStats
}

// func (b NodeBitset) String() string {
//
// }
