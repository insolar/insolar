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

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/routing"
	"github.com/pkg/errors"
)

// ServiceNetwork is facade for network.
type ServiceNetwork struct {
	hostNetwork  network.HostNetwork
	controller   network.Controller
	routingTable network.RoutingTable

	certificate  core.Certificate
	nodeNetwork  core.NodeNetwork
	pulseManager core.PulseManager
	coordinator  core.NetworkCoordinator
}

// NewServiceNetwork returns a new ServiceNetwork.
func NewServiceNetwork(conf configuration.Configuration) (*ServiceNetwork, error) {
	serviceNetwork := &ServiceNetwork{}

	routingTable, hostnetwork, controller, err := NewNetworkComponents(conf)
	if err != nil {
		log.Error("failed to create network components: %s", err.Error())
	}
	serviceNetwork.routingTable = routingTable
	serviceNetwork.hostNetwork = hostnetwork
	serviceNetwork.controller = controller
	return serviceNetwork, nil
}

// GetAddress returns host public address.
func (n *ServiceNetwork) GetAddress() string {
	return n.hostNetwork.PublicAddress()
}

// GetNodeID returns current node id.
func (n *ServiceNetwork) GetNodeID() core.RecordRef {
	return n.nodeNetwork.GetOrigin().ID()
}

// SendParcel sends a message from MessageBus.
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

// GetPrivateKey returns a private key.
// TODO: remove, use helper functions from certificate instead
func (n *ServiceNetwork) GetPrivateKey() *ecdsa.PrivateKey {
	return n.certificate.GetEcdsaPrivateKey()
}

// Start implements core.Component
func (n *ServiceNetwork) Start(ctx context.Context, components core.Components) error {
	n.inject(components)
	n.routingTable.Start(components)
	log.Infoln("Network starts listening")
	n.hostNetwork.Start()

	n.controller.Inject(components)

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
func (n *ServiceNetwork) Stop(ctx context.Context) error {
	n.hostNetwork.Stop()
	return nil
}

func (n *ServiceNetwork) bootstrap() {
	err := n.controller.Bootstrap()
	if err != nil {
		log.Errorln("Failed to bootstrap network", err.Error())
	}
}

func (n *ServiceNetwork) onPulse(pulse core.Pulse) {
	ctx := context.TODO()
	log.Infof("Got new pulse number: %d", pulse.PulseNumber)
	if n.pulseManager == nil {
		log.Error("PulseManager is not initialized")
		return
	}
	currentPulse, err := n.pulseManager.Current(ctx)
	if err != nil {
		log.Error(errors.Wrap(err, "Could not get current pulse"))
		return
	}
	if (pulse.PulseNumber > currentPulse.PulseNumber) &&
		(pulse.PulseNumber >= currentPulse.NextPulseNumber) {
		err = n.pulseManager.Set(ctx, pulse)
		if err != nil {
			log.Error(errors.Wrap(err, "Failed to set pulse"))
			return
		}
		log.Infof("Set new current pulse number: %d", pulse.PulseNumber)

		// deprecated
		/*
			go func(network *ServiceNetwork) {
				network.controller.ResendPulseToKnownHosts(pulse)
				if network.coordinator == nil {
					return
				}
				err := network.coordinator.WriteActiveNodes(ctx, pulse.PulseNumber, network.nodeNetwork.GetActiveNodes())
				if err != nil {
					log.Warn("Writing active nodes to ledger: " + err.Error())
				}
			}(n)
		*/

		// TODO: PLACE NEW CONSENSUS HERE
	}
}

// NewNetworkComponents create network.HostNetwork and network.Controller for new network
func NewNetworkComponents(conf configuration.Configuration) (network.RoutingTable, network.HostNetwork, network.Controller, error) {
	routingTable := routing.NewTable()
	internalTransport, err := hostnetwork.NewInternalTransport(conf)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "error creating internal transport")
	}
	hostNetwork := hostnetwork.NewHostTransport(internalTransport, routingTable)
	options := controller.ConfigureOptions(conf.Host)
	networkController := controller.NewNetworkController(options, internalTransport, routingTable, hostNetwork)
	return routingTable, hostNetwork, networkController, nil
}
