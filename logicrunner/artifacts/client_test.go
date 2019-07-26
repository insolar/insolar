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

	wmMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/platformpolicy"
)

func TestClientImplements(t *testing.T) {
	require.Implements(t, (*Client)(nil), &client{})
}

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

	var err error
	s.tmpDir1, err = ioutil.TempDir("", "bdb-test-")
	if err != nil {
		s.T().Error("Can't create TempDir", err)
	}

	s.badgerDB1, err = store.NewBadgerDB(s.tmpDir1)
	if err != nil {
		s.T().Error("Can't NewBadgerDB", err)
	}

	dropStorage := drop.NewDB(s.badgerDB1)
	s.dropAccessor = dropStorage
	s.dropModifier = dropStorage

	s.tmpDir2, err = ioutil.TempDir("", "bdb-test-")
	if err != nil {
		s.T().Error("Can't create TempDir", err)
	}

	s.badgerDB2, err = store.NewBadgerDB(s.tmpDir2)
	if err != nil {
		s.T().Error("Can't create NewBadgerDB", err)
	}

	s.cm.Inject(
		s.scheme,
		s.badgerDB2,
		s.jetStorage,
		s.nodeStorage,
		pulse.NewStorageMem(),
		s.dropAccessor,
		s.dropModifier,
	)

	err = s.cm.Init(s.ctx)
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
	// We don't call it explicitly since it's called by component manager
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

func (s *amSuite) TestLedgerArtifactManager_GetIncomingRequest_Success() {
	// Arrange
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	objectRef := gen.Reference()
	requestRef := gen.Reference()

	jc := jet.NewCoordinatorMock(mc)

	pulseAccessor := pulse.NewAccessorMock(s.T())
	pulseAccessor.LatestMock.Return(*insolar.GenesisPulse, nil)

	req := record.IncomingRequest{
		Method: "test",
	}

	finalResponse := &payload.Request{
		RequestID: *requestRef.Record(),
		Request:   record.Wrap(req),
	}
	reqMsg, err := payload.NewMessage(finalResponse)
	require.NoError(s.T(), err)

	sender := bus.NewSenderMock(s.T())
	sender.SendRoleFunc = func(_ context.Context, msg *wmMessage.Message, role insolar.DynamicRole, n insolar.Reference) (r <-chan *wmMessage.Message, r1 func()) {
		require.Equal(s.T(), insolar.DynamicRoleLightExecutor, role)

		getReq := payload.GetRequest{}
		err := getReq.Unmarshal(msg.Payload)
		require.NoError(s.T(), err)

		require.Equal(s.T(), *requestRef.Record(), getReq.RequestID)
		require.Equal(s.T(), objectRef, n)

		meta := payload.Meta{Payload: reqMsg.Payload}
		buf, err := meta.Marshal()
		require.NoError(s.T(), err)
		reqMsg.Payload = buf
		ch := make(chan *wmMessage.Message, 1)
		ch <- reqMsg
		return ch, func() {}
	}

	am := NewClient(nil)
	am.JetCoordinator = jc
	am.PulseAccessor = pulseAccessor
	am.sender = sender

	// Act
	res, err := am.GetIncomingRequest(inslogger.TestContext(s.T()), objectRef, requestRef)

	// Assert
	require.NoError(s.T(), err)
	require.Equal(s.T(), "test", res.Method)
}

func (s *amSuite) TestLedgerArtifactManager_GetPendings_Success() {
	// Arrange
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	objectRef := gen.Reference()
	requestRef := gen.Reference()

	jc := jet.NewCoordinatorMock(mc)

	pulseAccessor := pulse.NewAccessorMock(s.T())
	pulseAccessor.LatestMock.Return(*insolar.GenesisPulse, nil)

	resultIDs := &payload.IDs{
		IDs: []insolar.ID{*requestRef.Record()},
	}
	resMsg, err := payload.NewMessage(resultIDs)
	require.NoError(s.T(), err)

	sender := bus.NewSenderMock(s.T())
	sender.SendRoleFunc = func(p context.Context, msg *wmMessage.Message, role insolar.DynamicRole, ref insolar.Reference) (r <-chan *wmMessage.Message, r1 func()) {
		getPendings := payload.GetPendings{}
		err := getPendings.Unmarshal(msg.Payload)
		require.NoError(s.T(), err)

		require.Equal(s.T(), *objectRef.Record(), getPendings.ObjectID)

		meta := payload.Meta{Payload: resMsg.Payload}
		buf, err := meta.Marshal()
		require.NoError(s.T(), err)
		resMsg.Payload = buf
		ch := make(chan *wmMessage.Message, 1)
		ch <- resMsg
		return ch, func() {}
	}

	am := NewClient(nil)
	am.JetCoordinator = jc
	am.PulseAccessor = pulseAccessor
	am.sender = sender

	// Act
	res, err := am.GetPendings(inslogger.TestContext(s.T()), objectRef)

	// Assert
	require.NoError(s.T(), err)
	require.Equal(s.T(), []insolar.Reference{requestRef}, res)
}

func (s *amSuite) TestLedgerArtifactManager_HasPendings_Success() {
	// Arrange
	mc := minimock.NewController(s.T())
	defer mc.Finish()
	objectRef := gen.Reference()

	resultHas := &payload.PendingsInfo{
		HasPendings: true,
	}
	resMsg, err := payload.NewMessage(resultHas)
	require.NoError(s.T(), err)

	sender := bus.NewSenderMock(s.T())
	sender.SendRoleFunc = func(p context.Context, msg *wmMessage.Message, role insolar.DynamicRole, ref insolar.Reference) (r <-chan *wmMessage.Message, r1 func()) {
		hasPendings := payload.HasPendings{}
		err := hasPendings.Unmarshal(msg.Payload)
		require.NoError(s.T(), err)

		require.Equal(s.T(), *objectRef.Record(), hasPendings.ObjectID)

		meta := payload.Meta{Payload: resMsg.Payload}
		buf, err := meta.Marshal()
		require.NoError(s.T(), err)
		resMsg.Payload = buf
		ch := make(chan *wmMessage.Message, 1)
		ch <- resMsg
		return ch, func() {}
	}

	am := NewClient(sender)

	// Act
	res, err := am.HasPendings(inslogger.TestContext(s.T()), objectRef)

	// Assert
	require.NoError(s.T(), err)
	require.Equal(s.T(), true, res)
}
