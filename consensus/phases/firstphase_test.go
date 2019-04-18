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

package phases

import (
	"crypto"
	"testing"

	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/network/node"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/merkle"
	"github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFirstPhase_HandlePulse(t *testing.T) {
	firstPhase := &FirstPhaseImpl{}

	node := node.NewNode(insolar.Reference{}, insolar.StaticRoleUnknown, nil, "127.0.0.1:5432", "")
	nodeKeeper := nodenetwork.NewNodeKeeper(node)
	nodeKeeper.SetInitialSnapshot([]insolar.NetworkNode{node})

	pulseCalculatorMock := merkle.NewCalculatorMock(t)
	communicatorMock := NewCommunicatorMock(t)
	consensusNetworkMock := network.NewConsensusNetworkMock(t)
	terminationHandler := testutils.NewTerminationHandlerMock(t)
	messageBus := testutils.NewMessageBusLockerMock(t)
	cryptoServ := testutils.NewCryptographyServiceMock(t)
	cryptoServ.SignFunc = func(p []byte) (r *insolar.Signature, r1 error) {
		signature := insolar.SignatureFromBytes(nil)
		return &signature, nil
	}
	cryptoServ.VerifyFunc = func(p crypto.PublicKey, p1 insolar.Signature, p2 []byte) (r bool) {
		return true
	}

	cm := component.Manager{}
	cm.Inject(cryptoServ, nodeKeeper, firstPhase, pulseCalculatorMock, communicatorMock, consensusNetworkMock, terminationHandler, messageBus)

	require.NotNil(t, firstPhase.Calculator)
	require.NotNil(t, firstPhase.NodeKeeper)
	activeNodes := firstPhase.NodeKeeper.GetAccessor().GetActiveNodes()
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

func Test_getNodeState(t *testing.T) {
	n := node.NewNode(testutils.RandomRef(), insolar.StaticRoleVirtual, nil, "127.0.0.1:0", "")
	assert.Equal(t, packets.Legit, getNodeState(n, insolar.FirstPulseNumber))
	n.(node.MutableNode).SetState(insolar.NodeLeaving)
	n.(node.MutableNode).SetLeavingETA(insolar.FirstPulseNumber + 10)
	assert.Equal(t, packets.Legit, getNodeState(n, insolar.FirstPulseNumber))
	n.(node.MutableNode).SetLeavingETA(insolar.FirstPulseNumber - 10)
	assert.Equal(t, packets.TimedOut, getNodeState(n, insolar.FirstPulseNumber))
}
