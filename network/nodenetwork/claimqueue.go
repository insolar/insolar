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
	"context"
	"sync"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/instrumentation/instracer"
)

type claimQueue struct {
	data []packets.ReferendumClaim
	lock sync.RWMutex
}

func newClaimQueue() *claimQueue {
	return &claimQueue{data: make([]packets.ReferendumClaim, 0)}
}

func (cq *claimQueue) Pop() packets.ReferendumClaim {
	ctx, span := instracer.StartSpan(context.Background(), "claimQueue.Pop wait lock")
	cq.lock.Lock()
	span.End()
	ctx, span = instracer.StartSpan(ctx, "claimQueue.Pop lock")
	defer span.End()
	defer cq.lock.Unlock()

	if len(cq.data) == 0 {
		return nil
	}
	result := cq.data[0]
	cq.data = cq.data[1:]
	return result
}

func (cq *claimQueue) Front() packets.ReferendumClaim {
	ctx, span := instracer.StartSpan(context.Background(), "claimQueue.Front wait lock")
	cq.lock.RLock()
	span.End()
	ctx, span = instracer.StartSpan(ctx, "claimQueue.Front lock")
	defer span.End()
	defer cq.lock.RUnlock()

	if len(cq.data) == 0 {
		return nil
	}
	return cq.data[0]
}

func (cq *claimQueue) Length() int {
	ctx, span := instracer.StartSpan(context.Background(), "claimQueue.Length wait lock")
	cq.lock.RLock()
	span.End()
	ctx, span = instracer.StartSpan(ctx, "claimQueue.Length lock")
	defer span.End()
	defer cq.lock.RUnlock()

	return len(cq.data)
}

func (cq *claimQueue) Push(claim packets.ReferendumClaim) {
	ctx, span := instracer.StartSpan(context.Background(), "claimQueue.Push wait lock")
	cq.lock.Lock()
	span.End()
	ctx, span = instracer.StartSpan(ctx, "claimQueue.Push lock")
	defer span.End()
	defer cq.lock.Unlock()

	cq.data = append(cq.data, claim)
}
