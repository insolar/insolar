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
	"net/http"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsDescriptor struct {
	Name      string
	Collector interface{}
}

// Metrics is a component
type Metrics struct {
	registry    *prometheus.Registry
	httpHandler http.Handler
	server      *http.Server
}

// NewMetrics creates new Metrics component
func NewMetrics(cfg configuration.Metrics) (Metrics, error) {
	m := Metrics{registry: prometheus.NewRegistry()}
	m.httpHandler = promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{ErrorLog: &errorLogger{}})

	m.server = &http.Server{Addr: cfg.ListenAddress}
	err := m.registry.Register(nodeCount)

	return m, err
}

// Start is implementation of core.Component interface
func (m *Metrics) Start(components core.Components) error {
	http.Handle("/metrics", m.httpHandler)
	go m.server.ListenAndServe()

	return nil
}

// Stop is implementation of core.Component interface
func (m *Metrics) Stop() error {
	const timeOut = 3
	log.Infoln("Shutting down metrics server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancel()
	err := m.server.Shutdown(ctx)
	if err != nil {
		return errors.Wrap(err, "Can't gracefully stop metrics server")
	}

	return nil
}

// errorLogger wrapper for error logs
type errorLogger struct {
}

// Println is wrapper method for ErrorLn
func (e *errorLogger) Println(v ...interface{}) {
	log.Errorln(v)
}

func newCounter(name, subsystem, help string) prometheus.Counter {
	return prometheus.NewCounter(prometheus.CounterOpts{
		Name:      name,
		Help:      help,
		Namespace: "insolar",
		Subsystem: subsystem,
	})
}

func newGauge(name, subsystem, help string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      name,
		Help:      help,
		Namespace: "insolar",
		Subsystem: subsystem,
	})
}
