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

// Storage configures Ledger's storage.
type Storage struct {
	// DataDirectory is a directory where database's files live.
	DataDirectory string
	// TxRetriesOnConflict defines how many retries on transaction conflicts
	// storage update methods should do.
	TxRetriesOnConflict int
}

// PulseManager holds configuration for PulseManager.
type PulseManager struct {
	// SplitThreshold is a drop size threshold in bytes to perform split.
	SplitThreshold uint64
}

// Backoff configures retry backoff algorithm
type Backoff struct {
	Factor float64
	// Jitter eases contention by randomizing backoff steps
	Jitter bool
	// Min and Max are the minimum and maximum values of the counter
	Min, Max time.Duration
	// MaxAttempts holds max count of attempts for a instance of Backoff
	MaxAttempts int
}

// Exporter holds configuration of Exporter
type Exporter struct {
	// ExportLag is lag in second before we start to export pulse
	ExportLag uint32
}

// Ledger holds configuration for ledger.
type Ledger struct {
	// Storage defines storage configuration.
	Storage Storage
	// PulseManager holds configuration for PulseManager.
	PulseManager PulseManager

	// common/sharable values:

	// LightChainLimit is maximum pulse difference (NOT number of pulses)
	// between current and the latest replicated on heavy.
	//
	// IMPORTANT: It should be the same on ALL nodes.
	LightChainLimit int

	// Exporter holds configuration of Exporter
	Exporter Exporter

	// PendingRequestsLimit holds a number of pending requests, what can be stored in the system
	// before they are declined
	PendingRequestsLimit int
}

// NewLedger creates new default Ledger configuration.
func NewLedger() Ledger {
	return Ledger{
		Storage: Storage{
			DataDirectory:       "./data",
			TxRetriesOnConflict: 3,
		},

		PulseManager: PulseManager{
			SplitThreshold: 10 * 100, // 10 megabytes.
		},
		LightChainLimit: 5, // 5 pulses

		Exporter: Exporter{
			ExportLag: 40, // 40 seconds
		},

		PendingRequestsLimit: 1000,
	}
}
