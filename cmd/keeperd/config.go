// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
