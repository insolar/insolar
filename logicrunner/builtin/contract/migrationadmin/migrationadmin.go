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
	"strings"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/rootdomain"
)

// MigrationAdmin manage and change status for  migration daemon.
type MigrationAdmin struct {
	foundation.BaseContract
	MigrationDaemons     foundation.StableMap
	MigrationAdminMember insolar.Reference
	Lokup                int64
	Vesting              int64
}

const (
	StatusActive     = "active"
	StatusInactivate = "inactive"
)

func (mA *MigrationAdmin) MigrationAdminCall(params map[string]interface{}, nameMethod string, caller insolar.Reference) (interface{}, error) {

	switch nameMethod {
	case "addAddresses":
		return mA.addMigrationAddressesCall(params, caller)

	case "activateDaemon":
		return mA.activateDaemonCall(params, caller)

	case "deactivateDaemon":
		return mA.deactivateDaemonCall(params, caller)

	case "checkDaemon":
		return mA.checkDaemonCall(params, caller)
	}
	return nil, fmt.Errorf("unknown method: migration.'%s'", nameMethod)
}

func (mA *MigrationAdmin) activateDaemonCall(params map[string]interface{}, memberRef insolar.Reference) (interface{}, error) {
	migrationDaemon, ok := params["reference"].(string)
	if !ok && len(migrationDaemon) == 0 {
		return nil, fmt.Errorf("incorect input: failed to get 'reference' param")
	}
	return nil, mA.ActivateDaemon(strings.TrimSpace(migrationDaemon), memberRef)
}

func (mA *MigrationAdmin) deactivateDaemonCall(params map[string]interface{}, memberRef insolar.Reference) (interface{}, error) {
	migrationDaemon, ok := params["reference"].(string)

	if !ok && len(migrationDaemon) == 0 {
		return nil, fmt.Errorf("incorect input: failed to get 'reference' param")
	}
	return nil, mA.DeactivateDaemon(strings.TrimSpace(migrationDaemon), memberRef)
}

func (mA *MigrationAdmin) addMigrationAddressesCall(params map[string]interface{}, memberRef insolar.Reference) (interface{}, error) {
	migrationAddresses, ok := params["migrationAddresses"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'migrationAddresses' param")
	}

	rootDomain := rootdomain.GetObject(foundation.GetRootDomain())

	if memberRef != mA.MigrationAdminMember {
		return nil, fmt.Errorf("only migration daemon admin can call this method")
	}

	migrationAddressesStr := make([]string, len(migrationAddresses))

	for i, ba := range migrationAddresses {
		migrationAddress, ok := ba.(string)
		if !ok {
			return nil, fmt.Errorf("failed to 'migrationAddresses' param")
		}
		migrationAddressesStr[i] = migrationAddress
	}
	err := rootDomain.AddMigrationAddresses(migrationAddressesStr)
	if err != nil {
		return nil, fmt.Errorf("failed to add migration address: %s", err.Error())
	}

	return nil, nil
}

type CheckDaemonResponse struct {
	Status string `json:"status"`
}

func (mA *MigrationAdmin) checkDaemonCall(params map[string]interface{}, caller insolar.Reference) (interface{}, error) {
	migrationDaemon, ok := params["reference"].(string)

	_, err := mA.CheckDaemon(strings.TrimSpace(caller.String()))
	if caller != mA.MigrationAdminMember && err != nil {
		return nil, fmt.Errorf(" permission denied to information about migration daemons: %s", err.Error())
	}

	if !ok && len(migrationDaemon) == 0 {
		return nil, fmt.Errorf("incorect input: failed to get 'reference' param")
	}

	result, err := mA.CheckDaemon(strings.TrimSpace(migrationDaemon))
	if err != nil {
		return nil, fmt.Errorf(" check status migration daemon failed: %s", err.Error())
	}
	if result {
		return CheckDaemonResponse{Status: StatusActive}, nil
	}
	return CheckDaemonResponse{Status: StatusInactivate}, nil
}

// Return stable map migration daemon.
// ins:immutable
func (mA *MigrationAdmin) GetAllMigrationDaemon() (foundation.StableMap, error) {
	sizeMap := len(mA.MigrationDaemons)
	if sizeMap != insolar.GenesisAmountMigrationDaemonMembers {
		return foundation.StableMap{}, fmt.Errorf(" MigrationAdmin contains the wrong amount migration daemon %d", sizeMap)
	}
	return mA.MigrationDaemons, nil
}

// Activate migration daemon.
func (mA *MigrationAdmin) ActivateDaemon(daemonMember string, caller insolar.Reference) error {
	if caller != mA.MigrationAdminMember {
		return fmt.Errorf(" only migration admin can activate migration demons ")
	}
	switch mA.MigrationDaemons[daemonMember] {
	case StatusActive:
		return fmt.Errorf(" daemon member already activated - %s", daemonMember)
	case StatusInactivate:
		mA.MigrationDaemons[daemonMember] = StatusActive
		return nil
	default:
		return fmt.Errorf(" this referense is not daemon member ")
	}
}

// Deactivate migration daemon.
func (mA *MigrationAdmin) DeactivateDaemon(daemonMember string, caller insolar.Reference) error {
	if caller != mA.MigrationAdminMember {
		return fmt.Errorf(" only migration admin can deactivate migration demons ")
	}

	switch mA.MigrationDaemons[daemonMember] {
	case StatusActive:
		mA.MigrationDaemons[daemonMember] = StatusInactivate
		return nil
	case StatusInactivate:
		return fmt.Errorf(" daemon member already deactivated - %s", daemonMember)
	default:
		return fmt.Errorf(" this referense is not daemon member ")
	}
}

// Check this member is migration daemon or mot.
// ins:immutable
func (mA *MigrationAdmin) CheckDaemon(daemonMember string) (bool, error) {
	switch mA.MigrationDaemons[daemonMember] {
	case StatusActive:
		return true, nil
	case StatusInactivate:
		return false, nil
	default:
		return false, fmt.Errorf(" this reference is not daemon member %s", daemonMember)
	}
}

// Return only active daemons.
// ins:immutable
func (mA *MigrationAdmin) GetActiveDaemons() ([]string, error) {
	var activeDaemons []string
	for daemonsRef, status := range mA.MigrationDaemons {
		if status == StatusActive {
			activeDaemons = append(activeDaemons, daemonsRef)
		}
	}
	return activeDaemons, nil
}

func (mA MigrationAdmin) GetDepositParameters() (int64, int64, error) {
	return mA.Lokup, mA.Vesting, nil
}
