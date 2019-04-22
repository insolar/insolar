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
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// XXX: beware test result depends on it's file name.
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
	assert.Contains(t, b.String(), "log/caller_test.go:",
		"log contains call place")
	assert.NotContains(t, b.String(), "github.com/insolar/insolar",
		"log not contains package name")
}
