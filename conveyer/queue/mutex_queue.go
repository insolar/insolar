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
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/pkg/errors"
)

// MutexQueue is mutex-based realization of IQueue
type MutexQueue struct {
	locker sync.Mutex
	head   *QueueItem
}

// NewMutexQueue creates new instance of MutexQueue
func NewMutexQueue() IQueue {
	queue := &MutexQueue{
		head: &emptyQueueItem,
	}
	return queue
}

func (q *MutexQueue) sinkPush(newNode *QueueItem) error {
	q.locker.Lock()
	defer q.locker.Unlock()

	if q.isQueueBlockedUnsafe() {
		return errors.New("[ sinkPush ] Queue is blocked")
	}

	if q.HasSignal() {
		// do smth interesting
	}

	newNode.next = q.head
	newNode.index = q.head.index + 1

	// just max =(
	if q.head.biggestQueueSignal > newNode.itemType {
		newNode.biggestQueueSignal = q.head.biggestQueueSignal
	} else {
		newNode.biggestQueueSignal = newNode.itemType
	}

	q.head = newNode

	return nil
}

// SinkPush is implementation for IQueue
func (q *MutexQueue) SinkPush(data interface{}) error {

	newNode := &QueueItem{
		payload: data,
	}

	return q.sinkPush(newNode)
}

// SinkPushAll is implementation for IQueue
func (q *MutexQueue) SinkPushAll(data []interface{}) error {
	inputSize := len(data)
	lastElement := &QueueItem{}
	newHead := lastElement

	lastElement.payload = data[inputSize-1]
	for i := inputSize - 2; i >= 0; i-- {
		lastElement.next = &QueueItem{}
		lastElement = lastElement.next
		lastElement.payload = data[i]
	}

	q.locker.Lock()
	defer q.locker.Unlock()

	if q.isQueueBlockedUnsafe() {
		return errors.New("[ SinkPushAll ] Queue is blocked")
	}

	if q.HasSignal() {
		// do smth interesting
	}

	nextIndex := q.head.index + 1
	tmpHead := newHead
	for i := inputSize - 1; i >= 0; i-- {
		tmpHead.index = nextIndex + uint(i)
		tmpHead = tmpHead.next
	}

	lastElement.next = q.head
	// lastNew.biggestQueueSignal = max(head.biggestQueueSignal, lastNew.type) // TODO: ? What is that ?

	q.head = newHead

	return nil
}

func (q *MutexQueue) checkAndGetHeadUnsafe() *QueueItem {
	if q.isQueueBlockedUnsafe() || q.isQueueEmptyUnsafe() {
		return nil
	}

	return q.head
}

func (q *MutexQueue) isQueueBlockedUnsafe() bool {
	return q.head == nil
}

func (q *MutexQueue) isQueueEmptyUnsafe() bool {
	return q.head == &emptyQueueItem
}

// get pointer to head and unfold linked list to slice:
//  all signals will be at the begging of the slice
func convertSublistToArray(localHead *QueueItem) []OutputElement {
	result := make([]OutputElement, localHead.index)

	current := localHead
	signalCurrentIndex := 0
	messageCurrentIndex := localHead.index - 1
	for i := uint(0); i < localHead.index; i++ {
		element := OutputElement{
			data:     current.payload,
			itemType: current.itemType,
		}
		if current.isSignal() {
			result[signalCurrentIndex] = element
			signalCurrentIndex++
		} else {
			result[messageCurrentIndex] = element
			messageCurrentIndex--
		}
		current = current.next
	}

	return result
}

// SinkPushAll is implementation for IQueue
func (q *MutexQueue) RemoveAll() []OutputElement {

	var localHead *QueueItem
	q.locker.Lock()
	localHead = q.checkAndGetHeadUnsafe()
	if localHead == nil {
		q.locker.Unlock()
		return []OutputElement{}
	}
	q.head = &emptyQueueItem
	q.locker.Unlock()

	return convertSublistToArray(localHead)
}

func (q *MutexQueue) BlockAndRemoveAll() []OutputElement {
	var localHead *QueueItem
	q.locker.Lock()
	localHead = q.checkAndGetHeadUnsafe()
	if localHead == nil {
		q.head = nil
		q.locker.Unlock()
		return []OutputElement{}
	}
	q.head = nil
	q.locker.Unlock()

	return convertSublistToArray(localHead)
}

// SinkPushAll is implementation for IQueue
func (q *MutexQueue) Unblock() bool {
	q.locker.Lock()
	defer q.locker.Unlock()

	if q.head != nil {
		return false
	}

	q.head = &emptyQueueItem
	return true
}

// SinkPushAll is implementation for IQueue
func (q *MutexQueue) PushSignal(signalType uint32, callback SyncDone) error {
	if signalType == 0 {
		return errors.New(fmt.Sprintf("[ PushSignal ] Unsupported signalType: %d", signalType))
	}

	newNode := &QueueItem{
		payload:            callback,
		biggestQueueSignal: signalType,
		itemType:           signalType,
	}

	return q.sinkPush(newNode)
}

// SinkPushAll is implementation for IQueue
// If queue is locked then it returns false
// Now it uses unsafe pointer. This function will be called very frequently, that is why we use atomic here
// But to be sure, this should be benchmarked in comparison with simple lock
func (q *MutexQueue) HasSignal() bool {
	head := (*QueueItem)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head))))
	return !q.isQueueBlockedUnsafe() && !q.isQueueEmptyUnsafe() && head.hasSignal()

}
