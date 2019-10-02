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
	"strings"
	"unsafe"

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

type MarshallerFactory interface {
	CreateLogObjectMarshaller(o reflect.Value) insolar.LogObjectMarshaller
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
			if ok, s := tryStrValue(vv); ok {
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

func quickTag(prefix string, tag reflect.StructTag) (string, bool) {
	if len(tag) < len(prefix)+1 {
		return "", false
	}
	if strings.HasPrefix(string(tag), prefix) && strings.HasSuffix(string(tag), `"`) {
		return string(tag[len(prefix) : len(tag)-1]), true
	}
	return "", false
}

func printFields(sv reflect.Value, output insolar.LogObjectWriter) string {
	if !sv.CanAddr() {
		s2 := reflect.New(sv.Type()).Elem()
		s2.Set(sv)
		sv = s2
	}

	hasMsg := false
	msg := ""

	st := sv.Type()
	for i := 0; i < sv.NumField(); i++ {
		fv := sv.Field(i)
		ft := st.Field(i)
		if fv.Kind() == reflect.Struct {
			if ft.Anonymous {
				s := printFields(fv, output)
				if !hasMsg {
					msg = s
				}
				continue
			}
			//if tag, hasFmt := quickTag(`fmt:"`, ft.Tag); hasFmt {
			//	zz
			//}
			continue
		}
		tag, hasFmt := quickTag(`fmt:"`, ft.Tag)

		switch ft.Name {
		case "", "_":
			// unreadable field(s)
		case "msg", "message", "Msg", "Message":
			hasMsg = true
			if ok, s, iv := tryReflectStrValue(fv); ok {
				if !hasFmt {
					msg = s
				} else {
					msg = fmt.Sprintf(tag, s)
				}
			} else {
				if !hasFmt {
					tag = "%v"
				}
				msg = fmt.Sprintf(tag, iv)
			}
		default:
			iv := prepareReflectValue(fv)
			if hasFmt {
				output.AddField(ft.Name, fmt.Sprintf(tag, iv))
				continue
			}

			if rawTag, ok := quickTag(`raw:"`, ft.Tag); ok {
				output.AddRawJSON(ft.Name, []byte(fmt.Sprintf(rawTag, iv)))
			} else {
				output.AddField(ft.Name, iv)
			}
		}
	}
	if hasMsg {
		return msg
	}
	if ok, s := tryMsgValue(toInterface(sv)); ok {
		return s
	}
	return ""
}

func toInterface(v reflect.Value) interface{} {
	if !v.CanInterface() {
		v = reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
	}
	return v.Interface()
}

func tryReflectStrValue(v reflect.Value) (bool, string, interface{}) {
	switch k := v.Kind(); k {
	case reflect.Invalid:
		return false, "", nil
	case reflect.String:
		return true, v.String(), nil
	}
	iv := toInterface(v)
	if ok, s := tryStrValue(iv); ok {
		return true, s, nil
	}
	return false, "", iv
}

func prepareReflectValue(v reflect.Value) interface{} {
	switch k := v.Kind(); k {
	case reflect.String:
		return v.String()
	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
	case reflect.Invalid:
		return nil
	}
	iv := toInterface(v)
	if ok, s := tryStrValue(iv); ok {
		return s
	}
	return iv
}

func tryStrValue(v interface{}) (bool, string) {
	switch vv := v.(type) {
	case string:
		return true, vv
	case *string:
		if vv == nil {
			return false, ""
		}
		return true, *vv
	case func() string:
		return true, vv()
	default:
		return tryMsgValue(v)
	}
}

func tryMsgValue(v interface{}) (bool, string) {
	switch vv := v.(type) {
	case LogStringer:
		return true, vv.LogString()
	case fmt.Stringer:
		return true, vv.String()
	default:
		return false, ""
	}
}
