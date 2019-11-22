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
	"io"
	"strconv"
	"strings"

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

func NewZerologAdapter(pCfg logadapter.ParsedLogConfig, msgFmt logadapter.MsgFormatConfig) (insolar.Logger, error) {

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

func (m zerologMarshaller) AddField(key string, v interface{}) {
	m.event.Interface(key, v)
}

func (m zerologMarshaller) AddStrField(key string, v string) {
	m.event.Str(key, v)
}

func (m zerologMarshaller) AddRawJSON(key string, b []byte) {
	m.event.RawJSON(key, b)
}

/* ============================ */

var _ insolar.Logger = &zerologAdapter{}

type zerologAdapter struct {
	logger    zerolog.Logger
	dynFields logadapter.DynFieldMap
	config    *logadapter.Config
}

func (z *zerologAdapter) WithFields(fields map[string]interface{}) insolar.Logger {
	zCtx := z.logger.With()
	for key, value := range fields {
		zCtx = zCtx.Interface(key, value)
	}

	zCopy := *z
	zCopy.logger = zCtx.Logger()
	return &zCopy
}

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
	z.config.Metrics.OnFilteredEvent(m.metrics, level)
	return event
}

func (z *zerologAdapter) EmbeddedEvent(level insolar.LogLevel, args ...interface{}) {
	event := z.newEvent(level)

	if event == nil {
		collector := z.config.Metrics.GetMetricsCollector()
		if collector != nil {
			if obj := z.config.MsgFormat.PrepareMutedLogObject(args...); obj != nil {
				obj.MarshalMutedLogObject(collector)
			}
		}
		return
	}

	obj, msgStr := z.config.MsgFormat.FmtLogObject(args...)
	if obj != nil {
		collector := z.config.Metrics.GetMetricsCollector()
		msgStr = obj.MarshalTextLogObject(zerologMarshaller{event}, collector)
	}
	event.Msg(msgStr)
}

func (z *zerologAdapter) EmbeddedEventf(level insolar.LogLevel, fmt string, args ...interface{}) {
	event := z.newEvent(level)
	if event == nil {
		return
	}
	event.Msg(z.config.MsgFormat.Sformatf(fmt, args...))
}

func (z *zerologAdapter) EmbeddedFlush(msg string) {
	if len(msg) > 0 {
		z.newEvent(insolar.WarnLevel).Msg(msg)
	}
	_ = z.config.LoggerOutput.Flush()
}

func (z *zerologAdapter) EmbeddedTrace() insolar.LogLevel {
	return insolar.DebugLevel
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
var _ insolar.GlobalLogAdapterFactory = &zerologFactory{}

type zerologFactory struct {
	writeDelayPreferTrim bool
}

func (zf zerologFactory) CreateGlobalLogAdapter() insolar.GlobalLogAdapter {
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
) (insolar.Logger, error) {

	instruments := params.Config.Instruments
	skipFrames := int(instruments.SkipFrameCountBaseline) + int(instruments.SkipFrameCount)
	callerMode := instruments.CallerMode

	cfg := params.Config
	if instruments.MetricsMode == insolar.NoLogMetrics {
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

	if callerMode == insolar.CallerFieldWithFuncName {
		la.logger = la.logger.Hook(newCallerHook(2 + skipFrames))
	}
	lc := la.logger.With()

	// only add hooks, DON'T set the context as it can be replaced below
	lc = lc.Timestamp()
	if callerMode == insolar.CallerField {
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

func (zf zerologFactory) copyLogger(template *zerologAdapter, params logadapter.CopyLoggerParams) insolar.Logger {

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

func (zf zerologFactory) CreateNewLogger(params logadapter.NewLoggerParams) (insolar.Logger, error) {
	output := zf.createOutputWrapper(params.Config, params.Reqs)
	return zf.createNewLogger(output, params, nil)
}

func (zf zerologFactory) CanReuseMsgBuffer() bool {
	// zerolog does recycling of []byte buffers
	return false
}

/* =========================== */

var zerologGlobalAdapter insolar.GlobalLogAdapter = &zerologGlobal{}

type zerologGlobal struct {
}

func (zerologGlobal) SetGlobalLoggerFilter(level insolar.LogLevel) {
	zerolog.SetGlobalLevel(ToZerologLevel(level))
}

func (zerologGlobal) GetGlobalLoggerFilter() insolar.LogLevel {
	return FromZerologLevel(zerolog.GlobalLevel())
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

func (zf zerologTemplate) CopyTemplateLogger(params logadapter.CopyLoggerParams) insolar.Logger {
	return zf.copyLogger(zf.template, params)
}

func (zf zerologTemplate) CreateNewLogger(params logadapter.NewLoggerParams) (insolar.Logger, error) {
	output := zf.createOutputWrapper(params.Config, params.Reqs)
	return zf.createNewLogger(output, params, zf.template)
}

/* ========================================= */

var _ zerolog.LevelWriter = &zerologAdapterOutput{}

type zerologAdapterOutput struct {
	insolar.LoggerOutput
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
	insolar.LoggerOutput
}

func (z zerologAdapterLLOutput) WriteLevel(level zerolog.Level, b []byte) (int, error) {
	return z.LoggerOutput.LowLatencyWrite(FromZerologLevel(level), b)
}

func (z zerologAdapterLLOutput) Write(b []byte) (int, error) {
	panic("unexpected") // zerolog writes only to WriteLevel
}

/* ========================================= */

func newDynFieldsHook(dynFields logadapter.DynFieldMap) zerolog.Hook {
	return dynamicFieldsHook{dynFields}
}

type dynamicFieldsHook struct {
	dynFields logadapter.DynFieldMap
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
