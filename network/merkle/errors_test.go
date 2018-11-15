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
	"crypto"
	"encoding/hex"
	"testing"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/ledgertestutils"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar/pulsartestutils"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/certificate"
	"github.com/insolar/insolar/testutils/nodekeeper"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type calculatorErrorSuite struct {
	suite.Suite

	pulse       *core.Pulse
	nodeNetwork core.NodeNetwork
	service     core.CryptographyService

	calculator Calculator
}

func (t *calculatorErrorSuite) TestGetNodeProofError() {
	ph, np, err := t.calculator.GetPulseProof(&PulseEntry{Pulse: t.pulse})

	t.Assert().Error(err)
	t.Assert().Contains(err.Error(), "[ GetPulseProof ] Failed to sign node info hash")
	t.Assert().Nil(np)
	t.Assert().Nil(ph)
}

func (t *calculatorErrorSuite) TestGetGlobuleProofCalculateError() {
	pulseEntry := &PulseEntry{Pulse: t.pulse}

	prevCloudHash, _ := hex.DecodeString(
		"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	)

	globuleEntry := &GlobuleEntry{
		PulseEntry:    pulseEntry,
		PulseHash:     nil,
		ProofSet:      nil,
		PrevCloudHash: prevCloudHash,
		GlobuleIndex:  0,
	}
	gh, gp, err := t.calculator.GetGlobuleProof(globuleEntry)

	t.Assert().Error(err)
	t.Assert().Contains(err.Error(), "[ GetGlobuleProof ] Failed to calculate node root")
	t.Assert().Nil(gh)
	t.Assert().Nil(gp)
}

func (t *calculatorErrorSuite) TestGetGlobuleProofSignError() {
	pulseEntry := &PulseEntry{Pulse: t.pulse}

	prevCloudHash, _ := hex.DecodeString(
		"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	)

	globuleEntry := &GlobuleEntry{
		PulseEntry: pulseEntry,
		PulseHash:  nil,
		ProofSet: map[core.Node]*PulseProof{
			t.nodeNetwork.GetOrigin(): {},
		},
		PrevCloudHash: prevCloudHash,
		GlobuleIndex:  0,
	}
	gh, gp, err := t.calculator.GetGlobuleProof(globuleEntry)

	t.Assert().Error(err)
	t.Assert().Contains(err.Error(), "[ GetGlobuleProof ] Failed to sign globule hash")
	t.Assert().Nil(gh)
	t.Assert().Nil(gp)
}

func (t *calculatorErrorSuite) TestGetCloudProofSignError() {
	ch, cp, err := t.calculator.GetCloudProof(&CloudEntry{
		ProofSet: []*GlobuleProof{
			{},
		},
		PrevCloudHash: nil,
	})

	t.Assert().Error(err)
	t.Assert().Contains(err.Error(), "[ GetCloudProof ] Failed to sign cloud hash")
	t.Assert().Nil(ch)
	t.Assert().Nil(cp)
}

func (t *calculatorErrorSuite) TestGetCloudProofCalculateError() {
	ch, cp, err := t.calculator.GetCloudProof(&CloudEntry{
		ProofSet:      nil,
		PrevCloudHash: nil,
	})

	t.Assert().Error(err)
	t.Assert().Contains(err.Error(), "[ GetCloudProof ] Failed to calculate cloud hash")
	t.Assert().Nil(ch)
	t.Assert().Nil(cp)
}

func TestCalculatorError(t *testing.T) {
	c := certificate.GetTestCertificate()
	// FIXME: TmpLedger is deprecated. Use mocks instead.
	l, clean := ledgertestutils.TmpLedger(t, "", core.Components{})

	calculator := &calculator{}

	cm := component.Manager{}

	key, _ := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	assert.NotNil(t, key)

	service := testutils.NewCryptographyServiceMock(t)
	service.SignFunc = func(p []byte) (r *core.Signature, r1 error) {
		return nil, errors.New("Sign error")
	}
	service.GetPublicKeyFunc = func() (r crypto.PublicKey, r1 error) {
		return "key", nil
	}
	scheme := platformpolicy.NewPlatformCryptographyScheme()

	nk := nodekeeper.GetTestNodekeeper(service)
	cm.Inject(nk, l.ArtifactManager, c, calculator, service, scheme)

	assert.NotNil(t, calculator.ArtifactManager)
	assert.NotNil(t, calculator.NodeNetwork)
	assert.NotNil(t, calculator.CryptographyService)
	assert.NotNil(t, calculator.PlatformCryptographyScheme)

	err := cm.Init(context.Background())
	assert.NoError(t, err)

	pulse := &core.Pulse{
		PulseNumber:     core.PulseNumber(1337),
		NextPulseNumber: core.PulseNumber(1347),
		Entropy:         pulsartestutils.MockEntropyGenerator{}.GenerateEntropy(),
	}

	s := &calculatorErrorSuite{
		Suite:       suite.Suite{},
		calculator:  calculator,
		pulse:       pulse,
		nodeNetwork: nk,
		service:     service,
	}
	suite.Run(t, s)

	clean()
}

func TestCalculatorLedgerError(t *testing.T) {
	calculator := &calculator{}

	cm := component.Manager{}

	key, _ := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	assert.NotNil(t, key)

	service := testutils.NewCryptographyServiceMock(t)
	service.SignFunc = func(p []byte) (r *core.Signature, r1 error) {
		return nil, errors.New("Sign error")
	}
	service.GetPublicKeyFunc = func() (r crypto.PublicKey, r1 error) {
		return "key", nil
	}

	am := testutils.NewArtifactManagerMock(t)
	am.StateFunc = func() (r []byte, r1 error) {
		return nil, errors.New("State error")
	}

	scheme := platformpolicy.NewPlatformCryptographyScheme()
	nk := nodekeeper.GetTestNodekeeper(service)
	cm.Inject(nk, am, calculator, service, scheme)

	assert.NotNil(t, calculator.ArtifactManager)
	assert.NotNil(t, calculator.NodeNetwork)
	assert.NotNil(t, calculator.CryptographyService)
	assert.NotNil(t, calculator.PlatformCryptographyScheme)

	err := cm.Init(context.Background())
	assert.NoError(t, err)

	pulse := &core.Pulse{
		PulseNumber:     core.PulseNumber(1337),
		NextPulseNumber: core.PulseNumber(1347),
		Entropy:         pulsartestutils.MockEntropyGenerator{}.GenerateEntropy(),
	}

	ph, np, err := calculator.GetPulseProof(&PulseEntry{Pulse: pulse})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "[ GetPulseProof ] Failed to get node stateHash")
	assert.Nil(t, np)
	assert.Nil(t, ph)
}
