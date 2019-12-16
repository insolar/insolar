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

package logcommon

import "io"

type LoggerOutputGetter interface {
	GetLoggerOutput() LoggerOutput
}

type LoggerOutput interface {
	LogLevelWriter
	LowLatencyWrite(LogLevel, []byte) (int, error)
	IsLowLatencySupported() bool
}

type LogLevelWriter interface {
	io.WriteCloser
	LogLevelWrite(LogLevel, []byte) (int, error)
	Flush() error
}
