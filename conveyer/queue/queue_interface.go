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

// SyncDone is callback for signal
type SyncDone interface {
	done()
}

// QueueItem is one item if the queue
type QueueItem struct {
	itemType uint32
	signal   uint32
	index    uint
	payload  interface{}
	next     *QueueItem
}

func (qi *QueueItem) isSignal() bool {
	return qi.itemType == 0
}

var emptyQueueItem QueueItem

func init() {
	emptyQueueItem = QueueItem{
		itemType: 0,
		signal:   0,
		index:    0,
		payload:  nil,
		next:     nil,
	}
}

// IQueue is interface for queue
type IQueue interface {
	// SinkPush adds one item to queue
	SinkPush(data interface{}) bool
	// SinkPushAll adds list of items to queue
	SinkPushAll(data []interface{}) bool
	// RemoveAll removes all elements from queue and return them
	RemoveAll() []interface{}
	// BlockAndRemoveAll like RemoveAll + lock queue after that
	BlockAndRemoveAll() []interface{}
	// Unblock unlock queue after BlockAndRemoveAll. If not locked it returns false
	Unblock() bool
	// PushSignal adds signal to queue
	PushSignal(signalType uint32, callback SyncDone) bool
	// HasSignal is true if queue has at least ont signal. Must be atomic. Without mutex
	HasSignal() bool
}
