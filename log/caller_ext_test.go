// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package log_test

import (
	"bytes"
	"encoding/json"
	"runtime"
	"strconv"
	"testing"

	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestExtLog_ZerologCaller(t *testing.T) {
	l, err := log.NewLog(configuration.Log{
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
	assert.Regexp(t, "log/caller_ext_test.go:"+strconv.Itoa(line+1), lf.Caller, "log contains call place")
	assert.NotContains(t, "github.com/insolar/insolar", lf.Caller, "log not contains package name")
	assert.Equal(t, "", lf.Func, "log not contains func name")
}

// this test result depends on test name!
func TestExtLog_ZerologCallerWithFunc(t *testing.T) {
	l, err := log.NewLog(configuration.Log{
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
	assert.Regexp(t, "log/caller_ext_test.go:"+strconv.Itoa(line+1), lf.Caller, "log contains proper caller place")
	assert.NotContains(t, "github.com/insolar/insolar", lf.Caller, "log not contains package name")
	assert.Equal(t, "TestExtLog_ZerologCallerWithFunc", lf.Func, "log contains func name")
}

func TestExtLog_GlobalCaller(t *testing.T) {
	defer log.SaveGlobalLogger()()

	var b bytes.Buffer
	gl2, err := log.GlobalLogger().Copy().WithOutput(&b).WithCaller(insolar.CallerField).Build()
	require.NoError(t, err)
	log.SetGlobalLogger(gl2)

	err = log.SetLevel("info")
	require.NoError(t, err)

	_, _, line, _ := runtime.Caller(0)
	log.Info("test")
	log.Debug("test2shouldNotBeThere")

	s := b.String()
	lf := logFields(t, []byte(s))
	assert.Regexp(t, "log/caller_ext_test.go:"+strconv.Itoa(line+1), lf.Caller, "log contains proper call place")
	assert.Equal(t, "", lf.Func, "log not contains func name")
	assert.NotContains(t, s, "test2shouldNotBeThere")
}

func TestExtLog_GlobalCallerWithFunc(t *testing.T) {
	defer log.SaveGlobalLogger()()

	var b bytes.Buffer
	gl2, err := log.GlobalLogger().Copy().WithOutput(&b).WithCaller(insolar.CallerFieldWithFuncName).Build()
	require.NoError(t, err)
	log.SetGlobalLogger(gl2)

	err = log.SetLevel("info")
	require.NoError(t, err)

	_, _, line, _ := runtime.Caller(0)
	log.Info("test")
	log.Debug("test2shouldNotBeThere")

	s := b.String()
	lf := logFields(t, []byte(s))
	assert.Regexp(t, "log/caller_ext_test.go:"+strconv.Itoa(line+1), lf.Caller, "log contains proper call place")
	assert.Equal(t, "TestExtLog_GlobalCallerWithFunc", lf.Func, "log contains func name")
	assert.NotContains(t, s, "test2shouldNotBeThere")
}
