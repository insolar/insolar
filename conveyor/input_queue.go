///
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
///

package conveyor

import (
	"sync"
)

func NewInputQueue(bcast *sync.Cond) InputQueue {
	return InputQueue{bcast: bcast}
}

type InputQueueEvent func()

type InputQueue struct {
	mutex sync.Mutex
	bcast *sync.Cond

	hasSignal bool
	buffer    []InputQueueEvent
}

func (p *InputQueue) IsEmpty() bool {
	return p.bcast == nil && p.buffer == nil
}

func (p *InputQueue) Add(fn InputQueueEvent) {
	if fn == nil {
		panic("illegal value")
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.buffer = append(p.buffer, fn)
}

func (p *InputQueue) AddSignal(fn InputQueueEvent) {
	if fn == nil {
		panic("illegal value")
	}
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.buffer = append(p.buffer, fn)
	p.hasSignal = true

	if p.bcast != nil {
		// don't take a lock
		p.bcast.Broadcast()
	}
}

func (p *InputQueue) Flush() (events []InputQueueEvent, hasSignal bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	hasSignal = p.hasSignal
	p.hasSignal = false

	if len(p.buffer) == 0 {
		return nil, hasSignal
	}

	nextCap := cap(p.buffer)
	if nextCap > 128 && len(p.buffer)<<1 < nextCap {
		nextCap >>= 1
	}
	events = p.buffer
	p.buffer = make([]InputQueueEvent, 0, nextCap)

	return events, hasSignal
}
