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

import "github.com/prometheus/client_golang/prometheus"

var ParcelsSentTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: insolarNamespace,
		Name:      "parcels_sent_total",
		Help:      "Total number of parcels sent",
	},
	[]string{"messageType"},
)

var LocallyDeliveredParcelsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: insolarNamespace,
		Name:      "locally_delivered_parcels_total",
		Help:      "Total number of parcels delivered to the same machine",
	},
	[]string{"messageType"},
)

var ParcelsSentSizeBytes = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:  insolarNamespace,
		Name:       "parcels_sent_size_bytes",
		Help:       "Size of sent parcels",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
	},
	[]string{"messageType"},
)

var ParcelsReplySizeBytes = prometheus.NewSummaryVec(
	prometheus.SummaryOpts{
		Namespace:  insolarNamespace,
		Name:       "parcels_reply_size_bytes",
		Help:       "Size of replies to parcels",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
	},
	[]string{"messageType"},
)
