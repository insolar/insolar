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
	"context"
	"math/rand"
	"sync"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/delegationtoken"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/index"
	"github.com/insolar/insolar/ledger/storage/jet"
	"github.com/insolar/insolar/ledger/storage/record"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/testmessagebus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	domainID   = *genRandomID(0)
	domainRef  = *core.NewRecordRef(domainID, domainID)
	requestRef = *genRandomRef(0)
)

func genRandomID(pulse core.PulseNumber) *core.RecordID {
	buff := [core.RecordIDSize - core.PulseNumberSize]byte{}
	_, err := rand.Read(buff[:])
	if err != nil {
		panic(err)
	}
	return core.NewRecordID(pulse, buff[:])
}

func genRefWithID(id *core.RecordID) *core.RecordRef {
	return core.NewRecordRef(domainID, *id)
}

func genRandomRef(pulse core.PulseNumber) *core.RecordRef {
	return genRefWithID(genRandomID(pulse))
}

func getTestData(t *testing.T) (
	context.Context,
	*storage.DB,
	*LedgerArtifactManager,
	func(), // cleaner
) {
	scheme := platformpolicy.NewPlatformCryptographyScheme()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	jc := testutils.NewJetCoordinatorMock(mc)
	mb := testmessagebus.NewTestMessageBus(t)
	db.PlatformCryptographyScheme = scheme
	handler := MessageHandler{
		db:                         db,
		replayHandlers:             map[core.MessageType]core.MessageHandler{},
		PlatformCryptographyScheme: scheme,
		conf: &configuration.Ledger{LightChainLimit: 3},
	}

	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	provideMock := recentstorage.NewProviderMock(t)
	provideMock.GetStorageFunc = func(p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	handler.RecentStorageProvider = provideMock

	jc.AmIMock.Return(true, nil)

	handler.Bus = mb
	handler.JetCoordinator = jc
	err := handler.Init(ctx)
	require.NoError(t, err)
	am := LedgerArtifactManager{
		db:                         db,
		DefaultBus:                 mb,
		getChildrenChunkSize:       100,
		PlatformCryptographyScheme: scheme,
	}

	return ctx, db, &am, cleaner
}

func TestLedgerArtifactManager_RegisterRequest(t *testing.T) {
	t.Parallel()
	ctx, db, am, cleaner := getTestData(t)
	defer cleaner()

	parcel := message.Parcel{Msg: &message.GenesisRequest{Name: "4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.4FFB8zfQoGznSmzDxwv4njX1aR9ioL8GHSH17QXH2AFa"}}
	id, err := am.RegisterRequest(ctx, &parcel)
	assert.NoError(t, err)
	rec, err := db.GetRecord(ctx, *jet.NewID(0, nil), id)
	assert.NoError(t, err)
	assert.Equal(t, message.ParcelToBytes(&parcel), rec.(*record.CallRequest).Payload)
}

func TestLedgerArtifactManager_GetCodeWithCache(t *testing.T) {
	t.Parallel()

	code := []byte("test_code")
	ctx := context.Background()
	codeRef := testutils.RandomRef()

	mb := testutils.NewMessageBusMock(t)
	mb.SendFunc = func(p context.Context, p1 core.Message, p2 core.Pulse, p3 *core.MessageSendOptions) (r core.Reply, r1 error) {
		return &reply.Code{
			Code: code,
		}, nil
	}

	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	am := LedgerArtifactManager{
		DefaultBus:    mb,
		db:            db,
		codeCacheLock: &sync.Mutex{},
		codeCache:     make(map[core.RecordRef]*cacheEntry),
	}

	desc, err := am.GetCode(ctx, codeRef)
	receivedCode, err := desc.Code()
	require.NoError(t, err)
	require.Equal(t, code, receivedCode)

	mb.SendFunc = func(p context.Context, p1 core.Message, p2 core.Pulse, p3 *core.MessageSendOptions) (r core.Reply, r1 error) {
		t.Fatal("Func must not be called here")
		return nil, nil
	}

	desc, err = am.GetCode(ctx, codeRef)
	receivedCode, err = desc.Code()
	require.NoError(t, err)
	require.Equal(t, code, receivedCode)

}

func TestLedgerArtifactManager_DeclareType(t *testing.T) {
	t.Parallel()
	ctx, db, am, cleaner := getTestData(t)
	defer cleaner()

	typeDec := []byte{1, 2, 3}
	id, err := am.DeclareType(ctx, domainRef, requestRef, typeDec)
	assert.NoError(t, err)
	typeRec, err := db.GetRecord(ctx, *jet.NewID(0, nil), id)
	assert.NoError(t, err)
	assert.Equal(t, &record.TypeRecord{
		SideEffectRecord: record.SideEffectRecord{
			Domain:  domainRef,
			Request: requestRef,
		},
		TypeDeclaration: typeDec,
	}, typeRec)
}

func TestLedgerArtifactManager_DeployCode_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	ctx, db, am, cleaner := getTestData(t)
	defer cleaner()

	id, err := am.DeployCode(
		ctx,
		domainRef,
		requestRef,
		[]byte{1, 2, 3},
		core.MachineTypeBuiltin,
	)
	assert.NoError(t, err)
	codeRec, err := db.GetRecord(ctx, *jet.NewID(0, nil), id)
	assert.NoError(t, err)
	assert.Equal(t, codeRec, &record.CodeRecord{
		SideEffectRecord: record.SideEffectRecord{
			Domain:  domainRef,
			Request: requestRef,
		},
		Code:        record.CalculateIDForBlob(am.PlatformCryptographyScheme, core.GenesisPulse.PulseNumber, []byte{1, 2, 3}),
		MachineType: core.MachineTypeBuiltin,
	})
}

