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

package servicenetwork

import (
	"context"
	"testing"

	"github.com/insolar/insolar/consensus/phases"
	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/suite"
)

func (s *testSuite) TestNetworkConsensus3Times() {
	s.waitForConsensus(3)
}

func (s *testSuite) TestNodeConnect() {
	s.T().Skip()
	s.preInitNode(s.fixture().testNode)

	s.InitTestNode()
	s.StartTestNode()
	defer func() {
		s.StopTestNode()
	}()

	s.waitForConsensus(1)

	activeNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))

	s.waitForConsensus(1)

	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount()+1, len(activeNodes))

	s.waitForConsensus(2)

	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount()+1, len(activeNodes))
}

func (s *testSuite) TestNodeLeave() {
	s.T().Skip()
	s.preInitNode(s.fixture().testNode)

	s.InitTestNode()
	s.StartTestNode()
	defer func() {
		s.StopTestNode()
	}()

	s.waitForConsensus(1)

	activeNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))

	s.waitForConsensus(1)

	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount()+1, len(activeNodes))

	s.fixture().testNode.serviceNetwork.GracefulStop(context.Background())

	s.waitForConsensus(2)

	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func TestServiceNetworkOneBootstrap(t *testing.T) {
	s := NewTestSuite(1, 0)
	suite.Run(t, s)
}

func TestServiceNetworkManyBootstraps(t *testing.T) {
	s := NewTestSuite(15, 0)
	suite.Run(t, s)
}

func TestServiceNetworkManyNodes(t *testing.T) {
	t.Skip("tmp 123")

	s := NewTestSuite(3, 5)
	suite.Run(t, s)
}

// Full timeout test
type FullTimeoutPhaseManager struct {
}

func (ftpm *FullTimeoutPhaseManager) OnPulse(ctx context.Context, pulse *core.Pulse) error {
	return nil
}

func (s *testSuite) TestFullTimeOut() {
	if len(s.fixture().bootstrapNodes) < 3 {
		s.T().Skip("skip test for bootstrap nodes < 3")
	}

	wrapper := s.fixture().bootstrapNodes[1].serviceNetwork.PhaseManager.(*phaseManagerWrapper)
	wrapper.original = &FullTimeoutPhaseManager{}
	s.fixture().bootstrapNodes[1].serviceNetwork.PhaseManager = wrapper

	s.preInitNode(s.fixture().testNode)

	s.InitTestNode()
	s.StartTestNode()
	defer s.StopTestNode()

	s.waitForConsensus(1)

	activeNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))

	s.waitForConsensus(1)

	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount()+1-1, len(activeNodes))
}

// Partial timeout

func (s *testSuite) TestPartialTimeOut() {
	if len(s.fixture().bootstrapNodes) < 3 {
		s.T().Skip("skip test for bootstrap nodes < 3")
	}

	comm := s.fixture().bootstrapNodes[0].serviceNetwork.PhaseManager.(*phaseManagerWrapper).original.(*phases.Phases).FirstPhase.Communicator
	wrapper := &CommunicatorMock{comm, PartialPositive1Phase}
	s.fixture().bootstrapNodes[0].serviceNetwork.PhaseManager.(*phaseManagerWrapper).original.(*phases.Phases).FirstPhase.Communicator = wrapper

	s.preInitNode(s.fixture().testNode)

	s.InitTestNode()
	s.StartTestNode()
	defer s.StopTestNode()

	activeNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
	s.waitForConsensus(1)
	activeNodes = s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}
