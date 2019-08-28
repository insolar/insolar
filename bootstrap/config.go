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

package bootstrap

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Node contains info about discovery nodes.
type Node struct {
	Host string `mapstructure:"host"`
	Role string `mapstructure:"role"`
	// KeyName is used for generating keys file for discovery nodes.
	KeyName string `mapstructure:"key_name"`
	// CertName is used for generating cert file for discovery nodes.
	CertName string `mapstructure:"cert_name"`
	// KeysFile is used to set path to keys file in insolard config for every non-discovery nodes
	// (used only by generate_insolar_config.go).
	KeysFile string `mapstructure:"keys_file"`
}

// Contracts contains config for contract's plugins generation.
type Contracts struct {
	// Insgocc is the path to ingocc binary for plugins generation.
	Insgocc string
	// OutDir is the path to directory where plugins so files would be saved.
	OutDir string
}

// Config contains configuration required for bootstrap.
type Config struct {
	// MembersKeysDir is the root key place.
	MembersKeysDir string `mapstructure:"members_keys_dir"`
	// DiscoveryKeysDir is a default directory where save keys for discovery nodes.
	DiscoveryKeysDir string `mapstructure:"discovery_keys_dir"`
	// KeysNameFormat is the default key file name format for discovery nodes.
	KeysNameFormat string `mapstructure:"keys_name_format"`
	// ReuseKeys is a flag to reuse discovery nodes keys (don't use if your not understand how it works)
	ReuseKeys bool `mapstructure:"reuse_keys"`

	HeavyGenesisConfigFile string `mapstructure:"heavy_genesis_config_file"`
	HeavyGenesisPluginsDir string `mapstructure:"heavy_genesis_plugins_dir"`

	// RootBalance is a start balance for the root member's wallet.
	RootBalance string `mapstructure:"root_balance"`
	// MDBalance is a start balance for the migration admin member's wallet.
	MDBalance string `mapstructure:"md_balance"`
	// VestingPeriodInPulses - interval of count pulses for the full period of partial release.
	VestingPeriodInPulses int64 `mapstructure:"vesting_pulse_period"`
	// VestingPeriodInPulses - interval of count pulses for one step of partial release.
	VestingStepInPulses int64 `mapstructure:"vesting_pulse_step"`
	// LokupPeriodInPulses - interval of count pulses for the full period of hold.
	LokupPeriodInPulses int64 `mapstructure:"lokup_pulse_period"`
	Contracts           Contracts

	// Discovery settings.

	MajorityRule int `mapstructure:"majority_rule"`
	MinRoles     struct {
		Virtual       uint `mapstructure:"virtual"`
		HeavyMaterial uint `mapstructure:"heavy_material"`
		LightMaterial uint `mapstructure:"light_material"`
	} `mapstructure:"min_roles"`
	// DiscoveryNodes is a discovery nodes list.
	DiscoveryNodes []Node `mapstructure:"discovery_nodes"`

	// Nodes is used only by generate_insolar_config.go
	Nodes []Node `mapstructure:"nodes"`

	// PulsarPublicKeys is the pulsar's public keys for pulses validation
	// (not in use, just for future features).
	PulsarPublicKeys []string `mapstructure:"pulsar_public_keys"`
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
