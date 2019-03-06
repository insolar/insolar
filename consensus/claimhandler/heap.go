/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package claimhandler

import (
	"bytes"
	"container/heap"

	"github.com/insolar/insolar/consensus/packets"
)

// Queue implements heap.Interface.
type Queue []*Claim

type Claim struct {
	value    packets.ReferendumClaim
	priority []byte
	index    int
}

func (q *Queue) PushClaim(claim packets.ReferendumClaim, priority []byte) {
	item := &Claim{
		value:    claim,
		index:    q.Len(),
		priority: priority,
	}
	heap.Push(q, item)
}

func (q *Queue) Push(x interface{}) {
	item := x.(*Claim)
	*q = append(*q, item)
}

func (q *Queue) PopClaim() packets.ReferendumClaim {
	return heap.Pop(q).(packets.ReferendumClaim)
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
	return bytes.Compare(q[i].priority, q[j].priority) > 0
}
