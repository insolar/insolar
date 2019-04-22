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
	"github.com/insolar/insolar/insolar/flow/handler"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/consensus/phases"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/storage/pulse"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller"
	"github.com/insolar/insolar/network/controller/bootstrap"
	"github.com/insolar/insolar/network/gateway"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/merkle"
	"github.com/insolar/insolar/network/routing"
	"github.com/insolar/insolar/network/utils"
)

type resultSetter interface {
	SetResult(ctx context.Context, msg *message.Message)
}

const deliverWatermillMsg = "ServiceNetwork.processIncome"

// ServiceNetwork is facade for network.
type ServiceNetwork struct {
	cfg configuration.Configuration
	cm  *component.Manager

	// dependencies
	CertificateManager  insolar.CertificateManager  `inject:""`
	PulseManager        insolar.PulseManager        `inject:""`
	PulseAccessor       pulse.Accessor              `inject:""`
	CryptographyService insolar.CryptographyService `inject:""`
	NetworkCoordinator  insolar.NetworkCoordinator  `inject:""`
	NodeKeeper          network.NodeKeeper          `inject:""`
	TerminationHandler  insolar.TerminationHandler  `inject:""`
	GIL                 insolar.GlobalInsolarLock   `inject:""`
	ResultSetter        resultSetter                `inject:""`

	// subcomponents
	PhaseManager phases.PhaseManager `inject:"subcomponent"`
	Controller   network.Controller  `inject:"subcomponent"`

	Handler *handler.Handler

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

// Start implements component.Initer
func (n *ServiceNetwork) Init(ctx context.Context) error {
	hostNetwork, err := hostnetwork.NewHostNetwork(n.cfg, n.CertificateManager.GetCertificate().GetNodeRef().String())
	if err != nil {
		return errors.Wrap(err, "Failed to create internal transport")
	}

	consensusAddress := n.cfg.Host.Transport.Address
	if n.cfg.Host.Transport.FixedPublicAddress == "" {
		consensusAddress = n.NodeKeeper.GetOrigin().Address()
	}

	consensusNetwork, err := hostnetwork.NewConsensusNetwork(
		consensusAddress,
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
		hostNetwork,
		// use flaky network instead of hostNetwork to imitate network delays
		// NewFlakyNetwork(hostNetwork),
		merkle.NewCalculator(),
		consensusNetwork,
		phases.NewCommunicator(),
		phases.NewFirstPhase(),
		phases.NewSecondPhase(),
		phases.NewThirdPhase(),
		phases.NewPhaseManager(),
		bootstrap.NewSessionManager(),
		controller.NewNetworkController(),
		controller.NewRPCController(options),
		controller.NewPulseController(),
		bootstrap.NewBootstrapper(options, n.connectToNewNetwork),
		bootstrap.NewAuthorizationController(options),
		bootstrap.NewChallengeResponseController(options),
		bootstrap.NewNetworkBootstrapper(),
	)
	err = n.cm.Init(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to init internal components")
	}

	n.gateway = gateway.NewNoNetwork(n, n.GIL, n.NodeKeeper)
	n.gateway.Run(ctx)

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

	n.RemoteProcedureRegister(deliverWatermillMsg, n.processIncome)

	logger.Info("Service network started")
	return nil
}

func (n *ServiceNetwork) Leave(ctx context.Context, ETA insolar.PulseNumber) {
	logger := inslogger.FromContext(ctx)
	logger.Info("Gracefully stopping service network")

	n.NodeKeeper.GetClaimQueue().Push(&packets.NodeLeaveClaim{ETA: ETA})
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
	currentTime := time.Now()

	n.lock.Lock()
	defer n.lock.Unlock()

	if n.isGenesis {
		return
	}
	traceID := "pulse_" + strconv.FormatUint(uint64(newPulse.PulseNumber), 10)
	ctx, logger := inslogger.WithTraceField(ctx, traceID)
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
		go n.phaseManagerOnPulse(ctx, newPulse, currentTime)
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

	go n.phaseManagerOnPulse(ctx, newPulse, currentTime)
}

func (n *ServiceNetwork) shoudIgnorePulse(newPulse insolar.Pulse) bool {
	return n.isDiscovery && !n.NodeKeeper.GetConsensusInfo().IsJoiner() &&
		newPulse.PulseNumber <= n.Controller.GetLastIgnoredPulse()+insolar.PulseNumber(n.skip)
}

func (n *ServiceNetwork) phaseManagerOnPulse(ctx context.Context, newPulse insolar.Pulse, pulseStartTime time.Time) {
	logger := inslogger.FromContext(ctx)

	if err := n.PhaseManager.OnPulse(ctx, &newPulse, pulseStartTime); err != nil {
		errMsg := "Failed to pass consensus: " + err.Error()
		logger.Error(errMsg)
		n.TerminationHandler.Abort(errMsg)
	}
}

func (n *ServiceNetwork) connectToNewNetwork(ctx context.Context, node insolar.DiscoveryNode) {
	err := n.Controller.AuthenticateToDiscoveryNode(ctx, node)
	if err != nil {
		logger := inslogger.FromContext(ctx)
		logger.Errorf("Failed to authenticate a node: " + err.Error())
	}
}

// ProcessOutcome processes received message.
func (n *ServiceNetwork) ProcessOutcome(msg *message.Message) ([]*message.Message, error) {
	receiver := msg.Metadata.Get(insolar.ReceiverMetadataKey)
	if receiver == "" {
		return nil, errors.New("Receiver in msg.Metadata not set")
	}
	ref, err := insolar.NewReferenceFromBase58(receiver)
	if err != nil {
		return nil, errors.Wrap(err, "incorrect Receiver in msg.Metadata: ")
	}
	node := *ref
	// Short path when sending to self node. Skip serialization
	origin := n.NodeKeeper.GetOrigin()
	if node.Equal(origin.ID()) {
		err := n.Handler.Process(msg.Context(), msg)
		if err != nil {
			return nil, errors.Wrap(err, "error while start handle for msg:")
		}
		return nil, nil
	}
	msgBytes := MessageToBytes(msg)
	res, err := n.Controller.SendBytes(msg.Context(), node, deliverWatermillMsg, msgBytes)
	if err != nil {
		return nil, errors.Wrap(err, "error while sending watermillMsg to controller:")
	}
	rep, err := reply.Deserialize(bytes.NewBuffer(res))
	if err != nil {
		return nil, errors.Wrap(err, "error while deserialize reply:")
	}
	if rep.Type() != reply.TypeOK {
		return nil, errors.Errorf("reply is not ok: %s", rep)
	}
	return nil, nil
}

func isNextPulse(currentPulse, newPulse *insolar.Pulse) bool {
	return newPulse.PulseNumber > currentPulse.PulseNumber && newPulse.PulseNumber >= currentPulse.NextPulseNumber
}

func (n *ServiceNetwork) processIncome(ctx context.Context, args [][]byte) ([]byte, error) {
	if len(args) < 1 {
		return nil, errors.New("need exactly one argument when n.processIncome()")
	}
	msg, err := DeserializeMessage(bytes.NewBuffer(args[0]))
	if err != nil {
		return nil, err
	}
	// TODO: check pulse here

	if msg.Metadata.Get(insolar.TypeMetadataKey) == insolar.ReplyTypeMetadataValue {
		n.ResultSetter.SetResult(ctx, msg)
	} else {
		err = n.Handler.Process(ctx, msg)
		if err != nil {
			return nil, errors.Wrap(err, "error while start handle for msg:")
		}
	}

	rd, err := reply.Serialize(&reply.OK{})
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(rd)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
