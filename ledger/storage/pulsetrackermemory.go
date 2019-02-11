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
	memory map[core.PulseNumber]*Pulse
	mutex  sync.RWMutex
}

func NewPulseTrackerMemory() PulseTracker {
	return &pulseTrackerMemory{memory: make(map[core.PulseNumber]*Pulse)}
}

func (p *pulseTrackerMemory) GetPulse(ctx context.Context, num core.PulseNumber) (*Pulse, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	pulse, ok := p.memory[num]

	if !ok {
		return nil, ErrPulseNotFound
	}

	return pulse, nil
}

func (p *pulseTrackerMemory) GetPreviousPulse(ctx context.Context, num core.PulseNumber) (*Pulse, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	pulse, exists := p.memory[num]

	if !exists {
		return nil, ErrPulseNotFound
	}

	if pulse.Prev == nil {
		return nil, ErrPrevPulseNotFound
	}

	resultPulse, exists := p.memory[*pulse.Prev]

	if !exists {
		return nil, ErrPulseNotFound
	}

	return resultPulse, nil
}

func (p *pulseTrackerMemory) GetNthPrevPulse(ctx context.Context, n uint, from core.PulseNumber) (*Pulse, error) {
	panic("implement me")
}

func (p *pulseTrackerMemory) GetLatestPulse(ctx context.Context) (*Pulse, error) {
	panic("implement me")
}

func (p *pulseTrackerMemory) AddPulse(ctx context.Context, pulse core.Pulse) error {
	panic("implement me")
}
