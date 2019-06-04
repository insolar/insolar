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
	"time"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/ledger/light/recentstorage"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/testutils"
)

type handlerSuite struct {
	suite.Suite

	cm  *component.Manager
	ctx context.Context

	scheme      insolar.PlatformCryptographyScheme
	nodeStorage node.Accessor
	jetStorage  jet.Storage

	dropModifier drop.Modifier
	dropAccessor drop.Accessor

	blobModifier blob.Modifier
	blobAccessor blob.Accessor

	recordModifier object.RecordModifier
	recordAccessor object.RecordAccessor

	indexMemoryStor *object.InMemoryIndex
}

var domainID = *genRandomID(0)

func genRandomID(pulse insolar.PulseNumber) *insolar.ID {
	buff := [insolar.RecordIDSize - insolar.PulseNumberSize]byte{}
	_, err := rand.Read(buff[:])
	if err != nil {
		panic(err)
	}
	return insolar.NewID(pulse, buff[:])
}

func genRefWithID(id *insolar.ID) *insolar.Reference {
	return insolar.NewReference(*id)
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

	s.indexMemoryStor = object.NewInMemoryIndex()

	s.cm.Inject(
		s.scheme,
		s.indexMemoryStor,
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
}

type waiterMock struct {
}

func (*waiterMock) Wait(ctx context.Context, jetID insolar.ID, pulse insolar.PulseNumber) error {
	return nil
}

func (s *handlerSuite) TestMessageHandler_HandleGetDelegate_FetchesIndexFromHeavy() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()

	pendingMock := recentstorage.NewPendingStorageMock(s.T())
	pendingMock.GetRequestsForObjectMock.Return(nil)
	pendingMock.AddPendingRequestMock.Return()
	pendingMock.RemovePendingRequestMock.Return()

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetPendingStorageMock.Return(pendingMock)

	waiterMock := waiterMock{}

	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()
	jc := jet.NewCoordinatorMock(mc)

	h := NewMessageHandler(s.indexMemoryStor, s.indexMemoryStor, s.indexMemoryStor, &configuration.Ledger{
		LightChainLimit: 3,
	})
	h.JetStorage = s.jetStorage
	h.Nodes = s.nodeStorage

	h.RecentStorageProvider = provideMock
	idLock := object.NewIDLockerMock(s.T())
	idLock.LockMock.Return()
	idLock.UnlockMock.Return()
	h.IDLocker = idLock

	delegateType := *genRandomRef(0)
	delegate := *genRandomRef(0)
	objIndex := object.Lifeline{Delegates: []object.LifelineDelegate{{Key: delegateType, Value: delegate}}}
	msg := message.GetDelegate{
		Head:   *genRandomRef(0),
		AsType: delegateType,
	}

	fakeParcel := testutils.NewParcelMock(mc)
	fakeParcel.MessageMock.Return(&msg)
	fakeParcel.PulseMock.Return(insolar.FirstPulseNumber)

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
	h.HotDataWaiter = &waiterMock
	err := h.Init(s.ctx)
	require.NoError(s.T(), err)

	heavyRef := genRandomRef(0)
	jc.HeavyMock.Return(heavyRef, nil)

	rep, err := h.FlowDispatcher.WrapBusHandle(s.ctx, fakeParcel)
	require.NoError(s.T(), err)
	delegateRep, ok := rep.(*reply.Delegate)
	require.True(s.T(), ok)
	assert.Equal(s.T(), delegate, delegateRep.Head)

	idx, err := s.indexMemoryStor.ForID(s.ctx, insolar.FirstPulseNumber, *msg.Head.Record())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), objIndex.Delegates, idx.Delegates)
}

func (s *handlerSuite) TestMessageHandler_HandleHasPendingRequests() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()

	msg := message.GetPendingRequests{
		Object: *genRandomRef(0),
	}
	fakeParcel := testutils.NewParcelMock(mc)
	fakeParcel.MessageMock.Return(&msg)
	fakeParcel.PulseMock.Return(insolar.FirstPulseNumber - 1)

	pendingRequests := []insolar.ID{
		*genRandomID(insolar.FirstPulseNumber - 2),
		*genRandomID(insolar.FirstPulseNumber - 2),
	}

	recentStorageMock := recentstorage.NewPendingStorageMock(s.T())
	recentStorageMock.GetRequestsForObjectMock.Return(pendingRequests)

	jc := jet.NewCoordinatorMock(mc)
	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()

	h := NewMessageHandler(s.indexMemoryStor, s.indexMemoryStor, s.indexMemoryStor, &configuration.Ledger{})
	h.JetCoordinator = jc
	h.Bus = mb
	h.JetStorage = s.jetStorage
	h.Nodes = s.nodeStorage

	err := h.Init(s.ctx)
	require.NoError(s.T(), err)

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetPendingStorageMock.Return(recentStorageMock)

	h.RecentStorageProvider = provideMock

	rep, err := h.FlowDispatcher.WrapBusHandle(s.ctx, fakeParcel)
	require.NoError(s.T(), err)
	has, ok := rep.(*reply.HasPendingRequests)
	require.True(s.T(), ok)
	assert.True(s.T(), has.Has)
}

