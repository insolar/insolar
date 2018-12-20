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

package transport

import (
	"sync"

	"github.com/insolar/insolar/network/transport/packet"
)

type futureManagerImpl struct {
	mutex   sync.RWMutex
	futures map[packet.RequestID]Future
}

func newFutureManagerImpl() *futureManagerImpl {
	return &futureManagerImpl{
		futures: make(map[packet.RequestID]Future),
	}
}

func (fm *futureManagerImpl) Create(msg *packet.Packet) Future {
	future := NewFuture(msg.RequestID, msg.Receiver, msg, func(f Future) {
		fm.delete(f.ID())
	})

	fm.mutex.Lock()
	defer fm.mutex.Unlock()

	fm.futures[msg.RequestID] = future

	return future
}

func (fm *futureManagerImpl) Get(msg *packet.Packet) Future {
	fm.mutex.RLock()
	defer fm.mutex.RUnlock()

	return fm.futures[msg.RequestID]
}

func (fm *futureManagerImpl) delete(id packet.RequestID) {
	fm.mutex.Lock()
	defer fm.mutex.Unlock()

	delete(fm.futures, id)
}
