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

package genesis

import (
	"context"
	"crypto"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strconv"

	"github.com/insolar/insolar/bootstrap"
	"github.com/insolar/insolar/bootstrap/rootdomain"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/secrets"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

type nodeInfo struct {
	privateKey crypto.PrivateKey
	publicKey  string
	role       string
}

func (ni nodeInfo) reference() insolar.Reference {
	return rootdomain.GenesisRef(ni.publicKey)
}

// Generator is a component for generating RootDomain instance and genesis contracts.
type Generator struct {
	config *Config
	keyOut string
}

// NewGenerator creates new Generator.
func NewGenerator(
	config *Config,
	genesisKeyOut string,
) *Generator {
	return &Generator{
		config: config,
		keyOut: genesisKeyOut,
	}
}

// Run generates genesis data via headless bootstrap step.
//
// 1. builds Go plugins for genesis contracts
//    (gone when built-in contracts (INS-2308) would be implemented)
// 2. read root keys file and generates keys for discovery nodes
// 3. generates genesis config for heavy node.
func (g *Generator) Run(ctx context.Context) error {
	fmt.Printf("[ Genesis] config:\n%v\n", dumpAsJSON(g.config))

	inslog := inslogger.FromContext(ctx)
	inslog.Info("[ Genesis ] Starting  ...")

	inslog.Info("[ Genesis ] ReadKeysFile ...")
	pair, err := secrets.ReadKeysFile(g.config.RootKeysFile)
	if err != nil {
		return errors.Wrap(err, "[ Genesis ] couldn't get root keys")
	}
	publicKey := platformpolicy.MustPublicKeyToString(pair.Public)

	inslog.Info("[ Genesis ] generate plugins ...")
	err = g.generatePlugins()
	if err != nil {
		panic(errors.Wrap(err, "[ Genesis ] could't compile smart contracts via insgocc"))
	}
	inslog.Info("[ Genesis ] generate memory files ...")

	inslog.Info("[ Genesis ] create keys ...")
	discoveryNodes, err := createKeysInDir(
		ctx,
		g.config.DiscoveryKeysDir,
		g.config.KeysNameFormat,
		g.config.DiscoveryNodes,
		g.config.ReuseKeys,
	)
	if err != nil {
		return errors.Wrapf(err, "[ Genesis ] create keys step failed")
	}

	inslog.Info("[ Genesis ] create certificates ...")
	err = g.makeCertificates(ctx, discoveryNodes)
	if err != nil {
		return errors.Wrap(err, "[ Genesis ] generate discovery certificates failed")
	}

	inslog.Info("[ Genesis ] create heavy genesis config ...")
	contractsConfig := insolar.GenesisContractsConfig{
		RootBalance:   g.config.RootBalance,
		RootPublicKey: publicKey,
	}
	err = g.makeHeavyGenesisConfig(discoveryNodes, contractsConfig)
	if err != nil {
		return errors.Wrap(err, "[ Genesis ] generate heavy genesis config failed")
	}

	inslog.Info("[ Genesis ] Finished.")
	return nil
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

			RootDomainReference: bootstrap.ContractRootDomain.String(),
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
				return errors.Wrapf(err, "[ makeCertificates ] Can't SignNetworkPart for %s",
					dn.reference())
			}

			certs[i].BootstrapNodes[j].NodeSign, err = certs[i].SignNodePart(dn.privateKey)
			if err != nil {
				return errors.Wrapf(err, "[ makeCertificates ] Can't SignNodePart for %s",
					dn.reference())
			}
		}

		// save cert to disk
		cert, err := json.MarshalIndent(certs[i], "", "  ")
		if err != nil {
			return errors.Wrapf(err, "[ makeCertificates ] Can't MarshalIndent")
		}

		if len(node.CertName) == 0 {
			return errors.New("[ makeCertificates ] cert_name must not be empty for node number " + strconv.Itoa(i+1))
		}

		certFile := path.Join(g.keyOut, node.CertName)
		err = ioutil.WriteFile(certFile, cert, 0644)
		if err != nil {
			return errors.Wrapf(err, "[ makeCertificates ] filed create ceritificate: %v", certFile)
		}

		inslogger.FromContext(ctx).Debugf("[makeCertificates] write cert file to %v", certFile)
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
		PluginsDir:      g.config.HeavyGenesisPluginsDir,
		ContractsConfig: contractsConfig,
	}
	b, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return errors.Wrapf(err, "[ makeHeavyGenesisConfig ] failed to decode heavy config to json")
	}

	err = ioutil.WriteFile(g.config.HeavyGenesisConfigFile, b, 0600)
	return errors.Wrapf(err,
		"[ makeHeavyGenesisConfig ] failed to write heavy config %v", g.config.HeavyGenesisConfigFile)
}

func dumpAsJSON(data interface{}) string {
	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
