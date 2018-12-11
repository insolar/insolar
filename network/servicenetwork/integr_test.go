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
	"bytes"
	"context"
	"crypto"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/consensus/phases"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
	ctx            context.Context
	bootstrapNodes []networkNode
	networkNodes   []networkNode
	testNode       networkNode
	networkPort    int
}

func NewTestSuite(bootstrapCount, nodesCount int) *testSuite {
	return &testSuite{
		Suite:          suite.Suite{},
		ctx:            context.Background(),
		networkPort:    10001,
		bootstrapNodes: make([]networkNode, bootstrapCount),
		networkNodes:   make([]networkNode, nodesCount),
	}
}

// SetupSuite creates and run network with bootstrap and common nodes once before run all tests in the suite
func (s *testSuite) SetupSuite() {
	log.Infoln("SetupSuite")
	s.createBootstrapNodes()

	for i := 0; i < cap(s.networkNodes); i++ {
		s.networkNodes = append(s.networkNodes, s.createNetworkNode(s.T(), Disable))
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

	//TODO: wait for first consensus
	// active nodes count verification
	activeNodes := s.bootstrapNodes[0].serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.nodesCount(), len(activeNodes))
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
	Partitial
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
		err := s.testNode.componentManager.Init(s.ctx)
		s.NoError(err)
		err = s.testNode.componentManager.Start(s.ctx)
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

	componentManager *component.Manager
	serviceNetwork   *ServiceNetwork
}

// newNetworkNode returns networkNode initialized only with id and key pair
func newNetworkNode() networkNode {
	key, err := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	if err != nil {
		panic(err.Error())
	}
	return networkNode{
		id:                  testutils.RandomRef(),
		privateKey:          key,
		cryptographyService: cryptography.NewKeyBoundCryptographyService(key),
	}
}

func initCertificate(t *testing.T, nodes []certificate.BootstrapNode, key crypto.PublicKey, ref core.RecordRef) *certificate.CertificateManager {
	proc := platformpolicy.NewKeyProcessor()
	publicKey, err := proc.ExportPublicKey(key)
	assert.NoError(t, err)
	bytes.NewReader(publicKey)

	type сertInfo map[string]interface{}
	j := сertInfo{
		"public_key": string(publicKey[:]),
	}

	data, err := json.Marshal(j)

	cert, err := certificate.ReadCertificateFromReader(key, proc, bytes.NewReader(data))
	cert.Reference = ref.String()
	assert.NoError(t, err)
	cert.BootstrapNodes = nodes
	return certificate.NewCertificateManager(cert)
}

func initCrypto(t *testing.T, nodes []certificate.BootstrapNode, ref core.RecordRef) (*certificate.CertificateManager, core.CryptographyService) {
	key, err := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	assert.NoError(t, err)
	require.NotNil(t, key)
	cs := cryptography.NewKeyBoundCryptographyService(key)
	pubKey, err := cs.GetPublicKey()
	assert.NoError(t, err)
	mngr := initCertificate(t, nodes, pubKey, ref)

	return mngr, cs
}

func (s *testSuite) getBootstrapNodes(t *testing.T) []certificate.BootstrapNode {
	result := make([]certificate.BootstrapNode, 0)
	for _, b := range s.bootstrapNodes {
		node := certificate.NewBootstrapNode(
			b.serviceNetwork.CertificateManager.GetCertificate().GetPublicKey(),
			b.serviceNetwork.CertificateManager.GetCertificate().(*certificate.Certificate).PublicKey,
			b.serviceNetwork.cfg.Host.Transport.Address,
			b.serviceNetwork.NodeNetwork.GetOrigin().ID().String())
		result = append(result, *node)
	}
	return result
}

func (s *testSuite) createBootstrapNodes() {
	for i := 0; i < cap(s.bootstrapNodes); i++ {
		s.bootstrapNodes = append(s.bootstrapNodes, newNetworkNode())
	}

	//initCertificate()
	// genesis makeCertificates

	// VerifyAuthorizationCertificate()

	// generate node ids and key pairs
	// create each node with createNetworkNode

	s.bootstrapNodes = append(s.bootstrapNodes, s.createNetworkNode(s.T(), Disable))
}

