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

package conveyor

import (
	"fmt"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/network/consensus/common/pulse"
	"sync"
)

func NewPulseConveyor(config smachine.SlotMachineConfig) *PulseConveyor {

	bcast := sync.NewCond(&sync.Mutex{})
	return &PulseConveyor{
		slotMachine:     smachine.NewSlotMachine(config),
		past:            make(map[pulse.Number]*PulseSlot),
		signalBroadcast: bcast,
		antique:         PulseSlot{isPast: true, inputQueue: NewInputQueue(bcast)},
	}
}

type PulseConveyor struct {
	mutex sync.RWMutex

	future  *PulseSlot
	present *PulseSlot
	past    map[pulse.Number]*PulseSlot
	antique PulseSlot

	signalBroadcast *sync.Cond

	slotMachine  smachine.SlotMachine
	pulseService PulseServiceAdapter
}

func (p *PulseConveyor) ScanOnce(workCtl smachine.WorkerController) bool {
	if p.slotMachine.IsEmpty() {
		p.slotMachine.AddNew(smachine.NoLink(), &PastPulseSM{pulseSMTemplate{ps: &p.antique, psa: &p.pulseService}})
	}

	return p.slotMachine.ScanOnce(workCtl)
}

func (p *PulseConveyor) AddInput(pn pulse.Number, event InputQueueEvent) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	pulseSlot := p.present
	if pulseSlot == nil {
		return fmt.Errorf("uninitialized present pulse: pn=%v", pn)
	}

	presentPulse := p.present.pd.PulseNumber

	switch {
	case pn < presentPulse:
		pulseSlot = p.past[pn]
		if pulseSlot == nil {
			pulseSlot = &p.antique
		}
	case pn == presentPulse:
		// pulseSlot = p.present
	case pn < p.present.pd.GetNextPulseNumber():
		return fmt.Errorf("invalid pulse: pn=%v expected=%v", pn, p.present.pd.GetNextPulseNumber())
	default:
		pulseSlot = p.future
		if pulseSlot == nil {
			return fmt.Errorf("uninitialized future pulse: pn=%v", pn)
		}
	}
	pulseSlot.inputQueue.Add(event)
	return nil
}

func (p *PulseConveyor) PreparePulseChange() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

}

func (p *PulseConveyor) CommitPulseChange(pd pulse.Data) {
	pd.EnsurePulsarData()

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.present != nil {
		p.present.isPast = true
		p.present = nil
	}

	if p.future != nil {
		p.future.pd = pd
		p.present = p.future
	} else {
		p.present = &PulseSlot{
			pd:         pd,
			inputQueue: NewInputQueue(p.signalBroadcast),
		}
		p.slotMachine.AddNew(smachine.NoLink(), &PresentPulseSM{pulseSMTemplate{ps: p.present, psa: &p.pulseService}})
	}

	p.future = &PulseSlot{
		pd:         pd.CreateNextExpected(),
		inputQueue: NewInputQueue(p.signalBroadcast),
	}
	p.slotMachine.AddNew(smachine.NoLink(), &FuturePulseSM{pulseSMTemplate{ps: p.future, psa: &p.pulseService}})

	p.past[p.future.pd.PulseNumber] = p.future
}

func (p *PulseConveyor) CancelPulseChange() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

}

type PulseSlotState uint8

const (
	Uninitialized PulseSlotState = iota
	Future
	Present
	Past
)

type PulseSlot struct {
	pd          pulse.Data
	isPast      bool
	inputQueue  InputQueue
	slotMachine smachine.SlotMachine
}

func (p *PulseSlot) State() PulseSlotState {
	switch {
	case p.isPast:
		return Past
	case p.pd.IsEmpty():
		return Uninitialized
	case p.pd.IsExpectedPulse():
		return Future
	default:
		return Present
	}
}

func (p *PulseSlot) processEvents(ctx smachine.ExecutionContext, suspending bool) smachine.StateUpdate {
	events, _ := p.inputQueue.Flush()
	//if len(events) == 0 {
	//	// pass a signal to the state machine
	//	return ctx.WaitAny()
	//}

	for _, ev := range events {
		ev()
	}

	return ctx.Yield()
}

func (p *PulseSlot) processEventsAndOperations(ctx smachine.ExecutionContext) smachine.StateUpdate {
	p.processEvents(ctx, false)
	p.slotMachine.ScanOnceAsNested(ctx)
	return ctx.Yield()
}
