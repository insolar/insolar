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
	"errors"

	"github.com/spf13/viper"
)

// TODO: interface for configuration
//TODO: generate default config method + test golden file

type Configuration struct {
	Host  HostNetwork
	Node  NodeNetwork
	Log   Log
	Stats Stats

	viper *viper.Viper
}

// NewConfiguration creates new default configuration
func NewConfiguration() Configuration {
	return Configuration{
		Host:  NewHostNetwork(),
		Node:  NewNodeNetwork(),
		Log:   NewLog(),
		Stats: NewStats(),
		viper: viper.New(),
	}
}

func (c *Configuration) Load() error {
	return errors.New("not implemented")
}

func (c *Configuration) Save() error {
	return errors.New("not implemented")

	//viper.
}
