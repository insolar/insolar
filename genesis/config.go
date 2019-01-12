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

// It's very light check. It's not about majority rule
func hasMinimumRolesSet(conf *genesisConfig) error {
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

	err = hasMinimumRolesSet(conf)
	if err != nil {
		return nil, errors.Wrap(err, "[ parseGenesisConfig ]")
	}

	return conf, nil
}