func (s *testSuite) createNetworkNode(t *testing.T, timeOut PhaseTimeOut) networkNode {
	address := "127.0.0.1:" + strconv.Itoa(s.networkPort)
	s.networkPort += 2 // coz consensus transport port+=1

	origin := nodenetwork.NewNode(testutils.RandomRef(),
		core.StaticRoleVirtual,
		nil,
		address,
		"",
	)

	cfg := configuration.NewConfiguration()
	cfg.Host.Transport.Address = address

	scheme := platformpolicy.NewPlatformCryptographyScheme()
	serviceNetwork, err := NewServiceNetwork(cfg, scheme)
	assert.NoError(t, err)

	pulseManagerMock := testutils.NewPulseManagerMock(t)
	pulseManagerMock.CurrentMock.Set(func(p context.Context) (r *core.Pulse, r1 error) {
		return &core.Pulse{PulseNumber: 0}, nil
	})
	pulseManagerMock.SetMock.Set(func(p context.Context, p1 core.Pulse, p2 bool) (r error) {
		return nil
	})

	netCoordinator := testutils.NewNetworkCoordinatorMock(t)
	netCoordinator.ValidateCertMock.Set(func(p context.Context, p1 core.AuthorizationCertificate) (bool, error) {
		return true, nil
	})
	netCoordinator.WriteActiveNodesMock.Set(func(p context.Context, p1 core.PulseNumber, p2 []core.Node) (r error) {
		return nil
	})

	amMock := testutils.NewArtifactManagerMock(t)
	amMock.StateMock.Set(func() (r []byte, r1 error) {
		return make([]byte, 0), nil
	})

	certManager, cryptographyService := initCrypto(t, s.getBootstrapNodes(t), origin.ID())
	netSwitcher := testutils.NewNetworkSwitcherMock(t)
	realKeeper := nodenetwork.NewNodeKeeper(origin)
	var keeper network.NodeKeeper
	keeper = &nodeKeeperWrapper{realKeeper}

	var phaseManager phases.PhaseManager
	switch timeOut {
	case Disable:
		phaseManager = phases.NewPhaseManager()
	case Full:
		phaseManager = &FullTimeoutPhaseManager{}
	case Partitial:
		phaseManager = &PartialTimeoutPhaseManager{}
		keeper = &nodeKeeperWrapper{realKeeper}
	}

	cm := &component.Manager{}
	cm.Register(keeper, pulseManagerMock, netCoordinator, amMock, realKeeper)
	cm.Register(certManager, cryptographyService, phaseManager)
	cm.Inject(serviceNetwork, netSwitcher)

	return networkNode{componentManager: cm, serviceNetwork: serviceNetwork}
}

func (s *testSuite) TestNodeConnect() {
	phasesResult := make(chan error)
	s.testNode = s.createNetworkNode(s.T(), Disable)
	s.InitTestNode()
	s.StartTestNode()

	res := <-phasesResult
	s.NoError(res)
	activeNodes := s.testNode.serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.nodesCount()+1, len(activeNodes))
	// teardown
	<-time.After(time.Second * 5)
	s.StopTestNode()
}

func (s *testSuite) TestNodeLeave() {

	s.testNode = s.createNetworkNode(s.T(), Disable)

	phasesResult := make(chan error)
	s.InitTestNode()
	s.StartTestNode()
	res := <-phasesResult
	s.NoError(res)
	activeNodes := s.testNode.serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(2, len(activeNodes))

	// teardown
	<-time.After(time.Second * 5)

	res = <-phasesResult
	s.NoError(res)
	activeNodes = s.testNode.serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(s.nodesCount(), len(activeNodes))

	s.StopTestNode()
}

