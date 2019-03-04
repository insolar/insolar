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
	"crypto"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/consensus/claimhandler"
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/network/utils"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	testNetworkPort       = 10010
	pulseTimeMs     int32 = 5000
	reqTimeoutMs    int32 = 2000
	pulseDelta      int32 = 5
)

type fixture struct {
	ctx            context.Context
	bootstrapNodes []*networkNode
	networkNodes   []*networkNode
	pulsar         TestPulsar
}

func newFixture() *fixture {
	return &fixture{
		ctx:            context.Background(),
		bootstrapNodes: make([]*networkNode, 0),
		networkNodes:   make([]*networkNode, 0),
	}
}

type testSuite struct {
	suite.Suite
	fixtureMap     map[string]*fixture
	bootstrapCount int
	nodesCount     int
}

func NewTestSuite(bootstrapCount, nodesCount int) *testSuite {
	return &testSuite{
		Suite:          suite.Suite{},
		fixtureMap:     make(map[string]*fixture, 0),
		bootstrapCount: bootstrapCount,
		nodesCount:     nodesCount,
	}
}

func (s *testSuite) fixture() *fixture {
	return s.fixtureMap[s.T().Name()]
}

// SetupSuite creates and run network with bootstrap and common nodes once before run all tests in the suite
func (s *testSuite) SetupTest() {
	s.fixtureMap[s.T().Name()] = newFixture()
	var err error
	s.fixture().pulsar, err = NewTestPulsar(pulseTimeMs, reqTimeoutMs, pulseDelta)
	require.NoError(s.T(), err)

	log.Info("SetupTest")

	for i := 0; i < s.bootstrapCount; i++ {
		s.fixture().bootstrapNodes = append(s.fixture().bootstrapNodes, newNetworkNode())
	}

	for i := 0; i < s.nodesCount; i++ {
		s.fixture().networkNodes = append(s.fixture().networkNodes, newNetworkNode())
	}

	pulseReceivers := make([]string, 0)
	for _, node := range s.fixture().bootstrapNodes {
		pulseReceivers = append(pulseReceivers, node.host)
	}

	log.Info("Start test pulsar")
	err = s.fixture().pulsar.Start(s.fixture().ctx, pulseReceivers)
	require.NoError(s.T(), err)

	log.Info("Setup bootstrap nodes")
	s.SetupNodesNetwork(s.fixture().bootstrapNodes)

	<-time.After(time.Second * 2)
	activeNodes := s.fixture().bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	require.Equal(s.T(), len(s.fixture().bootstrapNodes), len(activeNodes))

	if len(s.fixture().networkNodes) > 0 {
		log.Info("Setup network nodes")
		s.SetupNodesNetwork(s.fixture().networkNodes)
		s.waitForConsensus(2)

		// active nodes count verification
		activeNodes1 := s.fixture().networkNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
		activeNodes2 := s.fixture().networkNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()

		require.Equal(s.T(), s.getNodesCount(), len(activeNodes1))
		require.Equal(s.T(), s.getNodesCount(), len(activeNodes2))
	}
	fmt.Println("=================== SetupTest() Done")
}

func (s *testSuite) SetupNodesNetwork(nodes []*networkNode) {
	for _, node := range nodes {
		s.preInitNode(node)
	}

	results := make(chan error, len(nodes))
	initNode := func(node *networkNode) {
		err := node.init(s.fixture().ctx)
		results <- err
	}
	startNode := func(node *networkNode) {
		err := node.componentManager.Start(s.fixture().ctx)
		results <- err
	}

	waitResults := func(results chan error, expected int) error {
		count := 0
		for {
			select {
			case err := <-results:
				count++
				s.NoError(err)
				if count == expected {
					return nil
				}
			case <-time.After(time.Second * 20):
				return errors.New("timeout")
			}
		}
	}

	log.Info("Init nodes")
	for _, node := range nodes {
		go initNode(node)
	}

	err := waitResults(results, len(nodes))
	s.NoError(err)

	log.Info("Start nodes")
	for _, node := range nodes {
		go startNode(node)
	}

	err = waitResults(results, len(nodes))
	s.NoError(err)
}

// TearDownSuite shutdowns all nodes in network, calls once after all tests in suite finished
func (s *testSuite) TearDownTest() {
	log.Info("=================== TearDownTest()")
	log.Info("Stop network nodes")
	for _, n := range s.fixture().networkNodes {
		err := n.componentManager.Stop(s.fixture().ctx)
		s.NoError(err)
	}
	log.Info("Stop bootstrap nodes")
	for _, n := range s.fixture().bootstrapNodes {
		err := n.componentManager.Stop(s.fixture().ctx)
		s.NoError(err)
	}
	log.Info("Stop test pulsar")
	s.fixture().pulsar.Stop(s.fixture().ctx)
}

func (s *testSuite) waitForConsensus(consensusCount int) {
	for i := 0; i < consensusCount; i++ {
		for _, n := range s.fixture().bootstrapNodes {
			err := <-n.consensusResult
			s.NoError(err)
		}

		for _, n := range s.fixture().networkNodes {
			err := <-n.consensusResult
			s.NoError(err)
		}
	}
}

func (s *testSuite) waitForConsensusExcept(consensusCount int, exception core.RecordRef) {
	for i := 0; i < consensusCount; i++ {
		for _, n := range s.fixture().bootstrapNodes {
			if n.id.Equal(exception) {
				continue
			}
			err := <-n.consensusResult
			s.NoError(err)
		}

		for _, n := range s.fixture().networkNodes {
			if n.id.Equal(exception) {
				continue
			}
			err := <-n.consensusResult
			s.NoError(err)
		}
	}
}

