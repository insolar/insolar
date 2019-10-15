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
	"sync"
	"sync/atomic"

	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachinev2"
	"github.com/insolar/insolar/conveyor/tools"
	"github.com/insolar/insolar/pulse"
)

type PreparedState = struct{}
type InputEvent = interface{}
type StateMachineFactoryFn = func(InputEvent, pulse.Number, PulseSlotState) smachine.StateMachine

func NewPulseConveyor(conveyorMachineConfig smachine.SlotMachineConfig, factoryFn StateMachineFactoryFn,
	slotMachineConfig smachine.SlotMachineConfig, registry injector.DependencyRegistry) *PulseConveyor {

	r := &PulseConveyor{
		slotConfig: PulseSlotConfig{
			config: slotMachineConfig,
		},
		past:      make(map[pulse.Number]*PulseSlotMachine),
		factoryFn: factoryFn,
	}
	//r.signalQueue = tools.NewSignalFuncQueue(&r.mutex, r.externalSignal.NextBroadcast)
	r.slotMachine = smachine.NewSlotMachine(conveyorMachineConfig,
		nil, r.externalSignal.NextBroadcast, registry)

	r.slotConfig.eventCallback = r.internalSignal.NextBroadcast
	r.slotConfig.signalCallback = nil // r.internalSignal.NextBroadcast
	r.slotConfig.parentRegistry = &r.slotMachine

	r.antique.isAntique = true
	r.antique.SlotMachine = smachine.NewSlotMachine(
		r.slotConfig.config,
		r.slotConfig.eventCallback,
		r.slotConfig.signalCallback,
		r.slotConfig.parentRegistry)

	return r
}

type PulseConveyor struct {
	presentPulse pulse.Number //atomic

	mutex sync.RWMutex

	slotConfig PulseSlotConfig

	future  *PulseSlotMachine
	present *PulseSlotMachine
	past    map[pulse.Number]*PulseSlotMachine
	antique PulseSlotMachine

	workerCtx      context.Context
	externalSignal tools.VersionedSignal
	internalSignal tools.VersionedSignal

	slotMachine smachine.SlotMachine

	//slotConfig   smachine.SlotMachineConfig
	factoryFn    StateMachineFactoryFn
	pulseService PulseServiceAdapter
}

func (p *PulseConveyor) GetPresentPulse() pulse.Number {
	return pulse.Number(atomic.LoadUint32((*uint32)(&p.presentPulse)))
}

func (p *PulseConveyor) GetPulseData(pn pulse.Number) pulse.Data {
	panic("unsupported")
}

func (p *PulseConveyor) getPulseSlot(pn pulse.Number) *PulseSlot {
	ps, ok := p.slotMachine.GetPublished(pn)
	if ok && ps == nil {
		panic("illegal state")
	}
	return ps.(*PulseSlot)
}

func (p *PulseConveyor) GetPulseSlot(pn pulse.Number) *PulseSlot {
	ips, ok := p.slotMachine.GetPublished(pn)
	if ok && ps == nil {
		panic("illegal state")
	}

}

func (p *PulseConveyor) AddInput(ctx context.Context, pn pulse.Number, event InputEvent) error {

	pulseSlotMachine, err := p.getPulseSlotMachine(ctx, pn)
	if err != nil {
		return err
	}

	pulseSlotMachine.AddNew(ctx, smachine.NoLink(),
		&inputEventSM{
			ps:        &pulseSlotMachine.pulseSlot,
			event:     event,
			factoryFn: p.factoryFn,
		})
	return nil
}

func (p *PulseConveyor) getPulseSlotMachine(ctx context.Context, pn pulse.Number) (*PulseSlotMachine, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	if p.present == nil {
		return nil, fmt.Errorf("uninitialized present pulse: pn=%v", pn)
	}

	presentPulseData := p.present.pulseSlot.PulseData()
	presentPulse := presentPulseData.PulseNumber
	npn := presentPulseData.GetNextPulseNumber()

	switch {
	case pn == presentPulse:
		return p.present, nil
	case pn >= presentPulseData.GetNextPulseNumber():
		if p.future == nil {
			return nil, fmt.Errorf("uninitialized future pulse: pn=%v", pn)
		}
		return p.future, nil
	case pn <= presentPulse:
		pulseSlot := p.past[pn]
		if pulseSlot == nil {
			pulseSlot = &p.antique
		}
		return pulseSlot, nil
	default:
		return nil, fmt.Errorf("invalid pulse: pn=%v expected=%v", pn, npn)
	}
}

func (p *PulseConveyor) PreparePulseChange(out chan<- PreparedState) {
	p.slotMachine.ScheduleSignal(func(worker smachine.FixedSlotWorker) {
		zzz
	})
}

func (p *PulseConveyor) CancelPulseChange() {
	p.slotMachine.ScheduleSignal(func(worker smachine.FixedSlotWorker) {
		zzz
	})
}

func (p *PulseConveyor) CommitPulseChange(pd pulse.Data) {
	pd.EnsurePulsarData()
	p.slotMachine.ScheduleSignal(func(worker smachine.FixedSlotWorker) {
		// set current pulse number
		current, ok := p.slotMachine.GetPublished(pd.PulseNumber)
		if !ok {

		}
		if ps, ok := current.(*PulseSlot); ok {
			pulseData
		}

		p.slotMachine.Migrate(worker)
		zzz
	})

	//p.signalQueue.Add(func(interface{}) {
	//	if p.present != nil {
	//		p.present.pulseSlot.pulseData.MakePast()
	//		p.present = nil
	//	}
	//
	//	if p.future != nil && !p.future.IsEmpty() {
	//		p.future.pulseSlot.pulseData.MakePresent(pd)
	//		p.present = p.future
	//		p.future = nil
	//	} else {
	//		p.present = newPresentPulseSlot(pd, p.slotConfig)
	//		p.slotMachine.AddNew(p.workerCtx, smachine.NoLink(),
	//			&PresentPulseSM{pulseSMTemplate{ps: p.present, psa: &p.pulseService}})
	//	}
	//
	//	futurePD := pd.CreateNextExpected()
	//	p.future = newFuturePulseSlot(futurePD, p.slotConfig)
	//
	//	p.slotMachine.AddNew(p.workerCtx, smachine.NoLink(),
	//		&PulseSM{pulseSMTemplate{ps: p.future, psa: &p.pulseService}})
	//	p.past[futurePD.PulseNumber] = p.future
	//
	//	p.pulseService.svc.onCommitPulseChange(pd)
	//	p.slotMachine.Migrate(true)
	//})
}

//func (p *PulseConveyor) StartWorker(ctx context.Context) {
//	if ctx == nil {
//		panic("illegal value")
//	}
//
//	p.mutex.Lock()
//	defer p.mutex.Unlock()
//
//	if p.workerCtx != nil {
//		panic("illegal state")
//	}
//	p.workerCtx = ctx
//
//	go p.workerConveyor()
//}
//
//func (p *PulseConveyor) workerConveyor() {
//	for {
//		worker := smachine.NewSimpleSlotWorker(p.externalSignal.Mark())
//
//		for !worker.HasSignal() {
//			select {
//			case <-p.workerCtx.Done():
//				p.externalSignal.NextBroadcast()
//				return
//			default:
//			}
//
//			mark := p.internalSignal.Mark()
//
//			for _, sig := range p.signalQueue.Flush() {
//				sig(nil)
//			}
//
//			p.slotMachine.ScanOnce(worker)
//
//			if !mark.HasSignal() {
//				// TODO we need a Yield to indicate "no work done"
//			}
//		}
//	}
//}
