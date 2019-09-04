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
	"sync"
	"sync/atomic"
)

type slotWorkerState = uint32

const (
	activeWorker slotWorkerState = 1 << iota
	detachableCall
	signalledWorker
	detachedWorker
)

type DetachableSlotWorker struct {
	// we only apply mutex for the following fields to ensure that tryDetach() has proper visibility
	cond  *sync.Cond
	state slotWorkerState

	mutex        sync.Mutex
	outerSignal  <-chan struct{}
	innerBreaker chan struct{}
}

func (p *DetachableSlotWorker) activate(outerSignal <-chan struct{}) {
	if outerSignal == nil {
		panic("illegal value")
	}
	if p.outerSignal != nil {
		panic("illegal state")
	}
	p.outerSignal = outerSignal
	//	p.innerBreaker = make(chan struct{})
}

func (p *DetachableSlotWorker) DetachableCall(fn DetachableFunc) (wasDetached bool, err error) {
	if !atomic.CompareAndSwapUint32(&p.state, activeWorker, detachableCall|activeWorker) {
		if atomic.LoadUint32(&p.state) == 0 {
			panic("illegal state - not initialized")
		} else {
			panic("illegal state - parallel access")
		}
	}

	p.startDetachableCall()
	defer func() {
		err = recoverToErr("slot execution has failed", recover(), err)

		if atomic.CompareAndSwapUint32(&p.state, detachableCall|activeWorker, activeWorker) { // fast path
			wasDetached = false
		} else {
			for {
				prev := atomic.LoadUint32(&p.state)
				if prev&(detachableCall|activeWorker) != detachableCall|activeWorker {
					panic("illegal state - parallel access")
				}

				if atomic.CompareAndSwapUint32(&p.state, prev, prev&^detachableCall) {
					wasDetached = prev&detachedWorker != 0
					break
				}
			}
		}
		p.endDetachableCall()
	}()

	fn(detachableWorkerContext{p})
	return true /* the worst case */, nil
}

func (p *DetachableSlotWorker) startDetachableCall() {
	//p.mutex.Lock()
	//if p.cond == nil {
	//	p.cond = sync.NewCond(&p.mutex)
	//}
	//p.mutex.Unlock()
}

func (p *DetachableSlotWorker) endDetachableCall() {
	p.mutex.Lock()
	if p.innerBreaker == nil { // fast path
		p.mutex.Unlock()
	}

	defer p.mutex.Unlock()
	close(p.innerBreaker)
}

func (p *DetachableSlotWorker) tryDetach() bool {

	p.mutex.Lock()
	defer p.mutex.Unlock()

	for {
		prev := atomic.LoadUint32(&p.state)
		switch {
		case prev == 0:
			panic("illegal state")
		case prev&activeWorker == 0:
			return true
		case prev&detachedWorker != 0:
			return true
		case prev&detachableCall == 0:
			if prev&signalledWorker != 0 || atomic.CompareAndSwapUint32(&p.state, prev, prev|signalledWorker) {
				return false
			}
			continue
		}
		if atomic.CompareAndSwapUint32(&p.state, prev, prev|detachedWorker) {
			break
		}
	}

	panic("not implemented")
	//switch {
	//case p.innerBreaker == nil:
	//	p.innerBreaker = make(chan struct{})
	//	p.outerSignal = breaker
	//	go p.workerDetachedBreaker()
	//case p.outerSignal != breaker:
	//	panic("illegal value")
	//}
	//
	//return true
}

func (p *DetachableSlotWorker) workerDetachedBreaker() {
outer:
	for {
		select {
		case <-p.outerSignal:
			p.setSignal()
			break outer
		case <-p.innerBreaker:
			select {
			case <-p.outerSignal:
				p.setSignal()
			default:
			}

			break outer
		}
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.cond != nil {
		p.cond.Broadcast()
	}
}

func (p *DetachableSlotWorker) hasSignal() bool {
	panic("not implemented")
	//return atomic.LoadUint32(&p.state) >= signalledWorking
}

func (p *DetachableSlotWorker) setSignal() {
	panic("not implemented")
	//outer:
	//	for {
	//		prev := atomic.LoadUint32(&p.state)
	//		switch prev {
	//		case readyToWork:
	//			if atomic.CompareAndSwapUint32(&p.state, readyToWork, signalledFinished) {
	//				return
	//			}
	//		case detachableWorking:
	//			if atomic.CompareAndSwapUint32(&p.state, detachableWorking, signalledWorking) {
	//				break outer
	//			}
	//		case signalledWorking, detachedWorking, signalledFinished:
	//			return
	//		default:
	//			panic("illegal state")
	//		}
	//	}
	//
	//	p.mutex.Lock()
	//	defer p.mutex.Unlock()
	//	if p.cond != nil {
	//		p.cond.Broadcast()
	//	}
}

func (p *DetachableSlotWorker) getCond() (bool, *sync.Cond) {

	panic("not implemented")
	//p.mutex.Lock()
	//defer p.mutex.Unlock()
	//
	//prev := atomic.LoadUint32(&p.state)
	//switch prev {
	//case detachableWorking, detachedWorking, signalledWorking:
	//	break
	//default:
	//	panic("illegal state")
	//}
	//
	//if p.cond != nil {
	//	p.cond = sync.NewCond(&p.mutex)
	//}
	//
	//return p.cond
}

func (p *DetachableSlotWorker) StartNested(state SlotMachineState) SlotWorker {
	return p
}

func (p *DetachableSlotWorker) FinishNested(state SlotMachineState) {
}

type detachableWorkerContext struct {
	w *DetachableSlotWorker
}

func (p detachableWorkerContext) StartNested(state SlotMachineState) SlotWorker {
	return p.w
}

func (p detachableWorkerContext) GetLoopLimit() uint32 {
	return 10
}

func (p detachableWorkerContext) HasSignal() bool {
	return p.w.hasSignal()
}

func (p detachableWorkerContext) GetCond() (bool, *sync.Cond) {
	return p.w.getCond()
}