func (s *handlerSuite) TestMessageHandler_HandleGetPendingRequestID() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()

	msg := message.GetPendingRequestID{
		ObjectID: *genRandomID(0),
	}
	fakeParcel := testutils.NewParcelMock(mc)
	fakeParcel.MessageMock.Return(&msg)
	fakeParcel.PulseMock.Return(insolar.FirstPulseNumber - 1)

	firstID := *genRandomID(insolar.FirstPulseNumber - 2)
	secondID := *genRandomID(insolar.FirstPulseNumber - 2)
	pendingRequests := []insolar.ID{
		firstID,
		secondID,
	}

	recentStorageMock := recentstorage.NewPendingStorageMock(s.T())
	recentStorageMock.GetRequestsForObjectMock.Return(pendingRequests)

	jc := jet.NewCoordinatorMock(mc)
	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()

	h := NewMessageHandler(s.indexMemoryStor, s.indexMemoryStor, s.indexMemoryStor, &configuration.Ledger{})
	h.JetCoordinator = jc
	h.Bus = mb
	h.JetStorage = s.jetStorage
	h.Nodes = s.nodeStorage

	err := h.Init(s.ctx)
	require.NoError(s.T(), err)

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetPendingStorageMock.Return(recentStorageMock)

	h.RecentStorageProvider = provideMock

	// call to object that has pending requests
	rep, err := h.FlowDispatcher.WrapBusHandle(s.ctx, fakeParcel)
	require.NoError(s.T(), err)
	result, ok := rep.(*reply.ID)
	require.True(s.T(), ok)
	assert.Equal(s.T(), firstID, result.ID)

	// call to object that hasn't pending requests
	noPendingRequests := make([]insolar.ID, 0)
	recentStorageMock.GetRequestsForObjectMock.Return(noPendingRequests)
	rep, err = h.FlowDispatcher.WrapBusHandle(s.ctx, fakeParcel)
	require.NoError(s.T(), err)
	replyError, ok := rep.(*reply.Error)
	require.True(s.T(), ok)
	assert.Equal(s.T(), &reply.Error{ErrType: reply.ErrNoPendingRequests}, replyError)
}

func (s *handlerSuite) TestMessageHandler_HandleRegisterChild_FetchesIndexFromHeavy() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	jetID := insolar.ID(*insolar.NewJetID(0, nil))

	pendingMock := recentstorage.NewPendingStorageMock(s.T())

	pendingMock.GetRequestsForObjectMock.Return(nil)
	pendingMock.AddPendingRequestMock.Return()
	pendingMock.RemovePendingRequestMock.Return()

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetPendingStorageMock.Return(pendingMock)

	mb := testutils.NewMessageBusMock(mc)
	mb.MustRegisterMock.Return()
	jc := jet.NewCoordinatorMock(mc)
	h := NewMessageHandler(s.indexMemoryStor, s.indexMemoryStor, s.indexMemoryStor, &configuration.Ledger{
		LightChainLimit: 2,
	})
	h.JetStorage = s.jetStorage
	h.Nodes = s.nodeStorage
	h.RecentStorageProvider = provideMock
	h.PCS = s.scheme
	h.RecordModifier = s.recordModifier

	idLockMock := object.NewIDLockerMock(s.T())
	idLockMock.LockMock.Return()
	idLockMock.UnlockMock.Return()
	h.IDLocker = idLockMock

	objIndex := object.Lifeline{LatestState: genRandomID(0), StateID: record.StateActivation}
	childRecord := record.Child{
		Ref: *genRandomRef(0),
	}

	virtChild := record.Wrap(childRecord)
	data, err := virtChild.Marshal()
	require.NoError(s.T(), err)
	hash := record.HashVirtual(s.scheme.ReferenceHasher(), virtChild)
	childID := insolar.NewID(0, hash)

	msg := message.RegisterChild{
		Record: data,
		Parent: *genRandomRef(0),
	}

	h.JetCoordinator = jc
	h.Bus = mb
	err = h.Init(s.ctx)
	require.NoError(s.T(), err)

	replyTo := make(chan bus.Reply, 1)
	registerChild := proc.NewRegisterChild(insolar.JetID(jetID), &msg, childID.Pulse(), objIndex, replyTo)
	registerChild.Dep.IDLocker = idLockMock
	registerChild.Dep.LifelineIndex = s.indexMemoryStor
	registerChild.Dep.JetCoordinator = jc
	registerChild.Dep.RecordModifier = s.recordModifier
	registerChild.Dep.LifelineStateModifier = s.indexMemoryStor
	registerChild.Dep.PCS = s.scheme

	err = registerChild.Proceed(s.ctx)
	require.NoError(s.T(), err)

	busRep := <-replyTo
	rep := busRep.Reply
	objRep, ok := rep.(*reply.ID)
	require.True(s.T(), ok)
	assert.Equal(s.T(), *childID, objRep.ID)

	idx, err := s.indexMemoryStor.ForID(s.ctx, 0, *msg.Parent.Record())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), childID, idx.ChildPointer)
}

