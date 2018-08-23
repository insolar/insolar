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
	"testing"

	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/leveldb/leveltestutils"
	"github.com/stretchr/testify/assert"
)

func genRandomRef() *record.Reference {
	return &record.Reference{Domain: record.ID{Pulse: record.PulseNum(rand.Int())}}
}

type preparedAMTestData struct {
	ledger     storage.LedgerStorer
	manager    ArtifactManager
	domainRef  *record.Reference
	requestRef *record.Reference
}

func prepareAMTestData(t *testing.T) (preparedAMTestData, func()) {
	ledger, cleaner := leveltestutils.TmpDB(t, "")

	return preparedAMTestData{
		ledger:     ledger,
		manager:    &LedgerArtifactManager{storer: ledger},
		domainRef:  genRandomRef(),
		requestRef: genRandomRef(),
	}, cleaner
}

func TestLedgerArtifactManager_DeployCode(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	codeMap := map[record.ArchType][]byte{1: {1}}
	ref, err := td.manager.DeployCode(*td.domainRef, *td.requestRef, codeMap)
	assert.NoError(t, err)
	codeRec, err := td.ledger.GetRecord(ref)
	assert.NoError(t, err)
	assert.Equal(t, codeRec, &record.CodeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  *td.domainRef,
					RequestRecord: *td.requestRef,
				},
			},
		},
		TargetedCode: codeMap,
	})
}

func TestLedgerArtifactManager_ActivateClass_VerifiesRecord(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	_, err := td.manager.ActivateClass(
		*td.domainRef, *td.requestRef, record.Reference{}, record.Memory{})
	assert.NotNil(t, err)
	notCodeRef, _ := td.ledger.SetRecord(&record.ClassActivateRecord{})
	_, err = td.manager.ActivateClass(*td.domainRef, *td.requestRef, *notCodeRef, record.Memory{})
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_ActivateClass_CreatesCorrectRecord(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	memory := record.Memory{1, 2, 3}
	codeRef, _ := td.ledger.SetRecord(&record.CodeRecord{})
	activateRef, err := td.manager.ActivateClass(*td.domainRef, *td.requestRef, *codeRef, memory)
	assert.Nil(t, err)
	activateRec, getErr := td.ledger.GetRecord(activateRef)
	assert.Nil(t, getErr)
	assert.Equal(t, activateRec, &record.ClassActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  *td.domainRef,
					RequestRecord: *td.requestRef,
				},
			},
		},
		CodeRecord:    *codeRef,
		DefaultMemory: memory,
	})
}

func TestLedgerArtifactManager_DeactivateClass_VerifiesRecord(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	_, err := td.manager.DeactivateClass(*td.domainRef, *td.requestRef, record.Reference{})
	assert.NotNil(t, err)

	notClassRef, _ := td.ledger.SetRecord(&record.CodeRecord{})
	_, err = td.manager.DeactivateClass(*td.domainRef, *td.requestRef, *notClassRef)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateClass_VerifiesClassIsActive(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	classRef, _ := td.ledger.SetRecord(&record.ClassActivateRecord{})
	deactivateRef, _ := td.ledger.SetRecord(&record.DeactivationRecord{})
	td.ledger.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *deactivateRef,
	})
	_, err := td.manager.DeactivateClass(*td.domainRef, *td.requestRef, *classRef)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateClass_CreatesCorrectRecord(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	classRef, _ := td.ledger.SetRecord(&record.ClassActivateRecord{})
	td.ledger.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *classRef,
	})

	deactivateRef, err := td.manager.DeactivateClass(*td.domainRef, *td.requestRef, *classRef)
	assert.NoError(t, err)
	deactivateRec, err := td.ledger.GetRecord(deactivateRef)
	assert.NoError(t, err)
	assert.Equal(t, deactivateRec, &record.DeactivationRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  *td.domainRef,
					RequestRecord: *td.requestRef,
				},
			},
			HeadRecord:    *classRef,
			AmendedRecord: *classRef,
		},
	})
}

