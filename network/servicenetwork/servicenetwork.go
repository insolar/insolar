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
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/serialization"
	"github.com/insolar/insolar/network/consensusv1/packets"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/network/rules"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/gateway"
	"github.com/insolar/insolar/network/gateway/bootstrap"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/routing"
	"github.com/insolar/insolar/network/transport"
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
	KeyStore            insolar.KeyStore                   `inject:""`
	TransportFactory    transport.Factory                  `inject:""`
	NodeKeeper          network.NodeKeeper                 `inject:""`
	TerminationHandler  insolar.TerminationHandler         `inject:""`
	ContractRequester   insolar.ContractRequester          `inject:""`

	// watermill support interfaces
	Pub message.Publisher `inject:""`

	// subcomponents
	//PhaseManager phases.PhaseManager      `inject:"subcomponent"`
	RPC   controller.RPCController `inject:"subcomponent"`
	Rules network.Rules            `inject:"subcomponent"`

	HostNetwork network.HostNetwork

	CurrentPulse insolar.Pulse
	Gatewayer    network.Gatewayer
	BaseGateway  *gateway.Base
	operableFunc insolar.NetworkOperableCallback

	pulseHandler        *adapters.PulseHandler
	datagramHandler     *adapters.DatagramHandler
	datagramTransport   transport.DatagramTransport
	consensusInstaller  consensus.Installer
	consensusController consensus.Controller

	ConsensusMode consensus.Mode

	lock sync.Mutex
}

var PULSETIMEOUT time.Duration

// NewServiceNetwork returns a new ServiceNetwork.
func NewServiceNetwork(conf configuration.Configuration, rootCm *component.Manager) (*ServiceNetwork, error) {
	serviceNetwork := &ServiceNetwork{cm: component.NewManager(rootCm), cfg: conf, ConsensusMode: consensus.Joiner}
	PULSETIMEOUT = time.Millisecond * time.Duration(conf.Pulsar.PulseTime)
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
		return errors.Wrap(err, "Failed to create hostnetwork")
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

	n.pulseHandler = adapters.NewPulseHandler()
	n.datagramHandler = adapters.NewDatagramHandler()
	datagramTransport, err := n.TransportFactory.CreateDatagramTransport(n.datagramHandler)
	if err != nil {
		return errors.Wrap(err, "Failed to create datagramTransport")
	}
	n.datagramTransport = datagramTransport

	n.cm.Inject(n,
		&routing.Table{},
		cert,
		transport.NewFactory(n.cfg.Host.Transport),
		hostNetwork,
		n.pulseHandler,
		n.datagramHandler,
		n.datagramTransport,

		controller.NewRPCController(options),
		controller.NewPulseController(),
		bootstrap.NewRequester(options),
		network.NewRules(),
		n.BaseGateway,
		n.Gatewayer,
		rules.NewRules(),
	)

	// sign origin
	origin := n.NodeKeeper.GetOrigin()
	digest, sign := getAnnounceSignature(
		origin,
		network.OriginIsDiscovery(cert),
		n.KeyProcessor,
		n.KeyStore,
		n.CryptographyScheme,
	)
	origin.(node.MutableNode).SetSignature(digest, sign)

	err = n.cm.Init(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to init internal components")
	}

	n.CurrentPulse = *insolar.GenesisPulse

	n.consensusInstaller = consensus.New(ctx, consensus.Dep{
		PrimingCloudStateHash: [64]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
		KeyProcessor:          n.KeyProcessor,
		Scheme:                n.CryptographyScheme,
		CertificateManager:    n.CertificateManager,
		KeyStore:              n.KeyStore,
		NodeKeeper:            n.NodeKeeper,
		StateGetter:           n,
		PulseChanger:          n,
		StateUpdater:          n,
		DatagramTransport:     n.datagramTransport,
	})

	return nil
}

// Start implements component.Starter
func (n *ServiceNetwork) Start(ctx context.Context) error {
	err := n.cm.Start(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to start component manager")
	}

	n.Gatewayer.Gateway().Run(ctx)

	n.consensusController = n.consensusInstaller.ControllerFor(n.ConsensusMode, n.pulseHandler, n.datagramHandler)
	n.consensusController.RegisterFinishedNotifier(func(_ member.OpMode, _ member.Power, effectiveSince insolar.PulseNumber) {
		n.Gatewayer.Gateway().OnConsensusFinished(effectiveSince)
	})
	n.BaseGateway.ConsensusController = n.consensusController

	n.RemoteProcedureRegister(deliverWatermillMsg, n.processIncoming)

	// logger.Info("Service network started")
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
	return n.cm.Stop(ctx)
}

func (n *ServiceNetwork) GetState() insolar.NetworkState {
	return n.Gatewayer.Gateway().GetState()
}

func (n *ServiceNetwork) SetOperableFunc(f insolar.NetworkOperableCallback) {
	n.operableFunc = f
}

// consensus handlers here

func (n *ServiceNetwork) ChangePulse(ctx context.Context, newPulse insolar.Pulse) {
	logger := inslogger.FromContext(ctx)

	logger.Infof("Got new pulse number: %d", newPulse.PulseNumber)
	ctx, span := instracer.StartSpan(ctx, "ServiceNetwork.Handlepulse")
	span.AddAttributes(
		trace.Int64Attribute("pulse.PulseNumber", int64(newPulse.PulseNumber)),
	)
	defer span.End()

	n.CurrentPulse = newPulse

	if err := n.NodeKeeper.MoveSyncToActive(ctx, newPulse.PulseNumber); err != nil {
		logger.Warn("MoveSyncToActive failed: ", err.Error())
	}

	if err := n.Gatewayer.Gateway().OnPulse(ctx, newPulse); err != nil {
		logger.Error(errors.Wrap(err, "Failed to call OnPulse on Gateway"))
	}
}

func (n *ServiceNetwork) UpdateState(ctx context.Context, pulseNumber insolar.PulseNumber, nodes []insolar.NetworkNode, cloudStateHash []byte) {
	err := n.NodeKeeper.Sync(ctx, nodes)
	if err != nil {
		inslogger.FromContext(ctx).Error(err)
	}
	n.NodeKeeper.SetCloudHash(cloudStateHash)
}

func (n *ServiceNetwork) State() []byte {
	nshBytes := make([]byte, 64)
	_, _ = rand.Read(nshBytes)
	return nshBytes
}

func (n *ServiceNetwork) RegisterConsensusFinishedNotifier(fn consensus.FinishedNotifier) {
	n.consensusController.RegisterFinishedNotifier(fn)
}

func getAnnounceSignature(
	node insolar.NetworkNode,
	isDiscovery bool,
	kp insolar.KeyProcessor,
	keystore insolar.KeyStore,
	scheme insolar.PlatformCryptographyScheme,
) ([]byte, insolar.Signature) {

	brief := serialization.NodeBriefIntro{}
	brief.ShortID = node.ShortID()
	brief.SetPrimaryRole(adapters.StaticRoleToPrimaryRole(node.Role()))
	if isDiscovery {
		brief.SpecialRoles = member.SpecialRoleDiscovery
	}
	brief.StartPower = 10

	addr, err := packets.NewNodeAddress(node.Address())
	if err != nil {
		panic(err)
	}
	copy(brief.Endpoint[:], addr[:])

	pk, err := kp.ExportPublicKeyBinary(node.PublicKey())
	if err != nil {
		panic(err)
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

	return digest, *sign
}
