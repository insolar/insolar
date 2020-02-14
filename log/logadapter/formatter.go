// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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

func GetDefaultLogMsgFormatter() MsgFormatConfig {
	return MsgFormatConfig{
		Sformat:  fmt.Sprint,
		Sformatf: fmt.Sprintf,
	}
}

type MsgFormatConfig struct {
	Sformat  FormatFunc
	Sformatf FormatfFunc
}

func (v MsgFormatConfig) TryLogObject(a ...interface{}) (insolar.LogObjectMarshaller, string) {
	if len(a) == 1 {
		switch v := a[0].(type) {
		case nil: // the most obvious case(s)
			break
		case string: // the most obvious case(s)
			return nil, v
		case insolar.LogObjectMarshaller:
			return v, ""
		default:
			vt := reflect.ValueOf(v)
			if vt.Kind() == reflect.Struct && len(vt.Type().Name()) == 0 {
				return defaultLogObjectMarshaller{vt}, ""
			}
		}
	}
	return nil, v.Sformat(a...)
}

func GetInlineLogObjectMarshaller(v reflect.Value) insolar.LogObjectMarshaller {
	if v.Kind() != reflect.Struct {
		panic("illegal value")
	}
	return defaultLogObjectMarshaller{v}
}

type defaultLogObjectMarshaller struct {
	v reflect.Value
}

func (v defaultLogObjectMarshaller) MarshalLogObject(output insolar.LogObjectWriter) string {
	return printFields(v.v, output)
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
