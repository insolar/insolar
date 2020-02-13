// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package gateway

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/network/rules"
	"github.com/insolar/insolar/network/storage"
	"github.com/insolar/insolar/network/transport"

	"github.com/pkg/errors"

	component "github.com/insolar/component-manager"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/gateway/bootstrap"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/platformpolicy"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
)

const (
	bootstrapTimeoutMessage = "Bootstrap timeout exceeded"
)

// Base is abstract class for gateways
type Base struct {
	component.Initer

	Self                network.Gateway
	Gatewayer           network.Gatewayer                  `inject:""`
	NodeKeeper          network.NodeKeeper                 `inject:""`
	ContractRequester   insolar.ContractRequester          `inject:""`
	CryptographyService insolar.CryptographyService        `inject:""`
	CryptographyScheme  insolar.PlatformCryptographyScheme `inject:""`
	CertificateManager  insolar.CertificateManager         `inject:""`
	HostNetwork         network.HostNetwork                `inject:""`
	PulseAccessor       storage.PulseAccessor              `inject:""`
	PulseAppender       storage.PulseAppender              `inject:""`
	PulseManager        insolar.PulseManager               `inject:""`
	BootstrapRequester  bootstrap.Requester                `inject:""`
	KeyProcessor        insolar.KeyProcessor               `inject:""`
	Aborter             network.Aborter                    `inject:""`
	TransportFactory    transport.Factory                  `inject:""`

	// nolint
	OriginProvider network.OriginProvider `inject:""`

	datagramHandler   *adapters.DatagramHandler
	datagramTransport transport.DatagramTransport

	ConsensusMode         consensus.Mode
	consensusInstaller    consensus.Installer
	ConsensusController   consensus.Controller
	consensusPulseHandler network.PulseHandler
	consensusStarted      uint32

	Options         *network.Options
	bootstrapTimer  *time.Timer // nolint
	bootstrapETA    time.Duration
	originCandidate *adapters.Candidate

	// Next request backoff.
	backoff time.Duration // nolint

	pulseWatchdog *pulseWatchdog
}

// NewGateway creates new gateway on top of existing
func (g *Base) NewGateway(ctx context.Context, state insolar.NetworkState) network.Gateway {
	inslogger.FromContext(ctx).Infof("NewGateway %s", state.String())
	switch state {
	case insolar.NoNetworkState:
		g.Self = newNoNetwork(g)
	case insolar.CompleteNetworkState:
		g.Self = newComplete(g)
	case insolar.JoinerBootstrap:
		g.Self = newJoinerBootstrap(g)
	case insolar.WaitConsensus:
		err := g.StartConsensus(ctx)
		if err != nil {
			g.FailState(ctx, fmt.Sprintf("Failed to start consensus: %s", err))
		}
		g.Self = newWaitConsensus(g)
	case insolar.WaitMajority:
		g.Self = newWaitMajority(g)
	case insolar.WaitMinRoles:
		g.Self = newWaitMinRoles(g)
	case insolar.WaitPulsar:
		g.Self = newWaitPulsar(g)
	default:
		inslogger.FromContext(ctx).Panic("Try to switch network to unknown state. Memory of process is inconsistent.")
	}
	return g.Self
}

func (g *Base) Init(ctx context.Context) error {
	g.pulseWatchdog = newPulseWatchdog(ctx, g.Gatewayer.Gateway(), g.Options.PulseWatchdogTimeout)

	g.HostNetwork.RegisterRequestHandler(
		types.Authorize, g.discoveryMiddleware(g.announceMiddleware(g.HandleNodeAuthorizeRequest)), // validate cert
	)
	g.HostNetwork.RegisterRequestHandler(
		types.Bootstrap, g.announceMiddleware(g.HandleNodeBootstrapRequest), // provide joiner claim
	)
	g.HostNetwork.RegisterRequestHandler(types.UpdateSchedule, g.HandleUpdateSchedule)
	g.HostNetwork.RegisterRequestHandler(types.Reconnect, g.HandleReconnect)

	g.bootstrapETA = g.Options.BootstrapTimeout

	return g.initConsensus(ctx)
}

