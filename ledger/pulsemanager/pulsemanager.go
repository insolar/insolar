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
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/ledger/jetcoordinator"
	"github.com/insolar/insolar/ledger/jetdrop"
	"github.com/insolar/insolar/ledger/storage"
)

// PulseManager implements core.PulseManager.
type PulseManager struct {
	db          *storage.DB
	lr          core.LogicRunner
	coordinator *jetcoordinator.JetCoordinator
	bus         core.MessageBus
}

// Current returns current pulse structure.
func (m *PulseManager) Current() (*core.Pulse, error) {
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
	err := m.db.SetEntropy(pulse.PulseNumber, pulse.Entropy)
	if err != nil {
		return err
	}

	drop, err := m.coordinator.CreateDrop(pulse.PulseNumber)
	if err != nil {
		return err
	}

	dropSerialized, err := jetdrop.Encode(drop)
	if err != nil {
		return err
	}

	_, err = m.bus.Send(&message.JetDrop{Drop: dropSerialized})
	if err != nil {
		return err
	}
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

// Link links external components.
func (m *PulseManager) Link(components core.Components) error {
	m.bus = components.MessageBus
	m.lr = components.LogicRunner
	return nil
}
