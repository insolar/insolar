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
	"time"

	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/conveyor/sworker"
	"github.com/insolar/insolar/conveyor/tools"
	"github.com/insolar/insolar/pulse"
)

type InputEvent = interface{}
type PulseEventFactoryFunc = func(pulse.Number, InputEvent) smachine.CreateFunc

type EventInputer interface {
	AddInput(ctx context.Context, pn pulse.Number, event InputEvent) error
}

type PreparedState = struct{}
type PreparePulseChangeChannel = chan<- PreparedState

type PulseChanger interface {
	PreparePulseChange(out PreparePulseChangeChannel) error
	CancelPulseChange() error
	CommitPulseChange(pd pulse.Data) error
}

type PulseConveyorConfig struct {
	ConveyorMachineConfig             smachine.SlotMachineConfig
	SlotMachineConfig                 smachine.SlotMachineConfig
	EventlessSleep                    time.Duration
	MinCachePulseAge, MaxPastPulseAge uint32
}

func NewPulseConveyor(
	ctx context.Context,
	config PulseConveyorConfig,
	factoryFn PulseEventFactoryFunc,
	registry injector.DependencyRegistry,
) *PulseConveyor {

	r := &PulseConveyor{
		workerCtx: ctx,
		slotConfig: PulseSlotConfig{
			config: config.SlotMachineConfig,
		},
		factoryFn:      factoryFn,
		eventlessSleep: config.EventlessSleep,
	}
	r.slotMachine = smachine.NewSlotMachine(config.ConveyorMachineConfig,
		r.internalSignal.NextBroadcast,
		combineCallbacks(r.externalSignal.NextBroadcast, r.internalSignal.NextBroadcast),
		registry)

	r.slotConfig.eventCallback = r.internalSignal.NextBroadcast
	r.slotConfig.parentRegistry = r.slotMachine

	// shared SlotId sequence
	r.slotConfig.config.SlotIdGenerateFn = r.slotMachine.CopyConfig().SlotIdGenerateFn

	r.pdm.Init(config.MinCachePulseAge, config.MaxPastPulseAge, 1)

	return r
}

type PulseConveyor struct {
	// immutable, provided, set at construction
	slotConfig     PulseSlotConfig
	eventlessSleep time.Duration
	factoryFn      PulseEventFactoryFunc
	workerCtx      context.Context

	// immutable, set at construction
	externalSignal tools.VersionedSignal
	internalSignal tools.VersionedSignal

	slotMachine   *smachine.SlotMachine
	machineWorker smachine.AttachableSlotWorker

	pdm PulseDataManager

	// mutable, set under SlotMachine synchronization
	presentMachine *PulseSlotMachine
	unpublishPulse pulse.Number
}

func (p *PulseConveyor) AddDependency(v interface{}) {
	p.slotMachine.AddDependency(v)
}

func (p *PulseConveyor) FindDependency(id string) (interface{}, bool) {
	return p.slotMachine.FindDependency(id)
}

func (p *PulseConveyor) PutDependency(id string, v interface{}) {
	p.slotMachine.PutDependency(id, v)
}

func (p *PulseConveyor) TryPutDependency(id string, v interface{}) bool {
	return p.slotMachine.TryPutDependency(id, v)
}

