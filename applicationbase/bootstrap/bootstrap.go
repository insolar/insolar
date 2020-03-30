// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package bootstrap

import (
	"context"
	"crypto"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/applicationbase/genesis"
	"github.com/insolar/insolar/applicationbase/genesisrefs"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// Generator is a component for generating bootstrap files required for discovery nodes bootstrap and heavy genesis.
type Generator struct {
	config          *Config
	contractsConfig map[string]interface{}
}

// NewGenerator parses config file and creates new generator on success.
func NewGenerator(configFile string, contractsConfig map[string]interface{}) (*Generator, error) {
	config, err := ParseConfig(configFile)
	if err != nil {
		return nil, err
	}

	return NewGeneratorWithConfig(config, contractsConfig), nil
}

// NewGeneratorWithConfig creates new Generator with provided config.
func NewGeneratorWithConfig(config *Config, contractsConfig map[string]interface{}) *Generator {
	if config.CertificateNameOffsetFromZero {
		certNamesStartFrom = 0
	}
	return &Generator{
		config:          config,
		contractsConfig: contractsConfig,
	}
}

// Run generates bootstrap data.
//
// 1. read root keys file and generates keys and certificates for discovery nodes.
// 2. generates genesis config for heavy node.
func (g *Generator) Run(ctx context.Context, properName bool) error {
	fmt.Printf("[ bootstrap ] config:\n%v\n", DumpAsJSON(g.config))

	inslog := inslogger.FromContext(ctx)

	inslog.Info("[ bootstrap ] create discovery keys ...")
	discoveryNodes, err := createKeysInDir(
		ctx,
		g.config.DiscoveryKeysDir,
		g.config.KeysNameFormat,
		g.config.DiscoveryNodes,
		g.config.ReuseKeys,
		properName,
	)
	if err != nil {
		return errors.Wrapf(err, "create keys step failed")
	}

	inslog.Info("[ bootstrap ] create discovery certificates ...")
	err = g.makeCertificates(ctx, discoveryNodes, discoveryNodes)
	if err != nil {
		return errors.Wrap(err, "generate discovery certificates failed")
	}

	if g.config.NotDiscoveryKeysDir != "" {
		inslog.Info("[ bootstrap ] create not discovery keys ...")
		nodes, err := createKeysInDir(
			ctx,
			g.config.NotDiscoveryKeysDir,
			g.config.KeysNameFormat,
			g.config.Nodes,
			g.config.ReuseKeys,
			properName,
		)
		if err != nil {
			return errors.Wrapf(err, "create keys step failed")
		}

		inslog.Info("[ bootstrap ] create not discovery certificates ...", nodes)
		err = g.makeCertificates(ctx, nodes, discoveryNodes)
		if err != nil {
			return errors.Wrap(err, "generate not discovery certificates failed")
		}
	}

	inslog.Info("[ bootstrap ] create heavy genesis config ...")
	err = g.makeHeavyGenesisConfig(discoveryNodes, g.contractsConfig)
	if err != nil {
		return errors.Wrap(err, "generate heavy genesis config failed")
	}

	return nil
}

type nodeInfo struct {
	privateKey crypto.PrivateKey
	publicKey  string
	role       string
	certName   string
}

func (ni nodeInfo) reference() insolar.Reference {
	return genesisrefs.GenesisRef(ni.publicKey)
}

func (g *Generator) makeCertificates(ctx context.Context, nodesInfo []nodeInfo, discoveryNodes []nodeInfo) error {
	certs := make([]certificate.Certificate, 0, len(g.config.DiscoveryNodes))
	for _, node := range nodesInfo {
		c := certificate.Certificate{
			AuthorizationCertificate: certificate.AuthorizationCertificate{
				PublicKey: node.publicKey,
				Role:      node.role,
				Reference: node.reference().String(),
			},
			MajorityRule: g.config.MajorityRule,
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

	err := os.MkdirAll(g.config.CertificatesOutDir, 0700)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrapf(err, "failed create dir %v: %v", g.config.CertificatesOutDir, err)
	}

	for i, node := range nodesInfo {
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

		if len(node.certName) == 0 {
			return errors.New("cert_name must not be empty for node number " + strconv.Itoa(i+1))
		}

		certFile := path.Join(g.config.CertificatesOutDir, node.certName)

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
	contractsConfig map[string]interface{},
) error {
	items := make([]genesis.DiscoveryNodeRegister, 0, len(g.config.DiscoveryNodes))
	for _, node := range discoveryNodes {
		items = append(items, genesis.DiscoveryNodeRegister{
			Role:      node.role,
			PublicKey: node.publicKey,
		})
	}

	cfg := map[string]interface{}{
		"DiscoveryNodes":  items,
		"ContractsConfig": contractsConfig,
	}

	b, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return errors.Wrapf(err, "failed to decode heavy config to json")
	}

	err = ioutil.WriteFile(g.config.HeavyGenesisConfigFile, b, 0600)
	return errors.Wrapf(err,
		"failed to write heavy config %v", g.config.HeavyGenesisConfigFile)
}

func DumpAsJSON(data interface{}) string {
	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
