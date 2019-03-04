/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package conveyor

import (
	"sync"

	"github.com/insolar/insolar/conveyor/queue"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

// EventSink allows to push events to conveyor
type EventSink interface {
	// SinkPush adds event to conveyor
	SinkPush(pulseNumber core.PulseNumber, data interface{}) error
	// SinkPushAll adds several events to conveyor
	SinkPushAll(pulseNumber core.PulseNumber, data []interface{}) error
}

// State is the states of conveyor
type State int

//go:generate stringer -type=State
const (
	Active = State(iota)
	PreparingPulse
	ShuttingDown
	Inactive
)

const (
	PendingPulseSignal  = 1
	ActivatePulseSignal = 2
)

// Control allows to control conveyor and pulse
type Control interface {
	// PreparePulse is preparing conveyor for working with provided pulse
	PreparePulse(pulse core.Pulse, callback queue.SyncDone) error
	// ActivatePulse is activate conveyor with prepared pulse
	ActivatePulse() error
	// GetState returns current state of conveyor
	GetState() State
	// IsOperational shows if conveyor is ready for work
	IsOperational() bool
	// InitiateShutdown shutting conveyor down and cancels tasks in adapters if force param set
	InitiateShutdown(force bool)
}

// Conveyor is responsible for all pulse-dependent processing logic
type Conveyor interface {
	EventSink
	Control
}

// PulseConveyor is realization of Conveyor
type PulseConveyor struct {
	slotMap              map[core.PulseNumber]*Slot
	futurePulseData      *core.Pulse
	newFuturePulseNumber *core.PulseNumber
	futurePulseNumber    *core.PulseNumber
	presentPulseNumber   *core.PulseNumber
	lock                 sync.RWMutex
	state                State
}

// NewPulseConveyor creates new instance of PulseConveyor
func NewPulseConveyor() Conveyor {
	c := &PulseConveyor{
		slotMap: make(map[core.PulseNumber]*Slot),
		state:   Inactive,
	}
	// antiqueSlot is slot for all pulses from past if conveyor dont have specific PastSlot for such pulse
	antiqueSlot := NewSlot(Antique, core.AntiquePulseNumber)
	c.slotMap[core.AntiquePulseNumber] = antiqueSlot
	return c
}

// GetState returns current state of conveyor
func (c *PulseConveyor) GetState() State {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.state
}

// IsOperational shows if conveyor is ready for work
func (c *PulseConveyor) IsOperational() bool {
	currentState := c.GetState()
	if currentState == Active || currentState == PreparingPulse {
		return true
	}
	return false
}

func (c *PulseConveyor) InitiateShutdown(force bool) {
	if force {
		// cancel all tasks in adapters
	}
}

func (c *PulseConveyor) unsafeGetSlot(pulseNumber core.PulseNumber) *Slot {
	slot, ok := c.slotMap[pulseNumber]
	if !ok {
		if c.futurePulseNumber == nil || pulseNumber > *c.futurePulseNumber {
			return nil
		}
		slot = c.slotMap[core.AntiquePulseNumber]
	}
	return slot
}

// SinkPush adds event to conveyor
func (c *PulseConveyor) SinkPush(pulseNumber core.PulseNumber, data interface{}) error {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if !c.IsOperational() {
		return errors.New("[ SinkPush ] conveyor is not operational now")
	}
	slot := c.unsafeGetSlot(pulseNumber)
	if slot == nil {
		return errors.Errorf("[ SinkPush ] can't get slot by pulse number %d", pulseNumber)
	}
	err := slot.inputQueue.SinkPush(data)
	return errors.Wrap(err, "[ SinkPush ] can't push to queue")
}

// SinkPushAll adds several events to conveyor
func (c *PulseConveyor) SinkPushAll(pulseNumber core.PulseNumber, data []interface{}) error {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if !c.IsOperational() {
		return errors.New("[ SinkPushAll ] conveyor is not operational now")
	}
	slot := c.unsafeGetSlot(pulseNumber)
	if slot == nil {
		return errors.Errorf("[ SinkPushAll ] can't get slot by pulse number %d", pulseNumber)
	}
	err := slot.inputQueue.SinkPushAll(data)
	return errors.Wrap(err, "[ SinkPushAll ] can't push to queue")
}

// PreparePulse is preparing conveyor for working with provided pulse
func (c *PulseConveyor) PreparePulse(pulse core.Pulse, callback queue.SyncDone) error {
	if !c.IsOperational() {
		return errors.New("[ PreparePulse ] conveyor is not operational now")
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.futurePulseData != nil {
		return errors.New("[ PreparePulse ] preparation was already done")
	}
	if c.futurePulseNumber == nil {
		c.slotMap[pulse.PulseNumber] = NewSlot(Future, pulse.PulseNumber)
		c.futurePulseNumber = &pulse.PulseNumber
	}
	if *c.futurePulseNumber != pulse.PulseNumber {
		return errors.New("[ PreparePulse ] received future pulse is different from expected")
	}

	futureSlot := c.slotMap[*c.futurePulseNumber]
	err := futureSlot.inputQueue.PushSignal(PendingPulseSignal, callback)
	if err != nil {
		log.Panicf("[ PreparePulse ] can't send signal to future slot (for pulse %d), error - %s", c.futurePulseNumber, err)
	}

	if c.presentPulseNumber != nil {
		presentSlot := c.slotMap[*c.presentPulseNumber]
		err := presentSlot.inputQueue.PushSignal(PendingPulseSignal, callback)
		if err != nil {
			log.Panicf("[ PreparePulse ] can't send signal to present slot (for pulse %d), error - %s", c.presentPulseNumber, err)
		}
	}

	c.futurePulseData = &pulse
	newFutureSlot := NewSlot(Unallocated, pulse.NextPulseNumber)
	c.slotMap[pulse.NextPulseNumber] = newFutureSlot
	c.newFuturePulseNumber = &pulse.NextPulseNumber
	c.state = PreparingPulse
	return nil
}

// PulseWithCallback contains info about new pulse and callback func
type PulseWithCallback interface {
	GetCallback() queue.SyncDone
	GetPulse() core.Pulse
	Done()
}

type pulseWithCallback struct {
	callback queue.SyncDone
	pulse    core.Pulse
}

// PulseWithCallback creates new instance of pulseWithCallback
func NewPulseWithCallback(callback queue.SyncDone, pulse core.Pulse) PulseWithCallback {
	return &pulseWithCallback{
		callback: callback,
		pulse:    pulse,
	}
}

// GetCallback returns callback
func (p *pulseWithCallback) GetCallback() queue.SyncDone {
	return p.callback
}

// GetCallback returns callback
func (p *pulseWithCallback) GetPulse() core.Pulse {
	return p.pulse
}

// Done calls .Done() on func in callback param
func (p *pulseWithCallback) Done() {
	p.callback.Done()
}

// ActivatePulse activates conveyor with prepared pulse
func (c *PulseConveyor) ActivatePulse() error {
	if !c.IsOperational() {
		return errors.New("[ ActivatePulse ] conveyor is not operational now")
	}

	c.lock.Lock()

	if c.futurePulseData == nil {
		c.lock.Unlock()
		return errors.New("[ ActivatePulse ] preparation missing")
	}

	c.presentPulseNumber = c.futurePulseNumber
	c.futurePulseNumber = c.newFuturePulseNumber

	wg := sync.WaitGroup{}

	futureSlot := c.slotMap[*c.futurePulseNumber]
	callback := NewPulseWithCallback(&wg, *c.futurePulseData)
	wg.Add(1)
	err := futureSlot.inputQueue.PushSignal(ActivatePulseSignal, callback)
	if err != nil {
		c.lock.Unlock()
		log.Panicf("[ ActivatePulse ] can't send signal to future slot (for pulse %d), error - %s", c.futurePulseNumber, err)
	}

	presentSlot := c.slotMap[*c.presentPulseNumber]
	wg.Add(1)
	err = presentSlot.inputQueue.PushSignal(ActivatePulseSignal, &wg)
	if err != nil {
		c.lock.Unlock()
		log.Panicf("[ ActivatePulse ] can't send signal to present slot (for pulse %d), error - %s", c.presentPulseNumber, err)
	}

	c.futurePulseData = nil
	c.state = Active
	c.lock.Unlock()
	wg.Wait()

	return nil
}

func (c *PulseConveyor) getSlotConfiguration(state SlotState) HandlersConfiguration { // nolint: unused
	return HandlersConfiguration{state: state}
}
