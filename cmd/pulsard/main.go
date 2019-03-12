/*
 *    Copyright 2019 Insolar Technologies
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
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/utils"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/pulsenetwork"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/transport/relay"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/pulsar/storage"
	"github.com/insolar/insolar/version"
)

type inputParams struct {
	configPath   string
	traceEnabled bool
}

func parseInputParams() inputParams {
	var rootCmd = &cobra.Command{Use: "insolard"}
	var result inputParams
	rootCmd.Flags().StringVarP(&result.configPath, "config", "c", "", "path to config file")
	rootCmd.Flags().BoolVarP(&result.traceEnabled, "trace", "t", false, "enable tracing")
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("Wrong input params:", err.Error())
	}

	return result
}

// Need to fix problem with start pulsar
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
		log.Warn("failed to load configuration from file: ", err.Error())
	}

	traceID := utils.RandTraceID()
	ctx, inslog := initLogger(context.Background(), cfgHolder.Configuration.Log, traceID)
	log.SetGlobalLogger(inslog)

	jaegerflush := func() {}
	if params.traceEnabled {
		jconf := cfgHolder.Configuration.Tracer.Jaeger
		log.Infof("Tracing enabled. Agent endpoint: '%s', collector endpoint: '%s'", jconf.AgentEndpoint, jconf.CollectorEndpoint)
		jaegerflush = instracer.ShouldRegisterJaeger(
			ctx,
			"pulsar",
			core.RecordRef{}.String(),
			jconf.AgentEndpoint,
			jconf.CollectorEndpoint,
			jconf.ProbabilityRate)
		ctx = instracer.SetBaggage(ctx, instracer.Entry{Key: "traceid", Value: traceID})
	}
	defer jaegerflush()

	cm, server, storage := initPulsar(ctx, cfgHolder.Configuration)
	server.ID = traceID

	go server.StartServer(ctx)
	pulseTicker, refreshTicker := runPulsar(ctx, server, cfgHolder.Configuration.Pulsar)

	defer func() {
		pulseTicker.Stop()
		refreshTicker.Stop()
		err = storage.Close()
		if err != nil {
			inslog.Error(err)
		}
		server.StopServer(ctx)
		err = cm.Stop(ctx)
		if err != nil {
			inslog.Error(err)
		}
	}()

	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	<-gracefulStop
}

func initPulsar(ctx context.Context, cfg configuration.Configuration) (*component.Manager, *pulsar.Pulsar, pulsarstorage.PulsarStorage) {
	fmt.Println("Starts with configuration:\n", configuration.ToString(cfg))
	fmt.Println("Version: ", version.GetFullVersion())

	keyStore, err := keystore.NewKeyStore(cfg.KeysPath)
	if err != nil {
		inslogger.FromContext(ctx).Fatal(err)
	}
	cryptographyScheme := platformpolicy.NewPlatformCryptographyScheme()
	cryptographyService := cryptography.NewCryptographyService()
	keyProcessor := platformpolicy.NewKeyProcessor()

	tp, err := transport.NewTransport(cfg.Pulsar.DistributionTransport, relay.NewProxy())
	if err != nil {
		inslogger.FromContext(ctx).Fatal(err)
	}

	pulseDistributor, err := pulsenetwork.NewDistributor(cfg.Pulsar.PulseDistributor)
	if err != nil {
		inslogger.FromContext(ctx).Fatal(err)
	}

	cm := &component.Manager{}
	cm.Register(cryptographyScheme, keyStore, keyProcessor, tp)
	cm.Inject(cryptographyService, pulseDistributor)

	if err = cm.Init(ctx); err != nil {
		inslogger.FromContext(ctx).Fatal(err)
	}

	if err = cm.Start(ctx); err != nil {
		inslogger.FromContext(ctx).Fatal(err)
	}

	storage, err := pulsarstorage.NewStorageBadger(cfg.Pulsar, nil)
	if err != nil {
		inslogger.FromContext(ctx).Fatal(err)
		panic(err)
	}
	switcher := &pulsar.StateSwitcherImpl{}
	server, err := pulsar.NewPulsar(
		cfg.Pulsar,
		cryptographyService,
		cryptographyScheme,
		keyProcessor,
		pulseDistributor,
		storage,
		&pulsar.RPCClientWrapperFactoryImpl{},
		&entropygenerator.StandardEntropyGenerator{},
		switcher,
		net.Listen,
	)

	if err != nil {
		inslogger.FromContext(ctx).Fatal(err)
		panic(err)
	}
	switcher.SetPulsar(server)

	return cm, server, storage
}

func runPulsar(ctx context.Context, server *pulsar.Pulsar, cfg configuration.Pulsar) (pulseTicker *time.Ticker, refreshTicker *time.Ticker) {
	server.CheckConnectionsToPulsars(ctx)

	nextPulseNumber := core.CalculatePulseNumber(time.Now())

	err := server.StartConsensusProcess(ctx, nextPulseNumber)
	if err != nil {
		inslogger.FromContext(ctx).Fatal(err)
		panic(err)
	}
	pulseTicker = time.NewTicker(time.Duration(cfg.PulseTime) * time.Millisecond)
	go func() {
		for range pulseTicker.C {
			err = server.StartConsensusProcess(ctx, core.PulseNumber(server.GetLastPulse().PulseNumber+core.PulseNumber(cfg.NumberDelta)))
			if err != nil {
				inslogger.FromContext(ctx).Fatal(err)
				panic(err)
			}
		}
	}()

	refreshTicker = time.NewTicker(1 * time.Second)
	go func() {
		for range refreshTicker.C {
			server.CheckConnectionsToPulsars(ctx)
		}
	}()

	return
}

func initLogger(ctx context.Context, cfg configuration.Log, traceid string) (context.Context, core.Logger) {
	inslog, err := log.NewLog(cfg)
	if err != nil {
		panic(err)
	}
	err = inslog.SetLevel(cfg.Level)
	if err != nil {
		inslog.Error(err.Error())
	}
	return inslogger.WithTraceField(inslogger.SetLogger(ctx, inslog), traceid)
}
