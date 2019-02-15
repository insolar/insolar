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

package claimhandler

import (
	"bytes"
	"container/heap"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
)

// Queue implements heap.Interface.
type Queue []*Claim

type Claim struct {
	value    packets.ReferendumClaim
	priority []byte
	index    int
}

func (q *Queue) Push(x interface{}) {
	item := Claim{
		value: x.(packets.ReferendumClaim),
		index: q.Len(),
	}
	item.index = q.Len()
	*q = append(*q, item)
	heap.Fix(q, item.index)
}

func (q *Queue) Pop() interface{} {
	l := q.Len()
	item := (*q)[l-1]
	*q = (*q)[0 : l-1]
	return item.value
}

func (q Queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q Queue) Len() int {
	return len(q)
}

// Less returns true if i > j cuz we need a greater to pop. Otherwise returns false.
func (q Queue) Less(i, j int) bool {
	if bytes.Compare(q[i].priority, q[j].priority) > 0 {
		return true
	}
	return false
}

func (q *Queue) getPriority(ref core.RecordRef, entropy core.Entropy) []byte {
	res := make([]byte, len(ref))
	for i := 0; i < len(ref); i++ {
		res[i] = ref[i] ^ entropy[i]
	}
	return res
}