func TestLedgerArtifactManager_UpdateClass_VerifiesRecord(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	_, err := td.manager.UpdateClass(
		*td.domainRef, *td.requestRef, record.Reference{}, record.Reference{}, nil)
	assert.NotNil(t, err)
	notClassRef, _ := td.ledger.SetRecord(&record.CodeRecord{})
	_, err = td.manager.UpdateClass(
		*td.domainRef, *td.requestRef, *notClassRef, record.Reference{}, nil)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateClass_VerifiesClassIsActive(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	classRef, _ := td.ledger.SetRecord(&record.ClassActivateRecord{})
	deactivateRef, _ := td.ledger.SetRecord(&record.DeactivationRecord{})
	codeRef, _ := td.ledger.SetRecord(&record.CodeRecord{})
	td.ledger.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *deactivateRef,
	})
	_, err := td.manager.UpdateClass(
		*td.domainRef, *td.requestRef, *classRef, *codeRef, nil)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateClass_CreatesCorrectRecord(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	classRef, _ := td.ledger.SetRecord(&record.ClassActivateRecord{})
	td.ledger.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *classRef,
	})
	codeRef, _ := td.ledger.SetRecord(&record.CodeRecord{})
	migrationRef, _ := td.ledger.SetRecord(&record.CodeRecord{SourceCode: "test"})
	migrationRefs := []record.Reference{*migrationRef}
	updateRef, err := td.manager.UpdateClass(
		*td.domainRef, *td.requestRef, *classRef, *codeRef, migrationRefs)
	assert.Nil(t, err)
	updateRec, getErr := td.ledger.GetRecord(updateRef)
	assert.Nil(t, getErr)
	assert.Equal(t, updateRec, &record.ClassAmendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  *td.domainRef,
					RequestRecord: *td.requestRef,
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
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	_, err := td.manager.ActivateObj(
		*td.domainRef, *td.requestRef, record.Reference{}, record.Memory{})
	assert.NotNil(t, err)
	notClassRef, _ := td.ledger.SetRecord(&record.ObjectActivateRecord{})
	_, err = td.manager.ActivateClass(
		*td.domainRef, *td.requestRef, *notClassRef, record.Memory{})
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_ActivateObj_CreatesCorrectRecord(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	memory := record.Memory{1, 2, 3}
	classRef, _ := td.ledger.SetRecord(&record.ClassActivateRecord{})
	td.ledger.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *classRef,
	})
	activateRef, err := td.manager.ActivateObj(
		*td.domainRef, *td.requestRef, *classRef, memory)
	assert.Nil(t, err)
	activateRec, err := td.ledger.GetRecord(activateRef)
	assert.Nil(t, err)
	assert.Equal(t, activateRec, &record.ObjectActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  *td.domainRef,
					RequestRecord: *td.requestRef,
				},
			},
		},
		ClassActivateRecord: *classRef,
		Memory:              memory,
	})
}

func TestLedgerArtifactManager_DeactivateObj_VerifiesRecord(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	_, err := td.manager.DeactivateClass(
		*td.domainRef, *td.requestRef, record.Reference{})
	assert.NotNil(t, err)
	notObjRef, _ := td.ledger.SetRecord(&record.ClassActivateRecord{})
	_, err = td.manager.DeactivateClass(*td.domainRef, *td.requestRef, *notObjRef)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateObj_VerifiesObjectIsActive(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	objRef, _ := td.ledger.SetRecord(&record.ObjectActivateRecord{})
	deactivateRef, _ := td.ledger.SetRecord(&record.DeactivationRecord{})
	td.ledger.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *deactivateRef,
	})
	_, err := td.manager.DeactivateObj(*td.domainRef, *td.requestRef, *objRef)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateObj_CreatesCorrectRecord(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	objRef, _ := td.ledger.SetRecord(&record.ObjectActivateRecord{})
	td.ledger.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *objRef,
	})
	deactivateRef, err := td.manager.DeactivateObj(*td.domainRef, *td.requestRef, *objRef)
	assert.Nil(t, err)
	deactivateRec, err := td.ledger.GetRecord(deactivateRef)
	assert.Nil(t, err)
	assert.Equal(t, deactivateRec, &record.DeactivationRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  *td.domainRef,
					RequestRecord: *td.requestRef,
				},
			},
			HeadRecord:    *objRef,
			AmendedRecord: *objRef,
		},
	})
}

