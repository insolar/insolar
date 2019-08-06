package gateway

import (
	"context"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
	mock "github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewGatewayer(t *testing.T) {
	t.Skip("fixme")
	gw := mock.NewGatewayMock(t)
	//gw.NeedLockMessageBusMock.Inspect() Expect().Return(false)

	gw.GetStateMock.Set(func() (n1 insolar.NetworkState) {
		return insolar.NoNetworkState
	})

	gw.NewGatewayMock.Set(func(ctx context.Context, s insolar.NetworkState) (g1 network.Gateway) {
		assert.Equal(t, insolar.WaitConsensus, s)
		return gw
	})

	gatewayer := NewGatewayer(gw, func(ctx context.Context, isNetworkOperable bool) {})
	assert.Equal(t, gw, gatewayer.Gateway())
	assert.Equal(t, insolar.NoNetworkState, gatewayer.Gateway().GetState())

	gatewayer.SwitchState(context.Background(), insolar.WaitConsensus)
}
