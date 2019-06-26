package gateway

import (
	"context"
	"errors"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"

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

	permit, err := g.authorize(ctx)
	if err != nil {
		// log warn
		g.Gatewayer.SwitchState(insolar.NoNetworkState)
	}

	// TODO: check authorize result and switch to JoinerBootstrap if other network is complete

	g.NodeKeeper.GetConsensusInfo().SetIsJoiner(false)

	pulse, err := g.PulseAccessor.Latest(ctx)
	if err != nil {
		pulse = insolar.Pulse{PulseNumber: 1}
	}

	resp, _ := g.BootstrapRequester.Bootstrap(ctx, permit, g.joinClaim, &pulse)

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

	// Authorize(utc) permit, check version
	// process response: trueAccept, redirect with permit, posibleAccept(regen shortId, updateScedule, update time utc)
	// check majority
	// handle reconect to other network
	// fake pulse

}

func (g *DiscoveryBootstrap) GetState() insolar.NetworkState {
	return insolar.DiscoveryBootstrap
}

func (g *DiscoveryBootstrap) authorize(ctx context.Context) (*packet.Permit, error) {
	cert := g.CertificateManager.GetCertificate()
	discoveryNodes := network.ExcludeOrigin(cert.GetDiscoveryNodes(), g.NodeKeeper.GetOrigin().ID())
	// todo: sort discoveryNodes

	for _, n := range discoveryNodes {
		if g.NodeKeeper.GetAccessor().GetActiveNode(*n.GetNodeRef()) != nil {
			inslogger.FromContext(ctx).Info("Skip discovery already in active list: ", n.GetNodeRef().String())
			continue
		}

		h, _ := host.NewHostN(n.GetHost(), *n.GetNodeRef())

		res, err := g.BootstrapRequester.Authorize(ctx, h, cert)
		if err != nil {
			inslogger.FromContext(ctx).Errorf("Error authorizing to discovery node %s: %s", h.String(), err.Error())
			continue
		}

		return res.Permit, nil
	}

	return nil, errors.New("Failed to authorize to any discovery node.")
}
