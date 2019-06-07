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
	"strconv"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
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
	"github.com/insolar/insolar/network/consensus/packets"
	"github.com/insolar/insolar/network/consensus/phases"
	"github.com/insolar/insolar/network/controller"
	"github.com/insolar/insolar/network/controller/bootstrap"
	"github.com/insolar/insolar/network/gateway"
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
	GIL                 insolar.GlobalInsolarLock   `inject:""`
	Pub                 message.Publisher           `inject:""`
	MessageBus          insolar.MessageBus          `inject:""`
	ContractRequester   insolar.ContractRequester   `inject:""`
	Sender              bus.Sender                  `inject:""`

	// subcomponents
	PhaseManager phases.PhaseManager `inject:"subcomponent"`
	Controller   network.Controller  `inject:"subcomponent"`

	isGenesis   bool
	isDiscovery bool
	skip        int

	lock sync.Mutex

	gateway   network.Gateway
	gatewayMu sync.RWMutex
}

// NewServiceNetwork returns a new ServiceNetwork.
func NewServiceNetwork(conf configuration.Configuration, rootCm *component.Manager, isGenesis bool) (*ServiceNetwork, error) {
	serviceNetwork := &ServiceNetwork{cm: component.NewManager(rootCm), cfg: conf, isGenesis: isGenesis, skip: conf.Service.Skip}
	return serviceNetwork, nil
}

func (n *ServiceNetwork) Gateway() network.Gateway {
	n.gatewayMu.RLock()
	defer n.gatewayMu.RUnlock()
	return n.gateway
}

func (n *ServiceNetwork) SetGateway(g network.Gateway) {
	n.gatewayMu.Lock()
	defer n.gatewayMu.Unlock()
	n.gateway = g
}

func (n *ServiceNetwork) GetState() insolar.NetworkState {
	return n.Gateway().GetState()
}

// SendMessage sends a message from MessageBus.
func (n *ServiceNetwork) SendMessage(nodeID insolar.Reference, method string, msg insolar.Parcel) ([]byte, error) {
	return n.Controller.SendMessage(nodeID, method, msg)
}

// SendCascadeMessage sends a message from MessageBus to a cascade of nodes
func (n *ServiceNetwork) SendCascadeMessage(data insolar.Cascade, method string, msg insolar.Parcel) error {
	return n.Controller.SendCascadeMessage(data, method, msg)
}

// RemoteProcedureRegister registers procedure for remote call on this host.
func (n *ServiceNetwork) RemoteProcedureRegister(name string, method insolar.RemoteProcedure) {
	n.Controller.RemoteProcedureRegister(name, method)
}

// Init implements component.Initer
func (n *ServiceNetwork) Init(ctx context.Context) error {
	hostNetwork, err := hostnetwork.NewHostNetwork(n.CertificateManager.GetCertificate().GetNodeRef().String())
	if err != nil {
		return errors.Wrap(err, "Failed to create hostnetwork")
	}

	consensusNetwork, err := hostnetwork.NewConsensusNetwork(
		n.CertificateManager.GetCertificate().GetNodeRef().String(),
		n.NodeKeeper.GetOrigin().ShortID(),
	)
	if err != nil {
		return errors.Wrap(err, "Failed to create consensus network.")
	}

	options := controller.ConfigureOptions(n.cfg)

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
		controller.NewNetworkController(),
		controller.NewRPCController(options),
		controller.NewPulseController(),
		bootstrap.NewBootstrapper(options, n.connectToNewNetwork),
		bootstrap.NewAuthorizationController(options),
		bootstrap.NewNetworkBootstrapper(),
	)
	err = n.cm.Init(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to init internal components")
	}

	if n.Gateway() == nil {
		n.gateway = gateway.NewNoNetwork(n, n.GIL, n.NodeKeeper, n.ContractRequester,
			n.CryptographyService, n.MessageBus, n.CertificateManager)
		n.gateway.Run(ctx)
		inslogger.FromContext(ctx).Debug("Launch network gateway")

	}

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
	_, err = n.Controller.Bootstrap(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to bootstrap network")
	}

	n.RemoteProcedureRegister(deliverWatermillMsg, n.processIncoming)

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
	if !n.isGenesis {
		logger.Info("ServiceNetwork.GracefulStop wait for accepting leaving claim")
		n.TerminationHandler.Leave(ctx, 0)
		logger.Info("ServiceNetwork.GracefulStop - leaving claim accepted")
	}

	return nil
}

