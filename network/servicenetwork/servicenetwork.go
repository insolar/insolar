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
	"strconv"

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

	fakePulsar *fakepulsar.FakePulsar
}

// NewServiceNetwork returns a new ServiceNetwork.
func NewServiceNetwork(conf configuration.Configuration, scheme core.PlatformCryptographyScheme) (*ServiceNetwork, error) {
	serviceNetwork := &ServiceNetwork{cfg: conf, CryptographyScheme: scheme}
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

	n.ConsensusNetwork, err = hostnetwork.NewConsensusNetwork(
		n.NodeNetwork.GetOrigin().ConsensusAddress(),
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

	err = n.MerkleCalculator.(component.Initer).Init(ctx)
	n.hostNetwork = hostnetwork.NewHostTransport(internalTransport, n.routingTable)
	options := controller.ConfigureOptions(n.cfg.Host)
	n.controller = controller.NewNetworkController(n, options, n.CertificateManager.GetCertificate(), internalTransport, n.routingTable, n.hostNetwork, n.CryptographyScheme)
	n.fakePulsar = fakepulsar.NewFakePulsar(n, n.cfg.Pulsar.PulseTime)
	log.Info("==== ServiceNetwork inited ====")

	return err
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
	log.Info("----- ServiceNetwork started ------")
	return nil
}

// Stop implements core.Component
func (n *ServiceNetwork) Stop(ctx context.Context) error {
	n.hostNetwork.Stop()
	n.ConsensusNetwork.Stop()
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
	currentPulse, err := n.PulseStorage.Current(ctx)
	if err != nil {
		logger.Error(errors.Wrap(err, "Could not get current pulse"))
		return
	}
	if (pulse.PulseNumber > currentPulse.PulseNumber) &&
		(pulse.PulseNumber >= currentPulse.NextPulseNumber) {
		err = n.PulseManager.Set(ctx, pulse, n.NetworkSwitcher.GetState() == core.CompleteNetworkState)
		if err != nil {
			logger.Error(errors.Wrap(err, "Failed to set pulse"))
			return
		}

		//err = n.NetworkSwitcher.OnPulse(ctx, pulse)
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
			log.Info("ServiceNetwork call PhaseManager.OnPulse")

			err = n.PhaseManager.OnPulse(ctx, &pulse)
			if err != nil {
				logger.Warn("phase manager fail: " + err.Error())
			}
		}(logger, n)
	} else {
		logger.Infof("Incorrect pulse number. Current: %d. New: %d", currentPulse.PulseNumber, pulse.PulseNumber)
	}
}

func (n *ServiceNetwork) isFakePulse(pulse *core.Pulse) bool {
	return (pulse.NextPulseNumber == 0) && (pulse.PulseNumber == 0)
}
