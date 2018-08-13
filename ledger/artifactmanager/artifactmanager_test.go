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
	"math/rand"
	"os"
	"testing"

	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/leveldb"
	"github.com/stretchr/testify/assert"
)

func genRandomRef() *record.Reference {
	return &record.Reference{Domain: record.ID{Pulse: record.PulseNum(rand.Int())}}
}

func getLedgerManager() (storage.LedgerStorer, ArtifactManager) {
	ledger, _ := leveldb.InitDB()
	manager := LedgerArtifactManager{storer: ledger}
	return ledger, &manager
}

func prepareTestArtifactManager() (storage.LedgerStorer, ArtifactManager, *record.Reference) {
	if err := leveldb.DropDB(); err != nil {
		os.Exit(1)
	}

	ledger, _ := leveldb.InitDB()
	manager := LedgerArtifactManager{storer: ledger}

	return ledger, &manager, genRandomRef()
}

func TestLedgerArtifactManager_DeployCode(t *testing.T) {
	ledger, manager, requestRef := prepareTestArtifactManager()
	codeMap := map[record.ArchType][]byte{1: {1}}
	ref, err := manager.DeployCode(*requestRef, codeMap)
	assert.NoError(t, err)
	codeRec, err := ledger.GetRecord(ref)
	assert.NoError(t, err)
	assert.Equal(t, codeRec, &record.CodeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: *requestRef,
				},
			},
		},
		TargetedCode: codeMap,
	})
}

func TestLedgerArtifactManager_ActivateClass_VerifiesRecord(t *testing.T) {
	ledger, manager, requestRef := prepareTestArtifactManager()
	_, err := manager.ActivateClass(*requestRef, record.Reference{}, record.Memory{})
	assert.NotNil(t, err)
	notCodeRef, _ := ledger.SetRecord(&record.ClassActivateRecord{})
	_, err = manager.ActivateClass(*requestRef, *notCodeRef, record.Memory{})
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_ActivateClass_CreatesCorrectRecord(t *testing.T) {
	ledger, manager, requestRef := prepareTestArtifactManager()
	memory := record.Memory{1, 2, 3}
	codeRef, _ := ledger.SetRecord(&record.CodeRecord{})
	activateRef, err := manager.ActivateClass(*requestRef, *codeRef, memory)
	assert.Nil(t, err)
	activateRec, getErr := ledger.GetRecord(activateRef)
	assert.Nil(t, getErr)
	assert.Equal(t, activateRec, &record.ClassActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: *requestRef,
				},
			},
		},
		CodeRecord:    *codeRef,
		DefaultMemory: memory,
	})
}

func TestLedgerArtifactManager_DeactivateClass_VerifiesRecord(t *testing.T) {
	ledger, manager, requestRef := prepareTestArtifactManager()
	_, err := manager.DeactivateClass(*requestRef, record.Reference{})
	assert.NotNil(t, err)

	notClassRef, _ := ledger.SetRecord(&record.CodeRecord{})
	_, err = manager.DeactivateClass(*requestRef, *notClassRef)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateClass_VerifiesClassIsActive(t *testing.T) {
	ledger, manager, requestRef := prepareTestArtifactManager()
	classRef, _ := ledger.SetRecord(&record.ClassActivateRecord{})
	deactivateRef, _ := ledger.SetRecord(&record.DeactivationRecord{})
	ledger.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *deactivateRef,
	})
	_, err := manager.DeactivateClass(*requestRef, *classRef)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateClass_CreatesCorrectRecord(t *testing.T) {
	ledger, manager, requestRef := prepareTestArtifactManager()
	classRef, _ := ledger.SetRecord(&record.ClassActivateRecord{})
	ledger.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *classRef,
	})

	deactivateRef, err := manager.DeactivateClass(*requestRef, *classRef)
	assert.NoError(t, err)
	deactivateRec, err := ledger.GetRecord(deactivateRef)
	assert.NoError(t, err)
	assert.Equal(t, deactivateRec, &record.DeactivationRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: *requestRef,
				},
			},
			HeadRecord:    *classRef,
			AmendedRecord: *classRef,
		},
	})
}

