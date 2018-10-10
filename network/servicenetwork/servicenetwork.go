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
	"io/ioutil"
	"strings"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/cascade"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/hostnetwork/hosthandler"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/pkg/errors"
)

// ServiceNetwork is facade for network.
type ServiceNetwork struct {
	nodeNetwork *nodenetwork.NodeNetwork
	hostNetwork hosthandler.HostHandler
}

// NewServiceNetwork returns a new ServiceNetwork.
func NewServiceNetwork(
	hostConf configuration.HostNetwork,
	nodeConf configuration.NodeNetwork,
) (*ServiceNetwork, error) {

	node := nodenetwork.NewNodeNetwork(nodeConf)
	if node == nil {
		return nil, errors.New("failed to create a node network")
	}

	cascade1 := &cascade.Cascade{}
	dht, err := hostnetwork.NewHostNetwork(hostConf, node, cascade1)
	if err != nil {
		return nil, err
	}

	service := &ServiceNetwork{nodeNetwork: node, hostNetwork: dht}
	f := func(data core.Cascade, method string, args [][]byte) error {
		return service.initCascadeSendMessage(data, true, method, args)
	}
	cascade1.SendMessage = f
	return service, nil
}

// GetAddress returns host public address.
func (network *ServiceNetwork) GetAddress() string {
	return network.hostNetwork.GetOriginHost().Address.String()
}

// GetNodeID returns current node id.
func (network *ServiceNetwork) GetNodeID() core.RecordRef {
	return network.nodeNetwork.GetID()
}

// SendMessage sends a message from MessageBus.
func (network *ServiceNetwork) SendMessage(nodeID core.RecordRef, method string, msg core.Message) ([]byte, error) {
	if msg == nil {
		return nil, errors.New("message is nil")
	}
	hostID := nodenetwork.ResolveHostID(nodeID)
	buff, err := messageToBytes(msg)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to serialize event")
	}

	log.Debugf("SendMessage with nodeID = %s method = %s, message reference = %s", nodeID.String(),
		method, msg.Target().String())

	metrics.NetworkMessageSentTotal.Inc()
	res, err := network.hostNetwork.RemoteProcedureCall(createContext(network.hostNetwork), hostID, method, [][]byte{buff})
	return res, err
}

// SendCascadeMessage sends a message from MessageBus to a cascade of nodes. Message reference is ignored
func (network *ServiceNetwork) SendCascadeMessage(data core.Cascade, method string, msg core.Message) error {
	if msg == nil {
		return errors.New("message is nil")
	}
	buff, err := messageToBytes(msg)
	if err != nil {
		return errors.Wrap(err, "Failed to serialize event")
	}

	return network.initCascadeSendMessage(data, false, method, [][]byte{buff})
}

func messageToBytes(msg core.Message) ([]byte, error) {
	reqBuff, err := message.Serialize(msg)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to serialize event")
	}
	return ioutil.ReadAll(reqBuff)
}

// RemoteProcedureRegister registers procedure for remote call on this host.
func (network *ServiceNetwork) RemoteProcedureRegister(name string, method core.RemoteProcedure) {
	network.hostNetwork.RemoteProcedureRegister(name, method)
}

// GetHostNetwork returns pointer to host network layer(DHT), temp method, refactoring needed
func (network *ServiceNetwork) GetHostNetwork() (hosthandler.HostHandler, hosthandler.Context) {
	return network.hostNetwork, createContext(network.hostNetwork)
}

func getPulseManager(components core.Components) (core.PulseManager, error) {
	if components.Ledger == nil {
		return nil, errors.New("no core.Ledger in components")
	}
	return components.Ledger.GetPulseManager(), nil
}

// Start implements core.Component
func (network *ServiceNetwork) Start(components core.Components) error {
	go network.listen()

	log.Infoln("Bootstrapping network...")
	network.bootstrap()

	log.Infoln("Searching for an active nodes...")
	err := network.hostNetwork.GetActiveNodes()
	if err != nil {
		return errors.Wrap(err, "failed to get active nodes")
	}

	err = network.hostNetwork.ObtainIP()
	if err != nil {
		return errors.Wrap(err, "Failed to ObtainIP")
	}

	err = network.hostNetwork.AnalyzeNetwork(createContext(network.hostNetwork))
	if err != nil {
		return errors.Wrap(err, "Failed to AnalyzeNetwork")
	}

	pm, err := getPulseManager(components)
	if err != nil {
		log.Error(err)
	} else {
		network.hostNetwork.GetNetworkCommonFacade().SetPulseManager(pm)
	}

	// TODO: may be merge bug, check copy-paste
	ctx := createContext(network.hostNetwork)
	err = network.hostNetwork.ObtainIP()
	if err != nil {
		return errors.Wrap(err, "Failed to ObtainIP")
	}

	err = network.hostNetwork.AnalyzeNetwork(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to AnalyzeNetwork")
	}
	// todo: end

	return nil
}

// Stop implements core.Component
func (network *ServiceNetwork) Stop() error {
	log.Infoln("Stop network")
	network.hostNetwork.Disconnect()
	return nil
}

func (network *ServiceNetwork) bootstrap() {
	err := network.hostNetwork.Bootstrap()
	if err != nil {
		log.Errorln("Failed to bootstrap network", err.Error())
		return
	}
	network.hostNetwork.GetHostsFromBootstrap()
}

func (network *ServiceNetwork) listen() {
	log.Infoln("Network starts listening")
	err := network.hostNetwork.Listen()
	if err != nil {
		log.Errorln("Listen failed:", err.Error())
	}
}

func createContext(handler hosthandler.HostHandler) hosthandler.Context {
	ctx, err := hostnetwork.NewContextBuilder(handler).SetDefaultHost().Build()
	if err != nil {
		log.Fatalln("Failed to create context:", err.Error())
	}
	return ctx
}

// initCascadeSendMessage initiates the RPC call on target host and sends messages to next cascade layers
func (network *ServiceNetwork) initCascadeSendMessage(data core.Cascade, findCurrentNode bool, method string, args [][]byte) error {
	if len(data.NodeIds) == 0 {
		return errors.New("node IDs list should not be empty")
	}
	if data.ReplicationFactor == 0 {
		return errors.New("replication factor should not be zero")
	}

	var nextNodes []core.RecordRef
	var err error

	if findCurrentNode {
		nodeID := network.nodeNetwork.GetID()
		nextNodes, err = cascade.CalculateNextNodes(data, &nodeID)
	} else {
		nextNodes, err = cascade.CalculateNextNodes(data, nil)
	}
	if err != nil {
		return errors.Wrap(err, "Failed to CalculateNextNodes")
	}
	if len(nextNodes) == 0 {
		return nil
	}

	var failedNodes []string
	for _, nextNode := range nextNodes {
		hostID := nodenetwork.ResolveHostID(nextNode)
		err = network.hostNetwork.CascadeSendMessage(data, hostID, method, args)
		if err != nil {
			log.Debugln("failed to send cascade message: ", err)
			failedNodes = append(failedNodes, nextNode.String())
		}
	}

	if len(failedNodes) > 0 {
		return errors.New("failed to send cascade message to nodes: " + strings.Join(failedNodes, ", "))
	}

	return nil
}
