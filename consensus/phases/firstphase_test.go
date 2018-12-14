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
)

func TestFirstPhase_HandlePulse(t *testing.T) {
	firstPhase := &FirstPhase{}
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

	cm := component.Manager{}
	cm.Inject(cryptoServ, nodeKeeperMock, firstPhase, pulseCalculatorMock, communicatorMock, consensusNetworkMock)

	assert.NotNil(t, firstPhase.Calculator)
	assert.NotNil(t, firstPhase.NodeKeeper)
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
