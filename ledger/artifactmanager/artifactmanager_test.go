/*
 *    Copyright 2019 Insolar Technologies
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
	"testing"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/ledger/storage/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/component"
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
)

type amSuite struct {
	suite.Suite

	cm      *component.Manager
	ctx     context.Context
	cleaner func()
	db      storage.DBContext

	scheme        core.PlatformCryptographyScheme
	pulseTracker  storage.PulseTracker
	nodeStorage   node.Accessor
	objectStorage storage.ObjectStorage
	jetStorage    storage.JetStorage
	dropStorage   storage.DropStorage
	genesisState  storage.GenesisState
}

func NewAmSuite() *amSuite {
	return &amSuite{
		Suite: suite.Suite{},
	}
}

// Init and run suite
func TestArtifactManager(t *testing.T) {
	suite.Run(t, NewAmSuite())
}

func (s *amSuite) BeforeTest(suiteName, testName string) {
	s.cm = &component.Manager{}
	s.ctx = inslogger.TestContext(s.T())

	db, cleaner := storagetest.TmpDB(s.ctx, s.T())
	s.cleaner = cleaner
	s.db = db
	s.scheme = platformpolicy.NewPlatformCryptographyScheme()
	s.jetStorage = storage.NewJetStorage()
	s.nodeStorage = node.NewStorage()
	s.pulseTracker = storage.NewPulseTracker()
	s.objectStorage = storage.NewObjectStorage()
	s.dropStorage = storage.NewDropStorage(10)
	s.genesisState = storage.NewGenesisInitializer()

	s.cm.Inject(
		s.scheme,
		s.db,
		s.jetStorage,
		s.nodeStorage,
		s.pulseTracker,
		s.objectStorage,
		s.dropStorage,
		s.genesisState,
	)

	err := s.cm.Init(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager init failed", err)
	}
	err = s.cm.Start(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager start failed", err)
	}
}

func (s *amSuite) AfterTest(suiteName, testName string) {
	err := s.cm.Stop(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager stop failed", err)
	}
	s.cleaner()
}

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

func getTestData(s *amSuite) (
	context.Context,
	storage.ObjectStorage,
	*LedgerArtifactManager,
) {
	mc := minimock.NewController(s.T())
	pulseStorage := storage.NewPulseStorage()
	pulseStorage.PulseTracker = s.pulseTracker

	pulse, err := s.pulseTracker.GetLatestPulse(s.ctx)
	require.NoError(s.T(), err)
	pulseStorage.Set(&pulse.Pulse)

	mb := testmessagebus.NewTestMessageBus(s.T())
	mb.PulseStorage = pulseStorage

	certificate := testutils.NewCertificateMock(s.T())
	certificate.GetRoleMock.Return(core.StaticRoleLightMaterial)

	handler := MessageHandler{
		replayHandlers:             map[core.MessageType]core.MessageHandler{},
		PlatformCryptographyScheme: s.scheme,
		conf:                       &configuration.Ledger{LightChainLimit: 3, PendingRequestsLimit: 10},
		certificate:                certificate,
	}

	handler.Nodes = s.nodeStorage
	handler.ObjectStorage = s.objectStorage
	handler.PulseTracker = s.pulseTracker
	handler.DBContext = s.db
	handler.JetStorage = s.jetStorage

	indexMock := recentstorage.NewRecentIndexStorageMock(s.T())
	pendingMock := recentstorage.NewPendingStorageMock(s.T())

	indexMock.AddObjectMock.Return()
	pendingMock.GetRequestsForObjectMock.Return(nil)
	pendingMock.AddPendingRequestMock.Return()
	pendingMock.RemovePendingRequestMock.Return()

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetIndexStorageMock.Return(indexMock)
	provideMock.GetPendingStorageMock.Return(pendingMock)
	provideMock.CountMock.Return(1)

	handler.RecentStorageProvider = provideMock

	handler.Bus = mb

	jc := testutils.NewJetCoordinatorMock(mc)
	jc.LightExecutorForJetMock.Return(&core.RecordRef{}, nil)
	jc.MeMock.Return(core.RecordRef{})
	jc.HeavyMock.Return(&core.RecordRef{}, nil)
	jc.NodeForJetMock.Return(&core.RecordRef{}, nil)
	jc.IsBeyondLimitMock.Return(false, nil)

	handler.JetCoordinator = jc

	err = handler.Init(s.ctx)
	require.NoError(s.T(), err)

	am := LedgerArtifactManager{
		DB:                         s.db,
		DefaultBus:                 mb,
		getChildrenChunkSize:       100,
		PlatformCryptographyScheme: s.scheme,
		PulseStorage:               pulseStorage,
		GenesisState:               s.genesisState,
	}

	return s.ctx, s.objectStorage, &am
}

func (s *amSuite) TestLedgerArtifactManager_RegisterRequest() {
	ctx, os, am := getTestData(s)

	parcel := message.Parcel{Msg: &message.GenesisRequest{Name: "4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.4FFB8zfQoGznSmzDxwv4njX1aR9ioL8GHSH17QXH2AFa"}}
	id, err := am.RegisterRequest(ctx, *am.GenesisRef(), &parcel)
	assert.NoError(s.T(), err)
	rec, err := os.GetRecord(ctx, *jet.NewID(0, nil), id)
	assert.NoError(s.T(), err)

	assert.Equal(
		s.T(),
		am.PlatformCryptographyScheme.IntegrityHasher().Hash(message.MustSerializeBytes(parcel.Msg)),
		rec.(*record.RequestRecord).MessageHash,
	)
}

func (s *amSuite) TestLedgerArtifactManager_GetCodeWithCache() {
	code := []byte("test_code")
	codeRef := testutils.RandomRef()

	mb := testutils.NewMessageBusMock(s.T())
	mb.SendFunc = func(p context.Context, p1 core.Message, p3 *core.MessageSendOptions) (r core.Reply, r1 error) {
		return &reply.Code{
			Code: code,
		}, nil
	}

	jc := testutils.NewJetCoordinatorMock(s.T())
	jc.LightExecutorForJetMock.Return(&core.RecordRef{}, nil)
	jc.MeMock.Return(core.RecordRef{})

	amPulseStorageMock := testutils.NewPulseStorageMock(s.T())
	amPulseStorageMock.CurrentFunc = func(p context.Context) (r *core.Pulse, r1 error) {
		pulse, err := s.pulseTracker.GetLatestPulse(p)
		require.NoError(s.T(), err)
		return &pulse.Pulse, err
	}

	am := LedgerArtifactManager{
		DefaultBus:                 mb,
		DB:                         s.db,
		PulseStorage:               amPulseStorageMock,
		JetCoordinator:             jc,
		PlatformCryptographyScheme: s.scheme,
		senders:                    newLedgerArtifactSenders(),
	}

	desc, err := am.GetCode(s.ctx, codeRef)
	receivedCode, err := desc.Code()
	require.NoError(s.T(), err)
	require.Equal(s.T(), code, receivedCode)

	mb.SendFunc = func(p context.Context, p1 core.Message, p3 *core.MessageSendOptions) (r core.Reply, r1 error) {
		s.T().Fatal("Func must not be called here")
		return nil, nil
	}

	desc, err = am.GetCode(s.ctx, codeRef)
	receivedCode, err = desc.Code()
	require.NoError(s.T(), err)
	require.Equal(s.T(), code, receivedCode)

}

func (s *amSuite) TestLedgerArtifactManager_DeclareType() {
	ctx, os, am := getTestData(s)

	typeDec := []byte{1, 2, 3}
	id, err := am.DeclareType(ctx, domainRef, requestRef, typeDec)
	assert.NoError(s.T(), err)
	typeRec, err := os.GetRecord(ctx, *jet.NewID(0, nil), id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), &record.TypeRecord{
		SideEffectRecord: record.SideEffectRecord{
			Domain:  domainRef,
			Request: requestRef,
		},
		TypeDeclaration: typeDec,
	}, typeRec)
}

func (s *amSuite) TestLedgerArtifactManager_DeployCode_CreatesCorrectRecord() {
	ctx, os, am := getTestData(s)

	id, err := am.DeployCode(
		ctx,
		domainRef,
		requestRef,
		[]byte{1, 2, 3},
		core.MachineTypeBuiltin,
	)
	assert.NoError(s.T(), err)
	codeRec, err := os.GetRecord(ctx, *jet.NewID(0, nil), id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), codeRec, &record.CodeRecord{
		SideEffectRecord: record.SideEffectRecord{
			Domain:  domainRef,
			Request: requestRef,
		},
		Code:        record.CalculateIDForBlob(am.PlatformCryptographyScheme, core.GenesisPulse.PulseNumber, []byte{1, 2, 3}),
		MachineType: core.MachineTypeBuiltin,
	})
}

func (s *amSuite) TestLedgerArtifactManager_ActivateObject_CreatesCorrectRecord() {
	ctx, os, am := getTestData(s)
	jetID := *jet.NewID(0, nil)

	memory := []byte{1, 2, 3}
	codeRef := genRandomRef(0)
	parentID, _ := os.SetRecord(
		ctx,
		jetID,
		core.GenesisPulse.PulseNumber,
		&record.ObjectActivateRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain: *genRandomRef(0),
			},
		},
	)
	err := os.SetObjectIndex(ctx, jetID, parentID, &index.ObjectLifeline{
		LatestState: parentID,
	})
	require.NoError(s.T(), err)

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
	assert.Nil(s.T(), err)
	activateRec, err := os.GetRecord(ctx, jetID, objDesc.StateID())
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), activateRec, &record.ObjectActivateRecord{
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

	idx, err := os.GetObjectIndex(ctx, jetID, parentID, false)
	assert.NoError(s.T(), err)
	childRec, err := os.GetRecord(ctx, jetID, idx.ChildPointer)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), objRef, childRec.(*record.ChildRecord).Ref)

	idx, err = os.GetObjectIndex(ctx, jetID, objRef.Record(), false)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), *objDesc.StateID(), *idx.LatestState)
	assert.Equal(s.T(), *objDesc.Parent(), idx.Parent)
}

func (s *amSuite) TestLedgerArtifactManager_DeactivateObject_CreatesCorrectRecord() {
	ctx, os, am := getTestData(s)
	jetID := *jet.NewID(0, nil)

	objID, _ := os.SetRecord(
		ctx,
		jetID,
		core.GenesisPulse.PulseNumber,
		&record.ObjectActivateRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain: *genRandomRef(0),
			},
		},
	)
	err := os.SetObjectIndex(ctx, jetID, objID, &index.ObjectLifeline{
		State:       record.StateActivation,
		LatestState: objID,
	})
	require.NoError(s.T(), err)
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
	assert.Nil(s.T(), err)
	deactivateRec, err := os.GetRecord(ctx, jetID, deactivateID)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), deactivateRec, &record.DeactivationRecord{
		SideEffectRecord: record.SideEffectRecord{
			Domain:  domainRef,
			Request: requestRef,
		},
		PrevState: *objID,
	})
}

func (s *amSuite) TestLedgerArtifactManager_UpdateObject_CreatesCorrectRecord() {
	ctx, os, am := getTestData(s)
	jetID := *jet.NewID(0, nil)

	objID, _ := os.SetRecord(
		ctx,
		jetID,
		core.GenesisPulse.PulseNumber,
		&record.ObjectActivateRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain: *genRandomRef(0),
			},
		},
	)
	err := os.SetObjectIndex(ctx, jetID, objID, &index.ObjectLifeline{
		State:       record.StateActivation,
		LatestState: objID,
	})
	require.NoError(s.T(), err)
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
	assert.Nil(s.T(), err)
	updateRec, err := os.GetRecord(ctx, jetID, obj.StateID())
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), updateRec, &record.ObjectAmendRecord{
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

func (s *amSuite) TestLedgerArtifactManager_GetObject_ReturnsCorrectDescriptors() {
	ctx, os, am := getTestData(s)
	jetID := *jet.NewID(0, nil)

	prototypeRef := genRandomRef(0)
	parentRef := genRandomRef(0)
	objRef := genRandomRef(0)
	_, err := os.SetRecord(
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
	require.NoError(s.T(), err)
	_, err = os.SetBlob(ctx, jetID, core.GenesisPulse.PulseNumber, []byte{3})
	require.NoError(s.T(), err)
	objectAmendID, _ := os.SetRecord(ctx, jetID, core.GenesisPulse.PulseNumber, &record.ObjectAmendRecord{
		SideEffectRecord: record.SideEffectRecord{
			Domain: domainRef,
		},
		ObjectStateRecord: record.ObjectStateRecord{
			Memory: record.CalculateIDForBlob(am.PlatformCryptographyScheme, core.GenesisPulse.PulseNumber, []byte{4}),
			Image:  *prototypeRef,
		},
	})
	_, err = os.SetBlob(ctx, jetID, core.GenesisPulse.PulseNumber, []byte{4})
	require.NoError(s.T(), err)

	objectIndex := index.ObjectLifeline{
		LatestState:  objectAmendID,
		ChildPointer: genRandomID(0),
		Parent:       *parentRef,
	}
	require.NoError(
		s.T(),
		os.SetObjectIndex(ctx, jetID, objRef.Record(), &objectIndex),
	)

	objDesc, err := am.GetObject(ctx, *objRef, nil, false)
	rObjDesc := objDesc.(*ObjectDescriptor)
	assert.NoError(s.T(), err)
	expectedObjDesc := &ObjectDescriptor{
		ctx:          rObjDesc.ctx,
		am:           am,
		head:         *objRef,
		state:        *objectAmendID,
		prototype:    prototypeRef,
		isPrototype:  false,
		childPointer: objectIndex.ChildPointer,
		memory:       []byte{4},
		parent:       *parentRef,
	}
	assert.Equal(s.T(), *expectedObjDesc, *rObjDesc)
}

func (s *amSuite) TestLedgerArtifactManager_GetObject_FollowsRedirect() {
	mc := minimock.NewController(s.T())
	am := NewArtifactManger()
	mb := testutils.NewMessageBusMock(mc)

	objRef := genRandomRef(0)
	nodeRef := genRandomRef(0)
	mb.SendFunc = func(c context.Context, m core.Message, o *core.MessageSendOptions) (r core.Reply, r1 error) {
		o = o.Safe()

		switch m.(type) {
		case *message.GetObjectIndex:
			return &reply.ObjectIndex{}, nil
		case *message.GetObject:
			if o.Receiver == nil {
				return &reply.GetObjectRedirectReply{
					Receiver: nodeRef,
					Token:    &delegationtoken.GetObjectRedirectToken{Signature: []byte{1, 2, 3}},
				}, nil
			}

			token, ok := o.Token.(*delegationtoken.GetObjectRedirectToken)
			assert.True(s.T(), ok)
			assert.Equal(s.T(), []byte{1, 2, 3}, token.Signature)
			assert.Equal(s.T(), nodeRef, o.Receiver)
			return &reply.Object{}, nil
		default:
			panic("unexpected call")
		}
	}
	am.DefaultBus = mb
	am.DB = s.db
	am.PulseStorage = makePulseStorage(s)

	_, err := am.GetObject(s.ctx, *objRef, nil, false)

	require.NoError(s.T(), err)
}

func (s *amSuite) TestLedgerArtifactManager_GetChildren() {
	// t.Parallel()
	ctx, os, am := getTestData(s)
	// defer cleaner()
	jetID := *jet.NewID(0, nil)

	parentID, _ := os.SetRecord(
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

	childMeta1, _ := os.SetRecord(
		ctx,
		jetID,
		core.GenesisPulse.PulseNumber,
		&record.ChildRecord{
			Ref: *child1Ref,
		})
	childMeta2, _ := os.SetRecord(
		ctx,
		jetID,
		core.GenesisPulse.PulseNumber,
		&record.ChildRecord{
			PrevChild: childMeta1,
			Ref:       *child2Ref,
		})
	childMeta3, _ := os.SetRecord(
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
		s.T(),
		os.SetObjectIndex(ctx, jetID, parentID, &parentIndex),
	)

	s.T().Run("returns correct children without pulse", func(t *testing.T) {
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

	s.T().Run("returns correct children with pulse", func(t *testing.T) {
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

	s.T().Run("returns correct children in many chunks", func(t *testing.T) {
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

	s.T().Run("doesn't fail when has no children to return", func(t *testing.T) {
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

func makePulseStorage(s *amSuite) core.PulseStorage {
	pulseStorage := storage.NewPulseStorage()
	pulseStorage.PulseTracker = s.pulseTracker
	pulse, err := s.pulseTracker.GetLatestPulse(s.ctx)
	require.NoError(s.T(), err)
	pulseStorage.Set(&pulse.Pulse)

	return pulseStorage
}

func (s *amSuite) TestLedgerArtifactManager_GetChildren_FollowsRedirect() {
	mc := minimock.NewController(s.T())
	am := NewArtifactManger()
	mb := testutils.NewMessageBusMock(mc)

	am.DB = s.db
	am.PulseStorage = makePulseStorage(s)

	objRef := genRandomRef(0)
	nodeRef := genRandomRef(0)
	mb.SendFunc = func(c context.Context, m core.Message, o *core.MessageSendOptions) (r core.Reply, r1 error) {
		o = o.Safe()
		if o.Receiver == nil {
			return &reply.GetChildrenRedirectReply{
				Receiver: nodeRef,
				Token:    &delegationtoken.GetChildrenRedirectToken{Signature: []byte{1, 2, 3}},
			}, nil
		}

		token, ok := o.Token.(*delegationtoken.GetChildrenRedirectToken)
		assert.True(s.T(), ok)
		assert.Equal(s.T(), []byte{1, 2, 3}, token.Signature)
		assert.Equal(s.T(), nodeRef, o.Receiver)
		return &reply.Children{}, nil
	}
	am.DefaultBus = mb

	_, err := am.GetChildren(s.ctx, *objRef, nil)
	require.NoError(s.T(), err)
}

func (s *amSuite) TestLedgerArtifactManager_HandleJetDrop() {
	s.T().Skip("jet drops are for validation and it doesn't work")

	ctx, os, am := getTestData(s)

	codeRecord := record.CodeRecord{
		Code: record.CalculateIDForBlob(am.PlatformCryptographyScheme, core.GenesisPulse.PulseNumber, []byte{1, 2, 3, 3, 2, 1}),
	}

	setRecordMessage := message.SetRecord{
		Record: record.SerializeRecord(&codeRecord),
	}

	jetID := *jet.NewID(0, nil)

	rep, err := am.DefaultBus.Send(ctx, &message.JetDrop{
		JetID: jetID,
		Messages: [][]byte{
			message.ToBytes(&setRecordMessage),
		},
		PulseNumber: core.GenesisPulse.PulseNumber,
	}, nil)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), reply.OK{}, *rep.(*reply.OK))

	id := record.NewRecordIDFromRecord(s.scheme, 0, &codeRecord)
	rec, err := os.GetRecord(ctx, jetID, id)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), codeRecord, *rec.(*record.CodeRecord))
}

func (s *amSuite) TestLedgerArtifactManager_RegisterValidation() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()

	mb := testmessagebus.NewTestMessageBus(s.T())
	mb.PulseStorage = makePulseStorage(s)
	jc := testutils.NewJetCoordinatorMock(mc)
	jc.IsBeyondLimitMock.Return(false, nil)
	jc.NodeForJetMock.Return(&core.RecordRef{}, nil)

	indexMock := recentstorage.NewRecentIndexStorageMock(s.T())
	pendingMock := recentstorage.NewPendingStorageMock(s.T())

	indexMock.AddObjectMock.Return()
	pendingMock.GetRequestsForObjectMock.Return(nil)
	pendingMock.AddPendingRequestMock.Return()
	pendingMock.RemovePendingRequestMock.Return()

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetIndexStorageMock.Return(indexMock)
	provideMock.GetPendingStorageMock.Return(pendingMock)
	provideMock.CountMock.Return(0)

	certificate := testutils.NewCertificateMock(s.T())
	certificate.GetRoleMock.Return(core.StaticRoleLightMaterial)

	handler := MessageHandler{
		replayHandlers:             map[core.MessageType]core.MessageHandler{},
		PlatformCryptographyScheme: s.scheme,
		conf:                       &configuration.Ledger{LightChainLimit: 3, PendingRequestsLimit: 10},
		certificate:                certificate,
	}

	handler.Bus = mb
	handler.JetCoordinator = jc
	handler.DBContext = s.db
	handler.ObjectStorage = s.objectStorage
	handler.PulseTracker = s.pulseTracker
	handler.Nodes = s.nodeStorage
	handler.JetStorage = s.jetStorage

	handler.RecentStorageProvider = provideMock

	err := handler.Init(s.ctx)
	require.NoError(s.T(), err)

	amPulseStorageMock := testutils.NewPulseStorageMock(s.T())
	amPulseStorageMock.CurrentFunc = func(p context.Context) (r *core.Pulse, r1 error) {
		pulse, err := s.pulseTracker.GetLatestPulse(p)
		require.NoError(s.T(), err)
		return &pulse.Pulse, err
	}

	am := LedgerArtifactManager{
		DB:                         s.db,
		DefaultBus:                 mb,
		getChildrenChunkSize:       100,
		PlatformCryptographyScheme: s.scheme,
		PulseStorage:               amPulseStorageMock,
		GenesisState:               s.genesisState,
	}

	objID, err := am.RegisterRequest(
		s.ctx,
		*am.GenesisRef(),
		&message.Parcel{
			Msg: &message.GenesisRequest{
				Name: "4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.4FFB8zfQoGznSmzDxwv4njX1aR9ioL8GHSH17QXH2AFa",
			},
		},
	)
	require.NoError(s.T(), err)
	objRef := genRefWithID(objID)

	desc, err := am.ActivateObject(
		s.ctx,
		domainRef,
		*objRef,
		*am.GenesisRef(),
		*genRandomRef(0),
		false,
		[]byte{1},
	)
	require.NoError(s.T(), err)
	stateID1 := desc.StateID()

	desc, err = am.GetObject(s.ctx, *objRef, nil, false)
	require.NoError(s.T(), err)
	require.Equal(s.T(), *stateID1, *desc.StateID())

	_, err = am.GetObject(s.ctx, *objRef, nil, true)
	require.Equal(s.T(), err, core.ErrStateNotAvailable)

	desc, err = am.GetObject(s.ctx, *objRef, nil, false)
	require.NoError(s.T(), err)
	desc, err = am.UpdateObject(
		s.ctx,
		domainRef,
		*genRandomRef(0),
		desc,
		[]byte{3},
	)
	require.NoError(s.T(), err)
	stateID3 := desc.StateID()
	err = am.RegisterValidation(s.ctx, *objRef, *stateID1, true, nil)
	require.NoError(s.T(), err)

	desc, err = am.GetObject(s.ctx, *objRef, nil, false)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), *stateID3, *desc.StateID())
	desc, err = am.GetObject(s.ctx, *objRef, nil, true)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), *stateID1, *desc.StateID())
}

func (s *amSuite) TestLedgerArtifactManager_RegisterResult() {
	ctx, os, am := getTestData(s)

	objID := core.RecordID{1, 2, 3}
	request := genRandomRef(0)
	requestID, err := am.RegisterResult(ctx, *core.NewRecordRef(core.RecordID{}, objID), *request, []byte{1, 2, 3})
	assert.NoError(s.T(), err)

	rec, err := os.GetRecord(ctx, *jet.NewID(0, nil), requestID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), record.ResultRecord{
		Object:  objID,
		Request: *request,
		Payload: []byte{1, 2, 3},
	}, *rec.(*record.ResultRecord))
}

func (s *amSuite) TestLedgerArtifactManager_RegisterRequest_JetMiss() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()

	cs := platformpolicy.NewPlatformCryptographyScheme()
	am := NewArtifactManger()
	am.PlatformCryptographyScheme = cs
	pulseStorageMock := testutils.NewPulseStorageMock(s.T())
	pulseStorageMock.CurrentFunc = func(ctx context.Context) (*core.Pulse, error) {
		return &core.Pulse{PulseNumber: core.FirstPulseNumber}, nil
	}

	am.PulseStorage = pulseStorageMock
	am.GenesisState = s.genesisState
	am.JetStorage = s.jetStorage

	s.T().Run("returns error on exceeding retry limit", func(t *testing.T) {
		mb := testutils.NewMessageBusMock(mc)
		am.DefaultBus = mb
		mb.SendMock.Return(&reply.JetMiss{JetID: *jet.NewID(5, []byte{1, 2, 3})}, nil)
		_, err := am.RegisterRequest(s.ctx, *am.GenesisRef(), &message.Parcel{Msg: &message.CallMethod{}})
		require.Error(t, err)
	})

	s.T().Run("returns no error and updates tree when jet miss", func(t *testing.T) {
		mb := testutils.NewMessageBusMock(mc)
		am.DefaultBus = mb
		retries := 3
		mb.SendFunc = func(c context.Context, m core.Message, o *core.MessageSendOptions) (r core.Reply, r1 error) {
			if retries == 0 {
				return &reply.ID{}, nil
			}
			retries--
			return &reply.JetMiss{JetID: *jet.NewID(4, []byte{0xD5})}, nil
		}
		_, err := am.RegisterRequest(s.ctx, *am.GenesisRef(), &message.Parcel{Msg: &message.CallMethod{}})
		require.NoError(t, err)

		jetID, actual := s.jetStorage.FindJet(
			s.ctx, core.FirstPulseNumber, *core.NewRecordID(0, []byte{0xD5}),
		)
		assert.Equal(t, *jet.NewID(4, []byte{0xD0}), *jetID)
		assert.True(t, actual)
	})
}

func (s *amSuite) TestLedgerArtifactManager_GetRequest_Success() {
	// Arrange
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	objectID := testutils.RandomID()
	requestID := testutils.RandomID()

	node := testutils.RandomRef()

	jc := testutils.NewJetCoordinatorMock(mc)
	jc.NodeForObjectMock.Return(&node, nil)

	pulseStorageMock := testutils.NewPulseStorageMock(mc)
	pulseStorageMock.CurrentMock.Return(core.GenesisPulse, nil)

	var parcel core.Parcel = &message.Parcel{PulseNumber: 123987}
	resRecord := record.RequestRecord{
		Parcel: message.ParcelToBytes(parcel),
	}
	finalResponse := &reply.Request{Record: record.SerializeRecord(&resRecord)}

	mb := testutils.NewMessageBusMock(s.T())
	mb.SendFunc = func(p context.Context, p1 core.Message, p2 *core.MessageSendOptions) (r core.Reply, r1 error) {
		switch mb.SendCounter {
		case 0:
			casted, ok := p1.(*message.GetPendingRequestID)
			require.Equal(s.T(), true, ok)
			require.Equal(s.T(), objectID, casted.ObjectID)
			return &reply.ID{ID: requestID}, nil
		case 1:
			casted, ok := p1.(*message.GetRequest)
			require.Equal(s.T(), true, ok)
			require.Equal(s.T(), requestID, casted.Request)
			require.Equal(s.T(), node, *p2.Receiver)
			return finalResponse, nil
		default:
			panic("test is totally broken")
		}
	}

	am := NewArtifactManger()
	am.JetCoordinator = jc
	am.DefaultBus = mb
	am.PulseStorage = pulseStorageMock

	// Act
	res, err := am.GetPendingRequest(inslogger.TestContext(s.T()), objectID)

	// Assert
	require.NoError(s.T(), err)
	require.Equal(s.T(), parcel, res)

}
