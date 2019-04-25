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

// +build networktest

package servicenetwork

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/consensus/claimhandler"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/network/nodenetwork"
)

var (
	consensusMin    = 5 // minimum count of participants that can survive when one node leaves
	consensusMinMsg = fmt.Sprintf("skip test for bootstrap nodes < %d", consensusMin)
)

func (s *testSuite) TestNetworkConsensus3Times() {
	s.waitForConsensus(3)
}

func (s *testSuite) TestNodeConnect() {
	testNode := s.newNetworkNode("testNode")
	s.preInitNode(testNode)

	s.InitNode(testNode)
	s.StartNode(testNode)
	defer func(s *testSuite) {
		s.StopNode(testNode)
	}(s)

	s.waitForConsensus(1)

	activeNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetAccessor().GetActiveNodes()
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
	testNode := s.newNetworkNode("testNode")
	s.preInitNode(testNode)
	testNode.serviceNetwork.NodeKeeper.GetOrigin().(node.MutableNode).SetVersion("ololo")
	s.InitNode(testNode)
	err := testNode.componentManager.Start(s.fixture().ctx)
	assert.Error(s.T(), err)
	log.Infof("Error: %s", err)
}

func (s *testSuite) TestManyNodesConnect() {
	s.T().Skip("test hangs in some situations, needs fix: INS-2200")
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	joinersCount := 10
	nodes := make([]*networkNode, 0)
	for i := 0; i < joinersCount; i++ {
		node := s.newNetworkNode(fmt.Sprintf("testNode_%d", i))
		nodes = append(nodes, node)
	}

	wg := sync.WaitGroup{}
	wg.Add(joinersCount)

	for _, node := range nodes {
		go func(wg *sync.WaitGroup, node *networkNode) {
			s.preInitNode(node)
			s.InitNode(node)
			s.StartNode(node)
			wg.Done()
		}(&wg, node)
	}

	wg.Wait()

	defer func() {
		for _, node := range nodes {
			s.StopNode(node)
		}
	}()

	s.waitForConsensus(5)

	joined := claimhandler.ApprovedJoinersCount(joinersCount, s.getNodesCount())
	activeNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetAccessor().GetActiveNodes()
	s.Equal(s.getNodesCount()+joined, len(activeNodes))
}

func (s *testSuite) TestNodeLeave() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	testNode := s.newNetworkNode("testNode")
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
	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetAccessor().GetActiveNodes()
	s.Equal(s.getNodesCount()+1, len(activeNodes))
}

func (s *testSuite) TestNodeLeaveAtETA() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	testNode := s.newNetworkNode("testNode")
	s.preInitNode(testNode)

	s.InitNode(testNode)
	s.StartNode(testNode)
	defer func(s *testSuite) {
		s.StopNode(testNode)
	}(s)

	// wait for node will be added at active list
	s.waitForConsensus(2)
	activeNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetAccessor().GetActiveNodes()
	s.Equal(s.getNodesCount()+1, len(activeNodes))

	pulse, err := s.fixture().bootstrapNodes[0].serviceNetwork.PulseAccessor.Latest(s.fixture().ctx)
	s.NoError(err)

	// next pulse will be last for this node
	testNode.serviceNetwork.Leave(s.fixture().ctx, pulse.NextPulseNumber)

	// node still active and working
	s.waitForConsensus(1)
	workingNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetWorkingNodes()
	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetAccessor().GetActiveNodes()
	s.Equal(s.getNodesCount()+1, len(activeNodes))
	s.Equal(s.getNodesCount()+1, len(workingNodes))

	// now node leaves, but it's still in active list
	s.waitForConsensus(1)
	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetAccessor().GetActiveNodes()
	workingNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(workingNodes))
	s.Equal(s.getNodesCount()+1, len(activeNodes))
}