func (g *Base) Stop(ctx context.Context) error {
	err := g.datagramTransport.Stop(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to stop datagram transport")
	}

	g.pulseWatchdog.Stop()
	return nil
}

func (g *Base) initConsensus(ctx context.Context) error {
	g.ConsensusMode = consensus.Joiner
	g.datagramHandler = adapters.NewDatagramHandler()
	datagramTransport, err := g.TransportFactory.CreateDatagramTransport(g.datagramHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create datagramTransport")
	}
	g.datagramTransport = datagramTransport

	proxy := &consensusProxy{g.Gatewayer}
	g.consensusInstaller = consensus.New(ctx, consensus.Dep{
		KeyProcessor:        g.KeyProcessor,
		Scheme:              g.CryptographyScheme,
		CertificateManager:  g.CertificateManager,
		KeyStore:            getKeyStore(g.CryptographyService),
		NodeKeeper:          g.NodeKeeper,
		StateGetter:         proxy,
		PulseChanger:        proxy,
		StateUpdater:        proxy,
		DatagramTransport:   g.datagramTransport,
		EphemeralController: g,
	})

	// transport start should be here because of TestComponents tests, couldn't createOriginCandidate with 0 port
	err = g.datagramTransport.Start(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to start datagram transport")
	}

	err = g.createOriginCandidate()
	if err != nil {
		return errors.Wrap(err, "failed to createOriginCandidate")
	}

	return nil
}

func (g *Base) createOriginCandidate() error {
	// sign origin
	origin := g.NodeKeeper.GetOrigin()
	mutableOrigin := origin.(node.MutableNode)
	mutableOrigin.SetAddress(g.datagramTransport.Address())

	digest, sign, err := getAnnounceSignature(
		origin,
		network.OriginIsDiscovery(g.CertificateManager.GetCertificate()),
		g.KeyProcessor,
		getKeyStore(g.CryptographyService),
		g.CryptographyScheme,
	)
	if err != nil {
		return err
	}
	mutableOrigin.SetSignature(digest, *sign)
	g.NodeKeeper.SetInitialSnapshot([]insolar.NetworkNode{origin})

	staticProfile := adapters.NewStaticProfile(origin, g.CertificateManager.GetCertificate(), g.KeyProcessor)
	g.originCandidate = adapters.NewCandidate(staticProfile, g.KeyProcessor)
	return nil
}

func (g *Base) StartConsensus(ctx context.Context) error {

	if g.NodeKeeper.GetOrigin().Role() == insolar.StaticRoleHeavyMaterial {
		g.ConsensusMode = consensus.ReadyNetwork
	}

	pulseHandler := adapters.NewPulseHandler()
	g.ConsensusController = g.consensusInstaller.ControllerFor(g.ConsensusMode, pulseHandler, g.datagramHandler)
	g.ConsensusController.RegisterFinishedNotifier(func(ctx context.Context, report network.Report) {
		g.Gatewayer.Gateway().OnConsensusFinished(ctx, report)
	})

	g.consensusPulseHandler = pulseHandler
	atomic.StoreUint32(&g.consensusStarted, 1)
	return nil
}

// ChangePulse process pulse from Consensus
func (g *Base) ChangePulse(ctx context.Context, pulse insolar.Pulse) {
	g.Gatewayer.Gateway().OnPulseFromConsensus(ctx, pulse)
}

func (g *Base) OnPulseFromPulsar(ctx context.Context, pu insolar.Pulse, originalPacket network.ReceivedPacket) {
	if atomic.LoadUint32(&g.consensusStarted) == 1 {
		g.consensusPulseHandler.HandlePulse(ctx, pu, originalPacket)
	}
}

