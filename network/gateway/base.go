//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package gateway

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log" // TODO remove before merge
	"github.com/insolar/insolar/network/gateway/bootstrap"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/platformpolicy"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
)

// Base is abstract class for gateways
type Base struct {
	component.Initer

	Self                  network.Gateway
	Gatewayer             network.Gatewayer               `inject:""`
	NodeKeeper            network.NodeKeeper              `inject:""`
	ContractRequester     insolar.ContractRequester       `inject:""`
	CryptographyService   insolar.CryptographyService     `inject:""`
	CertificateManager    insolar.CertificateManager      `inject:""`
	GIL                   insolar.GlobalInsolarLock       `inject:""`
	MessageBus            insolar.MessageBus              `inject:""`
	DiscoveryBootstrapper bootstrap.DiscoveryBootstrapper `inject:""`
	HostNetwork           network.HostNetwork             `inject:""`
	PulseAccessor         pulse.Accessor                  `inject:""`
}

// NewGateway creates new gateway on top of existing
func (g *Base) NewGateway(state insolar.NetworkState) network.Gateway {
	log.Infof("NewGateway %s", state.String())
	switch state {
	case insolar.NoNetworkState:
		g.Self = &NoNetwork{Base: g}
	case insolar.VoidNetworkState:
		g.Self = NewVoid(g)
	case insolar.JetlessNetworkState:
		g.Self = NewJetless(g)
	case insolar.AuthorizationNetworkState:
		g.Self = NewAuthorisation(g)
	case insolar.CompleteNetworkState:
		g.Self = NewComplete(g)
	default:
		panic("Try to switch network to unknown state. Memory of process is inconsistent.")
	}
	return g.Self
}

func (g *Base) Init(ctx context.Context) error {
	// TODO check in states
	g.HostNetwork.RegisterRequestHandler(types.Authorize, g.HandleNodeAuthorizeRequest) // validate cert
	g.HostNetwork.RegisterRequestHandler(types.Bootstrap, g.HandleNodeBootstrapRequest) // provide joiner claim

	return nil
}

func (g *Base) OnPulse(ctx context.Context, pu insolar.Pulse) error {
	if g.NodeKeeper == nil {
		return nil
	}
	// TODO switch state from another state
	// if g.NodeKeeper.IsBootstrapped() {
	// 	g.Gatewayer.SetGateway(g.Network.Gateway().NewGateway(insolar.CompleteNetworkState))
	// 	g.Gatewayer.Gateway().Run(ctx)
	// }
	return nil
}

// Auther casts us to Auther or obtain it in another way
func (g *Base) Auther() network.Auther {
	if ret, ok := g.Self.(network.Auther); ok {
		return ret
	}
	panic("Our network gateway suddenly is not an Auther")
}

// // Bootstrapper casts us to Bootstrapper or obtain it in another way
// func (g *Base) Bootstrapper() network.Bootstrapper {
// 	if ret, ok := g.Self.(network.Bootstrapper); ok {
// 		return ret
// 	}
// 	panic("Our network gateway suddenly is not an Bootstrapper")
// }

// GetCert method returns node certificate by requesting sign from discovery nodes
func (g *Base) GetCert(ctx context.Context, ref *insolar.Reference) (insolar.Certificate, error) {
	return nil, errors.New("GetCert() in non active mode")
}

// ValidateCert validates node certificate
func (g *Base) ValidateCert(ctx context.Context, certificate insolar.AuthorizationCertificate) (bool, error) {
	return false, errors.New("ValidateCert() in non active mode")
}

func (g *Base) FilterJoinerNodes(certificate insolar.Certificate, nodes []insolar.NetworkNode) []insolar.NetworkNode {
	dNodes := make(map[insolar.Reference]struct{}, len(certificate.GetDiscoveryNodes()))
	for _, dn := range certificate.GetDiscoveryNodes() {
		dNodes[*dn.GetNodeRef()] = struct{}{}
	}
	ret := []insolar.NetworkNode{}
	for _, n := range nodes {
		if _, ok := dNodes[n.ID()]; ok {
			ret = append(ret, n)
		}
	}
	return ret
}

// ============= Bootstrap =======

func (g *Base) ShoudIgnorePulse(context.Context, insolar.Pulse) bool {
	return false
}

// func (g *Base) HandleNodeRegisterRequest(ctx context.Context, request network.Packet) (network.Packet, error) {
// 	if request.GetRequest() == nil || request.GetRequest().GetRegister() == nil {
// 		return nil, errors.Errorf("process register: got invalid protobuf request message: %s", request)
// 	}
// 	data := request.GetRequest().GetRegister()
// 	if data.Version != g.NodeKeeper.GetOrigin().Version() {
// 		response := &packet.RegisterResponse{Code: packet.Denied,
// 			Error: fmt.Sprintf("Joiner version %s does not match discovery version %s",
// 				data.Version, g.NodeKeeper.GetOrigin().Version())}
// 		return g.HostNetwork.BuildResponse(ctx, request, response), nil
// 	}
// 	response := g.buildRegistrationResponse(bootstrap.SessionID(data.SessionID), data.JoinClaim)
// 	if response.Code != packet.Confirmed {
// 		return g.HostNetwork.BuildResponse(ctx, request, response), nil
// 	}
//
// 	// TODO: fix Short ID assignment logic
// 	if network.CheckShortIDCollision(g.NodeKeeper.GetAccessor().GetActiveNodes(), data.JoinClaim.ShortNodeID) {
// 		response = &packet.RegisterResponse{Code: packet.Denied,
// 			Error: "Short ID of the joiner node conflicts with active node short ID"}
// 		return g.HostNetwork.BuildResponse(ctx, request, response), nil
// 	}
//
// 	inslogger.FromContext(ctx).Infof("Added join claim from node %s", request.GetSender())
// 	g.NodeKeeper.GetClaimQueue().Push(data.JoinClaim)
// 	return g.HostNetwork.BuildResponse(ctx, request, response), nil
// }

