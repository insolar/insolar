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
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// Mock to imitate storage
type LedgerMock struct {
	Records       map[record.ID]record.Record
	ClassIndexes  map[record.ID]*index.ClassLifeline
	ObjectIndexes map[record.ID]*index.ObjectLifeline
}

func (mock *LedgerMock) GetRecord(id record.ID) (record.Record, error) {
	rec, ok := mock.Records[id]
	if !ok {
		return nil, errors.New("record not found")
	}
	return rec, nil
}

func (mock *LedgerMock) SetRecord(rec record.Record) (record.ID, error) {
	raw, _ := record.EncodeToRaw(rec)
	var id record.ID
	copy(id[:], raw.Hash())
	mock.Records[id] = rec
	return id, nil
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

// helper to workaround reference generation from id
func addRecord(mock *LedgerMock, rec record.Record) record.Reference {
	id, _ := mock.SetRecord(rec)
	return record.Reference{
		Domain: rec.Domain(),
		Record: id,
	}
}

func TestLedgerArtifactManager_DeployCode(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	codeMap := map[record.ArchType][]byte{1: {1}}
	ref, err := manager.DeployCode(requestRef, codeMap)
	assert.Nil(t, err)
	assert.Equal(t, ledger.Records[ref.Record], &record.CodeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
		},
		TargetedCode: codeMap,
	})
}

func TestLedgerArtifactManager_ActivateClass_VerifiesRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	_, err := manager.ActivateClass(requestRef, record.Reference{}, record.Memory{})
	assert.NotNil(t, err)
	notCodeRef := addRecord(&ledger, &record.ClassActivateRecord{})
	_, err = manager.ActivateClass(requestRef, notCodeRef, record.Memory{})
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_ActivateClass_CreatesCorrectRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	memory := record.Memory{1, 2, 3}
	codeRef := addRecord(&ledger, &record.CodeRecord{})
	activateRef, err := manager.ActivateClass(requestRef, codeRef, memory)
	assert.Nil(t, err)
	activateRec, getErr := ledger.GetRecord(activateRef.Record)
	assert.Nil(t, getErr)
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

func TestLedgerArtifactManager_DeactivateClass_VerifiesRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	_, err := manager.DeactivateClass(requestRef, record.Reference{})
	assert.NotNil(t, err)
	notClassRef := addRecord(&ledger, &record.CodeRecord{})
	_, err = manager.DeactivateClass(requestRef, notClassRef)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateClass_VerifiesClassIsActive(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	classRef := addRecord(&ledger, &record.ClassActivateRecord{})
	deactivateRef := addRecord(&ledger, &record.DeactivationRecord{})
	ledger.SetClassIndex(classRef.Record, &index.ClassLifeline{
		LatestStateID: deactivateRef.Record,
	})
	_, err := manager.DeactivateClass(requestRef, classRef)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateClass_CreatesCorrectRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	classRef := addRecord(&ledger, &record.ClassActivateRecord{})
	ledger.SetClassIndex(classRef.Record, &index.ClassLifeline{
		LatestStateID: classRef.Record,
	})
	deactivateRef, err := manager.DeactivateClass(requestRef, classRef)
	assert.Nil(t, err)
	deactivateRec, getErr := ledger.GetRecord(deactivateRef.Record)
	assert.Nil(t, getErr)
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

func TestLedgerArtifactManager_UpdateClass_VerifiesRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	_, err := manager.UpdateClass(requestRef, record.Reference{}, record.Reference{}, nil)
	assert.NotNil(t, err)
	notClassRef := addRecord(&ledger, &record.CodeRecord{})
	_, err = manager.UpdateClass(requestRef, notClassRef, record.Reference{}, nil)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateClass_VerifiesClassIsActive(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	classRef := addRecord(&ledger, &record.ClassActivateRecord{})
	deactivateRef := addRecord(&ledger, &record.DeactivationRecord{})
	codeRef := addRecord(&ledger, &record.CodeRecord{})
	ledger.SetClassIndex(classRef.Record, &index.ClassLifeline{
		LatestStateID: deactivateRef.Record,
	})
	_, err := manager.UpdateClass(requestRef, classRef, codeRef, nil)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateClass_CreatesCorrectRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	classRef := addRecord(&ledger, &record.ClassActivateRecord{})
	return
	ledger.SetClassIndex(classRef.Record, &index.ClassLifeline{
		LatestStateID: classRef.Record,
	})
	codeRef := addRecord(&ledger, &record.CodeRecord{})
	migrationRef := addRecord(&ledger, &record.CodeRecord{SourceCode: "test"})
	migrationRefs := []record.Reference{migrationRef}
	updateRef, err := manager.UpdateClass(requestRef, classRef, codeRef, migrationRefs)
	assert.Nil(t, err)
	updateRec, getErr := ledger.GetRecord(updateRef.Record)
	assert.Nil(t, getErr)
	assert.Equal(t, updateRec, &record.ClassAmendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: requestRef,
				},
			},
			AmendedRecord: classRef,
		},
		NewCode:    codeRef,
		Migrations: migrationRefs,
	})
}

