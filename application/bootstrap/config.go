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
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// ContractsConfig contains configuration required for bootstrap application logic.
type ContractsConfig struct {
	// MembersKeysDir is the root key place.
	MembersKeysDir string `mapstructure:"members_keys_dir" yaml:"members_keys_dir"`
	// RootBalance is a start balance for the root member's wallet.
	RootBalance string `mapstructure:"root_balance" yaml:"root_balance"`
	// MDBalance is a start balance for the migration admin member's wallet.
	MDBalance string `mapstructure:"md_balance" yaml:"md_balance"`
	// VestingPeriodInPulses - interval of count pulses for the full period of partial release.
	VestingPeriodInPulses int64 `mapstructure:"vesting_pulse_period" yaml:"vesting_pulse_period"`
	// VestingPeriodInPulses - interval of count pulses for one step of partial release.
	VestingStepInPulses int64 `mapstructure:"vesting_pulse_step" yaml:"vesting_pulse_step"`
	// LockupPeriodInPulses - interval of count pulses for the full period of hold.
	LockupPeriodInPulses int64 `mapstructure:"lockup_pulse_period" yaml:"lockup_pulse_period"`
	PKShardCount         int   `mapstructure:"pk_shard_count" yaml:"pk_shard_count"`
	MAShardCount         int   `mapstructure:"ma_shard_count" yaml:"ma_shard_count"`
}

// ParseContractsConfig parse bootstrap contracts config.
func ParseContractsConfig(path string) (*ContractsConfig, error) {
	var conf = &ContractsConfig{}
	v := viper.New()
	v.SetConfigFile(path)
	err := v.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't read config file")
	}
	err = v.Unmarshal(conf)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't unmarshal yaml to struct")
	}
	return conf, nil
}
