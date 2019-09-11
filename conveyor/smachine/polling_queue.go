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

package smachine

import "time"

type PollingQueue struct {
	prepared *poll

	polls   []*poll
	seqTail uint16
	seqLen  uint16
}

type poll struct {
	SlotQueue
	pollAfter time.Time
}

func (p *PollingQueue) Add(slot *Slot) {
	p.prepared.AddLast(slot)
}

func (p *PollingQueue) growPollingSlots() {
	sLen := len(p.polls)
	sizeInc := 10
	if sLen > 32 {
		sizeInc = sLen / 3
	}

	cp := make([]*poll, sLen, sLen+sizeInc)
	if p.seqTail != 0 {
		copy(cp, p.polls[p.seqTail:])
		copy(cp[p.seqTail:], p.polls[:p.seqTail])
		p.seqTail = 0
	} else {
		copy(cp, p.polls)
	}
	p.polls = cp

	bodies := make([]poll, sizeInc)
	for i := range bodies {
		bodies[i] = poll{SlotQueue: NewSlotQueue(PollingSlots)}
		p.polls = append(p.polls, &bodies[i])
	}
}

func (p *PollingQueue) FilterOut(scanTime time.Time, queue *SlotQueue) {
	if len(p.polls) == 0 {
		return
	}

	for {
		ps := p.polls[p.seqTail]

		if !ps.IsEmpty() && ps.pollAfter.After(scanTime) {
			break
		}

		queue.AppendAll(&ps.SlotQueue)

		if p.seqLen == 0 {
			return
		}
		p.seqLen--
		p.seqTail++

		if int(p.seqTail) >= len(p.polls) {
			p.seqTail = 0
		}
	}
}

func (p *PollingQueue) PrepareFor(pollTime time.Time) {

	seqHead := uint16(0)
	switch {
	case p.prepared == nil:
		if p.seqLen != 0 || p.seqTail != 0 {
			panic("illegal state")
		}
		if len(p.polls) == 0 {
			p.growPollingSlots()
		}

	case !p.prepared.IsEmpty():
		if p.prepared.pollAfter.Equal(pollTime) {
			return
		}

		p.seqLen++
		if int(p.seqLen) >= len(p.polls) {
			p.growPollingSlots()
		}
		seqHead = p.seqTail + p.seqLen
		if seqHead >= uint16(len(p.polls)) {
			seqHead -= uint16(len(p.polls))
		}
	}

	if !p.polls[seqHead].IsEmpty() {
		panic("illegal state")
	}
	p.prepared = p.polls[seqHead]
	p.prepared.pollAfter = pollTime
}

func (p *PollingQueue) GetNearestPollTime() time.Time {
	if p.seqLen > 0 {
		return p.polls[p.seqTail].pollAfter
	}

	if p.prepared == nil || p.prepared.IsEmpty() {
		return time.Time{}
	}

	return p.prepared.pollAfter
}
