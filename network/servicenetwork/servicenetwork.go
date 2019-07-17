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
	"sync"
	"time"

	"github.com/insolar/insolar/network/gateway"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensusv1/packets"
	"github.com/insolar/insolar/network/consensusv1/phases"
	"github.com/insolar/insolar/network/controller"
	"github.com/insolar/insolar/network/controller/bootstrap"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/merkle"
	"github.com/insolar/insolar/network/routing"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/utils"
)

const deliverWatermillMsg = "ServiceNetwork.processIncoming"

var ack = []byte{1}

// ServiceNetwork is facade for network.
type ServiceNetwork struct {
	cfg configuration.Configuration
	cm  *component.Manager

	// dependencies
	CertificateManager  insolar.CertificateManager  `inject:""`
	PulseManager        insolar.PulseManager        `inject:""`
	PulseAccessor       pulse.Accessor              `inject:""`
	CryptographyService insolar.CryptographyService `inject:""`
	NodeKeeper          network.NodeKeeper          `inject:""`
	TerminationHandler  insolar.TerminationHandler  `inject:""`
	Pub                 message.Publisher           `inject:""`
	ContractRequester   insolar.ContractRequester   `inject:""`

	// subcomponents
	PhaseManager phases.PhaseManager           `inject:"subcomponent"`
	Bootstrapper bootstrap.NetworkBootstrapper `inject:"subcomponent"`
	RPC          controller.RPCController      `inject:"subcomponent"`

	HostNetwork  network.HostNetwork
	OperableFunc func(ctx context.Context, operable bool)

	CurrentPulse insolar.Pulse

	isDiscovery bool
	skip        int

	lock sync.Mutex

	gateway   network.Gateway
	gatewayMu sync.RWMutex
}

// NewServiceNetwork returns a new ServiceNetwork.
func NewServiceNetwork(conf configuration.Configuration, rootCm *component.Manager) (*ServiceNetwork, error) {
	serviceNetwork := &ServiceNetwork{cm: component.NewManager(rootCm), cfg: conf, skip: conf.Service.Skip}
	return serviceNetwork, nil
}

func (n *ServiceNetwork) SetOperableFunc(f func(ctx context.Context, operable bool)) {
	n.OperableFunc = f
}

func (n *ServiceNetwork) Gateway() network.Gateway {
	n.gatewayMu.RLock()
	defer n.gatewayMu.RUnlock()
	return n.gateway
}

func (n *ServiceNetwork) SetGateway(g network.Gateway) {
	n.gatewayMu.Lock()

	if n.gateway != nil && n.gateway.GetState() == g.GetState() {
		log.Warn("Trying to set gateway to the same state")
		n.gatewayMu.Unlock()
		return
	}
	n.gateway = g
	n.gatewayMu.Unlock()
	ctx := context.Background()
	if n.gateway.NeedLockMessageBus() {
		n.OperableFunc(ctx, false)
	} else {
		n.OperableFunc(ctx, true)
	}
	n.gateway.Run(ctx)
}

func (n *ServiceNetwork) GetState() insolar.NetworkState {
	g := n.Gateway()
	if g == nil {
		return insolar.NoNetworkState
	}
	return g.GetState()
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

	consensusNetwork, err := hostnetwork.NewConsensusNetwork(
		n.CertificateManager.GetCertificate().GetNodeRef().String(),
		n.NodeKeeper.GetOrigin().ShortID(),
	)
	if err != nil {
		return errors.Wrap(err, "Failed to create consensus network.")
	}

	options := common.ConfigureOptions(n.cfg)

	cert := n.CertificateManager.GetCertificate()
	n.isDiscovery = utils.OriginIsDiscovery(cert)

	n.cm.Inject(n,
		&routing.Table{},
		cert,
		transport.NewFactory(n.cfg.Host.Transport),
		hostNetwork,
		merkle.NewCalculator(),
		consensusNetwork,
		phases.NewCommunicator(),
		phases.NewFirstPhase(),
		phases.NewSecondPhase(),
		phases.NewThirdPhase(),
		phases.NewPhaseManager(n.cfg.Service.Consensus),
		bootstrap.NewSessionManager(),
		controller.NewRPCController(options),
		controller.NewPulseController(),
		bootstrap.NewBootstrapper(options),
		bootstrap.NewAuthorizationController(options),
		bootstrap.NewNetworkBootstrapper(),
	)
	err = n.cm.Init(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to init internal components")
	}

	n.CurrentPulse = *insolar.GenesisPulse
	return nil
}