func TestLedgerArtifactManager_UpdateClass_VerifiesRecord(t *testing.T) {
	ledger, manager, requestRef := prepareTestArtifactManager()
	_, err := manager.UpdateClass(*requestRef, record.Reference{}, record.Reference{}, nil)
	assert.NotNil(t, err)
	notClassRef, _ := ledger.SetRecord(&record.CodeRecord{})
	_, err = manager.UpdateClass(*requestRef, *notClassRef, record.Reference{}, nil)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateClass_VerifiesClassIsActive(t *testing.T) {
	ledger, manager, requestRef := prepareTestArtifactManager()
	classRef, _ := ledger.SetRecord(&record.ClassActivateRecord{})
	deactivateRef, _ := ledger.SetRecord(&record.DeactivationRecord{})
	codeRef, _ := ledger.SetRecord(&record.CodeRecord{})
	ledger.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *deactivateRef,
	})
	_, err := manager.UpdateClass(*requestRef, *classRef, *codeRef, nil)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateClass_CreatesCorrectRecord(t *testing.T) {
	ledger, manager, requestRef := prepareTestArtifactManager()
	classRef, _ := ledger.SetRecord(&record.ClassActivateRecord{})
	ledger.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *classRef,
	})
	codeRef, _ := ledger.SetRecord(&record.CodeRecord{})
	migrationRef, _ := ledger.SetRecord(&record.CodeRecord{SourceCode: "test"})
	migrationRefs := []record.Reference{*migrationRef}
	updateRef, err := manager.UpdateClass(*requestRef, *classRef, *codeRef, migrationRefs)
	assert.Nil(t, err)
	updateRec, getErr := ledger.GetRecord(updateRef)
	assert.Nil(t, getErr)
	assert.Equal(t, updateRec, &record.ClassAmendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: *requestRef,
				},
			},
			HeadRecord:    *classRef,
			AmendedRecord: *classRef,
		},
		NewCode:    *codeRef,
		Migrations: migrationRefs,
	})
}

func TestLedgerArtifactManager_ActivateObj_VerifiesRecord(t *testing.T) {
	ledger, manager, requestRef := prepareTestArtifactManager()
	_, err := manager.ActivateObj(*requestRef, record.Reference{}, record.Memory{})
	assert.NotNil(t, err)
	notClassRef, _ := ledger.SetRecord(&record.ObjectActivateRecord{})
	_, err = manager.ActivateClass(*requestRef, *notClassRef, record.Memory{})
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_ActivateObj_CreatesCorrectRecord(t *testing.T) {
	ledger, manager, requestRef := prepareTestArtifactManager()
	memory := record.Memory{1, 2, 3}
	classRef, _ := ledger.SetRecord(&record.ClassActivateRecord{})
	ledger.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *classRef,
	})
	activateRef, err := manager.ActivateObj(*requestRef, *classRef, memory)
	assert.Nil(t, err)
	activateRec, err := ledger.GetRecord(activateRef)
	assert.Nil(t, err)
	assert.Equal(t, activateRec, &record.ObjectActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: *requestRef,
				},
			},
		},
		ClassActivateRecord: *classRef,
		Memory:              memory,
	})
}

