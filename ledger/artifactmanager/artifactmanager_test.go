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

	"github.com/pkg/errors"
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

func genRefWithID(id *record.ID) *core.RecordRef {
	ref := record.Reference{Record: *id}
	return ref.CoreRef()
}

type preparedAMTestData struct {
	db         *storage.DB
	manager    core.ArtifactManager
	domainRef  *record.Reference
	requestRef *record.Reference
}

type messageBusMock struct {
	handlers map[core.MessageType]core.MessageHandler
}

func NewMessageBusMock() *messageBusMock {
	return &messageBusMock{handlers: map[core.MessageType]core.MessageHandler{}}
}

func (mb *messageBusMock) Register(p core.MessageType, handler core.MessageHandler) error {
	_, ok := mb.handlers[p]
	if ok {
		return errors.New("handler for this type already exists")
	}

	mb.handlers[p] = handler
	return nil
}

func (mb *messageBusMock) Start(components core.Components) error {
	panic("implement me")
}

func (mb *messageBusMock) Stop() error {
	panic("implement me")
}

func (mb *messageBusMock) Send(m core.Message) (core.Reply, error) {
	handler, ok := mb.handlers[m.Type()]
	if !ok {
		return nil, errors.New("no handler for this message type")
	}

	return handler(m)
}

func (mb *messageBusMock) SendAsync(m core.Message) {
	panic("implement me")
}

