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

package sworker

import (
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/conveyor/tools"
)

//var _ smachine.AttachableSlotWorker = &AttachableWorker{}

type AttachableWorker struct {
	signalSource *tools.VersionedSignal
}

func (p *AttachableWorker) AttachTo(_ *smachine.SlotMachine, loopLimit uint32, fn smachine.AttachedFunc) (wasDetached bool) {
	//w := &SlotWorker{parent: p, outerSignal: p.signalSource.Mark(), loopLimit: loopLimit}
	//fn(w)
	return false
}

//var _ smachine.FixedSlotWorker = &SlotWorker{}

type SlotWorker struct {
	parent      *AttachableWorker
	outerSignal *tools.SignalVersion
	loopLimit   uint32
}

func (p *SlotWorker) HasSignal() bool {
	return p.outerSignal != nil && p.outerSignal.HasSignal()
}

func (*SlotWorker) IsDetached() bool {
	return false
}

func (p *SlotWorker) GetSignalMark() *tools.SignalVersion {
	return p.outerSignal
}

func (p *SlotWorker) OuterCall(*smachine.SlotMachine, smachine.NonDetachableFunc) (wasExecuted bool) {
	return false
}

func (p *SlotWorker) DetachableCall(fn smachine.DetachableFunc) (wasDetached bool) {
	//fn(&DetachableSimpleSlotWorker{p})
	return false
}