func TestLedgerArtifactManager_UpdateObj_VerifiesRecord(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	_, err := td.manager.UpdateObj(
		*td.domainRef, *td.requestRef, record.Reference{}, nil)
	assert.NotNil(t, err)
	notObjRef, _ := td.ledger.SetRecord(&record.CodeRecord{})
	_, err = td.manager.UpdateObj(
		*td.domainRef, *td.requestRef, *notObjRef, nil)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateObj_VerifiesObjectIsActive(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	objRef, _ := td.ledger.SetRecord(&record.ObjectActivateRecord{})
	deactivateRef, _ := td.ledger.SetRecord(&record.DeactivationRecord{})
	td.ledger.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *deactivateRef,
	})
	_, err := td.manager.UpdateObj(
		*td.domainRef, *td.requestRef, *objRef, nil)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateObj_CreatesCorrectRecord(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	objRef, _ := td.ledger.SetRecord(&record.ObjectActivateRecord{})
	td.ledger.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *objRef,
	})
	memory := record.Memory{1, 2, 3}
	updateRef, err := td.manager.UpdateObj(
		*td.domainRef, *td.requestRef, *objRef, memory)
	assert.Nil(t, err)
	updateRec, err := td.ledger.GetRecord(updateRef)
	assert.Nil(t, err)
	assert.Equal(t, updateRec, &record.ObjectAmendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  *td.domainRef,
					RequestRecord: *td.requestRef,
				},
			},
			HeadRecord:    *objRef,
			AmendedRecord: *objRef,
		},
		NewMemory: memory,
	})
}

func TestLedgerArtifactManager_AppendObjDelegate_VerifiesRecord(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	_, err := td.manager.AppendObjDelegate(
		*td.domainRef, *td.requestRef, record.Reference{}, nil)
	assert.NotNil(t, err)
	notObjRef, _ := td.ledger.SetRecord(&record.CodeRecord{})
	_, err = td.manager.AppendObjDelegate(
		*td.domainRef, *td.requestRef, *notObjRef, nil)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_AppendObjDelegate_VerifiesClassIsActive(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	objRef, _ := td.ledger.SetRecord(&record.ObjectActivateRecord{})
	deactivateRef, _ := td.ledger.SetRecord(&record.DeactivationRecord{})
	td.ledger.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *deactivateRef,
	})
	_, err := td.manager.AppendObjDelegate(
		*td.domainRef, *td.requestRef, *objRef, nil)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_AppendObjDelegate_CreatesCorrectRecord(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	objRef, _ := td.ledger.SetRecord(&record.ObjectActivateRecord{})
	td.ledger.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *objRef,
	})
	memory := record.Memory{1, 2, 3}
	appendRef, err := td.manager.AppendObjDelegate(
		*td.domainRef, *td.requestRef, *objRef, memory)
	assert.Nil(t, err)
	appendRec, _ := td.ledger.GetRecord(appendRef)
	objIndex, err := td.ledger.GetObjectIndex(objRef)
	assert.NoError(t, err)
	assert.Equal(t, objIndex.AppendRefs, []record.Reference{*appendRef})
	assert.Nil(t, err)
	assert.Equal(t, appendRec, &record.ObjectAppendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  *td.domainRef,
					RequestRecord: *td.requestRef,
				},
			},
			HeadRecord:    *objRef,
			AmendedRecord: *objRef,
		},
		AppendMemory: memory,
	})
}

func TestLedgerArtifactManager_GetLatestObj_VerifiesRecords(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	_, _, err := td.manager.GetLatestObj(
		record.Reference{}, record.Reference{}, record.Reference{})
	assert.NotNil(t, err)

	classRef, _ := td.ledger.SetRecord(&record.ClassActivateRecord{})
	objectRef, _ := td.ledger.SetRecord(&record.ObjectActivateRecord{})
	wrongRef, _ := td.ledger.SetRecord(&record.CodeRecord{})
	_, _, err = td.manager.GetLatestObj(*wrongRef, *classRef, *objectRef)
	assert.NotNil(t, err)
	_, _, err = td.manager.GetLatestObj(*objectRef, *wrongRef, *objectRef)
	assert.NotNil(t, err)
	_, _, err = td.manager.GetLatestObj(*objectRef, *classRef, *wrongRef)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_GetLatestObj_VerifiesClassIsActive(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	classRef, _ := td.ledger.SetRecord(&record.ClassActivateRecord{})
	classDeactivateRef, _ := td.ledger.SetRecord(&record.DeactivationRecord{})
	td.ledger.SetClassIndex(classRef, &index.ClassLifeline{LatestStateRef: *classDeactivateRef})
	objectRef, _ := td.ledger.SetRecord(&record.ObjectActivateRecord{})
	td.ledger.SetObjectIndex(objectRef, &index.ObjectLifeline{
		LatestStateRef: *objectRef,
		ClassRef:       *classRef,
	})
	_, _, err := td.manager.GetLatestObj(*objectRef, *classRef, *objectRef)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_GetLatestObj_ReturnsNilDescriptorsIfCurrentStateProvided(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	classRef, _ := td.ledger.SetRecord(&record.ClassActivateRecord{})
	classAmendRef, _ := td.ledger.SetRecord(&record.ClassAmendRecord{})
	td.ledger.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *classAmendRef,
	})
	objectRef, _ := td.ledger.SetRecord(&record.ObjectActivateRecord{})
	objectAmendRef, _ := td.ledger.SetRecord(&record.ObjectAmendRecord{})
	td.ledger.SetObjectIndex(objectRef, &index.ObjectLifeline{
		LatestStateRef: *objectAmendRef,
		ClassRef:       *classRef,
	})
	classDesc, objDesc, err := td.manager.GetLatestObj(
		*objectRef, *classAmendRef, *objectAmendRef)
	assert.Nil(t, err)
	assert.Nil(t, classDesc)
	assert.Nil(t, objDesc)
}

