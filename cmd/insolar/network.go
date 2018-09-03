/*
 *    Copyright 2018 INS Ecosystem
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

package main

import (
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/sirupsen/logrus"
)

// Network contains all Insolar network layers
type Network struct {
	HostNetwork *hostnetwork.DHT
	Node        *nodenetwork.Node
	ctx         hostnetwork.Context
}

// StartNetwork creates and starts network
func StartNetwork(cfg configuration.Configuration) (Network, error) {
	var n Network
	var err error

	logrus.Infoln("Starting network...")
	n.HostNetwork, err = hostnetwork.NewHostNetwork(cfg.Host)
	if err != nil {
		logrus.Errorln("Failed to create network:", err.Error())
	}

	n.ctx = n.createContext()
	go n.listen()

	if len(cfg.Host.BootstrapHosts) > 0 {
		logrus.Infoln("Bootstrapping network...")
		n.bootstrap()
	}

	err = n.HostNetwork.ObtainIP(n.ctx)
	if err != nil {
		logrus.Errorln(err)
	}

	err = n.HostNetwork.AnalyzeNetwork(n.ctx)
	if err != nil {
		logrus.Errorln(err)
	}

	n.Node = nodenetwork.NewNode("123", "MyDomain")
	return n, nil
}

func (n *Network) closeNetwork() {
	logrus.Infoln("Close network")
	n.HostNetwork.Disconnect()
}

func (n *Network) bootstrap() {
	err := n.HostNetwork.Bootstrap()
	if err != nil {
		logrus.Errorln("Failed to bootstrap network", err.Error())
	}
}

func (n *Network) listen() {
	func() {
		logrus.Infoln("Network starts listening")
		err := n.HostNetwork.Listen()
		if err != nil {
			logrus.Errorln("Listen failed:", err.Error())
		}
	}()
}

func (n *Network) createContext() hostnetwork.Context {
	ctx, err := hostnetwork.NewContextBuilder(n.HostNetwork).SetDefaultHost().Build()
	if err != nil {
		logrus.Fatalln("Failed to create context:", err.Error())
	}
	return ctx
}
