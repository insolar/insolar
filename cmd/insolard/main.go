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
	"os"
	"os/signal"
	"syscall"

	"github.com/insolar/insolar/core/utils"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/version"
)

/*
// ComponentManager is deprecated and will be removed after completly switching to component.Manager
type ComponentManager struct {
	components core.Components
}

func (cm *ComponentManager) GetAll() core.Components {
	return cm.components
}

// LinkAll - link dependency for all components
func (cm *ComponentManager) LinkAll(ctx context.Context) {
	inslog := inslogger.FromContext(ctx)
	v := reflect.ValueOf(cm.components)
	for i := 0; i < v.NumField(); i++ {

		if component, ok := v.Field(i).Interface().(core.Component); ok {
			componentName := v.Field(i).String()
			inslog.Infof("==== Old ComponentManager: Starting component `%s` ...", componentName)
			err := component.Start(ctx, cm.components)
			if err != nil {
				inslog.Fatalf("==== Old ComponentManager: failed to start component %s : %s", componentName, err.Error())
			}
			inslog.Infof("==== Old ComponentManager: Component `%s` successfully started", componentName)
		}
	}
}
*/
type inputParams struct {
	configPath        string
	isGenesis         bool
	genesisConfigPath string
	genesisKeyOut     string
	traceEnabled      bool
}

func parseInputParams() inputParams {
	var rootCmd = &cobra.Command{Use: "insolard"}
	var result inputParams
	rootCmd.Flags().StringVarP(&result.configPath, "config", "c", "", "path to config file")
	rootCmd.Flags().StringVarP(&result.genesisConfigPath, "genesis", "g", "", "path to genesis config file")
	rootCmd.Flags().StringVarP(&result.genesisKeyOut, "keyout", "", ".", "genesis certificates path")
	rootCmd.Flags().BoolVarP(&result.traceEnabled, "trace", "t", false, "enable tracing")
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("Wrong input params:", err)
	}

	if result.genesisConfigPath != "" {
		result.isGenesis = true
	}

	return result
}

func mergeConfigAndCertificate(ctx context.Context, cfg *configuration.Configuration, cert *certificate.Certificate) {
	inslog := inslogger.FromContext(ctx)
	if len(cfg.CertificatePath) == 0 {
		inslog.Info("No certificate path - No merge")
		return
	}

	cfg.Host.BootstrapHosts = []string{}
	for _, bn := range cert.BootstrapNodes {
		cfg.Host.BootstrapHosts = append(cfg.Host.BootstrapHosts, bn.Host)
	}
	cfg.Node.Node.ID = cert.Reference
	cfg.Host.MajorityRule = cert.MajorityRule

	inslog.Infof("Add %d bootstrap nodes. Set node id to %s. Set majority rule to %d",
		len(cfg.Host.BootstrapHosts), cfg.Node.Node.ID, cfg.Host.MajorityRule)
}

func main() {
	params := parseInputParams()

	jww.SetStdoutThreshold(jww.LevelDebug)
	cfgHolder := configuration.NewHolder()
	var err error
	if len(params.configPath) != 0 {
		err = cfgHolder.LoadFromFile(params.configPath)
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

	cfg := &cfgHolder.Configuration

	traceid := utils.RandTraceID()
	ctx, inslog := initLogger(context.Background(), cfg.Log, traceid)

	bootstrapComponents := InitBootstrapComponents(ctx, *cfg)
	cert := InitCertificate(
		ctx,
		*cfg,
		params.isGenesis,
		bootstrapComponents.CryptographyService,
		bootstrapComponents.KeyProcessor,
	)

	if !params.isGenesis {
		mergeConfigAndCertificate(ctx, cfg, cert)
	}
	cfg.Metrics.Namespace = "insolard"

	fmt.Print("Starts with configuration:\n", configuration.ToString(cfgHolder.Configuration))

	jaegerflush := func() {}
	if params.traceEnabled {
		jconf := cfg.Tracer.Jaeger
		jaegerflush = instracer.ShouldRegisterJaeger(ctx, "insolard", jconf.AgentEndpoint, jconf.CollectorEndpoint)
		ctx = instracer.SetBaggage(ctx, instracer.Entry{Key: "traceid", Value: traceid})
	}
	defer jaegerflush()

	cm, repl, err := InitComponents(
		ctx,
		*cfg,
		bootstrapComponents.CryptographyService,
		bootstrapComponents.PlatformCryptographyScheme,
		bootstrapComponents.KeyStore,
		bootstrapComponents.KeyProcessor,
		cert,
		params.isGenesis,
		params.genesisConfigPath,
		params.genesisKeyOut,
	)
	checkError(ctx, err, "failed to init components")

	err = cm.Init(ctx)
	checkError(ctx, err, "failed to init components")

	//cmOld.LinkAll(ctx)

	err = cm.Start(ctx)
	checkError(ctx, err, "failed to start components")

	defer func() {
		inslog.Warn("DEFER STOP APP")
		err = cm.Stop(ctx)
		checkError(ctx, err, "failed to stop components")
	}()

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	go func() {
		sig := <-gracefulStop
		inslog.Debugln("caught sig: ", sig)

		inslog.Warn("GRACEFULL STOP APP")
		err = cm.Stop(ctx)
		jaegerflush()
		checkError(ctx, err, "failed to graceful stop components")
		os.Exit(0)
	}()

	fmt.Println("Version: ", version.GetFullVersion())
	fmt.Println("Running interactive mode:")
	repl.Start(ctx)
}

func initLogger(ctx context.Context, cfg configuration.Log, traceid string) (context.Context, core.Logger) {
	inslog, err := log.NewLog(cfg)
	if err != nil {
		panic(err)
	}
	err = inslog.SetLevel(cfg.Level)
	if err != nil {
		inslog.Errorln(err.Error())
	}
	return inslogger.WithTraceField(inslogger.SetLogger(ctx, inslog), traceid)
}

func checkError(ctx context.Context, err error, message string) {
	if err == nil {
		return
	}
	inslog := inslogger.FromContext(ctx)
	log.WithSkipDelta(inslog, +1).Fatalf("%v: %v", message, err.Error())
}
