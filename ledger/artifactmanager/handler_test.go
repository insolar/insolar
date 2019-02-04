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
	"testing"

	"github.com/gojuno/minimock"
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
	"github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type handlerSuite struct {
	suite.Suite

	cm      *component.Manager
	ctx     context.Context
	cleaner func()
	db      storage.DBContext

	scheme        core.PlatformCryptographyScheme
	pulseTracker  storage.PulseTracker
	nodeStorage   storage.NodeStorage
	objectStorage storage.ObjectStorage
	jetStorage    storage.JetStorage
	dropStorage   storage.DropStorage
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

	db, cleaner := storagetest.TmpDB(s.ctx, s.T())
	s.cleaner = cleaner
	s.db = db
	s.scheme = platformpolicy.NewPlatformCryptographyScheme()
	s.jetStorage = storage.NewJetStorage()
	s.nodeStorage = storage.NewNodeStorage()
	s.pulseTracker = storage.NewPulseTracker()
	s.objectStorage = storage.NewObjectStorage()
	s.dropStorage = storage.NewDropStorage(10)

	s.cm.Inject(
		s.scheme,
		s.db,
		s.jetStorage,
		s.nodeStorage,
		s.pulseTracker,
		s.objectStorage,
		s.dropStorage,
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

func (s *handlerSuite) TestMessageHandler_HandleGetObject_Redirects() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	jetID := *jet.NewID(0, nil)

	tf := testutils.NewDelegationTokenFactoryMock(mc)
	jc := testutils.NewJetCoordinatorMock(mc)
	msg := message.GetObject{
		Head: *genRandomRef(core.FirstPulseNumber),
	}

	certificate := testutils.NewCertificateMock(s.T())
	certificate.GetRoleMock.Return(core.StaticRoleLightMaterial)

	tf.IssueGetObjectRedirectMock.Return(&delegationtoken.GetObjectRedirectToken{Signature: []byte{1, 2, 3}}, nil)

	h := NewMessageHandler(&configuration.Ledger{
		LightChainLimit: 2,
	}, certificate)
	h.JetStorage = s.jetStorage
	h.NodeStorage = s.nodeStorage
	h.DBContext = s.db
	h.PulseTracker = s.pulseTracker
	h.ObjectStorage = s.objectStorage

	recentStorageMock := recentstorage.NewRecentStorageMock(s.T())
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetStorageFunc = func(ctx context.Context, p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	h.RecentStorageProvider = provideMock
	h.JetCoordinator = jc
	h.DelegationTokenFactory = tf
	h.Bus = mb

	err := h.Init(s.ctx)
	require.NoError(s.T(), err)

	s.T().Run("fetches_index_from_heavy_when_no_index", func(t *testing.T) {
		idxState := genRandomID(core.FirstPulseNumber)
		objIndex := index.ObjectLifeline{
			LatestState: idxState,
		}
		lightRef := genRandomRef(0)
		heavyRef := genRandomRef(1)

		mb.SendFunc = func(c context.Context, gm core.Message, o *core.MessageSendOptions) (r core.Reply, r1 error) {
			if m, ok := gm.(*message.GetObjectIndex); ok {
				assert.Equal(t, msg.Head, m.Object)
				buf, err := index.EncodeObjectLifeline(&objIndex)
				require.NoError(t, err)
				return &reply.ObjectIndex{Index: buf}, nil
			}

			panic("unexpected call")
		}

		jc.LightExecutorForJetMock.Return(lightRef, nil)
		jc.HeavyMock.Return(heavyRef, nil)

		rep, err := h.handleGetObject(contextWithJet(s.ctx, jetID), &message.Parcel{
			Msg:         &msg,
			PulseNumber: core.FirstPulseNumber,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetObjectRedirectReply)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetObjectRedirectToken)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, lightRef, redirect.GetReceiver())
		assert.Equal(t, idxState, redirect.StateID)

		idx, err := s.objectStorage.GetObjectIndex(s.ctx, jetID, msg.Head.Record(), false)
		require.NoError(t, err)
		assert.Equal(t, objIndex.LatestState, idx.LatestState)
	})

	err = s.pulseTracker.AddPulse(s.ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 1})
	require.NoError(s.T(), err)
	s.T().Run("redirect to light when has index and state later than limit", func(t *testing.T) {
		lightRef := genRandomRef(0)
		jc.LightExecutorForJetMock.Set(nil)
		jc.LightExecutorForJetFunc = func(c context.Context, j core.RecordID, p core.PulseNumber) (*core.RecordRef, error) {
			switch p {
			case core.FirstPulseNumber:
				return lightRef, nil
			case core.FirstPulseNumber + 1:
				return &core.RecordRef{}, nil
			}
			panic("unexpected call")
		}
		stateID := genRandomID(core.FirstPulseNumber)
		err = s.objectStorage.SetObjectIndex(s.ctx, jetID, msg.Head.Record(), &index.ObjectLifeline{
			LatestState: stateID,
		})
		require.NoError(t, err)
		rep, err := h.handleGetObject(contextWithJet(s.ctx, jetID), &message.Parcel{
			Msg:         &msg,
			PulseNumber: core.FirstPulseNumber + 1,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetObjectRedirectReply)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetObjectRedirectToken)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, lightRef, redirect.GetReceiver())
		assert.Equal(t, stateID, redirect.StateID)
	})

	err = s.pulseTracker.AddPulse(s.ctx, core.Pulse{
		PulseNumber: core.FirstPulseNumber + 2,
	})
	require.NoError(s.T(), err)
	s.T().Run("redirect to heavy when has index and state earlier than limit", func(t *testing.T) {
		heavyRef := genRandomRef(0)
		jc.LightExecutorForJetMock.Return(&core.RecordRef{}, nil)
		jc.HeavyMock.Return(heavyRef, nil)
		stateID := genRandomID(core.FirstPulseNumber)

		err = s.objectStorage.SetObjectIndex(s.ctx, jetID, msg.Head.Record(), &index.ObjectLifeline{
			LatestState: stateID,
		})
		require.NoError(t, err)
		rep, err := h.handleGetObject(contextWithJet(s.ctx, jetID), &message.Parcel{
			Msg:         &msg,
			PulseNumber: core.FirstPulseNumber + 2,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetObjectRedirectReply)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetObjectRedirectToken)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, heavyRef, redirect.GetReceiver())
		assert.Equal(t, stateID, redirect.StateID)
	})
}

