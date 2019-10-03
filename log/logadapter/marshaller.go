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

package logadapter

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unsafe"

	"github.com/insolar/insolar/insolar"
)

var marshallerFactory MarshallerFactory = &defaultLogObjectMarshallerFactory{}

type defaultLogObjectMarshallerFactory struct {
	mutex       sync.RWMutex
	marshallers map[reflect.Type]*typeMarshaller
	reporters   map[reflect.Type]FieldReporterFunc
}

func (p *defaultLogObjectMarshallerFactory) RegisterFieldReporter(fieldType reflect.Type, fn FieldReporterFunc) {
	if fn == nil {
		panic("illegal value")
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.reporters == nil {
		p.reporters = make(map[reflect.Type]FieldReporterFunc)
	}
	p.reporters[fieldType] = fn
}

func (p *defaultLogObjectMarshallerFactory) CreateLogObjectMarshaller(o reflect.Value) insolar.LogObjectMarshaller {
	if o.Kind() != reflect.Struct {
		panic("illegal value")
	}
	t := p.getTypeMarshaller(o.Type())
	return defaultLogObjectMarshaller{t, t.prepareValue(o)} // do prepare for a repeated use of marshaller
}

func (p *defaultLogObjectMarshallerFactory) getFieldReporter(t reflect.Type) FieldReporterFunc {
	p.mutex.RLock()
	fr := p.reporters[t]
	p.mutex.RUnlock()
	return fr
}

func (p *defaultLogObjectMarshallerFactory) getTypeMarshaller(t reflect.Type) *typeMarshaller {
	p.mutex.RLock()
	tm := p.marshallers[t]
	p.mutex.RUnlock()
	if tm != nil {
		return tm
	}

	tm = p.buildTypeMarshaller(t) // do before lock to reduce in-lock time

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.marshallers == nil {
		p.marshallers = make(map[reflect.Type]*typeMarshaller)
	} else {
		tm2 := p.marshallers[t]
		if tm2 != nil {
			return tm2
		}
	}
	p.marshallers[t] = tm
	return tm
}

func (p *defaultLogObjectMarshallerFactory) buildTypeMarshaller(t reflect.Type) *typeMarshaller {
	n := t.NumField()
	if n <= 0 {
		return nil
	}

	tm := typeMarshaller{}

	if !tm.getFieldsOf(t, 0, p.getFieldReporter) {
		return nil
	}
	return &tm
}

type defaultLogObjectMarshaller struct {
	t *typeMarshaller
	v reflect.Value
}

func (v defaultLogObjectMarshaller) MarshalLogObject(output insolar.LogObjectWriter, collector insolar.LogObjectMetricCollector) string {
	return v.t.printFields(v.v, output, collector)
}

func (v defaultLogObjectMarshaller) MarshalMutedLogObject(collector insolar.LogObjectMetricCollector) {
	v.t.reportFields(v.v, collector)
}

type fieldValueGetterFunc func(value reflect.Value) interface{}

var fieldValueGetters = map[reflect.Kind]func(unexported bool, t reflect.Type) (bool, fieldValueGetterFunc){
	// ======== Simple values, are safe to read from unexported fields ============
	reflect.Bool: func(bool, reflect.Type) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			return value.Bool()
		}
	},
	reflect.Int: func(bool, reflect.Type) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			return int(value.Int())
		}
	},
	reflect.Int8: func(bool, reflect.Type) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			return int8(value.Int())
		}
	},
	reflect.Int16: func(bool, reflect.Type) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			return int16(value.Int())
		}
	},
	reflect.Int32: func(bool, reflect.Type) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			return int32(value.Int())
		}
	},
	reflect.Int64: func(bool, reflect.Type) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			return value.Int()
		}
	},
	reflect.Uint: func(bool, reflect.Type) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			return uint(value.Uint())
		}
	},
	reflect.Uint8: func(bool, reflect.Type) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			return uint8(value.Uint())
		}
	},
	reflect.Uint16: func(bool, reflect.Type) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			return uint16(value.Uint())
		}
	},
	reflect.Uint32: func(bool, reflect.Type) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			return uint32(value.Uint())
		}
	},
	reflect.Uint64: func(bool, reflect.Type) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			return value.Uint()
		}
	},
	reflect.Uintptr: func(bool, reflect.Type) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			return uintptr(value.Uint())
		}
	},
	reflect.Float32: func(bool, reflect.Type) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			return float32(value.Float())
		}
	},
	reflect.Float64: func(bool, reflect.Type) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			return value.Float()
		}
	},
	reflect.Complex64: func(bool, reflect.Type) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			return complex64(value.Complex())
		}
	},
	reflect.Complex128: func(bool, reflect.Type) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			return value.Complex()
		}
	},
	reflect.String: func(bool, reflect.Type) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			return value.String()
		}
	},

	// ============ Special handling for unexported fields ===========

	reflect.Ptr: func(unexported bool, t reflect.Type) (bool, fieldValueGetterFunc) {
		te := t.Elem()
		if te.Kind() == reflect.String {
			return false, func(value reflect.Value) interface{} {
				if value.IsNil() {
					return nil
				}
				return value.Elem().String()
			}
		}
		return defaultObjFieldGetterFactory(unexported, t)
	},

	reflect.Func: func(unexported bool, t reflect.Type) (bool, fieldValueGetterFunc) {
		if t.NumIn() == 0 && t.NumOut() == 1 && t.Out(0).Kind() == reflect.String {
			return unexported, func(value reflect.Value) interface{} {
				if value.IsNil() {
					return nil
				}
				fn := value.Interface().(func() string)
				return fn()
			}
		}
		return unexported, reflect.Value.Interface
	},

	reflect.Interface: func(unexported bool, _ reflect.Type) (b bool, getterFunc fieldValueGetterFunc) {
		return unexported, func(value reflect.Value) interface{} {
			if value.IsNil() {
				return nil
			}
			iv := value.Interface()
			switch vv := iv.(type) {
			case func() string:
				return vv()
			default:
				vv, _ = tryDefaultValuePrepare(vv)
				return vv
			}
		}
	},

	reflect.Struct: defaultObjFieldGetterFactory,
	reflect.Array:  defaultObjFieldGetterFactory,
	reflect.Map:    defaultObjFieldGetterFactory,
	reflect.Slice:  defaultObjFieldGetterFactory,
	reflect.Chan:   defaultObjFieldGetterFactory,

	// ============ Excluded ===================
	//reflect.UnsafePointer
}

