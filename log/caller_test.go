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

package log

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// XXX: Beware, test results there depends on test file name (caller_test.go)!

type loggerField struct {
	Caller string
	Func   string
}

func logFields(t *testing.T, b []byte) loggerField {
	var lf loggerField
	err := json.Unmarshal(b, &lf)
	require.NoError(t, err)
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
	l = l.WithOutput(&b)

	l.Info("test")

	lf := logFields(t, b.Bytes())
	assert.Regexp(t, "^log/caller_test.go:", lf.Caller, "log contains call place")
	assert.NotContains(t, "github.com/insolar/insolar", lf.Caller, "log not contains package name")
}

// XXX: test result depends on test name
func TestLog_ZerologCallerWithFunc(t *testing.T) {
	l, err := NewLog(configuration.Log{
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
	assert.Regexp(t, "^log/caller_test.go:", lf.Caller, "log contains call place")
	assert.NotContains(t, "github.com/insolar/insolar", lf.Caller, "log not contains package name")
	assert.Equal(t, "TestLog_ZerologCallerWithFunc", lf.Func, "log contains call place")
}

func TestLog_GlobalCaller(t *testing.T) {
	var b bytes.Buffer
	GlobalLogger = GlobalLogger.WithOutput(&b)

	Info("test")

	lf := logFields(t, b.Bytes())
	assert.Regexp(t, "^log/caller_test.go:", lf.Caller, "log contains call place")
}

// XXX: test result depends on test name
func TestLog_GlobalCallerWithFunc(t *testing.T) {
	var b bytes.Buffer
	GlobalLogger = GlobalLogger.WithOutput(&b)
	GlobalLogger = GlobalLogger.WithFuncName(true)

	Info("test")

	lf := logFields(t, b.Bytes())
	assert.Regexp(t, "^log/caller_test.go:", lf.Caller, "log contains call place")
	assert.Equal(t, "TestLog_GlobalCallerWithFunc", lf.Func, "log contains call place")
}
