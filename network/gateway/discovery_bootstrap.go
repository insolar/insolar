package gateway

import (
	"context"

	"github.com/insolar/insolar/insolar"
)

func newDiscoveryBootstrap(b *Base) *DiscoveryBootstrap {
	return &DiscoveryBootstrap{b}
}

// DiscoveryBootstrap void network state
type DiscoveryBootstrap struct {
	*Base
}

func (g *DiscoveryBootstrap) Run(ctx context.Context) {
	// var err error

	// cert := g.CertificateManager.GetCertificate()

	// TODO: shaffle discovery nodes

	// ping ?
	// Authorize(utc) permit, check version
	// process response: trueAccept, redirect with permit, posibleAccept(regen shortId, updateScedule, update time utc)
	// check majority
	// handle reconect to other network
	// fake pulse

	// if network.OriginIsDiscovery(cert) {
	// 	_, err = g.Bootstrapper.BootstrapDiscovery(ctx)
	// 	// if the network is up and complete, we return discovery nodes via consensus
	// 	if err == bootstrap.ErrReconnectRequired {
	// 		log.Debugf("[ Bootstrap ] Connecting discovery node %s as joiner", g.NodeKeeper.GetOrigin().ID())
	// 		g.NodeKeeper.GetOrigin().(node.MutableNode).SetState(insolar.NodePending)
	// 		g.bootstrapJoiner(ctx)
	// 		return
	// 	}
	//
	// }
}

func (g *DiscoveryBootstrap) GetState() insolar.NetworkState {
	return insolar.DiscoveryBootstrap
}

func (g *DiscoveryBootstrap) OnPulse(ctx context.Context, pu insolar.Pulse) error {
	return g.Base.OnPulse(ctx, pu)
}

func (g *DiscoveryBootstrap) ShoudIgnorePulse(context.Context, insolar.Pulse) bool {
	return false
}

func (g *DiscoveryBootstrap) bootstrapJoiner(ctx context.Context) {
	g.Gatewayer.SwitchState(insolar.JoinerBootstrap)
}