func TestLedgerArtifactManager_ActivateObj_VerifiesRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	_, err := manager.ActivateObj(requestRef, record.Reference{}, record.Memory{})
	assert.NotNil(t, err)
	notClassRef := addRecord(&ledger, &record.ObjectActivateRecord{})
	_, err = manager.ActivateClass(requestRef, notClassRef, record.Memory{})
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_ActivateObj_CreatesCorrectRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	memory := record.Memory{1, 2, 3}
	classRef := addRecord(&ledger, &record.ClassActivateRecord{})
	ledger.SetClassIndex(classRef.Record, &index.ClassLifeline{
		LatestStateID: classRef.Record,
	})
	activateRef, err := manager.ActivateObj(requestRef, classRef, memory)
	assert.Nil(t, err)
	activateRec, err := ledger.GetRecord(activateRef.Record)
	assert.Nil(t, err)
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

func TestLedgerArtifactManager_DeactivateObj_VerifiesRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	_, err := manager.DeactivateClass(requestRef, record.Reference{})
	assert.NotNil(t, err)
	notObjRef := addRecord(&ledger, &record.ClassActivateRecord{})
	_, err = manager.DeactivateClass(requestRef, notObjRef)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateObj_VerifiesObjectIsActive(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	objRef := addRecord(&ledger, &record.ObjectActivateRecord{})
	deactivateRef := addRecord(&ledger, &record.DeactivationRecord{})
	ledger.SetObjectIndex(objRef.Record, &index.ObjectLifeline{
		LatestStateID: deactivateRef.Record,
	})
	_, err := manager.DeactivateObj(requestRef, objRef)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateObj_CreatesCorrectRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	objRef := addRecord(&ledger, &record.ObjectActivateRecord{})
	ledger.SetObjectIndex(objRef.Record, &index.ObjectLifeline{
		LatestStateID: objRef.Record,
	})
	deactivateRef, err := manager.DeactivateObj(requestRef, objRef)
	assert.Nil(t, err)
	deactivateRec, err := ledger.GetRecord(deactivateRef.Record)
	assert.Nil(t, err)
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

func TestLedgerArtifactManager_UpdateObj_VerifiesRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	_, err := manager.UpdateObj(requestRef, record.Reference{}, nil)
	assert.NotNil(t, err)
	notObjRef := addRecord(&ledger, &record.CodeRecord{})
	_, err = manager.UpdateObj(requestRef, notObjRef, nil)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateObj_VerifiesObjectIsActive(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	objRef := addRecord(&ledger, &record.ObjectActivateRecord{})
	deactivateRef := addRecord(&ledger, &record.DeactivationRecord{})
	ledger.SetObjectIndex(objRef.Record, &index.ObjectLifeline{
		LatestStateID: deactivateRef.Record,
	})
	_, err := manager.UpdateObj(requestRef, objRef, nil)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateObj_CreatesCorrectRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	objRef := addRecord(&ledger, &record.ObjectActivateRecord{})
	ledger.SetObjectIndex(objRef.Record, &index.ObjectLifeline{
		LatestStateID: objRef.Record,
	})
	memory := record.Memory{1, 2, 3}
	updateRef, err := manager.UpdateObj(requestRef, objRef, memory)
	assert.Nil(t, err)
	updateRec, err := ledger.GetRecord(updateRef.Record)
	assert.Nil(t, err)
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

func TestLedgerArtifactManager_AppendObjDelegate_VerifiesRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	_, err := manager.AppendObjDelegate(requestRef, record.Reference{}, nil)
	assert.NotNil(t, err)
	notObjRef := addRecord(&ledger, &record.CodeRecord{})
	_, err = manager.AppendObjDelegate(requestRef, notObjRef, nil)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_AppendObjDelegate_VerifiesClassIsActive(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	objRef := addRecord(&ledger, &record.ObjectActivateRecord{})
	deactivateRef := addRecord(&ledger, &record.DeactivationRecord{})
	ledger.SetObjectIndex(objRef.Record, &index.ObjectLifeline{
		LatestStateID: deactivateRef.Record,
	})
	_, err := manager.AppendObjDelegate(requestRef, objRef, nil)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_AppendObjDelegate_CreatesCorrectRecord(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	objRef := addRecord(&ledger, &record.ObjectActivateRecord{})
	ledger.SetObjectIndex(objRef.Record, &index.ObjectLifeline{
		LatestStateID: objRef.Record,
	})
	memory := record.Memory{1, 2, 3}
	appendRef, err := manager.AppendObjDelegate(requestRef, objRef, memory)
	assert.Nil(t, err)
	appendRec, err := ledger.GetRecord(appendRef.Record)
	objIndex, ok := ledger.GetObjectIndex(objRef.Record)
	assert.True(t, ok)
	assert.Equal(t, objIndex.AppendIDs, []record.ID{appendRef.Record})
	assert.Nil(t, err)
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

func TestLedgerArtifactManager_GetLatestObj_VerifiesRecords(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	_, _, err := manager.GetLatestObj(record.Reference{}, record.Reference{}, record.Reference{})
	assert.NotNil(t, err)

	classRef := addRecord(&ledger, &record.ClassActivateRecord{})
	objectRef := addRecord(&ledger, &record.ObjectActivateRecord{})
	wrongRef := addRecord(&ledger, &record.CodeRecord{})
	_, _, err = manager.GetLatestObj(wrongRef, classRef, objectRef)
	assert.NotNil(t, err)
	_, _, err = manager.GetLatestObj(objectRef, wrongRef, objectRef)
	assert.NotNil(t, err)
	_, _, err = manager.GetLatestObj(objectRef, classRef, wrongRef)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_GetLatestObj_VerifiesClassIsActive(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	classRef := addRecord(&ledger, &record.ClassActivateRecord{})
	classDeactivateRef := addRecord(&ledger, &record.DeactivationRecord{})
	ledger.SetClassIndex(classRef.Record, &index.ClassLifeline{LatestStateID: classDeactivateRef.Record})
	objectRef := addRecord(&ledger, &record.ObjectActivateRecord{})
	ledger.SetObjectIndex(objectRef.Record, &index.ObjectLifeline{
		LatestStateID: objectRef.Record,
		ClassID:       classRef.Record,
	})
	_, _, err := manager.GetLatestObj(objectRef, classRef, objectRef)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_GetLatestObj_ReturnsNilDescriptorsIfCurrentStateProvided(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}
	manager := LedgerArtifactManager{storer: &ledger}
	classRef := addRecord(&ledger, &record.ClassActivateRecord{})
	classAmendRef := addRecord(&ledger, &record.ClassAmendRecord{})
	ledger.SetClassIndex(classRef.Record, &index.ClassLifeline{
		LatestStateID: classAmendRef.Record,
	})
	objectRef := addRecord(&ledger, &record.ObjectActivateRecord{})
	objectAmendRef := addRecord(&ledger, &record.ObjectAmendRecord{})
	ledger.SetObjectIndex(objectRef.Record, &index.ObjectLifeline{
		LatestStateID: objectAmendRef.Record,
		ClassID:       classRef.Record,
	})
	classDesc, objDesc, err := manager.GetLatestObj(objectRef, classAmendRef, objectAmendRef)
	assert.Nil(t, err)
	assert.Nil(t, classDesc)
	assert.Nil(t, objDesc)
}

func TestLedgerArtifactManager_GetLatestObj_ReturnsCorrectDescriptors(t *testing.T) {
	ledger := LedgerMock{
		Records:       map[record.ID]record.Record{},
		ClassIndexes:  map[record.ID]*index.ClassLifeline{},
		ObjectIndexes: map[record.ID]*index.ObjectLifeline{},
	}

	manager := LedgerArtifactManager{storer: &ledger}

	classRef := addRecord(&ledger, &record.ClassActivateRecord{DefaultMemory: record.Memory{1}})
	classRec, _ := ledger.GetRecord(classRef.Record)
	classRecCasted, _ := classRec.(*record.ClassActivateRecord)
	classAmendRef := addRecord(&ledger, &record.ClassAmendRecord{NewCode: record.Reference{Record: record.ID{2}}})
	classAmendRec, _ := ledger.GetRecord(classAmendRef.Record)
	classAmendRecCasted, _ := classAmendRec.(*record.ClassAmendRecord)
	classIndex := index.ClassLifeline{
		LatestStateID: classAmendRef.Record,
	}
	ledger.SetClassIndex(classRef.Record, &classIndex)

	objectRef := addRecord(&ledger, &record.ObjectActivateRecord{Memory: record.Memory{3}})
	objectRec, _ := ledger.GetRecord(objectRef.Record)
	objectRecCasted, _ := objectRec.(*record.ObjectActivateRecord)
	objectAmendRef := addRecord(&ledger, &record.ObjectAmendRecord{NewMemory: record.Memory{4}})
	objectAmendRec, _ := ledger.GetRecord(objectAmendRef.Record)
	objectAmendRecCasted, _ := objectAmendRec.(*record.ObjectAmendRecord)
	objectIndex := index.ObjectLifeline{
		LatestStateID: objectAmendRef.Record,
		ClassID:       classRef.Record,
	}
	ledger.SetObjectIndex(objectRef.Record, &objectIndex)

	classDesc, objectDesc, err := manager.GetLatestObj(objectRef, classRef, objectRef)
	assert.NoError(t, err)
	assert.Equal(t, *classDesc, ClassDescriptor{
		StateRef: record.Reference{
			Domain: classRef.Domain,
			Record: classAmendRef.Record,
		},

		manager:           &manager,
		fromState:         classRef,
		activateRecord:    classRecCasted,
		latestAmendRecord: classAmendRecCasted,
		lifelineIndex:     &classIndex,
	})
	assert.Equal(t, *objectDesc, ObjectDescriptor{
		StateRef: record.Reference{
			Domain: objectRef.Domain,
			Record: objectAmendRef.Record,
		},

		manager:           &manager,
		activateRecord:    objectRecCasted,
		latestAmendRecord: objectAmendRecCasted,
		lifelineIndex:     &objectIndex,
	})
}
