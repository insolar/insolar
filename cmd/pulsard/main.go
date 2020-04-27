// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	component "github.com/insolar/component-manager"
	"github.com/insolar/insconfig"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/pulsenetwork"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/pulse"
	"github.com/insolar/insolar/version"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
)

type inputParams struct {
	configPath string
}

func (i inputParams) GetConfigPath() string {
	return i.configPath
}

func parseInputParams() inputParams {
	var rootCmd = &cobra.Command{Use: "pulsard"}
	var result inputParams
	rootCmd.Flags().StringVarP(&result.configPath, "config", "c", "", "path to config file")
	rootCmd.AddCommand(version.GetCommand("pulsard"))
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("Wrong input params:", err.Error())
	}

	return result
}

// Need to fix problem with start pulsar
func main() {
	_cfg := configuration.PulsarConfiguration{}
	cfg := &_cfg
	params := parseInputParams()

	cfgParams := insconfig.Params{
		EnvPrefix:        "pulsard",
		ConfigPathGetter: params,
	}
	insConfigurator := insconfig.New(cfgParams)
	if err := insConfigurator.Load(cfg); err != nil {
		panic(err)
	}

	jww.SetStdoutThreshold(jww.LevelDebug)
	var err error

	ctx := context.Background()
	ctx, inslog := inslogger.InitNodeLogger(ctx, cfg.Log, "", "pulsar")

	jaegerflush := func() {}
	if cfg.Tracer.Jaeger.AgentEndpoint != "" {
		jconf := cfg.Tracer.Jaeger
		log.Infof("Tracing enabled. Agent endpoint: '%s', collector endpoint: '%s'", jconf.AgentEndpoint, jconf.CollectorEndpoint)
		jaegerflush = instracer.ShouldRegisterJaeger(
			ctx,
			"pulsar",
			"pulsar",
			jconf.AgentEndpoint,
			jconf.CollectorEndpoint,
			jconf.ProbabilityRate)
	}
	defer jaegerflush()

	m := metrics.NewMetrics(cfg.Metrics, metrics.GetInsolarRegistry("pulsar"), "pulsar")
	err = m.Init(ctx)
	if err != nil {
		log.Fatal("Couldn't init metrics:", err)
		os.Exit(1)
	}
	err = m.Start(ctx)
	if err != nil {
		log.Fatal("Couldn't start metrics:", err)
		os.Exit(1)
	}

	cm, server := initPulsar(ctx, cfg)
	pulseTicker := runPulsar(ctx, server, cfg.Pulsar)

	defer func() {
		pulseTicker.Stop()
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

func initPulsar(ctx context.Context, cfg *configuration.PulsarConfiguration) (*component.Manager, *pulsar.Pulsar) {
	fmt.Println("Version: ", version.GetFullVersion())
	fmt.Println("Starts with configuration:\n", configuration.ToString(cfg))

	keyStore, err := keystore.NewKeyStore(cfg.KeysPath)
	if err != nil {
		panic(err)
	}
	cryptographyScheme := platformpolicy.NewPlatformCryptographyScheme()
	cryptographyService := cryptography.NewCryptographyService()
	keyProcessor := platformpolicy.NewKeyProcessor()

	pulseDistributor, err := pulsenetwork.NewDistributor(cfg.Pulsar.PulseDistributor)
	if err != nil {
		panic(err)
	}

	cm := component.NewManager(nil)
	cm.Register(cryptographyScheme, keyStore, keyProcessor, transport.NewFactory(cfg.Pulsar.DistributionTransport))
	cm.Inject(cryptographyService, pulseDistributor)

	if err = cm.Init(ctx); err != nil {
		panic(err)
	}

	if err = cm.Start(ctx); err != nil {
		panic(err)
	}

	server := pulsar.NewPulsar(
		cfg.Pulsar,
		cryptographyService,
		cryptographyScheme,
		keyProcessor,
		pulseDistributor,
		&entropygenerator.StandardEntropyGenerator{},
	)

	return cm, server
}

func runPulsar(ctx context.Context, server *pulsar.Pulsar, cfg configuration.Pulsar) *time.Ticker {
	nextPulseNumber := pulse.OfNow()
	err := server.Send(ctx, nextPulseNumber)
	if err != nil {
		panic(err)
	}

	pulseTicker := time.NewTicker(time.Duration(cfg.PulseTime) * time.Millisecond)
	go func() {
		for range pulseTicker.C {
			err := server.Send(ctx, server.LastPN()+insolar.PulseNumber(cfg.NumberDelta))
			if err != nil {
				panic(err)
			}
		}
	}()

	return pulseTicker
}
