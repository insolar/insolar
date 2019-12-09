//
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
//

package filaments

import (
	"math"

	"github.com/insolar/insolar/ledger-v2/dropbag"
	"github.com/insolar/insolar/reference"
)

func NewLatest(entry *WriteEntry, lifelineRoot reference.Holder, filamentRoot reference.Holder,
	filamentSection dropbag.JetSectionId,
) *LocalSegment {
	switch {
	case entry == nil:
		panic("illegal value")
	case filamentRoot == nil:
		panic("illegal value")
	case entry.FilamentSeq == 0:
		panic("illegal value")
	case entry.Prev != nil:
		panic("illegal value")
	}
	ls := &LocalSegment{
		filamentRoot, lifelineRoot, filamentSection,
		entry, NewAtomicEntry(entry),
	}
	entry.Latest = ls

	return ls
}

type LocalSegment struct {
	filamentRoot reference.Holder // not-nil, time-based record ref // TODO make specific holder type
	lifelineRoot reference.Holder

	filamentSection dropbag.JetSectionId

	localEarliest *WriteEntry

	possibleLatest AtomicEntry
}

func (p *LocalSegment) checkEntry(entry, latest *WriteEntry) (*WriteEntry, bool) {
	switch {
	case entry == nil:
		panic("illegal value")
	case latest == nil:
		panic("illegal state")
	case latest == entry:
		return latest, false
	case entry.EventSeq == latest.EventSeq:
		panic("illegal state")
	default:
		return latest, latest.EventSeq < entry.EventSeq
	}
}

func (p *LocalSegment) Update(entry *WriteEntry) {
	latest, isNext := p.checkEntry(entry, p.possibleLatest.Get())
	switch {
	case entry.Prev == nil:
		panic("illegal value")
	case entry.Latest == nil:
		entry.Latest = p
	case entry.Latest != p:
		panic("illegal value")
	case !isNext:
		return
	}

	for prev := entry.Prev; prev != nil; prev = prev.Prev {
		switch {
		case prev == entry:
			panic("illegal state")
		case prev == latest:
			// it is safe to set even if p.possibleLatest was already updated - will only bring minor inefficiency
			p.possibleLatest.Set(entry)
			return
		}
	}

	panic("illegal value")
}

func (p *LocalSegment) AddNext(entry *WriteEntry) {
	latest, _ := p.checkEntry(entry, p.GetLatest())
	switch {
	case entry.Latest != nil:
		panic("illegal value")
	case entry.Prev != nil:
		panic("illegal value")
	}

	entry.Latest = p
	for {
		switch {
		case latest.FilamentSeq == math.MaxUint64:
			panic("illegal state - FilamentSeq overflow")
		case latest.FilamentSeq == 0:
			panic("illegal state")
		case latest.Next.Get() == nil:
			entry.FilamentSeq = latest.FilamentSeq + 1
			entry.Prev = latest

			if latest.Next.CmpAndSwap(nil, entry) {
				// it is safe to set even if p.possibleLatest was already updated - will only bring minor inefficiency
				p.possibleLatest.Set(entry)
				return
			}
		}
		latest = p.GetLatest()
	}
}

func (p *LocalSegment) GetLatest() *WriteEntry {
	for {
		latest := p.possibleLatest.Get()

		next := latest.Next.Get()
		if next == nil {
			return latest
		}

		for {
			n := next.Next.Get()
			if n != nil {
				continue
			}
			if p.possibleLatest.CmpAndSwap(latest, next) {
				return next
			}
			break
		}
	}
}
