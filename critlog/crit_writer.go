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
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
)

/*
	Purpose:
	CriticalWriter provides logging for time-critical components. Any writes will return without waiting on an underlying writer.
	Writes are made in a separate goroutine, but only one goroutine is allowed to go through.
	All other writes will be dropped, but will be counted and later logged.

	Usage:
	- CriticalWriter MUST be applied before any write that can impede performance
	- CriticalWriter SHOULD be applied AFTER a buffer to avoid excessive loss of messages

	Examples:
		logger -> a buffered writer -> CriticalWriter -> ... -> output writer
		logger -> CriticalWriter -> ... -> output writer
*/

func NewCriticalWriter(ctx context.Context, w io.Writer, bufLen int) LevelWriteCloser {
	r := &criticalWriter{w: AsLevelWriter(w), spinCount: 10}
	r.start(ctx, bufLen)
	return r
}

func NewCriticalWriterExt(ctx context.Context, w io.Writer, bufLen int, skippedFormatter SkippedFormatterFunc, spinWaitCount int) LevelWriteCloser {
	r := &criticalWriter{w: AsLevelWriter(w), skippedFn: skippedFormatter, spinCount: spinWaitCount}
	r.start(ctx, bufLen)
	return r
}

type TimeCriticalWriter interface {
	IsTimeCriticalWriter() bool
}
type SkippedFormatterFunc func(missed uint32) []byte

var _ zerolog.LevelWriter = &criticalWriter{}
var _ TimeCriticalWriter = &criticalWriter{}

type criticalWriter struct {
	buffer    chan bufferEntry
	w         zerolog.LevelWriter
	skippedFn SkippedFormatterFunc
	spinCount int
	missCount uint32 // atomic
}

type bufferEntry struct {
	p     []byte
	level zerolog.Level
	wg    *sync.WaitGroup
}

func (w *criticalWriter) waitFlush(closeBuf bool) (ok bool) {
	defer func() {
		ok = recover() == nil
	}()

	if closeBuf {
		close(w.buffer)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	w.buffer <- bufferEntry{wg: &wg}
	wg.Wait()

	return true
}

func (w *criticalWriter) Flush() (err error) {
	if w.waitFlush(false) {
		if f, ok := w.w.(Flusher); ok {
			return f.Flush()
		}
		return errors.New("unsupported: Flush")
	}
	return errors.New("closed")
}

func (w *criticalWriter) Close() (err error) {
	if w.waitFlush(true) {
		if f, ok := w.w.(io.Closer); ok {
			return f.Close()
		}
		return errors.New("unsupported: Close")
	}
	return errors.New("closed")
}

func (w *criticalWriter) Sync() (err error) {
	if w.waitFlush(false) {
		if f, ok := w.w.(Syncer); ok {
			return f.Sync()
		}
		return errors.New("unsupported: Sync")
	}
	return errors.New("closed")
}

func (w *criticalWriter) IsTimeCriticalWriter() bool {
	return true
}

const writeWithoutLevel zerolog.Level = math.MaxUint8

func (w *criticalWriter) Write(p []byte) (int, error) {
	return len(p), w.writeLevel(bufferEntry{p: p, level: writeWithoutLevel})
}

func (w *criticalWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	if level == writeWithoutLevel {
		panic("illegal value")
	}
	return len(p), w.writeLevel(bufferEntry{p: p, level: level})
}

func (w *criticalWriter) writeLevel(entry bufferEntry) (err error) {
	defer func() {
		_ = recover()
		//recovered := recover()
		//if recovered != nil {
		//	err = fmt.Errorf("%v", recovered)
		//}
	}()

	for i := 0; ; i++ {
		select {
		case w.buffer <- entry:
			return nil
		default:
		}

		if i > 0 && i >= w.spinCount {
			c := atomic.LoadUint32(&w.missCount)
			if atomic.CompareAndSwapUint32(&w.missCount, c, c+1) {
				return nil
			}
		}
		runtime.Gosched()
	}
}

func (w *criticalWriter) reportMissed(missed uint32) {
	var skippedMsg []byte
	if w.skippedFn != nil {
		skippedMsg = w.skippedFn(missed)
		if len(skippedMsg) == 0 {
			return
		}
	} else {
		skippedMsg = ([]byte)(fmt.Sprintf("critical logger dropped %d messages", missed))
	}
	_, _ = w.w.WriteLevel(zerolog.WarnLevel, skippedMsg)
}

func (w *criticalWriter) start(ctx context.Context, bufLen int) {
	w.buffer = make(chan bufferEntry, bufLen)
	go w.worker(ctx)
}

func (w *criticalWriter) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case entry, ok := <-w.buffer:
			if !ok {
				return
			}
			if entry.wg != nil {
				entry.wg.Done()
				continue
			}
			w.writeEntry(entry)
		}
	}
}

func (w *criticalWriter) writeEntry(entry bufferEntry) {
	if entry.level == writeWithoutLevel {
		_, _ = w.w.Write(entry.p)
	} else {
		_, _ = w.w.WriteLevel(entry.level, entry.p)
	}

	for {
		c := atomic.LoadUint32(&w.missCount)
		if c == 0 {
			return
		}
		if atomic.CompareAndSwapUint32(&w.missCount, c, 0) {
			w.reportMissed(c)
			return
		}
	}
}
