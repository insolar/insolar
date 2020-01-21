// Copyright 2020 Insolar Network Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package critlog

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log/logoutput"
	"github.com/insolar/insolar/network/consensus/common/args"
)

func TestBackpressureBuffer_stop(t *testing.T) {
	t.SkipNow()

	for _, c := range constructors {
		t.Run(c.name, func(t *testing.T) {

			// BufferCloseOnStop
			var internal *internalBackpressureBuffer
			init := func() {
				bb := c.fn(&testWriter{})
				internal = bb.internalBackpressureBuffer
				bb.StartWorker(context.Background())
			}

			init()
			init = nil // make sure that *BackpressureBuffer is lost

			time.Sleep(time.Millisecond) // and make sure that the worker has started
			require.Equal(t, int32(1+1<<16), atomic.LoadInt32(&internal.writerCounts))

			runtime.GC()                 // the init func() and *BackpressureBuffer are lost, finalizer is created
			time.Sleep(time.Millisecond) // make sure the finalizer was executed
			runtime.GC()                 // finalizer is released

			time.Sleep(time.Millisecond) // not needed, just to be safe
			runtime.GC()

			require.Equal(t, int32(0), atomic.LoadInt32(&internal.writerCounts))

			be := <-internal.buffer
			// depletionMark mark in the buffer indicates that any worker will stop
			require.Equal(t, depletionMark, be.flushMark)
			require.Equal(t, 0, len(internal.buffer))
		})
	}
}

func TestBackpressureBuffer_parallel_write_limits_on_buffer(t *testing.T) {
	t.SkipNow()

	for repeat := 1; repeat > 0; repeat-- {
		for parWriters := 1; parWriters <= 20; parWriters++ {
			for useWorker := 0; useWorker <= 1; useWorker++ {
				t.Run(fmt.Sprintf("buffer parWriters=%v useWorker=%v", parWriters, useWorker != 0), func(t *testing.T) {
					testBackpressureBufferLimit(t, parWriters, true, useWorker != 0)
				})
			}
		}
	}
}

func TestBackpressureBuffer_parallel_write_limits_on_bypass(t *testing.T) {
	t.SkipNow()

	for repeat := 1; repeat > 0; repeat-- {
		for parWriters := 1; parWriters <= 20; parWriters++ {
			for useWorker := 0; useWorker <= 1; useWorker++ {
				t.Run(fmt.Sprintf("bypass parWriters=%v useWorker=%v", parWriters, useWorker != 0), func(t *testing.T) {
					testBackpressureBufferLimit(t, parWriters, false, useWorker != 0)
				})
			}
		}
	}
}

