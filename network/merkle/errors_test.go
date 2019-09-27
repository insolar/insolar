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
	"crypto"
	"encoding/hex"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar/pulsartestutils"
	"github.com/insolar/insolar/testutils"
	network2 "github.com/insolar/insolar/testutils/network"
)

type calculatorErrorSuite struct {
	suite.Suite

	pulse          *insolar.Pulse
	originProvider network.OriginProvider
	service        insolar.CryptographyService

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
		ProofSet: map[insolar.NetworkNode]*PulseProof{
			t.originProvider.GetOrigin(): {},
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
	calculator := &calculator{}

	cm := component.Manager{}

	key, _ := platformpolicy.NewKeyProcessor().GeneratePrivateKey()
	require.NotNil(t, key)

	service := testutils.NewCryptographyServiceMock(t)
	service.SignMock.Set(func(p []byte) (r *insolar.Signature, r1 error) {
		return nil, errors.New("Sign error")
	})
	service.GetPublicKeyMock.Set(func() (r crypto.PublicKey, r1 error) {
		return "key", nil
	})
	scheme := platformpolicy.NewPlatformCryptographyScheme()

	ps := pulse.NewStorageMem()

	op := network2.NewOriginProviderMock(t)
	op.GetOriginMock.Set(func() insolar.NetworkNode {
		return createOrigin()
	})

	th := testutils.NewTerminationHandlerMock(t)

	am := staterMock{
		stateFunc: func() []byte {
			return []byte{1, 2, 3}
		},
	}
	jc := jet.NewCoordinatorMock(t)

	cm.Inject(th, op, jc, &am, calculator, service, scheme, ps)

	require.NotNil(t, calculator.Stater)
	require.NotNil(t, calculator.OriginProvider)
	require.NotNil(t, calculator.CryptographyService)
	require.NotNil(t, calculator.PlatformCryptographyScheme)

	err := cm.Init(context.Background())
	require.NoError(t, err)

	pulseObject := &insolar.Pulse{
		PulseNumber:     insolar.PulseNumber(1337),
		NextPulseNumber: insolar.PulseNumber(1347),
		Entropy:         pulsartestutils.MockEntropyGenerator{}.GenerateEntropy(),
	}

	s := &calculatorErrorSuite{
		Suite:          suite.Suite{},
		calculator:     calculator,
		pulse:          pulseObject,
		originProvider: op,
		service:        service,
	}
	suite.Run(t, s)
}

type staterMock struct {
	stateFunc func() []byte
}

func (m staterMock) State() []byte {
	return m.stateFunc()
}
