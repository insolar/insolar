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

	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/stats"
)

type ComparedState uint8

const (
	ComparedSame ComparedState = iota
	ComparedLessTrustedHere
	ComparedLessTrustedThere
	ComparedMissingThere
	ComparedDoubtedMissingHere
	ComparedMissingHere
	maxComparedTypes
)

func FmtComparedStat(v uint8) string {
	switch ComparedState(v) {
	case ComparedSame:
		return "≡"
	case ComparedLessTrustedHere:
		return "≤"
	case ComparedLessTrustedThere:
		return "≥"
	case ComparedMissingThere:
		return ">"
	case ComparedDoubtedMissingHere:
		return "≨"
	case ComparedMissingHere:
		return "<"
	default:
		return fmt.Sprintf("%d", v)
	}
}

func NewComparedBitsetRow(columnCount int) ComparedBitsetRow {
	return ComparedBitsetRow{innerRow{stats.NewStatRow(uint8(maxComparedTypes)-1, columnCount)}}
}

type ComparedBitsetRow struct {
	innerRow
}

func (r *ComparedBitsetRow) Get(column int) ComparedState {
	return ComparedState(r.innerRow.Get(column))
}

func (r *ComparedBitsetRow) HasValues(value ComparedState) bool {
	return r.innerRow.HasValues(uint8(value))
}

func (r *ComparedBitsetRow) HasAllValues(value ComparedState) bool {
	return r.innerRow.HasAllValues(uint8(value))
}

func (r *ComparedBitsetRow) HasAllValuesOf(value0, value1 ComparedState) bool {
	return r.innerRow.HasAllValuesOf(uint8(value0), uint8(value1))
}

func (r *ComparedBitsetRow) GetSummaryByValue(value ComparedState) uint16 {
	return r.innerRow.GetSummaryByValue(uint8(value))
}

func (r *ComparedBitsetRow) Set(column int, value ComparedState) ComparedState {
	return ComparedState(r.innerRow.Set(column, uint8(value)))
}

func (r ComparedBitsetRow) String() string {
	return r.StringFull()
}

func (r ComparedBitsetRow) StringFull() string {
	return r.innerRow.Row.StringFmt(FmtComparedStat, true)
}

func (r ComparedBitsetRow) StringSummary() string {
	return r.innerRow.Row.StringSummaryFmt(FmtComparedStat)
}

func CompareToStatRow(b member.StateBitset, otherDataBitset member.StateBitset) ComparedBitsetRow {

	if otherDataBitset.Len() != b.Len() {
		// TODO handle different bitset size
		panic("different bitset size")
	}

	bitStats := NewComparedBitsetRow(b.Len())

	for i, fHere := range b {
		fThere := otherDataBitset[i]
		var bitStat ComparedState

		switch {
		case fHere == fThere:
			// all the same, proceed as it is
			bitStat = ComparedSame
		case fThere.IsTimeout():
			// we can skip this NSH and recalculate
			bitStat = ComparedMissingThere

		case fHere.IsTimeout():
			if fThere.IsTrusted() {
				bitStat = ComparedMissingHere
			} else {
				bitStat = ComparedDoubtedMissingHere
			}
		case fHere.IsTrusted() == fThere.IsTrusted():
			// NB! fraud is considered as "doubted"
			bitStat = ComparedSame
		case fThere.IsTrusted():
			// we don't trust, other one does
			bitStat = ComparedLessTrustedHere
		default: // fHere.IsTrusted()
			// we trust, other one doesn't
			bitStat = ComparedLessTrustedThere
		}
		bitStats.Set(i, bitStat)
	}

	return bitStats
}
