// +build networktest

/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package servicenetwork

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/insolar/insolar/consensus/phases"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	consensusMin    = 5 // minimum count of participants that can survive when one node leaves
	consensusMinMsg = fmt.Sprintf("skip test for bootstrap nodes < %d", consensusMin)
)

func (s *testSuite) TestNetworkConsensus3Times() {
	s.waitForConsensus(3)
}

func (s *testSuite) TestNodeConnect() {
	testNode := newNetworkNode()
	s.preInitNode(testNode)

	s.InitNode(testNode)
	s.StartNode(testNode)
	defer func(s *testSuite) {
		s.StopNode(testNode)
	}(s)

	s.waitForConsensus(1)

	activeNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))

	s.waitForConsensus(1)

	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))

	s.waitForConsensus(2)

	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount()+1, len(activeNodes))
	activeNodes = testNode.serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount()+1, len(activeNodes))
}

func (s *testSuite) TestNodeConnectInvalidVersion() {
	testNode := newNetworkNode()
	s.preInitNode(testNode)
	testNode.serviceNetwork.NodeKeeper.GetOrigin().(nodenetwork.MutableNode).SetVersion("ololo")
	s.InitNode(testNode)
	err := testNode.componentManager.Start(s.fixture().ctx)
	assert.Error(s.T(), err)
	log.Infof("Error: %s", err)
}

func (s *testSuite) TestManyNodesConnect() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	nodesCount := 10
	nodes := make([]*networkNode, 0)
	for i := 0; i < nodesCount; i++ {
		nodes = append(nodes, newNetworkNode())
		s.preInitNode(nodes[i])
		s.InitNode(nodes[i])
	}

	wg := sync.WaitGroup{}
	wg.Add(nodesCount)

	for _, node := range nodes {
		go func(wg *sync.WaitGroup, node *networkNode) {
			s.StartNode(node)
			wg.Done()
		}(&wg, node)
	}

	wg.Wait()

	defer func(s *testSuite) {
		for _, node := range nodes {
			s.StopNode(node)
		}
	}(s)

	s.waitForConsensus(5)

	joined := nodesCount
	if s.getMaxJoinCount() < nodesCount {
		joined = s.getMaxJoinCount()
	}
	activeNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount()+joined, len(activeNodes))
}

func (s *testSuite) TestNodeLeave() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	testNode := newNetworkNode()
	s.preInitNode(testNode)

	s.InitNode(testNode)
	s.StartNode(testNode)
	defer func(s *testSuite) {
		s.StopNode(testNode)
	}(s)

	s.waitForConsensus(2)

	activeNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))

	s.waitForConsensus(1)

	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount()+1, len(activeNodes))

	testNode.serviceNetwork.Leave(context.Background(), 0)

	s.waitForConsensus(2)

	// one active node becomes "not working"
	workingNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(workingNodes))

	// but all nodes are active
	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount()+1, len(activeNodes))
}

func (s *testSuite) TestNodeLeaveAtETA() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	testNode := newNetworkNode()
	s.preInitNode(testNode)

	s.InitNode(testNode)
	s.StartNode(testNode)
	defer func(s *testSuite) {
		s.StopNode(testNode)
	}(s)

	// wait for node will be added at active list
	s.waitForConsensus(2)
	activeNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount()+1, len(activeNodes))

	pulse, err := s.fixture().bootstrapNodes[0].serviceNetwork.PulseStorage.Current(s.fixture().ctx)
	s.NoError(err)

	// next pulse will be last for this node
	testNode.serviceNetwork.Leave(s.fixture().ctx, pulse.NextPulseNumber)

	// node still active and working
	s.waitForConsensus(1)
	workingNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetWorkingNodes()
	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount()+1, len(activeNodes))
	s.Equal(s.getNodesCount()+1, len(workingNodes))

	// now node leaves, but it's still in active list
	s.waitForConsensus(1)
	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	workingNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(workingNodes))
	s.Equal(s.getNodesCount()+1, len(activeNodes))
}