func prepareAMTestData(t *testing.T) (preparedAMTestData, func()) {
	db, cleaner := storagetest.TmpDB(t, "")

	mb := NewMessageBusMock()
	components := core.Components{MessageBus: mb}
	handler := MessageHandler{db: db}
	handler.Link(components)

	return preparedAMTestData{
		db: db,
		manager: &LedgerArtifactManager{
			db:         db,
			messageBus: mb,
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
	typeRec, err := td.db.GetRecord(&ref.Record)
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
	codeRec, err := td.db.GetRecord(&ref.Record)
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
	activateRec, getErr := td.db.GetRecord(&activateRef.Record)
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

	notClassID, _ := td.db.SetRecord(&record.CodeRecord{})
	_, err = td.manager.DeactivateClass(*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRefWithID(notClassID))
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateClass_VerifiesClassIsActive(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	classID, _ := td.db.SetRecord(&record.ClassActivateRecord{})
	deactivateRef, _ := td.db.SetRecord(&record.DeactivationRecord{})
	td.db.SetClassIndex(classID, &index.ClassLifeline{
		LatestState: *deactivateRef,
	})
	_, err := td.manager.DeactivateClass(*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRefWithID(classID))
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateClass_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	classID, _ := td.db.SetRecord(&record.ClassActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.db.SetClassIndex(classID, &index.ClassLifeline{
		LatestState: *classID,
	})

	deactivateCoreID, err := td.manager.DeactivateClass(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRefWithID(classID),
	)
	assert.NoError(t, err)
	deactivateID := record.Bytes2ID(deactivateCoreID[:])
	deactivateRec, err := td.db.GetRecord(&deactivateID)
	assert.NoError(t, err)
	assert.Equal(t, deactivateRec, &record.DeactivationRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  *td.domainRef,
					RequestRecord: *td.requestRef,
				},
			},
			AmendedRecord: *classID,
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
	notClassID, _ := td.db.SetRecord(&record.CodeRecord{})
	_, err = td.manager.UpdateClass(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRefWithID(notClassID), *genRandomRef().CoreRef(), nil,
	)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateClass_VerifiesClassIsActive(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	classID, _ := td.db.SetRecord(&record.ClassActivateRecord{})
	deactivateID, _ := td.db.SetRecord(&record.DeactivationRecord{})
	codeRef, _ := td.db.SetRecord(&record.CodeRecord{})
	td.db.SetClassIndex(classID, &index.ClassLifeline{
		LatestState: *deactivateID,
	})
	_, err := td.manager.UpdateClass(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRefWithID(classID), *genRefWithID(codeRef), nil)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateClass_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	classID, _ := td.db.SetRecord(&record.ClassActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.db.SetClassIndex(classID, &index.ClassLifeline{
		LatestState: *classID,
	})
	codeID, _ := td.db.SetRecord(&record.CodeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
		TargetedCode: map[core.MachineType][]byte{core.MachineTypeBuiltin: {}},
	})
	migrationID, _ := td.db.SetRecord(&record.CodeRecord{
		StorageRecord: record.StorageRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
		TargetedCode: map[core.MachineType][]byte{core.MachineTypeBuiltin: {}},
	})
	migrationRefs := []record.Reference{{Domain: td.domainRef.Domain, Record: *migrationID}}
	migrationCoreRefs := []core.RecordRef{*migrationRefs[0].CoreRef()}
	updateCoreID, err := td.manager.UpdateClass(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRefWithID(classID), *genRefWithID(codeID),
		migrationCoreRefs,
	)
	assert.NoError(t, err)
	updateID := record.Bytes2ID(updateCoreID[:])
	updateRec, getErr := td.db.GetRecord(&updateID)
	assert.Nil(t, getErr)
	assert.Equal(t, updateRec, &record.ClassAmendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  *td.domainRef,
					RequestRecord: *td.requestRef,
				},
			},
			AmendedRecord: *classID,
		},
		NewCode:    record.Reference{Domain: td.requestRef.Domain, Record: *codeID},
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
	notClassID, _ := td.db.SetRecord(&record.ObjectActivateRecord{})
	_, err = td.manager.ActivateObject(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRefWithID(notClassID), *genRandomRef().CoreRef(), []byte{},
	)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_ActivateObject_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	memory := []byte{1, 2, 3}
	classID, _ := td.db.SetRecord(&record.ClassActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.db.SetClassIndex(classID, &index.ClassLifeline{
		LatestState: *classID,
	})
	parentID, _ := td.db.SetRecord(&record.ObjectActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.db.SetObjectIndex(parentID, &index.ObjectLifeline{
		ClassRef:    record.Reference{Record: *classID},
		LatestState: *parentID,
	})

	activateCoreRef, err := td.manager.ActivateObject(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRefWithID(classID), *genRefWithID(parentID), memory,
	)
	assert.Nil(t, err)
	activateRef := record.Core2Reference(*activateCoreRef)
	activateRec, err := td.db.GetRecord(&activateRef.Record)
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
		ClassActivateRecord: record.Reference{Domain: td.requestRef.Domain, Record: *classID},
		Memory:              memory,
		Parent:              record.Reference{Domain: td.requestRef.Domain, Record: *parentID},
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
	notClassID, _ := td.db.SetRecord(&record.ObjectActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	_, err = td.manager.ActivateObjectDelegate(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRefWithID(notClassID), *genRefWithID(notClassID),
		[]byte{},
	)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_ActivateObjectDelegate_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	memory := []byte{1, 2, 3}
	classID, _ := td.db.SetRecord(&record.ClassActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.db.SetClassIndex(classID, &index.ClassLifeline{
		LatestState: *classID,
	})
	parentID, _ := td.db.SetRecord(&record.ObjectActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.db.SetObjectIndex(parentID, &index.ObjectLifeline{
		ClassRef:    record.Reference{Domain: td.requestRef.Domain, Record: *classID},
		LatestState: *parentID,
	})

	activateCoreRef, err := td.manager.ActivateObjectDelegate(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRefWithID(classID), *genRefWithID(parentID), memory,
	)
	assert.Nil(t, err)
	activateRef := record.Core2Reference(*activateCoreRef)
	activateRec, err := td.db.GetRecord(&activateRef.Record)
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
		ClassActivateRecord: record.Reference{Domain: td.domainRef.Domain, Record: *classID},
		Memory:              memory,
		Parent:              record.Reference{Domain: td.domainRef.Domain, Record: *parentID},
		Delegate:            true,
	})

	delegate, err := td.manager.GetDelegate(*genRefWithID(parentID), *genRefWithID(classID))
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
	notObjID, _ := td.db.SetRecord(&record.ClassActivateRecord{})
	_, err = td.manager.DeactivateClass(*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRefWithID(notObjID))
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeActivateObject_VerifiesObjectIsActive(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	objRef, _ := td.db.SetRecord(&record.ObjectActivateRecord{})
	deactivateID, _ := td.db.SetRecord(&record.DeactivationRecord{})
	td.db.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestState: *deactivateID,
	})
	_, err := td.manager.DeactivateObject(*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRefWithID(deactivateID))
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_DeactivateObject_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	objID, _ := td.db.SetRecord(&record.ObjectActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.db.SetObjectIndex(objID, &index.ObjectLifeline{
		LatestState: *objID,
	})
	deactivateCoreID, err := td.manager.DeactivateObject(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRefWithID(objID),
	)
	assert.Nil(t, err)
	deactivateID := record.Bytes2ID(deactivateCoreID[:])
	deactivateRec, err := td.db.GetRecord(&deactivateID)
	assert.Nil(t, err)
	assert.Equal(t, deactivateRec, &record.DeactivationRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  *td.domainRef,
					RequestRecord: *td.requestRef,
				},
			},
			AmendedRecord: *objID,
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
	notObjID, _ := td.db.SetRecord(&record.CodeRecord{})
	_, err = td.manager.UpdateObject(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRefWithID(notObjID), nil,
	)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateObject_VerifiesObjectIsActive(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	objRef, _ := td.db.SetRecord(&record.ObjectActivateRecord{})
	deactivateID, _ := td.db.SetRecord(&record.DeactivationRecord{})
	td.db.SetObjectIndex(objRef, &index.ObjectLifeline{
		LatestState: *deactivateID,
	})
	_, err := td.manager.UpdateObject(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRefWithID(deactivateID), nil,
	)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateObject_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	objID, _ := td.db.SetRecord(&record.ObjectActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *genRandomRef(),
				},
			},
		},
	})
	td.db.SetObjectIndex(objID, &index.ObjectLifeline{
		LatestState: *objID,
	})
	memory := []byte{1, 2, 3}
	updateCoreID, err := td.manager.UpdateObject(
		*td.domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRefWithID(objID), memory)
	assert.Nil(t, err)
	updateID := record.Bytes2ID(updateCoreID[:])
	updateRec, err := td.db.GetRecord(&updateID)
	assert.Nil(t, err)
	assert.Equal(t, updateRec, &record.ObjectAmendRecord{
		AmendRecord: record.AmendRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord:  *td.domainRef,
					RequestRecord: *td.requestRef,
				},
			},
			AmendedRecord: *objID,
		},
		NewMemory: memory,
	})
}