// Stop implements insolar.Component
func (n *ServiceNetwork) Stop(ctx context.Context) error {
	inslogger.FromContext(ctx).Info("Stopping network component manager...")
	return n.cm.Stop(ctx)
}

func (n *ServiceNetwork) HandlePulse(ctx context.Context, newPulse insolar.Pulse) {
	pulseTime := time.Unix(0, newPulse.PulseTimestamp)
	logger := inslogger.FromContext(ctx)

	n.lock.Lock()
	defer n.lock.Unlock()

	if n.isGenesis {
		return
	}

	// Because we want to set InsTraceID (it's our custom traceID)
	// Because @egorikas didn't have enough time for sending `insTraceID` from pulsar
	// We calculate it 2 times, first time on a pulsar's side. Second time on a network's side
	insTraceID := "pulse_" + strconv.FormatUint(uint64(newPulse.PulseNumber), 10)
	ctx = inslogger.ContextWithTrace(ctx, insTraceID)

	logger.Infof("Got new pulse number: %d", newPulse.PulseNumber)
	ctx, span := instracer.StartSpan(ctx, "ServiceNetwork.Handlepulse")
	span.AddAttributes(
		trace.Int64Attribute("pulse.PulseNumber", int64(newPulse.PulseNumber)),
	)
	defer span.End()

	if !n.NodeKeeper.IsBootstrapped() {
		n.Controller.SetLastIgnoredPulse(newPulse.NextPulseNumber)
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

	// Ignore insolar.ErrNotFound because
	// sometimes we can't fetch current pulse in new nodes
	// (for fresh bootstrapped light-material with in-memory pulse-tracker)
	if currentPulse, err := n.PulseAccessor.Latest(ctx); err != nil {
		if err != pulse.ErrNotFound {
			currentPulse = *insolar.GenesisPulse
		}
	} else {
		if !isNextPulse(&currentPulse, &newPulse) {
			logger.Infof("Incorrect pulse number. Current: %+v. New: %+v", currentPulse, newPulse)
			return
		}
	}

	if err := n.Gateway().OnPulse(ctx, newPulse); err != nil {
		logger.Error(errors.Wrap(err, "Failed to call OnPulse on Gateway"))
	}

	logger.Debugf("Before set new current pulse number: %d", newPulse.PulseNumber)
	err := n.PulseManager.Set(ctx, newPulse, n.Gateway().GetState() == insolar.CompleteNetworkState)
	if err != nil {
		logger.Fatalf("Failed to set new pulse: %s", err.Error())
	}
	logger.Infof("Set new current pulse number: %d", newPulse.PulseNumber)

	go n.phaseManagerOnPulse(ctx, newPulse, pulseTime)
}

func (n *ServiceNetwork) shoudIgnorePulse(newPulse insolar.Pulse) bool {
	return n.isDiscovery && !n.NodeKeeper.GetConsensusInfo().IsJoiner() &&
		newPulse.PulseNumber <= n.Controller.GetLastIgnoredPulse()+insolar.PulseNumber(n.skip)
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

func (n *ServiceNetwork) connectToNewNetwork(ctx context.Context, address string) {
	n.NodeKeeper.GetClaimQueue().Push(&packets.ChangeNetworkClaim{Address: address})
	logger := inslogger.FromContext(ctx)

	node, err := findNodeByAddress(address, n.CertificateManager.GetCertificate().GetDiscoveryNodes())
	if err != nil {
		logger.Warnf("Failed to find a discovery node: ", err)
	}

	err = n.Controller.AuthenticateToDiscoveryNode(ctx, node)
	if err != nil {
		logger.Errorf("Failed to authenticate a node: " + err.Error())
	}
}

// SendMessageHandler async sends message with confirmation of delivery.
func (n *ServiceNetwork) SendMessageHandler(msg *message.Message) ([]*message.Message, error) {
	ctx := inslogger.ContextWithTrace(context.Background(), msg.Metadata.Get(bus.MetaTraceID))
	logger := inslogger.FromContext(ctx)
	msgType, err := payload.UnmarshalType(msg.Payload)
	if err != nil {
		logger.Error("failed to extract message type")
	}

	err = n.sendMessage(ctx, msg)
	if err != nil {
		n.replyError(ctx, msg, err)
		return nil, nil
	}

	logger.WithFields(map[string]interface{}{
		"msg_type":       msgType.String(),
		"correlation_id": middleware.MessageCorrelationID(msg),
	}).Info("Network sent message")

	return nil, nil
}

func (n *ServiceNetwork) sendMessage(ctx context.Context, msg *message.Message) error {
	node, err := n.wrapMeta(msg)
	if err != nil {
		return errors.Wrap(err, "failed to send message")
	}

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
	res, err := n.Controller.SendBytes(ctx, node, deliverWatermillMsg, msgBytes)
	if err != nil {
		return errors.Wrap(err, "error while sending watermillMsg to controller")
	}
	if !bytes.Equal(res, ack) {
		return errors.Errorf("reply is not ack: %s", res)
	}
	return nil
}

func (n *ServiceNetwork) replyError(ctx context.Context, msg *message.Message, repErr error) {
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"correlation_id": middleware.MessageCorrelationID(msg),
	})
	errMsg, err := payload.NewMessage(&payload.Error{Text: repErr.Error()})
	if err != nil {
		logger.Error(errors.Wrapf(err, "failed to create error as reply (%s)", repErr.Error()))
		return
	}
	wrapper := payload.Meta{
		Payload: msg.Payload,
		Sender:  n.NodeKeeper.GetOrigin().ID(),
	}
	buf, err := wrapper.Marshal()
	if err != nil {
		logger.Error(errors.Wrapf(err, "failed to wrap error message (%s)", repErr.Error()))
		return
	}
	msg.Payload = buf
	n.Sender.Reply(ctx, msg, errMsg)
}