func (s *handlerSuite) TestMessageHandler_HandleGetChildren_Redirects() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	jetID := *jet.NewID(0, nil)

	tf := testutils.NewDelegationTokenFactoryMock(mc)
	tf.IssueGetChildrenRedirectMock.Return(&delegationtoken.GetChildrenRedirectToken{Signature: []byte{1, 2, 3}}, nil)
	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()
	jc := testutils.NewJetCoordinatorMock(mc)

	recentStorageMock := recentstorage.NewRecentStorageMock(s.T())
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	certificate := testutils.NewCertificateMock(s.T())
	certificate.GetRoleMock.Return(core.StaticRoleLightMaterial)

	msg := message.GetChildren{
		Parent: *genRandomRef(0),
	}
	h := NewMessageHandler(&configuration.Ledger{
		LightChainLimit: 2,
	}, certificate)
	h.JetCoordinator = jc
	h.DelegationTokenFactory = tf
	h.Bus = mb
	h.JetStorage = s.jetStorage
	h.NodeStorage = s.nodeStorage
	h.DBContext = s.db
	h.PulseTracker = s.pulseTracker
	h.ObjectStorage = s.objectStorage

	err := h.Init(s.ctx)
	require.NoError(s.T(), err)
	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetStorageFunc = func(ctx context.Context, p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	h.RecentStorageProvider = provideMock

	err = s.pulseTracker.AddPulse(s.ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 1})
	require.NoError(s.T(), err)

	s.T().Run("redirects to heavy when no index", func(t *testing.T) {
		objIndex := index.ObjectLifeline{
			LatestState:  genRandomID(core.FirstPulseNumber),
			ChildPointer: genRandomID(core.FirstPulseNumber),
		}
		mb.SendFunc = func(c context.Context, gm core.Message, o *core.MessageSendOptions) (r core.Reply, r1 error) {
			if m, ok := gm.(*message.GetObjectIndex); ok {
				assert.Equal(t, msg.Parent, m.Object)
				buf, err := index.EncodeObjectLifeline(&objIndex)
				require.NoError(t, err)
				return &reply.ObjectIndex{Index: buf}, nil
			}

			panic("unexpected call")
		}
		jc.LightExecutorForJetMock.Return(&core.RecordRef{}, nil)
		heavyRef := genRandomRef(0)
		lightRef := genRandomRef(1)

		jc.HeavyMock.Return(heavyRef, nil)
		jc.LightExecutorForJetMock.Return(lightRef, nil)
		rep, err := h.handleGetChildren(contextWithJet(s.ctx, jetID), &message.Parcel{
			Msg:         &msg,
			PulseNumber: core.FirstPulseNumber + 1,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetChildrenRedirectReply)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetChildrenRedirectToken)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, lightRef, redirect.GetReceiver())

		idx, err := s.objectStorage.GetObjectIndex(s.ctx, jetID, msg.Parent.Record(), false)
		require.NoError(t, err)
		assert.Equal(t, objIndex.LatestState, idx.LatestState)
	})

	s.T().Run("redirect to light when has index and child later than limit", func(t *testing.T) {
		lightRef := genRandomRef(0)
		jc.LightExecutorForJetMock.Set(nil)
		jc.LightExecutorForJetFunc = func(c context.Context, j core.RecordID, p core.PulseNumber) (*core.RecordRef, error) {
			switch p {
			case core.FirstPulseNumber:
				return lightRef, nil
			case core.FirstPulseNumber + 1:
				return &core.RecordRef{}, nil
			}
			panic("unexpected call")
		}
		err = s.objectStorage.SetObjectIndex(s.ctx, jetID, msg.Parent.Record(), &index.ObjectLifeline{
			ChildPointer: genRandomID(core.FirstPulseNumber),
		})
		require.NoError(t, err)
		rep, err := h.handleGetChildren(contextWithJet(s.ctx, jetID), &message.Parcel{
			Msg:         &msg,
			PulseNumber: core.FirstPulseNumber + 1,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetChildrenRedirectReply)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetChildrenRedirectToken)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, lightRef, redirect.GetReceiver())
	})

	s.T().Run("redirect to heavy when has index and child earlier than limit", func(t *testing.T) {
		err = s.pulseTracker.AddPulse(s.ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 2})
		require.NoError(t, err)
		heavyRef := genRandomRef(0)
		jc.LightExecutorForJetMock.Return(&core.RecordRef{}, nil)
		jc.HeavyMock.Return(heavyRef, nil)
		err = s.objectStorage.SetObjectIndex(s.ctx, jetID, msg.Parent.Record(), &index.ObjectLifeline{
			ChildPointer: genRandomID(core.FirstPulseNumber),
		})
		require.NoError(t, err)
		rep, err := h.handleGetChildren(contextWithJet(s.ctx, jetID), &message.Parcel{
			Msg:         &msg,
			PulseNumber: core.FirstPulseNumber + 2,
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
	jetID := *jet.NewID(0, nil)

	recentStorageMock := recentstorage.NewRecentStorageMock(s.T())
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	certificate := testutils.NewCertificateMock(s.T())
	certificate.GetRoleMock.Return(core.StaticRoleLightMaterial)

	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()
	jc := testutils.NewJetCoordinatorMock(mc)

	h := NewMessageHandler(&configuration.Ledger{
		LightChainLimit: 3,
	}, certificate)
	h.JetStorage = s.jetStorage
	h.NodeStorage = s.nodeStorage
	h.DBContext = s.db
	h.PulseTracker = s.pulseTracker
	h.ObjectStorage = s.objectStorage

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetStorageFunc = func(ctx context.Context, p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	h.RecentStorageProvider = provideMock

	delegateType := *genRandomRef(0)
	delegate := *genRandomRef(0)
	objIndex := index.ObjectLifeline{Delegates: map[core.RecordRef]core.RecordRef{delegateType: delegate}}
	msg := message.GetDelegate{
		Head:   *genRandomRef(0),
		AsType: delegateType,
	}

	mb.SendFunc = func(c context.Context, gm core.Message, o *core.MessageSendOptions) (r core.Reply, r1 error) {
		if m, ok := gm.(*message.GetObjectIndex); ok {
			assert.Equal(s.T(), msg.Head, m.Object)
			buf, err := index.EncodeObjectLifeline(&objIndex)
			require.NoError(s.T(), err)
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

	idx, err := s.objectStorage.GetObjectIndex(s.ctx, jetID, msg.Head.Record(), false)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), objIndex.Delegates, idx.Delegates)
}

func (s *handlerSuite) TestMessageHandler_HandleUpdateObject_FetchesIndexFromHeavy() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	jetID := *jet.NewID(0, nil)

	recentStorageMock := recentstorage.NewRecentStorageMock(s.T())
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	certificate := testutils.NewCertificateMock(s.T())
	certificate.GetRoleMock.Return(core.StaticRoleLightMaterial)

	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()
	jc := testutils.NewJetCoordinatorMock(mc)

	h := NewMessageHandler(&configuration.Ledger{
		LightChainLimit: 3,
	}, certificate)
	h.JetStorage = s.jetStorage
	h.NodeStorage = s.nodeStorage
	h.DBContext = s.db
	h.PulseTracker = s.pulseTracker
	h.ObjectStorage = s.objectStorage

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetStorageFunc = func(ctx context.Context, p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	h.RecentStorageProvider = provideMock

	objIndex := index.ObjectLifeline{LatestState: genRandomID(0), State: record.StateActivation}
	amendRecord := record.ObjectAmendRecord{
		PrevState: *objIndex.LatestState,
	}
	amendHash := s.scheme.ReferenceHasher()
	_, err := amendRecord.WriteHashData(amendHash)
	require.NoError(s.T(), err)

	msg := message.UpdateObject{
		Record: record.SerializeRecord(&amendRecord),
		Object: *genRandomRef(0),
	}

	mb.SendFunc = func(c context.Context, gm core.Message, o *core.MessageSendOptions) (r core.Reply, r1 error) {
		if m, ok := gm.(*message.GetObjectIndex); ok {
			assert.Equal(s.T(), msg.Object, m.Object)
			buf, err := index.EncodeObjectLifeline(&objIndex)
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
		PulseNumber: core.FirstPulseNumber,
	})
	require.NoError(s.T(), err)
	objRep, ok := rep.(*reply.Object)
	require.True(s.T(), ok)

	idx, err := s.objectStorage.GetObjectIndex(s.ctx, jetID, msg.Object.Record(), false)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), objRep.State, *idx.LatestState)
}

func (s *handlerSuite) TestMessageHandler_HandleUpdateObject_UpdateIndexState() {
	// Arrange
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	jetID := *jet.NewID(0, nil)

	recentStorageMock := recentstorage.NewRecentStorageMock(s.T())
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetStorageFunc = func(ctx context.Context, p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	certificate := testutils.NewCertificateMock(s.T())
	certificate.GetRoleMock.Return(core.StaticRoleLightMaterial)

	h := NewMessageHandler(&configuration.Ledger{
		LightChainLimit: 3,
	}, certificate)
	h.JetStorage = s.jetStorage
	h.NodeStorage = s.nodeStorage
	h.DBContext = s.db
	h.PulseTracker = s.pulseTracker
	h.ObjectStorage = s.objectStorage
	h.RecentStorageProvider = provideMock

	objIndex := index.ObjectLifeline{
		LatestState:  genRandomID(0),
		State:        record.StateActivation,
		LatestUpdate: 0,
	}
	amendRecord := record.ObjectAmendRecord{
		PrevState: *objIndex.LatestState,
	}
	amendHash := s.db.GetPlatformCryptographyScheme().ReferenceHasher()
	_, err := amendRecord.WriteHashData(amendHash)
	require.NoError(s.T(), err)

	msg := message.UpdateObject{
		Record: record.SerializeRecord(&amendRecord),
		Object: *genRandomRef(0),
	}
	err = s.objectStorage.SetObjectIndex(s.ctx, jetID, msg.Object.Record(), &objIndex)
	require.NoError(s.T(), err)

	// Act
	rep, err := h.handleUpdateObject(contextWithJet(s.ctx, jetID), &message.Parcel{
		Msg:         &msg,
		PulseNumber: core.FirstPulseNumber,
	})
	require.NoError(s.T(), err)
	_, ok := rep.(*reply.Object)
	require.True(s.T(), ok)

	// Arrange
	idx, err := s.objectStorage.GetObjectIndex(s.ctx, jetID, msg.Object.Record(), false)
	require.NoError(s.T(), err)
	require.Equal(s.T(), core.FirstPulseNumber, int(idx.LatestUpdate))
}

func (s *handlerSuite) TestMessageHandler_HandleGetObjectIndex() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	jetID := *jet.NewID(0, nil)
	msg := message.GetObjectIndex{
		Object: *genRandomRef(0),
	}
	recentStorageMock := recentstorage.NewRecentStorageMock(s.T())
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	certificate := testutils.NewCertificateMock(s.T())
	certificate.GetRoleMock.Return(core.StaticRoleLightMaterial)

	jc := testutils.NewJetCoordinatorMock(mc)

	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()

	h := NewMessageHandler(&configuration.Ledger{
		LightChainLimit: 3,
	}, certificate)
	h.JetCoordinator = jc
	h.Bus = mb
	h.JetStorage = s.jetStorage
	h.NodeStorage = s.nodeStorage
	h.DBContext = s.db
	h.PulseTracker = s.pulseTracker
	h.ObjectStorage = s.objectStorage

	err := h.Init(s.ctx)
	require.NoError(s.T(), err)

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetStorageFunc = func(ctx context.Context, p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	h.RecentStorageProvider = provideMock

	objectIndex := index.ObjectLifeline{LatestState: genRandomID(0)}
	err = s.objectStorage.SetObjectIndex(s.ctx, jetID, msg.Object.Record(), &objectIndex)
	require.NoError(s.T(), err)

	rep, err := h.handleGetObjectIndex(contextWithJet(s.ctx, jetID), &message.Parcel{
		Msg: &msg,
	})
	require.NoError(s.T(), err)
	indexRep, ok := rep.(*reply.ObjectIndex)
	require.True(s.T(), ok)
	decodedIndex, err := index.DecodeObjectLifeline(indexRep.Index)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), objectIndex, *decodedIndex)
}

