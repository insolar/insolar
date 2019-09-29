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

package logadapter

import (
	"strconv"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
)

const TimestampFormat = "2006-01-02T15:04:05.000000000Z07:00"

type ParsedLogConfig struct {
	OutputType insolar.LogOutput
	LogLevel   insolar.LogLevel
	//GlobalLevel insolar.LogLevel

	OutputParam string

	Output      OutputConfig
	Instruments InstrumentationConfig

	SkipFrameBaselineAdjustment int8
}

const defaultLowLatencyBufferSize = 100

func DefaultLoggerSettings() ParsedLogConfig {
	r := ParsedLogConfig{}
	r.Instruments.MetricsMode = insolar.LogMetricsEventCount | insolar.LogMetricsWriteDelayReport | insolar.LogMetricsWriteDelayField
	r.Instruments.CallerMode = insolar.CallerField
	return r
}

func ParseLogConfig(cfg configuration.Log) (plc ParsedLogConfig, err error) {
	return ParseLogConfigWithDefaults(cfg, DefaultLoggerSettings())
}

func ParseLogConfigWithDefaults(cfg configuration.Log, defaults ParsedLogConfig) (plc ParsedLogConfig, err error) {
	plc = defaults

	plc.OutputType, err = insolar.ParseOutput(cfg.OutputType, insolar.DefaultLogOutput)
	if err != nil {
		return
	}
	plc.OutputParam = cfg.OutputParams

	plc.Output.Format, err = insolar.ParseFormat(cfg.Formatter, insolar.DefaultLogFormat)
	if err != nil {
		return
	}

	plc.LogLevel, err = insolar.ParseLevel(cfg.Level)
	if err != nil {
		return
	}

	if len(cfg.OutputParallelLimit) > 0 {
		plc.Output.ParallelWriters, err = strconv.Atoi(cfg.OutputParallelLimit)
		if err != nil {
			return
		}
	} else {
		plc.Output.ParallelWriters = 0
	}

	//plc.GlobalLevel, err = insolar.ParseLevel(cfg.GlobalLevel)
	//if err != nil {
	//	plc.GlobalLevel = insolar.NoLevel
	//}

	switch {
	case cfg.LLBufferSize < 0:
		// LL buffer is disabled
		plc.Output.BufferSize = cfg.BufferSize
	case cfg.LLBufferSize > 0:
		plc.Output.BufferSize = cfg.LLBufferSize
	default:
		plc.Output.BufferSize = defaultLowLatencyBufferSize
	}

	if plc.Output.BufferSize < cfg.BufferSize {
		plc.Output.BufferSize = cfg.BufferSize
	}
	plc.Output.EnableRegularBuffer = cfg.BufferSize > 0

	return plc, nil
}
