// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package nodeset

import (
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/stats"
)

func NewMemberBitsetRow(columnCount int) MemberBitsetRow {
	return MemberBitsetRow{innerRow{stats.NewStatRow(uint8(member.MaxBitsetEntry)-1, columnCount)}}
}

type innerRow struct {
	stats.Row
}

/* TODO unused */
type MemberBitsetRow struct {
	innerRow
}

func (r *MemberBitsetRow) Get(column int) member.BitsetEntry {
	return member.BitsetEntry(r.innerRow.Get(column))
}

func (r *MemberBitsetRow) HasValues(value member.BitsetEntry) bool {
	return r.innerRow.HasValues(uint8(value))
}

func (r *MemberBitsetRow) HasAllValues(value member.BitsetEntry) bool {
	return r.innerRow.HasAllValues(uint8(value))
}

func (r *MemberBitsetRow) HasAllValuesOf(value0, value1 member.BitsetEntry) bool {
	return r.innerRow.HasAllValuesOf(uint8(value0), uint8(value1))
}

func (r *MemberBitsetRow) GetSummaryByValue(value member.BitsetEntry) uint16 {
	return r.innerRow.GetSummaryByValue(uint8(value))
}

func (r *MemberBitsetRow) Set(column int, value member.BitsetEntry) member.BitsetEntry {
	return member.BitsetEntry(r.innerRow.Set(column, uint8(value)))
}
