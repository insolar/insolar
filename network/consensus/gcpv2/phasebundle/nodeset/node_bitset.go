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
