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

package main

import (
	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/bootstrap"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/messagebus"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/nodekeeper"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/networkcoordinator"
	"github.com/insolar/insolar/pulsar"
	"github.com/insolar/insolar/pulsar/entropygenerator"
)

// InitComponents creates and links all insolard components
func InitComponents(cfg configuration.Configuration) (*component.Manager, *ComponentManager, *Repl, error) {

	var cert *certificate.Certificate
	var err error
	if isBootstrap {
		cert, err = certificate.NewCertificatesWithKeys(cfg.KeysPath)
		checkError("failed to start Certificate ( bootstrap mode ): ", err)
	} else {
		cert, err = certificate.NewCertificate(cfg.KeysPath, cfg.CertificatePath)
		checkError("failed to start Certificate: ", err)
	}

	nodeNetwork, err := nodekeeper.NewNodeNetwork(cfg)
	checkError("failed to start NodeNetwork: ", err)

	logicRunner, err := logicrunner.NewLogicRunner(&cfg.LogicRunner)
	checkError("failed to start LogicRunner: ", err)

	ledger, err := ledger.NewLedger(cfg.Ledger)
	checkError("failed to start Ledger: ", err)

	nw, err := servicenetwork.NewServiceNetwork(cfg)
	checkError("failed to start Network: ", err)

	messageBus, err := messagebus.NewMessageBus(cfg)
	checkError("failed to start MessageBus: ", err)

	bootstrapper, err := bootstrap.NewBootstrapper(cfg.Bootstrap)
	checkError("failed to start Bootstrapper: ", err)

	apiRunner, err := api.NewRunner(&cfg.APIRunner)
	checkError("failed to start ApiRunner: ", err)

	metricsHandler, err := metrics.NewMetrics(cfg.Metrics)
	checkError("failed to start Metrics: ", err)

	networkCoordinator, err := networkcoordinator.New()
	checkError("failed to start NetworkCoordinator: ", err)

	// move to logic runner ??
	err = logicRunner.OnPulse(*pulsar.NewPulse(cfg.Pulsar.NumberDelta, 0, &entropygenerator.StandardEntropyGenerator{}))
	checkError("failed init pulse for LogicRunner: ", err)

	cm := component.Manager{}
	cm.Register(
		cert,
		nodeNetwork,
		logicRunner,
		ledger,
		nw,
		messageBus,
		bootstrapper,
		apiRunner,
		metricsHandler,
		networkCoordinator,
	)

	cmOld := ComponentManager{components: core.Components{
		Certificate:        cert,
		NodeNetwork:        nodeNetwork,
		LogicRunner:        logicRunner,
		Ledger:             ledger,
		Network:            nw,
		MessageBus:         messageBus,
		Metrics:            metricsHandler,
		Bootstrapper:       bootstrapper,
		APIRunner:          apiRunner,
		NetworkCoordinator: networkCoordinator,
	}}

	return &cm, &cmOld, &Repl{Manager: ledger.GetPulseManager(), Service: nw}, nil
}
