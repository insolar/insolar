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
	"sync/atomic"

	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachinev2"
	"github.com/insolar/insolar/conveyor/tools"
	"github.com/insolar/insolar/pulse"
)

type PreparedState = struct{}
type InputEvent = interface{}
type PulseEventFactoryFunc = func(pulse.Number, InputEvent) smachine.CreateFunc

func NewPulseConveyor(conveyorMachineConfig smachine.SlotMachineConfig, factoryFn PulseEventFactoryFunc,
	slotMachineConfig smachine.SlotMachineConfig, registry injector.DependencyRegistry) *PulseConveyor {

	r := &PulseConveyor{
		slotConfig: PulseSlotConfig{
			config: slotMachineConfig,
		},
		factoryFn: factoryFn,
	}
	//r.signalQueue = tools.NewSignalFuncQueue(&r.mutex, r.externalSignal.NextBroadcast)
	r.slotMachine = smachine.NewSlotMachine(conveyorMachineConfig,
		nil, r.externalSignal.NextBroadcast, registry)

	r.slotConfig.eventCallback = r.internalSignal.NextBroadcast
	r.slotConfig.signalCallback = nil // r.internalSignal.NextBroadcast
	r.slotConfig.parentRegistry = &r.slotMachine

	//r.antique.isAntique = true
	//r.antique.SlotMachine = smachine.NewSlotMachine(
	//	r.slotConfig.config,
	//	r.slotConfig.eventCallback,
	//	r.slotConfig.signalCallback,
	//	r.slotConfig.parentRegistry)

	return r
}

type PulseConveyor struct {
	// immutable, provided, set at construction
	slotConfig PulseSlotConfig
	factoryFn  PulseEventFactoryFunc
	workerCtx  context.Context
	//slotConfig   smachine.SlotMachineConfig

	// immutable, set at construction
	externalSignal tools.VersionedSignal
	internalSignal tools.VersionedSignal
	slotMachine    smachine.SlotMachine

	// mutable
	presentAndFuturePulse uint64 //atomic

	// mutable, set under SlotMachine synchronization
	presentMachine *PulseSlotMachine
}

const uninitializedFuture = pulse.LocalRelative

func (p *PulseConveyor) GetPresentPulse() (present pulse.Number, nearestFuture pulse.Number) {
	v := atomic.LoadUint64(&p.presentAndFuturePulse)
	if v == 0 {
		return pulse.Unknown, uninitializedFuture
	}
	return pulse.Number(v), pulse.Number(v >> 32)
}

func (p *PulseConveyor) setPresentPulse(pd pulse.Data) (prevPresent pulse.Number, prevFuture pulse.Number, err error) {
	for {
		prev := atomic.LoadUint64(&p.presentAndFuturePulse)
		presentPN := pd.PulseNumber
		futurePN := pd.GetNextPulseNumber()

		if prev != 0 {
			expectedPN := pulse.Number(prev >> 32)
			if pd.PulseNumber < expectedPN {
				return pulse.Number(prev), pulse.Number(prev >> 32),
					fmt.Errorf("illegal pulse data: pn=%v, expected=%v", presentPN, expectedPN)
			}
		}
		// TODO store pulse.Data
		if atomic.CompareAndSwapUint64(&p.presentAndFuturePulse, prev, uint64(presentPN)|uint64(futurePN)<<32) {
			return pulse.Number(prev), pulse.Number(prev >> 32), nil
		}
	}
}

func (p *PulseConveyor) GetPulseData(pn pulse.Number) (pulse.Data, bool) {
	panic("unimplemented")
}

func (p *PulseConveyor) HasPulseData(pn pulse.Number) bool {
	panic("unimplemented")
}

func (p *PulseConveyor) AddInput(ctx context.Context, pn pulse.Number, event InputEvent) error {

	pulseSlotMachine, targetPN, pulseState, err := p.mapToPulseSlotMachine(pn)
	switch {
	case err != nil:
		return err
	case pulseSlotMachine == nil || pulseState == 0:
		return fmt.Errorf("slotMachine is missing: pn=%v", pn)
	}

	createFn := p.factoryFn(targetPN, event)

	switch {
	case createFn == nil:
		return fmt.Errorf("unrecognized event: pn=%v event=%v", targetPN, event)
	case pulseState == Future:
		// event for future need special handling
		pulseSlotMachine.machine.AddNew(ctx, smachine.NoLink(),
			&futureEventSM{pn: targetPN, ps: &pulseSlotMachine.pulseSlot, createFn: createFn})
	default:
		// TODO Functions to control max future and min past pulse
		// TODO here we should only check for recent data - very old data should be retrieved via a special SM
		if !p.HasPulseData(targetPN) {
			return fmt.Errorf("unknown data for pulse : pn=%v event=%v", targetPN, event)
		}
		if _, ok := pulseSlotMachine.machine.AddNewByFunc(ctx, smachine.NoLink(), createFn); !ok {
			return fmt.Errorf("ignored event: pn=%v event=%v", targetPN, event)
		}
	}
	return nil
}

