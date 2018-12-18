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
	"context"
	"sync"
	"testing"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/require"
)

func mockSwitcherWorkAround(t *testing.T, isBootstrapped bool) *network.SwitcherWorkAroundMock {
	swaMock := network.NewSwitcherWorkAroundMock(t)
	swaMock.IsBootstrappedFunc = func() bool {
		return isBootstrapped
	}
	return swaMock
}

func mockMessageBusLocker(t *testing.T) *messageBusLockerMock {
	mblMock := NewmessageBusLockerMock(t)
	mblMock.UnlockFunc = func(p context.Context) {}
	return mblMock
}

func TestNewNetworkSwitcher(t *testing.T) {
	nodeNet := network.NewNodeNetworkMock(t)
	switcherWorkAround := mockSwitcherWorkAround(t, false)
	messageBusLocker := mockMessageBusLocker(t)

	switcher, err := NewNetworkSwitcher()
	require.NoError(t, err)

	cm := &component.Manager{}
	cm.Inject(nodeNet, switcherWorkAround, messageBusLocker, switcher)

	require.Equal(t, nodeNet, switcher.NodeNetwork)
	require.Equal(t, switcherWorkAround, switcher.SwitcherWorkAround)
	require.Equal(t, messageBusLocker, switcher.MBLocker)
	require.Equal(t, core.NoNetworkState, switcher.state)
	require.Equal(t, sync.RWMutex{}, switcher.stateLock)
}

func TestGetState(t *testing.T) {
	switcher, err := NewNetworkSwitcher()
	require.NoError(t, err)

	state := switcher.GetState()
	require.Equal(t, core.NoNetworkState, state)
}

func TestOnPulseNoChange(t *testing.T) {
	switcher, err := NewNetworkSwitcher()
	require.NoError(t, err)
	switcherWorkAround := mockSwitcherWorkAround(t, false)
	nodeNet := network.NewNodeNetworkMock(t)
	messageBusLocker := mockMessageBusLocker(t)

	cm := &component.Manager{}
	cm.Inject(switcherWorkAround, switcher, nodeNet, messageBusLocker)

	err = switcher.OnPulse(context.Background(), core.Pulse{})
	require.NoError(t, err)
	require.Equal(t, core.NoNetworkState, switcher.state)
	require.Equal(t, uint64(0), messageBusLocker.UnlockCounter)
}

func TestOnPulseStateChanged(t *testing.T) {
	switcher, err := NewNetworkSwitcher()
	require.NoError(t, err)
	switcherWorkAround := mockSwitcherWorkAround(t, true)
	nodeNet := network.NewNodeNetworkMock(t)
	messageBusLocker := mockMessageBusLocker(t)

	cm := &component.Manager{}
	cm.Inject(switcherWorkAround, switcher, nodeNet, messageBusLocker)

	err = switcher.OnPulse(context.Background(), core.Pulse{})
	require.NoError(t, err)
	require.Equal(t, core.CompleteNetworkState, switcher.state)
	require.Equal(t, uint64(1), messageBusLocker.UnlockCounter)
}

func TestGetStateAfterStateChanged(t *testing.T) {
	switcher, err := NewNetworkSwitcher()
	require.NoError(t, err)
	switcherWorkAround := mockSwitcherWorkAround(t, true)
	nodeNet := network.NewNodeNetworkMock(t)
	messageBusLocker := mockMessageBusLocker(t)

	cm := &component.Manager{}
	cm.Inject(switcherWorkAround, switcher, nodeNet, messageBusLocker)

	err = switcher.OnPulse(context.Background(), core.Pulse{})
	require.NoError(t, err)
	require.Equal(t, core.CompleteNetworkState, switcher.state)

	state := switcher.GetState()
	require.Equal(t, core.CompleteNetworkState, state)
}
