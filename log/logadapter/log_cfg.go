// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