// Start implements component.Starter
func (n *ServiceNetwork) Start(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	log.Info("Starting network component manager...")
	err := n.cm.Start(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to start component manager")
	}

	log.Info("Bootstrapping network...")
	_, err = n.Bootstrapper.Bootstrap(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to bootstrap network")
	}

	n.RemoteProcedureRegister(deliverWatermillMsg, n.processIncoming)

	n.SetGateway(gateway.NewNoNetwork(n, n.PulseManager, n.NodeKeeper, n.ContractRequester,
		n.CryptographyService, n.HostNetwork, n.CertificateManager))

	logger.Info("Service network started")
	return nil
}

func (n *ServiceNetwork) Leave(ctx context.Context, eta insolar.PulseNumber) {
	logger := inslogger.FromContext(ctx)
	logger.Info("Gracefully stopping service network")

	n.NodeKeeper.GetClaimQueue().Push(&packets.NodeLeaveClaim{ETA: eta})
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
	inslogger.FromContext(ctx).Info("Stopping network component manager...")
	return n.cm.Stop(ctx)
}

func (n *ServiceNetwork) HandlePulse(ctx context.Context, newPulse insolar.Pulse, _ network.ReceivedPacket) {
	pulseTime := time.Unix(0, newPulse.PulseTimestamp)
	logger := inslogger.FromContext(ctx)

	n.lock.Lock()
	defer n.lock.Unlock()

	ctx = utils.NewPulseContext(ctx, uint32(newPulse.PulseNumber))
	logger.Infof("Got new pulse number: %d", newPulse.PulseNumber)
	ctx, span := instracer.StartSpan(ctx, "ServiceNetwork.Handlepulse")
	span.AddAttributes(
		trace.Int64Attribute("pulse.PulseNumber", int64(newPulse.PulseNumber)),
	)
	defer span.End()

	if !n.NodeKeeper.IsBootstrapped() {
		n.Bootstrapper.SetLastPulse(newPulse.NextPulseNumber)
		return
	}
	if n.shoudIgnorePulse(newPulse) {
		log.Infof("Ignore pulse %d: network is not yet initialized", newPulse.PulseNumber)
		return
	}

	if n.NodeKeeper.GetConsensusInfo().IsJoiner() {
		// do not set pulse because otherwise we will set invalid active list
		// pass consensus, prepare valid active list and set it on next pulse
		go n.phaseManagerOnPulse(ctx, newPulse, pulseTime)
		return
	}

	//TODO: use network pulsestorage here

	if n.CurrentPulse.PulseNumber != insolar.GenesisPulse.PulseNumber && !isNextPulse(&n.CurrentPulse, &newPulse) {
		logger.Infof("Incorrect pulse number. Current: %+v. New: %+v", n.CurrentPulse, newPulse)
		return
	}

	n.ChangePulse(ctx, newPulse)

	go n.phaseManagerOnPulse(ctx, newPulse, pulseTime)
}

func (n *ServiceNetwork) ChangePulse(ctx context.Context, newPulse insolar.Pulse) {
	logger := inslogger.FromContext(ctx)
	n.CurrentPulse = newPulse

	if err := n.NodeKeeper.MoveSyncToActive(ctx, newPulse.PulseNumber); err != nil {
		logger.Warn("MoveSyncToActive failed: ", err.Error())
	}

	if err := n.Gateway().OnPulse(ctx, newPulse); err != nil {
		logger.Error(errors.Wrap(err, "Failed to call OnPulse on Gateway"))
	}
}