func (s *testSuite) TestNodeComeAfterAnotherNodeSendLeaveETA() {
	s.T().Skip("fix testcase in TESTNET 2.0")
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	leavingNode := newNetworkNode()
	s.preInitNode(leavingNode)

	s.InitNode(leavingNode)
	s.StartNode(leavingNode)
	defer func(s *testSuite) {
		s.StopNode(leavingNode)
	}(s)

	// wait for node will be added at active list
	s.waitForConsensus(2)
	activeNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount()+1, len(activeNodes))

	pulse, err := s.fixture().bootstrapNodes[0].serviceNetwork.PulseStorage.Current(s.fixture().ctx)
	s.NoError(err)

	// leaving in 3 pulses
	pulseDelta := pulse.NextPulseNumber - pulse.PulseNumber
	leavingNode.serviceNetwork.Leave(s.fixture().ctx, pulse.PulseNumber+3*pulseDelta)

	// wait for leavingNode will be marked as leaving
	s.waitForConsensus(1)

	newNode := newNetworkNode()
	s.preInitNode(newNode)

	s.InitNode(newNode)
	s.StartNode(newNode)
	defer func(s *testSuite) {
		s.StopNode(newNode)
	}(s)

	// wait for newNode will be added at active list, its a last pulse for leavingNode
	s.waitForConsensus(2)

	// newNode doesn't have workingNodes
	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	workingNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetWorkingNodes()
	newNodeWorkingNodes := newNode.serviceNetwork.NodeKeeper.GetWorkingNodes()

	s.Equal(s.getNodesCount()+2, len(activeNodes))
	s.Equal(s.getNodesCount()+1, len(workingNodes))
	s.Equal(0, len(newNodeWorkingNodes))

	// newNode have to have same working node list as other nodes, but it doesn't because it miss leaving claim
	s.waitForConsensus(1)
	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	workingNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetWorkingNodes()
	newNodeWorkingNodes = newNode.serviceNetwork.NodeKeeper.GetWorkingNodes()

	s.Equal(s.getNodesCount()+2, len(activeNodes))
	s.Equal(s.getNodesCount()+1, len(workingNodes))
	// TODO: fix this testcase
	s.Equal(s.getNodesCount()+1, len(newNodeWorkingNodes))

	// leaveNode leaving, newNode still ok
	s.waitForConsensus(1)
	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	workingNodes = newNode.serviceNetwork.NodeKeeper.GetWorkingNodes()
	newNodeWorkingNodes = newNode.serviceNetwork.NodeKeeper.GetWorkingNodes()

	s.Equal(s.getNodesCount()+1, len(activeNodes))
	s.Equal(s.getNodesCount()+1, len(workingNodes))
	s.Equal(s.getNodesCount()+1, len(newNodeWorkingNodes))
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

func (ftpm *FullTimeoutPhaseManager) OnPulse(ctx context.Context, pulse *core.Pulse, pulseStartTime time.Time) error {
	return nil
}

func (s *testSuite) TestFullTimeOut() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	// TODO: make this set operation thread-safe somehow (race detector does not like this code)
	wrapper := s.fixture().bootstrapNodes[1].serviceNetwork.PhaseManager.(*phaseManagerWrapper)
	wrapper.original = &FullTimeoutPhaseManager{}
	s.fixture().bootstrapNodes[1].serviceNetwork.PhaseManager = wrapper

	s.waitForConsensus(2)

	activeNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount()-1, len(activeNodes))
}

// Partial timeout

func (s *testSuite) TestPartialPositive1PhaseTimeOut() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	setCommunicatorMock(s.fixture().bootstrapNodes, PartialPositive1Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *testSuite) TestPartialPositive2PhaseTimeOut() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	setCommunicatorMock(s.fixture().bootstrapNodes, PartialPositive2Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *testSuite) TestPartialNegative1PhaseTimeOut() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	setCommunicatorMock(s.fixture().bootstrapNodes, PartialNegative1Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount()-1, len(activeNodes))
}