var prepareObjTypes = []struct {
	t  reflect.Type
	fn func(interface{}) (interface{}, bool)
}{
	{reflect.TypeOf((*LogStringer)(nil)).Elem(), func(value interface{}) (interface{}, bool) {
		if vv, ok := value.(LogStringer); ok {
			return vv.LogString(), true
		}
		return value, false
	}},
	{reflect.TypeOf((*fmt.Stringer)(nil)).Elem(), func(value interface{}) (interface{}, bool) {
		if vv, ok := value.(fmt.Stringer); ok {
			return vv.String(), true
		}
		return value, false
	}},
}

func defaultStrValuePrepare(iv interface{}) (string, bool) {
	switch vv := iv.(type) {
	case string:
		return vv, true
	case *string:
		if vv == nil {
			return "", false
		}
		return *vv, true
	case func() string:
		return vv(), true
	}
	if vv, ok := tryDefaultValuePrepare(iv); ok {
		return vv.(string), true
	}
	return "", false
}

func defaultObjFieldGetterFactory(unexported bool, t reflect.Type) (bool, fieldValueGetterFunc) {
	for _, f := range prepareObjTypes {
		if t.Implements(f.t) {
			fn := f.fn
			if t.Kind() == reflect.Struct {
				return unexported, func(value reflect.Value) interface{} {
					vv, _ := fn(value.Interface())
					return vv
				}
			}

			return unexported, func(value reflect.Value) interface{} {
				if value.IsNil() {
					return nil
				}
				vv, _ := fn(value.Interface())
				return vv
			}
		}
	}
	return unexported, reflect.Value.Interface
}

