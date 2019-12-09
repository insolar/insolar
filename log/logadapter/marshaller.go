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
	"sync"

	"github.com/insolar/insolar/ledger-v2/reflectkit"

	"github.com/insolar/insolar/log/logcommon"
)

type FieldReporterFunc func(collector logcommon.LogObjectMetricCollector, fieldName string, v interface{})

type MarshallerFactory interface {
	CreateLogObjectMarshaller(o reflect.Value) logcommon.LogObjectMarshaller
	RegisterFieldReporter(fieldType reflect.Type, fn FieldReporterFunc)
}

func GetDefaultLogMsgMarshallerFactory() MarshallerFactory {
	return marshallerFactory
}

var marshallerFactory MarshallerFactory = &defaultLogObjectMarshallerFactory{}

type defaultLogObjectMarshallerFactory struct {
	mutex       sync.RWMutex
	marshallers map[reflect.Type]*typeMarshaller
	reporters   map[reflect.Type]FieldReporterFunc
	forceAddr   bool // enforce use of address/pointer-based access to fields
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

func (p *defaultLogObjectMarshallerFactory) CreateLogObjectMarshaller(o reflect.Value) logcommon.LogObjectMarshaller {
	if o.Kind() != reflect.Struct {
		panic("illegal value")
	}
	t := p.getTypeMarshaller(o.Type())
	if t == nil {
		return nil
	}
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

	tm := typeMarshaller{printNeedsAddr: p.forceAddr, reportNeedsAddr: p.forceAddr}

	if !tm.getFieldsOf(t, 0, p.getFieldReporter) {
		return nil
	}
	return &tm
}

type defaultLogObjectMarshaller struct {
	t *typeMarshaller
	v reflect.Value
}

func (v defaultLogObjectMarshaller) MarshalLogObject(output logcommon.LogObjectWriter, collector logcommon.LogObjectMetricCollector) string {
	return v.t.printFields(v.v, output, collector)
}

func (v defaultLogObjectMarshaller) MarshalMutedLogObject(collector logcommon.LogObjectMetricCollector) {
	if collector == nil {
		return
	}
	v.t.reportFields(v.v, collector)
}

type typeMarshaller struct {
	loggerFields    []logField
	metricFields    []reportField
	msgField        logField
	printNeedsAddr  bool
	reportNeedsAddr bool
}

type reportField struct {
	fieldGet reflectkit.FieldGetterFunc
	reportFn FieldReporterFunc
	name     string
}

type logField struct {
	fieldGet reflectkit.FieldGetterFunc
	reportFn FieldReporterFunc
	getterFn reflectkit.ValueToReceiverFunc
	receiver fieldFmtReceiver
}

type fieldDesc struct {
	reflect.StructField
	getterFn   reflectkit.ValueToReceiverFunc
	index      int
	reportFn   FieldReporterFunc
	outputRecv fieldFmtReceiver
}

func (p *typeMarshaller) getFieldsOf(t reflect.Type, baseOffset uintptr, getReporterFn func(reflect.Type) FieldReporterFunc) bool {
	n := t.NumField()
	var msgGetter fieldDesc
	valueGetters := make([]fieldDesc, 0, n)

	for i := 0; i < n; i++ {
		tf := t.Field(i)
		fieldName := tf.Name

		k := tf.Type.Kind()
		valueGetterFactory := reflectkit.ValueToReceiverFactory(k, marshallerSpecialTypes)
		if valueGetterFactory == nil {
			continue
		}
		unexported := len(tf.PkgPath) != 0

		fd := fieldDesc{StructField: tf, index: i}

		fd.reportFn = getReporterFn(fd.Type)
		tagType, fmtStr := singleTag(fd.Tag)

		msgField := false
		switch {
		case tagType == fmtTagText && fieldName == "_":
			msgField = true
		case fd.reportFn != nil:
			//
		case tagType == fmtTagSkip:
			continue
		case fieldName == "" || fieldName[0] == '_':
			continue
		case !fd.Anonymous:
			//
		case tagType == fmtTagText:
			msgField = true
		default:
			switch k := fd.Type.Kind(); {
			case k == reflect.String:
				msgField = fieldName == "string"
			case k > reflect.Array: // any other non-literals
				continue
			}
		}

		if !msgField {
			switch fieldName {
			case "msg", "Msg", "message", "Message":
				msgField = true
			}
		}

		var needsAddr bool
		switch {
		case tagType != fmtTagSkip && tagType != fmtTagText:
			needsAddr, fd.getterFn = valueGetterFactory(unexported, tf.Type, tagType.IsOpt())
		case msgField || tagType == fmtTagText:
			fd.getterFn = func(_ reflect.Value, out reflectkit.TypedReceiver) {
				out.ReceiveZero(reflect.String)
			}
		default:
			fd.getterFn = func(reflect.Value, reflectkit.TypedReceiver) {}
		}

		fd.outputRecv = fieldFmtReceiver{fmtStr: fmtStr, fmtTag: tagType, k: fd.StructField.Name}

		switch {
		case msgField && msgGetter.getterFn == nil:
			msgGetter = fd
		case msgField && fieldName == "_":
			fd.outputRecv.k = fmt.Sprintf("_txtTag%d", i)
			fallthrough
		default:
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

	if len(valueGetters) == 0 && msgGetter.getterFn == nil {
		return false
	}

	p.loggerFields = make([]logField, 0, len(valueGetters))

	for _, fd := range valueGetters {
		fieldGetter := reflectkit.FieldValueGetter(fd.index, fd.StructField, p.printNeedsAddr, baseOffset)

		if fd.outputRecv.fmtTag != fmtTagSkip {
			p.loggerFields = append(p.loggerFields, logField{
				fieldGetter, fd.reportFn, fd.getterFn, fd.outputRecv})
		}

		if fd.reportFn != nil {
			reportFieldGetter := fieldGetter
			if p.reportNeedsAddr != p.printNeedsAddr {
				reportFieldGetter = reflectkit.FieldValueGetter(fd.index, fd.StructField, p.reportNeedsAddr, baseOffset)
			}
			p.metricFields = append(p.metricFields, reportField{
				reportFieldGetter, fd.reportFn, fd.Name})
		}
	}

	if msgGetter.getterFn == nil {
		p.msgField = logField{}
		return true
	}

	fieldGetter := reflectkit.FieldValueGetter(msgGetter.index, msgGetter.StructField, p.printNeedsAddr, baseOffset)
	p.msgField = logField{
		fieldGetter, msgGetter.reportFn, msgGetter.getterFn, msgGetter.outputRecv}
	p.msgField.receiver.k = ""

	return true
}

func (p *typeMarshaller) prepareValue(value reflect.Value) reflect.Value {
	return p._prepareValue(value, p.printNeedsAddr)
}

func (p *typeMarshaller) _prepareValue(value reflect.Value, needsAddr bool) reflect.Value {
	if !needsAddr {
		return value
	}
	return reflectkit.MakeAddressable(value)
}

func (f logField) printField(v reflect.Value, w logcommon.LogObjectWriter, c logcommon.LogObjectMetricCollector) {
	fieldValue := f.fieldGet(v)
	if c != nil && f.reportFn != nil {
		f.reportFn(c, f.receiver.k, fieldValue.Interface())
	}
	f.receiver.w = w
	f.getterFn(fieldValue, f.receiver)
}

func (p *typeMarshaller) printFields(value reflect.Value, writer logcommon.LogObjectWriter, collector logcommon.LogObjectMetricCollector) string {
	value = p._prepareValue(value, p.printNeedsAddr) // double check

	for _, field := range p.loggerFields {
		field.printField(value, writer, collector)
	}

	if field := p.msgField; field.getterFn == nil {
		return ""
	} else {
		sc := stringCapturer{}
		field.printField(value, &sc, collector)
		return sc.v
	}
}

func (p *typeMarshaller) reportFields(value reflect.Value, collector logcommon.LogObjectMetricCollector) {
	if len(p.metricFields) == 0 {
		return
	}

	value = p._prepareValue(value, p.reportNeedsAddr) // double check

	for _, field := range p.metricFields {
		fieldValue := field.fieldGet(value)
		field.reportFn(collector, field.name, fieldValue.Interface())
	}
}

type fmtTagType uint8

const (
	fmtTagDefault fmtTagType = iota
	fmtTagOptional

	fmtTagText
	fmtTagSkip // + opt

	fmtTagFormatRaw
	fmtTagFormatRawOpt // + opt

	fmtTagFormatValue
	fmtTagFormatValueOpt // + opt
)

func (v fmtTagType) IsOpt() bool {
	return v&fmtTagOptional != 0
}

func (v fmtTagType) IsRaw() bool {
	return v&^1 == fmtTagFormatRaw
}

func (v fmtTagType) HasFmt() bool {
	return v >= fmtTagFormatRaw
}

func singleTag(tag reflect.StructTag) (fmtTagType, string) {
	tagType := fmtTagDefault
	if _, v, ok := ParseStructTag(tag, func(name, _ string) bool {
		switch name {
		case "fmt+opt", "opt+fmt":
			tagType = fmtTagFormatValueOpt
		case "fmt":
			tagType = fmtTagFormatValue
		case "raw+opt", "opt+raw":
			tagType = fmtTagFormatRawOpt
		case "raw":
			tagType = fmtTagFormatRaw
		case "skip":
			tagType = fmtTagSkip
		case "txt":
			tagType = fmtTagText
		case "opt":
			tagType = fmtTagOptional
		default:
			return false
		}
		return true
	}); ok {
		return tagType, v
	}
	return fmtTagDefault, ""
}

func marshallerSpecialTypes(t reflect.Type, checkZero bool) reflectkit.IfaceToReceiverFunc {
	var prepFn valuePrepareFn

	switch kind := t.Kind(); kind {
	case reflect.Interface:
		prepFn = prepareValue

	case reflect.Ptr:
		if te := t.Elem(); te.Kind() == reflect.String {
			return func(value interface{}, _ reflect.Kind, out reflectkit.TypedReceiver) {
				if vv := value.(*string); vv == nil {
					out.ReceiveNil(reflect.String)
				} else if v := *vv; !checkZero || v != "" {
					out.ReceiveString(reflect.String, v)
				} else {
					out.ReceiveZero(reflect.String)
				}
			}
		}
		fallthrough

	default:
		prepFn = findPrepareValueFn(t)
		if prepFn == nil {
			return nil
		}
	}

	return func(value interface{}, kind reflect.Kind, out reflectkit.TypedReceiver) {
		switch s, k, isNil := prepFn(value); {
		case k == reflect.Invalid:
			out.ReceiveElse(kind, value, false)
		case isNil:
			out.ReceiveNil(kind)
		case !checkZero || s != "":
			out.ReceiveString(kind, s)
		default:
			out.ReceiveZero(kind)
		}
	}
}
