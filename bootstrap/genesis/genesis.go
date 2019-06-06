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

	"github.com/insolar/insolar/application/contract/member"
	"github.com/insolar/insolar/application/contract/nodedomain"
	rootdomaincontract "github.com/insolar/insolar/application/contract/rootdomain"
	walletcontract "github.com/insolar/insolar/application/contract/wallet"
	"github.com/insolar/insolar/bootstrap"
	"github.com/insolar/insolar/bootstrap/rootdomain"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/secrets"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/artifact"
	"github.com/insolar/insolar/platformpolicy"

	"github.com/pkg/errors"
)

var contractNames = []string{
	insolar.GenesisNameRootDomain,
	insolar.GenesisNameNodeDomain,
	insolar.GenesisNameNodeRecord,
	insolar.GenesisNameMember,
	insolar.GenesisNameWallet,
	insolar.GenesisNameDeposit,
}

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
	config          *Config
	artifactManager artifact.Manager

	keyOut string
}

// NewGenerator creates new Generator.
func NewGenerator(
	config *Config,
	am artifact.Manager,
	genesisKeyOut string,
) *Generator {
	return &Generator{
		config:          config,
		artifactManager: am,

		keyOut: genesisKeyOut,
	}
}

// Run generates genesis data via headless bootstrap step.
//
// 1. builds genesis Go-contracts
// 2. stores smart contracts objects on ledger
// 3. read root keys file and generates keys for discovery nodes
//    generates and saves certificates to files for all discovery nodes
//
// when built-in contracts (INS-2308) would be implemented:
// * step 1 should be gone
// * step 2 should be moved to heavy (INS-2265)
//
func (g *Generator) Run(ctx context.Context) error {
	fmt.Printf("[ Genesis] config:\n%v\n", dumpAsJSON(g.config))

	inslog := inslogger.FromContext(ctx)
	inslog.Info("[ Genesis ] Starting  ...")

	inslog.Info("[ Genesis ] newContractBuilder ...")
	cb := newContractBuilder(g.artifactManager)
	defer cb.clean()

	inslog.Info("[ Genesis ] buildSmartContracts ...")
	prototypes, err := cb.buildPrototypes(ctx, contractNames)
	if err != nil {
		panic(errors.Wrap(err, "[ Genesis ] couldn't build contracts"))
	}

	inslog.Info("[ Genesis ] Read root keys file ...")
	rootKeys, err := secrets.ReadKeysFile(g.config.RootKeysFile)
	if err != nil {
		return errors.Wrap(err, "[ Genesis ] couldn't get root keys")
	}
	rootPubStr := platformpolicy.MustPublicKeyToString(rootKeys.Public)

	inslog.Info("[ Genesis ] Read md admin keys file ...")
	mdAdminKeys, err := secrets.ReadKeysFile(g.config.MDAdminKeysFile)
	if err != nil {
		return errors.Wrap(err, "[ Genesis ] couldn't get md admin keys")
	}
	mdAdminPubStr := platformpolicy.MustPublicKeyToString(mdAdminKeys.Public)

	inslog.Info("[ Genesis ] Read oracles keys file ...")
	oraclePubStrs := map[string]string{}
	for _, o := range g.config.OracleKeysFiles {
		oracleKeys, err := secrets.ReadKeysFile(o.KeysFile)
		if err != nil {
			return errors.Wrap(err, "[ Genesis ] couldn't get md admin keys")
		}
		oraclePubStrs[o.Name] = platformpolicy.MustPublicKeyToString(oracleKeys.Public)
	}

	err = g.activateSmartContracts(ctx, rootPubStr, mdAdminPubStr, prototypes)
	if err != nil {
		panic(errors.Wrap(err, "[ Genesis ] could't activate smart contracts"))
	}

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
	err = g.makeHeavyGenesisConfig(discoveryNodes)
	if err != nil {
		return errors.Wrap(err, "[ Genesis ] generate heavy genesis config failed")
	}

	inslog.Info("[ Genesis ] Finished.")
	return nil
}

