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
	yaml "gopkg.in/yaml.v2"
)

func baseDir() string {
	return defaults.LaunchnetDir()
}

func writePulsarConfgi(outputDir string) {
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
	raw, err := yaml.Marshal(bootstrap.Config{})
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path.Join(outputDir, "bootstrap_default.yaml"), raw, 0644)
	if err != nil {
		panic(err)
	}
}

func writeNodeConfgi(outputDir string) {
	cfg := configuration.NewConfiguration()
	raw, err := yaml.Marshal(cfg)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path.Join(outputDir, "node_default.yaml"), raw, 0644)
	if err != nil {
		panic(err)
	}
}
