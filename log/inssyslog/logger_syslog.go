// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