// nodesCount returns count of nodes in network without testNode
func (s *testSuite) getNodesCount() int {
	return len(s.fixture().bootstrapNodes) + len(s.fixture().networkNodes)
}

func (s *testSuite) getMaxJoinCount() int {
	return int(float64(s.getNodesCount()) * claimhandler.NodesToJoinPercent)
}

func (s *testSuite) InitNode(node *networkNode) {
	if node.componentManager != nil {
		err := node.init(s.fixture().ctx)
		s.NoError(err)
	}
}

func (s *testSuite) StartNode(node *networkNode) {
	if node.componentManager != nil {
		err := node.componentManager.Start(s.fixture().ctx)
		s.NoError(err)
	}
}

func (s *testSuite) StopNode(node *networkNode) {
	if node.componentManager != nil {
		err := node.componentManager.Stop(s.fixture().ctx)
		s.NoError(err)
	}
}

type networkNode struct {
	id                  core.RecordRef
	role                core.StaticRole
	privateKey          crypto.PrivateKey
	cryptographyService core.CryptographyService
	host                string

	componentManager *component.Manager
	serviceNetwork   *ServiceNetwork
	consensusResult  chan error
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
		role:                RandomRole(),
		privateKey:          key,
		cryptographyService: cryptography.NewKeyBoundCryptographyService(key),
		host:                address,
		consensusResult:     make(chan error, 30),
	}
}

// init calls Init for node component manager and wraps PhaseManager
func (n *networkNode) init(ctx context.Context) error {
	err := n.componentManager.Init(ctx)
	n.serviceNetwork.PhaseManager = &phaseManagerWrapper{original: n.serviceNetwork.PhaseManager, result: n.consensusResult}
	n.serviceNetwork.NodeKeeper = &nodeKeeperWrapper{original: n.serviceNetwork.NodeKeeper}
	return err
}

func (s *testSuite) initCrypto(node *networkNode) (*certificate.CertificateManager, core.CryptographyService) {
	pubKey, err := node.cryptographyService.GetPublicKey()
	s.NoError(err)

	// init certificate

	proc := platformpolicy.NewKeyProcessor()
	publicKey, err := proc.ExportPublicKeyPEM(pubKey)
	s.NoError(err)

	cert := &certificate.Certificate{}
	cert.PublicKey = string(publicKey[:])
	cert.Reference = node.id.String()
	cert.Role = node.role.String()
	cert.BootstrapNodes = make([]certificate.BootstrapNode, 0)

	for _, b := range s.fixture().bootstrapNodes {
		pubKey, _ := b.cryptographyService.GetPublicKey()
		pubKeyBuf, err := proc.ExportPublicKeyPEM(pubKey)
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

func RandomRole() core.StaticRole {
	i := rand.Int()%3 + 1
	return core.StaticRole(i)
}

type terminationHandler struct {
	NodeID core.RecordRef
}

func (t *terminationHandler) Abort() {
	log.Errorf("Abort node: %s", t.NodeID)
}

type pulseManagerMock struct {
	pulse core.Pulse
	lock  sync.Mutex

	keeper network.NodeKeeper
}

func newPulseManagerMock(keeper network.NodeKeeper) *pulseManagerMock {
	return &pulseManagerMock{pulse: *core.GenesisPulse, keeper: keeper}
}

func (p *pulseManagerMock) Current(ctx context.Context) (*core.Pulse, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	return &p.pulse, nil
}

func (p *pulseManagerMock) Set(ctx context.Context, pulse core.Pulse, persist bool) error {
	p.lock.Lock()
	p.pulse = pulse
	p.lock.Unlock()

	return p.keeper.MoveSyncToActive(ctx)
}

// preInitNode inits previously created node with mocks and external dependencies
func (s *testSuite) preInitNode(node *networkNode) {
	cfg := configuration.NewConfiguration()
	cfg.Pulsar.PulseTime = pulseTimeMs // pulse 5 sec for faster tests
	cfg.Host.Transport.Address = node.host
	cfg.Service.Skip = 5

	node.componentManager = &component.Manager{}
	node.componentManager.Register(platformpolicy.NewPlatformCryptographyScheme())
	serviceNetwork, err := NewServiceNetwork(cfg, node.componentManager, false)
	s.NoError(err)

	netCoordinator := testutils.NewNetworkCoordinatorMock(s.T())
	netCoordinator.ValidateCertMock.Set(func(p context.Context, p1 core.AuthorizationCertificate) (bool, error) {
		return true, nil
	})

	netCoordinator.IsStartedMock.Set(func() (r bool) {
		return true
	})

	amMock := testutils.NewArtifactManagerMock(s.T())
	amMock.StateMock.Set(func() (r []byte, r1 error) {
		return make([]byte, packets.HashLength), nil
	})

	pubKey, _ := node.cryptographyService.GetPublicKey()

	origin := nodenetwork.NewNode(node.id, node.role, pubKey, node.host, "")
	certManager, cryptographyService := s.initCrypto(node)

	realKeeper := nodenetwork.NewNodeKeeper(origin)
	terminationHandler := &terminationHandler{NodeID: origin.ID()}

	realKeeper.SetState(core.WaitingNodeNetworkState)
	if len(certManager.GetCertificate().GetDiscoveryNodes()) == 0 || utils.OriginIsDiscovery(certManager.GetCertificate()) {
		realKeeper.SetState(core.ReadyNodeNetworkState)
		realKeeper.AddActiveNodes([]core.Node{origin})
	}

	node.componentManager.Register(terminationHandler, realKeeper, newPulseManagerMock(realKeeper), netCoordinator, amMock)
	node.componentManager.Register(certManager, cryptographyService)
	node.componentManager.Inject(serviceNetwork, NewTestNetworkSwitcher())
	node.serviceNetwork = serviceNetwork
}
