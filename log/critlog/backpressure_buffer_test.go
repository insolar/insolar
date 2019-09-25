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
	"bytes"
	"context"
	"fmt"
	"github.com/insolar/insolar/insolar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBackpressureBuffer_stop(t *testing.T) {
	// BufferCloseOnStop
	var internal *internalBackpressureBuffer
	init := func() {
		bb := NewBackpressureBuffer(&bytes.Buffer{}, 10, 0, 0, 0, nil)
		internal = bb.internalBackpressureBuffer

		bb.StartWorker(context.Background())
	}

	init()
	init = nil // make sure that *BackpressureBuffer is lost

	time.Sleep(time.Millisecond) // and make sure that the worker has started
	require.Equal(t, uint32(1), atomic.LoadUint32(&internal.pendingWrites))

	runtime.GC()                 // the init func() and *BackpressureBuffer are lost, finalizer is created
	time.Sleep(time.Millisecond) // make sure the finalizer was executed
	runtime.GC()                 // finalizer is released

	time.Sleep(time.Millisecond) // not needed, just to be safe
	runtime.GC()

	require.Equal(t, uint32(0), atomic.LoadUint32(&internal.pendingWrites))

	be := <-internal.buffer
	// depletionMark mark in the buffer indicates that any worker will stop
	require.Equal(t, depletionMark, be.flushMark)
	require.Equal(t, 0, len(internal.buffer))
}

type chanWriter struct {
	out      chan<- []byte
	total    uint32
	parallel uint32
}

func (c *chanWriter) Write(p []byte) (int, error) {
	atomic.AddUint32(&c.total, 1)
	atomic.AddUint32(&c.parallel, 1)
	//fmt.Println("before out <- ", string(p), "\n", string(debug.Stack()))
	c.out <- p
	//fmt.Println("after out <- ", string(p))
	atomic.AddUint32(&c.parallel, ^uint32(0))
	return len(p), nil
}

func (c *chanWriter) Close() (err error) {
	close(c.out)
	return nil
}

func TestBackpressureBuffer_parallel_write_limits(t *testing.T) {
	for parWriters := 1; parWriters < 10; parWriters++ {
		for useWorker := 0; useWorker <= 1; useWorker++ {
			t.Run(fmt.Sprintf("parWriters=%v useWorker=%v", parWriters, useWorker != 0), func(t *testing.T) {
				testBackpressureBufferLimit(t, parWriters, 10, useWorker != 0)
			})
		}
	}
}

