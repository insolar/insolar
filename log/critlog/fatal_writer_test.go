// Copyright 2020 Insolar Network Ltd.
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

package critlog

import (
	"strings"
	"testing"

	"github.com/insolar/insolar/log/logoutput"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testWriter struct {
	strings.Builder
	closed, flushed, noFlush bool
}

func (w *testWriter) Close() error {
	w.closed = true
	return nil
}

func (w *testWriter) Flush() error {
	if w.closed {
		return errors.New("closed")
	}
	if w.noFlush {
		return errors.New("unsupported")
	}
	w.flushed = true
	return nil
}

func TestFatalDirectWriter_mute_on_fatal(t *testing.T) {
	tw := testWriter{}
	writer := NewFatalDirectWriter(logoutput.NewAdapter(&tw, false, nil, func() error {
		_ = tw.Flush()
		panic("fatal")
	}))
	// We don't want to lock the writer on fatal in tests.
	var err error

	_, err = writer.LogLevelWrite(insolar.WarnLevel, []byte("WARN must pass\n"))
	require.NoError(t, err)
	_, err = writer.LogLevelWrite(insolar.ErrorLevel, []byte("ERROR must pass\n"))
	require.NoError(t, err)

	assert.False(t, tw.flushed)
	_, err = writer.LogLevelWrite(insolar.PanicLevel, []byte("PANIC must pass\n"))
	require.NoError(t, err)
	assert.True(t, tw.flushed)

	tw.flushed = false
	require.PanicsWithValue(t, "fatal", func() {
		_, _ = writer.LogLevelWrite(insolar.FatalLevel, []byte("FATAL must pass\n"))
	})
	assert.True(t, tw.flushed)
	assert.False(t, tw.closed)

	// MUST hang. Tested by logoutput.Adapter
	//_, _ = writer.LogLevelWrite(insolar.WarnLevel, []byte("WARN must NOT pass\n"))
	//_, _ = writer.LogLevelWrite(insolar.ErrorLevel, []byte("ERROR must NOT pass\n"))
	//_, _ = writer.LogLevelWrite(insolar.PanicLevel, []byte("PANIC must NOT pass\n"))
	testLog := tw.String()
	assert.Contains(t, testLog, "WARN must pass")
	assert.Contains(t, testLog, "ERROR must pass")
	assert.Contains(t, testLog, "FATAL must pass")
	//assert.NotContains(t, testLog, "must NOT pass")
}

func TestFatalDirectWriter_close_on_fatal_without_flush(t *testing.T) {
	tw := testWriter{}
	tw.noFlush = true

	writer := NewFatalDirectWriter(logoutput.NewAdapter(&tw, false, nil, nil))
	var err error

	_, err = writer.LogLevelWrite(insolar.WarnLevel, []byte("WARN must pass\n"))
	require.NoError(t, err)

	_, err = writer.LogLevelWrite(insolar.FatalLevel, []byte("FATAL must pass\n"))
	require.NoError(t, err)
	assert.False(t, tw.flushed)
	assert.True(t, tw.closed)

	// MUST hang. Tested by logoutput.Adapter
	//_, _ = writer.LogLevelWrite(insolar.WarnLevel, []byte("WARN must NOT pass\n"))
	//_, _ = writer.LogLevelWrite(insolar.ErrorLevel, []byte("ERROR must NOT pass\n"))
	//_, _ = writer.LogLevelWrite(insolar.PanicLevel, []byte("PANIC must NOT pass\n"))
	testLog := tw.String()
	assert.Contains(t, testLog, "WARN must pass")
	assert.Contains(t, testLog, "FATAL must pass")
	//assert.NotContains(t, testLog, "must NOT pass")
}
