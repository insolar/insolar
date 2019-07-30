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

package servicenetwork

import (
	"bytes"
	"context"
	"crypto/rand"

	"github.com/insolar/insolar/cryptography"

	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/serialization"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/network/rules"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/network/controller"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/gateway"
	"github.com/insolar/insolar/network/gateway/bootstrap"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/routing"
	"github.com/insolar/insolar/network/transport"
	"github.com/pkg/errors"
)

// ServiceNetwork is facade for network.
type ServiceNetwork struct {
	cfg configuration.Configuration
	cm  *component.Manager

	// dependencies
	CertificateManager  insolar.CertificateManager         `inject:""`
	PulseManager        insolar.PulseManager               `inject:""`
	PulseAccessor       pulse.Accessor                     `inject:""`
	PulseAppender       pulse.Appender                     `inject:""`
	CryptographyService insolar.CryptographyService        `inject:""`
	CryptographyScheme  insolar.PlatformCryptographyScheme `inject:""`
	KeyProcessor        insolar.KeyProcessor               `inject:""`
	NodeKeeper          network.NodeKeeper                 `inject:""`
	TerminationHandler  insolar.TerminationHandler         `inject:""`
	ContractRequester   insolar.ContractRequester          `inject:""`

	// watermill support interfaces
	Pub message.Publisher `inject:""`

	// subcomponents
	RPC              controller.RPCController `inject:"subcomponent"`
	Rules            network.Rules            `inject:"subcomponent"`
	TransportFactory transport.Factory        `inject:"subcomponent"`

	HostNetwork network.HostNetwork

	CurrentPulse insolar.Pulse
	Gatewayer    network.Gatewayer
	BaseGateway  *gateway.Base
	operableFunc insolar.NetworkOperableCallback

	datagramHandler   *adapters.DatagramHandler
	datagramTransport transport.DatagramTransport

	consensusInstaller  consensus.Installer
	consensusController consensus.Controller

	ConsensusMode consensus.Mode
}

//var PULSETIMEOUT time.Duration

// NewServiceNetwork returns a new ServiceNetwork.
func NewServiceNetwork(conf configuration.Configuration, rootCm *component.Manager) (*ServiceNetwork, error) {
	serviceNetwork := &ServiceNetwork{cm: component.NewManager(rootCm), cfg: conf, ConsensusMode: consensus.Joiner}
	//PULSETIMEOUT = time.Millisecond * time.Duration(conf.Pulsar.PulseTime)
	return serviceNetwork, nil
}

// SendMessage sends a message from MessageBus.
func (n *ServiceNetwork) SendMessage(nodeID insolar.Reference, method string, msg insolar.Parcel) ([]byte, error) {
	return n.RPC.SendMessage(nodeID, method, msg)
}

// SendCascadeMessage sends a message from MessageBus to a cascade of nodes
func (n *ServiceNetwork) SendCascadeMessage(data insolar.Cascade, method string, msg insolar.Parcel) error {
	return n.RPC.SendCascadeMessage(data, method, msg)
}

// RemoteProcedureRegister registers procedure for remote call on this host.
func (n *ServiceNetwork) RemoteProcedureRegister(name string, method insolar.RemoteProcedure) {
	n.RPC.RemoteProcedureRegister(name, method)
}