func (n *ServiceNetwork) shoudIgnorePulse(newPulse insolar.Pulse) bool {
	return n.isDiscovery && !n.NodeKeeper.GetConsensusInfo().IsJoiner() &&
		newPulse.PulseNumber <= n.Bootstrapper.GetLastPulse()+insolar.PulseNumber(n.skip)
}

func (n *ServiceNetwork) phaseManagerOnPulse(ctx context.Context, newPulse insolar.Pulse, pulseStartTime time.Time) {
	logger := inslogger.FromContext(ctx)

	if !n.cfg.Service.ConsensusEnabled {
		logger.Warn("Consensus is disabled")
		return
	}

	if err := n.PhaseManager.OnPulse(ctx, &newPulse, pulseStartTime); err != nil {
		errMsg := "Failed to pass consensus: " + err.Error()
		logger.Error(errMsg)
		n.SetGateway(n.Gateway().NewGateway(insolar.NoNetworkState))
	}
}

// SendMessageHandler async sends message with confirmation of delivery.
func (n *ServiceNetwork) SendMessageHandler(msg *message.Message) ([]*message.Message, error) {
	ctx := inslogger.ContextWithTrace(context.Background(), msg.Metadata.Get(bus.MetaTraceID))
	parentSpan, err := instracer.Deserialize([]byte(msg.Metadata.Get(bus.MetaSpanData)))
	if err == nil {
		ctx = instracer.WithParentSpan(ctx, parentSpan)
	} else {
		inslogger.FromContext(ctx).Error(err)
	}
	err = n.sendMessage(ctx, msg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send message")
	}
	return nil, nil
}

func (n *ServiceNetwork) sendMessage(ctx context.Context, msg *message.Message) error {
	meta := payload.Meta{}
	err := meta.Unmarshal(msg.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unwrap message")
	}
	if meta.Receiver.IsEmpty() {
		return errors.New("failed to send message: Receiver in meta message not set")
	}

	node := meta.Receiver

	// Short path when sending to self node. Skip serialization
	origin := n.NodeKeeper.GetOrigin()
	if node.Equal(origin.ID()) {
		err := n.Pub.Publish(bus.TopicIncoming, msg)
		if err != nil {
			return errors.Wrap(err, "error while publish msg to TopicIncoming")
		}
		return nil
	}
	msgBytes, err := messageToBytes(msg)
	if err != nil {
		return errors.Wrap(err, "error while converting message to bytes")
	}
	res, err := n.RPC.SendBytes(ctx, node, deliverWatermillMsg, msgBytes)
	if err != nil {
		return errors.Wrap(err, "error while sending watermillMsg to controller")
	}
	if !bytes.Equal(res, ack) {
		return errors.Errorf("reply is not ack: %s", res)
	}
	return nil
}

func isNextPulse(currentPulse, newPulse *insolar.Pulse) bool {
	return newPulse.PulseNumber > currentPulse.PulseNumber && newPulse.PulseNumber >= currentPulse.NextPulseNumber
}

func (n *ServiceNetwork) processIncoming(ctx context.Context, args []byte) ([]byte, error) {
	logger := inslogger.FromContext(ctx)
	msg, err := deserializeMessage(bytes.NewBuffer(args))
	if err != nil {
		err = errors.Wrap(err, "error while deserialize msg from buffer")
		logger.Error(err)
		return nil, err
	}
	logger = inslogger.FromContext(ctx)
	if inslogger.TraceID(ctx) != msg.Metadata.Get(bus.MetaTraceID) {
		logger.Errorf("traceID from context (%s) is different from traceID from message Metadata (%s)", inslogger.TraceID(ctx), msg.Metadata.Get(bus.MetaTraceID))
	}
	// TODO: check pulse here

	err = n.Pub.Publish(bus.TopicIncoming, msg)
	if err != nil {
		err = errors.Wrap(err, "error while publish msg to TopicIncoming")
		logger.Error(err)
		return nil, err
	}

	return ack, nil
}
