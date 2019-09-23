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
	"github.com/stretchr/testify/require"
	"runtime"
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
	// tillDepletion mark in the buffer indicates that any worker will stop
	require.Equal(t, tillDepletion, be.flushMark)
}
