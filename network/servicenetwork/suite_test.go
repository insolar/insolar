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
	"crypto"
	"strconv"
	"strings"
	"time"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/network/utils"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var testNetworkPort = 10010

type testSuite struct {
	suite.Suite
	ctx            context.Context
	bootstrapNodes []*networkNode
	networkNodes   []*networkNode
	testNode       *networkNode
}

func NewTestSuite(bootstrapCount, nodesCount int) *testSuite {
	s := &testSuite{
		Suite:          suite.Suite{},
		ctx:            context.Background(),
		bootstrapNodes: make([]*networkNode, 0),
		networkNodes:   make([]*networkNode, 0),
	}

	for i := 0; i < bootstrapCount; i++ {
		s.bootstrapNodes = append(s.bootstrapNodes, newNetworkNode())
	}

	for i := 0; i < nodesCount; i++ {
		s.networkNodes = append(s.networkNodes, newNetworkNode())
	}

	s.testNode = newNetworkNode()
	return s
}

// SetupSuite creates and run network with bootstrap and common nodes once before run all tests in the suite
func (s *testSuite) SetupSuite() {
	log.Infoln("SetupSuite")
	for _, node := range s.bootstrapNodes {
		s.initNode(node, Disable)
	}

	for _, node := range s.networkNodes {
		s.initNode(node, Disable)
	}

	log.Infoln("Init bootstrap nodes")
	for _, n := range s.bootstrapNodes {
		err := n.componentManager.Init(s.ctx)
		s.NoError(err)
	}
	log.Infoln("Init network nodes")
	for _, n := range s.networkNodes {
		err := n.componentManager.Init(s.ctx)
		s.NoError(err)
	}

	log.Infoln("Start bootstrap nodes")
	for _, n := range s.bootstrapNodes {
		err := n.componentManager.Start(s.ctx)
		s.NoError(err)
	}
	log.Infoln("Start network nodes")
	for _, n := range s.networkNodes {
		err := n.componentManager.Start(s.ctx)
		s.NoError(err)
	}

	<-time.After(time.Second * 2)
	//TODO: wait for first consensus
	// active nodes count verification
	activeNodes := s.bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	require.Equal(s.T(), s.nodesCount(), len(activeNodes))
}

// TearDownSuite shutdowns all nodes in network, calls once after all tests in suite finished
func (s *testSuite) TearDownSuite() {
	log.Infoln("TearDownSuite")
	log.Infoln("Stop network nodes")
	for _, n := range s.networkNodes {
		err := n.componentManager.Stop(s.ctx)
		s.NoError(err)
	}
	log.Infoln("Stop bootstrap nodes")
	for _, n := range s.bootstrapNodes {
		err := n.componentManager.Stop(s.ctx)
		s.NoError(err)
	}
}

// nodesCount returns count of nodes in network without testNode
func (s *testSuite) nodesCount() int {
	return len(s.bootstrapNodes) + len(s.networkNodes)
}

type PhaseTimeOut uint8

const (
	Disable = PhaseTimeOut(iota + 1)
	Partial
	Full
)

func (s *testSuite) InitTestNode() {
	if s.testNode.componentManager != nil {
		err := s.testNode.componentManager.Init(s.ctx)
		s.NoError(err)
	}
}

func (s *testSuite) StartTestNode() {
	if s.testNode.componentManager != nil {
		err := s.testNode.componentManager.Start(s.ctx)
		s.NoError(err)
	}
}

func (s *testSuite) StopTestNode() {
	if s.testNode.componentManager != nil {
		err := s.testNode.componentManager.Stop(s.ctx)
		s.NoError(err)
	}
}

type networkNode struct {
	id                  core.RecordRef
	privateKey          crypto.PrivateKey
	cryptographyService core.CryptographyService
	host                string

	componentManager *component.Manager
	serviceNetwork   *ServiceNetwork
}

// newNetworkNode returns networkNode initialized only with id, host address and key pair
func newNetworkNode() *networkNode {
	key, err := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	if err != nil {
		panic(err.Error())
	}
	address := "127.0.0.1:" + strconv.Itoa(testNetworkPort)
	testNetworkPort += 2 // coz consensus transport port+=1

	return &networkNode{
		id:                  testutils.RandomRef(),
		privateKey:          key,
		cryptographyService: cryptography.NewKeyBoundCryptographyService(key),
		host:                address,
	}
}

