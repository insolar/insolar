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

	certificate  core.Certificate
	nodeNetwork  core.NodeNetwork
	pulseManager core.PulseManager
	coordinator  core.NetworkCoordinator
}

// NewServiceNetwork returns a new ServiceNetwork.
func NewServiceNetwork(conf configuration.Configuration) (*ServiceNetwork, error) {
	network := &ServiceNetwork{}

	// workaround before DI
	cert, err := certificate.NewCertificate(conf.KeysPath, conf.CertificatePath)
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
func (n *ServiceNetwork) GetAddress() string {
	return n.hostNetwork.PublicAddress()
}

// GetNodeID returns current node id.
func (n *ServiceNetwork) GetNodeID() core.RecordRef {
	return n.nodeNetwork.GetOrigin().NodeID
}

// SendMessage sends a message from MessageBus.
func (n *ServiceNetwork) SendMessage(nodeID core.RecordRef, method string, msg core.Message) ([]byte, error) {
	return n.controller.SendMessage(nodeID, method, msg)
}

// SendCascadeMessage sends a message from MessageBus to a cascade of nodes
func (n *ServiceNetwork) SendCascadeMessage(data core.Cascade, method string, msg core.Message) error {
	return n.controller.SendCascadeMessage(data, method, msg)
}

// RemoteProcedureRegister registers procedure for remote call on this host.
func (n *ServiceNetwork) RemoteProcedureRegister(name string, method core.RemoteProcedure) {
	n.controller.RemoteProcedureRegister(name, method)
}

// GetHostNetwork returns pointer to host network layer(DHT), temp method, refactoring needed
// TODO: replace with GetNetworkHelper that returns a component with all needed data for interactive/rest API
func (n *ServiceNetwork) GetHostNetwork() (hosthandler.HostHandler, hosthandler.Context) {
	hostNetwork := n.hostNetwork.(*dhtnetwork.Wrapper).HostNetwork
	return hostNetwork, dhtnetwork.CreateDHTContext(hostNetwork)
}

// GetPrivateKey returns a private key.
// TODO: remove, use helper functions from certificate instead
func (n *ServiceNetwork) GetPrivateKey() *ecdsa.PrivateKey {
	return n.certificate.GetEcdsaPrivateKey()
}

// Start implements core.Component
func (n *ServiceNetwork) Start(insctx core.Context, components core.Components) error {
	n.inject(components)
	go n.listen()

	n.controller.Inject(components)
	n.consensus.SetNodeKeeper(components.NodeNetwork.(network.NodeKeeper))

	log.Infoln("Bootstrapping network...")
	n.bootstrap()

	err := n.controller.AnalyzeNetwork()
	if err != nil {
		log.Error(err)
	}

	err = n.controller.Authorize()
	if err != nil {
		return errors.Wrap(err, "error authorizing node")
	}

	return nil
}

func (n *ServiceNetwork) inject(components core.Components) {
	n.certificate = components.Certificate
	n.nodeNetwork = components.NodeNetwork
	n.pulseManager = components.Ledger.GetPulseManager()
	n.coordinator = components.NetworkCoordinator
}

// Stop implements core.Component
func (n *ServiceNetwork) Stop(insctx core.Context) error {
	return n.hostNetwork.Disconnect()
}

func (n *ServiceNetwork) bootstrap() {
	err := n.controller.Bootstrap()
	if err != nil {
		log.Errorln("Failed to bootstrap n", err.Error())
	}
}

func (n *ServiceNetwork) listen() {
	log.Infoln("Network starts listening")
	err := n.hostNetwork.Listen()
	if err != nil {
		log.Errorln("Listen failed:", err.Error())
	}
}

func (n *ServiceNetwork) onPulse(pulse core.Pulse) {
	if n.pulseManager == nil {
		log.Error("PulseManager is not initialized")
		return
	}
	currentPulse, err := n.pulseManager.Current()
	if err != nil {
		log.Error(errors.Wrap(err, "Could not get current pulse"))
		return
	}
	if (pulse.PulseNumber > currentPulse.PulseNumber) &&
		(pulse.PulseNumber >= currentPulse.NextPulseNumber) {
		err = n.pulseManager.Set(pulse)
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
			err := network.coordinator.WriteActiveNodes(pulse.PulseNumber, network.nodeNetwork.GetActiveNodes())
			if err != nil {
				log.Warn("Writing active nodes to ledger: " + err.Error())
			}
		}(n)

		// TODO: create adequate cancelable context without dht values (after switching to new n)
		ctx := context.WithValue(context.Background(), dhtnetwork.CtxTableIndex, dhtnetwork.DefaultHostID)
		n.doConsensus(ctx, pulse)
	}
}

func (n *ServiceNetwork) doConsensus(ctx hosthandler.Context, pulse core.Pulse) {
	if !n.consensus.IsPartOfConsensus() {
		log.Debug("Node is not active and does not participate in consensus")
		return
	}
	log.Debugf("Initiating consensus for pulse %d", pulse.PulseNumber)
	go n.consensus.ProcessPulse(ctx, pulse)
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
