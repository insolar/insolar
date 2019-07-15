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

func LocalToConsensusStatRow(b member.StateBitset) *ConsensusStatRow {

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

	return &nodeStats
}
