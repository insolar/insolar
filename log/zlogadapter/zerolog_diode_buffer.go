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

package zlogadapter

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log/critlog"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"io"
	"sync/atomic"
	"time"
	"unsafe"
)

type SkippedFormatterFunc func(missed int) []byte

func NewDiodeBufferedLevelWriter(output io.Writer, bufSize int, bufPollInterval time.Duration,
	dropBufOnFatal bool, skippedFn SkippedFormatterFunc,
) *DiodeBufferedLevelWriter {
	bw := DiodeBufferedLevelWriter{
		output:          outputGuard{critlog.FlushBypass{Writer: output}},
		bufSize:         bufSize,
		bufPollInterval: bufPollInterval,
		skippedFn:       skippedFn,
		dropBufOnFatal:  dropBufOnFatal,
	}
	bw.buffer = (unsafe.Pointer)(bw.newBuffer())
	return &bw
}

var _ insolar.LogLevelWriter = &DiodeBufferedLevelWriter{}
var _ zerolog.LevelWriter = &DiodeBufferedLevelWriter{}
var _ io.WriteCloser = &DiodeBufferedLevelWriter{}

type DiodeBufferedLevelWriter struct {
	output          outputGuard
	bufSize         int
	bufPollInterval time.Duration
	unlockPostFatal bool
	dropBufOnFatal  bool
	skippedFn       SkippedFormatterFunc

	buffer unsafe.Pointer // *diode.Writer
	state  uint32         // atomic
}

type outputGuard struct {
	critlog.FlushBypass
}

func (p *outputGuard) Flush() error {
	return p.FlushOrSync()
}

func (p *outputGuard) Close() error {
	// fence out the underlying
	return nil
}

func (p *outputGuard) writeLevel(level insolar.LogLevel, b []byte) (n int, err error) {
	if lw, ok := p.Writer.(insolar.LogLevelWriter); ok {
		return lw.LogLevelWrite(level, b)
	}
	if lw, ok := p.Writer.(zerolog.LevelWriter); ok {
		return lw.WriteLevel(ToZerologLevel(level), b)
	}
	return p.Writer.Write(b)
}

/* =================================== */

func newDiodeBuffer(output *outputGuard, bufSize int, bufPollInterval time.Duration, skippedFn SkippedFormatterFunc) *diode.Writer {
	var alertFn func(int)
	if skippedFn != nil {
		alertFn = func(missed int) {
			writeMissedMsg(output, skippedFn, missed)
		}
	}

	nb := diode.NewWriter(output, bufSize, bufPollInterval, alertFn)
	return &nb
}

func writeMissedMsg(output *outputGuard, skippedFn SkippedFormatterFunc, missed int) {
	msg := skippedFn(missed)
	if len(msg) > 0 {
		_, _ = output.writeLevel(insolar.WarnLevel, msg)
	}
}

func (p *DiodeBufferedLevelWriter) newBuffer() *diode.Writer {
	return newDiodeBuffer(&p.output, p.bufSize, p.bufPollInterval, p.skippedFn)
}

func (p *DiodeBufferedLevelWriter) getBuffer() *diode.Writer {
	return (*diode.Writer)(atomic.LoadPointer(&p.buffer))
}

func (p *DiodeBufferedLevelWriter) dropBuffer() *diode.Writer {
	buf := (*diode.Writer)(atomic.SwapPointer(&p.buffer, nil))
	if buf == nil {
		return nil
	}
	_ = buf.Close()
	return buf
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
	if p.dropBuffer() == nil {
		return nil
	}
	return p.output.FlushBypass.Close()
}

func (p *DiodeBufferedLevelWriter) Flush() error {
	buf := p.replaceBuffer()
	if buf == nil {
		return errors.New("closed")
	}
	_ = buf.Close()
	_ = p.output.Flush()
	return nil
}

func (p *DiodeBufferedLevelWriter) Write(b []byte) (n int, err error) {
	if p.isFatal() {
		return p.onFatal(insolar.NoLevel, b)
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

func (p *DiodeBufferedLevelWriter) LogLevelWrite(level insolar.LogLevel, b []byte) (n int, err error) {
	if p.isFatal() {
		return p.onFatal(level, b)
	}

	switch level {
	case insolar.FatalLevel:
		if !p.setFatal() {
			return p.onFatal(level, b)
		}

		if p.dropBufOnFatal {
			if p.dropBuffer() != nil && p.skippedFn != nil {
				writeMissedMsg(&p.output, p.skippedFn, -1)
			}
		} else {
			_ = p.Close()
		}
		// direct write to the underlying
		return p.output.writeLevel(level, b)

	case insolar.PanicLevel:
		n, err = p._write(b)
		if err != nil {
			return
		}
		return n, p.Flush()
	default:
		return p._write(b)
	}
}

func (p *DiodeBufferedLevelWriter) WriteLevel(level zerolog.Level, b []byte) (n int, err error) {
	return p.LogLevelWrite(FromZerologLevel(level), b)
}

func (p *DiodeBufferedLevelWriter) setFatal() bool {
	return atomic.CompareAndSwapUint32(&p.state, 0, 1)
}

func (p *DiodeBufferedLevelWriter) isFatal() bool {
	return atomic.LoadUint32(&p.state) != 0
}

func (p *DiodeBufferedLevelWriter) onFatal(_ insolar.LogLevel, bytes []byte) (int, error) {
	if p.unlockPostFatal {
		return len(bytes), nil
	}
	select {}
}
