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

package stats

import (
	"fmt"
	"math"
	"strings"
)

type RowValueFormatFunc func(v uint8) string

type Row struct {
	rowIndex      int
	values        []uint8
	summary       []uint16
	customOptions uint32
}

func NewStatRow(maxValue uint8, columns int) Row {
	if columns > math.MaxUint16 {
		panic("too many columns")
	}
	return Row{rowIndex: -1, values: make([]uint8, columns), summary: make([]uint16, maxValue+1)}
}

func (r *Row) GetCustomOptions() uint32 {
	return r.customOptions
}

func (r *Row) SetCustomOptions(v uint32) {
	r.customOptions = v
}

func (r *Row) Len() int {
	return len(r.values)
}

func (r *Row) IsEmpty() bool {
	return len(r.values) == 0 && len(r.summary) == 0
}

func (r *Row) Set(column int, value uint8) uint8 {
	r.ensureUpdateable()
	prev := r.values[column]
	if prev != value {
		r.summary[prev]--
		r.values[column] = value
		r.summary[value]++
	}
	return prev
}

func (r *Row) ensureForTable(t *StatTable) {
	r.ensureUpdateable()
	if r.ColumnCount() != t.ColumnCount() {
		panic("column count mismatched")
	}
	if r.MaxValue() != t.MaxValue() {
		panic("max value mismatched")
	}
}

func (r *Row) ensureUpdateable() {
	if !r.CanUpdate() {
		panic("row is in the table or uninitialized")
	}
}

func (r *Row) Get(column int) uint8 {
	return r.values[column]
}

func (r *Row) ColumnCount() int {
	return len(r.values)
}

func (r *Row) MaxValue() uint8 {
	return uint8(len(r.summary) - 1)
}

func (r *Row) GetRowIndex() int {
	if r.rowIndex < 0 {
		return -1
	}
	return r.rowIndex
}

func (r *Row) HasValues(value uint8) bool {
	return r.GetSummaryByValue(value) > 0
}

func (r *Row) HasAllValues(value uint8) bool {
	if value != 0 {
		return r.summary[value] == uint16(len(r.values))
	}
	return r.summary[0] == 0 // zero is reverse-counted
}

func (r *Row) HasAllValuesOf(value0, value1 uint8) bool {
	if value0 == 0 || value1 == 0 {
		return r.summary[value0]+r.summary[value1] == 0 // zero is reverse-counted
	}
	return r.summary[value0]+r.summary[value1] == uint16(len(r.values))
}

func (r *Row) GetSummaryByValue(value uint8) uint16 {
	if value != 0 {
		return r.summary[value]
	}
	return r.summary[0] + uint16(len(r.values)) // zero is reverse-counted
}

func (r *Row) GetSummary() []uint16 {
	v := append(make([]uint16, 0, len(r.summary)), r.summary...)
	v[0] += uint16(len(r.values)) // zero is reverse-counted
	return v
}

func (r *Row) CanUpdate() bool {
	return r.rowIndex < 0
}

func (r *Row) String() string {
	return fmt.Sprintf("%v∑%v", r.values, r.GetSummary())
}

func (r *Row) StringFmt(fmtFn RowValueFormatFunc, summaryPrefixes bool) string {
	if fmtFn == nil {
		fmtFn = defaultValueFmt
	}

	builder := strings.Builder{}
	builder.WriteRune('[')
	for i, v := range r.values {
		if i > 0 {
			builder.WriteRune(' ')
		}
		builder.WriteString(fmtFn(v))
	}
	builder.WriteRune(']')
	builder.WriteRune('∑')
	if summaryPrefixes {
		stringSummary16Fmt(r.GetSummary(), &builder, fmtFn)
	} else {
		builder.WriteString(fmt.Sprintf("%v", r.GetSummary()))
	}
	return builder.String()
}

func (r *Row) StringSummaryFmt(fmtFn RowValueFormatFunc) string {
	if fmtFn == nil {
		fmtFn = defaultValueFmt
	}

	builder := strings.Builder{}
	stringSummary16Fmt(r.GetSummary(), &builder, fmtFn)
	return builder.String()
}

func (r *Row) Equals(o *Row) bool {
	if r == nil || o == nil {
		return false
	}
	if r == o {
		return true
	}
	return r.equals(o)
}

func (r *Row) equals(o *Row) bool {
	for i, tS := range r.summary {
		if tS != o.summary[i] {
			return false
		}
	}
	for j, tC := range r.values {
		if tC != o.values[j] {
			return false
		}
	}
	return true
}

func stringSummary16Fmt(summary []uint16, builder *strings.Builder, fmtFn RowValueFormatFunc) {
	builder.WriteRune('[')
	for i, v := range summary {
		if i > 0 {
			builder.WriteRune(' ')
		}

		if v > 0 {
			builder.WriteString(fmtFn(uint8(i)))
			builder.WriteString(fmt.Sprintf("%v", v))
		} else {
			builder.WriteByte(' ')
			builder.WriteByte(' ')
		}
	}
	builder.WriteRune(']')
}

func stringSummary32Fmt(summary []uint32, builder *strings.Builder, fmtFn RowValueFormatFunc) {
	builder.WriteRune('[')
	for i, v := range summary {
		if i > 0 {
			builder.WriteRune(' ')
		}

		if v > 0 {
			builder.WriteString(fmtFn(uint8(i)))
			builder.WriteString(fmt.Sprintf("%v", v))
		} else {
			builder.WriteByte(' ')
			builder.WriteByte(' ')
		}
	}
	builder.WriteRune(']')
}

func defaultValueFmt(v uint8) string {
	return fmt.Sprintf("%v", v)
}
