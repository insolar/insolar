//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package log

import (
	"context"
	"fmt"
	"github.com/insolar/insolar/log/critlog"
	"github.com/insolar/insolar/log/inssyslog"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"io"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/insolar/insolar/configuration"
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
	zerolog.TimeFieldFormat = timestampFormat
	zerolog.CallerMarshalFunc = trimInsolarPrefix
	initLevelMappings()
}

var _ insolar.Logger = &zerologAdapter{}

type zerologAdapter struct {
	logger      zerolog.Logger
	output      io.Writer
	outputWraps outputWrapFlag
	config      zerologAdapterConfig
}

type zerologAdapterConfig struct {
	bareOutput io.Writer

	format     insolar.LogFormat
	bufferSize int

	callerMode             insolar.CallerFieldMode
	skipFrameCountBaseline int
	skipFrameCount         int

	dynLevel     insolar.LogLevelReader
	pollInterval time.Duration
}

type outputWrapFlag uint32

const (
	outputWrappedWithBuffer outputWrapFlag = 1 << iota
	outputWrappedWithCritical
	outputWrappedWithFormatter
)

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
		ctx, err := tag.New(context.Background(), tag.Insert(tagLevel, insolar.LogLevel(i).String()))
		if err != nil {
			panic(err)
		}
		zerologLevelMapping[i].metrics = ctx
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

const zerologSkipFrameCount = 3 + defaultCallerSkipFrameCount

func newZerologAdapter(cfg configuration.Log, pCfg parsedLogConfig) (insolar.Logger, error) {

	zb := zerologBuilder{
		level: ToZerologLevel(pCfg.LogLevel),
		zerologAdapterConfig: zerologAdapterConfig{
			format:                 pCfg.LogFormat,
			bufferSize:             cfg.BufferSize,
			skipFrameCountBaseline: zerologSkipFrameCount + pCfg.SkipFrameBaselineDelta,
		},
	}

	var err error
	zb.bareOutput, err = selectOutput(pCfg.OutputType)
	if err != nil {
		return nil, err
	}

	return zb.Build()
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
	stats.Record(m.metrics, statLogCalls.M(1))
	event := m.fn(&z.logger)
	if event == nil {
		return nil
	}
	if z.config.dynLevel != nil && z.config.dynLevel.GetLogLevel() > level {
		return nil
	}
	stats.Record(m.metrics, statLogWrites.M(1))
	return event
}

func (z *zerologAdapter) EmbeddedEvent(level insolar.LogLevel, args ...interface{}) {
	z.newEvent(level).Msg(fmt.Sprint(args...))
}

func (z *zerologAdapter) EmbeddedEventf(level insolar.LogLevel, fmt string, args ...interface{}) {
	z.newEvent(level).Msgf(fmt, args...)
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
	return zerologBuilder{z, z.logger.GetLevel(), z.config}
}

func (z *zerologAdapter) Level(lvl insolar.LogLevel) insolar.Logger {
	zCopy := *z
	zCopy.logger = z.logger.Level(ToZerologLevel(lvl))
	return &zCopy
}

func (z *zerologAdapter) Embeddable() insolar.EmbeddedLogger {
	return z
}

/* =========================== */

type zerologBuilder struct {
	template *zerologAdapter
	level    zerolog.Level
	zerologAdapterConfig
}

func (z zerologBuilder) GetOutput() io.Writer {
	return z.bareOutput
}

func (z zerologBuilder) GetLogLevel() insolar.LogLevel {
	return FromZerologLevel(z.level)
}

func (z zerologBuilder) WithOutput(w io.Writer) insolar.LoggerBuilder {
	z.bareOutput = w
	return z
}

func (z zerologBuilder) WithLevel(level insolar.LogLevel) insolar.LoggerBuilder {
	z.level = getLevelMapping(level).zl
	z.dynLevel = nil
	return z
}

func (z zerologBuilder) WithDynamicLevel(level insolar.LogLevelReader) insolar.LoggerBuilder {
	z.dynLevel = level
	z.level = zerolog.DebugLevel
	return z
}

func (z zerologBuilder) WithFormat(format insolar.LogFormat) insolar.LoggerBuilder {
	z.format = format
	return z
}

func (z zerologBuilder) WithCaller(mode insolar.CallerFieldMode) insolar.LoggerBuilder {
	z.callerMode = mode
	return z
}

func (z zerologBuilder) WithSkipFrameCount(skipFrameCount int) insolar.LoggerBuilder {
	z.skipFrameCount = skipFrameCount
	return z
}

//func (z zerologBuilder) WithSkipFrameCount(delta int) insolar.LoggerBuilder {
//	z.skipFrameCount = delta
//	return z
//}

func (z zerologBuilder) Build() (insolar.Logger, error) {
	zNew := zerologAdapter{config: z.zerologAdapterConfig}

	switch {
	case z.bareOutput == nil:
		return nil, errors.New("output is nil")
	case z.template != nil && z.template.config.bareOutput == z.bareOutput:
		if z.template.config == z.zerologAdapterConfig {
			return z.template, nil
		}
		if z.template.config.format == z.format {
			zNew.output = z.template.output
			break
		}
		fallthrough
	default:
		err := zNew.prepareOutput()
		if err != nil {
			return nil, err
		}
	}

	zNew.prepareLogger(z.level)

	return &zNew, nil
}

func (z *zerologAdapter) prepareOutput() error {
	var err error
	z.output, err = selectFormatter(z.config.format, z.config.bareOutput)
	if err != nil {
		return err
	}
	if z.output != z.config.bareOutput {
		z.outputWraps |= outputWrappedWithFormatter
	}

	if z.config.bufferSize > 0 {
		dropBufOnFatal := z.config.bufferSize > 1000

		z.output = NewDiodeBufferedLevelWriter(z.output, z.config.bufferSize,
			z.config.pollInterval, dropBufOnFatal,
			func(missed int) []byte {
				return ([]byte)(fmt.Sprintf("logger dropped %d messages", missed))
			},
		)
		z.outputWraps |= outputWrappedWithBuffer | outputWrappedWithCritical
	} else {
		z.output = critlog.NewFatalDirectWriter(z.output)
	}

	return nil
}

func (z *zerologAdapter) prepareLogger(level zerolog.Level) {

	ls := zerolog.New(z.output).Level(level)
	if z.config.callerMode == insolar.CallerFieldWithFuncName {
		ls = ls.Hook(newCallerHook(z.config.skipFrameCountBaseline + 2 + z.config.skipFrameCount))
	}

	lc := ls.With().Timestamp()

	if z.config.callerMode == insolar.CallerField {
		lc = lc.CallerWithSkipFrameCount(z.config.skipFrameCountBaseline + z.config.skipFrameCount)
	}

	z.logger = lc.Logger()
}
