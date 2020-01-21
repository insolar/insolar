// Copyright 2020 Insolar Network Ltd.
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

package main

import (
	"time"

	"github.com/insolar/insolar/configuration"
)

type Config struct {
	Log     configuration.Log
	Metrics configuration.Metrics
	Keeper  KeeperConfig
}

func NewConfig() Config {
	return Config{
		Log:     configuration.NewLog(),
		Metrics: configuration.NewMetrics(),
		Keeper:  NewKeeperConfig(),
	}
}

type KeeperConfig struct {
	ListenAddress string
	FakeTrue      bool
	PollPeriod    time.Duration
	QueryURL      string
	Queries       []string
	MaxMetricLag  time.Duration
}

func NewKeeperConfig() KeeperConfig {
	return KeeperConfig{
		ListenAddress: ":12012",
		FakeTrue:      false,
		PollPeriod:    5 * time.Second,
		QueryURL:      "https://prometheus.insolar.io/api/v1/query?query=",
		Queries:       make([]string, 0),
		MaxMetricLag:  2 * time.Minute,
	}
}