func (s *testSuite) initCrypto(node *networkNode, ref core.RecordRef) (*certificate.CertificateManager, core.CryptographyService) {
	pubKey, err := node.cryptographyService.GetPublicKey()
	s.NoError(err)

	// init certificate

	proc := platformpolicy.NewKeyProcessor()
	publicKey, err := proc.ExportPublicKey(pubKey)
	s.NoError(err)

	cert := &certificate.Certificate{}
	cert.PublicKey = string(publicKey[:])
	cert.Reference = ref.String()
	cert.Role = "virtual"
	cert.BootstrapNodes = make([]certificate.BootstrapNode, 0)

	for _, b := range s.bootstrapNodes {
		pubKey, _ := b.cryptographyService.GetPublicKey()
		pubKeyBuf, err := proc.ExportPublicKey(pubKey)
		s.NoError(err)

		bootstrapNode := certificate.NewBootstrapNode(
			pubKey,
			string(pubKeyBuf[:]),
			b.host,
			b.id.String())

		cert.BootstrapNodes = append(cert.BootstrapNodes, *bootstrapNode)
	}

	// dump cert and read it again from json for correct private files initialization
	jsonCert, err := cert.Dump()
	s.NoError(err)
	log.Infof("cert: %s", jsonCert)

	cert, err = certificate.ReadCertificateFromReader(pubKey, proc, strings.NewReader(jsonCert))
	s.NoError(err)
	return certificate.NewCertificateManager(cert), node.cryptographyService
}

// initNode inits previously created node
func (s *testSuite) initNode(node *networkNode, timeOut PhaseTimeOut) {
	cfg := configuration.NewConfiguration()
	cfg.Host.Transport.Address = node.host

	scheme := platformpolicy.NewPlatformCryptographyScheme()
	serviceNetwork, err := NewServiceNetwork(cfg, scheme)
	s.NoError(err)

	pulseStorageMock := testutils.NewPulseStorageMock(s.T())
	pulseStorageMock.CurrentMock.Set(func(p context.Context) (r *core.Pulse, r1 error) {
		return &core.Pulse{PulseNumber: 0}, nil

	})

	pulseManagerMock := testutils.NewPulseManagerMock(s.T())
	pulseManagerMock.SetMock.Set(func(p context.Context, p1 core.Pulse, p2 bool) (r error) {
		return nil
	})

	netCoordinator := testutils.NewNetworkCoordinatorMock(s.T())
	netCoordinator.ValidateCertMock.Set(func(p context.Context, p1 core.AuthorizationCertificate) (bool, error) {
		return true, nil
	})
	netCoordinator.WriteActiveNodesMock.Set(func(p context.Context, p1 core.PulseNumber, p2 []core.Node) (r error) {
		return nil
	})

	amMock := testutils.NewArtifactManagerMock(s.T())
	amMock.StateMock.Set(func() (r []byte, r1 error) {
		return make([]byte, packets.HashLength), nil
	})

	pubKey, _ := node.cryptographyService.GetPublicKey()

	origin := nodenetwork.NewNode(node.id, core.StaticRoleVirtual, pubKey, node.host, "")
	certManager, cryptographyService := s.initCrypto(node, origin.ID())
	netSwitcher := testutils.NewNetworkSwitcherMock(s.T())
	netSwitcher.GetStateMock.Set(func() (r core.NetworkState) {
		return core.CompleteNetworkState
	})

	realKeeper := nodenetwork.NewNodeKeeper(origin)

	if len(certManager.GetCertificate().GetDiscoveryNodes()) == 0 || utils.OriginIsDiscovery(certManager.GetCertificate()) {
		realKeeper.SetState(network.Ready)
		realKeeper.AddActiveNodes([]core.Node{origin})
	}

	var keeper network.NodeKeeper
	keeper = &nodeKeeperWrapper{realKeeper}

	node.componentManager = &component.Manager{}
	node.componentManager.Register(keeper, pulseManagerMock, pulseStorageMock, netCoordinator, amMock, realKeeper)
	node.componentManager.Register(certManager, cryptographyService)
	node.componentManager.Inject(serviceNetwork, netSwitcher)
	node.serviceNetwork = serviceNetwork
	/*
		var phaseManager phases.PhaseManager
		switch timeOut {
		case Disable:
			phaseManager = &phaseManagerWrapper{original: node.serviceNetwork.PhaseManager}
		case Full:
			phaseManager = &FullTimeoutPhaseManager{}
		case Partial:
			phaseManager = &PartialTimeoutPhaseManager{}
			keeper = &nodeKeeperWrapper{realKeeper}
		}

		node.serviceNetwork.PhaseManager = phaseManager
	*/
}
