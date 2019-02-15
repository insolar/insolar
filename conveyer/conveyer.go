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
	"context"
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
	SinkPush(addr core.PulseNumber, data interface{}) bool
	SinkPushAll(addr core.PulseNumber, data []interface{}) bool
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
	GetState() State
	IsOperational() bool
}

type Conveyer interface {
	EventSink
	Control
}

type PulseConveyer struct {
	PulseStorage core.PulseStorage `inject:""`
	slotMap      map[core.PulseNumber]Slot
	lock         sync.RWMutex
	state        State
}

func NewPulseConveyer() *PulseConveyer {
	return &PulseConveyer{
		slotMap: make(map[core.PulseNumber]Slot),
		state:   Inactive,
	}
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
	inputQueue NonBlockingQueue
	pulseState PulseState
}

func NewSlot(pulseState PulseState) Slot {
	return Slot{
		pulseState: pulseState,
		inputQueue: &Queue{},
	}
}

func (c *PulseConveyer) readLock() {
	c.lock.RLock()
}

func (c *PulseConveyer) readUnlock() {
	c.lock.RUnlock()
}

func (c *PulseConveyer) writeLock() {
	c.lock.Lock()
}

func (c *PulseConveyer) writeUnlock() {
	c.lock.Unlock()
}

func (c *PulseConveyer) GetState() State {
	c.readLock()
	defer c.readUnlock()
	return c.state
}

func (c *PulseConveyer) IsOperational() bool {
	currentState := c.GetState()
	if currentState == Active || currentState == PreparingPulse {
		return true
	}
	return false
}

func (c *PulseConveyer) getSlot(addr core.PulseNumber) (Slot, error) {
	slot, ok := c.slotMap[addr]
	if !ok {
		ctx := context.Background()
		currentPulse, err := c.PulseStorage.Current(ctx)
		if err != nil {
			return Slot{}, err
		}
		if addr >= currentPulse.PulseNumber {
			return Slot{}, errors.New("unknown pulse")
		}
		slot = c.slotMap[AntiqueSlotPulse]
	}
	return slot, nil
}

func (c *PulseConveyer) SinkPush(addr core.PulseNumber, data interface{}) bool {
	c.readLock()
	defer c.readUnlock()
	slot, err := c.getSlot(addr)
	if err != nil {
		return false
	}
	ok := slot.inputQueue.SinkPush(data)
	return ok
}

func (c *PulseConveyer) SinkPushAll(addr core.PulseNumber, data []interface{}) bool {
	c.readLock()
	defer c.readUnlock()
	slot, err := c.getSlot(addr)
	if err != nil {
		return false
	}
	ok := slot.inputQueue.SinkPushAll(data)
	return ok
}

// Start creates Present and Future Slots.
func (c *PulseConveyer) Start(ctx context.Context) error {
	pulse, err := c.PulseStorage.Current(ctx)
	if err != nil {
		return err
	}

	c.writeLock()
	defer c.writeUnlock()

	presentSlot := NewSlot(Present)
	c.slotMap[pulse.PulseNumber] = presentSlot

	futureSlot := NewSlot(Future)
	c.slotMap[pulse.NextPulseNumber] = futureSlot

	antiqueSlot := NewSlot(Antique)
	c.slotMap[AntiqueSlotPulse] = antiqueSlot

	c.state = Active
	return nil
}