func (s *testSuite) TestNodeComeAfterAnotherNodeSendLeaveETA() {
	s.T().Skip("fix testcase in TESTNET 2.0")
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	leavingNode := s.newNetworkNode("leavingNode")
	s.preInitNode(leavingNode)

	s.InitNode(leavingNode)
	s.StartNode(leavingNode)
	defer func(s *testSuite) {
		s.StopNode(leavingNode)
	}(s)

	// wait for node will be added at active list
	s.waitForConsensus(2)
	activeNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetAccessor().GetActiveNodes()
	s.Equal(s.getNodesCount()+1, len(activeNodes))

	pulse, err := s.fixture().bootstrapNodes[0].serviceNetwork.PulseAccessor.Latest(s.fixture().ctx)
	s.NoError(err)

	// leaving in 3 pulses
	pulseDelta := pulse.NextPulseNumber - pulse.PulseNumber
	leavingNode.serviceNetwork.Leave(s.fixture().ctx, pulse.PulseNumber+3*pulseDelta)

	// wait for leavingNode will be marked as leaving
	s.waitForConsensus(1)

	newNode := s.newNetworkNode("testNode")
	s.preInitNode(newNode)

	s.InitNode(newNode)
	s.StartNode(newNode)
	defer func(s *testSuite) {
		s.StopNode(newNode)
	}(s)

	// wait for newNode will be added at active list, its a last pulse for leavingNode
	s.waitForConsensus(2)

	// newNode doesn't have workingNodes
	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetAccessor().GetActiveNodes()
	workingNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetWorkingNodes()
	newNodeWorkingNodes := newNode.serviceNetwork.NodeKeeper.GetWorkingNodes()

	s.Equal(s.getNodesCount()+2, len(activeNodes))
	s.Equal(s.getNodesCount()+1, len(workingNodes))
	s.Equal(0, len(newNodeWorkingNodes))

	// newNode have to have same working node list as other nodes, but it doesn't because it miss leaving claim
	s.waitForConsensus(1)
	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetAccessor().GetActiveNodes()
	workingNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetWorkingNodes()
	newNodeWorkingNodes = newNode.serviceNetwork.NodeKeeper.GetWorkingNodes()

	s.Equal(s.getNodesCount()+2, len(activeNodes))
	s.Equal(s.getNodesCount()+1, len(workingNodes))
	// TODO: fix this testcase
	s.Equal(s.getNodesCount()+1, len(newNodeWorkingNodes))

	// leaveNode leaving, newNode still ok
	s.waitForConsensus(1)
	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetAccessor().GetActiveNodes()
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

func (s *testSuite) TestFullTimeOut() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	s.SetCommunicationPolicy(FullTimeout)

	s.waitForConsensus(2)
	activeNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount()-1, len(activeNodes))
}

// Partial timeout

func (s *testSuite) TestPartialPositive1PhaseTimeOut() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	s.SetCommunicationPolicy(PartialPositive1Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *testSuite) TestPartialPositive2PhaseTimeOut() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	s.SetCommunicationPolicy(PartialPositive2Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *testSuite) TestPartialNegative1PhaseTimeOut() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	s.SetCommunicationPolicy(PartialNegative1Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount()-1, len(activeNodes))
}

func (s *testSuite) TestPartialNegative2PhaseTimeOut() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	s.SetCommunicationPolicy(PartialNegative2Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *testSuite) TestPartialNegative3PhaseTimeOut() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	s.SetCommunicationPolicy(PartialNegative3Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *testSuite) TestPartialPositive3PhaseTimeOut() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	s.SetCommunicationPolicy(PartialPositive3Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *testSuite) TestPartialNegative23PhaseTimeOut() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	s.SetCommunicationPolicy(PartialNegative23Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *testSuite) TestPartialPositive23PhaseTimeOut() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	s.SetCommunicationPolicy(PartialPositive23Phase)

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
	for i := 1; i < s.getNodesCount(); i++ {
		activeNodes := s.fixture().bootstrapNodes[i].serviceNetwork.NodeKeeper.GetWorkingNodes()
		s.Equal(s.getNodesCount()-1, len(activeNodes))
	}
}

func flushNodeKeeper(keeper network.NodeKeeper) {
	keeper.SetIsBootstrapped(false)
	keeper.GetConsensusInfo().(*nodenetwork.ConsensusInfo).Flush(false)
	keeper.SetCloudHash(nil)
	keeper.SetInitialSnapshot([]insolar.NetworkNode{})
	keeper.GetClaimQueue().Clear()
	keeper.GetOrigin().(node.MutableNode).SetState(insolar.NodeReady)
}

func (s *testSuite) TestDiscoveryRestart() {
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	s.waitForConsensus(2)

	log.Info("Discovery node stopping...")
	err := s.fixture().bootstrapNodes[0].serviceNetwork.Stop(context.Background())
	flushNodeKeeper(s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper)
	log.Info("Discovery node stopped...")
	s.Require().NoError(err)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount()-1, len(activeNodes))

	log.Info("Discovery node starting...")
	err = s.fixture().bootstrapNodes[0].serviceNetwork.Start(context.Background())
	log.Info("Discovery node started")
	s.Require().NoError(err)

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
	flushNodeKeeper(s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper)
	log.Info("Discovery node stopped...")
	s.Require().NoError(err)

	go func(s *testSuite) {
		log.Info("Discovery node starting...")
		err = s.fixture().bootstrapNodes[0].serviceNetwork.Start(context.Background())
		log.Info("Discovery node started")
		s.Require().NoError(err)
	}(s)

	s.waitForConsensusExcept(4, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetAccessor().GetActiveNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
	s.waitForConsensusExcept(1, s.fixture().bootstrapNodes[0].id)
	activeNodes = s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
	activeNodes = s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}
