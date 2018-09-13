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

package main

import (
	"net/http"

	"github.com/insolar/insolar/configuration"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func handleStats(cfg configuration.Stats) {
	logrus.Println("handleStats")
	var cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current temperature of the CPU.",
	})

	var nodeCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "network_host_count",
		Help: "Insolar network host count",
	})

	cpuTemp.Set(65.3)
	nodeCount.Set(float64(77))
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(nodeCount)

	http.Handle("/metrics", promhttp.Handler())
	logrus.Fatal(http.ListenAndServe(cfg.ListenAddress, nil))
}
