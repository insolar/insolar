//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package artifactmanager

import (
	"bytes"
	"context"
	"crypto/rand"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar/record"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/delegationtoken"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/blob"
	"github.com/insolar/insolar/ledger/storage/drop"
	"github.com/insolar/insolar/ledger/storage/node"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/testutils"
)

type handlerSuite struct {
	suite.Suite

	cm      *component.Manager
	ctx     context.Context
	cleaner func()
	db      storage.DBContext

	scheme      insolar.PlatformCryptographyScheme
	nodeStorage node.Accessor
	jetStorage  jet.Storage

	dropModifier drop.Modifier
	dropAccessor drop.Accessor

	blobModifier blob.Modifier
	blobAccessor blob.Accessor

	recordModifier object.RecordModifier
	recordAccessor object.RecordAccessor

	indexAccessor object.IndexAccessor
	indexModifier object.IndexModifier
}

var (
	domainID = *genRandomID(0)
)

func genRandomID(pulse insolar.PulseNumber) *insolar.ID {
	buff := [insolar.RecordIDSize - insolar.PulseNumberSize]byte{}
	_, err := rand.Read(buff[:])
	if err != nil {
		panic(err)
	}
	return insolar.NewID(pulse, buff[:])
}

func genRefWithID(id *insolar.ID) *insolar.Reference {
	return insolar.NewReference(domainID, *id)
}

func genRandomRef(pulse insolar.PulseNumber) *insolar.Reference {
	return genRefWithID(genRandomID(pulse))
}

func NewHandlerSuite() *handlerSuite {
	return &handlerSuite{
		Suite: suite.Suite{},
	}
}

// Init and run suite
func TestHandlerSuite(t *testing.T) {
	suite.Run(t, NewHandlerSuite())
}

