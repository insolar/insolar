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
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/insolar/insolar/core/utils"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/version"
)

type inputParams struct {
	configPath        string
	isGenesis         bool
	genesisConfigPath string
	genesisKeyOut     string
	traceEnabled      bool
	measurementsFile  string
}

func parseInputParams() inputParams {
	var rootCmd = &cobra.Command{Use: "insolard"}
	var result inputParams
	rootCmd.Flags().StringVarP(&result.configPath, "config", "c", "", "path to config file")
	rootCmd.Flags().StringVarP(&result.genesisConfigPath, "genesis", "g", "", "path to genesis config file")
	rootCmd.Flags().StringVarP(&result.genesisKeyOut, "keyout", "", ".", "genesis certificates path")
	rootCmd.Flags().BoolVarP(&result.traceEnabled, "trace", "t", false, "enable tracing")
	rootCmd.Flags().StringVarP(&result.measurementsFile, "measure", "m", "", "enable execution time logging to the given file")
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("Wrong input params:", err)
	}

	if result.genesisConfigPath != "" {
		result.isGenesis = true
	}

	return result
}

func removeLedgerDataDir(ctx context.Context, cfg *configuration.Configuration) {
	_, err := exec.Command(
		"rm", "-rfv",
		cfg.Ledger.Storage.DataDirectory,
	).CombinedOutput()
	checkError(ctx, err, "failed to delete ledger storage data directory")
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
	cfg.Metrics.Namespace = "insolard"

	traceID := utils.RandTraceID()
	ctx, inslog := initLogger(context.Background(), cfg.Log, traceID)
	log.SetGlobalLogger(inslog)

	if params.isGenesis {
		removeLedgerDataDir(ctx, cfg)
		cfg.Ledger.PulseManager.HeavySyncEnabled = false
	}

	bootstrapComponents := initBootstrapComponents(ctx, *cfg)
	certManager := initCertificateManager(
		ctx,
		*cfg,
		params.isGenesis,
		bootstrapComponents.CryptographyService,
		bootstrapComponents.KeyProcessor,
	)

	fmt.Print("Starts with configuration:\n", configuration.ToString(cfgHolder.Configuration))

	cleanup := func() { /* by default - do nothing */ }
	if params.measurementsFile != "" {
		cleanup, err = utils.EnableExecutionTimeMeasurement(params.measurementsFile)
		if err != nil {
			log.Warnln("failed to enable execution time measurement:", err.Error())
		}
	}
	defer cleanup()

	jaegerflush := func() {}
	if params.traceEnabled {
		jconf := cfg.Tracer.Jaeger
		jaegerflush = instracer.ShouldRegisterJaeger(ctx, "insolard", jconf.AgentEndpoint, jconf.CollectorEndpoint)
		ctx = instracer.SetBaggage(ctx, instracer.Entry{Key: "traceid", Value: traceID})
	}
	defer jaegerflush()

	cm, err := initComponents(
		ctx,
		*cfg,
		bootstrapComponents.CryptographyService,
		bootstrapComponents.PlatformCryptographyScheme,
		bootstrapComponents.KeyStore,
		bootstrapComponents.KeyProcessor,
		certManager,
		params.isGenesis,
		params.genesisConfigPath,
		params.genesisKeyOut,
	)
	checkError(ctx, err, "failed to init components")

	err = cm.Init(ctx)
	checkError(ctx, err, "failed to init components")

	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	var waitChannel = make(chan bool)

	go func() {
		sig := <-gracefulStop
		inslog.Debugln("caught sig: ", sig)

		inslog.Warn("GRACEFULL STOP APP")
		err = cm.Stop(ctx)
		jaegerflush()
		checkError(ctx, err, "failed to graceful stop components")
		os.Exit(0)
	}()

	err = cm.Start(ctx)
	checkError(ctx, err, "failed to start components")
	fmt.Println("Version: ", version.GetFullVersion())
	fmt.Println("All components were started")
	<-waitChannel
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