func TestLedgerArtifactManager_GetClass_ReturnsCorrectDescriptors(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	codeRef := *genRandomRef()
	classID, _ := td.db.SetRecord(&record.ClassActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *td.domainRef,
				},
			},
		},
	})
	classAmendID, _ := td.db.SetRecord(&record.ClassAmendRecord{
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
		LatestState: *classAmendID,
	}
	td.db.SetClassIndex(classID, &classIndex)

	classRef := genRefWithID(classID)
	classDesc, err := td.manager.GetClass(*classRef, nil)
	assert.NoError(t, err)
	expectedClassDesc := &ClassDescriptor{
		am:    td.manager,
		head:  *classRef,
		state: *getReference(td.requestRef.CoreRef(), classAmendID),
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

	classID, _ := td.db.SetRecord(&record.ClassActivateRecord{
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
		LatestState: *classAmendRef,
	}
	td.db.SetClassIndex(classID, &classIndex)

	objectID, _ := td.db.SetRecord(&record.ObjectActivateRecord{
		ActivationRecord: record.ActivationRecord{
			StatefulResult: record.StatefulResult{
				ResultRecord: record.ResultRecord{
					DomainRecord: *td.domainRef,
				},
			},
		},
		Memory: []byte{3},
	})
	objectAmendID, _ := td.db.SetRecord(&record.ObjectAmendRecord{
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
		LatestState: *objectAmendID,
		ClassRef:    record.Reference{Domain: td.requestRef.Domain, Record: *classID},
		Children:    []record.Reference{*genRandomRef(), *genRandomRef()},
	}
	td.db.SetObjectIndex(objectID, &objectIndex)

	objDesc, err := td.manager.GetObject(*genRefWithID(objectID), nil)
	assert.NoError(t, err)
	expectedChildren := make([]core.RecordRef, 0, len(objectIndex.Children))
	for _, c := range objectIndex.Children {
		expectedChildren = append(expectedChildren, *c.CoreRef())
	}
	expectedObjDesc := &ObjectDescriptor{
		am: td.manager,

		head:     *getReference(td.requestRef.CoreRef(), objectID),
		state:    *getReference(td.requestRef.CoreRef(), objectAmendID),
		class:    *getReference(td.requestRef.CoreRef(), classID),
		memory:   []byte{4},
		children: expectedChildren,
	}

	assert.Equal(t, *expectedObjDesc, *objDesc.(*ObjectDescriptor))
}
