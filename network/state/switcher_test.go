/*
 *    Copyright 2018 INS Ecosystem
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

package state

import (
	"testing"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/require"
)

func mockNodeNetwork(t *testing.T) *network.NodeNetworkMock {
	nnMock := network.NewNodeNetworkMock(t)
	return nnMock
}

func mockSwitcherWorkAround(t *testing.T, isBootstrapped bool) *network.SwitcherWorkAroundMock {
	swaMock := network.NewSwitcherWorkAroundMock(t)
	swaMock.IsBootstrappedFunc = func() bool {
		return isBootstrapped
	}
	return swaMock
}

func TestNewNetworkSwitcher(t *testing.T) {
	nodeNet := network.NewNodeNetworkMock(t)
	switcherWorkAround := mockSwitcherWorkAround(t, false)

	switcher, err := NewNetworkSwitcher()
	require.NoError(t, err)

	cm := &component.Manager{}
	cm.Inject(nodeNet, switcherWorkAround, switcher)

	require.Equal(t, nodeNet, switcher.NodeNetwork)
	require.Equal(t, switcherWorkAround, switcher.SwitcherWorkAround)
}
