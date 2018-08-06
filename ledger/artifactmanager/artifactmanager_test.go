/*
 *    Copyright 2018 INS Ecosystem
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package artifactmanager

import (
	"os"
	"testing"

	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/stretchr/testify/assert"
)

// Mock to imitate storage
type LedgerMock struct {
	Records       map[record.ID]record.Record
	ClassIndexes  map[record.ID]*index.ClassLifeline
	ObjectIndexes map[record.ID]*index.ObjectLifeline
}

func (mock *LedgerMock) GetRecord(k record.Key) (record.Record, bool) {
	rec, ok := mock.Records[record.Key2ID(k)]
	return rec, ok
}

func (mock *LedgerMock) AddRecord(rec record.Record) (record.Reference, error) {
	raw, _ := record.EncodeToRaw(rec)
	raw.Hash()
	var id record.ID
	copy(raw.Hash(), id[:])
	mock.Records[record.ID{}] = rec
	return record.Reference{
		Domain: record.ID{},
		Record: id,
	}, nil
}

func (mock *LedgerMock) GetClassIndex(id record.ID) (*index.ClassLifeline, bool) {
	idx, ok := mock.ClassIndexes[id]
	return idx, ok
}

func (mock *LedgerMock) SetClassIndex(id record.ID, idx *index.ClassLifeline) error {
	mock.ClassIndexes[id] = idx
	return nil
}

func (mock *LedgerMock) GetObjectIndex(id record.ID) (*index.ObjectLifeline, bool) {
	idx, ok := mock.ObjectIndexes[id]
	return idx, ok
}

func (mock *LedgerMock) SetObjectIndex(id record.ID, idx *index.ObjectLifeline) error {
	mock.ObjectIndexes[id] = idx
	return nil
}

func (mock *LedgerMock) Close() error {
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestDeployCodeCreatesRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	ref, err := manager.DeployCode(requestRef)
	assert.Nil(t, err)
	assert.Equal(t, ledger.Records[ref.Record], &record.CodeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
		},
	})
}

func TestActivateClassVerifiesCodeReference(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	_, err := manager.ActivateClass(requestRef, record.Reference{}, record.Memory{})
	assert.NotNil(t, err)
	notCodeRef, _ := ledger.AddRecord(&record.ClassActivateRecord{})
	_, err = manager.ActivateClass(requestRef, notCodeRef, record.Memory{})
	assert.NotNil(t, err)
}

func TestActivateClassCreatesActivateRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	memory := record.Memory{1, 2, 3}
	codeRef, _ := ledger.AddRecord(&record.CodeRecord{})
	activateRef, err := manager.ActivateClass(requestRef, codeRef, memory)
	assert.Nil(t, err)
	activateRec, isFound := ledger.GetRecord(record.ID2Key(activateRef.Record))
	assert.Equal(t, isFound, true)
	assert.Equal(t, activateRec, &record.ClassActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
		},
		CodeRecord:    codeRef,
		DefaultMemory: memory,
	})
}

func TestDeactivateClassVerifiesClassReference(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	_, err := manager.DeactivateClass(requestRef, record.Reference{})
	assert.NotNil(t, err)
	notClassRef, _ := ledger.AddRecord(&record.CodeRecord{})
	_, err = manager.DeactivateClass(requestRef, notClassRef)
	assert.NotNil(t, err)
}

func TestDeactivateClassVerifiesClassIsActive(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	classRef, _ := ledger.AddRecord(&record.ClassActivateRecord{})
	deactivateRef, _ := ledger.AddRecord(&record.DeactivationRecord{})
	ledger.SetClassIndex(classRef.Record, &index.ClassLifeline{
		LatestStateID: deactivateRef.Record,
	})
	_, err := manager.DeactivateClass(requestRef, classRef)
	assert.NotNil(t, err)
}

func TestDeactivateClassCreatesDeactivateRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	classRef, _ := ledger.AddRecord(&record.ClassActivateRecord{})
	ledger.SetClassIndex(classRef.Record, &index.ClassLifeline{
		LatestStateID: classRef.Record,
	})
	deactivateRef, err := manager.DeactivateClass(requestRef, classRef)
	assert.Nil(t, err)
	deactivateRec, isFound := ledger.GetRecord(record.ID2Key(deactivateRef.Record))
	assert.Equal(t, isFound, true)
	assert.Equal(t, deactivateRec, &record.DeactivationRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
			AmendedRecord: classRef,
		},
	})
}

func TestUpdateClassVerifiesClassReference(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	_, err := manager.UpdateClass(requestRef, record.Reference{}, []record.MemoryMigrationCode{})
	assert.NotNil(t, err)
	notClassRef, _ := ledger.AddRecord(&record.CodeRecord{})
	_, err = manager.UpdateClass(requestRef, notClassRef, []record.MemoryMigrationCode{})
	assert.NotNil(t, err)
}

func TestUpdateClassVerifiesClassIsActive(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	classRef, _ := ledger.AddRecord(&record.ClassActivateRecord{})
	deactivateRef, _ := ledger.AddRecord(&record.DeactivationRecord{})
	ledger.SetClassIndex(classRef.Record, &index.ClassLifeline{
		LatestStateID: deactivateRef.Record,
	})
	_, err := manager.UpdateClass(requestRef, classRef, nil)
	assert.NotNil(t, err)
}

func TestUpdateClassCreatesAmendRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	classRef, _ := ledger.AddRecord(&record.ClassActivateRecord{})
	ledger.SetClassIndex(classRef.Record, &index.ClassLifeline{
		LatestStateID: classRef.Record,
	})
	updateRef, err := manager.UpdateClass(requestRef, classRef, nil)
	assert.Nil(t, err)
	updateRec, isFound := ledger.GetRecord(record.ID2Key(updateRef.Record))
	assert.Equal(t, isFound, true)
	assert.Equal(t, updateRec, &record.ClassAmendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
			AmendedRecord: classRef,
		},
	})
}

func TestActivateObjVerifiesClassReference(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	_, err := manager.ActivateObj(requestRef, record.Reference{}, record.Memory{})
	assert.NotNil(t, err)
	notClassRef, _ := ledger.AddRecord(&record.ObjectActivateRecord{})
	_, err = manager.ActivateClass(requestRef, notClassRef, record.Memory{})
	assert.NotNil(t, err)
}

func TestActivateObjCreatesActivateRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	memory := record.Memory{1, 2, 3}
	classRef, _ := ledger.AddRecord(&record.ClassActivateRecord{})
	ledger.SetClassIndex(classRef.Record, &index.ClassLifeline{})
	activateRef, err := manager.ActivateObj(requestRef, classRef, memory)
	assert.Nil(t, err)
	activateRec, isFound := ledger.GetRecord(record.ID2Key(activateRef.Record))
	assert.Equal(t, isFound, true)
	assert.Equal(t, activateRec, &record.ObjectActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
		},
		ClassActivateRecord: classRef,
		Memory:              memory,
	})
}

func TestDeactivateObjVerifiesObjReference(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	_, err := manager.DeactivateClass(requestRef, record.Reference{})
	assert.NotNil(t, err)
	notObjRef, _ := ledger.AddRecord(&record.ClassActivateRecord{})
	_, err = manager.DeactivateClass(requestRef, notObjRef)
	assert.NotNil(t, err)
}

func TestDeactivateObjVerifiesObjectIsActive(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	objRef, _ := ledger.AddRecord(&record.ObjectActivateRecord{})
	deactivateRef, _ := ledger.AddRecord(&record.DeactivationRecord{})
	ledger.SetObjectIndex(objRef.Record, &index.ObjectLifeline{
		LatestStateID: deactivateRef.Record,
	})
	_, err := manager.DeactivateObj(requestRef, objRef)
	assert.NotNil(t, err)
}

func TestDeactivateObjCreatesDeactivateRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	objRef, _ := ledger.AddRecord(&record.ObjectActivateRecord{})
	ledger.SetObjectIndex(objRef.Record, &index.ObjectLifeline{
		LatestStateID: objRef.Record,
	})
	deactivateRef, err := manager.DeactivateObj(requestRef, objRef)
	assert.Nil(t, err)
	deactivateRec, isFound := ledger.GetRecord(record.ID2Key(deactivateRef.Record))
	assert.Equal(t, isFound, true)
	assert.Equal(t, deactivateRec, &record.DeactivationRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
			AmendedRecord: objRef,
		},
	})
}

func TestUpdateObjVerifiesObjectReference(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	_, err := manager.UpdateObj(requestRef, record.Reference{}, nil)
	assert.NotNil(t, err)
	notObjRef, _ := ledger.AddRecord(&record.CodeRecord{})
	_, err = manager.UpdateObj(requestRef, notObjRef, nil)
	assert.NotNil(t, err)
}

func TestUpdateObjVerifiesObjectIsActive(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	objRef, _ := ledger.AddRecord(&record.ObjectActivateRecord{})
	deactivateRef, _ := ledger.AddRecord(&record.DeactivationRecord{})
	ledger.SetObjectIndex(objRef.Record, &index.ObjectLifeline{
		LatestStateID: deactivateRef.Record,
	})
	_, err := manager.UpdateClass(requestRef, objRef, nil)
	assert.NotNil(t, err)
}

func TestUpdateObjCreatesAmendRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	objRef, _ := ledger.AddRecord(&record.ObjectActivateRecord{})
	ledger.SetObjectIndex(objRef.Record, &index.ObjectLifeline{
		LatestStateID: objRef.Record,
	})
	memory := record.Memory{1, 2, 3}
	updateRef, err := manager.UpdateObj(requestRef, objRef, memory)
	assert.Nil(t, err)
	updateRec, isFound := ledger.GetRecord(record.ID2Key(updateRef.Record))
	assert.Equal(t, isFound, true)
	assert.Equal(t, updateRec, &record.ObjectAmendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
			AmendedRecord: objRef,
		},
		NewMemory: memory,
	})
}

func TestAppendObjDelegateVerifiesObjRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	_, err := manager.AppendObjDelegate(requestRef, record.Reference{}, nil)
	assert.NotNil(t, err)
	notObjRef, _ := ledger.AddRecord(&record.CodeRecord{})
	_, err = manager.AppendObjDelegate(requestRef, notObjRef, nil)
	assert.NotNil(t, err)
}

func TestAppendObjDelegateVerifiesObjectIsActive(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	objRef, _ := ledger.AddRecord(&record.ObjectActivateRecord{})
	deactivateRef, _ := ledger.AddRecord(&record.DeactivationRecord{})
	ledger.SetObjectIndex(objRef.Record, &index.ObjectLifeline{
		LatestStateID: deactivateRef.Record,
	})
	_, err := manager.AppendObjDelegate(requestRef, objRef, nil)
	assert.NotNil(t, err)
}

func TestAppendObjDelegateCreatesAmendRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	objRef, _ := ledger.AddRecord(&record.ObjectActivateRecord{})
	ledger.SetObjectIndex(objRef.Record, &index.ObjectLifeline{
		LatestStateID: objRef.Record,
	})
	memory := record.Memory{1, 2, 3}
	appendRef, err := manager.AppendObjDelegate(requestRef, objRef, memory)
	assert.Nil(t, err)
	appendRec, isFound := ledger.GetRecord(record.ID2Key(appendRef.Record))
	objIndex, ok := ledger.GetObjectIndex(objRef.Record)
	assert.True(t, ok)
	assert.Equal(t, objIndex.AppendIDs, []record.ID{appendRef.Record})
	assert.Equal(t, isFound, true)
	assert.Equal(t, appendRec, &record.ObjectAppendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
			AmendedRecord: objRef,
		},
		AppendMemory: memory,
	})
}
