package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// NetworkPacketSentTotal is total number of sent packets metric
var NetworkPacketSentTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name:      "packet_sent_total",
	Help:      "Total number of sent packets",
	Namespace: insolarNamespace,
	Subsystem: "network",
}, []string{"packetType"})

// NetworkPacketReceivedTotal is is total number of received packets metric
var NetworkPacketReceivedTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name:      "packet_received_total",
	Help:      "Total number of received packets",
	Namespace: insolarNamespace,
	Subsystem: "network",
}, []string{"packetType"})

// NetworkSentSize is total sent bytes
var NetworkSentSize = prometheus.NewSummary(prometheus.SummaryOpts{
	Name:      "sent_bytes",
	Help:      "Sent by transport",
	Namespace: insolarNamespace,
	Subsystem: "network",
})

// NetworkRecvSize is total received bytes
var NetworkRecvSize = prometheus.NewSummary(prometheus.SummaryOpts{
	Name:      "recv_bytes",
	Help:      "Received by transport",
	Namespace: insolarNamespace,
	Subsystem: "network",
})

// NetworkPacketTimeoutTotal is is total number of timed out packets metric
var NetworkPacketTimeoutTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name:      "packet_timeout_total",
	Help:      "Total number of timed out packets",
	Namespace: insolarNamespace,
	Subsystem: "network",
}, []string{"packetType"})

// NetworkFutures is current network transport futures count metric
var NetworkFutures = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name:      "futures",
	Help:      "Current network transport futures count",
	Namespace: insolarNamespace,
	Subsystem: "network",
}, []string{"packetType"})

// NetworkConnections is current network transport connections count metric
var NetworkConnections = prometheus.NewGauge(prometheus.GaugeOpts{
	Name:      "connections",
	Help:      "Connections count in connection pool",
	Namespace: insolarNamespace,
	Subsystem: "network",
})
