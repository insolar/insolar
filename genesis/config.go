package genesis

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type discovery struct {
	Host     string `mapstructure:"host"`
	Role     string `mapstructure:"role"`
	KeysFile string `mapstructure:"keys_file"`
	CertName string `mapstructure:"cert_name"`
}

type genesisConfig struct {
	RootKeysFile string `mapstructure:"root_keys_file"`
	RootBalance  uint   `mapstructure:"root_balance"`
	MajorityRule int    `mapstructure:"majority_rule"`
	MinRoles     struct {
		Virtual       uint `mapstructure:"virtual"`
		HeavyMaterial uint `mapstructure:"heavy_material"`
		LightMaterial uint `mapstructure:"light_material"`
	} `mapstructure:"min_roles"`
	PulsarPublicKeys []string    `mapstructure:"pulsar_public_keys"`
	DiscoveryNodes   []discovery `mapstructure:"discovery_nodes"`
}

func parseGenesisConfig(path string) (*genesisConfig, error) {
	var conf = &genesisConfig{}
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
	return conf, nil
}