func (s *handlerSuite) TestMessageHandler_HandleHasPendingRequests() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	msg := message.GetPendingRequests{
		Object: *genRandomRef(0),
	}
	pendingRequests := []core.RecordID{
		*genRandomID(core.FirstPulseNumber),
		*genRandomID(core.FirstPulseNumber),
	}

	recentStorageMock := recentstorage.NewRecentStorageMock(s.T())
	recentStorageMock.GetRequestsForObjectMock.Return(pendingRequests)

	certificate := testutils.NewCertificateMock(s.T())
	certificate.GetRoleMock.Return(core.StaticRoleLightMaterial)

	jetID := *jet.NewID(0, nil)
	jc := testutils.NewJetCoordinatorMock(mc)
	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()

	h := NewMessageHandler(&configuration.Ledger{}, certificate)
	h.JetCoordinator = jc
	h.Bus = mb
	h.JetStorage = s.jetStorage
	h.NodeStorage = s.nodeStorage
	h.DBContext = s.db
	h.PulseTracker = s.pulseTracker
	h.ObjectStorage = s.objectStorage

	err := h.Init(s.ctx)
	require.NoError(s.T(), err)

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetStorageFunc = func(ctx context.Context, p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	h.RecentStorageProvider = provideMock

	rep, err := h.handleHasPendingRequests(contextWithJet(s.ctx, jetID), &message.Parcel{
		Msg:         &msg,
		PulseNumber: core.FirstPulseNumber + 1,
	})
	require.NoError(s.T(), err)
	has, ok := rep.(*reply.HasPendingRequests)
	require.True(s.T(), ok)
	assert.True(s.T(), has.Has)
}

