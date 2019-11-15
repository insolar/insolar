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

package logadapter

import (
	"context"
	"errors"
	"fmt"
	"github.com/insolar/insolar/log/logoutput"
	"io"
	"math"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log/critlog"
	"github.com/insolar/insolar/log/logmetrics"
)

type Config struct {
	BuildConfig

	BareOutput   BareOutput
	LoggerOutput insolar.LoggerOutput

	Metrics   *logmetrics.MetricsHelper
	MsgFormat MsgFormatConfig
}

type BareOutput struct {
	Writer         io.Writer
	FlushFn        logoutput.LogFlushFunc
	ProtectedClose bool
}

type BuildConfig struct {
	DynLevel    insolar.LogLevelGetter
	Output      OutputConfig
	Instruments InstrumentationConfig
}

type OutputConfig struct {
	BufferSize      int
	ParallelWriters int
	Format          insolar.LogFormat

	// allow buffer for regular events
	EnableRegularBuffer bool
}

func (v OutputConfig) CanReuseOutputFor(config OutputConfig) bool {
	return v.Format == config.Format &&
		(v.BufferSize > 0 || config.BufferSize <= 0)
}

type InstrumentationConfig struct {
	Recorder               insolar.LogMetricsRecorder
	MetricsMode            insolar.LogMetricsMode
	CallerMode             insolar.CallerFieldMode
	SkipFrameCountBaseline uint8
	SkipFrameCount         int8
}

const writeDelayFieldFlags = insolar.LogMetricsWriteDelayReport | insolar.LogMetricsWriteDelayField

func (v InstrumentationConfig) CanReuseOutputFor(config InstrumentationConfig) bool {
	vTWD := v.MetricsMode&writeDelayFieldFlags != 0
	cTWD := config.MetricsMode&writeDelayFieldFlags != 0

	if v.Recorder != config.Recorder {
		return !cTWD && !vTWD
	}

	return vTWD == cTWD || vTWD && !cTWD
}

type FactoryRequirementFlags uint8

const (
	RequiresLowLatency FactoryRequirementFlags = 1 << iota
	RequiresParentFields
)

type Factory interface {
	PrepareBareOutput(output BareOutput, metrics *logmetrics.MetricsHelper, config BuildConfig) (io.Writer, error)
	CreateNewLogger(level insolar.LogLevel, config Config, reqs FactoryRequirementFlags, dynFields map[string]func() interface{}) (insolar.Logger, error)
	CanReuseMsgBuffer() bool
}

type Template interface {
	Factory
	GetTemplateConfig() Config
	GetTemplateLogger() insolar.Logger
}

func NewBuilderWithTemplate(template Template, level insolar.LogLevel) LoggerBuilder {
	config := template.GetTemplateConfig()
	return LoggerBuilder{
		factory:     template,
		hasTemplate: true,
		level:       level,
		Config:      config,
	}
}

func NewBuilder(factory Factory, config Config, level insolar.LogLevel) LoggerBuilder {
	return LoggerBuilder{
		factory: factory,
		level:   level,
		Config:  config,
	}
}

var _ insolar.GlobalLogAdapterFactory = &LoggerBuilder{}

type LoggerBuilder struct {
	factory     Factory
	hasTemplate bool
	level       insolar.LogLevel
	fields      map[string]interface{}
	dynFields   map[string]func() interface{}
	Config
}

func (z LoggerBuilder) CreateGlobalLogAdapter() insolar.GlobalLogAdapter {
	if f, ok := z.factory.(insolar.GlobalLogAdapterFactory); ok {
		return f.CreateGlobalLogAdapter()
	}
	return nil
}

func (z LoggerBuilder) GetOutput() io.Writer {
	return z.BareOutput.Writer
}

func (z LoggerBuilder) GetLoggerOutput() insolar.LoggerOutput {
	return z.Config.LoggerOutput
}

func (z LoggerBuilder) GetLogLevel() insolar.LogLevel {
	return z.level
}

func (z LoggerBuilder) WithOutput(w io.Writer) insolar.LoggerBuilder {

	type flusher interface {
		Flush() error
	}
	type syncer interface {
		Sync() error
	}

	z.BareOutput = BareOutput{Writer: w}
	switch ww := w.(type) {
	case flusher:
		z.BareOutput.FlushFn = ww.Flush
	case syncer:
		z.BareOutput.FlushFn = ww.Sync
	}

	return z
}

