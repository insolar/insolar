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

package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/server"
	"github.com/insolar/insolar/version"
)

type inputParams struct {
	configPath        string
	genesisConfigPath string
}

func parseInputParams() inputParams {
	var rootCmd = &cobra.Command{Use: "insolard"}
	var result inputParams
	rootCmd.Flags().StringVarP(&result.configPath, "config", "c", "", "path to config file")
	rootCmd.Flags().StringVarP(&result.genesisConfigPath, "heavy-genesis", "", "", "path to genesis config for heavy node")
	rootCmd.AddCommand(version.GetCommand("insolard"))
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("Wrong input params:", err)
	}

	return result
}

// toLaunch is an array of routines created in instrumentations
var toLaunch []func(inputParams) error

func main() {
	params := parseInputParams()
	jww.SetStdoutThreshold(jww.LevelDebug)

	role, err := readRole(params.configPath)
	if err != nil {
		log.Fatal(errors.Wrap(err, "readRole failed"))
	}

	for _, f := range toLaunch {
		if err := f(params); err != nil {
			log.Warnf("Error when launch startup routine: %s", err)
		}
	}

	switch role {
	case insolar.StaticRoleHeavyMaterial:
		s := server.NewHeavyServer(params.configPath, params.genesisConfigPath)
		s.Serve()
	case insolar.StaticRoleLightMaterial:
		s := server.NewLightServer(params.configPath)
		s.Serve()
	case insolar.StaticRoleVirtual:
		s := server.NewVirtualServer(params.configPath)
		s.Serve()
	}
}

func readRole(path string) (insolar.StaticRole, error) {
	var err error
	cfg := configuration.NewHolder()
	if len(path) != 0 {
		err = cfg.LoadFromFile(path)
	} else {
		err = cfg.Load()
	}
	if err != nil {
		return insolar.StaticRoleUnknown, errors.Wrap(err, "failed to load configuration from file")
	}

	data, err := ioutil.ReadFile(filepath.Clean(cfg.Configuration.CertificatePath))
	if err != nil {
		return insolar.StaticRoleUnknown, errors.Wrapf(
			err,
			"failed to read certificate from: %s",
			cfg.Configuration.CertificatePath,
		)
	}
	cert := certificate.AuthorizationCertificate{}
	err = json.Unmarshal(data, &cert)
	if err != nil {
		return insolar.StaticRoleUnknown, errors.Wrap(err, "failed to parse certificate json")
	}
	return cert.GetRole(), nil
}
