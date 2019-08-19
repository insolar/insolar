///
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
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

	cond *sync.Cond
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
	if p.cond == nil {
		p.cond = sync.NewCond(&sync.Mutex{})
	}
	return p.cond
}

func (p *SlotWorker) HasSignal() bool {
	return false
}

func (p *SlotWorker) GetLoopLimit() int {
	return 5
}

func (p *SlotWorker) detachableCall(fn func()) (wasDetached bool, err error) {
	defer func() {
		err = recoverToErr("slot execution has failed", recover(), err)
	}()

	fn()
	return false, nil
}
