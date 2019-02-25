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

package phases

import (
	"crypto"
	"testing"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/merkle"
	"github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFirstPhase_HandlePulse(t *testing.T) {
	firstPhase := &FirstPhaseImpl{}
	nodeKeeperMock := network.NewNodeKeeperMock(t)
	pulseCalculatorMock := merkle.NewCalculatorMock(t)
	communicatorMock := NewCommunicatorMock(t)
	consensusNetworkMock := network.NewConsensusNetworkMock(t)

	cryptoServ := testutils.NewCryptographyServiceMock(t)
	cryptoServ.SignFunc = func(p []byte) (r *core.Signature, r1 error) {
		signature := core.SignatureFromBytes(nil)
		return &signature, nil
	}
	cryptoServ.VerifyFunc = func(p crypto.PublicKey, p1 core.Signature, p2 []byte) (r bool) {
		return true
	}

	nodeKeeperMock.GetActiveNodesMock.Set(func() (r []core.Node) {
		return []core.Node{nodenetwork.NewNode(core.RecordRef{}, core.StaticRoleUnknown, nil, "127.0.0.1:5432", "")}

	})

	networkSwitcherMock := testutils.NewNetworkSwitcherMock(t)
	certificateManagerMock := testutils.NewCertificateManagerMock(t)

	cm := component.Manager{}
	cm.Inject(
		cryptoServ,
		nodeKeeperMock,
		firstPhase,
		pulseCalculatorMock,
		communicatorMock,
		consensusNetworkMock,
		networkSwitcherMock,
		certificateManagerMock,
	)

	require.NotNil(t, firstPhase.Calculator)
	require.NotNil(t, firstPhase.NodeKeeper)
	activeNodes := firstPhase.NodeKeeper.GetActiveNodes()
	assert.Equal(t, 1, len(activeNodes))
}

func Test_consensusReached(t *testing.T) {
	assert.True(t, consensusReachedBFT(5, 6))
	assert.False(t, consensusReachedBFT(4, 6))

	assert.True(t, consensusReachedBFT(201, 300))
	assert.False(t, consensusReachedBFT(200, 300))

	assert.True(t, consensusReachedMajority(4, 6))
	assert.False(t, consensusReachedMajority(3, 6))

	assert.True(t, consensusReachedMajority(151, 300))
	assert.False(t, consensusReachedMajority(150, 300))
}
