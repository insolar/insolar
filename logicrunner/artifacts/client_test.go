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
	"io/ioutil"
	"math/rand"
	"os"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/internal/ledger/store"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/delegationtoken"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
)

type amSuite struct {
	suite.Suite

	cm  *component.Manager
	ctx context.Context

	scheme insolar.PlatformCryptographyScheme

	nodeStorage  node.Accessor
	jetStorage   jet.Storage
	dropModifier drop.Modifier
	dropAccessor drop.Accessor

	tmpDir1 string
	tmpDir2 string

	badgerDB1 *store.BadgerDB
	badgerDB2 *store.BadgerDB
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

	s.scheme = platformpolicy.NewPlatformCryptographyScheme()
	s.jetStorage = jet.NewStore()
	s.nodeStorage = node.NewStorage()

	s.tmpDir1, _ = ioutil.TempDir("", "bdb-test-")

	s.badgerDB1, _ = store.NewBadgerDB(s.tmpDir1)

	dropStorage := drop.NewDB(s.badgerDB1)
	s.dropAccessor = dropStorage
	s.dropModifier = dropStorage

	s.tmpDir2, _ = ioutil.TempDir("", "bdb-test-")

	s.badgerDB2, _ = store.NewBadgerDB(s.tmpDir2)

	s.cm.Inject(
		s.scheme,
		s.badgerDB2,
		s.jetStorage,
		s.nodeStorage,
		pulse.NewStorageMem(),
		s.dropAccessor,
		s.dropModifier,
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

	os.RemoveAll(s.tmpDir1)
	os.RemoveAll(s.tmpDir2)
	s.badgerDB1.Stop(s.ctx)
	// s.badgerDB2.Stop(s.ctx)
}

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

func (s *amSuite) TestLedgerArtifactManager_GetChildren_FollowsRedirect() {
	mc := minimock.NewController(s.T())
	am := NewClient(nil)
	mb := testutils.NewMessageBusMock(mc)

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

	pa := pulse.NewAccessorMock(s.T())
	pa.LatestMock.Return(*insolar.GenesisPulse, nil)
	am.PulseAccessor = pa

	_, err := am.GetChildren(s.ctx, *objRef, nil)
	require.NoError(s.T(), err)
}

func (s *amSuite) TestLedgerArtifactManager_GetRequest_Success() {
	// Arrange
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	objectID := testutils.RandomID()
	requestID := testutils.RandomID()

	node := testutils.RandomRef()

	jc := jet.NewCoordinatorMock(mc)
	jc.NodeForObjectMock.Return(&node, nil)

	pulseAccessor := pulse.NewAccessorMock(s.T())
	pulseAccessor.LatestMock.Return(*insolar.GenesisPulse, nil)

	req := record.IncomingRequest{
		Method: "test",
	}
	virtRec := record.Wrap(req)
	data, err := virtRec.Marshal()
	require.NoError(s.T(), err)
	finalResponse := &reply.Request{Record: data}

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

	am := NewClient(nil)
	am.JetCoordinator = jc
	am.DefaultBus = mb
	am.PulseAccessor = pulseAccessor

	// Act
	_, res, err := am.GetPendingRequest(inslogger.TestContext(s.T()), objectID)

	// Assert
	require.NoError(s.T(), err)
	require.Equal(s.T(), "test", res.Message().(*message.CallMethod).Method)
}
