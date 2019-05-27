package tests

import "github.com/insolar/insolar/insolar"

func (s *testSuite) AssertNodesIsActive(target *networkNode, nodes []insolar.NetworkNode) {
	// allNodes := append(s.fixture().bootstrapNodes, s.fixture().networkNodes...)
	// activeNodes := target.serviceNetwork.NodeKeeper.GetAccessor().GetActiveNodes()
	// activeNodes

}

func (s *testSuite) AssertNodesIsWorking(target *networkNode, nodes []insolar.NetworkNode) {

}

func (s *testSuite) AssertActiveNodesCountDelta(delta int) {
	activeNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetAccessor().GetActiveNodes()
	s.Equal(s.getNodesCount()+delta, len(activeNodes))
}

func (s *testSuite) AssertWorkingNodesCountDelta(delta int) {
	workingNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetAccessor().GetWorkingNodes()
	s.Equal(s.getNodesCount()+delta, len(workingNodes))
}
