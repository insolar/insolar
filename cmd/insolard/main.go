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

func main() {
	var (
		configPath        string
		genesisConfigPath string
	)

	var rootCmd = &cobra.Command{
		Use: "insolard",
		Run: func(_ *cobra.Command, _ []string) {
			runInsolardServer(configPath, genesisConfigPath)
		},
	}
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "path to config file")
	rootCmd.Flags().StringVarP(&genesisConfigPath, "heavy-genesis", "", "", "path to genesis config for heavy node")
	rootCmd.AddCommand(version.GetCommand("insolard"))
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("insolard execution failed:", err)
	}
}

// psAgentLauncher is a stub for gops agent launcher (available with 'debug' build tag)
var psAgentLauncher = func() error { return nil }

func runInsolardServer(configPath string, genesisConfigPath string) {
	jww.SetStdoutThreshold(jww.LevelDebug)

	holder, err := readConfig(configPath)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to load configuration"))
	}
	role, err := readRole(holder)
	if err != nil {
		log.Fatal(errors.Wrap(err, "readRole failed"))
	}

	if err := psAgentLauncher(); err != nil {
		log.Warnf("Failed to launch gops agent: %s", err)
	}

	switch role {
	case insolar.StaticRoleHeavyMaterial:
		s := server.NewHeavyServer(holder, genesisConfigPath)
		s.Serve()
	case insolar.StaticRoleLightMaterial:
		s := server.NewLightServer(holder)
		s.Serve()
	case insolar.StaticRoleVirtual:
		s := server.NewVirtualServer(holder)
		s.Serve()
	}
}

func readConfig(path string) (*configuration.Holder, error) {
	cfg := configuration.NewHolder(path)
	err := cfg.Load()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load configuration")
	}
	return cfg, nil
}

func readRole(holder *configuration.Holder) (insolar.StaticRole, error) {
	data, err := ioutil.ReadFile(filepath.Clean(holder.Configuration.CertificatePath))
	if err != nil {
		return insolar.StaticRoleUnknown, errors.Wrapf(
			err,
			"failed to read certificate from: %s",
			holder.Configuration.CertificatePath,
		)
	}
	cert := certificate.AuthorizationCertificate{}
	err = json.Unmarshal(data, &cert)
	if err != nil {
		return insolar.StaticRoleUnknown, errors.Wrap(err, "failed to parse certificate json")
	}
	return cert.GetRole(), nil
}
