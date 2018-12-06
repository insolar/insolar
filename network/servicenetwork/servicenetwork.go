/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package servicenetwork

import (
	"context"
	"fmt"
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
	"github.com/insolar/insolar/network/fakepulsar"
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

	fakePulsar *fakepulsar.FakePulsar
}

// NewServiceNetwork returns a new ServiceNetwork.
func NewServiceNetwork(conf configuration.Configuration, scheme core.PlatformCryptographyScheme) (*ServiceNetwork, error) {
	serviceNetwork := &ServiceNetwork{cfg: conf, CryptographyScheme: scheme}
	return serviceNetwork, nil
}

// GetAddress returns host public address.
func (n *ServiceNetwork) GetAddress() string {
	return n.hostNetwork.PublicAddress()
}

// GetNodeID returns current node id.
func (n *ServiceNetwork) GetNodeID() core.RecordRef {
	return n.NodeNetwork.GetOrigin().ID()
}

// GetGlobuleID returns current globule id.
func (n *ServiceNetwork) GetGlobuleID() core.GlobuleID {
	return 0
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
	if len(parts) != 2 {
		return address, errors.New("failed to get port from address")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return address, err
	}

	if port != 0 {
		port++
	}
	return fmt.Sprintf("%s:%d", parts[0], port), nil
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
		n.PhaseManager,
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
	n.fakePulsar = fakepulsar.NewFakePulsar(n.HandlePulse, n.cfg.Pulsar.PulseTime)
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

	n.fakePulsar.Start(ctx)

	return nil
}

// Stop implements core.Component
func (n *ServiceNetwork) Stop(ctx context.Context) error {
	n.hostNetwork.Stop()
	return nil
}

func (n *ServiceNetwork) HandlePulse(ctx context.Context, pulse core.Pulse) {
	if !n.isFakePulse(&pulse) {
		n.fakePulsar.Stop(ctx)
	}
	traceID := "pulse_" + strconv.FormatUint(uint64(pulse.PulseNumber), 10)
	ctx, logger := inslogger.WithTraceField(ctx, traceID)
	logger.Infof("Got new pulse number: %d", pulse.PulseNumber)
	if n.PulseManager == nil {
		logger.Error("PulseManager is not initialized")
		return
	}
	currentPulse, err := n.PulseManager.Current(ctx)
	if err != nil {
		logger.Error(errors.Wrap(err, "Could not get current pulse"))
		return
	}
	if (pulse.PulseNumber > currentPulse.PulseNumber) &&
		(pulse.PulseNumber >= currentPulse.NextPulseNumber) {
		err = n.PulseManager.Set(ctx, pulse, false)
		if err != nil {
			logger.Error(errors.Wrap(err, "Failed to set pulse"))
			return
		}

		// TODO: I don't know why I put it here. If you know better place for that, move it there please
		err = n.NetworkSwitcher.OnPulse(ctx, pulse)
		if err != nil {
			logger.Error(errors.Wrap(err, "Failed to call OnPulse on NetworkSwitcher"))
			return
		}

		logger.Infof("Set new current pulse number: %d", pulse.PulseNumber)
		go func(logger core.Logger, network *ServiceNetwork) {
			if network.NetworkCoordinator == nil {
				return
			}
			err := network.NetworkCoordinator.WriteActiveNodes(ctx, pulse.PulseNumber, network.NodeNetwork.GetActiveNodes())
			if err != nil {
				logger.Warn("Error writing active nodes to ledger: " + err.Error())
			}
			// TODO: make PhaseManager works and uncomment this
			// err = n.PhaseManager.OnPulse(ctx, &pulse)
			// if err != nil {
			// 	logger.Warn("phase manager fail: " + err.Error())
			// }
		}(logger, n)
	} else {
		logger.Infof("Incorrect pulse number. Current: %d. New: %d", currentPulse.PulseNumber, pulse.PulseNumber)
	}
}

func (n *ServiceNetwork) isFakePulse(pulse *core.Pulse) bool {
	return (pulse.NextPulseNumber == 0) && (pulse.PulseNumber == 0)
}