func (p *PulseConveyor) GetPublishedGlobalAlias(key interface{}) smachine.SlotLink {
	return p.slotMachine.GetPublishedGlobalAlias(key)
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
	var createDefaults smachine.CreateDefaultValues

	switch {
	case createFn == nil:
		return fmt.Errorf("unrecognized event: pn=%v event=%v", targetPN, event)

	case pulseState == Future:
		// event for future needs special handling - it must wait until the pulse will actually arrive
		pulseSlotMachine.innerMachine.AddNew(ctx,
			newFutureEventSM(targetPN, &pulseSlotMachine.pulseSlot, createFn), createDefaults)
		return nil

	case pulseState == Antique:
		// Antique events have individual pulse slots, while being executed in a single SlotMachine
		if cps, ok := p.pdm.getCachedPulseSlot(targetPN); ok {
			createDefaults.PutOverride(injector.GetDefaultInjectionId(cps), cps)
			break // add SM
		}

		if !p.pdm.IsRecentPastRange(pn) {
			// for non-recent past HasPulseData() can be incorrect / incomplete
			// we must use a longer procedure to get PulseData and utilize SM for it
			pulseSlotMachine.innerMachine.AddNew(ctx,
				newAntiqueEventSM(targetPN, &pulseSlotMachine.pulseSlot, createFn), createDefaults)
			return nil
		}
		fallthrough

	case !p.pdm.TouchPulseData(targetPN): // make sure - for PAST and PRESENT we must always have the data ...
		return fmt.Errorf("unknown data for pulse: pn=%v event=%v", targetPN, event)
	}

	if _, ok := pulseSlotMachine.innerMachine.AddNewByFunc(ctx, createFn, createDefaults); !ok {
		return fmt.Errorf("ignored event: pn=%v event=%v", targetPN, event)
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
		// present slot must be present
		panic("illegal state")
	case !pn.IsTimePulse():
		return nil, 0, 0, fmt.Errorf("pulse number is invalid: pn=%v", pn)
	case pn < presentPN:
		// this can be either be a past/antique slot, or a part of the present range
		if psm := p.getPulseSlotMachine(pn); psm != nil {
			return psm, pn, Past, nil
		}

		// check if the pulse is within PRESENT range (as it may include some skipped pulses)
		if psm := p.getPulseSlotMachine(presentPN); psm == nil {
			// present slot must be present
			panic("illegal state")
		} else {
			switch ps, ok := psm.pulseSlot._isAcceptedPresent(presentPN, pn); {
			case ps == Past:
				// pulse has changed - then we handle the packet as usual
				break
			case !ok:
				return nil, 0, 0, fmt.Errorf("pulse number is not allowed: pn=%v", pn)
			case ps != Present:
				panic("illegal state")
			}
		}

		if !p.pdm.isAllowedPastSpan(presentPN, pn) {
			return nil, 0, 0, fmt.Errorf("pulse number is too far in past: pn=%v, present=%v", pn, presentPN)
		}
		return p.getAntiquePulseSlotMachine(), pn, Antique, nil
	case pn < futurePN:
		return nil, 0, 0, fmt.Errorf("pulse number is unexpected: pn=%v", pn)
	default: // pn >= futurePN
		if !p.pdm.isAllowedFutureSpan(presentPN, futurePN, pn) {
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
	case prevDelta > math.MaxUint16:
		prevDelta = math.MaxUint16
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
	psm.activate(p.workerCtx, p.slotMachine.AddNew)
	psm.setPulseForUnpublish(p.slotMachine, pn)

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

func (p *PulseConveyor) PreparePulseChange(out PreparePulseChangeChannel) error {
	return p.sendSignal(func(ctx smachine.MachineCallContext) {
		if p.presentMachine == nil {
			// wrong - first pulse can only be committed but not prepared
			panic("illegal state")
		}
		p.pdm.setPreparingPulse(out)
		if !ctx.BargeInNow(p.presentMachine.SlotLink(), out, p.presentMachine.preparePulseChange) {
			panic("present slot is busy")
		}
	})
}

func (p *PulseConveyor) CancelPulseChange() error {
	return p.sendSignal(func(ctx smachine.MachineCallContext) {
		if p.presentMachine == nil {
			// wrong - first pulse can only be committed but not prepared
			panic("illegal state")
		}
		p.pdm.unsetPreparingPulse()
		if !ctx.BargeInNow(p.presentMachine.SlotLink(), nil, p.presentMachine.cancelPulseChange) {
			panic("present slot is busy")
		}
	})
}

func (p *PulseConveyor) CommitPulseChange(pr pulse.Range) error {
	return p.sendSignal(func(ctx smachine.MachineCallContext) {
		p.pdm.unsetPreparingPulse()
		ctx.Migrate(func() {
			p._promotePulseSlots(ctx, pr)
		})
	})
}

func (p *PulseConveyor) _promotePulseSlots(ctx smachine.MachineCallContext, pr pulse.Range) {
	pd := pr.RightBoundData()
	pd.EnsurePulsarData()
	prevPresentPN, prevFuturePN := p.pdm.GetPresentPulse()

	if p.presentMachine == nil {
		switch {
		case prevFuturePN != uninitializedFuture:
			panic("illegal state")
		case p.getPulseSlotMachine(prevPresentPN) != nil:
			panic("illegal state")
		}
	} else {
		switch {
		case p.getPulseSlotMachine(prevPresentPN) != p.presentMachine:
			panic("illegal state")
		case pr.LeftBoundNumber() != prevFuturePN:
			panic("illegal state")
		case prevPresentPN.Next(pr.LeftPrevDelta()) != pr.LeftBoundNumber():
			panic("illegal state")
		}
		p.presentMachine.setPast()
	}
	pr.EnumNonArticulatedData(func(data pulse.Data) bool {
		p.pdm.putPulseData(data) // add to the recent cache
		return false
	})

	if p.unpublishPulse.IsTimePulse() {
		// we know what we do - right!?
		p.slotMachine.TryUnsafeUnpublish(p.unpublishPulse)
		p.unpublishPulse = pulse.Unknown
	}

	prevFuture := p.getPulseSlotMachine(prevFuturePN)

	republishPresent := false
	activatePresent := false

	if prevFuture != nil {
		prevFuture.setPresent(pr)
		p.presentMachine = prevFuture

		if prevFuturePN != pd.PulseNumber {
			// new pulse is different than expected at the previous cycle, so we have to remove the pulse number alias
			// to avoids unnecessary synchronization - the previous alias will be unpublished on commit of a next pulse
			p.unpublishPulse = prevFuturePN
			republishPresent = true
		}
	} else {
		psm := p.newPulseSlotMachine()
		psm.setPresent(pr)
		p.presentMachine = psm
		activatePresent = true
	}

	if republishPresent || activatePresent {
		p.presentMachine.setPulseForUnpublish(p.slotMachine, pd.PulseNumber)

		if _, ok := p.slotMachine.TryPublish(pd.PulseNumber, p.presentMachine); !ok {
			panic("illegal state")
		}
	}

	if activatePresent {
		p.presentMachine.activate(p.workerCtx, ctx.AddNew)
	}
	p.pdm.setPresentPulse(pd) // reroutes incoming events
}

func (p *PulseConveyor) StopNoWait() {
	p.slotMachine.Stop()
}

func (p *PulseConveyor) StartWorker(emergencyStop <-chan struct{}, completedFn func()) {

	if p.machineWorker != nil {
		panic("illegal state")
	}
	p.machineWorker = sworker.NewAttachableSimpleSlotWorker()
	go p.runWorker(emergencyStop, completedFn)
}

func (p *PulseConveyor) runWorker(emergencyStop <-chan struct{}, completedFn func()) {
	if emergencyStop != nil {
		go func() {
			select {
			case <-emergencyStop:
				p.slotMachine.Stop()
				p.externalSignal.NextBroadcast()
				return
			}
		}()
	}

	if completedFn != nil {
		defer completedFn()
	}

	for {
		var (
			repeatNow    bool
			nextPollTime time.Time
		)
		eventMark := p.internalSignal.Mark()
		p.machineWorker.AttachTo(p.slotMachine, p.externalSignal.Mark(), math.MaxUint32, func(worker smachine.AttachedSlotWorker) {
			repeatNow, nextPollTime = p.slotMachine.ScanOnce(smachine.ScanDefault, worker)
		})

		select {
		case <-emergencyStop:
			return
		default:
			// pass
		}

		if !p.slotMachine.IsActive() {
			break
		}

		if repeatNow || eventMark.HasSignal() {
			continue
		}

		select {
		case <-emergencyStop:
			return
		case <-eventMark.Channel():
		case <-func() <-chan time.Time {
			switch {
			case !nextPollTime.IsZero():
				return time.After(time.Until(nextPollTime))
			case p.eventlessSleep > 0 && p.eventlessSleep < math.MaxInt64:
				return time.After(p.eventlessSleep)
			}
			return nil
		}():
		}
	}

	p.slotMachine.RunToStop(p.machineWorker, tools.NewNeverSignal())
	p.presentMachine = nil
}
