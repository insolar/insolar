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
	"crypto/ecdsa"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/dhtnetwork"
	"github.com/insolar/insolar/network/dhtnetwork/hosthandler"
	"github.com/pkg/errors"
)

// ServiceNetwork is facade for network.
type ServiceNetwork struct {
	hostNetwork network.HostNetwork
	controller  network.Controller
	consensus   consensus.Processor

	certificate         core.Certificate
	activeNodeComponent core.ActiveNodeComponent
	pulseManager        core.PulseManager
	coordinator         core.NetworkCoordinator
}

// NewServiceNetwork returns a new ServiceNetwork.
func NewServiceNetwork(conf configuration.Configuration) (*ServiceNetwork, error) {
	network := &ServiceNetwork{}

	// workaround before DI
	cert, err := certificate.NewCertificate(conf.KeysPath)
	if err != nil {
		log.Warnf("failed to read certificate: %s", err.Error())
	}
	hostnetwork, err := NewHostNetwork(conf, cert, network.onPulse)
	if err != nil {
		log.Error("failed to create hostnetwork: %s", err.Error())
	}
	controller, err := NewNetworkController(conf, hostnetwork)
	if err != nil {
		log.Error("failed to create network controller: %s", err.Error())
	}
	network.hostNetwork = hostnetwork
	network.controller = controller
	network.certificate = cert
	network.consensus = NewConsensus(network.hostNetwork)
	return network, nil
}

// GetAddress returns host public address.
func (network *ServiceNetwork) GetAddress() string {
	return network.hostNetwork.PublicAddress()
}

// GetNodeID returns current node id.
func (network *ServiceNetwork) GetNodeID() core.RecordRef {
	return network.activeNodeComponent.GetID()
}

// SendMessage sends a message from MessageBus.
func (network *ServiceNetwork) SendMessage(nodeID core.RecordRef, method string, msg core.Message) ([]byte, error) {
	return network.controller.SendMessage(nodeID, method, msg)
}

// SendCascadeMessage sends a message from MessageBus to a cascade of nodes
func (network *ServiceNetwork) SendCascadeMessage(data core.Cascade, method string, msg core.Message) error {
	return network.controller.SendCascadeMessage(data, method, msg)
}

// RemoteProcedureRegister registers procedure for remote call on this host.
func (network *ServiceNetwork) RemoteProcedureRegister(name string, method core.RemoteProcedure) {
	network.controller.RemoteProcedureRegister(name, method)
}

// GetHostNetwork returns pointer to host network layer(DHT), temp method, refactoring needed
// TODO: replace with GetNetworkHelper that returns a component with all needed data for interactive/rest API
func (network *ServiceNetwork) GetHostNetwork() (hosthandler.HostHandler, hosthandler.Context) {
	n := network.hostNetwork.(*dhtnetwork.Wrapper)
	return n.HostNetwork, dhtnetwork.CreateDHTContext(n.HostNetwork)
}

// GetPrivateKey returns a private key.
// TODO: remove, use helper functions from certificate instead
func (network *ServiceNetwork) GetPrivateKey() *ecdsa.PrivateKey {
	return network.certificate.GetEcdsaPrivateKey()
}

// Start implements core.Component
func (network *ServiceNetwork) Start(components core.Components) error {
	network.inject(components)
	go network.listen()

	network.controller.Inject(components)

	log.Infoln("Bootstrapping network...")
	network.bootstrap()

	err := network.controller.AnalyzeNetwork()
	if err != nil {
		log.Error(err)
	}

	err = network.controller.Authorize()
	if err != nil {
		return errors.Wrap(err, "error authorizing node")
	}

	return nil
}

func (network *ServiceNetwork) inject(components core.Components) {
	network.certificate = components.Certificate
	network.activeNodeComponent = components.ActiveNodeComponent
	network.pulseManager = components.Ledger.GetPulseManager()
	network.coordinator = components.NetworkCoordinator
}

// Stop implements core.Component
func (network *ServiceNetwork) Stop() error {
	return network.hostNetwork.Disconnect()
}

func (network *ServiceNetwork) bootstrap() {
	err := network.controller.Bootstrap()
	if err != nil {
		log.Errorln("Failed to bootstrap network", err.Error())
	}
}

func (network *ServiceNetwork) listen() {
	log.Infoln("Network starts listening")
	err := network.hostNetwork.Listen()
	if err != nil {
		log.Errorln("Listen failed:", err.Error())
	}
}

func (network *ServiceNetwork) onPulse(pulse core.Pulse) {
	if network.pulseManager == nil {
		log.Error("PulseManager is not initialized")
		return
	}
	currentPulse, err := network.pulseManager.Current()
	if err != nil {
		log.Error(errors.Wrap(err, "Could not get current pulse"))
		return
	}
	if (pulse.PulseNumber > currentPulse.PulseNumber) &&
		(pulse.PulseNumber >= currentPulse.NextPulseNumber) {
		err = network.pulseManager.Set(pulse)
		if err != nil {
			log.Error(errors.Wrap(err, "Failed to set pulse"))
			return
		}
		log.Infof("Set new current pulse number: %d", pulse.PulseNumber)
		go func(network *ServiceNetwork) {
			network.controller.ResendPulseToKnownHosts(pulse)
			if network.coordinator == nil {
				return
			}
			err := network.coordinator.WriteActiveNodes(pulse.PulseNumber, network.activeNodeComponent.GetActiveNodes())
			if err != nil {
				log.Warn("Writing active nodes to ledger: " + err.Error())
			}
		}(network)

		// TODO: create adequate cancelable context without dht values (after switching to new network)
		ctx := context.WithValue(context.Background(), dhtnetwork.CtxTableIndex, dhtnetwork.DefaultHostID)
		network.doConsensus(ctx, pulse)
	}
}

func (network *ServiceNetwork) doConsensus(ctx hosthandler.Context, pulse core.Pulse) {
	if !network.consensus.IsPartOfConsensus() {
		log.Debug("Node is not active and does not participate in consensus")
		return
	}
	log.Debugf("Initiating consensus for pulse %d", pulse.PulseNumber)
	go network.consensus.ProcessPulse(ctx, pulse)
}

// NewHostNetwork create new HostNetwork. Certificate in new network should be removed and pulseCallback should be passed to NewNetworkController.
func NewHostNetwork(conf configuration.Configuration, certificate core.Certificate, pulseCallback network.OnPulse) (network.HostNetwork, error) {
	return dhtnetwork.NewDhtHostNetwork(conf, certificate, pulseCallback)
}

// NewNetworkController create new network.Controller. In new network it should read conf.
func NewNetworkController(conf configuration.Configuration, network network.HostNetwork) (network.Controller, error) {
	return dhtnetwork.NewDhtNetworkController(network)
}

func NewConsensus(network network.HostNetwork) consensus.Processor {
	return dhtnetwork.NewNetworkConsensus(network)
}
