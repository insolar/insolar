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
	"github.com/fortytw2/leaktest"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/node"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	consensusMin    = 5 // minimum count of participants that can survive when one node leaves
	consensusMinMsg = fmt.Sprintf("skip test for bootstrap nodes < %d", consensusMin)
)

func serviceNetworkManyBootstraps(t *testing.T) *consensusSuite {
	cs := newConsensusSuite(t, 5, 0)
	cs.SetupTest()

	return cs
}

// Consensus suite tests

func TestNetworkConsensusManyTimes(t *testing.T) {
	defer leaktest.Check(t)()

	s := serviceNetworkManyBootstraps(t)
	defer s.TearDownTest()

	s.waitForConsensus(9)
	s.AssertActiveNodesCountDelta(0)
}

func TestJoinerNodeConnect(t *testing.T) {
	defer leaktest.Check(t)()

	s := serviceNetworkManyBootstraps(t)
	defer s.TearDownTest()

	testNode := s.newNetworkNode("JoinerNode")
	s.preInitNode(testNode)

	s.InitNode(testNode)
	s.StartNode(testNode)
	defer func(s *consensusSuite) {
		s.StopNode(testNode)
	}(s)

	s.waitForConsensus(1)

	s.AssertActiveNodesCountDelta(0)

	s.waitForConsensus(2)

	s.AssertActiveNodesCountDelta(1)
	require.Equal(s.t, insolar.CompleteNetworkState, testNode.serviceNetwork.Gatewayer.Gateway().GetState())
}

func TestNodeConnectInvalidVersion(t *testing.T) {
	t.Skip("Behavior changed, fix assertion in test needed")
	s := serviceNetworkManyBootstraps(t)
	defer s.TearDownTest()

	testNode := s.newNetworkNode("testNode")
	s.preInitNode(testNode)
	testNode.serviceNetwork.NodeKeeper.GetOrigin().(node.MutableNode).SetVersion("ololo")
	require.Equal(t, "ololo", testNode.serviceNetwork.NodeKeeper.GetOrigin().Version())

	s.InitNode(testNode)
	err := testNode.componentManager.Start(s.fixture().ctx)
	defer func(s *consensusSuite) {
		s.StopNode(testNode)
	}(s)
	assert.NoError(t, err)

	s.waitForConsensus(2)
	s.AssertActiveNodesCountDelta(0)

}

func TestNodeLeave(t *testing.T) {
	defer leaktest.Check(t)()

	s := serviceNetworkManyBootstraps(t)
	defer s.TearDownTest()

	s.CheckBootstrapCount()

	testNode := s.newNetworkNode("testNode")
	s.preInitNode(testNode)

	s.InitNode(testNode)
	s.StartNode(testNode)

	s.waitForConsensus(2)

	s.AssertActiveNodesCountDelta(1)
	s.AssertWorkingNodesCountDelta(1)
	require.Equal(s.t, insolar.CompleteNetworkState, testNode.serviceNetwork.Gatewayer.Gateway().GetState())

	s.StopNode(testNode)

	s.waitForConsensus(3)

	s.AssertWorkingNodesCountDelta(0)
	s.AssertActiveNodesCountDelta(0)
}

func TestNodeGracefulLeave(t *testing.T) {
	defer leaktest.Check(t)()

	s := serviceNetworkManyBootstraps(t)
	defer s.TearDownTest()

	s.CheckBootstrapCount()

	testNode := s.newNetworkNode("testNode")
	s.preInitNode(testNode)

	s.InitNode(testNode)
	s.StartNode(testNode)

	s.waitForConsensus(2)

	s.AssertActiveNodesCountDelta(1)
	s.AssertWorkingNodesCountDelta(1)
	require.Equal(s.t, insolar.CompleteNetworkState, testNode.serviceNetwork.Gatewayer.Gateway().GetState())

	s.GracefulStop(testNode)

	s.waitForConsensus(3)

	s.AssertWorkingNodesCountDelta(0)
	s.AssertActiveNodesCountDelta(0)
}

