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

import "sort"

type SortedSlotDependencies struct {
	key          string
	items        []*SlotDep
	removedCount uint32
	hasAdded     bool
}

func (p *SortedSlotDependencies) Add(slot *Slot, weight int32, fn BroadcastReceiveFunc) SlotDependency {
	if slot == nil {
		panic("illegal value")
	}
	p.hasAdded = len(p.items) > 0
	r := &SlotDep{s: slot, fn: fn, c: p, w: weight}
	p.items = append(p.items, r)
	return r
}

func (p *SortedSlotDependencies) GetHead() SlotDependency {
	switch {
	case p.hasAdded:
		p.hasAdded = false
		break
	case len(p.items) == 0:
		return nil
	case p.items[0].s != nil:
		return p.items[0]
	}
	p.sort()
	if len(p.items) == 0 {
		return nil
	}
	return p.items[0]
}

func (p *SortedSlotDependencies) sort() {
	if uint32(len(p.items)) == p.removedCount {
		p.removedCount = 0
		p.items = nil
		return
	}
	sort.Stable(depSortHelper{p.items})
	if p.removedCount == 0 {
		return
	}
	p.items = p.items[:uint32(len(p.items))-p.removedCount]
	p.removedCount = 0
}

var _ sort.Interface = &depSortHelper{}

type depSortHelper struct {
	items []*SlotDep
}

func (d depSortHelper) Len() int {
	return len(d.items)
}

func (d depSortHelper) Less(i, j int) bool {
	switch {
	case d.items[j].s == nil:
		return d.items[i].s != nil
	case d.items[i].s == nil:
		return false
	default:
		return d.items[i].w > d.items[j].w
	}
}

func (d depSortHelper) Swap(i, j int) {
	d.items[i], d.items[j] = d.items[j], d.items[i]
}

var _ SlotDependency = &SlotDep{}

type SlotDep struct {
	w  int32
	s  *Slot
	fn BroadcastReceiveFunc
	c  *SortedSlotDependencies
}

func (s *SlotDep) OnSlotWorking() bool {
	panic("implement me")
}

func (s *SlotDep) OnStepChanged() bool {
	panic("implement me")
}

func (s *SlotDep) GetWeight() int32 {
	return s.w
}

func (s *SlotDep) Remove() {
	s.c.removedCount++
	s.s = nil
	s.fn = nil
}

func (s *SlotDep) GetKey() string {
	return s.c.key
}

func (s *SlotDep) OnSlotDisposed() {
	panic("implement me")
}

func (s *SlotDep) OnBroadcast(payload interface{}) (accepted, wakeup bool) {
	if s.fn == nil {
		return false, false
	}
	ac := asyncResultContext{slot: s.s}
	return ac.executeBroadcast(payload, s.fn)
}
