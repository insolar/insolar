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

package networkcoordinator

import (
	"context"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

func TestNewNetworkCoordinator(t *testing.T) {
	nc, err := New()
	require.NoError(t, err)
	require.Equal(t, &NetworkCoordinator{}, nc)
}

func TestNetworkCoordinator_Start(t *testing.T) {
	nc, err := New()
	require.NoError(t, err)
	mb := testutils.NewMessageBusMock(t)
	mb.MustRegisterFunc = func(p core.MessageType, handler core.MessageHandler) {
		require.Equal(t, p, core.NetworkCoordinatorNodeSignRequest)
	}
	nc.MessageBus = mb
	ctx := context.Background()
	err = nc.Start(ctx)
	require.NoError(t, err)
}

func TestNetworkCoordinator_GetCoordinator_Zero(t *testing.T) {
	nc, err := New()
	require.NoError(t, err)
	ns := testutils.NewNetworkSwitcherMock(t)
	ns.GetStateFunc = func() core.NetworkState {
		return core.NoNetworkState
	}
	nc.NetworkSwitcher = ns
	crd := nc.getCoordinator()
	require.Equal(t, nc.zeroCoordinator, crd)
}

func TestNetworkCoordinator_GetCoordinator_Real(t *testing.T) {
	nc, err := New()
	require.NoError(t, err)
	ns := testutils.NewNetworkSwitcherMock(t)
	ns.GetStateFunc = func() core.NetworkState {
		return core.CompleteNetworkState
	}
	nc.NetworkSwitcher = ns
	crd := nc.getCoordinator()
	require.Equal(t, nc.realCoordinator, crd)
}

func test_getNode() *reply.CallMethod {
	node, err := core.MarshalArgs(struct {
		PublicKey string
		Role      core.StaticRole
	}{
		PublicKey: "test_node_public_key",
		Role:      core.StaticRoleVirtual,
	}, nil)
	if err != nil {
		log.Fatal(err)
	}
	return &reply.CallMethod{
		Result: []byte(node),
	}
}
