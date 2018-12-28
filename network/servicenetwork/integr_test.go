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

	s := NewTestSuite(5, 10)
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

	// TODO: make this set operation thread-safe somehow (race detector does not like this code)
	wrapper := s.fixture().bootstrapNodes[1].serviceNetwork.PhaseManager.(*phaseManagerWrapper)
	wrapper.original = &FullTimeoutPhaseManager{}
	s.fixture().bootstrapNodes[1].serviceNetwork.PhaseManager = wrapper

	s.waitForConsensus(2)

	activeNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount()-1, len(activeNodes))
}

// Partial timeout

func (s *testSuite) TestPartialPositive1PhaseTimeOut() {
	if len(s.fixture().bootstrapNodes) < 3 {
		s.T().Skip("skip test for bootstrap nodes < 3")
	}

	setCommunicatorMock(s.fixture().bootstrapNodes, PartialPositive1Phase)

	s.waitForConsensus(2)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *testSuite) TestPartialPositive2PhaseTimeOut() {
	if len(s.fixture().bootstrapNodes) < 3 {
		s.T().Skip("skip test for bootstrap nodes < 3")
	}

	setCommunicatorMock(s.fixture().bootstrapNodes, PartialPositive2Phase)

	s.waitForConsensus(2)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *testSuite) TestPartialNegative1PhaseTimeOut() {
	if len(s.fixture().bootstrapNodes) < 3 {
		s.T().Skip("skip test for bootstrap nodes < 3")
	}

	setCommunicatorMock(s.fixture().bootstrapNodes, PartialNegative1Phase)

	s.waitForConsensus(2)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *testSuite) TestPartialNegative2PhaseTimeOut() {
	if len(s.fixture().bootstrapNodes) < 3 {
		s.T().Skip("skip test for bootstrap nodes < 3")
	}

	setCommunicatorMock(s.fixture().bootstrapNodes, PartialNegative2Phase)

	s.waitForConsensus(2)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func setCommunicatorMock(nodes []*networkNode, opt CommunicatorTestOpt) {
	ref := nodes[0].id
	timedOutNodesCount := 0
	switch opt {
	case PartialNegative1Phase:
		fallthrough
	case PartialNegative2Phase:
		timedOutNodesCount = int(float64(len(nodes)) * 0.6)
	case PartialPositive1Phase:
		fallthrough
	case PartialPositive2Phase:
		timedOutNodesCount = int(float64(len(nodes)) * 0.2)
	}
	// TODO: make these set operations thread-safe somehow (race detector does not like this code)
	for i := 1; i <= timedOutNodesCount; i++ {
		comm := nodes[i].serviceNetwork.PhaseManager.(*phaseManagerWrapper).original.(*phases.Phases).FirstPhase.Communicator
		wrapper := &CommunicatorMock{communicator: comm, ignoreFrom: ref, testOpt: opt}
		nodes[i].serviceNetwork.PhaseManager.(*phaseManagerWrapper).original.(*phases.Phases).FirstPhase.Communicator = wrapper
	}
}