func TestLedgerArtifactManager_DeactivateObj_VerifiesRecord(t *testing.T) {
	ledger, manager, requestRef := prepareTestArtifactManager()
	_, err := manager.DeactivateClass(*requestRef, record.Reference{})
	assert.NotNil(t, err)
	notObjRef, _ := ledger.SetRecord(&record.ClassActivateRecord{})
	_, err = manager.DeactivateClass(*requestRef, *notObjRef)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateObj_VerifiesObjectIsActive(t *testing.T) {
	ledger, manager, requestRef := prepareTestArtifactManager()
	objRef, _ := ledger.SetRecord(&record.ObjectActivateRecord{})
	deactivateRef, _ := ledger.SetRecord(&record.DeactivationRecord{})
	ledger.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *deactivateRef,
	})
	_, err := manager.DeactivateObj(*requestRef, *objRef)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateObj_CreatesCorrectRecord(t *testing.T) {
	ledger, manager, requestRef := prepareTestArtifactManager()
	objRef, _ := ledger.SetRecord(&record.ObjectActivateRecord{})
	ledger.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *objRef,
	})
	deactivateRef, err := manager.DeactivateObj(*requestRef, *objRef)
	assert.Nil(t, err)
	deactivateRec, err := ledger.GetRecord(deactivateRef)
	assert.Nil(t, err)
	assert.Equal(t, deactivateRec, &record.DeactivationRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: *requestRef,
				},
			},
			HeadRecord:    *objRef,
			AmendedRecord: *objRef,
		},
	})
}

func TestLedgerArtifactManager_UpdateObj_VerifiesRecord(t *testing.T) {
	ledger, manager, requestRef := prepareTestArtifactManager()
	_, err := manager.UpdateObj(*requestRef, record.Reference{}, nil)
	assert.NotNil(t, err)
	notObjRef, _ := ledger.SetRecord(&record.CodeRecord{})
	_, err = manager.UpdateObj(*requestRef, *notObjRef, nil)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateObj_VerifiesObjectIsActive(t *testing.T) {
	ledger, manager, requestRef := prepareTestArtifactManager()
	objRef, _ := ledger.SetRecord(&record.ObjectActivateRecord{})
	deactivateRef, _ := ledger.SetRecord(&record.DeactivationRecord{})
	ledger.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *deactivateRef,
	})
	_, err := manager.UpdateObj(*requestRef, *objRef, nil)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateObj_CreatesCorrectRecord(t *testing.T) {
	ledger, manager, requestRef := prepareTestArtifactManager()
	objRef, _ := ledger.SetRecord(&record.ObjectActivateRecord{})
	ledger.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *objRef,
	})
	memory := record.Memory{1, 2, 3}
	updateRef, err := manager.UpdateObj(*requestRef, *objRef, memory)
	assert.Nil(t, err)
	updateRec, err := ledger.GetRecord(updateRef)
	assert.Nil(t, err)
	assert.Equal(t, updateRec, &record.ObjectAmendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: *requestRef,
				},
			},
			HeadRecord:    *objRef,
			AmendedRecord: *objRef,
		},
		NewMemory: memory,
	})
}

func TestLedgerArtifactManager_AppendObjDelegate_VerifiesRecord(t *testing.T) {
	ledger, manager, requestRef := prepareTestArtifactManager()
	_, err := manager.AppendObjDelegate(*requestRef, record.Reference{}, nil)
	assert.NotNil(t, err)
	notObjRef, _ := ledger.SetRecord(&record.CodeRecord{})
	_, err = manager.AppendObjDelegate(*requestRef, *notObjRef, nil)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_AppendObjDelegate_VerifiesClassIsActive(t *testing.T) {
	ledger, manager, requestRef := prepareTestArtifactManager()
	objRef, _ := ledger.SetRecord(&record.ObjectActivateRecord{})
	deactivateRef, _ := ledger.SetRecord(&record.DeactivationRecord{})
	ledger.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *deactivateRef,
	})
	_, err := manager.AppendObjDelegate(*requestRef, *objRef, nil)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_AppendObjDelegate_CreatesCorrectRecord(t *testing.T) {
	ledger, manager, requestRef := prepareTestArtifactManager()
	objRef, _ := ledger.SetRecord(&record.ObjectActivateRecord{})
	ledger.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *objRef,
	})
	memory := record.Memory{1, 2, 3}
	appendRef, err := manager.AppendObjDelegate(*requestRef, *objRef, memory)
	assert.Nil(t, err)
	appendRec, err := ledger.GetRecord(appendRef)
	objIndex, ok := ledger.GetObjectIndex(objRef)
	assert.True(t, ok)
	assert.Equal(t, objIndex.AppendRefs, []record.Reference{*appendRef})
	assert.Nil(t, err)
	assert.Equal(t, appendRec, &record.ObjectAppendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					RequestRecord: *requestRef,
				},
			},
			HeadRecord:    *objRef,
			AmendedRecord: *objRef,
		},
		AppendMemory: memory,
	})
}