func testBackpressureBufferLimit(t *testing.T, parWriters int, hasBuffer bool, startWorker bool) {

	ch := make(chan []byte)
	cw := &chanWriter{out: ch}

	missedFn := func(missed int) (level insolar.LogLevel, i []byte) {
		return insolar.WarnLevel, []byte(fmt.Sprintf("missed %d", missed))
	}

	allocatedBufSize := int(args.Prime(parWriters)) * 2
	var bb *BackpressureBuffer
	if hasBuffer {
		bb = NewBackpressureBuffer(wrapOutput(cw), allocatedBufSize, uint8(parWriters), 0, missedFn)
	} else {
		bb = NewBackpressureBufferWithBypass(wrapOutput(cw), allocatedBufSize, uint8(parWriters), 0, missedFn)
	}

	producersCount := allocatedBufSize + parWriters*2 + 1

	bufSize := 0
	if hasBuffer {
		bufSize = allocatedBufSize
	}

	wgStarted := sync.WaitGroup{}

	wgFinished := Semaphore{}
	wgFinished.Add(bufSize)

	var produceIndex int
	var producersDone uint32

	nextProduce := func(i int) {
		msg := fmt.Sprintf("test msg %d", i)
		wgStarted.Done()
		//go fmt.Println("before: ", msg)
		n, err := bb.Write([]byte(msg))
		//go fmt.Println(" after: ", msg)

		if n != len(msg) || err != nil {
			panic("write was wrong")
		}

		atomic.AddUint32(&producersDone, 1)
		wgFinished.Done()
	}

	cw.wgBefore.Add(parWriters)

	// fill up all writing slots
	// NB! some writers will also capture +1 event, so we need to put more
	for i := bufSize + parWriters*2; i > 0; i-- {
		if produceIndex >= producersCount {
			break
		}
		wgStarted.Add(1)
		go nextProduce(produceIndex)
		produceIndex++
	}
	wgStarted.Wait()

	cw.wgBefore.Wait()
	cw.wgBefore.Add(producersCount - parWriters)

	require.Equal(t, parWriters, int(atomic.LoadUint32(&cw.total)), "io.Writer is hit by exactly the number of write slots")
	require.Equal(t, parWriters, int(atomic.LoadUint32(&cw.parallel)), "io.Writer is hit by exactly the number of write slots")
	require.Equal(t, parWriters, bb.getPendingWrites(), "not all write slots are occupied")
	require.Equal(t, 0, int(atomic.LoadUint32(&bb.missCount)), "no misses")

	entriesPickedByWriters := bufSize - len(bb.buffer)
	// fill up whole buffer
	for i := entriesPickedByWriters; i > 0; i-- {
		if produceIndex >= producersCount {
			break
		}
		wgStarted.Add(1)
		go nextProduce(produceIndex)
		produceIndex++
	}
	wgStarted.Wait()
	wgFinished.Wait()
	wgFinished.Add(producersCount - bufSize)

	require.Equal(t, bufSize, len(bb.buffer), "buffer is full")

	//there could be up-to parWriters difference because each writer can pick something from a queue
	require.LessOrEqual(t, bufSize, int(atomic.LoadUint32(&producersDone)))
	require.GreaterOrEqual(t, bufSize+parWriters*2, int(atomic.LoadUint32(&producersDone)), "producers that hit output")

	require.Equal(t, parWriters, int(atomic.LoadUint32(&cw.total)), "io.Writer is hit by exactly the number of write slots")
	require.Equal(t, parWriters, int(atomic.LoadUint32(&cw.parallel)), "io.Writer is hit by exactly the number of write slots")
	require.Equal(t, parWriters, bb.getPendingWrites(), "not all write slots are occupied")
	require.Equal(t, 0, int(atomic.LoadUint32(&bb.missCount)), "no misses")

	// add up remaining producers
	for i := producersCount - produceIndex; i > 0; i-- {
		wgStarted.Add(1)
		go nextProduce(produceIndex)
		produceIndex++
	}
	wgStarted.Wait()

	wgLLFinished := sync.WaitGroup{}

	expectedMiss := producersCount - cap(bb.buffer) + len(bb.buffer)
	for i := 0; i < producersCount; i++ {
		wgLLFinished.Add(1)
		msg := fmt.Sprintf("test LL msg %d", i)
		go func() {
			//go fmt.Println("before: ", msg)
			n, err := bb.LowLatencyWrite(insolar.InfoLevel, []byte(msg))
			//go fmt.Println(" after: ", msg)

			require.Equal(t, n, len(msg))
			require.NoError(t, err)

			wgLLFinished.Done()
		}()
	}
	wgLLFinished.Wait()

	require.Equal(t, expectedMiss, int(atomic.LoadUint32(&bb.missCount)), "all LL writes were missed")
	require.Equal(t, parWriters, bb.getPendingWrites(), "all write slots are still occupied after LL")
	require.Equal(t, parWriters, int(atomic.LoadUint32(&cw.total)), "io.Writer is hit by exactly the number of write slots after LL")
	require.Equal(t, parWriters, int(atomic.LoadUint32(&cw.parallel)), "io.Writer is hit by exactly the number of write slots after LL")

	cw.wgBefore.Add(producersCount - expectedMiss + 1)
	cw.wgAfter.Add(producersCount*2 - expectedMiss + 1)

	go func() {
		for p := range ch {
			//go fmt.Println("   out: ", string(p))
			runtime.KeepAlive(p)
		}
	}()

	if startWorker {
		bb.StartWorker(context.Background())
	}

	require.NoError(t, bb.Flush(), "flush")
	require.Equal(t, 0, int(atomic.LoadUint32(&bb.missCount)), "miss counter was flushed")

	// NB! Flush() may NOT be able to clean up whole buffer when there are many pending writers, so we will retry
	for i := 0; i <= 9; i++ {
		if len(bb.buffer) == 0 && int(atomic.LoadUint32(&cw.parallel)) == 0 {
			break
		}
		require.NoError(t, bb.Flush(), "repeated flush error")
		time.Sleep(time.Duration(i+1) * 5 * time.Millisecond)
	}

	wgFinished.Wait()
	//cw.wgBefore.Wait()
	//cw.wgAfter.Wait()
	//
	//require.Equal(t, 0, len(bb.buffer), "buffer is flushed and no marks left")
	//require.Equal(t, producersCount*2+1-expectedMiss, int(atomic.LoadUint32(&cw.total)), "producers + miss message")
	//require.Equal(t, 0, int(atomic.LoadUint32(&cw.parallel)), "io.Writer is hit by exactly the number of write slots after LL")

	expectedWrites := 0
	if startWorker {
		expectedWrites++
	}
	require.Equal(t, expectedWrites, bb.getPendingWrites(), "no writers but worker")

	require.NoError(t, bb.Close(), "close")
	require.Equal(t, 1, len(bb.buffer), "depletion mark")

	require.Error(t, bb.Close(), "must be closed")
	require.Equal(t, 1, len(bb.buffer), "depletion mark")

	require.Equal(t, producersCount*2+1-expectedMiss, int(atomic.LoadUint32(&cw.total)), "no more messages")
	require.Equal(t, 0, bb.getPendingWrites())
	require.Equal(t, 1, len(bb.buffer), "depletion mark")
	require.Panics(t, func() {
		ch <- nil
	}, "send on closed channel")
}

