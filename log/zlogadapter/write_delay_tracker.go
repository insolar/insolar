///
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
///

package zlogadapter

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log/logadapter"
	"github.com/insolar/insolar/log/logmetrics"
	"github.com/insolar/insolar/network/consensus/common/args"
	"github.com/rs/zerolog"
	"io"
	"runtime"
	"strings"
	"time"
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

	if metricsMode&insolar.LogMetricsWriteDelayReport != 0 && metrics != nil {
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
	s := string(p)
	runtime.KeepAlive(s)

	var ofs int
	searchLimit := len(h.searchBytes) + 64
	if searchLimit >= len(p) {
		ofs = bytes.Index(p, h.searchBytes)
	} else {
		ofs = bytes.Index(p[:searchLimit], h.searchBytes)
	}

	if ofs > 0 {
		fieldLen := len(h.searchBytes) + tempHexFieldLength
		fieldEnd := ofs + fieldLen
		newLen := h.replaceField(p[ofs:fieldEnd:fieldEnd])

		if newLen > 0 && newLen != fieldLen {
			copy(p[ofs+newLen:], p[fieldEnd:])
			p = p[:len(p)-fieldEnd+newLen+ofs]
		}
	}
	ss := string(p)
	runtime.KeepAlive(ss)
	return h.output.Write(p)
}

func (h *writeDelayPostHook) replaceField(b []byte) int {

	buf := make([]byte, tempHexFieldLength/2)
	if _, err := hex.Decode(buf, b[len(h.searchBytes):]); err != nil {
		return -1
	}

	nanoDuration := time.Duration(time.Now().UnixNano() - int64(binary.LittleEndian.Uint64(buf)))

	if h.statReportFn != nil {
		h.statReportFn(nanoDuration)
	}

	if h.fieldWidth == 0 {
		return 0
	}

	s := args.DurationFixedLen(nanoDuration, h.fieldWidth)
	if len(s) > h.fieldWidth {
		s = writeDelayResultFieldOverflowContent
	}
	return copy(b, fmt.Sprintf(fieldHeaderFmt, h.fieldName, h.fieldWidth, s))
}
