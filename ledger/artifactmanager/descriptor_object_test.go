package artifactmanager

import (
	"testing"

	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/stretchr/testify/assert"
)

func prepareObjectDescriptorTest() (
	*LedgerMock, *LedgerArtifactManager, *record.ObjectActivateRecord, record.Reference,
) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{
		storer:   &ledger,
		archPref: []record.ArchType{1},
	}
	rec := record.ObjectActivateRecord{Memory: record.Memory{1}}
	ref := addRecord(&ledger, &rec)

	return &ledger, &manager, &rec, ref
}

func TestObjectDescriptor_GetMemory(t *testing.T) {
	ledger, manager, objRec, objRef := prepareObjectDescriptorTest()
	amendRec := record.ObjectAmendRecord{NewMemory: record.Memory{2}}
	amendRef := addRecord(ledger, &amendRec)
	idx := index.ObjectLifeline{
		LatestStateID: amendRef.Record,
	}
	ledger.SetObjectIndex(objRef.Record, &idx)

	desc := ObjectDescriptor{
		manager:           manager,
		activateRecord:    objRec,
		latestAmendRecord: nil,
		lifelineIndex:     &idx,
	}
	mem, err := desc.GetMemory()
	assert.NoError(t, err)
	assert.Equal(t, record.Memory{1}, mem)

	desc = ObjectDescriptor{
		manager:           manager,
		activateRecord:    objRec,
		latestAmendRecord: &amendRec,
		lifelineIndex:     &idx,
	}
	mem, err = desc.GetMemory()
	assert.NoError(t, err)
	assert.Equal(t, record.Memory{2}, mem)
}

func TestObjectDescriptor_GetDelegates(t *testing.T) {
	ledger, manager, objRec, objRef := prepareObjectDescriptorTest()
	appendRec1 := record.ObjectAppendRecord{AppendMemory: record.Memory{2}}
	appendRec2 := record.ObjectAppendRecord{AppendMemory: record.Memory{3}}
	appendRef1 := addRecord(ledger, &appendRec1)
	appendRef2 := addRecord(ledger, &appendRec2)
	idx := index.ObjectLifeline{
		LatestStateID: objRef.Record,
		AppendIDs:     []record.ID{appendRef1.Record, appendRef2.Record},
	}
	ledger.SetObjectIndex(objRef.Record, &idx)

	desc := ObjectDescriptor{
		manager:           manager,
		activateRecord:    objRec,
		latestAmendRecord: nil,
		lifelineIndex:     &idx,
	}

	appends, err := desc.GetDelegates()
	assert.NoError(t, err)
	assert.Equal(t, []record.Memory{{2}, {3}}, appends)
}
