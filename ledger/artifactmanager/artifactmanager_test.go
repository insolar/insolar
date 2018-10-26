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

	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/cryptohelpers/hash"
	"github.com/insolar/insolar/inscontext"
	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/testutils/testmessagebus"
)

var (
	domainID  = *genRandomID(0)
	domainRef = record.Reference{
		Record: domainID,
		Domain: domainID,
	}
)

func genRandomID(pulse core.PulseNumber) *record.ID {
	id := record.ID{
		Pulse: pulse,
		Hash:  []byte{byte(rand.Int() % 256)},
	}
	zeroFilledID := record.Bytes2ID(id.CoreID()[:]) // Double conversion hack to fill missing length with zeros.
	return &zeroFilledID
}

func genRandomRef(pulse core.PulseNumber) *record.Reference {
	return &record.Reference{
		Record: *genRandomID(pulse),
		Domain: domainID,
	}
}

func genRefWithID(id *record.ID) *core.RecordRef {
	ref := record.Reference{Record: *id, Domain: domainID}
	return ref.CoreRef()
}

type preparedAMTestData struct {
	db         *storage.DB
	manager    *LedgerArtifactManager
	requestRef *record.Reference
}

func prepareAMTestData(t *testing.T) (preparedAMTestData, func()) {
	db, cleaner := storagetest.TmpDB(t, "")

	mb := testmessagebus.NewTestMessageBus()
	components := core.Components{MessageBus: mb}
	handler := MessageHandler{db: db, jetDropHandlers: map[core.MessageType]internalHandler{}}
	handler.Link(components)

	return preparedAMTestData{
		db: db,
		manager: &LedgerArtifactManager{
			db:                   db,
			messageBus:           mb,
			getChildrenChunkSize: 100,
		},
		requestRef: genRandomRef(0),
	}, cleaner
}

func TestLedgerArtifactManager_RegisterRequest(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	ctx := inscontext.TODO()

	msg := message.BootstrapRequest{Name: "my little message"}
	coreID, err := td.manager.RegisterRequest(ctx, &msg)
	assert.NoError(t, err)
	id := record.Bytes2ID(coreID[:])
	rec, err := td.db.GetRecord(&id)
	assert.NoError(t, err)
	assert.Equal(t, message.MustSerializeBytes(&msg), rec.(*record.CallRequest).Payload)
}

func TestLedgerArtifactManager_DeclareType(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	ctx := inscontext.TODO()

	typeDec := []byte{1, 2, 3}
	coreID, err := td.manager.DeclareType(ctx, *domainRef.CoreRef(), *td.requestRef.CoreRef(), typeDec)
	assert.NoError(t, err)
	id := record.Bytes2ID(coreID[:])
	typeRec, err := td.db.GetRecord(&id)
	assert.NoError(t, err)
	assert.Equal(t, &record.TypeRecord{
		ResultRecord: record.ResultRecord{
			Domain:  domainRef,
			Request: *td.requestRef,
		},
		TypeDeclaration: typeDec,
	}, typeRec)
}

func TestLedgerArtifactManager_DeployCode_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	ctx := inscontext.TODO()

	coreID, err := td.manager.DeployCode(
		ctx,
		*domainRef.CoreRef(), *td.requestRef.CoreRef(), []byte{1, 2, 3}, core.MachineTypeBuiltin,
	)
	assert.NoError(t, err)
	id := record.Bytes2ID(coreID[:])
	codeRec, err := td.db.GetRecord(&id)
	assert.NoError(t, err)
	assert.Equal(t, codeRec, &record.CodeRecord{
		ResultRecord: record.ResultRecord{
			Domain:  domainRef,
			Request: *td.requestRef,
		},
		Code:        []byte{1, 2, 3},
		MachineType: core.MachineTypeBuiltin,
	})
}

