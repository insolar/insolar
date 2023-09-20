package configuration

import (
	"time"
)

// Metrics holds configuration for metrics publishing.
type Metrics struct {
	ListenAddress string
	Namespace     string
	ZpagesEnabled bool
	// ReportingPeriod defines exporter reporting period
	// if zero, exporter uses default value (1s)
	ReportingPeriod time.Duration
}

// NewMetrics creates new default configuration for metrics publishing.
func NewMetrics() Metrics {
	return Metrics{
		ListenAddress: "0.0.0.0:9091",
		Namespace:     "insolar",
		ZpagesEnabled: true,
	}
}