func (z LoggerBuilder) WithBuffer(bufferSize int, bufferForAll bool) insolar.LoggerBuilder {
	z.Output.BufferSize = bufferSize
	z.Output.EnableRegularBuffer = bufferForAll
	return z
}

func (z LoggerBuilder) WithLevel(level insolar.LogLevel) insolar.LoggerBuilder {
	z.level = level
	z.DynLevel = nil
	return z
}

func (z LoggerBuilder) WithDynamicLevel(level insolar.LogLevelGetter) insolar.LoggerBuilder {
	z.DynLevel = level
	z.level = insolar.DebugLevel
	return z
}

func (z LoggerBuilder) WithFormat(format insolar.LogFormat) insolar.LoggerBuilder {
	z.Output.Format = format
	return z
}

func (z LoggerBuilder) WithCaller(mode insolar.CallerFieldMode) insolar.LoggerBuilder {
	z.Instruments.CallerMode = mode
	return z
}

func (z LoggerBuilder) WithMetrics(mode insolar.LogMetricsMode) insolar.LoggerBuilder {
	if mode&insolar.LogMetricsResetMode != 0 {
		z.Instruments.MetricsMode = 0
		mode &^= insolar.LogMetricsResetMode
	}
	z.Instruments.MetricsMode |= mode
	return z
}

func (z LoggerBuilder) WithMetricsRecorder(recorder insolar.LogMetricsRecorder) insolar.LoggerBuilder {
	z.Instruments.Recorder = recorder
	return z
}

func (z LoggerBuilder) WithSkipFrameCount(skipFrameCount int) insolar.LoggerBuilder {
	if skipFrameCount < math.MinInt8 || skipFrameCount > math.MaxInt8 {
		panic("illegal value")
	}
	z.Instruments.SkipFrameCount = int8(skipFrameCount)
	return z
}

func (z LoggerBuilder) WithFields(fields map[string]interface{}) insolar.LoggerBuilder {
	if z.fields == nil {
		z.fields = make(map[string]interface{}, len(fields))
	}
	for k, v := range fields {
		delete(z.dynFields, k)
		z.fields[k] = v
	}
	return z
}

func (z LoggerBuilder) WithField(k string, v interface{}) insolar.LoggerBuilder {
	if z.fields == nil {
		z.fields = make(map[string]interface{})
	}
	delete(z.dynFields, k)
	z.fields[k] = v
	return z
}

func (z LoggerBuilder) WithDynamicField(k string, fn func() interface{}) insolar.LoggerBuilder {
	if fn == nil {
		panic("illegal value")
	}
	if z.dynFields == nil {
		z.dynFields = make(map[string]func() interface{})
	}
	delete(z.fields, k)
	z.dynFields[k] = fn
	return z
}

func (z LoggerBuilder) Build() (insolar.Logger, error) {
	return z.build(false)
}

func (z LoggerBuilder) BuildLowLatency() (insolar.Logger, error) {
	return z.build(true)
}

func (z LoggerBuilder) build(needsLowLatency bool) (insolar.Logger, error) {

	var metrics *logmetrics.MetricsHelper

	if z.Config.Instruments.MetricsMode != insolar.NoLogMetrics {
		metrics = logmetrics.NewMetricsHelper(z.Config.Instruments.Recorder)
	}

	var output insolar.LoggerOutput

	switch {
	case z.BareOutput.Writer == nil:
		return nil, errors.New("output is nil")
	case z.hasTemplate:
		template := z.factory.(Template)
		origConfig := template.GetTemplateConfig()

		sameBareOutput := false
		switch {
		case z.BareOutput.Writer == origConfig.LoggerOutput: // users can be crazy
			fallthrough
		case z.BareOutput.Writer == origConfig.BareOutput.Writer:
			// keep the original settings if writer wasn't changed
			z.BareOutput = origConfig.BareOutput
			sameBareOutput = true
		}

		if origConfig.BuildConfig == z.Config.BuildConfig && sameBareOutput {
			// config and output are identical - we can reuse the original logger
			if needsLowLatency && !origConfig.LoggerOutput.IsLowLatencySupported() {
				// ... LL support is missing - shall not reuse the original logger
				break
			}
			return template.GetTemplateLogger(), nil
		}
		if lo, ok := z.BareOutput.Writer.(insolar.LoggerOutput); ok {
			// something strange, but we can also work this way
			output = lo
			break
		}
		if sameBareOutput &&
			origConfig.Output.CanReuseOutputFor(z.Output) &&
			origConfig.Instruments.CanReuseOutputFor(z.Instruments) {

			// same output, and it can be reused with the new settings
			output = origConfig.LoggerOutput
			break
		}
	}
	if output == nil || needsLowLatency && !output.IsLowLatencySupported() {
		var err error
		output, err = z.prepareOutput(metrics, needsLowLatency)
		if err != nil {
			return nil, err
		}
	}

	z.Config.Metrics = metrics
	z.Config.LoggerOutput = output

	requirements := FactoryRequirementFlags(0) | RequiresParentFields
	if needsLowLatency {
		requirements |= RequiresLowLatency
	}

	logger, err := z.factory.CreateNewLogger(z.level, z.Config, requirements, z.dynFields)

	if len(z.fields) > 0 && logger != nil && err == nil {
		logger = logger.WithFields(z.fields)
	}
	return logger, err
}

