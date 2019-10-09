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
	"github.com/insolar/insolar/conveyor/smachine/tools"
	smachine "github.com/insolar/insolar/conveyor/smachinev2"

	"sync"
)

func NewSimpleSlotWorker(outerSignal tools.SignalVersion) *SimpleSlotWorker {
	return &SimpleSlotWorker{outerSignal: outerSignal}
}

var _ smachine.FixedSlotWorker = &SimpleSlotWorker{}

type SimpleSlotWorker struct {
	outerSignal tools.SignalVersion
	innerSignal func()
	cond        *sync.Cond
}

func (p *SimpleSlotWorker) ActivateLinkedList(linkedList *smachine.Slot, hotWait bool) {
	panic("implement me")
}

func (p *SimpleSlotWorker) HasSignal() bool {
	return false
}

func (*SimpleSlotWorker) IsDetached() bool {
	return false
}

func (p *SimpleSlotWorker) DetachableCall(fn smachine.DetachableFunc) (wasDetached bool) {
	fn(&DetachableSimpleSlotWorker{p})
	return false
}

func (p *DetachableSimpleSlotWorker) GetCond() (bool, *sync.Cond) {
	if p.cond == nil {
		p.cond = sync.NewCond(&sync.Mutex{})
	}
	return true, p.cond
}

var _ smachine.DetachableSlotWorker = &DetachableSimpleSlotWorker{}

type DetachableSimpleSlotWorker struct {
	*SimpleSlotWorker
}

func (p *DetachableSimpleSlotWorker) CanLoopOrHasSignal(loopCount int) (canLoop, hasSignal bool) {
	return loopCount < 100, false
}

func (p *DetachableSimpleSlotWorker) NonDetachableCall(fn smachine.NonDetachableFunc) (wasExecuted bool) {
	fn(p.SimpleSlotWorker)
	return true
}

//func (p *SimpleSlotWorker) FinishNested(state SlotMachineState) {
//}
//
//func (p *SimpleSlotWorker) DetachableCall(fn DetachableFunc) (wasDetached bool, err error) {
//	wCtx := simpleWorkerContext{p}
//
//	defer func() {
//		err = recoverSlotPanic("slot execution has failed", recover(), err)
//		wCtx.w = nil
//	}()
//
//	fn(wCtx)
//	return false, nil
//}
//
//func (p *SimpleSlotWorker) HasSignal() bool {
//	return p.outerSignal.HasSignal()
//}
//
//func (p *SimpleSlotWorker) getCond() *sync.Cond {
//	if p.cond == nil {
//		p.cond = sync.NewCond(&sync.Mutex{})
//	}
//	return p.cond
//}
//
//func (p *SimpleSlotWorker) wakeUpAfterSharedAccess(slot *Slot, link SlotLink) context.CancelFunc {
//	return func() {
//		if !link.IsValid() {
//			return
//		}
//		m := link.s.machine
//		m._applyInplaceUpdate(link.s, true, activateSlot)
//	}
//}
//
//type simpleWorkerContext struct {
//	w *SimpleSlotWorker
//}
//
//func (p simpleWorkerContext) AttachTo(slot *Slot, link SlotLink, wakeUpOnUse bool) (SharedAccessReport, context.CancelFunc) {
//	switch {
//	case !link.IsValid():
//		return SharedSlotAbsent, nil
//	case slot == link.s:
//		// no need to wakeup ourselves
//		return SharedSlotAvailableAlways, nil
//	}
//	if link.s.machine == nil {
//		panic("illegal state")
//	}
//
//	isRemote := slot.machine != link.s.machine
//	if link.s.isWorking() {
//		if isRemote {
//			return SharedSlotRemoteBusy, nil
//		}
//		return SharedSlotLocalBusy, nil
//	}
//
//	var finishFn context.CancelFunc
//	if wakeUpOnUse {
//		finishFn = p.w.wakeUpAfterSharedAccess(slot, link)
//	}
//
//	if isRemote {
//		return SharedSlotRemoteAvailable, finishFn
//	}
//	return SharedSlotAvailableAlways, finishFn
//}
//
//func (p simpleWorkerContext) CanLoopOrHasSignal(loopCount uint32) (canLoop, hasSignal bool) {
//	return loopCount < 10, p.w.HasSignal()
//}
//
//func (p simpleWorkerContext) StartNested(state SlotMachineState) SlotWorker {
//	return p.w
//}
//
//func (p simpleWorkerContext) HasSignal() bool {
//	return p.w.HasSignal()
//}
//
//func (p simpleWorkerContext) GetCond() (bool, *sync.Cond) {
//	if p.HasSignal() {
//		return false, nil
//	}
//	return true, p.w.getCond()
//}
