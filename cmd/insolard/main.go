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
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/bootstrap"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
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
func (cm *componentManager) linkAll() {
	v := reflect.ValueOf(cm.components)
	for i := 0; i < v.NumField(); i++ {
		componentName := v.Field(i).String()
		log.Infof("Starting component `%s` ...", componentName)
		err := v.Field(i).Interface().(core.Component).Start(cm.components)
		if err != nil {
			log.Fatalf("failed to start component %s : %s", componentName, err.Error())
		}

		log.Infof("Component `%s` successfully started", componentName)
	}
}

// stopAll - reverse order stop all components
func (cm *componentManager) stopAll() {
	v := reflect.ValueOf(cm.components)
	for i := v.NumField() - 1; i >= 0; i-- {
		err := v.Field(i).Interface().(core.Component).Stop()
		log.Infoln("Stop component: ", v.String())
		if err != nil {
			log.Errorf("failed to stop component %s : %s", v.String(), err.Error())
		}
	}
}

var (
	configPath string
)

func parseInputParams() {
	var rootCmd = &cobra.Command{Use: "insolard"}
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "path to config file")
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("Wrong input params:", err)
	}
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

	initLogger(cfgHolder.Configuration.Log)

	fmt.Print("Starts with configuration:\n", configuration.ToString(cfgHolder.Configuration))

	cm := componentManager{}

	cm.components.ActiveNodeComponent = nodekeeper.NewActiveNodeComponent(cfgHolder.Configuration)

	cm.components.LogicRunner, err = logicrunner.NewLogicRunner(&cfgHolder.Configuration.LogicRunner)
	if err != nil {
		log.Fatalln("failed to start LogicRunner: ", err.Error())
	}

	cm.components.Ledger, err = ledger.NewLedger(cfgHolder.Configuration.Ledger)
	if err != nil {
		log.Fatalln("failed to start Ledger: ", err.Error())
	}

	nw, err := servicenetwork.NewServiceNetwork(cfgHolder.Configuration)
	if err != nil {
		log.Fatalln("failed to start Network: ", err.Error())
	}
	cm.components.Network = nw

	cm.components.MessageBus, err = messagebus.NewMessageBus(cfgHolder.Configuration)
	if err != nil {
		log.Fatalln("failed to start LogicRunner: ", err.Error())
	}

	cm.components.Bootstrapper, err = bootstrap.NewBootstrapper(cfgHolder.Configuration)
	if err != nil {
		log.Fatalln("failed to start Bootstrapper: ", err.Error())
	}

	cm.components.APIRunner, err = api.NewRunner(&cfgHolder.Configuration.APIRunner)
	if err != nil {
		log.Fatalln("failed to start ApiRunner: ", err.Error())
	}

	cm.components.Metrics, err = metrics.NewMetrics(cfgHolder.Configuration.Metrics)
	if err != nil {
		log.Fatalln("failed to start Metrics: ", err.Error())
	}

	cm.components.NetworkCoordinator, err = networkcoordinator.New()
	if err != nil {
		log.Fatalln("failed to start NetworkCoordinator: ", err.Error())
	}

	cm.linkAll()
	err = cm.components.LogicRunner.OnPulse(*pulsar.NewPulse(cfgHolder.Configuration.Pulsar.NumberDelta, 0, &entropygenerator.StandardEntropyGenerator{}))
	if err != nil {
		log.Fatalln("failed init pulse for LogicRunner: ", err.Error())
	}

	defer func() {
		cm.stopAll()
	}()

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		log.Debugln("caught sig: ", sig)

		cm.stopAll()
		os.Exit(0)
	}()

	fmt.Println("Version: ", version.GetFullVersion())
	fmt.Println("Running interactive mode:")
	repl(nw)
}

func initLogger(cfg configuration.Log) {
	err := log.SetLevel(cfg.Level)
	if err != nil {
		log.Errorln(err.Error())
	}
}
