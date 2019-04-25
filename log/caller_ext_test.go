//
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
//

package log_test

import (
	"bytes"
	"encoding/json"
	"testing"

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
	l = l.WithOutput(&b)

	l.Info("test")

	lf := logFields(t, b.Bytes())
	assert.Regexp(t, "^log/caller_ext_test.go:", lf.Caller, "log contains call place")
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
	l = l.WithFuncName(true)
	l = l.WithOutput(&b)

	l.Info("test")

	lf := logFields(t, b.Bytes())
	assert.Regexp(t, "^log/caller_ext_test.go:", lf.Caller, "log contains proper caller place")
	assert.NotContains(t, "github.com/insolar/insolar", lf.Caller, "log not contains package name")
	assert.Equal(t, "TestExtLog_ZerologCallerWithFunc", lf.Func, "log contains func name")
}

func TestExtLog_GlobalCaller(t *testing.T) {
	gl := log.GlobalLogger
	defer func() { log.GlobalLogger = gl }()

	var b bytes.Buffer
	log.GlobalLogger = log.GlobalLogger.WithOutput(&b)

	log.SetLevel("info")
	log.Info("test")

	lf := logFields(t, b.Bytes())
	assert.Regexp(t, "^log/caller_ext_test.go:", lf.Caller, "log contains proper call place")
	assert.Equal(t, "", lf.Func, "log not contains func name")
}

// this test result depends on test name!
func TestExtLog_GlobalCallerWithFunc(t *testing.T) {
	gl := log.GlobalLogger
	defer func() { log.GlobalLogger = gl }()

	var b bytes.Buffer
	log.GlobalLogger = log.GlobalLogger.WithOutput(&b)
	log.GlobalLogger = log.GlobalLogger.WithFuncName(true)

	log.SetLevel("info")
	log.Info("test")

	lf := logFields(t, b.Bytes())
	assert.Regexp(t, "^log/caller_ext_test.go:", lf.Caller, "log contains call place")
	assert.Equal(t, "TestExtLog_GlobalCallerWithFunc", lf.Func, "log contains call place")
}