func (s *handlerSuite) TestMessageHandler_HandleGetCode_Redirects() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()

	tf := testutils.NewDelegationTokenFactoryMock(mc)
	jc := testutils.NewJetCoordinatorMock(mc)
	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()

	recentStorageMock := recentstorage.NewRecentStorageMock(s.T())
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	certificate := testutils.NewCertificateMock(s.T())
	certificate.GetRoleMock.Return(core.StaticRoleLightMaterial)

	tf.IssueGetCodeRedirectMock.Return(&delegationtoken.GetCodeRedirectToken{Signature: []byte{1, 2, 3}}, nil)

	h := NewMessageHandler(&configuration.Ledger{
		LightChainLimit: 2,
	}, certificate)
	h.JetCoordinator = jc
	h.DelegationTokenFactory = tf
	h.Bus = mb
	h.JetStorage = s.jetStorage
	h.NodeStorage = s.nodeStorage
	h.DBContext = s.db
	h.PulseTracker = s.pulseTracker
	h.ObjectStorage = s.objectStorage
	err := h.Init(s.ctx)
	require.NoError(s.T(), err)

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetStorageFunc = func(ctx context.Context, p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	h.RecentStorageProvider = provideMock

	jetID := *jet.NewID(0, nil)
	msg := message.GetCode{
		Code: *genRandomRef(core.FirstPulseNumber),
	}

	s.T().Run("redirects to light before limit threshold", func(t *testing.T) {
		err := s.pulseTracker.AddPulse(s.ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 1})
		require.NoError(t, err)
		lightRef := genRandomRef(0)
		jc.LightExecutorForJetMock.Set(nil)
		jc.LightExecutorForJetFunc = func(c context.Context, j core.RecordID, p core.PulseNumber) (*core.RecordRef, error) {
			return lightRef, nil
		}
		rep, err := h.handleGetCode(s.ctx, &message.Parcel{
			Msg:         &msg,
			PulseNumber: core.FirstPulseNumber + 1,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetCodeRedirectReply)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetCodeRedirectToken)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, lightRef, redirect.GetReceiver())
	})

	s.T().Run("redirects to heavy after limit threshold", func(t *testing.T) {
		err = s.pulseTracker.AddPulse(s.ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 2})
		require.NoError(t, err)
		heavyRef := genRandomRef(0)
		jc.HeavyMock.Return(heavyRef, nil)
		rep, err := h.handleGetCode(contextWithJet(s.ctx, jetID), &message.Parcel{
			Msg:         &msg,
			PulseNumber: core.FirstPulseNumber + 2,
		})
		require.NoError(t, err)
		redirect, ok := rep.(*reply.GetCodeRedirectReply)
		require.True(t, ok)
		token, ok := redirect.Token.(*delegationtoken.GetCodeRedirectToken)
		assert.Equal(t, []byte{1, 2, 3}, token.Signature)
		assert.Equal(t, heavyRef, redirect.GetReceiver())
	})
}

