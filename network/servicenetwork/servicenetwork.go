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

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/consensus/packets"
	"github.com/insolar/insolar/consensus/phases"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller"
	"github.com/insolar/insolar/network/controller/bootstrap"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/merkle"
	"github.com/insolar/insolar/network/routing"
	"github.com/insolar/insolar/network/utils"
	"github.com/pkg/errors"
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

	// subcomponents
	PhaseManager phases.PhaseManager `inject:"subcomponent"`
	Controller   network.Controller  `inject:"subcomponent"`

	// fakePulsar *fakepulsar.FakePulsar
	isGenesis   bool
	isDiscovery bool
	skip        int
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
		return address, errors.New("failed to get port from address")
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

	// workaround for Consensus transport, port+=1 of default transport
	n.cfg.Host.Transport.Address, err = incrementPort(n.cfg.Host.Transport.Address)
	if err != nil {
		return errors.Wrap(err, "failed to increment port.")
	}

	consensusNetwork, err := hostnetwork.NewConsensusNetwork(
		n.NodeKeeper.GetOrigin().ConsensusAddress(),
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
		phases.NewNaiveCommunicator(),
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

	// n.fakePulsar = fakepulsar.NewFakePulsar(n, time.Duration(n.cfg.Pulsar.PulseTime)*time.Millisecond)
	return nil
}

// Start implements component.Starter
func (n *ServiceNetwork) Start(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	logger.Infoln("Network starts listening...")
	n.routingTable.Inject(n.NodeKeeper)
	n.hostNetwork.Start(ctx)

	log.Info("Starting network component manager...")
	err := n.cm.Start(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to bootstrap network")
	}

	log.Infoln("Bootstrapping network...")
	_, err = n.Controller.Bootstrap(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to bootstrap network")
	}
	// if !n.isGenesis {
	// 	n.fakePulsar.Start(ctx, result.FirstPulseTime)
	// }
	logger.Info("Service network started")
	return nil
}

func (n *ServiceNetwork) GracefulStop(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	logger.Info("Gracefully stopping service network")

	n.NodeKeeper.AddPendingClaim(&packets.NodeLeaveClaim{})
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
	if n.isGenesis {
		return
	}
	traceID := "pulse_" + strconv.FormatUint(uint64(newPulse.PulseNumber), 10)
	ctx, logger := inslogger.WithTraceField(ctx, traceID)
	logger.Infof("Got new pulse number: %d", newPulse.PulseNumber)

	if !n.NodeKeeper.IsBootstrapped() {
		n.Controller.SetLastIgnoredPulse(newPulse.NextPulseNumber)
		return
	}
	if n.isDiscovery && newPulse.PulseNumber <= n.Controller.GetLastIgnoredPulse()+core.PulseNumber(n.skip) {
		log.Infof("Ignore pulse %d: network is not yet initialized", newPulse.PulseNumber)
		return
	}

	currentPulse, err := n.PulseStorage.Current(ctx)
	if err != nil {
		logger.Error(errors.Wrap(err, "Could not get current pulse"))
		return
	}

	// Working on early network state, ready for fake pulses
	// if isFakePulse(&newPulse) && !fakePulseAllowed(n.NetworkSwitcher.GetState()) {
	// 	logger.Infof("Got fake pulse on invalid network state. Current: %+v. New: %+v", currentPulse, newPulse)
	// 	return
	// }

	if !isNextPulse(currentPulse, &newPulse) {
		logger.Infof("Incorrect newPulse number. Current: %+v. New: %+v", currentPulse, newPulse)
		return
	}

	// Got real pulse
	// if isFakePulse(currentPulse) && !isFakePulse(&newPulse) {
	// 	n.fakePulsar.Stop(ctx)
	// }

	err = n.NetworkSwitcher.OnPulse(ctx, newPulse)
	if err != nil {
		logger.Error(errors.Wrap(err, "Failed to call OnPulse on NetworkSwitcher"))
		return
	}

	err = n.PulseManager.Set(ctx, newPulse, n.NetworkSwitcher.GetState() == core.CompleteNetworkState)
	if err != nil {
		logger.Error(errors.Wrap(err, "Failed to set newPulse"))
		return
	}

	logger.Infof("Set new current pulse number: %d", newPulse.PulseNumber)
	go n.phaseManagerOnPulse(ctx, newPulse)
}

func (n *ServiceNetwork) phaseManagerOnPulse(ctx context.Context, newPulse core.Pulse) {
	logger := inslogger.FromContext(ctx)

	if err := n.PhaseManager.OnPulse(ctx, &newPulse); err != nil {
		logger.Warn("phase manager fail: " + err.Error())
	}
}

func isFakePulse(newPulse *core.Pulse) bool {
	return newPulse.EpochPulseNumber == -1
}

func isNewEpoch(currentPulse, newPulse *core.Pulse) bool {
	return newPulse.EpochPulseNumber > currentPulse.EpochPulseNumber
}

func isNextPulse(currentPulse, newPulse *core.Pulse) bool {
	return newPulse.PulseNumber > currentPulse.PulseNumber && newPulse.PulseNumber >= currentPulse.NextPulseNumber
}

func fakePulseStarted(currentPulse, newPulse *core.Pulse) bool {
	return isFakePulse(newPulse) && currentPulse.EpochPulseNumber > -1
}

func fakePulseAllowed(state core.NetworkState) bool {
	return state == core.VoidNetworkState || state == core.NoNetworkState
}
