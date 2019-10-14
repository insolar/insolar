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

import "github.com/insolar/insolar/conveyor/injector"

func (s *Slot) activateSlot(worker FixedSlotWorker) {
	s.machine.updateSlotQueue(s, worker, activateSlot)
}

func (p SlotLink) activateSlot(worker FixedSlotWorker) {
	if p.IsValid() {
		p.s.activateSlot(worker)
	}
}

func (s *Slot) releaseDependency(worker FixedSlotWorker) {
	dep := s.dependency
	if dep == nil {
		return
	}
	s.dependency = nil
	dep.Release(func(link SlotLink) {
		s.machine.activateDependantByLink(link, worker)
	})
}

// MUST be a busy-holder to use
func (s *Slot) wakeUpSlot(worker DetachableSlotWorker) bool {
	if s.slotFlags&slotWokenUp != 0 || s.QueueType().IsActiveOrPolling() {
		return false
	}
	s.slotFlags |= slotWokenUp

	if !worker.NonDetachableCall(s.activateSlot) {
		s.machine.syncQueue.AddAsyncUpdate(s.NewLink(), SlotLink.activateSlot)
	}
	return true
}

func buildShadowMigrator(c injector.ReadOnlyContainer, defFn ShadowMigrateFunc) ShadowMigrateFunc {
	count := c.Count()
	if defFn != nil {
		count++
	}
	shadowMigrates := make([]ShadowMigrateFunc, 0, count)

	c.FilterLocalDependencies(func(id string, v interface{}) bool {
		if smFn, ok := v.(ShadowMigrator); ok {
			shadowMigrates = append(shadowMigrates, smFn.ShadowMigrate)
		}
		return false
	})

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
