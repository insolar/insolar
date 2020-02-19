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
}

// NewExporter creates new default configuration for export.
func NewExporter() Exporter {
	return Exporter{
		Addr: ":5678",
	}
}
