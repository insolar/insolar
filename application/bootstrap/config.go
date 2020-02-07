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

// Node contains info about discovery nodes.
type Node struct {
	Host string `mapstructure:"host" yaml:"host"`
	Role string `mapstructure:"role" yaml:"role"`
	// KeyName is used for generating keys file for discovery nodes.
	KeyName string `mapstructure:"key_name" yaml:"key_name"`
	// CertName is used for generating cert file for discovery nodes.
	CertName string `mapstructure:"cert_name" yaml:"cert_name"`
	// KeysFile is used to set path to keys file in insolard config for every non-discovery nodes
	// (used only by generate_insolar_config.go).
	KeysFile string `mapstructure:"keys_file" yaml:"keys_file"`
}

// Config contains configuration required for bootstrap.
type Config struct {
	// MembersKeysDir is the root key place.
	MembersKeysDir string `mapstructure:"members_keys_dir" yaml:"members_keys_dir"`
	// DiscoveryKeysDir is a default directory where save keys for discovery nodes.
	DiscoveryKeysDir string `mapstructure:"discovery_keys_dir" yaml:"discovery_keys_dir"`
	// NotDiscoveryKeysDir is a default directory where save keys for discovery nodes.
	NotDiscoveryKeysDir string `mapstructure:"not_discovery_keys_dir" yaml:"not_discovery_keys_dir"`
	// KeysNameFormat is the default key file name format for discovery nodes.
	KeysNameFormat string `mapstructure:"keys_name_format" yaml:"keys_name_format"`
	// ReuseKeys is a flag to reuse discovery nodes keys (don't use if your not understand how it works)
	ReuseKeys bool `mapstructure:"reuse_keys" yaml:"reuse_keys"`

	HeavyGenesisConfigFile string `mapstructure:"heavy_genesis_config_file" yaml:"heavy_genesis_config_file"`

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

	// Discovery settings.

	MajorityRule int `mapstructure:"majority_rule" yaml:"majority_rule"`
	MinRoles     struct {
		Virtual       uint `mapstructure:"virtual" yaml:"virtual"`
		HeavyMaterial uint `mapstructure:"heavy_material" yaml:"heavy_material"`
		LightMaterial uint `mapstructure:"light_material" yaml:"light_material"`
	} `mapstructure:"min_roles" yaml:"min_roles"`
	// DiscoveryNodes is a discovery nodes list.
	DiscoveryNodes []Node `mapstructure:"discovery_nodes" yaml:"discovery_nodes"`

	// Nodes is used only by generate_insolar_config.go
	Nodes []Node `mapstructure:"nodes" yaml:"nodes"`

	// PulsarPublicKeys is the pulsar's public keys for pulses validation
	// (not in use, just for future features).
	PulsarPublicKeys []string `mapstructure:"pulsar_public_keys" yaml:"pulsar_public_keys"`
}

// hasMinimumRolesSet does basic check (it's not about majority rule).
func hasMinimumRolesSet(conf *Config) error {
	minRequiredRolesSet := map[string]bool{
		"virtual":        true,
		"heavy_material": true,
		"light_material": true,
	}

	for _, discNode := range conf.DiscoveryNodes {
		delete(minRequiredRolesSet, discNode.Role)
	}

	for _, node := range conf.Nodes {
		delete(minRequiredRolesSet, node.Role)
	}

	if len(minRequiredRolesSet) != 0 {
		var missingRoles string
		for role := range minRequiredRolesSet {
			missingRoles += role + ", "
		}
		return errors.New("no required roles in config: " + missingRoles)
	}

	return nil
}

// ParseConfig parse bootstrap config.
func ParseConfig(path string) (*Config, error) {
	var conf = &Config{}
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

	err = hasMinimumRolesSet(conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
