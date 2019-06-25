package gateway

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/instracer"
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
	ctx, span := instracer.StartSpan(ctx, "NetworkBootstrapper.bootstrapJoiner")
	defer span.End()
	g.NodeKeeper.GetConsensusInfo().SetIsJoiner(true)

	if g.permit == nil {
		// log warn
		g.Gatewayer.SwitchState(insolar.NoNetworkState)
	}

	claim, _ := g.NodeKeeper.GetOriginJoinClaim()
	pulse, _ := g.PulseAccessor.Latest(ctx)

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

func (g *JoinerBootstrap) ShoudIgnorePulse(context.Context, insolar.Pulse) bool {
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
