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

package inssyslog

import (
	critlog2 "github.com/insolar/insolar/log/critlog"
	"github.com/rs/zerolog"
)

/*
	This code replicates zerolog/syslog.go to implement io.Closer
	Here we can't use zerolog/syslog.go as it is based on syslog package that breaks OS compatibility.
*/

// SyslogLevelWriter wraps a SyslogWriter and call the right syslog level
// method matching the zerolog level.
func NewSyslogLevelWriter(w SyslogWriteCloser) critlog2.LevelWriteCloser {
	return &syslogWriter{w}
}

type syslogWriter struct {
	w SyslogWriteCloser
}

func (sw *syslogWriter) Close() error {
	return sw.w.Close()
}

func (sw *syslogWriter) Write(p []byte) (n int, err error) {
	return sw.w.Write(p)
}

// WriteLevel implements LevelWriter interface.
func (sw *syslogWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	switch level {
	case zerolog.DebugLevel:
		err = sw.w.Debug(string(p))
	case zerolog.InfoLevel:
		err = sw.w.Info(string(p))
	case zerolog.WarnLevel:
		err = sw.w.Warning(string(p))
	case zerolog.ErrorLevel:
		err = sw.w.Err(string(p))
	case zerolog.FatalLevel:
		err = sw.w.Emerg(string(p))
	case zerolog.PanicLevel:
		err = sw.w.Crit(string(p))
	case zerolog.NoLevel:
		err = sw.w.Info(string(p))
	default:
		panic("invalid level")
	}
	n = len(p)
	return
}
