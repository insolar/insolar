//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package storage

import (
	"context"
	"sync"

	"github.com/insolar/insolar/insolar"
)

// PulseStorage implements insolar.PulseStorage
type PulseStorage struct {
	PulseTracker PulseTracker `inject:""`
	rwLock       sync.RWMutex
	currentPulse *insolar.Pulse
}

// NewPulseStorage creates new pulse storage
func NewPulseStorage() *PulseStorage {
	return &PulseStorage{}
}

// Current returns current pulse of the system
func (ps *PulseStorage) Current(ctx context.Context) (*insolar.Pulse, error) {
	ps.rwLock.RLock()

	if ps.currentPulse == nil {
		ps.rwLock.RUnlock()

		ps.rwLock.Lock()
		defer ps.rwLock.Unlock()

		if ps.currentPulse == nil {
			currentPulse, err := ps.PulseTracker.GetLatestPulse(ctx)
			if err != nil {
				return nil, err
			}
			ps.currentPulse = &currentPulse.Pulse
		}

		return ps.currentPulse, nil
	}

	defer ps.rwLock.RUnlock()
	return ps.currentPulse, nil
}

func (ps *PulseStorage) Set(pulse *insolar.Pulse) {
	ps.currentPulse = pulse
}

// Lock takes lock on parent's pulse storage
func (ps *PulseStorage) Lock() {
	ps.rwLock.Lock()
}

// Unlock takes unlock on parent's pulse storage
func (ps *PulseStorage) Unlock() {
	ps.rwLock.Unlock()
}
