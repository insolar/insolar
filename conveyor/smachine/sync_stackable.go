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

// methods of this interfaces can be protected by mutex
type dependencyStackController interface {
	ReleaseStacked(releasedBy *dependencyQueueEntry, flags SlotDependencyFlags)
}

type dependencyStackEntry struct {
	controller dependencyStackController
}

func (p *dependencyStackEntry) ActivateStack(activateBy *dependencyQueueEntry, link StepLink) PostponedDependency {
	return nil
}

func (p *dependencyStackEntry) ReleasedBy(entry *dependencyQueueEntry, flags SlotDependencyFlags) {
	if p == nil {
		return
	}
	p.controller.ReleaseStacked(entry, flags)
}
