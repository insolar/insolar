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
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/insolar/insolar/log/inssyslog"
	"github.com/insolar/insolar/log/logadapter"
	"github.com/insolar/insolar/log/logmetrics"
	"github.com/insolar/insolar/network/consensus/common/args"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/insolar/insolar/insolar"
)

var insolarPrefix = "github.com/insolar/insolar/"

func trimInsolarPrefix(file string, line int) string {
	var skip = 0
	if idx := strings.Index(file, insolarPrefix); idx != -1 {
		skip = idx + len(insolarPrefix)
	}
	return file[skip:] + ":" + strconv.Itoa(line)
}

func init() {
	zerolog.TimeFieldFormat = logadapter.TimestampFormat
	zerolog.CallerMarshalFunc = trimInsolarPrefix
	initLevelMappings()
}

type zerologMapping struct {
	zl      zerolog.Level
	fn      func(*zerolog.Logger) *zerolog.Event
	metrics context.Context
}

func (v zerologMapping) IsEmpty() bool {
	return v.fn == nil
}

var zerologLevelMapping = []zerologMapping{
	insolar.NoLevel:    {zl: zerolog.NoLevel, fn: (*zerolog.Logger).Debug},
	insolar.DebugLevel: {zl: zerolog.DebugLevel, fn: (*zerolog.Logger).Debug},
	insolar.InfoLevel:  {zl: zerolog.InfoLevel, fn: (*zerolog.Logger).Info},
	insolar.WarnLevel:  {zl: zerolog.WarnLevel, fn: (*zerolog.Logger).Warn},
	insolar.ErrorLevel: {zl: zerolog.ErrorLevel, fn: (*zerolog.Logger).Error},
	insolar.FatalLevel: {zl: zerolog.FatalLevel, fn: (*zerolog.Logger).Fatal},
	insolar.PanicLevel: {zl: zerolog.PanicLevel, fn: (*zerolog.Logger).Panic},
}

var zerologReverseMapping []insolar.LogLevel

func initLevelMappings() {
	var zLevelMax zerolog.Level
	for i := range zerologLevelMapping {
		if zerologLevelMapping[i].IsEmpty() {
			continue
		}
		if zLevelMax < zerologLevelMapping[i].zl {
			zLevelMax = zerologLevelMapping[i].zl
		}
		zerologLevelMapping[i].metrics = logmetrics.GetLogLevelContext(insolar.LogLevel(i))
	}

	zerologReverseMapping = make([]insolar.LogLevel, zLevelMax+1)
	for i := range zerologReverseMapping {
		zerologReverseMapping[i] = insolar.NoLevel
	}

	for i := range zerologLevelMapping {
		if zerologLevelMapping[i].IsEmpty() {
			zerologLevelMapping[i] = zerologLevelMapping[insolar.NoLevel]
		} else {
			zl := zerologLevelMapping[i].zl
			if zerologReverseMapping[zl] != insolar.NoLevel {
				panic("duplicate level mapping")
			}
			zerologReverseMapping[zl] = insolar.LogLevel(i)
		}
	}
}

func getLevelMapping(insLevel insolar.LogLevel) zerologMapping {
	if int(insLevel) > len(zerologLevelMapping) {
		return zerologLevelMapping[insolar.NoLevel]
	}
	return zerologLevelMapping[insLevel]
}

func ToZerologLevel(insLevel insolar.LogLevel) zerolog.Level {
	return getLevelMapping(insLevel).zl
}

func FromZerologLevel(zLevel zerolog.Level) insolar.LogLevel {
	if int(zLevel) > len(zerologReverseMapping) {
		return zerologReverseMapping[zerolog.NoLevel]
	}
	return zerologReverseMapping[zLevel]
}

func selectOutput(output insolar.LogOutput) (io.WriteCloser, error) {
	switch output {
	case insolar.StdErrOutput:
		// we open a separate file handle as it will be closed, so it should not interfere with os.Stderr
		return os.NewFile(uintptr(syscall.Stderr), "/dev/stderr"), nil
	case insolar.SysLogOutput:
		return inssyslog.ConnectDefaultSyslog("insolar") // breaks dependency on windows
	default:
		return nil, errors.New("unknown output " + output.String())
	}
}

func selectFormatter(format insolar.LogFormat, output io.Writer) (io.Writer, error) {
	switch format {
	case insolar.TextFormat:
		return newDefaultTextOutput(output), nil
	case insolar.JSONFormat:
		return output, nil
	default:
		return nil, errors.New("unknown formatter " + format.String())
	}
}

const zerologSkipFrameCount = 4

func NewZerologAdapter(pCfg logadapter.ParsedLogConfig, msgFmt logadapter.MsgFormatConfig) (insolar.Logger, error) {

	bareOutput, err := selectOutput(pCfg.OutputType)
	if err != nil {
		return nil, err
	}

	sfb := zerologSkipFrameCount + pCfg.SkipFrameBaselineAdjustment
	if sfb < 0 {
		sfb = 0
	}

	zc := logadapter.Config{}
	zc.Output = pCfg.Output
	zc.MsgFormat = msgFmt
	zc.Instruments.SkipFrameCountBaseline = uint8(sfb)

	zb := logadapter.NewBuilder(zerologFactory{}, bareOutput, zc, pCfg.LogLevel)

	return zb.Build()
}

