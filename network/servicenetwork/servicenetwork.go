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
	"bytes"
	"encoding/gob"
	"io/ioutil"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/pkg/errors"
)

// ServiceNetwork is facade for network.
type ServiceNetwork struct {
	nodeNetwork *nodenetwork.NodeNetwork
	hostNetwork *hostnetwork.DHT
}

// NewServiceNetwork returns a new ServiceNetwork.
func NewServiceNetwork(
	hostConf configuration.HostNetwork,
	nodeConf configuration.NodeNetwork) (*ServiceNetwork, error) {

	node := nodenetwork.NewNodeNetwork(nodeConf)
	if node == nil {
		return nil, errors.New("failed to create a node network")
	}

	dht, err := hostnetwork.NewHostNetwork(hostConf, node)
	if err != nil {
		return nil, err
	}

	err = dht.ObtainIP(createContext(dht))
	if err != nil {
		return nil, err
	}

	err = dht.AnalyzeNetwork(createContext(dht))
	if err != nil {
		return nil, err
	}

	return &ServiceNetwork{nodeNetwork: node, hostNetwork: dht}, nil
}

// GetAddress returns host public address.
func (network *ServiceNetwork) GetAddress() string {
	ctx, err := hostnetwork.NewContextBuilder(network.hostNetwork).SetDefaultHost().Build()
	if err != nil {
		log.Fatalln("Failed to create context:", err.Error())
	}
	return network.hostNetwork.GetOriginHost(ctx).Address.String()
}

// SendMessage sends a message from MessageRouter.
func (network *ServiceNetwork) SendMessage(nodeID core.RecordRef, method string, msg core.Message) ([]byte, error) {
	if msg == nil {
		return nil, errors.New("message is nil")
	}
	hostID := network.nodeNetwork.ResolveHostID(nodeID)
	buff, err := messageToBytes(msg)
	if err != nil {
		return nil, err
	}

	log.Debugln("SendMessage with nodeID = %s method = %s, message reference = %s", nodeID.String(),
		method, msg.GetReference().String())

	res, err := network.hostNetwork.RemoteProcedureCall(createContext(network.hostNetwork), hostID, method, [][]byte{buff})
	return res, err
}

// SendCascadeMessage sends a message from MessageRouter to a cascade of nodes. Message reference is ignored
func (network *ServiceNetwork) SendCascadeMessage(data core.Cascade, method string, msg core.Message) error {
	if msg == nil {
		return errors.New("message is nil")
	}
	buff, err := messageToBytes(msg)
	if err != nil {
		return err
	}

	log.Debugln("SendCascadeMessage with cascade NodeIds = %v method = %s, message reference = %s", data.NodeIds,
		method, msg.GetReference().String())

	return network.hostNetwork.InitCascadeSendMessage(data, nil, createContext(network.hostNetwork), method, [][]byte{buff})
}

func messageToBytes(msg core.Message) ([]byte, error) {
	reqBuff, err := msg.Serialize()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(reqBuff)
}

// Serialize converts Message or Response to byte slice.
func Serialize(value interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(value)
	if err != nil {
		return nil, err
	}
	res := buffer.Bytes()
	return res, err
}

// RemoteProcedureRegister registers procedure for remote call on this host.
func (network *ServiceNetwork) RemoteProcedureRegister(name string, method core.RemoteProcedure) {
	network.hostNetwork.RemoteProcedureRegister(name, method)
}

// GetHostNetwork returns pointer to host network layer(DHT), temp method, refactoring needed
func (network *ServiceNetwork) GetHostNetwork() (*hostnetwork.DHT, hostnetwork.Context) {
	return network.hostNetwork, createContext(network.hostNetwork)
}

// Start implements core.Component
func (network *ServiceNetwork) Start(components core.Components) error {
	go network.listen()
	log.Infoln("Bootstrapping network...")
	network.bootstrap()

	ctx := createContext(network.hostNetwork)
	err := network.hostNetwork.ObtainIP(ctx)
	if err != nil {
		return err
	}

	err = network.hostNetwork.AnalyzeNetwork(ctx)
	if err != nil {
		return err
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
	}
}

func (network *ServiceNetwork) listen() {
	func() {
		log.Infoln("Network starts listening")
		err := network.hostNetwork.Listen()
		if err != nil {
			log.Errorln("Listen failed:", err.Error())
		}
	}()
}

func createContext(dht *hostnetwork.DHT) hostnetwork.Context {

	ctx, err := hostnetwork.NewContextBuilder(dht).SetDefaultHost().Build()
	if err != nil {
		log.Fatalln("Failed to create context:", err.Error())
	}
	return ctx
}
