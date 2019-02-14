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
)

//go:generate minimock -i github.com/insolar/insolar/conveyer.NonBlockingQueue -o ./ -s _mock.go
// will be change to Ivan realization
type NonBlockingQueue interface {
	SinkPush(data interface{}) bool
	SinkPushAll(data []interface{}) bool
	RemoveAll() [](interface{})
}

type EventSink interface {
	sinkPush(addr core.PulseNumber, data interface{}) bool
	sinkPushAll(addr core.PulseNumber, data []interface{}) bool
}

type Conveyer interface {
	EventSink
	readLock()
	readUnlock()
	writeLock()
	writeUnlock()
}

type PulseConveyer struct {
	PulseStorage core.PulseStorage `inject:""`
	slotMap      map[core.PulseNumber]Slot
	lock         sync.RWMutex
}

func NewPulseConveyer() *PulseConveyer {
	return &PulseConveyer{
		slotMap: make(map[core.PulseNumber]Slot),
	}
}

type PulseState int

const (
	Unassigned = PulseState(iota)
	Future
	Present
	Past
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

func (c *PulseConveyer) sinkPush(addr core.PulseNumber, data interface{}) bool {
	c.readLock()
	defer c.readUnlock()
	slot, ok := c.slotMap[addr]
	if !ok {
		// TODO: create if some cases new slot?
		return false
	}
	ok = slot.inputQueue.SinkPush(data)
	if !ok {
		return false
	}
	return true
}

func (c *PulseConveyer) sinkPushAll(addr core.PulseNumber, data []interface{}) bool {
	c.readLock()
	defer c.readUnlock()
	slot, ok := c.slotMap[addr]
	if !ok {
		// TODO: create if some cases new slot?
		return false
	}
	for _, d := range data {
		ok = slot.inputQueue.SinkPush(d)
		if !ok {
			return false
		}
	}
	return true
}

// Start creates Present and Future Slots.
func (c *PulseConveyer) Start(ctx context.Context) error {
	pulse, err := c.PulseStorage.Current(ctx)
	if err != nil {
		return err
	}
	presentSlot := NewSlot(Present)
	c.slotMap[pulse.PulseNumber] = presentSlot
	futureSlot := NewSlot(Future)
	c.slotMap[pulse.NextPulseNumber] = futureSlot
	return nil
}

// Stop stops PulseConveyer.
func (c *PulseConveyer) Stop(ctx context.Context) error {
	// TODO: process every slot? just get out and delete all?
	return nil
}