func (p *PulseConveyor) mapToPulseSlotMachine(pn pulse.Number) (*PulseSlotMachine, pulse.Number, PulseSlotState, error) {
	presentPN, futurePN := p.GetPresentPulse()

	switch {
	case presentPN.IsUnknown():
		// when no present pulse - all pulses go to future
		return p.getFuturePulseSlotMachine(futurePN), pn, Future, nil
	case pn.IsUnknownOrEqualTo(presentPN):
		if psm := p.getPulseSlotMachine(presentPN); psm != nil {
			return psm, presentPN, Present, nil
		}
		// present slot must be present ;)
		panic("illegal state")
	case pn < presentPN:
		if psm := p.getPulseSlotMachine(pn); psm != nil {
			return psm, pn, Past, nil
		}
		if pn.IsTimePulse() {
			return p.getAntiquePulseSlotMachine(), pn, Antique, nil
		}
		return nil, 0, 0, fmt.Errorf("pulse number is invalid: pn=%v", pn)
	case pn < futurePN:
		return nil, 0, 0, fmt.Errorf("pulse number is unexpected: pn=%v", pn)
	default: // pn >= futurePN
		return p.getFuturePulseSlotMachine(futurePN), pn, Future, nil
	}
}

func (p *PulseConveyor) getPulseSlotMachine(pn pulse.Number) *PulseSlotMachine {
	if psv, ok := p.slotMachine.GetPublished(pn); ok {
		if psm, ok := psv.(*PulseSlotMachine); ok {
			return psm
		}
		panic("illegal state")
	}
	return nil
}

func (p *PulseConveyor) getFuturePulseSlotMachine(pn pulse.Number) *PulseSlotMachine {
	if psm := p.getPulseSlotMachine(pn); psm != nil {
		return psm
	}
	return p.createAndPublishPulseSlotMachine(pn, Future)
}

func (p *PulseConveyor) getAntiquePulseSlotMachine() *PulseSlotMachine {
	if psm := p.getPulseSlotMachine(0); psm != nil {
		return psm
	}
	return p.createAndPublishPulseSlotMachine(0, Antique)
}

func (p *PulseConveyor) createAndPublishPulseSlotMachine(pn pulse.Number, mode PulseSlotState) *PulseSlotMachine {
	var psm *PulseSlotMachine
	// TODO new PulseSlotMachine for the given mode
	psv, _ := p.slotMachine.TryPublish(pn, psm)
	psm = psv.(*PulseSlotMachine)
	if psm == nil {
		panic("illegal state")
	}
	return psm
}

func (p *PulseConveyor) sendSignal(fn smachine.MachineCallFunc) error {
	result := make(chan error, 1)
	p.slotMachine.ScheduleCall(func(ctx smachine.MachineCallContext) {
		defer func() {
			result <- smachine.RecoverSlotPanicWithStack("signal", recover(), nil)
			close(result)
		}()
		fn(ctx)
	}, true)
	return <-result
}

func (p *PulseConveyor) PreparePulseChange(out chan<- PreparedState) error {
	return p.sendSignal(func(ctx smachine.MachineCallContext) {
		if p.presentMachine == nil {
			// wrong - we cant produce a state for first pulse
			panic("illegal state")
		}
		ctx.BargeInNow(p.presentMachine.SlotLink(), out, p.presentMachine.preparePulseChange)
	})
}

func (p *PulseConveyor) CancelPulseChange() error {
	return p.sendSignal(func(ctx smachine.MachineCallContext) {
		if p.presentMachine == nil {
			return
		}
		ctx.BargeInNow(p.presentMachine.SlotLink(), nil, p.presentMachine.cancelPulseChange)
	})
}

func (p *PulseConveyor) CommitPulseChange(pd pulse.Data) error {
	pd.EnsurePulsarData()
	// TODO check pulse data

	return p.sendSignal(func(ctx smachine.MachineCallContext) {
		//defer func() {
		//	result <- recover()
		//	close(result)
		//}()
		//if p.presentMachine == nil {
		//	return
		//}
		// TODO allocate and swap slots
		ctx.Migrate()
	})
}
