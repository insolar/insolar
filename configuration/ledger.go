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
	// DataDirectoryNewDB is a directory where new database's files live.
	DataDirectoryNewDB string
	// TxRetriesOnConflict defines how many retries on transaction conflicts
	// storage update methods should do.
	TxRetriesOnConflict int
}

// PulseManager holds configuration for PulseManager.
type PulseManager struct {
	// // HeavySyncEnabled enables replication to heavy (could be disabled for testing purposes)
	// HeavySyncEnabled bool
	// // HeavySyncMessageLimit soft limit of single message for replication to heavy.
	// HeavySyncMessageLimit int
	// // Backoff configures retry backoff algorithm for Heavy Sync
	// HeavyBackoff Backoff
	// SplitThreshold is a drop size threshold in bytes to perform split.
	SplitThreshold uint64
}

// LightToHeavySync holds settings for a light to heavy sync process
type LightToHeavySync struct {
	// Backoff holds a backoff configuration for failed sendings of payload from a light to a heavy
	Backoff Backoff
	// RetryLoopDuration holds a value of a light's sync process frequency
	RetryLoopDuration time.Duration
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

// RecentStorage holds configuration for RecentStorage
type RecentStorage struct {
	// Default TTL is a value of default ttl for redirects
	DefaultTTL int
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
	// RecentStorage holds configuration for RecentStorage
	RecentStorage RecentStorage

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

	// LightToHeavySync holds settings for a light to heavy sync process
	LightToHeavySync LightToHeavySync
}

// NewLedger creates new default Ledger configuration.
func NewLedger() Ledger {
	return Ledger{
		Storage: Storage{
			DataDirectory:       "./data",
			DataDirectoryNewDB:  "./new-data",
			TxRetriesOnConflict: 3,
		},

		PulseManager: PulseManager{
			// HeavySyncEnabled:      true,
			// HeavySyncMessageLimit: 1 << 20, // 1Mb
			// HeavyBackoff:
			SplitThreshold: 10 * 100, // 10 megabytes.
		},

		LightToHeavySync: LightToHeavySync{
			Backoff: Backoff{
				Jitter:      true,
				Min:         200 * time.Millisecond,
				Max:         2 * time.Second,
				Factor:      2,
				MaxAttempts: 10,
			},
			RetryLoopDuration: 1 * time.Second,
		},

		RecentStorage: RecentStorage{
			DefaultTTL: 10,
		},

		LightChainLimit: 5, // 5 pulses

		Exporter: Exporter{
			ExportLag: 40, // 40 seconds
		},

		PendingRequestsLimit: 1000,
	}
}
