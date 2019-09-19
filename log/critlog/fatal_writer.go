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

package critlog

import (
	"github.com/rs/zerolog"
	"io"
	"sync/atomic"
)

func NewFatalDirectWriter(output io.Writer) *FatalDirectWriter {
	return &FatalDirectWriter{
		output: outputGuard{output},
	}
}

var _ zerolog.LevelWriter = &FatalDirectWriter{}
var _ io.WriteCloser = &FatalDirectWriter{}

type FatalDirectWriter struct {
	output          outputGuard
	state           uint32 // atomic
	unlockPostFatal bool
}

func (p *FatalDirectWriter) Close() error {
	return p.output.close()
}

func (p *FatalDirectWriter) Flush() error {
	_ = p.output.flush()
	return nil
}

func (p *FatalDirectWriter) Write(b []byte) (n int, err error) {
	if p.isFatal() {
		return p.onFatal(zerolog.NoLevel, b)
	}
	return p.output.Write(b)
}

func (p *FatalDirectWriter) WriteLevel(level zerolog.Level, b []byte) (n int, err error) {
	if p.isFatal() {
		return p.onFatal(level, b)
	}

	switch level {
	case zerolog.FatalLevel:
		if !p.setFatal() {
			return p.onFatal(level, b)
		}
		n, _ = p.output.writeLevel(level, b)
		return n, p.Close()

	case zerolog.PanicLevel:
		n, err = p.output.writeLevel(level, b)
		if err != nil {
			_ = p.Flush()
			return n, err
		}
		return n, p.Flush()
	default:
		return p.output.writeLevel(level, b)
	}
}

func (p *FatalDirectWriter) setFatal() bool {
	return atomic.CompareAndSwapUint32(&p.state, 0, 1)
}

func (p *FatalDirectWriter) isFatal() bool {
	return atomic.LoadUint32(&p.state) != 0
}

func (p *FatalDirectWriter) onFatal(level zerolog.Level, bytes []byte) (int, error) {
	if p.unlockPostFatal {
		return len(bytes), nil
	}
	select {}
}
