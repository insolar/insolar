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
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/instrumentation/pprof"
)

const insolarNamespace = "insolar"
const insgorundNamespace = "insgorund"

// Metrics is a component which serve metrics data to Prometheus.
type Metrics struct {
	server   *http.Server
	listener net.Listener
}

// NewMetrics creates new Metrics component.
func NewMetrics(ctx context.Context, cfg configuration.Metrics, registry *prometheus.Registry) (*Metrics, error) {
	errlogger := &errorLogger{inslogger.FromContext(ctx)}
	promhandler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{ErrorLog: errlogger})

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhandler)
	pprof.Handle(mux)
	if cfg.ZpagesEnabled {
		// https://opencensus.io/zpages/
		zpages.Handle(mux, "/debug")
	}

	m := &Metrics{
		server: &http.Server{
			Addr:    cfg.ListenAddress,
			Handler: mux,
		},
	}

	_, err := insmetrics.RegisterPrometheus(ctx, cfg.Namespace, registry, cfg.ReportingPeriod)
	if err != nil {
		errlogger.Println(err.Error())
	}

	return m, nil
}

// ErrBind special case for Start method.
// We can use it for easier check in metrics creation code.
var ErrBind = errors.New("failed to bind")

// Start is implementation of core.Component interface.
func (m *Metrics) Start(ctx context.Context) error {
	inslog := inslogger.FromContext(ctx)

	listener, err := net.Listen("tcp", m.server.Addr)
	if err != nil {
		if opErr, ok := err.(*net.OpError); ok {
			if opErr.Op == "listen" && IsAddrInUse(opErr) {
				return ErrBind
			}
		}
		return errors.Wrap(err, "Failed to listen at address")
	}
	m.listener = listener
	inslog.Infoln("Started metrics server", m.AddrString())

	go func() {
		inslog.Debugln("metrics server starting on", m.server.Addr)
		err := m.server.Serve(listener)
		if err == nil {
			return
		}
		if IsServerClosed(err) {
			return
		}
		inslog.Errorln("falied to start metrics server", err)
	}()

	return nil
}

// Stop is implementation of core.Component interface.
func (m *Metrics) Stop(ctx context.Context) error {
	const timeOut = 3
	inslogger.FromContext(ctx).Info("Shutting down metrics server")
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeOut)*time.Second)
	defer cancel()
	err := m.server.Shutdown(ctxWithTimeout)
	if err != nil {
		return errors.Wrap(err, "Can't gracefully stop metrics server")
	}

	return nil
}

// AddrString returns listener address.
func (m *Metrics) AddrString() string {
	return m.listener.Addr().String()
}

// errorLogger wrapper for error logs.
type errorLogger struct {
	core.Logger
}

// Println is wrapper method for ErrorLn.
func (e *errorLogger) Println(v ...interface{}) {
	e.Error(v)
}

// IsAddrInUse checks error text for well known phrase.
func IsAddrInUse(err error) bool {
	return strings.Contains(err.Error(), "address already in use")
}

// IsServerClosed checks error text for well known phrase.
func IsServerClosed(err error) bool {
	return strings.Contains(err.Error(), "http: Server closed")
}
