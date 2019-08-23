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

// +build slowtest

package log

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

var logLevelEnvVarName = "INSOLAR_LOG_LEVEL"

func testWithEnvVar(t *testing.T) {
	val := strings.ToLower(os.Getenv(logLevelEnvVarName))

	assert.Containsf(t,
		capture(func() { Warn("HelloWorld") }),
		"HelloWorld", "Warn on level=%v by is set", val)
	assert.Containsf(t,
		capture(func() { Info("HelloWorld") }),
		"HelloWorld", "Info on level=%v is set", val)

	if val == "debug" {
		assert.Containsf(t,
			capture(func() { Debug("HelloWorld") }),
			"HelloWorld", "Debug should work on level %v", val)
	} else {
		assert.NotContainsf(t, capture(func() { Debug("HelloWorld") }),
			"HelloWorld", "Debug should not work on level %v", val)
	}
}

func TestLog_GlobalLogger_Env(t *testing.T) {
	if os.Getenv("__TestLoggerWithEnv__") == "1" {
		testWithEnvVar(t)
		return
	}

	levels := []string{"", "debug"}
	for _, val := range levels {
		name := val
		if name == "" {
			name = "empty"
		}
		t.Run(name, func(t *testing.T) {
			env := []string{"__TestLoggerWithEnv__=1"}
			for _, e := range os.Environ() {
				if strings.HasPrefix(e, logLevelEnvVarName+"=") {
					e = logLevelEnvVarName + "=" + val
				}
				env = append(env, e)
			}

			cmd := exec.Command(os.Args[0], "-test.run=TestLog_GlobalLogger_Env")
			cmd.Env = env
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if e, ok := err.(*exec.ExitError); ok && !e.Success() {
				exitCode := 0
				if status, ok := e.Sys().(syscall.WaitStatus); ok {
					exitCode = status.ExitStatus()
				}
				t.Fatalf("%v with env var %v=%v failed (status=%v, code=%v)",
					os.Args[0], logLevelEnvVarName, val, e.String(), exitCode)
			}
		})
	}
}
