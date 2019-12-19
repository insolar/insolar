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

package linsadapter

import (
	"fmt"
	"os"
	"sort"

	"github.com/insolar/insolar/log/logadapter"
	"github.com/insolar/insolar/log/logcommon"
)

var _ logcommon.EmbeddedLogger = &logInsAdapter{}

//var _ logcommon.EmbeddedLoggerAssistant = &logInsAdapter{}

type logInsAdapter struct {
	config  *logadapter.Config
	encoder logcommon.EncoderManager
	writer  logcommon.LogLevelWriter

	parentStatic  *[]byte
	staticFields  []byte
	dynamicHooks  map[string]logcommon.LogObjectMarshallerFunc
	dynamicFields logcommon.DynFieldMap
}

func (v logInsAdapter) sendEvent(level logcommon.LogLevel, event logcommon.LogObjectWriter, msg string) {
	var eventBuf, extraBuf []byte

	switch {
	case len(v.dynamicFields) > 0 || len(v.dynamicHooks) > 0:
		extra := v.encoder.CreatePartEncoder(nil)

		for _, hook := range v.dynamicHooks {
			// TODO handle panic
			hook(extra)
		}

		if len(v.dynamicFields) == 1 {
			for key, val := range v.dynamicFields {
				addDynField(extra, key, val)
			}
		} else {
			keys := make([]string, 0, len(v.dynamicFields))
			for key := range v.dynamicFields {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			for _, key := range keys {
				addDynField(extra, key, v.dynamicFields[key])
			}
		}

		if event == nil {
			// avoids unnecessary allocation of another writer
			extra.AddStrField("Msg", msg, logcommon.LogFieldFormat{})
		} else {
			event.AddStrField("Msg", msg, logcommon.LogFieldFormat{})
			eventBuf = v.encoder.FlushPartEncoder(event)
		}
		extraBuf = v.encoder.FlushPartEncoder(extra)
	case event == nil:
		event = v.encoder.CreatePartEncoder(nil)
		fallthrough
	default:
		event.AddStrField("Msg", msg, logcommon.LogFieldFormat{})
		eventBuf = v.encoder.FlushPartEncoder(event)
	}

	switch err := v._sendEvent(level, v.parentStatic, v.staticFields, extraBuf, eventBuf); {
	case err == nil:
	case v.config.ErrorFn != nil:
		v.config.ErrorFn(err)
	default:
		_, _ = fmt.Fprintf(os.Stderr, "inslog: could not write event: %v\n", err)
	}
}

func (v logInsAdapter) _sendEvent(level logcommon.LogLevel, parent *[]byte, static, extra, event []byte) error {
	if parent == nil {
		return v.encoder.WriteParts(level, [][]byte{static, extra, event}, v.writer)
	}
	return v.encoder.WriteParts(level, [][]byte{*parent, static, extra, event}, v.writer)
}

func addDynField(w logcommon.LogObjectWriter, key string, valFn logcommon.DynFieldFunc) {
	// TODO handle panic
	val := valFn()
	w.AddIntfField(key, val, logcommon.LogFieldFormat{})
}

func (v logInsAdapter) NewEventStruct(level logcommon.LogLevel) func(interface{}) {
	if !v.Is(level) {
		return nil
	}
	return func(arg interface{}) {
		obj, msgStr := v.config.MsgFormat.FmtLogStruct(arg)
		if obj == nil {
			v.sendEvent(level, nil, msgStr)
			return
		}

		event := v.encoder.CreatePartEncoder(nil)
		collector := v.config.Metrics.GetMetricsCollector()
		msgStr = obj.MarshalLogObject(event, collector)
		v.sendEvent(level, event, msgStr)
	}
}

func (v logInsAdapter) NewEvent(level logcommon.LogLevel) func(args []interface{}) {
	if !v.Is(level) {
		return nil
	}
	return func(args []interface{}) {
		if len(args) != 1 {
			msgStr := v.config.MsgFormat.FmtLogObject(args...)
			v.sendEvent(level, nil, msgStr)
			return
		}

		obj, msgStr := v.config.MsgFormat.FmtLogStructOrObject(args[0])
		if obj == nil {
			v.sendEvent(level, nil, msgStr)
			return
		}

		event := v.encoder.CreatePartEncoder(nil)
		collector := v.config.Metrics.GetMetricsCollector()
		msgStr = obj.MarshalLogObject(event, collector)
		v.sendEvent(level, event, msgStr)
	}
}

func (v logInsAdapter) NewEventFmt(level logcommon.LogLevel) func(fmt string, args []interface{}) {
	if !v.Is(level) {
		return nil
	}
	return func(fmt string, args []interface{}) {
		v.sendEvent(level, nil, v.config.MsgFormat.Sformatf(fmt, args...))
	}
}

func (v logInsAdapter) EmbeddedFlush(msg string) {
	if len(msg) > 0 {
		v.sendEvent(logcommon.WarnLevel, nil, msg)
	}
	_ = v.config.LoggerOutput.Flush()
}

func (v logInsAdapter) Is(level logcommon.LogLevel) bool {
	return v.config.LevelFn(level)
}

func (v logInsAdapter) Copy() logcommon.LoggerBuilder {
	panic("implement me")
}

func (v logInsAdapter) WithFieldMarshaller(fields logcommon.LogObjectMarshallerFunc) logcommon.Logger {
	panic("implement me")
}
