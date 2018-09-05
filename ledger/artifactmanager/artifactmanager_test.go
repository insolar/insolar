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

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"

	"github.com/insolar/insolar/ledger/storage/storagetest"
)

func genRandomRef() *record.Reference {
	var coreRef core.RecordRef
	coreRef[record.RefIDSize-1] = byte(rand.Int() % 256)
	ref := record.Core2Reference(coreRef)
	return &ref
}

type preparedAMTestData struct {
	ledger     *storage.Store
	manager    core.ArtifactManager
	domainRef  *record.Reference
	requestRef *record.Reference
}

func prepareAMTestData(t *testing.T) (preparedAMTestData, func()) {
	ledger, cleaner := storagetest.TmpStore(t, "")

	return preparedAMTestData{
		ledger:     ledger,
		manager:    &LedgerArtifactManager{store: ledger},
		domainRef:  genRandomRef(),
		requestRef: genRandomRef(),
	}, cleaner
}

func TestLedgerArtifactManager_DeclareType(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	typeDec := []byte{1, 2, 3}
	coreRef, err := td.manager.DeclareType(*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), typeDec)
	assert.NoError(t, err)
	ref := record.Core2Reference(*coreRef)
	typeRec, err := td.ledger.GetRecord(&ref)
	assert.NoError(t, err)
	assert.Equal(t, typeRec, &record.TypeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  *td.domainRef,
					RequestRecord: *td.requestRef,
				},
			},
		},
		TypeDeclaration: typeDec,
	})
}

func TestLedgerArtifactManager_DeployCode_VerifiesRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	codeMap := map[core.MachineType][]byte{1: {1}}
	_, err := td.manager.DeployCode(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), []core.RecordRef{*genRandomRef().CoreRef()}, codeMap,
	)
	assert.Error(t, err)
	notTypeRef, _ := td.ledger.SetRecord(&record.CodeRecord{})
	_, err = td.manager.DeployCode(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), []core.RecordRef{*notTypeRef.CoreRef()}, codeMap,
	)
	assert.Error(t, err)
}

func TestLedgerArtifactManager_DeployCode_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	codeMap := map[core.MachineType][]byte{1: {1}}
	typeRef, _ := td.ledger.SetRecord(&record.TypeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	coreRef, err := td.manager.DeployCode(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), []core.RecordRef{*typeRef.CoreRef()}, codeMap,
	)
	assert.NoError(t, err)
	ref := record.Core2Reference(*coreRef)
	codeRec, err := td.ledger.GetRecord(&ref)
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
		Types:        []record.Reference{*typeRef},
		TargetedCode: codeMap,
	})
}

