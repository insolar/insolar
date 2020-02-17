// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
