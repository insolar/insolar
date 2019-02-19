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

// Control allows to control conveyor and pulse
type Control interface {
	// PreparePulse is preparing conveyor for working with provided pulse
	PreparePulse(pulse core.Pulse) error
	// ActivatePulse is activate conveyor with prepared pulse
	ActivatePulse() error
	// GetState returns current state of conveyor
	GetState() State
	// IsOperational shows if conveyor is ready for work
	IsOperational() bool
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

func (c *PulseConveyor) unsafeGetSlot(pulseNumber core.PulseNumber) *Slot {
	slot, ok := c.slotMap[pulseNumber]
	if !ok {
		if c.futurePulseNumber == nil || pulseNumber > *c.futurePulseNumber {
			return nil
		}
		slot = c.slotMap[AntiqueSlotPulse]
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
// TODO: add callback param
func (c *PulseConveyor) PreparePulse(pulse core.Pulse) error {
	if !c.IsOperational() {
		return errors.New("[ PreparePulse ] conveyor is not operational now")
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

// ActivatePulse activates conveyor with prepared pulse
func (c *PulseConveyor) ActivatePulse() error {
	if !c.IsOperational() {
		return errors.New("[ ActivatePulse ] conveyor is not operational now")
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
