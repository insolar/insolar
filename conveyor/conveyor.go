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
	"context"
	"fmt"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/conveyor/smachine/tools"
	"github.com/insolar/insolar/pulse"
	"sync"
)

type PreparedState = struct{}
type InputEvent = interface{}
type StateMachineFactoryFn = func(InputEvent, PulseSlotState) smachine.StateMachine

func NewPulseConveyor(conveyorMachineConfig smachine.SlotMachineConfig, factoryFn StateMachineFactoryFn,
	slotMachineConfig smachine.SlotMachineConfig, injector smachine.DependencyInjector) *PulseConveyor {

	r := &PulseConveyor{
		past:       make(map[pulse.Number]*PulseSlot),
		antique:    PulseSlot{isPast: true},
		slotConfig: slotMachineConfig,
		factoryFn:  factoryFn,
		injector:   injector,
	}
	r.signalQueue = tools.NewSignalFuncQueue(&r.mutex, r.externalSignal.NextBroadcast)
	r.slotMachine = smachine.NewSlotMachine(conveyorMachineConfig, nil, nil)

	return r
}

type PulseConveyor struct {
	mutex sync.RWMutex

	future  *PulseSlot
	present *PulseSlot
	past    map[pulse.Number]*PulseSlot
	antique PulseSlot

	workerCtx      context.Context
	signalQueue    tools.SyncQueue
	externalSignal tools.VersionedSignal
	internalSignal tools.VersionedSignal

	slotMachine smachine.SlotMachine

	slotConfig   smachine.SlotMachineConfig
	injector     smachine.DependencyInjector
	factoryFn    StateMachineFactoryFn
	pulseService PulseServiceAdapter
}

func (m *PulseConveyor) GetAdapters() *smachine.SharedRegistry {
	return m.slotMachine.GetAdapters()
}

func (p *PulseConveyor) AddInput(ctx context.Context, pn pulse.Number, event InputEvent) error {

	pulseSlot, err := p.getPulseSlot(ctx, pn)
	if err != nil {
		return err
	}

	pulseSlot.slots.AddAsyncNew(ctx, smachine.NoLink(), &inputEventSM{ps: pulseSlot, event: event, factoryFn: p.factoryFn})
	return nil
}

func (p *PulseConveyor) getPulseSlot(ctx context.Context, pn pulse.Number) (*PulseSlot, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	pulseSlot := p.present
	if pulseSlot == nil {
		return nil, fmt.Errorf("uninitialized present pulse: pn=%v", pn)
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
		return nil, fmt.Errorf("invalid pulse: pn=%v expected=%v", pn, p.present.pd.GetNextPulseNumber())
	default:
		pulseSlot = p.future
		if pulseSlot == nil {
			return nil, fmt.Errorf("uninitialized future pulse: pn=%v", pn)
		}
	}
	return pulseSlot, nil
}

func (p *PulseConveyor) PreparePulseChange(out chan<- PreparedState) {
	p.signalQueue.Add(func(interface{}) {
		p.pulseService.svc.onPreparePulseChange(out)
	})
}

func (p *PulseConveyor) CancelPulseChange() {
	p.signalQueue.Add(func(interface{}) {
		p.pulseService.svc.onCancelPulseChange()
	})
}

func (p *PulseConveyor) CommitPulseChange(pd pulse.Data) {
	pd.EnsurePulsarData()

	p.signalQueue.Add(func(interface{}) {
		if p.present != nil {
			p.present.isPast = true
			p.present = nil
		}

		if p.future != nil {
			p.future.pd = pd
			p.present = p.future
			p.future = nil
		} else {
			p.present = newPulseSlot(pd,
				p.slotConfig, p.injector, p.GetAdapters())

			p.slotMachine.AddNew(p.workerCtx, smachine.NoLink(),
				&PresentPulseSM{pulseSMTemplate{ps: p.present, psa: &p.pulseService}})
		}

		p.future = newPulseSlot(pd.CreateNextExpected(),
			p.slotConfig, p.injector, p.GetAdapters())

		p.slotMachine.AddNew(p.workerCtx, smachine.NoLink(),
			&FuturePulseSM{pulseSMTemplate{ps: p.future, psa: &p.pulseService}})

		p.past[p.future.pd.PulseNumber] = p.future

		p.pulseService.svc.onCommitPulseChange(pd)
		p.slotMachine.Migrate(true)
	})
}

func (p *PulseConveyor) StartWorker(ctx context.Context) {
	if ctx == nil {
		panic("illegal value")
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.workerCtx != nil {
		panic("illegal state")
	}
	p.workerCtx = ctx

	go p.workerConveyor()
}

func (p *PulseConveyor) workerConveyor() {
	for {
		worker := smachine.NewSimpleSlotWorker(p.externalSignal.Mark())

		for !worker.HasSignal() {
			select {
			case <-p.workerCtx.Done():
				p.externalSignal.NextBroadcast()
				return
			default:
			}

			mark := p.internalSignal.Mark()

			for _, sig := range p.signalQueue.Flush() {
				sig(nil)
			}

			p.slotMachine.ScanOnce(worker)

			if !mark.HasSignal() {
				// TODO we need a Yield to indicate "no work done"
			}
		}
	}
}

type PulseSlotState uint8

const (
	Uninitialized PulseSlotState = iota
	Future
	Present
	Past
)

func newPulseSlot(pd pulse.Data, config smachine.SlotMachineConfig, injector smachine.DependencyInjector,
	adapters *smachine.SharedRegistry) *PulseSlot {

	r := &PulseSlot{pd: pd, slots: smachine.NewSlotMachine(config, injector, adapters)}
	r.slots.SetContainerState(r)
	return r
}

type PulseSlot struct {
	pd     pulse.Data
	isPast bool
	slots  smachine.SlotMachine
}

func (p *PulseSlot) SlotMachineState() {
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
