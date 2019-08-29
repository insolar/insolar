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
	"sync/atomic"
)

func NoLink() SlotLink {
	return SlotLink{}
}

type SlotLink struct {
	id SlotID
	s  *Slot
}

func (p *SlotLink) SlotID() SlotID {
	return p.id
}

func (p *SlotLink) IsEmpty() bool {
	return p.s == nil
}

func (p *SlotLink) IsValid() bool {
	if p.s == nil {
		return false
	}
	return p.id == SlotID(atomic.LoadUint64(&p.s.idAndStep))
}

type StepLink struct {
	SlotLink
	step uint32
}

func (p *StepLink) IsAtStep() bool {
	if p.s == nil {
		return false
	}
	id, step := p.s.GetAtomicIDAndStep()
	return p.id == id && (p.step == 0 || p.step == step)
}
