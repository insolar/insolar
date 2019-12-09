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

package logcommon

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

type LogLevel uint8

// NoLevel means it should be ignored
const (
	NoLevel LogLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel

	LogLevelCount = int(iota)
)

func (l LogLevel) Equal(other LogLevel) bool {
	return l == other
}

func (l LogLevel) String() string {
	switch l {
	case NoLevel:
		return ""
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	}
	return ""
}

func ParseLevel(levelStr string) (LogLevel, error) {
	switch strings.ToLower(levelStr) {
	case NoLevel.String():
		return NoLevel, nil
	case DebugLevel.String():
		return DebugLevel, nil
	case InfoLevel.String():
		return InfoLevel, nil
	case WarnLevel.String():
		return WarnLevel, nil
	case ErrorLevel.String():
		return ErrorLevel, nil
	case FatalLevel.String():
		return FatalLevel, nil
	case PanicLevel.String():
		return PanicLevel, nil
	}
	return NoLevel, fmt.Errorf("unknown Level String: '%s', defaulting to NoLevel", levelStr)
}

type LoggerOutputGetter interface {
	GetLoggerOutput() LoggerOutput
}

type LoggerOutput interface {
	LogLevelWriter
	LowLatencyWrite(LogLevel, []byte) (int, error)
	IsLowLatencySupported() bool
}

// Presence of this interface indicates that this object can be used as a log event
type LogObject interface {
	// should return nil to use default (external) marshaller
	GetLogObjectMarshaller() LogObjectMarshaller
}

var _ LogObject = &LogObjectTemplate{}

type LogObjectTemplate struct{}

func (*LogObjectTemplate) GetLogObjectMarshaller() LogObjectMarshaller {
	return nil
}

type LogObjectFields struct {
	Msg    string
	Fields map[string]interface{}
}

func (v LogObjectFields) MarshalLogObject(w LogObjectWriter, _ LogObjectMetricCollector) string {
	for k, v := range v.Fields {
		w.AddIntfField(k, v, LogFieldFormat{})
	}
	return v.Msg
}

type LogObjectMarshaller interface {
	MarshalLogObject(LogObjectWriter, LogObjectMetricCollector) string
}

type MutedLogObjectMarshaller interface {
	MarshalMutedLogObject(LogObjectMetricCollector)
}

type LogObjectMetricCollector interface {
	LogObjectMetricCollector()
	//ReportMetricSample(metricType uint32, reporterFieldName string, value interface{})
}

type LogFieldFormat struct {
	Fmt    string
	Kind   reflect.Kind
	HasFmt bool
}

func (f LogFieldFormat) IsInt() bool {
	return f.Kind >= reflect.Int && f.Kind <= reflect.Int64
}

func (f LogFieldFormat) IsUint() bool {
	return f.Kind >= reflect.Uint && f.Kind <= reflect.Uintptr
}

func (f LogFieldFormat) ToString(v interface{}, defFmt string) string {
	if f.HasFmt {
		return fmt.Sprintf(f.Fmt, v)
	}
	return fmt.Sprintf(defFmt, v)
}

type LogObjectWriter interface {
	AddIntField(key string, v int64, fmt LogFieldFormat)
	AddUintField(key string, v uint64, fmt LogFieldFormat)
	AddBoolField(key string, v bool, fmt LogFieldFormat)
	AddFloatField(key string, v float64, fmt LogFieldFormat)
	AddComplexField(key string, v complex128, fmt LogFieldFormat)
	AddStrField(key string, v string, fmt LogFieldFormat)
	AddIntfField(key string, v interface{}, fmt LogFieldFormat)
	AddRawJSONField(key string, v interface{}, fmt LogFieldFormat)
}

type LogLevelWriter interface {
	io.WriteCloser
	LogLevelWrite(LogLevel, []byte) (int, error)
	Flush() error
}
