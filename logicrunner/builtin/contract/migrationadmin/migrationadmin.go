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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/migrationdaemon"
	"github.com/insolar/insolar/logicrunner/builtin/proxy/migrationshard"
)

// MigrationAdmin manage and change status for  migration daemon.
type MigrationAdmin struct {
	foundation.BaseContract

	MigrationAdminMember   insolar.Reference
	MigrationAddressShards [insolar.GenesisAmountMigrationAddressShards]insolar.Reference
	VestingParams          *VestingParams
}

type VestingParams struct {
	Lockup      int64 `json:"lockupInPulses"`
	Vesting     int64 `json:"vestingInPulses"`
	VestingStep int64 `json:"vestingStepInPulses"`
}

type CheckDaemonResponse struct {
	Status string `json:"status"`
}

const (
	StatusActive     = "active"
	StatusInactivate = "inactive"
)

//Call internal function migration admin from api.
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

func (mA *MigrationAdmin) getMigrationDamon(params map[string]interface{}, caller insolar.Reference) (*migrationdaemon.MigrationDaemon, error) {

	migrationDaemonMember, ok := params["reference"].(string)
	if !ok && len(migrationDaemonMember) == 0 {
		return nil, fmt.Errorf("incorect input: failed to get 'reference' param")
	}

	migrationDaemonMemberRef, err := insolar.NewReferenceFromBase58(migrationDaemonMember)
	if err != nil {
		return nil, fmt.Errorf(" failed to parse params.Reference")
	}

	migrationDaemonContractRef, err := foundation.GetMigrationDaemon(*migrationDaemonMemberRef)
	if err != nil {
		return nil, fmt.Errorf(" get migration daemon contract from foundation failed, %s ", err.Error())
	}
	if migrationDaemonContractRef.IsEmpty() {
		return nil, fmt.Errorf(" the member is not migration daemon, %s ", migrationDaemonMemberRef)
	}
	migrationDaemonContract := migrationdaemon.GetObject(migrationDaemonContractRef)

	return migrationDaemonContract, nil
}

func (mA *MigrationAdmin) activateDaemonCall(params map[string]interface{}, caller insolar.Reference) (interface{}, error) {
	if caller != mA.MigrationAdminMember {
		return nil, fmt.Errorf(" only migration admin can activate migration demons ")
	}
	migrationDaemonContract, err := mA.getMigrationDamon(params, caller)
	if err != nil {
		return nil, err
	}
	status, err := migrationDaemonContract.GetActivationStatus()
	if err != nil {
		return nil, err
	}
	if status {
		return nil, fmt.Errorf("daemon member already activated")
	}
	err = migrationDaemonContract.SetActivationStatus(true)
	return nil, err
}

func (mA *MigrationAdmin) deactivateDaemonCall(params map[string]interface{}, memberRef insolar.Reference) (interface{}, error) {
	if memberRef != mA.MigrationAdminMember {
		return nil, fmt.Errorf(" only migration admin can deactivate migration demons ")
	}
	migrationDaemonContract, err := mA.getMigrationDamon(params, memberRef)
	if err != nil {
		return nil, err
	}
	status, err := migrationDaemonContract.GetActivationStatus()
	if err != nil {
		return nil, err
	}
	if !status {
		return nil, fmt.Errorf("daemon member already deactivated")
	}
	err = migrationDaemonContract.SetActivationStatus(false)
	return nil, err
}

func (mA *MigrationAdmin) addMigrationAddressesCall(params map[string]interface{}, memberRef insolar.Reference) (interface{}, error) {
	migrationAddresses, ok := params["migrationAddresses"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("incorect input: failed to get 'migrationAddresses' param")
	}

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
	err := mA.addMigrationAddresses(migrationAddressesStr)
	if err != nil {
		return nil, fmt.Errorf("failed to add migration address: %s", err.Error())
	}

	return nil, nil
}

func (mA *MigrationAdmin) checkDaemonCall(params map[string]interface{}, caller insolar.Reference) (interface{}, error) {

	if caller != mA.MigrationAdminMember && !foundation.IsMigrationDaemonMember(caller) {
		return nil, fmt.Errorf(" permission denied to information about migration daemons")
	}
	migrationDaemonContract, err := mA.getMigrationDamon(params, caller)
	if err != nil {
		return nil, err
	}
	status, err := migrationDaemonContract.GetActivationStatus()
	if err != nil {
		return nil, err
	}
	if status {
		return CheckDaemonResponse{Status: StatusActive}, nil
	}
	return CheckDaemonResponse{Status: StatusInactivate}, nil
}

func (mA *MigrationAdmin) GetDepositParameters() (*VestingParams, error) {
	return mA.VestingParams, nil
}

