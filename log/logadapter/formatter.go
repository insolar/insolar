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

	"github.com/insolar/insolar/insolar"
)

type LogStringer interface {
	LogString() string
}

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

type FieldReporterFunc func(fieldName string, v interface{})

type MarshallerFactory interface {
	CreateLogObjectMarshaller(o reflect.Value) insolar.LogObjectMarshaller
	RegisterFieldReporter(fieldType reflect.Type, fn FieldReporterFunc)
}

func GetDefaultLogMsgMarshallerFactory() MarshallerFactory {
	return marshallerFactory
}

func (v MsgFormatConfig) TryLogObject(a ...interface{}) (insolar.LogObjectMarshaller, string) {
	if len(a) == 1 {
		switch vv := a[0].(type) {
		case nil: // the most obvious case
			break
		case string: // the most obvious case
			return nil, vv
		case *string: // handled separately to avoid unnecessary reflect on nil
			if vv == nil {
				break
			}
			return nil, *vv
		case insolar.LogObject:
			m := vv.GetLogObjectMarshaller()
			if m != nil {
				return m, ""
			}
			vr := reflect.ValueOf(vv)
			if vr.Kind() == reflect.Ptr {
				vr = vr.Elem()
			}
			return v.MFactory.CreateLogObjectMarshaller(vr), ""
		case insolar.LogObjectMarshaller:
			return vv, ""
		default:
			if s, ok := defaultStrValuePrepare(vv); ok {
				return nil, s
			}

			vr := reflect.ValueOf(vv)
			switch k := vr.Kind(); k {
			case reflect.Ptr:
				if vr.IsNil() {
					break
				}
				vr = vr.Elem()
				k = vr.Kind()
				if k != reflect.Struct {
					break
				}
				fallthrough
			case reflect.Struct:
				if len(vr.Type().Name()) == 0 {
					return v.MFactory.CreateLogObjectMarshaller(vr), ""
				}
			}
		}
	}
	return nil, v.Sformat(a...)
}
