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
	"errors"
	"github.com/insolar/insolar/insolar"
	"io"
	"sync/atomic"
)

func NewFatalDirectWriter(output io.Writer) *FatalDirectWriter {
	return &FatalDirectWriter{
		output: OutputHelper{output},
	}
}

type Flusher interface {
	Flush() error
}

type Syncer interface {
	Sync() error
}

var _ insolar.LogLevelWriter = &FatalDirectWriter{}
var _ io.WriteCloser = &FatalDirectWriter{}

type FatalDirectWriter struct {
	output          OutputHelper
	state           uint32 // atomic
	unlockPostFatal bool
}

func (p *FatalDirectWriter) Close() error {
	return p.output.DoClose()
}

func (p *FatalDirectWriter) Flush() error {
	_ = p.output.DoFlush()
	return nil
}

func (p *FatalDirectWriter) Write(b []byte) (n int, err error) {
	if p.isFatal() {
		return p.onFatal(insolar.NoLevel, b)
	}
	return p.output.Write(b)
}

func (p *FatalDirectWriter) LogLevelWrite(level insolar.LogLevel, b []byte) (n int, err error) {
	if p.isFatal() {
		return p.onFatal(level, b)
	}

	switch level {
	case insolar.FatalLevel:
		if !p.setFatal() {
			return p.onFatal(level, b)
		}
		n, _ = p.output.DoWriteLevel(level, b)
		return n, p.Close()

	case insolar.PanicLevel:
		n, err = p.output.DoWriteLevel(level, b)
		if err != nil {
			_ = p.Flush()
			return n, err
		}
		return n, p.Flush()
	default:
		return p.output.DoWriteLevel(level, b)
	}
}

func (p *FatalDirectWriter) setFatal() bool {
	return atomic.CompareAndSwapUint32(&p.state, 0, 1)
}

func (p *FatalDirectWriter) isFatal() bool {
	return atomic.LoadUint32(&p.state) != 0
}

func (p *FatalDirectWriter) onFatal(_ insolar.LogLevel, bytes []byte) (int, error) {
	if p.unlockPostFatal {
		return len(bytes), nil
	}
	select {}
}

/* =============================== */

type OutputHelper struct {
	io.Writer
}

func (p *OutputHelper) DoFlush() (err error) {
	if f, ok := p.Writer.(Flusher); ok {
		err = f.Flush()
		if err == nil {
			return nil
		}
	}
	if f, ok := p.Writer.(Syncer); ok {
		err = f.Sync()
		if err == nil {
			return nil
		}
	}
	if err != nil {
		return err
	}
	return errors.New("unsupported: Flush")
}

func (p *OutputHelper) DoClose() error {
	if f, ok := p.Writer.(io.Closer); ok {
		return f.Close()
	}
	return errors.New("unsupported: Close")
}

func (p *OutputHelper) DoWriteLevel(level insolar.LogLevel, b []byte) (n int, err error) {
	if lw, ok := p.Writer.(insolar.LogLevelWriter); ok {
		return lw.LogLevelWrite(level, b)
	}
	return p.Writer.Write(b)
}
