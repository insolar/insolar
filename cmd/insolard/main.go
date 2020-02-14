// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	appbuiltin "github.com/insolar/insolar/application/builtin"
	"github.com/insolar/insolar/logicrunner/builtin"
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
		genesisOnly       bool
	)

	var rootCmd = &cobra.Command{
		Use: "insolard",
		Run: func(_ *cobra.Command, _ []string) {
			runInsolardServer(configPath, genesisConfigPath, genesisOnly)
		},
	}
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "path to config file")
	rootCmd.Flags().StringVarP(&genesisConfigPath, "heavy-genesis", "", "", "path to genesis config for heavy node")
	rootCmd.Flags().BoolVarP(&genesisOnly, "genesis-only", "", false, "run only genesis and then terminate")
	rootCmd.AddCommand(version.GetCommand("insolard"))
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("insolard execution failed:", err)
	}
}

// psAgentLauncher is a stub for gops agent launcher (available with 'debug' build tag)
var psAgentLauncher = func() error { return nil }

func runInsolardServer(configPath string, genesisConfigPath string, genesisOnly bool) {
	jww.SetStdoutThreshold(jww.LevelDebug)

	role, err := readRole(configPath)
	if err != nil {
		log.Fatal(errors.Wrap(err, "readRole failed"))
	}

	if err := psAgentLauncher(); err != nil {
		log.Warnf("Failed to launch gops agent: %s", err)
	}

	switch role {
	case insolar.StaticRoleHeavyMaterial:
		s := server.NewHeavyServer(configPath, genesisConfigPath, genesisOnly)
		s.Serve()
	case insolar.StaticRoleLightMaterial:
		s := server.NewLightServer(configPath)
		s.Serve()
	case insolar.StaticRoleVirtual:
		builtinContracts := builtin.BuiltinContracts{
			CodeRegistry:         appbuiltin.InitializeContractMethods(),
			CodeRefRegistry:      appbuiltin.InitializeCodeRefs(),
			CodeDescriptors:      appbuiltin.InitializeCodeDescriptors(),
			PrototypeDescriptors: appbuiltin.InitializePrototypeDescriptors(),
		}
		s := server.NewVirtualServer(configPath, builtinContracts)
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
