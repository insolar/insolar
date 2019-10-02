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
}

func (p *defaultLogObjectMarshallerFactory) CreateLogObjectMarshaller(o reflect.Value) insolar.LogObjectMarshaller {
	if o.Kind() != reflect.Struct {
		panic("illegal value")
	}
	t := p.getTypeMarshaller(o.Type())
	return defaultLogObjectMarshaller{t, o}
}

func (p *defaultLogObjectMarshallerFactory) getTypeMarshaller(t reflect.Type) *typeMarshaller {
	p.mutex.RLock()
	tm := p.marshallers[t]
	p.mutex.RUnlock()
	if tm != nil {
		return tm
	}

	tm = p.buildTypeMarshaller(t)

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

type fieldValueGetterFunc func(value reflect.Value) interface{}

var fieldValueGetters = map[reflect.Kind]func(unexported bool, t reflect.Type) (bool, fieldValueGetterFunc){
	reflect.Bool: func(bool, reflect.Type) (bool, fieldValueGetterFunc) {
		return false, func(value reflect.Value) interface{} {
			return value.Bool()
		}
	},
	//reflect.Int
	//reflect.Int8
	//reflect.Int16
	//reflect.Int32
	//reflect.Int64
	//reflect.Uint
	//reflect.Uint8
	//reflect.Uint16
	//reflect.Uint32
	//reflect.Uint64
	//reflect.Uintptr
	//reflect.Float32
	//reflect.Float64
	//reflect.Complex64
	//reflect.Complex128
	//reflect.String

	// ============ Special Handling ===========
	reflect.Ptr: func(_ bool, t reflect.Type) (bool, fieldValueGetterFunc) {
		te := t.Elem()
		switch te.Kind() {
		case reflect.String:
			return false, func(value reflect.Value) interface{} {
				if value.IsNil() {
					return nil
				}
				return value.Elem().String()
			}
		default:
			return false, nil
		}
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

	reflect.Interface: func(unexported bool, t reflect.Type) (b bool, getterFunc fieldValueGetterFunc) {
		return unexported, func(value reflect.Value) interface{} {
			if value.IsNil() {
				return nil
			}
			iv := value.Interface()
			switch vv := iv.(type) {
			case func() string:
				return vv()
			default:
				return defaultValuePrepare(vv)
			}
		}
	},

	reflect.Struct: defaultObjFieldGetterFactory,
	reflect.Array:  defaultObjFieldGetterFactory,
	reflect.Map:    defaultObjFieldGetterFactory,
	reflect.Slice:  defaultObjFieldGetterFactory,

	// ============ Excluded ===================
	//reflect.Chan
	//reflect.UnsafePointer
}

func defaultValuePrepare(iv interface{}) interface{} {
	switch vv := iv.(type) {
	case LogStringer:
		return vv.LogString()
	case fmt.Stringer:
		return vv.String()
	default:
		return vv
	}
}

var filterReflectTypes = []reflect.Type{
	reflect.TypeOf((*LogStringer)(nil)).Elem(),
	reflect.TypeOf((*fmt.Stringer)(nil)).Elem(),
}

func defaultObjFieldGetterFactory(unexported bool, t reflect.Type) (b bool, getterFunc fieldValueGetterFunc) {
	for _, ft := range filterReflectTypes {

	}

	return unexported, func(value reflect.Value) interface{} {
		return defaultValuePrepare(value.Interface())
	}
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

func getFieldsOf(t reflect.Type, baseOffset uintptr) (bool, []fieldMarshallerFunc, fieldMarshallerMsgFunc) {
	n := t.NumField()

	type fieldDesc struct {
		reflect.StructField
		fn    fieldValueGetterFunc
		index int
	}

	var msgGetter fieldDesc
	valueGetters := make([]fieldDesc, 0, n)
	needsAddr := false

	for i := 0; i < n; i++ {
		tf := t.Field(i)
		fieldName := tf.Name

		if fieldName == "" || fieldName[0] == '_' || tf.Anonymous || strings.HasPrefix(string(tf.Tag), `skip:"`) {
			continue
		}

		k := tf.Type.Kind()
		valueGetterFactory := fieldValueGetters[k]
		if valueGetterFactory == nil {
			continue
		}
		unexported := len(tf.PkgPath) != 0
		addrReq, valueGetter := valueGetterFactory(unexported, tf.Type) // default handler
		if valueGetter == nil {
			continue
		}
		if addrReq {
			needsAddr = true
		}

		//if tf.Anonymous {
		//	switch k {
		//	case reflect.Ptr:
		//		et := tf.Type.Elem()
		//		if et.Kind() == reflect.Struct {
		//			p.addFieldsOf(et)
		//			continue
		//		}
		//	case reflect.Struct:
		//		p.addFieldsOf(tf.Type)
		//		continue
		//	}
		//}

		fd := fieldDesc{tf, valueGetter, i}
		switch fieldName {
		case "msg", "Msg", "message", "Message":
			msgGetter = fd
		default:
			valueGetters = append(valueGetters, fd)
		}
	}

	if len(valueGetters) == 0 && msgGetter.fn == nil {
		return false, nil, nil
	}

	fields := make([]fieldMarshallerFunc, len(valueGetters))

	for i, fd := range valueGetters {
		fieldGetter := getFieldGetter(fd.index, fd.StructField, needsAddr, baseOffset)
		valueGetter := fd.fn
		fieldName := fd.Name

		switch tagType, fmtStr := singleTag(fd.Tag); tagType {
		case "fmt":
			fields[i] = func(obj reflect.Value, writer insolar.LogObjectWriter) {
				f := fieldGetter(obj)
				v := valueGetter(f)
				s := fmt.Sprintf(fmtStr, v)
				writer.AddField(fieldName, s)
			}
		case "raw":
			fields[i] = func(obj reflect.Value, writer insolar.LogObjectWriter) {
				f := fieldGetter(obj)
				v := valueGetter(f)
				s := fmt.Sprintf(fmtStr, v)
				writer.AddRawJSON(fieldName, []byte(s))
			}
		default:
			fields[i] = func(obj reflect.Value, writer insolar.LogObjectWriter) {
				f := fieldGetter(obj)
				v := valueGetter(f)
				writer.AddField(fieldName, v)
			}
		}
	}

	if msgGetter.fn == nil {
		return needsAddr, fields, nil
	}

	fieldGetter := getFieldGetter(msgGetter.index, msgGetter.StructField, needsAddr, baseOffset)
	valueGetter := msgGetter.fn

	switch tagType, fmtStr := singleTag(msgGetter.Tag); tagType {
	case "fmt":
		return needsAddr, fields, func(obj reflect.Value) string {
			f := fieldGetter(obj)
			v := valueGetter(f)
			s := fmt.Sprintf(fmtStr, v)
			return s
		}
	default:
		return needsAddr, fields, func(obj reflect.Value) string {
			f := fieldGetter(obj)
			v := valueGetter(f)
			s := fmt.Sprintf("%v", v)
			return s
		}
	}
}

func offsetFieldGetter(v reflect.Value, fieldOffset uintptr, fieldType reflect.Type) reflect.Value {
	return reflect.NewAt(fieldType, unsafe.Pointer(v.UnsafeAddr()+fieldOffset)).Elem()
}

func (p *defaultLogObjectMarshallerFactory) buildTypeMarshaller(t reflect.Type) *typeMarshaller {

	n := t.NumField()
	if n <= 0 {
		return nil
	}

	tm := typeMarshaller{}
	tm.needsAddr, tm.fields, tm.msgField = getFieldsOf(t, 0)
	if len(tm.fields) == 0 && tm.msgField == nil {
		return nil
	}
	return &tm
}

type defaultLogObjectMarshaller struct {
	t *typeMarshaller
	v reflect.Value
}

func (v defaultLogObjectMarshaller) MarshalLogObject(output insolar.LogObjectWriter) string {
	return v.t.printFields(v.v, output)
}

type fieldMarshallerFunc func(value reflect.Value, writer insolar.LogObjectWriter)
type fieldMarshallerMsgFunc func(value reflect.Value) string

type typeMarshaller struct {
	fields    []fieldMarshallerFunc
	msgField  fieldMarshallerMsgFunc
	needsAddr bool
}

func (p *typeMarshaller) printFields(value reflect.Value, writer insolar.LogObjectWriter) string {
	if p.needsAddr && !value.CanAddr() {
		s2 := reflect.New(value.Type()).Elem()
		s2.Set(value)
		value = s2
	}

	for _, fn := range p.fields {
		fn(value, writer)
	}
	if p.msgField == nil {
		return ""
	}
	return p.msgField(value)
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
