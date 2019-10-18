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

package example

import (
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/conveyor/tools"

	"sync"
)

func NewSimpleSlotWorker(outerSignal *tools.SignalVersion) *SimpleSlotWorker {
	return &SimpleSlotWorker{outerSignal: outerSignal}
}

var _ smachine.FixedSlotWorker = &SimpleSlotWorker{}

type SimpleSlotWorker struct {
	outerSignal *tools.SignalVersion
	innerSignal func()
	cond        *sync.Cond
}

func (p *SimpleSlotWorker) HasSignal() bool {
	return false
}

func (*SimpleSlotWorker) IsDetached() bool {
	return false
}

func (p *SimpleSlotWorker) GetSignalMark() *tools.SignalVersion {
	return p.outerSignal
}

func (p *SimpleSlotWorker) OuterCall(*smachine.SlotMachine, smachine.NonDetachableFunc) (wasExecuted bool) {
	panic("unsupported")
}

func (p *SimpleSlotWorker) DetachableCall(fn smachine.DetachableFunc) (wasDetached bool) {
	fn(&DetachableSimpleSlotWorker{p})
	return false
}

var _ smachine.DetachableSlotWorker = &DetachableSimpleSlotWorker{}

type DetachableSimpleSlotWorker struct {
	*SimpleSlotWorker
}

func (p *DetachableSimpleSlotWorker) TryDetach(flags smachine.LongRunFlags) {
	panic("unsupported")
}

func (p *DetachableSimpleSlotWorker) NonDetachableOuterCall(_ *smachine.SlotMachine, fn smachine.NonDetachableFunc) (wasExecuted bool) {
	fn(&NonDetachableSimpleSlotWorker{p.SimpleSlotWorker})
	return true
}

func (p *DetachableSimpleSlotWorker) CanLoopOrHasSignal(loopCount int) (canLoop, hasSignal bool) {
	return loopCount < 100, false
}

func (p *DetachableSimpleSlotWorker) NonDetachableCall(fn smachine.NonDetachableFunc) (wasExecuted bool) {
	fn(&NonDetachableSimpleSlotWorker{p.SimpleSlotWorker})
	return true
}

type NonDetachableSimpleSlotWorker struct {
	*SimpleSlotWorker
}

func (p *NonDetachableSimpleSlotWorker) DetachableCall(fn smachine.DetachableFunc) (wasDetached bool) {
	panic("not allowed")
}
