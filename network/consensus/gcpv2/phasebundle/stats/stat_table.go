package stats

import (
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
)

type StatTable struct {
	columns  []Column
	rows     []*Row
	rowCount int
	summary  []uint32
}

func NewStatTable(maxValue uint8, columns int) StatTable {
	if columns > math.MaxUint16 {
		panic("too many columns")
	}
	r := StatTable{columns: make([]Column, columns), rows: make([]*Row, 0, columns)}
	for i := 0; i < columns; i++ {
		r.columns[i].colIndex = uint16(i)
		r.columns[i].summary = make([]uint16, maxValue+1)
	}
	r.summary = make([]uint32, maxValue+1)
	return r
}

func (t *StatTable) NewRow() *Row {
	nr := NewStatRow(t.MaxValue(), t.ColumnCount())
	return &nr
}

func (t *StatTable) AddRow(row *Row) int {
	row.ensureForTable(t)
	row.rowIndex = len(t.rows)
	t.rows = append(t.rows, row)
	t.rowCount++

	for i, v := range row.values {
		t.columns[i].summary[v]++
		t.summary[v]++
	}
	return row.rowIndex
}

func (t *StatTable) PutRow(rowIndex int, row *Row) {
	row.ensureForTable(t)
	switch {
	case rowIndex == len(t.rows):
		t.rows = append(t.rows, row)
	case rowIndex > len(t.rows):
		t.rows = append(t.rows, make([]*Row, rowIndex-len(t.rows)+1)...)
		t.rows[rowIndex] = row
	case t.rows[rowIndex] != nil:
		panic("row is in use")
	default:
		t.rows[rowIndex] = row
	}
	row.rowIndex = rowIndex
	t.rowCount++

	for i, v := range row.values {
		t.columns[i].summary[v]++
		t.summary[v]++
	}
}

func (t *StatTable) GetRow(rowIndex int) (row *Row, ok bool) {
	if rowIndex >= len(t.rows) {
		return nil, false
	}
	row = t.rows[rowIndex]
	if row == nil {
		return nil, false
	}
	rowCopy := *row
	return &rowCopy, true
}

func (t *StatTable) RemoveRow(rowIndex int) (ok bool) {
	if rowIndex >= len(t.rows) {
		return false
	}
	row := t.rows[rowIndex]
	if row == nil {
		return false
	}
	t.rows[rowIndex] = nil
	for i, v := range row.values {
		t.columns[i].summary[v]--
		t.summary[v]--
	}
	t.rowCount--
	return true
}

func (t *StatTable) RowCount() int {
	return t.rowCount
}

func (t *StatTable) ColumnCount() int {
	return len(t.columns)
}

func (t *StatTable) GetSummaryByValue(value uint8) uint32 {
	return t.summary[value]
}

func (t *StatTable) GetSummary() []uint32 {
	return append(make([]uint32, 0, len(t.summary)), t.summary...)
}

func (t *StatTable) MaxValue() uint8 {
	return uint8(len(t.summary) - 1)
}

func (t *StatTable) GetColumn(colIndex int) *Column {
	return &t.columns[colIndex]
}

func (t *StatTable) String() string {
	return fmt.Sprintf("stats[v=%d, c=%d, r=%d/%d]", t.MaxValue()+1, t.ColumnCount(), t.RowCount(), len(t.rows))
}

func (t *StatTable) AsText(header string) string {
	return t.TableFmt(header, nil)
}

func (t *StatTable) Equals(o *StatTable) bool {
	if t == nil || o == nil {
		return false
	}
	if t == o {
		return true
	}
	if t.rowCount != o.rowCount || len(t.columns) != len(o.columns) || len(t.summary) != len(o.summary) ||
		len(t.rows) != len(o.rows) {
		return false
	}
	for i, tS := range t.summary {
		if tS != o.summary[i] {
			return false
		}
	}
	for j, tC := range t.columns {
		if !tC.Equals(o.columns[j]) {
			return false
		}
	}
	for j, tR := range t.rows {
		oR := o.rows[j]
		if tR == oR { // for both nil
			continue
		}
		if tR == nil || oR == nil || !tR.equals(oR) {
			return false
		}
	}
	return true
}

func (t *StatTable) TableFmt(header string, fmtFn RowValueFormatFunc) string {
	widths := make([]int, t.ColumnCount())
	builder := strings.Builder{}
	builder.WriteString(header)
	if fmtFn != nil {
		builder.WriteString("\nLEGEND [")
		for i := uint8(0); i <= t.MaxValue(); i++ {
			if i != 0 {
				builder.WriteRune(' ')
			}
			builder.WriteString(fmtFn(i))
		}
		builder.WriteString(fmt.Sprintf("] RowCount=%d", t.rowCount))
	}
	builder.WriteString("\n###")
	for i, c := range t.columns {
		s := fmt.Sprintf("|%03d%+v", c.colIndex, c.summary)
		widths[i] = utf8.RuneCountInString(s)
		builder.WriteString(s)
	}
	builder.WriteString("|∑")
	stringSummary32Fmt(t.summary, &builder, fmtFn)
	builder.WriteByte('\n')
	for i, r := range t.rows {
		if r == nil {
			continue
		}
		builder.WriteString(fmt.Sprintf("%03d", i))
		for j, v := range r.values {
			if fmtFn == nil {
				builder.WriteString(fmt.Sprintf("|%*d", widths[j]-1, v))
			} else {
				builder.WriteString(fmt.Sprintf("|%*s", widths[j]-1, fmtFn(v)))
			}
		}
		builder.WriteString("|∑")

		if fmtFn == nil {
			builder.WriteString(fmt.Sprintf("%+v\n", r.GetSummary()))
		} else {
			stringSummary16Fmt(r.GetSummary(), &builder, fmtFn)
			builder.WriteByte('\n')
		}
	}
	builder.WriteByte('\n')
	return builder.String()
}
