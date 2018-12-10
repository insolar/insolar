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
	"github.com/insolar/insolar/consensus/packets"
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

type nodeKeeperWrapperPartial struct {
	realKeeper network.NodeKeeper
}

func (nkw *nodeKeeperWrapperPartial) SetCloudHash(hash []byte) {
	nkw.realKeeper.SetCloudHash(hash)
}
func (nkw *nodeKeeperWrapperPartial) AddActiveNodes(nodes []core.Node) {
	nkw.realKeeper.AddActiveNodes(nodes)
}
func (nkw *nodeKeeperWrapperPartial) GetActiveNodeByShortID(shortID core.ShortNodeID) core.Node {
	return nkw.realKeeper.GetActiveNodeByShortID(shortID)
}

func (nkw *nodeKeeperWrapperPartial) SetState(state network.NodeKeeperState) {
	nkw.realKeeper.SetState(state)
}

func (nkw *nodeKeeperWrapperPartial) GetState() network.NodeKeeperState {
	return nkw.realKeeper.GetState()
}

func (nkw *nodeKeeperWrapperPartial) GetOriginClaim() (*packets.NodeJoinClaim, error) {
	return nkw.realKeeper.GetOriginClaim()
}

func (nkw *nodeKeeperWrapperPartial) NodesJoinedDuringPreviousPulse() bool {
	return nkw.realKeeper.NodesJoinedDuringPreviousPulse()
}

func (nkw *nodeKeeperWrapperPartial) AddPendingClaim(claim packets.ReferendumClaim) bool {
	return nkw.realKeeper.AddPendingClaim(claim)
}

func (nkw *nodeKeeperWrapperPartial) GetClaimQueue() network.ClaimQueue {
	return nkw.realKeeper.GetClaimQueue()
}

func (nkw *nodeKeeperWrapperPartial) GetUnsyncList() network.UnsyncList {
	return nkw.realKeeper.GetUnsyncList()
}

func (nkw *nodeKeeperWrapperPartial) GetSparseUnsyncList(length int) network.UnsyncList {
	return nkw.realKeeper.GetSparseUnsyncList(length)
}

func (nkw *nodeKeeperWrapperPartial) Sync(list network.UnsyncList) {
	nkw.realKeeper.Sync(list)
}

func (nkw *nodeKeeperWrapperPartial) MoveSyncToActive() {
	nkw.realKeeper.MoveSyncToActive()
}

func (nkw *nodeKeeperWrapperPartial) GetActiveNodes() []core.Node {
	tmp := nkw.realKeeper.GetActiveNodes()
	tmp = tmp[:len(tmp)-2]
	return tmp
}

func (nkw *nodeKeeperWrapperPartial) GetOrigin() core.Node {
	return nkw.realKeeper.GetOrigin()
}

func (nkw *nodeKeeperWrapperPartial) GetActiveNode(ref core.RecordRef) core.Node {
	return nkw.realKeeper.GetActiveNode(ref)
}

func (nkw *nodeKeeperWrapperPartial) GetActiveNodesByRole(role core.DynamicRole) []core.RecordRef {
	return nkw.realKeeper.GetActiveNodesByRole(role)
}

func (nkw *nodeKeeperWrapperPartial) GetCloudHash() []byte {
	return nkw.realKeeper.GetCloudHash()
}

func (nkw *nodeKeeperWrapperPartial) IsBootstrapped() bool {
	return nkw.realKeeper.IsBootstrapped()
}

func (nkw *nodeKeeperWrapperPartial) SetIsBootstrapped(isBootstrap bool) {
	nkw.realKeeper.SetIsBootstrapped(isBootstrap)
}

type testSuite struct {
	suite.Suite
	ctx            context.Context
	bootstrapNodes []networkNode
	networkNodes   []networkNode
	testNode       networkNode
	networkPort    int
}

func NewTestSuite() *testSuite {
	return &testSuite{
		Suite:        suite.Suite{},
		ctx:          context.Background(),
		networkNodes: make([]networkNode, 0),
		networkPort:  10001,
	}
}

type PhaseTimeOut uint8

const (
	Disable = PhaseTimeOut(iota + 1)
	Partitial
	Full
)

func (s *testSuite) InitNodes() {
	for _, n := range s.bootstrapNodes {
		err := n.componentManager.Init(s.ctx)
		s.NoError(err)
	}
	log.Info("========== Bootstrap nodes inited")
	<-time.After(time.Second * 1)

	if s.testNode.componentManager != nil {
		err := s.testNode.componentManager.Init(s.ctx)
		s.NoError(err)
	}
}

func (s *testSuite) StartNodes() {
	for _, n := range s.bootstrapNodes {
		err := n.componentManager.Start(s.ctx)
		s.NoError(err)
	}
	log.Info("========== Bootstrap nodes started")
	<-time.After(time.Second * 1)

	if s.testNode.componentManager != nil {
		err := s.testNode.componentManager.Init(s.ctx)
		s.NoError(err)
		err = s.testNode.componentManager.Start(s.ctx)
		s.NoError(err)
	}

}

func (s *testSuite) StopNodes() {
	for _, n := range s.networkNodes {
		err := n.componentManager.Stop(s.ctx)
		s.NoError(err)
	}

	if s.testNode.componentManager != nil {
		err := s.testNode.componentManager.Stop(s.ctx)
		s.NoError(err)
	}
}

type networkNode struct {
	componentManager *component.Manager
	serviceNetwork   *ServiceNetwork
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
	mngr := certificate.NewCertificateManager(cert)
	return mngr
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
	netCoordinator := testutils.NewNetworkCoordinatorMock(t)
	netCoordinator.ValidateCertMock.Set(func(p context.Context, p1 core.AuthorizationCertificate) (bool, error) {
		return true, nil
	})

	amMock := testutils.NewArtifactManagerMock(t)

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
		keeper = &nodeKeeperWrapperPartial{realKeeper}
	}

	cm := &component.Manager{}
	cm.Register(keeper, pulseManagerMock, netCoordinator, amMock, realKeeper)
	cm.Register(certManager, cryptographyService, phaseManager)
	cm.Inject(serviceNetwork, netSwitcher)

	return networkNode{cm, serviceNetwork}
}

func (s *testSuite) TestNodeConnect() {
	s.T().Skip("will be available after phase result fix !")
	phasesResult := make(chan error)
	bootstrapNode1 := s.createNetworkNode(s.T(), Disable)
	s.bootstrapNodes = append(s.bootstrapNodes, bootstrapNode1)

	s.testNode = s.createNetworkNode(s.T(), Disable)

	s.InitNodes()
	s.StartNodes()
	res := <-phasesResult
	s.NoError(res)
	activeNodes := s.testNode.serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(2, len(activeNodes))
	// teardown
	<-time.After(time.Second * 5)
	s.StopNodes()
}

func TestServiceNetworkIntegration(t *testing.T) {
	s := NewTestSuite()
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

	s.InitNodes()
	s.StartNodes()
	res := <-phasesResult
	s.NoError(res)
	activeNodes := s.testNode.serviceNetwork.NodeKeeper.GetActiveNodes()
	s.Equal(1, len(activeNodes))
	// teardown
	<-time.After(time.Second * 5)
	s.StopNodes()
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

	s.InitNodes()
	s.StartNodes()
	res := <-phasesResult
	s.NoError(res)
	// activeNodes := s.testNode.serviceNetwork.NodeKeeper.GetActiveNodes()
	// s.Equal(1, len(activeNodes))	// TODO: do test check
	// teardown
	<-time.After(time.Second * 5)
	s.StopNodes()
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