// Init implements component.Initer
func (n *ServiceNetwork) Init(ctx context.Context) error {
	hostNetwork, err := hostnetwork.NewHostNetwork(n.CertificateManager.GetCertificate().GetNodeRef().String())
	if err != nil {
		return errors.Wrap(err, "failed to create hostnetwork")
	}
	n.HostNetwork = hostNetwork

	options := common.ConfigureOptions(n.cfg)

	cert := n.CertificateManager.GetCertificate()

	n.BaseGateway = &gateway.Base{}
	n.Gatewayer = gateway.NewGatewayer(n.BaseGateway.NewGateway(ctx, insolar.NoNetworkState), func(ctx context.Context, isNetworkOperable bool) {
		if n.operableFunc != nil {
			n.operableFunc(ctx, isNetworkOperable)
		}
	})

	n.cm.Inject(n,
		&routing.Table{},
		cert,
		transport.NewFactory(n.cfg.Host.Transport),
		hostNetwork,
		controller.NewRPCController(options),
		controller.NewPulseController(),
		bootstrap.NewRequester(options),
		network.NewRules(),
		n.BaseGateway,
		n.Gatewayer,
		rules.NewRules(),
	)

	n.datagramHandler = adapters.NewDatagramHandler()
	datagramTransport, err := n.TransportFactory.CreateDatagramTransport(n.datagramHandler)
	if err != nil {
		return errors.Wrap(err, "failed to create datagramTransport")
	}
	n.datagramTransport = datagramTransport

	// sign origin
	origin := n.NodeKeeper.GetOrigin()
	// TODO: hack
	ks := n.CryptographyService.(*cryptography.NodeCryptographyService).KeyStore
	digest, sign, err := getAnnounceSignature(
		origin,
		network.OriginIsDiscovery(cert),
		n.KeyProcessor,
		ks,
		n.CryptographyScheme,
	)
	if err != nil {
		return errors.Wrap(err, "failed to getAnnounceSignature")
	}

	origin.(node.MutableNode).SetSignature(digest, *sign)
	n.NodeKeeper.SetInitialSnapshot([]insolar.NetworkNode{origin})

	err = n.cm.Init(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to init internal components")
	}

	n.CurrentPulse = *insolar.GenesisPulse

	n.consensusInstaller = consensus.New(ctx, consensus.Dep{
		KeyProcessor:        n.KeyProcessor,
		Scheme:              n.CryptographyScheme,
		CertificateManager:  n.CertificateManager,
		KeyStore:            ks,
		NodeKeeper:          n.NodeKeeper,
		StateGetter:         n,
		PulseChanger:        n,
		StateUpdater:        n,
		DatagramTransport:   n.datagramTransport,
		EphemeralController: n,
	})

	return nil
}

func (n *ServiceNetwork) initConsensus() {

	if n.NodeKeeper.GetOrigin().Role() == insolar.StaticRoleHeavyMaterial {
		n.ConsensusMode = consensus.ReadyNetwork
	}

	pulseHandler := adapters.NewPulseHandler()
	n.consensusController = n.consensusInstaller.ControllerFor(n.ConsensusMode, pulseHandler, n.datagramHandler)
	n.consensusController.RegisterFinishedNotifier(func(report network.Report) {
		n.Gatewayer.Gateway().OnConsensusFinished(report.PulseNumber)
	})
	n.BaseGateway.ConsensusController = n.consensusController
	n.BaseGateway.ConsensusPulseHandler = pulseHandler
}

// Start implements component.Starter
func (n *ServiceNetwork) Start(ctx context.Context) error {
	err := n.datagramTransport.Start(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to start datagram transport")
	}

	err = n.cm.Start(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to start component manager")
	}

	//if !n.cfg.Service.ConsensusEnabled {
	//	cert := n.CertificateManager.GetCertificate()
	//	nodes := make([]insolar.NetworkNode, len(cert.GetDiscoveryNodes()))
	//	for i, dn := range cert.GetDiscoveryNodes() {
	//		nodes[i] = node.NewNode(*dn.GetNodeRef(), dn.GetRole(), dn.GetPublicKey(), dn.GetHost(), "")
	//		nodes[i].(node.MutableNode).SetEvidence(cryptkit.NewSignedDigest(
	//			cryptkit.NewDigest(longbits.NewBits512FromBytes(dn.GetBriefDigest()), adapters.SHA3512Digest),
	//			cryptkit.NewSignature(longbits.NewBits512FromBytes(dn.GetBriefSign()), adapters.SHA3512Digest.SignedBy(adapters.SECP256r1Sign)),
	//		))
	//	}
	//	n.operableFunc(ctx, false)
	//	n.NodeKeeper.SetInitialSnapshot(nodes)
	//	n.Gatewayer.SwitchState(ctx, insolar.CompleteNetworkState)
	//	n.ConsensusMode = consensus.ReadyNetwork
	//}

	n.initConsensus()
	n.Gatewayer.Gateway().Run(ctx)

	n.RemoteProcedureRegister(deliverWatermillMsg, n.processIncoming)

	return nil
}

func (n *ServiceNetwork) Leave(ctx context.Context, eta insolar.PulseNumber) {
	logger := inslogger.FromContext(ctx)
	logger.Info("Gracefully stopping service network")

	// TODO: fix leave
	n.consensusController.Leave(0)
}