func (z LoggerBuilder) prepareOutput(metrics *logmetrics.MetricsHelper, needsLowLatency bool) (insolar.LoggerOutput, error) {

	outputWriter, err := z.factory.PrepareBareOutput(z.BareOutput, metrics, z.Config.BuildConfig)
	if err != nil {
		return nil, err
	}

	outputAdapter := logoutput.NewAdapter(outputWriter, z.BareOutput.ProtectedClose,
		z.BareOutput.FlushFn, z.BareOutput.FlushFn)

	if z.Config.Output.ParallelWriters < 0 || z.Config.Output.ParallelWriters > math.MaxUint8 {
		return nil, errors.New("argument ParallelWriters is out of bounds")
	}

	if z.Config.Output.BufferSize > 0 {
		if z.Config.Output.ParallelWriters > 0 && z.Config.Output.ParallelWriters*2 < z.Config.Output.BufferSize {
			// to limit write parallelism - buffer must be active
			return nil, errors.New("write parallelism limiter requires BufferSize >= ParallelWriters*2 ")
		}

		flags := critlog.BufferWriteDelayFairness | critlog.BufferTrackWriteDuration

		if z.Config.Output.BufferSize > 1000 {
			flags |= critlog.BufferDropOnFatal
		}

		if z.factory.CanReuseMsgBuffer() {
			flags |= critlog.BufferReuse
		}

		missedFn := z.loggerMissedEvent(insolar.WarnLevel, metrics)

		var bpb *critlog.BackpressureBuffer
		switch {
		case z.Config.Output.EnableRegularBuffer:
			pw := uint8(insolar.DefaultOutputParallelLimit)
			if z.Config.Output.ParallelWriters != 0 {
				pw = uint8(z.Config.Output.ParallelWriters)
			}
			bpb = critlog.NewBackpressureBuffer(outputAdapter, z.Config.Output.BufferSize, pw, flags, missedFn)
		case z.Config.Output.ParallelWriters == 0 || z.Config.Output.ParallelWriters == math.MaxUint8:
			bpb = critlog.NewBackpressureBufferWithBypass(outputAdapter, z.Config.Output.BufferSize,
				0 /* no limit */, flags, missedFn)
		default:
			bpb = critlog.NewBackpressureBufferWithBypass(outputAdapter, z.Config.Output.BufferSize,
				uint8(z.Config.Output.ParallelWriters), flags, missedFn)
		}

		bpb.StartWorker(context.Background())
		return bpb, nil
	}

	if needsLowLatency {
		return nil, errors.New("low latency buffer was disabled but is required")
	}

	fdw := critlog.NewFatalDirectWriter(outputAdapter)
	return fdw, nil
}

func (z LoggerBuilder) loggerMissedEvent(level insolar.LogLevel, metrics *logmetrics.MetricsHelper) critlog.MissedEventFunc {
	return func(missed int) (insolar.LogLevel, []byte) {
		metrics.OnWriteSkip(missed)
		return level, ([]byte)(
			fmt.Sprintf(`{"level":"%v","message":"logger dropped %d messages"}`, level.String(), missed))
	}
}