func (g *Generator) activateRootDomain(
	ctx context.Context,
	rootDomainProto insolar.Reference,
) error {
	inslog := inslogger.FromContext(ctx)

	data, err := insolar.Serialize(&rootdomaincontract.RootDomain{
		RootMember: bootstrap.ContractRootMember,
		//OracleMembers:    bootstrap.ContractOracleMembers,
		MDAdminMember:     bootstrap.ContractMDAdminMember,
		MDWallet:          bootstrap.ContractMDWallet,
		BurnAddressMap:    map[string]insolar.Reference{},
		PublicKeyMap:      map[string]insolar.Reference{},
		FreeBurnAddresses: []string{},
		NodeDomain:        bootstrap.ContractNodeDomain,
	})
	if err != nil {
		return errors.Wrap(err, "[ activateRootDomain ] serialization failed")
	}

	_, err = g.artifactManager.RegisterRequest(
		ctx,
		record.Request{
			CallType: record.CTGenesis,
			Method:   insolar.GenesisNameRootDomain,
		},
	)
	if err != nil {
		panic(errors.Wrap(err, "[ activateRootDomain ] Couldn't create rootdomain instance"))
	}

	rootDomainRef := rootdomain.RootDomain.Ref()
	rootDomainDesc, err := g.artifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		rootDomainRef,
		insolar.GenesisRecord.Ref(),
		rootDomainProto,
		false,
		data,
	)
	if err != nil {
		return errors.Wrap(err, "[ activateRootDomain ] Couldn't create rootdomain instance")
	}

	_, err = g.artifactManager.RegisterResult(ctx, insolar.GenesisRecord.Ref(), rootDomainRef, nil)
	if err != nil {
		return errors.Wrap(err, "[ activateRootDomain ] Couldn't create rootdomain instance")
	}

	inslog.Infof("[ activateRootDomain ] create rootdomain with reference=%v", rootDomainDesc.HeadRef().String())

	return nil
}

func (g *Generator) activateNodeDomain(
	ctx context.Context, nodeDomainProto insolar.Reference,
) error {
	nd, _ := nodedomain.NewNodeDomain()

	instanceData, err := insolar.Serialize(nd)
	if err != nil {
		return errors.Wrap(err, "[ activateNodeDomain ] node domain serialization")
	}

	contractID, err := g.artifactManager.RegisterRequest(
		ctx,
		record.Request{
			CallType: record.CTGenesis,
			Method:   insolar.GenesisNameNodeDomain,
		},
	)

	if err != nil {
		return errors.Wrap(err, "[ activateNodeDomain ] couldn't create nodedomain instance")
	}
	contract := insolar.NewReference(rootdomain.RootDomain.ID(), *contractID)

	inslogger.FromContext(ctx).Infof("[ activateNodeDomain ] Ref: %v", contract)

	_, err = g.artifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		*contract,
		bootstrap.ContractRootDomain,
		nodeDomainProto,
		false,
		instanceData,
	)
	if err != nil {
		return errors.Wrap(err, "[ activateNodeDomain ] couldn't create nodedomain instance")
	}
	_, err = g.artifactManager.RegisterResult(ctx, bootstrap.ContractRootDomain, *contract, nil)
	if err != nil {
		return errors.Wrap(err, "[ activateNodeDomain ] couldn't create nodedomain instance")
	}

	inslogger.FromContext(ctx).Infof("[ activateNodeDomain ] %v contract ref=%v", bootstrap.ContractNodeDomain, contract)

	return nil
}

func (g *Generator) activateRootMember(
	ctx context.Context,
	rootPubKey string,
	memberContractProto insolar.Reference,
) error {
	m, err := member.New("RootMember", rootPubKey)
	if err != nil {
		return errors.Wrap(err, "[ activateRootMember ] root member constructor failed")
	}

	instanceData, err := insolar.Serialize(m)
	if err != nil {
		return errors.Wrap(err, "[ activateRootMember ] root member serialization failed")
	}

	contractID, err := g.artifactManager.RegisterRequest(
		ctx,
		record.Request{
			CallType: record.CTGenesis,
			Method:   "root" + insolar.GenesisNameMember,
		},
	)

	if err != nil {
		return errors.Wrap(err, "[ activateRootMember ] request register root rootMember instance failed")
	}
	contract := insolar.NewReference(rootdomain.RootDomain.ID(), *contractID)
	memberDesc, err := g.artifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		*contract,
		bootstrap.ContractRootDomain,
		memberContractProto,
		false,
		instanceData,
	)
	if err != nil {
		return errors.Wrap(err, "[ activateRootMember ] activation root rootMember instance failed")
	}
	_, err = g.artifactManager.RegisterResult(ctx, bootstrap.ContractRootDomain, *contract, nil)
	if err != nil {
		return errors.Wrap(err, "[ activateRootMember ] result register root rootMember instance failed")
	}

	inslogger.FromContext(ctx).Infof("[ activateRootMember ] register by ref=%v", memberDesc.HeadRef().String())
	return nil
}