func (n *ServiceNetwork) GracefulStop(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	// node leaving from network
	// all components need to do what they want over net in gracefulStop

	logger.Info("ServiceNetwork.GracefulStop wait for accepting leaving claim")
	n.TerminationHandler.Leave(ctx, 0)
	logger.Info("ServiceNetwork.GracefulStop - leaving claim accepted")

	return nil
}

// Stop implements insolar.Component
func (n *ServiceNetwork) Stop(ctx context.Context) error {
	err := n.datagramTransport.Stop(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to stop datagram transport")
	}

	return n.cm.Stop(ctx)
}

func (n *ServiceNetwork) GetState() insolar.NetworkState {
	return n.Gatewayer.Gateway().GetState()
}

func (n *ServiceNetwork) SetOperableFunc(f insolar.NetworkOperableCallback) {
	n.operableFunc = f
}

// HandlePulse process pulse from PulseController
func (n *ServiceNetwork) HandlePulse(ctx context.Context, pulse insolar.Pulse, originalPacket network.ReceivedPacket) {

	//pulseTime := time.Unix(0, newPulse.PulseTimestamp)
	//logger := inslogger.FromContext(ctx)
	//
	//n.lock.Lock()
	//defer n.lock.Unlock()
	//done := make(chan struct{})
	//defer close(done)
	//go func() {
	//	select {
	//	case <-time.After(n.pulseTimeout):
	//		log.Error("Node stopped due to long pulse processing")
	//	case <-done:
	//	}
	//}()

	n.Gatewayer.Gateway().OnPulseFromPulsar(ctx, pulse, originalPacket)
}

// consensus handlers here

// ChangePulse process pulse from Consensus
func (n *ServiceNetwork) ChangePulse(ctx context.Context, pulse insolar.Pulse) {
	n.CurrentPulse = pulse
	n.Gatewayer.Gateway().OnPulseFromConsensus(ctx, pulse)
}

func (n *ServiceNetwork) UpdateState(ctx context.Context, pulseNumber insolar.PulseNumber, nodes []insolar.NetworkNode, cloudStateHash []byte) {
	n.Gatewayer.Gateway().UpdateState(ctx, pulseNumber, nodes, cloudStateHash)
}

func (n *ServiceNetwork) State() []byte {
	nshBytes := make([]byte, 64)
	_, _ = rand.Read(nshBytes)
	return nshBytes
}

func getAnnounceSignature(
	node insolar.NetworkNode,
	isDiscovery bool,
	kp insolar.KeyProcessor,
	keystore insolar.KeyStore,
	scheme insolar.PlatformCryptographyScheme,
) ([]byte, *insolar.Signature, error) {

	brief := serialization.NodeBriefIntro{}
	brief.ShortID = node.ShortID()
	brief.SetPrimaryRole(adapters.StaticRoleToPrimaryRole(node.Role()))
	if isDiscovery {
		brief.SpecialRoles = member.SpecialRoleDiscovery
	}
	brief.StartPower = 10

	addr, err := endpoints.NewIPAddress(node.Address())
	if err != nil {
		return nil, nil, err
	}
	copy(brief.Endpoint[:], addr[:])

	pk, err := kp.ExportPublicKeyBinary(node.PublicKey())
	if err != nil {
		return nil, nil, err
	}

	copy(brief.NodePK[:], pk)

	buf := &bytes.Buffer{}
	err = brief.SerializeTo(nil, buf)
	if err != nil {
		panic(err)
	}

	data := buf.Bytes()
	data = data[:len(data)-64]

	key, err := keystore.GetPrivateKey("")
	if err != nil {
		panic(err)
	}

	digest := scheme.IntegrityHasher().Hash(data)
	sign, err := scheme.DigestSigner(key).Sign(digest)
	if err != nil {
		panic(err)
	}

	return digest, sign, nil
}

// RegisterConsensusFinishedNotifier for integrtest TODO: remove
func (n *ServiceNetwork) RegisterConsensusFinishedNotifier(fn network.OnConsensusFinished) {
	n.consensusController.RegisterFinishedNotifier(fn)
}

func (n *ServiceNetwork) GetCert(ctx context.Context, ref *insolar.Reference) (insolar.Certificate, error) {
	return n.Gatewayer.Gateway().Auther().GetCert(ctx, ref)
}

func (n *ServiceNetwork) EphemeralMode() bool {
	return n.Gatewayer.Gateway().EphemeralMode()
	//return false
}
