/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package messagebus

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/insolar/insolar/core"
)

var (
	// ErrFutureTimeout is returned when the operation timeout is exceeded.
	ErrFutureTimeout = errors.New("can't wait for result: timeout")
	// ErrFutureChannelClosed is returned when the input channel is closed.
	ErrFutureChannelClosed = errors.New("can't wait for result: channel closed")
)

// CancelCallback is a callback function executed when cancelling ConveyorFuture.
type CancelCallback func(core.ConveyorFuture)

type future struct {
	result         chan core.Reply
	id             uint32
	finished       uint32
	cancelCallback CancelCallback
}

// NewFuture creates new ConveyorFuture.
func NewFuture(id uint32, cancelCallback CancelCallback) core.ConveyorFuture {
	return &future{
		result:         make(chan core.Reply, 1),
		id:             id,
		cancelCallback: cancelCallback,
	}
}

// ID returns RequestID of packet.
func (future *future) ID() uint32 {
	return future.id
}

// Result returns result packet channel.
func (future *future) Result() <-chan core.Reply {
	return future.result
}

// SetResult write packet to the result channel.
func (future *future) SetResult(res core.Reply) {
	if atomic.CompareAndSwapUint32(&future.finished, 0, 1) {
		future.result <- res
		future.finish()
	}
}

// GetResult gets the future result from Result() channel with a timeout set to `duration`.
func (future *future) GetResult(duration time.Duration) (core.Reply, error) {
	select {
	case result, ok := <-future.Result():
		if !ok {
			return nil, ErrFutureChannelClosed
		}
		return result, nil
	case <-time.After(duration):
		future.Cancel()
		return nil, ErrFutureTimeout
	}
}

// Cancel allows to cancel ConveyorFuture processing.
func (future *future) Cancel() {
	if atomic.CompareAndSwapUint32(&future.finished, 0, 1) {
		future.finish()
	}
}

func (future *future) finish() {
	close(future.result)
	future.cancelCallback(future)
}

type futureManager struct {
	mutex   sync.RWMutex
	futures map[uint32]core.ConveyorFuture
}

func newFutureManager() *futureManager {
	return &futureManager{
		futures: make(map[uint32]core.ConveyorFuture),
	}
}

// Create implements FutureManager interface
func (fm *futureManager) Create() core.ConveyorFuture {
	id := uint32(1)
	future := NewFuture(id, func(f core.ConveyorFuture) {
		fm.delete(f.ID())
	})

	fm.mutex.Lock()
	defer fm.mutex.Unlock()

	fm.futures[id] = future

	return future
}

// Get implements FutureManager interface
func (fm *futureManager) Get(id uint32) core.ConveyorFuture {
	fm.mutex.RLock()
	defer fm.mutex.RUnlock()

	return fm.futures[id]
}

func (fm *futureManager) delete(id uint32) {
	fm.mutex.Lock()
	defer fm.mutex.Unlock()

	delete(fm.futures, id)
}

// FutureManager is store and create ConveyorFuture instances
type FutureManager interface {
	Get(id uint32) core.ConveyorFuture
	Create() core.ConveyorFuture
}

func NewFutureManager() FutureManager {
	return newFutureManager()
}

// ConveyorPendingMessage is message for conveyor witch can pending for response
type ConveyorPendingMessage struct {
	Msg core.Parcel
	F   core.ConveyorFuture
}