func (s *handlerSuite) TestMessageHandler_HandleRegisterChild_FetchesIndexFromHeavy() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	jetID := *jet.NewID(0, nil)

	recentStorageMock := recentstorage.NewRecentStorageMock(s.T())
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	certificate := testutils.NewCertificateMock(s.T())
	certificate.GetRoleMock.Return(core.StaticRoleLightMaterial)

	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()
	jc := testutils.NewJetCoordinatorMock(mc)
	h := NewMessageHandler(&configuration.Ledger{
		LightChainLimit: 2,
	}, certificate)
	h.JetStorage = s.jetStorage
	h.NodeStorage = s.nodeStorage
	h.DBContext = s.db
	h.PulseTracker = s.pulseTracker
	h.ObjectStorage = s.objectStorage

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetStorageFunc = func(ctx context.Context, p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	h.RecentStorageProvider = provideMock

	objIndex := index.ObjectLifeline{LatestState: genRandomID(0), State: record.StateActivation}
	childRecord := record.ChildRecord{
		Ref:       *genRandomRef(0),
		PrevChild: nil,
	}
	amendHash := s.scheme.ReferenceHasher()
	_, err := childRecord.WriteHashData(amendHash)
	require.NoError(s.T(), err)
	childID := core.NewRecordID(0, amendHash.Sum(nil))

	msg := message.RegisterChild{
		Record: record.SerializeRecord(&childRecord),
		Parent: *genRandomRef(0),
	}

	mb.SendFunc = func(c context.Context, gm core.Message, o *core.MessageSendOptions) (r core.Reply, r1 error) {
		if m, ok := gm.(*message.GetObjectIndex); ok {
			assert.Equal(s.T(), msg.Parent, m.Object)
			buf, err := index.EncodeObjectLifeline(&objIndex)
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

	idx, err := s.objectStorage.GetObjectIndex(s.ctx, jetID, msg.Parent.Record(), false)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), childID, idx.ChildPointer)
}

