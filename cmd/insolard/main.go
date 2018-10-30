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
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/version"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

// componentManager is deprecated and will be removed after completly switching to component.Manager
type componentManager struct {
	components core.Components
}

// linkAll - link dependency for all components
func (cm *componentManager) linkAll(ctx context.Context) {
	v := reflect.ValueOf(cm.components)
	for i := 0; i < v.NumField(); i++ {

		if component, ok := v.Field(i).Interface().(core.Component); ok {
			componentName := v.Field(i).String()
			log.Infof("==== Old ComponentManager: Starting component `%s` ...", componentName)
			err := component.Start(ctx, cm.components)
			if err != nil {
				log.Fatalf("==== Old ComponentManager: failed to start component %s : %s", componentName, err.Error())
			}
			log.Infof("==== Old ComponentManager: Component `%s` successfully started", componentName)
		}
	}
}

// stopAll - reverse order stop all components
func (cm *componentManager) stopAll(ctx context.Context) {
	v := reflect.ValueOf(cm.components)
	for i := v.NumField() - 1; i >= 0; i-- {

		if component, ok := v.Field(i).Interface().(core.Component); ok {
			componentName := v.Field(i).String()
			err := component.Stop(ctx)
			log.Infoln("==== Old ComponentManager: Stop component: ", componentName)
			if err != nil {
				log.Errorf("==== Old ComponentManager: failed to stop component %s : %s", v.String(), err.Error())
			}
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

func registerCurrentNode(cfgHolder *configuration.Holder, cert core.Certificate, nc core.NetworkCoordinator) {
	roles := []string{"virtual", "heavy_material", "light_material"}
	host := cfgHolder.Configuration.Host.Transport.Address
	publicKey, err := cert.GetPublicKey()
	checkError("failed to get public key: ", err)

	ctx := context.TODO()
	rawCertificate, err := nc.RegisterNode(ctx, publicKey, 0, 0, roles, host)
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

	if !isBootstrap {
		mergeConfigAndCertificate(&cfgHolder.Configuration)
	}

	initLogger(cfgHolder.Configuration.Log)

	fmt.Print("Starts with configuration:\n", configuration.ToString(cfgHolder.Configuration))

	ctx := inslogger.ContextWithTrace(context.Background(), api.RandTraceID())

	cm, cmOld, repl, err := InitComponents(cfgHolder.Configuration)
	checkError("failed to init components", err)

	cmOld.linkAll(ctx)

	//err = cm.Start(ctx)
	checkError("Failed to start components", err)

	defer func() {
		log.Warn("DEFER STOP APP")
		err = cm.Stop(ctx)
		checkError("Failed to stop components", err)
	}()

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		log.Debugln("caught sig: ", sig)

		log.Warn("GRACEFULL STOP APP")
		err = cm.Stop(ctx)
		checkError("Failed to graceful stop components", err)
		os.Exit(0)
	}()

	// move to bootstrap component
	if isBootstrap {
		registerCurrentNode(cfgHolder, cmOld.components.Certificate, cmOld.components.NetworkCoordinator)
		log.Info("It's bootstrap mode, that is why gracefully stop daemon by sending SIGINT")
		gracefulStop <- syscall.SIGINT
	}

	fmt.Println("Version: ", version.GetFullVersion())
	fmt.Println("Running interactive mode:")
	repl.Start()
}

func initLogger(cfg configuration.Log) {
	err := log.SetLevel(cfg.Level)
	if err != nil {
		log.Errorln(err.Error())
	}
}
