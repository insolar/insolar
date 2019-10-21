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
	"sync/atomic"
)

// Very simple implementation of a slot worker. No support for detachments.
func NewAttachableSimpleSlotWorker() *AttachableSimpleSlotWorker {
	return &AttachableSimpleSlotWorker{}
}

var _ smachine.AttachableSlotWorker = &AttachableSimpleSlotWorker{}

type AttachableSimpleSlotWorker struct {
	exclusive uint32
}

func (v *AttachableSimpleSlotWorker) AttachAsNested(m *smachine.SlotMachine, outer smachine.DetachableSlotWorker,
	loopLimit uint32, fn smachine.AttachedFunc) (wasDetached bool) {

	if !atomic.CompareAndSwapUint32(&v.exclusive, 0, 1) {
		panic("is attached")
	}
	defer atomic.StoreUint32(&v.exclusive, 0)

	w := &SimpleSlotWorker{outerSignal: outer.GetSignalMark(), loopLimitFn: outer.CanLoopOrHasSignal,
		machine: m, loopLimit: int(loopLimit)}

	w.init()
	fn(w)
	return false
}

func (v *AttachableSimpleSlotWorker) AttachTo(m *smachine.SlotMachine, signal *tools.SignalVersion,
	loopLimit uint32, fn smachine.AttachedFunc) (wasDetached bool) {

	if !atomic.CompareAndSwapUint32(&v.exclusive, 0, 1) {
		panic("is attached")
	}
	defer atomic.StoreUint32(&v.exclusive, 0)

	w := &SimpleSlotWorker{outerSignal: signal, machine: m, loopLimit: int(loopLimit)}

	w.init()
	fn(w)
	return false
}

var _ smachine.FixedSlotWorker = &SimpleSlotWorker{}

type SimpleSlotWorker struct {
	outerSignal *tools.SignalVersion
	loopLimitFn smachine.LoopLimiterFunc // NB! MUST correlate with outerSignal
	loopLimit   int

	machine *smachine.SlotMachine

	dsw DetachableSimpleSlotWorker
	nsw NonDetachableSimpleSlotWorker
}

func (p *SimpleSlotWorker) init() {
	p.dsw.SimpleSlotWorker = p
	p.nsw.SimpleSlotWorker = p
}

func (p *SimpleSlotWorker) HasSignal() bool {
	return p.outerSignal != nil && p.outerSignal.HasSignal()
}

func (*SimpleSlotWorker) IsDetached() bool {
	return false
}

func (p *SimpleSlotWorker) GetSignalMark() *tools.SignalVersion {
	return p.outerSignal
}

func (p *SimpleSlotWorker) CanLoopOrHasSignal(loopCount int) (canLoop, hasSignal bool) {
	switch {
	case p.loopLimitFn != nil:
		canLoop, hasSignal = p.loopLimitFn(loopCount)
		if loopCount >= p.loopLimit {
			canLoop = false
		}
		return canLoop, hasSignal

	case p.outerSignal.HasSignal():
		return false, true
	default:
		return loopCount < p.loopLimit, false
	}
}

func (p *SimpleSlotWorker) OuterCall(*smachine.SlotMachine, smachine.NonDetachableFunc) (wasExecuted bool) {
	return false
}

func (p *SimpleSlotWorker) DetachableCall(fn smachine.DetachableFunc) (wasDetached bool) {
	fn(&p.dsw)
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
	//fn(&p.nsw)
	return false
}

func (p *DetachableSimpleSlotWorker) NonDetachableCall(fn smachine.NonDetachableFunc) (wasExecuted bool) {
	fn(&NonDetachableSimpleSlotWorker{p.SimpleSlotWorker})
	return true
}

type NonDetachableSimpleSlotWorker struct {
	*SimpleSlotWorker
}

func (p *NonDetachableSimpleSlotWorker) DetachableCall(fn smachine.DetachableFunc) (wasDetached bool) {
	panic("not allowed") // this method shouldn't be accessible through interface
}
