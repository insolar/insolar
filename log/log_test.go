/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package log

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/testutils"
)

func capture(f func()) string {
	var buf bytes.Buffer
	SetOutput(&buf)
	f()
	SetOutput(os.Stderr)
	return buf.String()
}

func assertHelloWorld(t *testing.T, out string) {
	assert.Contains(t, out, " msg=HelloWorld")
}

func TestLog_GlobalLogger(t *testing.T) {

	assert.NoError(t, SetLevel("debug"))

	assertHelloWorld(t, capture(func() { Debug("HelloWorld") }))
	assertHelloWorld(t, capture(func() { Debugln("HelloWorld") }))
	assertHelloWorld(t, capture(func() { Debugf("%s", "HelloWorld") }))
	assertHelloWorld(t, capture(func() { Info("HelloWorld") }))
	assertHelloWorld(t, capture(func() { Infoln("HelloWorld") }))
	assertHelloWorld(t, capture(func() { Infof("%s", "HelloWorld") }))

	assertHelloWorld(t, capture(func() { Warn("HelloWorld") }))
	assertHelloWorld(t, capture(func() { Warnln("HelloWorld") }))
	assertHelloWorld(t, capture(func() { Warnf("%s", "HelloWorld") }))

	assertHelloWorld(t, capture(func() { Error("HelloWorld") }))
	assertHelloWorld(t, capture(func() { Errorln("HelloWorld") }))
	assertHelloWorld(t, capture(func() { Errorf("%s", "HelloWorld") }))

	assert.Panics(t, func() { Panic("HelloWorld") })
	assert.Panics(t, func() { Panicln("HelloWorld") })
	assert.Panics(t, func() { Panicf("%s", "HelloWorld") })

	// can't catch os.exit() to test Fatal
	// Fatal("HelloWorld")
	// Fatalln("HelloWorld")
	// Fatalf("%s", "HelloWorld")
}

func TestLog_NewLog_Config(t *testing.T) {
	invalidtests := map[string]configuration.Log{
		"InvalidAdapter": configuration.Log{Level: "Debug", Adapter: "invalid"},
		"InvalidLevel":   configuration.Log{Level: "Invalid", Adapter: "logrus"},
	}

	for name, test := range invalidtests {
		t.Run(name, func(t *testing.T) {
			logger, err := NewLog(test)
			assert.Nil(t, logger)
			assert.Error(t, err)
		})
	}

	validtests := map[string]configuration.Log{
		"WithAdapter": configuration.Log{Level: "Debug", Adapter: "logrus"},
	}
	for name, test := range validtests {
		t.Run(name, func(t *testing.T) {
			logger, err := NewLog(test)
			assert.NotNil(t, logger)
			assert.NoError(t, err)
		})
	}
}

func TestLog_GlobalLogger_Level(t *testing.T) {
	got := GetLevel()
	assert.NoError(t, SetLevel("error"))
	assert.Error(t, SetLevel("errorrr"))
	assert.Equal(t, "error", GetLevel())
	assert.NoError(t, SetLevel(got))
	assert.Equal(t, got, GetLevel())
}

func TestLog_AddFields(t *testing.T) {
	errtxt1 := "~CHECK_ERROR_OUTPUT_WITH_FIELDS~"
	errtxt2 := "~CHECK_ERROR_OUTPUT_WITHOUT_FIELDS~"

	var (
		fieldname  = "TraceID"
		fieldvalue = "Trace100500"
	)
	tt := []struct {
		name    string
		fieldfn func(la logrusAdapter) core.Logger
	}{
		{
			name: "WithFields",
			fieldfn: func(la logrusAdapter) core.Logger {
				fields := map[string]interface{}{fieldname: fieldvalue}
				return la.WithFields(fields)
			},
		},
		{
			name: "WithField",
			fieldfn: func(la logrusAdapter) core.Logger {
				return la.WithField(fieldname, fieldvalue)
			},
		},
	}

	for _, tItem := range tt {
		t.Run(tItem.name, func(t *testing.T) {
			recorder := testutils.NewRecoder()

			la := newLogrusAdapter()
			la.SetOutput(recorder)

			tItem.fieldfn(la).Error(errtxt1)
			la.Error(errtxt2)

			recitems := recorder.Items()
			assert.Contains(t, recitems[0], errtxt1)
			assert.Contains(t, recitems[1], errtxt2)
			assert.Contains(t, recitems[0], fieldvalue)
			assert.NotContains(t, recitems[1], fieldvalue)
		})
	}
}
