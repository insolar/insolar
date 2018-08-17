package jetcoordinator

import (
	"errors"

	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
)

type JetCoordinator struct {
	storage storage.LedgerStorer
}

func (jc *JetCoordinator) Pulse(newPulse record.PulseNum) error {
	if newPulse-jc.storage.GetCurrentPulse() != 1 {
		return errors.New("wrong pulse")
	}

	drop, err := CreateJetDrop(jc.storage, jc.storage.GetCurrentPulse(), newPulse)
	if err != nil {
		return err
	}

	jc.storage.SetDrop(newPulse, drop)

	return nil
}
