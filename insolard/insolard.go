package insolard

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/applicationbase/genesis"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/builtin"
	"github.com/insolar/insolar/server"
)

func RunInsolarNode(parentDomain string, scVersion int64, apiOptions api.Options, initStates InitStates, builtinContracts builtin.BuiltinContracts) {
	var (
		configPath        string
		genesisConfigPath string
		heavyDB           string
		genesisOnly       bool
	)

	var cmdHeavy = &cobra.Command{
		Use:   "heavy --config=path --heavy-genesis=path",
		Short: "starts heavy node",
		Run: func(cmd *cobra.Command, args []string) {
			runHeavyNode(configPath, genesisConfigPath, heavyDB, parentDomain, genesisOnly, scVersion, apiOptions, initStates)
		},
	}
	cmdHeavy.Flags().StringVarP(&genesisConfigPath, "heavy-genesis", "", "", "path to genesis config for heavy node")
	if err := cmdHeavy.MarkFlagRequired("heavy-genesis"); err != nil {
		log.Fatal("MarkFlagRequired failed:", err)
	}
	cmdHeavy.Flags().StringVarP(&heavyDB, "database", "", "", "sets database type for heavy node, available badger/postgres")
	if err := cmdHeavy.MarkFlagRequired("database"); err != nil {
		log.Fatal("MarkFlagRequired failed:", err)
	}
	cmdHeavy.Flags().BoolVarP(&genesisOnly, "genesis-only", "", false, "run only genesis and then terminate")

	var cmdLight = &cobra.Command{
		Use:   "light --config=path",
		Short: "starts light node",
		Run: func(cmd *cobra.Command, args []string) {
			runLightNode(configPath, apiOptions)
		},
	}

	var cmdVirtual = &cobra.Command{
		Use:   "virtual --config=path",
		Short: "starts virtual node",
		Run: func(cmd *cobra.Command, args []string) {
			runVirtualNode(configPath, builtinContracts, apiOptions)
		},
	}

	var rootCmd = &cobra.Command{
		Use: "insolard",
	}
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to config file")
	if err := rootCmd.MarkPersistentFlagRequired("config"); err != nil {
		log.Fatal("MarkFlagRequired failed:", err)
	}
	rootCmd.AddCommand(cmdHeavy, cmdLight, cmdVirtual)
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("insolard execution failed:", err)
	}
}

type InitStates func(string) ([]genesis.ContractState, error)

func runHeavyNode(configPath, genesisConfigPath, db, parentDomain string, genesisOnly bool, scVersion int64, apiOptions api.Options, initStates InitStates) {
	var holder configuration.ConfigHolder
	var err error

	switch db {
	case configuration.DbTypeBadger:
		holder, err = readHeavyBadgerConfig(configPath)
	case configuration.DbTypePg:
		holder, err = readHeavyPgConfig(configPath)
	default:
		log.Fatal("db type is not supported")
	}
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to load configuration"))
	}

	role, err := readRole(holder)
	if err != nil {
		log.Fatal(errors.Wrap(err, "readRole failed"))
	}
	if role != insolar.StaticRoleHeavyMaterial {
		log.Fatal(errors.New("role in cert is not heavy"))
	}

	if err := psAgentLauncher(); err != nil {
		log.Warnf("Failed to launch gops agent: %s", err)
	}

	states, _ := initStates(genesisConfigPath)
	s := server.NewHeavyServer(
		holder,
		genesisConfigPath,
		genesis.Options{
			States:       states,
			ParentDomain: parentDomain,
		},
		genesisOnly,
		apiOptions,
		scVersion,
	)
	s.Serve()
}

func runVirtualNode(configPath string, builtinContracts builtin.BuiltinContracts, apiOptions api.Options) {
	jww.SetStdoutThreshold(jww.LevelDebug)

	holder, err := readVirtualConfig(configPath)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to load configuration"))
	}

	role, err := readRole(holder)
	if err != nil {
		log.Fatal(errors.Wrap(err, "readRole failed"))
	}
	if role != insolar.StaticRoleVirtual {
		log.Fatal(errors.New("role in cert is not virtual executor"))
	}

	if err := psAgentLauncher(); err != nil {
		log.Warnf("Failed to launch gops agent: %s", err)
	}

	s := server.NewVirtualServer(holder, builtinContracts, apiOptions)
	s.Serve()
}

func runLightNode(configPath string, apiOptions api.Options) {
	jww.SetStdoutThreshold(jww.LevelDebug)

	holder, err := readLightConfig(configPath)
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to load configuration"))
	}

	role, err := readRole(holder)
	if err != nil {
		log.Fatal(errors.Wrap(err, "readRole failed"))
	}
	if role != insolar.StaticRoleLightMaterial {
		log.Fatal(errors.New("role in cert is not light material"))
	}

	if err := psAgentLauncher(); err != nil {
		log.Warnf("Failed to launch gops agent: %s", err)
	}

	s := server.NewLightServer(holder, apiOptions)
	s.Serve()
}

// psAgentLauncher is a stub for gops agent launcher (available with 'debug' build tag)
var psAgentLauncher = func() error { return nil }

func readHeavyBadgerConfig(path string) (*configuration.HeavyBadgerHolder, error) {
	cfg := configuration.NewHeavyBadgerHolder(path)
	err := cfg.Load()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load configuration")
	}
	return cfg, nil
}

func readHeavyPgConfig(path string) (*configuration.HeavyPgHolder, error) {
	cfg := configuration.NewHeavyPgHolder(path)
	err := cfg.Load()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load configuration")
	}
	return cfg, nil
}

func readLightConfig(path string) (*configuration.LightHolder, error) {
	cfg := configuration.NewLightHolder(path)
	err := cfg.Load()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load configuration")
	}
	return cfg, nil
}

func readVirtualConfig(path string) (*configuration.VirtualHolder, error) {
	cfg := configuration.NewVirtualHolder(path)
	err := cfg.Load()
	if err != nil {
		return nil, errors.Wrap(err, "failed to load configuration")
	}
	return cfg, nil
}

func readRole(holder configuration.ConfigHolder) (insolar.StaticRole, error) {
	data, err := ioutil.ReadFile(filepath.Clean(holder.GetGenericConfig().CertificatePath))
	if err != nil {
		return insolar.StaticRoleUnknown, errors.Wrapf(
			err,
			"failed to read certificate from: %s",
			holder.GetGenericConfig().CertificatePath,
		)
	}
	cert := certificate.AuthorizationCertificate{}
	err = json.Unmarshal(data, &cert)
	if err != nil {
		return insolar.StaticRoleUnknown, errors.Wrap(err, "failed to parse certificate json")
	}
	return cert.GetRole(), nil
}