func TestLedgerArtifactManager_ActivateClass_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	ctx := inscontext.TODO()

	codeID, err := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.CodeRecord{})
	codeRef := record.Reference{Record: *codeID, Domain: domainID}
	classRef := genRandomRef(0)
	activateCoreID, err := td.manager.ActivateClass(
		ctx,
		*domainRef.CoreRef(),
		*classRef.CoreRef(),
		*codeRef.CoreRef(),
		core.MachineTypeBuiltin,
	)
	assert.Nil(t, err)
	activateID := record.Bytes2ID(activateCoreID[:])
	activateRec, getErr := td.db.GetRecord(&activateID)
	assert.Nil(t, getErr)
	assert.Equal(t, activateRec, &record.ClassActivateRecord{
		ResultRecord: record.ResultRecord{
			Domain:  domainRef,
			Request: *classRef,
		},
		ClassStateRecord: record.ClassStateRecord{
			MachineType: core.MachineTypeBuiltin,
			Code:        codeRef,
		},
	})
	idx, err := td.db.GetClassIndex(&classRef.Record, false)
	assert.NoError(t, err)
	assert.Equal(t, activateID, *idx.LatestState)
}

func TestLedgerArtifactManager_DeactivateClass_VerifiesClassIsActive(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	ctx := inscontext.TODO()

	classID, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.ClassActivateRecord{})
	err := td.db.SetClassIndex(classID, &index.ClassLifeline{
		State: record.StateDeactivation,
	})
	assert.NoError(t, err)
	_, err = td.manager.DeactivateClass(
		ctx,
		*domainRef.CoreRef(),
		*td.requestRef.CoreRef(),
		*genRefWithID(classID),
		*genRandomID(0).CoreID(),
	)
	assert.Error(t, err)
}

func TestLedgerArtifactManager_DeactivateClass_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	ctx := inscontext.TODO()

	classID, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.ClassActivateRecord{
		ResultRecord: record.ResultRecord{
			Domain: *genRandomRef(0),
		},
	})
	td.db.SetClassIndex(classID, &index.ClassLifeline{
		State:       record.StateActivation,
		LatestState: classID,
	})

	deactivateCoreID, err := td.manager.DeactivateClass(
		ctx,
		*domainRef.CoreRef(), *td.requestRef.CoreRef(), *genRefWithID(classID), *classID.CoreID(),
	)
	assert.NoError(t, err)
	deactivateID := record.Bytes2ID(deactivateCoreID[:])
	deactivateRec, err := td.db.GetRecord(&deactivateID)
	assert.NoError(t, err)
	assert.Equal(t, deactivateRec, &record.DeactivationRecord{
		ResultRecord: record.ResultRecord{
			Domain:  domainRef,
			Request: *td.requestRef,
		},
		PrevState: *classID,
	})
}

func TestLedgerArtifactManager_UpdateClass_VerifiesClassIsActive(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	ctx := inscontext.TODO()

	classID, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.ClassActivateRecord{})
	deactivateID, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.DeactivationRecord{})
	codeRef, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.CodeRecord{})
	err := td.db.SetClassIndex(classID, &index.ClassLifeline{
		State:       record.StateDeactivation,
		LatestState: deactivateID,
	})
	assert.NoError(t, err)
	_, err = td.manager.UpdateClass(
		ctx,
		*domainRef.CoreRef(),
		*td.requestRef.CoreRef(),
		*genRefWithID(classID),
		*genRefWithID(codeRef),
		core.MachineTypeBuiltin,
		*deactivateID.CoreID(),
	)
	assert.NotNil(t, err)
}

