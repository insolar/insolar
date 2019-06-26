package gateway

import (
	"context"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/hostnetwork/packet"
)

func newJoinerBootstrap(b *Base) *JoinerBootstrap {
	return &JoinerBootstrap{b}
}

// JoinerBootstrap void network state
type JoinerBootstrap struct {
	*Base
}

func (g *JoinerBootstrap) Run(ctx context.Context) {

	permit, err := g.authorize(ctx)
	if err != nil {
		// log warn
		g.Gatewayer.SwitchState(insolar.NoNetworkState)
	}

	g.NodeKeeper.GetConsensusInfo().SetIsJoiner(true)

	claim, _ := g.NodeKeeper.GetOriginJoinClaim()
	pulse, _ := g.PulseAccessor.Latest(ctx)

	resp, _ := g.BootstrapRequester.Bootstrap(ctx, permit, claim, &pulse)

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

	// todo bootstrap request
	// g.startBootstrap()
	// wait for consensus

	// ch := g.GetDiscoveryNodesChannel(ctx, discoveryNodes, 1)
	// result := bootstrap.WaitResultFromChannel(ctx, ch, time.Second) // TODO: options
	// if result == nil {
	// 	return nil, nil, errors.New("Failed to bootstrap to any of discovery nodes")
	// }
	// discovery := network.FindDiscoveryByRef(cert, result.Host.NodeID)
	// //return result, &bootstrap.DiscoveryNode{result.Host, discovery}, nil
	//
	// // if err != nil {
	// // 	// todo change state
	// // 	return nil, errors.Wrap(err, "Error bootstrapping to discovery node")
	// // }
	// return result, g.AuthenticateToDiscoveryNode(ctx, discovery)
}

func (g *JoinerBootstrap) GetState() insolar.NetworkState {
	return insolar.JoinerBootstrap
}

func (g *JoinerBootstrap) OnPulse(ctx context.Context, pu insolar.Pulse) error {
	return g.Base.OnPulse(ctx, pu)
}

func (g *JoinerBootstrap) ShouldIgnorePulse(context.Context, insolar.Pulse) bool {
	return false
}

func (bc *JoinerBootstrap) startBootstrap(ctx context.Context, perm *packet.Permit) (*packet.BootstrapResponse, error) {

	claim, err := bc.NodeKeeper.GetOriginJoinClaim()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get a join claim")
	}
	lastPulse, err := bc.PulseAccessor.Latest(ctx)
	if err != nil {
		lastPulse = *insolar.GenesisPulse
	}

	// BOOTSTRAP request --------
	resp, err := bc.BootstrapRequester.Bootstrap(ctx, perm, claim, &lastPulse)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to Bootstrap: %s", err.Error())
	}

	// TODO:

	return resp, nil
}

func (g *JoinerBootstrap) authorize(ctx context.Context) (*packet.Permit, error) {
	cert := g.CertificateManager.GetCertificate()
	discoveryNodes := network.ExcludeOrigin(cert.GetDiscoveryNodes(), g.NodeKeeper.GetOrigin().ID())
	// todo: shuffle discoveryNodes

	for _, n := range discoveryNodes {
		h, _ := host.NewHostN(n.GetHost(), *n.GetNodeRef())

		res, err := g.BootstrapRequester.Authorize(ctx, h, cert)
		if err != nil {
			inslogger.FromContext(ctx).Errorf("Error authorizing to host %s: %s", h.String(), err.Error())
			continue
		}
		// TODO: check majority and res.NetworkState

		return res.Permit, nil
	}

	return nil, errors.New("Failed to authorize to any discovery node.")
}