func (s *testSuite) TestPartialNegative2PhaseTimeOut() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	setCommunicatorMock(s.fixture().bootstrapNodes, PartialNegative2Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *testSuite) TestPartialNegative3PhaseTimeOut() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	setCommunicatorMock(s.fixture().bootstrapNodes, PartialNegative3Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *testSuite) TestPartialPositive3PhaseTimeOut() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	setCommunicatorMock(s.fixture().bootstrapNodes, PartialPositive3Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *testSuite) TestPartialNegative23PhaseTimeOut() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	setCommunicatorMock(s.fixture().bootstrapNodes, PartialNegative23Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *testSuite) TestPartialPositive23PhaseTimeOut() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	setCommunicatorMock(s.fixture().bootstrapNodes, PartialPositive23Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *testSuite) TestDiscoveryDown() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	s.StopNode(s.fixture().bootstrapNodes[0])

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount()-1, len(activeNodes))
}

func (s *testSuite) TestDiscoveryRestart() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	s.waitForConsensus(2)

	log.Info("Discovery node stopping...")
	err := s.fixture().bootstrapNodes[0].serviceNetwork.Stop(context.Background())
	s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.(*nodeKeeperWrapper).Wipe(true)
	log.Info("Discovery node stopped...")
	require.NoError(s.T(), err)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount()-1, len(activeNodes))

	log.Info("Discovery node starting...")
	err = s.fixture().bootstrapNodes[0].serviceNetwork.Start(context.Background())
	log.Info("Discovery node started")
	require.NoError(s.T(), err)

	s.waitForConsensusExcept(3, s.fixture().bootstrapNodes[0].id)
	activeNodes = s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *testSuite) TestDiscoveryRestartNoWait() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	s.waitForConsensus(2)

	log.Info("Discovery node stopping...")
	err := s.fixture().bootstrapNodes[0].serviceNetwork.Stop(context.Background())
	s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.(*nodeKeeperWrapper).Wipe(true)
	log.Info("Discovery node stopped...")
	require.NoError(s.T(), err)

	go func(s *testSuite) {
		log.Info("Discovery node starting...")
		err = s.fixture().bootstrapNodes[0].serviceNetwork.Start(context.Background())
		log.Info("Discovery node started")
		require.NoError(s.T(), err)
	}(s)

	s.waitForConsensusExcept(4, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
	s.waitForConsensusExcept(1, s.fixture().bootstrapNodes[0].id)
	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
	activeNodes = s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func setCommunicatorMock(nodes []*networkNode, opt CommunicatorTestOpt) {
	ref := nodes[0].id
	timedOutNodesCount := 0
	switch opt {
	case PartialNegative1Phase, PartialNegative2Phase, PartialNegative3Phase, PartialNegative23Phase:
		timedOutNodesCount = int(float64(len(nodes)) * 0.6)
	case PartialPositive1Phase, PartialPositive2Phase, PartialPositive3Phase, PartialPositive23Phase:
		timedOutNodesCount = int(float64(len(nodes)) * 0.2)
	}
	// TODO: make these set operations thread-safe somehow (race detector does not like this code)
	for i := 1; i <= timedOutNodesCount; i++ {
		comm := nodes[i].serviceNetwork.PhaseManager.(*phaseManagerWrapper).original.(*phases.Phases).FirstPhase.(*phases.FirstPhaseImpl).Communicator
		wrapper := &CommunicatorMock{communicator: comm, ignoreFrom: ref, testOpt: opt}
		phasemanager := nodes[i].serviceNetwork.PhaseManager.(*phaseManagerWrapper).original.(*phases.Phases)
		phasemanager.FirstPhase.(*phases.FirstPhaseImpl).Communicator = wrapper
		phasemanager.SecondPhase.(*phases.SecondPhaseImpl).Communicator = wrapper
		phasemanager.ThirdPhase.(*phases.ThirdPhaseImpl).Communicator = wrapper
	}
}
