//
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
//

package zlogadapter

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/insolar/insolar/log/logcommon"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log/logadapter"
	"github.com/insolar/insolar/log/logmetrics"
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
	zerolog.TimeFieldFormat = insolar.TimestampFormat
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
	insolar.NoLevel:    {zl: zerolog.NoLevel, fn: (*zerolog.Logger).Log},
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

func NewZerologAdapter(pCfg insolar.ParsedLogConfig, msgFmt logadapter.MsgFormatConfig) (insolar.Logger, error) {

	zc := logadapter.Config{}

	var err error
	zc.BareOutput, err = logadapter.OpenLogBareOutput(pCfg.OutputType, pCfg.OutputParam)
	if err != nil {
		return nil, err
	}
	if zc.BareOutput.Writer == nil {
		return nil, errors.New("output is nil")
	}

	sfb := zerologSkipFrameCount + pCfg.SkipFrameBaselineAdjustment
	if sfb < 0 {
		sfb = 0
	}

	zc.Output = pCfg.Output
	zc.Instruments = pCfg.Instruments
	zc.MsgFormat = msgFmt
	zc.Instruments.SkipFrameCountBaseline = uint8(sfb)
	//zc.TraceLevel = insolar.InfoLevel

	zb := logadapter.NewBuilder(zerologFactory{}, zc, pCfg.LogLevel)

	return zb.Build()
}

/* ============================ */

type zerologMarshaller struct {
	event *zerolog.Event
}

func (m zerologMarshaller) AddIntField(key string, v int64, fFmt logcommon.LogFieldFormat) {
	if fFmt.HasFmt {
		m.event.Str(key, fmt.Sprintf(fFmt.Fmt, v))
	} else {
		m.event.Int64(key, v)
	}
}

func (m zerologMarshaller) AddUintField(key string, v uint64, fFmt logcommon.LogFieldFormat) {
	if fFmt.HasFmt {
		m.event.Str(key, fmt.Sprintf(fFmt.Fmt, v))
	} else {
		m.event.Uint64(key, v)
	}
}

func (m zerologMarshaller) AddBoolField(key string, v bool, fFmt logcommon.LogFieldFormat) {
	if fFmt.HasFmt {
		m.event.Str(key, fmt.Sprintf(fFmt.Fmt, v))
	} else {
		m.event.Bool(key, v)
	}
}

func (m zerologMarshaller) AddFloatField(key string, v float64, fFmt logcommon.LogFieldFormat) {
	if fFmt.HasFmt {
		m.event.Str(key, fmt.Sprintf(fFmt.Fmt, v))
	} else {
		m.event.Float64(key, v)
	}
}

func (m zerologMarshaller) AddComplexField(key string, v complex128, fFmt logcommon.LogFieldFormat) {
	if fFmt.HasFmt {
		m.event.Str(key, fmt.Sprintf(fFmt.Fmt, v))
	} else {
		m.event.Str(key, fmt.Sprint(v))
	}
}

func (m zerologMarshaller) AddStrField(key string, v string, fFmt logcommon.LogFieldFormat) {
	if fFmt.HasFmt {
		m.event.Str(key, fmt.Sprintf(fFmt.Fmt, v))
	} else {
		m.event.Str(key, v)
	}
}

func (m zerologMarshaller) AddIntfField(key string, v interface{}, fFmt logcommon.LogFieldFormat) {
	if fFmt.HasFmt {
		m.event.Str(key, fmt.Sprintf(fFmt.Fmt, v))
	} else {
		m.event.Interface(key, v)
	}
}

func (m zerologMarshaller) AddRawJSONField(key string, v interface{}, fFmt logcommon.LogFieldFormat) {
	m.event.RawJSON(key, []byte(fmt.Sprintf(fFmt.Fmt, v)))
}

/* ============================ */

var _ logcommon.EmbeddedLogger = &zerologAdapter{}
var _ logcommon.EmbeddedLoggerAssistant = &zerologAdapter{}

type zerologAdapter struct {
	logger    zerolog.Logger
	dynFields logcommon.DynFieldMap
	config    *logadapter.Config
}

