/*
 *    Copyright 2018 Insolar
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
	coreRef[core.RecordIDSize-1] = byte(rand.Int() % 256)
	ref := record.Core2Reference(coreRef)
	return &ref
}

type preparedAMTestData struct {
	db         *storage.DB
	manager    core.ArtifactManager
	domainRef  *record.Reference
	requestRef *record.Reference
}

type messageBusMock struct {
	handler *MessageHandler
}

func (m *messageBusMock) Send(e core.Message) (core.Reply, error) {
	return m.handler.Handle(e)
}

func (m *messageBusMock) SendAsync(e core.Message) {
	m.handler.Handle(e)
}

func prepareAMTestData(t *testing.T) (preparedAMTestData, func()) {
	db, cleaner := storagetest.TmpDB(t, "")

	return preparedAMTestData{
		db: db,
		manager: &LedgerArtifactManager{
			db:         db,
			messageBus: &messageBusMock{handler: &MessageHandler{db: db}},
		},
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
	typeRec, err := td.db.GetRecord(&ref)
	assert.NoError(t, err)
	assert.Equal(t, &record.TypeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  *td.domainRef,
					RequestRecord: *td.requestRef,
				},
			},
		},
		TypeDeclaration: typeDec,
	}, typeRec)
}

func TestLedgerArtifactManager_DeployCode_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	codeMap := map[core.MachineType][]byte{1: {1}}
	coreRef, err := td.manager.DeployCode(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), codeMap,
	)
	assert.NoError(t, err)
	ref := record.Core2Reference(*coreRef)
	codeRec, err := td.db.GetRecord(&ref)
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

func TestLedgerArtifactManager_ActivateClass_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	activateCoreRef, err := td.manager.ActivateClass(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(),
	)
	activateRef := record.Core2Reference(*activateCoreRef)
	assert.Nil(t, err)
	activateRec, getErr := td.db.GetRecord(&activateRef)
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

	notClassRef, _ := td.db.SetRecord(&record.CodeRecord{})
	_, err = td.manager.DeactivateClass(*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *notClassRef.CoreRef())
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateClass_VerifiesClassIsActive(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	classRef, _ := td.db.SetRecord(&record.ClassActivateRecord{})
	deactivateRef, _ := td.db.SetRecord(&record.DeactivationRecord{})
	td.db.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *deactivateRef,
	})
	_, err := td.manager.DeactivateClass(*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *classRef.CoreRef())
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateClass_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	classRef, _ := td.db.SetRecord(&record.ClassActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.db.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *classRef,
	})

	deactivateCoreRef, err := td.manager.DeactivateClass(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *classRef.CoreRef(),
	)
	assert.NoError(t, err)
	deactivateRef := record.Core2Reference(*deactivateCoreRef)
	deactivateRec, err := td.db.GetRecord(&deactivateRef)
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
	notClassRef, _ := td.db.SetRecord(&record.CodeRecord{})
	_, err = td.manager.UpdateClass(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *notClassRef.CoreRef(), *genRandomRef().CoreRef(), nil,
	)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateClass_VerifiesClassIsActive(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	classRef, _ := td.db.SetRecord(&record.ClassActivateRecord{})
	deactivateRef, _ := td.db.SetRecord(&record.DeactivationRecord{})
	codeRef, _ := td.db.SetRecord(&record.CodeRecord{})
	td.db.SetClassIndex(classRef, &index.ClassLifeline{
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

	classRef, _ := td.db.SetRecord(&record.ClassActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.db.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *classRef,
	})
	codeRef, _ := td.db.SetRecord(&record.CodeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
		TargetedCode: map[core.MachineType][]byte{core.MachineTypeBuiltin: {}},
	})
	migrationRef, _ := td.db.SetRecord(&record.CodeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
		TargetedCode: map[core.MachineType][]byte{core.MachineTypeBuiltin: {}},
	})
	migrationRefs := []record.Reference{*migrationRef}
	migrationCoreRefs := []core.RecordRef{*migrationRef.CoreRef()}
	updateCoreRef, err := td.manager.UpdateClass(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *classRef.CoreRef(), *codeRef.CoreRef(), migrationCoreRefs)
	assert.NoError(t, err)
	updateRef := record.Core2Reference(*updateCoreRef)
	updateRec, getErr := td.db.GetRecord(&updateRef)
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

func TestLedgerArtifactManager_ActivateObject_VerifiesRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	_, err := td.manager.ActivateObject(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRandomRef().CoreRef(), *genRandomRef().CoreRef(),
		[]byte{},
	)
	assert.NotNil(t, err)
	notClassRef, _ := td.db.SetRecord(&record.ObjectActivateRecord{})
	_, err = td.manager.ActivateObject(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *notClassRef.CoreRef(), *genRandomRef().CoreRef(), []byte{},
	)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_ActivateObject_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	memory := []byte{1, 2, 3}
	classRef, _ := td.db.SetRecord(&record.ClassActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.db.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *classRef,
	})
	parentRef, _ := td.db.SetRecord(&record.ObjectActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.db.SetObjectIndex(parentRef, &index.ObjectLifeline{
		ClassRef:       *classRef,
		LatestStateRef: *parentRef,
	})

	activateCoreRef, err := td.manager.ActivateObject(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *classRef.CoreRef(), *parentRef.CoreRef(), memory,
	)
	assert.Nil(t, err)
	activateRef := record.Core2Reference(*activateCoreRef)
	activateRec, err := td.db.GetRecord(&activateRef)
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
		Parent:              *parentRef,
		Delegate:            false,
	})
}

func TestLedgerArtifactManager_ActivateObjectDelegate_VerifiesRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	_, err := td.manager.ActivateObjectDelegate(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRandomRef().CoreRef(), *genRandomRef().CoreRef(),
		[]byte{},
	)
	assert.NotNil(t, err)
	notClassRef, _ := td.db.SetRecord(&record.ObjectActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	_, err = td.manager.ActivateObjectDelegate(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *notClassRef.CoreRef(), *notClassRef.CoreRef(), []byte{},
	)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_ActivateObjectDelegate_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	memory := []byte{1, 2, 3}
	classRef, _ := td.db.SetRecord(&record.ClassActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.db.SetClassIndex(classRef, &index.ClassLifeline{
		LatestStateRef: *classRef,
	})
	parentRef, _ := td.db.SetRecord(&record.ObjectActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.db.SetObjectIndex(parentRef, &index.ObjectLifeline{
		ClassRef:       *classRef,
		LatestStateRef: *parentRef,
		Delegates:      map[core.RecordRef]record.Reference{},
	})

	activateCoreRef, err := td.manager.ActivateObjectDelegate(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *classRef.CoreRef(), *parentRef.CoreRef(), memory,
	)
	assert.Nil(t, err)
	activateRef := record.Core2Reference(*activateCoreRef)
	activateRec, err := td.db.GetRecord(&activateRef)
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
		Parent:              *parentRef,
		Delegate:            true,
	})

	delegate, err := td.manager.GetDelegate(*parentRef.CoreRef(), *classRef.CoreRef())
	assert.NoError(t, err)
	assert.Equal(t, activateCoreRef, delegate)
}

func TestLedgerArtifactManager_DeActivateObject_VerifiesRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	_, err := td.manager.DeactivateClass(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRandomRef().CoreRef())
	assert.NotNil(t, err)
	notObjRef, _ := td.db.SetRecord(&record.ClassActivateRecord{})
	_, err = td.manager.DeactivateClass(*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *notObjRef.CoreRef())
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeActivateObject_VerifiesObjectIsActive(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	objRef, _ := td.db.SetRecord(&record.ObjectActivateRecord{})
	deactivateRef, _ := td.db.SetRecord(&record.DeactivationRecord{})
	td.db.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *deactivateRef,
	})
	_, err := td.manager.DeactivateObject(*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *objRef.CoreRef())
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateObject_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	objRef, _ := td.db.SetRecord(&record.ObjectActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.db.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *objRef,
	})
	deactivateCoreRef, err := td.manager.DeactivateObject(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *objRef.CoreRef(),
	)
	assert.Nil(t, err)
	deactivateRef := record.Core2Reference(*deactivateCoreRef)
	deactivateRec, err := td.db.GetRecord(&deactivateRef)
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

func TestLedgerArtifactManager_UpdateObject_VerifiesRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	_, err := td.manager.UpdateObject(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRandomRef().CoreRef(), nil)
	assert.NotNil(t, err)
	notObjRef, _ := td.db.SetRecord(&record.CodeRecord{})
	_, err = td.manager.UpdateObject(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *notObjRef.CoreRef(), nil,
	)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateObject_VerifiesObjectIsActive(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	objRef, _ := td.db.SetRecord(&record.ObjectActivateRecord{})
	deactivateRef, _ := td.db.SetRecord(&record.DeactivationRecord{})
	td.db.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *deactivateRef,
	})
	_, err := td.manager.UpdateObject(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *objRef.CoreRef(), nil,
	)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateObject_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	objRef, _ := td.db.SetRecord(&record.ObjectActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.db.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestStateRef: *objRef,
	})
	memory := []byte{1, 2, 3}
	updateCoreRef, err := td.manager.UpdateObject(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *objRef.CoreRef(), memory)
	assert.Nil(t, err)
	updateRef := record.Core2Reference(*updateCoreRef)
	updateRec, err := td.db.GetRecord(&updateRef)
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

func TestLedgerArtifactManager_GetClass_ReturnsCorrectDescriptors(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	codeRef := *genRandomRef()
	classRef, _ := td.db.SetRecord(&record.ClassActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *td.domainRef,
				},
			},
		},
	})
	classAmendRef, _ := td.db.SetRecord(&record.ClassAmendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  *td.domainRef,
					RequestRecord: *td.requestRef,
				},
			},
		},
		NewCode: codeRef,
	})
	classIndex := index.ClassLifeline{
		LatestStateRef: *classAmendRef,
	}
	td.db.SetClassIndex(classRef, &classIndex)

	classDesc, err := td.manager.GetClass(*classRef.CoreRef(), nil)
	assert.NoError(t, err)
	expectedClassDesc := &ClassDescriptor{
		am:    td.manager,
		head:  *classRef.CoreRef(),
		state: *classAmendRef.CoreRef(),
		code:  codeRef.CoreRef(),
	}

	assert.Equal(t, *expectedClassDesc, *classDesc.(*ClassDescriptor))
}

func TestLedgerArtifactManager_GetObject_VerifiesRecords(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	_, err := td.manager.GetObject(*genRandomRef().CoreRef(), nil)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_GetLatestObj_ReturnsCorrectDescriptors(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	classRef, _ := td.db.SetRecord(&record.ClassActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *td.domainRef,
				},
			},
		},
	})
	classAmendRef, _ := td.db.SetRecord(&record.ClassAmendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  *td.domainRef,
					RequestRecord: *td.requestRef,
				},
			},
		},
	})
	classIndex := index.ClassLifeline{
		LatestStateRef: *classAmendRef,
	}
	td.db.SetClassIndex(classRef, &classIndex)

	objectRef, _ := td.db.SetRecord(&record.ObjectActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *td.domainRef,
				},
			},
		},
		Memory: []byte{3},
	})
	objectAmendRef, _ := td.db.SetRecord(&record.ObjectAmendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *td.domainRef,
				},
			},
		},
		NewMemory: []byte{4},
	})
	objectIndex := index.ObjectLifeline{
		LatestStateRef: *objectAmendRef,
		ClassRef:       *classRef,
		Children:       []record.Reference{*genRandomRef(), *genRandomRef()},
	}
	td.db.SetObjectIndex(objectRef, &objectIndex)

	objDesc, err := td.manager.GetObject(*objectRef.CoreRef(), nil)
	assert.NoError(t, err)
	expectedChildren := make([]core.RecordRef, 0, len(objectIndex.Children))
	for _, c := range objectIndex.Children {
		expectedChildren = append(expectedChildren, *c.CoreRef())
	}
	expectedObjDesc := &ObjectDescriptor{
		am: td.manager,

		head:     *objectRef.CoreRef(),
		state:    *objectAmendRef.CoreRef(),
		class:    *classRef.CoreRef(),
		memory:   []byte{4},
		children: expectedChildren,
	}

	assert.Equal(t, *expectedObjDesc, *objDesc.(*ObjectDescriptor))
}
