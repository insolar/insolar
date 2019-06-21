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

	"github.com/insolar/insolar/network/gateway"
	"github.com/insolar/insolar/network/gateway/bootstrap"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/network/controller/common"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus/packets"
	"github.com/insolar/insolar/network/consensus/phases"
	"github.com/insolar/insolar/network/controller"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/merkle"
	"github.com/insolar/insolar/network/routing"
	"github.com/insolar/insolar/network/transport"
)

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
	ContractRequester   insolar.ContractRequester   `inject:""` // todo: remove
	Sender              bus.Sender                  `inject:""`

	// subcomponents
	PhaseManager phases.PhaseManager      `inject:"subcomponent"`
	RPC          controller.RPCController `inject:"subcomponent"`

	HostNetwork network.HostNetwork

	gatewayer    network.Gatewayer
	operableFunc insolar.NetworkOperableCallback

	isGenesis bool
	// isDiscovery bool

	lock sync.Mutex
}

// NewServiceNetwork returns a new ServiceNetwork.
func NewServiceNetwork(conf configuration.Configuration, rootCm *component.Manager, isGenesis bool) (*ServiceNetwork, error) {
	serviceNetwork := &ServiceNetwork{cm: component.NewManager(rootCm), cfg: conf, isGenesis: isGenesis}
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

	consensusNetwork, err := hostnetwork.NewConsensusNetwork(
		n.CertificateManager.GetCertificate().GetNodeRef().String(),
		n.NodeKeeper.GetOrigin().ShortID(),
	)
	if err != nil {
		return errors.Wrap(err, "Failed to create consensus network.")
	}

	options := common.ConfigureOptions(n.cfg)

	cert := n.CertificateManager.GetCertificate()

	baseGateway := &gateway.Base{}
	n.gatewayer = gateway.NewGatewayer(baseGateway.NewGateway(insolar.NoNetworkState), func(ctx context.Context, isNetworkOperable bool) {
		if n.operableFunc != nil {
			n.operableFunc(ctx, isNetworkOperable)
		}
	})

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
		controller.NewRPCController(options),
		controller.NewPulseController(),
		bootstrap.NewRequester(options),
		baseGateway,
		n.gatewayer,
	)
	err = n.cm.Init(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to init internal components")
	}

	return nil
}

// Start implements component.Starter
func (n *ServiceNetwork) Start(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	logger.Info("Starting network component manager...")
	err := n.cm.Start(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to start component manager")
	}

	n.gatewayer.Gateway().Run(ctx)

	n.RemoteProcedureRegister(deliverWatermillMsg, n.processIncoming)

	// logger.Info("Service network started")
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

	if n.gatewayer.Gateway().ShoudIgnorePulse(ctx, newPulse) {
		return
	}
	//todo call gatewayer
	/*
		if !n.NodeKeeper.IsBootstrapped() {
			n.Controller.SetLastIgnoredPulse(newPulse.NextPulseNumber)
			return
		}
		if n.shoudIgnorePulse(newPulse) {
			log.Infof("Ignore pulse %d: network is not yet initialized", newPulse.PulseNumber)
			return
		}
	*/

	// todo move to nonetwork gateway
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

	if err := n.gatewayer.Gateway().OnPulse(ctx, newPulse); err != nil {
		logger.Error(errors.Wrap(err, "Failed to call OnPulse on Gateway"))
	}

	logger.Debugf("Before set new current pulse number: %d", newPulse.PulseNumber)
	err := n.PulseManager.Set(ctx, newPulse, n.gatewayer.Gateway().GetState() == insolar.CompleteNetworkState)
	if err != nil {
		logger.Fatalf("Failed to set new pulse: %s", err.Error())
	}
	logger.Infof("Set new current pulse number: %d", newPulse.PulseNumber)

	go n.phaseManagerOnPulse(ctx, newPulse, pulseTime)
}

// func (n *ServiceNetwork) shoudIgnorePulse(newPulse insolar.Pulse) bool {
// 	return n.isDiscovery && !n.NodeKeeper.GetConsensusInfo().IsJoiner() &&
// 		newPulse.PulseNumber <= n.Controller.GetLastIgnoredPulse()+insolar.PulseNumber(n.skip)
// }

func (n *ServiceNetwork) phaseManagerOnPulse(ctx context.Context, newPulse insolar.Pulse, pulseStartTime time.Time) {
	logger := inslogger.FromContext(ctx)

	if !n.cfg.Service.ConsensusEnabled {
		logger.Warn("Consensus is disabled")
		return
	}

	if err := n.PhaseManager.OnPulse(ctx, &newPulse, pulseStartTime); err != nil {
		errMsg := "Failed to pass consensus: " + err.Error()
		logger.Error(errMsg)
		n.gatewayer.SwitchState(insolar.NoNetworkState)
	}
}

func isNextPulse(currentPulse, newPulse *insolar.Pulse) bool {
	return newPulse.PulseNumber > currentPulse.PulseNumber && newPulse.PulseNumber >= currentPulse.NextPulseNumber
}

func (n *ServiceNetwork) GetState() insolar.NetworkState {
	return n.gatewayer.Gateway().GetState()
}

func (n *ServiceNetwork) SetOperableFunc(f insolar.NetworkOperableCallback) {
	n.operableFunc = f
}
