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

package insolar

import (
	"fmt"
	"strings"

	"github.com/insolar/insolar/log/logadapter"
	"github.com/insolar/insolar/log/logcommon"
)

type LogLevel = logcommon.LogLevel

const (
	NoLevel    = logcommon.NoLevel
	DebugLevel = logcommon.DebugLevel
	InfoLevel  = logcommon.InfoLevel
	WarnLevel  = logcommon.WarnLevel
	ErrorLevel = logcommon.ErrorLevel
	FatalLevel = logcommon.FatalLevel
	PanicLevel = logcommon.PanicLevel
)

func ParseLevel(levelStr string) (LogLevel, error) {
	return logcommon.ParseLevel(levelStr)
}

//go:generate minimock -i github.com/insolar/insolar/insolar.Logger -o ./ -s _mock.go -g

type Logger = logcommon.Logger

type LogFormat = logcommon.LogFormat

const (
	TextFormat = logcommon.TextFormat
	JSONFormat = logcommon.JSONFormat
)
const DefaultLogFormat = TextFormat

type LogOutput = logadapter.LogOutput

const (
	StdErrOutput = logadapter.StdErrOutput
	SysLogOutput = logadapter.SysLogOutput
)

const DefaultLogOutput = StdErrOutput

type LoggerBuilder = logcommon.LoggerBuilder

type CallerFieldMode = logcommon.CallerFieldMode

const (
	NoCallerField           = logcommon.NoCallerField
	CallerField             = logcommon.CallerField
	CallerFieldWithFuncName = logcommon.CallerFieldWithFuncName
)

const NoLogMetrics = logcommon.NoLogMetrics
const (
	LogMetricsEventCount       = logcommon.LogMetricsEventCount
	LogMetricsWriteDelayReport = logcommon.LogMetricsWriteDelayReport
	LogMetricsWriteDelayField  = logcommon.LogMetricsWriteDelayField
	LogMetricsResetMode        = logcommon.LogMetricsResetMode
)

func ParseFormat(formatStr string, defValue LogFormat) (LogFormat, error) {
	switch strings.ToLower(formatStr) {
	case "", "default":
		return defValue, nil
	case TextFormat.String():
		return TextFormat, nil
	case JSONFormat.String():
		return JSONFormat, nil
	}
	return defValue, fmt.Errorf("unknown Format: '%s', replaced with '%s'", formatStr, defValue)
}

func ParseOutput(outputStr string, defValue LogOutput) (LogOutput, error) {
	switch strings.ToLower(outputStr) {
	case "", "default":
		return defValue, nil
	case StdErrOutput.String():
		return StdErrOutput, nil
	case SysLogOutput.String():
		return SysLogOutput, nil
	}
	return defValue, fmt.Errorf("unknown Output: '%s', replaced with '%s'", outputStr, defValue)
}
