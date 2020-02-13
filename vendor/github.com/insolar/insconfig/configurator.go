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

package insconfig

import (
	goflag "flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	flag "github.com/spf13/pflag"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// this should be implemented by local config struct
type ConfigStruct interface {
	GetConfig() interface{}
}

type Params struct {
	ConfigStruct ConfigStruct
	// Prefix for environment variables
	EnvPrefix string
	// Custom viper decoding hooks
	ViperHooks []mapstructure.DecodeHookFunc
	// Should return config path
	ConfigPathGetter ConfigPathGetter
	// If set then return error on file not found
	FileRequired bool
}

type ConfigPathGetter interface {
	GetConfigPath() string
}

// Adds "--config" flag and read path from it
type DefaultConfigPathGetter struct {
	// For go flags compatibility
	GoFlags *goflag.FlagSet
	// For spf13/pflags compatibility
	PFlags *flag.FlagSet
}

func (g DefaultConfigPathGetter) GetConfigPath() string {
	if g.GoFlags != nil {
		flag.CommandLine.AddGoFlagSet(g.GoFlags)
	}
	if g.PFlags != nil {
		flag.CommandLine.AddFlagSet(g.PFlags)
	}
	configPath := flag.String("config", "", "path to config")
	flag.Parse()
	return *configPath
}

type insConfigurator struct {
	params Params
	viper  *viper.Viper
}

func NewInsConfigurator(params Params) insConfigurator {
	return insConfigurator{
		params: params,
		viper:  viper.New(),
	}
}

// Loads configuration from path and making checks
func (i *insConfigurator) Load() (ConfigStruct, error) {
	if i.params.EnvPrefix == "" {
		return nil, errors.New("EnvPrefix should be defined")
	}

	configPath := i.params.ConfigPathGetter.GetConfigPath()
	return i.load(configPath, i.params.FileRequired)
}

func (i *insConfigurator) load(path string, required bool) (ConfigStruct, error) {

	i.viper.AutomaticEnv()
	i.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	i.viper.SetEnvPrefix(i.params.EnvPrefix)

	i.viper.SetConfigFile(path)
	if err := i.viper.ReadInConfig(); err != nil {
		if required {
			return nil, err
		}
		fmt.Printf("failed to load config from '%s'\n", path)
	}
	actual := i.params.ConfigStruct.GetConfig()
	i.params.ViperHooks = append(i.params.ViperHooks, mapstructure.StringToTimeDurationHookFunc(), mapstructure.StringToSliceHookFunc(","))
	err := i.viper.UnmarshalExact(actual, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		i.params.ViperHooks...,
	)))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal config file into configuration structure")
	}
	configStructKeys, err := i.checkAllValuesIsSet()
	if err != nil {
		return nil, err
	}

	if err := i.checkNoExtraENVValues(configStructKeys); err != nil {
		return nil, err
	}

	// Second Unmarshal needed because of bug https://github.com/spf13/viper/issues/761
	// This should be evaluated after manual values overriding is done
	err = i.viper.UnmarshalExact(actual, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		i.params.ViperHooks...,
	)))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal config file into configuration structure 2")
	}
	return actual.(ConfigStruct), nil
}

func (i *insConfigurator) checkNoExtraENVValues(structKeys []string) error {
	prefixLen := len(i.params.EnvPrefix)
	for _, e := range os.Environ() {
		if len(e) > prefixLen && e[0:prefixLen]+"_" == strings.ToUpper(i.params.EnvPrefix)+"_" {
			kv := strings.SplitN(e, "=", 2)
			key := strings.ReplaceAll(strings.Replace(strings.ToLower(kv[0]), i.params.EnvPrefix+"_", "", 1), "_", ".")
			found := false
			for _, val := range structKeys {
				if strings.ToLower(val) == key {
					found = true
					// This manually sets value from ENV and overrides everything, this temporarily fix issue https://github.com/spf13/viper/issues/761
					i.viper.Set(key, kv[1])
					break
				}
			}
			if !found {
				return errors.New(fmt.Sprintf("Value not found in config: %s", key))
			}
		}
	}
	return nil
}

func (i *insConfigurator) checkAllValuesIsSet() ([]string, error) {
	names := deepFieldNames(i.params.ConfigStruct, "")
	for _, keyName := range names {
		if !i.viper.IsSet(keyName) {
			// Due to a bug https://github.com/spf13/viper/issues/447 we can't use InConfig, so
			if !stringInSlice(strings.ToLower(keyName), i.viper.AllKeys()) {
				return nil, errors.New(fmt.Sprintf("Value not found in config: %s", keyName))
			}
			// Value of this key is "null" but it set in config file
		}
	}
	return names, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func deepFieldNames(iface interface{}, prefix string) []string {
	names := make([]string, 0)
	ifv := reflect.ValueOf(iface)

	for i := 0; i < ifv.NumField(); i++ {
		v := ifv.Field(i)

		switch v.Kind() {
		case reflect.Struct:
			subPrefix := ""
			if prefix != "" {
				subPrefix = prefix + "."
			}
			names = append(names, deepFieldNames(v.Interface(), subPrefix+ifv.Type().Field(i).Name)...)
		default:
			prefWithPoint := ""
			if prefix != "" {
				prefWithPoint = prefix + "."
			}
			names = append(names, prefWithPoint+ifv.Type().Field(i).Name)
		}
	}

	return names
}

// todo clean password
func (i *insConfigurator) ToString(c ConfigStruct) string {
	out, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Sprintf("failed to marshal config structure: %v", err)
	}
	return string(out)
}
