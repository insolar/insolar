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

func formatStrValue(a interface{}) string {
	v := reflect.ValueOf(a)
	if ok, s, iv := tryReflectStrValue(v); ok {
		return s
	} else {
		if iv == nil {
			return fmt.Sprint(nil)
		}
		return fmt.Sprintf("%v", iv)
	}
}

func formatValue(a interface{}) interface{} {
	v := reflect.ValueOf(a)
	return prepareReflectValue(v)
}

func TestTryFormatStrValue(t *testing.T) {
	s := "test"
	require.Equal(t, "test", formatStrValue(&s))
	require.Equal(t, "<nil>", formatStrValue((*string)(nil)))
	require.Equal(t, "test", formatStrValue("test"))
	require.Equal(t, "test", formatStrValue(args.LazyFmt("%s", "test")))
	require.Equal(t, "test", formatStrValue(func() string { return "test" }))
	require.Equal(t, "123", formatStrValue(123))
	require.Equal(t, "<nil>", formatStrValue(nil))
}

func TestTryFormatValue(t *testing.T) {
	s := "test"
	require.Equal(t, "test", formatValue(&s))
	require.Equal(t, nil, formatValue((*string)(nil)))
	require.Equal(t, "test", formatValue("test"))
	require.Equal(t, "test", formatValue(args.LazyFmt("%s", "test")))
	require.Equal(t, "test", formatValue(func() string { return "test" }))
	require.Equal(t, 123, formatValue(123))
	require.Equal(t, nil, formatValue(nil))
}

func TestPrintFields(t *testing.T) {
	s := "test2"
	require.Equal(t,
		"f0:  99:string,f1:999:int,f2:test_raw,f3:test2:string,f4:nil,f5:stringer_test:string,f6:func_result:string,msg:message title",
		printFieldsOut(struct {
			msg string
			f0  int `fmt:"%4d"`
			f1  int
			f2  string `raw:"%s"`
			f3  *string
			f4  *string
			f5  interface{}
			f6  func() string
		}{
			"message title",
			99, 999, "test_raw", &s, nil,
			args.LazyFmt("stringer_test"),
			func() string { return "func_result" },
		}))
}

func printFieldsOut(v interface{}) string {
	o := output{}
	msg := printFields(reflect.ValueOf(v), &o)
	o.buf.WriteString("msg:")
	o.buf.WriteString(msg)
	return o.buf.String()
}

var _ insolar.LogObjectWriter = &output{}

type output struct {
	buf strings.Builder
}

func (p *output) AddField(k string, v interface{}) {
	if v == nil {
		p.buf.WriteString(fmt.Sprintf("%s:nil,", k))
	} else {
		p.buf.WriteString(fmt.Sprintf("%s:%v:%s,", k, v, reflect.TypeOf(v).Name()))
	}
}

func (output) AddFieldMap(m map[string]interface{}) {
	panic("implement me")
}

func (p *output) AddRawJSON(k string, b []byte) {
	p.buf.WriteString(fmt.Sprintf("%s:%s,", k, b))
}
