package stats

import (
	"fmt"
	"strings"
)

type Column struct {
	colIndex uint16
	summary  []uint16
}

func (c *Column) ColumnIndex() uint16 {
	return c.colIndex
}

func (c *Column) GetSummaryByValue(value uint8) uint16 {
	return c.summary[value]
}

func (c *Column) StringFmt(fmtFn RowValueFormatFunc) string {
	if fmtFn == nil {
		fmtFn = defaultValueFmt
	}

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%03d", c.colIndex))
	stringSummary16Fmt(c.summary, &builder, fmtFn)
	return builder.String()
}

func (c Column) String() string {
	return fmt.Sprintf("%03d%v", c.colIndex, c.summary)
}

func (c *Column) GetSummary() []uint16 {
	return append(make([]uint16, 0, len(c.summary)), c.summary...)
}

func (c Column) Equals(o Column) bool {
	for i, tS := range c.summary {
		if tS != o.summary[i] {
			return false
		}
	}
	return true
}