func (g *Base) OnPulseFromConsensus(ctx context.Context, pu insolar.Pulse) {
	g.pulseWatchdog.Reset()
	g.NodeKeeper.MoveSyncToActive(ctx, pu.PulseNumber)
	err := g.PulseAppender.AppendPulse(ctx, pu)
	if err != nil {
		inslogger.FromContext(ctx).Panic("failed to append pulse: ", err.Error())
	}

	nodes := g.NodeKeeper.GetAccessor(pu.PulseNumber).GetActiveNodes()
	inslogger.FromContext(ctx).Debugf("OnPulseFromConsensus: %d : epoch %d : nodes %d", pu.PulseNumber, pu.EpochPulseNumber, len(nodes))
}

// UpdateState called then Consensus done
func (g *Base) UpdateState(ctx context.Context, pulseNumber insolar.PulseNumber, nodes []insolar.NetworkNode, cloudStateHash []byte) {
	g.NodeKeeper.Sync(ctx, pulseNumber, nodes)
}

func (g *Base) BeforeRun(ctx context.Context, pulse insolar.Pulse) {}

// Auther casts us to Auther or obtain it in another way
func (g *Base) Auther() network.Auther {
	if ret, ok := g.Self.(network.Auther); ok {
		return ret
	}
	panic("Our network gateway suddenly is not an Auther")
}

// Bootstrapper casts us to Bootstrapper or obtain it in another way
func (g *Base) Bootstrapper() network.Bootstrapper {
	if ret, ok := g.Self.(network.Bootstrapper); ok {
		return ret
	}
	panic("Our network gateway suddenly is not an Bootstrapper")
}

// GetCert method returns node certificate by requesting sign from discovery nodes
func (g *Base) GetCert(ctx context.Context, ref *insolar.Reference) (insolar.Certificate, error) {
	return nil, errors.New("GetCert() in non active mode")
}

// ValidateCert validates node certificate
func (g *Base) ValidateCert(ctx context.Context, authCert insolar.AuthorizationCertificate) (bool, error) {
	return certificate.VerifyAuthorizationCertificate(g.CryptographyService, g.CertificateManager.GetCertificate().GetDiscoveryNodes(), authCert)
}

// ============= Bootstrap =======

func (g *Base) checkCanAnnounceCandidate(ctx context.Context) error {
	// 1. Current node is heavy:
	// 		could announce candidate when network is initialized
	// 		NB: announcing in WaitConsensus state is allowed
	// 2. Otherwise:
	// 		could announce candidate when heavy node found in *active* list and initial consensus passed
	// 		NB: announcing in WaitConsensus state is *NOT* allowed

	state := g.Gatewayer.Gateway().GetState()
	origin := g.OriginProvider.GetOrigin()

	if origin.Role() == insolar.StaticRoleHeavyMaterial && state >= insolar.WaitConsensus {
		return nil
	}

	bootstrapPulse := GetBootstrapPulse(ctx, g.PulseAccessor)
	nodes := g.NodeKeeper.GetAccessor(bootstrapPulse.PulseNumber).GetActiveNodes()

	var hasHeavy bool
	for _, n := range nodes {
		if n.Role() == insolar.StaticRoleHeavyMaterial {
			hasHeavy = true
			break
		}
	}

	if hasHeavy && state > insolar.WaitConsensus {
		return nil
	}

	return errors.Errorf(
		"can't announce candidate: role=%v pulse=%d hasHeavy=%t state=%s",
		origin.Role(),
		bootstrapPulse.PulseNumber,
		hasHeavy,
		state,
	)
}

func (g *Base) announceMiddleware(handler network.RequestHandler) network.RequestHandler {
	return func(ctx context.Context, request network.ReceivedPacket) (network.Packet, error) {
		if err := g.checkCanAnnounceCandidate(ctx); err != nil {
			return nil, err
		}
		return handler(ctx, request)
	}
}

