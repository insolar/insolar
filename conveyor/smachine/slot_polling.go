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

import (
	"sort"
	"time"
)

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

func (p *PollingQueue) AddToLatest(slot *Slot) {
	p.prepared.AddLast(slot)
}

func (p *PollingQueue) AddToLatestBefore(waitUntil time.Time, slot *Slot) bool {

	switch nextPoll := p.GetPreparedPollTime(); {
	case nextPoll.IsZero():
		return false
	case !waitUntil.Before(nextPoll):
		p.prepared.AddLast(slot)
		return true
	case p.seqLen == 0:
		return false
	case waitUntil.Before(p.polls[p.seqTail].pollAfter):
		return false
	case p.seqLen == 1:
		p.polls[p.seqTail].AddLast(slot)
		return true
	}

	base, count := int(p.seqTail)+1, int(p.seqLen)-1
	switch {
	case int(p.seqTail+p.seqLen) <= len(p.polls):
		// continuous range
	case waitUntil.Before(p.polls[0].pollAfter): // wrapped range - lets just split it in halves
		count = len(p.polls) - base
	default:
		count -= len(p.polls) - base
		base = 0
	}

	pos := base + sort.Search(count, func(i int) bool {
		return !waitUntil.Before(p.polls[base+i].pollAfter)
	})
	if pos == 0 {
		pos = len(p.polls) - 1
	} else {
		pos--
	}

	p.polls[pos].AddLast(slot)
	return true
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
		bodies[i].initSlotQueue(PollingSlots)
		p.polls = append(p.polls, &bodies[i])
	}
}

func (p *PollingQueue) FilterOut(scanTime time.Time, addSlots func(*SlotQueue)) {
	if len(p.polls) == 0 {
		return
	}

	for {
		ps := p.polls[p.seqTail]

		if !ps.IsEmpty() && ps.pollAfter.After(scanTime) {
			break
		}

		addSlots(&ps.SlotQueue)

		if !ps.SlotQueue.IsEmpty() || p.seqLen == 0 {
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
	default: // reuse the empty prepared
		p.prepared.pollAfter = pollTime
		return
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

func (p *PollingQueue) GetPreparedPollTime() time.Time {
	if p.prepared == nil {
		return time.Time{}
	}
	return p.prepared.pollAfter
}
