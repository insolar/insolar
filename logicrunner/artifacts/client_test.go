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

	wmMessage "github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
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
		RequestID: *requestRef.GetLocal(),
		Request:   record.Wrap(&req),
	}
	reqMsg, err := payload.NewMessage(finalResponse)
	require.NoError(s.T(), err)

	sender := bus.NewSenderMock(s.T())
	sender.SendRoleMock.Set(func(_ context.Context, msg *wmMessage.Message, role insolar.DynamicRole, n insolar.Reference) (r <-chan *wmMessage.Message, r1 func()) {
		require.Equal(s.T(), insolar.DynamicRoleLightExecutor, role)

		getReq := payload.GetRequest{}
		err := getReq.Unmarshal(msg.Payload)
		require.NoError(s.T(), err)

		require.Equal(s.T(), *requestRef.GetLocal(), getReq.RequestID)
		require.Equal(s.T(), objectRef, n)

		meta := payload.Meta{Payload: reqMsg.Payload}
		buf, err := meta.Marshal()
		require.NoError(s.T(), err)
		reqMsg.Payload = buf
		ch := make(chan *wmMessage.Message, 1)
		ch <- reqMsg
		return ch, func() {}
	})

	am := NewClient(nil)
	am.JetCoordinator = jc
	am.PulseAccessor = pulseAccessor
	am.sender = sender

	// Act
	request, err := am.GetAbandonedRequest(inslogger.TestContext(s.T()), objectRef, requestRef)

	// Assert
	require.NoError(s.T(), err)
	require.Equal(s.T(), "test", request.(*record.IncomingRequest).Method)
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
		IDs: []insolar.ID{*requestRef.GetLocal()},
	}
	resMsg, err := payload.NewMessage(resultIDs)
	require.NoError(s.T(), err)

	sender := bus.NewSenderMock(s.T())
	sender.SendRoleMock.Set(func(p context.Context, msg *wmMessage.Message, role insolar.DynamicRole, ref insolar.Reference) (r <-chan *wmMessage.Message, r1 func()) {
		getPendings := payload.GetPendings{}
		err := getPendings.Unmarshal(msg.Payload)
		require.NoError(s.T(), err)

		require.Equal(s.T(), *objectRef.GetLocal(), getPendings.ObjectID)

		meta := payload.Meta{Payload: resMsg.Payload}
		buf, err := meta.Marshal()
		require.NoError(s.T(), err)
		resMsg.Payload = buf
		ch := make(chan *wmMessage.Message, 1)
		ch <- resMsg
		return ch, func() {}
	})

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
	sender.SendRoleMock.Set(func(p context.Context, msg *wmMessage.Message, role insolar.DynamicRole, ref insolar.Reference) (r <-chan *wmMessage.Message, r1 func()) {
		hasPendings := payload.HasPendings{}
		err := hasPendings.Unmarshal(msg.Payload)
		require.NoError(s.T(), err)

		require.Equal(s.T(), *objectRef.GetLocal(), hasPendings.ObjectID)

		meta := payload.Meta{Payload: resMsg.Payload}
		buf, err := meta.Marshal()
		require.NoError(s.T(), err)
		resMsg.Payload = buf
		ch := make(chan *wmMessage.Message, 1)
		ch <- resMsg
		return ch, func() {}
	})

	am := NewClient(sender)

	// Act
	res, err := am.HasPendings(inslogger.TestContext(s.T()), objectRef)

	// Assert
	require.NoError(s.T(), err)
	require.Equal(s.T(), true, res)
}
