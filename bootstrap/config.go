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

// Node contains info about discovery nodes
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
	Insgocc string
	OutDir  string
}

// Config contains all genesis config
type Config struct {
	// RootKeysFile is the root key place.
	RootKeysFile string `mapstructure:"root_keys_file"`
	// DiscoveryKeysDir is a default directory where save keys for discovery nodes.
	DiscoveryKeysDir string `mapstructure:"discovery_keys_dir"`
	// KeysNameFormat is the default key file name format for discovery nodes.
	KeysNameFormat string `mapstructure:"keys_name_format"`
	// ReuseKeys is a flag to reuse discovery nodes keys (don't use if your not understand how it works)
	ReuseKeys bool `mapstructure:"reuse_keys"`

	HeavyGenesisConfigFile string `mapstructure:"heavy_genesis_config_file"`
	HeavyGenesisPluginsDir string `mapstructure:"heavy_genesis_plugins_dir"`

	// RootBalance is a start balance for the root member's wallet.
	RootBalance uint `mapstructure:"root_balance"`
	Contracts   Contracts

	// Discovery settings.

	MajorityRule int `mapstructure:"majority_rule"`
	MinRoles     struct {
		Virtual       uint `mapstructure:"virtual"`
		HeavyMaterial uint `mapstructure:"heavy_material"`
		LightMaterial uint `mapstructure:"light_material"`
	} `mapstructure:"min_roles"`
	// DiscoveryNodes is a discovery nodes list.
	DiscoveryNodes []Node `mapstructure:"discovery_nodes"`

	// Nodes is not need on genesis and only used by generate_insolar_config.go
	Nodes []Node `mapstructure:"nodes"`

	// PulsarPublicKeys is the pulsar's public keys for pulses  validation
	// (not in use, just for future features).
	PulsarPublicKeys []string `mapstructure:"pulsar_public_keys"`
}

// It's very light check. It's not about majority rule
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
		return errors.New("[ hasMinimumRolesSet ] No required roles in genesis config: " + missingRoles)
	}

	return nil
}

// ParseGenesisConfig parse genesis config
func ParseGenesisConfig(path string) (*Config, error) {
	var conf = &Config{}
	v := viper.New()
	v.SetConfigFile(path)
	err := v.ReadInConfig()
	if err != nil {
		return nil, errors.Wrap(err, "[ parseGenesisConfig ] couldn't read config file")
	}
	err = v.Unmarshal(conf)
	if err != nil {
		return nil, errors.Wrap(err, "[ parseGenesisConfig ] couldn't unmarshal yaml to struct")
	}

	err = hasMinimumRolesSet(conf)
	if err != nil {
		return nil, errors.Wrap(err, "[ parseGenesisConfig ]")
	}

	return conf, nil
}
