// Copyright 2020 Insolar Network Ltd.
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

package bootstrap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/insolar/insolar/applicationbase/bootstrap"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/application/appfoundation"
	"github.com/insolar/insolar/insolar/secrets"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

func readMigrationAddresses(config *ContractsConfig) ([][]string, error) {
	file := filepath.Join(config.MembersKeysDir, "migration_addresses.json")
	result := make([][]string, config.MAShardCount)
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return result, errors.Wrapf(err, " couldn't read migration addresses file %v", file)
	}

	var ma []string
	err = json.NewDecoder(bytes.NewReader(b)).Decode(&ma)
	if err != nil {
		return result, errors.Wrapf(err, "fail unmarshal migration addresses data")
	}

	for _, a := range ma {
		if appfoundation.IsEthereumAddress(a) {
			address := foundation.TrimAddress(a)
			i := foundation.GetShardIndex(address, config.MAShardCount)
			result[i] = append(result[i], address)
		}
	}
	return result, nil
}

// CreateGenesisContractsConfig generates bootstrap data.
//
// 1. read application-related keys files.
// 2. generates genesis contracts config for heavy node.
func CreateGenesisContractsConfig(ctx context.Context, configFile string) (map[string]interface{}, error) {
	config, err := ParseContractsConfig(configFile)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[ bootstrap ] config:\n%v\n", bootstrap.DumpAsJSON(config))

	inslog := inslogger.FromContext(ctx)

	inslog.Info("[ bootstrap ] read keys files")
	rootPublicKey, err := secrets.GetPublicKeyFromFile(filepath.Join(config.MembersKeysDir, "root_member_keys.json"))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get root keys")
	}

	feePublicKey, err := secrets.GetPublicKeyFromFile(filepath.Join(config.MembersKeysDir, "fee_member_keys.json"))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get fees keys")
	}

	migrationAdminPublicKey, err := secrets.GetPublicKeyFromFile(
		filepath.Join(config.MembersKeysDir, "migration_admin_member_keys.json"))
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get migration admin keys")
	}
	migrationDaemonPublicKeys := []string{}
	for i := 0; i < application.GenesisAmountMigrationDaemonMembers; i++ {
		k, err := secrets.GetPublicKeyFromFile(config.MembersKeysDir + GetMigrationDaemonPath(i))
		if err != nil {
			return nil, errors.Wrap(err, "couldn't get migration daemon keys")
		}
		migrationDaemonPublicKeys = append(migrationDaemonPublicKeys, k)
	}

	networkIncentivesPublicKeys := []string{}
	for i := 0; i < application.GenesisAmountNetworkIncentivesMembers; i++ {
		k, err := secrets.GetPublicKeyFromFile(
			filepath.Join(config.MembersKeysDir, GetFundPath(i, "network_incentives_")))
		if err != nil {
			return nil, errors.Wrap(err, "couldn't get network incentives keys")
		}
		networkIncentivesPublicKeys = append(networkIncentivesPublicKeys, k)
	}

	applicationIncentivesPublicKeys := []string{}
	for i := 0; i < application.GenesisAmountApplicationIncentivesMembers; i++ {
		k, err := secrets.GetPublicKeyFromFile(
			filepath.Join(config.MembersKeysDir, GetFundPath(i, "application_incentives_")))
		if err != nil {
			return nil, errors.Wrap(err, "couldn't get application incentives keys")
		}
		applicationIncentivesPublicKeys = append(applicationIncentivesPublicKeys, k)
	}

	foundationPublicKeys := []string{}
	for i := 0; i < application.GenesisAmountFoundationMembers; i++ {
		k, err := secrets.GetPublicKeyFromFile(
			filepath.Join(config.MembersKeysDir, GetFundPath(i, "foundation_")))
		if err != nil {
			return nil, errors.Wrap(err, "couldn't get foundation keys")
		}
		foundationPublicKeys = append(foundationPublicKeys, k)
	}

	enterprisePublicKeys := []string{}
	for i := 0; i < application.GenesisAmountEnterpriseMembers; i++ {
		k, err := secrets.GetPublicKeyFromFile(
			filepath.Join(config.MembersKeysDir, GetFundPath(i, "enterprise_")))
		if err != nil {
			return nil, errors.Wrap(err, "couldn't get enterprise keys")
		}
		enterprisePublicKeys = append(enterprisePublicKeys, k)
	}

	if config.MAShardCount <= 0 {
		panic(fmt.Sprintf("[genesis] store contracts failed: setup ma_shard_count parameter, current value %v", config.MAShardCount))
	}

	inslog.Info("[ bootstrap ] read migration addresses ...")
	migrationAddresses, err := readMigrationAddresses(config)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get migration addresses")
	}

	vestingStep := config.VestingStepInPulses
	if vestingStep == 0 {
		vestingStep = 60 * 60 * 24
	}

	return map[string]interface{}{
		"RootBalance":                     config.RootBalance,
		"MDBalance":                       config.MDBalance,
		"RootPublicKey":                   rootPublicKey,
		"FeePublicKey":                    feePublicKey,
		"MigrationAdminPublicKey":         migrationAdminPublicKey,
		"MigrationDaemonPublicKeys":       migrationDaemonPublicKeys,
		"NetworkIncentivesPublicKeys":     networkIncentivesPublicKeys,
		"ApplicationIncentivesPublicKeys": applicationIncentivesPublicKeys,
		"FoundationPublicKeys":            foundationPublicKeys,
		"EnterprisePublicKeys":            enterprisePublicKeys,
		"MigrationAddresses":              migrationAddresses,
		"VestingPeriodInPulses":           config.VestingPeriodInPulses,
		"LockupPeriodInPulses":            config.LockupPeriodInPulses,
		"VestingStepInPulses":             vestingStep,
		"MAShardCount":                    config.MAShardCount,
		"PKShardCount":                    config.PKShardCount,
	}, nil
}

// GetMigrationDaemonPath generate key file name for migration daemon
func GetMigrationDaemonPath(i int) string {
	return "migration_daemon_" + strconv.Itoa(i) + "_member_keys.json"
}

// GetFundPath generate key file name for composite name
func GetFundPath(i int, prefix string) string {
	return prefix + strconv.Itoa(i) + "_member_keys.json"
}
