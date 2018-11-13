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
	"testing"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/testutils/merkle"
	"github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/assert"
)

func TestFirstPhase_HandlePulse(t *testing.T) {
	firstPhase := &FirstPhase{}
	nodeNetworkMock := network.NewNodeNetworkMock(t)
	pulseCalculatorMock := merkle.NewCalculatorMock(t)
	communicatorMock := network.NewCommunicatorMock(t)
	consensusNetworkMock := network.NewConsensusNetworkMock(t)

	nodeNetworkMock.GetActiveNodesMock.Set(func() (r []core.Node) {
		return []core.Node{nodenetwork.NewNode(core.RecordRef{}, nil, nil, 0, "", "")}
	})

	cm := component.Manager{}
	cm.Register(nodeNetworkMock, firstPhase, pulseCalculatorMock, communicatorMock, consensusNetworkMock)

	assert.NotNil(t, firstPhase.Calculator)
	assert.NotNil(t, firstPhase.NodeNetwork)
	activeNodes := firstPhase.NodeNetwork.GetActiveNodes()
	assert.Equal(t, 1, len(activeNodes))
}
