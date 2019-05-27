///
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
///

package inslogger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"runtime"
	"strconv"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Beware, test results there depends on test file and package names!

var (
	pkgRelName   = "instrumentation/inslogger/"
	testFileName = "inslogger_ext_test.go"
	callerRe     = "^" + pkgRelName + testFileName + ":"
)

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

func logFromContextWithSkip(ctx context.Context) insolar.Logger {
	// skip testing.go & runtime wrappers
	return inslogger.FromContext(ctx).WithSkipFrameCount(-2)
}

func TestExt_Global(t *testing.T) {

	l := logFromContextWithSkip(context.Background())
	l, _ = l.WithLevel("info")
	var b bytes.Buffer
	l = l.WithOutput(&b)

	_, _, line, _ := runtime.Caller(0)
	l.Info("test")

	lf := logFields(t, b.Bytes())
	assert.Regexp(t, callerRe+strconv.Itoa(line+1), lf.Caller, "log contains call place")
	assert.Equal(t, "", lf.Func, "log not contains func name")
}

func TestExt_Global_WithFunc(t *testing.T) {
	l := logFromContextWithSkip(context.Background())
	var b bytes.Buffer
	l = l.WithOutput(&b)
	l = l.WithFuncName(true)
	l, _ = l.WithLevel("info")

	_, _, line, _ := runtime.Caller(0)
	l.Info("test")

	lf := logFields(t, b.Bytes())
	assert.Regexp(t, callerRe+strconv.Itoa(line+1), lf.Caller, "log contains call place")
	assert.Equal(t, t.Name(), lf.Func, "log not contains func name")
}

func TestExt_Log(t *testing.T) {
	logPut, err := log.NewLog(configuration.Log{
		Level:     "info",
		Adapter:   "zerolog",
		Formatter: "json",
	})
	require.NoError(t, err, "log creation")
	ctx := inslogger.SetLogger(context.TODO(), logPut)

	l := inslogger.FromContext(ctx)
	var b bytes.Buffer
	l = l.WithOutput(&b)

	_, _, line, _ := runtime.Caller(0)
	l.Info("test")

	lf := logFields(t, b.Bytes())
	assert.Regexp(t, callerRe+strconv.Itoa(line+1), lf.Caller, "log contains call place")
	assert.Equal(t, "", lf.Func, "log not contains func name")
}

func TestExt_Log_WithFunc(t *testing.T) {
	logPut, err := log.NewLog(configuration.Log{
		Level:     "info",
		Adapter:   "zerolog",
		Formatter: "json",
	})
	require.NoError(t, err, "log creation")
	ctx := inslogger.SetLogger(context.TODO(), logPut)

	l := inslogger.FromContext(ctx)
	var b bytes.Buffer
	l = l.WithOutput(&b)
	l = l.WithFuncName(true)
	l, _ = l.WithLevel("info")

	_, _, line, _ := runtime.Caller(0)
	l.Info("test")

	lf := logFields(t, b.Bytes())
	assert.Regexp(t, callerRe+strconv.Itoa(line+1), lf.Caller,
		"log contains call place")
	assert.Equal(t, t.Name(), lf.Func, "log not contains func name")
}

func TestExt_Log_SubCall(t *testing.T) {
	logPut, err := log.NewLog(configuration.Log{
		Level:     "info",
		Adapter:   "zerolog",
		Formatter: "json",
	})
	require.NoError(t, err, "log creation")
	ctx := inslogger.SetLogger(context.TODO(), logPut.WithFuncName(true))

	lf, line := logCaller(ctx, t)
	assert.Regexp(t, callerRe+line, lf.Caller, "log contains call place")
	assert.Equal(t, "logCaller", lf.Func, "log not contains func name")
}

func logCaller(ctx context.Context, t *testing.T) (loggerField, string) {
	l := inslogger.FromContext(ctx)
	var b bytes.Buffer
	l = l.WithOutput(&b)

	_, _, line, _ := runtime.Caller(0)
	l.Info("test")

	return logFields(t, b.Bytes()), strconv.Itoa(line + 1)
}

func TestExt_Global_SubCall(t *testing.T) {
	lf, line := logCallerGlobal(context.Background(), t)
	assert.Regexp(t, callerRe+line, lf.Caller, "log contains call place")
}

func logCallerGlobal(ctx context.Context, t *testing.T) (loggerField, string) {
	l := inslogger.FromContext(ctx).WithSkipFrameCount(-2)
	l, _ = l.WithLevel("info")

	var b bytes.Buffer
	l = l.WithOutput(&b)

	_, _, line, _ := runtime.Caller(0)
	l.Info("test")
	return logFields(t, b.Bytes()), strconv.Itoa(line + 1)
}
