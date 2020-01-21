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

package logoutput

import (
	"errors"
	"strings"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testWriter struct {
	strings.Builder
	closeCount, flushCount int
	noFlush                bool
}

func (w *testWriter) Close() error {
	w.closeCount++
	if w.closeCount > 1 {
		return errClosed
	}
	return nil
}

func (w *testWriter) Flush() error {
	w.flushCount++
	if w.closeCount > 1 {
		return errClosed
	}
	if w.noFlush {
		return errors.New("unsupported")
	}
	return nil
}

var _ insolar.LogLevelWriter = &testLevelWriter{}

type testLevelWriter struct {
	testWriter
}

func (p *testLevelWriter) Write([]byte) (int, error) {
	panic("unexpected")
}

func (p *testLevelWriter) LogLevelWrite(level insolar.LogLevel, b []byte) (int, error) {
	return p.testWriter.Write([]byte(level.String() + string(b)))
}

func TestAdapter_fatal_close_on_no_flush(t *testing.T) {
	tw := testWriter{}
	tw.noFlush = true

	writer := NewAdapter(&tw, false, nil, nil)
	writer.setState(adapterPanicOnFatal)

	var err error

	require.Equal(t, 0, tw.flushCount)
	err = writer.DirectFlushFatal()
	require.NoError(t, err)

	require.Equal(t, 1, tw.flushCount)
	require.Equal(t, 1, tw.closeCount)

	require.PanicsWithValue(t, "fatal lock", func() {
		_ = writer.Flush()
	})
	require.Equal(t, 1, tw.flushCount)
}

func TestAdapter_fatal(t *testing.T) {
	tw := testWriter{}
	writer := NewAdapter(&tw, false, nil, nil)
	writer.setState(adapterPanicOnFatal)

	var err error

	_, err = writer.LogLevelWrite(insolar.WarnLevel, []byte("NORM must pass\n"))
	require.NoError(t, err)
	assert.Contains(t, tw.String(), "NORM must pass")

	require.Equal(t, 0, tw.flushCount)
	err = writer.Flush()
	require.NoError(t, err)
	require.Equal(t, 1, tw.flushCount)

	require.False(t, writer.IsFatal())
	require.True(t, writer.SetFatal())
	require.True(t, writer.IsFatal())
	require.False(t, writer.SetFatal())
	require.True(t, writer.IsFatal())

	require.Panics(t, func() {
		_, _ = writer.LogLevelWrite(insolar.WarnLevel, []byte("must NOT pass\n"))
	})
	require.Panics(t, func() {
		_ = writer.Flush()
	})
	require.Equal(t, 1, tw.flushCount)
	err = writer.DirectFlushFatal()
	require.NoError(t, err)
	require.Equal(t, 2, tw.flushCount)

	require.True(t, writer.IsClosed())
	require.Equal(t, 0, tw.closeCount)
	require.Panics(t, func() {
		_ = writer.Close()
	})
	require.Equal(t, 0, tw.closeCount)
	require.True(t, writer.IsClosed())

	assert.NotContains(t, tw.String(), "must NOT pass")

	_, err = writer.DirectLevelWrite(insolar.WarnLevel, []byte("DIRECT must pass\n"))
	require.NoError(t, err)
	assert.Contains(t, tw.String(), "DIRECT must pass")
}

func TestAdapter_fatal_close(t *testing.T) {
	tw := testWriter{}
	writer := NewAdapter(&tw, false, nil, nil)
	writer.setState(adapterPanicOnFatal)

	require.True(t, writer.SetFatal())

	require.False(t, writer.IsClosed())
	require.Equal(t, 0, tw.closeCount)
	require.Panics(t, func() {
		_ = writer.Close()
	})
	require.Equal(t, 1, tw.closeCount)
	require.True(t, writer.IsClosed())
}

func TestAdapter_fatal_direct_close(t *testing.T) {
	tw := testWriter{}
	writer := NewAdapter(&tw, false, nil, nil)
	writer.setState(adapterPanicOnFatal)

	var err error

	require.True(t, writer.SetFatal())

	require.Equal(t, 0, tw.closeCount)
	err = writer.DirectClose()
	require.NoError(t, err)
	require.Equal(t, 1, tw.closeCount)
	require.True(t, writer.IsClosed())
	err = writer.DirectClose()
	require.Error(t, err) // underlying's error
	require.Equal(t, 2, tw.closeCount)

	require.Panics(t, func() {
		_ = writer.Close()
	})
	require.Equal(t, 2, tw.closeCount)
	require.True(t, writer.IsClosed())
}

func TestAdapter_fatal_flush(t *testing.T) {
	tw := testWriter{}
	writer := NewAdapter(&tw, false, nil, nil)
	writer.setState(adapterPanicOnFatal)

	var err error

	_, err = writer.LogLevelWrite(insolar.WarnLevel, []byte("must pass\n"))
	require.NoError(t, err)
	assert.Contains(t, tw.String(), "must pass")

	require.False(t, writer.IsFatal())
	require.Equal(t, 0, tw.flushCount)
	err = writer.DirectFlushFatal()
	require.NoError(t, err)
	require.Equal(t, 1, tw.flushCount)
	require.True(t, writer.IsFatal())
}

func TestAdapter_fatal_flush_helper(t *testing.T) {
	tw := testWriter{}
	flushHelperCount := 0
	writer := NewAdapter(&tw, false, nil, func() error {
		flushHelperCount++
		return nil
	})
	writer.setState(adapterPanicOnFatal)

	var err error

	_, err = writer.LogLevelWrite(insolar.WarnLevel, []byte("must pass\n"))
	require.NoError(t, err)
	assert.Contains(t, tw.String(), "must pass")

	require.False(t, writer.IsFatal())
	require.Equal(t, 0, tw.flushCount)
	require.Equal(t, 0, flushHelperCount)
	err = writer.Flush()
	require.NoError(t, err)
	require.Equal(t, 1, tw.flushCount)
	require.Equal(t, 0, flushHelperCount)
	err = writer.Flush()
	require.NoError(t, err)
	require.Equal(t, 2, tw.flushCount)
	require.Equal(t, 0, flushHelperCount)

	err = writer.DirectFlushFatal()
	require.NoError(t, err)
	require.Equal(t, 3, tw.flushCount)
	require.Equal(t, 1, flushHelperCount)
	require.True(t, writer.IsFatal())

	err = writer.DirectFlushFatal()
	require.Error(t, err)
	require.Equal(t, 3, tw.flushCount)
	require.Equal(t, 1, flushHelperCount)
}

func TestAdapter_close(t *testing.T) {
	tw := testWriter{}
	writer := NewAdapter(&tw, false, nil, nil)

	var err error

	_, err = writer.LogLevelWrite(insolar.WarnLevel, []byte("must pass\n"))
	require.NoError(t, err)
	assert.Contains(t, tw.String(), "must pass")

	require.False(t, writer.IsClosed())
	require.Equal(t, 0, tw.closeCount)
	err = writer.Close()
	require.NoError(t, err)
	require.Equal(t, 1, tw.closeCount)
	require.True(t, writer.IsClosed())
	require.False(t, writer.SetClosed())
	require.True(t, writer.IsClosed())
	err = writer.Close()
	require.Error(t, err)

	_, err = writer.LogLevelWrite(insolar.WarnLevel, []byte("must NOT pass\n"))
	require.Error(t, err)

	assert.NotContains(t, tw.String(), "must NOT pass")

	err = writer.DirectClose()
	require.Error(t, err)
	require.Equal(t, 2, tw.closeCount)

	_, err = writer.DirectLevelWrite(insolar.WarnLevel, []byte("DIRECT must pass\n"))
	require.NoError(t, err)
	assert.Contains(t, tw.String(), "DIRECT must pass")
}

func TestAdapter_close_protect(t *testing.T) {
	tw := testWriter{}
	writer := NewAdapter(&tw, true, nil, nil)

	var err error

	_, err = writer.LogLevelWrite(insolar.WarnLevel, []byte("must pass\n"))
	require.NoError(t, err)
	assert.Contains(t, tw.String(), "must pass")

	require.False(t, writer.IsClosed())
	require.Equal(t, 0, tw.closeCount)
	err = writer.Close()
	require.NoError(t, err)
	require.Equal(t, 0, tw.closeCount)
	require.True(t, writer.IsClosed())
	require.False(t, writer.SetClosed())
	require.True(t, writer.IsClosed())
	err = writer.Close()
	require.Error(t, err)
	require.Equal(t, 0, tw.closeCount)

	_, err = writer.LogLevelWrite(insolar.WarnLevel, []byte("must NOT pass\n"))
	require.Error(t, err)

	assert.NotContains(t, tw.String(), "must NOT pass")
}

func TestAdapter_direct_close_protect(t *testing.T) {
	tw := testWriter{}
	writer := NewAdapter(&tw, true, nil, nil)

	var err error

	_, err = writer.LogLevelWrite(insolar.WarnLevel, []byte("must pass\n"))
	require.NoError(t, err)
	assert.Contains(t, tw.String(), "must pass")

	require.False(t, writer.IsClosed())
	require.Equal(t, 0, tw.closeCount)
	err = writer.DirectClose()
	require.NoError(t, err)
	require.Equal(t, 0, tw.closeCount)
	require.True(t, writer.IsClosed())
	require.False(t, writer.SetClosed())
	require.True(t, writer.IsClosed())
	err = writer.DirectClose()
	require.NoError(t, err)
	require.Equal(t, 0, tw.closeCount)

	_, err = writer.LogLevelWrite(insolar.WarnLevel, []byte("must NOT pass\n"))
	require.Error(t, err)

	assert.NotContains(t, tw.String(), "must NOT pass")
}

func TestAdapter_level_write(t *testing.T) {
	tw := testLevelWriter{}
	writer := NewAdapter(&tw, false, nil, nil)

	_, err := writer.LogLevelWrite(insolar.WarnLevel, []byte(" must pass\n"))
	require.NoError(t, err)
	assert.Contains(t, tw.String(), "warn must pass")
}
