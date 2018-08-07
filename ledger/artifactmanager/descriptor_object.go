package artifactmanager

import (
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/pkg/errors"
)

// ObjectDescriptor represents meta info required to fetch all object data
type ObjectDescriptor struct {
	StateRef record.Reference

	manager           *LedgerArtifactManager
	activateRecord    *record.ObjectActivateRecord
	latestAmendRecord *record.ObjectAmendRecord
	lifelineIndex     *index.ObjectLifeline
}

func (d *ObjectDescriptor) GetMemory() (record.Memory, error) {
	if d.latestAmendRecord != nil {
		return d.latestAmendRecord.NewMemory, nil
	}

	return d.activateRecord.Memory, nil
}

func (d *ObjectDescriptor) GetDelegates() ([]record.Memory, error) {
	var delegates []record.Memory
	for _, appendID := range d.lifelineIndex.AppendIDs {
		rec, err := d.manager.storer.GetRecord(appendID)
		if err != nil {
			return nil, errors.Wrap(err, "invalid append reference in object index")
		}
		appendRec, ok := rec.(*record.ObjectAppendRecord)
		if !ok {
			return nil, errors.New("invalid append reference in object index")
		}
		delegates = append(delegates, appendRec.AppendMemory)
	}

	return delegates, nil
}
