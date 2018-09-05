package storage

import (
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
)

type Ledger interface {
	GetRecord(ref *record.Reference) (record.Record, error)
	SetRecord(rec record.Record) (*record.Reference, error)
	GetClassIndex(ref *record.Reference) (*index.ClassLifeline, error)
	SetClassIndex(ref *record.Reference, idx *index.ClassLifeline) error
	GetObjectIndex(ref *record.Reference) (*index.ObjectLifeline, error)
	SetObjectIndex(ref *record.Reference, idx *index.ObjectLifeline) error
}