func TestLedgerArtifactManager_UpdateClass_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	ctx := inscontext.TODO()

	classID, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.ClassActivateRecord{
		ResultRecord: record.ResultRecord{
			Domain: domainRef,
		},
	})
	td.db.SetClassIndex(classID, &index.ClassLifeline{
		State:       record.StateActivation,
		LatestState: classID,
	})
	codeID, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.CodeRecord{
		ResultRecord: record.ResultRecord{
			Domain: domainRef,
		},
		Code: []byte{1},
	})
	updateCoreID, err := td.manager.UpdateClass(
		ctx,
		*domainRef.CoreRef(),
		*td.requestRef.CoreRef(),
		*genRefWithID(classID),
		*genRefWithID(codeID),
		core.MachineTypeBuiltin,
		*classID.CoreID(),
	)
	assert.NoError(t, err)
	updateID := record.Bytes2ID(updateCoreID[:])
	updateRec, getErr := td.db.GetRecord(&updateID)
	assert.Nil(t, getErr)
	assert.Equal(t, updateRec, &record.ClassAmendRecord{
		ResultRecord: record.ResultRecord{
			Domain:  domainRef,
			Request: *td.requestRef,
		},
		ClassStateRecord: record.ClassStateRecord{
			Code:        record.Reference{Domain: td.requestRef.Domain, Record: *codeID},
			MachineType: core.MachineTypeBuiltin,
		},
		PrevState: *classID,
	})
}

func TestLedgerArtifactManager_ActivateObject_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	ctx := inscontext.TODO()
	memory := []byte{1, 2, 3}
	classID, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.ClassActivateRecord{
		ResultRecord: record.ResultRecord{
			Domain: *genRandomRef(0),
		},
	})
	td.db.SetClassIndex(classID, &index.ClassLifeline{
		LatestState: classID,
	})
	parentID, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.ObjectActivateRecord{
		ResultRecord: record.ResultRecord{
			Domain: *genRandomRef(0),
		},
	})
	td.db.SetObjectIndex(parentID, &index.ObjectLifeline{
		ClassRef:    record.Reference{Record: *classID},
		LatestState: parentID,
	})

	objRef := *genRandomRef(0)
	objDesc, err := td.manager.ActivateObject(
		ctx,
		*domainRef.CoreRef(),
		*objRef.CoreRef(),
		*genRefWithID(classID),
		*genRefWithID(parentID),
		false,
		memory,
	)
	assert.Nil(t, err)
	activateID := record.Bytes2ID(objDesc.StateID()[:])
	activateRec, err := td.db.GetRecord(&activateID)
	assert.Nil(t, err)
	assert.Equal(t, activateRec, &record.ObjectActivateRecord{
		ResultRecord: record.ResultRecord{
			Domain:  domainRef,
			Request: objRef,
		},
		ObjectStateRecord: record.ObjectStateRecord{
			Memory: memory,
		},
		Class:    record.Reference{Domain: td.requestRef.Domain, Record: *classID},
		Parent:   record.Reference{Domain: td.requestRef.Domain, Record: *parentID},
		Delegate: false,
	})

	idx, err := td.db.GetObjectIndex(parentID, false)
	assert.NoError(t, err)
	childRec, err := td.db.GetRecord(idx.ChildPointer)
	assert.NoError(t, err)
	assert.Equal(t, objRef, childRec.(*record.ChildRecord).Ref)

	idx, err = td.db.GetObjectIndex(&objRef.Record, false)
	assert.NoError(t, err)
	assert.Equal(t, activateID, *idx.LatestState)
}

func TestLedgerArtifactManager_DeactivateObject_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	ctx := inscontext.TODO()
	objID, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.ObjectActivateRecord{
		ResultRecord: record.ResultRecord{
			Domain: *genRandomRef(0),
		},
	})
	td.db.SetObjectIndex(objID, &index.ObjectLifeline{
		State:       record.StateActivation,
		LatestState: objID,
	})
	deactivateCoreID, err := td.manager.DeactivateObject(
		ctx,
		*domainRef.CoreRef(),
		*td.requestRef.CoreRef(),
		&ObjectDescriptor{head: *genRefWithID(objID), state: *objID.CoreID()},
	)
	assert.Nil(t, err)
	deactivateID := record.Bytes2ID(deactivateCoreID[:])
	deactivateRec, err := td.db.GetRecord(&deactivateID)
	assert.Nil(t, err)
	assert.Equal(t, deactivateRec, &record.DeactivationRecord{
		ResultRecord: record.ResultRecord{
			Domain:  domainRef,
			Request: *td.requestRef,
		},
		PrevState: *objID,
	})
}

