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
