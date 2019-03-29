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
	"encoding/hex"
	"sync"

	"github.com/insolar/insolar/conveyor/queue"
	"github.com/insolar/insolar/conveyor/slot"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

// PulseConveyor is realization of Conveyor
type PulseConveyor struct {
	slotMap            map[insolar.PulseNumber]TaskPusher
	futurePulseData    *insolar.Pulse
	futurePulseNumber  *insolar.PulseNumber
	presentPulseNumber *insolar.PulseNumber
	lock               sync.RWMutex
	state              insolar.ConveyorState
}

// NewPulseConveyor creates new instance of PulseConveyor
func NewPulseConveyor() (insolar.Conveyor, error) {
	c := &PulseConveyor{
		slotMap: make(map[insolar.PulseNumber]TaskPusher),
		state:   insolar.ConveyorInactive,
	}
	// antiqueSlot is slot for all pulses from past if conveyor dont have specific PastSlot for such pulse
	antiqueSlot := slot.NewWorkingSlot(slot.Antique, insolar.AntiquePulseNumber, c.removeSlot)

	c.slotMap[insolar.AntiquePulseNumber] = antiqueSlot
	return c, nil
}

func (c *PulseConveyor) removeSlot(number insolar.PulseNumber) {
	c.lock.Lock()
	defer c.lock.Unlock()

	err := c.slotMap[number].PushSignal(slot.CancelSignal, nil)
	if err != nil {
		panic("[ removeSlot ] Can't PushSignal CancelSignal: " + err.Error())
	}
	delete(c.slotMap, number)
}

// GetState returns current state of conveyor
func (c *PulseConveyor) GetState() insolar.ConveyorState {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.unsafeGetState()
}

// IsOperational shows if conveyor is ready for work
func (c *PulseConveyor) IsOperational() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	return c.unsafeIsOperational()
}

func (c *PulseConveyor) unsafeIsOperational() bool {
	currentState := c.unsafeGetState()
	return currentState == insolar.ConveyorActive || currentState == insolar.ConveyorPreparingPulse
}

func (c *PulseConveyor) unsafeGetState() insolar.ConveyorState {
	return c.state
}

func (c *PulseConveyor) unsafeGetSlot(pulseNumber insolar.PulseNumber) TaskPusher {
	slot, ok := c.slotMap[pulseNumber]
	if !ok {
		if c.futurePulseNumber == nil || pulseNumber > *c.futurePulseNumber {
			return nil
		}
		slot = c.slotMap[insolar.AntiquePulseNumber]
	}
	return slot
}

// InitiateShutdown starts shutdown process
func (c *PulseConveyor) InitiateShutdown(force bool) {
	c.lock.Lock()
	c.state = insolar.ConveyorShuttingDown
	c.lock.Unlock()
	if force { // nolint
		// TODO: cancel all tasks in adapters
	}
}

// SinkPush adds event to conveyor
func (c *PulseConveyor) SinkPush(pulseNumber insolar.PulseNumber, data interface{}) error {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if !c.unsafeIsOperational() {
		return errors.New("[ SinkPush ] conveyor is not operational now")
	}
	slot := c.unsafeGetSlot(pulseNumber)
	if slot == nil {
		return errors.Errorf("[ SinkPush ] can't get slot by pulse number %d", pulseNumber)
	}
	err := slot.SinkPush(data)
	return errors.Wrap(err, "[ SinkPush ] can't push to queue")
}

// SinkPushAll adds several events to conveyor
func (c *PulseConveyor) SinkPushAll(pulseNumber insolar.PulseNumber, data []interface{}) error {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if !c.unsafeIsOperational() {
		return errors.New("[ SinkPushAll ] conveyor is not operational now")
	}
	slot := c.unsafeGetSlot(pulseNumber)
	if slot == nil {
		return errors.Errorf("[ SinkPushAll ] can't get slot by pulse number %d", pulseNumber)
	}
	err := slot.SinkPushAll(data)
	return errors.Wrap(err, "[ SinkPushAll ] can't push to queue")
}

// PreparePulse is preparing conveyor for working with provided pulse
func (c *PulseConveyor) PreparePulse(pulse insolar.Pulse, callback queue.SyncDone) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.state == insolar.ConveyorShuttingDown {
		return errors.New("[ PreparePulse ] conveyor is shut down")
	}

	if c.futurePulseData != nil {
		return errors.New("[ PreparePulse ] preparation was already done")
	}
	if c.futurePulseNumber == nil {
		c.slotMap[pulse.PulseNumber] = slot.NewWorkingSlot(slot.Future, pulse.PulseNumber, c.removeSlot)
		c.futurePulseNumber = &pulse.PulseNumber
	}
	if *c.futurePulseNumber != pulse.PulseNumber {
		return errors.New("[ PreparePulse ] received future pulse is different from expected")
	}

	expectedNumCallbacks := 1
	if c.presentPulseNumber != nil {
		expectedNumCallbacks++
	}

	barrierCallback := newBarrierCallback(expectedNumCallbacks, callback)

	futureSlot := c.slotMap[*c.futurePulseNumber]
	err := futureSlot.PushSignal(slot.PendingPulseSignal, barrierCallback)
	if err != nil {
		log.Panicf("[ PreparePulse ] can't send signal to future slot (for pulse %d), error - %s", c.futurePulseNumber, err)
	}

	if c.presentPulseNumber != nil {
		presentSlot := c.slotMap[*c.presentPulseNumber]

		err := presentSlot.PushSignal(slot.PendingPulseSignal, barrierCallback)
		if err != nil {
			log.Panicf("[ PreparePulse ] can't send signal to present slot (for pulse %d), error - %s", c.presentPulseNumber, err)
		}
	}

	c.futurePulseData = &pulse
	c.state = insolar.ConveyorPreparingPulse
	return nil
}