func (s *handlerSuite) BeforeTest(suiteName, testName string) {
	s.cm = &component.Manager{}
	s.ctx = inslogger.TestContext(s.T())

	tmpDB, _, cleaner := storagetest.TmpDB(s.ctx, s.T())
	s.cleaner = cleaner
	s.db = tmpDB
	s.scheme = testutils.NewPlatformCryptographyScheme()
	s.jetStorage = jet.NewStore()
	s.nodeStorage = node.NewStorage()

	storageDB := store.NewMemoryMockDB()
	dropStorage := drop.NewDB(storageDB)
	s.dropAccessor = dropStorage
	s.dropModifier = dropStorage

	blobStorage := blob.NewStorageMemory()
	s.blobAccessor = blobStorage
	s.blobModifier = blobStorage

	recordStorage := object.NewRecordMemory()
	s.recordModifier = recordStorage
	s.recordAccessor = recordStorage

	idxStor := object.NewIndexMemory()
	s.indexAccessor = idxStor
	s.indexModifier = idxStor

	s.cm.Inject(
		s.scheme,
		s.db,
		idxStor,
		store.NewMemoryMockDB(),
		s.jetStorage,
		s.nodeStorage,
		s.dropAccessor,
		s.dropModifier,
		s.recordAccessor,
		s.recordModifier,
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

func (s *handlerSuite) AfterTest(suiteName, testName string) {
	err := s.cm.Stop(s.ctx)
	if err != nil {
		s.T().Error("ComponentManager stop failed", err)
	}
	s.cleaner()
}

func (s *handlerSuite) TestMessageHandler_HandleGetChildren_Redirects() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	jetID := insolar.ID(*insolar.NewJetID(0, nil))

	tf := testutils.NewDelegationTokenFactoryMock(mc)
	tf.IssueGetChildrenRedirectMock.Return(&delegationtoken.GetChildrenRedirectToken{Signature: []byte{1, 2, 3}}, nil)
	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()
	jc := testutils.NewJetCoordinatorMock(mc)

	indexMock := recentstorage.NewRecentIndexStorageMock(s.T())
	pendingMock := recentstorage.NewPendingStorageMock(s.T())

	indexMock.AddObjectMock.Return()
	pendingMock.GetRequestsForObjectMock.Return(nil)
	pendingMock.AddPendingRequestMock.Return()
	pendingMock.RemovePendingRequestMock.Return()

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetIndexStorageMock.Return(indexMock)
	provideMock.GetPendingStorageMock.Return(pendingMock)

	msg := message.GetChildren{
		Parent: *genRandomRef(0),
	}
	h := NewMessageHandler(&configuration.Ledger{
		LightChainLimit: 2,
	})
	h.JetCoordinator = jc
	h.DelegationTokenFactory = tf
	h.Bus = mb
	h.JetStorage = s.jetStorage
	h.Nodes = s.nodeStorage
	h.DBContext = s.db
	h.IndexAccessor = s.indexAccessor
	h.IndexModifier = s.indexModifier
	h.RecordAccessor = s.recordAccessor

	locker := storage.NewIDLockerMock(s.T())
	locker.LockMock.Return()
	locker.UnlockMock.Return()
	h.IDLocker = locker

	err := h.Init(s.ctx)
	require.NoError(s.T(), err)

	h.RecentStorageProvider = provideMock

	s.T().Run("redirects to heavy when no index", func(t *testing.T) {
		objIndex := object.Lifeline{
			LatestState:  genRandomID(insolar.FirstPulseNumber),
			ChildPointer: genRandomID(insolar.FirstPulseNumber),
		}
		mb.SendFunc = func(c context.Context, gm insolar.Message, o *insolar.MessageSendOptions) (r insolar.Reply, r1 error) {
			if m, ok := gm.(*message.GetObjectIndex); ok {
				assert.Equal(t, msg.Parent, m.Object)
				buf := object.EncodeIndex(objIndex)
				require.NoError(t, err)
				return &reply.ObjectIndex{Index: buf}, nil
			}

			panic("unexpected call")
		}
		heavyRef := genRandomRef(0)

		jc.HeavyMock.Return(heavyRef, nil)
		jc.IsBeyondLimitMock.Return(true, nil)
		rep, err := h.handleGetChildren(contextWithJet(s.ctx, jetID), &message.Parcel{
			Msg:         &msg,
			PulseNumber: insolar.FirstPulseNumber + 1,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetChildrenRedirectReply)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetChildrenRedirectToken)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, heavyRef, redirect.GetReceiver())

		idx, err := s.indexAccessor.ForID(s.ctx, *msg.Parent.Record())
		require.NoError(t, err)
		assert.Equal(t, objIndex.LatestState, idx.LatestState)
	})

	s.T().Run("redirect to light when has index and child later than limit", func(t *testing.T) {
		lightRef := genRandomRef(0)
		jc.IsBeyondLimitMock.Return(false, nil)
		jc.NodeForJetMock.Return(lightRef, nil)
		err = s.indexModifier.Set(s.ctx, *msg.Parent.Record(), object.Lifeline{
			ChildPointer: genRandomID(insolar.FirstPulseNumber),
			JetID:        insolar.JetID(jetID),
		})
		require.NoError(t, err)
		rep, err := h.handleGetChildren(contextWithJet(s.ctx, jetID), &message.Parcel{
			Msg:         &msg,
			PulseNumber: insolar.FirstPulseNumber + 1,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetChildrenRedirectReply)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetChildrenRedirectToken)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, lightRef, redirect.GetReceiver())
	})

	s.T().Run("redirect to heavy when has index and child earlier than limit", func(t *testing.T) {
		heavyRef := genRandomRef(0)
		jc.IsBeyondLimitMock.Return(false, nil)
		jc.NodeForJetMock.Return(heavyRef, nil)
		err = s.indexModifier.Set(s.ctx, *msg.Parent.Record(), object.Lifeline{
			ChildPointer: genRandomID(insolar.FirstPulseNumber),
			JetID:        insolar.JetID(jetID),
		})
		require.NoError(t, err)
		rep, err := h.handleGetChildren(contextWithJet(s.ctx, jetID), &message.Parcel{
			Msg:         &msg,
			PulseNumber: insolar.FirstPulseNumber + 2,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetChildrenRedirectReply)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetChildrenRedirectToken)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, heavyRef, redirect.GetReceiver())
	})
}

