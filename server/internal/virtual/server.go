//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package virtual

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/version"
)

type Server struct {
	cfgPath string
	trace   bool
}

func New(cfgPath string, trace bool) *Server {
	return &Server{
		cfgPath: cfgPath,
		trace:   trace,
	}
}

func (s *Server) Serve() {
	cfgHolder := configuration.NewHolder()
	var err error
	if len(s.cfgPath) != 0 {
		err = cfgHolder.LoadFromFile(s.cfgPath)
	} else {
		err = cfgHolder.Load()
	}
	if err != nil {
		log.Warn("failed to load configuration from file: ", err.Error())
	}

	cfg := &cfgHolder.Configuration
	cfg.Metrics.Namespace = "insolard"

	traceID := "main_" + utils.RandTraceID()
	ctx, inslog := initLogger(context.Background(), cfg.Log, traceID)
	log.SetGlobalLogger(inslog)
	fmt.Println("Starts with configuration:\n", configuration.ToString(cfgHolder.Configuration))

	bootstrapComponents := initBootstrapComponents(ctx, *cfg)
	certManager := initCertificateManager(
		ctx,
		*cfg,
		false,
		bootstrapComponents.CryptographyService,
		bootstrapComponents.KeyProcessor,
	)

	jaegerflush := func() {}
	if s.trace {
		jconf := cfg.Tracer.Jaeger
		log.Infof("Tracing enabled. Agent endpoint: '%s', collector endpoint: '%s'\n", jconf.AgentEndpoint, jconf.CollectorEndpoint)
		jaegerflush = instracer.ShouldRegisterJaeger(
			ctx,
			certManager.GetCertificate().GetRole().String(),
			certManager.GetCertificate().GetNodeRef().String(),
			jconf.AgentEndpoint,
			jconf.CollectorEndpoint,
			jconf.ProbabilityRate)
		ctx = instracer.SetBaggage(ctx, instracer.Entry{Key: "traceid", Value: traceID})
	}
	defer jaegerflush()

	cm, th := initComponents(
		ctx,
		*cfg,
		bootstrapComponents.CryptographyService,
		bootstrapComponents.PlatformCryptographyScheme,
		bootstrapComponents.KeyStore,
		bootstrapComponents.KeyProcessor,
		certManager,
		false,
	)

	ctx, inslog = inslogger.WithField(ctx, "nodeid", certManager.GetCertificate().GetNodeRef().String())
	ctx, inslog = inslogger.WithField(ctx, "role", certManager.GetCertificate().GetRole().String())
	ctx = inslogger.SetLogger(ctx, inslog)
	log.SetGlobalLogger(inslog)

	err = cm.Init(ctx)
	checkError(ctx, err, "failed to init components")

	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	var waitChannel = make(chan bool)

	go func() {
		sig := <-gracefulStop
		inslog.Debug("caught sig: ", sig)

		inslog.Warn("GRACEFULL STOP APP")
		th.Leave(ctx, 10)
		inslog.Info("main leave ends ")
		err = cm.GracefulStop(ctx)
		checkError(ctx, err, "failed to graceful stop components")

		err = cm.Stop(ctx)
		checkError(ctx, err, "failed to stop components")
		close(waitChannel)
	}()

	err = cm.Start(ctx)
	checkError(ctx, err, "failed to start components")
	fmt.Println("Version: ", version.GetFullVersion())
	fmt.Println("All components were started")
	<-waitChannel
}

func initLogger(ctx context.Context, cfg configuration.Log, traceid string) (context.Context, insolar.Logger) {
	inslog, err := log.NewLog(cfg)
	if err != nil {
		panic(err)
	}

	if newInslog, err := inslog.WithLevel(cfg.Level); err != nil {
		inslog.Error(err.Error())
	} else {
		inslog = newInslog
	}

	ctx = inslogger.SetLogger(ctx, inslog)
	ctx, inslog = inslogger.WithTraceField(ctx, traceid)
	return ctx, inslog
}

func checkError(ctx context.Context, err error, message string) {
	if err == nil {
		return
	}
	inslogger.FromContext(ctx).Fatalf("%v: %v", message, err.Error())
}