func tryDefaultValuePrepare(iv interface{}) (interface{}, bool) {
	for _, f := range prepareObjTypes {
		if vv, ok := f.fn(iv); ok {
			return vv, true
		}
	}
	return iv, false
}

func getFieldGetter(index int, fd reflect.StructField, useAddr bool, baseOffset uintptr) func(reflect.Value) reflect.Value {
	if !useAddr {
		return func(value reflect.Value) reflect.Value {
			return value.Field(index)
		}
	}

	fieldOffset := fd.Offset + baseOffset
	fieldType := fd.Type

	return func(value reflect.Value) reflect.Value {
		return offsetFieldGetter(value, fieldOffset, fieldType)
	}
}

func offsetFieldGetter(v reflect.Value, fieldOffset uintptr, fieldType reflect.Type) reflect.Value {
	return reflect.NewAt(fieldType, unsafe.Pointer(v.UnsafeAddr()+fieldOffset)).Elem()
}

type fieldMarshallerFunc func(value reflect.Value, writer insolar.LogObjectWriter, collector insolar.LogObjectMetricCollector)
type fieldReportFunc func(value reflect.Value, collector insolar.LogObjectMetricCollector)
type fieldMarshallerMsgFunc func(value reflect.Value) string

type typeMarshaller struct {
	fields          []fieldMarshallerFunc
	reporters       []fieldReportFunc
	msgField        fieldMarshallerMsgFunc
	printNeedsAddr  bool
	reportNeedsAddr bool
}

