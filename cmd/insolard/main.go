// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/insolar/insolar/application"
	appbuiltin "github.com/insolar/insolar/application/builtin"
	"github.com/insolar/insolar/applicationbase/genesis"
	"github.com/insolar/insolar/logicrunner/builtin"

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
		states, _ := initStates(genesisConfigPath)
		s := server.NewHeavyServer(
			holder,
			genesisConfigPath,
			genesis.Options{
				States:       states,
				ParentDomain: application.GenesisNameRootDomain,
			},
			genesisOnly,
		)
		s.Serve()
	case insolar.StaticRoleLightMaterial:
		s := server.NewLightServer(holder)
		s.Serve()
	case insolar.StaticRoleVirtual:
		builtinContracts := builtin.BuiltinContracts{
			CodeRegistry:         appbuiltin.InitializeContractMethods(),
			CodeRefRegistry:      appbuiltin.InitializeCodeRefs(),
			CodeDescriptors:      appbuiltin.InitializeCodeDescriptors(),
			PrototypeDescriptors: appbuiltin.InitializePrototypeDescriptors(),
		}
		s := server.NewVirtualServer(holder, builtinContracts)
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
