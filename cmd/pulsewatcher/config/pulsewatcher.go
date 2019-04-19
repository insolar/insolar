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

package pulsewatcher

import (
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Nodes    []string
	Interval time.Duration
	Timeout  time.Duration
}

func WriteConfig(file string, conf Config) error {
	data, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, data, 0644)
}

func ReadConfig(file string) (*Config, error) {
	var conf Config
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
