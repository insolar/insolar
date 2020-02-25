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

type FieldReporterFunc func(collector insolar.LogObjectMetricCollector, fieldName string, v interface{})

type MarshallerFactory interface {
	CreateLogObjectMarshaller(o reflect.Value) insolar.LogObjectMarshaller
	RegisterFieldReporter(fieldType reflect.Type, fn FieldReporterFunc)
}

func GetDefaultLogMsgMarshallerFactory() MarshallerFactory {
	return marshallerFactory
}

func (v MsgFormatConfig) FmtLogObject(a ...interface{}) (insolar.LogObjectMarshaller, string) {
	if len(a) == 1 {
		switch vv := a[0].(type) {
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
		case string: // the most obvious case
			return nil, vv
		case *string: // handled separately to avoid unnecessary reflect
			if vv == nil {
				break
			}
			return nil, *vv
		case nil:
			break
		default:
			if s, ok := defaultStrValuePrepare(vv); ok {
				return nil, s
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
					return v.MFactory.CreateLogObjectMarshaller(vr), ""
				}
			}
		}
	}
	return nil, v.Sformat(a...)
}

func (v MsgFormatConfig) PrepareMutedLogObject(a ...interface{}) insolar.LogObjectMarshaller {
	if len(a) != 1 {
		return nil
	}

	switch vv := a[0].(type) {
	case insolar.LogObject:
		m := vv.GetLogObjectMarshaller()
		if m != nil {
			return m
		}
		vr := reflect.ValueOf(vv)
		if vr.Kind() == reflect.Ptr {
			vr = vr.Elem()
		}
		return v.MFactory.CreateLogObjectMarshaller(vr)
	case insolar.LogObjectMarshaller:
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
