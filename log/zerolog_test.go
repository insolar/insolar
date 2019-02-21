/*
 *    Copyright 2019 Insolar Technologies
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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/configuration"
)

func TestZeroLogAdapter_CallerInfo(t *testing.T) {
	log, err := NewLog(configuration.Log{Level: "info", Adapter: "zerolog", Formatter: "json"})
	require.NoError(t, err)
	require.NotNil(t, log)

	var buf bytes.Buffer
	log.SetOutput(&buf)

	log.Error("test")

	require.Contains(t, buf.String(), "zerolog_test.go:36")
}