func (n *ServiceNetwork) wrapMeta(msg *message.Message) (insolar.Reference, error) {
	receiver := msg.Metadata.Get(bus.MetaReceiver)
	if receiver == "" {
		return insolar.Reference{}, errors.New("Receiver in msg.Metadata not set")
	}
	receiverRef, err := insolar.NewReferenceFromBase58(receiver)
	if err != nil {
		return insolar.Reference{}, errors.Wrap(err, "incorrect Receiver in msg.Metadata")
	}

	latestPulse, err := n.PulseAccessor.Latest(context.Background())
	if err != nil {
		return insolar.Reference{}, errors.Wrap(err, "failed to fetch pulse")
	}
	wrapper := payload.Meta{
		Payload:  msg.Payload,
		Receiver: *receiverRef,
		Sender:   n.NodeKeeper.GetOrigin().ID(),
		Pulse:    latestPulse.PulseNumber,
	}
	buf, err := wrapper.Marshal()
	if err != nil {
		return insolar.Reference{}, errors.Wrap(err, "failed to wrap message")
	}
	msg.Payload = buf

	return *receiverRef, nil
}

func findNodeByAddress(address string, nodes []insolar.DiscoveryNode) (insolar.DiscoveryNode, error) {
	for _, node := range nodes {
		if node.GetHost() == address {
			return node, nil
		}
	}
	return nil, errors.New("Failed to find a discovery node with address: " + address)
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
	ctx = inslogger.ContextWithTrace(ctx, msg.Metadata.Get(bus.MetaTraceID))
	logger = inslogger.FromContext(ctx)
	// TODO: check pulse here

	msgType, err := payload.UnmarshalTypeFromMeta(msg.Payload)
	if err != nil {
		logger.Error("failed to extract message type")
	}
	logger.WithFields(map[string]interface{}{
		"msg_type":       msgType.String(),
		"correlation_id": middleware.MessageCorrelationID(msg),
	}).Info("Network received message")

	err = n.Pub.Publish(bus.TopicIncoming, msg)
	if err != nil {
		err = errors.Wrap(err, "error while publish msg to TopicIncoming")
		logger.Error(err)
		return nil, err
	}

	return ack, nil
}
