// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package log

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
)

func capture(f func()) string {
	defer SaveGlobalLogger()()

	var buf bytes.Buffer

	gl, err := GlobalLogger().Copy().WithOutput(&buf).Build()
	if err != nil {
		panic(err)
	}
	SetGlobalLogger(gl)

	f()

	return buf.String()
}

func assertHelloWorld(t *testing.T, out string) {
	assert.Contains(t, out, "HelloWorld")
}

func TestLog_GlobalLogger_redirection(t *testing.T) {
	defer SaveGlobalLogger()()

	SetLogLevel(insolar.InfoLevel)

	originalG := GlobalLogger()

	var buf bytes.Buffer
	newGL, err := GlobalLogger().Copy().WithOutput(&buf).WithBuffer(10, false).Build()
	require.NoError(t, err)

	SetGlobalLogger(newGL)
	newCopyLL, err := GlobalLogger().Copy().BuildLowLatency()
	require.NoError(t, err)

	originalG.Info("viaOriginalGlobal")
	newGL.Info("viaNewInstance")
	GlobalLogger().Info("viaNewGlobal")
	newCopyLL.Info("viaNewLLCopyOfGlobal")

	s := buf.String()
	require.Contains(t, s, "viaOriginalGlobal")
	require.Contains(t, s, "viaNewInstance")
	require.Contains(t, s, "viaNewGlobal")
	require.Contains(t, s, "viaNewLLCopyOfGlobal")
}

func TestLog_GlobalLogger(t *testing.T) {

	assert.NoError(t, SetLevel("debug"))

	assertHelloWorld(t, capture(func() { Debug("HelloWorld") }))
	assertHelloWorld(t, capture(func() { Debugf("%s", "HelloWorld") }))

	assertHelloWorld(t, capture(func() { Info("HelloWorld") }))
	assertHelloWorld(t, capture(func() { Infof("%s", "HelloWorld") }))

	assertHelloWorld(t, capture(func() { Warn("HelloWorld") }))
	assertHelloWorld(t, capture(func() { Warnf("%s", "HelloWorld") }))

	assertHelloWorld(t, capture(func() { Error("HelloWorld") }))
	assertHelloWorld(t, capture(func() { Errorf("%s", "HelloWorld") }))

	assert.Panics(t, func() { capture(func() { Panic("HelloWorld") }) })
	assert.Panics(t, func() { capture(func() { Panicf("%s", "HelloWorld") }) })

	// cyclic run of this test changes loglevel, so revert it back
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	// can't catch os.exit() to test Fatal
	// Fatal("HelloWorld")
	// Fatalln("HelloWorld")
	// Fatalf("%s", "HelloWorld")
}

func TestLog_NewLog_Config(t *testing.T) {
	invalidtests := map[string]configuration.Log{
		"InvalidAdapter":   configuration.Log{Level: "Debug", Adapter: "invalid", Formatter: "text"},
		"InvalidLevel":     configuration.Log{Level: "Invalid", Adapter: "zerolog", Formatter: "text"},
		"InvalidFormatter": configuration.Log{Level: "Debug", Adapter: "zerolog", Formatter: "invalid"},
	}

	for name, test := range invalidtests {
		t.Run(name, func(t *testing.T) {
			logger, err := NewLog(test)
			assert.Nil(t, logger)
			assert.Error(t, err)
		})
	}

	validtests := map[string]configuration.Log{
		"WithAdapter": configuration.Log{Level: "Debug", Adapter: "zerolog", Formatter: "text"},
	}
	for name, test := range validtests {
		t.Run(name, func(t *testing.T) {
			logger, err := NewLog(test)
			assert.NotNil(t, logger)
			assert.NoError(t, err)
		})
	}
}

func TestLog_GlobalLogger_Level(t *testing.T) {
	defer SaveGlobalLogger()()

	assert.NoError(t, SetLevel("error"))
	assert.Error(t, SetLevel("errorrr"))
}

func TestLog_GlobalLogger_FilterLevel(t *testing.T) {
	defer SaveGlobalLoggerAndFilter(true)()

	assert.NoError(t, SetLevel("debug"))
	assert.NoError(t, SetGlobalLevelFilter(insolar.DebugLevel))
	assertHelloWorld(t, capture(func() { Debug("HelloWorld") }))
	assert.NoError(t, SetGlobalLevelFilter(insolar.InfoLevel))
	assert.Equal(t, "", capture(func() { Debug("HelloWorld") }))
}

func TestLog_GlobalLogger_Save(t *testing.T) {
	assert.NotNil(t, GlobalLogger()) // ensure initialization

	restoreFn := SaveGlobalLoggerAndFilter(true)
	level := GlobalLogger().Copy().GetLogLevel()
	filter := GetGlobalLevelFilter()

	if level != insolar.PanicLevel {
		SetLogLevel(level + 1)
	} else {
		SetLogLevel(insolar.DebugLevel)
	}
	assert.NotEqual(t, level, GlobalLogger().Copy().GetLogLevel())

	if filter != insolar.PanicLevel {
		assert.NoError(t, SetGlobalLevelFilter(filter+1))
	} else {
		assert.NoError(t, SetGlobalLevelFilter(insolar.DebugLevel))
	}
	assert.NotEqual(t, filter, GetGlobalLevelFilter())

	restoreFn()

	assert.Equal(t, level, GlobalLogger().Copy().GetLogLevel())
	assert.Equal(t, filter, GetGlobalLevelFilter())
}