// GetMemberByMigrationAddress gets member reference by burn address.
// ins:immutable
func (mA *MigrationAdmin) GetMemberByMigrationAddress(migrationAddress string) (*insolar.Reference, error) {
	trimmedMigrationAddress := foundation.TrimAddress(migrationAddress)
	i := foundation.GetShardIndex(trimmedMigrationAddress, insolar.GenesisAmountMigrationAddressShards)
	if i >= len(mA.MigrationAddressShards) {
		return nil, fmt.Errorf("incorect shard index")
	}
	s := migrationshard.GetObject(mA.MigrationAddressShards[i])
	refStr, err := s.GetRef(trimmedMigrationAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get reference in shard")
	}
	ref, err := insolar.NewReferenceFromBase58(refStr)
	if err != nil {
		return nil, errors.Wrap(err, "bad member reference for this migration address")
	}

	return ref, nil
}

// AddMigrationAddresses adds migration addresses to list.
// ins:immutable
func (mA *MigrationAdmin) addMigrationAddresses(migrationAddresses []string) error {
	newMA := [insolar.GenesisAmountMigrationAddressShards][]string{}
	for _, ma := range migrationAddresses {
		trimmedMigrationAddress := foundation.TrimAddress(ma)
		i := foundation.GetShardIndex(trimmedMigrationAddress, insolar.GenesisAmountMigrationAddressShards)
		if i >= len(newMA) {
			return fmt.Errorf("incorect migration shard index")
		}
		newMA[i] = append(newMA[i], trimmedMigrationAddress)
	}

	for i, ma := range newMA {
		if len(ma) == 0 {
			continue
		}
		s := migrationshard.GetObject(mA.MigrationAddressShards[i])
		err := s.AddFreeMigrationAddresses(ma)
		if err != nil {
			return errors.New("failed to add migration addresses to shard")
		}
	}

	return nil
}

// AddMigrationAddress adds migration address to list.
// ins:immutable
func (mA *MigrationAdmin) addMigrationAddress(migrationAddress string) error {
	trimmedMigrationAddress := foundation.TrimAddress(migrationAddress)
	i := foundation.GetShardIndex(trimmedMigrationAddress, insolar.GenesisAmountMigrationAddressShards)
	if i >= len(mA.MigrationAddressShards) {
		return fmt.Errorf("incorect migration shard index")
	}
	s := migrationshard.GetObject(mA.MigrationAddressShards[i])
	err := s.AddFreeMigrationAddresses([]string{trimmedMigrationAddress})
	if err != nil {
		return errors.New("failed to add migration address to shard")
	}

	return nil
}

// ins:immutable
func (mA *MigrationAdmin) GetFreeMigrationAddress(publicKey string) (string, error) {
	trimmedPublicKey := foundation.TrimPublicKey(publicKey)
	shardIndex := foundation.GetShardIndex(trimmedPublicKey, insolar.GenesisAmountPublicKeyShards)
	if shardIndex >= len(mA.MigrationAddressShards) {
		return "", fmt.Errorf("incorect migration address shard index")
	}

	for i := shardIndex; i < len(mA.MigrationAddressShards); i++ {
		mas := migrationshard.GetObject(mA.MigrationAddressShards[i])
		ma, err := mas.GetFreeMigrationAddress()

		if err == nil {
			return ma, nil
		}

		if err != nil {
			if !strings.Contains(err.Error(), "no more migration address left") {
				return "", errors.Wrap(err, "failed to set reference in migration address shard")
			}
		}
	}

	for i := 0; i < shardIndex; i++ {
		mas := migrationshard.GetObject(mA.MigrationAddressShards[i])
		ma, err := mas.GetFreeMigrationAddress()

		if err == nil {
			return ma, nil
		}

		if err != nil {
			if !strings.Contains(err.Error(), "no more migration address left") {
				return "", errors.Wrap(err, "failed to set reference in migration address shard")
			}
		}
	}

	return "", errors.New("no more migration addresses left in any shard")
}

// AddNewMemberToMaps adds new member to MigrationAddressMap.
// ins:immutable
func (mA *MigrationAdmin) AddNewMigrationAddressToMaps(migrationAddress string, memberRef insolar.Reference) error {
	trimmedMigrationAddress := foundation.TrimAddress(migrationAddress)
	shardIndex := foundation.GetShardIndex(trimmedMigrationAddress, insolar.GenesisAmountPublicKeyShards)
	if shardIndex >= len(mA.MigrationAddressShards) {
		return fmt.Errorf("incorect migration address shard index")
	}
	mas := migrationshard.GetObject(mA.MigrationAddressShards[shardIndex])
	err := mas.SetRef(migrationAddress, memberRef.String())
	if err != nil {
		return errors.Wrap(err, "failed to set reference in migration address shard")
	}

	return nil
}
