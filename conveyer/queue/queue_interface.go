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

// SyncDone is callback for biggestQueueSignal
type SyncDone interface {
	done()
}

// OutputElement represent one element returned from queue
type OutputElement struct {
	data     interface{}
	itemType uint32
}

type queueItem struct {
	itemType           uint32
	biggestQueueSignal uint32
	index              uint
	payload            interface{}
	next               *queueItem
}

func (qi *queueItem) isSignal() bool {
	return qi.itemType != 0
}

func (qi *queueItem) hasSignal() bool {
	return qi.biggestQueueSignal != 0
}

var emptyQueueItem queueItem

func init() {
	emptyQueueItem = queueItem{
		itemType:           0,
		biggestQueueSignal: 0,
		index:              0,
		payload:            nil,
		next:               nil,
	}
}

func maxSignal(signal1 uint32, signal2 uint32) uint32 {
	if signal1 > signal2 {
		return signal1
	}

	return signal2
}

// IQueue is interface for queue
//go:generate minimock -i github.com/insolar/insolar/conveyer/queue.IQueue -o ./ -s _mock.go
type IQueue interface {
	// SinkPush adds one item to queue
	SinkPush(data interface{}) error
	// SinkPushAll adds list of items to queue
	SinkPushAll(data []interface{}) error
	// RemoveAll removes all elements from queue and return them
	RemoveAll() []OutputElement
	// BlockAndRemoveAll like RemoveAll + lock queue after that
	BlockAndRemoveAll() []OutputElement
	// Unblock unlock queue after BlockAndRemoveAll. If not locked it returns false
	Unblock() bool
	// PushSignal adds biggestQueueSignal to queue
	PushSignal(signalType uint32, callback SyncDone) error
	// HasSignal is true if queue has at least one biggestQueueSignal. Must be atomic. Without mutex
	HasSignal() bool
}