var _ insolar.Logger = &zerologAdapter{}

type zerologAdapter struct {
	logger zerolog.Logger
	config *logadapter.Config
}

// WithFields return copy of adapter with predefined fields.
func (z *zerologAdapter) WithFields(fields map[string]interface{}) insolar.Logger {
	zCtx := z.logger.With()
	for key, value := range fields {
		zCtx = zCtx.Interface(key, value)
	}

	zCopy := *z
	zCopy.logger = zCtx.Logger()
	return &zCopy
}

// WithField return copy of adapter with predefined single field.
func (z *zerologAdapter) WithField(key string, value interface{}) insolar.Logger {
	zCopy := *z
	zCopy.logger = z.logger.With().Interface(key, value).Logger()
	return &zCopy
}

func (z *zerologAdapter) newEvent(level insolar.LogLevel) *zerolog.Event {
	m := getLevelMapping(level)
	z.config.Metrics.OnNewEvent(m.metrics, level)
	event := m.fn(&z.logger)
	if event == nil {
		return nil
	}
	if z.config.DynLevel != nil && z.config.DynLevel.GetLogLevel() > level {
		return nil
	}
	z.config.Metrics.OnFilteredEvent(m.metrics, level)
	return event
}

func (z *zerologAdapter) EmbeddedEvent(level insolar.LogLevel, args ...interface{}) {
	event := z.newEvent(level)
	if event != nil { // avoid unnecessary call to fmt.Sprint
		event.Msg(z.config.MsgFormat.Sformat(args...))
	}
}

func (z *zerologAdapter) EmbeddedEventf(level insolar.LogLevel, fmt string, args ...interface{}) {
	event := z.newEvent(level)
	if event != nil { // avoid unnecessary call to fmt.Sprintf
		event.Msg(z.config.MsgFormat.Sformatf(fmt, args...))
	}
}

func (z *zerologAdapter) Event(level insolar.LogLevel, args ...interface{}) {
	z.EmbeddedEvent(level, args...)
}

func (z *zerologAdapter) Eventf(level insolar.LogLevel, fmt string, args ...interface{}) {
	z.EmbeddedEventf(level, fmt, args...)
}

func (z *zerologAdapter) Debug(args ...interface{}) {
	z.EmbeddedEvent(insolar.DebugLevel, args...)
}

func (z *zerologAdapter) Debugf(format string, args ...interface{}) {
	z.EmbeddedEventf(insolar.DebugLevel, format, args...)
}

func (z *zerologAdapter) Info(args ...interface{}) {
	z.EmbeddedEvent(insolar.InfoLevel, args...)
}

func (z *zerologAdapter) Infof(format string, args ...interface{}) {
	z.EmbeddedEventf(insolar.InfoLevel, format, args...)
}

func (z *zerologAdapter) Warn(args ...interface{}) {
	z.EmbeddedEvent(insolar.WarnLevel, args...)
}

func (z *zerologAdapter) Warnf(format string, args ...interface{}) {
	z.EmbeddedEventf(insolar.WarnLevel, format, args...)
}

func (z *zerologAdapter) Error(args ...interface{}) {
	z.EmbeddedEvent(insolar.ErrorLevel, args...)
}

func (z *zerologAdapter) Errorf(format string, args ...interface{}) {
	z.EmbeddedEventf(insolar.ErrorLevel, format, args...)
}

func (z *zerologAdapter) Fatal(args ...interface{}) {
	z.EmbeddedEvent(insolar.FatalLevel, args...)
}

func (z *zerologAdapter) Fatalf(format string, args ...interface{}) {
	z.EmbeddedEventf(insolar.FatalLevel, format, args...)
}

func (z *zerologAdapter) Panic(args ...interface{}) {
	z.EmbeddedEvent(insolar.PanicLevel, args...)
}

func (z *zerologAdapter) Panicf(format string, args ...interface{}) {
	z.EmbeddedEventf(insolar.PanicLevel, format, args...)
}

func (z *zerologAdapter) Is(level insolar.LogLevel) bool {
	return z.newEvent(level) != nil
}

func (z *zerologAdapter) Copy() insolar.LoggerBuilder {
	return logadapter.NewBuilderWithTemplate(zerologTemplate{template: z}, FromZerologLevel(z.logger.GetLevel()))
}

func (z *zerologAdapter) Level(lvl insolar.LogLevel) insolar.Logger {
	zCopy := *z
	zCopy.logger = z.logger.Level(ToZerologLevel(lvl))
	return &zCopy
}

func (z *zerologAdapter) Embeddable() insolar.EmbeddedLogger {
	return z
}

