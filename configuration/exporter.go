// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package configuration

// Exporter holds exporter configuration.
// Is was assumed, that exporter will be used for exporting data for observer
// Exporter is grpc-base service
type Exporter struct {
	// Addr specifies address where exporter server starts
	Addr         string
	Auth         Auth
	CheckVersion bool
	// RateLimit specifies in/out limits for the Exporter API
	RateLimit RateLimit
}

// Auth specifies parameters for a token-based authorization of an observer
type Auth struct {
	// Allows turn off the JWT validation.
	// Set 'false' only for testing and when observer and exporter runs within the same secured environment
	Required bool
	// Used for validation an issuer of the JWT
	Issuer string
	// 512-bit secret for JWT validation
	Secret string
}

type RateLimit struct {
	// Required allows turn off the rate-limiting
	Required bool
	// In specifies the number of requests per second
	In Limits
	// Out specifies the number of responses per second for server_stream RPCs
	Out Limits
}

type Limits struct {
	// Global specifies a limit for all inbound requests or all outbound responses per second
	Global int
	// PerClient specifies a key - full name of gRPC method and value - number of requests or responses per second
	PerClient Handlers
}

type Handlers struct {
	RecordExport            int
	PulseExport             int
	PulseTopSyncPulse       int
	PulseNextFinalizedPulse int
}

func (h Handlers) Limit(method string) int {
	switch method {
	case "/exporter.RecordExporter/Export":
		return h.RecordExport
	case "/exporter.PulseExporter/Export":
		return h.PulseExport
	case "/exporter.PulseExporter/TopSyncPulse":
		return h.PulseTopSyncPulse
	case "/exporter.PulseExporter/NextFinalizedPulse":
		return h.PulseNextFinalizedPulse
	default:
		return 0
	}
}

// NewExporter creates new default configuration for export.
func NewExporter() Exporter {
	return Exporter{
		Addr: ":5678",
		Auth: Auth{
			Required: true,
			Issuer:   "Insolar-auth-service",
			Secret:   "",
		},
	}
}