func TestLedgerArtifactManager_ActivateObject_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	ctx, db, am, cleaner := getTestData(t)
	defer cleaner()
	jetID := *jet.NewID(0, nil)

	memory := []byte{1, 2, 3}
	codeRef := genRandomRef(0)
	parentID, _ := db.SetRecord(
		ctx,
		jetID,
		core.GenesisPulse.PulseNumber,
		&record.ObjectActivateRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain: *genRandomRef(0),
			},
		},
	)
	err := db.SetObjectIndex(ctx, jetID, parentID, &index.ObjectLifeline{
		LatestState: parentID,
	})
	require.NoError(t, err)

	objRef := *genRandomRef(0)
	objDesc, err := am.ActivateObject(
		ctx,
		domainRef,
		objRef,
		*genRefWithID(parentID),
		*codeRef,
		false,
		memory,
	)
	assert.Nil(t, err)
	activateRec, err := db.GetRecord(ctx, jetID, objDesc.StateID())
	assert.Nil(t, err)
	assert.Equal(t, activateRec, &record.ObjectActivateRecord{
		SideEffectRecord: record.SideEffectRecord{
			Domain:  domainRef,
			Request: objRef,
		},
		ObjectStateRecord: record.ObjectStateRecord{
			Memory:      record.CalculateIDForBlob(am.PlatformCryptographyScheme, core.GenesisPulse.PulseNumber, memory),
			Image:       *codeRef,
			IsPrototype: false,
		},
		Parent:     *genRefWithID(parentID),
		IsDelegate: false,
	})

	idx, err := db.GetObjectIndex(ctx, jetID, parentID, false)
	assert.NoError(t, err)
	childRec, err := db.GetRecord(ctx, jetID, idx.ChildPointer)
	assert.NoError(t, err)
	assert.Equal(t, objRef, childRec.(*record.ChildRecord).Ref)

	idx, err = db.GetObjectIndex(ctx, jetID, objRef.Record(), false)
	assert.NoError(t, err)
	assert.Equal(t, *objDesc.StateID(), *idx.LatestState)
	assert.Equal(t, *objDesc.Parent(), idx.Parent)
}

func TestLedgerArtifactManager_DeactivateObject_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	ctx, db, am, cleaner := getTestData(t)
	defer cleaner()
	jetID := *jet.NewID(0, nil)

	objID, _ := db.SetRecord(
		ctx,
		jetID,
		core.GenesisPulse.PulseNumber,
		&record.ObjectActivateRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain: *genRandomRef(0),
			},
		},
	)
	err := db.SetObjectIndex(ctx, jetID, objID, &index.ObjectLifeline{
		State:       record.StateActivation,
		LatestState: objID,
	})
	require.NoError(t, err)
	deactivateID, err := am.DeactivateObject(
		ctx,
		domainRef,
		requestRef,
		&ObjectDescriptor{
			ctx:   ctx,
			head:  *genRefWithID(objID),
			state: *objID,
		},
	)
	assert.Nil(t, err)
	deactivateRec, err := db.GetRecord(ctx, jetID, deactivateID)
	assert.Nil(t, err)
	assert.Equal(t, deactivateRec, &record.DeactivationRecord{
		SideEffectRecord: record.SideEffectRecord{
			Domain:  domainRef,
			Request: requestRef,
		},
		PrevState: *objID,
	})
}

