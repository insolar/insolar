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

package nodenetwork

import (
	"sync"

	"github.com/insolar/insolar/consensus/packets"
)

// ClaimQueue is the queue that contains consensus claims.
type ClaimQueue interface {
	// Pop takes claim from the queue.
	Pop() packets.ReferendumClaim
	// Front returns claim from the queue without removing it from the queue.
	Front() packets.ReferendumClaim
	// Length returns the length of the queue
	Length() int
}

type claimQueue struct {
	data []packets.ReferendumClaim
	lock sync.RWMutex
}

func newClaimQueue() *claimQueue {
	return &claimQueue{data: make([]packets.ReferendumClaim, 0)}
}

func (cq *claimQueue) Pop() packets.ReferendumClaim {
	cq.lock.Lock()
	defer cq.lock.Unlock()

	if len(cq.data) == 0 {
		return nil
	}
	result := cq.data[0]
	cq.data = cq.data[1:]
	return result
}

func (cq *claimQueue) Front() packets.ReferendumClaim {
	cq.lock.RLock()
	defer cq.lock.RUnlock()

	if len(cq.data) == 0 {
		return nil
	}
	return cq.data[0]
}

func (cq *claimQueue) Length() int {
	cq.lock.RLock()
	defer cq.lock.RUnlock()

	return len(cq.data)
}

func (cq *claimQueue) Push(claim packets.ReferendumClaim) {
	cq.lock.Lock()
	defer cq.lock.Unlock()

	cq.data = append(cq.data, claim)
}
