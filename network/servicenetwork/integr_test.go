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

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
	ctx          context.Context
	networkNodes []networkNode
}

func NewTestSuite() *testSuite {
	return &testSuite{
		Suite:        suite.Suite{},
		ctx:          context.Background(),
		networkNodes: make([]networkNode, 0),
	}
}

func (s *testSuite) StartNodes() {
	for _, n := range s.networkNodes {
		err := n.componentManager.Init(s.ctx)
		s.NoError(err)
		err = n.componentManager.Start(s.ctx)
		s.NoError(err)
	}
}

func (s *testSuite) StopNodes() {
	for _, n := range s.networkNodes {
		err := n.componentManager.Stop(s.ctx)
		s.NoError(err)
	}
}

type networkNode struct {
	componentManager *component.Manager
	serviceNetwork   *ServiceNetwork
}

func initCrypto(t *testing.T) (*certificate.Certificate, core.CryptographyService) {
	key, _ := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	require.NotNil(t, key)
	cs := cryptography.NewKeyBoundCryptographyService(key)
	kp := platformpolicy.NewKeyProcessor()
	pk, _ := cs.GetPublicKey()
	cert, err := certificate.NewCertificatesWithKeys(pk, kp)
	require.NoError(t, err)

	return cert, cs
}

func createNetworkNode(t *testing.T) networkNode {
	address := "127.0.0.1:0"
	consensusAddr := "127.0.0.1:0"

	origin := nodenetwork.NewNode(testutils.RandomRef(), core.StaticRoleUnknown, nil, 0, address, "")
	keeper := nodenetwork.NewNodeKeeper(origin)

	cfg := configuration.NewConfiguration()
	cfg.Node.Node.ID = origin.ID().String()
	cfg.Host.Transport.Address = address
	cfg.Host.ConsensusTransport.Address = consensusAddr

	//cfg.Host.BootstrapHosts = append(cfg.Host.BootstrapHosts, "127.0.0.1:0")

	scheme := platformpolicy.NewPlatformCryptographyScheme()
	serviceNetwork, err := NewServiceNetwork(cfg, scheme)
	assert.NoError(t, err)

	pulseManagerMock := testutils.NewPulseManagerMock(t)
	netCoordinator := testutils.NewNetworkCoordinatorMock(t)
	amMock := testutils.NewArtifactManagerMock(t)

	cm := &component.Manager{}
	cm.Register(keeper, pulseManagerMock, netCoordinator, amMock)
	cm.Register(initCrypto(t))
	cm.Inject(serviceNetwork)

	serviceNetwork.NodeKeeper = keeper

	return networkNode{cm, serviceNetwork}
}

func (s *testSuite) TestSendConsensusPhase() {
	s.StartNodes()
	s.StopNodes()
	//activeNodes := s.networkNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	//s.Equal(1, len(activeNodes))
}

func TestNewServiceNetwork2(t *testing.T) {
	s := NewTestSuite()
	node1 := createNetworkNode(t)
	node2 := createNetworkNode(t)
	node3 := createNetworkNode(t)

	s.networkNodes = append(s.networkNodes, node1, node2, node3)

	suite.Run(t, s)

}
