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
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"io"
	"sync/atomic"
	"time"
	"unsafe"
)

func NewDiodeBufferedLevelWriter(output io.Writer, bufSize int, bufPollInterval time.Duration,
	dropBufOnFatal bool, skippedFn SkippedFormatterFunc,
) DiodeBufferedLevelWriter {
	return DiodeBufferedLevelWriter{
		output:          outputGuard{output},
		bufSize:         bufSize,
		bufPollInterval: bufPollInterval,
		skippedFn:       skippedFn,
		dropBufOnFatal:  dropBufOnFatal,
	}
}

var _ zerolog.LevelWriter = &DiodeBufferedLevelWriter{}
var _ io.WriteCloser = &DiodeBufferedLevelWriter{}

type DiodeBufferedLevelWriter struct {
	output          outputGuard
	bufSize         int
	bufPollInterval time.Duration
	lockPostFatal   bool
	dropBufOnFatal  bool
	skippedFn       SkippedFormatterFunc

	buffer unsafe.Pointer // *diode.Writer
	state  uint32         // atomic
}

type outputGuard struct {
	io.Writer
}

func (p *outputGuard) Close() error {
	// fence out the underlying
	return nil
}

func (p *outputGuard) flush() (err error) {
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

func (p *outputGuard) close() error {
	if f, ok := p.Writer.(io.Closer); ok {
		return f.Close()
	}
	return errors.New("unsupported: Close")
}

func (p *outputGuard) writeLevel(level zerolog.Level, b []byte) (n int, err error) {
	if lw, ok := p.Writer.(zerolog.LevelWriter); ok {
		return lw.WriteLevel(level, b)
	}
	return p.Writer.Write(b)
}

/* =================================== */

func (p *DiodeBufferedLevelWriter) newBuffer() *diode.Writer {

	var alertFn func(int)
	if p.skippedFn != nil {
		alertFn = func(missed int) {
			msg := p.skippedFn(uint32(missed))
			if len(msg) > 0 {
				_, _ = p.output.writeLevel(zerolog.WarnLevel, msg)
			}
		}
	}

	nb := diode.NewWriter(&p.output, p.bufSize, p.bufPollInterval, alertFn)
	return &nb
}

func (p *DiodeBufferedLevelWriter) getBuffer() *diode.Writer {
	return (*diode.Writer)(atomic.LoadPointer(&p.buffer))
}

func (p *DiodeBufferedLevelWriter) dropBuffer() *diode.Writer {
	return (*diode.Writer)(atomic.SwapPointer(&p.buffer, nil))
}

func (p *DiodeBufferedLevelWriter) replaceBuffer() *diode.Writer {
	var newBuffer *diode.Writer
	for {
		prev := atomic.LoadPointer(&p.buffer)
		if prev == nil {
			return nil
		}

		if newBuffer == nil {
			newBuffer = p.newBuffer()
		}

		if atomic.CompareAndSwapPointer(&p.buffer, prev, (unsafe.Pointer)(newBuffer)) {
			return (*diode.Writer)(prev)
		}
	}
}

func (p *DiodeBufferedLevelWriter) Close() error {
	buf := p.dropBuffer()
	if buf == nil {
		return nil
	}
	_ = buf.Close()
	return p.output.close()
}

func (p *DiodeBufferedLevelWriter) Flush() error {
	buf := p.replaceBuffer()
	if buf == nil {
		return errors.New("closed")
	}
	_ = buf.Close()
	_ = p.output.flush()
	return nil
}

func (p *DiodeBufferedLevelWriter) Write(b []byte) (n int, err error) {
	if p.isFatal() {
		return p.onFatal(zerolog.NoLevel, b)
	}
	return p._write(b)
}

func (p *DiodeBufferedLevelWriter) _write(b []byte) (n int, err error) {
	buf := p.getBuffer()
	if buf == nil {
		return 0, errors.New("closed")
	}
	return buf.Write(b)
}

func (p *DiodeBufferedLevelWriter) WriteLevel(level zerolog.Level, b []byte) (n int, err error) {
	if p.isFatal() {
		return p.onFatal(level, b)
	}

	switch level {
	case zerolog.FatalLevel:
		if !p.setFatal() {
			return p.onFatal(level, b)
		}

		if p.dropBufOnFatal {
			p.dropBuffer()
		} else {
			err = p.Close()
		}
		// direct write to the underlying
		return p.output.writeLevel(level, b)

	case zerolog.PanicLevel:
		n, err = p._write(b)
		if err != nil {
			return
		}
		return n, p.Flush()
	default:
		return p._write(b)
	}
}

func (p *DiodeBufferedLevelWriter) setFatal() bool {
	return atomic.CompareAndSwapUint32(&p.state, 0, 1)
}

func (p *DiodeBufferedLevelWriter) isFatal() bool {
	return atomic.LoadUint32(&p.state) != 0
}

func (p *DiodeBufferedLevelWriter) onFatal(level zerolog.Level, bytes []byte) (int, error) {
	if p.lockPostFatal {
		select {}
	}
	return len(bytes), nil
}
