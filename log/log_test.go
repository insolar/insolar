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
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/stretchr/testify/assert"
)

func TestLog_GlobalLogger(t *testing.T) {

	assert.NoError(t, SetLevel("debug"))

	Debug("HelloWorld")
	Debugln("HelloWorld")
	Debugf("%s", "HelloWorld")

	Info("HelloWorld")
	Infoln("HelloWorld")
	Infof("%s", "HelloWorld")

	Warn("HelloWorld")
	Warnln("HelloWorld")
	Warnf("%s", "HelloWorld")

	Error("HelloWorld")
	Errorln("HelloWorld")
	Errorf("%s", "HelloWorld")

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
