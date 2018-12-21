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
	"strconv"
	"testing"
	"time"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
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

func NewTestSuite() *testSuite {
	return &testSuite{
		Suite:        suite.Suite{},
		ctx:          context.Background(),
		networkNodes: make([]networkNode, 0),
		networkPort:  10001,
	}
}
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

func (s *testSuite) createNetworkNode(t *testing.T) networkNode {
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
	keeper := &nodeKeeperWrapper{realKeeper}

	cm := &component.Manager{}
	cm.Register(keeper, pulseManagerMock, netCoordinator, amMock, realKeeper)
	cm.Register(certManager, cryptographyService)
	cm.Inject(serviceNetwork, netSwitcher)

	serviceNetwork.NodeKeeper = keeper

	return networkNode{cm, serviceNetwork}
}