func (g *Base) discoveryMiddleware(handler network.RequestHandler) network.RequestHandler {
	return func(ctx context.Context, request network.ReceivedPacket) (network.Packet, error) {
		if !network.OriginIsDiscovery(g.CertificateManager.GetCertificate()) {
			return nil, errors.New("only discovery nodes could authorize other nodes, this is not a discovery node")
		}
		return handler(ctx, request)
	}
}

func (g *Base) HandleNodeBootstrapRequest(ctx context.Context, request network.ReceivedPacket) (network.Packet, error) {
	if request.GetRequest() == nil || request.GetRequest().GetBootstrap() == nil {
		return nil, errors.Errorf("process bootstrap: got invalid protobuf request message: %s", request)
	}

	data := request.GetRequest().GetBootstrap()

	bootstrapPulse := GetBootstrapPulse(ctx, g.PulseAccessor)

	if network.CheckShortIDCollision(g.NodeKeeper.GetAccessor(bootstrapPulse.PulseNumber).GetActiveNodes(), data.CandidateProfile.ShortID) {
		return g.HostNetwork.BuildResponse(ctx, request, &packet.BootstrapResponse{Code: packet.UpdateShortID}), nil
	}

	err := bootstrap.ValidatePermit(data.Permit, g.CertificateManager.GetCertificate(), g.CryptographyService)
	if err != nil {
		inslogger.FromContext(ctx).Warnf("Rejected bootstrap request from node %s: %s", request.GetSender(), err.Error())
		return g.HostNetwork.BuildResponse(ctx, request, &packet.BootstrapResponse{Code: packet.Reject}), nil
	}

	type candidate struct {
		profiles.StaticProfile
		profiles.StaticProfileExtension
	}

	profile := adapters.Candidate(data.CandidateProfile).StaticProfile(g.KeyProcessor)

	err = g.ConsensusController.AddJoinCandidate(candidate{profile, profile.GetExtension()})
	if err != nil {
		inslogger.FromContext(ctx).Warnf("Rejected bootstrap request from node %s: %s", request.GetSender(), err.Error())
		return g.HostNetwork.BuildResponse(ctx, request, &packet.BootstrapResponse{Code: packet.Reject}), nil
	}

	inslogger.FromContext(ctx).Infof("=== AddJoinCandidate id = %d, address = %s ", data.CandidateProfile.ShortID, data.CandidateProfile.Address)

	return g.HostNetwork.BuildResponse(ctx, request,
		&packet.BootstrapResponse{
			Code:       packet.Accepted,
			Pulse:      *pulse.ToProto(&bootstrapPulse),
			ETASeconds: uint32(g.bootstrapETA.Seconds()),
		}), nil
}

// validateTimestamp returns true if difference between timestamp ant current UTC < delta
func validateTimestamp(timestamp int64, delta time.Duration) bool {
	return time.Now().UTC().Sub(time.Unix(timestamp, 0)) < delta
}

