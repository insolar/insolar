/*
 *    Copyright 2018 Insolar
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
	"container/list"
	"sync"
	"time"
)

type ExpiryQueue struct {
	mutex sync.Mutex
	items *list.List
	ttl   time.Duration
}

type Item struct {
	data     interface{}
	expireAt time.Time
}

func NewExpiryQueue(ttlTime time.Duration) *ExpiryQueue {
	eq := &ExpiryQueue{
		items: list.New(),
		ttl:   ttlTime,
	}
	eq.items.Init()
	go eq.startExpirationProcessing()
	return eq
}

func (eq *ExpiryQueue) PopValues() []interface{} {
	response := []interface{}{}
	if eq.items.Len() > 0 {
		eq.mutex.Lock()
		for {
			e := eq.items.Front()
			if e == nil {
				break
			}
			elem := eq.items.Remove(e)
			response = append(response, elem.(*Item).data)

		}
		eq.mutex.Unlock()
	}
	return response
}

func (eq *ExpiryQueue) Pop() interface{} {
	if eq.items.Len() > 0 {
		eq.mutex.Lock()
		elem := eq.items.Remove(eq.items.Front())
		eq.mutex.Unlock()
		return elem.(*Item).data
	}
	return nil
}

func (eq *ExpiryQueue) Push(msg interface{}) *list.Element {
	item := &Item{
		data:     msg,
		expireAt: time.Now().Add(eq.ttl),
	}
	return eq.items.PushBack(item)
}

func (eq *ExpiryQueue) startExpirationProcessing() {
	for {
		var sleepTime time.Duration
		if eq.items.Len() > 0 {

			eq.mutex.Lock()
			for {
				e := eq.items.Front()
				if e == nil || e.Value == nil {
					break
				}
				elem := e.Value.(*Item)
				sleepTime = time.Until(elem.expireAt)
				if sleepTime <= 0 {
					eq.items.Remove(e)
				} else {
					break
				}
			}
			eq.mutex.Unlock()

			time.Sleep(sleepTime)
		} else {
			time.Sleep(eq.ttl)
		}

	}
}
