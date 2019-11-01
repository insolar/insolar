package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
)

type x struct {
	s string
}

func (x x) GetLogObjectMarshaller() insolar.LogObjectMarshaller {
	return &mar{x}
}

type mar struct {
	x x
}

func (m mar) MarshalTextLogObject(lw insolar.LogObjectWriter, lmc insolar.LogObjectMetricCollector) string {
	return m.x.s
}

func (m mar) MarshalBinaryLogObject(lw insolar.LogObjectWriter, lmc insolar.LogObjectMetricCollector) string {
	panic(lw)
}

func (m mar) MarshalMutedLogObject(lw insolar.LogObjectMetricCollector) {
	panic(lw)
}

type y struct {
	insolar.LogObjectTemplate
	s string
}

var logstring = "opaopa"

func TestLogFieldsMarshaler(t *testing.T) {
	for _, obj := range []interface{}{
		x{logstring}, &x{logstring},
		mar{x{logstring}}, &mar{x{logstring}},
		logstring, &logstring, func() string { return logstring },
	} {
		buf := bytes.Buffer{}
		lg, _ := log.GlobalLogger().Copy().WithOutput(&buf).Build()

		_, file, line, _ := runtime.Caller(0)
		lg.WithField("testfield", 200.200).Warn(obj)

		c := make(map[string]interface{})
		err := json.Unmarshal(buf.Bytes(), &c)
		require.NoError(t, err, "unmarshal")

		require.Equal(t, "warn", c["level"], "right message")
		require.Equal(t, "opaopa", c["message"], "right message")
		require.Regexp(t, zerolog.CallerMarshalFunc(file, line+1), c["caller"], "right caller line")
		ltime, err := time.Parse(time.RFC3339Nano, c["time"].(string))
		require.NoError(t, err, "parseable time")
		ldur := time.Now().Sub(ltime)
		require.True(t, ldur > 0, "worktime is greater than zero")
		require.True(t, ldur < time.Second, "worktime lesser than second")
		require.Equal(t, 200.200, c["testfield"], "customfield")
		require.NotNil(t, c["writeDuration"], "duration exists")
	}
}

func TestLogLevels(t *testing.T) {
	buf := bytes.Buffer{}
	lg, _ := log.GlobalLogger().Copy().WithOutput(&buf).Build()

	lg.Level(insolar.FatalLevel).Warn(logstring)
	require.Nil(t, buf.Bytes(), "do not log warns at panic level")

	lg.Warn(logstring)
	require.NotNil(t, buf.Bytes(), "previous logger saves it's level")

	if false {
		log.SetGlobalLevelFilter(insolar.PanicLevel)
		lg.Warn(logstring)
		require.Nil(t, buf.String(), "do not log warns at global panic level")
	}
}

func TestLogOther(t *testing.T) {
	buf := bytes.Buffer{}
	lg, _ := log.GlobalLogger().Copy().WithOutput(&buf).Build()
	c := make(map[string]interface{})
	lg.Warn(nil)
	require.NoError(t, json.Unmarshal(buf.Bytes(), &c))
	require.Equal(t, "<nil>", c["message"], "nil")

	buf.Reset()
	lg.Warn(100.1)
	require.NoError(t, json.Unmarshal(buf.Bytes(), &c))
	require.Equal(t, "100.1", c["message"], "nil")

	buf.Reset()
	lg.Warn(y{s: logstring})
	require.NoError(t, json.Unmarshal(buf.Bytes(), &c))
	require.Equal(t, fmt.Sprintf("{{} %s}", logstring), c["message"], "nil")

	buf.Reset()
	lg.Warn(&y{s: logstring})
	require.NoError(t, json.Unmarshal(buf.Bytes(), &c))
	require.Equal(t, fmt.Sprintf("{{} %s}", logstring), c["message"], "nil")

	// ???
	buf.Reset()
	lg.Warn(struct{ s string }{logstring})
	require.NoError(t, json.Unmarshal(buf.Bytes(), &c))
	require.Equal(t, fmt.Sprintf("{{} %s}", logstring), c["message"], "nil")

}
