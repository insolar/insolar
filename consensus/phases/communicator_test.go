//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package phases

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/platformpolicy/keys"
	"github.com/insolar/insolar/testutils"
	networkUtils "github.com/insolar/insolar/testutils/network"
)

type communicatorSuite struct {
	suite.Suite
	componentManager component.Manager
	communicator     Communicator
	originNode       insolar.NetworkNode
	participants     []insolar.NetworkNode
	hostNetworkMock  *networkUtils.HostNetworkMock

	consensusNetworkMock *networkUtils.ConsensusNetworkMock
	pulseHandlerMock     *networkUtils.PulseHandlerMock
}

func NewSuite() *communicatorSuite {
	return &communicatorSuite{
		Suite:        suite.Suite{},
		communicator: NewCommunicator(),
		participants: nil,
	}
}

func (s *communicatorSuite) SetupTest() {
	s.consensusNetworkMock = networkUtils.NewConsensusNetworkMock(s.T())
	s.pulseHandlerMock = networkUtils.NewPulseHandlerMock(s.T())
	s.originNode = makeRandomNode()
	nodeN := networkUtils.NewNodeKeeperMock(s.T())

	cryptoServ := testutils.NewCryptographyServiceMock(s.T())
	cryptoServ.SignFunc = func(p []byte) (r *insolar.Signature, r1 error) {
		signature := insolar.SignatureFromBytes(nil)
		return &signature, nil
	}
	cryptoServ.VerifyFunc = func(p keys.PublicKey, p1 insolar.Signature, p2 []byte) (r bool) {
		return true
	}

	s.consensusNetworkMock.RegisterPacketHandlerMock.Set(func(p packets.PacketType, p1 network.ConsensusPacketHandler) {

	})

	s.consensusNetworkMock.StartMock.Set(func(context.Context) error { return nil })

	s.consensusNetworkMock.GetNodeIDMock.Set(func() (r insolar.Reference) {
		return s.originNode.ID()
	})

	s.pulseHandlerMock.HandlePulseMock.Set(func(p context.Context, p1 insolar.Pulse) {

	})

	s.componentManager.Inject(nodeN, cryptoServ, s.communicator, s.consensusNetworkMock, s.pulseHandlerMock)
	err := s.componentManager.Start(context.TODO())
	s.NoError(err)
}

func makeRandomNode() insolar.NetworkNode {
	return node.NewNode(testutils.RandomRef(), insolar.StaticRoleUnknown, nil, "127.0.0.1:5432", "")
}

func (s *communicatorSuite) TestExchangeData() {
	s.Assert().NotNil(s.communicator)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	result, err := s.communicator.ExchangePhase1(ctx, nil, s.participants, &packets.Phase1Packet{})
	s.Assert().NoError(err)
	s.NotEqual(0, len(result))
}

func TestNaiveCommunicator(t *testing.T) {
	suite.Run(t, NewSuite())
}