func (s *handlerSuite) TestMessageHandler_HandleGetDelegate_FetchesIndexFromHeavy() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	jetID := insolar.ID(*insolar.NewJetID(0, nil))

	indexMock := recentstorage.NewRecentIndexStorageMock(s.T())
	pendingMock := recentstorage.NewPendingStorageMock(s.T())

	indexMock.AddObjectMock.Return()
	pendingMock.GetRequestsForObjectMock.Return(nil)
	pendingMock.AddPendingRequestMock.Return()
	pendingMock.RemovePendingRequestMock.Return()

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetIndexStorageMock.Return(indexMock)
	provideMock.GetPendingStorageMock.Return(pendingMock)

	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()
	jc := testutils.NewJetCoordinatorMock(mc)

	h := NewMessageHandler(&configuration.Ledger{
		LightChainLimit: 3,
	})
	h.JetStorage = s.jetStorage
	h.Nodes = s.nodeStorage
	h.DBContext = s.db
	h.IndexModifier = s.indexModifier
	h.IndexAccessor = s.indexAccessor

	h.RecentStorageProvider = provideMock
	idLock := storage.NewIDLockerMock(s.T())
	idLock.LockMock.Return()
	idLock.UnlockMock.Return()
	h.IDLocker = idLock

	delegateType := *genRandomRef(0)
	delegate := *genRandomRef(0)
	objIndex := object.Lifeline{Delegates: map[insolar.Reference]insolar.Reference{delegateType: delegate}}
	msg := message.GetDelegate{
		Head:   *genRandomRef(0),
		AsType: delegateType,
	}

	mb.SendFunc = func(c context.Context, gm insolar.Message, o *insolar.MessageSendOptions) (r insolar.Reply, r1 error) {
		if m, ok := gm.(*message.GetObjectIndex); ok {
			assert.Equal(s.T(), msg.Head, m.Object)
			buf := object.EncodeIndex(objIndex)
			return &reply.ObjectIndex{Index: buf}, nil
		}

		panic("unexpected call")
	}

	h.JetCoordinator = jc
	h.Bus = mb
	err := h.Init(s.ctx)
	require.NoError(s.T(), err)

	heavyRef := genRandomRef(0)
	jc.HeavyMock.Return(heavyRef, nil)
	rep, err := h.handleGetDelegate(contextWithJet(s.ctx, jetID), &message.Parcel{
		Msg: &msg,
	})
	require.NoError(s.T(), err)
	delegateRep, ok := rep.(*reply.Delegate)
	require.True(s.T(), ok)
	assert.Equal(s.T(), delegate, delegateRep.Head)

	idx, err := s.indexAccessor.ForID(s.ctx, *msg.Head.Record())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), objIndex.Delegates, idx.Delegates)
}