func (g *Generator) activateOracleMembers(
	ctx context.Context,
	oraclePubKeys map[string]string,
	memberContractProto insolar.Reference,
) error {

	for name, key := range oraclePubKeys {
		m, err := member.New(name, key)
		if err != nil {
			return errors.Wrap(err, "[ activateOracleMembers ] oracle member constructor failed")
		}

		instanceData, err := insolar.Serialize(m)
		if err != nil {
			return errors.Wrap(err, "[ activateOracleMembers ] oracle member serialization failed")
		}

		contractID, err := g.artifactManager.RegisterRequest(
			ctx,
			record.Request{
				CallType: record.CTGenesis,
				Method:   name + insolar.GenesisNameMember,
			},
		)

		if err != nil {
			return errors.Wrap(err, "[ activateOracleMembers ] request register oracle instance failed")
		}
		contract := insolar.NewReference(rootdomain.RootDomain.ID(), *contractID)
		memberDesc, err := g.artifactManager.ActivateObject(
			ctx,
			insolar.Reference{},
			*contract,
			bootstrap.ContractRootDomain,
			memberContractProto,
			false,
			instanceData,
		)
		if err != nil {
			return errors.Wrap(err, "[ activateOracleMembers ] activation oracle instance failed")
		}
		_, err = g.artifactManager.RegisterResult(ctx, bootstrap.ContractRootDomain, *contract, nil)
		if err != nil {
			return errors.Wrap(err, "[ activateOracleMembers ] result register oracle instance failed")
		}

		inslogger.FromContext(ctx).Infof("[ activateOracleMembers ] register by ref=%v", memberDesc.HeadRef().String())
	}
	return nil
}

func (g *Generator) activateMDAdminMember(
	ctx context.Context,
	mdAdminPubKey string,
	memberContractProto insolar.Reference,
) error {
	m, err := member.New("MDAdminMember", mdAdminPubKey)
	if err != nil {
		return errors.Wrap(err, "[ activateMDAdminMember ] mdAdmin member constructor failed")
	}

	instanceData, err := insolar.Serialize(m)
	if err != nil {
		return errors.Wrap(err, "[ activateMDAdminMember ] mdAdmin member serialization failed")
	}

	contractID, err := g.artifactManager.RegisterRequest(
		ctx,
		record.Request{
			CallType: record.CTGenesis,
			Method:   "mdadmin" + insolar.GenesisNameMember,
		},
	)

	if err != nil {
		return errors.Wrap(err, "[ activateMDAdminMember ] request register mdAdmin instance failed")
	}
	contract := insolar.NewReference(rootdomain.RootDomain.ID(), *contractID)
	memberDesc, err := g.artifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		*contract,
		bootstrap.ContractRootDomain,
		memberContractProto,
		false,
		instanceData,
	)
	if err != nil {
		return errors.Wrap(err, "[ activateMDAdminMember ] activation mdAdmin instance failed")
	}
	_, err = g.artifactManager.RegisterResult(ctx, bootstrap.ContractRootDomain, *contract, nil)
	if err != nil {
		return errors.Wrap(err, "[ activateMDAdminMember ] result register mdAdmin instance failed")
	}

	inslogger.FromContext(ctx).Infof("[ activateMDAdminMember ] register by ref=%v", memberDesc.HeadRef().String())
	return nil
}

func (g *Generator) activateRootWallet(
	ctx context.Context, walletContractProto insolar.Reference,
) error {

	w, err := walletcontract.New(g.config.RootBalance)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ] failed to create wallet instance")
	}

	instanceData, err := insolar.Serialize(w)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ] failed to serialize wallet data")
	}

	contractID, err := g.artifactManager.RegisterRequest(
		ctx,
		record.Request{
			CallType: record.CTGenesis,
			Method:   "root" + insolar.GenesisNameWallet,
		},
	)

	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ] couldn't register root wallet")
	}
	contract := insolar.NewReference(rootdomain.RootDomain.ID(), *contractID)
	_, err = g.artifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		*contract,
		bootstrap.ContractRootMember,
		walletContractProto,
		true,
		instanceData,
	)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ] couldn't activate root wallet")
	}
	_, err = g.artifactManager.RegisterResult(ctx, bootstrap.ContractRootDomain, *contract, nil)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ] couldn't register root wallet")
	}

	return nil
}

