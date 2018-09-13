package pulsemanager

import (
	"fmt"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/jetcoordinator"
	"github.com/insolar/insolar/ledger/storage"
)

type PulseManager struct {
	db          *storage.DB
	coordinator *jetcoordinator.JetCoordinator
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
	current := m.db.GetCurrentPulse()
	if pulse.PulseNumber-current != 1 {
		panic(fmt.Sprintf("Wrong pulse, got %v, but current is %v\n", pulse, current))
	}

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

	return nil
}

func NewPulseManager(db *storage.DB, coordinator *jetcoordinator.JetCoordinator) (*PulseManager, error) {
	pm := PulseManager{
		db:          db,
		coordinator: coordinator,
	}
	return &pm, nil
}
