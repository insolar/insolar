// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package args

import (
	"fmt"
	"math"
	"strings"
	"time"
)

/*
	NB! "µs" require 3 bytes, not 2 bytes, hence the difference with "ms" on a number of decimal positions
*/
func DurationFixedLen(d time.Duration, expectedLen int) string {
	return durationFixedLen(d, d, expectedLen)
}

func durationFixedLen(d, base time.Duration, expectedLen int) string {

	base, baseUnit := metricDurationBase(base)
	if base < time.Minute {
		return fmtMetric(d, base, baseUnit, expectedLen)
	}

	return fmtAboveSeconds(d, expectedLen)
}

func fmtMetric(d time.Duration, base time.Duration, baseUnit string, expectedLen int) string {
	w := expectedLen - len(baseUnit)
	v := float64(d) / float64(base)

	if w < 1 {
		w = 1
	}
	vRounded := math.Round(v)
	if vRounded >= math.Pow10(w) {
		return durationFixedLen(d, base*1000, expectedLen)
	}

	if w < 3 || vRounded >= math.Pow10(w-2) {
		return fmt.Sprintf("%.0f%s", v, baseUnit)
	}

	b := strings.Builder{}
	b.Grow(w + len(baseUnit))

	vInt, vFrac := math.Modf(v)

	b.WriteString(fmt.Sprintf("%d.", uint64(vInt)))
	decimalCount := w - b.Len()
	decimals := uint64(vFrac * math.Pow10(decimalCount))
	b.WriteString(fmt.Sprintf("%0*d", decimalCount, decimals))
	b.WriteString(baseUnit)
	return b.String()
}

func metricDurationBase(d time.Duration) (time.Duration, string) {
	if d < 10*time.Minute {
		switch {
		case d > 500*time.Millisecond:
			return time.Second, "s"
		case d > 500*time.Microsecond:
			return time.Millisecond, "ms"
		default:
			return time.Microsecond, "µs"
		}
	}
	return time.Minute, "m"
}

func fmtAboveSeconds(d time.Duration, expectedLen int) string {
	minutes := d / time.Minute
	seconds := (d - minutes*time.Minute) / time.Second

	if minutes < 600 {
		return fmtPortions(uint64(seconds), "s", uint64(minutes), "m", expectedLen)
	}

	hours := minutes / 60
	minutes -= hours * 60
	return fmtPortions(uint64(minutes), "m", uint64(hours), "h", expectedLen)
}

func fmtPortions(valueLo uint64, unitLo string, valueHi uint64, unitHi string, expectedLen int) string {
	if valueLo > 59 {
		panic("illegal value")
	}

	switch {
	case valueHi > 0:
		break
	case expectedLen-len(unitLo) < 2:
		if valueLo > 9 {
			if valueLo > 30 {
				valueLo = 0
				valueHi++
			}
			break
		}
		fallthrough
	default:
		return fmt.Sprintf("%d%s", valueLo, unitLo)
	}

	vHi := fmt.Sprintf("%d%s", valueHi, unitHi)
	if len(vHi)+2+len(unitLo) > expectedLen {
		if valueLo > 30 {
			return fmt.Sprintf("%d%s", valueHi+1, unitHi)
		}
		return vHi
	}
	return vHi + fmt.Sprintf("%02d%s", valueLo, unitLo)
}

func BitCountToMaxDecimalCount(bitLen int) int {

	if bitLen == 0 {
		return 0
	}

	const k = 3321928 // 3.3219280948873623478703194294894 = Ln(10) / Ln(2)
	return int(1 + (uint64(bitLen-1)*1000000)/k)
}