func (s *handlerSuite) TestMessageHandler_HandleUpdateObject_FetchesIndexFromHeavy() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	jetID := insolar.ID(*insolar.NewJetID(0, nil))

	indexMock := recentstorage.NewRecentIndexStorageMock(s.T())
	pendingMock := recentstorage.NewPendingStorageMock(s.T())

	indexMock.AddObjectMock.Return()
	pendingMock.GetRequestsForObjectMock.Return(nil)
	pendingMock.AddPendingRequestMock.Return()
	pendingMock.RemovePendingRequestMock.Return()

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetIndexStorageMock.Return(indexMock)
	provideMock.GetPendingStorageMock.Return(pendingMock)

	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()
	jc := testutils.NewJetCoordinatorMock(mc)

	h := NewMessageHandler(&configuration.Ledger{
		LightChainLimit: 3,
	})
	h.JetStorage = s.jetStorage
	h.Nodes = s.nodeStorage
	h.DBContext = s.db
	h.IndexAccessor = s.indexAccessor
	h.IndexModifier = s.indexModifier
	h.PlatformCryptographyScheme = s.scheme
	h.RecentStorageProvider = provideMock
	h.RecordModifier = s.recordModifier

	blobStorage := blob.NewStorageMemory()
	h.BlobModifier = blobStorage
	h.BlobAccessor = blobStorage

	idLockMock := storage.NewIDLockerMock(s.T())
	idLockMock.LockMock.Return()
	idLockMock.UnlockMock.Return()
	h.IDLocker = idLockMock

	objIndex := object.Lifeline{LatestState: genRandomID(0), State: object.StateActivation}
	amendRecord := object.AmendRecord{
		PrevState: *objIndex.LatestState,
	}
	amendHash := s.scheme.ReferenceHasher()
	_, err := amendRecord.WriteHashData(amendHash)
	require.NoError(s.T(), err)

	msg := message.UpdateObject{
		Record: object.EncodeVirtual(&amendRecord),
		Object: *genRandomRef(0),
	}

	mb.SendFunc = func(c context.Context, gm insolar.Message, o *insolar.MessageSendOptions) (r insolar.Reply, r1 error) {
		if m, ok := gm.(*message.GetObjectIndex); ok {
			assert.Equal(s.T(), msg.Object, m.Object)
			buf := object.EncodeIndex(objIndex)
			require.NoError(s.T(), err)
			return &reply.ObjectIndex{Index: buf}, nil
		}

		panic("unexpected call")
	}

	h.JetCoordinator = jc
	h.Bus = mb
	err = h.Init(s.ctx)
	require.NoError(s.T(), err)
	heavyRef := genRandomRef(0)
	jc.HeavyMock.Return(heavyRef, nil)
	rep, err := h.handleUpdateObject(contextWithJet(s.ctx, jetID), &message.Parcel{
		Msg:         &msg,
		PulseNumber: insolar.FirstPulseNumber,
	})
	require.NoError(s.T(), err)
	objRep, ok := rep.(*reply.Object)
	require.True(s.T(), ok)

	idx, err := s.indexAccessor.ForID(s.ctx, *msg.Object.Record())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), objRep.State, *idx.LatestState)
}

func (s *handlerSuite) TestMessageHandler_HandleUpdateObject_UpdateIndexState() {
	// Arrange
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	jetID := insolar.ID(*insolar.NewJetID(0, nil))

	indexMock := recentstorage.NewRecentIndexStorageMock(s.T())
	pendingMock := recentstorage.NewPendingStorageMock(s.T())

	indexMock.AddObjectMock.Return()
	pendingMock.GetRequestsForObjectMock.Return(nil)
	pendingMock.AddPendingRequestMock.Return()
	pendingMock.RemovePendingRequestMock.Return()

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetIndexStorageMock.Return(indexMock)
	provideMock.GetPendingStorageMock.Return(pendingMock)

	h := NewMessageHandler(&configuration.Ledger{
		LightChainLimit: 3,
	})
	h.JetStorage = s.jetStorage
	h.Nodes = s.nodeStorage
	h.DBContext = s.db
	h.IndexModifier = s.indexModifier
	h.IndexAccessor = s.indexAccessor
	h.RecentStorageProvider = provideMock
	h.PlatformCryptographyScheme = s.scheme
	h.RecordModifier = s.recordModifier

	blobStorage := blob.NewStorageMemory()
	h.BlobModifier = blobStorage
	h.BlobAccessor = blobStorage

	idLockMock := storage.NewIDLockerMock(s.T())
	idLockMock.LockMock.Return()
	idLockMock.UnlockMock.Return()
	h.IDLocker = idLockMock

	objIndex := object.Lifeline{
		LatestState:  genRandomID(0),
		State:        object.StateActivation,
		LatestUpdate: 0,
		JetID:        insolar.JetID(jetID),
	}
	amendRecord := object.AmendRecord{
		PrevState: *objIndex.LatestState,
	}
	amendHash := s.scheme.ReferenceHasher()
	_, err := amendRecord.WriteHashData(amendHash)
	require.NoError(s.T(), err)

	msg := message.UpdateObject{
		Record: object.EncodeVirtual(&amendRecord),
		Object: *genRandomRef(0),
	}
	err = s.indexModifier.Set(s.ctx, *msg.Object.Record(), objIndex)
	require.NoError(s.T(), err)

	// Act
	rep, err := h.handleUpdateObject(contextWithJet(s.ctx, jetID), &message.Parcel{
		Msg:         &msg,
		PulseNumber: insolar.FirstPulseNumber,
	})
	require.NoError(s.T(), err)
	_, ok := rep.(*reply.Object)
	require.True(s.T(), ok)

	// Arrange
	idx, err := s.indexAccessor.ForID(s.ctx, *msg.Object.Record())
	require.NoError(s.T(), err)
	require.Equal(s.T(), insolar.FirstPulseNumber, int(idx.LatestUpdate))
}

