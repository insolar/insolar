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
	"fmt"
	"testing"

	"github.com/insolar/insolar/configuration"
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

func TestExt_Global(t *testing.T) {

	l := inslogger.FromContext(context.TODO())
	var b bytes.Buffer
	l = l.WithOutput(&b)

	l.Info("test")

	fmt.Println(b.String())

	lf := logFields(t, b.Bytes())
	assert.Regexp(t, callerRe, lf.Caller, "log contains call place")
	assert.Equal(t, "", lf.Func, "log not contains func name")
}

func TestExt_Global_WithFunc(t *testing.T) {
	l := inslogger.FromContext(context.TODO())
	var b bytes.Buffer
	l = l.WithOutput(&b)
	l = l.WithFuncName(true)

	l.Info("test")

	fmt.Println(b.String())

	lf := logFields(t, b.Bytes())
	assert.Regexp(t, callerRe, lf.Caller, "log contains call place")
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

	l.Info("test")

	fmt.Println(b.String())

	lf := logFields(t, b.Bytes())
	assert.Regexp(t, callerRe, lf.Caller, "log contains call place")
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

	l.Info("test")
	// l.Debug("test")

	fmt.Println(b.String())

	lf := logFields(t, b.Bytes())
	assert.Regexp(t, callerRe, lf.Caller, "log contains call place")
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

	lf := logCaller(ctx, t)
	assert.Regexp(t, callerRe, lf.Caller, "log contains call place")
	assert.Equal(t, "logCaller", lf.Func, "log not contains func name")
}

func TestExt_Global_SubCall(t *testing.T) {
	lf := logCaller(context.Background(), t)
	assert.Regexp(t, callerRe, lf.Caller, "log contains call place")
}

func logCaller(ctx context.Context, t *testing.T) loggerField {
	l := inslogger.FromContext(ctx)
	var b bytes.Buffer
	l = l.WithOutput(&b)
	l.Info("test")
	// fmt.Print("logCaller:\n", b.String())
	return logFields(t, b.Bytes())
}
