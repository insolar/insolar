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
	"fmt"
)

const UnknownSlotID SlotID = 0

type SlotID uint32

func (id SlotID) IsUnknown() bool {
	return id == UnknownSlotID
}

func NoLink() SlotLink {
	return SlotLink{}
}

func NoStepLink() StepLink {
	return StepLink{}
}

/* A lazy link to a slot, that can detect if a slot is dead */
type SlotLink struct {
	id SlotID
	s  *Slot
}

func (p SlotLink) String() string {
	if p.s == nil {
		if p.id == 0 {
			return "<nil>"
		}
		return fmt.Sprintf("noslot-%d", p.id)
	}
	return fmt.Sprintf("slot-%d", p.id)
}

func (p SlotLink) SlotID() SlotID {
	return p.id
}

func (p SlotLink) IsEmpty() bool {
	return p.s == nil
}

func (p SlotLink) IsValid() bool {
	if p.s == nil {
		return false
	}
	id, _, _ := p.s.GetState()
	return p.id == id
}

func (p SlotLink) isValidAndBusy() bool {
	if p.s == nil {
		return false
	}
	id, _, isBusy := p.s.GetState()
	return p.id == id && isBusy
}

func (p SlotLink) getIsValidAndBusy() (isValid, isBusy bool) {
	if p.s == nil {
		return false, false
	}
	id, _, isBusy := p.s.GetState()
	return p.id == id, isBusy
}

func (p SlotLink) tryStartWorking() (s *Slot, isStarted bool, prevStepNo uint32) {
	if p.s != nil {
		if _, isStarted, prevStepNo = p.s._tryStartWithId(p.id, 1); isStarted {
			return p.s, true, prevStepNo
		}
	}
	return nil, false, 0
}

type StepLink struct {
	SlotLink
	step uint32
}

func (p StepLink) AnyStep() StepLink {
	if p.step != 0 {
		p.step = 0
	}
	return p
}

func (p StepLink) String() string {
	if p.step == 0 {
		return p.SlotLink.String()
	}
	return fmt.Sprintf("%s-step-%d", p.SlotLink.String(), p.id)
}

func (p StepLink) IsAtStep() bool {
	if p.s == nil {
		return false
	}
	id, step, _ := p.s.GetState()
	return p.id == id && (p.step == 0 || p.step == step)
}

func (p StepLink) isValidAndAtExactStep() (valid, atExactStep bool) {
	if p.s == nil {
		return false, false
	}
	id, step, _ := p.s.GetState()
	return p.id == id, p.step == step
}

func (p StepLink) getIsValidBusyAndAtStep() (isValid, isBusy, atExactStep bool) {
	if p.s == nil {
		return false, false, false
	}
	id, step, isBusy := p.s.GetState()
	return p.id == id, isBusy, p.step == 0 || p.step == step
}
