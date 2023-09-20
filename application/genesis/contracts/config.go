package contracts

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
