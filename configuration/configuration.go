/*
 *    Copyright 2018 INS Ecosystem
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
	"github.com/spf13/viper"
)

// Configuration contains configuration params for all Insolar components
type Configuration struct {
	Host  HostNetwork
	Node  NodeNetwork
	Log   Log
	Stats Stats
}

// Holder provides methods to manage configuration
type Holder struct {
	Configuration Configuration
	viper *viper.Viper
}

// NewConfiguration creates new default configuration
func NewConfiguration() Configuration {
	cfg := Configuration{
		Host:  NewHostNetwork(),
		Node:  NewNodeNetwork(),
		Log:   NewLog(),
		Stats: NewStats(),
	}
	
	return cfg
}

// NewHolder creates new Holder with default configuration
func NewHolder() Holder {
	cfg := NewConfiguration()
	holder := Holder{cfg, viper.New()}

	holder.viper.SetConfigName("insolar")
	holder.viper.AddConfigPath("$HOME/.insolar")
	holder.viper.AddConfigPath(".")
	holder.viper.SetConfigType("yml")

	holder.viper.SetDefault("insolar", cfg)
	return holder
}

// Load method reads configuration from default file path
func (c *Holder) Load() error {
	err := c.viper.ReadInConfig()
	if err != nil {
		return err
	}

	return c.viper.UnmarshalKey("insolar", &c.Configuration)}

// LoadFromFile method reads configuration from particular file path
func (c *Holder) LoadFromFile(path string) error {
	c.viper.SetConfigFile(path)
	return c.Load()
}

// Save method writes configuration to default file path
func (c *Holder) Save() error {
	return c.viper.WriteConfig()
}

// SaveAs method writes configuration to particular file path
func (c *Holder) SaveAs(path string) error {
	return c.viper.WriteConfigAs(path)
}