func TestLedgerArtifactManager_UpdateObject_CreatesCorrectRecord(t *testing.T) {
	t.Parallel()
	ctx, db, am, cleaner := getTestData(t)
	defer cleaner()
	jetID := *jet.NewID(0, nil)

	objID, _ := db.SetRecord(
		ctx,
		jetID,
		core.GenesisPulse.PulseNumber,
		&record.ObjectActivateRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain: *genRandomRef(0),
			},
		},
	)
	err := db.SetObjectIndex(ctx, jetID, objID, &index.ObjectLifeline{
		State:       record.StateActivation,
		LatestState: objID,
	})
	require.NoError(t, err)
	memory := []byte{1, 2, 3}
	prototype := genRandomRef(0)
	obj, err := am.UpdateObject(
		ctx,
		domainRef,
		requestRef,
		&ObjectDescriptor{
			ctx:       ctx,
			head:      *genRefWithID(objID),
			state:     *objID,
			prototype: prototype,
		},
		memory,
	)
	assert.Nil(t, err)
	updateRec, err := db.GetRecord(ctx, jetID, obj.StateID())
	assert.Nil(t, err)
	assert.Equal(t, updateRec, &record.ObjectAmendRecord{
		SideEffectRecord: record.SideEffectRecord{
			Domain:  domainRef,
			Request: requestRef,
		},
		ObjectStateRecord: record.ObjectStateRecord{
			Memory:      record.CalculateIDForBlob(am.PlatformCryptographyScheme, core.GenesisPulse.PulseNumber, memory),
			Image:       *prototype,
			IsPrototype: false,
		},
		PrevState: *objID,
	})
}

func TestLedgerArtifactManager_GetObject_ReturnsCorrectDescriptors(t *testing.T) {
	t.Parallel()
	ctx, db, am, cleaner := getTestData(t)
	defer cleaner()
	jetID := *jet.NewID(0, nil)

	prototypeRef := genRandomRef(0)
	parentRef := genRandomRef(0)
	objRef := genRandomRef(0)
	_, err := db.SetRecord(
		ctx,
		jetID,
		core.GenesisPulse.PulseNumber,
		&record.ObjectActivateRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain: domainRef,
			},
			ObjectStateRecord: record.ObjectStateRecord{
				Memory: record.CalculateIDForBlob(am.PlatformCryptographyScheme, core.GenesisPulse.PulseNumber, []byte{3}),
			},
			Parent: *parentRef,
		},
	)
	require.NoError(t, err)
	_, err = db.SetBlob(ctx, jetID, core.GenesisPulse.PulseNumber, []byte{3})
	require.NoError(t, err)
	objectAmendID, _ := db.SetRecord(ctx, jetID, core.GenesisPulse.PulseNumber, &record.ObjectAmendRecord{
		SideEffectRecord: record.SideEffectRecord{
			Domain: domainRef,
		},
		ObjectStateRecord: record.ObjectStateRecord{
			Memory: record.CalculateIDForBlob(am.PlatformCryptographyScheme, core.GenesisPulse.PulseNumber, []byte{4}),
			Image:  *prototypeRef,
		},
	})
	_, err = db.SetBlob(ctx, jetID, core.GenesisPulse.PulseNumber, []byte{4})
	require.NoError(t, err)

	objectIndex := index.ObjectLifeline{
		LatestState:  objectAmendID,
		ChildPointer: genRandomID(0),
		Parent:       *parentRef,
	}
	require.NoError(
		t,
		db.SetObjectIndex(ctx, jetID, objRef.Record(), &objectIndex),
	)

	objDesc, err := am.GetObject(ctx, *objRef, nil, false)
	assert.NoError(t, err)
	expectedObjDesc := &ObjectDescriptor{
		ctx: ctx,
		am:  am,

		head:         *objRef,
		state:        *objectAmendID,
		prototype:    prototypeRef,
		isPrototype:  false,
		childPointer: objectIndex.ChildPointer,
		memory:       []byte{4},
		parent:       *parentRef,
	}

	assert.Equal(t, *expectedObjDesc, *objDesc.(*ObjectDescriptor))
}

