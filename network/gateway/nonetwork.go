package gateway

// TODO: spans, metrics

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
)

// NewNoNetwork this initial constructor have special signature to be called outside
func NewNoNetwork(n network.Gatewayer, mb messageBusLocker) *NoNetwork {
	return &NoNetwork{c: &commons{Network: n, MBLocker: mb}}
}

// NoNetwork initial state
type NoNetwork struct {
	c *commons
}

func (g *NoNetwork) Run() {

}

func (g *NoNetwork) GetState() insolar.NetworkState {
	return insolar.NoNetworkState
}

func (g *NoNetwork) OnPulse(context.Context, insolar.Pulse) error {
	panic("oops")
}
