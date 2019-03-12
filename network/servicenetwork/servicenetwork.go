/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package servicenetwork

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/consensus/phases"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller"
	"github.com/insolar/insolar/network/controller/bootstrap"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/merkle"
	"github.com/insolar/insolar/network/routing"
	"github.com/insolar/insolar/network/utils"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// ServiceNetwork is facade for network.
type ServiceNetwork struct {
	cfg configuration.Configuration
	cm  *component.Manager

	hostNetwork  network.HostNetwork  // TODO: should be injected
	routingTable network.RoutingTable // TODO: should be injected

	// dependencies
	CertificateManager  core.CertificateManager         `inject:""`
	PulseManager        core.PulseManager               `inject:""`
	PulseStorage        core.PulseStorage               `inject:""`
	CryptographyService core.CryptographyService        `inject:""`
	NetworkCoordinator  core.NetworkCoordinator         `inject:""`
	CryptographyScheme  core.PlatformCryptographyScheme `inject:""`
	NodeKeeper          network.NodeKeeper              `inject:""`
	NetworkSwitcher     core.NetworkSwitcher            `inject:""`
	TerminationHandler  core.TerminationHandler         `inject:""`

	// subcomponents
	PhaseManager phases.PhaseManager `inject:"subcomponent"`
	Controller   network.Controller  `inject:"subcomponent"`

	isGenesis   bool
	isDiscovery bool
	skip        int

	lock sync.Mutex
}

// NewServiceNetwork returns a new ServiceNetwork.
func NewServiceNetwork(conf configuration.Configuration, rootCm *component.Manager, isGenesis bool) (*ServiceNetwork, error) {
	serviceNetwork := &ServiceNetwork{cm: component.NewManager(rootCm), cfg: conf, isGenesis: isGenesis, skip: conf.Service.Skip}
	return serviceNetwork, nil
}

// SendMessage sends a message from MessageBus.
func (n *ServiceNetwork) SendMessage(nodeID core.RecordRef, method string, msg core.Parcel) ([]byte, error) {
	return n.Controller.SendMessage(nodeID, method, msg)
}

// SendCascadeMessage sends a message from MessageBus to a cascade of nodes
func (n *ServiceNetwork) SendCascadeMessage(data core.Cascade, method string, msg core.Parcel) error {
	return n.Controller.SendCascadeMessage(data, method, msg)
}

// RemoteProcedureRegister registers procedure for remote call on this host.
func (n *ServiceNetwork) RemoteProcedureRegister(name string, method core.RemoteProcedure) {
	n.Controller.RemoteProcedureRegister(name, method)
}

// incrementPort increments port number if it not equals 0
func incrementPort(address string) (string, error) {
	parts := strings.Split(address, ":")
	if len(parts) < 2 {
		return address, errors.New("failed to get port from address " + address)
	}
	port, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return address, err
	}

	if port != 0 {
		port++
	}

	parts = append(parts[:len(parts)-1], strconv.Itoa(port))
	return strings.Join(parts, ":"), nil
}