func TestNodeLeaveAtETA(t *testing.T) {
	t.Skip("FIXME")
	s := serviceNetworkManyBootstraps(t)
	defer s.TearDownTest()

	s.CheckBootstrapCount()

	testNode := s.newNetworkNode("testNode")
	s.preInitNode(testNode)

	s.InitNode(testNode)
	leaveApprovedCalls := 0
	testNode.terminationHandler.OnLeaveApprovedMock.Set(func(p context.Context) {
		leaveApprovedCalls++
		s.StopNode(testNode)
	})
	s.StartNode(testNode)
	defer func(s *consensusSuite) {
		s.StopNode(testNode)
	}(s)

	// wait for node will be added at active and working lists
	s.waitForConsensus(3)
	s.AssertActiveNodesCountDelta(1)
	s.AssertWorkingNodesCountDelta(1)

	pulse, err := s.fixture().bootstrapNodes[0].serviceNetwork.PulseAccessor.GetLatestPulse(s.fixture().ctx)
	require.NoError(t, err)

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
	require.Equal(t, 1, leaveApprovedCalls)

	s.waitForConsensus(1)
	s.AssertActiveNodesCountDelta(0)
	s.AssertWorkingNodesCountDelta(0)
}

func TestNodeComeAfterAnotherNodeSendLeaveETA(t *testing.T) {
	t.Skip("FIXME")
	s := serviceNetworkManyBootstraps(t)
	defer s.TearDownTest()

	t.Skip("fix testcase in TESTNET 2.0")
	if len(s.fixture().bootstrapNodes) < consensusMin {
		t.Skip(consensusMinMsg)
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
	activeNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetAccessor(0).GetActiveNodes()
	require.Equal(t, s.getNodesCount()+1, len(activeNodes))

	pulse, err := s.fixture().bootstrapNodes[0].serviceNetwork.PulseAccessor.GetLatestPulse(s.fixture().ctx)
	require.NoError(t, err)

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
	activeNodes = s.fixture().bootstrapNodes[0].GetActiveNodes()
	workingNodes := s.fixture().bootstrapNodes[0].GetWorkingNodes()
	newNodeWorkingNodes := newNode.GetWorkingNodes()

	require.Equal(t, s.getNodesCount()+2, len(activeNodes))
	require.Equal(t, s.getNodesCount()+1, len(workingNodes))
	require.Equal(t, 0, len(newNodeWorkingNodes))

	// newNode have to have same working node list as other nodes, but it doesn't because it miss leaving claim
	s.waitForConsensus(1)
	activeNodes = s.fixture().bootstrapNodes[0].GetActiveNodes()
	workingNodes = s.fixture().bootstrapNodes[0].GetWorkingNodes()
	newNodeWorkingNodes = newNode.GetWorkingNodes()

	require.Equal(t, s.getNodesCount()+2, len(activeNodes))
	require.Equal(t, s.getNodesCount()+1, len(workingNodes))
	// TODO: fix this testcase
	require.Equal(t, s.getNodesCount()+1, len(newNodeWorkingNodes))

	// leaveNode leaving, newNode still ok
	s.waitForConsensus(1)
	activeNodes = s.fixture().bootstrapNodes[0].GetActiveNodes()
	workingNodes = newNode.GetWorkingNodes()
	newNodeWorkingNodes = newNode.GetWorkingNodes()

	require.Equal(t, s.getNodesCount()+1, len(activeNodes))
	require.Equal(t, s.getNodesCount()+1, len(workingNodes))
	require.Equal(t, s.getNodesCount()+1, len(newNodeWorkingNodes))
}

func TestDiscoveryDown(t *testing.T) {
	t.Skip("FIXME")
	s := serviceNetworkManyBootstraps(t)
	defer s.TearDownTest()

	s.CheckBootstrapCount()
	s.StopNode(s.fixture().bootstrapNodes[0])
	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	for i := 1; i < s.getNodesCount(); i++ {
		activeNodes := s.fixture().bootstrapNodes[i].GetWorkingNodes()
		require.Equal(t, s.getNodesCount()-1, len(activeNodes))
	}
}

func flushNodeKeeper(keeper network.NodeKeeper) {
	// keeper.SetIsBootstrapped(false)
	// keeper.SetCloudHash(nil)
	keeper.SetInitialSnapshot([]insolar.NetworkNode{})
	keeper.GetOrigin().(node.MutableNode).SetState(insolar.NodeReady)
}

func TestDiscoveryRestart(t *testing.T) {
	t.Skip("FIXME")
	s := serviceNetworkManyBootstraps(t)
	defer s.TearDownTest()

	s.CheckBootstrapCount()

	s.waitForConsensus(2)

	log.Info("Discovery node stopping...")
	err := s.fixture().bootstrapNodes[0].serviceNetwork.Stop(context.Background())
	flushNodeKeeper(s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper)
	log.Info("Discovery node stopped...")
	require.NoError(t, err)

	s.waitForConsensusExcept(2, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].GetWorkingNodes()
	require.Equal(t, s.getNodesCount()-1, len(activeNodes))

	log.Info("Discovery node starting...")
	err = s.fixture().bootstrapNodes[0].serviceNetwork.Start(context.Background())
	log.Info("Discovery node started")
	require.NoError(t, err)

	s.waitForConsensusExcept(3, s.fixture().bootstrapNodes[0].id)
	activeNodes = s.fixture().bootstrapNodes[1].GetWorkingNodes()
	require.Equal(t, s.getNodesCount(), len(activeNodes))
	activeNodes = s.fixture().bootstrapNodes[0].GetWorkingNodes()
	require.Equal(t, s.getNodesCount(), len(activeNodes))
}

func TestDiscoveryRestartNoWait(t *testing.T) {
	t.Skip("FIXME")
	s := serviceNetworkManyBootstraps(t)
	defer s.TearDownTest()

	s.CheckBootstrapCount()

	s.waitForConsensus(2)

	log.Info("Discovery node stopping...")
	err := s.fixture().bootstrapNodes[0].serviceNetwork.Stop(context.Background())
	flushNodeKeeper(s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper)
	log.Info("Discovery node stopped...")
	require.NoError(t, err)

	go func(s *consensusSuite) {
		log.Info("Discovery node starting...")
		err = s.fixture().bootstrapNodes[0].serviceNetwork.Start(context.Background())
		log.Info("Discovery node started")
		require.NoError(t, err)
	}(s)

	s.waitForConsensusExcept(4, s.fixture().bootstrapNodes[0].id)
	activeNodes := s.fixture().bootstrapNodes[1].GetActiveNodes()
	require.Equal(t, s.getNodesCount(), len(activeNodes))
	s.waitForConsensusExcept(1, s.fixture().bootstrapNodes[0].id)
	activeNodes = s.fixture().bootstrapNodes[0].GetWorkingNodes()
	require.Equal(t, s.getNodesCount(), len(activeNodes))
	activeNodes = s.fixture().bootstrapNodes[1].GetWorkingNodes()
	require.Equal(t, s.getNodesCount(), len(activeNodes))
}

// func (s *consensusSuite) TestJoinerSplitPackets() {
//	s.CheckBootstrapCount()
//
//	testNode := s.newNetworkNode("testNode")
//	s.SetCommunicationPolicyForNode(testNode.id, SplitCase)
//	s.preInitNode(testNode)
//
//	s.InitNode(testNode)
//	s.StartNode(testNode)
//	defer func(s *consensusSuite) {
//		s.StopNode(testNode)
//	}(s)
//
//	s.waitForConsensus(1)
//
//	s.AssertActiveNodesCountDelta(0)
//
//	s.waitForConsensus(1)
//
//	s.AssertActiveNodesCountDelta(1)
//	s.AssertWorkingNodesCountDelta(0)
//
//	s.waitForConsensus(2)
//
//	s.AssertActiveNodesCountDelta(1)
//	s.AssertWorkingNodesCountDelta(1)
// }
