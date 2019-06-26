package gateway

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/hostnetwork/packet"
)

func newDiscoveryBootstrap(b *Base) *DiscoveryBootstrap {
	return &DiscoveryBootstrap{b}
}

// DiscoveryBootstrap void network state
type DiscoveryBootstrap struct {
	*Base
}

func (g *DiscoveryBootstrap) Run(ctx context.Context) {
	g.NodeKeeper.GetConsensusInfo().SetIsJoiner(false)

	if g.permit == nil {
		// log warn
		g.Gatewayer.SwitchState(insolar.NoNetworkState)
	}

	claim, _ := g.NodeKeeper.GetOriginJoinClaim()
	pulse, err := g.PulseAccessor.Latest(ctx)
	if err != nil {
		pulse = insolar.Pulse{PulseNumber: 1}
	}

	resp, _ := g.BootstrapRequester.Bootstrap(ctx, g.permit, claim, &pulse)

	if resp.Code == packet.Reject {
		g.Gatewayer.SwitchState(insolar.NoNetworkState)
		return
	}

	if resp.Code == packet.Accepted {
		//  ConsensusWaiting, ETA
		g.bootstrapETA = insolar.PulseNumber(resp.ETA)
		g.Gatewayer.SwitchState(insolar.WaitConsensus)
		return
	}

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
