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

package smachine

import (
	"sync/atomic"
	"unsafe"
)

type dependencyStackController interface {
	IsReleaseOnStepping(link SlotLink, flags SlotDependencyFlags) bool
	IsReleaseOnWorking(link SlotLink, flags SlotDependencyFlags) bool

	// MUST NOT return nil
	AcquireParent(child *dependencyQueueEntry, flags SlotDependencyFlags) *dependencyQueueEntry
}

var _ SlotDependency = &dependencyStackEntry{}
var _ PostponedDependency = &dependencyStackEntry{}

type dependencyStackEntry struct {
	controller dependencyStackController
	child      *dependencyQueueEntry
	parent     *dependencyQueueEntry // atomic
	flags      dependencyStackFlags
}

var emptyQueueEntry = &dependencyQueueEntry{}

type dependencyStackFlags uint16

const (
	// indicates that the child was pre-acquired
	dependencyPartialRelease = 1 << iota
	slotFlagsOffset          = iota
)

func (p *dependencyStackEntry) _unsafeParent() *unsafe.Pointer {
	return (*unsafe.Pointer)((unsafe.Pointer)(&p.parent))
}

func (p *dependencyStackEntry) getParent() *dependencyQueueEntry {
	return (*dependencyQueueEntry)(atomic.LoadPointer(p._unsafeParent()))
}

func (p *dependencyStackEntry) getFlags() SlotDependencyFlags {
	return SlotDependencyFlags(p.flags >> slotFlagsOffset)
}

func (p *dependencyStackEntry) Release() (SlotDependency, []PostponedDependency, []StepLink) {
	if p.flags&dependencyPartialRelease != 0 {
		_, d, s := p.releaseParent()
		return p.child, d, s
	}
	d, s := p.ReleaseAll()
	return nil, d, s
}

func (p *dependencyStackEntry) releaseParent() (hasParent bool, d []PostponedDependency, s []StepLink) {
	// prevents postponed activation
	if atomic.CompareAndSwapPointer(p._unsafeParent(), nil, (unsafe.Pointer)(emptyQueueEntry)) {
		return false, nil, nil
	}

	parent := p.getParent()
	if parent == nil {
		panic("illegal state")
	}
	d, s = parent.ReleaseAll()
	return true, d, s
}

func (p *dependencyStackEntry) ReleaseAll() ([]PostponedDependency, []StepLink) {
	hasParent, d, s := p.releaseParent()
	if !hasParent {
		return p.child.ReleaseAll()
	}

	dc, sc := p.child.ReleaseAll()
	if d == nil {
		d = dc
	} else if dc != nil {
		d = append(d, dc...)
	}
	if s == nil {
		s = sc
	} else if sc != nil {
		s = append(s, sc...)
	}
	return d, s
}

func (p *dependencyStackEntry) IsReleaseOnWorking() bool {
	return p.controller.IsReleaseOnWorking(p.child.link, p.getFlags())
}

func (p *dependencyStackEntry) IsReleaseOnStepping() bool {
	return p.controller.IsReleaseOnStepping(p.child.link, p.getFlags())
}

func (p *dependencyStackEntry) IsCompatibleWith(flags SlotDependencyFlags) bool {
	f := p.getFlags()
	return f&flags == flags
}

func (p *dependencyStackEntry) ActivateStack(activateBy *dependencyQueueEntry, link StepLink) PostponedDependency {
	if p == nil {
		return nil
	}
	switch parent := p.getParent(); {
	case activateBy != p.child:
		// only child can call ActivateStack
		panic("illegal state")
	case parent == nil:
		// set the parent outside of this call to avoid deadlocks
	case parent != emptyQueueEntry:
		// parent can only be set when child was activated, and child is only activated once
		// so the only valid option for this case is release
		panic("illegal state")
	}
	return p
}

func (p *dependencyStackEntry) PostponedActivate(appendTo []StepLink) []StepLink {
	switch parent := p.getParent(); {
	case parent == nil:
		// set the parent
	case parent != emptyQueueEntry:
		// parent can only be set when child was activated, and child is only activated once
		// so the only valid option for this case is release
		panic("illegal state")
	}

	if !p.child.isOpen() {
		// this method can only be invoked when the child is open
		panic("illegal state")
	}

	switch parent := p.controller.AcquireParent(p.child, p.getFlags()); {
	case parent == nil:
		panic("illegal state")

	case !atomic.CompareAndSwapPointer(p._unsafeParent(), nil, (unsafe.Pointer)(parent)):
		d, s := parent.ReleaseAll()
		s = PostponedList(d).PostponedActivate(s)

		// a sync object can activate a holder only once
		// so this can only happen when a holder was released
		if parent = p.getParent(); parent != emptyQueueEntry {
			panic("illegal state") // TODO made all sync errors irrecoverable by SlotMachine?
		}
		return append(appendTo, s...)

	case parent.isOpen():
		if step, ok := p.child.link.GetStepLink(); ok {
			return append(appendTo, step)
		}
	}
	return appendTo
}

var _ PostponedDependency = &PostponedList{}

type PostponedList []PostponedDependency

func (p PostponedList) PostponedActivate(appendTo []StepLink) []StepLink {
	for _, d := range p {
		appendTo = d.PostponedActivate(appendTo)
	}
	return appendTo
}
