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

package tests

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus/claimhandler"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/network/nodenetwork"
)

var (
	consensusMin    = 5 // minimum count of participants that can survive when one node leaves
	consensusMinMsg = fmt.Sprintf("skip test for bootstrap nodes < %d", consensusMin)
)

func TestServiceNetworkManyBootstraps(t *testing.T) {
	s := newConsensusSuite(12, 0)
	suite.Run(t, s)
}

func TestServiceNetworkManyNodes(t *testing.T) {
	t.Skip("Long time setup, wait for mock bootstrap")
	s := newConsensusSuite(5, 10)
	suite.Run(t, s)
}

// Consensus suite tests

func (s *consensusSuite) TestNetworkConsensus3Times() {
	s.waitForConsensus(3)
}

func (s *consensusSuite) TestNodeConnect() {
	testNode := s.newNetworkNode("testNode")
	s.preInitNode(testNode)

	s.InitNode(testNode)
	s.StartNode(testNode)
	defer func(s *consensusSuite) {
		s.StopNode(testNode)
	}(s)

	s.waitForConsensus(1)

	s.AssertActiveNodesCountDelta(0)

	s.waitForConsensus(1)

	s.AssertActiveNodesCountDelta(1)
	s.AssertWorkingNodesCountDelta(0)

	s.waitForConsensus(2)

	s.AssertActiveNodesCountDelta(1)
	s.AssertWorkingNodesCountDelta(1)
}

func (s *consensusSuite) TestNodeConnectInvalidVersion() {
	testNode := s.newNetworkNode("testNode")
	s.preInitNode(testNode)
	testNode.serviceNetwork.NodeKeeper.GetOrigin().(node.MutableNode).SetVersion("ololo")
	s.InitNode(testNode)
	err := testNode.componentManager.Start(s.fixture().ctx)
	assert.Error(s.T(), err)
	log.Infof("Error: %s", err)
}

func (s *consensusSuite) TestManyNodesConnect() {
	s.T().Skip("test hangs in some situations, needs fix: INS-2200")

	s.CheckBootstrapCount()

	joinersCount := 10
	nodes := make([]*networkNode, 0)
	for i := 0; i < joinersCount; i++ {
		n := s.newNetworkNode(fmt.Sprintf("testNode_%d", i))
		nodes = append(nodes, n)
	}

	wg := sync.WaitGroup{}
	wg.Add(joinersCount)

	for _, n := range nodes {
		go func(wg *sync.WaitGroup, node *networkNode) {
			s.preInitNode(node)
			s.InitNode(node)
			s.StartNode(node)
			wg.Done()
		}(&wg, n)
	}

	wg.Wait()

	defer func() {
		for _, n := range nodes {
			s.StopNode(n)
		}
	}()

	s.waitForConsensus(5)

	joined := claimhandler.ApprovedJoinersCount(joinersCount, s.getNodesCount())
	activeNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetAccessor().GetActiveNodes()
	s.Equal(s.getNodesCount()+joined, len(activeNodes))
}

func (s *consensusSuite) TestNodeLeave() {
	s.CheckBootstrapCount()

	testNode := s.newNetworkNode("testNode")
	s.preInitNode(testNode)

	s.InitNode(testNode)
	s.StartNode(testNode)
	defer func(s *consensusSuite) {
		s.StopNode(testNode)
	}(s)

	s.waitForConsensus(2)
	s.AssertActiveNodesCountDelta(1)
	s.AssertWorkingNodesCountDelta(0)

	// node become working after 3 consensuses
	s.waitForConsensus(1)
	s.AssertActiveNodesCountDelta(1)
	s.AssertWorkingNodesCountDelta(1)

	testNode.serviceNetwork.Leave(context.Background(), 0)

	s.waitForConsensus(2)

	// one active node becomes "not working"
	s.AssertWorkingNodesCountDelta(0)

	// but all nodes are active
	s.AssertActiveNodesCountDelta(1)
}

