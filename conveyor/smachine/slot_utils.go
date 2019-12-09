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

func (s *Slot) activateSlot(worker FixedSlotWorker) {
	s.machine.updateSlotQueue(s, worker, activateSlot)
}

func (p SlotLink) activateSlot(worker FixedSlotWorker) {
	if p.IsValid() {
		p.s.activateSlot(worker)
	}
}

func (p StepLink) activateSlotStep(worker FixedSlotWorker) {
	if p.IsAtStep() {
		p.s.activateSlot(worker)
	}
}

func (p StepLink) activateSlotStepWithSlotLink(_ SlotLink, worker FixedSlotWorker) {
	p.activateSlotStep(worker)
}

func (p SlotLink) activateOnNonBusy(waitOn SlotLink, worker DetachableSlotWorker) bool {
	switch {
	case !p.IsValid():
		// requester is dead, don't wait anymore and don't wake it up
		return true
	case worker == nil:
		// too many retries - have to wake up the requester
	case waitOn.isValidAndBusy():
		// someone got it already, this callback should be added back to the queue
		return false
	}

	// lets try to wake up with synchronous operations
	m := p.getActiveMachine()
	switch {
	case m == nil:
		// requester is dead, don't wait anymore and don't wake it up
		return true
	case worker == nil:
		//
	case !waitOn.isMachine(m):
		if worker.NonDetachableOuterCall(m, p.s.activateSlot) {
			return true
		}
	case worker.NonDetachableCall(p.s.activateSlot):
		return true
	}

	// last resort - try to wake up via asynchronous
	m.syncQueue.AddAsyncUpdate(p, SlotLink.activateSlot)
	return true
}

func buildShadowMigrator(localInjects []interface{}, defFn ShadowMigrateFunc) ShadowMigrateFunc {
	count := len(localInjects)
	if defFn != nil {
		count++
	}
	shadowMigrates := make([]ShadowMigrateFunc, 0, count)

	for _, v := range localInjects {
		if smFn, ok := v.(ShadowMigrator); ok {
			shadowMigrates = append(shadowMigrates, smFn.ShadowMigrate)
		}
	}

	switch {
	case len(shadowMigrates) == 0:
		return defFn
	case defFn != nil:
		shadowMigrates = append(shadowMigrates, defFn)
	}
	if len(shadowMigrates)+1 < cap(shadowMigrates) { // allow only a minimal oversize
		shadowMigrates = append([]ShadowMigrateFunc(nil), shadowMigrates...)
	}

	return func(start, delta uint32) {
		for _, fn := range shadowMigrates {
			fn(start, delta)
		}
	}
}

func (s *Slot) _releaseDependency() []StepLink {
	dep := s.dependency
	s.dependency = nil
	replace, postponed, released := dep.Release()
	s.dependency = replace

	released = PostponedList(postponed).PostponedActivate(released)
	return released
}

var _ PostponedDependency = &PostponedList{}

type PostponedList []PostponedDependency

func (p PostponedList) PostponedActivate(appendTo []StepLink) []StepLink {
	for _, d := range p {
		appendTo = d.PostponedActivate(appendTo)
	}
	return appendTo
}
