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

package merkle

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar/pulsartestutils"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/nodekeeper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type calculatorSuite struct {
	suite.Suite

	pulse       *core.Pulse
	nodeNetwork core.NodeNetwork
	service     core.CryptographyService

	calculator Calculator
}

func (t *calculatorSuite) TestGetNodeProof() {
	// t.Suite.T().Skip("skipped")
	ph, np, err := t.calculator.GetPulseProof(&PulseEntry{Pulse: t.pulse})

	t.Assert().NoError(err)
	t.Assert().NotNil(np)

	key, err := t.service.GetPublicKey()
	t.Assert().NoError(err)

	t.Assert().True(t.calculator.IsValid(np, ph, key))
}

func (t *calculatorSuite) TestGetGlobuleProof() {
	// t.Suite.T().Skip("skipped")
	pulseEntry := &PulseEntry{Pulse: t.pulse}
	ph, pp, err := t.calculator.GetPulseProof(pulseEntry)
	t.Assert().NoError(err)

	prevCloudHash, _ := hex.DecodeString(
		"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	)

	globuleEntry := &GlobuleEntry{
		PulseEntry: pulseEntry,
		PulseHash:  ph,
		ProofSet: map[core.Node]*PulseProof{
			t.nodeNetwork.GetOrigin(): pp,
		},
		PrevCloudHash: prevCloudHash,
		GlobuleID:     0,
	}
	gh, gp, err := t.calculator.GetGlobuleProof(globuleEntry)

	t.Assert().NoError(err)
	t.Assert().NotNil(gp)

	key, err := t.service.GetPublicKey()
	t.Assert().NoError(err)

	valid := t.calculator.IsValid(gp, gh, key)
	t.Assert().True(valid)
}

func (t *calculatorSuite) TestGetCloudProof() {
	// t.Suite.T().Skip("skipped")
	pulseEntry := &PulseEntry{Pulse: t.pulse}
	ph, pp, err := t.calculator.GetPulseProof(pulseEntry)
	t.Assert().NoError(err)

	prevCloudHash, _ := hex.DecodeString(
		"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	)

	globuleEntry := &GlobuleEntry{
		PulseEntry: pulseEntry,
		PulseHash:  ph,
		ProofSet: map[core.Node]*PulseProof{
			t.nodeNetwork.GetOrigin(): pp,
		},
		PrevCloudHash: prevCloudHash,
		GlobuleID:     0,
	}
	_, gp, err := t.calculator.GetGlobuleProof(globuleEntry)

	ch, cp, err := t.calculator.GetCloudProof(&CloudEntry{
		ProofSet:      []*GlobuleProof{gp},
		PrevCloudHash: prevCloudHash,
	})

	t.Assert().NoError(err)
	t.Assert().NotNil(gp)

	key, err := t.service.GetPublicKey()
	t.Assert().NoError(err)

	valid := t.calculator.IsValid(cp, ch, key)
	t.Assert().True(valid)
}

func TestNewCalculator(t *testing.T) {
	c := NewCalculator()
	require.NotNil(t, c)
}

func TestCalculator(t *testing.T) {
	calculator := &calculator{}

	key, _ := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	require.NotNil(t, key)

	service := cryptography.NewKeyBoundCryptographyService(key)
	scheme := platformpolicy.NewPlatformCryptographyScheme()
	nk := nodekeeper.GetTestNodekeeper(service)

	am := testutils.NewArtifactManagerMock(t)
	am.StateFunc = func() (r []byte, r1 error) {
		return []byte("state"), nil
	}

	cm := component.Manager{}
	cm.Inject(nk, am, calculator, service, scheme)

	require.NotNil(t, calculator.ArtifactManager)
	require.NotNil(t, calculator.NodeNetwork)
	require.NotNil(t, calculator.CryptographyService)
	require.NotNil(t, calculator.PlatformCryptographyScheme)

	err := cm.Init(context.Background())
	require.NoError(t, err)

	pulse := &core.Pulse{
		PulseNumber:     core.PulseNumber(1337),
		NextPulseNumber: core.PulseNumber(1347),
		Entropy:         pulsartestutils.MockEntropyGenerator{}.GenerateEntropy(),
	}

	s := &calculatorSuite{
		Suite:       suite.Suite{},
		calculator:  calculator,
		pulse:       pulse,
		nodeNetwork: nk,
		service:     service,
	}
	suite.Run(t, s)
}