func TestLedgerArtifactManager_UpdateObject_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	ctx := inscontext.TODO()
	objID, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.ObjectActivateRecord{
		ResultRecord: record.ResultRecord{
			Domain: *genRandomRef(0),
		},
	})
	td.db.SetObjectIndex(objID, &index.ObjectLifeline{
		State:       record.StateActivation,
		LatestState: objID,
	})
	memory := []byte{1, 2, 3}
	obj, err := td.manager.UpdateObject(
		ctx,
		*domainRef.CoreRef(),
		*td.requestRef.CoreRef(),
		&ObjectDescriptor{head: *genRefWithID(objID), state: *objID.CoreID()},
		memory,
	)
	assert.Nil(t, err)
	updateID := record.Bytes2ID(obj.StateID()[:])
	updateRec, err := td.db.GetRecord(&updateID)
	assert.Nil(t, err)
	assert.Equal(t, updateRec, &record.ObjectAmendRecord{
		ResultRecord: record.ResultRecord{
			Domain:  domainRef,
			Request: *td.requestRef,
		},
		ObjectStateRecord: record.ObjectStateRecord{
			Memory: memory,
		},
		PrevState: *objID,
	})
}

func TestLedgerArtifactManager_GetClass_ReturnsCorrectDescriptors(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	ctx := inscontext.TODO()

	codeRef := *genRandomRef(0)
	classID, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.ClassActivateRecord{
		ResultRecord: record.ResultRecord{
			Domain: domainRef,
		},
	})
	classAmendID, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.ClassAmendRecord{
		ResultRecord: record.ResultRecord{
			Domain:  domainRef,
			Request: *td.requestRef,
		},
		ClassStateRecord: record.ClassStateRecord{
			Code: codeRef,
		},
	})
	classIndex := index.ClassLifeline{
		LatestState: classAmendID,
	}
	td.db.SetClassIndex(classID, &classIndex)

	classRef := genRefWithID(classID)
	classDesc, err := td.manager.GetClass(ctx, *classRef, nil)
	assert.NoError(t, err)
	expectedClassDesc := &ClassDescriptor{
		am:    td.manager,
		head:  *classRef,
		state: *classAmendID.CoreID(),
		code:  codeRef.CoreRef(),
	}

	assert.Equal(t, *expectedClassDesc, *classDesc.(*ClassDescriptor))
}

func TestLedgerArtifactManager_GetObject_VerifiesRecords(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	ctx := inscontext.TODO()

	objID := genRandomID(0)
	_, err := td.manager.GetObject(ctx, *genRefWithID(objID), nil, false)
	assert.NotNil(t, err)

	deactivateID, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.DeactivationRecord{})
	objectIndex := index.ObjectLifeline{
		LatestState: deactivateID,
		ClassRef:    *genRandomRef(0),
	}
	td.db.SetObjectIndex(objID, &objectIndex)

	_, err = td.manager.GetObject(ctx, *genRefWithID(objID), nil, false)
	assert.Equal(t, core.ErrDeactivated, err)
}

