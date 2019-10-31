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
	"fmt"
)

func NewSyncHierarchy(parent, child SyncLink) SyncLink {
	// needs no mutex
	return SyncLink{newHierarchySync(parent.controller, child.controller)}
}

func newHierarchySync(parentCtl, childCtl DependencyController) *hierarchySync {
	panic("not implemented")
}

var _ DependencyController = &hierarchySync{}
var _ dependencyStackController = &hierarchySync{}

type hierarchySync struct {
	name                string
	parentCtl, childCtl DependencyController
}

// ============== DependencyController ================

func (p *hierarchySync) CheckState() Decision {
	if d := p.childCtl.CheckState(); !d.IsPassed() {
		return d
	}
	return p.parentCtl.CheckState()
}

func (p *hierarchySync) CheckDependency(dep SlotDependency) Decision {
	if entry, ok := dep.(*dependencyStackEntry); ok {
		switch {
		case entry.controller != p:
			return Impossible
		case !entry.child.link.IsValid(): // just to make sure
			return Impossible
		}

		if d := p.childCtl.CheckDependency(entry.child); !d.IsPassed() {
			return d
		}
		if parent := entry.getParent(); parent != nil {
			return p.parentCtl.CheckDependency(entry.parent)
		}
		return NotPassed
	}
	return Impossible
}

func (p *hierarchySync) UseDependency(dep SlotDependency, flags SlotDependencyFlags) (Decision, SlotDependency) {
	if entry, ok := dep.(*dependencyStackEntry); ok {
		switch {
		case entry.controller != p:
			return Impossible, nil
		case !entry.child.link.IsValid(): // just to make sure
			return Impossible, nil
		case !entry.IsCompatibleWith(flags):
			return Impossible, nil
		}
		if d := p.childCtl.CheckDependency(entry.child); !d.IsPassed() {
			return d, nil
		}
		if parent := entry.getParent(); parent != nil {
			return p.parentCtl.CheckDependency(entry.parent), nil
		}
		return NotPassed, nil
	} else if d := p.childCtl.CheckDependency(dep); d.IsValid() {
		// partial acquire
		//child := dep.(*dependencyQueueEntry)
		//if d.IsPassed() {
		//	// create with parent
		//	dd, parent := p.parentCtl.CreateDependency(child.link, flags)
		//	return dd.GetDecision(), &dependencyStackEntry{p, child, parent, flags }
		//}
		panic("not implemented")
		//
		//p.childCtl.AttachAsChild(dep, *dependencyStackEntry)
		//return NotPassed, &dependencyStackEntry{p, child, nil, flags }
	} else {
		return Impossible, nil
	}
}

func (p *hierarchySync) CreateDependency(holder SlotLink, flags SlotDependencyFlags) (BoolDecision, SlotDependency) {
	panic("not implemented")
	//cd, child := p.childCtl.CreateDependency(holder, flags)
	//if !cd {
	//	//p.childCtl.AttachAsChild(dep, *dependencyStackEntry)
	//	//return NotPassed, &dependencyStackEntry{p, child, nil, flags }
	//
	//}
	//pd, parent := p.parentCtl.CreateDependency(holder, flags)
	//return pd, &dependencyStackEntry{p, child, parent, flags }
}

func (p *hierarchySync) newStackEntry() {

}

func (p *hierarchySync) GetLimit() (limit int, isAdjustable bool) {
	return -1, false
}

func (p *hierarchySync) AdjustLimit(limit int, absolute bool) (deps []StepLink, activate bool) {
	panic("illegal state")
}

func (p *hierarchySync) GetCounts() (active, inactive int) {
	return -1, -1
}

func (p *hierarchySync) GetName() string {
	if len(p.name) != 0 {
		return p.name
	}
	return fmt.Sprintf("sync-hierarchy-%p(%s -> %s)", p, p.childCtl.GetName(), p.parentCtl.GetName())
}

// ========= dependencyStackController ============

func (p *hierarchySync) IsReleaseOnStepping(link SlotLink, flags SlotDependencyFlags) bool {
	panic("implement me")
}

func (p *hierarchySync) IsReleaseOnWorking(link SlotLink, flags SlotDependencyFlags) bool {
	panic("implement me")
}

func (p *hierarchySync) AcquireParent(child *dependencyQueueEntry, flags SlotDependencyFlags) *dependencyQueueEntry {
	panic("implement me")
}
