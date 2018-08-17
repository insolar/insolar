package jetcoordinator

import (
	"errors"

	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
)

// JetCoordinator is responsible for all jet interactions
type JetCoordinator struct {
	storage storage.LedgerStorer
}

// Pulse creates new jet drop and ends current slot. This should be called when receiving a new pulse from pulsar.
func (jc *JetCoordinator) Pulse(newPulse record.PulseNum) error {
	if newPulse-jc.storage.GetCurrentPulse() != 1 {
		return errors.New("wrong pulse")
	}
	// TODO: increment stored pulse number and wait for all records from previous pulse to store
	drop, err := CreateJetDrop(jc.storage, jc.storage.GetCurrentPulse(), newPulse)
	if err != nil {
		return err
	}

	return jc.storage.SetDrop(newPulse, drop)
}
