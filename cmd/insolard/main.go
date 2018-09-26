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
	"github.com/insolar/insolar/eventbus"
	"github.com/insolar/insolar/ledger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/version"
	jww "github.com/spf13/jwalterweatherman"
)

type componentManager struct {
	components core.Components
}

// linkAll - link dependency for all components
func (cm *componentManager) linkAll() {
	v := reflect.ValueOf(cm.components)
	for i := 0; i < v.NumField(); i++ {
		err := v.Field(i).Interface().(core.Component).Start(cm.components)
		if err != nil {
			log.Errorf("failed to start component %s : %s", v.Field(i).String(), err.Error())
		}
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

func main() {
	jww.SetStdoutThreshold(jww.LevelDebug)
	cfgHolder := configuration.NewHolder()
	err := cfgHolder.Load()
	if err != nil {
		log.Warnln("Failed to load configuration from file: ", err.Error())
	}

	err = cfgHolder.LoadEnv()
	if err != nil {
		log.Warnln("Failed to load configuration from env:", err.Error())
	}

	initLogger(cfgHolder.Configuration.Log)

	fmt.Print("Starts with configuration:\n", configuration.ToString(cfgHolder.Configuration))

	cm := componentManager{}
	nw, err := servicenetwork.NewServiceNetwork(cfgHolder.Configuration.Host, cfgHolder.Configuration.Node)
	if err != nil {
		log.Fatalln("Failed to start Network: ", err.Error())
	}
	cm.components.Network = nw

	cm.components.Ledger, err = ledger.NewLedger(cfgHolder.Configuration.Ledger)
	if err != nil {
		log.Fatalln("Failed to start Ledger: ", err.Error())
	}

	cm.components.LogicRunner, err = logicrunner.NewLogicRunner(cfgHolder.Configuration.LogicRunner)
	if err != nil {
		log.Fatalln("Failed to start LogicRunner: ", err.Error())
	}

	cm.components.EventBus, err = eventbus.New(cfgHolder.Configuration)
	if err != nil {
		log.Fatalln("Failed to start LogicRunner: ", err.Error())
	}

	cm.components.Bootstrapper, err = bootstrap.NewBootstrapper(cfgHolder.Configuration)
	if err != nil {
		log.Fatalln("Failed to start Bootstrapper: ", err.Error())
	}

	cm.components.APIRunner, err = api.NewRunner(&cfgHolder.Configuration.APIRunner)
	if err != nil {
		log.Fatalln("Failed to start ApiRunner: ", err.Error())
	}

	cm.components.Metrics, err = metrics.NewMetrics(cfgHolder.Configuration.Metrics)
	if err != nil {
		log.Fatalln("Failed to start Metrics: ", err.Error())
	}

	cm.linkAll()

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
