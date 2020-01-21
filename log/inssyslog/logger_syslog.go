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

package inssyslog

import (
	"io"
	"regexp"

	"github.com/insolar/insolar/insolar"
)

/*
	This code replicates zerolog/syslog.go to implement io.Closer
	Here we can't use zerolog/syslog.go as it is based on syslog package that breaks OS compatibility.
*/

type LogLevelWriteCloser interface {
	insolar.LogLevelWriter
}

// SyslogWriter is an interface matching a syslog.Writer struct.
type SyslogWriteCloser interface {
	io.Closer
	io.Writer
	Debug(m string) error
	Info(m string) error
	Warning(m string) error
	Err(m string) error
	Emerg(m string) error
	Crit(m string) error
}

const DefaultSyslogNetwork = "udp"

var addrRegex = regexp.MustCompile(`^((ip|tcp|udp)(|4|6)|unix|unixgram|unixpacket):`)

func toNetworkAndAddress(s string) (string, string) {
	indexes := addrRegex.FindStringSubmatchIndex(s)
	if len(indexes) == 0 {
		return DefaultSyslogNetwork, s
	}
	return s[:indexes[3]], s[indexes[3]+1:]
}

func ConnectSyslogByParam(outputParam, tag string) (LogLevelWriteCloser, error) {
	if len(outputParam) == 0 || outputParam == "localhost" {
		return ConnectDefaultSyslog(tag)
	}

	nw, addr := toNetworkAndAddress(outputParam)
	return ConnectRemoteSyslog(nw, addr, tag)
}

// SyslogLevelWriter wraps a SyslogWriter and call the right syslog level
// method matching the zerolog level.
func NewSyslogLevelWriter(w SyslogWriteCloser) LogLevelWriteCloser {
	return &syslogWriter{w}
}

type syslogWriter struct {
	w SyslogWriteCloser
}

func (sw *syslogWriter) Flush() error {
	return nil
}

func (sw *syslogWriter) Close() error {
	return sw.w.Close()
}

func (sw *syslogWriter) Write(p []byte) (n int, err error) {
	return sw.w.Write(p)
}

// WriteLevel implements LevelWriter interface.
func (sw *syslogWriter) LogLevelWrite(level insolar.LogLevel, p []byte) (n int, err error) {
	switch level {
	case insolar.DebugLevel:
		err = sw.w.Debug(string(p))
	case insolar.InfoLevel:
		err = sw.w.Info(string(p))
	case insolar.WarnLevel:
		err = sw.w.Warning(string(p))
	case insolar.ErrorLevel:
		err = sw.w.Err(string(p))
	case insolar.FatalLevel:
		err = sw.w.Emerg(string(p))
	case insolar.PanicLevel:
		err = sw.w.Crit(string(p))
	case insolar.NoLevel:
		err = sw.w.Info(string(p))
	default:
		panic("invalid level")
	}
	n = len(p)
	return
}