func TestServiceNetworkIntegration(t *testing.T) {
	s := NewTestSuite(1, 0)
	suite.Run(t, s)
}

// Full timeout test

type FullTimeoutPhaseManager struct {
}

func (ftpm *FullTimeoutPhaseManager) OnPulse(ctx context.Context, pulse *core.Pulse) error {
	return nil
}

func (s *testSuite) TestFullTimeOut() {
	s.T().Skip("will be available after phase result fix !")
	networkNodesCount := 5
	phasesResult := make(chan error)
	bootstrapNode1 := s.createNetworkNode(s.T(), Disable)
	s.bootstrapNodes = append(s.bootstrapNodes, bootstrapNode1)

	s.testNode = s.createNetworkNode(s.T(), Full)

	for i := 0; i < networkNodesCount; i++ {
		s.networkNodes = append(s.networkNodes, s.createNetworkNode(s.T(), Disable))
	}

	s.InitTestNode()
	s.StartTestNode()
	res := <-phasesResult
	s.NoError(res)
	activeNodes := s.testNode.serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(1, len(activeNodes))
	// teardown
	<-time.After(time.Second * 5)
	s.StopTestNode()
}

// Partitial timeout

func (s *testSuite) TestPartialTimeOut() {
	networkNodesCount := 5
	phasesResult := make(chan error)
	bootstrapNode1 := s.createNetworkNode(s.T(), Disable)
	s.bootstrapNodes = append(s.bootstrapNodes, bootstrapNode1)

	s.testNode = s.createNetworkNode(s.T(), Partitial)

	for i := 0; i < networkNodesCount; i++ {
		s.networkNodes = append(s.networkNodes, s.createNetworkNode(s.T(), Disable))
	}

	s.InitTestNode()
	s.StartTestNode()
	res := <-phasesResult
	s.NoError(res)
	// activeNodes := s.testNode.serviceNetwork.NodeKeeper.GetActiveNodes()
	// s.Equal(1, len(activeNodes))	// TODO: do test check
	// teardown
	<-time.After(time.Second * 5)
	s.StopTestNode()
}

type PartialTimeoutPhaseManager struct {
	FirstPhase  *phases.FirstPhase
	SecondPhase *phases.SecondPhase
	ThirdPhase  *phases.ThirdPhase
	Keeper      network.NodeKeeper `inject:""`
}

func (ftpm *PartialTimeoutPhaseManager) OnPulse(ctx context.Context, pulse *core.Pulse) error {
	var err error

	pulseDuration, err := getPulseDuration(pulse)
	if err != nil {
		return errors.Wrap(err, "[ OnPulse ] Failed to get pulse duration")
	}

	var tctx context.Context
	var cancel context.CancelFunc

	tctx, cancel = contextTimeout(ctx, *pulseDuration, 0.2)
	defer cancel()

	firstPhaseState, err := ftpm.FirstPhase.Execute(tctx, pulse)

	if err != nil {
		return errors.Wrap(err, "[ TestCase.OnPulse ] failed to execute a phase")
	}

	tctx, cancel = contextTimeout(ctx, *pulseDuration, 0.2)
	defer cancel()

	secondPhaseState, err := ftpm.SecondPhase.Execute(tctx, firstPhaseState)
	checkError(err)

	fmt.Println(secondPhaseState) // TODO: remove after use
	checkError(ftpm.ThirdPhase.Execute(ctx, secondPhaseState))

	return nil
}

func contextTimeout(ctx context.Context, duration time.Duration, k float64) (context.Context, context.CancelFunc) {
	timeout := time.Duration(k * float64(duration))
	timedCtx, cancelFund := context.WithTimeout(ctx, timeout)
	return timedCtx, cancelFund
}

func getPulseDuration(pulse *core.Pulse) (*time.Duration, error) {
	duration := time.Duration(pulse.PulseNumber-pulse.PrevPulseNumber) * time.Second
	return &duration, nil
}

func checkError(err error) {
	if err != nil {
		log.Error(err)
	}
}