// PulseWithCallback contains info about new pulse and callback func
type PulseWithCallback interface {
	queue.SyncDone
	GetPulse() insolar.Pulse
}

type pulseWithCallback struct {
	callback queue.SyncDone
	pulse    insolar.Pulse
}

// PulseWithCallback creates new instance of pulseWithCallback
func NewPulseWithCallback(callback queue.SyncDone, pulse insolar.Pulse) PulseWithCallback {
	return &pulseWithCallback{
		callback: callback,
		pulse:    pulse,
	}
}

// GetCallback returns callback
func (p *pulseWithCallback) GetPulse() insolar.Pulse {
	return p.pulse
}

// SetResult calls .SetResult() on func in callback param
func (p *pulseWithCallback) SetResult(result interface{}) {
	p.callback.SetResult(result)
}

type waitGroupSyncDone struct {
	sync.WaitGroup
}

// SetResult implements SyncDone
func (sd *waitGroupSyncDone) SetResult(result interface{}) {
	sd.Done()
}

// ActivatePulse activates conveyor with prepared pulse
func (c *PulseConveyor) ActivatePulse() error {
	c.lock.Lock()

	if c.state == insolar.ConveyorShuttingDown {
		c.lock.Unlock()
		return errors.New("[ ActivatePulse ] conveyor is shut down")
	}

	if c.futurePulseData == nil {
		c.lock.Unlock()
		return errors.New("[ ActivatePulse ] preparation missing")
	}

	wg := waitGroupSyncDone{}
	numCallbacks := 1
	if c.presentPulseNumber != nil {
		numCallbacks++
	}
	wg.Add(numCallbacks)

	futureSlot := c.slotMap[*c.futurePulseNumber]
	callback := NewPulseWithCallback(&wg, *c.futurePulseData)

	err := futureSlot.PushSignal(slot.ActivatePulseSignal, callback)
	if err != nil {
		c.lock.Unlock()
		log.Panicf("[ ActivatePulse ] can't send signal to future slot (for pulse %d), error - %s", c.futurePulseNumber, err)
	}

	if c.presentPulseNumber != nil {
		presentSlot := c.slotMap[*c.presentPulseNumber]
		err = presentSlot.PushSignal(slot.ActivatePulseSignal, &wg)
		if err != nil {
			c.lock.Unlock()
			log.Panicf("[ ActivatePulse ] can't send signal to present slot (for pulse %d), error - %s", c.presentPulseNumber, err)
		}
	}

	c.presentPulseNumber = c.futurePulseNumber
	c.futurePulseNumber = &c.futurePulseData.NextPulseNumber

	c.slotMap[*c.futurePulseNumber] = slot.NewWorkingSlot(slot.Future, *c.futurePulseNumber, c.removeSlot)

	c.futurePulseData = nil
	c.state = insolar.ConveyorActive
	c.lock.Unlock()
	wg.Wait()

	return nil
}

// BarrierCallback wait for required number of SetResult.
// After that it invokes SetResult on given callback and forward there last result from all SetResults
type BarrierCallback struct {
	wg     *sync.WaitGroup
	result interface{}
}

const defaultHash = "0c60ae04fbb17fe36f4e84631a5b8f3cd6d0cd46e80056bdfec97fd305f764daadef8ae1adc89b203043d7e2af1fb341df0ce5f66dfe3204ec3a9831532a8e4c"

func newBarrierCallback(num int, callback queue.SyncDone) *BarrierCallback {
	var wg sync.WaitGroup
	wg.Add(num)

	bc := &BarrierCallback{
		wg: &wg,
	}

	go func(bc *BarrierCallback) {
		wg.Wait()
		// TODO: this situation (no present pulse) must be handled in different way
		if num == 1 && bc.result == nil {
			log.Info("There is no present pulse and future pulse callback returned nil, set []byte{1, 2, 3} as result")
			hash, _ := hex.DecodeString(defaultHash)
			bc.result = hash
		}
		callback.SetResult(bc.result)
	}(bc)

	return bc
}

// SetResult saves result if it's not nil and invoke waitGroup.Done
func (c *BarrierCallback) SetResult(result interface{}) {
	if result != nil {
		c.result = result
	}
	c.wg.Done()
}
