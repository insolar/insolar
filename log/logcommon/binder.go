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

package logcommon

type EncoderManager interface {
	CreatePartEncoder([]byte) LogObjectWriter
	FlushPartEncoder(LogObjectWriter) []byte
	WriteParts(level LogLevel, parts [][]byte, writer LogLevelWriter) error
}

type TextBinder struct {
	prefix []byte

	parentCtx []byte

	context []byte

	dynCtx []byte

	msgFields []byte

	msgText []byte
}
