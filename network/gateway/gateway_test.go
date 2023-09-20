package gateway

import (
	"context"
	"testing"

	"github.com/insolar/insolar/certificate"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/reply"

	"github.com/insolar/insolar/network"
	testnet "github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/testutils"

	"github.com/insolar/insolar/insolar"
)

func emtygateway(t *testing.T) network.Gateway {
	// todo use mockPulseManager(t)
	return newNoNetwork(&Base{})
}

func TestSwitch(t *testing.T) {
	t.Skip("fixme")
	ctx := context.Background()

	// nodekeeper := testnet.NewNodeKeeperMock(t)
	nodekeeper := testnet.NewNodeKeeperMock(t)
	nodekeeper.MoveSyncToActiveMock.Set(func(ctx context.Context, number insolar.PulseNumber) {})
	gatewayer := testnet.NewGatewayerMock(t)
	// pm := mockPulseManager(t)

	ge := emtygateway(t)

	require.NotNil(t, ge)
	require.Equal(t, "NoNetworkState", ge.GetState().String())

	ge.Run(ctx, *insolar.EphemeralPulse)

	gatewayer.GatewayMock.Set(func() (g1 network.Gateway) {
		return ge
	})
	gatewayer.SwitchStateMock.Set(func(ctx context.Context, state insolar.NetworkState, pulse insolar.Pulse) {
		ge = ge.NewGateway(ctx, state)
	})
	gilreleased := false

	ge.OnPulseFromPulsar(ctx, insolar.Pulse{}, nil)

	require.Equal(t, "CompleteNetworkState", ge.GetState().String())
	require.False(t, gilreleased)
	cref := gen.Reference()

	for _, state := range []insolar.NetworkState{insolar.NoNetworkState,
		insolar.JoinerBootstrap, insolar.CompleteNetworkState} {
		ge = ge.NewGateway(ctx, state)
		require.Equal(t, state, ge.GetState())
		ge.Run(ctx, *insolar.EphemeralPulse)
		au := ge.Auther()

		_, err := au.GetCert(ctx, &cref)
		require.Error(t, err)

		_, err = au.ValidateCert(ctx, &certificate.Certificate{})
		require.Error(t, err)

		ge.OnPulseFromPulsar(ctx, insolar.Pulse{}, nil)

	}

}

func TestDumbComplete_GetCert(t *testing.T) {
	t.Skip("fixme")
	ctx := context.Background()

	// nodekeeper := testnet.NewNodeKeeperMock(t)
	nodekeeper := testnet.NewNodeKeeperMock(t)
	nodekeeper.MoveSyncToActiveMock.Set(func(ctx context.Context, number insolar.PulseNumber) {})

	gatewayer := testnet.NewGatewayerMock(t)

	CR := testutils.NewContractRequesterMock(t)
	CM := testutils.NewCertificateManagerMock(t)
	ge := emtygateway(t)
	// pm := mockPulseManager(t)

	// ge := newNoNetwork(gatewayer, pm,
	//	nodekeeper, CR,
	//	testutils.NewCryptographyServiceMock(t),
	//	testnet.NewHostNetworkMock(t),
	//	CM)

	require.NotNil(t, ge)
	require.Equal(t, "NoNetworkState", ge.GetState().String())

	ge.Run(ctx, *insolar.EphemeralPulse)

	gatewayer.GatewayMock.Set(func() (r network.Gateway) { return ge })
	gatewayer.SwitchStateMock.Set(func(ctx context.Context, state insolar.NetworkState, pulse insolar.Pulse) {
		ge = ge.NewGateway(ctx, state)
	})
	gilreleased := false

	ge.OnPulseFromPulsar(ctx, insolar.Pulse{}, nil)

	require.Equal(t, "CompleteNetworkState", ge.GetState().String())
	require.False(t, gilreleased)

	cref := gen.Reference()

	CR.CallMock.Set(func(ctx context.Context, ref *insolar.Reference, method string, argsIn []interface{}, p insolar.PulseNumber,
	) (r insolar.Reply, r2 *insolar.Reference, r1 error) {
		require.Equal(t, &cref, ref)
		require.Equal(t, "GetNodeInfo", method)
		repl, _ := insolar.Serialize(struct {
			PublicKey string
			Role      insolar.StaticRole
		}{"LALALA", insolar.StaticRoleVirtual})
		return &reply.CallMethod{
			Result: repl,
		}, nil, nil
	})

	CM.GetCertificateMock.Set(func() (r insolar.Certificate) { return &certificate.Certificate{} })
	cert, err := ge.Auther().GetCert(ctx, &cref)

	require.NoError(t, err)
	require.NotNil(t, cert)
	require.Equal(t, cert, &certificate.Certificate{})
}
