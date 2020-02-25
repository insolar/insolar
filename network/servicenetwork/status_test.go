// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package servicenetwork

import (
	"context"
	"testing"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/version"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/network"

	testutils "github.com/insolar/insolar/testutils/network"
)

func TestGetNetworkStatus(t *testing.T) {
	sn := &ServiceNetwork{}
	gwer := testutils.NewGatewayerMock(t)
	gw := testutils.NewGatewayMock(t)
	ins := insolar.NetworkState(1)
	gw.GetStateMock.Set(func() insolar.NetworkState { return ins })
	gwer.GatewayMock.Set(func() network.Gateway { return gw })
	sn.Gatewayer = gwer

	pa := testutils.NewPulseAccessorMock(t)
	ppn := insolar.PulseNumber(2)
	pulse := insolar.Pulse{PulseNumber: 2}
	pa.GetLatestPulseMock.Set(func(context.Context) (insolar.Pulse, error) { return pulse, nil })
	sn.PulseAccessor = pa

	nk := testutils.NewNodeKeeperMock(t)
	a := testutils.NewAccessorMock(t)
	activeLen := 1
	active := make([]insolar.NetworkNode, activeLen)
	a.GetActiveNodesMock.Set(func() []insolar.NetworkNode { return active })

	workingLen := 2
	working := make([]insolar.NetworkNode, workingLen)
	a.GetWorkingNodesMock.Set(func() []insolar.NetworkNode { return working })

	nk.GetAccessorMock.Set(func(insolar.PulseNumber) network.Accessor { return a })

	nn := testutils.NewNetworkNodeMock(t)
	nk.GetOriginMock.Set(func() insolar.NetworkNode { return nn })

	sn.NodeKeeper = nk

	ns := sn.GetNetworkStatus()
	require.Equal(t, ins, ns.NetworkState)

	require.Equal(t, nn, ns.Origin)

	require.Equal(t, activeLen, ns.ActiveListSize)

	require.Equal(t, workingLen, ns.WorkingListSize)

	require.Len(t, ns.Nodes, activeLen)

	require.Equal(t, ppn, ns.Pulse.PulseNumber)

	require.Equal(t, version.Version, ns.Version)

	pa.GetLatestPulseMock.Set(func(context.Context) (insolar.Pulse, error) { return pulse, errors.New("test") })
	ns = sn.GetNetworkStatus()
	require.Equal(t, ins, ns.NetworkState)

	require.Equal(t, nn, ns.Origin)

	require.Equal(t, activeLen, ns.ActiveListSize)

	require.Equal(t, workingLen, ns.WorkingListSize)

	require.Len(t, ns.Nodes, activeLen)

	require.Equal(t, insolar.GenesisPulse.PulseNumber, ns.Pulse.PulseNumber)

	require.Equal(t, version.Version, ns.Version)
}