func (s *handlerSuite) TestMessageHandler_HandleGetObjectIndex() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	jetID := insolar.ID(*insolar.NewJetID(0, nil))
	msg := message.GetObjectIndex{
		Object: *genRandomRef(0),
	}
	indexMock := recentstorage.NewRecentIndexStorageMock(s.T())
	pendingMock := recentstorage.NewPendingStorageMock(s.T())

	indexMock.AddObjectMock.Return()
	pendingMock.GetRequestsForObjectMock.Return(nil)
	pendingMock.AddPendingRequestMock.Return()
	pendingMock.RemovePendingRequestMock.Return()

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetIndexStorageMock.Return(indexMock)
	provideMock.GetPendingStorageMock.Return(pendingMock)

	jc := testutils.NewJetCoordinatorMock(mc)

	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()

	h := NewMessageHandler(&configuration.Ledger{
		LightChainLimit: 3,
	})
	h.JetCoordinator = jc
	h.Bus = mb
	h.JetStorage = s.jetStorage
	h.Nodes = s.nodeStorage
	h.DBContext = s.db
	h.IndexAccessor = s.indexAccessor
	h.IndexModifier = s.indexModifier

	idLock := storage.NewIDLockerMock(s.T())
	idLock.LockMock.Return()
	idLock.UnlockMock.Return()
	h.IDLocker = idLock

	err := h.Init(s.ctx)
	require.NoError(s.T(), err)

	h.RecentStorageProvider = provideMock

	objectIndex := object.Lifeline{LatestState: genRandomID(0), JetID: insolar.JetID(jetID), Delegates: map[insolar.Reference]insolar.Reference{}}
	err = s.indexModifier.Set(s.ctx, *msg.Object.Record(), objectIndex)
	require.NoError(s.T(), err)

	rep, err := h.handleGetObjectIndex(contextWithJet(s.ctx, jetID), &message.Parcel{
		Msg: &msg,
	})
	require.NoError(s.T(), err)
	indexRep, ok := rep.(*reply.ObjectIndex)
	require.True(s.T(), ok)
	decodedIndex := object.MustDecodeIndex(indexRep.Index)
	assert.Equal(s.T(), objectIndex, decodedIndex)
}

