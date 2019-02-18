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

package conveyer

import (
	"sync"

	"github.com/insolar/insolar/conveyer/queue"
	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

// EventSink allows to push events to conveyer
type EventSink interface {
	// SinkPush adds event to conveyer
	SinkPush(pulseNumber core.PulseNumber, data interface{}) error
	// SinkPushAll adds several events to conveyer
	SinkPushAll(pulseNumber core.PulseNumber, data []interface{}) error
}

// State is the states of conveyer
type State int

//go:generate stringer -type=State
const (
	Active = State(iota)
	PreparingPulse
	ShuttingDown
	Inactive
)

// Control allows to control conveyer and pulse
type Control interface {
	// PreparePulse is preparing conveyer for working with provided pulse
	PreparePulse(pulse core.Pulse) error
	// ActivatePulse is activate conveyer with prepared pulse
	ActivatePulse() error
	// GetState returns current state of conveyer
	GetState() State
	// IsOperational shows if conveyer is ready for work
	IsOperational() bool
}

// Conveyer is responsible for all pulse-dependent processing logic
type Conveyer interface {
	EventSink
	Control
}

// PulseConveyer is realization of Conveyer
type PulseConveyer struct {
	slotMap              map[core.PulseNumber]*Slot
	futurePulseData      *core.Pulse
	newFuturePulseNumber *core.PulseNumber
	futurePulseNumber    *core.PulseNumber
	presentPulseNumber   *core.PulseNumber
	lock                 sync.RWMutex
	state                State
}

// NewPulseConveyer creates new instance of PulseConveyer
func NewPulseConveyer() Conveyer {
	c := &PulseConveyer{
		slotMap: make(map[core.PulseNumber]*Slot),
		state:   Inactive,
	}
	antiqueSlot := NewSlot(Antique, AntiqueSlotPulse)
	c.slotMap[AntiqueSlotPulse] = antiqueSlot
	return c
}

// PulseState is the states of pulse inside slot
type PulseState int

const AntiqueSlotPulse = core.PulseNumber(0)

//go:generate stringer -type=PulseState
const (
	Unallocated = PulseState(iota)
	Future
	Present
	Past
	Antique
)

// Slot holds info about specific pulse and events for it
type Slot struct {
	inputQueue  queue.IQueue
	pulseState  PulseState
	pulseNumber core.PulseNumber
}

// NewSlot creates new instance of Slot
func NewSlot(pulseState PulseState, pulseNumber core.PulseNumber) *Slot {
	return &Slot{
		pulseState:  pulseState,
		inputQueue:  queue.NewMutexQueue(),
		pulseNumber: pulseNumber,
	}
}

// GetState returns current state of conveyer
func (c *PulseConveyer) GetState() State {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.state
}

// IsOperational shows if conveyer is ready for work
func (c *PulseConveyer) IsOperational() bool {
	currentState := c.GetState()
	if currentState == Active || currentState == PreparingPulse {
		return true
	}
	return false
}

func (c *PulseConveyer) unsafeGetSlot(pulseNumber core.PulseNumber) *Slot {
	slot, ok := c.slotMap[pulseNumber]
	if !ok {
		if c.futurePulseNumber == nil || pulseNumber > *c.futurePulseNumber {
			return nil
		}
		slot = c.slotMap[AntiqueSlotPulse]
	}
	return slot
}

// SinkPush adds event to conveyer
func (c *PulseConveyer) SinkPush(pulseNumber core.PulseNumber, data interface{}) error {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if !c.IsOperational() {
		return errors.New("[ SinkPush ] conveyer is not operational now")
	}
	slot := c.unsafeGetSlot(pulseNumber)
	if slot == nil {
		return errors.Errorf("[ SinkPush ] can't get slot by pulse number %d", pulseNumber)
	}
	return errors.Wrap(slot.inputQueue.SinkPush(data), "[ SinkPush ] can't push to queue")
}

// SinkPushAll adds several events to conveyer
func (c *PulseConveyer) SinkPushAll(pulseNumber core.PulseNumber, data []interface{}) error {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if !c.IsOperational() {
		return errors.New("[ SinkPushAll ] conveyer is not operational now")
	}
	slot := c.unsafeGetSlot(pulseNumber)
	if slot == nil {
		return errors.Errorf("[ SinkPushAll ] can't get slot by pulse number %d", pulseNumber)
	}
	return errors.Wrap(slot.inputQueue.SinkPushAll(data), "[ SinkPushAll ] can't push to queue")
}

// PreparePulse is preparing conveyer for working with provided pulse
// TODO: add callback param
func (c *PulseConveyer) PreparePulse(pulse core.Pulse) error {
	if !c.IsOperational() {
		return errors.New("[ PreparePulse ] conveyer is not operational now")
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.futurePulseData != nil {
		return errors.New("[ PreparePulse ] preparation was already done")
	}
	if c.futurePulseNumber == nil {
		futureSlot := NewSlot(Future, pulse.PulseNumber)
		c.slotMap[pulse.PulseNumber] = futureSlot
		c.futurePulseNumber = &pulse.PulseNumber
	}
	if *c.futurePulseNumber != pulse.PulseNumber {
		return errors.New("[ PreparePulse ] received future pulse is different from expected")
	}
	// TODO: add sending signal to slots queues

	c.futurePulseData = &pulse
	newFutureSlot := NewSlot(Unallocated, pulse.NextPulseNumber)
	c.slotMap[pulse.NextPulseNumber] = newFutureSlot
	c.newFuturePulseNumber = &pulse.NextPulseNumber
	c.state = PreparingPulse
	return nil
}

// ActivatePulse is activate conveyer with prepared pulse
func (c *PulseConveyer) ActivatePulse() error {
	if !c.IsOperational() {
		return errors.New("[ ActivatePulse ] conveyer is not operational now")
	}

	c.lock.Lock()
	defer c.lock.Unlock()

	if c.futurePulseData == nil {
		return errors.New("[ ActivatePulse ] preparation missing")
	}

	c.futurePulseData = nil

	c.presentPulseNumber = c.futurePulseNumber
	c.futurePulseNumber = c.newFuturePulseNumber
	// TODO: add sending signal to slots queues
	c.state = Active

	return nil
}