// func (g *Base) buildRegistrationResponse(sessionID bootstrap.SessionID, claim *packets.NodeJoinClaim) *packet.RegisterResponse {
// 	session, err := g.getSession(sessionID, claim)
// 	if err != nil {
// 		return &packet.RegisterResponse{Code: packet.Denied, Error: err.Error()}
// 	}
// 	if node := g.NodeKeeper.GetAccessor().GetActiveNode(claim.NodeRef); node != nil {
// 		retryIn := session.TTL / 2
//
// 		keyProc := platformpolicy.NewKeyProcessor()
// 		// little hack: ignoring error, because it never fails in current implementation
// 		nodeKey, _ := keyProc.ExportPublicKeyBinary(node.PublicKey())
//
// 		log.Warnf("Joiner node (ID: %s, PK: %s) conflicts with node (ID: %s, PK: %s) in active list, sending request to reconnect in %v",
// 			claim.NodeRef, base58.Encode(claim.NodePK[:]), node.ID(), base58.Encode(nodeKey), retryIn)
//
// 		statsErr := stats.RecordWithTags(context.Background(), []tag.Mutator{
// 			tag.Upsert(tagNodeRef, claim.NodeRef.String()),
// 		}, statBootstrapReconnectRequired.M(1))
// 		if statsErr != nil {
// 			log.Warn("Failed to record reconnection retries metric: " + statsErr.Error())
// 		}
//
// 		g.SessionManager.ProlongateSession(sessionID, session)
// 		return &packet.RegisterResponse{Code: packet.Retry, RetryIn: int64(retryIn)}
// 	}
// 	return &packet.RegisterResponse{Code: packet.Confirmed}
// }

// func (g *Base) getSession(sessionID bootstrap.SessionID, claim *packets.NodeJoinClaim) (*bootstrap.Session, error) {
// 	session, err := g.SessionManager.ReleaseSession(sessionID)
// 	if err != nil {
// 		return nil, errors.Wrapf(err, "Error getting session %d for authorization", sessionID)
// 	}
// 	if !claim.NodeRef.Equal(session.NodeID) {
// 		return nil, errors.New("Claim node ID is not equal to session node ID")
// 	}
// 	// TODO: check claim signature
// 	return session, nil
// }

func (g *Base) HandleNodeBootstrapRequest(ctx context.Context, request network.Packet) (network.Packet, error) {
	if request.GetRequest() == nil || request.GetRequest().GetBootstrap() == nil {
		return nil, errors.Errorf("process bootstrap: got invalid protobuf request message: %s", request)
	}

	code := packet.Accepted
	data := request.GetRequest().GetBootstrap()
	var shortID insolar.ShortNodeID
	if network.CheckShortIDCollision(g.NodeKeeper.GetAccessor().GetActiveNodes(), data.JoinClaim.ShortNodeID) {
		shortID = network.GenerateUniqueShortID(g.NodeKeeper.GetAccessor().GetActiveNodes(), data.JoinClaim.GetNodeID())
	} else {
		shortID = data.JoinClaim.ShortNodeID
	}
	lastPulse, err := g.PulseAccessor.Latest(ctx)
	if err != nil {
		lastPulse = *insolar.GenesisPulse
	}

	var perm *packet.Permit
	// if g.nodeShouldReconnectAsJoiner(data.JoinClaim.NodeRef) {
	// 	code = packet.ReconnectRequired
	// } else if data.Permit == nil {
	// 	code = packet.Redirected
	// 	perm, err = g.CreatePermit(data.JoinClaim.NodePK[:])
	// 	if err != nil {
	// 		err = errors.Wrapf(err, "failed to generate permission")
	// 		return g.rejectBootstrapRequest(ctx, request, err.Error()), nil
	// 	}
	// } else {
	// 	err = g.checkPermission(data.Permit)
	// 	if err != nil {
	// 		err = errors.Wrapf(err, "failed to check permission")
	// 		return g.rejectBootstrapRequest(ctx, request, err.Error()), nil
	// 	}
	// }

	networkSize := uint32(len(g.NodeKeeper.GetAccessor().GetActiveNodes()))
	return g.HostNetwork.BuildResponse(ctx, request,
		&packet.BootstrapResponse{
			Code: code,
			// TODO: calculate ETA
			AssignShortID:    uint32(shortID),
			UpdateSincePulse: lastPulse.PulseNumber,
			NetworkSize:      networkSize,
		}), nil
}

func (bc *Base) nodeShouldReconnectAsJoiner(nodeID insolar.Reference) bool {
	// TODO:
	return bc.Gatewayer.Gateway().GetState() == insolar.CompleteNetworkState &&
		network.IsDiscovery(nodeID, bc.CertificateManager.GetCertificate())
}

func (bc *Base) rejectBootstrapRequest(ctx context.Context, request network.Packet, reason string) network.Packet {
	inslogger.FromContext(ctx).Errorf("Rejected bootstrap request from node %s: %s", request.GetSender(), reason)
	return bc.HostNetwork.BuildResponse(ctx, request, &packet.BootstrapResponse{Code: packet.Rejected, RejectReason: reason})
}