func (s *handlerSuite) TestMessageHandler_HandleRegisterChild_IndexStateUpdated() {
	// Arrange
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	jetID := *jet.NewID(0, nil)

	recentStorageMock := recentstorage.NewRecentStorageMock(s.T())
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetStorageFunc = func(ctx context.Context, p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	certificate := testutils.NewCertificateMock(s.T())
	certificate.GetRoleMock.Return(core.StaticRoleLightMaterial)

	h := NewMessageHandler(&configuration.Ledger{
		LightChainLimit: 2,
	}, certificate)
	h.JetStorage = s.jetStorage
	h.NodeStorage = s.nodeStorage
	h.DBContext = s.db
	h.PulseTracker = s.pulseTracker
	h.ObjectStorage = s.objectStorage
	h.RecentStorageProvider = provideMock

	objIndex := index.ObjectLifeline{
		LatestState:  genRandomID(0),
		State:        record.StateActivation,
		LatestUpdate: core.FirstPulseNumber,
	}
	childRecord := record.ChildRecord{
		Ref:       *genRandomRef(0),
		PrevChild: nil,
	}
	msg := message.RegisterChild{
		Record: record.SerializeRecord(&childRecord),
		Parent: *genRandomRef(0),
	}

	err := s.objectStorage.SetObjectIndex(s.ctx, jetID, msg.Parent.Record(), &objIndex)
	require.NoError(s.T(), err)

	// Act
	_, err = h.handleRegisterChild(contextWithJet(s.ctx, jetID), &message.Parcel{
		Msg:         &msg,
		PulseNumber: core.FirstPulseNumber + 100,
	})
	require.NoError(s.T(), err)

	// Assert
	idx, err := s.objectStorage.GetObjectIndex(s.ctx, jetID, msg.Parent.Record(), false)
	require.NoError(s.T(), err)
	require.Equal(s.T(), int(idx.LatestUpdate), core.FirstPulseNumber+100)
}

const testDropSize uint64 = 100

func addDropSizeToDB(s *handlerSuite, jetID core.RecordID) {
	dropSizeData := &jet.DropSize{
		JetID:    jetID,
		PulseNo:  core.FirstPulseNumber,
		DropSize: testDropSize,
	}

	cryptoServiceMock := testutils.NewCryptographyServiceMock(s.T())
	cryptoServiceMock.SignFunc = func(p []byte) (r *core.Signature, r1 error) {
		signature := core.SignatureFromBytes(nil)
		return &signature, nil
	}

	hasher := platformpolicy.NewPlatformCryptographyScheme().IntegrityHasher()
	_, err := dropSizeData.WriteHashData(hasher)
	require.NoError(s.T(), err)

	signature, err := cryptoServiceMock.Sign(hasher.Sum(nil))
	require.NoError(s.T(), err)

	dropSizeData.Signature = signature.Bytes()

	err = s.dropStorage.AddDropSize(s.ctx, dropSizeData)
	require.NoError(s.T(), err)
}

func (s *handlerSuite) TestMessageHandler_HandleHotRecords() {
	mc := minimock.NewController(s.T())
	jetID := testutils.RandomJet()

	err := s.pulseTracker.AddPulse(s.ctx, core.Pulse{PulseNumber: core.FirstPulseNumber + 1})
	require.NoError(s.T(), err)

	jc := testutils.NewJetCoordinatorMock(mc)

	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()

	firstID := core.NewRecordID(core.FirstPulseNumber, []byte{1, 2, 3})
	secondId := record.NewRecordIDFromRecord(s.scheme, core.FirstPulseNumber, &record.CodeRecord{})

	firstIndex, _ := index.EncodeObjectLifeline(&index.ObjectLifeline{
		LatestState: firstID,
	})
	err = s.objectStorage.SetObjectIndex(s.ctx, jetID, firstID, &index.ObjectLifeline{
		LatestState: firstID,
	})

	dropSizeHistory, err := s.dropStorage.GetDropSizeHistory(s.ctx, jetID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), jet.DropSizeHistory{}, dropSizeHistory)
	addDropSizeToDB(s, jetID)

	dropSizeHistory, err = s.dropStorage.GetDropSizeHistory(s.ctx, jetID)
	require.NoError(s.T(), err)

	obj := core.RecordID{}
	hotIndexes := &message.HotData{
		Jet:         *core.NewRecordRef(core.DomainID, jetID),
		PulseNumber: core.FirstPulseNumber,
		RecentObjects: map[core.RecordID]*message.HotIndex{
			*firstID: {
				Index: firstIndex,
				TTL:   320,
			},
		},
		PendingRequests: map[core.RecordID]map[core.RecordID][]byte{
			obj: {
				*secondId: record.SerializeRecord(&record.CodeRecord{}),
			},
		},
		Drop:               jet.JetDrop{Pulse: core.FirstPulseNumber, Hash: []byte{88}},
		DropJet:            jetID,
		JetDropSizeHistory: dropSizeHistory,
	}

	recentStorageMock := recentstorage.NewRecentStorageMock(s.T())
	recentStorageMock.AddPendingRequestFunc = func(ctx context.Context, o, p core.RecordID) {
		require.Equal(s.T(), o, obj)
		require.Equal(s.T(), p, *secondId)
	}
	recentStorageMock.AddObjectWithTLLFunc = func(ctx context.Context, p core.RecordID, ttl int) {
		require.Equal(s.T(), p, *firstID)
		require.Equal(s.T(), 320, ttl)
	}
	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetStorageFunc = func(ctx context.Context, p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	certificate := testutils.NewCertificateMock(s.T())
	certificate.GetRoleMock.Return(core.StaticRoleLightMaterial)

	h := NewMessageHandler(&configuration.Ledger{}, certificate)
	h.JetCoordinator = jc
	h.RecentStorageProvider = provideMock
	h.Bus = mb
	h.JetStorage = s.jetStorage
	h.NodeStorage = s.nodeStorage
	h.DBContext = s.db
	h.PulseTracker = s.pulseTracker
	h.ObjectStorage = s.objectStorage
	h.DropStorage = s.dropStorage

	err = h.Init(s.ctx)
	require.NoError(s.T(), err)

	res, err := h.handleHotRecords(s.ctx, &message.Parcel{Msg: hotIndexes})

	require.NoError(s.T(), err)
	require.Equal(s.T(), res, &reply.OK{})

	savedDrop, err := h.DropStorage.GetDrop(s.ctx, jetID, core.FirstPulseNumber)
	require.NoError(s.T(), err)
	require.Equal(s.T(), &jet.JetDrop{Pulse: core.FirstPulseNumber, Hash: []byte{88}}, savedDrop)

	// check drop size list
	dropSizeHistory, err = s.dropStorage.GetDropSizeHistory(s.ctx, jetID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), testDropSize, dropSizeHistory[0].DropSize)
	require.Equal(s.T(), jetID, dropSizeHistory[0].JetID)
	require.Equal(s.T(), core.FirstPulseNumber, int(dropSizeHistory[0].PulseNo))

	recentStorageMock.MinimockFinish()

}