func (z *zerologAdapter) WithFields(fields map[string]interface{}) insolar.Logger {
	zCtx := z.logger.With()
	for key, value := range fields {
		zCtx = zCtx.Interface(key, value)
	}

	zCopy := *z
	zCopy.logger = zCtx.Logger()
	return logcommon.WrapEmbeddedLogger(&zCopy)
}

func (z *zerologAdapter) WithField(key string, value interface{}) logcommon.Logger {
	zCopy := *z
	zCopy.logger = z.logger.With().Interface(key, value).Logger()
	return logcommon.WrapEmbeddedLogger(&zCopy)
}

func (z *zerologAdapter) newEvent(level logcommon.LogLevel) *zerolog.Event {
	m := getLevelMapping(level)
	z.config.Metrics.OnNewEvent(m.metrics, level)
	event := m.fn(&z.logger)
	if event == nil {
		return nil
	}
	z.config.Metrics.OnFilteredEvent(m.metrics, level)
	return event
}

func (z *zerologAdapter) NewEventStruct(level logcommon.LogLevel) func(interface{}) {
	switch event := z.newEvent(level); {
	case event == nil:
		//collector := z.config.Metrics.GetMetricsCollector()
		//if collector != nil {
		//	if obj := z.config.MsgFormat.PrepareMutedLogObject(arg); obj != nil {
		//		obj.MarshalMutedLogObject(collector)
		//	}
		//}
		return nil
	default:
		return func(arg interface{}) {
			obj, msgStr := z.config.MsgFormat.FmtLogStruct(arg)
			if obj != nil {
				collector := z.config.Metrics.GetMetricsCollector()
				msgStr = obj.MarshalLogObject(zerologMarshaller{event}, collector)
			}
			event.Msg(msgStr)
		}
	}
}

func (z *zerologAdapter) NewEvent(level logcommon.LogLevel) func(args []interface{}) {
	switch event := z.newEvent(level); {
	case event == nil:
		return nil
	default:
		return func(args []interface{}) {
			if len(args) != 1 {
				msgStr := z.config.MsgFormat.FmtLogObject(args...)
				event.Msg(msgStr)
				return
			}

			obj, msgStr := z.config.MsgFormat.FmtLogStructOrObject(args[0])
			if obj != nil {
				collector := z.config.Metrics.GetMetricsCollector()
				msgStr = obj.MarshalLogObject(zerologMarshaller{event}, collector)
			}
			event.Msg(msgStr)
		}
	}
}

func (z *zerologAdapter) NewEventFmt(level logcommon.LogLevel) func(fmt string, args []interface{}) {
	event := z.newEvent(level)
	if event == nil {
		return nil
	}
	return func(fmt string, args []interface{}) {
		event.Msg(z.config.MsgFormat.Sformatf(fmt, args...))
	}
}

func (z *zerologAdapter) EmbeddedFlush(msg string) {
	if len(msg) > 0 {
		z.newEvent(logcommon.WarnLevel).Msg(msg)
	}
	_ = z.config.LoggerOutput.Flush()
}

func (z *zerologAdapter) Is(level logcommon.LogLevel) bool {
	return z.newEvent(level) != nil
}

func (z *zerologAdapter) Copy() logcommon.LoggerBuilder {
	return logadapter.NewBuilderWithTemplate(zerologTemplate{template: z}, FromZerologLevel(z.logger.GetLevel()))
}

func (z *zerologAdapter) GetLoggerOutput() logcommon.LoggerOutput {
	return z.config.LoggerOutput
}

/* =========================== */

var _ logadapter.Factory = &zerologFactory{}
var _ logcommon.GlobalLogAdapterFactory = &zerologFactory{}

type zerologFactory struct {
	writeDelayPreferTrim bool
}

func (zf zerologFactory) CreateGlobalLogAdapter() logcommon.GlobalLogAdapter {
	return zerologGlobalAdapter
}

func (zf zerologFactory) PrepareBareOutput(output logadapter.BareOutput, metrics *logmetrics.MetricsHelper, config logadapter.BuildConfig) (io.Writer, error) {
	outputWriter, err := selectFormatter(config.Output.Format, output.Writer)

	if err != nil {
		return nil, err
	}

	if ok, name, reportFn := getWriteDelayConfig(metrics, config); ok {
		outputWriter = newWriteDelayPostHook(outputWriter, name, zf.writeDelayPreferTrim, reportFn)
	}

	return outputWriter, nil
}

