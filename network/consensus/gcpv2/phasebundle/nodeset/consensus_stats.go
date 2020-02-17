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

type ConsensusStat uint8

const (
	ConsensusStatUnknown ConsensusStat = iota
	ConsensusStatTrusted
	ConsensusStatDoubted
	ConsensusStatMissingThere
	ConsensusStatFraudSuspect
	ConsensusStatFraud
	maxConsensusStat
)

const ConsensusStatMissingHere = ConsensusStatUnknown

func FmtConsensusStat(v uint8) string {
	switch ConsensusStat(v) {
	case ConsensusStatUnknown:
		return "U"
	case ConsensusStatTrusted:
		return "T"
	case ConsensusStatDoubted:
		return "D"
	case ConsensusStatMissingThere:
		return "Ø"
	case ConsensusStatFraudSuspect:
		return "f"
	case ConsensusStatFraud:
		return "F"
	default:
		return fmt.Sprintf("%d", v)
	}
}

type innerStatTable struct {
	stats.StatTable
}

type ConsensusStatTable struct {
	innerStatTable
}

func NewConsensusStatTable(nodeCount int) ConsensusStatTable {
	return ConsensusStatTable{innerStatTable{stats.NewStatTable(uint8(maxConsensusStat)-1, nodeCount)}}
}

func (t *ConsensusStatTable) NewRow() *ConsensusStatRow {
	nr := NewConsensusStatRow(t.ColumnCount())
	return &nr
}

func (t *ConsensusStatTable) AddRow(row *ConsensusStatRow) int {
	return t.innerStatTable.AddRow(&row.Row)
}

func (t *ConsensusStatTable) PutRow(rowIndex int, row *ConsensusStatRow) {
	t.innerStatTable.PutRow(rowIndex, &row.Row)
}

func (t *ConsensusStatTable) GetRow(rowIndex int) (*ConsensusStatRow, bool) {
	row, ok := t.innerStatTable.GetRow(rowIndex)
	if !ok {
		return nil, false
	}
	return &ConsensusStatRow{innerRow{*row}}, true
}

func (t *ConsensusStatTable) GetColumn(colIndex int) *ConsensusStatColumn {
	return &ConsensusStatColumn{innerStatColumn{t.StatTable.GetColumn(colIndex)}}
}

func (t *ConsensusStatTable) AsText(header string) string {
	return t.TableFmt(header, FmtConsensusStat)
}

func (t *ConsensusStatTable) EqualsTyped(o *ConsensusStatTable) bool {
	return o != nil && t.StatTable.Equals(&o.StatTable)
}

type innerStatColumn struct {
	*stats.Column
}

type ConsensusStatColumn struct {
	innerStatColumn
}

func (c *ConsensusStatColumn) GetSummaryByValue(value ConsensusStat) uint16 {
	return c.innerStatColumn.GetSummaryByValue(uint8(value))
}

func (c *ConsensusStatColumn) String() string {
	return c.StringFmt(FmtConsensusStat)
}

func NewConsensusStatRow(columnCount int) ConsensusStatRow {
	return ConsensusStatRow{innerRow{stats.NewStatRow(uint8(maxConsensusStat)-1, columnCount)}}
}

type ConsensusStatRow struct {
	innerRow
}

func (r *ConsensusStatRow) Get(column int) ConsensusStat {
	return ConsensusStat(r.innerRow.Get(column))
}

func (r *ConsensusStatRow) HasValues(value ConsensusStat) bool {
	return r.innerRow.HasValues(uint8(value))
}

func (r *ConsensusStatRow) HasAllValues(value ConsensusStat) bool {
	return r.innerRow.HasAllValues(uint8(value))
}

func (r *ConsensusStatRow) HasAllValuesOf(value0, value1 ConsensusStat) bool {
	return r.innerRow.HasAllValuesOf(uint8(value0), uint8(value1))
}

func (r *ConsensusStatRow) GetSummaryByValue(value ConsensusStat) uint16 {
	return r.innerRow.GetSummaryByValue(uint8(value))
}

func (r *ConsensusStatRow) Set(column int, value ConsensusStat) ConsensusBitsetEntry {
	return ConsensusBitsetEntry(r.innerRow.Set(column, uint8(value)))
}

func (r ConsensusStatRow) String() string {
	return r.innerRow.Row.StringFmt(FmtConsensusStat, true)
}

func StateToConsensusStatRow(b member.StateBitset) ConsensusStatRow {

	nodeStats := NewConsensusStatRow(b.Len())
	for i, v := range b {
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

	return nodeStats
}
