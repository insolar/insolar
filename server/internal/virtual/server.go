// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package virtual

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/builtin"
	"github.com/insolar/insolar/server/internal"
	"github.com/insolar/insolar/version"
)

type Server struct {
	cfgHolder        *configuration.VirtualHolder
	builtinContracts builtin.BuiltinContracts
	apiOptions       api.Options
}

func New(cfgHolder *configuration.VirtualHolder, builtinContracts builtin.BuiltinContracts, apiOptions api.Options) *Server {
	return &Server{
		cfgHolder:        cfgHolder,
		builtinContracts: builtinContracts,
		apiOptions:       apiOptions,
	}
}

func (s *Server) Serve() {
	var err error
	cfg := *s.cfgHolder.Configuration

	fmt.Println("Version: ", version.GetFullVersion())
	fmt.Println("Starts with configuration:\n", configuration.ToString(s.cfgHolder.Configuration))

	ctx := context.Background()
	bootstrapComponents := initBootstrapComponents(ctx, cfg)
	certManager := initCertificateManager(
		ctx,
		cfg,
		bootstrapComponents.CryptographyService,
		bootstrapComponents.KeyProcessor,
	)

	nodeRole := certManager.GetCertificate().GetRole().String()
	nodeRef := certManager.GetCertificate().GetNodeRef().String()

	traceID := utils.RandTraceID() + "_main"
	ctx, logger := inslogger.InitNodeLogger(ctx, cfg.Log, nodeRef, nodeRole)
	log.InitTicker()

	if cfg.Tracer.Jaeger.AgentEndpoint != "" {
		jaegerFlush := internal.Jaeger(ctx, cfg.Tracer.Jaeger, traceID, nodeRef, nodeRole)
		defer jaegerFlush()
	}

	cm, stopWatermill := initComponents(
		ctx,
		cfg,
		bootstrapComponents.CryptographyService,
		bootstrapComponents.PlatformCryptographyScheme,
		bootstrapComponents.KeyStore,
		bootstrapComponents.KeyProcessor,
		certManager,
		s.builtinContracts,
		s.apiOptions,
	)

	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	var waitChannel = make(chan bool)

	go func() {
		sig := <-gracefulStop
		logger.Debug("caught sig: ", sig)

		logger.Warn("GRACEFUL STOP APP")
		// th.Leave(ctx, 10) TODO: is actual ??
		logger.Info("main leave ends ")

		err = cm.GracefulStop(ctx)
		checkError(ctx, err, "failed to graceful stop components")

		stopWatermill()

		err = cm.Stop(ctx)
		checkError(ctx, err, "failed to stop components")
		close(waitChannel)
	}()

	err = cm.Start(ctx)
	checkError(ctx, err, "failed to start components")
	fmt.Println("All components were started")
	<-waitChannel
}

func checkError(ctx context.Context, err error, message string) {
	if err == nil {
		return
	}
	inslogger.FromContext(ctx).Fatalf("%v: %v", message, err.Error())
}
