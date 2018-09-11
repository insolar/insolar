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
	"log"

	"io/ioutil"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ServiceNetwork is facade for network.
type ServiceNetwork struct {
	nodeNetwork *nodenetwork.Nodenetwork
	hostNetwork *hostnetwork.DHT
}

// NewServiceNetwork returns a new ServiceNetwork.
func NewServiceNetwork(
	hostConf configuration.HostNetwork,
	nodeConf configuration.NodeNetwork) (*ServiceNetwork, error) {

	dht, err := hostnetwork.NewHostNetwork(hostConf)
	if err != nil {
		return nil, err
	}
	node := nodenetwork.NewNodeNetwork(nodeConf)
	if node == nil {
		return nil, errors.New("failed to create a node network")
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
func (network *ServiceNetwork) SendMessage(method string, msg core.Message) ([]byte, error) {
	if msg == nil {
		return nil, errors.New("message is nil")
	}
	hostID, err := network.nodeNetwork.GetReferenceHostID(msg.GetReference().String())
	if err != nil {
		return nil, err
	}
	reqBuff, err := msg.Serialize()
	if err != nil {
		return nil, err
	}
	buff, err := ioutil.ReadAll(reqBuff)
	if err != nil {
		return nil, err
	}
	res, err := network.hostNetwork.RemoteProcedureCall(createContext(network.hostNetwork), hostID, method, [][]byte{buff})
	return res, err
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
	logrus.Infoln("Bootstrapping network...")
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
	logrus.Infoln("Stop network")
	network.hostNetwork.Disconnect()
	return nil
}

func (network *ServiceNetwork) bootstrap() {
	err := network.hostNetwork.Bootstrap()
	if err != nil {
		logrus.Errorln("Failed to bootstrap network", err.Error())
	}
}

func (network *ServiceNetwork) listen() {
	func() {
		logrus.Infoln("Network starts listening")
		err := network.hostNetwork.Listen()
		if err != nil {
			logrus.Errorln("Listen failed:", err.Error())
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
