package artifactmanager

import (
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/pkg/errors"
)

// ClassDescriptor represents meta info required to fetch all class data
type ClassDescriptor struct {
	StateRef record.Reference

	manager           *LedgerArtifactManager
	fromState         record.Reference
	activateRecord    *record.ClassActivateRecord
	latestAmendRecord *record.ClassAmendRecord
	lifelineIndex     *index.ClassLifeline
}

func (d *ClassDescriptor) getRecordCode(codeID record.ID) ([]byte, error) {
	codeRec, err := d.manager.getCodeRecord(codeID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve code record")
	}
	code, err := codeRec.GetCode(d.manager.archPref)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve code")
	}

	return code, nil
}

func (d *ClassDescriptor) GetCode() ([]byte, error) {
	codeRef := d.activateRecord.CodeRecord
	if d.latestAmendRecord != nil {
		codeRef = d.latestAmendRecord.NewCode
	}
	code, err := d.getRecordCode(codeRef.Record)
	if err != nil {
		return nil, err
	}

	return code, nil
}

func (d *ClassDescriptor) GetMigrations() ([][]byte, error) {
	var amends []*record.ClassAmendRecord
	for i := len(d.lifelineIndex.AmendIDs) - 1; i >= 0; i-- {
		amendID := d.lifelineIndex.AmendIDs[i]
		if d.fromState.Record == amendID {
			break
		}
		rec, err := d.manager.storer.GetRecord(amendID)
		if err != nil {
			return nil, errors.Wrap(err, "inconsistent class index")
		}
		amendRec, ok := rec.(*record.ClassAmendRecord)
		if !ok {
			return nil, errors.New("inconsistent class index")
		}
		amends = append(amends, amendRec)
	}
	sortedAmends := make([]*record.ClassAmendRecord, len(amends))
	for i, amend := range amends {
		sortedAmends[len(amends)-i-1] = amend
	}

	var migrations [][]byte
	for _, amendRec := range sortedAmends {
		for _, codeID := range amendRec.Migrations {
			code, err := d.getRecordCode(codeID.Record)
			if err != nil {
				return nil, errors.Wrap(err, "invalid migration reference in amend record")
			}
			migrations = append(migrations, code)
		}
	}

	return migrations, nil
}
