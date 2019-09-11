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
	id, _ := p.s.GetAtomicIDAndStep()
	return p.id == id
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

func (p StepLink) IsAtStep() bool {
	if p.s == nil {
		return false
	}
	id, step := p.s.GetAtomicIDAndStep()
	return p.id == id && (p.step == 0 || p.step == step)
}

func (p StepLink) isValidAndAtExactStep() (valid, atExactStep bool) {
	if p.s == nil {
		return false, false
	}
	id, step := p.s.GetAtomicIDAndStep()
	return p.id == id, p.step == step
}

type SharedDataFunc func(interface{})

type SharedDataLink struct {
	link   StepLink
	wakeup bool
	data   interface{}
}

func (v SharedDataLink) PrepareAccess(fn SharedDataFunc) SharedDataAccessor {
	return SharedDataAccessor{v, fn}
}

type SharedDataAccessor struct {
	link     SharedDataLink
	accessFn SharedDataFunc
}

type SharedAccessReport uint8

const (
	SharedSlotAbsent SharedAccessReport = iota
	_
	SharedSlotLocalAvailable
	SharedSlotLocalBusy
	SharedSlotRemoteAvailable
	SharedSlotRemoteBusy
)

func (v SharedAccessReport) IsAvailable() bool {
	return v == SharedSlotLocalAvailable || v == SharedSlotRemoteAvailable
}

func (v SharedAccessReport) IsRemote() bool {
	return v == SharedSlotRemoteBusy || v == SharedSlotRemoteAvailable
}

func (v SharedAccessReport) IsAbsent() bool {
	return v == SharedSlotAbsent
}
