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
	"bytes"
	"os"
	"testing"

	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/stretchr/testify/assert"
)

// Mock to imitate storage
type LedgerMock struct {
	Records map[record.ID]record.Record
	Indexes map[record.ID]*index.Lifeline
}

func (mock *LedgerMock) GetRecord(k record.Key) (record.Record, bool) {
	rec, ok := mock.Records[record.Key2ID(k)]
	return rec, ok
}

func (mock *LedgerMock) AddRecord(rec record.Record) (record.Reference, error) {
	buf := bytes.Buffer{}
	// TODO: implement properly after merge with INS-1-storage-records-persist
	// rec.WriteHash(&buf)
	buf.Write([]byte{1})
	var id record.ID
	copy(buf.Bytes()[0:record.IDSize], id[:])
	mock.Records[record.ID{}] = rec
	return record.Reference{
		Domain: record.ID{},
		Record: id,
	}, nil
}

func (mock *LedgerMock) GetIndex(id record.ID) (*index.Lifeline, bool) {
	idx, ok := mock.Indexes[id]
	return idx, ok
}

func (mock *LedgerMock) SetIndex(id record.ID, idx *index.Lifeline) error {
	mock.Indexes[id] = idx
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
	ledger := LedgerMock{map[record.ID]record.Record{}, map[record.ID]*index.Lifeline{}}
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
	ledger := LedgerMock{map[record.ID]record.Record{}, map[record.ID]*index.Lifeline{}}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	_, err := manager.ActivateClass(requestRef, record.Reference{}, record.Memory{})
	assert.NotNil(t, err)
	notCodeRef, _ := ledger.AddRecord(&record.ClassActivateRecord{})
	_, err = manager.ActivateClass(requestRef, notCodeRef, record.Memory{})
	assert.NotNil(t, err)
}

func TestActivateClassCreatesActivateRecord(t *testing.T) {
	ledger := LedgerMock{map[record.ID]record.Record{}, map[record.ID]*index.Lifeline{}}
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
	ledger := LedgerMock{map[record.ID]record.Record{}, map[record.ID]*index.Lifeline{}}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	_, err := manager.DeactivateClass(requestRef, record.Reference{})
	assert.NotNil(t, err)
	notClassRef, _ := ledger.AddRecord(&record.CodeRecord{})
	_, err = manager.DeactivateClass(requestRef, notClassRef)
	assert.NotNil(t, err)
}

func TestDeactivateClassVerifiesClassIsActive(t *testing.T) {
	ledger := LedgerMock{map[record.ID]record.Record{}, map[record.ID]*index.Lifeline{}}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	classRef, _ := ledger.AddRecord(&record.ClassActivateRecord{})
	deactivateRef, _ := ledger.AddRecord(&record.DeactivationRecord{})
	ledger.SetIndex(classRef.Record, &index.Lifeline{
		LatestStateID: deactivateRef.Record,
	})
	_, err := manager.DeactivateClass(requestRef, classRef)
	assert.NotNil(t, err)
}

func TestDeactivateClassCreatesDeactivateRecord(t *testing.T) {
	ledger := LedgerMock{map[record.ID]record.Record{}, map[record.ID]*index.Lifeline{}}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	classRef, _ := ledger.AddRecord(&record.ClassActivateRecord{})
	ledger.SetIndex(classRef.Record, &index.Lifeline{
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
	ledger := LedgerMock{map[record.ID]record.Record{}, map[record.ID]*index.Lifeline{}}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	_, err := manager.UpdateClass(requestRef, record.Reference{}, []record.MemoryMigrationCode{})
	assert.NotNil(t, err)
	notClassRef, _ := ledger.AddRecord(&record.CodeRecord{})
	_, err = manager.UpdateClass(requestRef, notClassRef, []record.MemoryMigrationCode{})
	assert.NotNil(t, err)
}

func TestUpdateClassVerifiesClassIsActive(t *testing.T) {
	ledger := LedgerMock{map[record.ID]record.Record{}, map[record.ID]*index.Lifeline{}}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	migrations := []record.MemoryMigrationCode{
		record.MemoryMigrationCode{
			MigrationCodeRecord: record.Reference{Record: record.ID{1, 2, 3}},
		},
	}
	classRef, _ := ledger.AddRecord(&record.ClassActivateRecord{})
	deactivateRef, _ := ledger.AddRecord(&record.DeactivationRecord{})
	ledger.SetIndex(classRef.Record, &index.Lifeline{
		LatestStateID: deactivateRef.Record,
	})
	_, err := manager.UpdateClass(requestRef, classRef, migrations)
	assert.NotNil(t, err)
}

func TestUpdateClassCreatesAmendRecord(t *testing.T) {
	ledger := LedgerMock{map[record.ID]record.Record{}, map[record.ID]*index.Lifeline{}}
	manager := LedgerArtifactManager{storer: &ledger}
	requestRef := record.Reference{Domain: record.ID{1}, Record: record.ID{2}}
	migrations := []record.MemoryMigrationCode{
		record.MemoryMigrationCode{
			MigrationCodeRecord: record.Reference{Record: record.ID{1, 2, 3}},
		},
	}
	classRef, _ := ledger.AddRecord(&record.ClassActivateRecord{})
	ledger.SetIndex(classRef.Record, &index.Lifeline{
		LatestStateID: classRef.Record,
	})
	updateRef, err := manager.UpdateClass(requestRef, classRef, migrations)
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
