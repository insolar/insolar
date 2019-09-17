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
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"runtime"
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

func NewCriticalWriter(w io.Writer) zerolog.LevelWriter {
	return &criticalWriter{w: AsLevelWriter(w), spinWaitCount: 10}
}

func NewCriticalWriterExt(w io.Writer, skippedFormatter SkippedFormatterFunc, spinWaitCount int) zerolog.LevelWriter {
	return &criticalWriter{w: AsLevelWriter(w), skippedFormatter: skippedFormatter, spinWaitCount: spinWaitCount}
}

type SkippedFormatterFunc func(missed uint32) []byte

var _ zerolog.LevelWriter = &criticalWriter{}

type criticalWriter struct {
	w                zerolog.LevelWriter
	skippedFormatter SkippedFormatterFunc
	spinWaitCount    int
	state            uint32 // atomic
}

func (w *criticalWriter) Write(p []byte) (int, error) {
	w.spinWrite(func() {
		_, _ = w.w.Write(p)
	})
	return len(p), nil
}

func (w *criticalWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	w.spinWrite(func() {
		_, _ = w.w.WriteLevel(level, p)
	})
	return len(p), nil
}

func (w *criticalWriter) spinWrite(fn func()) {
	for i := 0; ; i++ {
		c := atomic.LoadUint32(&w.state)

		switch {
		case c == 0:
			if !atomic.CompareAndSwapUint32(&w.state, 0, 1) {
				continue
			}
		case i > 0 && i >= w.spinWaitCount:
			if !atomic.CompareAndSwapUint32(&w.state, c, c+1) {
				continue
			}
			return
		default:
			runtime.Gosched()
			continue
		}

		break
	}

	go func() {
		fn()

		for {
			c := atomic.LoadUint32(&w.state)
			if atomic.CompareAndSwapUint32(&w.state, c, 0) {
				if c > 1 {
					w.reportMissed(c - 1)
				}
				return
			}
		}
	}()
}

func (w *criticalWriter) reportMissed(missed uint32) {
	var skippedMsg []byte
	if w.skippedFormatter != nil {
		skippedMsg = w.skippedFormatter(missed)
		if len(skippedMsg) == 0 {
			return
		}
	} else {
		skippedMsg = ([]byte)(fmt.Sprintf("events were skipped by critical write: n=%d", missed))
	}
	_, _ = w.w.WriteLevel(zerolog.WarnLevel, skippedMsg)
}
