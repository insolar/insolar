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
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/instrumentation/pprof"
	"github.com/insolar/insolar/log"
)

const insolarNamespace = "insolar"

// Metrics is a component which serve metrics data to Prometheus.
type Metrics struct {
	registry    *prometheus.Registry
	httpHandler http.Handler
	server      *http.Server
}

// NewMetrics creates new Metrics component.
func NewMetrics(cfg configuration.Metrics) (*Metrics, error) {
	m := Metrics{registry: prometheus.NewRegistry()}
	m.httpHandler = promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{ErrorLog: &errorLogger{}})

	m.server = &http.Server{Addr: cfg.ListenAddress}

	// default system collectors
	m.registry.MustRegister(prometheus.NewProcessCollector(os.Getpid(), cfg.Namespace))
	m.registry.MustRegister(prometheus.NewGoCollector())

	// insolar collectors
	m.registry.MustRegister(NetworkMessageSentTotal)
	m.registry.MustRegister(NetworkFutures)
	m.registry.MustRegister(NetworkPacketSentTotal)
	m.registry.MustRegister(NetworkPacketReceivedTotal)

	insmetrics.RegisterPrometheus(cfg.Namespace, m.registry)

	return &m, nil
}

// Start is implementation of core.Component interface.
func (m *Metrics) Start(ctx context.Context) error {
	inslogger.FromContext(ctx).Infoln("Starting metrics server", m.server.Addr)

	http.Handle("/metrics", m.httpHandler)
	pprof.Handle(http.DefaultServeMux)

	listener, err := net.Listen("tcp", m.server.Addr)
	if err != nil {
		return errors.Wrap(err, "Failed to listen at address")
	}

	go func() {
		err := m.server.Serve(listener)
		if err != nil && err.Error() != "http: Server closed" {
			log.Errorln(err, "falied to start metrics server")
			return
		}
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

// errorLogger wrapper for error logs.
type errorLogger struct {
}

// Println is wrapper method for ErrorLn.
func (e *errorLogger) Println(v ...interface{}) {
	log.Errorln("Metrics error logger:", v)
}