func (s *consensusSuite) TestNodeLeaveAtETA() {
	s.CheckBootstrapCount()

	testNode := s.newNetworkNode("testNode")
	s.preInitNode(testNode)

	s.InitNode(testNode)
	testNode.terminationHandler.OnLeaveApprovedFinished()
	testNode.terminationHandler.OnLeaveApprovedFunc = func(p context.Context) {
		s.StopNode(testNode)
	}
	s.StartNode(testNode)
	defer func(s *consensusSuite) {
		s.StopNode(testNode)
	}(s)

	// wait for node will be added at active and working lists
	s.waitForConsensus(3)
	s.AssertActiveNodesCountDelta(1)
	s.AssertWorkingNodesCountDelta(1)

	pulse, err := s.fixture().bootstrapNodes[0].serviceNetwork.PulseAccessor.Latest(s.fixture().ctx)
	s.NoError(err)

	// next pulse will be last for this node
	// leaving in 3 pulses
	pulseDelta := pulse.NextPulseNumber - pulse.PulseNumber
	testNode.serviceNetwork.Leave(s.fixture().ctx, pulse.PulseNumber+3*pulseDelta)

	// node still active and working
	s.waitForConsensus(1)
	s.AssertActiveNodesCountDelta(1)
	s.AssertWorkingNodesCountDelta(1)

	// now node leaves, but it's still in active list
	s.waitForConsensus(1)
	s.AssertActiveNodesCountDelta(1)
	s.AssertWorkingNodesCountDelta(0)
	s.True(testNode.terminationHandler.OnLeaveApprovedFinished())

	s.waitForConsensus(1)
	s.AssertActiveNodesCountDelta(0)
	s.AssertWorkingNodesCountDelta(0)
}

func (s *consensusSuite) TestNodeComeAfterAnotherNodeSendLeaveETA() {
	s.T().Skip("fix testcase in TESTNET 2.0")
	if len(s.fixture().bootstrapNodes) < consensusMin {
		s.T().Skip(consensusMinMsg)
	}

	leavingNode := s.newNetworkNode("leavingNode")
	s.preInitNode(leavingNode)

	s.InitNode(leavingNode)
	s.StartNode(leavingNode)
	defer func(s *consensusSuite) {
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
	defer func(s *consensusSuite) {
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

func (s *consensusSuite) TestFullTimeOut() {
	s.CheckBootstrapCount()

	s.SetCommunicationPolicy(FullTimeout)

	s.waitForConsensus(2)
	s.AssertWorkingNodesCountDelta(-1)
}

// Partial timeout

func (s *consensusSuite) TestPartialPositive1PhaseTimeOut() {
	s.CheckBootstrapCount()

	s.SetCommunicationPolicy(PartialPositive1Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *consensusSuite) TestPartialPositive2PhaseTimeOut() {
	s.CheckBootstrapCount()

	s.SetCommunicationPolicy(PartialPositive2Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *consensusSuite) TestPartialNegative1PhaseTimeOut() {
	s.CheckBootstrapCount()

	s.SetCommunicationPolicy(PartialNegative1Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount()-1, len(activeNodes))
}

func (s *consensusSuite) TestPartialNegative2PhaseTimeOut() {
	s.CheckBootstrapCount()

	s.SetCommunicationPolicy(PartialNegative2Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *consensusSuite) TestPartialNegative3PhaseTimeOut() {
	s.CheckBootstrapCount()

	s.SetCommunicationPolicy(PartialNegative3Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *consensusSuite) TestPartialPositive3PhaseTimeOut() {
	s.CheckBootstrapCount()

	s.SetCommunicationPolicy(PartialPositive3Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *consensusSuite) TestPartialNegative23PhaseTimeOut() {
	s.CheckBootstrapCount()

	s.SetCommunicationPolicy(PartialNegative23Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *consensusSuite) TestPartialPositive23PhaseTimeOut() {
	s.CheckBootstrapCount()

	s.SetCommunicationPolicy(PartialPositive23Phase)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].serviceNetwork.NodeKeeper.GetWorkingNodes()
	s.Equal(s.getNodesCount(), len(activeNodes))
}

func (s *consensusSuite) TestDiscoveryDown() {
	s.CheckBootstrapCount()

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

func (s *consensusSuite) TestDiscoveryRestart() {
	s.CheckBootstrapCount()

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

func (s *consensusSuite) TestDiscoveryRestartNoWait() {
	s.CheckBootstrapCount()

	s.waitForConsensus(2)

	log.Info("Discovery node stopping...")
	err := s.fixture().bootstrapNodes[0].serviceNetwork.Stop(context.Background())
	flushNodeKeeper(s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper)
	log.Info("Discovery node stopped...")
	s.Require().NoError(err)

	go func(s *consensusSuite) {
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

func (s *consensusSuite) TestJoinerSplitPackets() {
	s.CheckBootstrapCount()

	testNode := s.newNetworkNode("testNode")
	s.SetCommunicationPolicyForNode(testNode.id, SplitCase)
	s.preInitNode(testNode)

	s.InitNode(testNode)
	s.StartNode(testNode)
	defer func(s *consensusSuite) {
		s.StopNode(testNode)
	}(s)

	s.waitForConsensus(1)

	s.AssertActiveNodesCountDelta(0)

	s.waitForConsensus(1)

	s.AssertActiveNodesCountDelta(1)
	s.AssertWorkingNodesCountDelta(0)

	s.waitForConsensus(2)

	s.AssertActiveNodesCountDelta(1)
	s.AssertWorkingNodesCountDelta(1)
}