func TestLedgerArtifactManager_GetLatestObj_VerifiesRecords(t *testing.T) {
	ledger, manager, _ := prepareTestArtifactManager()
	_, _, err := manager.GetLatestObj(record.Reference{}, record.Reference{}, record.Reference{})
	assert.NotNil(t, err)

	classRef, _ := ledger.SetRecord(&record.ClassActivateRecord{})
	objectRef, _ := ledger.SetRecord(&record.ObjectActivateRecord{})
	wrongRef, _ := ledger.SetRecord(&record.CodeRecord{})
	_, _, err = manager.GetLatestObj(*wrongRef, *classRef, *objectRef)
	assert.NotNil(t, err)
	_, _, err = manager.GetLatestObj(*objectRef, *wrongRef, *objectRef)
	assert.NotNil(t, err)
	_, _, err = manager.GetLatestObj(*objectRef, *classRef, *wrongRef)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_GetLatestObj_VerifiesClassIsActive(t *testing.T) {
	ledger, manager, _ := prepareTestArtifactManager()
	classRef, _ := ledger.SetRecord(&record.ClassActivateRecord{})
	classDeactivateRef, _ := ledger.SetRecord(&record.DeactivationRecord{})
	ledger.SetClassIndex(classRef, &index.ClassLifeline{LatestStateRef: *classDeactivateRef})
	objectRef, _ := ledger.SetRecord(&record.ObjectActivateRecord{})
	ledger.SetObjectIndex(objectRef, &index.ObjectLifeline{
		LatestStateRef: *objectRef,
		ClassRef:       *classRef,
	})
	_, _, err := manager.GetLatestObj(*objectRef, *classRef, *objectRef)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_GetLatestObj_ReturnsNilDescriptorsIfCurrentStateProvided(t *testing.T) {
	ledger, manager, _ := prepareTestArtifactManager()
	classRef, _ := ledger.SetRecord(&record.ClassActivateRecord{})
	classAmendRef, _ := ledger.SetRecord(&record.ClassAmendRecord{})
	ledger.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *classAmendRef,
	})
	objectRef, _ := ledger.SetRecord(&record.ObjectActivateRecord{})
	objectAmendRef, _ := ledger.SetRecord(&record.ObjectAmendRecord{})
	ledger.SetObjectIndex(objectRef, &index.ObjectLifeline{
		LatestStateRef: *objectAmendRef,
		ClassRef:       *classRef,
	})
	classDesc, objDesc, err := manager.GetLatestObj(*objectRef, *classAmendRef, *objectAmendRef)
	assert.Nil(t, err)
	assert.Nil(t, classDesc)
	assert.Nil(t, objDesc)
}

