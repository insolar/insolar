package critlog

import (
	"strings"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testWriter struct {
	strings.Builder
	closed, flushed bool
}

func (w *testWriter) Close() error {
	w.closed = true
	return nil
}

func (w *testWriter) Flush() error {
	w.flushed = true
	return nil
}

func TestFatalDirectWriter_mute_on_fatal(t *testing.T) {
	tw := testWriter{}
	writer := NewFatalDirectWriter(&tw)
	// We don't want to lock the writer on fatal in tests.
	writer.fatal.unlockPostFatal = true
	var err error

	_, err = writer.LogLevelWrite(insolar.WarnLevel, []byte("WARN must pass\n"))
	require.NoError(t, err)
	_, err = writer.LogLevelWrite(insolar.ErrorLevel, []byte("ERROR must pass\n"))
	require.NoError(t, err)

	_, err = writer.LogLevelWrite(insolar.PanicLevel, []byte("PANIC must pass\n"))
	require.NoError(t, err)
	assert.True(t, tw.flushed)

	_, err = writer.LogLevelWrite(insolar.FatalLevel, []byte("FATAL must pass\n"))
	require.NoError(t, err)
	assert.True(t, tw.closed)

	_, err = writer.LogLevelWrite(insolar.WarnLevel, []byte("WARN must NOT pass\n"))
	require.NoError(t, err)
	_, err = writer.LogLevelWrite(insolar.ErrorLevel, []byte("ERROR must NOT pass\n"))
	require.NoError(t, err)
	_, err = writer.LogLevelWrite(insolar.PanicLevel, []byte("PANIC must NOT pass\n"))
	require.NoError(t, err)

	testLog := tw.String()
	assert.Contains(t, testLog, "WARN must pass")
	assert.Contains(t, testLog, "ERROR must pass")
	assert.Contains(t, testLog, "FATAL must pass")
	assert.NotContains(t, testLog, "must NOT pass")
}
