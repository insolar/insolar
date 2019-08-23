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

package bootstrap

import (
	"bytes"
	"context"
	"crypto"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strconv"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/genesisrefs"
	"github.com/insolar/insolar/insolar/secrets"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

// Generator is a component for generating bootstrap files required for discovery nodes bootstrap and heavy genesis.
type Generator struct {
	config             *Config
	certificatesOutDir string
}

// NewGenerator parses config file and creates new generator on success.
func NewGenerator(configFile, certificatesOutDir string) (*Generator, error) {
	config, err := ParseConfig(configFile)
	if err != nil {
		return nil, err
	}

	return NewGeneratorWithConfig(config, certificatesOutDir), nil
}

// NewGeneratorWithConfig creates new Generator with provided config.
func NewGeneratorWithConfig(config *Config, certificatesOutDir string) *Generator {
	return &Generator{
		config:             config,
		certificatesOutDir: certificatesOutDir,
	}
}

func readMigrationAddresses(file string) ([insolar.GenesisAmountMigrationAddressShards][]string, error) {
	var result [insolar.GenesisAmountMigrationAddressShards][]string
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return result, errors.Wrapf(err, " couldn't read migration addresses file %v", file)
	}

	var ma []string
	err = json.NewDecoder(bytes.NewReader(b)).Decode(&ma)
	if err != nil {
		return result, errors.Wrapf(err, "fail unmarshal migration addresses data")
	}

	for _, a := range ma {
		i := foundation.GetShardIndex(a, insolar.GenesisAmountMigrationAddressShards)
		result[i] = append(result[i])
	}

	return result, nil
}

// Run generates bootstrap data.
//
// 1. builds Go plugins for genesis contracts
//    (gone when built-in contracts (INS-2308) would be implemented)
// 2. read root keys file and generates keys and certificates for discovery nodes.
// 3. generates genesis config for heavy node.
func (g *Generator) Run(ctx context.Context) error {
	fmt.Printf("[ bootstrap ] config:\n%v\n", dumpAsJSON(g.config))

	inslog := inslogger.FromContext(ctx)

	inslog.Info("[ bootstrap ] read keys files")
	rootPublicKey, err := secrets.GetPublicKeyFromFile(g.config.MembersKeysDir + "root_member_keys.json")
	if err != nil {
		return errors.Wrap(err, "couldn't get root keys")
	}
	migrationAdminPublicKey, err := secrets.GetPublicKeyFromFile(g.config.MembersKeysDir + "migration_admin_member_keys.json")
	if err != nil {
		return errors.Wrap(err, "couldn't get migration admin keys")
	}
	migrationDaemonPublicKeys := []string{}
	for i := 0; i < insolar.GenesisAmountMigrationDaemonMembers; i++ {
		k, err := secrets.GetPublicKeyFromFile(g.config.MembersKeysDir + GetMigrationDaemonPath(i))
		if err != nil {
			return errors.Wrap(err, "couldn't get migration daemon keys")
		}
		migrationDaemonPublicKeys = append(migrationDaemonPublicKeys, k)
	}

	inslog.Info("[ bootstrap ] read migration addresses ...")
	migrationAddresses, err := readMigrationAddresses(g.config.MembersKeysDir + "migration_addresses.json")
	if err != nil {
		return errors.Wrap(err, "couldn't get migration addresses")
	}

	inslog.Info("[ bootstrap ] create keys ...")
	discoveryNodes, err := createKeysInDir(
		ctx,
		g.config.DiscoveryKeysDir,
		g.config.KeysNameFormat,
		g.config.DiscoveryNodes,
		g.config.ReuseKeys,
	)
	if err != nil {
		return errors.Wrapf(err, "create keys step failed")
	}

	inslog.Info("[ bootstrap ] create certificates ...")
	err = g.makeCertificates(ctx, discoveryNodes)
	if err != nil {
		return errors.Wrap(err, "generate discovery certificates failed")
	}

	inslog.Info("[ bootstrap ] create heavy genesis config ...")
	contractsConfig := insolar.GenesisContractsConfig{
		RootBalance:               g.config.RootBalance,
		MDBalance:                 g.config.MDBalance,
		RootPublicKey:             rootPublicKey,
		MigrationAdminPublicKey:   migrationAdminPublicKey,
		MigrationDaemonPublicKeys: migrationDaemonPublicKeys,
		MigrationAddresses:        migrationAddresses,
		VestingPeriodInPulses:     g.config.VestingPeriodInPulses,
		LokupPeriodInPulses:       g.config.LokupPeriodInPulses,
	}
	err = g.makeHeavyGenesisConfig(discoveryNodes, contractsConfig)
	if err != nil {
		return errors.Wrap(err, "generate heavy genesis config failed")
	}

	return nil
}