func TestLedgerArtifactManager_ActivateClass_VerifiesRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	_, err := td.manager.ActivateClass(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRandomRef().CoreRef(), []byte{})
	assert.NotNil(t, err)
	notCodeRef, _ := td.ledger.SetRecord(&record.ClassActivateRecord{})
	_, err = td.manager.ActivateClass(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *notCodeRef.CoreRef(), []byte{},
	)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_ActivateClass_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	memory := []byte{1, 2, 3}
	codeRef, _ := td.ledger.SetRecord(&record.CodeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	activateCoreRef, err := td.manager.ActivateClass(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *codeRef.CoreRef(), memory,
	)
	activateRef := record.Core2Reference(*activateCoreRef)
	assert.Nil(t, err)
	activateRec, getErr := td.ledger.GetRecord(&activateRef)
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
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	_, err := td.manager.DeactivateClass(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRandomRef().CoreRef(),
	)
	assert.NotNil(t, err)

	notClassRef, _ := td.ledger.SetRecord(&record.CodeRecord{})
	_, err = td.manager.DeactivateClass(*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *notClassRef.CoreRef())
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateClass_VerifiesClassIsActive(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	classRef, _ := td.ledger.SetRecord(&record.ClassActivateRecord{})
	deactivateRef, _ := td.ledger.SetRecord(&record.DeactivationRecord{})
	td.ledger.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *deactivateRef,
	})
	_, err := td.manager.DeactivateClass(*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *classRef.CoreRef())
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateClass_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	classRef, _ := td.ledger.SetRecord(&record.ClassActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.ledger.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *classRef,
	})

	deactivateCoreRef, err := td.manager.DeactivateClass(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *classRef.CoreRef(),
	)
	assert.NoError(t, err)
	deactivateRef := record.Core2Reference(*deactivateCoreRef)
	deactivateRec, err := td.ledger.GetRecord(&deactivateRef)
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
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	_, err := td.manager.UpdateClass(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRandomRef().CoreRef(), *genRandomRef().CoreRef(), nil,
	)
	assert.NotNil(t, err)
	notClassRef, _ := td.ledger.SetRecord(&record.CodeRecord{})
	_, err = td.manager.UpdateClass(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *notClassRef.CoreRef(), *genRandomRef().CoreRef(), nil,
	)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateClass_VerifiesClassIsActive(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	classRef, _ := td.ledger.SetRecord(&record.ClassActivateRecord{})
	deactivateRef, _ := td.ledger.SetRecord(&record.DeactivationRecord{})
	codeRef, _ := td.ledger.SetRecord(&record.CodeRecord{})
	td.ledger.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *deactivateRef,
	})
	_, err := td.manager.UpdateClass(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *classRef.CoreRef(), *codeRef.CoreRef(), nil)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateClass_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	classRef, _ := td.ledger.SetRecord(&record.ClassActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.ledger.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *classRef,
	})
	codeRef, _ := td.ledger.SetRecord(&record.CodeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	migrationRef, _ := td.ledger.SetRecord(&record.CodeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	migrationRefs := []record.Reference{*migrationRef}
	migrationCoreRefs := []core.RecordRef{*migrationRef.CoreRef()}
	updateCoreRef, err := td.manager.UpdateClass(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *classRef.CoreRef(), *codeRef.CoreRef(), migrationCoreRefs)
	assert.NoError(t, err)
	updateRef := record.Core2Reference(*updateCoreRef)
	updateRec, getErr := td.ledger.GetRecord(&updateRef)
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
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	_, err := td.manager.ActivateObj(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRandomRef().CoreRef(), []byte{})
	assert.NotNil(t, err)
	notClassRef, _ := td.ledger.SetRecord(&record.ObjectActivateRecord{})
	_, err = td.manager.ActivateClass(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *notClassRef.CoreRef(), []byte{})
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_ActivateObj_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	memory := []byte{1, 2, 3}
	classRef, _ := td.ledger.SetRecord(&record.ClassActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.ledger.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *classRef,
	})
	activateCoreRef, err := td.manager.ActivateObj(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *classRef.CoreRef(), memory)
	assert.Nil(t, err)
	activateRef := record.Core2Reference(*activateCoreRef)
	activateRec, err := td.ledger.GetRecord(&activateRef)
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
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	_, err := td.manager.DeactivateClass(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRandomRef().CoreRef())
	assert.NotNil(t, err)
	notObjRef, _ := td.ledger.SetRecord(&record.ClassActivateRecord{})
	_, err = td.manager.DeactivateClass(*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *notObjRef.CoreRef())
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateObj_VerifiesObjectIsActive(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	objRef, _ := td.ledger.SetRecord(&record.ObjectActivateRecord{})
	deactivateRef, _ := td.ledger.SetRecord(&record.DeactivationRecord{})
	td.ledger.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *deactivateRef,
	})
	_, err := td.manager.DeactivateObj(*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *objRef.CoreRef())
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateObj_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	objRef, _ := td.ledger.SetRecord(&record.ObjectActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.ledger.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *objRef,
	})
	deactivateCoreRef, err := td.manager.DeactivateObj(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *objRef.CoreRef(),
	)
	assert.Nil(t, err)
	deactivateRef := record.Core2Reference(*deactivateCoreRef)
	deactivateRec, err := td.ledger.GetRecord(&deactivateRef)
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
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	_, err := td.manager.UpdateObj(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRandomRef().CoreRef(), nil)
	assert.NotNil(t, err)
	notObjRef, _ := td.ledger.SetRecord(&record.CodeRecord{})
	_, err = td.manager.UpdateObj(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *notObjRef.CoreRef(), nil,
	)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateObj_VerifiesObjectIsActive(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	objRef, _ := td.ledger.SetRecord(&record.ObjectActivateRecord{})
	deactivateRef, _ := td.ledger.SetRecord(&record.DeactivationRecord{})
	td.ledger.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *deactivateRef,
	})
	_, err := td.manager.UpdateObj(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *objRef.CoreRef(), nil,
	)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateObj_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	objRef, _ := td.ledger.SetRecord(&record.ObjectActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.ledger.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *objRef,
	})
	memory := []byte{1, 2, 3}
	updateCoreRef, err := td.manager.UpdateObj(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *objRef.CoreRef(), memory)
	assert.Nil(t, err)
	updateRef := record.Core2Reference(*updateCoreRef)
	updateRec, err := td.ledger.GetRecord(&updateRef)
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
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	_, err := td.manager.AppendObjDelegate(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRandomRef().CoreRef(), nil,
	)
	assert.NotNil(t, err)
	notObjRef, _ := td.ledger.SetRecord(&record.CodeRecord{})
	_, err = td.manager.AppendObjDelegate(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *notObjRef.CoreRef(), nil,
	)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_AppendObjDelegate_VerifiesClassIsActive(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	objRef, _ := td.ledger.SetRecord(&record.ObjectActivateRecord{})
	deactivateRef, _ := td.ledger.SetRecord(&record.DeactivationRecord{})
	td.ledger.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *deactivateRef,
	})
	_, err := td.manager.AppendObjDelegate(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *objRef.CoreRef(), nil,
	)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_AppendObjDelegate_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	objRef, _ := td.ledger.SetRecord(&record.ObjectActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.ledger.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *objRef,
	})
	memory := []byte{1, 2, 3}
	appendCoreRef, err := td.manager.AppendObjDelegate(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *objRef.CoreRef(), memory,
	)
	assert.Nil(t, err)
	appendRef := record.Core2Reference(*appendCoreRef)
	appendRec, _ := td.ledger.GetRecord(&appendRef)
	objIndex, err := td.ledger.GetObjectIndex(objRef)
	assert.NoError(t, err)
	assert.Equal(t, objIndex.AppendRefs, []record.Reference{appendRef})
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

