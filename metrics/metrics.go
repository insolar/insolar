// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package metrics

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opencensus.io/zpages"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/instrumentation/pprof"
	"github.com/insolar/insolar/log"
)

const insolarNamespace = "insolar"
const insgorundNamespace = "insgorund"

// Metrics is a component which serve metrics data to Prometheus.
type Metrics struct {
	config   configuration.Metrics
	registry *prometheus.Registry
	nodeRole string

	handler  http.Handler
	server   *http.Server
	listener net.Listener
}

// NewMetrics creates new Metrics instance.
func NewMetrics(
	cfg configuration.Metrics,
	registry *prometheus.Registry,
	nodeRole string,
) *Metrics {
	m := &Metrics{
		config:   cfg,
		registry: registry,
		nodeRole: nodeRole,
	}
	return m
}

// Init inits metrics instance. Creates and registers all handlers.
func (m *Metrics) Init(ctx context.Context) error {
	mux := http.NewServeMux()

	errLogger := &errorLogger{inslogger.FromContext(ctx)}
	promHandler := promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{ErrorLog: errLogger})
	mux.Handle("/metrics", promHandler)
	mux.Handle("/_status", newProcStatus())
	mux.Handle("/debug/loglevel", log.NewLoglevelChangeHandler())
	pprof.Handle(mux)
	if m.config.ZpagesEnabled {
		// https://opencensus.io/zpages/
		zpages.Handle(mux, "/debug")
	}

	_, err := insmetrics.RegisterPrometheus(
		m.config.Namespace,
		m.registry,
		m.config.ReportingPeriod,
		errLogger,
		m.nodeRole,
	)
	if err != nil {
		errLogger.Error(err.Error())
		return err
	}

	m.handler = mux
	return nil
}

// Handler returns http handler, created on initialization stage.
func (m *Metrics) Handler() http.Handler {
	return m.handler
}

// ErrBind special case for Start method.
// We can use it for easier check in metrics creation code.
var ErrBind = errors.New("Failed to bind")

// Start is implementation of insolar.Component interface.
func (m *Metrics) Start(ctx context.Context) error {
	inslog := inslogger.FromContext(ctx)
	m.server = &http.Server{
		Addr:    m.config.ListenAddress,
		Handler: m.handler,
	}

	listener, err := net.Listen("tcp", m.server.Addr)
	if err != nil {
		if opErr, ok := err.(*net.OpError); ok {
			if opErr.Op == "listen" && IsAddrInUse(opErr) {
				return errors.Wrapf(ErrBind, "addr=%v", m.server.Addr)
			}
		}
		return errors.Wrap(err, "Failed to listen at address")
	}
	m.listener = listener
	inslog.Info("Started metrics server: ", m.listener.Addr().String())

	go func() {
		inslog.Debug("Metrics server starting on ", m.server.Addr)
		if err := m.server.Serve(listener); err != http.ErrServerClosed {
			inslog.Error("Failed to start metrics server ", err)
		}
	}()

	return nil
}

// Stop is implementation of insolar.Component interface.
func (m *Metrics) Stop(ctx context.Context) error {
	if m.server == nil {
		return nil
	}

	const timeOut = 3
	inslogger.FromContext(ctx).Info("Shutting down metrics server")
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeOut)*time.Second)
	defer cancel()

	err := m.server.Shutdown(ctxWithTimeout)
	return errors.Wrap(err, "Can't gracefully stop metrics server")
}

// errorLogger wrapper for error logs.
type errorLogger struct {
	insolar.Logger
}

// Println is wrapper method for Error method.
func (e *errorLogger) Println(v ...interface{}) {
	e.Error(v)
}

// IsAddrInUse checks error text for well known phrase.
func IsAddrInUse(err error) bool {
	return strings.Contains(err.Error(), "address already in use")
}
