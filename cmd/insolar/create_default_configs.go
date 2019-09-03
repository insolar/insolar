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

package main

import (
	"io/ioutil"
	"path"

	"github.com/insolar/insolar/bootstrap"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar/defaults"
	"gopkg.in/yaml.v2"
)

func baseDir() string {
	return defaults.LaunchnetDir()
}

type PulsarConfig struct {
	Pulsar configuration.Pulsar
	Tracer configuration.Tracer
	Log    configuration.Log
}

func writePulsarConfgi(cfg configuration.Configuration, outputDir string) {
	pcfg := PulsarConfig{
		Pulsar: cfg.Pulsar,
		Tracer: cfg.Tracer,
		Log:    cfg.Log,
	}
	raw, err := yaml.Marshal(pcfg)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path.Join(outputDir, "pulsar_default.yaml"), raw, 0644)
	if err != nil {
		panic(err)
	}
}

func writeBootstrapConfig(outputDir string) {
	cfg := bootstrap.Config{
		MembersKeysDir:         path.Join(baseDir(), "configs"),
		DiscoveryKeysDir:       path.Join(baseDir(), "reusekeys", "discovery"),
		NotDiscoveryKeysDir:    path.Join(baseDir(), "reusekeys", "nodes"),
		KeysNameFormat:         "/node_%02d.json",
		ReuseKeys:              false,
		HeavyGenesisConfigFile: path.Join(baseDir(), "configs", "heavy_genesis.json"),
		HeavyGenesisPluginsDir: path.Join(baseDir(), "plugins"),
		RootBalance:            "0",
		MDBalance:              "50000000000000000000",
		VestingPeriodInPulses:  10,
		VestingStepInPulses:    10,
		LockupPeriodInPulses:   20,
		MAShardCount:           10,
		PKShardCount:           10,
		Contracts: bootstrap.Contracts{
			Insgocc: path.Join("bin", "insgocc"),
			OutDir:  path.Join(baseDir(), "plugins"),
		},
		MajorityRule: 5,
		MinRoles: struct {
			Virtual       uint `mapstructure:"virtual"`
			HeavyMaterial uint `mapstructure:"heavy_material"`
			LightMaterial uint `mapstructure:"light_material"`
		}{
			Virtual:       2,
			HeavyMaterial: 1,
			LightMaterial: 2,
		},
		DiscoveryNodes: []bootstrap.Node{
			{
				Host:     "127.0.0.1:13831",
				Role:     "heavy_material",
				CertName: "discovery_cert_1.json",
			},
			{
				Host:     "127.0.0.1:23832",
				Role:     "virtual",
				CertName: "discovery_cert_2.json",
			},
			{
				Host:     "127.0.0.1:33833",
				Role:     "light_material",
				CertName: "discovery_cert_3.json",
			},
			{
				Host:     "127.0.0.1:43834",
				Role:     "virtual",
				CertName: "discovery_cert_4.json",
			},
			{
				Host:     "127.0.0.1:53835",
				Role:     "light_material",
				CertName: "discovery_cert_5.json",
			},
		},
		Nodes:            nil,
		PulsarPublicKeys: nil,
	}

	raw, err := yaml.Marshal(cfg)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path.Join(outputDir, "bootstrap_default.yaml"), raw, 0644)
	if err != nil {
		panic(err)
	}
}

type NodeConfiguration struct {
	Host            configuration.HostNetwork
	Service         configuration.ServiceNetwork
	Ledger          configuration.Ledger
	Log             configuration.Log
	Metrics         configuration.Metrics
	LogicRunner     configuration.LogicRunner
	APIRunner       configuration.APIRunner
	AdminAPIRunner  configuration.APIRunner
	KeysPath        string
	CertificatePath string
	Tracer          configuration.Tracer
	Introspection   configuration.Introspection
	Exporter        configuration.Exporter
	Bus             configuration.Bus
}

func writeNodeConfgi(cfg configuration.Configuration, outputDir string) {
	pcfg := NodeConfiguration{
		Host:            cfg.Host,
		Service:         cfg.Service,
		Ledger:          cfg.Ledger,
		Log:             cfg.Log,
		Metrics:         cfg.Metrics,
		LogicRunner:     cfg.LogicRunner,
		APIRunner:       cfg.APIRunner,
		AdminAPIRunner:  cfg.AdminAPIRunner,
		KeysPath:        cfg.KeysPath,
		CertificatePath: cfg.CertificatePath,
		Tracer:          cfg.Tracer,
		Introspection:   cfg.Introspection,
		Exporter:        cfg.Exporter,
		Bus:             cfg.Bus,
	}
	raw, err := yaml.Marshal(pcfg)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path.Join(outputDir, "node_default.yaml"), raw, 0644)
	if err != nil {
		panic(err)
	}
}
