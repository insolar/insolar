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

package conveyor

import (
	"context"
	"fmt"
	"math"

	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
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
	pdm            PulseDataManager

	// mutable, set under SlotMachine synchronization
	presentMachine *PulseSlotMachine
	unpublishPulse pulse.Number
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
		pulseSlotMachine.innerMachine.AddNew(ctx, smachine.NoLink(),
			newFutureEventSM(targetPN, &pulseSlotMachine.pulseSlot, createFn))
		break
	case pulseState == Antique:
		if !p.pdm.IsRecentPastRange(pn) {
			// for non-recent past HasPulseData() can be incorrect / incomplete
			// we must use a longer procedure to get PulseData and utilize SM for it
			pulseSlotMachine.innerMachine.AddNew(ctx, smachine.NoLink(),
				newAntiqueEventSM(targetPN, &pulseSlotMachine.pulseSlot, createFn))
			break
		}
		fallthrough
	default:
		if !p.pdm.HasPulseData(targetPN) {
			return fmt.Errorf("unknown data for pulse : pn=%v event=%v", targetPN, event)
		}
		if _, ok := pulseSlotMachine.innerMachine.AddNewByFunc(ctx, smachine.NoLink(), createFn); !ok {
			return fmt.Errorf("ignored event: pn=%v event=%v", targetPN, event)
		}
	}
	return nil
}

func (p *PulseConveyor) mapToPulseSlotMachine(pn pulse.Number) (*PulseSlotMachine, pulse.Number, PulseSlotState, error) {
	presentPN, futurePN := p.pdm.GetPresentPulse()

	switch {
	case presentPN.IsUnknown():
		// when no present pulse - all pulses go to future
		return p.getFuturePulseSlotMachine(presentPN, futurePN), pn, Future, nil
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
		if !pn.IsTimePulse() {
			return nil, 0, 0, fmt.Errorf("pulse number is invalid: pn=%v", pn)
		}
		if !p.pdm.IsAllowedPastSpan(presentPN, pn) {
			return nil, 0, 0, fmt.Errorf("pulse number is too far in past: pn=%v, present=%v", pn, presentPN)
		}
		return p.getAntiquePulseSlotMachine(), pn, Antique, nil
	case pn < futurePN:
		return nil, 0, 0, fmt.Errorf("pulse number is unexpected: pn=%v", pn)
	default: // pn >= futurePN
		if !pn.IsTimePulse() {
			return nil, 0, 0, fmt.Errorf("pulse number is invalid: pn=%v", pn)
		}
		if !p.pdm.IsAllowedFutureSpan(futurePN, pn) {
			return nil, 0, 0, fmt.Errorf("pulse number is too far in future: pn=%v, expected=%v", pn, futurePN)
		}
		return p.getFuturePulseSlotMachine(presentPN, futurePN), pn, Future, nil
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

func (p *PulseConveyor) getFuturePulseSlotMachine(presentPN, futurePN pulse.Number) *PulseSlotMachine {
	if psm := p.getPulseSlotMachine(futurePN); psm != nil {
		return psm
	}
	psm := p.newPulseSlotMachine()

	prevDelta := futurePN - presentPN
	switch {
	case presentPN.IsUnknown():
		prevDelta = 0
	case prevDelta >= math.MaxUint16:
		prevDelta = math.MaxUint16 - 1
	}

	psm.setFuture(pulse.NewExpectedPulsarData(futurePN, uint16(prevDelta)))
	return p._publishPulseSlotMachine(futurePN, psm)
}

func (p *PulseConveyor) getAntiquePulseSlotMachine() *PulseSlotMachine {
	if psm := p.getPulseSlotMachine(0); psm != nil {
		return psm
	}
	psm := p.newPulseSlotMachine()
	psm.setAntique()
	return p._publishPulseSlotMachine(0, psm)
}

func (p *PulseConveyor) _publishPulseSlotMachine(pn pulse.Number, psm *PulseSlotMachine) *PulseSlotMachine {
	if psv, ok := p.slotMachine.TryPublish(pn, psm); !ok {
		psm = psv.(*PulseSlotMachine)
		if psm == nil {
			panic("illegal state")
		}
		return psm
	}
	psm.activate(p.workerCtx, &p.slotMachine)
	psm.setPulseForUnpublish(&p.slotMachine, pn)

	return psm
}

func (p *PulseConveyor) newPulseSlotMachine() *PulseSlotMachine {
	return NewPulseSlotMachine(p.slotConfig, &p.pdm)
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

type PreparePulseChangeChannel = chan<- PreparedState

func (p *PulseConveyor) PreparePulseChange(out PreparePulseChangeChannel) error {
	return p.sendSignal(func(ctx smachine.MachineCallContext) {
		if p.presentMachine == nil {
			// wrong - first pulse can only be committed
			panic("illegal state")
		}
		if !ctx.BargeInNow(p.presentMachine.SlotLink(), out, p.presentMachine.preparePulseChange) {
			panic("present slot is busy")
		}
	})
}

func (p *PulseConveyor) CancelPulseChange() error {
	return p.sendSignal(func(ctx smachine.MachineCallContext) {
		if p.presentMachine == nil {
			// wrong - first pulse can only be committed
			panic("illegal state")
		}
		if !ctx.BargeInNow(p.presentMachine.SlotLink(), nil, p.presentMachine.cancelPulseChange) {
			panic("present slot is busy")
		}
	})
}

func (p *PulseConveyor) CommitPulseChange(pd pulse.Data) error {
	return p.sendSignal(func(ctx smachine.MachineCallContext) {
		p._promotePulseSlots(ctx, pd)
		ctx.Migrate()
	})
}

func (p *PulseConveyor) _promotePulseSlots(ctx smachine.MachineCallContext, pd pulse.Data) {
	pd.EnsurePulsarData()
	prevPresentPN, prevFuturePN := p.pdm.GetPresentPulse()

	if p.presentMachine == nil {
		if prevFuturePN != uninitializedFuture {
			panic("illegal state")
		}
		if p.getPulseSlotMachine(prevPresentPN) != nil {
			panic("illegal state")
		}
	} else {
		if p.getPulseSlotMachine(prevPresentPN) != p.presentMachine {
			panic("illegal state")
		}
		p.presentMachine.setPast()
	}

	if p.unpublishPulse.IsTimePulse() {
		// we know what we do - right!?
		p.slotMachine.TryUnsafeUnpublish(pd.PulseNumber)

		p.unpublishPulse = pulse.Unknown
	}

	prevFuture := p.getPulseSlotMachine(prevFuturePN)

	republishPresent := false
	activatePresent := false

	if prevFuture != nil {
		prevFuture.setPresent(pd)
		p.presentMachine = prevFuture

		if prevFuturePN != pd.PulseNumber {
			// to avoid unnecessary synchronization the alias will be unpublished on commit of a next pulse
			p.unpublishPulse = prevFuturePN
			republishPresent = true
		}
	} else {
		psm := p.newPulseSlotMachine()
		psm.setPresent(pd)
		p.presentMachine = psm
		activatePresent = true
	}

	if republishPresent || activatePresent {
		p.presentMachine.setPulseForUnpublish(&p.slotMachine, pd.PulseNumber)

		if _, ok := p.slotMachine.TryPublish(pd.PulseNumber, p.presentMachine); !ok {
			panic("illegal state")
		}
	}

	if activatePresent {
		p.presentMachine.activateWithCtx(p.workerCtx, ctx)
	}
	p.pdm.setPresentPulse(pd) // reroutes incoming events
}
