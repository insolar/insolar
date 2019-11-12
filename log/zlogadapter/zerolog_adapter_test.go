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
	"sync"
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
	require.Contains(t, s, "zerolog_adapter_test.go:47")
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
	require.Contains(t, s, "zerolog_adapter_test.go:73")
}

func TestZeroLogAdapter_InheritFields(t *testing.T) {
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
	log, err = log.Copy().WithOutput(&buf).WithCaller(insolar.CallerField).WithField("field1", "value1").Build()
	require.NoError(t, err)

	log = log.WithField("field2", "value2")

	var buf2 bytes.Buffer
	log, err = log.Copy().WithOutput(&buf2).Build()
	require.NoError(t, err)

	log.Error("test")

	s := buf2.String()
	require.Contains(t, s, "value1")
	require.Contains(t, s, "value2")
}

func TestZeroLogAdapter_Fatal(t *testing.T) {
	zc := logadapter.Config{}

	var buf bytes.Buffer
	wg := sync.WaitGroup{}
	wg.Add(1)
	zc.BareOutput = logadapter.BareOutput{
		Writer: &buf,
		FlushFn: func() error {
			wg.Done()
			select {} // hang up to stop zerolog's call to os.Exit
		},
	}
	zc.Output = logadapter.OutputConfig{Format: insolar.DefaultLogFormat}
	zc.MsgFormat = logadapter.GetDefaultLogMsgFormatter()
	zc.Instruments.SkipFrameCountBaseline = 0

	zb := logadapter.NewBuilder(zerologFactory{}, zc, insolar.InfoLevel)
	log, err := zb.Build()

	require.NoError(t, err)
	require.NotNil(t, log)

	log.Error("errorMsgText")
	go log.Fatal("fatalMsgText") // it will hang on flush
	wg.Wait()

	s := buf.String()
	require.Contains(t, s, "errorMsgText")
	require.Contains(t, s, "fatalMsgText")
}

func TestZeroLogAdapter_Panic(t *testing.T) {
	zc := logadapter.Config{}

	var buf bytes.Buffer
	wg := sync.WaitGroup{}
	wg.Add(1)
	zc.BareOutput = logadapter.BareOutput{
		Writer: &buf,
		FlushFn: func() error {
			wg.Done()
			return nil
		},
	}
	zc.Output = logadapter.OutputConfig{Format: insolar.DefaultLogFormat}
	zc.MsgFormat = logadapter.GetDefaultLogMsgFormatter()
	zc.Instruments.SkipFrameCountBaseline = 0

	zb := logadapter.NewBuilder(zerologFactory{}, zc, insolar.InfoLevel)
	log, err := zb.Build()

	require.NoError(t, err)
	require.NotNil(t, log)

	log.Error("errorMsgText")
	require.PanicsWithValue(t, "panicMsgText", func() {
		log.Panic("panicMsgText")
	})
	wg.Wait()
	log.Error("errorNextMsgText")

	s := buf.String()
	require.Contains(t, s, "errorMsgText")
	require.Contains(t, s, "panicMsgText")
	require.Contains(t, s, "errorNextMsgText")
}
