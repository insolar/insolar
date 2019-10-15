//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package configuration

import (
	"time"
)

// Exporter holds exporter configuration.
// Is was assumed, that exporter will be used for exporting data for observer
// Exporter is grpc-base service
type Exporter struct {
	// Addr specifies address where exporter server starts
	Addr string
	// DurationBetweenRequests specifies time limit between requests to exporter
	DurationBetweenRequests time.Duration
	// Number of pulses after which client can see finalized pulse
	ExportDelay int
}

// NewExporter creates new default configuration for export.
func NewExporter() Exporter {
	return Exporter{
		Addr:                    ":5678",
		DurationBetweenRequests: 10 * time.Second,
	}
}