func (s *handlerSuite) TestMessageHandler_HandleHasPendingRequests() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	msg := message.GetPendingRequests{
		Object: *genRandomRef(0),
	}
	pendingRequests := []insolar.ID{
		*genRandomID(insolar.FirstPulseNumber),
		*genRandomID(insolar.FirstPulseNumber),
	}

	recentStorageMock := recentstorage.NewPendingStorageMock(s.T())
	recentStorageMock.GetRequestsForObjectMock.Return(pendingRequests)

	jetID := insolar.ID(*insolar.NewJetID(0, nil))
	jc := testutils.NewJetCoordinatorMock(mc)
	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()

	h := NewMessageHandler(&configuration.Ledger{})
	h.JetCoordinator = jc
	h.Bus = mb
	h.JetStorage = s.jetStorage
	h.Nodes = s.nodeStorage
	h.DBContext = s.db
	h.IndexModifier = s.indexModifier
	h.IndexAccessor = s.indexAccessor

	err := h.Init(s.ctx)
	require.NoError(s.T(), err)

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetPendingStorageMock.Return(recentStorageMock)

	h.RecentStorageProvider = provideMock

	rep, err := h.handleHasPendingRequests(contextWithJet(s.ctx, jetID), &message.Parcel{
		Msg:         &msg,
		PulseNumber: insolar.FirstPulseNumber + 1,
	})
	require.NoError(s.T(), err)
	has, ok := rep.(*reply.HasPendingRequests)
	require.True(s.T(), ok)
	assert.True(s.T(), has.Has)
}

func (s *handlerSuite) TestMessageHandler_HandleRegisterChild_FetchesIndexFromHeavy() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	jetID := insolar.ID(*insolar.NewJetID(0, nil))

	indexMock := recentstorage.NewRecentIndexStorageMock(s.T())
	pendingMock := recentstorage.NewPendingStorageMock(s.T())

	indexMock.AddObjectMock.Return()
	pendingMock.GetRequestsForObjectMock.Return(nil)
	pendingMock.AddPendingRequestMock.Return()
	pendingMock.RemovePendingRequestMock.Return()

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetIndexStorageMock.Return(indexMock)
	provideMock.GetPendingStorageMock.Return(pendingMock)

	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()
	jc := testutils.NewJetCoordinatorMock(mc)
	h := NewMessageHandler(&configuration.Ledger{
		LightChainLimit: 2,
	})
	h.JetStorage = s.jetStorage
	h.Nodes = s.nodeStorage
	h.DBContext = s.db
	h.IndexModifier = s.indexModifier
	h.IndexAccessor = s.indexAccessor
	h.RecentStorageProvider = provideMock
	h.PlatformCryptographyScheme = s.scheme
	h.RecordModifier = s.recordModifier

	idLockMock := storage.NewIDLockerMock(s.T())
	idLockMock.LockMock.Return()
	idLockMock.UnlockMock.Return()
	h.IDLocker = idLockMock

	objIndex := object.Lifeline{LatestState: genRandomID(0), State: object.StateActivation}
	childRecord := object.ChildRecord{
		Ref:       *genRandomRef(0),
		PrevChild: nil,
	}
	amendHash := s.scheme.ReferenceHasher()
	_, err := childRecord.WriteHashData(amendHash)
	require.NoError(s.T(), err)
	childID := insolar.NewID(0, amendHash.Sum(nil))

	msg := message.RegisterChild{
		Record: object.EncodeVirtual(&childRecord),
		Parent: *genRandomRef(0),
	}

	mb.SendFunc = func(c context.Context, gm insolar.Message, o *insolar.MessageSendOptions) (r insolar.Reply, r1 error) {
		if m, ok := gm.(*message.GetObjectIndex); ok {
			assert.Equal(s.T(), msg.Parent, m.Object)
			buf := object.EncodeIndex(objIndex)
			require.NoError(s.T(), err)
			return &reply.ObjectIndex{Index: buf}, nil
		}

		panic("unexpected call")
	}

	h.JetCoordinator = jc
	h.Bus = mb
	err = h.Init(s.ctx)
	require.NoError(s.T(), err)
	heavyRef := genRandomRef(0)
	jc.HeavyMock.Return(heavyRef, nil)
	rep, err := h.handleRegisterChild(contextWithJet(s.ctx, jetID), &message.Parcel{
		Msg: &msg,
	})
	require.NoError(s.T(), err)
	objRep, ok := rep.(*reply.ID)
	require.True(s.T(), ok)
	assert.Equal(s.T(), *childID, objRep.ID)

	idx, err := s.indexAccessor.ForID(s.ctx, *msg.Parent.Record())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), childID, idx.ChildPointer)
}

