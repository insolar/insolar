///
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
///

package logtext

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/insolar/insolar/log/logcommon"
)

var _ logcommon.EncoderManager = JsonEncoderManager{}

type JsonEncoderManager struct{}

func (JsonEncoderManager) CreatePartEncoder(b []byte) logcommon.LogObjectWriter {
	return &jsonEncoder{b}
}

func (JsonEncoderManager) FlushPartEncoder(w logcommon.LogObjectWriter) []byte {
	return w.(*jsonEncoder).dst
}

var _ logcommon.LogObjectWriter = &jsonEncoder{}

type jsonEncoder struct {
	dst []byte
}

func (p *jsonEncoder) appendKey(key string) {
	p.dst = AppendString(p.dst, key)
}

func (p *jsonEncoder) appendStrf(f string, a ...interface{}) {
	p.dst = AppendString(p.dst, fmt.Sprintf(f, a...))
}

func (p *jsonEncoder) AddIntField(key string, v int64, fFmt logcommon.LogFieldFormat) {
	p.appendKey(key)
	if fFmt.HasFmt {
		p.appendStrf(fFmt.Fmt, v)
	} else {
		p.dst = strconv.AppendInt(p.dst, v, 10)
	}
}

func (p *jsonEncoder) AddUintField(key string, v uint64, fFmt logcommon.LogFieldFormat) {
	p.appendKey(key)
	switch {
	case fFmt.Kind == reflect.Uintptr:
		if !fFmt.HasFmt {
			fFmt.Fmt = "%v"
		}
		p.appendStrf(fFmt.Fmt, uintptr(v))
	case fFmt.HasFmt:
		p.appendStrf(fFmt.Fmt, v)
	default:
		p.dst = strconv.AppendUint(p.dst, uint64(v), 10)
	}
}

func (p *jsonEncoder) AddBoolField(key string, v bool, fFmt logcommon.LogFieldFormat) {
	p.appendKey(key)
	if fFmt.HasFmt {
		p.appendStrf(fFmt.Fmt, v)
	} else {
		p.dst = strconv.AppendBool(p.dst, v)
	}
}

func (p *jsonEncoder) AddFloatField(key string, v float64, fFmt logcommon.LogFieldFormat) {
	p.appendKey(key)
	if fFmt.HasFmt {
		if fFmt.Kind == reflect.Float32 {
			p.appendStrf(fFmt.Fmt, float32(v))
		} else {
			p.appendStrf(fFmt.Fmt, v)
		}
	} else {
		bits := 64
		if fFmt.Kind == reflect.Float32 {
			bits = 32
		}
		p.dst = AppendFloat(p.dst, v, bits)
	}
}

func (p *jsonEncoder) AddComplexField(key string, v complex128, fFmt logcommon.LogFieldFormat) {
	p.appendKey(key)
	if fFmt.HasFmt {
		p.appendStrf(fFmt.Fmt, v)
	} else {
		bits := 64
		if fFmt.Kind == reflect.Complex64 {
			bits = 32
		}
		p.dst = append(p.dst, '[')
		p.dst = AppendFloat(p.dst, real(v), bits)
		p.dst = append(p.dst, ',')
		p.dst = AppendFloat(p.dst, imag(v), bits)
		p.dst = append(p.dst, ']')
	}
}

func (p *jsonEncoder) AddStrField(key string, v string, fFmt logcommon.LogFieldFormat) {
	p.appendKey(key)
	if fFmt.HasFmt {
		p.appendStrf(fFmt.Fmt, v)
	} else {
		p.dst = AppendString(p.dst, v)
	}
}

func (p *jsonEncoder) AddIntfField(key string, v interface{}, fFmt logcommon.LogFieldFormat) {
	p.appendKey(key)
	if fFmt.HasFmt {
		p.appendStrf(fFmt.Fmt, v)
	} else {
		marshaled, err := json.Marshal(v)
		if err != nil {
			p.appendStrf("marshaling error: %v", err)
		} else {
			p.dst = append(p.dst, marshaled...)
		}
	}
}

func (p *jsonEncoder) AddRawJSONField(key string, v interface{}, fFmt logcommon.LogFieldFormat) {
	p.appendKey(key)
	if fFmt.HasFmt {
		p.dst = append(p.dst, fmt.Sprintf(fFmt.Fmt, v)...)
	} else {
		switch vv := v.(type) {
		case string:
			p.dst = append(p.dst, vv...)
		case []byte:
			p.dst = append(p.dst, vv...)
		default:
			marshaled, err := json.Marshal(vv)
			if err != nil {
				p.appendStrf("marshaling error: %v", err)
			} else {
				p.dst = append(p.dst, marshaled...)
			}
		}
	}
}
