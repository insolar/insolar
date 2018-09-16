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
	"os"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

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

	m.registry.MustRegister(prometheus.NewProcessCollector(os.Getpid(), ""))
	m.registry.MustRegister(prometheus.NewGoCollector())

	return m, nil
}

// Start is implementation of core.Component interface
func (m *Metrics) Start(components core.Components) error {
	log.Infoln("Starting metrics server")
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

// AddCounter adds new counter to metrics registry
func (m *Metrics) AddCounter(name, componentName, help string) (prometheus.Counter, error) {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name:      name,
		Help:      help,
		Namespace: "insolar",
		Subsystem: componentName,
	})

	log.Debugln("Register counter: " + name)
	err := m.registry.Register(counter)
	if err != nil {
		return nil, err
	}
	return counter, nil
}

// AddGauge adds new gauge to metrics registry
func (m *Metrics) AddGauge(name, componentName, help string) (prometheus.Gauge, error) {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      name,
		Help:      help,
		Namespace: "insolar",
		Subsystem: componentName,
	})

	log.Debugln("Register gauge: " + name)
	err := m.registry.Register(gauge)
	if err != nil {
		return nil, err
	}
	return gauge, nil
}
