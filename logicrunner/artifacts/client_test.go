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

package artifacts

import (
	"context"
	"math/rand"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/delegationtoken"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/db"
	"github.com/insolar/insolar/ledger/storage/drop"
	"github.com/insolar/insolar/ledger/storage/genesis"
	"github.com/insolar/insolar/ledger/storage/node"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/insolar/insolar/ledger/storage/storagetest"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
)

type amSuite struct {
	suite.Suite

	cm      *component.Manager
	ctx     context.Context
	cleaner func()
	db      storage.DBContext

	scheme        insolar.PlatformCryptographyScheme
	pulseTracker  storage.PulseTracker
	nodeStorage   node.Accessor
	objectStorage storage.ObjectStorage
	jetStorage    jet.Storage
	dropModifier  drop.Modifier
	dropAccessor  drop.Accessor
	genesisState  genesis.GenesisState
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

	tempDB, cleaner := storagetest.TmpDB(s.ctx, s.T())
	s.cleaner = cleaner
	s.db = tempDB
	s.scheme = platformpolicy.NewPlatformCryptographyScheme()
	s.jetStorage = jet.NewStore()
	s.nodeStorage = node.NewStorage()
	s.pulseTracker = storage.NewPulseTracker()
	s.objectStorage = storage.NewObjectStorage()

	dbStore := db.NewDBWithBadger(tempDB.GetBadgerDB())
	dropStorage := drop.NewStorageDB(dbStore)
	s.dropAccessor = dropStorage
	s.dropModifier = dropStorage
	s.genesisState = genesis.NewGenesisInitializer()

	s.cm.Inject(
		s.scheme,
		s.db,
		db.NewMemoryMockDB(),
		s.jetStorage,
		s.nodeStorage,
		s.pulseTracker,
		s.objectStorage,
		s.dropAccessor,
		s.dropModifier,
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

func (s *amSuite) TestLedgerArtifactManager_GetCodeWithCache() {
	code := []byte("test_code")
	codeRef := testutils.RandomRef()

	mb := testutils.NewMessageBusMock(s.T())
	mb.SendFunc = func(p context.Context, p1 insolar.Message, p3 *insolar.MessageSendOptions) (r insolar.Reply, r1 error) {
		return &reply.Code{
			Code: code,
		}, nil
	}

	jc := testutils.NewJetCoordinatorMock(s.T())
	jc.LightExecutorForJetMock.Return(&insolar.Reference{}, nil)
	jc.MeMock.Return(insolar.Reference{})

	amPulseStorageMock := testutils.NewPulseStorageMock(s.T())
	amPulseStorageMock.CurrentFunc = func(p context.Context) (r *insolar.Pulse, r1 error) {
		pulse, err := s.pulseTracker.GetLatestPulse(p)
		require.NoError(s.T(), err)
		return &pulse.Pulse, err
	}

	am := client{
		DefaultBus:                 mb,
		PulseStorage:               amPulseStorageMock,
		JetCoordinator:             jc,
		PlatformCryptographyScheme: s.scheme,
		senders:                    newLedgerArtifactSenders(),
	}

	desc, err := am.GetCode(s.ctx, codeRef)
	receivedCode, err := desc.Code()
	require.NoError(s.T(), err)
	require.Equal(s.T(), code, receivedCode)

	mb.SendFunc = func(p context.Context, p1 insolar.Message, p3 *insolar.MessageSendOptions) (r insolar.Reply, r1 error) {
		s.T().Fatal("Func must not be called here")
		return nil, nil
	}

	desc, err = am.GetCode(s.ctx, codeRef)
	receivedCode, err = desc.Code()
	require.NoError(s.T(), err)
	require.Equal(s.T(), code, receivedCode)

}

func (s *amSuite) TestLedgerArtifactManager_GetObject_FollowsRedirect() {
	mc := minimock.NewController(s.T())
	am := NewClient()
	mb := testutils.NewMessageBusMock(mc)

	objRef := genRandomRef(0)
	nodeRef := genRandomRef(0)
	mb.SendFunc = func(c context.Context, m insolar.Message, o *insolar.MessageSendOptions) (r insolar.Reply, r1 error) {
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
	am.PulseStorage = makePulseStorage(s)

	_, err := am.GetObject(s.ctx, *objRef, nil, false)

	require.NoError(s.T(), err)
}

func makePulseStorage(s *amSuite) insolar.PulseStorage {
	pulseStorage := storage.NewPulseStorage()
	pulseStorage.PulseTracker = s.pulseTracker
	pulse, err := s.pulseTracker.GetLatestPulse(s.ctx)
	require.NoError(s.T(), err)
	pulseStorage.Set(&pulse.Pulse)

	return pulseStorage
}

func (s *amSuite) TestLedgerArtifactManager_GetChildren_FollowsRedirect() {
	mc := minimock.NewController(s.T())
	am := NewClient()
	mb := testutils.NewMessageBusMock(mc)

	am.PulseStorage = makePulseStorage(s)

	objRef := genRandomRef(0)
	nodeRef := genRandomRef(0)
	mb.SendFunc = func(c context.Context, m insolar.Message, o *insolar.MessageSendOptions) (r insolar.Reply, r1 error) {
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

func (s *amSuite) TestLedgerArtifactManager_RegisterRequest_JetMiss() {
	mc := minimock.NewController(s.T())
	defer mc.Finish()

	cs := platformpolicy.NewPlatformCryptographyScheme()
	am := NewClient()
	am.PlatformCryptographyScheme = cs
	pulseStorageMock := testutils.NewPulseStorageMock(s.T())
	pulseStorageMock.CurrentFunc = func(ctx context.Context) (*insolar.Pulse, error) {
		return &insolar.Pulse{PulseNumber: insolar.FirstPulseNumber}, nil
	}

	am.PulseStorage = pulseStorageMock
	am.JetStorage = s.jetStorage

	s.T().Run("returns error on exceeding retry limit", func(t *testing.T) {
		mb := testutils.NewMessageBusMock(mc)
		am.DefaultBus = mb
		mb.SendMock.Return(&reply.JetMiss{
			JetID: insolar.ID(*insolar.NewJetID(5, []byte{1, 2, 3})),
		}, nil)
		_, err := am.RegisterRequest(s.ctx, *am.GenesisRef(), &message.Parcel{Msg: &message.CallMethod{}})
		require.Error(t, err)
	})

	s.T().Run("returns no error and updates tree when jet miss", func(t *testing.T) {
		b_1101 := byte(0xD0)
		b_11010101 := byte(0xD5)
		mb := testutils.NewMessageBusMock(mc)
		am.DefaultBus = mb
		retries := 3
		mb.SendFunc = func(c context.Context, m insolar.Message, o *insolar.MessageSendOptions) (r insolar.Reply, r1 error) {
			if retries == 0 {
				return &reply.ID{}, nil
			}
			retries--
			return &reply.JetMiss{JetID: insolar.ID(*insolar.NewJetID(4, []byte{b_11010101}))}, nil
		}
		_, err := am.RegisterRequest(s.ctx, *am.GenesisRef(), &message.Parcel{Msg: &message.CallMethod{}})
		require.NoError(t, err)

		jetID, actual := s.jetStorage.ForID(
			s.ctx, insolar.FirstPulseNumber, *insolar.NewID(0, []byte{0xD5}),
		)

		assert.Equal(t, insolar.NewJetID(4, []byte{b_1101}), &jetID, "proper jet ID for record")
		assert.True(t, actual, "jet ID is actual in tree")
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
	pulseStorageMock.CurrentMock.Return(insolar.GenesisPulse, nil)

	var parcel insolar.Parcel = &message.Parcel{PulseNumber: 123987}
	resRecord := object.RequestRecord{
		Parcel: message.ParcelToBytes(parcel),
	}
	finalResponse := &reply.Request{Record: object.SerializeRecord(&resRecord)}

	mb := testutils.NewMessageBusMock(s.T())
	mb.SendFunc = func(p context.Context, p1 insolar.Message, p2 *insolar.MessageSendOptions) (r insolar.Reply, r1 error) {
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

	am := NewClient()
	am.JetCoordinator = jc
	am.DefaultBus = mb
	am.PulseStorage = pulseStorageMock

	// Act
	res, err := am.GetPendingRequest(inslogger.TestContext(s.T()), objectID)

	// Assert
	require.NoError(s.T(), err)
	require.Equal(s.T(), parcel, res)

}
