///
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
///

package zlogadapter

import (
	"bytes"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log/logadapter"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZeroLogAdapter_CallerInfoWithFunc(t *testing.T) {
	pCfg := logadapter.ParsedLogConfig{
		OutputType: insolar.DefaultLogOutput,
		LogLevel:   insolar.InfoLevel,
		Output: logadapter.OutputConfig{
			Format: insolar.DefaultLogFormat,
		},
	}
	msgFmt := logadapter.GetDefaultLogMsgFormatter()

	log, err := NewZerologAdapter(pCfg, msgFmt)
	require.NoError(t, err)
	require.NotNil(t, log)

	var buf bytes.Buffer
	log, err = log.Copy().WithOutput(&buf).WithCaller(insolar.CallerFieldWithFuncName).Build()
	require.NoError(t, err)

	log.Error("test")

	s := buf.String()
	require.Contains(t, s, "zerolog_adapter_test.go:46")
	require.Contains(t, s, "TestZeroLogAdapter_CallerInfoWithFunc")
}

func TestZeroLogAdapter_CallerInfo(t *testing.T) {
	pCfg := logadapter.ParsedLogConfig{
		OutputType: insolar.DefaultLogOutput,
		LogLevel:   insolar.InfoLevel,
		Output: logadapter.OutputConfig{
			Format: insolar.DefaultLogFormat,
		},
	}
	msgFmt := logadapter.GetDefaultLogMsgFormatter()

	log, err := NewZerologAdapter(pCfg, msgFmt)

	require.NoError(t, err)
	require.NotNil(t, log)

	var buf bytes.Buffer
	log, err = log.Copy().WithOutput(&buf).WithCaller(insolar.CallerField).Build()
	require.NoError(t, err)

	log.Error("test")

	s := buf.String()
	require.Contains(t, s, "zerolog_adapter_test.go:72")
}