func TestLedgerArtifactManager_GetLatestObj_ReturnsCorrectDescriptors(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	ledgerManager, _ := td.manager.(*LedgerArtifactManager)

	classRef, _ := td.ledger.SetRecord(&record.ClassActivateRecord{DefaultMemory: record.Memory{1}})
	classRec, _ := td.ledger.GetRecord(classRef)
	classRecCasted, _ := classRec.(*record.ClassActivateRecord)
	classAmendRef, _ := td.ledger.SetRecord(&record.ClassAmendRecord{NewCode: *genRandomRef()})
	classAmendRec, _ := td.ledger.GetRecord(classAmendRef)
	classAmendRecCasted, _ := classAmendRec.(*record.ClassAmendRecord)
	classIndex := index.ClassLifeline{
		LatestStateRef: *classAmendRef,
	}
	td.ledger.SetClassIndex(classRef, &classIndex)

	objectRef, _ := td.ledger.SetRecord(&record.ObjectActivateRecord{Memory: record.Memory{3}})
	objectRec, _ := td.ledger.GetRecord(objectRef)
	objectRecCasted, _ := objectRec.(*record.ObjectActivateRecord)
	objectAmendRef, _ := td.ledger.SetRecord(&record.ObjectAmendRecord{NewMemory: record.Memory{4}})
	objectAmendRec, _ := td.ledger.GetRecord(objectAmendRef)
	objectAmendRecCasted, _ := objectAmendRec.(*record.ObjectAmendRecord)
	objectIndex := index.ObjectLifeline{
		LatestStateRef: *objectAmendRef,
		ClassRef:       *classRef,
	}
	td.ledger.SetObjectIndex(objectRef, &objectIndex)

	classDesc, objectDesc, err := td.manager.GetLatestObj(
		*objectRef, *classRef, *objectRef)
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
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	_, _, err := td.manager.GetExactObj(record.Reference{}, record.Reference{})
	assert.Error(t, err)
}

func TestLedgerArtifactManager_GetExactObj_ReturnsCorrectData(t *testing.T) {
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	td.manager.SetArchPref([]record.ArchType{1})

	codeRec := record.CodeRecord{TargetedCode: map[record.ArchType][]byte{
		1: {1},
	}}
	codeRef, _ := td.ledger.SetRecord(&codeRec)
	classRef, _ := td.ledger.SetRecord(&record.ClassActivateRecord{})
	classAmendRef, _ := td.ledger.SetRecord(&record.ClassAmendRecord{
		NewCode: *codeRef,
		AmendRecord: record.AmendRecord{
			HeadRecord: *classRef,
		},
	})

	memoryRec := record.Memory{4}
	objectRef, _ := td.ledger.SetRecord(&record.ObjectActivateRecord{Memory: record.Memory{3}})
	objectAmendRef, _ := td.ledger.SetRecord(&record.ObjectAmendRecord{
		NewMemory: memoryRec,
		AmendRecord: record.AmendRecord{
			HeadRecord: *objectRef,
		},
	})

	_, _, err := td.manager.GetExactObj(*classAmendRef, *objectAmendRef)
	assert.Error(t, err)

	td.ledger.SetObjectIndex(objectRef, &index.ObjectLifeline{
		LatestStateRef: *objectAmendRef,
		ClassRef:       *classRef,
	})

	code, memory, err := td.manager.GetExactObj(*classAmendRef, *objectAmendRef)
	assert.NoError(t, err)
	assert.Equal(t, []byte{1}, code)
	assert.Equal(t, memoryRec, memory)
}