func TestBackpressureBuffer_mute_on_fatal(t *testing.T) {

	for _, c := range constructors {
		tw := testWriter{}
		writer := c.fn(&tw)

		t.Run(c.name, func(t *testing.T) {
			// We don't want to lock the writer on fatal in tests.

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
			require.PanicsWithValue(t, "fatal", func() {
				_, _ = writer.LogLevelWrite(insolar.FatalLevel, []byte("FATAL must pass\n"))
			})
			assert.True(t, tw.flushed)
			assert.False(t, tw.closed)

			// MUST hang. Tested by logoutput.Adapter
			//_, _ = writer.LogLevelWrite(insolar.WarnLevel, []byte("WARN must NOT pass\n"))
			//_, _ = writer.LogLevelWrite(insolar.ErrorLevel, []byte("ERROR must NOT pass\n"))
			//_, _ = writer.LogLevelWrite(insolar.PanicLevel, []byte("PANIC must NOT pass\n"))
			//
			testLog := tw.String()
			assert.Contains(t, testLog, "WARN must pass")
			assert.Contains(t, testLog, "ERROR must pass")
			assert.Contains(t, testLog, "FATAL must pass")
			//assert.NotContains(t, testLog, "must NOT pass")
		})
	}
}

var constructors = []struct {
	name string
	fn   func(output io.WriteCloser) *BackpressureBuffer
}{
	{name: "unlimited_bypass", fn: func(output io.WriteCloser) *BackpressureBuffer {
		return NewBackpressureBufferWithBypass(wrapOutput(output), 10, 0, 0, nil)
	}},
	{name: "limited_bypass", fn: func(output io.WriteCloser) *BackpressureBuffer {
		return NewBackpressureBufferWithBypass(wrapOutput(output), 10, 5, 0, nil)
	}},
	{name: "limited_buffer", fn: func(output io.WriteCloser) *BackpressureBuffer {
		return NewBackpressureBuffer(wrapOutput(output), 10, 5, 0, nil)
	}},
}

func wrapOutput(output io.WriteCloser) *logoutput.Adapter {
	return logoutput.NewAdapter(output, false, nil, func() error {
		if tw, ok := output.(*testWriter); ok {
			_ = tw.Flush()
		} else {
			_ = output.Close()
		}
		panic("fatal")
	})
}

type chanWriter struct {
	out      chan<- []byte
	total    uint32
	parallel uint32
	wgBefore sync.WaitGroup
	wgAfter  Semaphore
}

func (c *chanWriter) Write(p []byte) (int, error) {
	atomic.AddUint32(&c.total, 1)
	// maxParallel :=
	atomic.AddUint32(&c.parallel, 1)
	//fmt.Println("before: ", string(p))//, "\n", string(debug.Stack()))
	c.wgBefore.Done()
	c.out <- p
	//fmt.Println(" after: ", string(p))
	atomic.AddUint32(&c.parallel, ^uint32(0))
	c.wgAfter.Done()
	return len(p), nil
}

func (c *chanWriter) Close() (err error) {
	close(c.out)
	return nil
}

type Semaphore struct {
	sync.Mutex
	cond    *sync.Cond
	counter int32
}

func (p *Semaphore) init() {
	if p.cond == nil {
		p.cond = sync.NewCond(&p.Mutex)
	}
}

func (p *Semaphore) Done() {
	p.Add(-1)
}

func (p *Semaphore) Add(increment int) int {
	p.Lock()
	p.init()

	before := p.counter
	p.counter += int32(increment)
	if before > 0 && p.counter <= 0 {
		p.cond.Broadcast()
	}
	v := int(p.counter)
	p.Unlock()

	return v
}

func (p *Semaphore) Wait() {
	p.Lock()
	p.init()

	for {
		if p.counter <= 0 {
			break
		}
		p.cond.Wait()
	}
	p.Unlock()
}
