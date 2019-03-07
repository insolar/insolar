/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package configuration

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

// Configuration contains configuration params for all Insolar components
type Configuration struct {
	Host            HostNetwork
	Service         ServiceNetwork
	Ledger          Ledger
	Log             Log
	Metrics         Metrics
	LogicRunner     LogicRunner
	APIRunner       APIRunner
	Pulsar          Pulsar
	VersionManager  VersionManager
	KeysPath        string
	CertificatePath string
	Tracer          Tracer
}

// Holder provides methods to manage configuration
type Holder struct {
	Configuration Configuration
	viper         *viper.Viper
}

// NewConfiguration creates new default configuration
func NewConfiguration() Configuration {
	cfg := Configuration{
		Host:            NewHostNetwork(),
		Service:         NewServiceNetwork(),
		Ledger:          NewLedger(),
		Log:             NewLog(),
		Metrics:         NewMetrics(),
		LogicRunner:     NewLogicRunner(),
		APIRunner:       NewAPIRunner(),
		Pulsar:          NewPulsar(),
		VersionManager:  NewVersionManager(),
		KeysPath:        "./",
		CertificatePath: "",
		Tracer:          NewTracer(),
	}

	return cfg
}

// MustInit wrapper around Init function which panics on error.
func (c *Holder) MustInit(required bool) *Holder {
	_, err := c.Init(required)
	if err != nil {
		panic(err)
	}
	return c
}

// Init init all configuration data from config file and environment.
//
// Does not fail on not found config file if the 'required' flag set to false.
func (c *Holder) Init(required bool) (*Holder, error) {
	err := c.Load()
	if err != nil {
		if required {
			return nil, err
		}
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}
	return c, nil
}

func (c *Holder) registerDefaultValue(val reflect.Value, parts ...string) {
	variablePath := strings.ToLower(strings.Join(parts, "."))

	c.viper.SetDefault(variablePath, val.Interface())
}

func (c *Holder) registerDifferentValue(val reflect.Value, parts ...string) {
	variablePath := strings.Join(parts, ".")
	previousValue := c.viper.Get(variablePath)

	if !reflect.DeepEqual(previousValue, val.Interface()) {
		c.viper.Set(variablePath, val.Interface())
	}
}

func (c *Holder) recurseCallInLeaf(cb func(reflect.Value, ...string), iface interface{}, parts ...string) {
	fldV := reflect.ValueOf(iface)
	fldT := reflect.TypeOf(iface)

	for fldPos := 0; fldPos < fldV.NumField(); fldPos++ {
		fldName, fldValue := fldT.Field(fldPos).Name, fldV.Field(fldPos)

		path := append(parts, fldName)

		switch fldValue.Kind() {
		case reflect.Struct:
			c.recurseCallInLeaf(cb, fldValue.Interface(), path...)
		default:
			cb(fldValue, path...)
		}
	}
}

// NewHolder creates new Holder with default configuration
func NewHolder() *Holder {
	cfg := NewConfiguration()
	holder := &Holder{Configuration: cfg, viper: viper.New()}

	holder.viper.SetConfigName(".insolar")
	holder.viper.AddConfigPath(".")
	holder.viper.SetConfigType("yml")

	holder.recurseCallInLeaf(holder.registerDefaultValue, cfg)

	holder.viper.AutomaticEnv()
	holder.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	holder.viper.SetEnvPrefix("insolar")
	return holder
}

// Load method reads configuration from default file path
func (c *Holder) Load() error {
	err := c.viper.ReadInConfig()
	if err != nil {
		return err
	}

	return c.viper.Unmarshal(&c.Configuration)
}

// LoadFromFile method reads configuration from particular file path
func (c *Holder) LoadFromFile(path string) error {
	c.viper.SetConfigFile(path)
	return c.Load()
}

// Save method writes configuration to default file path
func (c *Holder) Save() error {
	c.recurseCallInLeaf(c.registerDifferentValue, c.Configuration)
	return c.viper.WriteConfig()
}

// SaveAs method writes configuration to particular file path
func (c *Holder) SaveAs(path string) error {
	c.recurseCallInLeaf(c.registerDifferentValue, c.Configuration)
	return c.viper.WriteConfigAs(path)
}

// ToString converts any configuration struct to yaml string
func ToString(in interface{}) string {
	d, err := yaml.Marshal(in)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return string(d)
}
