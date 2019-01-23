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
	"github.com/insolar/insolar/consensus/phases"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/merkle"
	"github.com/insolar/insolar/network/routing"
	"github.com/pkg/errors"
)

// ServiceNetwork is facade for network.
type ServiceNetwork struct {
	cfg configuration.Configuration

	hostNetwork  network.HostNetwork  // TODO: should be injected
	controller   network.Controller   // TODO: should be injected
	routingTable network.RoutingTable // TODO: should be injected

	// dependencies
	CertificateManager  core.CertificateManager         `inject:""`
	NodeNetwork         core.NodeNetwork                `inject:""`
	PulseManager        core.PulseManager               `inject:""`
	PulseStorage        core.PulseStorage               `inject:""`
	CryptographyService core.CryptographyService        `inject:""`
	NetworkCoordinator  core.NetworkCoordinator         `inject:""`
	ArtifactManager     core.ArtifactManager            `inject:""`
	CryptographyScheme  core.PlatformCryptographyScheme `inject:""`
	NodeKeeper          network.NodeKeeper              `inject:""`
	NetworkSwitcher     core.NetworkSwitcher            `inject:""`

	// subcomponents
	PhaseManager     phases.PhaseManager      // `inject:""`
	MerkleCalculator merkle.Calculator        // `inject:""`
	ConsensusNetwork network.ConsensusNetwork // `inject:""`
	PulseHandler     network.PulseHandler
	Communicator     phases.Communicator

	// fakePulsar *fakepulsar.FakePulsar
	isGenesis bool
	skip      int
}

// NewServiceNetwork returns a new ServiceNetwork.
func NewServiceNetwork(conf configuration.Configuration, scheme core.PlatformCryptographyScheme, isGenesis bool) (*ServiceNetwork, error) {
	serviceNetwork := &ServiceNetwork{cfg: conf, CryptographyScheme: scheme, isGenesis: isGenesis, skip: conf.Service.Skip}
	return serviceNetwork, nil
}

// SendMessage sends a message from MessageBus.
func (n *ServiceNetwork) SendMessage(nodeID core.RecordRef, method string, msg core.Parcel) ([]byte, error) {
	return n.controller.SendMessage(nodeID, method, msg)
}

// SendCascadeMessage sends a message from MessageBus to a cascade of nodes
func (n *ServiceNetwork) SendCascadeMessage(data core.Cascade, method string, msg core.Parcel) error {
	return n.controller.SendCascadeMessage(data, method, msg)
}

// RemoteProcedureRegister registers procedure for remote call on this host.
func (n *ServiceNetwork) RemoteProcedureRegister(name string, method core.RemoteProcedure) {
	n.controller.RemoteProcedureRegister(name, method)
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

	n.PhaseManager = phases.NewPhaseManager()
	n.MerkleCalculator = merkle.NewCalculator()
	n.Communicator = phases.NewNaiveCommunicator()
	n.PulseHandler = n // self

	firstPhase := &phases.FirstPhase{}
	secondPhase := &phases.SecondPhase{}
	thirdPhase := &phases.ThirdPhase{}

	// inject workaround
	n.PhaseManager.(*phases.Phases).FirstPhase = firstPhase
	n.PhaseManager.(*phases.Phases).SecondPhase = secondPhase
	n.PhaseManager.(*phases.Phases).ThirdPhase = thirdPhase

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

	n.ConsensusNetwork, err = hostnetwork.NewConsensusNetwork(
		n.cfg.Host.Transport.Address,
		n.CertificateManager.GetCertificate().GetNodeRef().String(),
		n.NodeNetwork.GetOrigin().ShortID(),
		n.routingTable,
	)
	if err != nil {
		return errors.Wrap(err, "Failed to create consensus network.")
	}

	cm := component.Manager{}
	cm.Register(n.CertificateManager, n.NodeNetwork, n.PulseManager, n.CryptographyService, n.NetworkCoordinator,
		n.ArtifactManager, n.CryptographyScheme, n.PulseHandler)

	cm.Inject(n.NodeKeeper,
		n.MerkleCalculator,
		n.ConsensusNetwork,
		n.Communicator,
		firstPhase,
		secondPhase,
		thirdPhase,
	)

	n.hostNetwork = hostnetwork.NewHostTransport(internalTransport, n.routingTable)
	options := controller.ConfigureOptions(n.cfg.Host)
	n.controller = controller.NewNetworkController(n, options, n.CertificateManager.GetCertificate(), internalTransport, n.routingTable, n.hostNetwork, n.CryptographyScheme)
	// n.fakePulsar = fakepulsar.NewFakePulsar(n.HandlePulse, n.cfg.Pulsar.PulseTime)
	return nil
}

