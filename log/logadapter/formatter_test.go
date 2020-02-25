// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package logadapter

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/args"
)

type stringerStruct struct {
	s string
}

func (v stringerStruct) String() string {
	return "wrong" // must take LogString() first
}

func (v stringerStruct) LogString() string {
	return v.s
}

type stringerRefStruct struct {
	s string
}

func (p *stringerRefStruct) String() string {
	return p.s
}

type stubStruct struct {
}

const sampleStructAsString = "f0:  99:string,f1:999:int,f2:test_raw,f3:test2:string,f4:nil,f5:stringer_test:string,f6:func_result:string,f7:stringerVal:string,f8:stringerRef:string,f9:nil,f10:{}:stubStruct,msg:message title"

func createSampleStruct() interface{} {
	s := "test2"
	return struct {
		msg string
		f0  int `fmt:"%4d"`
		f1  int
		f2  string `raw:"%s"`
		f3  *string
		f4  *string
		f5  interface{}
		f6  func() string
		f7  stringerStruct
		f8  *stringerRefStruct
		f9  *stringerRefStruct
		f10 stubStruct // no special handling
	}{
		"message title",
		99, 999, "test_raw", &s, nil,
		args.LazyFmt("stringer_test"),
		func() string { return "func_result" },
		stringerStruct{"stringerVal"},
		&stringerRefStruct{"stringerRef"},
		nil,
		stubStruct{},
	}
}

func TestTryLogObject_Many(t *testing.T) {
	f := GetDefaultLogMsgFormatter()

	require.Equal(t,
		"{message title 99} 888",
		f.testTryLogObject(struct {
			msg string
			f1  int
		}{
			"message title",
			99,
		}, 888))
}

func TestTryLogObject_Str(t *testing.T) {
	f := GetDefaultLogMsgFormatter()

	require.Equal(t, "text", f.testTryLogObject("text"))
	s := "text"
	require.Equal(t, "text", f.testTryLogObject(s))
	require.Equal(t, "text", f.testTryLogObject(&s))
	ps := &s
	require.Equal(t, "text", f.testTryLogObject(ps))
	ps = nil
	require.Equal(t, "<nil>", f.testTryLogObject(ps))

	require.Equal(t, "text", f.testTryLogObject(func() string { return "text" }))
}

func TestTryLogObject_SingleUnnamed(t *testing.T) {
	f := GetDefaultLogMsgFormatter()

	require.Equal(t, sampleStructAsString, f.testTryLogObject(createSampleStruct()))
}

func TestTryLogObject_SingleNamed(t *testing.T) {
	f := GetDefaultLogMsgFormatter()

	type SomeType struct {
		i int
	}

	require.Equal(t,
		"{7}",
		f.testTryLogObject(SomeType{7}))
}

var _ insolar.LogObject = SomeLogObjectValue{}

type SomeLogObjectValue struct {
	IntVal int
	Msg    string
}

func (SomeLogObjectValue) GetLogObjectMarshaller() insolar.LogObjectMarshaller {
	return nil
}

var _ insolar.LogObject = &SomeLogObjectPtr{}

type SomeLogObjectPtr struct {
	IntVal int
	Msg    string
}

func (*SomeLogObjectPtr) GetLogObjectMarshaller() insolar.LogObjectMarshaller {
	return nil
}

var _ insolar.LogObject = SomeLogObjectWithTemplate{}
var _ insolar.LogObject = &SomeLogObjectWithTemplate{}

type SomeLogObjectWithTemplate struct {
	*insolar.LogObjectTemplate
	IntVal int
}

type SomeLogObjectWithTemplateAndMsg struct {
	*insolar.LogObjectTemplate `txt:"TemplateAndMsg"`
	IntVal                     int
}

type SomeLogObjectWithTemplateAndMsg2 struct {
	*insolar.LogObjectTemplate
	IntVal int
	_      struct{} `txt:"TemplateAndMsg"`
}

func TestTryLogObject_SingleLogObject(t *testing.T) {
	f := GetDefaultLogMsgFormatter()

	require.Equal(t,
		"IntVal:7:int,msg:msgText",
		f.testTryLogObject(SomeLogObjectValue{7, "msgText"}))

	require.Equal(t,
		"IntVal:7:int,msg:msgText",
		f.testTryLogObject(&SomeLogObjectPtr{7, "msgText"}))

	require.Equal(t,
		"{7 msgText}",
		f.testTryLogObject(SomeLogObjectPtr{7, "msgText"})) // function has ptr receiver and is not taken

	require.Equal(t,
		"IntVal:7:int,msg:",
		f.testTryLogObject(&SomeLogObjectWithTemplate{nil, 7}))

	require.Equal(t,
		"IntVal:7:int,msg:",
		f.testTryLogObject(SomeLogObjectWithTemplate{nil, 7}))

	require.Equal(t,
		"IntVal:7:int,msg:TemplateAndMsg",
		f.testTryLogObject(SomeLogObjectWithTemplateAndMsg{nil, 7}))

	require.Equal(t,
		"IntVal:7:int,msg:TemplateAndMsg",
		f.testTryLogObject(SomeLogObjectWithTemplateAndMsg2{nil, 7, struct{}{}}))
}

type SomeLogObjectWithMsg struct {
	*insolar.LogObjectTemplate
	IntVal int    `opt:""`
	Msg    string `txt:"fixedObjectMessage"`
}

func TestTryLogObject_ConstMsg(t *testing.T) {
	f := GetDefaultLogMsgFormatter()

	require.Equal(t,
		"Data:0:int,msg:inlineConstantText",
		f.testTryLogObject(struct {
			Data int
			Msg  string `txt:"inlineConstantText"`
		}{}))

	require.Equal(t,
		"IntVal:78:int,msg:fixedObjectMessage",
		f.testTryLogObject(SomeLogObjectWithMsg{IntVal: 78}))
}

func TestTryLogObject_Optionals(t *testing.T) {
	f := GetDefaultLogMsgFormatter()

	require.Equal(t,
		"msg:fixedObjectMessage",
		f.testTryLogObject(SomeLogObjectWithMsg{}))

	require.Equal(t,
		"msg:inlineConstantText",
		f.testTryLogObject(struct {
			Data int    `anotherTag:"" opt:""`
			Msg  string `txt:"inlineConstantText"`
		}{}))
}

func (v MsgFormatConfig) testTryLogObject(a ...interface{}) string {
	m, s := v.FmtLogObject(a...)
	if m == nil {
		return s
	}
	o := output{}
	msg := m.MarshalTextLogObject(&o, nil)
	o.buf.WriteString("msg:")
	o.buf.WriteString(msg)
	return o.buf.String()
}

var _ insolar.LogObjectWriter = &output{}

type output struct {
	buf strings.Builder
}

func (p *output) AddStrField(key string, v string) {
	p.AddField(key, v)
}

func (p *output) AddField(k string, v interface{}) {
	if v == nil {
		p.buf.WriteString(fmt.Sprintf("%s:nil,", k))
	} else {
		p.buf.WriteString(fmt.Sprintf("%s:%v:%s,", k, v, reflect.TypeOf(v).Name()))
	}
}

func (p *output) AddRawJSON(k string, b []byte) {
	p.buf.WriteString(fmt.Sprintf("%s:%s,", k, b))
}