func (s *handlerSuite) TestMessageHandler_HandleValidationCheck() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	jetID := *jet.NewID(0, nil)

	recentStorageMock := recentstorage.NewRecentStorageMock(s.T())
	recentStorageMock.AddPendingRequestMock.Return()
	recentStorageMock.AddObjectMock.Return()
	recentStorageMock.RemovePendingRequestMock.Return()

	nodeMock := network.NewNodeMock(s.T())
	nodeMock.RoleMock.Return(core.StaticRoleLightMaterial)
	nodeNetworkMock := network.NewNodeNetworkMock(s.T())
	nodeNetworkMock.GetOriginMock.Return(nodeMock)

	jc := testutils.NewJetCoordinatorMock(mc)

	certificate := testutils.NewCertificateMock(s.T())
	certificate.GetRoleMock.Return(core.StaticRoleLightMaterial)

	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()
	h := NewMessageHandler(&configuration.Ledger{
		LightChainLimit: 3,
	}, certificate)
	h.JetCoordinator = jc
	h.Bus = mb
	h.JetStorage = s.jetStorage
	h.NodeStorage = s.nodeStorage
	h.DBContext = s.db
	h.PulseTracker = s.pulseTracker
	h.ObjectStorage = s.objectStorage

	err := h.Init(s.ctx)
	require.NoError(s.T(), err)

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetStorageFunc = func(ctx context.Context, p core.RecordID) (r recentstorage.RecentStorage) {
		return recentStorageMock
	}

	h.RecentStorageProvider = provideMock

	s.T().Run("returns not ok when not valid", func(t *testing.T) {
		validatedStateID, err := s.objectStorage.SetRecord(s.ctx, jetID, 0, &record.ObjectAmendRecord{})
		require.NoError(t, err)

		msg := message.ValidationCheck{
			Object:              *genRandomRef(0),
			ValidatedState:      *validatedStateID,
			LatestStateApproved: genRandomID(0),
		}

		rep, err := h.handleValidationCheck(contextWithJet(s.ctx, jetID), &message.Parcel{
			Msg: &msg,
		})
		require.NoError(t, err)
		_, ok := rep.(*reply.NotOK)
		assert.True(t, ok)
	})

	s.T().Run("returns ok when valid", func(t *testing.T) {
		approvedStateID := *genRandomID(0)
		validatedStateID, err := s.objectStorage.SetRecord(s.ctx, jetID, 0, &record.ObjectAmendRecord{
			PrevState: approvedStateID,
		})
		require.NoError(t, err)

		msg := message.ValidationCheck{
			Object:              *genRandomRef(0),
			ValidatedState:      *validatedStateID,
			LatestStateApproved: &approvedStateID,
		}

		rep, err := h.handleValidationCheck(contextWithJet(s.ctx, jetID), &message.Parcel{
			Msg: &msg,
		})
		require.NoError(t, err)
		_, ok := rep.(*reply.OK)
		assert.True(t, ok)
	})
}

