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

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

//go:generate minimock -i github.com/insolar/insolar/conveyer.NonBlockingQueue -o ./ -s _mock.go
// will be change to Ivan realization
type NonBlockingQueue interface {
	SinkPush(data interface{}) bool
	SinkPushAll(data []interface{}) bool
	RemoveAll() [](interface{})
}

type EventSink interface {
	SinkPush(pulseNumber core.PulseNumber, data interface{}) bool
	SinkPushAll(pulseNumber core.PulseNumber, data []interface{}) bool
}

type State int

//go:generate stringer -type=State
const (
	Active = State(iota)
	PreparingPulse
	ShuttingDown
	Inactive
)

type Control interface {
	PreparePulse(pulse core.Pulse) error
	ActivatePulse() error
	GetState() State
	IsOperational() bool
}

type Conveyer interface {
	EventSink
	Control
}

type PulseConveyer struct {
	slotMap              map[core.PulseNumber]*Slot
	futurePulseData      *core.Pulse
	newFuturePulseNumber *core.PulseNumber
	futurePulseNumber    *core.PulseNumber
	presentPulseNumber   *core.PulseNumber
	lock                 sync.RWMutex
	state                State
}

func NewPulseConveyer() *PulseConveyer {
	c := &PulseConveyer{
		slotMap: make(map[core.PulseNumber]*Slot),
		state:   Inactive,
	}
	antiqueSlot := NewSlot(Antique, AntiqueSlotPulse)
	c.slotMap[AntiqueSlotPulse] = antiqueSlot
	return c
}

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

type Slot struct {
	inputQueue  NonBlockingQueue
	pulseState  PulseState
	pulseNumber core.PulseNumber
}

func NewSlot(pulseState PulseState, pulseNumber core.PulseNumber) *Slot {
	return &Slot{
		pulseState:  pulseState,
		inputQueue:  &Queue{},
		pulseNumber: pulseNumber,
	}
}

// GetState returns current state of Conveyer
func (c *PulseConveyer) GetState() State {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.state
}

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

func (c *PulseConveyer) SinkPush(pulseNumber core.PulseNumber, data interface{}) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	slot := c.unsafeGetSlot(pulseNumber)
	if slot == nil {
		return false
	}
	return slot.inputQueue.SinkPush(data)
}

func (c *PulseConveyer) SinkPushAll(pulseNumber core.PulseNumber, data []interface{}) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	slot := c.unsafeGetSlot(pulseNumber)
	if slot == nil {
		return false
	}
	return slot.inputQueue.SinkPushAll(data)
}

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
	return nil
}

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

	return nil
}
