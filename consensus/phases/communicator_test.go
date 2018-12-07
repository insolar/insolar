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

package phases

import (
	"context"
	"crypto"
	"testing"
	"time"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/network/transport/packet/types"
	"github.com/insolar/insolar/testutils"
	networkUtils "github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/suite"
)

type communicatorSuite struct {
	suite.Suite
	componentManager component.Manager
	communicator     Communicator
	originNode       core.Node
	participants     []core.Node
	hostNetworkMock  *networkUtils.HostNetworkMock

	consensusNetworkMock *networkUtils.ConsensusNetworkMock
	pulseHandlerMock     *networkUtils.PulseHandlerMock
}

func NewSuite() *communicatorSuite {
	return &communicatorSuite{
		Suite:        suite.Suite{},
		communicator: NewNaiveCommunicator(),
		participants: nil,
	}
}

func (s *communicatorSuite) SetupTest() {
	s.consensusNetworkMock = networkUtils.NewConsensusNetworkMock(s.T())
	s.pulseHandlerMock = networkUtils.NewPulseHandlerMock(s.T())
	s.originNode = makeRandomNode()
	nodeN := networkUtils.NewNodeNetworkMock(s.T())

	cryptoServ := testutils.NewCryptographyServiceMock(s.T())
	cryptoServ.SignFunc = func(p []byte) (r *core.Signature, r1 error) {
		signature := core.SignatureFromBytes(nil)
		return &signature, nil
	}
	cryptoServ.VerifyFunc = func(p crypto.PublicKey, p1 core.Signature, p2 []byte) (r bool) {
		return true
	}

	s.consensusNetworkMock.RegisterRequestHandlerMock.Set(func(p types.PacketType, p1 network.ConsensusRequestHandler) {
	})

	s.consensusNetworkMock.NewRequestBuilderMock.Set(func() (r network.RequestBuilder) {
		return &hostnetwork.Builder{}
	})

	s.consensusNetworkMock.GetNodeIDMock.Set(func() (r core.RecordRef) {
		return s.originNode.ID()
	})

	s.pulseHandlerMock.HandlePulseMock.Set(func(p context.Context, p1 core.Pulse) {

	})

	s.componentManager.Inject(nodeN, cryptoServ, s.communicator, s.consensusNetworkMock, s.pulseHandlerMock)
	err := s.componentManager.Start(context.TODO())
	s.NoError(err)
}

func makeRandomNode() core.Node {
	return nodenetwork.NewNode(testutils.RandomRef(), core.StaticRoleUnknown, nil, "127.0.0.1", "")
}

func (s *communicatorSuite) TestExchangeData() {
	s.Assert().NotNil(s.communicator)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	result, _, err := s.communicator.ExchangePhase1(ctx, s.participants, &packets.Phase1Packet{})
	s.Assert().NoError(err)
	s.NotEqual(0, len(result))
}

func TestNaiveCommunicator(t *testing.T) {
	suite.Run(t, NewSuite())
}
