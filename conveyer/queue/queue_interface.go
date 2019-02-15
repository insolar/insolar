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

type SyncDone interface {
	done()
}

type QueueItem struct {
	itemType uint
	signal   uint
	index    uint
	payload  interface{}
	next     *QueueItem
}

// Why do we need both this and isSignal functions?
func (qi *QueueItem) hasSignal() bool {
	return qi.signal != 0
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

type IQueue interface {
	SinkPush(data interface{}) bool
	SinkPushAll(data []interface{}) bool
	RemoveAll() []interface{}
	BlockAndRemoveAll() []interface{}
	Unblock() bool
	PushSignal(signalType uint, callback SyncDone) bool
	HasSignal() bool
}