func (s *handlerSuite) TestMessageHandler_HandleRegisterChild_IndexStateUpdated() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	jetID := insolar.ID(*insolar.NewJetID(0, nil))

	pendingMock := recentstorage.NewPendingStorageMock(s.T())

	pendingMock.GetRequestsForObjectMock.Return(nil)
	pendingMock.AddPendingRequestMock.Return()
	pendingMock.RemovePendingRequestMock.Return()

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetPendingStorageMock.Return(pendingMock)

	h := NewMessageHandler(s.indexMemoryStor, s.indexMemoryStor, s.indexMemoryStor, &configuration.Ledger{
		LightChainLimit: 2,
	})
	h.JetStorage = s.jetStorage
	h.Nodes = s.nodeStorage
	h.LifelineIndex = s.indexMemoryStor
	h.LifelineStateModifier = s.indexMemoryStor
	h.RecentStorageProvider = provideMock
	h.PCS = s.scheme
	h.RecordModifier = s.recordModifier

	idLockMock := object.NewIDLockerMock(s.T())
	idLockMock.LockMock.Return()
	idLockMock.UnlockMock.Return()
	h.IDLocker = idLockMock

	objIndex := object.Lifeline{
		LatestState:  genRandomID(0),
		StateID:      record.StateActivation,
		LatestUpdate: insolar.FirstPulseNumber,
		JetID:        insolar.JetID(jetID),
	}
	childRecord := record.Child{
		Ref: *genRandomRef(0),
	}

	virtRec := record.Wrap(childRecord)
	data, err := virtRec.Marshal()
	require.NoError(s.T(), err)

	msg := message.RegisterChild{
		Record: data,
		Parent: *genRandomRef(0),
	}

	pulse := gen.PulseNumber()
	err = s.indexMemoryStor.Set(s.ctx, pulse, *msg.Parent.Record(), objIndex)
	require.NoError(s.T(), err)

	replyTo := make(chan bus.Reply, 1)

	registerChild := proc.NewRegisterChild(insolar.JetID(jetID), &msg, pulse, objIndex, replyTo)
	registerChild.Dep.IDLocker = idLockMock
	registerChild.Dep.LifelineIndex = s.indexMemoryStor
	registerChild.Dep.JetCoordinator = jet.NewCoordinatorMock(mc)
	registerChild.Dep.RecordModifier = s.recordModifier
	registerChild.Dep.LifelineStateModifier = s.indexMemoryStor
	registerChild.Dep.PCS = s.scheme

	err = registerChild.Proceed(s.ctx)
	require.NoError(s.T(), err)

	idx, err := s.indexMemoryStor.ForID(s.ctx, pulse, *msg.Parent.Record())
	require.NoError(s.T(), err)
	require.Equal(s.T(), idx.LatestUpdate, pulse)
}

