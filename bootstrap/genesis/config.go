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

package genesis

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Node contains info about discovery nodes
type Node struct {
	Host     string `mapstructure:"host"`
	Role     string `mapstructure:"role"`
	KeysFile string `mapstructure:"keys_file"`
	CertName string `mapstructure:"cert_name"`
}

// Config contains all genesis config
type Config struct {
	RootKeysFile     string `mapstructure:"root_keys_file"`
	DiscoveryKeysDir string `mapstructure:"discovery_keys_dir"`
	KeysNameFormat   string `mapstructure:"keys_name_format"`
	ReuseKeys        bool   `mapstructure:"reuse_keys"`
	RootBalance      uint   `mapstructure:"root_balance"`
	MajorityRule     int    `mapstructure:"majority_rule"`
	MinRoles         struct {
		Virtual       uint `mapstructure:"virtual"`
		HeavyMaterial uint `mapstructure:"heavy_material"`
		LightMaterial uint `mapstructure:"light_material"`
	} `mapstructure:"min_roles"`
	PulsarPublicKeys []string `mapstructure:"pulsar_public_keys"`
	DiscoveryNodes   []Node   `mapstructure:"discovery_nodes"`
	Nodes            []Node   `mapstructure:"nodes"`
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
