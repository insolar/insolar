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
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/bootstrap"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/inscontext"
	"github.com/insolar/insolar/ledger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/messagebus"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/nodekeeper"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/networkcoordinator"
	"github.com/insolar/insolar/pulsar"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/version"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

type componentManager struct {
	components core.Components
}

// linkAll - link dependency for all components
func (cm *componentManager) linkAll(ctx core.Context) {
	v := reflect.ValueOf(cm.components)
	for i := 0; i < v.NumField(); i++ {
		componentName := v.Field(i).String()
		log.Infof("Starting component `%s` ...", componentName)
		err := v.Field(i).Interface().(core.Component).Start(ctx, cm.components)
		if err != nil {
			log.Fatalf("failed to start component %s : %s", componentName, err.Error())
		}

		log.Infof("Component `%s` successfully started", componentName)
	}
}

// stopAll - reverse order stop all components
func (cm *componentManager) stopAll(ctx core.Context) {
	v := reflect.ValueOf(cm.components)
	for i := v.NumField() - 1; i >= 0; i-- {
		err := v.Field(i).Interface().(core.Component).Stop(ctx)
		log.Infoln("Stop component: ", v.String())
		if err != nil {
			log.Errorf("failed to stop component %s : %s", v.String(), err.Error())
		}
	}
}

var (
	configPath               string
	isBootstrap              bool
	bootstrapCertificatePath string
)

func parseInputParams() {
	var rootCmd = &cobra.Command{Use: "insolard"}
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "path to config file")
	rootCmd.Flags().BoolVarP(&isBootstrap, "bootstrap", "b", false, "is bootstrap mode")
	rootCmd.Flags().StringVarP(&bootstrapCertificatePath, "cert_out", "r", "", "path to write bootstrap certificate")
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("Wrong input params:", err)
	}

	if isBootstrap && len(bootstrapCertificatePath) == 0 {
		log.Fatal("flag '--cert_out|-r' must not be empty, if '--bootstrap|-b' exists")
	}
}

func registerCurrentNode(cfgHolder *configuration.Holder, cert *certificate.Certificate, nc core.NetworkCoordinator) {
	roles := []string{"virtual", "heavy_material", "light_material"}
	host := cfgHolder.Configuration.Host.Transport.Address
	publicKey, err := cert.GetPublicKey()
	checkError("failed to get public key: ", err)

	rawCertificate, err := nc.RegisterNode(publicKey, 0, 0, roles, host)
	checkError("Can't register node: ", err)

	err = ioutil.WriteFile(bootstrapCertificatePath, rawCertificate, 0644)
	checkError("Can't write certificate: ", err)
}

func checkError(msg string, err error) {
	if err != nil {
		log.Fatalln(msg, err)
		os.Exit(1)
	}
}

func mergeConfigAndCertificate(cfg *configuration.Configuration) {
	if len(cfg.CertificatePath) == 0 {
		log.Info("[ mergeConfigAndCertificate ] No certificate path - No merge")
		return
	}
	cert, err := certificate.NewCertificate(cfg.KeysPath, cfg.CertificatePath)
	checkError("[ mergeConfigAndCertificate ] Can't create certificate", err)

	cfg.Host.BootstrapHosts = []string{}
	for _, bn := range cert.BootstrapNodes {
		cfg.Host.BootstrapHosts = append(cfg.Host.BootstrapHosts, bn.Host)
	}
	cfg.Node.Node.ID = cert.Reference
	cfg.Host.MajorityRule = cert.MajorityRule

	log.Infof("[ mergeConfigAndCertificate ] Add %d bootstrap nodes. Set node id to %s. Set majority rule to %d",
		len(cfg.Host.BootstrapHosts), cfg.Node.Node.ID, cfg.Host.MajorityRule)
}

func main() {
	parseInputParams()

	jww.SetStdoutThreshold(jww.LevelDebug)
	cfgHolder := configuration.NewHolder()
	var err error
	if len(configPath) != 0 {
		err = cfgHolder.LoadFromFile(configPath)
	} else {
		err = cfgHolder.Load()
	}
	if err != nil {
		log.Warnln("failed to load configuration from file: ", err.Error())
	}

	err = cfgHolder.LoadEnv()
	if err != nil {
		log.Warnln("failed to load configuration from env:", err.Error())
	}

	config := &cfgHolder.Configuration

	if !isBootstrap {
		mergeConfigAndCertificate(config)
	}

	initLogger(config.Log)

	fmt.Print("Starts with configuration:\n", configuration.ToString(config))

	ctx := inscontext.WithTraceID(inscontext.Background(), api.RandTraceID())

	cm := componentManager{}
	cert, err := certificate.NewCertificate(config.KeysPath, config.CertificatePath)
	checkError("failed to start Certificate: ", err)
	cm.components.Certificate = cert

	cm.components.NodeNetwork, err = nodekeeper.NewNodeNetwork(*config)
	checkError("failed to start NodeNetwork: ", err)

	cm.components.LogicRunner, err = logicrunner.NewLogicRunner(&config.LogicRunner)
	checkError("failed to start LogicRunner: ", err)

	cm.components.Ledger, err = ledger.NewLedger(config.Ledger)
	checkError("failed to start Ledger: ", err)

	nw, err := servicenetwork.NewServiceNetwork(*config)
	checkError("failed to start Network: ", err)
	cm.components.Network = nw

	cm.components.MessageBus, err = messagebus.NewMessageBus(*config)
	checkError("failed to start MessageBus: ", err)

	cm.components.Bootstrapper, err = bootstrap.NewBootstrapper(config.Bootstrap)
	checkError("failed to start Bootstrapper: ", err)

	cm.components.APIRunner, err = api.NewRunner(&config.APIRunner)
	checkError("failed to start ApiRunner: ", err)

	cm.components.Metrics, err = metrics.NewMetrics(config.Metrics)
	checkError("failed to start Metrics: ", err)

	cm.components.NetworkCoordinator, err = networkcoordinator.New()
	checkError("failed to start NetworkCoordinator: ", err)

	err = cm.components.LogicRunner.OnPulse(*pulsar.NewPulse(config.Pulsar.NumberDelta, 0, &entropygenerator.StandardEntropyGenerator{}))
	checkError("failed init pulse for LogicRunner: ", err)

	cm.linkAll(ctx)

	defer func() {
		cm.stopAll(ctx)
	}()

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		log.Debugln("caught sig: ", sig)

		cm.stopAll(ctx)
		os.Exit(0)
	}()

	if isBootstrap {
		registerCurrentNode(cfgHolder, cert, cm.components.NetworkCoordinator)
		log.Info("It's bootstrap mode, that is why gracefully stop daemon by sending SIGINT")
		gracefulStop <- syscall.SIGINT
	}

	fmt.Println("Version: ", version.GetFullVersion())
	fmt.Println("Running interactive mode:")
	repl(nw, cm.components.Ledger.GetPulseManager())
}

func initLogger(cfg configuration.Log) {
	err := log.SetLevel(cfg.Level)
	if err != nil {
		log.Errorln(err.Error())
	}
}
