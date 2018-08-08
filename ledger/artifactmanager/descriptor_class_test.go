package artifactmanager

import (
	"testing"

	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/stretchr/testify/assert"
)

func prepareClassDescriptorTest() (*LedgerMock, *LedgerArtifactManager, *record.ClassActivateRecord, record.Reference) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{
		storer:   &ledger,
		archPref: []record.ArchType{1},
	}
	rec := record.ClassActivateRecord{}
	ref := addRecord(&ledger, &rec)

	return &ledger, &manager, &rec, ref
}

func TestClassDescriptor_GetCode(t *testing.T) {
	ledger, manager, classRec, classRef := prepareClassDescriptorTest()
	codeRef := addRecord(ledger, &record.CodeRecord{TargetedCode: map[record.ArchType][]byte{
		1: {1, 2, 3},
	}})
	amendRec := record.ClassAmendRecord{NewCode: codeRef}
	amendRef := addRecord(ledger, &amendRec)
	idx := index.ClassLifeline{
		LatestStateID: amendRef.Record,
	}
	ledger.SetClassIndex(classRef.Record, &idx)

	desc := ClassDescriptor{
		manager:           manager,
		activateRecord:    classRec,
		latestAmendRecord: &amendRec,
		lifelineIndex:     &idx,
	}

	code, err := desc.GetCode()
	assert.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, code)
}

func TestClassDescriptor_GetMigrations(t *testing.T) {
	ledger, manager, classRec, classRef := prepareClassDescriptorTest()
	codeRef1 := addRecord(ledger, &record.CodeRecord{TargetedCode: map[record.ArchType][]byte{
		record.ArchType(1): {1},
	}})
	codeRef2 := addRecord(ledger, &record.CodeRecord{TargetedCode: map[record.ArchType][]byte{
		record.ArchType(1): {2},
	}})
	codeRef3 := addRecord(ledger, &record.CodeRecord{TargetedCode: map[record.ArchType][]byte{
		record.ArchType(1): {3},
	}})
	codeRef4 := addRecord(ledger, &record.CodeRecord{TargetedCode: map[record.ArchType][]byte{
		record.ArchType(1): {4},
	}})

	amendRec3 := record.ClassAmendRecord{Migrations: []record.Reference{codeRef4}}
	amendRef1 := addRecord(ledger, &record.ClassAmendRecord{Migrations: []record.Reference{codeRef1}})
	amendRef2 := addRecord(ledger, &record.ClassAmendRecord{Migrations: []record.Reference{codeRef2, codeRef3}})
	amendRef3 := addRecord(ledger, &amendRec3)
	idx := index.ClassLifeline{
		LatestStateID: amendRef2.Record,
		AmendIDs:      []record.ID{amendRef1.Record, amendRef2.Record, amendRef3.Record},
	}
	ledger.SetClassIndex(classRef.Record, &idx)

	desc := ClassDescriptor{
		manager:           manager,
		fromState:         amendRef1,
		activateRecord:    classRec,
		latestAmendRecord: &amendRec3,
		lifelineIndex:     &idx,
	}

	migrations, err := desc.GetMigrations()
	assert.NoError(t, err)
	assert.Equal(t, [][]byte{{2}, {3}, {4}}, migrations)
}