// Start implements component.Initer
func (n *ServiceNetwork) Init(ctx context.Context) error {
	n.routingTable = &routing.Table{}
	internalTransport, err := hostnetwork.NewInternalTransport(n.cfg, n.CertificateManager.GetCertificate().GetNodeRef().String())
	if err != nil {
		return errors.Wrap(err, "Failed to create internal transport")
	}

	var consensusAddress string
	if n.cfg.Host.Transport.FixedPublicAddress != "" {
		// workaround for Consensus transport, port+=1 of default transport
		consensusAddress, err = incrementPort(n.cfg.Host.Transport.Address)
		if err != nil {
			return errors.Wrap(err, "failed to increment port.")
		}
	} else {
		consensusAddress = n.NodeKeeper.GetOrigin().ConsensusAddress()
	}

	consensusNetwork, err := hostnetwork.NewConsensusNetwork(
		consensusAddress,
		n.CertificateManager.GetCertificate().GetNodeRef().String(),
		n.NodeKeeper.GetOrigin().ShortID(),
		n.routingTable,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to create consensus network.")
	}

	n.hostNetwork = hostnetwork.NewHostTransport(internalTransport, n.routingTable)
	options := controller.ConfigureOptions(n.cfg)

	cert := n.CertificateManager.GetCertificate()
	n.isDiscovery = utils.OriginIsDiscovery(cert)

	n.cm.Inject(n,
		cert,
		n.NodeKeeper,
		merkle.NewCalculator(),
		consensusNetwork,
		phases.NewCommunicator(),
		phases.NewFirstPhase(),
		phases.NewSecondPhase(),
		phases.NewThirdPhase(),
		phases.NewPhaseManager(),
		bootstrap.NewSessionManager(),
		controller.NewNetworkController(n.hostNetwork),
		controller.NewRPCController(options, n.hostNetwork),
		controller.NewPulseController(n.hostNetwork, n.routingTable),
		bootstrap.NewBootstrapper(options, internalTransport),
		bootstrap.NewAuthorizationController(options, internalTransport),
		bootstrap.NewChallengeResponseController(options, internalTransport),
		bootstrap.NewNetworkBootstrapper(),
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

	logger.Info("Network starts listening...")
	n.routingTable.Inject(n.NodeKeeper)
	n.hostNetwork.Start(ctx)

	log.Info("Starting network component manager...")
	err := n.cm.Start(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to bootstrap network")
	}

	log.Info("Bootstrapping network...")
	_, err = n.Controller.Bootstrap(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to bootstrap network")
	}

	logger.Info("Service network started")
	return nil
}

func (n *ServiceNetwork) Leave(ctx context.Context, ETA core.PulseNumber) {
	logger := inslogger.FromContext(ctx)
	logger.Info("Gracefully stopping service network")

	n.NodeKeeper.AddPendingClaim(&packets.NodeLeaveClaim{ETA: ETA})
}

// Stop implements core.Component
func (n *ServiceNetwork) Stop(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	logger.Info("Stopping network components")
	if err := n.cm.Stop(ctx); err != nil {
		log.Errorf("Error while stopping network components: %s", err.Error())
	}
	logger.Info("Stopping host network")
	n.hostNetwork.Stop()
	return nil
}

func (n *ServiceNetwork) HandlePulse(ctx context.Context, newPulse core.Pulse) {
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
	if n.isDiscovery && newPulse.PulseNumber <= n.Controller.GetLastIgnoredPulse()+core.PulseNumber(n.skip) {
		log.Infof("Ignore pulse %d: network is not yet initialized", newPulse.PulseNumber)
		return
	}

	if n.NodeKeeper.GetState() == core.WaitingNodeNetworkState {
		// do not set pulse because otherwise we will set invalid active list
		// pass consensus, prepare valid active list and set it on next pulse
		go n.phaseManagerOnPulse(ctx, newPulse, currentTime)
		return
	}

	// Ignore core.ErrNotFound because
	// sometimes we can't fetch current pulse in new nodes
	// (for fresh bootstrapped light-material with in-memory pulse-tracker)
	if currentPulse, err := n.PulseStorage.Current(ctx); err != nil {
		if err != core.ErrNotFound {
			logger.Fatalf("Could not get current pulse: %s", err.Error())
		}
	} else {
		if !isNextPulse(currentPulse, &newPulse) {
			logger.Infof("Incorrect pulse number. Current: %+v. New: %+v", currentPulse, newPulse)
			return
		}
	}

	err := n.NetworkSwitcher.OnPulse(ctx, newPulse)
	if err != nil {
		logger.Error(errors.Wrap(err, "Failed to call OnPulse on NetworkSwitcher"))
	}

	logger.Debugf("Before set new current pulse number: %d", newPulse.PulseNumber)
	err = n.PulseManager.Set(ctx, newPulse, n.NetworkSwitcher.GetState() == core.CompleteNetworkState)
	if err != nil {
		logger.Fatalf("Failed to set new pulse: %s", err.Error())
	}
	logger.Infof("Set new current pulse number: %d", newPulse.PulseNumber)

	go n.phaseManagerOnPulse(ctx, newPulse, currentTime)
}

func (n *ServiceNetwork) phaseManagerOnPulse(ctx context.Context, newPulse core.Pulse, pulseStartTime time.Time) {
	logger := inslogger.FromContext(ctx)

	if err := n.PhaseManager.OnPulse(ctx, &newPulse, pulseStartTime); err != nil {
		logger.Error("Failed to pass consensus: " + err.Error())
		n.TerminationHandler.Abort()
	}
}

func isNextPulse(currentPulse, newPulse *core.Pulse) bool {
	return newPulse.PulseNumber > currentPulse.PulseNumber && newPulse.PulseNumber >= currentPulse.NextPulseNumber
}