func (s *handlerSuite) TestMessageHandler_HandleJetDrop_SaveJet() {
	// Arrange
	mc := minimock.NewController(s.T())
	defer mc.Finish()

	jetID := jet.NewID(0, []byte{2})
	msg := message.JetDrop{
		JetID: *jetID,
	}
	expectedSetId := jet.IDSet{
		*jetID: struct{}{},
	}

	certificate := testutils.NewCertificateMock(s.T())
	certificate.GetRoleMock.Return(core.StaticRoleLightMaterial)

	h := NewMessageHandler(&configuration.Ledger{
		LightChainLimit: 3,
	}, certificate)
	h.JetStorage = s.jetStorage
	h.NodeStorage = s.nodeStorage
	h.DBContext = s.db
	h.PulseTracker = s.pulseTracker
	h.ObjectStorage = s.objectStorage

	// Act
	response, err := h.handleJetDrop(s.ctx, &message.Parcel{Msg: &msg})
	require.NoError(s.T(), err)

	idSet, err := s.jetStorage.GetJets(s.ctx)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), idSet)

	// Assert
	require.Equal(s.T(), &reply.OK{}, response)
	for id := range expectedSetId {
		require.True(s.T(), idSet.Has(id))
	}
}

func (s *handlerSuite) TestMessageHandler_HandleJetDrop_SaveJet_ExistingMap() {
	// Arrange
	// ctx := inslogger.TestContext(t)
	mc := minimock.NewController(s.T())
	// db, cleaner := storagetest.TmpDB(ctx, t)
	defer mc.Finish()

	jetID := jet.NewID(0, []byte{2})
	secondJetID := jet.NewID(0, []byte{3})
	msg := message.JetDrop{
		JetID: *jetID,
	}
	secondMsg := message.JetDrop{
		JetID: *secondJetID,
	}
	expectedSetId := jet.IDSet{
		*jetID:       struct{}{},
		*secondJetID: struct{}{},
	}

	certificate := testutils.NewCertificateMock(s.T())
	certificate.GetRoleMock.Return(core.StaticRoleLightMaterial)

	h := NewMessageHandler(&configuration.Ledger{
		LightChainLimit: 3,
	}, certificate)
	h.JetStorage = s.jetStorage
	h.NodeStorage = s.nodeStorage
	h.DBContext = s.db
	h.PulseTracker = s.pulseTracker
	h.ObjectStorage = s.objectStorage

	// Act
	response, err := h.handleJetDrop(s.ctx, &message.Parcel{Msg: &msg})
	require.NoError(s.T(), err)
	require.Equal(s.T(), &reply.OK{}, response)

	secondResponse, err := h.handleJetDrop(s.ctx, &message.Parcel{Msg: &secondMsg})
	require.NoError(s.T(), err)
	require.Equal(s.T(), &reply.OK{}, secondResponse)

	idSet, err := s.jetStorage.GetJets(s.ctx)
	require.NoError(s.T(), err)
	require.NotNil(s.T(), idSet)

	// Assert
	for id := range expectedSetId {
		require.True(s.T(), idSet.Has(id))
	}
}

func (s *handlerSuite) TestMessageHandler_HandleGetRequest() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()

	jetID := *jet.NewID(0, nil)

	req := record.RequestRecord{
		Payload: []byte{1, 2, 3},
		Object:  *genRandomID(0),
	}
	reqID, err := s.objectStorage.SetRecord(s.ctx, jetID, core.FirstPulseNumber, &req)

	msg := message.GetRequest{
		Request: *reqID,
	}
	certificate := testutils.NewCertificateMock(s.T())
	certificate.GetRoleMock.Return(core.StaticRoleLightMaterial)

	h := NewMessageHandler(&configuration.Ledger{}, certificate)
	h.ObjectStorage = s.objectStorage

	rep, err := h.handleGetRequest(contextWithJet(s.ctx, jetID), &message.Parcel{
		Msg:         &msg,
		PulseNumber: core.FirstPulseNumber + 1,
	})
	require.NoError(s.T(), err)
	reqReply, ok := rep.(*reply.Request)
	require.True(s.T(), ok)
	assert.Equal(s.T(), req, *record.DeserializeRecord(reqReply.Record).(*record.RequestRecord))
}