func (z *zerologAdapter) GetLoggerOutput() insolar.LoggerOutput {
	return z.config.LoggerOutput
}

/* =========================== */

var _ logadapter.Factory = &zerologFactory{}

type zerologFactory struct {
}

func (zf zerologFactory) getWriteDelayHookParams(metrics *logmetrics.MetricsHelper,
	config logadapter.BuildConfig) (needsHook bool, fieldName string, preferTrim bool, reportFn logmetrics.DurationReportFunc) {

	metricsMode := config.Instruments.MetricsMode
	if metricsMode&(insolar.LogMetricsWriteDelayField|insolar.LogMetricsWriteDelayReport) == 0 {
		return
	}

	if metricsMode&insolar.LogMetricsWriteDelayField != 0 {
		fieldName = "writeDuration"
	}

	if metricsMode&insolar.LogMetricsWriteDelayReport != 0 {
		reportFn = metrics.GetOnWriteDurationReport()
	}

	if len(fieldName) == 0 { // && reportFn == nil {
		return
	}

	return true, fieldName, false, nil
}

func (zf zerologFactory) PrepareBareOutput(output io.Writer, metrics *logmetrics.MetricsHelper, config logadapter.BuildConfig) (io.Writer, error) {
	var err error
	output, err = selectFormatter(config.Output.Format, output)

	if err != nil {
		return nil, err
	}

	if ok, name, trim, reportFn := zf.getWriteDelayHookParams(metrics, config); ok {
		output = newWriteDelayPostHook(output, name, trim, reportFn)
	}

	return output, nil
}

func (zf zerologFactory) CanReuseMsgBuffer() bool {
	// zerolog does recycling of []byte buffers
	return false
}

func (zf zerologFactory) CreateNewLowLatencyLogger(level insolar.LogLevel, config logadapter.Config) (insolar.Logger, error) {
	return zf.createNewLogger(zerologAdapterLLOutput{config.LoggerOutput}, level, config)
}

func (zf zerologFactory) CreateNewLogger(level insolar.LogLevel, config logadapter.Config) (insolar.Logger, error) {
	return zf.createNewLogger(zerologAdapterOutput{config.LoggerOutput}, level, config)
}

func (zf zerologFactory) createNewLogger(output zerolog.LevelWriter, level insolar.LogLevel, config logadapter.Config) (insolar.Logger, error) {

	ls := zerolog.New(output).Level(ToZerologLevel(level))

	if ok, name, trim, _ := zf.getWriteDelayHookParams(config.Metrics, config.BuildConfig); ok {
		// MUST be the first Hook
		ls = ls.Hook(newWriteDelayPreHook(name, trim))
	}

	lc := ls.With().Timestamp()

	skipFrames := int(config.Instruments.SkipFrameCountBaseline) + int(config.Instruments.SkipFrameCount)
	callerMode := config.Instruments.CallerMode

	if callerMode == insolar.CallerField {
		lc = lc.CallerWithSkipFrameCount(skipFrames)
	}
	ls = lc.Logger()
	if callerMode == insolar.CallerFieldWithFuncName {
		ls = ls.Hook(newCallerHook(2 + skipFrames))
	}

	if config.Instruments.MetricsMode == insolar.NoLogMetrics {
		config.Metrics = nil
	}

	return &zerologAdapter{logger: ls, config: &config}, nil
}

/* =========================== */

var _ logadapter.Template = &zerologTemplate{}

type zerologTemplate struct {
	zerologFactory
	template *zerologAdapter
}

func (zf zerologTemplate) GetLoggerOutput() insolar.LoggerOutput {
	return zf.template.GetLoggerOutput()
}

func (zf zerologTemplate) GetTemplateConfig() logadapter.Config {
	return *zf.template.config
}

func (zf zerologTemplate) GetTemplateLogger() insolar.Logger {
	return zf.template
}

/* ========================================= */

var _ zerolog.LevelWriter = &zerologAdapterOutput{}

type zerologAdapterOutput struct {
	insolar.LoggerOutput
}

func (z zerologAdapterOutput) WriteLevel(level zerolog.Level, b []byte) (int, error) {
	return z.LoggerOutput.LogLevelWrite(FromZerologLevel(level), b)
}

/* ========================================= */

var _ zerolog.LevelWriter = &zerologAdapterLLOutput{}

type zerologAdapterLLOutput struct {
	insolar.LoggerOutput
}

func (z zerologAdapterLLOutput) WriteLevel(level zerolog.Level, b []byte) (int, error) {
	return z.LoggerOutput.LowLatencyWrite(FromZerologLevel(level), b)
}

/* ========================================= */

const internalTempFieldName = "_TWD_"
const fieldHeaderFmt = `,"%s":"%*v`
const tempHexFieldLength = 16 // HEX for Uint64
const writeDelayResultFieldOverflowContent = "ovrflw"
const writeDelayResultFieldMinWidth = len(writeDelayResultFieldOverflowContent)

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

	ofs := -1
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
