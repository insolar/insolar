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
	Addr string
	Auth Auth
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