func TestLedgerArtifactManager_GetLatestObj_ReturnsCorrectDescriptors(t *testing.T) {
	ledger, manager, _ := prepareTestArtifactManager()
	ledgerManager, _ := manager.(*LedgerArtifactManager)

	classRef, _ := ledger.SetRecord(&record.ClassActivateRecord{DefaultMemory: record.Memory{1}})
	classRec, _ := ledger.GetRecord(classRef)
	classRecCasted, _ := classRec.(*record.ClassActivateRecord)
	classAmendRef, _ := ledger.SetRecord(&record.ClassAmendRecord{NewCode: *genRandomRef()})
	classAmendRec, _ := ledger.GetRecord(classAmendRef)
	classAmendRecCasted, _ := classAmendRec.(*record.ClassAmendRecord)
	classIndex := index.ClassLifeline{
		LatestStateRef: *classAmendRef,
	}
	ledger.SetClassIndex(classRef, &classIndex)

	objectRef, _ := ledger.SetRecord(&record.ObjectActivateRecord{Memory: record.Memory{3}})
	objectRec, _ := ledger.GetRecord(objectRef)
	objectRecCasted, _ := objectRec.(*record.ObjectActivateRecord)
	objectAmendRef, _ := ledger.SetRecord(&record.ObjectAmendRecord{NewMemory: record.Memory{4}})
	objectAmendRec, _ := ledger.GetRecord(objectAmendRef)
	objectAmendRecCasted, _ := objectAmendRec.(*record.ObjectAmendRecord)
	objectIndex := index.ObjectLifeline{
		LatestStateRef: *objectAmendRef,
		ClassRef:       *classRef,
	}
	ledger.SetObjectIndex(objectRef, &objectIndex)

	classDesc, objectDesc, err := manager.GetLatestObj(*objectRef, *classRef, *objectRef)
	assert.NoError(t, err)
	assert.Equal(t, *classDesc, ClassDescriptor{
		StateRef: record.Reference{
			Domain: classRef.Domain,
			Record: classAmendRef.Record,
		},

		manager:           ledgerManager,
		fromState:         *classRef,
		activateRecord:    classRecCasted,
		latestAmendRecord: classAmendRecCasted,
		lifelineIndex:     &classIndex,
	})
	assert.Equal(t, *objectDesc, ObjectDescriptor{
		StateRef: record.Reference{
			Domain: objectRef.Domain,
			Record: objectAmendRef.Record,
		},

		manager:           ledgerManager,
		activateRecord:    objectRecCasted,
		latestAmendRecord: objectAmendRecCasted,
		lifelineIndex:     &objectIndex,
	})
}

func TestLedgerArtifactManager_GetExactObj_VerifiesRecords(t *testing.T) {
	_, manager, _ := prepareTestArtifactManager()
	_, _, err := manager.GetExactObj(record.Reference{}, record.Reference{})
	assert.Error(t, err)
}

func TestLedgerArtifactManager_GetExactObj_ReturnsCorrectData(t *testing.T) {
	ledger, manager, _ := prepareTestArtifactManager()
	manager.SetArchPref([]record.ArchType{1})

	codeRec := record.CodeRecord{TargetedCode: map[record.ArchType][]byte{
		1: {1},
	}}
	codeRef, _ := ledger.SetRecord(&codeRec)
	classRef, _ := ledger.SetRecord(&record.ClassActivateRecord{})
	classAmendRef, _ := ledger.SetRecord(&record.ClassAmendRecord{
		NewCode: *codeRef,
		AmendRecord: record.AmendRecord{
			HeadRecord: *classRef,
		},
	})

	memoryRec := record.Memory{4}
	objectRef, _ := ledger.SetRecord(&record.ObjectActivateRecord{Memory: record.Memory{3}})
	objectAmendRef, _ := ledger.SetRecord(&record.ObjectAmendRecord{
		NewMemory: memoryRec,
		AmendRecord: record.AmendRecord{
			HeadRecord: *objectRef,
		},
	})

	_, _, err := manager.GetExactObj(*classAmendRef, *objectAmendRef)
	assert.Error(t, err)

	ledger.SetObjectIndex(objectRef, &index.ObjectLifeline{
		LatestStateRef: *objectAmendRef,
		ClassRef:       *classRef,
	})

	code, memory, err := manager.GetExactObj(*classAmendRef, *objectAmendRef)
	assert.NoError(t, err)
	assert.Equal(t, []byte{1}, code)
	assert.Equal(t, memoryRec, memory)
}
