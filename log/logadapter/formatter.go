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
	"fmt"
	"reflect"

	"github.com/insolar/insolar/log/logcommon"
)

type FormatFunc func(...interface{}) string
type FormatfFunc func(string, ...interface{}) string

type MsgFormatConfig struct {
	Sformat  FormatFunc
	Sformatf FormatfFunc
	MFactory MarshallerFactory
}

func GetDefaultLogMsgFormatter() MsgFormatConfig {
	return MsgFormatConfig{
		Sformat:  fmt.Sprint,
		Sformatf: fmt.Sprintf,
		MFactory: GetDefaultLogMsgMarshallerFactory(),
	}
}

func (v MsgFormatConfig) fmtLogStruct(a interface{}) (logcommon.LogObjectMarshaller, *string) {
	switch vv := a.(type) {
	case logcommon.LogObject:
		m := vv.GetLogObjectMarshaller()
		if m != nil {
			return m, nil
		}
		vr := reflect.ValueOf(vv)
		if vr.Kind() == reflect.Ptr {
			vr = vr.Elem()
		}
		return v.MFactory.CreateLogObjectMarshaller(vr), nil
	case logcommon.LogObjectMarshaller:
		return vv, nil
	case string: // the most obvious case
		return nil, &vv
	case *string: // handled separately to avoid unnecessary reflect
		return nil, vv
	case nil:
		return nil, nil
	default:
		if s, t, isNil := prepareValue(vv); t != reflect.Invalid {
			if isNil {
				return nil, nil
			}
			return nil, &s
		}

		vr := reflect.ValueOf(vv)
		switch vr.Kind() {
		case reflect.Ptr:
			if vr.IsNil() {
				break
			}
			vr = vr.Elem()
			if vr.Kind() != reflect.Struct {
				break
			}
			fallthrough
		case reflect.Struct:
			if len(vr.Type().Name()) == 0 { // only unnamed objects are handled by default
				return v.MFactory.CreateLogObjectMarshaller(vr), nil
			}
		}
	}
	return nil, nil
}

func (v MsgFormatConfig) FmtLogStruct(a interface{}) (logcommon.LogObjectMarshaller, string) {
	if m, s := v.fmtLogStruct(a); m != nil {
		return m, ""
	} else if s != nil {
		return m, *s
	}
	return nil, v.Sformat(a)
}

func (v MsgFormatConfig) FmtLogStructOrObject(a interface{}) (logcommon.LogObjectMarshaller, string) {
	if m, s := v.fmtLogStruct(a); m != nil {
		return m, ""
	} else if s != nil {
		return m, *s
	}
	return nil, v.Sformat(a)
}

func (v MsgFormatConfig) FmtLogObject(a ...interface{}) string {
	return v.Sformat(a...)
}

func (v MsgFormatConfig) PrepareMutedLogObject(a ...interface{}) logcommon.LogObjectMarshaller {
	if len(a) != 1 {
		return nil
	}

	switch vv := a[0].(type) {
	case logcommon.LogObject:
		m := vv.GetLogObjectMarshaller()
		if m != nil {
			return m
		}
		vr := reflect.ValueOf(vv)
		if vr.Kind() == reflect.Ptr {
			vr = vr.Elem()
		}
		return v.MFactory.CreateLogObjectMarshaller(vr)
	case logcommon.LogObjectMarshaller:
		return vv
	case string, *string, nil: // the most obvious case(s) - avoid reflect
		return nil
	default:
		vr := reflect.ValueOf(vv)
		switch vr.Kind() {
		case reflect.Ptr:
			if vr.IsNil() {
				return nil
			}
			vr = vr.Elem()
			if vr.Kind() != reflect.Struct {
				return nil
			}
			fallthrough
		case reflect.Struct:
			if len(vr.Type().Name()) == 0 { // only unnamed objects are handled by default
				return v.MFactory.CreateLogObjectMarshaller(vr)
			}
		}
	}
	return nil
}