func TestLedgerArtifactManager_GetObject_FollowsRedirect(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	am := NewArtifactManger(nil)
	mb := testutils.NewMessageBusMock(mc)

	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	objRef := genRandomRef(0)
	nodeRef := genRandomRef(0)
	mb.SendFunc = func(c context.Context, m core.Message, _ core.Pulse, o *core.MessageSendOptions) (r core.Reply, r1 error) {
		o = o.Safe()
		if o.Receiver == nil {
			return &reply.GetObjectRedirect{
				Receiver: nodeRef,
				Token:    &delegationtoken.GetObjectRedirect{Signature: []byte{1, 2, 3}},
			}, nil
		}

		token, ok := o.Token.(*delegationtoken.GetObjectRedirect)
		assert.True(t, ok)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, nodeRef, o.Receiver)
		return &reply.Object{}, nil
	}
	am.DefaultBus = mb
	am.db = db

	_, err := am.GetObject(ctx, *objRef, nil, false)

	require.NoError(t, err)
}

func TestLedgerArtifactManager_GetChildren(t *testing.T) {
	t.Parallel()
	ctx, db, am, cleaner := getTestData(t)
	defer cleaner()
	jetID := *jet.NewID(0, nil)

	parentID, _ := db.SetRecord(
		ctx,
		jetID,
		core.GenesisPulse.PulseNumber,
		&record.ObjectActivateRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain: domainRef,
			},
			ObjectStateRecord: record.ObjectStateRecord{
				Memory: record.CalculateIDForBlob(am.PlatformCryptographyScheme, core.GenesisPulse.PulseNumber, []byte{0}),
			},
		})
	child1Ref := genRandomRef(1)
	child2Ref := genRandomRef(1)
	child3Ref := genRandomRef(2)

	childMeta1, _ := db.SetRecord(
		ctx,
		jetID,
		core.GenesisPulse.PulseNumber,
		&record.ChildRecord{
			Ref: *child1Ref,
		})
	childMeta2, _ := db.SetRecord(
		ctx,
		jetID,
		core.GenesisPulse.PulseNumber,
		&record.ChildRecord{
			PrevChild: childMeta1,
			Ref:       *child2Ref,
		})
	childMeta3, _ := db.SetRecord(
		ctx,
		jetID,
		core.GenesisPulse.PulseNumber,
		&record.ChildRecord{
			PrevChild: childMeta2,
			Ref:       *child3Ref,
		})

	parentIndex := index.ObjectLifeline{
		LatestState:  parentID,
		ChildPointer: childMeta3,
	}
	require.NoError(
		t,
		db.SetObjectIndex(ctx, jetID, parentID, &parentIndex),
	)

	t.Run("returns correct children without pulse", func(t *testing.T) {
		i, err := am.GetChildren(ctx, *genRefWithID(parentID), nil)
		require.NoError(t, err)
		child, err := i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child3Ref, *child)
		child, err = i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child2Ref, *child)
		child, err = i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child1Ref, *child)
		hasNext := i.HasNext()
		assert.False(t, hasNext)
		_, err = i.Next()
		assert.Error(t, err)
	})

	t.Run("returns correct children with pulse", func(t *testing.T) {
		pn := core.PulseNumber(1)
		i, err := am.GetChildren(ctx, *genRefWithID(parentID), &pn)
		require.NoError(t, err)
		child, err := i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child2Ref, *child)
		child, err = i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child1Ref, *child)
		hasNext := i.HasNext()
		assert.NoError(t, err)
		assert.False(t, hasNext)
		_, err = i.Next()
		assert.Error(t, err)
	})

	t.Run("returns correct children in many chunks", func(t *testing.T) {
		am.getChildrenChunkSize = 1
		i, err := am.GetChildren(ctx, *genRefWithID(parentID), nil)
		require.NoError(t, err)
		child, err := i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child3Ref, *child)
		child, err = i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child2Ref, *child)
		child, err = i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child1Ref, *child)
		hasNext := i.HasNext()
		assert.NoError(t, err)
		assert.False(t, hasNext)
		_, err = i.Next()
		assert.Error(t, err)
	})

	t.Run("doesn't fail when has no children to return", func(t *testing.T) {
		am.getChildrenChunkSize = 1
		pn := core.PulseNumber(3)
		i, err := am.GetChildren(ctx, *genRefWithID(parentID), &pn)
		require.NoError(t, err)
		child, err := i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child3Ref, *child)
		child, err = i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child2Ref, *child)
		child, err = i.Next()
		assert.NoError(t, err)
		assert.Equal(t, *child1Ref, *child)
		hasNext := i.HasNext()
		assert.NoError(t, err)
		assert.False(t, hasNext)
		_, err = i.Next()
		assert.Error(t, err)
	})
}

