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
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log/critlog"
	"github.com/insolar/insolar/log/logmetrics"
	"io"
	"math"
)

type Config struct {
	BuildConfig
	MsgFormat    MsgFormatConfig
	LoggerOutput insolar.LoggerOutput
	Metrics      *logmetrics.MetricsHelper
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
	EnableBuffer   bool
	EnableFairness bool
}

func (v OutputConfig) CanBeReusedAs(config OutputConfig) bool {
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

type Factory interface {
	PrepareBareOutput(output io.Writer, metrics *logmetrics.MetricsHelper, config BuildConfig) (io.Writer, error)

	CreateNewLowLatencyLogger(level insolar.LogLevel, config Config) (insolar.Logger, error)
	CreateNewLogger(level insolar.LogLevel, config Config) (insolar.Logger, error)

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
		bareOutput:  config.LoggerOutput.GetBareOutput(),
		level:       level,
		Config:      config,
	}
}

func NewBuilder(factory Factory, bareOutput io.Writer, config Config, level insolar.LogLevel) LoggerBuilder {
	return LoggerBuilder{
		factory:    factory,
		bareOutput: bareOutput,
		level:      level,
		Config:     config,
	}
}

type LoggerBuilder struct {
	factory     Factory
	bareOutput  io.Writer
	hasTemplate bool
	level       insolar.LogLevel
	Config
}

func (z LoggerBuilder) GetOutput() io.Writer {
	return z.bareOutput
}

func (z LoggerBuilder) GetLoggerOutput() insolar.LoggerOutput {
	return z.Config.LoggerOutput
}

func (z LoggerBuilder) GetLogLevel() insolar.LogLevel {
	return z.level
}

func (z LoggerBuilder) WithOutput(w io.Writer) insolar.LoggerBuilder {
	z.bareOutput = w
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
	z.Instruments.MetricsMode = mode
	return z
}

func (z LoggerBuilder) WithMetricsRecorder(recorder insolar.LogMetricsRecorder) LoggerBuilder {
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
	case z.bareOutput == nil:
		return nil, errors.New("output is nil")
	case z.hasTemplate:
		template := z.factory.(Template)
		origConfig := template.GetTemplateConfig()

		if origConfig.BuildConfig == z.Config.BuildConfig {
			if z.bareOutput == origConfig.LoggerOutput || z.bareOutput == origConfig.LoggerOutput.GetBareOutput() {
				if needsLowLatency && !origConfig.LoggerOutput.IsLowLatencySupported() {
					break
				}
				return template.GetTemplateLogger(), nil
			}
		}
		if lo, ok := z.bareOutput.(insolar.LoggerOutput); ok {
			output = lo
			break
		}
		if z.bareOutput == origConfig.LoggerOutput.GetBareOutput() && origConfig.Output.CanBeReusedAs(z.Output) {
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
	if needsLowLatency {
		return z.factory.CreateNewLowLatencyLogger(z.level, z.Config)
	}
	return z.factory.CreateNewLogger(z.level, z.Config)
}

func (z LoggerBuilder) prepareOutput(metrics *logmetrics.MetricsHelper, needsLowLatency bool) (insolar.LoggerOutput, error) {

	output, err := z.factory.PrepareBareOutput(z.bareOutput, metrics, z.Config.BuildConfig)
	if err != nil {
		return nil, err
	}

	if z.Config.Output.BufferSize > 0 {
		var flags critlog.BackpressureBufferFlags

		if z.Config.Output.BufferSize > 100 {
			flags |= critlog.BufferDropOnFatal
		}

		if z.factory.CanReuseMsgBuffer() {
			flags |= critlog.BufferReuse
		}

		if z.Config.Output.EnableFairness {
			flags |= critlog.BufferWriteDelayFairness | critlog.BufferTrackWriteDuration
		}

		pw := uint8(5)
		if z.Config.Output.ParallelWriters > 0 && z.Config.Output.ParallelWriters <= math.MaxInt8 {
			pw = uint8(z.Config.Output.ParallelWriters)
		}

		if !z.Config.Output.EnableBuffer {
			flags |= critlog.BufferDirectForRegular
		}

		missedFn := z.loggerMissedEvent(insolar.WarnLevel)

		bpb := critlog.NewBackpressureBuffer(output,
			z.Config.Output.BufferSize, 0, pw, flags, missedFn)
		bpb.StartWorker(context.Background())
		return bpb, nil
	} else {
		if needsLowLatency {
			return nil, errors.New("low latency buffer is disabled but was required")
		}
	}
	return critlog.NewFatalDirectWriter(output), nil
}

func (z LoggerBuilder) loggerMissedEvent(level insolar.LogLevel) critlog.MissedEventFunc {
	return func(missed int) (level insolar.LogLevel, bytes []byte) {
		return level, ([]byte)(
			fmt.Sprintf(`{"level":"%v","message":"logger dropped %d messages"}`, level.String(), missed))
	}
}
