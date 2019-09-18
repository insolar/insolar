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

// +build windows

package inssyslog

import (
	"github.com/insolar/insolar/log/critlog"
	"github.com/pkg/errors"
	"io"
)

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

func ConnectDefaultSyslog(tag string) (critlog.LevelWriteCloser, error) {
	return nil, errors.New("not implemented for Windows")
}

func ConnectRemoteSyslog(network, raddr string, tag string) (critlog.LevelWriteCloser, error) {
	return nil, errors.New("not implemented for Windows")
}