// Start implements component.Starter
func (n *ServiceNetwork) Start(ctx context.Context) error {
	log.Infoln("Network starts listening...")
	n.hostNetwork.Start(ctx)

	n.controller.Inject(n.CryptographyService, n.NetworkCoordinator, n.NodeKeeper)
	n.routingTable.Inject(n.NodeKeeper)

	log.Infoln("Bootstrapping network...")
	err := n.controller.Bootstrap(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to bootstrap network")
	}

	// n.fakePulsar.Start(ctx)

	return nil
}

// Stop implements core.Component
func (n *ServiceNetwork) Stop(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)

	logger.Info("Stopping host network")
	n.hostNetwork.Stop()
	logger.Info("Stopping consensus network")
	n.ConsensusNetwork.Stop()
	return nil
}

func (n *ServiceNetwork) HandlePulse(ctx context.Context, pulse core.Pulse) {
	// if !n.isFakePulse(&pulse) {
	// 	n.fakePulsar.Stop(ctx)
	// }
	if n.isGenesis {
		return
	}

	traceID := "pulse_" + strconv.FormatUint(uint64(pulse.PulseNumber), 10)

	ctx, logger := inslogger.WithTraceField(ctx, traceID)
	logger.Infof("Got new pulse number: %d", pulse.PulseNumber)
	if n.PulseManager == nil {
		logger.Error("PulseManager is not initialized")
		return
	}
	if !n.NodeKeeper.IsBootstrapped() {
		n.controller.SetLastIgnoredPulse(pulse.NextPulseNumber)
		return
	}
	if pulse.PulseNumber <= n.controller.GetLastIgnoredPulse()+core.PulseNumber(n.skip) {
		log.Infof("Ignore pulse %d: network is not yet initialized", pulse.PulseNumber)
		return
	}
	currentPulse, err := n.PulseStorage.Current(ctx)
	if err != nil {
		logger.Error(errors.Wrap(err, "Could not get current pulse"))
		return
	}
	if (pulse.PulseNumber > currentPulse.PulseNumber) &&
		(pulse.PulseNumber >= currentPulse.NextPulseNumber) {

		err = n.NetworkSwitcher.OnPulse(ctx, pulse)
		if err != nil {
			logger.Error(errors.Wrap(err, "Failed to call OnPulse on NetworkSwitcher"))
			return
		}

		err = n.PulseManager.Set(ctx, pulse, n.NetworkSwitcher.GetState() == core.CompleteNetworkState)
		if err != nil {
			logger.Error(errors.Wrap(err, "Failed to set pulse"))
			return
		}

		logger.Infof("Set new current pulse number: %d", pulse.PulseNumber)
		// go func(logger core.Logger, network *ServiceNetwork) {
		// 	TODO: make PhaseManager works and uncomment this (after NETD18-75)
		// 	err = n.PhaseManager.OnPulse(ctx, &pulse)
		// 	if err != nil {
		// 		logger.Warn("phase manager fail: " + err.Error())
		// 	}
		// }(logger, n)
	} else {
		logger.Infof("Incorrect pulse number. Current: %d. New: %d", currentPulse.PulseNumber, pulse.PulseNumber)
	}
}

func (n *ServiceNetwork) isFakePulse(pulse *core.Pulse) bool {
	return (pulse.NextPulseNumber == 0) && (pulse.PulseNumber == 0)
}