func (g *Base) HandleNodeAuthorizeRequest(ctx context.Context, request network.ReceivedPacket) (network.Packet, error) {
	if request.GetRequest() == nil || request.GetRequest().GetAuthorize() == nil {
		return nil, errors.Errorf("process authorize: got invalid protobuf request message: %s", request)
	}
	data := request.GetRequest().GetAuthorize().AuthorizeData
	o := g.OriginProvider.GetOrigin()

	if data.Version != o.Version() {
		return nil, errors.Errorf("wrong version in AuthorizeRequest, actual network version is: %s", o.Version())
	}

	// TODO: move time.Minute to config
	if !validateTimestamp(data.Timestamp, time.Minute) {
		return g.HostNetwork.BuildResponse(ctx, request, &packet.AuthorizeResponse{
			Code:      packet.WrongTimestamp,
			Timestamp: time.Now().UTC().Unix(),
		}), nil
	}

	cert, err := certificate.Deserialize(data.Certificate, platformpolicy.NewKeyProcessor())
	if err != nil {
		return g.HostNetwork.BuildResponse(ctx, request, &packet.AuthorizeResponse{Code: packet.WrongMandate, Error: err.Error()}), nil
	}

	valid, err := g.Gatewayer.Gateway().Auther().ValidateCert(ctx, cert)
	if err != nil || !valid {
		if err == nil {
			err = errors.New("Certificate validation failed")
		}

		inslogger.FromContext(ctx).Warn("AuthorizeRequest with invalid cert: ", err.Error())
		return g.HostNetwork.BuildResponse(ctx, request, &packet.AuthorizeResponse{Code: packet.WrongMandate, Error: err.Error()}), nil
	}

	// TODO: get random reconnectHost
	// nodes := g.NodeKeeper.GetAccessor().GetActiveNodes()

	// workaround bootstrap to the origin node
	reconnectHost, err := host.NewHostNS(o.Address(), o.ID(), o.ShortID())
	if err != nil {
		err = errors.Wrap(err, "Failed to get reconnectHost")
		inslogger.FromContext(ctx).Warn(err.Error())
		return nil, err
	}

	pubKey, err := g.KeyProcessor.ExportPublicKeyPEM(o.PublicKey())
	if err != nil {
		err = errors.Wrap(err, "Failed to export public key")
		inslogger.FromContext(ctx).Warn(err.Error())
		return nil, err
	}

	permit, err := bootstrap.CreatePermit(g.OriginProvider.GetOrigin().ID(),
		reconnectHost,
		pubKey,
		g.CryptographyService,
	)
	if err != nil {
		return nil, err
	}

	bootstrapPulse := GetBootstrapPulse(ctx, g.PulseAccessor)
	discoveryCount := len(network.FindDiscoveriesInNodeList(
		g.NodeKeeper.GetAccessor(bootstrapPulse.PulseNumber).GetActiveNodes(),
		g.CertificateManager.GetCertificate(),
	))

	return g.HostNetwork.BuildResponse(ctx, request, &packet.AuthorizeResponse{
		Code:           packet.Success,
		Timestamp:      time.Now().UTC().Unix(),
		Permit:         permit,
		DiscoveryCount: uint32(discoveryCount),
		Pulse:          pulse.ToProto(&bootstrapPulse),
	}), nil
}

func (g *Base) HandleUpdateSchedule(ctx context.Context, request network.ReceivedPacket) (network.Packet, error) {
	// TODO:
	return g.HostNetwork.BuildResponse(ctx, request, &packet.UpdateScheduleResponse{}), nil
}

func (g *Base) HandleReconnect(ctx context.Context, request network.ReceivedPacket) (network.Packet, error) {
	if request.GetRequest() == nil || request.GetRequest().GetReconnect() == nil {
		return nil, errors.Errorf("process reconnect: got invalid protobuf request message: %s", request)
	}

	// check permit, if permit from Discovery node
	// request.GetRequest().GetReconnect().Permit

	// TODO:
	return g.HostNetwork.BuildResponse(ctx, request, &packet.ReconnectResponse{}), nil
}

func (g *Base) OnConsensusFinished(ctx context.Context, report network.Report) {
	inslogger.FromContext(ctx).Infof("OnConsensusFinished for pulse %d", report.PulseNumber)
}

func (g *Base) EphemeralMode(nodes []insolar.NetworkNode) bool {
	_, majorityErr := rules.CheckMajorityRule(g.CertificateManager.GetCertificate(), nodes)
	minRoleErr := rules.CheckMinRole(g.CertificateManager.GetCertificate(), nodes)

	return majorityErr != nil || minRoleErr != nil
}

func (g *Base) FailState(ctx context.Context, reason string) {
	o := g.OriginProvider.GetOrigin()
	wrapReason := fmt.Sprintf("Abort node with address: %s role: %s state: %s, reason: %s",
		o.Address(),
		o.Role().String(),
		g.Gatewayer.Gateway().GetState().String(),
		reason,
	)
	g.Aborter.Abort(ctx, wrapReason)
}
