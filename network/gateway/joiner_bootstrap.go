package gateway

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/gateway/bootstrap"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
)

const (
	registrationRetries = 20
)

func NewJoinerBootstrap(b *Base) *JoinerBootstrap {
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

	cert := g.CertificateManager.GetCertificate()
	discoveryNodes := cert.GetDiscoveryNodes()
	// TODO: shuffle nodes

	if network.OriginIsDiscovery(cert) {
		discoveryNodes = network.ExcludeOrigin(discoveryNodes, g.NodeKeeper.GetOrigin().ID())
	}
	if len(discoveryNodes) == 0 {
		panic("ZeroBootstrap. There are 0 discovery nodes to connect to")
	}

	d := discoveryNodes[0]
	host, _ := host.NewHostN(d.GetHost(), *d.GetNodeRef())
	discovery := &bootstrap.DiscoveryNode{host, d}

	data, err := g.Authorize(ctx, discovery, g.CertificateManager.GetCertificate())
	if err != nil {
		panic("Error authorizing on discovery node")
	}

	// todo bootstrap request
	g.startBootstrap()
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

func (bc *JoinerBootstrap) startBootstrap(ctx context.Context, perm *packet.Permission) (*network.BootstrapResult, error) {

	h, _ := host.NewHostN(perm.Payload.ReconnectTo, perm.Payload.DiscoveryRef)

	claim, err := bc.NodeKeeper.GetOriginJoinClaim()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get a join claim")
	}
	lastPulse, err := bc.PulseAccessor.Latest(ctx)
	if err != nil {
		lastPulse = *insolar.GenesisPulse
	}

	request := &packet.BootstrapRequest{
		JoinClaim:     claim,
		LastNodePulse: lastPulse.PulseNumber,
		Permission:    perm,
	}

	// BOOTSTRAP request --------
	future, err := bc.HostNetwork.SendRequestToHost(ctx, types.Bootstrap, request, h)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to send bootstrap request to %s", h.String())
	}

	response, err := future.WaitResponse(bc.options.BootstrapTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get response to bootstrap request from %s", h.String())
	}
	if response.GetResponse() == nil || response.GetResponse().GetBootstrap() == nil {
		return nil, errors.Errorf("Failed to get response to bootstrap request from address %s: "+
			"got incorrect response: %s", h.String(), response)
	}
}

// func (ac *JoinerBootstrap) register(ctx context.Context, discoveryNode *bootstrap.DiscoveryNode,
// 	sessionID bootstrap.SessionID, attempt int) error {
//
// 	if attempt == 0 {
// 		inslogger.FromContext(ctx).Infof("Registering on host: %s", discoveryNode.Host)
// 	} else {
// 		inslogger.FromContext(ctx).Infof("Registering on host: %s; attempt: %d", discoveryNode.Host, attempt+1)
// 	}
//
// 	ctx, span := instracer.StartSpan(ctx, "AuthorizationController.Register")
// 	span.AddAttributes(
// 		trace.StringAttribute("node", discoveryNode.Node.GetNodeRef().String()),
// 	)
// 	defer span.End()
// 	originClaim, err := ac.NodeKeeper.GetOriginJoinClaim()
// 	if err != nil {
// 		return errors.Wrap(err, "Failed to get origin claim")
// 	}
// 	request := &packet.RegisterRequest{
// 		Version:   ac.NodeKeeper.GetOrigin().Version(),
// 		SessionID: uint64(sessionID),
// 		JoinClaim: originClaim,
// 	}
// 	future, err := ac.HostNetwork.SendRequestToHost(ctx, types.Register, request, discoveryNode.Host)
// 	if err != nil {
// 		return errors.Wrapf(err, "Error sending register request")
// 	}
// 	response, err := future.WaitResponse(ac.options.PacketTimeout)
// 	if err != nil {
// 		return errors.Wrapf(err, "Error getting response for register request")
// 	}
// 	if response.GetResponse() == nil || response.GetResponse().GetRegister() == nil {
// 		return errors.Errorf("Register failed: got incorrect response: %s", response)
// 	}
// 	data := response.GetResponse().GetRegister()
// 	if data.Code == packet.Denied {
// 		return errors.New("Register rejected: " + data.Error)
// 	}
// 	if data.Code == packet.Retry {
// 		if attempt >= registrationRetries {
// 			return errors.Errorf("Exceeded maximum number of registration retries (%d)", registrationRetries)
// 		}
// 		log.Warnf("Failed to register on discovery node %s. Reason: node %s is already in network active list. "+
// 			"Retrying registration in %v", discoveryNode.Host, ac.NodeKeeper.GetOrigin().ID(), data.RetryIn)
// 		time.Sleep(time.Duration(data.RetryIn))
// 		return ac.register(ctx, discoveryNode, sessionID, attempt+1)
// 	}
// 	return nil
// }
