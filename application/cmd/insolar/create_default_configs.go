// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"io/ioutil"
	"path"

	yaml "gopkg.in/yaml.v2"

	"github.com/insolar/insolar/application/genesis/contracts"
	bootstrapbase "github.com/insolar/insolar/applicationbase/bootstrap"
	pulsewatcher "github.com/insolar/insolar/cmd/pulsewatcher/config"
	"github.com/insolar/insolar/configuration"
)

func writePulsarConfig(outputDir string) {
	pcfg := configuration.NewPulsarConfiguration()
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
	rawBase, err := yaml.Marshal(bootstrapbase.Config{})
	if err != nil {
		panic(err)
	}

	rawApp, err := yaml.Marshal(contracts.ContractsConfig{})
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path.Join(outputDir, "bootstrap_default.yaml"), append(rawBase, rawApp...), 0644)
	if err != nil {
		panic(err)
	}
}

func writeNodeConfig(outputDir string) {
	cfg := configuration.NewGenericConfiguration()
	raw, err := yaml.Marshal(cfg)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path.Join(outputDir, "node_default.yaml"), raw, 0644)
	if err != nil {
		panic(err)
	}
}

func writePulseWatcher(outputDir string) {
	raw, err := yaml.Marshal(pulsewatcher.Config{})
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path.Join(outputDir, "pulsewatcher_default.yaml"), raw, 0644)
	if err != nil {
		panic(err)
	}
}
