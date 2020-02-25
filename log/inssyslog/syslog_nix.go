// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build !windows,!nacl,!plan9

package inssyslog

import (
	"log/syslog"
)

const defaultSyslogPriority = syslog.LOG_LOCAL0 | syslog.LOG_DEBUG

func ConnectDefaultSyslog(tag string) (LogLevelWriteCloser, error) {
	w, err := syslog.New(defaultSyslogPriority, tag)
	if err != nil {
		return nil, err
	}
	return NewSyslogLevelWriter(w), nil
}

func ConnectRemoteSyslog(network, raddr string, tag string) (LogLevelWriteCloser, error) {
	w, err := syslog.Dial(network, raddr, defaultSyslogPriority, tag)
	if err != nil {
		return nil, err
	}
	return NewSyslogLevelWriter(w), nil
}
