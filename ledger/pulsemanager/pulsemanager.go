/*
 *    Copyright 2018 Insolar
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

package pulsemanager

import (
	"sync"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/jetcoordinator"
	"github.com/insolar/insolar/ledger/storage"
)

// PulseManager implements core.PulseManager.
type PulseManager struct {
	lock        sync.RWMutex
	db          *storage.DB
	lr          core.LogicRunner
	coordinator *jetcoordinator.JetCoordinator
}

// Current returns current pulse structure.
func (m *PulseManager) Current() (*core.Pulse, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	pulseNum := m.db.GetCurrentPulse()
	entropy, err := m.db.GetEntropy(pulseNum)
	if err != nil {
		return nil, err
	}
	pulse := core.Pulse{
		PulseNumber: pulseNum,
		Entropy:     *entropy,
	}
	return &pulse, nil
}

// Set set's new pulse and closes current jet drop.
func (m *PulseManager) Set(pulse core.Pulse) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	err := m.db.SetEntropy(pulse.PulseNumber, pulse.Entropy)
	if err != nil {
		return err
	}
	drop, err := m.coordinator.CreateDrop(pulse.PulseNumber)
	if err != nil {
		return err

	}

	_ = drop // TODO: send drop to the validators
	m.db.SetCurrentPulse(pulse.PulseNumber)
	return m.lr.OnPulse(pulse)
}

// NewPulseManager creates PulseManager instance.
func NewPulseManager(db *storage.DB, coordinator *jetcoordinator.JetCoordinator) (*PulseManager, error) {
	pm := PulseManager{
		db:          db,
		coordinator: coordinator,
	}
	return &pm, nil
}

func (m *PulseManager) Link(c core.Components) error {
	m.lr = c.LogicRunner
	return nil
}
