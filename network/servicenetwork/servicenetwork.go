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
	"strings"
	"time"

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
	"github.com/insolar/insolar/network/fakepulsar"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/merkle"
	"github.com/insolar/insolar/network/routing"
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

	fakePulsar *fakepulsar.FakePulsar
	isGenesis  bool
	skip       int
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

	n.cm.Inject(n,
		n.CertificateManager.GetCertificate(),
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

	n.fakePulsar = fakepulsar.NewFakePulsar(n, time.Duration(n.cfg.Pulsar.PulseTime)*time.Millisecond)
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
	result, err := n.Controller.Bootstrap(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to bootstrap network")
	}
	if !n.isGenesis {
		n.fakePulsar.Start(ctx, result.FirstPulseTime)
	}
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
	if newPulse.PulseNumber <= n.Controller.GetLastIgnoredPulse()+core.PulseNumber(n.skip) {
		log.Infof("Ignore pulse %d: network is not yet initialized", newPulse.PulseNumber)
		return
	}

	currentPulse, err := n.PulseStorage.Current(ctx)
	if err != nil {
		logger.Error(errors.Wrap(err, "Could not get current pulse"))
		return
	}

	// Working on early network state, ready for fake pulses
	// TODO: !!!
	// if isFakePulse(&newPulse) && !fakePulseAllowed(n.NetworkSwitcher.GetState()) {
	// 	logger.Infof("Got fake pulse on invalid network state. Current: %+v. New: %+v", currentPulse, newPulse)
	// 	return
	// }

	if !isNextPulse(currentPulse, &newPulse) && !isNewEpoch(currentPulse, &newPulse) && !fakePulseStarted(currentPulse, &newPulse) {
		logger.Infof("Incorrect newPulse number. Current: %+v. New: %+v", currentPulse, newPulse)
		return
	}

	// Got real pulse
	if isFakePulse(currentPulse) && !isFakePulse(&newPulse) {
		n.fakePulsar.Stop(ctx)
	}

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
