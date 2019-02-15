/*
 *    Copyright 2019 Insolar
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
	"fmt"
	"sync/atomic"
	"unsafe"
)

// TODO: not completely implemented
type LockFreeQueue struct {
	head *QueueItem
}

func NewLockFreeQueue() IQueue {
	queue := &LockFreeQueue{
		head: &emptyQueueItem,
	}
	return queue
}

func (q *LockFreeQueue) SinkPush(data interface{}) bool {

	newNode := &QueueItem{
		payload: data,
	}

	newNodeAdded := false

	for !newNodeAdded {
		head := (*QueueItem)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head))))
		if head == nil {
			return false
		}

		if head.hasSignal() {
			// do smth interesting
		}

		newNode.next = head
		newNode.index = head.index + 1
		// lastNew.signal = max(head.signal, lastNew.type) // TODO:

		newNodeAdded = atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head)), unsafe.Pointer(head), unsafe.Pointer(newNode))
	}

	//fmt.Println("Pushing: ", data)

	return true
}

func (q *LockFreeQueue) SinkPushAll(data []interface{}) bool {
	inputSize := len(data)
	lastElement := &QueueItem{}
	newHead := lastElement
	for i := 0; i < inputSize-1; i++ {
		lastElement.payload = data[i]
		lastElement.next = &QueueItem{}
		lastElement = lastElement.next
	}
	lastElement.payload = data[inputSize-1]

	newNodeAdded := false

	for !newNodeAdded {
		head := (*QueueItem)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head))))
		if head == nil {
			return false
		}

		if head.hasSignal() {
			// do smth interesting
		}

		nextIndex := head.index + 1
		tmpHead := newHead
		for i := inputSize - 1; i >= 0; i-- {
			tmpHead.index = nextIndex + uint(i)
			tmpHead = tmpHead.next
		}

		lastElement.next = head
		// lastNew.signal = max(head.signal, lastNew.type) // TODO:

		newNodeAdded = atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head)), unsafe.Pointer(head), unsafe.Pointer(newHead))
	}

	fmt.Println(" ALL Pushing: ", data)

	return true
}

func (q *LockFreeQueue) RemoveAll() []interface{} {
	removed := false
	var head *QueueItem
	for !removed {
		head = (*QueueItem)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head))))
		if head == nil || head == &emptyQueueItem {
			return nil
		}

		removed = atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head)), unsafe.Pointer(head), unsafe.Pointer(&emptyQueueItem))
	}

	result := make([]interface{}, 0, head.index)

	current := head
	for i := uint(0); i < head.index; i++ {
		result = append(result, current.payload)
		current = current.next
	}

	//fmt.Println(" REMOVING ALL Pushing: ", result)

	return result
}

func (q *LockFreeQueue) BlockAndRemoveAll() []interface{} {
	return nil
}

func (q *LockFreeQueue) Unblock() bool {
	return true
}

func (q *LockFreeQueue) PushSignal(signalType uint, callback SyncDone) bool {
	return q.SinkPush(signalType)
}

func (q *LockFreeQueue) HasSignal() bool {
	return false
}
