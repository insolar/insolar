// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package light

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/server/internal"
	"github.com/insolar/insolar/version"
)

type Server struct {
	cfgHolder  *configuration.LightHolder
	apiOptions api.Options
}

func New(cfgHolder *configuration.LightHolder, apiOptions api.Options) *Server {
	return &Server{
		cfgHolder:  cfgHolder,
		apiOptions: apiOptions,
	}
}

func (s *Server) Serve() {
	cfg := *s.cfgHolder.Configuration

	fmt.Println("Version: ", version.GetFullVersion())
	fmt.Println("Starts with configuration:\n", configuration.ToString(s.cfgHolder.Configuration))

	var (
		ctx         = context.Background()
		mainTraceID = utils.RandTraceID() + "_main"
		logger      insolar.Logger
	)
	{
		var (
			nodeRole      = "light_material"
			nodeReference = ""
		)
		certManager, err := initTemporaryCertificateManager(ctx, &cfg)
		if err != nil {
			log.Warn("Failed to initialize nodeRef, nodeRole fields: ", err.Error())
		} else {
			nodeRole = certManager.GetCertificate().GetRole().String()
			nodeReference = certManager.GetCertificate().GetNodeRef().String()
		}

		ctx, logger = inslogger.InitNodeLogger(ctx, cfg.Log, nodeReference, nodeRole)
		log.InitTicker()
	}

	cmp, err := newComponents(ctx, cfg, s.apiOptions)
	fatal(ctx, err, "failed to create components")

	if cfg.Tracer.Jaeger.AgentEndpoint != "" {
		jaegerFlush := internal.Jaeger(ctx, cfg.Tracer.Jaeger, mainTraceID, cmp.NodeRef, cmp.NodeRole)
		defer jaegerFlush()
	}

	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	var waitChannel = make(chan bool)

	go func() {
		sig := <-gracefulStop
		logger.Debug("caught sig: ", sig)

		logger.Warn("GRACEFUL STOP APP")
		err = cmp.Stop(ctx)
		fatal(ctx, err, "failed to graceful stop components")
		close(waitChannel)
	}()

	err = cmp.Start(ctx)
	fatal(ctx, err, "failed to start components")
	fmt.Println("All components were started")
	<-waitChannel
}

func fatal(ctx context.Context, err error, message string) {
	if err == nil {
		return
	}
	inslogger.FromContext(ctx).Fatalf("%v: %v", message, err.Error())
}