func (s *handlerSuite) TestMessageHandler_HandleRegisterChild_IndexStateUpdated() {
	// Arrange
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	jetID := insolar.ID(*insolar.NewJetID(0, nil))

	indexMock := recentstorage.NewRecentIndexStorageMock(s.T())
	pendingMock := recentstorage.NewPendingStorageMock(s.T())

	indexMock.AddObjectMock.Return()
	pendingMock.GetRequestsForObjectMock.Return(nil)
	pendingMock.AddPendingRequestMock.Return()
	pendingMock.RemovePendingRequestMock.Return()

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetIndexStorageMock.Return(indexMock)
	provideMock.GetPendingStorageMock.Return(pendingMock)

	h := NewMessageHandler(&configuration.Ledger{
		LightChainLimit: 2,
	})
	h.JetStorage = s.jetStorage
	h.Nodes = s.nodeStorage
	h.DBContext = s.db
	h.IndexModifier = s.indexModifier
	h.IndexAccessor = s.indexAccessor
	h.RecentStorageProvider = provideMock
	h.PlatformCryptographyScheme = s.scheme
	h.RecordModifier = s.recordModifier

	idLockMock := storage.NewIDLockerMock(s.T())
	idLockMock.LockMock.Return()
	idLockMock.UnlockMock.Return()
	h.IDLocker = idLockMock

	objIndex := object.Lifeline{
		LatestState:  genRandomID(0),
		State:        object.StateActivation,
		LatestUpdate: insolar.FirstPulseNumber,
		JetID:        insolar.JetID(jetID),
	}
	childRecord := object.ChildRecord{
		Ref:       *genRandomRef(0),
		PrevChild: nil,
	}
	msg := message.RegisterChild{
		Record: object.EncodeVirtual(&childRecord),
		Parent: *genRandomRef(0),
	}

	err := s.indexModifier.Set(s.ctx, *msg.Parent.Record(), objIndex)
	require.NoError(s.T(), err)

	// Act
	_, err = h.handleRegisterChild(contextWithJet(s.ctx, jetID), &message.Parcel{
		Msg:         &msg,
		PulseNumber: insolar.FirstPulseNumber + 100,
	})
	require.NoError(s.T(), err)

	// Assert
	idx, err := s.indexAccessor.ForID(s.ctx, *msg.Parent.Record())
	require.NoError(s.T(), err)
	require.Equal(s.T(), int(idx.LatestUpdate), insolar.FirstPulseNumber+100)
}

