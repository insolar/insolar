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

package critlog

import (
	"github.com/insolar/insolar/insolar"
	"io"
	"sync/atomic"
)

func NewFatalDirectWriter(output io.Writer) *FatalDirectWriter {
	return &FatalDirectWriter{
		output: NewFlushBypass(output),
	}
}

var _ insolar.LogLevelWriter = &FatalDirectWriter{}
var _ io.WriteCloser = &FatalDirectWriter{}

type FatalDirectWriter struct {
	output FlushBypass
	fatal  FatalHelper
}

func (p *FatalDirectWriter) Close() error {
	return p.output.Close()
}

func (p *FatalDirectWriter) Flush() error {
	return p.output.FlushOrSync()
}

func (p *FatalDirectWriter) Write(b []byte) (n int, err error) {
	if p.fatal.IsFatal() {
		return p.fatal.PostFatalWrite(insolar.NoLevel, b)
	}
	return p.output.Write(b)
}

func (p *FatalDirectWriter) LowLatencyWrite(level insolar.LogLevel, b []byte) (int, error) {
	return p.LogLevelWrite(level, b)
}

func (p *FatalDirectWriter) IsLowLatencySupported() bool {
	return false
}

func (p *FatalDirectWriter) GetBareOutput() io.Writer {
	return p.output.Writer
}

func (p *FatalDirectWriter) LogLevelWrite(level insolar.LogLevel, b []byte) (n int, err error) {
	if p.fatal.IsFatal() {
		return p.fatal.PostFatalWrite(level, b)
	}

	switch level {
	case insolar.FatalLevel:
		if !p.fatal.SetFatal() {
			return p.fatal.PostFatalWrite(level, b)
		}
		n, _ = p.output.LogLevelWrite(level, b)
		if ok, err := p.output.DoFlushOrSync(); ok && err == nil {
			return n, nil
		}

		return n, p.output.Close()

	case insolar.PanicLevel:
		n, err = p.output.LogLevelWrite(level, b)
		if err == nil {
			_, _ = p.output.DoFlushOrSync()
		}
		return n, err

	default:
		return p.output.LogLevelWrite(level, b)
	}
}

/* =============================== */

type FatalHelper struct {
	state           uint32 // atomic
	unlockPostFatal bool   // for test usage
}

func (p *FatalHelper) SetFatal() bool {
	return atomic.CompareAndSwapUint32(&p.state, 0, 1)
}

func (p *FatalHelper) IsFatal() bool {
	return atomic.LoadUint32(&p.state) != 0
}

func (p *FatalHelper) LockFatal() {
	if p.unlockPostFatal {
		return
	}
	select {} // lock it down forever
}

func (p *FatalHelper) PostFatalWrite(_ insolar.LogLevel, b []byte) (int, error) {
	p.LockFatal()
	return len(b), nil
}