func checkNewLoggerOutput(output zerolog.LevelWriter) zerolog.LevelWriter {
	if output == nil {
		panic("illegal value")
	}
	//
	return output
}

func (zf zerologFactory) createNewLogger(output zerolog.LevelWriter, params logadapter.NewLoggerParams, template *zerologAdapter,
) (logcommon.EmbeddedLogger, error) {

	instruments := params.Config.Instruments
	skipFrames := int(instruments.SkipFrameCountBaseline) + int(instruments.SkipFrameCount)
	callerMode := instruments.CallerMode

	cfg := params.Config
	if instruments.MetricsMode == logcommon.NoLogMetrics {
		cfg.Metrics = nil
	}

	la := zerologAdapter{
		// NB! We have to create a new logger and pass the context separately
		// Otherwise, zerolog will also copy hooks - which we need to get rid of some.
		logger: zerolog.New(checkNewLoggerOutput(output)).Level(ToZerologLevel(params.Level)),
		config: &cfg,
	}

	if ok, name, _ := getWriteDelayConfig(params.Config.Metrics, params.Config.BuildConfig); ok {
		la.logger = la.logger.Hook(newWriteDelayPreHook(name, zf.writeDelayPreferTrim))
	}

	{ // replacement and inheritance for dynFields
		switch newFields := params.DynFields; {
		case template != nil && params.Reqs&logadapter.RequiresParentDynFields != 0 && len(template.dynFields) > 0:
			prevFields := template.dynFields

			if len(newFields) > 0 {
				for k, v := range prevFields {
					if _, ok := newFields[k]; !ok {
						newFields[k] = v
					}
				}
			} else {
				newFields = prevFields
			}
			fallthrough
		case len(newFields) > 0:
			la.dynFields = newFields
			la.logger = la.logger.Hook(newDynFieldsHook(newFields))
		}
	}

	if callerMode == logcommon.CallerFieldWithFuncName {
		la.logger = la.logger.Hook(newCallerHook(2 + skipFrames))
	}
	lc := la.logger.With()

	// only add hooks, DON'T set the context as it can be replaced below
	lc = lc.Timestamp()
	if callerMode == logcommon.CallerField {
		lc = lc.CallerWithSkipFrameCount(skipFrames)
	}

	if template != nil && params.Reqs&logadapter.RequiresParentCtxFields != 0 {
		la.logger = lc.Logger()     // save hooks
		lc = template.logger.With() // get a copy of the inherited context
	}
	for k, v := range params.Fields {
		lc = lc.Interface(k, v)
	}

	la.logger.UpdateContext(func(zerolog.Context) zerolog.Context {
		return lc
	})

	return &la, nil
}

func (zf zerologFactory) copyLogger(template *zerologAdapter, params logadapter.CopyLoggerParams) logcommon.EmbeddedLogger {

	if params.Reqs&logadapter.RequiresParentDynFields == 0 {
		// have to reset hooks, but zerolog can't reset hooks
		// so we have to create the logger from scratch
		return nil
	}

	hasUpdates := false
	la := *template

	if newFields := params.AppendDynFields; len(newFields) > 0 {
		if prevFields := la.dynFields; len(prevFields) > 0 {
			// NB! avoid modification of newFields when nil can be returned
			for k := range newFields {
				if _, ok := prevFields[k]; ok {
					// key collision
					// have to reset hooks, but zerolog can't reset hooks
					// so we have to create the logger from scratch
					return nil
				}
			}
			for k, v := range prevFields {
				newFields[k] = v
			}
		}
		la.dynFields = newFields
		la.logger = la.logger.Hook(newDynFieldsHook(newFields))
		hasUpdates = true
	}

	newLevel := ToZerologLevel(params.Level)
	if template.logger.GetLevel() != newLevel {
		la.logger = la.logger.Level(newLevel)
		hasUpdates = true
	}

	{
		hasCtxUpdates := false
		var lc zerolog.Context

		if params.Reqs&logadapter.RequiresParentCtxFields == 0 {
			// have to reset logger context
			lc = zerolog.New(nil).With()
			hasCtxUpdates = true
		}

		if len(params.AppendFields) > 0 {
			if !hasCtxUpdates {
				lc = la.logger.With()
			}
			for k, v := range params.AppendFields {
				lc = lc.Interface(k, v)
			}
			hasCtxUpdates = true
		}

		if hasCtxUpdates {
			la.logger.UpdateContext(func(zerolog.Context) zerolog.Context {
				return lc
			})
			hasUpdates = true
		}
	}

	if !hasUpdates {
		return template
	}
	return &la
}