func TestLedgerArtifactManager_GetChildren_FollowsRedirect(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	am := NewArtifactManger(nil)
	mb := testutils.NewMessageBusMock(mc)

	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()

	am.db = db

	objRef := genRandomRef(0)
	nodeRef := genRandomRef(0)
	mb.SendFunc = func(c context.Context, m core.Message, cp core.Pulse, o *core.MessageSendOptions) (r core.Reply, r1 error) {
		o = o.Safe()
		if o.Receiver == nil {
			return &reply.GetChildrenRedirect{
				Receiver: nodeRef,
				Token:    &delegationtoken.GetChildrenRedirect{Signature: []byte{1, 2, 3}},
			}, nil
		}

		token, ok := o.Token.(*delegationtoken.GetChildrenRedirect)
		assert.True(t, ok)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, nodeRef, o.Receiver)
		return &reply.Children{}, nil
	}
	am.DefaultBus = mb

	_, err := am.GetChildren(ctx, *objRef, nil)
	require.NoError(t, err)
}

func TestLedgerArtifactManager_HandleJetDrop(t *testing.T) {
	t.Parallel()
	ctx, db, am, cleaner := getTestData(t)
	defer cleaner()

	codeRecord := record.CodeRecord{
		Code: record.CalculateIDForBlob(am.PlatformCryptographyScheme, core.GenesisPulse.PulseNumber, []byte{1, 2, 3, 3, 2, 1}),
	}

	setRecordMessage := message.SetRecord{
		Record: record.SerializeRecord(&codeRecord),
	}

	jetID := *jet.NewID(0, nil)

	rep, err := am.DefaultBus.Send(
		ctx,
		&message.JetDrop{
			JetID: jetID,
			Messages: [][]byte{
				message.ToBytes(&setRecordMessage),
			},
			PulseNumber: core.GenesisPulse.PulseNumber,
		},
		*core.GenesisPulse,
		nil,
	)
	assert.NoError(t, err)
	assert.Equal(t, reply.OK{}, *rep.(*reply.OK))

	id := record.NewRecordIDFromRecord(db.PlatformCryptographyScheme, 0, &codeRecord)
	rec, err := db.GetRecord(ctx, jetID, id)
	require.NoError(t, err)
	assert.Equal(t, codeRecord, *rec.(*record.CodeRecord))
}

