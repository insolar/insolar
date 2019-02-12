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

package storage

import (
	"context"
	"sync"

	"github.com/insolar/insolar/core"
)

type pulseTrackerMemory struct {
	memory      map[core.PulseNumber]*Pulse
	latestPulse core.PulseNumber
	mutex       sync.RWMutex
}

func NewPulseTrackerMemory() PulseTracker {
	return &pulseTrackerMemory{memory: make(map[core.PulseNumber]*Pulse)}
}

func (p *pulseTrackerMemory) GetPulse(ctx context.Context, num core.PulseNumber) (*Pulse, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.getPulse(ctx, num)
}

func (p *pulseTrackerMemory) GetPreviousPulse(ctx context.Context, num core.PulseNumber) (*Pulse, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.getNthPrevPulse(ctx, 1, num)
}

func (p *pulseTrackerMemory) GetNthPrevPulse(ctx context.Context, n uint, from core.PulseNumber) (*Pulse, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.getNthPrevPulse(ctx, n, from)
}

func (p *pulseTrackerMemory) GetLatestPulse(ctx context.Context) (*Pulse, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	return p.getLatestPulse(ctx)
}

func (p *pulseTrackerMemory) AddPulse(ctx context.Context, pulse core.Pulse) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	pn := pulse.PulseNumber

	if pn < p.latestPulse {
		return ErrLesserPulse
	}

	if pn == p.latestPulse {
		return ErrOverride
	}

	var (
		previousPulseNumber  core.PulseNumber
		previousSerialNumber int
	)

	// Save new pulse.
	newPulse := Pulse{
		Prev:         &previousPulseNumber,
		SerialNumber: previousSerialNumber + 1,
		Pulse:        pulse,
	}

	p.memory[pn] = &newPulse
	p.latestPulse = pn

	return nil
}

func (p *pulseTrackerMemory) getPulse(ctx context.Context, num core.PulseNumber) (*Pulse, error) {
	pulse, ok := p.memory[num]

	if !ok {
		return nil, ErrPulseNotFound
	}

	return pulse, nil
}

func (p *pulseTrackerMemory) getNthPrevPulse(ctx context.Context, n uint, num core.PulseNumber) (*Pulse, error) {
	pulse, err := p.getPulse(ctx, num)
	if err != nil {
		return nil, err
	}

	for n > 0 {
		if pulse.Prev == nil {
			return nil, ErrPrevPulseNotFound
		}

		pulse, err = p.getPulse(ctx, *pulse.Prev)

		if err != nil {
			return nil, ErrPulseNotFound
		}
		n--
	}
	return pulse, nil
}

func (p *pulseTrackerMemory) getLatestPulse(ctx context.Context) (*Pulse, error) {
	if p.latestPulse == 0 {
		return nil, ErrEmptyLatestPulse
	}

	return p.getPulse(ctx, p.latestPulse)
}