func (zf zerologFactory) createOutputWrapper(config logadapter.Config, reqs logadapter.FactoryRequirementFlags) zerolog.LevelWriter {
	if reqs&logadapter.RequiresLowLatency != 0 {
		return zerologAdapterLLOutput{config.LoggerOutput}
	}
	return zerologAdapterOutput{config.LoggerOutput}
}

func (zf zerologFactory) CreateNewLogger(params logadapter.NewLoggerParams) (logcommon.EmbeddedLogger, error) {
	output := zf.createOutputWrapper(params.Config, params.Reqs)
	return zf.createNewLogger(output, params, nil)
}

func (zf zerologFactory) CanReuseMsgBuffer() bool {
	// zerolog does recycling of []byte buffers
	return false
}

/* =========================== */

var zerologGlobalAdapter logcommon.GlobalLogAdapter = &zerologGlobal{}

type zerologGlobal struct {
}

func (zerologGlobal) SetGlobalLoggerFilter(level logcommon.LogLevel) {
	zerolog.SetGlobalLevel(ToZerologLevel(level))
}

func (zerologGlobal) GetGlobalLoggerFilter() logcommon.LogLevel {
	return FromZerologLevel(zerolog.GlobalLevel())
}

/* =========================== */

var _ logadapter.Template = &zerologTemplate{}

type zerologTemplate struct {
	zerologFactory
	template *zerologAdapter
}

func (zf zerologTemplate) GetLoggerOutput() logcommon.LoggerOutput {
	return zf.template.GetLoggerOutput()
}

func (zf zerologTemplate) GetTemplateConfig() logadapter.Config {
	return *zf.template.config
}

func (zf zerologTemplate) CopyTemplateLogger(params logadapter.CopyLoggerParams) logcommon.EmbeddedLogger {
	return zf.copyLogger(zf.template, params)
}

func (zf zerologTemplate) CreateNewLogger(params logadapter.NewLoggerParams) (logcommon.EmbeddedLogger, error) {
	output := zf.createOutputWrapper(params.Config, params.Reqs)
	return zf.createNewLogger(output, params, zf.template)
}

/* ========================================= */

var _ zerolog.LevelWriter = &zerologAdapterOutput{}

type zerologAdapterOutput struct {
	logcommon.LoggerOutput
}

func (z zerologAdapterOutput) WriteLevel(level zerolog.Level, b []byte) (int, error) {
	return z.LoggerOutput.LogLevelWrite(FromZerologLevel(level), b)
}

func (z zerologAdapterOutput) Write(b []byte) (int, error) {
	panic("unexpected") // zerolog writes only to WriteLevel
}

/* ========================================= */

var _ zerolog.LevelWriter = &zerologAdapterLLOutput{}

type zerologAdapterLLOutput struct {
	logcommon.LoggerOutput
}

func (z zerologAdapterLLOutput) WriteLevel(level zerolog.Level, b []byte) (int, error) {
	return z.LoggerOutput.LowLatencyWrite(FromZerologLevel(level), b)
}

func (z zerologAdapterLLOutput) Write(b []byte) (int, error) {
	panic("unexpected") // zerolog writes only to WriteLevel
}

/* ========================================= */

func newDynFieldsHook(dynFields logcommon.DynFieldMap) zerolog.Hook {
	return dynamicFieldsHook{dynFields}
}

type dynamicFieldsHook struct {
	dynFields logcommon.DynFieldMap
}

func (v dynamicFieldsHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	for k, fn := range v.dynFields {
		if fn == nil {
			continue
		}
		vv := fn()
		if vv == nil {
			continue
		}
		e.Interface(k, vv)
	}
}
