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
	"github.com/insolar/insolar/testutils/nodekeeper"
	"github.com/insolar/insolar/testutils/terminationhandler"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
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
		GlobuleID:     0,
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
		GlobuleID:     0,
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
	// FIXME: TmpLedger is deprecated. Use mocks instead.
	l, _, clean := ledgertestutils.TmpLedger(t, "", core.StaticRoleLightMaterial, core.Components{}, true)

	calculator := &calculator{}

	cm := component.Manager{}

	key, _ := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	require.NotNil(t, key)

	service := testutils.NewCryptographyServiceMock(t)
	service.SignFunc = func(p []byte) (r *core.Signature, r1 error) {
		return nil, errors.New("Sign error")
	}
	service.GetPublicKeyFunc = func() (r crypto.PublicKey, r1 error) {
		return "key", nil
	}
	scheme := platformpolicy.NewPlatformCryptographyScheme()

	pulseManager := testutils.NewPulseStorageMock(t)

	nk := nodekeeper.GetTestNodekeeper(service)
	th := terminationhandler.NewTestTerminationHandler()

	jc := testutils.NewJetCoordinatorMock(t)

	cm.Inject(th, nk, jc, l.ArtifactManager, calculator, service, scheme, pulseManager)

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
	require.NotNil(t, key)

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
	th := terminationhandler.NewTestTerminationHandler()
	cm.Inject(th, nk, am, calculator, service, scheme)

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

	ph, np, err := calculator.GetPulseProof(&PulseEntry{Pulse: pulse})

	require.Error(t, err)
	require.Contains(t, err.Error(), "[ GetPulseProof ] Failed to get node stateHash")
	require.Nil(t, np)
	require.Nil(t, ph)
}