func (p *typeMarshaller) getFieldsOf(t reflect.Type, baseOffset uintptr, getReporterFn func(reflect.Type) FieldReporterFunc) bool {
	n := t.NumField()

	type fieldDesc struct {
		reflect.StructField
		getFn    fieldValueGetterFunc
		reportFn FieldReporterFunc
		index    int
	}

	var msgGetter fieldDesc
	valueGetters := make([]fieldDesc, 0, n)

	for i := 0; i < n; i++ {
		tf := t.Field(i)
		fieldName := tf.Name

		if tf.Anonymous || fieldName == "" || fieldName[0] == '_' || strings.HasPrefix(string(tf.Tag), `skip:"`) {
			continue
		}

		k := tf.Type.Kind()
		valueGetterFactory := fieldValueGetters[k]
		if valueGetterFactory == nil {
			continue
		}
		unexported := len(tf.PkgPath) != 0

		fd := fieldDesc{StructField: tf, index: i}
		needsAddr := false
		needsAddr, fd.getFn = valueGetterFactory(unexported, tf.Type) // default handler
		if fd.getFn == nil {
			continue
		}

		switch fieldName {
		case "msg", "Msg", "message", "Message":
			msgGetter = fd
		default:
			fd.reportFn = getReporterFn(fd.Type)
			valueGetters = append(valueGetters, fd)
		}

		if needsAddr {
			p.printNeedsAddr = true
			if fd.reportFn != nil {
				p.reportNeedsAddr = true
			}
		}
	}

	if p.reportNeedsAddr && !p.printNeedsAddr {
		panic("illegal state")
	}

	if len(valueGetters) == 0 && msgGetter.getFn == nil {
		return false
	}

	p.fields = make([]fieldMarshallerFunc, len(valueGetters))

	for i, fd := range valueGetters {
		fieldGetter := getFieldGetter(fd.index, fd.StructField, p.printNeedsAddr, baseOffset)
		valueGetter := fd.getFn
		fieldName := fd.Name
		fieldReporter := fd.reportFn

		switch tagType, fmtStr := singleTag(fd.Tag); tagType {
		case "fmt":
			p.fields[i] = func(obj reflect.Value, writer insolar.LogObjectWriter, collector insolar.LogObjectMetricCollector) {
				f := fieldGetter(obj)
				v := valueGetter(f)
				if fieldReporter != nil && collector != nil {
					fieldReporter(collector, fieldName, v)
				}
				s := fmt.Sprintf(fmtStr, v)
				writer.AddField(fieldName, s)
			}
		case "raw":
			p.fields[i] = func(obj reflect.Value, writer insolar.LogObjectWriter, collector insolar.LogObjectMetricCollector) {
				f := fieldGetter(obj)
				v := valueGetter(f)
				if fieldReporter != nil && collector != nil {
					fieldReporter(collector, fieldName, v)
				}
				s := fmt.Sprintf(fmtStr, v)
				writer.AddRawJSON(fieldName, []byte(s))
			}
		default:
			p.fields[i] = func(obj reflect.Value, writer insolar.LogObjectWriter, collector insolar.LogObjectMetricCollector) {
				f := fieldGetter(obj)
				v := valueGetter(f)
				if fieldReporter != nil && collector != nil {
					fieldReporter(collector, fieldName, v)
				}
				writer.AddField(fieldName, v)
			}
		}

		if fieldReporter != nil {
			reportFieldGetter := fieldGetter
			if p.reportNeedsAddr != p.printNeedsAddr {
				reportFieldGetter = getFieldGetter(fd.index, fd.StructField, p.reportNeedsAddr, baseOffset)
			}

			p.reporters = append(p.reporters, func(obj reflect.Value, collector insolar.LogObjectMetricCollector) {
				if collector == nil {
					return
				}
				f := reportFieldGetter(obj)
				v := valueGetter(f)
				fieldReporter(collector, fieldName, v)
			})
		}
	}

	if msgGetter.getFn == nil {
		p.msgField = nil
		return true
	}

	fieldGetter := getFieldGetter(msgGetter.index, msgGetter.StructField, p.printNeedsAddr, baseOffset)
	valueGetter := msgGetter.getFn

	switch tagType, fmtStr := singleTag(msgGetter.Tag); tagType {
	case "fmt":
		p.msgField = func(obj reflect.Value) string {
			f := fieldGetter(obj)
			v := valueGetter(f)
			s := fmt.Sprintf(fmtStr, v)
			return s
		}
	default:
		p.msgField = func(obj reflect.Value) string {
			f := fieldGetter(obj)
			v := valueGetter(f)
			s := fmt.Sprintf("%v", v)
			return s
		}
	}
	return true
}

func (p *typeMarshaller) prepareValue(value reflect.Value) reflect.Value {
	return p._prepareValue(value, p.printNeedsAddr)
}

func (p *typeMarshaller) _prepareValue(value reflect.Value, needsAddr bool) reflect.Value {
	if !needsAddr || value.CanAddr() {
		return value
	}
	valueCopy := reflect.New(value.Type()).Elem()
	valueCopy.Set(value)
	return valueCopy
}

func (p *typeMarshaller) printFields(value reflect.Value, writer insolar.LogObjectWriter, collector insolar.LogObjectMetricCollector) string {
	value = p._prepareValue(value, p.printNeedsAddr) // double check

	for _, fn := range p.fields {
		fn(value, writer, collector)
	}
	if p.msgField == nil {
		return ""
	}
	return p.msgField(value)
}

func (p *typeMarshaller) reportFields(value reflect.Value, collector insolar.LogObjectMetricCollector) {
	if len(p.reporters) == 0 {
		return
	}

	value = p._prepareValue(value, p.reportNeedsAddr) // double check

	for _, fn := range p.reporters {
		fn(value, collector)
	}
}

func singleTag(tag reflect.StructTag) (string, string) {
	if len(tag) <= 3 {
		return "", ""
	}

	colon := strings.IndexByte(string(tag), ':')
	if colon <= 0 || colon+2 >= len(tag) || tag[colon+1] != '"' || tag[len(tag)-1] != '"' {
		return "", ""
	}

	return string(tag[:colon]), string(tag[colon+2 : len(tag)-1])
}
