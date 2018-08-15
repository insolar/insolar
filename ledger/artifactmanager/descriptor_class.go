package artifactmanager

import (
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/pkg/errors"
)

// ClassDescriptor represents meta info required to fetch all class data.
type ClassDescriptor struct {
	StateRef record.Reference

	manager           *LedgerArtifactManager
	fromState         record.Reference
	activateRecord    *record.ClassActivateRecord
	latestAmendRecord *record.ClassAmendRecord
	lifelineIndex     *index.ClassLifeline
}

// GetCode fetches the latest class code known to storage. Code will be fetched according to architecture preferences
// set via SetArchPref in artifact manager. If preferences are not provided, an error will be returned.
func (d *ClassDescriptor) GetCode() ([]byte, error) {
	codeRef := d.activateRecord.CodeRecord
	if d.latestAmendRecord != nil {
		codeRef = d.latestAmendRecord.NewCode
	}
	code, err := d.manager.getCodeRecordCode(codeRef)
	if err != nil {
		return nil, err
	}

	return code, nil
}

// GetMigrations fetches all migrations from provided to artifact manager state to the last state known to storage. VM
// is responsible for applying these migrations and updating objects.
func (d *ClassDescriptor) GetMigrations() ([][]byte, error) {
	var amends []*record.ClassAmendRecord
	// Search for provided state in class amends from the end of the list.
	// Record keys are hashes and are not incremental, so we can't say if provided state should be before or after.
	for i := len(d.lifelineIndex.AmendRefs) - 1; i >= 0; i-- {
		amendRef := d.lifelineIndex.AmendRefs[i]
		if d.fromState.IsEqual(amendRef) {
			break // Provided state is found. It means we now have all the amends we need.
		}
		rec, err := d.manager.storer.GetRecord(&amendRef)
		if err != nil {
			return nil, errors.Wrap(err, "inconsistent class index")
		}
		amendRec, ok := rec.(*record.ClassAmendRecord)
		if !ok {
			return nil, errors.New("inconsistent class index")
		}
		amends = append(amends, amendRec)
	}
	// Reverse found amends again (we appended them from the end) so they'll have the original order.
	sortedAmends := make([]*record.ClassAmendRecord, len(amends))
	for i, amend := range amends {
		sortedAmends[len(amends)-i-1] = amend
	}

	// Flatten the migrations list from amends.
	var migrations [][]byte
	for _, amendRec := range sortedAmends {
		for _, codeRef := range amendRec.Migrations {
			code, err := d.manager.getCodeRecordCode(codeRef)
			if err != nil {
				return nil, errors.Wrap(err, "invalid migration reference in amend record")
			}
			migrations = append(migrations, code)
		}
	}

	return migrations, nil
}
