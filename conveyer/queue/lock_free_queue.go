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

package queue

import (
	"sync/atomic"
	"unsafe"

	"github.com/pkg/errors"
)

// TODO: not completely implemented
type LockFreeQueue struct {
	head *queueItem
}

// NewLockFreeQueue creates lf-queue
func NewLockFreeQueue() IQueue {
	queue := &LockFreeQueue{
		head: &emptyQueueItem,
	}
	return queue
}

// SinkPush is mutex-based realization of IQueue
func (q *LockFreeQueue) SinkPush(data interface{}) error {

	newNode := &queueItem{
		payload: data,
	}

	newNodeAdded := false

	for !newNodeAdded {
		head := (*queueItem)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head))))
		if head == nil {
			return errors.New("[ SinkPush ] Queue is blocked")
		}

		if q.HasSignal() {
			// TODO: no specific logic yet. Will be added later
		}

		newNode.next = head
		newNode.index = head.index + 1
		newNode.biggestQueueSignal = maxSignal(q.head.biggestQueueSignal, newNode.itemType)

		newNodeAdded = atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head)), unsafe.Pointer(head), unsafe.Pointer(newNode))
	}

	return nil
}

// SinkPushAll is mutex-based realization of IQueue
func (q *LockFreeQueue) SinkPushAll(data []interface{}) error {
	inputSize := len(data)
	lastElement := &queueItem{}
	newHead := lastElement
	for i := 0; i < inputSize-1; i++ {
		lastElement.payload = data[i]
		lastElement.next = &queueItem{}
		lastElement = lastElement.next
	}
	lastElement.payload = data[inputSize-1]

	newNodeAdded := false

	for !newNodeAdded {
		head := (*queueItem)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head))))
		if head == nil {
			return errors.New("[ SinkPush ] Queue is blocked")
		}

		if q.HasSignal() {
			// do smth interesting
		}

		nextIndex := head.index + 1
		tmpHead := newHead
		for i := inputSize - 1; i >= 0; i-- {
			tmpHead.index = nextIndex + uint(i)
			tmpHead = tmpHead.next
		}

		lastElement.next = head

		newNodeAdded = atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head)), unsafe.Pointer(head), unsafe.Pointer(newHead))
	}

	return nil
}

// // RemoveAll is mutex-based realization of IQueue
func (q *LockFreeQueue) RemoveAll() []OutputElement {
	removed := false
	var head *queueItem
	for !removed {
		head = (*queueItem)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head))))
		if head == nil || head == &emptyQueueItem {
			return nil
		}

		removed = atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head)), unsafe.Pointer(head), unsafe.Pointer(&emptyQueueItem))
	}

	result := make([]OutputElement, 0, head.index)

	current := head
	for i := uint(0); i < head.index; i++ {
		element := OutputElement{
			data:     current.payload,
			itemType: current.itemType,
		}
		result = append(result, element)
		current = current.next
	}

	return result
}

// BlockAndRemoveAll is mutex-based realization of IQueue
func (q *LockFreeQueue) BlockAndRemoveAll() []OutputElement {
	return nil
}

// Unblock is mutex-based realization of IQueue
func (q *LockFreeQueue) Unblock() bool {
	return true
}

// PushSignal is mutex-based realization of IQueue
func (q *LockFreeQueue) PushSignal(signalType uint32, callback SyncDone) error {
	return q.SinkPush(signalType)
}

// HasSignal is mutex-based realization of IQueue
func (q *LockFreeQueue) HasSignal() bool {
	return atomic.LoadUint32(&q.head.biggestQueueSignal) != 0
}