func testBackpressureBufferLimit(t *testing.T, parWriters, bufSize int, startWorker bool) {

	ch := make(chan []byte)
	cw := &chanWriter{out: ch}
	bb := NewBackpressureBuffer(cw, bufSize, 0, uint8(parWriters), 0,
		func(missed int) (level insolar.LogLevel, i []byte) {
			return insolar.WarnLevel, []byte(fmt.Sprintf("missed %d", missed))
		})

	producersCount := bufSize + parWriters*2 + 1

	wgStarted := sync.WaitGroup{}
	wgFinished := sync.WaitGroup{}
	wgStarted.Add(producersCount)
	wgFinished.Add(producersCount)

	var producersDone uint32
	for i := 0; i < producersCount; i++ {
		msg := fmt.Sprintf("test msg %d\n", i)
		go func() {
			wgStarted.Done()
			n, err := bb.Write([]byte(msg))

			if n != len(msg) || err != nil {
				panic("write was wrong")
			}

			atomic.AddUint32(&producersDone, 1)
			wgFinished.Done()
		}()
	}

	wgStarted.Wait()

	for i := 0; i <= 9; i++ {
		if parWriters == int(atomic.LoadUint32(&cw.parallel)) && len(bb.buffer) == bufSize &&
			parWriters == int(atomic.LoadUint32(&bb.pendingWrites)) {
			break
		}
		time.Sleep(time.Duration(i+1) * 5 * time.Millisecond)
	}

	require.Equal(t, parWriters, int(atomic.LoadUint32(&bb.pendingWrites)), "not all write slots are occupied")
	require.Equal(t, parWriters, int(atomic.LoadUint32(&cw.total)), "io.Writer is hit by exactly the number of write slots")
	require.Equal(t, parWriters, int(atomic.LoadUint32(&cw.parallel)), "io.Writer is hit by exactly the number of write slots")
	require.Equal(t, 0, int(atomic.LoadUint32(&bb.missCount)), "no misses")
	require.Equal(t, bufSize, len(bb.buffer), "buffer is full")

	require.LessOrEqual(t, bufSize, int(atomic.LoadUint32(&producersDone)), "only producers that fit the buffer")
	// there could be up-to 2*parWriters difference because each writer can pick something from a queue
	require.GreaterOrEqual(t, bufSize+parWriters*2, int(atomic.LoadUint32(&producersDone)), "producers that hit output")

	producersLLDone := uint32(0)
	for i := 0; i < producersCount; i++ {
		msg := fmt.Sprintf("test ll msg %d\n", i)
		go func() {
			n, err := bb.LowLatencyWrite(insolar.InfoLevel, []byte(msg))

			require.Equal(t, n, len(msg))
			require.NoError(t, err)

			atomic.AddUint32(&producersLLDone, 1)
		}()
	}

	for i := 0; i <= 9; i++ {
		if producersCount == int(atomic.LoadUint32(&producersLLDone)) {
			break
		}
		time.Sleep(time.Duration(i+1) * 5 * time.Millisecond)
	}

	require.Equal(t, producersCount, int(atomic.LoadUint32(&producersLLDone)), "all LL producers are done")
	require.Equal(t, producersCount, int(atomic.LoadUint32(&bb.missCount)), "all LL write were missed")
	require.Equal(t, parWriters, int(atomic.LoadUint32(&bb.pendingWrites)), "all write slots are still occupied after LL")
	require.Equal(t, parWriters, int(atomic.LoadUint32(&cw.total)), "io.Writer is hit by exactly the number of write slots after LL")
	require.Equal(t, parWriters, int(atomic.LoadUint32(&cw.parallel)), "io.Writer is hit by exactly the number of write slots after LL")

	go func() {
		for range ch {
		}
	}()

	if startWorker {
		bb.StartWorker(context.Background())
	}

	// NB! Flush() may NOT be able to clean up whole buffer when there are many pending writers, so we will retry
	require.NoError(t, bb.Flush(), "flush")
	for i := 0; i <= 9; i++ {
		require.NoError(t, bb.Flush(), "flush")
		if producersCount == int(atomic.LoadUint32(&producersDone)) && len(bb.buffer) == 0 &&
			producersCount+1 == int(atomic.LoadUint32(&cw.total)) && int(atomic.LoadUint32(&cw.parallel)) == 0 {
			break
		}
		time.Sleep(time.Duration(i+1) * 5 * time.Millisecond)
	}

	wgFinished.Wait()

	require.Equal(t, producersCount, int(atomic.LoadUint32(&producersDone)), "all writers are done")
	require.Equal(t, 0, len(bb.buffer), "buffer is flushed and no marks left")
	require.Equal(t, 0, int(atomic.LoadUint32(&bb.missCount)), "miss counter was flushed")
	require.Equal(t, producersCount+1, int(atomic.LoadUint32(&cw.total)), "producers + miss message")
	require.Equal(t, 0, int(atomic.LoadUint32(&cw.parallel)), "io.Writer is hit by exactly the number of write slots after LL")

	expectedWrites := 0
	if startWorker {
		expectedWrites++
	}
	require.Equal(t, expectedWrites, int(atomic.LoadUint32(&bb.pendingWrites)), "no writers but worker")

	require.NoError(t, bb.Close(), "close")
	require.Errorf(t, bb.Close(), "closed", "must be closed")

	// make sure that the worker will enough time to find a mark and put it back
	for i := 0; i <= 10; i++ {
		if int(atomic.LoadUint32(&bb.pendingWrites)) == 0 && len(bb.buffer) == 1 {
			break
		}
		time.Sleep(time.Duration(i+1) * 5 * time.Millisecond)
	}

	require.Equal(t, producersCount+1, int(atomic.LoadUint32(&cw.total)), "no more messages")
	require.Equal(t, 0, int(atomic.LoadUint32(&bb.pendingWrites)))
	require.Equal(t, 1, len(bb.buffer), "depletion mark")
	require.Panics(t, func() {
		ch <- nil
	}, "send on closed channel")
}

