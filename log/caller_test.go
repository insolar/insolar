package log

import (
	"bytes"
	"encoding/json"
	"runtime"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
)

// Beware, test results there depends on test file name (caller_test.go)!

type loggerField struct {
	Caller string
	Func   string
}

func logFields(t *testing.T, b []byte) loggerField {
	var lf loggerField
	err := json.Unmarshal(b, &lf)
	require.NoErrorf(t, err, "failed decode: '%v'", string(b))
	return lf
}

func TestLog_ZerologCaller(t *testing.T) {
	l, err := NewLog(configuration.Log{
		Level:     "info",
		Adapter:   "zerolog",
		Formatter: "json",
	})
	require.NoError(t, err, "log creation")

	var b bytes.Buffer
	l, err = l.Copy().WithOutput(&b).WithCaller(insolar.CallerField).Build()
	require.NoError(t, err)

	_, _, line, _ := runtime.Caller(0)
	l.Info("test")

	lf := logFields(t, b.Bytes())
	assert.Regexp(t, "log/caller_test.go:"+strconv.Itoa(line+1), lf.Caller, "log contains call place")
	assert.NotContains(t, "github.com/insolar/insolar", lf.Caller, "log not contains package name")
	assert.Equal(t, "", lf.Func, "log not contains func name")
}

// this test result depends on test name!
func TestLog_ZerologCallerWithFunc(t *testing.T) {
	l, err := NewLog(configuration.Log{
		Level:     "info",
		Adapter:   "zerolog",
		Formatter: "json",
	})
	require.NoError(t, err, "log creation")

	var b bytes.Buffer
	l, err = l.Copy().WithOutput(&b).WithCaller(insolar.CallerFieldWithFuncName).Build()
	require.NoError(t, err)

	_, _, line, _ := runtime.Caller(0)
	l.Info("test")

	lf := logFields(t, b.Bytes())
	assert.Regexp(t, "log/caller_test.go:"+strconv.Itoa(line+1), lf.Caller, "log contains proper caller place")
	assert.NotContains(t, "github.com/insolar/insolar", lf.Caller, "log not contains package name")
	assert.Equal(t, "TestLog_ZerologCallerWithFunc", lf.Func, "log contains func name")
}

func TestLog_GlobalCaller(t *testing.T) {
	defer SaveGlobalLogger()()

	var b bytes.Buffer
	gl2, err := GlobalLogger().Copy().WithOutput(&b).WithCaller(insolar.CallerField).Build()
	require.NoError(t, err)
	SetGlobalLogger(gl2)

	err = SetLevel("info")
	require.NoError(t, err)

	_, _, line, _ := runtime.Caller(0)
	Info("test")
	Debug("test2shouldNotBeThere")

	s := b.String()
	lf := logFields(t, []byte(s))
	assert.Regexp(t, "log/caller_test.go:"+strconv.Itoa(line+1), lf.Caller, "log contains proper call place")
	assert.Equal(t, "", lf.Func, "log not contains func name")
	assert.NotContains(t, s, "test2shouldNotBeThere")
}

func TestLog_GlobalCallerWithFunc(t *testing.T) {
	defer SaveGlobalLogger()()

	var b bytes.Buffer
	gl2, err := GlobalLogger().Copy().WithOutput(&b).WithCaller(insolar.CallerFieldWithFuncName).Build()
	require.NoError(t, err)
	SetGlobalLogger(gl2)

	err = SetLevel("info")
	require.NoError(t, err)

	_, _, line, _ := runtime.Caller(0)
	Info("test")
	Debug("test2shouldNotBeThere")

	s := b.String()
	lf := logFields(t, []byte(s))
	assert.Regexp(t, "log/caller_test.go:"+strconv.Itoa(line+1), lf.Caller, "log contains proper call place")
	assert.Equal(t, "TestLog_GlobalCallerWithFunc", lf.Func, "log contains func name")
	assert.NotContains(t, s, "test2shouldNotBeThere")
}
