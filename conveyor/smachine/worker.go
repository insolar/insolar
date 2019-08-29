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
)

type WorkerController interface {
	startWorkerDetachment(worker *SlotWorker)
	endWorkerDetachment(worker *SlotWorker)
}

type SlotWorker struct {
	workCtl WorkerController
	machine *SlotMachine

	mutex sync.Mutex
	cond  *sync.Cond
}

//func (p *SlotWorker) startSyncCall(ctx *slotContext) int32 {
//	if p.detachTimer != nil {
//		panic("illegal state")
//	}
//	lastState := atomic.LoadInt32(&p.detachedWorker)
//
//	timeout := p.machine.config.BeforeDetach
//	if timeout == 0 || timeout == math.MaxInt64 {
//		return lastState
//	}
//
//	p.detachTimer = time.AfterFunc(timeout, func() {
//		p.workCtl.startWorkerDetachment(p)
//		if !atomic.CompareAndSwapInt32(&p.detachedWorker, lastState, -1) {
//			p.workCtl.endWorkerDetachment(p)
//			return
//		}
//
//	})
//	return lastState
//}
//
//func (p *SlotWorker) endSyncCall(lastState int32) {
//	atomic.CompareAndSwapInt32(&p.detachedWorker, lastState, lastState+1)
//	p.detachTimer.Stop()
//	p.detachTimer = nil
//}

func (p *SlotWorker) getCond() *sync.Cond {
	if p.cond != nil {
		return p.cond
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.cond == nil {
		p.cond = sync.NewCond(&p.mutex)
	}
	return p.cond
}

func (p *SlotWorker) HasSignal() bool {
	return false
}

func (p *SlotWorker) GetLoopLimit() uint32 {
	return 5
}

func (p *SlotWorker) detachableCall(fn func()) (wasDetached bool, err error) {
	defer func() {
		err = recoverToErr("slot execution has failed", recover(), err)
	}()

	fn()
	return false, nil
}