type nodeInfo struct {
	privateKey crypto.PrivateKey
	publicKey  string
	role       string
}

func (ni nodeInfo) reference() insolar.Reference {
	return genesisrefs.GenesisRef(ni.publicKey)
}

func (g *Generator) makeCertificates(ctx context.Context, discoveryNodes []nodeInfo) error {
	certs := make([]certificate.Certificate, 0, len(g.config.DiscoveryNodes))
	for _, node := range discoveryNodes {
		c := certificate.Certificate{
			AuthorizationCertificate: certificate.AuthorizationCertificate{
				PublicKey: node.publicKey,
				Role:      node.role,
				Reference: node.reference().String(),
			},
			MajorityRule: g.config.MajorityRule,

			RootDomainReference: genesisrefs.ContractRootDomain.String(),
		}
		c.MinRoles.Virtual = g.config.MinRoles.Virtual
		c.MinRoles.HeavyMaterial = g.config.MinRoles.HeavyMaterial
		c.MinRoles.LightMaterial = g.config.MinRoles.LightMaterial
		c.BootstrapNodes = []certificate.BootstrapNode{}

		for j, n2 := range discoveryNodes {
			host := g.config.DiscoveryNodes[j].Host
			c.BootstrapNodes = append(c.BootstrapNodes, certificate.BootstrapNode{
				PublicKey: n2.publicKey,
				Host:      host,
				NodeRef:   n2.reference().String(),
				NodeRole:  n2.role,
			})
		}

		certs = append(certs, c)
	}

	var err error
	for i, node := range g.config.DiscoveryNodes {
		for j := range g.config.DiscoveryNodes {
			dn := discoveryNodes[j]

			certs[i].BootstrapNodes[j].NetworkSign, err = certs[i].SignNetworkPart(dn.privateKey)
			if err != nil {
				return errors.Wrapf(err, "can't SignNetworkPart for %s",
					dn.reference())
			}

			certs[i].BootstrapNodes[j].NodeSign, err = certs[i].SignNodePart(dn.privateKey)
			if err != nil {
				return errors.Wrapf(err, "can't SignNodePart for %s",
					dn.reference())
			}
		}

		// save cert to disk
		cert, err := json.MarshalIndent(certs[i], "", "  ")
		if err != nil {
			return errors.Wrapf(err, "can't MarshalIndent")
		}

		if len(node.CertName) == 0 {
			return errors.New("cert_name must not be empty for node number " + strconv.Itoa(i+1))
		}

		certFile := path.Join(g.certificatesOutDir, node.CertName)
		err = ioutil.WriteFile(certFile, cert, 0600)
		if err != nil {
			return errors.Wrapf(err, "failed to create certificate: %v", certFile)
		}

		inslogger.FromContext(ctx).Infof("[ bootstrap ] write certificate file: %v", certFile)
	}
	return nil
}

func (g *Generator) makeHeavyGenesisConfig(
	discoveryNodes []nodeInfo,
	contractsConfig insolar.GenesisContractsConfig,
) error {
	items := make([]insolar.DiscoveryNodeRegister, 0, len(g.config.DiscoveryNodes))
	for _, node := range discoveryNodes {
		items = append(items, insolar.DiscoveryNodeRegister{
			Role:      node.role,
			PublicKey: node.publicKey,
		})
	}
	cfg := &insolar.GenesisHeavyConfig{
		DiscoveryNodes:  items,
		ContractsConfig: contractsConfig,
	}
	b, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return errors.Wrapf(err, "failed to decode heavy config to json")
	}

	err = ioutil.WriteFile(g.config.HeavyGenesisConfigFile, b, 0600)
	return errors.Wrapf(err,
		"failed to write heavy config %v", g.config.HeavyGenesisConfigFile)
}

func dumpAsJSON(data interface{}) string {
	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

// GetMigrationDaemonPath generate key file name for migration daemon
func GetMigrationDaemonPath(i int) string {
	return "migration_daemon_" + strconv.Itoa(i) + "_member_keys.json"
}