func (s *handlerSuite) TestMessageHandler_HandleHotRecords() {
	mc := minimock.NewController(s.T())
	jetID := gen.JetID()

	jc := jet.NewCoordinatorMock(mc)

	firstID := insolar.NewID(insolar.FirstPulseNumber, []byte{1, 2, 3})

	codeRec := record.Code{}
	virtCodeRec := record.Wrap(codeRec)
	hash := record.HashVirtual(s.scheme.ReferenceHasher(), virtCodeRec)

	secondID := insolar.NewID(insolar.FirstPulseNumber, hash)

	codeRec = record.Code{}
	virtCodeRec = record.Wrap(codeRec)
	hash = record.HashVirtual(s.scheme.ReferenceHasher(), virtCodeRec)

	thirdID := insolar.NewID(insolar.FirstPulseNumber-1, hash)

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
	err := s.indexMemoryStor.Set(s.ctx, insolar.FirstPulseNumber, *firstID, object.Lifeline{
		LatestState: firstID,
		JetID:       insolar.JetID(jetID),
	})

	hotIndexes := &message.HotData{
		Jet:         *insolar.NewReference(insolar.ID(jetID)),
		PulseNumber: insolar.FirstPulseNumber,
		HotIndexes: []message.HotIndex{
			{
				Index:    firstIndex,
				LastUsed: insolar.PulseNumber(234),
				ObjID:    *firstID,
			},
		},
		PendingRequests: map[insolar.ID]recentstorage.PendingObjectContext{
			*secondID: {},
			*thirdID:  {Active: true},
		},
		Drop: drop.Drop{Pulse: insolar.FirstPulseNumber, Hash: []byte{88}, JetID: jetID},
	}

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

	idxStateModifierMock := object.NewLifelineStateModifierMock(s.T())
	bucketMock := object.NewIndexBucketModifierMock(s.T())
	idxMock := object.NewLifelineIndexMock(s.T())

	bucketMock.SetBucketFunc = func(ctx context.Context, pn insolar.PulseNumber, ib object.IndexBucket) (r error) {
		require.Equal(s.T(), *firstID, ib.ObjID)
		require.Equal(s.T(), insolar.FirstPulseNumber, int(pn))
		require.Equal(s.T(), *firstID, *ib.Lifeline.LatestState)

		return nil
	}

	idxMock.SetFunc = func(p context.Context, p1 insolar.PulseNumber, p2 insolar.ID, p3 object.Lifeline) (r error) {
		require.Equal(s.T(), *firstID, p2)
		require.Equal(s.T(), insolar.FirstPulseNumber, int(p1))
		require.Equal(s.T(), *firstID, *p3.LatestState)

		return nil
	}

	provideMock := recentstorage.NewProviderMock(s.T())
	provideMock.GetPendingStorageMock.Return(pendingMock)

	h := NewMessageHandler(idxMock, bucketMock, idxStateModifierMock, &configuration.Ledger{})
	h.JetCoordinator = jc
	h.RecentStorageProvider = provideMock
	h.Bus = mb
	h.JetStorage = s.jetStorage
	h.Nodes = s.nodeStorage
	h.DropModifier = s.dropModifier

	jr := testutils.NewJetReleaserMock(s.T())
	jr.UnlockMock.Return(nil)
	h.JetReleaser = jr

	err = h.Init(s.ctx)
	require.NoError(s.T(), err)

	replyTo := make(chan bus.Reply, 1)
	p := proc.NewHotData(hotIndexes, replyTo)
	p.Dep.DropModifier = h.DropModifier
	p.Dep.RecentStorageProvider = h.RecentStorageProvider
	p.Dep.MessageBus = h.Bus
	p.Dep.IndexBucketModifier = h.IndexBucketModifier
	p.Dep.JetStorage = h.JetStorage
	p.Dep.JetFetcher = h.jetTreeUpdater
	p.Dep.JetReleaser = h.JetReleaser
	err = p.Proceed(s.ctx)
	require.NoError(s.T(), err)

	resWrapper := <-replyTo
	res := resWrapper.Reply
	require.Equal(s.T(), res, &reply.OK{})

	savedDrop, err := s.dropAccessor.ForPulse(s.ctx, jetID, insolar.FirstPulseNumber)
	require.NoError(s.T(), err)
	require.Equal(s.T(), drop.Drop{Pulse: insolar.FirstPulseNumber, Hash: []byte{88}, JetID: jetID}, savedDrop)

	mc.Wait(1 * time.Minute)
	pendingMock.MinimockFinish()
}

func (s *handlerSuite) TestMessageHandler_HandleGetRequest() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()

	jetID := insolar.ID(*insolar.NewJetID(0, nil))

	req := record.Request{
		Object: genRandomRef(0),
	}

	virtRec := record.Wrap(req)
	hash := record.HashVirtual(s.scheme.ReferenceHasher(), virtRec)
	reqID := insolar.NewID(insolar.FirstPulseNumber, hash)

	rec := record.Material{
		Virtual: &virtRec,
		JetID:   insolar.JetID(jetID),
	}
	err := s.recordModifier.Set(s.ctx, *reqID, rec)
	require.NoError(s.T(), err)

	h := NewMessageHandler(s.indexMemoryStor, s.indexMemoryStor, s.indexMemoryStor, &configuration.Ledger{})
	h.RecordAccessor = s.recordAccessor

	replyTo := make(chan bus.Reply, 1)
	procGetRequest := proc.NewGetRequest(*reqID, replyTo)
	procGetRequest.Dep.RecordAccessor = s.recordAccessor

	err = procGetRequest.Proceed(s.ctx)

	require.NoError(s.T(), err)
	res := <-replyTo
	reqReply, ok := (res.Reply).(*reply.Request)
	require.True(s.T(), ok)

	vRec := record.Virtual{}
	err = vRec.Unmarshal(reqReply.Record)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), &req, record.Unwrap(&vRec))
}