func TestLedgerArtifactManager_GetLatestObj_ReturnsCorrectDescriptors(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	ctx := inscontext.TODO()

	classID, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.ClassActivateRecord{
		ResultRecord: record.ResultRecord{
			Domain: domainRef,
		},
	})
	classAmendRef, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.ClassAmendRecord{
		ResultRecord: record.ResultRecord{
			Domain:  domainRef,
			Request: *td.requestRef,
		},
	})
	classIndex := index.ClassLifeline{
		LatestState: classAmendRef,
	}
	td.db.SetClassIndex(classID, &classIndex)

	objectID, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.ObjectActivateRecord{
		ResultRecord: record.ResultRecord{
			Domain: domainRef,
		},
		ObjectStateRecord: record.ObjectStateRecord{
			Memory: []byte{3},
		},
	})
	objectAmendID, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.ObjectAmendRecord{
		ResultRecord: record.ResultRecord{
			Domain: domainRef,
		},
		ObjectStateRecord: record.ObjectStateRecord{
			Memory: []byte{4},
		},
	})
	objectIndex := index.ObjectLifeline{
		LatestState: objectAmendID,
		ClassRef:    record.Reference{Domain: td.requestRef.Domain, Record: *classID},
	}
	td.db.SetObjectIndex(objectID, &objectIndex)

	objDesc, err := td.manager.GetObject(ctx, *genRefWithID(objectID), nil, false)
	assert.NoError(t, err)
	expectedObjDesc := &ObjectDescriptor{
		am: td.manager,

		head:   *getReference(td.requestRef.CoreRef(), objectID),
		state:  *objectAmendID.CoreID(),
		class:  *getReference(td.requestRef.CoreRef(), classID),
		memory: []byte{4},
	}

	assert.Equal(t, *expectedObjDesc, *objDesc.(*ObjectDescriptor))
}

func TestLedgerArtifactManager_GetChildren(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	ctx := inscontext.TODO()

	parentID, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.ObjectActivateRecord{
		ResultRecord: record.ResultRecord{
			Domain: domainRef,
		},
		ObjectStateRecord: record.ObjectStateRecord{
			Memory: []byte{0},
		},
	})
	child1Ref := genRandomRef(1)
	child2Ref := genRandomRef(1)
	child3Ref := genRandomRef(2)

	childMeta1, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.ChildRecord{
		Ref: *child1Ref,
	})
	childMeta2, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.ChildRecord{
		PrevChild: childMeta1,
		Ref:       *child2Ref,
	})
	childMeta3, _ := td.db.SetRecord(core.GenesisPulse.PulseNumber, &record.ChildRecord{
		PrevChild: childMeta2,
		Ref:       *child3Ref,
	})

	parentIndex := index.ObjectLifeline{
		LatestState:  parentID,
		ChildPointer: childMeta3,
	}
	td.db.SetObjectIndex(parentID, &parentIndex)

	t.Run("returns correct children without pulse", func(t *testing.T) {
		i, err := td.manager.GetChildren(ctx, *genRefWithID(parentID), nil)
		assert.NoError(t, err)
		child, err := i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child3Ref.CoreRef(), *child)
		child, err = i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child2Ref.CoreRef(), *child)
		child, err = i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child1Ref.CoreRef(), *child)
		hasNext := i.HasNext()
		assert.False(t, hasNext)
		_, err = i.Next()
		assert.Error(t, err)
	})

	t.Run("returns correct children with pulse", func(t *testing.T) {
		pn := core.PulseNumber(1)
		i, err := td.manager.GetChildren(ctx, *genRefWithID(parentID), &pn)
		assert.NoError(t, err)
		child, err := i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child2Ref.CoreRef(), *child)
		child, err = i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child1Ref.CoreRef(), *child)
		hasNext := i.HasNext()
		assert.NoError(t, err)
		assert.False(t, hasNext)
		_, err = i.Next()
		assert.Error(t, err)
	})

	t.Run("returns correct children in many chunks", func(t *testing.T) {
		td.manager.getChildrenChunkSize = 1
		i, err := td.manager.GetChildren(ctx, *genRefWithID(parentID), nil)
		assert.NoError(t, err)
		child, err := i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child3Ref.CoreRef(), *child)
		child, err = i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child2Ref.CoreRef(), *child)
		child, err = i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child1Ref.CoreRef(), *child)
		hasNext := i.HasNext()
		assert.NoError(t, err)
		assert.False(t, hasNext)
		_, err = i.Next()
		assert.Error(t, err)
	})

	t.Run("doesn't fail when has no children to return", func(t *testing.T) {
		td.manager.getChildrenChunkSize = 1
		pn := core.PulseNumber(3)
		i, err := td.manager.GetChildren(ctx, *genRefWithID(parentID), &pn)
		assert.NoError(t, err)
		child, err := i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child3Ref.CoreRef(), *child)
		child, err = i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child2Ref.CoreRef(), *child)
		child, err = i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child1Ref.CoreRef(), *child)
		hasNext := i.HasNext()
		assert.NoError(t, err)
		assert.False(t, hasNext)
		_, err = i.Next()
		assert.Error(t, err)
	})
}

