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
	"crypto/ecdsa"
	"strings"
	"time"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/cascade"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/dhtnetwork"
	"github.com/insolar/insolar/network/dhtnetwork/hosthandler"
	"github.com/insolar/insolar/network/dhtnetwork/resolver"
	"github.com/pkg/errors"
)

// ServiceNetwork is facade for network.
type ServiceNetwork struct {
	hostNetwork hosthandler.HostHandler
	nodeKeeper  consensus.NodeKeeper
	certificate core.Certificate
}

// NewServiceNetwork returns a new ServiceNetwork.
func NewServiceNetwork(conf configuration.Configuration) (*ServiceNetwork, error) {

	// workaround before DI
	cert, err := certificate.NewCertificate(conf.KeysPath)
	if err != nil {
		log.Warnf("failed to read certificate: %s", err.Error())
	}

	cascade1 := &cascade.Cascade{}
	dht, err := dhtnetwork.NewHostNetwork(conf, cascade1, cert)
	if err != nil {
		return nil, err
	}

	service := &ServiceNetwork{hostNetwork: dht}
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
	return network.hostNetwork.GetNodeID()
}

// SendMessage sends a message from MessageBus.
func (network *ServiceNetwork) SendMessage(nodeID core.RecordRef, method string, msg core.Message) ([]byte, error) {
	start := time.Now()
	if msg == nil {
		return nil, errors.New("message is nil")
	}
	hostID := resolver.ResolveHostID(nodeID)
	buff, err := message.ToBytes(msg)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to serialize event")
	}

	log.Debugf("SendMessage with nodeID = %s method = %s, message reference = %s", nodeID.String(),
		method, msg.Target().String())

	metrics.NetworkMessageSentTotal.Inc()
	res, err := network.hostNetwork.RemoteProcedureCall(createContext(network.hostNetwork), hostID, method, [][]byte{buff})
	log.Debugf("Inside SendMessage: type - '%s', target - %s, caller - %s, targetRole - %s, time - %s",
		msg.Type(), msg.Target(), msg.GetCaller(), msg.TargetRole(), time.Since(start))
	return res, err
}

// SendCascadeMessage sends a message from MessageBus to a cascade of nodes. Message reference is ignored
func (network *ServiceNetwork) SendCascadeMessage(data core.Cascade, method string, msg core.Message) error {
	if msg == nil {
		return errors.New("message is nil")
	}
	buff, err := message.ToBytes(msg)
	if err != nil {
		return errors.Wrap(err, "Failed to serialize event")
	}

	return network.initCascadeSendMessage(data, false, method, [][]byte{buff})
}

// RemoteProcedureRegister registers procedure for remote call on this host.
func (network *ServiceNetwork) RemoteProcedureRegister(name string, method core.RemoteProcedure) {
	network.hostNetwork.RemoteProcedureRegister(name, method)
}

// GetHostNetwork returns pointer to host network layer(DHT), temp method, refactoring needed
func (network *ServiceNetwork) GetHostNetwork() (hosthandler.HostHandler, hosthandler.Context) {
	return network.hostNetwork, createContext(network.hostNetwork)
}

// GetPrivateKey returns a private key.
func (network *ServiceNetwork) GetPrivateKey() *ecdsa.PrivateKey {
	return network.hostNetwork.GetPrivateKey()
}

func getPulseManager(components core.Components) (core.PulseManager, error) {
	if components.Ledger == nil {
		return nil, errors.New("no core.Ledger in components")
	}
	return components.Ledger.GetPulseManager(), nil
}

// Start implements core.Component
func (network *ServiceNetwork) Start(components core.Components) error {
	network.certificate = components.Certificate
	go network.listen()

	if components.ActiveNodeComponent == nil {
		log.Error("active node component is nil")
	} else {
		nodeKeeper := components.ActiveNodeComponent.(consensus.NodeKeeper)
		network.nodeKeeper = nodeKeeper
		network.hostNetwork.SetNodeKeeper(nodeKeeper)
	}

	if components.NetworkCoordinator == nil {
		log.Error("network coordinator is nil")
	} else {
		network.hostNetwork.GetNetworkCommonFacade().SetNetworkCoordinator(components.NetworkCoordinator)
	}

	log.Infoln("Bootstrapping network...")
	network.bootstrap()

	pm, err := getPulseManager(components)
	if err != nil {
		log.Error(err)
	} else {
		network.hostNetwork.GetNetworkCommonFacade().SetPulseManager(pm)
	}

	ctx := createContext(network.hostNetwork)
	err = network.hostNetwork.ObtainIP()
	if err != nil {
		return errors.Wrap(err, "Failed to ObtainIP")
	}

	err = network.hostNetwork.AnalyzeNetwork(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to AnalyzeNetwork")
	}

	err = network.hostNetwork.StartAuthorize()
	if err != nil {
		return errors.Wrap(err, "error authorizing node")
		// log.Errorln(err.Error())
	}

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
	ctx, err := dhtnetwork.NewContextBuilder(handler).SetDefaultHost().Build()
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
		nodeID := network.hostNetwork.GetNodeID()
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
		hostID := resolver.ResolveHostID(nextNode)
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