func TestLog_AddFields(t *testing.T) {
	errtxt1 := "~CHECK_ERROR_OUTPUT_WITH_FIELDS~"
	errtxt2 := "~CHECK_ERROR_OUTPUT_WITHOUT_FIELDS~"

	var (
		fieldname  = "TraceID"
		fieldvalue = "Trace100500"
	)
	tt := []struct {
		name    string
		fieldfn func(la insolar.Logger) insolar.Logger
	}{
		{
			name: "WithFields",
			fieldfn: func(la insolar.Logger) insolar.Logger {
				fields := map[string]interface{}{fieldname: fieldvalue}
				return la.WithFields(fields)
			},
		},
		{
			name: "WithField",
			fieldfn: func(la insolar.Logger) insolar.Logger {
				return la.WithField(fieldname, fieldvalue)
			},
		},
	}

	for _, tItem := range tt {
		t.Run(tItem.name, func(t *testing.T) {
			la, err := NewLog(configuration.NewLog())
			assert.NoError(t, err)

			var b bytes.Buffer
			logger, err := la.Copy().WithOutput(&b).Build()
			assert.NoError(t, err)

			tItem.fieldfn(logger).Error(errtxt1)
			logger.Error(errtxt2)

			var recitems []string
			for {
				line, err := b.ReadBytes('\n')
				if err != nil && err != io.EOF {
					require.NoError(t, err)
				}

				recitems = append(recitems, string(line))
				if err == io.EOF {
					break
				}
			}
			assert.Contains(t, recitems[0], errtxt1)
			assert.Contains(t, recitems[1], errtxt2)
			assert.Contains(t, recitems[0], fieldvalue)
			assert.NotContains(t, recitems[1], fieldvalue)
		})
	}
}

var adapters = []string{"zerolog"}

func TestLog_Timestamp(t *testing.T) {
	for _, adapter := range adapters {
		adapter := adapter
		t.Run(adapter, func(t *testing.T) {
			logger, err := NewLog(configuration.Log{Level: "info", Adapter: adapter, Formatter: "json"})
			require.NoError(t, err)
			require.NotNil(t, logger)

			var buf bytes.Buffer
			logger, err = logger.Copy().WithOutput(&buf).Build()
			require.NoError(t, err)

			logger.Error("test")

			assert.Regexp(t, regexp.MustCompile("[0-9][0-9]:[0-9][0-9]:[0-9][0-9]"), buf.String())
		})
	}
}

func TestLog_WriteDuration(t *testing.T) {
	for _, adapter := range adapters {
		adapter := adapter
		t.Run(adapter, func(t *testing.T) {
			logger, err := NewLog(configuration.Log{Level: "info", Adapter: adapter, Formatter: "json"})
			require.NoError(t, err)
			require.NotNil(t, logger)

			var buf bytes.Buffer
			logger, err = logger.Copy().WithOutput(&buf).WithMetrics(insolar.LogMetricsResetMode).Build()
			require.NoError(t, err)

			logger2, err := logger.Copy().WithMetrics(insolar.LogMetricsWriteDelayField).Build()
			require.NoError(t, err)

			logger3, err := logger.Copy().WithMetrics(insolar.LogMetricsResetMode).Build()
			require.NoError(t, err)

			logger.Error("test")
			assert.NotContains(t, buf.String(), `,"writeDuration":"`)
			logger3.Error("test")
			assert.NotContains(t, buf.String(), `,"writeDuration":"`)
			logger2.Error("test2")
			s := buf.String()
			assert.Contains(t, s, `,"writeDuration":"`)
		})
	}
}

func TestLog_DynField(t *testing.T) {
	for _, adapter := range adapters {
		adapter := adapter
		t.Run(adapter, func(t *testing.T) {
			logger, err := NewLog(configuration.Log{Level: "info", Adapter: adapter, Formatter: "json"})
			require.NoError(t, err)
			require.NotNil(t, logger)

			const skipConstant = "---skip---"
			dynFieldValue := skipConstant
			var buf bytes.Buffer
			logger, err = logger.Copy().WithOutput(&buf).WithDynamicField("dynField1", func() interface{} {
				if dynFieldValue == skipConstant {
					return nil
				}
				return dynFieldValue
			}).Build()
			require.NoError(t, err)

			logger.Error("test1")
			assert.NotContains(t, buf.String(), `"dynField1":`)
			dynFieldValue = ""
			logger.Error("test2")
			assert.Contains(t, buf.String(), `"dynField1":""`)
			dynFieldValue = "some text"
			logger.Error("test3")
			assert.Contains(t, buf.String(), `"dynField1":"some text"`)
		})
	}
}

func TestMain(m *testing.M) {
	l, err := GlobalLogger().Copy().WithFormat(insolar.JSONFormat).WithLevel(insolar.DebugLevel).Build()
	if err != nil {
		panic(err)
	}
	SetGlobalLogger(l)
	_ = SetGlobalLevelFilter(insolar.DebugLevel)
	exitCode := m.Run()
	os.Exit(exitCode)
}
