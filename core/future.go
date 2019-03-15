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

package core

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var (
	// ErrFutureTimeout is returned when the operation timeout is exceeded.
	ErrFutureTimeout = errors.New("can't wait for result: timeout")
	// ErrFutureChannelClosed is returned when the input channel is closed.
	ErrFutureChannelClosed = errors.New("can't wait for result: channel closed")
)

// Future is ConveyorPendingMessage response future.
type Future interface {

	// ID returns number.
	ID() uint32

	// Result is a channel to listen for future result.
	Result() <-chan Reply

	// SetResult makes packet to appear in result channel.
	SetResult(res Reply)

	// GetResult gets the future result from Result() channel with a timeout set to `duration`.
	GetResult(duration time.Duration) (Reply, error)

	// Cancel closes all channels and cleans up underlying structures.
	Cancel()
}

// CancelCallback is a callback function executed when cancelling Future.
type CancelCallback func(Future)

type future struct {
	result         chan Reply
	id             uint32
	finished       uint32
	cancelCallback CancelCallback
}

// NewFuture creates new Future.
func NewFuture(id uint32, cancelCallback CancelCallback) Future {
	return &future{
		result:         make(chan Reply, 1),
		id:             id,
		cancelCallback: cancelCallback,
	}
}

// ID returns RequestID of packet.
func (future *future) ID() uint32 {
	return future.id
}

// Result returns result packet channel.
func (future *future) Result() <-chan Reply {
	return future.result
}

// SetResult write packet to the result channel.
func (future *future) SetResult(res Reply) {
	if atomic.CompareAndSwapUint32(&future.finished, 0, 1) {
		future.result <- res
		future.finish()
	}
}

// GetResult gets the future result from Result() channel with a timeout set to `duration`.
func (future *future) GetResult(duration time.Duration) (Reply, error) {
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

// Cancel allows to cancel Future processing.
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
	futures map[uint32]Future
}

func newFutureManager() *futureManager {
	return &futureManager{
		futures: make(map[uint32]Future),
	}
}

// Create implements FutureManager interface
func (fm *futureManager) Create() Future {
	id := uint32(1)
	future := NewFuture(id, func(f Future) {
		fm.delete(f.ID())
	})

	fm.mutex.Lock()
	defer fm.mutex.Unlock()

	fm.futures[id] = future

	return future
}

// Get implements FutureManager interface
func (fm *futureManager) Get(id uint32) Future {
	fm.mutex.RLock()
	defer fm.mutex.RUnlock()

	return fm.futures[id]
}

func (fm *futureManager) delete(id uint32) {
	fm.mutex.Lock()
	defer fm.mutex.Unlock()

	delete(fm.futures, id)
}

// FutureManager is store and create Future instances
type FutureManager interface {
	Get(id uint32) Future
	Create() Future
}

func NewFutureManager() FutureManager {
	return newFutureManager()
}