func TestLedgerArtifactManager_HandleJetDrop(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()

	codeRecord := record.CodeRecord{
		Code: []byte{1, 2, 3, 3, 2, 1},
	}
	recHash := hash.NewIDHash()
	_, err := codeRecord.WriteHashData(recHash)
	assert.NoError(t, err)
	latestPulse, err := td.db.GetLatestPulseNumber()
	assert.NoError(t, err)
	id := record.ID{
		Pulse: latestPulse,
		Hash:  recHash.Sum(nil),
	}

	setRecordMessage := message.SetRecord{
		Record: record.SerializeRecord(&codeRecord),
	}
	messageBytes, err := message.ToBytes(&setRecordMessage)
	assert.NoError(t, err)

	rep, err := td.manager.messageBus.Send(
		inscontext.TODO(),
		&message.JetDrop{
			Messages: [][]byte{
				messageBytes,
			},
			PulseNumber: core.GenesisPulse.PulseNumber,
		},
	)
	assert.NoError(t, err)
	assert.Equal(t, reply.OK{}, *rep.(*reply.OK))

	rec, err := td.db.GetRecord(&id)
	assert.NoError(t, err)
	assert.Equal(t, codeRecord, *rec.(*record.CodeRecord))
}

func TestLedgerArtifactManager_RegisterValidation(t *testing.T) {
	t.Parallel()
	td, cleaner := prepareAMTestData(t)
	defer cleaner()
	ctx := inscontext.TODO()

	objCoreID, err := td.manager.RegisterRequest(ctx, &message.BootstrapRequest{Name: "object"})
	objID := record.Bytes2ID(objCoreID[:])
	objRef := genRefWithID(&objID)
	assert.NoError(t, err)

	desc, err := td.manager.ActivateObject(
		ctx,
		*domainRef.CoreRef(),
		*objRef,
		*genRandomRef(0).CoreRef(),
		*td.manager.GenesisRef(),
		false,
		[]byte{1},
	)
	assert.NoError(t, err)
	stateID1 := desc.StateID()

	desc, err = td.manager.GetObject(ctx, *objRef, nil, false)
	assert.NoError(t, err)
	assert.Equal(t, *stateID1, *desc.StateID())

	_, err = td.manager.GetObject(ctx, *objRef, nil, true)
	assert.Equal(t, err, core.ErrStateNotAvailable)

	desc, err = td.manager.UpdateObject(
		ctx,
		*domainRef.CoreRef(),
		*genRandomRef(0).CoreRef(),
		desc,
		[]byte{2},
	)
	assert.NoError(t, err)
	stateID2 := desc.StateID()

	desc, err = td.manager.GetObject(ctx, *objRef, nil, false)
	assert.NoError(t, err)
	desc, err = td.manager.UpdateObject(
		ctx,
		*domainRef.CoreRef(),
		*genRandomRef(0).CoreRef(),
		desc,
		[]byte{3},
	)
	assert.NoError(t, err)
	stateID3 := desc.StateID()
	err = td.manager.RegisterValidation(ctx, *objRef, *stateID2, true, nil)
	assert.NoError(t, err)

	desc, err = td.manager.GetObject(ctx, *objRef, nil, false)
	assert.NoError(t, err)
	assert.Equal(t, *stateID3, *desc.StateID())
	desc, err = td.manager.GetObject(ctx, *objRef, nil, true)
	assert.NoError(t, err)
	assert.Equal(t, *stateID2, *desc.StateID())
}
