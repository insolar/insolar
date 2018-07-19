package ledger

import (
	"github.com/insolar/insolar/ledger/record"
)

type Ledger interface {
	Get(id record.RecordHash) (bool, record.Record)
	Set(record record.Record) error
	Update(id record.RecordHash, record record.Record) error
}

func NewLedger() (Ledger, error) {
	return newLedger()
}