func (g *Generator) activateMDWallet(
	ctx context.Context, walletContractProto insolar.Reference,
) error {

	w, err := walletcontract.New(g.config.MDBalance)
	if err != nil {
		return errors.Wrap(err, "[ activateMDWallet ] failed to create wallet instance")
	}

	instanceData, err := insolar.Serialize(w)
	if err != nil {
		return errors.Wrap(err, "[ activateMDWallet ] failed to serialize wallet data")
	}

	contractID, err := g.artifactManager.RegisterRequest(
		ctx,
		record.Request{
			CallType: record.CTGenesis,
			Method:   "md" + insolar.GenesisNameWallet,
		},
	)

	if err != nil {
		return errors.Wrap(err, "[ activateMDWallet ] couldn't register root wallet")
	}
	contract := insolar.NewReference(rootdomain.RootDomain.ID(), *contractID)
	_, err = g.artifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		*contract,
		bootstrap.ContractRootMember,
		walletContractProto,
		true,
		instanceData,
	)
	if err != nil {
		return errors.Wrap(err, "[ activateMDWallet ] couldn't activate root wallet")
	}
	_, err = g.artifactManager.RegisterResult(ctx, bootstrap.ContractRootDomain, *contract, nil)
	if err != nil {
		return errors.Wrap(err, "[ activateMDWallet ] couldn't register root wallet")
	}

	return nil
}

func (g *Generator) activateSmartContracts(
	ctx context.Context,
	rootPubKey string,
	mdPubKey string,
	prototypes prototypes,
) error {
	var err error

	err = g.activateRootDomain(ctx, *prototypes[insolar.GenesisNameRootDomain])
	if err != nil {
		return errors.Wrap(err, "failed to store root domain contract")
	}

	err = g.activateNodeDomain(ctx, *prototypes[insolar.GenesisNameNodeDomain])
	if err != nil {
		return errors.Wrap(err, "failed to store node domain contract")
	}

	err = g.activateRootMember(ctx, rootPubKey, *prototypes[insolar.GenesisNameMember])
	if err != nil {
		return errors.Wrap(err, "failed to store root member contract")
	}

	err = g.activateRootWallet(ctx, *prototypes[insolar.GenesisNameWallet])
	if err != nil {
		return errors.Wrap(err, "failed to store root wallet contract")
	}

	//err = g.activateOracleMembers(ctx, oraclePubKeys, *prototypes[insolar.GenesisNameMember])
	//if err != nil {
	//	return errors.Wrap(err, "failed to store oracle members contracts")
	//}

	err = g.activateMDAdminMember(ctx, mdPubKey, *prototypes[insolar.GenesisNameMember])
	if err != nil {
		return errors.Wrap(err, "failed to store md admin member contract")
	}

	err = g.activateMDWallet(ctx, *prototypes[insolar.GenesisNameWallet])
	if err != nil {
		return errors.Wrap(err, "failed to store md wallet contract")
	}

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

func (g *Generator) makeHeavyGenesisConfig(discoveryNodes []nodeInfo) error {
	items := make([]insolar.DiscoveryNodeRegister, 0, len(g.config.DiscoveryNodes))
	for _, node := range discoveryNodes {
		items = append(items, insolar.DiscoveryNodeRegister{
			Role:      node.role,
			PublicKey: node.publicKey,
		})
	}
	cfg := &insolar.GenesisHeavyConfig{
		DiscoveryNodes: items,
	}
	b, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		return errors.Wrapf(err, "[ makeHeavyGenesisConfig ] failed to decode heavy config to json")
	}

	err = ioutil.WriteFile(g.config.HeavyGeneisConfigFile, b, 0640)
	return errors.Wrapf(err, "[ makeHeavyGenesisConfig ] failed to write heavy config "+g.config.HeavyGeneisConfigFile)
}

func dumpAsJSON(data interface{}) string {
	b, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
