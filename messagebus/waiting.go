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

import "sync"

const MaxNextPulseMessagePool = 1000
const MaxNextPulseReplyPool = 1000

type waitingPool struct {
	waiting chan interface{}
	locker  sync.Mutex
	counter uint32
	limit   uint32
}

func newWaitingPool(limit uint32) waitingPool {
	return waitingPool{
		waiting: make(chan interface{}),
		limit:   limit,
	}
}

func (wp *waitingPool) clear() {
	tmp := wp.waiting
	wp.waiting = make(chan interface{})
	wp.locker.Lock()
	wp.counter = 0
	wp.locker.Unlock()
	close(tmp)
}

func (wp *waitingPool) acquire() bool {
	wp.locker.Lock()
	defer wp.locker.Unlock()

	if wp.counter >= wp.limit {
		return false
	}

	wp.counter++
	return true
}
