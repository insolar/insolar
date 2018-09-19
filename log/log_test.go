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

func TestLog_NewLog_InvalidConfig(t *testing.T) {
	cfg := configuration.NewLog()
	cfg.Adapter = "invalid"

	tests := map[string]configuration.Log{
		"InvalidAdapter": configuration.Log{Level: "Debug", Adapter: "invalid"},
		"InvalidLevel":   configuration.Log{Level: "Invalid", Adapter: "logrus"},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			logger, err := NewLog(test)
			assert.Nil(t, logger)
			assert.Error(t, err)
		})
	}
}
