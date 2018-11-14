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
	"testing"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/network/transport/packet/types"
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

	nodenetwork.NewNode()
	s.consensusNetworkMock = networkUtils.NewConsensusNetworkMock(s.T())
	s.pulseHandlerMock = networkUtils.NewPulseHandlerMock(s.T())

	s.consensusNetworkMock.RegisterRequestHandlerMock.Set(func(p types.PacketType, p1 network.ConsensusRequestHandler) {

	})

	s.consensusNetworkMock.NewRequestBuilderMock.Set(func() (r network.RequestBuilder) {
		return &hostnetwork.Builder{}
	})

	s.componentManager.Register(s.communicator, s.consensusNetworkMock, s.pulseHandlerMock)
	err := s.componentManager.Start(nil)
	s.NoError(err)
}

func (s *communicatorSuite) TestExchangeData() {
	s.Assert().NotNil(s.communicator)
	_, err := s.communicator.ExchangeData(context.Background(), s.participants, packets.Phase1Packet{})
	s.Assert().NoError(err)
}

func TestNaiveCommunicator(t *testing.T) {
	suite.Run(t, NewSuite())
}
