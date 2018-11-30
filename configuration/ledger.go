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

package configuration

import (
	"github.com/insolar/insolar/core"
)

// Storage configures Ledger's storage.
type Storage struct {
	// DataDirectory is a directory where database's files live.
	DataDirectory string
	// TxRetriesOnConflict defines how many retries on transaction conflicts
	// storage update methods should do.
	TxRetriesOnConflict int
}

// JetCoordinator holds configuration for JetCoordinator.
type JetCoordinator struct {
	RoleCounts map[int]int
}

// ArtifactManager holds configuration for ArtifactManager.
type ArtifactManager struct {
	// Maximum pulse difference (NOT number of pulses) between current and the latest replicated on heavy.
	// IMPORTANT: It should be the same on ALL nodes.
	LightChainLimit core.PulseNumber
}

// PulseManager holds configuration for PulseManager.
type PulseManager struct {
	// HeavySyncEnabled enables replication to heavy (could be disabled for testing purposes)
	HeavySyncEnabled bool
	// HeavySyncMessageLimit soft limit of single message for replication to heavy.
	HeavySyncMessageLimit int
}

// RecentStorage holds configuration for RecentStorage
type RecentStorage struct {
	// Default TTL is a value of default ttl for redirects
	DefaultTTL int
}

// Ledger holds configuration for ledger.
type Ledger struct {
	// Storage defines storage configuration.
	Storage Storage
	// JetCoordinator defines jet coordinator configuration.
	JetCoordinator JetCoordinator
	// ArtifactManager holds configuration for ArtifactManager.
	ArtifactManager ArtifactManager
	// PulseManager holds configuration for PulseManager.
	PulseManager PulseManager
	// RecentStorage holds configuration for RecentStorage
	RecentStorage RecentStorage
}

// NewLedger creates new default Ledger configuration.
func NewLedger() Ledger {
	return Ledger{
		Storage: Storage{
			DataDirectory:       "./data",
			TxRetriesOnConflict: 3,
		},

		JetCoordinator: JetCoordinator{
			RoleCounts: map[int]int{
				int(core.DynamicRoleVirtualExecutor):  1,
				int(core.DynamicRoleHeavyExecutor):    1,
				int(core.DynamicRoleLightExecutor):    1,
				int(core.DynamicRoleVirtualValidator): 1,
				int(core.DynamicRoleLightValidator):   1,
			},
		},

		ArtifactManager: ArtifactManager{
			LightChainLimit: 10 * 30, // 30 pulses
		},

		PulseManager: PulseManager{
			HeavySyncEnabled:      true,
			HeavySyncMessageLimit: 1 << 20, // 1Mb
		},

		RecentStorage: RecentStorage{
			DefaultTTL: 50,
		},
	}
}
