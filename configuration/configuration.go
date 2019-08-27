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

package configuration

import (
	"fmt"
	"path/filepath"
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
	AdminAPIRunner  APIRunner
	Pulsar          Pulsar
	VersionManager  VersionManager
	KeysPath        string
	CertificatePath string
	Tracer          Tracer
	Introspection   Introspection
	Exporter        Exporter
	Bus             Bus
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
		APIRunner:       NewAPIRunner(false),
		AdminAPIRunner:  NewAPIRunner(true),
		Pulsar:          NewPulsar(),
		VersionManager:  NewVersionManager(),
		KeysPath:        "./",
		CertificatePath: "",
		Tracer:          NewTracer(),
		Introspection:   NewIntrospection(),
		Exporter:        NewExporter(),
		Bus:             NewBus(),
	}

	return cfg
}

// MustInit wrapper around Init function which panics on error.
func (h *Holder) MustInit(required bool) *Holder {
	_, err := h.Init(required)
	if err != nil {
		panic(err)
	}
	return h
}

// Init init all configuration data from config file and environment.
//
// Does not fail on not found config file if the 'required' flag set to false.
func (h *Holder) Init(required bool) (*Holder, error) {
	err := h.Load()
	if err != nil {
		if required {
			return nil, err
		}
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		// read env vars if config file is not required and viper failed to load it.
		h.viper.AutomaticEnv()
		if err = h.viper.Unmarshal(&h.Configuration); err != nil {
			return nil, err
		}
	}
	return h, nil
}

func (h *Holder) registerDefaultValue(val reflect.Value, parts ...string) {
	variablePath := strings.ToLower(strings.Join(parts, "."))

	h.viper.SetDefault(variablePath, val.Interface())
}

func (h *Holder) registerDifferentValue(val reflect.Value, parts ...string) {
	variablePath := strings.Join(parts, ".")
	previousValue := h.viper.Get(variablePath)

	if !reflect.DeepEqual(previousValue, val.Interface()) {
		h.viper.Set(variablePath, val.Interface())
	}
}

func (h *Holder) recurseCallInLeaf(cb func(reflect.Value, ...string), iface interface{}, parts ...string) {
	fldV := reflect.ValueOf(iface)
	fldT := reflect.TypeOf(iface)

	for fldPos := 0; fldPos < fldV.NumField(); fldPos++ {
		fldName, fldValue := fldT.Field(fldPos).Name, fldV.Field(fldPos)

		path := append(parts, fldName)

		switch fldValue.Kind() {
		case reflect.Struct:
			h.recurseCallInLeaf(cb, fldValue.Interface(), path...)
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

	return holder.defaults()
}

// NewHolderWithFilePaths creates new holder with possible configuration files paths.
func NewHolderWithFilePaths(files ...string) *Holder {
	cfg := NewConfiguration()
	holder := &Holder{Configuration: cfg, viper: viper.New()}

	holder.viper.SetConfigType("yml")
	for _, f := range files {
		dir, file := filepath.Split(f)
		if len(dir) == 0 {
			dir = "."
		}
		file = file[:len(file)-len(filepath.Ext(file))]

		holder.viper.AddConfigPath(dir)
		holder.viper.SetConfigName(file)
	}

	return holder.defaults()
}

func (h *Holder) defaults() *Holder {
	h.recurseCallInLeaf(h.registerDefaultValue, h.Configuration)

	h.viper.AutomaticEnv()
	h.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	h.viper.SetEnvPrefix("insolar")
	return h
}

// Load method reads configuration from default file path
func (h *Holder) Load() error {
	err := h.viper.ReadInConfig()
	if err != nil {
		return err
	}

	return h.viper.Unmarshal(&h.Configuration)
}

// LoadFromFile method reads configuration from particular file path
func (h *Holder) LoadFromFile(path string) error {
	h.viper.SetConfigFile(path)
	return h.Load()
}

// Save method writes configuration to default file path
func (h *Holder) Save() error {
	h.recurseCallInLeaf(h.registerDifferentValue, h.Configuration)
	return h.viper.WriteConfig()
}

// SaveAs method writes configuration to particular file path
func (h *Holder) SaveAs(path string) error {
	h.recurseCallInLeaf(h.registerDifferentValue, h.Configuration)
	return h.viper.WriteConfigAs(path)
}

// ToString converts any configuration struct to yaml string
func ToString(in interface{}) string {
	d, err := yaml.Marshal(in)
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return string(d)
}
