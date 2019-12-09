package logadapter

import (
	"reflect"

	"github.com/insolar/insolar/ledger-v2/reflectkit"
	"github.com/insolar/insolar/log/logcommon"
)

var _ reflectkit.TypedReceiver = fieldFmtReceiver{}

type fieldFmtReceiver struct {
	w logcommon.LogObjectWriter

	k      string
	fmtStr string
	fmtTag fmtTagType
}

func (f fieldFmtReceiver) def(t reflect.Kind) bool {
	switch f.fmtTag {
	case fmtTagText:
		f.w.AddStrField(f.k, f.fmtStr, logcommon.LogFieldFormat{Kind: t})
		return true
	}
	return false
}

func (f fieldFmtReceiver) fmt(t reflect.Kind) logcommon.LogFieldFormat {
	return logcommon.LogFieldFormat{Fmt: f.fmtStr, Kind: t, HasFmt: f.fmtTag.HasFmt()}
}

func (f fieldFmtReceiver) ReceiveBool(t reflect.Kind, v bool) {
	switch {
	case f.def(t):
		return
	case f.fmtTag.IsRaw():
		f.w.AddRawJSONField(f.k, v, f.fmt(t))
	default:
		f.w.AddBoolField(f.k, v, f.fmt(t))
	}
}

func (f fieldFmtReceiver) ReceiveInt(t reflect.Kind, v int64) {
	switch {
	case f.def(t):
		return
	case f.fmtTag.IsRaw():
		f.w.AddRawJSONField(f.k, v, f.fmt(t))
	default:
		f.w.AddIntField(f.k, v, f.fmt(t))
	}
}

func (f fieldFmtReceiver) ReceiveUint(t reflect.Kind, v uint64) {
	switch {
	case f.def(t):
		return
	case f.fmtTag.IsRaw():
		f.w.AddRawJSONField(f.k, v, f.fmt(t))
	default:
		f.w.AddUintField(f.k, v, f.fmt(t))
	}
}

func (f fieldFmtReceiver) ReceiveFloat(t reflect.Kind, v float64) {
	switch {
	case f.def(t):
		return
	case f.fmtTag.IsRaw():
		f.w.AddRawJSONField(f.k, v, f.fmt(t))
	default:
		f.w.AddFloatField(f.k, v, f.fmt(t))
	}
}

func (f fieldFmtReceiver) ReceiveComplex(t reflect.Kind, v complex128) {
	switch {
	case f.def(t):
		return
	case f.fmtTag.IsRaw():
		f.w.AddRawJSONField(f.k, v, f.fmt(t))
	default:
		f.w.AddComplexField(f.k, v, f.fmt(t))
	}
}

func (f fieldFmtReceiver) ReceiveString(t reflect.Kind, v string) {
	switch {
	case f.def(t):
		return
	case f.fmtTag.IsRaw():
		f.w.AddRawJSONField(f.k, v, f.fmt(t))
	default:
		f.w.AddStrField(f.k, v, f.fmt(t))
	}
}

func (f fieldFmtReceiver) ReceiveZero(t reflect.Kind) {
	f.def(t)
}

func (f fieldFmtReceiver) ReceiveNil(t reflect.Kind) {
	switch {
	case f.def(t) || f.fmtTag.IsOpt():
		return
	case f.fmtTag.IsRaw():
		f.w.AddRawJSONField(f.k, nil, f.fmt(t))
	default:
		f.w.AddIntfField(f.k, nil, f.fmt(t))
	}
}

func (f fieldFmtReceiver) ReceiveIface(t reflect.Kind, v interface{}) {
	switch {
	case f.def(t):
		return
	case f.fmtTag.IsRaw():
		f.w.AddRawJSONField(f.k, v, f.fmt(t))
	default:
		f.w.AddIntfField(f.k, v, f.fmt(t))
	}
}

func (f fieldFmtReceiver) ReceiveElse(t reflect.Kind, v interface{}, isZero bool) {
	switch {
	case f.def(t) || f.fmtTag.IsOpt() && isZero:
		return
	case f.fmtTag.IsRaw():
		f.w.AddRawJSONField(f.k, v, f.fmt(t))
	default:
		f.w.AddIntfField(f.k, v, f.fmt(t))
	}
}

type stringCapturer struct {
	v string
}

func (p *stringCapturer) AddComplexField(key string, v complex128, fmt logcommon.LogFieldFormat) {
	p.v = fmt.ToString(v, "%v")
}

func (p *stringCapturer) AddRawJSONField(_ string, v interface{}, fmt logcommon.LogFieldFormat) {
	p.v = fmt.ToString(v, "%v")
}

func (p *stringCapturer) AddIntField(_ string, v int64, fmt logcommon.LogFieldFormat) {
	p.v = fmt.ToString(v, "%v")
}

func (p *stringCapturer) AddUintField(_ string, v uint64, fmt logcommon.LogFieldFormat) {
	p.v = fmt.ToString(v, "%v")
}

func (p *stringCapturer) AddBoolField(_ string, v bool, fmt logcommon.LogFieldFormat) {
	p.v = fmt.ToString(v, "%v")
}

func (p *stringCapturer) AddFloatField(_ string, v float64, fmt logcommon.LogFieldFormat) {
	p.v = fmt.ToString(v, "%v")
}

func (p *stringCapturer) AddStrField(_ string, v string, fmt logcommon.LogFieldFormat) {
	if fmt.HasFmt {
		p.v = fmt.ToString(v, "%v")
	} else {
		p.v = v
	}
}

func (p *stringCapturer) AddIntfField(_ string, v interface{}, fmt logcommon.LogFieldFormat) {
	p.v = fmt.ToString(v, "%v")
}
