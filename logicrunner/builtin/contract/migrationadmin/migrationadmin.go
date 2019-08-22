/*
 *
 *  Copyright  2019. Insolar Technologies GmbH
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package migrationadmin

import (
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

// MigrationAdmin manage and change status for  migration daemon.
type MigrationAdmin struct {
	foundation.BaseContract
	MigrationDaemon      foundation.StableMap
	MigrationAdminMember insolar.Reference
}

const (
	StatusActive     = "ACTIVE"
	StatusInactivate = "INACTIVE"
)

// Create new Migration admin in genesis.
func New(migrationDaemons [insolar.GenesisAmountMigrationDaemonMembers]insolar.Reference, migrationAdminMember insolar.Reference) (*MigrationAdmin, error) {
	daemonMigration := make(foundation.StableMap)
	for i := 0; i < insolar.GenesisAmountMigrationDaemonMembers; i++ {
		daemonMigration[migrationDaemons[i].String()] = StatusInactivate
	}
	return &MigrationAdmin{MigrationDaemon: daemonMigration, MigrationAdminMember: migrationAdminMember}, nil
}

// Return stable map migration daemon.
// ins:immutable
func (mA MigrationAdmin) GetAllMigrationDaemon() (foundation.StableMap, error) {
	sizeMap := len(mA.MigrationDaemon)
	if sizeMap != insolar.GenesisAmountMigrationDaemonMembers {
		return foundation.StableMap{}, fmt.Errorf(" MigrationAdmin contains the wrong amount migration daemon %d", sizeMap)
	}
	return mA.MigrationDaemon, nil
}

// Activate migration daemon.
func (mA MigrationAdmin) ActivateDaemon(daemonMember string, caller insolar.Reference) error {
	if caller != mA.MigrationAdminMember {
		return fmt.Errorf(" only migration admin can activate migration demons ")
	}
	switch mA.MigrationDaemon[daemonMember] {
	case StatusActive:
		return fmt.Errorf(" daemon member already activated - %s", daemonMember)
	case StatusInactivate:
		mA.MigrationDaemon[daemonMember] = StatusActive
		return nil
	default:
		return fmt.Errorf(" this referense is not daemon member ")
	}
}

// Deactivate migration daemon.
func (mA MigrationAdmin) DeactivateDaemon(daemonMember string, caller insolar.Reference) error {
	if caller != mA.MigrationAdminMember {
		return fmt.Errorf(" only migration admin can deactivate migration demons ")
	}
	switch mA.MigrationDaemon[daemonMember] {
	case StatusActive:
		mA.MigrationDaemon[daemonMember] = StatusInactivate
		return nil
	case StatusInactivate:
		return fmt.Errorf(" daemon member already deactivated - %s", daemonMember)
	default:
		return fmt.Errorf(" this referense is not daemon member ")
	}
}

// Check this member is migration daemon or mot.
// ins:immutable
func (mA MigrationAdmin) CheckActiveDaemon(daemonMember string) (bool, error) {
	status := mA.MigrationDaemon[daemonMember]
	if status == StatusActive {
		return true, nil
	}
	return false, fmt.Errorf(" this referense is not  active daemon member %s", daemonMember)
}

// Return only active daemons.
// ins:immutable
func (mA MigrationAdmin) GetActiveDaemons() ([]string, error) {
	var activeDaemons []string
	for daemonsRef, status := range mA.MigrationDaemon {
		if status == StatusActive {
			activeDaemons = append(activeDaemons, daemonsRef)
		}
	}
	return activeDaemons, nil
}