func TestLedgerArtifactManager_GetCode(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	codeRef := genRandomRef()

	desc, err := td.manager.GetCode(*codeRef.CoreRef())
	assert.NoError(t, err)
	assert.Equal(t, CodeDescriptor{
		manager: td.manager.(*LedgerArtifactManager),
		ref:     codeRef,
	}, *desc.(*CodeDescriptor))
}

func TestLedgerArtifactManager_GetLatestObj_VerifiesRecords(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	_, err := td.manager.GetLatestObj(*genRandomRef().CoreRef())
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_GetLatestObj_VerifiesClassIsActive(t *testing.T) {
	t.Parallel()
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
	_, err := td.manager.GetLatestObj(*objectRef.CoreRef())
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_GetLatestObj_ReturnsCorrectDescriptors(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	ledgerManager, _ := td.manager.(*LedgerArtifactManager)

	classRef, _ := td.ledger.SetRecord(&record.ClassActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	classRec, _ := td.ledger.GetRecord(classRef)
	classRecCasted, _ := classRec.(*record.ClassActivateRecord)
	classAmendRef, _ := td.ledger.SetRecord(&record.ClassAmendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  *td.domainRef,
					RequestRecord: *td.requestRef,
				},
			},
		},
	})
	classAmendRec, _ := td.ledger.GetRecord(classAmendRef)
	classAmendRecCasted, _ := classAmendRec.(*record.ClassAmendRecord)
	classIndex := index.ClassLifeline{
		LatestStateRef: *classAmendRef,
	}
	td.ledger.SetClassIndex(classRef, &classIndex)

	objectRef, _ := td.ledger.SetRecord(&record.ObjectActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
		Memory: []byte{3},
	})
	objectRec, _ := td.ledger.GetRecord(objectRef)
	objectRecCasted, _ := objectRec.(*record.ObjectActivateRecord)
	objectAmendRef, _ := td.ledger.SetRecord(&record.ObjectAmendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *td.domainRef,
				},
			},
		},
		NewMemory: []byte{4},
	})
	objectAmendRec, _ := td.ledger.GetRecord(objectAmendRef)
	objectAmendRecCasted, _ := objectAmendRec.(*record.ObjectAmendRecord)
	objectIndex := index.ObjectLifeline{
		LatestStateRef: *objectAmendRef,
		ClassRef:       *classRef,
	}
	td.ledger.SetObjectIndex(objectRef, &objectIndex)

	objDesc, err := td.manager.GetLatestObj(*objectRef.CoreRef())
	assert.NoError(t, err)
	expectedObjDesc := &ObjectDescriptor{
		manager: ledgerManager,

		headRef: objectRef,
		stateRef: &record.Reference{
			Domain: objectRef.Domain,
			Record: objectAmendRef.Record,
		},

		headRecord:    objectRecCasted,
		stateRecord:   objectAmendRecCasted,
		lifelineIndex: &objectIndex,

		classDescriptor: &ClassDescriptor{
			manager: ledgerManager,

			headRef: classRef,
			stateRef: &record.Reference{
				Domain: classRef.Domain,
				Record: classAmendRef.Record,
			},

			headRecord:    classRecCasted,
			stateRecord:   classAmendRecCasted,
			lifelineIndex: &classIndex,
		},
	}

	assert.Equal(t, *expectedObjDesc, *objDesc.(*ObjectDescriptor))
}
