// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package bootstrap

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// certNamesStartFrom is an offset certificate name index starts from.
var certNamesStartFrom = 1

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
	// DiscoveryKeysDir is a default directory where save keys for discovery nodes.
	DiscoveryKeysDir string `mapstructure:"discovery_keys_dir" yaml:"discovery_keys_dir"`
	// NotDiscoveryKeysDir is a default directory where save keys for discovery nodes.
	NotDiscoveryKeysDir string `mapstructure:"not_discovery_keys_dir" yaml:"not_discovery_keys_dir"`
	// CertificatesOutDir is a directory where to save generated cert files.
	CertificatesOutDir string `mapstructure:"certificates_out_dir" yaml:"certificates_out_dir"`
	// CertificateNameOffsetFromZero specifies if starting index number starts from zero.
	// Mostly exists because of launchnet, where node names uses numeric suffix started from 1.
	// TODO: get rid of launchnet and this parameter - @nordicdyno 17.02.2020
	CertificateNameOffsetFromZero bool `mapstructure:"certificate_name_offset_from_zero" yaml:"certificate_name_offset_from_zero"`
	// KeysNameFormat is the default key file name format for discovery nodes.
	KeysNameFormat string `mapstructure:"keys_name_format" yaml:"keys_name_format"`
	// ReuseKeys is a flag to reuse discovery nodes keys (don't use if your not understand how it works)
	ReuseKeys bool `mapstructure:"reuse_keys" yaml:"reuse_keys"`

	HeavyGenesisConfigFile string `mapstructure:"heavy_genesis_config_file" yaml:"heavy_genesis_config_file"`

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
