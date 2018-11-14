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
	"fmt"
	"testing"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar/pulsartestutils"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/ledger"
	"github.com/insolar/insolar/testutils/nodekeeper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func (t *calculatorHashesSuite) TestGetPulseHash() {
	pulseEntry := &PulseEntry{Pulse: t.pulse}
	ph, _, err := t.calculator.GetPulseProof(pulseEntry)
	t.Assert().NoError(err)

	expectedHash, _ := hex.DecodeString(
		"bd18c009950389026c5c6f85c838b899d188ec0d667f77948aa72a49747c3ed31835b1bdbb8bd1d1de62846b5f308ae3eac5127c7d36d7d5464985004122cc90",
	)

	t.Assert().Equal(OriginHash(expectedHash), ph)
}

func (t *calculatorHashesSuite) TestGetGlobuleHash() {
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
		GlobuleIndex:  0,
	}
	gh, _, err := t.calculator.GetGlobuleProof(globuleEntry)
	t.Assert().NoError(err)

	expectedHash, _ := hex.DecodeString(
		"42b62f955387238087acb76648d5f7a22e18c69de455ee1af50cd62573cca180861b3150c49fc2c5368f81b9acc604e28d72bd28ce1467720b83887d3580fadd",
	)

	fmt.Println(hex.EncodeToString(gh))

	t.Assert().Equal(OriginHash(expectedHash), gh)
}

func (t *calculatorHashesSuite) TestGetCloudHash() {
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
		GlobuleIndex:  0,
	}
	_, gp, err := t.calculator.GetGlobuleProof(globuleEntry)

	ch, _, err := t.calculator.GetCloudProof(&CloudEntry{
		ProofSet:      []*GlobuleProof{gp},
		PrevCloudHash: prevCloudHash,
	})

	t.Assert().NoError(err)

	expectedHash, _ := hex.DecodeString(
		"2d4022fedad75de096a5e978bca69ccf4fb836de3066be7ac15c1a07c6590f8b3e43cbdb9404b118974bfb764e3560648fe5e4969003bd5a1b50fb4e42b1399f",
	)

	fmt.Println(hex.EncodeToString(ch))

	t.Assert().Equal(OriginHash(expectedHash), ch)
}

type calculatorHashesSuite struct {
	suite.Suite

	pulse       *core.Pulse
	nodeNetwork core.NodeNetwork
	service     core.CryptographyService

	calculator Calculator
}

func TestCalculatorHashes(t *testing.T) {
	calculator := &calculator{}

	key, _ := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	assert.NotNil(t, key)

	service := testutils.NewCryptographyServiceMock(t)
	service.SignFunc = func(p []byte) (r *core.Signature, r1 error) {
		signature := core.SignatureFromBytes([]byte("signature"))
		return &signature, nil
	}
	service.GetPublicKeyFunc = func() (r crypto.PublicKey, r1 error) {
		return "key", nil
	}

	am := ledger.NewArtifactManagerMock(t)
	am.StateFunc = func() (r []byte, r1 error) {
		return []byte("state"), nil
	}

	scheme := platformpolicy.NewPlatformCryptographyScheme()
	nk := nodekeeper.GetTestNodekeeper(service)

	cm := component.Manager{}
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

	s := &calculatorHashesSuite{
		Suite:       suite.Suite{},
		calculator:  calculator,
		pulse:       pulse,
		nodeNetwork: nk,
		service:     service,
	}
	suite.Run(t, s)
}