func TestLedgerArtifactManager_RegisterValidation(t *testing.T) {
	t.Parallel()
	scheme := platformpolicy.NewPlatformCryptographyScheme()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	defer mc.Finish()

	mb := testmessagebus.NewTestMessageBus(t)
	jc := testutils.NewJetCoordinatorMock(mc)

	recentStorageMock := recentstorage.NewRecentStorageMock(t)
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()

	handler := MessageHandler{
		db:                         db,
		replayHandlers:             map[core.MessageType]core.MessageHandler{},
		PlatformCryptographyScheme: scheme,
		conf: &configuration.Ledger{LightChainLimit: 3},
	}

	handler.Bus = mb
	handler.JetCoordinator = jc

	provideMock := recentstorage.NewProviderMock(t)
	provideMock.GetStorageFunc = func(p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	handler.RecentStorageProvider = provideMock

	err := handler.Init(ctx)
	require.NoError(t, err)
	am := LedgerArtifactManager{
		db:                         db,
		DefaultBus:                 mb,
		getChildrenChunkSize:       100,
		PlatformCryptographyScheme: scheme,
	}

	jc.QueryRoleMock.Return([]core.RecordRef{*genRandomRef(0)}, nil)
	jc.AmIMock.Return(true, nil)

	objID, err := am.RegisterRequest(
		ctx,
		&message.Parcel{
			Msg: &message.GenesisRequest{
				Name: "4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.4FFB8zfQoGznSmzDxwv4njX1aR9ioL8GHSH17QXH2AFa",
			},
		},
	)
	require.NoError(t, err)
	objRef := genRefWithID(objID)

	desc, err := am.ActivateObject(
		ctx,
		domainRef,
		*objRef,
		*am.GenesisRef(),
		*genRandomRef(0),
		false,
		[]byte{1},
	)
	require.NoError(t, err)
	stateID1 := desc.StateID()

	desc, err = am.GetObject(ctx, *objRef, nil, false)
	require.NoError(t, err)
	require.Equal(t, *stateID1, *desc.StateID())

	_, err = am.GetObject(ctx, *objRef, nil, true)
	require.Equal(t, err, core.ErrStateNotAvailable)

	desc, err = am.GetObject(ctx, *objRef, nil, false)
	require.NoError(t, err)
	desc, err = am.UpdateObject(
		ctx,
		domainRef,
		*genRandomRef(0),
		desc,
		[]byte{3},
	)
	require.NoError(t, err)
	stateID3 := desc.StateID()
	err = am.RegisterValidation(ctx, *objRef, *stateID1, true, nil)
	require.NoError(t, err)

	desc, err = am.GetObject(ctx, *objRef, nil, false)
	assert.NoError(t, err)
	assert.Equal(t, *stateID3, *desc.StateID())
	desc, err = am.GetObject(ctx, *objRef, nil, true)
	assert.NoError(t, err)
	assert.Equal(t, *stateID1, *desc.StateID())
}

func TestLedgerArtifactManager_RegisterResult(t *testing.T) {
	t.Parallel()
	ctx, db, am, cleaner := getTestData(t)
	defer cleaner()

	request := genRandomRef(0)
	requestID, err := am.RegisterResult(ctx, *request, []byte{1, 2, 3})
	assert.NoError(t, err)

	rec, err := db.GetRecord(ctx, *jet.NewID(0, nil), requestID)
	assert.NoError(t, err)
	assert.Equal(t, record.ResultRecord{Request: *request, Payload: []byte{1, 2, 3}}, *rec.(*record.ResultRecord))
}

func TestLedgerArtifactManager_RegisterRequest_JetMiss(t *testing.T) {
	t.Parallel()
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)
	db, cleaner := storagetest.TmpDB(ctx, t)
	defer cleaner()
	defer mc.Finish()

	cs := testutils.NewPlatformCryptographyScheme()
	am := NewArtifactManger(db)
	am.PlatformCryptographyScheme = cs

	t.Run("returns error on exceeding retry limit", func(t *testing.T) {
		mb := testutils.NewMessageBusMock(mc)
		am.DefaultBus = mb
		mb.SendMock.Return(&reply.JetMiss{JetID: *jet.NewID(5, []byte{1, 2, 3})}, nil)
		_, err := am.RegisterRequest(ctx, &message.Parcel{Msg: &message.CallMethod{}})
		require.Error(t, err)
	})

	t.Run("returns no error and updates tree when jet miss", func(t *testing.T) {
		mb := testutils.NewMessageBusMock(mc)
		am.DefaultBus = mb
		retries := 3
		mb.SendFunc = func(c context.Context, m core.Message, p core.Pulse, o *core.MessageSendOptions) (r core.Reply, r1 error) {
			if retries == 0 {
				return &reply.ID{}, nil
			}
			retries--
			return &reply.JetMiss{JetID: *jet.NewID(4, []byte{0xD5})}, nil
		}
		_, err := am.RegisterRequest(ctx, &message.Parcel{Msg: &message.CallMethod{}})
		require.NoError(t, err)

		tree, err := db.GetJetTree(ctx, core.FirstPulseNumber)
		require.NoError(t, err)
		jetID := tree.Find(*core.NewRecordID(0, []byte{0xD5}))
		assert.Equal(t, *jet.NewID(4, []byte{0xD0}), *jetID)
	})
}
