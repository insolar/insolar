// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package nodeset

import (
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/phasebundle/stats"
)

/* MUST be based on BitsetEntry to reuse serialization */
type ConsensusBitsetEntry member.BitsetEntry

const (
	CbsIncluded  = ConsensusBitsetEntry(member.BeHighTrust)
	CbsSuspected = ConsensusBitsetEntry(member.BeBaselineTrust)
	CbsExcluded  = ConsensusBitsetEntry(member.BeTimeout)
	CbsFraud     = ConsensusBitsetEntry(member.BeFraud)
)

func FmtConsensusBitsetEntry(v uint8) string {
	return member.FmtBitsetEntry(v)
}

func (s ConsensusBitsetEntry) String() string {
	return FmtConsensusBitsetEntry(uint8(s))
}

func NewConsensusBitsetRow(columnCount int) ConsensusBitsetRow {
	return ConsensusBitsetRow{innerRow{stats.NewStatRow(uint8(member.MaxBitsetEntry)-1, columnCount)}}
}

type ConsensusBitsetRow struct {
	innerRow
}

func (r *ConsensusBitsetRow) Get(column int) ConsensusBitsetEntry {
	return ConsensusBitsetEntry(r.innerRow.Get(column))
}

func (r *ConsensusBitsetRow) HasValues(value ConsensusBitsetEntry) bool {
	return r.innerRow.HasValues(uint8(value))
}

func (r *ConsensusBitsetRow) HasAllValues(value ConsensusBitsetEntry) bool {
	return r.innerRow.HasAllValues(uint8(value))
}

func (r *ConsensusBitsetRow) HasAllValuesOf(value0, value1 ConsensusBitsetEntry) bool {
	return r.innerRow.HasAllValuesOf(uint8(value0), uint8(value1))
}

func (r *ConsensusBitsetRow) GetSummaryByValue(value ConsensusBitsetEntry) uint16 {
	return r.innerRow.GetSummaryByValue(uint8(value))
}

func (r *ConsensusBitsetRow) Set(column int, value ConsensusBitsetEntry) ConsensusBitsetEntry {
	return ConsensusBitsetEntry(r.innerRow.Set(column, uint8(value)))
}

func (r ConsensusBitsetRow) String() string {
	return r.innerRow.Row.StringFmt(FmtConsensusBitsetEntry, true)
}
