//
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
//

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
