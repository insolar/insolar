// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package zlogadapter

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/rs/zerolog"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log/logadapter"
	"github.com/insolar/insolar/log/logmetrics"
	"github.com/insolar/insolar/network/consensus/common/args"
)

const internalTempFieldName = "_TWD_"
const fieldHeaderFmt = `,"%s":"%*v`
const tempHexFieldLength = 16 // HEX for Uint64
const writeDelayResultFieldOverflowContent = "ovrflw"
const writeDelayResultFieldMinWidth = len(writeDelayResultFieldOverflowContent)
const writeDelayPreferTrim = false

func getWriteDelayConfig(metrics *logmetrics.MetricsHelper,
	config logadapter.BuildConfig) (needsHook bool, fieldName string, reportFn logmetrics.DurationReportFunc) {

	metricsMode := config.Instruments.MetricsMode
	if metricsMode&(insolar.LogMetricsWriteDelayField|insolar.LogMetricsWriteDelayReport) == 0 {
		return
	}

	if metricsMode&insolar.LogMetricsWriteDelayField != 0 {
		fieldName = "writeDuration"
	}

	if metrics != nil {
		reportFn = metrics.GetOnWriteDurationReport()
	}

	return len(fieldName) != 0 || reportFn != nil, fieldName, reportFn
}

func getWriteDelayHookParams(fieldName string, preferTrim bool) (fieldWidth int, searchField string) {
	searchField = internalTempFieldName
	if len(fieldName) != 0 {
		fieldWidth = writeDelayResultFieldMinWidth
		paddingLen := (len(fieldName) + fieldWidth) - (len(searchField) + tempHexFieldLength)

		if paddingLen < 0 {
			// we have more space than needed
			if !preferTrim {
				// ensure proper wipe out of temporary field data
				fieldWidth -= paddingLen
			}
		} else {
			if paddingLen > len(fieldName) {
				searchField += fieldName + strings.Repeat("_", paddingLen-len(fieldName))
			} else {
				searchField += fieldName[:paddingLen]
			}
		}
	}
	return
}

func newWriteDelayPreHook(fieldName string, preferTrim bool) *writeDelayHook {
	_, searchField := getWriteDelayHookParams(fieldName, preferTrim)
	return &writeDelayHook{searchField: searchField}
}

type writeDelayHook struct {
	searchField string
}

func (h *writeDelayHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	buf := make([]byte, tempHexFieldLength/2)
	binary.LittleEndian.PutUint64(buf, uint64(time.Now().UnixNano()))

	e.Hex(h.searchField, buf)
}

func newWriteDelayPostHook(output io.Writer, fieldName string, preferTrim bool, statReportFn logmetrics.DurationReportFunc) *writeDelayPostHook {
	fieldWidth, searchField := getWriteDelayHookParams(fieldName, preferTrim)
	return &writeDelayPostHook{
		output:       output,
		searchBytes:  []byte(fmt.Sprintf(fieldHeaderFmt, searchField, 0, "")),
		fieldName:    fieldName,
		fieldWidth:   fieldWidth,
		statReportFn: statReportFn,
	}
}

type writeDelayPostHook struct {
	output       io.Writer
	searchBytes  []byte
	fieldName    string
	fieldWidth   int
	statReportFn func(d time.Duration)
}

func (h *writeDelayPostHook) Write(p []byte) (n int, err error) {
	var ofs int
	//searchLimit := len(h.searchBytes) + 32
	//if searchLimit >= len(p) {
	ofs = bytes.Index(p, h.searchBytes)
	//} else {
	//	ofs = bytes.Index(p[:searchLimit], h.searchBytes)
	//}

	//s := string(p)
	//runtime.KeepAlive(s)

	if ofs < 0 {
		return h.output.Write(p)
	}

	fieldLen := len(h.searchBytes) + tempHexFieldLength
	fieldEnd := ofs + fieldLen
	newLen, startedAt := h.replaceField(p[ofs:fieldEnd:fieldEnd])

	if newLen > 0 && newLen != fieldLen {
		copy(p[ofs+newLen:], p[fieldEnd:])
		p = p[:len(p)-fieldEnd+newLen+ofs]
	}
	n, err = h.output.Write(p)

	if h.statReportFn != nil && startedAt > 0 {
		nanoDuration := time.Duration(time.Now().UnixNano() - startedAt)
		h.statReportFn(nanoDuration)
	}

	return n, err
}

func (h *writeDelayPostHook) replaceField(b []byte) (int, int64) {

	buf := make([]byte, tempHexFieldLength/2)
	if _, err := hex.Decode(buf, b[len(h.searchBytes):]); err != nil {
		return -1, 0
	}
	startedAt := int64(binary.LittleEndian.Uint64(buf))
	nanoDuration := time.Duration(time.Now().UnixNano() - startedAt)

	if h.fieldWidth == 0 {
		return 0, startedAt
	}

	s := args.DurationFixedLen(nanoDuration, h.fieldWidth)
	w := h.fieldWidth
	if len(s) > w {
		s = writeDelayResultFieldOverflowContent
	} else {
		w -= len(s) - utf8.RuneCountInString(s)
	}
	rs := fmt.Sprintf(fieldHeaderFmt, h.fieldName, w, s)
	return copy(b, rs), startedAt
}
