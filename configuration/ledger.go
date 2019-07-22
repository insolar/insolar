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

// Replica holds configuration for Replicator.
type Replica struct {
	// Role defines that should do replicator (subscribe on replica parent or send notifications to replica targets).
	Role string
	// Port that should be opened for replica communication.
	Port uint32
	// ParentAddress is an address to connect to replica parent.
	ParentAddress string
	// ParentPubKey is a public key that replica parent will use to sign replication records.
	ParentPubKey string
	// ScopesToReplicate is a list of DB scopes identifiers that define what should to replicate.
	ScopesToReplicate []byte
	// Attempts are a maximum count of attempts to connect to replica parent.
	Attempts int
	// DelayForAttempt is a time between connection attempt.
	DelayForAttempt time.Duration
	// DefaultBatchSize is a preferable count of records in a pull batch.
	DefaultBatchSize uint32
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

	// Replica holds configuration for Replicator.
	Replica Replica
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

		Replica: Replica{
			Role:              "root",
			Port:              20111,
			ParentAddress:     "127.0.0.1:20111",
			ParentPubKey:      "",
			ScopesToReplicate: []byte{2},
			Attempts:          60,
			DelayForAttempt:   1 * time.Second,
			DefaultBatchSize:  uint32(1000),
		},
	}
}
