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
