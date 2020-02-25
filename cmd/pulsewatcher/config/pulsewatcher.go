// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pulsewatcher

import (
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"
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
