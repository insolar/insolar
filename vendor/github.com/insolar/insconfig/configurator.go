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
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Params for config parsing
type Params struct {
	// EnvPrefix is a prefix for environment variables
	EnvPrefix string
	// ViperHooks is custom viper decoding hooks
	ViperHooks []mapstructure.DecodeHookFunc
	// ConfigPathGetter should return config path
	ConfigPathGetter ConfigPathGetter
	// FileNotRequired - do not return error on file not found
	FileNotRequired bool
}

// ConfigPathGetter - implement this if you don't want to use config path from --config flag
type ConfigPathGetter interface {
	GetConfigPath() string
}

type insConfigurator struct {
	params Params
	viper  *viper.Viper
}

// New creates new insConfigurator with params
func New(params Params) insConfigurator {
	return insConfigurator{
		params: params,
		viper:  viper.New(),
	}
}

// Load loads configuration from path, env and makes checks
// configStruct is a pointer to your config
func (i *insConfigurator) Load(configStruct interface{}) error {
	if i.params.EnvPrefix == "" {
		return errors.New("EnvPrefix should be defined")
	}

	configPath := i.params.ConfigPathGetter.GetConfigPath()
	return i.load(configPath, configStruct)
}

func (i *insConfigurator) load(path string, configStruct interface{}) error {

	i.viper.AutomaticEnv()
	i.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	i.viper.SetEnvPrefix(i.params.EnvPrefix)

	i.viper.SetConfigFile(path)
	if err := i.viper.ReadInConfig(); err != nil {
		if !i.params.FileNotRequired {
			return err
		}
		fmt.Printf("failed to load config from '%s'\n", path)
	}
	i.params.ViperHooks = append(i.params.ViperHooks, mapstructure.StringToTimeDurationHookFunc(), mapstructure.StringToSliceHookFunc(","))
	err := i.viper.UnmarshalExact(configStruct, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		i.params.ViperHooks...,
	)))
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal config file into configuration structure")
	}
	configStructKeys, err := i.checkAllValuesIsSet(configStruct)
	if err != nil {
		return err
	}

	if err := i.checkNoExtraENVValues(configStructKeys); err != nil {
		return err
	}

	// Second Unmarshal needed because of bug https://github.com/spf13/viper/issues/761
	// This should be evaluated after manual values overriding is done
	err = i.viper.UnmarshalExact(configStruct, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		i.params.ViperHooks...,
	)))
	if err != nil {
		return errors.Wrapf(err, "failed to unmarshal config file into configuration structure 2")
	}
	return nil
}

func (i *insConfigurator) checkNoExtraENVValues(structKeys []string) error {
	var errorKeys []string
	prefixLen := len(i.params.EnvPrefix)
	for _, e := range os.Environ() {
		if len(e) > prefixLen && e[0:prefixLen]+"_" == strings.ToUpper(i.params.EnvPrefix)+"_" {
			kv := strings.SplitN(e, "=", 2)
			key := strings.ReplaceAll(strings.Replace(strings.ToLower(kv[0]), i.params.EnvPrefix+"_", "", 1), "_", ".")
			if stringInSlice(key, structKeys) {
				// This manually sets value from ENV and overrides everything, this temporarily fix issue https://github.com/spf13/viper/issues/761
				i.viper.Set(key, kv[1])
			} else {
				errorKeys = append(errorKeys, key)
			}
		}
	}
	if len(errorKeys) > 0 {
		return errors.New(fmt.Sprintf("Wrong config keys found in ENV: %s", strings.Join(errorKeys, ", ")))
	}
	return nil
}

func (i *insConfigurator) checkAllValuesIsSet(configStruct interface{}) ([]string, error) {
	var errorKeys []string
	names := deepFieldNames(configStruct, "")
	allKeys := i.viper.AllKeys()
	for _, keyName := range names {
		if !i.viper.IsSet(keyName) {
			// Due to a bug https://github.com/spf13/viper/issues/447 we can't use InConfig, so
			if !stringInSlice(keyName, allKeys) {
				errorKeys = append(errorKeys, keyName)
			}
			// Value of this key is "null" but it's set in config file
		}
	}
	if len(errorKeys) > 0 {
		return nil, errors.New(fmt.Sprintf("Keys is not defined in config: %s", strings.Join(errorKeys, ", ")))
	}
	return names, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.ToLower(b) == strings.ToLower(a) {
			return true
		}
	}
	return false
}

func deepFieldNames(iface interface{}, prefix string) []string {
	names := make([]string, 0)
	v := reflect.ValueOf(iface)
	ifv := reflect.Indirect(v)
	s := ifv.Type()

	for i := 0; i < s.NumField(); i++ {
		v := ifv.Field(i)
		tagValue := ifv.Type().Field(i).Tag.Get("mapstructure")
		tagParts := strings.Split(tagValue, ",")

		// If "squash" is specified in the tag, we squash the field down.
		squash := false
		for _, tag := range tagParts[1:] {
			if tag == "squash" {
				squash = true
				break
			}
		}

		switch v.Kind() {
		case reflect.Struct:
			newPrefix := ""
			currPrefix := ""
			if !squash {
				currPrefix = ifv.Type().Field(i).Name
			}
			if prefix != "" {
				newPrefix = strings.Join([]string{prefix, currPrefix}, ".")
			} else {
				newPrefix = currPrefix
			}

			names = append(names, deepFieldNames(v.Interface(), newPrefix)...)
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

// ToYaml returns yaml marshalled struct
func (i *insConfigurator) ToYaml(c interface{}) string {
	// todo clean password
	out, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Sprintf("failed to marshal config structure: %v", err)
	}
	return string(out)
}
