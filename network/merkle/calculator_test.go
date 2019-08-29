//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package merkle

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/node"
	network2 "github.com/insolar/insolar/testutils/network"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar/pulsartestutils"
	"github.com/insolar/insolar/testutils"
)

func createOrigin() insolar.NetworkNode {
	ref, _ := insolar.NewReferenceFromBase58("14K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.17ZQboaH24PH42sqZKUvoa7UBrpuuubRtShp6CKNuWGZa")
	return node.NewNode(*ref, insolar.StaticRoleVirtual, nil, "127.0.0.1:5432", "")
}

type calculatorSuite struct {
	suite.Suite

	pulse          *insolar.Pulse
	originProvider network.OriginProvider
	service        insolar.CryptographyService

	calculator Calculator
}

func (t *calculatorSuite) TestGetNodeProof() {
	ph, np, err := t.calculator.GetPulseProof(&PulseEntry{Pulse: t.pulse})

	t.Assert().NoError(err)
	t.Assert().NotNil(np)

	key, err := t.service.GetPublicKey()
	t.Assert().NoError(err)

	t.Assert().True(t.calculator.IsValid(np, ph, key))
}

func (t *calculatorSuite) TestGetGlobuleProof() {
	pulseEntry := &PulseEntry{Pulse: t.pulse}
	ph, pp, err := t.calculator.GetPulseProof(pulseEntry)
	t.Assert().NoError(err)

	prevCloudHash, _ := hex.DecodeString(
		"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	)

	globuleEntry := &GlobuleEntry{
		PulseEntry: pulseEntry,
		PulseHash:  ph,
		ProofSet: map[insolar.NetworkNode]*PulseProof{
			t.originProvider.GetOrigin(): pp,
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
	pulseEntry := &PulseEntry{Pulse: t.pulse}
	ph, pp, err := t.calculator.GetPulseProof(pulseEntry)
	t.Assert().NoError(err)

	prevCloudHash, _ := hex.DecodeString(
		"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	)

	globuleEntry := &GlobuleEntry{
		PulseEntry: pulseEntry,
		PulseHash:  ph,
		ProofSet: map[insolar.NetworkNode]*PulseProof{
			t.originProvider.GetOrigin(): pp,
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
	op := network2.NewOriginProviderMock(t)
	op.GetOriginMock.Set(func() insolar.NetworkNode {
		return createOrigin()
	})

	th := testutils.NewTerminationHandlerMock(t)
	am := staterMock{
		stateFunc: func() []byte {
			return []byte("state")
		},
	}

	cm := component.Manager{}
	cm.Inject(th, op, &am, calculator, service, scheme)

	require.NotNil(t, calculator.Stater)
	require.NotNil(t, calculator.OriginProvider)
	require.NotNil(t, calculator.CryptographyService)
	require.NotNil(t, calculator.PlatformCryptographyScheme)

	err := cm.Init(context.Background())
	require.NoError(t, err)

	pulse := &insolar.Pulse{
		PulseNumber:     insolar.PulseNumber(1337),
		NextPulseNumber: insolar.PulseNumber(1347),
		Entropy:         pulsartestutils.MockEntropyGenerator{}.GenerateEntropy(),
	}

	s := &calculatorSuite{
		Suite:          suite.Suite{},
		calculator:     calculator,
		pulse:          pulse,
		originProvider: op,
		service:        service,
	}
	suite.Run(t, s)
}
