// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package heavy

import (
	"github.com/prometheus/client_golang/prometheus"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"

	"github.com/insolar/insolar/ledger/heavy/exporter"
)

const keyTagClientID = "exporter_client_id"

var (
	statBadgerStartTime = stats.Float64(
		"badger_start_time",
		"Time of last badger starting",
		stats.UnitMilliseconds,
	)
	statContractVersionClient = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "contract_version_exporter_client",
		Help: "What version of contracts is used by the exporter client",
	}, []string{keyTagClientID})

	statHeavyVersionClient = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "heavy_version_exporter_client",
		Help: "What version of heavy protocol is used by the exporter client",
	}, []string{keyTagClientID})
)

func init() {
	err := view.Register(
		&view.View{
			Name:        statBadgerStartTime.Name(),
			Description: statBadgerStartTime.Description(),
			Measure:     statBadgerStartTime,
			Aggregation: view.LastValue(),
		},
	)
	if err != nil {
		panic(err)
	}
}

func setPlatformVersionMetrics(actualVersionContract int64) {
	statContractVersionClient.WithLabelValues("heavy_exporter").Set(float64(actualVersionContract))
	statHeavyVersionClient.WithLabelValues("heavy_exporter").Set(float64(exporter.AllowedOnHeavyVersion))
}
