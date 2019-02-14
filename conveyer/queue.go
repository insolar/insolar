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

package conveyer

import (
	"fmt"
	"sync"
)

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

type Queue struct {
	head *QueueItem
}

func NewQueue() *Queue {
	queue := &Queue{
		head: &emptyQueueItem,
	}

	return queue
}

func (q *Queue) SinkPush(data interface{}) bool {
	fmt.Println("Pushing: ", data)
	newNode := &QueueItem{
		payload: data,
		next:    q.head,
	}

	q.head = newNode

	return true
}

func (q *Queue) SinkPushAll(data []interface{}) bool {
	for el := range data {
		res := q.SinkPush(el)
		if !res {
			fmt.Println("[ SinkPushAll ] Couldn't sinkPush of", el)
		}

	}

	return true
}

func (q *Queue) RemoveAll() [](interface{}) {
	return []interface{}{}
}

func (q *Queue) Unblock() bool {
	return true
}

func (q *Queue) PushSignal(signalType uint, callback SyncDone) bool {
	return q.SinkPush(signalType)
}

func (q *Queue) Next() (*QueueItem, error) {
	if *q.head == emptyQueueItem {
		return nil, fmt.Errorf("Empty queue")
	}
	retElement := q.head
	q.head = q.head.next

	return retElement, nil
}

func (q *Queue) HasSignal() bool {
	return false
}

func main() {
	queue := NewQueue()

	parallel := 200
	wg := sync.WaitGroup{}
	wg.Add(parallel)

	numIterations := 20

	for i := 0; i < parallel; i++ {
		go func(wg *sync.WaitGroup, q *Queue) {
			fmt.Println("START")
			for i := 0; i < numIterations; i++ {
				q.SinkPush(i)
			}
			wg.Done()
		}(&wg, queue)
	}

	wg.Wait()

	numElement := 0
	for true {
		element, err := queue.Next()
		if err != nil {
			fmt.Printf("Got error %s. Stop.\n", err)
			break
		}
		numElement++

		fmt.Println("Next element: ", element.payload)
	}

	fmt.Printf("Num Elements: %d . Must be: %d\n", numElement, parallel*numIterations)

}