func TestBackpressureBuffer_mute_on_fatal(t *testing.T) {
	tw := testWriter{}
	writer := NewBackpressureBuffer(&tw, 10, 0, 0, 0, nil)
	// We don't want to lock the writer on fatal in tests.
	writer.fatal.unlockPostFatal = true
	tw.flushSupported = true

	var err error

	_, err = writer.LogLevelWrite(insolar.WarnLevel, []byte("WARN must pass\n"))
	require.NoError(t, err)
	_, err = writer.LogLevelWrite(insolar.ErrorLevel, []byte("ERROR must pass\n"))
	require.NoError(t, err)

	assert.False(t, tw.flushed)
	_, err = writer.LogLevelWrite(insolar.PanicLevel, []byte("PANIC must pass\n"))
	require.NoError(t, err)
	assert.True(t, tw.flushed)

	tw.flushed = false
	_, err = writer.LogLevelWrite(insolar.FatalLevel, []byte("FATAL must pass\n"))
	require.NoError(t, err)
	assert.True(t, tw.flushed)
	assert.False(t, tw.closed)

	_, err = writer.LogLevelWrite(insolar.WarnLevel, []byte("WARN must NOT pass\n"))
	require.NoError(t, err)
	_, err = writer.LogLevelWrite(insolar.ErrorLevel, []byte("ERROR must NOT pass\n"))
	require.NoError(t, err)
	_, err = writer.LogLevelWrite(insolar.PanicLevel, []byte("PANIC must NOT pass\n"))
	require.NoError(t, err)

	testLog := tw.String()
	assert.Contains(t, testLog, "WARN must pass")
	assert.Contains(t, testLog, "ERROR must pass")
	assert.Contains(t, testLog, "FATAL must pass")
	assert.NotContains(t, testLog, "must NOT pass")
}

func TestBackpressureBuffer_close_on_no_flush(t *testing.T) {
	tw := testWriter{}
	writer := NewBackpressureBuffer(&tw, 10, 0, 0, 0, nil)
	// We don't want to lock the writer on fatal in tests.
	writer.fatal.unlockPostFatal = true
	tw.flushSupported = false

	var err error

	_, err = writer.LogLevelWrite(insolar.WarnLevel, []byte("WARN must pass\n"))
	require.NoError(t, err)
	_, err = writer.LogLevelWrite(insolar.ErrorLevel, []byte("ERROR must pass\n"))
	require.NoError(t, err)

	_, err = writer.LogLevelWrite(insolar.PanicLevel, []byte("PANIC must pass\n"))
	require.NoError(t, err)
	assert.False(t, tw.flushed)

	_, err = writer.LogLevelWrite(insolar.FatalLevel, []byte("FATAL must pass\n"))
	require.NoError(t, err)
	assert.False(t, tw.flushed)
	assert.True(t, tw.closed)

	_, err = writer.LogLevelWrite(insolar.WarnLevel, []byte("WARN must NOT pass\n"))
	require.NoError(t, err)
	_, err = writer.LogLevelWrite(insolar.ErrorLevel, []byte("ERROR must NOT pass\n"))
	require.NoError(t, err)
	_, err = writer.LogLevelWrite(insolar.PanicLevel, []byte("PANIC must NOT pass\n"))
	require.NoError(t, err)

	testLog := tw.String()
	assert.Contains(t, testLog, "WARN must pass")
	assert.Contains(t, testLog, "ERROR must pass")
	assert.Contains(t, testLog, "FATAL must pass")
	assert.NotContains(t, testLog, "must NOT pass")
}