func (s *handlerSuite) TestMessageHandler_HandleHotRecords() {
	mc := minimock.NewController(s.T())
	jetID := gen.JetID()

	jc := testutils.NewJetCoordinatorMock(mc)

	firstID := insolar.NewID(insolar.FirstPulseNumber, []byte{1, 2, 3})
	secondID := object.NewRecordIDFromRecord(s.scheme, insolar.FirstPulseNumber, &object.CodeRecord{})
	thirdID := object.NewRecordIDFromRecord(s.scheme, insolar.FirstPulseNumber-1, &object.CodeRecord{})

	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()
	mb.SendFunc = func(p context.Context, p1 insolar.Message, p2 *insolar.MessageSendOptions) (r insolar.Reply, r1 error) {
		parsedMsg, ok := p1.(*message.AbandonedRequestsNotification)
		require.Equal(s.T(), true, ok)
		require.Equal(s.T(), *secondID, parsedMsg.Object)
		return &reply.OK{}, nil
	}

	firstIndex := object.EncodeIndex(object.Lifeline{
		LatestState: firstID,
	})
	err := s.indexModifier.Set(s.ctx, *firstID, object.Lifeline{
		LatestState: firstID,
		JetID:       insolar.JetID(jetID),
	})

	hotIndexes := &message.HotData{
		Jet:         *insolar.NewReference(insolar.DomainID, insolar.ID(jetID)),
		PulseNumber: insolar.FirstPulseNumber,
		RecentObjects: map[insolar.ID]message.HotIndex{
			*firstID: {
				Index: firstIndex,
				TTL:   320,
			},
		},
		PendingRequests: map[insolar.ID]recentstorage.PendingObjectContext{
			*secondID: {},
			*thirdID:  {Active: true},
		},
		Drop: drop.Drop{Pulse: insolar.FirstPulseNumber, Hash: []byte{88}, JetID: jetID},
	}

	indexMock := recentstorage.NewRecentIndexStorageMock(s.T())
	pendingMock := recentstorage.NewPendingStorageMock(s.T())

	pendingMock.SetContextToObjectFunc = func(p context.Context, p1 insolar.ID, p2 recentstorage.PendingObjectContext) {

		if bytes.Equal(p1.Bytes(), secondID.Bytes()) {
			require.Equal(s.T(), false, p2.Active)
			return
		}
		if bytes.Equal(p1.Bytes(), thirdID.Bytes()) {
			require.Equal(s.T(), false, p2.Active)
			return
		}
		s.T().Fail()
	}
	indexMock.AddObjectWithTLLFunc = func(ctx context.Context, p insolar.ID, ttl int) {
		require.Equal(s.T(), p, *firstID)
		require.Equal(s.T(), 320, ttl)
	}
	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetPendingStorageMock.Return(pendingMock)
	provideMock.GetIndexStorageMock.Return(indexMock)

	h := NewMessageHandler(&configuration.Ledger{})
	h.JetCoordinator = jc
	h.RecentStorageProvider = provideMock
	h.Bus = mb
	h.JetStorage = s.jetStorage
	h.Nodes = s.nodeStorage
	h.DBContext = s.db
	h.IndexModifier = s.indexModifier
	h.IndexAccessor = s.indexAccessor
	h.DropModifier = s.dropModifier

	err = h.Init(s.ctx)
	require.NoError(s.T(), err)

	res, err := h.handleHotRecords(s.ctx, &message.Parcel{Msg: hotIndexes})

	require.NoError(s.T(), err)
	require.Equal(s.T(), res, &reply.OK{})

	savedDrop, err := s.dropAccessor.ForPulse(s.ctx, jetID, insolar.FirstPulseNumber)
	require.NoError(s.T(), err)
	require.Equal(s.T(), drop.Drop{Pulse: insolar.FirstPulseNumber, Hash: []byte{88}, JetID: jetID}, savedDrop)

	indexMock.MinimockFinish()
	pendingMock.MinimockFinish()
}

func (s *handlerSuite) TestMessageHandler_HandleGetRequest() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()

	jetID := insolar.ID(*insolar.NewJetID(0, nil))

	req := object.RequestRecord{
		MessageHash: []byte{1, 2, 3},
		Object:      *genRandomID(0),
	}

	reqID := object.NewRecordIDFromRecord(s.scheme, insolar.FirstPulseNumber, &req)
	rec := record.MaterialRecord{
		Record: &req,
		JetID:  insolar.JetID(jetID),
	}
	err := s.recordModifier.Set(s.ctx, *reqID, rec)
	require.NoError(s.T(), err)

	msg := message.GetRequest{
		Request: *reqID,
	}

	h := NewMessageHandler(&configuration.Ledger{})
	h.RecordAccessor = s.recordAccessor

	rep, err := h.handleGetRequest(contextWithJet(s.ctx, jetID), &message.Parcel{
		Msg:         &msg,
		PulseNumber: insolar.FirstPulseNumber + 1,
	})
	require.NoError(s.T(), err)
	reqReply, ok := rep.(*reply.Request)
	require.True(s.T(), ok)
	vrec, _ := object.DecodeVirtual(reqReply.Record)
	assert.Equal(s.T(), req, *vrec.(*object.RequestRecord))
}
